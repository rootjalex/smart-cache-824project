package cache

import (
	"os"

	"./cache"
	"./datastore"
)

// ------------------------------ Abstract Base Task

// no need for mutex since only task runs at a time
type AbstractBaseTask struct {
	clients   []*cache.Client
	datastore *datastore.DataStore
	master    *cache.CacheMaster
}

// TODO: datastore instead of files
func NewAbstractBaseTask(numClients int, numCaches int, replicationFactor int, cacheType cache.CacheType, cacheSize int, datastore *datastore.DataStore, ms int) *AbstractBaseTask {
	// make clients
	clients := make([]*cache.Client, numClients)
	for i := range clients {
		// TODO: implement and call Client constructor
		clients[i] = &cache.Client{}
		// TODO: set their workloads somehow
	}

	// make cache master
	cacheMaster := cache.StartTask(clients, cacheType, cacheSize, numCaches, replicationFactor, datastore, ms)

	// TODO: add chache size

	// ends := cacheMaster.getCacheEnds()
	// assign end e for each client c

	return &AbstractBaseTask{
		clients:   clients,
		datastore: datastore,
		master:    cacheMaster,
	}
}

func (w *AbstractBaseTask) Launch() {

}

// ------------------------------ ML Task

type MLTask struct {
	t *AbstractBaseTask
}

func NewMLTask(clients []cache.Client, files []*os.File) *MLTask {
	ml := &MLTask{}
	// ml.aw = NewAbstractBaseTask(clients, files)
	return ml
}

// ML
