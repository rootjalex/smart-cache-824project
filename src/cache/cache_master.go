package cache

import (
	"sync"
	"time"

	"../datastore"
)

/************************************************
Cache Master API

Initialization:
    m = StartTask(
            clients      []Client
            cacheType   CacheType - specification for prefetch and eviction policies
            numCaches         int - number of cache machines to use
            replication       int - replication factor
            datastore   Datastore
        )
    Initialize a cache master with client list, and replication factor (r)

syncCaches:


*************************************************/

type CacheMaster struct {
	mu          sync.Mutex
	clients     []*Client
	caches      map[int]*Cache
	cacheType   CacheType
	replication int // replication factor
	numCaches   int // number of caches
	n           int // number of pieces of data
	datastore   *datastore.DataStore
	hash        *Hash
	ms          int // how often caches are synced

}

func StartTask(clients []*Client, cacheType CacheType, cacheSize int, numCaches int, replication int, datastore *datastore.DataStore, ms int) (map[int]*Cache, *Hash) {
	// k: number of caches
	// r: replication factor for data desired
	// this is trivial (can store everything) if cacheSize >= nr/k (where n is
	// size of datastore)
	m := &CacheMaster{}
	m.clients = clients
	m.numCaches = numCaches
	m.replication = replication
	m.datastore = datastore
	m.n = datastore.Size()
	m.ms = ms
	m.caches = map[int]*Cache{}
	for i := 0; i < m.numCaches; i++ {
		c := Cache{}
		c.Init(i, cacheSize, cacheType, m.datastore)
		m.caches[i] = &c
	}
	m.hash = MakeHash(m.numCaches, m.datastore.GetFileNames(), m.n, m.replication, m.clients)

	go m.syncCaches(ms)

	return m.caches, m.hash
}

func (m *CacheMaster) requestCacheState(cacheId int, args *GetCacheStateArgs, reply *GetCacheStateReply) bool {
	//ok := m.caches[cacheId].Call("Cache.GetState", args, reply)
	ok := m.caches[cacheId].GetState(args, reply)
	return ok
}

func (m *CacheMaster) updateCacheState(cacheId int, args *UpdateCacheArgs, reply *UpdateCacheReply) bool {
	//ok := m.caches[cacheId].Call("Cache.UpdateState", args, reply)
	ok := m.caches[cacheId].UpdateState(args, reply)
	return ok
}

func (m *CacheMaster) syncCaches(ms int) {
	for {
		// for group := range(m.hash.getGroups()){
		// }
		for i := 0; i < m.numCaches; i++ {
			args := GetCacheStateArgs{}
			reply := GetCacheStateReply{}
			m.requestCacheState(i, &args, &reply)
		}

		// cast int to duration for multiplication to work
		time.Sleep(time.Duration(ms) * time.Millisecond)
	}
}
