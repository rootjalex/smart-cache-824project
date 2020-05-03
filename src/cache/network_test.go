package cache

type NetTestConfig struct {
	cn CacheNet
}

var CONFIG NetTestConfig = NetTestConfig{
	cn: CacheNet{},
}

// func TestBasicNetworkCall(t *testing.T) {
// 	fmt.Printf("TestBasicNetworkCall ...\n")
// 	failed := false
// 	data := datastore.MakeDataStore()

// 	// add files to datastore
// 	for j := 0; j < 5; j++ {
// 		filename := "fake_" + strconv.Itoa(j) + ".txt"
// 		data.Make(filename)
// 	}

// 	// this copies data, so can't adjust later
// 	var lruCache Cache
// 	lruCache.Init(1, CACHE_SIZE, LRU, data)
// 	sockname := CONFIG.cn.startCacheRPCServer(&lruCache)

// 	for j := 0; j < 5; j++ {
// 		filename := "fake_" + strconv.Itoa(j) + ".txt"
// 		args := RequestFileArgs{Filename: filename}
// 		reply := RequestFileReply{}

// 		ok := call(sockname, "Cache.FetchRPC", &args, &reply)
// 		if !ok {
// 			t.Errorf("Could not open %s from cache", filename)
// 			failed = true
// 		}
// 	}

// 	if failed {
// 		fmt.Printf("\t... FAILED\n")
// 	} else {
// 		fmt.Printf("\t... PASSED\n")
// 	}
// }

// func TestDoubleNetworkCall(t *testing.T) {
// 	fmt.Printf("TestDoubleNetworkCall ...\n")
// 	failed := false
// 	data := datastore.MakeDataStore()

// 	// add files to datastore
// 	for j := 0; j < 5; j++ {
// 		filename := "fake_" + strconv.Itoa(j) + ".txt"
// 		data.Make(filename)
// 	}

// 	// this copies data, so can't adjust later
// 	var lruCacheFirst Cache
// 	lruCacheFirst.Init(1, CACHE_SIZE, LRU, data)
// 	socknameFirst := CONFIG.cn.startCacheRPCServer(&lruCacheFirst)

// 	var lruCacheSecond Cache
// 	lruCacheSecond.Init(2, CACHE_SIZE, LRU, data)
// 	socknameSecond := CONFIG.cn.startCacheRPCServer(&lruCacheSecond)

// 	for j := 0; j < 5; j++ {
// 		filename := "fake_" + strconv.Itoa(j) + ".txt"
// 		args := RequestFileArgs{Filename: filename}
// 		reply := RequestFileReply{}

// 		// first cache
// 		okFirst := call(socknameFirst, "Cache.FetchRPC", &args, &reply)

// 		// second cache
// 		okSecond := call(socknameSecond, "Cache.FetchRPC", &args, &reply)

// 		if !okFirst {
// 			t.Errorf("Could not open %s from cache 1", filename)
// 			failed = true
// 		}
// 		if !okSecond {
// 			t.Errorf("Could not open %s from cache 2", filename)
// 			failed = true
// 		}
// 	}

// 	if failed {
// 		fmt.Printf("\t... FAILED\n")
// 	} else {
// 		fmt.Printf("\t... PASSED\n")
// 	}
// }
