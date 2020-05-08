package task

import (
	"sync"
	"time"

	"../cache"
	"../config"
	// "../utils"
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
	workload    *Workload
	id          int
	startTime   time.Time
}

func (c *Client) BootstrapClient(cachedIDMap map[int]*cache.Cache, hash *cache.Hash, workload *Workload) {
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

func (c *Client) Run() []config.DataType {
	// log.Printf("client running. %+v", c)
	fetched := []config.DataType{}
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

func (c *Client) fetchItemGroup(itemGroup []string) []config.DataType {
	// var wg sync.WaitGroup
	items := make([]config.DataType, 0)

	// fetch each item in the group
	for _, itemName := range itemGroup {
		res := c.fetchItem(itemName)
		items = append(items, res)
	}
	// make client wait to simulate computation
	time.Sleep(config.CLIENT_COMPUTATION_TIME)
	return items
}

func (c *Client) fetchItem(itemName string) config.DataType {
	cacheIds := c.hash.GetCaches(itemName, c.id)
	for _, cacheID := range cacheIds {
		cache := c.cachedIDMap[cacheID]
		// log.Printf("cache: %v", cache)
		item, err := cache.Fetch(itemName)
		if err == nil {
			return item
		}
	}
	// arriving here means all caches are dead
	panic("ALL MY CACHES ARE DEAD. PUSH ME TO THE EDGE")
}
