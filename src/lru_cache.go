package cache

import (
	"os"
	"log"
)

/********************************
LRU Cache supports the following external API to users

c = Make()
    Initializes a LRUCache
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
}

func (c *LRUCache) Fetch(name string) (*os.File, error) {
	c.checkSize()
	file, ok := c.cache[name]
	var err error

    c.mu.Lock()
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
    c.mu.Unlock()
	c.checkSize()
	return file, err
}

func HandleRequests(c *LRUCache) {
    // TODO: listen on a channel for data requests?
}


func (c *LRUCache) Report() (int, int) {
	c.checkSize()
    c.mu.Lock()
    defer c.mu.Unlock()
	return c.hits, c.misses
}

func Make(cacheSize int) *LRUCache {
	c := &LRUCache{}
    //lc.Init()
    //Init(lc)
    return c
}

// Is there a reason these are RPC handler functions and not regular functions?
func (c *LRUCache) Init(cacheSize int) {
	c.misses = 0
	c.hits = 0
    c.cacheSize = cacheSize
	c.cache = make(map[string]*os.File)
	c.timestamp = 0
	c.heap.Init()
	c.checkSize()
    go HandleRequests(c)
}

func (c *LRUCache) replace(name string, file *os.File) {
	c.checkSize()
	c.cache[name] = file
	c.heap.Insert(name, c.timestamp)
	if c.heap.n > c.cacheSize {
		// must evict
		evict := c.heap.ExtractMin()
		delete(c.cache, evict)
	}
}

// sanity check
func (c *LRUCache) checkSize() {
	if c.heap.n != len(c.cache) {
		log.Fatalf("Sizes don't match, heap size: %d and cache size: %d", c.heap.n, len(c.cache))
	}
}
