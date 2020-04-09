package cache

import (
	"testing"
)

func TestBasicHeap(t *testing.T) {
	var heap MinHeap
	heap.Init()
	
	heap.Insert("first", -1)
	heap.Insert("fourth", 5)
	heap.Insert("second", 2)
	heap.Insert("last", 12394)
	heap.Insert("third", 4)

	label := heap.ExtractMin()
	if label != "first" {
		t.Errorf("Expected 'first', got %s", label)
	}
	label = heap.ExtractMin()
	if label != "second" {
		t.Errorf("Expected 'second', got %s", label)
	}
	label = heap.ExtractMin()
	if label != "third" {
		t.Errorf("Expected 'third', got %s", label)
	}
	heap.Insert("fifth", 100)
	label = heap.ExtractMin()
	if label != "fourth" {
		t.Errorf("Expected 'fourth', got %s", label)
	}
	label = heap.ExtractMin()
	if label != "fifth" {
		t.Errorf("Expected 'fifth', got %s", label)
	}
	label = heap.ExtractMin()
	if label != "last" {
		t.Errorf("Expected 'last', got %s", label)
	}
	if heap.n != 0 {
		t.Errorf("Expected 0 items left, got %d items", heap.n)
	}
}