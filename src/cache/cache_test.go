package cache

import (
	"testing"
	"fmt"
	"strconv"
	"../datastore"
)

func TestBasicLRUFail(t *testing.T) {
	fmt.Printf("TestBasicLRUFail ...\n")
	failed := false
	misses := 0

	data := datastore.MakeDataStore()

	// add files to datastore
	for j := 0; j < (CACHE_SIZE + 1); j++ {
		filename := "fake_" + strconv.Itoa(j) + ".txt"
		data.Make(filename)
	}

	// this copies data, so can't adjust later
	var cache Cache
	cache.Init(CACHE_SIZE, LRU, data)

	for i := 0; i < 2; i++ {
		for j := 0; j < (CACHE_SIZE + 1); j++ {
			filename := "fake_" + strconv.Itoa(j) + ".txt"
			_, err := cache.Fetch(filename)
			if err != nil {
				t.Errorf("Could not open %s from cache", filename)
				failed = true
			}
		}
	}

	hits, misses := cache.Report()
	if hits != 0 || misses != (2 * CACHE_SIZE + 2) {
		t.Errorf("Expected 0 hits and %d miss, got %d hits and %d misses.", (2 * CACHE_SIZE + 2), hits, misses)
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
	for j := 0; j < CACHE_SIZE; j++ {
		filename := "fake_" + strconv.Itoa(j) + ".txt"
		data.Make(filename)
	}

	var cache Cache
	cache.Init(CACHE_SIZE, LRU, data)
	if CACHE_SIZE > 100 {
		fmt.Printf("\tignoring, CACHE_SIZE too big\n")
		return
	}

	for i := 0; i < 2; i++ {
		for j := 0; j < CACHE_SIZE; j++ {
			filename := "fake_" + strconv.Itoa(j) + ".txt"
			_, err := cache.Fetch(filename)
			if err != nil {
				t.Errorf("Could not open %s from cache", filename)
				failed = true
			}
		}
	}

	hits, misses := cache.Report()
	if hits != CACHE_SIZE || misses != CACHE_SIZE {
		t.Errorf("Expected %d hits and %d miss, got %d hits and %d misses.", CACHE_SIZE, CACHE_SIZE, hits, misses)
		failed = true
	}

	if failed {
		fmt.Printf("\t... FAILED\n")
	} else {
		fmt.Printf("\t... PASSED\n")
	}
}
