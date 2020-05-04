package cache

import (
	"fmt"
	"reflect"
	"strconv"
	"testing"
	"../datastore"
	"../utils"
	"../config"
)

func TestBasicLRUFail(t *testing.T) {
	fmt.Printf("TestBasicLRUFail ...\n")
	failed := false
	misses := 0

	data := datastore.MakeDataStore()

	// add files to datastore
	for j := 0; j < (config.CACHE_SIZE + 1); j++ {
		filename := "fake_" + strconv.Itoa(j) + ".txt"
        data.Make(filename, config.DataType(filename))
	}

	// this copies data, so can't adjust later
	var cache Cache
	cache.Init(1, config.CACHE_SIZE, config.LRU, data)

	for i := 0; i < 2; i++ {
		for j := 0; j < (config.CACHE_SIZE + 1); j++ {
			filename := "fake_" + strconv.Itoa(j) + ".txt"
			_, err := cache.Fetch(filename)
			if err != nil {
				t.Errorf("Could not open %s from cache", filename)
				failed = true
			}
		}
	}

	hits, misses := cache.Report()
	if hits != 0 || misses != (2*config.CACHE_SIZE+2) {
		t.Errorf("Expected 0 hits and %d miss, got %d hits and %d misses.", (2*config.CACHE_SIZE + 2), hits, misses)
		failed = true
	}

	if failed {
		fmt.Printf("\t... FAILED\n")
	} else {
		fmt.Printf("\t... PASSED\n")
	}
}

func TestBasicLRUSuccess(t *testing.T) {
	fmt.Printf("TestBasicLRUSuccess ...\n")
	failed := false
	misses := 0
	data := datastore.MakeDataStore()

	// add files to datastore
	for j := 0; j < config.CACHE_SIZE; j++ {
		filename := "fake_" + strconv.Itoa(j) + ".txt"
        data.Make(filename, config.DataType(filename))
	}

	var cache Cache
	cache.Init(1, config.CACHE_SIZE, config.LRU, data)
	if config.CACHE_SIZE > 100 {
		fmt.Printf("\tignoring, CACHE_SIZE too big\n")
		return
	}

	for i := 0; i < 2; i++ {
		for j := 0; j < config.CACHE_SIZE; j++ {
			filename := "fake_" + strconv.Itoa(j) + ".txt"
			_, err := cache.Fetch(filename)
			if err != nil {
				t.Errorf("Could not open %s from cache", filename)
				failed = true
			}
		}
	}

	hits, misses := cache.Report()
	if hits != config.CACHE_SIZE || misses != config.CACHE_SIZE {
		t.Errorf("Expected %d hits and %d miss, got %d hits and %d misses.", config.CACHE_SIZE, config.CACHE_SIZE, hits, misses)
		failed = true
	}

	if failed {
		fmt.Printf("\t... FAILED\n")
	} else {
		fmt.Printf("\t... PASSED\n")
	}
}

func TestHashsplitAmongstGroups(t *testing.T) {
	fmt.Printf("TestHashsplitAmongstGroups ...\n")
	failed := false

	// case 0
	n := 4
	numGroups := 2
	groups := splitAmongstGroups(n, numGroups)
	expected := make([]int, n)
	expected[0] = 0
	expected[1] = 0
	expected[2] = 1
	expected[3] = 1

	if !utils.IntArrayEqual(expected, groups) {
		failed = true
		t.Errorf("Got %v but expected %v with n: %v and numGroups: %v", groups, expected, n, groups)
	}

	// case 1
	n = 5
	numGroups = 2
	groups = splitAmongstGroups(n, numGroups)
	expected = make([]int, n)
	expected[0] = 0
	expected[1] = 0
	expected[2] = 1
	expected[3] = 1
	expected[4] = 0

	if !utils.IntArrayEqual(expected, groups) {
		failed = true
		t.Errorf("Got %v but expected %v with n: %v and numGroups: %v", groups, expected, n, groups)
	}

	// case 2
	n = 4
	numGroups = 3
	groups = splitAmongstGroups(n, numGroups)
	expected = make([]int, n)
	expected[0] = 0
	expected[1] = 1
	expected[2] = 2
	expected[3] = 0

	if !utils.IntArrayEqual(expected, groups) {
		failed = true
		t.Errorf("Got %v but expected %v with n: %v and numGroups: %v", groups, expected, n, groups)
	}

	// case 3
	n = 9
	numGroups = 3
	groups = splitAmongstGroups(n, numGroups)
	expected = make([]int, n)
	expected[0] = 0
	expected[1] = 0
	expected[2] = 0
	expected[3] = 1
	expected[4] = 1
	expected[5] = 1
	expected[6] = 2
	expected[7] = 2
	expected[8] = 2

	if !utils.IntArrayEqual(expected, groups) {
		failed = true
		t.Errorf("Got %v but expected %v with n: %v and numGroups: %v", groups, expected, n, groups)
	}

	// case 4
	n = 11
	numGroups = 3
	groups = splitAmongstGroups(n, numGroups)
	expected = make([]int, n)
	expected[0] = 0
	expected[1] = 0
	expected[2] = 0
	expected[3] = 1
	expected[4] = 1
	expected[5] = 1
	expected[6] = 2
	expected[7] = 2
	expected[8] = 2
	expected[9] = 0
	expected[10] = 1

	if !utils.IntArrayEqual(expected, groups) {
		failed = true
		t.Errorf("Got %v but expected %v with n: %v and numGroups: %v", groups, expected, n, groups)
	}

	// case 5
	n = 4
	numGroups = 5
	groups = splitAmongstGroups(n, numGroups)
	expected = make([]int, n)
	expected[0] = 0
	expected[1] = 1
	expected[2] = 2
	expected[3] = 3

	if !utils.IntArrayEqual(expected, groups) {
		failed = true
		t.Errorf("Got %v but expected %v with n: %v and numGroups: %v", groups, expected, n, groups)
	}
	if failed {
		fmt.Printf("\t... FAILED\n")
	} else {
		fmt.Printf("\t... PASSED\n")
	}
}

func TestHashmakeFileGroups(t *testing.T) {
	fmt.Printf("TestHashmakeFileGroups ...\n")
	failed := false

	// case 0
	filenames := []string{"a", "b", "c", "d"}
	numGroups := 2
	groups := makeFileGroups(filenames, len(filenames), numGroups, 1)
	expected := make(map[string]int)
	expected["a"] = 0
	expected["b"] = 0
	expected["c"] = 1
	expected["d"] = 1

	if !reflect.DeepEqual(expected, groups) {
		t.Errorf("Got %v but expected %v with filenames: %v and numGroups: %v", groups, expected, filenames, groups)
		failed = true
	}

	// case 1
	filenames = []string{"a", "b", "c", "d"}
	numGroups = 3
	groups = makeFileGroups(filenames, len(filenames), numGroups, 1)
	expected = make(map[string]int)
	expected["a"] = 0
	expected["b"] = 1
	expected["c"] = 0
	expected["d"] = 2

	// case 1
	filenames = []string{"ab", "b", "c", "d"}
	numGroups = 3
	groups = makeFileGroups(filenames, len(filenames), numGroups, 1)
	expected = make(map[string]int)
	expected["ab"] = 0
	expected["b"] = 1
	expected["c"] = 0
	expected["d"] = 2

	if !reflect.DeepEqual(expected, groups) {
		t.Errorf("Got %v but expected %v with filenames: %v and numGroups: %v", groups, expected, filenames, numGroups)
		failed = true
	}
	if failed {
		fmt.Printf("\t... FAILED\n")
	} else {
		fmt.Printf("\t... PASSED\n")
	}
}
