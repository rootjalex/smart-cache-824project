package cache

import (
	"os"
)

type Cache interface {
	Fetch(name string) (*os.File, error)
	Report() (int, int) // Cache hits, cache misses
}

type BadCache struct {
	misses int
}

func (c *BadCache) Fetch(name string) (*os.File, error) { 
	file, err := os.Open(name)
	c.misses++
	return file, err
 }


func (c *BadCache) Report() (int, int) {
	return 0, c.misses
}