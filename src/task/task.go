package task

import (
	"sync"
	"time"

	"../cache"
	"../config"
	"../datastore"
)

// ------------------------------ Abstract Base Task

// no need for mutex since only task runs at a time
type AbstractBaseTask struct {
	mu        sync.Mutex
	Clients   []*Client
	datastore *datastore.DataStore
	caches    map[int]*cache.Cache
	hash      *cache.Hash
}

func NewAbstractBaseTask(workload Workload, numClients int, numCaches int, replicationFactor int, cacheType config.CacheType, cacheSize int, datastore *datastore.DataStore, ms int) *AbstractBaseTask {
	// make clients
	clients := make([]*Client, numClients)
	for i := range clients {
		clients[i] = Init(i)
	}
	// make cache master
	clientIds := make([]int, len(clients))
	for i := 0; i < len(clients); i++ {
		clientIds[i] = clients[i].GetID()
	}
	caches, hash := StartTask(clientIds, cacheType, cacheSize, numCaches, replicationFactor, datastore, ms)

	// bootstrap clients
	for i := range clients {
		w := workload.FreshCopy()
		clients[i].BootstrapClient(caches, hash, &w)
	}

	return &AbstractBaseTask{
		Clients:   clients,
		datastore: datastore,
		caches:    caches,
		hash:      hash,
	}
}

func (t *AbstractBaseTask) Launch() (map[int][]config.DataType, time.Duration) {
	preFetchTime := time.Now()
	clientIDToFetchedFiles := make(map[int][]config.DataType)

	// run all clients in parallel, wait until all are done
	// aggregate client fetch results
	var wg sync.WaitGroup
	for _, c := range t.Clients {
		wg.Add(1)
		go func(client *Client) {
			fetched := client.Run()
			// log.Printf("fetched from client: %+v", fetched)
			t.mu.Lock()
			clientIDToFetchedFiles[client.GetID()] = fetched
			t.mu.Unlock()
			wg.Done()
		}(c)
	}
	wg.Wait()

	return clientIDToFetchedFiles, time.Since(preFetchTime)
}

// ------------------------------ ML Task

type MLTask struct {
	abstractTask *AbstractBaseTask
}

func NewMLTask(batchSize int, numIterations int, numClients int, numCaches int, replicationFactor int, cacheType config.CacheType, cacheSize int, datastore *datastore.DataStore, ms int) *MLTask {
	// make ML workload
	itemNames := datastore.GetFileNames()
	mlWkld := NewMLWorkload(itemNames, batchSize, numIterations)

	// make abstract task
	t := NewAbstractBaseTask(mlWkld, numClients, numCaches, replicationFactor, cacheType, cacheSize, datastore, ms)
	return &MLTask{
		abstractTask: t,
	}
}

func (ml *MLTask) Launch() (map[int][]config.DataType, time.Duration) {
	return ml.abstractTask.Launch()
}
