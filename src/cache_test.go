package cache

import (
	"testing"
	"fmt"
	"strconv"
	"os"
)

func TestHelloBadCache(t *testing.T) {
	fmt.Printf("TestHelloBadCache ...\n")
	failed := false
	want := "input/hello.txt"
	var cache BadCache
	_, err := cache.Fetch(want)
	if err != nil {
		t.Errorf("Could not open %s", want)
		failed = true
	}
	hits, misses := cache.Report()
	if hits != 0 || misses != 1 {
		t.Errorf("Expected 0 hits and 1 miss, got %d hits and %d misses.", hits, misses)
		failed = true
	}

	if failed {
		fmt.Printf("\t... FAILED\n")
	} else {
		fmt.Printf("\t... PASSED\n")
	}
}

// assumes CACHE_SIZE <= 100
func TestBasicLRUFail(t *testing.T) {
	fmt.Printf("TestBasicLRUFail ...\n")
	failed := false
	misses := 0
	if CACHE_SIZE > 100 {
		fmt.Printf("\tignoring, CACHE_SIZE too big\n")
		return
	}

	var cache Cache
	cache.Init(CACHE_SIZE)

	for i := 0; i < 2; i++ {
		for j := 0; j < (CACHE_SIZE + 1); j++ {
			filename := "input/" + strconv.Itoa(j) + ".txt"
			f, _ := os.Create(filename) // ignore error if already exists
			f.Close()
			_, err := cache.Fetch(filename)
			if err != nil {
				t.Errorf("Could not open %s from cache", filename)
				failed = true
			}
		}
	}

	// clear up created files 
	for j := 0; j < (CACHE_SIZE + 1); j++ {
		filename := "input/" + strconv.Itoa(j) + ".txt"
		os.Remove(filename)
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

// assumes CACHE_SIZE <= 100
func TestBasicLRUSuccess(t *testing.T) {
	fmt.Printf("TestBasicLRUSuccess ...\n")
	failed := false
	misses := 0
	if CACHE_SIZE > 100 {
		fmt.Printf("\tignoring, CACHE_SIZE too big\n")
		return
	}

	var cache Cache
	cache.Init(CACHE_SIZE)

	for i := 0; i < 2; i++ {
		for j := 0; j < CACHE_SIZE; j++ {
			filename := "input/" + strconv.Itoa(j) + ".txt"
			f, _ := os.Create(filename) // ignore error if already exists
			f.Close()
			_, err := cache.Fetch(filename)
			if err != nil {
				t.Errorf("Could not open %s from cache", filename)
				failed = true
			}
		}
	}

	// clear up created files 
	for j := 0; j < CACHE_SIZE; j++ {
		filename := "input/" + strconv.Itoa(j) + ".txt"
		os.Remove(filename)
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
