package cache

import (
	"os"
)

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

func (c *BadCache) Init() {
	c.misses = 0
}