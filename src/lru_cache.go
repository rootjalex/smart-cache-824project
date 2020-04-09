package cache

import (
	"os"
	"log"
)

type LRUCache struct {
	misses 		int
	hits 		int
	cache 		map[string]*os.File	
	heap		MinHeap
	timestamp 	int64 // for controlling LRU heap
}

func (c *LRUCache) Fetch(name string) (*os.File, error) { 
	c.checkSize()
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
	c.checkSize()
	return file, err
}


func (c *LRUCache) Report() (int, int) {
	c.checkSize()
	return c.hits, c.misses
}

func (c *LRUCache) Init() {
	c.misses = 0
	c.hits = 0
	c.cache = make(map[string]*os.File)
	c.timestamp = 0
	c.heap.Init()
	c.checkSize()
}

func (c *LRUCache) replace(name string, file *os.File) {
	c.checkSize()
	c.cache[name] = file
	c.heap.Insert(name, c.timestamp)
	if c.heap.n > CACHE_SIZE {
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