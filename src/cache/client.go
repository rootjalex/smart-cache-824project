package cache

import (
	"sync"
	"time"

	"../datastore"
	"../task"
)

/************************************************
Client supports
* tracking hits and misses
* gets back a cache master / hash function?
*************************************************/

type Client struct {
	mu          sync.Mutex
	cachedIDMap map[int]*Cache
	hash        Hash
	workload    *task.Workload
	id          int
	startTime   time.Time
}

func (c *Client) BootstrapClient(cachedIDMap map[int]*Cache, hash Hash, workload *task.Workload) {
	c.cachedIDMap = cachedIDMap
	c.hash = hash
	c.workload = workload
	c.startTime = time.Now()
}

func Init(id int) *Client {
	c := &Client{}
	c.id = id
	return c
}

func (c *Client) Run() []datastore.DataType {
	fetched := []datastore.DataType{}
	for c.workload.HasNextItemGroup() {
		nextItemGroup := c.workload.GetNextItemGroup()
		fetchedItems := c.fetchItemGroup(nextItemGroup)
		fetched = append(fetched, fetchedItems...)
	}
	return fetched
}

func (c *Client) GetID() int {
	return c.id
}

// ----------------------------------------------- UTILS

func (c *Client) fetchItemGroup(itemGroup []string) []datastore.DataType {
	var wg sync.WaitGroup
	items := make([]datastore.DataType, len(itemGroup))

	// fetch each item in the group asynchronously
	for _, itemName := range itemGroup {
		go func(item string) {
			wg.Add(1)
			res := c.fetchItem(item)
			c.mu.Lock()
			items = append(items, res)
			c.mu.Unlock()
			wg.Done()
		}(itemName)
	}
	// wait for all the fetchers to return
	wg.Wait()
	return items
}

func (c *Client) fetchItem(itemName string) datastore.DataType {
	for _, cacheID := range c.hash.GetCaches(itemName, c.id) {
		item, err := c.cachedIDMap[cacheID].Fetch(itemName)
		if err == nil {
			// utils.DPrintf("Fetched %+v-->%+v")
			return item
		}
	}
	// arriving here means all caches are dead
	panic("ALL MY CACHES ARE DEAD. PUSH ME TO THE EDGE")
}
