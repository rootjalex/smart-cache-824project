package cache

import (
	"sync"
	"../heap"
	"../markov"
	"../datastore"
	"errors"
)

/********************************
Cache supports the following external API to users

c.Init(cacheSize int, cacheType CacheType, data *datastore.DataStore)
	Initializes a cache with eviction policy and prefetch defined by cache type
	Copies underlying datastore
c.Report() (hits, misses)
    Get a report of the hits and misses  TODO: Do we want a version number or
    timestamp mechanism of any form here?
c.Fetch(name string) (datastore.DataType, error)

*********************************/
type Cache struct {
	mu          sync.Mutex          // Lock to protect shared access to cache
	misses		int
	hits		int
	cache		map[string]datastore.DataType
	heap		heap.MinHeap
	timestamp	int64 // for controlling LRU heap
	cacheSize	int
	chain 		*markov.MarkovChain
	cacheType	CacheType
	data 		*datastore.DataStore
}

// copies underlying datastore
func (c *Cache) Init(cacheSize int, cacheType CacheType, data *datastore.DataStore) {
	c.cacheType = cacheType
	c.misses = 0
	c.hits = 0
    c.cacheSize = cacheSize
	c.cache = make(map[string]datastore.DataType)
	c.timestamp = 0
	c.data = data.Copy()
	
	if cacheType == LRU || cacheType == MarkovEviction {
		// only LRU caches should use heap
		c.heap.Init()
	}
	if cacheType != LRU {
		// all other caches need a MarkovChain
		c.chain = markov.MakeMarkovChain()
	}
}

func (c *Cache) Fetch(name string) (datastore.DataType, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	file, ok := c.cache[name]
	var err error

	if ok {
		c.hits++
		err = nil

		// TODO: THIS IS BAAD PRACTICE BUT WILL SUFFICE FOR NOW
		if c.cacheType == LRU || c.cacheType == MarkovEviction {
			// only LRU caches should use heap
			c.heap.ChangeKey(name, c.timestamp)
		}
		if c.cacheType != LRU {
			// all other caches need a MarkovChain
			c.chain.Access(name)
		}
	} else {
		c.AddToCache(name)
		c.misses++
		file, ok = c.cache[name]
		if !ok {
			// failed again - should not happen
			err = errors.New("failed")
		}
	}
	c.timestamp++

	go c.Prefetch(name)
	return file, err
}


func (c *Cache) Report() (int, int) {
    c.mu.Lock()
    defer c.mu.Unlock()
	return c.hits, c.misses
}

// TODO: REPLACEMENT POLICY FOR MARKOV CHAIN
// assumes mu is Locked
func (c *Cache) replace(name string, file datastore.DataType) {
	c.cache[name] = file
	c.heap.Insert(name, c.timestamp)
	if c.heap.Size > c.cacheSize {
		// must evict
		evict := c.heap.ExtractMin()
		delete(c.cache, evict)
	}
}

func (c *Cache) Prefetch(filename string) {
    c.mu.Lock()
	defer c.mu.Unlock()
	if c.cacheType == MarkovPrefetch || c.cacheType == MarkovBoth {
		files := c.chain.Predict(filename, PREFETCH_SIZE)
		for _, file :=  range files {
			c.AddToCache(file)
		}
	}
}

func (c *Cache) AddToCache(filename string) bool {
	// assumes c.mu is held
	_, ok := c.cache[filename]

	if !ok {
		file, err := c.data.Get(filename)
		c.replace(filename, file) // handles insertion into heap
		return err
	}
	return ok
}