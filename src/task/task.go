package task

import (
	"../cache"
	"../datastore"
	"../config"
)

// ------------------------------ Abstract Base Task

// no need for mutex since only task runs at a time
type AbstractBaseTask struct {
	clients   []*Client
	datastore *datastore.DataStore
	caches    map[int]*cache.Cache
	hash      *cache.Hash
}

// TODO: datastore instead of files
func NewAbstractBaseTask(numClients int, numCaches int, replicationFactor int, cacheType config.CacheType, cacheSize int, datastore *datastore.DataStore, ms int) *AbstractBaseTask {
	// make clients
	clients := make([]*Client, numClients)
	for i := range clients {
		// TODO: implement and call Client constructor
		clients[i] = &Client{}
		// TODO: set their workloads somehow
	}

	// make cache master
	clientIds := make([]int, len(clients))
	for i := 0; i < len(clients); i++ {
		clientIds[i] = clients[i].GetID()
	}
	caches, hash := StartTask(clientIds, cacheType, cacheSize, numCaches, replicationFactor, datastore, ms)

	// TODO: add chache size

	// ends := cacheMaster.getCacheEnds()
	// assign end e for each client c

	return &AbstractBaseTask{
		clients:   clients,
		datastore: datastore,
		caches:    caches,
		hash:      hash,
	}
}

func (w *AbstractBaseTask) Launch() {

}

// ------------------------------ ML Task

// type MLTask struct {
// 	t *AbstractBaseTask
// }

// func NewMLTask(clients []*Client, files []*os.File) *MLTask {
// 	ml := &MLTask{}
// 	// ml.aw = NewAbstractBaseTask(clients, files)
// 	return ml
// }

// ML
