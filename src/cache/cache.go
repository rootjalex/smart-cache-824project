package cache

import (
	"sync"
	"errors"
	"../heap"
	"../markov"
	"../datastore"
	"../config"
    "log"
)

/********************************
Cache supports the following external API to users

c.Init(cacheSize int, cacheType CacheType, data *datastore.DataStore)
	Initializes a cache with eviction policy and prefetch defined by cache type
	Copies underlying datastore
c.Report() (hits, misses)
    Get a report of the hits and misses  TODO: Do we want a version number or
    timestamp mechanism of any form here?
c.Fetch(name string) (config.DataType, error)

*********************************/
type Cache struct {
	mu          sync.Mutex          // Lock to protect shared access to cache
    id          int
	misses		int
	hits		int
    cache	    map[string]config.DataType
	heap		*heap.MinHeapInt64
	timestamp	int64 // for controlling LRU heap
	cacheSize	int
	chain		*markov.MarkovChain
	cacheType	config.CacheType
	data		*datastore.DataStore
    alive       bool
}

// client -> cache (Request a file)
type RequestFileArgs struct {
	Filename string
}

type RequestFileReply struct {
	File	config.DataType
}

// copies underlying datastore
func (c *Cache) Init(id int, cacheSize int, cacheType config.CacheType, data *datastore.DataStore) {
	c.cacheType = cacheType
    c.id = id
	c.misses = 0
	c.hits = 0
    c.cacheSize = cacheSize
	c.cache = make(map[string]config.DataType)
	c.timestamp = 0
	c.data = data.Copy()

	if cacheType == config.LRU || cacheType == config.MarkovEviction {
		// only LRU caches should use heap
		c.heap = heap.MakeMinHeapInt64()
	}
	if cacheType != config.LRU {
		// all other caches need a MarkovChain
		c.chain = markov.MakeMarkovChain()
	}
}

func (c *Cache) GetId() int {
    return c.id
}

func (c *Cache) Fetch(name string) (config.DataType, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	file, ok := c.cache[name]
	var err error

	if ok {
		c.hits++
		err = nil

		// TODO: THIS IS BAAD PRACTICE BUT WILL SUFFICE FOR NOW
		if c.cacheType == config.LRU || c.cacheType == config.MarkovEviction {
			// only LRU caches should use heap
			c.heap.ChangeKey(name, c.timestamp)
		}
		if c.cacheType != config.LRU {
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
            log.Printf("Cache failed again - should not happen")
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
func (c *Cache) replace(name string, file config.DataType) {
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
	if c.cacheType == config.MarkovPrefetch || c.cacheType == config.MarkovBoth {
		files := c.chain.Predict(filename, config.PREFETCH_SIZE)
		for _, file :=  range files {
			c.AddToCache(file)
		}
	}
}

func (c *Cache) GetState(prevChain *markov.MarkovChain) *markov.MarkovChain {
    // TODO: Need some if statements around cache type here
    c.mu.Lock()
    defer c.mu.Unlock()
    return markov.ChainSub(c.chain, prevChain)
}

func (c *Cache) UpdateState(chain *markov.MarkovChain) {
    // TODO: Need some if statements around cache type here
    c.mu.Lock()
    defer c.mu.Unlock()
    c.chain = chain.Copy()
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


func (c *Cache) FetchRPC(args *RequestFileArgs, reply *RequestFileReply) error {
	var err error
	reply.File, err = c.Fetch(args.Filename)
	return err
}
