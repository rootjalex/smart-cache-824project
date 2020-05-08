package main

import (
	// "strconv"
	// "testing"

	"./config"
	// "../datastore"
	"./utils"

	"fmt"
	"log"

	"./task"
)

// ------------------------------------------------------------
// ------------------------------------------------------------
// ------------------------------------------------------------ BENCHMARK MAKER
// ------------------------------------------------------------
// ------------------------------------------------------------

func MakeMLBenchmark(cacheType config.CacheType, nFiles int, batchSize int, nIterations int, nClients int, nCaches int, rFactor, cacheSize int, ms int) {
	failed := false
	datastore, _, _, fileContents := task.MakeDatastore(nFiles)
	// make and launch new ML task
	mlTask := task.NewMLTask(batchSize, nIterations, nClients, nCaches, rFactor, cacheType, cacheSize, datastore, ms)
	clientFetchMap, taskDuration := mlTask.Launch()
	fmt.Printf("\tTask Duration: %+v\n", taskDuration)

	// check that all files fetched per client are the expected files
	for clientID, fetchedFiles := range clientFetchMap {
		repeatedFileContents := utils.DataTypeSliceExtendMany(fileContents, nIterations)
		if !utils.DataTypeArraySetsEqual(fetchedFiles, repeatedFileContents) {
			log.Printf("Fetched file contents for cleint %v does not match datastore file contents", clientID)
			failed = true
		}
	}
	task.PrintFailure(failed)
}

func MakeRandomBenchmark(cacheType config.CacheType, nFiles int, batchSize int, nClients int, nCaches int, rFactor, cacheSize int, ms int) {
	failed := false
	datastore, _, _, fileContents := task.MakeDatastore(nFiles)
	// make and launch new random task
	mlTask := task.NewRandomTask(batchSize, nClients, nCaches, rFactor, cacheType, cacheSize, datastore, ms)
	clientFetchMap, taskDuration := mlTask.Launch()
	fmt.Printf("\tTask Duration: %+v\n", taskDuration)

	// check that all files fetched per client are the expected files
	for clientID, fetchedFiles := range clientFetchMap {
		repeatedFileContents := utils.DataTypeSliceExtendMany(fileContents, 1)
		if !utils.DataTypeArraySetsEqual(fetchedFiles, repeatedFileContents) {
			log.Printf("Fetched file contents for cleint %v does not match datastore file contents", clientID)
			failed = true
		}
	}
	task.PrintFailure(failed)
}

// ------------------------------------------------------------
// ------------------------------------------------------------
// ------------------------------------------------------------ RANDOM
// ------------------------------------------------------------
// ------------------------------------------------------------

func TestSmallRandomTaskLRU() {
	fmt.Println("TestSmallRandomTaskLRU...")
	MakeRandomBenchmark(
		config.LRU, // LRU
		config.NFILES_SMALL,
		config.BATCH_SMALL,
		config.NCLIENTS_SMALL,
		config.NCACHES_SMALL,
		config.RFACTOR_SMALL,
		config.CACHE_SIZE,
		config.SYNC_MS,
	)
}

func TestSmallRandomTaskMarkov() {
	fmt.Println("TestSmallRandomTaskMarkov...")
	MakeRandomBenchmark(
		config.MarkovPrefetch, // MARKOV
		config.NFILES_SMALL,
		config.BATCH_SMALL,
		config.NCLIENTS_SMALL,
		config.NCACHES_SMALL,
		config.RFACTOR_SMALL,
		config.CACHE_SIZE,
		config.SYNC_MS,
	)
}

func TestMediumRandomTaskLRU() {
	fmt.Println("TestMediumRandomTaskLRU...")
	MakeRandomBenchmark(
		config.LRU, // LRU
		config.NFILES_MED,
		config.BATCH_MED,
		config.NCLIENTS_MED,
		config.NCACHES_MED,
		config.RFACTOR_MED,
		config.CACHE_SIZE,
		config.SYNC_MS,
	)
}

func TestMediumRandomTaskMarkov() {
	fmt.Println("TestMediumRandomTaskMarkov...")
	MakeRandomBenchmark(
		config.MarkovPrefetch, // MARKOV
		config.NFILES_MED,
		config.BATCH_MED,
		config.NCLIENTS_MED,
		config.NCACHES_MED,
		config.RFACTOR_MED,
		config.CACHE_SIZE,
		config.SYNC_MS,
	)
}

// ------------------------------------------------------------
// ------------------------------------------------------------
// ------------------------------------------------------------ ML
// ------------------------------------------------------------
// ------------------------------------------------------------

func TestSmallMLTaskMarkov() {
	fmt.Println("TestSmallMLTaskMarkov...")
	MakeMLBenchmark(
		config.MarkovPrefetch, // MARKOV
		config.NFILES_SMALL,
		config.BATCH_SMALL,
		config.ITERS_SMALL,
		config.NCLIENTS_SMALL,
		config.NCACHES_SMALL,
		config.RFACTOR_SMALL,
		config.CACHE_SIZE,
		config.SYNC_MS,
	)
}

func TestSmallMLTaskLRU() {
	fmt.Println("TestSmallMLTaskLRU...")
	MakeMLBenchmark(
		config.LRU, // LRU
		config.NFILES_SMALL,
		config.BATCH_SMALL,
		config.ITERS_SMALL,
		config.NCLIENTS_SMALL,
		config.NCACHES_SMALL,
		config.RFACTOR_SMALL,
		config.CACHE_SIZE,
		config.SYNC_MS,
	)
}

func TestMediumMLTaskMarkov() {
	fmt.Println("TestMediumMLTaskMarkov...")
	MakeMLBenchmark(
		config.MarkovPrefetch, // MARKOV
		config.NFILES_MED,
		config.BATCH_MED,
		config.ITERS_MED,
		config.NCLIENTS_MED,
		config.NCACHES_MED,
		config.RFACTOR_MED,
		config.CACHE_SIZE,
		config.SYNC_MS,
	)
}

func TestMediumMLTaskLRU() {
	fmt.Println("TestMediumMLTaskLRU...")
	MakeMLBenchmark(
		config.LRU, // LRU
		config.NFILES_MED,
		config.BATCH_MED,
		config.ITERS_MED,
		config.NCLIENTS_MED,
		config.NCACHES_MED,
		config.RFACTOR_MED,
		config.CACHE_SIZE,
		config.SYNC_MS,
	)
}

func TestMediumLowRepMLTaskMarkov() {
	fmt.Println("TestMediumLowRepMLTaskMarkov...")
	MakeMLBenchmark(
		config.MarkovPrefetch, // MARKOV
		config.NFILES_MED,
		config.BATCH_MED,
		config.ITERS_MED,
		config.NCLIENTS_MED,
		config.NCACHES_MED,
		config.RFACTOR_SMALL,
		config.CACHE_SIZE,
		config.SYNC_MS,
	)
}

func TestMediumLowRepMLTaskLRU() {
	fmt.Println("TestMediumLowRepMLTaskLRU...")
	MakeMLBenchmark(
		config.LRU, // LRU
		config.NFILES_MED,
		config.BATCH_MED,
		config.ITERS_MED,
		config.NCLIENTS_MED,
		config.NCACHES_MED,
		config.RFACTOR_SMALL,
		config.CACHE_SIZE,
		config.SYNC_MS,
	)
}


// ------------------------------------------------------------
// ------------------------------------------------------------
// ------------------------------------------------------------ RANDOM
// ------------------------------------------------------------
// ------------------------------------------------------------

func main() {
	// Random Benchmarks
	TestSmallRandomTaskLRU()
	TestSmallRandomTaskMarkov()
	TestMediumRandomTaskLRU()
	TestMediumRandomTaskMarkov()

	// ML Benchmarks
	TestSmallMLTaskMarkov()
	TestSmallMLTaskLRU()
	TestMediumMLTaskMarkov()
	TestMediumMLTaskLRU()
	TestMediumLowRepMLTaskMarkov()
	TestMediumLowRepMLTaskLRU()
}
