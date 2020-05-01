package cache

import (
//    "time"
)

const CACHE_SIZE = 1000

const PREFETCH_SIZE = 10

//const SEED = time.Now().UnixNano()
const SEED = 1

type CacheType int
const (
    LRU            CacheType = 0
    MarkovPrefetch CacheType = 1
    MarkovEviction CacheType = 2
    MarkovBoth     CacheType = 3
)