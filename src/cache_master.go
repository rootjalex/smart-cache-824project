package cache

import (
    "sync"
)
/************************************************
Cache Master API

Initialization: 
    m = StartTask(
            clients []Client,
            cacheType CacheType - specification for prefetch and eviction policies
            k int - number of cache machines to use
            r int - replication factor
            datastore Datastore
        )
    Initialize a cache master with client list, and replication factor (r)



*************************************************/


type CacheMaster struct {
    mu          sync.Mutex
    clients     []Client
    caches      []Cache
    cacheType   CacheType
    r           int // replication factor
    k           int // number of caches
    n           int // number of pieces of data
    datastore   *DataStore
    hash        *Hash

}

type CacheType int
const (
    LRU            CacheType = 0
    MarkovPrefetch CacheType = 1
    MarkovEviction CacheType = 2
    MarkovBoth     CacheType = 3

)


func StartTask(clients []Client, cacheType CacheType, cacheSize int, k int, r int,
               datastore *DataStore) ([]Cache, *Hash) {
    // k: number of caches
    // r: replication factor for data desired
    // this is trivial (can store everything) if cacheSize >= nr/k (where n is
    // size of datastore)
    m := &CacheMaster{}
    m.clients = clients
    m.k = k
    m.r = r
    m.n = datastore.Size()
    m.datastore = datastore
    m.caches = make([]Cache, k)
    for i := 0; i < k; i++ {
        if cacheType == LRU {
            m.caches[i] = &LRUCache{}
        }
        m.caches[i].Init(cacheSize)
    }
    m.hash = MakeHash(k, m.n)
    return m.caches, m.hash
}



