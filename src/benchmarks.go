package main

import (
	// "strconv"
	// "testing"

	"./config"
	// "../datastore"
	// "./utils"

	"fmt"
	// "log"

	"./task"
)

// ------------------------------------------------------------
// ------------------------------------------------------------
// ------------------------------------------------------------ BENCHMARK MAKER
// ------------------------------------------------------------
// ------------------------------------------------------------


func RunMLBenchmark(params task.MLParams, cacheParams task.CacheParams, nClients int) {
	fb := task.MakeBaseBatch(params.MinFileSleep, params.MaxFileSleep, params.MinBatchSleep, params.MaxBatchSleep)

	// fills with batches
	task.MakeMLFiles(params.NBatches, params.BatchLength, params.NIterations, fb)

	fb_list := make([]task.FileBatches, nClients)

	for i := 0; i < nClients; i++ {
		fb_list[i] = fb.Copy()
	}
	datastore := task.CreateDatastore(fb_list)

	// Markov run-through
	cacheParams.Datastore = datastore
	cacheParams.CacheType = config.LRU
	runner := task.CreateRunner(fb_list, nClients, cacheParams)

	_, t := runner.Run()
	fmt.Printf("ML Benchmark %v completed in %v with LRU\n", params.Name, t)
	hits, misses, calls := runner.ReportData()
	fmt.Printf("hits: %v, misses: %v, queries: %v\n", hits, misses, calls)

	cacheParams.Datastore = datastore.Copy()
	cacheParams.CacheType = config.MarkovPrefetch
	runnerLRU := task.CreateRunner(fb_list, nClients, cacheParams)

	_, t = runnerLRU.Run()
	fmt.Printf("ML Benchmark %v completed in %v with Markov\n", params.Name, t)
	hits, misses, calls = runnerLRU.ReportData()
	fmt.Printf("hits: %v, misses: %v, queries: %v\n", hits, misses, calls)
}


func RunRandomBenchmark(params task.RandomParams, cacheParams task.CacheParams, nClients int) {

	fb_list := make([]task.FileBatches, nClients)

	for i := 0; i < nClients; i++ {
		fb := task.MakeBaseBatch(params.MinFileSleep, params.MaxFileSleep, params.MinBatchSleep, params.MaxBatchSleep)
		task.MakeRandomFiles(params.NBatches, params.MinBatchLength, params.MaxBatchLength, params.MaxFileCount, fb)
		fb_list[i] = fb.Copy()
	}
	datastore := task.CreateDatastore(fb_list)

	// Markov run-through
	cacheParams.Datastore = datastore
	cacheParams.CacheType = config.MarkovPrefetch
	runner := task.CreateRunner(fb_list, nClients, cacheParams)

	_, t := runner.Run()
	fmt.Printf("Random Benchmark %v completed in %v with MarkovPrefetching\n", params.Name, t)
	hits, misses, calls := runner.ReportData()
	fmt.Printf("hits: %v, misses: %v, queries: %v\n", hits, misses, calls)

	cacheParams.CacheType = config.LRU
	runner = task.CreateRunner(fb_list, nClients, cacheParams)

	_, t = runner.Run()
	fmt.Printf("Random Benchmark %v completed in %v with LRU\n", params.Name, t)
	hits, misses, calls = runner.ReportData()
	fmt.Printf("hits: %v, misses: %v, queries: %v\n", hits, misses, calls)
}


func RunWebBenchmark(params task.WebParams, cacheParams task.CacheParams, nClients int) {

	fb_list := make([]task.FileBatches, nClients)

	for i := 0; i < nClients; i++ {
		fb := task.MakeBaseBatch(params.MinFileSleep, params.MaxFileSleep, params.MinBatchSleep, params.MaxBatchSleep)
		// fills with batches
		task.MakeWebFiles(params.NBatches, params.NPatterns, params.MinPatternLength, params.MaxPatternLength, params.BatchLength, params.MaxFileCount, fb)
		fb_list[i] = fb.Copy()
	}
	datastore := task.CreateDatastore(fb_list)

	// Markov run-through
	cacheParams.Datastore = datastore
	cacheParams.CacheType = config.MarkovPrefetch
	runner := task.CreateRunner(fb_list, nClients, cacheParams)

	_, t := runner.Run()
	fmt.Printf("Web Benchmark %v completed in %v with MarkovPrefetching\n", params.Name, t)
	hits, misses, calls := runner.ReportData()
	fmt.Printf("hits: %v, misses: %v, queries: %v\n", hits, misses, calls)

	cacheParams.CacheType = config.LRU
	runner = task.CreateRunner(fb_list, nClients, cacheParams)

	_, t = runner.Run()
	fmt.Printf("Web Benchmark %v completed in %v with LRU\n", params.Name, t)
	hits, misses, calls = runner.ReportData()
	fmt.Printf("hits: %v, misses: %v, queries: %v\n", hits, misses, calls)
}

func TestSmallRandomTask() {
	params := task.RandomParams {
		MinFileSleep: 10,
		MaxFileSleep: 15,
		MinBatchSleep: 1,
		MaxBatchSleep: 2,
		NBatches: 300,
		MinBatchLength: 1,
		MaxBatchLength: 2,
        MaxFileCount: 300,
        BatchLength: 1,
		NIterations: 10,
		Name: "Small",
	}
	cacheParams := task.CacheParams {
		NCaches: 6,
		RFactor: config.RFACTOR_SMALL,
		CacheSize: config.CACHE_SIZE,
		Sync_time: 100,
	}
	nClients := 10
	RunRandomBenchmark(params, cacheParams, nClients)
}



func TestSmallMLTask() {
	params := task.MLParams {
		MinFileSleep: 5,
		MaxFileSleep: 8,
		MinBatchSleep: 10,
		MaxBatchSleep: 15,
		NBatches: 100,
		BatchLength: 1,
		NIterations: 10,
		Name: "Small",
	}
	cacheParams := task.CacheParams {
		NCaches: 6,
		RFactor: config.RFACTOR_SMALL,
		CacheSize: config.CACHE_SIZE,
		Sync_time: 100,
	}
	nClients := 10
	RunMLBenchmark(params, cacheParams, nClients)
}

func TestSmallWebTask() {
	params := task.WebParams {
		MinFileSleep: config.MIN_PATTERN_WAIT,
		MaxFileSleep: config.MAX_PATTERN_WAIT,
		MinBatchSleep: 20,
        MaxBatchSleep: 25,
		NBatches: 300,
        NPatterns: config.NUM_PATTERNS_SMALL,
        MinPatternLength: config.MIN_PATTERN_LENGTH,
        MaxPatternLength: config.MAX_PATTERN_LENGTH,
        MaxFileCount: 300,
		BatchLength: 1,
		Name: "Small Web",
	}
	cacheParams := task.CacheParams {
		NCaches: 2,
		RFactor: config.RFACTOR_SMALL,
		CacheSize: config.CACHE_SIZE,
		Sync_time: 100,
	}
    nClients := 16
	RunWebBenchmark(params, cacheParams, nClients)
}

func web() {
	TestSmallWebTask()
	return
}

func random() {
	TestSmallRandomTask()
	return
}

func ml() {
	TestSmallMLTask()
	return
}

func main() {
	// Web Benchmarks
    web()

	// Random Benchmarks
    random()


	// ML Benchmarks
	ml()
}
