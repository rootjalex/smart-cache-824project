package task

import (
	"sync"
	"time"

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
	items := make([]config.DataType, 0)

	// fetch each item in the group
	for _, itemName := range itemGroup {
		// at the end of a web workload pattern, we wait
		if itemName == config.PATTERN_END_MARKER {
			utils.WaitRandomMillis(config.MIN_PATTERN_WAIT, config.MAX_PATTERN_WAIT)
		} else {
			res := c.fetchItem(itemName)
			items = append(items, res)
		}
	}
	// make client wait to simulate computation
	time.Sleep(config.CLIENT_COMPUTATION_TIME)
	return items
}

func (c *Client) fetchItem(itemName string) config.DataType {
	cacheIds := c.hash.GetCaches(itemName, c.id)
	for _, cacheID := range cacheIds {
		cache := c.cachedIDMap[cacheID]
		item, err := cache.Fetch(itemName)
		if err == nil {
			return item
		}
	}
	// arriving here means all caches are dead
	panic("ALL MY CACHES ARE DEAD. PUSH ME TO THE EDGE")
}
