package cache

import (
	"fmt"
	"log"
	"strconv"
	"testing"

	"../datastore"
	"../task"
)

func TestClientSimpleWorkload(t *testing.T) {
	fmt.Println("TestClientSimpleWorkload ...")

	numFiles := 100
	numClients := 1
	numCaches := 1
	syncCachesEveryMS := 100
	replicationFactor := 1

	// make datastore
	data := datastore.MakeDataStore()
	for i := 0; i < numFiles; i++ {
		filename := "fake_" + strconv.Itoa(i) + ".txt"
		data.Make(filename)
	}

	// make basic workload
	fileNames := data.GetFileNames()
	itemGroupIndices := make([][]int, len(fileNames))
	for i := 0; i < len(fileNames); i++ {
		itemGroupIndices[i] = []int{i}
	}
	w := &task.Workload{ItemNames: fileNames, ItemGroupIndices: itemGroupIndices}

	// make clients backbone
	clients := make([]*Client, numClients)
	for i := 0; i < numClients; i++ {
		clients[i] = Init(i)
	}
	cachedIDMap, hash := StartTask(clients, LRU, CACHE_SIZE, numCaches, replicationFactor, data, syncCachesEveryMS)
	for i := 0; i < numClients; i++ {
		clients[i].BootstrapClient(cachedIDMap, *hash, w)
	}

	for _, c := range clients {
		log.Println(c.Run())
	}
}
