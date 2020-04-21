package cache

import (
	"os"
)

type Cache interface {
	Init(cacheSize int)
	Fetch(name string) (*os.File, error)
	Report() (int, int) // Cache hits, cache misses
}
