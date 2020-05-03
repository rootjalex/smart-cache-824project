package task

import "log"

type Workload struct {
	ItemNames        []string
	ItemGroupIndices [][]int // slice of sequence of indices representing item names to access
	curr             int     // index of the current group of item indices
}

func NewMLWorkload(itemNames []string, batchSize int, numIterations int) Workload {
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
				batch = append(batch, k)
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
		log.Println(allBatches)
	}

	// shuffle the batches
	// rand.Seed(time.Now().UnixNano())
	// rand.Shuffle(len(allBatches), func(i, j int) { allBatches[i], allBatches[j] = allBatches[j], allBatches[i] })

	return Workload{
		ItemNames:        itemNames,
		ItemGroupIndices: allBatches,
	}
}

// ------------------------------------------------------

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
