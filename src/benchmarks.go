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

type MLParams struct {
	minFileSleep	int
	maxFileSleep	int
	minBatchSleep	int
	maxBatchSleep	int
	nBatches 		int 
	batchLength 	int
	nIterations 	int
	name 			string
}

func RunMLBenchmark(params MLParams, cacheParams task.CacheParams, nClients int) {
	fb := task.MakeBaseBatch(params.minFileSleep, params.maxFileSleep, params.minBatchSleep, params.maxBatchSleep)

	// fills with batches
	task.MakeMLFiles(params.nBatches, params.batchLength, params.nIterations, fb)

	fb_list := make([]task.FileBatches, nClients)

	for i := 0; i < nClients; i++ {
		fb_list[i] = fb.Copy()
	}
	datastore := task.CreateDatastore(fb_list)

	// Markov run-through
	cacheParams.Datastore = datastore
	cacheParams.CacheType = config.LRU
	// fmt.Printf("%v", fb_list)
	// fmt.Printf("%v\n", cacheParams)
	runner := task.CreateRunner(fb_list, nClients, cacheParams)

	_, t := runner.Run()
	fmt.Printf("ML Benchmark %v completed in %v with LRU\n", params.name, t)
	hits, misses, calls := runner.ReportData()
	fmt.Printf("hits: %v, misses: %v, queries: %v\n", hits, misses, calls)

	cacheParams.Datastore = datastore.Copy()
	cacheParams.CacheType = config.MarkovPrefetch
	runnerLRU := task.CreateRunner(fb_list, nClients, cacheParams)

	_, t = runnerLRU.Run()
	fmt.Printf("ML Benchmark %v completed in %v with Markov\n", params.name, t)
	hits, misses, calls = runnerLRU.ReportData()
	fmt.Printf("hits: %v, misses: %v, queries: %v\n", hits, misses, calls)
}

type RandomParams struct {
	minFileSleep	int
	maxFileSleep	int
	minBatchSleep	int
	maxBatchSleep	int
	nBatches 		int
	minBatchLength	int
	maxBatchLength	int
	maxFileCount 	int
	batchLength 	int
	nIterations 	int
	name 			string
}

func RunRandomBenchmark(params RandomParams, cacheParams task.CacheParams, nClients int) {

	fb_list := make([]task.FileBatches, nClients)

	for i := 0; i < nClients; i++ {
		fb := task.MakeBaseBatch(params.minFileSleep, params.maxFileSleep, params.minBatchSleep, params.maxBatchSleep)
		task.MakeRandomFiles(params.nBatches, params.minBatchLength, params.maxBatchLength, params.maxFileCount, fb)
		fb_list[i] = fb.Copy()
	}
	datastore := task.CreateDatastore(fb_list)

	// Markov run-through
	cacheParams.Datastore = datastore
	cacheParams.CacheType = config.MarkovPrefetch
	runner := task.CreateRunner(fb_list, nClients, cacheParams)

	_, t := runner.Run()
	fmt.Printf("Random Benchmark %v completed in %v with MarkovPrefetching\n", params.name, t)

	cacheParams.CacheType = config.LRU
	runner = task.CreateRunner(fb_list, nClients, cacheParams)

	_, t = runner.Run()
	fmt.Printf("Random Benchmark %v completed in %v with LRU\n", params.name, t)
}

type WebParams struct {
	minFileSleep		int
	maxFileSleep		int
	minBatchSleep		int
	maxBatchSleep		int
	nBatches 			int
	nPatterns 			int
	minPatternLength	int
	maxPatternLength	int
	maxFileCount 		int
	batchLength 		int
	name 				string
}

func RunWebBenchmark(params WebParams, cacheParams task.CacheParams, nClients int) {

	fb_list := make([]task.FileBatches, nClients)

	for i := 0; i < nClients; i++ {
		fb := task.MakeBaseBatch(params.minFileSleep, params.maxFileSleep, params.minBatchSleep, params.maxBatchSleep)
		// fills with batches
		task.MakeWebFiles(params.nBatches, params.nPatterns, params.minPatternLength, params.maxPatternLength, params.batchLength, params.maxFileCount, fb)
		fb_list[i] = fb.Copy()
	}
	datastore := task.CreateDatastore(fb_list)

	// Markov run-through
	cacheParams.Datastore = datastore
	cacheParams.CacheType = config.MarkovPrefetch
	runner := task.CreateRunner(fb_list, nClients, cacheParams)

	_, t := runner.Run()
	fmt.Printf("Web Benchmark %v completed in %v with MarkovPrefetching\n", params.name, t)

	cacheParams.CacheType = config.LRU
	runner = task.CreateRunner(fb_list, nClients, cacheParams)

	_, t = runner.Run()
	fmt.Printf("Web Benchmark %v completed in %v with LRU\n", params.name, t)
}
// func MakeMLBenchmark(cacheType config.CacheType, nFiles int, batchSize int, nIterations int, nClients int, nCaches int, rFactor, cacheSize int, ms int) {
// 	failed := false
// 	datastore, _, _, fileContents := task.MakeDatastore(nFiles)
// 	// make and launch new ML task
// 	mlTask := task.NewMLTask(batchSize, nIterations, nClients, nCaches, rFactor, cacheType, cacheSize, datastore, ms)
// 	clientFetchMap, taskDuration := mlTask.Launch()
// 	fmt.Printf("\tTask Duration: %+v\n", taskDuration)

// 	// check that all files fetched per client are the expected files
// 	for clientID, fetchedFiles := range clientFetchMap {
// 		repeatedFileContents := utils.DataTypeSliceExtendMany(fileContents, nIterations)
// 		if !utils.DataTypeArraySetsEqual(fetchedFiles, repeatedFileContents) {
// 			log.Printf("Fetched file contents for cleint %v does not match datastore file contents", clientID)
// 			failed = true
// 		}
// 	}
// 	task.PrintFailure(failed)
// }

// func MakeRandomBenchmark(cacheType config.CacheType, nFiles int, batchSize int, nClients int, nCaches int, rFactor, cacheSize int, ms int) {
// 	failed := false
// 	datastore, _, _, fileContents := task.MakeDatastore(nFiles)
// 	// make and launch new random task
// 	mlTask := task.NewRandomTask(batchSize, nClients, nCaches, rFactor, cacheType, cacheSize, datastore, ms)
// 	clientFetchMap, taskDuration := mlTask.Launch()
// 	fmt.Printf("\tTask Duration: %+v\n", taskDuration)

// 	// check that all files fetched per client are the expected files
// 	for clientID, fetchedFiles := range clientFetchMap {
// 		repeatedFileContents := utils.DataTypeSliceExtendMany(fileContents, 1)
// 		if !utils.DataTypeArraySetsEqual(fetchedFiles, repeatedFileContents) {
// 			log.Printf("Fetched file contents for cleint %v does not match datastore file contents", clientID)
// 			failed = true
// 		}
// 	}
// 	task.PrintFailure(failed)
// }

// func MakeWebBenchmark(cacheType config.CacheType, numPatterns int, patternReplication int, nFiles int, batchSize int, nClients int, nCaches int, rFactor, cacheSize int, ms int) {
// 	failed := false
// 	datastore, _, _, _ := task.MakeDatastore(nFiles)
// 	// make and launch new web task
// 	webTask := task.NewWebTask(numPatterns, config.MIN_PATTERN_LENGTH, config.MAX_PATTERN_LENGTH, patternReplication, nClients, nCaches, rFactor, cacheType, cacheSize, datastore, ms)
// 	_, taskDuration := webTask.Launch()
// 	fmt.Printf("\tTask Duration: %+v\n", taskDuration)
// 	task.PrintFailure(failed)
// }

// ------------------------------------------------------------
// ------------------------------------------------------------
// ------------------------------------------------------------ WEB
// ------------------------------------------------------------
// ------------------------------------------------------------

func TestSmallMLTask() {
	params := MLParams {
		minFileSleep: 5,
		maxFileSleep: 10,
		minBatchSleep: 10,
		maxBatchSleep: 20, 
		nBatches: 5,
		batchLength: 40,
		nIterations: 20,
		name: "Small",
	}
	cacheParams := task.CacheParams {
		NCaches: 2,
		RFactor: config.RFACTOR_SMALL,
		CacheSize: config.CACHE_SIZE,
		Sync_time: 100,
	}
	nClients := 10
	RunMLBenchmark(params, cacheParams, nClients)
}
// func TestSmallWebTaskLRU() {
// 	fmt.Println("TestSmallWebTaskLRU...")
// 	MakeWebBenchmark(
// 		config.LRU, // LRU
// 		config.NUM_PATTERNS_SMALL,
// 		config.PATTERN_REPLICATION_SMALL,
// 		config.NFILES_SMALL,
// 		config.BATCH_SMALL,
// 		config.NCLIENTS_SMALL,
// 		config.NCACHES_SMALL,
// 		config.RFACTOR_SMALL,
// 		config.CACHE_SIZE,
// 		config.SYNC_MS,
// 	)
// }

// func TestSmallWebTaskMarkov() {
// 	fmt.Println("TestSmallWebTaskMarkov...")
// 	MakeWebBenchmark(
// 		config.MarkovPrefetch, // MARKOV
// 		config.NUM_PATTERNS_SMALL,
// 		config.PATTERN_REPLICATION_SMALL,
// 		config.NFILES_SMALL,
// 		config.BATCH_SMALL,
// 		config.NCLIENTS_SMALL,
// 		config.NCACHES_SMALL,
// 		config.RFACTOR_SMALL,
// 		config.CACHE_SIZE,
// 		config.SYNC_MS,
// 	)
// }

// func TestMediumWebTaskLRU() {
// 	fmt.Println("TestMediumWebTaskLRU...")
// 	MakeWebBenchmark(
// 		config.LRU, // LRU
// 		config.NUM_PATTERNS_MED,
// 		config.PATTERN_REPLICATION_MED,
// 		config.NFILES_MED,
// 		config.BATCH_MED,
// 		config.NCLIENTS_MED,
// 		config.NCACHES_MED,
// 		config.RFACTOR_MED,
// 		config.CACHE_SIZE,
// 		config.SYNC_MS,
// 	)
// }

// func TestMediumWebTaskMarkov() {
// 	fmt.Println("TestMediumWebTaskMarkov...")
// 	MakeWebBenchmark(
// 		config.MarkovPrefetch, // MARKOV
// 		config.NUM_PATTERNS_MED,
// 		config.PATTERN_REPLICATION_MED,
// 		config.NFILES_MED,
// 		config.BATCH_MED,
// 		config.NCLIENTS_MED,
// 		config.NCACHES_MED,
// 		config.RFACTOR_MED,
// 		config.CACHE_SIZE,
// 		config.SYNC_MS,
// 	)
// }

// // ------------------------------------------------------------
// // ------------------------------------------------------------
// // ------------------------------------------------------------ RANDOM
// // ------------------------------------------------------------
// // ------------------------------------------------------------

// func TestSmallRandomTaskLRU() {
// 	fmt.Println("TestSmallRandomTaskLRU...")
// 	MakeRandomBenchmark(
// 		config.LRU, // LRU
// 		config.NFILES_SMALL,
// 		config.BATCH_SMALL,
// 		config.NCLIENTS_SMALL,
// 		config.NCACHES_SMALL,
// 		config.RFACTOR_SMALL,
// 		config.CACHE_SIZE,
// 		config.SYNC_MS,
// 	)
// }

// func TestSmallRandomTaskMarkov() {
// 	fmt.Println("TestSmallRandomTaskMarkov...")
// 	MakeRandomBenchmark(
// 		config.MarkovPrefetch, // MARKOV
// 		config.NFILES_SMALL,
// 		config.BATCH_SMALL,
// 		config.NCLIENTS_SMALL,
// 		config.NCACHES_SMALL,
// 		config.RFACTOR_SMALL,
// 		config.CACHE_SIZE,
// 		config.SYNC_MS,
// 	)
// }

// func TestMediumRandomTaskLRU() {
// 	fmt.Println("TestMediumRandomTaskLRU...")
// 	MakeRandomBenchmark(
// 		config.LRU, // LRU
// 		config.NFILES_MED,
// 		config.BATCH_MED,
// 		config.NCLIENTS_MED,
// 		config.NCACHES_MED,
// 		config.RFACTOR_MED,
// 		config.CACHE_SIZE,
// 		config.SYNC_MS,
// 	)
// }

// func TestMediumRandomTaskMarkov() {
// 	fmt.Println("TestMediumRandomTaskMarkov...")
// 	MakeRandomBenchmark(
// 		config.MarkovPrefetch, // MARKOV
// 		config.NFILES_MED,
// 		config.BATCH_MED,
// 		config.NCLIENTS_MED,
// 		config.NCACHES_MED,
// 		config.RFACTOR_MED,
// 		config.CACHE_SIZE,
// 		config.SYNC_MS,
// 	)
// }

// // ------------------------------------------------------------
// // ------------------------------------------------------------
// // ------------------------------------------------------------ ML
// // ------------------------------------------------------------
// // ------------------------------------------------------------

// func TestSmallMLTaskMarkov() {
// 	fmt.Println("TestSmallMLTaskMarkov...")
// 	MakeMLBenchmark(
// 		config.MarkovPrefetch, // MARKOV
// 		config.NFILES_SMALL,
// 		config.BATCH_SMALL,
// 		config.ITERS_SMALL,
// 		config.NCLIENTS_SMALL,
// 		config.NCACHES_SMALL,
// 		config.RFACTOR_SMALL,
// 		config.CACHE_SIZE,
// 		config.SYNC_MS,
// 	)
// }

// func TestSmallMLTaskLRU() {
// 	fmt.Println("TestSmallMLTaskLRU...")
// 	MakeMLBenchmark(
// 		config.LRU, // LRU
// 		config.NFILES_SMALL,
// 		config.BATCH_SMALL,
// 		config.ITERS_SMALL,
// 		config.NCLIENTS_SMALL,
// 		config.NCACHES_SMALL,
// 		config.RFACTOR_SMALL,
// 		config.CACHE_SIZE,
// 		config.SYNC_MS,
// 	)
// }

// func TestMediumMLTaskMarkov() {
// 	fmt.Println("TestMediumMLTaskMarkov...")
// 	MakeMLBenchmark(
// 		config.MarkovPrefetch, // MARKOV
// 		config.NFILES_MED,
// 		config.BATCH_MED,
// 		config.ITERS_MED,
// 		config.NCLIENTS_MED,
// 		config.NCACHES_MED,
// 		config.RFACTOR_MED,
// 		config.CACHE_SIZE,
// 		config.SYNC_MS,
// 	)
// }

// func TestMediumMLTaskLRU() {
// 	fmt.Println("TestMediumMLTaskLRU...")
// 	MakeMLBenchmark(
// 		config.LRU, // LRU
// 		config.NFILES_MED,
// 		config.BATCH_MED,
// 		config.ITERS_MED,
// 		config.NCLIENTS_MED,
// 		config.NCACHES_MED,
// 		config.RFACTOR_MED,
// 		config.CACHE_SIZE,
// 		config.SYNC_MS,
// 	)
// }

// func TestMediumLowRepMLTaskMarkov() {
// 	fmt.Println("TestMediumLowRepMLTaskMarkov...")
// 	MakeMLBenchmark(
// 		config.MarkovPrefetch, // MARKOV
// 		config.NFILES_MED,
// 		config.BATCH_MED,
// 		config.ITERS_MED,
// 		config.NCLIENTS_MED,
// 		config.NCACHES_MED,
// 		config.RFACTOR_SMALL,
// 		config.CACHE_SIZE,
// 		config.SYNC_MS,
// 	)
// }

// func TestMediumLowRepMLTaskLRU() {
// 	fmt.Println("TestMediumLowRepMLTaskLRU...")
// 	MakeMLBenchmark(
// 		config.LRU, // LRU
// 		config.NFILES_MED,
// 		config.BATCH_MED,
// 		config.ITERS_MED,
// 		config.NCLIENTS_MED,
// 		config.NCACHES_MED,
// 		config.RFACTOR_SMALL,
// 		config.CACHE_SIZE,
// 		config.SYNC_MS,
// 	)
// }

// ------------------------------------------------------------
// ------------------------------------------------------------
// ------------------------------------------------------------ RANDOM
// ------------------------------------------------------------
// ------------------------------------------------------------
func web() {
	// TestSmallWebTaskLRU()
	// TestSmallWebTaskMarkov()
	// TestMediumWebTaskLRU()
	// TestMediumWebTaskMarkov()
	return
}

func random() {
	// TestSmallRandomTaskLRU()
	// TestSmallRandomTaskMarkov()
	// TestMediumRandomTaskLRU()
	// TestMediumRandomTaskMarkov()
	return
}

func ml() {
	// TestSmallMLTaskMarkov()
	// TestSmallMLTaskLRU()
	// TestMediumMLTaskMarkov()
	// TestMediumMLTaskLRU()
	// TestMediumLowRepMLTaskMarkov()
	// TestMediumLowRepMLTaskLRU()
	return
}

func main() {
	// Web Benchmarks
    // web()

	// // Random Benchmarks
    // random()

	// ML Benchmarks
	// ml()
	TestSmallMLTask()
}
