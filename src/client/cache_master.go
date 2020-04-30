package client

import (
    "sync"
    "../cache"
)
/************************************************
Cache Master supports

m = StartTask(clients []Client, caches []Cache, r int)
    Initialize a cache master with a cache list, client list, and replication factor (r)





*************************************************/


type CacheMaster struct {
    mu       sync.Mutex
    clients  []Client
    caches   []cache.Cache

}


func StartTask(clients []Client, k int, r int, cacheType cache.CacheType) *CacheMaster {
    m := &CacheMaster{}

    return m
}



