package task

import (
	"math/rand"
	"strconv"
	"time"
	"sync"
	"log"
	"../cache"
	"../utils"
	"../datastore"
	"../config"
)

type WorkType int 

const (
	WEB 	WorkType = 0
	ML 		WorkType = 1
	RANDOM 	WorkType = 2
)

// Control params such as time to wait and holds all fetches to be made
type FileBatches struct {
	batches 		[][]string  // list of batches of files to fetch
	minFileSleep	int
	maxFileSleep	int
	minBatchSleep	int
	maxBatchSleep	int
	workType 		WorkType
}

func (fb *FileBatches) Copy() FileBatches {
	fc := FileBatches{
		minFileSleep: fb.minFileSleep, 
		minBatchSleep: fb.minBatchSleep, 
		maxFileSleep: fb.maxFileSleep,
		maxBatchSleep: fb.maxBatchSleep,
	}

	batch_copy := make([][]string, len(fb.batches))
	for i, batch := range fb.batches {
		batch_copy[i] = make([]string, len(batch))
		for j, filename := range batch {
			batch_copy[i][j] = filename
		}
	}
	fc.batches = batch_copy
	return fc
}

type CacheParams struct {
	NCaches 		int
	RFactor 		int
	CacheType 		config.CacheType
	CacheSize 		int
	Datastore 		*datastore.DataStore
	Sync_time 		int
}

type Runner struct {
	mu        sync.Mutex
	Clients   []*Client
	Datastore *datastore.DataStore
	Caches    map[int]*cache.Cache
	Hash      *cache.Hash
}

func CreateRunner(fb []FileBatches, nClients int, params CacheParams) *Runner {
	runner := &Runner{}

	if len(fb) != nClients {
		log.Fatalf("[ CreateRunner ] %v FileBatches but %v Clients", len(fb), nClients)
	}

	// create clients and make a list of their UIDs
	runner.Clients = make([]*Client, nClients)
	clientIDs := make([]int, nClients)
	for i := range runner.Clients {
		runner.Clients[i] = Init(i)
		clientIDs[i] = runner.Clients[i].GetID()
	}

	runner.Caches, runner.Hash = StartTask(clientIDs, params)

	// set up hashes, caches, and file batches
	for i := range runner.Clients {
		runner.Clients[i].BootstrapClient(runner.Caches, runner.Hash, fb[i])
	}

	return runner
} 

// runs all tasks for all clients
func (r *Runner) Run() (map[int][]config.DataType, time.Duration) {
	clientToFiles := make(map[int][]config.DataType)

	startTime := time.Now()
	// run all clients in parallel, wait until all are done
	// aggregate client fetch results
	var wg sync.WaitGroup
	for i, c := range r.Clients {
		wg.Add(1)
		// utils.WaitRandomMillis(1000, 2000)
		go func(client *Client, nc int) {
			utils.DPrintf("Entering lambda Client %v...", nc)
			utils.DPrintf("Leaving lambda Client %v...", nc)
			fetched := client.Run()
			r.mu.Lock()
			clientToFiles[client.GetID()] = fetched
			r.mu.Unlock()
			wg.Done()
		}(c, i)
	}
	wg.Wait()
	elapsedTime := time.Since(startTime)
	return clientToFiles, elapsedTime
}

func (r *Runner) ReportData() (int64,  int64, int64) {
	var hits 	int64
	var misses	int64
	var calls 	int64
	hits = 0
	misses = 0
	calls = 0
	for _, cache := range r.Caches {
		h, m, c := cache.Report()
		hits += h
		misses += m
		calls += c
	}
	return hits, misses, calls
}

func CreateDatastore(fb []FileBatches) *datastore.DataStore {
	datastore := datastore.MakeDataStore()

	for _, files := range fb {
		for _, batch := range files.batches {
			for _, file := range batch {
				datastore.Make(file, config.DataType("good"))
			}
		}

		// exit early if all the same, otherwise make sure all files are made
		if files.workType == ML {
			return datastore
		}
	}
	return datastore
}

func MakeBaseBatch(minFileSleep int, maxFileSleep int, minBatchSleep int, maxBatchSleep int) *FileBatches {
	fb := &FileBatches{
		minFileSleep: minFileSleep,
		maxFileSleep: maxFileSleep,
		minBatchSleep: minBatchSleep,
		maxBatchSleep: maxBatchSleep, 
	}
	return fb
}

func MakeMLFiles(nBatches int, batchLength int, nIterations int, fb *FileBatches) {
	fb.workType = ML
	fb.batches = make([][]string, nBatches * nIterations)

	for c := 0; c < nIterations; c++ {
		i := 1
		for j := 0; j < nBatches; j++ {
			// create a batch
			fb.batches[(c * nBatches) + j] = make([]string, batchLength)
			// fill the batch
			for k := 0; k < batchLength; k++ {
				fb.batches[(c * nBatches) + j][k] = "MNIST_" + strconv.Itoa(i) + ".png"
				i++
				// could do index math here but I am lazy rn
			}
		}
	}
	// log.Printf("%v", fb.batches)
}

func MakeRandomFiles(nBatches int, minBatchLength int, maxBatchLength int, maxFileCount int, fb *FileBatches) {
	fb.workType = RANDOM
	fb.batches = make([][]string, nBatches)

	rand.Seed(time.Now().UnixNano())

	for j := 0; j < nBatches; j++ {
		// create a batch
		batchLength := minBatchLength + rand.Intn(maxBatchLength - minBatchLength)
		fb.batches[j] = make([]string, batchLength)
		// fill the batch
		for k := 0; k < batchLength; k++ {
			// do not use 0 files
			i := rand.Intn(maxFileCount) + 1
			fb.batches[j][k] = "file_" + strconv.Itoa(i) + ".txt"
			i++
		}
	}
}

func MakeWebFiles(nBatches int, nPatterns int, minPatternLength int, maxPatternLength int, 
	batchSize int, maxFileCount int, fb *FileBatches) {
	fb.workType = WEB
	fb.batches = make([][]string, nBatches)

	rand.Seed(time.Now().UnixNano())

	// make patterns
	patterns := make([][]string, nPatterns)
	i := 1
	for j := 0; j < nPatterns; j++ {
		// create a random pattern
		pLength := minPatternLength + rand.Intn(maxPatternLength - minPatternLength)
		patterns[j] = make([]string, pLength)
		for k := range patterns[j] {
			patterns[j][k] = "random_" + strconv.Itoa(i) + ".html"
			i++
		}
	}

	for j := 0; j < nBatches; j++ {
		// create a batch
		fb.batches[j] = make([]string, 0)
		// fill the batch
		for len(fb.batches[j]) < batchSize {
			// append another randomly chosen pattern
			i := rand.Intn(nPatterns)
			fb.batches[j] = utils.JoinStrings(fb.batches[j], patterns[i])
		}
	}
}
