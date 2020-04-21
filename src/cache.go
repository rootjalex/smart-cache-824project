package cache

import (
	"os"
)

/********************************
Cache supports the following external API to users

c.Init(cacheSize int, cacheType CacheType)
    Initializes a cache with eviction policy and prefetch defined by cache type
c.Report() (hits, misses)
    Get a report of the hits and misses  TODO: Do we want a version number or
    timestamp mechanism of any form here?
c.Fetch(name string) (*os.File, error)

*********************************/
type Cache interface {
	Init(cacheSize int)
	Fetch(name string) (*os.File, error)
	Report() (int, int) // Cache hits, cache misses
}
