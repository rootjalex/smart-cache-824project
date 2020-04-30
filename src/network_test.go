package cache

import (
	"fmt"
	"strconv"
	"testing"

	"./cache"
	"./datastore"
)

type NetTestConfig struct {
	cn CacheNet
}

var CONFIG NetTestConfig = NetTestConfig{
	cn: CacheNet{},
}

func TestBasicNetworkCall(t *testing.T) {
	fmt.Printf("TestBasicNetworkCall ...\n")
	failed := false
	data := datastore.MakeDataStore()

	// add files to datastore
	for j := 0; j < 5; j++ {
		filename := "fake_" + strconv.Itoa(j) + ".txt"
		data.Make(filename)
	}

	// this copies data, so can't adjust later
	var lruCache cache.Cache
	lruCache.Init(cache.CACHE_SIZE, cache.LRU, data)
	sockname := CONFIG.cn.startCacheRPCServer(&lruCache)

	for j := 0; j < 5; j++ {
		filename := "fake_" + strconv.Itoa(j) + ".txt"
		args := cache.RequestFileArgs{Filename: filename}
		reply := cache.RequestFileReply{}

		ok := call(sockname, "Cache.FetchRPC", &args, &reply)
		if !ok {
			t.Errorf("Could not open %s from cache", filename)
			failed = true
		}
	}

	if failed {
		fmt.Printf("\t... FAILED\n")
	} else {
		fmt.Printf("\t... PASSED\n")
	}
}

func TestDoubleNetworkCall(t *testing.T) {
	fmt.Printf("TestDoubleNetworkCall ...\n")
	failed := false
	data := datastore.MakeDataStore()

	// add files to datastore
	for j := 0; j < 5; j++ {
		filename := "fake_" + strconv.Itoa(j) + ".txt"
		data.Make(filename)
	}

	// this copies data, so can't adjust later
	var lruCacheFirst cache.Cache
	lruCacheFirst.Init(cache.CACHE_SIZE, cache.LRU, data)
	socknameFirst := CONFIG.cn.startCacheRPCServer(&lruCacheFirst)

	var lruCacheSecond cache.Cache
	lruCacheSecond.Init(cache.CACHE_SIZE, cache.LRU, data)
	socknameSecond := CONFIG.cn.startCacheRPCServer(&lruCacheSecond)

	for j := 0; j < 5; j++ {
		filename := "fake_" + strconv.Itoa(j) + ".txt"
		args := cache.RequestFileArgs{Filename: filename}
		reply := cache.RequestFileReply{}

		// first cache
		okFirst := call(socknameFirst, "Cache.FetchRPC", &args, &reply)

		// second cache
		okSecond := call(socknameSecond, "Cache.FetchRPC", &args, &reply)

		if !okFirst {
			t.Errorf("Could not open %s from cache 1", filename)
			failed = true
		}
		if !okSecond {
			t.Errorf("Could not open %s from cache 2", filename)
			failed = true
		}
	}

	if failed {
		fmt.Printf("\t... FAILED\n")
	} else {
		fmt.Printf("\t... PASSED\n")
	}
}
