package task

import (
	"sync"
	"time"
	"log"
	"../cache"
	"../config"
    "../utils"
)

/************************************************
Client supports
* tracking hits and misses
* gets back a cache master / hash function?
*************************************************/

type Client struct {
	mu          sync.Mutex
	cachedIDMap map[int]*cache.Cache
	hash        *cache.Hash
	// workload    *Workload
	files		FileBatches
	id          int
	startTime   time.Time
}

func (c *Client) BootstrapClient(cachedIDMap map[int]*cache.Cache, hash *cache.Hash, files FileBatches) {
	c.cachedIDMap = cachedIDMap
	c.hash = hash
	c.files = files.Copy()
	// c.workload = workload
	c.startTime = time.Now()
}

func Init(id int) *Client {
	c := &Client{}
	c.id = id
	return c
}

func (c *Client) Run() []config.DataType {
	// totalFiles := 0
	// for _, igi := range c.workload.ItemGroupIndices {
	// 	totalFiles += len(igi)
	// }
	// log.Printf("client %v fetching %+v items", c.id, totalFiles)

	fetched := []config.DataType{}
	for _, batch := range c.files.batches {
		fetchedItems := c.fetchBatch(batch)
		fetched = append(fetched, fetchedItems...)
	}
	return fetched
}

func (c *Client) GetID() int {
	return c.id
}

// ----------------------------------------------- UTILS

func (c *Client) fetchBatch(itemGroup []string) []config.DataType {
	items := make([]config.DataType, 0)

	// fetch each item in the group
	for _, itemName := range itemGroup {
		res := c.fetchItem(itemName)
		if res != config.DataType("good") {
			log.Printf("FAILED TO FETCH %v", itemName)
		}
		items = append(items, res)
		utils.WaitRandomMillis(c.files.minFileSleep, c.files.maxFileSleep)
		// at the end of a web workload pattern, we wait
	}

	utils.WaitRandomMillis(c.files.minBatchSleep, c.files.maxBatchSleep)
	return items
}

func (c *Client) fetchItem(itemName string) config.DataType {
	cacheIds := c.hash.GetCaches(itemName, c.id)
	for _, cacheID := range cacheIds {
		cache := c.cachedIDMap[cacheID]
		item, err := cache.Fetch(itemName, c.GetID())
		if err == nil {
			return item
		}
	}
	// arriving here means all caches are dead
	panic("ALL MY CACHES ARE DEAD. PUSH ME TO THE EDGE")
}
