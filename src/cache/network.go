package cache

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"strconv"
	"time"
)

type CacheNet struct {
	registeredCacheMaster bool
	registeredCache       bool
}

func (cn *CacheNet) startCacheRPCServer(c *Cache) string {
	if !cn.registeredCache {
		rpc.Register(c)
		cn.registeredCache = true
		rpc.HandleHTTP()
	}
	//l, e := net.Listen("tcp", ":1234")
	sockname := cacheSock(string(c.id))
	fmt.Printf("SOCKNAME: %v\n", sockname)

	os.Remove(sockname)
	l, e := net.Listen("unix", sockname)
	if e != nil {
		log.Fatal("listen error:", e)
	}
	go http.Serve(l, nil)
	return sockname
}

func (cn *CacheNet) startCacheMasterRPCServer(cm *CacheMaster) string {
	if !cn.registeredCacheMaster {
		rpc.Register(cm)
		cn.registeredCacheMaster = true
		rpc.HandleHTTP()
	}
	//l, e := net.Listen("tcp", ":1234")
	sockname := cacheSock(time.Now().String())
	fmt.Printf("SOCKNAME: %v\n", sockname)

	os.Remove(sockname)
	l, e := net.Listen("unix", sockname)
	if e != nil {
		log.Fatal("listen error:", e)
	}
	go http.Serve(l, nil)
	return sockname
}

// ---------------------------------- CACHE MASTER

// assumes there is only one cache master
func cacheMasterSock() string {
	s := "/var/tmp/824-smart-cache-master-"
	s += strconv.Itoa(os.Getuid())
	return s
}

func startCacheMasterRPCServer(m *CacheMaster) {
	rpc.Register(m)
	rpc.HandleHTTP()
	//l, e := net.Listen("tcp", ":1234")
	sockname := cacheMasterSock()
	os.Remove(sockname)
	l, e := net.Listen("unix", sockname)
	if e != nil {
		log.Fatal("listen error:", e)
	}
	go http.Serve(l, nil)
}

// ---------------------------------- CACHE INSTANCE

func cacheSock(cacheUID string) string {
	s := "/var/tmp/824-smart-cache-master-"
	s += cacheUID
	return s
}

func startCacheRPCServer(c *Cache) string {
	rpc.Register(c)
	rpc.HandleHTTP()
	//l, e := net.Listen("tcp", ":1234")
	sockname := cacheSock(time.Now().String())
	fmt.Printf("SOCKNAME: %v\n", sockname)

	os.Remove(sockname)
	l, e := net.Listen("unix", sockname)
	if e != nil {
		log.Fatal("listen error:", e)
	}
	go http.Serve(l, nil)
	return sockname
}

// ---------------------------------- SENDING

//
// send an RPC request to the master, wait for the response.
// usually returns true.
// returns false if something goes wrong.
//
func call(sockname string, rpcname string, args interface{}, reply interface{}) bool {
	// c, err := rpc.DialHTTP("tcp", "127.0.0.1"+":1234")
	c, err := rpc.DialHTTP("unix", sockname)
	if err != nil {
		log.Fatal("dialing:", err)
	}
	defer c.Close()

	err = c.Call(rpcname, args, reply)
	if err == nil {
		return true
	}

	fmt.Println(err)
	return false
}
