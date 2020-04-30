package cache

import "os"

// ------------------------------ Abstract Base Task

// no need for mutex since only task runs at a time
type AbstractBaseTask struct {
	clients []Client
	files   []*os.File
	master  *CacheMaster
}

// TODO: datastore instead of files
func NewAbstractBaseTask(numClients int, numCaches int, replicationFactor int, cacheType CacheType, cacheSize int, files []*os.File) *AbstractBaseTask {
	// make clients
	clients := make([]Client, numClients)
	for i := range clients {
		// TODO: implement and call Client constructor
		clients[i] = Client{}
		// TODO: set their workloads somehow
	}

	// make cache master
	cacheMaster := StartTask(clients, numCaches, replicationFactor, cacheType)

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

func NewMLTask(clients []Client, files []*os.File) *MLTask {
	ml := &MLTask{}
	// ml.aw = NewAbstractBaseTask(clients, files)
	return ml
}

// ML
