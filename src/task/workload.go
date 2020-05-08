package task

import (
	"math/rand"
	"time"
)

// ------------------------------------------------------ WORKLOAD GENERATOR

type WorkloadGenerator struct {
	wkld *Workload
}

// generates the same workload every time
func NewMLWorkloadGenerator(itemNames []string, batchSize int, numIterations int) WorkloadGenerator {
	ml := newMLWorkload(itemNames, batchSize, numIterations)
	return WorkloadGenerator{
		wkld: &ml,
	}
}

// generates different workload order every time
func NewRandomWorkloadGenerator(itemNames []string, batchSize int) WorkloadGenerator {
	r := newRandomWorkload(itemNames, batchSize)
	return WorkloadGenerator{
		wkld: &r,
	}
}

// generates different workload every time
func NewWebWorkloadGenerator(itemNames []string, numPatterns int, minPatternLength int, maxPatternLength int, replicationFactor int) WorkloadGenerator {
	w := newWebWorkload(itemNames, numPatterns, minPatternLength, maxPatternLength, replicationFactor)
	return WorkloadGenerator{
		wkld: &w,
	}
}

// ------------------------------------------------------ WORKLOAD GENERATOR METHODS

func (wg *WorkloadGenerator) GenerateWorkload() Workload {
	return wg.wkld.FreshCopy()
}

// ------------------------------------------------------ WORKLOAD

type Workload struct {
	ItemNames        []string
	ItemGroupIndices [][]int // slice of sequence of indices representing item names to access
	curr             int     // index of the current group of item indices

	workloadName      string
	batchSize         int
	numIterations     int
	numPatterns       int
	minPatternLength  int
	maxPatternLength  int
	replicationFactor int
}

func newMLWorkload(itemNames []string, batchSize int, numIterations int) Workload {
	var allBatches [][]int

	if batchSize >= len(itemNames) {
		// case: batch size greater than items ==> batch is all the items
		batch := make([]int, len(itemNames))
		for i := 0; i < len(itemNames); i++ {
			batch[i] = i
		}
		// replicate the batch for the specified number of iterations
		allBatches = make([][]int, numIterations)
		for i := 0; i < numIterations; i++ {
			allBatches[i] = make([]int, len(batch))
			copy(allBatches[i], batch)
		}
	} else {
		// case: batch smaller than items ==> multiple batches
		batches := [][]int{}
		i, j := 0, batchSize-1
		for i < len(itemNames) {
			batch := []int{}
			for k := i; k <= j; k++ {
				if k < len(itemNames) {
					batch = append(batch, k)
				}
			}
			i += batchSize
			j += batchSize
			batches = append(batches, batch)
		}
		// replicate the batches for the specified number of iterations
		allBatches = [][]int{}
		for i := 0; i < numIterations; i++ {
			for _, batch := range batches {
				b := make([]int, len(batch))
				copy(b, batch)
				allBatches = append(allBatches, b)
			}
		}
	}

	return Workload{
		ItemNames:        itemNames,
		ItemGroupIndices: allBatches,
		batchSize:        batchSize,
		numIterations:    numIterations,
		workloadName:     "ml",
	}
}

func newRandomWorkload(itemNames []string, batchSize int) Workload {
	batchIndices := [][]int{}

	// construct and add the batches to indices
	if batchSize >= len(itemNames) {
		// case: batch size greater than items ==> batch is all the items
		batch := make([]int, len(itemNames))
		for i := range batch {
			batch[i] = i
		}
		batchIndices = append(batchIndices, batch)
	} else {
		// case: batch smaller than items ==> multiple batches
		i, j := 0, batchSize-1
		for i < len(itemNames) {
			batch := []int{}
			for k := i; k <= j; k++ {
				if k < len(itemNames) {
					batch = append(batch, k)
				}
			}
			i += batchSize
			j += batchSize
			batchIndices = append(batchIndices, batch)
		}
	}
	// randomize order of each batch of indices and the whole bag of batches
	rand.Seed(time.Now().UnixNano())
	for k := range batchIndices {
		rand.Shuffle(len(batchIndices[k]), func(i, j int) { batchIndices[k][i], batchIndices[k][j] = batchIndices[k][j], batchIndices[k][i] })
	}
	rand.Shuffle(len(batchIndices), func(i, j int) { batchIndices[i], batchIndices[j] = batchIndices[j], batchIndices[i] })

	return Workload{
		ItemNames:        itemNames,
		ItemGroupIndices: batchIndices,
		batchSize:        batchSize,
		workloadName:     "random",
	}
}

func newWebWorkload(itemNames []string, numPatterns int, minPatternLength int, maxPatternLength int, replicationFactor int) Workload {
	rand.Seed(time.Now().UnixNano())

	// populate different kinds of patterns
	patterns := [][]int{}
	for i := 0; i < numPatterns; i++ {
		pLength := minPatternLength + rand.Intn(maxPatternLength-minPatternLength)
		p := []int{}
		pStart := len(itemNames) / (i + 1)
		for j := pStart; j < pStart+pLength; j++ {
			p = append(p, j%len(itemNames))
		}
		patterns = append(patterns, p)
	}
	// log.Printf("Patterns %+v", patterns)
	// randomly pick patterns and extend to the big pattern
	bigPattern := []int{}
	for i := 0; i < replicationFactor; i++ {
		pi := rand.Intn(len(patterns))
		bigPattern = append(bigPattern, patterns[pi]...)
	}
	// log.Printf("bigPattern %+v", bigPattern)
	return Workload{
		ItemNames:         itemNames,
		ItemGroupIndices:  [][]int{bigPattern},
		numPatterns:       numPatterns,
		minPatternLength:  minPatternLength,
		maxPatternLength:  maxPatternLength,
		replicationFactor: replicationFactor,
		workloadName:      "web",
	}
}

// ------------------------------------------------------ WORKLOAD METHODS

func (wkld *Workload) FreshCopy() Workload {
	switch wkld.workloadName {
	case "random":
		return newRandomWorkload(wkld.ItemNames, wkld.batchSize)
	case "ml":
		return newMLWorkload(wkld.ItemNames, wkld.batchSize, wkld.numIterations)
	case "web":
		return newWebWorkload(wkld.ItemNames, wkld.numPatterns, wkld.minPatternLength, wkld.maxPatternLength, wkld.replicationFactor)
	default:
		return Workload{
			ItemNames:         wkld.ItemNames,
			ItemGroupIndices:  wkld.ItemGroupIndices,
			workloadName:      wkld.workloadName,
			batchSize:         wkld.batchSize,
			numIterations:     wkld.numIterations,
			numPatterns:       wkld.numPatterns,
			minPatternLength:  wkld.minPatternLength,
			maxPatternLength:  wkld.maxPatternLength,
			replicationFactor: wkld.replicationFactor,
		}
	}
}

func (wkld *Workload) HasNextItemGroup() bool {
	return 0 <= wkld.curr && wkld.curr < len(wkld.ItemGroupIndices)
}

func (wkld *Workload) GetNextItemGroup() []string {
	itemNameGroup := []string{}
	for _, j := range wkld.ItemGroupIndices[wkld.curr] {
		itemNameGroup = append(itemNameGroup, wkld.ItemNames[j])
	}
	wkld.curr++
	return itemNameGroup
}
