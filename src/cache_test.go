package cache

import (
	"testing"
	"fmt"
)

func TestHello(t *testing.T) {
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