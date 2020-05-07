package main

import (
	// "strconv"
	// "testing"

	"./config"
	// "../datastore"
	"./utils"

	"log"
	"fmt"
	"./task"
)


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

func TestSmallMLTaskMarkov() {
	fmt.Println("TestSmallMLTaskMarkov...")

	// Datastore
	numFiles := 200

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

func TestSmallMLTaskLRU() {
	fmt.Println("TestSmallMLTaskLRU...")

	// Datastore
	numFiles := 200

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

func TestModestMLTaskLRU() {
	fmt.Println("TestModestMLTaskLRU...")
	// Datastore
	numFiles := 500

	// ML parameters
	batchSize := 16
	numIterations := 10

	// Task parameters
	numClients := 10
	numCaches := 4
	replicationFactor := 2
	cacheType := config.LRU
	cacheSize := config.CACHE_SIZE
	ms := 100

	MakeMLBenchmark(cacheType, numFiles, batchSize, numIterations, numClients, numCaches, replicationFactor, cacheSize, ms)
}

func TestModestMLTaskMarkov() {
	fmt.Println("TestModestMLTaskMarkov...")
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

func main() {
	TestSmallMLTaskMarkov()
	TestSmallMLTaskLRU()
	TestModestMLTaskMarkov()
	TestModestMLTaskLRU()
	TestMediumMLTaskMarkov()
	TestMediumMLTaskLRU()
}