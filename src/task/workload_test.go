package task

// import (
// 	"fmt"
// 	"reflect"
// 	"strconv"
// 	"testing"

// 	"../config"
// 	"../utils"
// )

// func TestBasicWorkload(t *testing.T) {
// 	fmt.Println("TestBasicWorkload ...")

// 	itemNames := []string{"a", "b", "c", "d"}
// 	itemGroupIndices := [][]int{
// 		[]int{0, 1},
// 		[]int{2, 3},
// 		[]int{0, 1, 2, 3},
// 	}
// 	wg := WorkloadGenerator{wkld: &Workload{ItemNames: itemNames, ItemGroupIndices: itemGroupIndices}}
// 	w := wg.GenerateWorkload()

// 	// check the {0, 1} ==> {a, b} group
// 	assertWorkloadHasNextItemGroup(t, &w, []string{"a", "b"})

// 	// check the {2, 3} ==> {c, d} group
// 	assertWorkloadHasNextItemGroup(t, &w, []string{"c", "d"})

// 	// check the {0, 1, 2, 3} ==> {a, b, c, d} group
// 	assertWorkloadHasNextItemGroup(t, &w, []string{"a", "b", "c", "d"})

// 	// check that there are no more item groups
// 	assertNoMoreItemGroups(t, &w)
// }

// func TestBasicMLWorkloadSmallBatchSmallIters(t *testing.T) {
// 	fmt.Println("TestBasicMLWorkloadSmallBatchSmallIters ...")

// 	numFiles := 1000

// 	itemNames := make([]string, numFiles)
// 	for i := 0; i < numFiles; i++ {
// 		itemNames[i] = "imagenet-" + strconv.Itoa(i+1)
// 	}
// 	wg := NewMLWorkloadGenerator(itemNames, 1, 1)
// 	w := wg.GenerateWorkload()

// 	// check that each file makes it in its own group
// 	for i := 0; i < numFiles; i++ {
// 		assertWorkloadHasNextItemGroup(t, &w, []string{itemNames[i]})
// 	}

// 	// check that there are no more item groups
// 	assertNoMoreItemGroups(t, &w)
// }

// func TestBasicMLWorkloadLargeBatchLargeIters(t *testing.T) {
// 	fmt.Println("TestBasicMLWorkloadLargeBatchLargeIters ...")

// 	numFiles := 51
// 	batchSize := 10
// 	numIterations := 10

// 	itemNames := make([]string, numFiles)
// 	for i := 0; i < numFiles; i++ {
// 		itemNames[i] = "imagenet-" + strconv.Itoa(i+1)
// 	}
// 	wg := NewMLWorkloadGenerator(itemNames, batchSize, numIterations)
// 	w := wg.GenerateWorkload()

// 	for i := 0; i < numIterations; i++ {
// 		// 1-10
// 		assertWorkloadHasNextItemGroup(t, &w, []string{"imagenet-1", "imagenet-2", "imagenet-3", "imagenet-4", "imagenet-5", "imagenet-6", "imagenet-7", "imagenet-8", "imagenet-9", "imagenet-10"})
// 		// 11-20
// 		assertWorkloadHasNextItemGroup(t, &w, []string{"imagenet-11", "imagenet-12", "imagenet-13", "imagenet-14", "imagenet-15", "imagenet-16", "imagenet-17", "imagenet-18", "imagenet-19", "imagenet-20"})
// 		// 21-30
// 		assertWorkloadHasNextItemGroup(t, &w, []string{"imagenet-21", "imagenet-22", "imagenet-23", "imagenet-24", "imagenet-25", "imagenet-26", "imagenet-27", "imagenet-28", "imagenet-29", "imagenet-30"})
// 		// 31-40
// 		assertWorkloadHasNextItemGroup(t, &w, []string{"imagenet-31", "imagenet-32", "imagenet-33", "imagenet-34", "imagenet-35", "imagenet-36", "imagenet-37", "imagenet-38", "imagenet-39", "imagenet-40"})
// 		// 41-50
// 		assertWorkloadHasNextItemGroup(t, &w, []string{"imagenet-41", "imagenet-42", "imagenet-43", "imagenet-44", "imagenet-45", "imagenet-46", "imagenet-47", "imagenet-48", "imagenet-49", "imagenet-50"})
// 		// 51
// 		assertWorkloadHasNextItemGroup(t, &w, []string{"imagenet-51"})
// 	}

// 	// check that there are no more item groups
// 	assertNoMoreItemGroups(t, &w)
// }

// func TestRandomWorkloadAndGenerator(t *testing.T) {
// 	fmt.Println("TestRandomWorkloadAndGenerator ...")

// 	// param and file setup
// 	numFiles := 10
// 	batchSize := 4
// 	itemNames := make([]string, numFiles)
// 	for i := 0; i < numFiles; i++ {
// 		itemNames[i] = "" + strconv.Itoa(i+1)
// 	}
// 	// create two workload generators
// 	wg1 := NewRandomWorkloadGenerator(itemNames, batchSize)
// 	wg2 := NewRandomWorkloadGenerator(itemNames, batchSize)

// 	// check multiple workloads generated from the same generator
// 	workloadsPerGenerator := 5
// 	itemGroups1 := [][]string{}
// 	itemGroups2 := [][]string{}
// 	for i := 0; i < workloadsPerGenerator; i++ {
// 		w1, items1 := wg1.GenerateWorkload(), []string{}
// 		w2, items2 := wg2.GenerateWorkload(), []string{}
// 		for w1.HasNextItemGroup() && w2.HasNextItemGroup() {
// 			ig1, ig2 := w1.GetNextItemGroup(), w2.GetNextItemGroup()
// 			items1 = append(items1, ig1...)
// 			items2 = append(items2, ig2...)
// 		}
// 		itemGroups1 = append(itemGroups1, items1)
// 		itemGroups2 = append(itemGroups2, items2)

// 		// items must be the same ones we created
// 		if !utils.StringArraySetsEqual(itemNames, items1) {
// 			t.Errorf("Expected items from workload %+v to equal the files %+v", w1, itemNames)
// 		}
// 		if !utils.StringArraySetsEqual(itemNames, items2) {
// 			t.Errorf("Expected items from workload %+v to equal the files %+v", w2, itemNames)
// 		}
// 		// check that there are no more item groups
// 		assertNoMoreItemGroups(t, &w1)
// 		assertNoMoreItemGroups(t, &w2)
// 	}
// 	// check that item groups not equal
// 	for i := 1; i < workloadsPerGenerator; i++ {
// 		if reflect.DeepEqual(itemGroups1[i-1], itemGroups1[i]) || reflect.DeepEqual(itemGroups2[i-1], itemGroups2[i]) {
// 			t.Errorf("Expected different item groups")
// 		}
// 	}
// }

// func TestWebWorkloadAndGenerator(t *testing.T) {
// 	fmt.Println("TestWebWorkloadAndGenerator ...")

// 	numFiles := 100
// 	itemNames := make([]string, numFiles)
// 	for i := 0; i < numFiles; i++ {
// 		itemNames[i] = "x" + strconv.Itoa(i+1)
// 	}

// 	wg := NewWebWorkloadGenerator(itemNames, config.NUM_PATTERNS_SMALL, config.MIN_PATTERN_LENGTH, config.MAX_PATTERN_LENGTH, config.PATTERN_REPLICATION_SMALL)
// 	w := wg.GenerateWorkload()
// 	// log.Printf("wkld %+v", w)
// 	for w.HasNextItemGroup() {
// 		// log.Println(w.GetNextItemGroup())
// 		continue
// 	}

// 	// just checkt h
// }

// func assertWorkloadHasNextItemGroup(t *testing.T, w *Workload, itemGroup []string) {
// 	// check has group
// 	if !w.HasNextItemGroup() {
// 		t.Errorf("Workload says there is not a next item when there should be. \n\tWORKLOAD %+v", w)
// 	}
// 	// check obtained group
// 	obtainedGroup := w.GetNextItemGroup()
// 	if !reflect.DeepEqual(obtainedGroup, itemGroup) {
// 		t.Errorf("Expected workload to produce item group %+v but got %+v. \n\tWORKLOAD %+v", itemGroup, obtainedGroup, w)
// 	}
// }

// func assertNoMoreItemGroups(t *testing.T, w *Workload) {
// 	if w.HasNextItemGroup() {
// 		t.Errorf("Expected workload to not have any more item groups. \n\tWORKLOAD %+v", w)
// 	}
// }
