package task

import (
	"strconv"
	"testing"

	"../config"
	"../datastore"
	"../utils"
)

func TestSmallMLTaskLRU(t *testing.T) {
	t.Log("TestSmallMLTask...")

	// Datastore
	numFiles := 1000
	datastore, _, _, fileContents := makeDatastore(numFiles)

	// ML parameters
	batchSize := 16
	numIterations := 50

	// Task parameters
	numClients := 5
	numCaches := 2
	replicationFactor := 1
	cacheType := config.LRU
	cacheSize := config.CACHE_SIZE
	ms := 100

	// make and launch new ML task
	mlTask := NewMLTask(batchSize, numIterations, numClients, numCaches, replicationFactor, cacheType, cacheSize, datastore, ms)
	clientFetchMap, taskDuration := mlTask.Launch()
	t.Logf("\tTask Duration: %+v,", taskDuration)

	// check that all files fetched per client are the expected files
	for clientID, fetchedFiles := range clientFetchMap {
		repeatedFileContents := utils.DataTypeSliceExtendMany(fileContents, numIterations)
		if !utils.DataTypeArraySetsEqual(fetchedFiles, repeatedFileContents) {
			t.Errorf("Fetched file contents for cleint %v does not match datastore file contents", clientID)
		}
	}
}

func TestModestMLTaskLRU(t *testing.T) {
	t.Log("TestModestMLTask...")

	// Datastore
	numFiles := 1000
	datastore, _, _, fileContents := makeDatastore(numFiles)

	// ML parameters
	batchSize := 32
	numIterations := 200

	// Task parameters
	numClients := 20
	numCaches := 5
	replicationFactor := 2
	cacheType := config.LRU
	cacheSize := config.CACHE_SIZE
	ms := 100

	// make and launch new ML task
	mlTask := NewMLTask(batchSize, numIterations, numClients, numCaches, replicationFactor, cacheType, cacheSize, datastore, ms)
	clientFetchMap, taskDuration := mlTask.Launch()
	t.Logf("\tTask Duration: %+v,", taskDuration)

	// check that all files fetched per client are the expected files
	for clientID, fetchedFiles := range clientFetchMap {
		repeatedFileContents := utils.DataTypeSliceExtendMany(fileContents, numIterations)
		if !utils.DataTypeArraySetsEqual(fetchedFiles, repeatedFileContents) {
			t.Errorf("Fetched file contents for cleint %v does not match datastore file contents", clientID)
		}
	}
}

func makeDatastore(numDatastoreFiles int) (*datastore.DataStore, map[string]config.DataType, []string, []config.DataType) {
	ds := datastore.MakeDataStore()
	files := make(map[string]config.DataType)
	fileNames, fileContents := []string{}, []config.DataType{}

	// add files to the datastore
	for i := 0; i < numDatastoreFiles; i++ {
		fileName := "MNIST_" + strconv.Itoa(i+1) + ".png"
		fileContent := (config.DataType)("image_" + strconv.Itoa(i+1))
		ds.Make(fileName, fileContent)

		files[fileName] = fileContent
		fileNames = append(fileNames, fileName)
		fileContents = append(fileContents, fileContent)
	}

	return ds, files, fileNames, fileContents
}
