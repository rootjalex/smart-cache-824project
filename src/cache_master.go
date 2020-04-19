package cache

/************************************************
Cache Master supports

m = StartTask(clients []Client, caches []Cache, r int)
    Initialize a cache master with a cache list, client list, and replication factor (r)





*************************************************/


type CacheMaster {
    mu       sync.Mutex
    clients  []Client
    caches   []Cache

}

type CacheType int
const (
    LRU            CacheType = 0
    MarkovPrefetch CacheType = 1
    MarkovEviction CacheType = 2
    MarkovBoth     CacheType = 3

)


func StartTask(clients []Client, k int, r int, cacheType CacheType) *CacheMaster {
    m := CacheMaster{}
    m

    return m
}



