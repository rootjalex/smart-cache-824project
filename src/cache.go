package cache

import (
	"os"
)

type PrefetchInterface interface {
	Prefetch(filename string, n int) []string
	Access(filename string)
}

type EvictionInterface interface {
	Evict(filename string) string
	Access(filename string)
}

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
	PrefetchInterface
	EvictionInterface
	Init(cacheSize int)
	Fetch(name string) (*os.File, error)
	Report() (int, int) // Cache hits, cache misses
}


type Thing struct {
	w int 
}

func (t *Thing) What() int {
	return t.w
}

type OtherThing struct {
	Thing
}

// var t OtherThing 
// t.w = 4
// fmt.Println(t.What()) -> 4