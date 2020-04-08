package cache

import (
	"testing"
)

func TestHello(t *testing.T) {
	want := "input/hello.txt"
	var cache BadCache
	_, err := cache.Fetch(want)
	if err != nil {
		t.Errorf("Could not open %s", want)
	}
	hits, misses := cache.Report()
	if hits != 0 || misses != 1 {
		t.Errorf("Expected 0 hits and 1 miss, got %d hits and %d misses.", hits, misses)
	}
}