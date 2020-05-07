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
		config.SYNC_MS_SMALL,
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
		config.SYNC_MS_SMALL,
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
		config.SYNC_MS_SMALL,
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
		config.SYNC_MS_SMALL,
	)
}

func TestMediumMLTaskMarkov() {
	fmt.Println("TestMediumMLTaskMarkov...")

	// Datastore
	numFiles := 1000

	// ML parameters
	batchSize := 16
	numIterations := 10

	// Task parameters
	numClients := 5
	numCaches := 2
	replicationFactor := 1
	cacheType := config.MarkovPrefetch
	cacheSize := config.CACHE_SIZE
	ms := 100

	MakeMLBenchmark(cacheType, numFiles, batchSize, numIterations, numClients, numCaches, replicationFactor, cacheSize, ms)
}

func TestMediumMLTaskLRU() {
	fmt.Println("TestMediumMLTaskLRU...")

	// Datastore
	numFiles := 1000

	// ML parameters
	batchSize := 16
	numIterations := 10

	// Task parameters
	numClients := 5
	numCaches := 2
	replicationFactor := 1
	cacheType := config.LRU
	cacheSize := config.CACHE_SIZE
	ms := 100

	MakeMLBenchmark(cacheType, numFiles, batchSize, numIterations, numClients, numCaches, replicationFactor, cacheSize, ms)
}

func TestSmallerModestMLTaskLRU() {
	fmt.Println("TestSmallerModestMLTaskLRU...")
	// Datastore
	numFiles := 500

	// ML parameters
	batchSize := 16
	numIterations := 10

	// Task parameters
	numClients := 10
	numCaches := 4
	replicationFactor := 2
	cacheType := config.MarkovPrefetch
	cacheSize := config.CACHE_SIZE
	ms := 100

	MakeMLBenchmark(cacheType, numFiles, batchSize, numIterations, numClients, numCaches, replicationFactor, cacheSize, ms)
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

	// ML Benchmarks
	// TestSmallMLTaskMarkov()
	// TestSmallMLTaskLRU()
	// TestMediumMLTaskMarkov()
	// TestMediumMLTaskLRU()
	// TestSmallerModestMLTaskLRU()
}
