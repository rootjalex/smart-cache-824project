package cache

import "os"

// ------------------------------ Abstract Base Task

// no need for mutex since only task runs at a time
type AbstractBaseTask struct {
	clients []Client
	files   []*os.File
	master  *CacheMaster
}

func NewAbstractBaseTask(numClients int, numCaches int, replicationFactor int, cacheType CacheType, files []*os.File) *AbstractBaseTask {
	// make clients
	clients := make([]Client, numClients)
	for i := range clients {
		// TODO: implement and call Client constructor
		clients[i] = Client{}
	}

	// make cache master
	cacheMaster := StartTask(clients, numCaches, replicationFactor, cacheType)

	return &AbstractBaseTask{
		clients: clients,
		files:   files,
		master:  cacheMaster,
	}
}

func (w *AbstractBaseTask) Init() {

}

func (w *AbstractBaseTask) registerWithMaster() {

}

// ------------------------------ ML Task

type MLTask struct {
	t *AbstractBaseTask
}

func NewMLTask(clients []Client, files []*os.File) *MLTask {
	ml := &MLTask{}
	ml.aw = NewAbstractBaseTask(clients, files)
	return ml
}
