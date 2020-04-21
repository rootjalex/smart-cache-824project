package cache

import (
	"os"
	"sync"
)

/********************************
LRU Cache supports the following external API to users

c.Init(cacheSize int)
    Initializes a cache with LRU eviction policy and no prefetching
c.Report() (hits, misses)
    Get a report of the hits and misses  TODO: Do we want a version number or
    timestamp mechanism of any form here?

*********************************/

type LRUCache struct {
	mu          sync.Mutex          // Lock to protect shared access to cache
	misses		int
	hits		int
	cache		map[string]*os.File
	heap		MinHeap
	timestamp	int64 // for controlling LRU heap
	cacheSize	int
}


func (c *LRUCache) Fetch(name string) (*os.File, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	file, ok := c.cache[name]
	var err error

	if ok {
		c.hits++
		err = nil
		c.heap.ChangeKey(name, c.timestamp)
	} else {
		file, err = os.Open(name)
		c.replace(name, file) // handles insertion into heap
		c.misses++
	}
	c.timestamp++
	return file, err
}


func (c *LRUCache) Report() (int, int) {
    c.mu.Lock()
    defer c.mu.Unlock()
	return c.hits, c.misses
}


func (c *LRUCache) Init(cacheSize int){
	c.misses = 0
	c.hits = 0
    c.cacheSize = cacheSize
	c.cache = make(map[string]*os.File)
	c.timestamp = 0
	c.heap.Init()
}


// assumes mu is Locked
func (c *LRUCache) replace(name string, file *os.File) {
	c.cache[name] = file
	c.heap.Insert(name, c.timestamp)
	if c.heap.n > c.cacheSize {
		// must evict
		evict := c.heap.ExtractMin()
		delete(c.cache, evict)
	}
}
