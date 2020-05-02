package cache

import (
	"sync"
	"time"

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

func Init(id int) *Client {
	c := &Client{}
	c.id = id
	return c
}

func (c *Client) Run() {
	for c.workload.HasNextItemGroup() {
		itemGroup := c.workload.GetNextItemGroup()
		for _, itemName := range itemGroup {
			go c.fetchItem(itemName)
		}
	}
}

func (c *Client) GetID() int {
	return c.id
}

// ----------------------------------------------- UTILS

func (c *Client) fetchItem(itemName string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// c.hash.
}
