package cache

import (
	"os"
	"./client"
	"./cache"
)

// ------------------------------ Abstract Base Task

// no need for mutex since only task runs at a time
type AbstractBaseTask struct {
	clients []client.Client
	files   []*os.File
	master  *client.CacheMaster
}

// TODO: datastore instead of files
func NewAbstractBaseTask(numClients int, numCaches int, replicationFactor int, cacheType cache.CacheType, cacheSize int, files []*os.File) *AbstractBaseTask {
	// make clients
	clients := make([]client.Client, numClients)
	for i := range clients {
		// TODO: implement and call Client constructor
		clients[i] = client.Client{}
		// TODO: set their workloads somehow
	}

	// make cache master
	cacheMaster := client.StartTask(clients, numCaches, replicationFactor, cacheType)

	// TODO: add chache size

	// ends := cacheMaster.getCacheEnds()
	// assign end e for each client c

	return &AbstractBaseTask{
		clients: clients,
		files:   files,
		master:  cacheMaster,
	}
}

func (w *AbstractBaseTask) Launch() {

}

// ------------------------------ ML Task

type MLTask struct {
	t *AbstractBaseTask
}

// TODO: use datastore.DataType instead of os.File
func NewMLTask(clients []client.Client, files []*os.File) *MLTask {
	ml := &MLTask{}
	// ml.aw = NewAbstractBaseTask(clients, files)
	return ml
}

// ML
