package task

import (
	"sync"
	"time"
	"../datastore"
	"../markov"
	"../config"
	"../cache"
)

/************************************************
Cache Master API

Initialization:
    m = StartTask(
            clientIds       []int
            cacheType   CacheType - specification for prefetch and eviction policies
            numCaches         int - number of cache machines to use
            replication       int - replication factor
            datastore   Datastore
        )
    Initialize a cache master with client list, and replication factor (r)

syncCaches


*************************************************/

type CacheMaster struct {
	mu          sync.Mutex
	clientIds   []int
	caches      map[int]*cache.Cache
	cacheType   config.CacheType
	replication int // replication factor
	numCaches   int // number of caches
	n           int // number of pieces of data
	datastore   *datastore.DataStore
	hash        *cache.Hash
	ms          int // how often caches are synced
    chain       *markov.MarkovChain

}

func StartTask(clientIds []int, cacheType config.CacheType, cacheSize int, numCaches int,
	replication int, datastore *datastore.DataStore, ms int) (map[int]*cache.Cache, *cache.Hash) {
	// k: number of caches
	// r: replication factor for data desired
	// this is trivial (can store everything) if cacheSize >= nr/k (where n is
	// size of datastore)
	m := &CacheMaster{}
	m.clientIds = clientIds
	m.numCaches = numCaches
	m.replication = replication
	m.datastore = datastore
	m.n = datastore.Size()
    m.chain = markov.MakeMarkovChain()
	m.ms = ms
	m.caches = map[int]*cache.Cache{}
	for i := 0; i < m.numCaches; i++ {
		c := cache.Cache{}
		c.Init(i, cacheSize, cacheType, m.datastore.Copy())
		m.caches[i] = &c
	}


	m.hash = cache.MakeHash(m.numCaches, m.datastore.GetFileNames(), m.n, m.replication, m.clientIds)

    if (cacheType != config.LRU) {
        go m.syncCaches(ms)
    }

	return m.caches, m.hash
}


func (m *CacheMaster) syncGroup(groupID int) {
    cacheIDs := m.hash.GetCachesInGroup(groupID)
    newChain := m.chain.Copy()
    for _, id := range cacheIDs {
        newChain = markov.ChainAdd(newChain, m.caches[id].GetState(m.chain))
    }
    m.chain = newChain

    for _, id := range cacheIDs {
        m.caches[id].UpdateState(m.chain)
    }
}

func (m *CacheMaster) syncCaches(ms int) {
    for {
        // for group := range(m.hash.getGroups()){
        // }
        for groupId := 0; groupId < m.hash.NumGroups; groupId++ {
            go m.syncGroup(groupId)
        }

        // cast int to duration for multiplication to work
        time.Sleep(time.Duration(ms)*time.Millisecond)
    }
}


