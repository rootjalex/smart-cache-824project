package task

import (
	"sync"
	"time"

	"../cache"
	"../config"
    "log"
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
	var wg sync.WaitGroup
	items := make([]config.DataType, len(itemGroup))

	// fetch each item in the group asynchronously
	for _, itemName := range itemGroup {
		wg.Add(1)
		go func(item string) {
			res := c.fetchItem(item)
            log.Printf("result of fetch: %v", res)
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

func (c *Client) fetchItem(itemName string) config.DataType {
    cacheIds := c.hash.GetCaches(itemName, c.id)
    log.Printf("inside fetchItem %v can look at cacheIds: %v", itemName, cacheIds)
	for _, cacheID := range cacheIds {
        cache := c.cachedIDMap[cacheID]
        log.Printf("cache: %v", cache)
		item, err := cache.Fetch(itemName)
		if err == nil {
            utils.DPrintf("Fetched %+v-->%+v", itemName, item)
			return item
		}
	}
	// arriving here means all caches are dead
	panic("ALL MY CACHES ARE DEAD. PUSH ME TO THE EDGE")
}
