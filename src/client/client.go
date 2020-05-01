package client

/************************************************
Client supports
* tracking hits and misses
* gets back a cache master / hash function?
*************************************************/


// placeholder for client struct during development
type Client struct {
    id   int

}

func Init(id int) *Client {
    c := &Client{}
    c.id = id
    return c
}

func (c *Client) GetID() int {
    return c.id
}
