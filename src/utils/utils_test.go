package utils

import (
	"fmt"
	"reflect"
	"testing"
)

func TestArraySetsEqual(t *testing.T) {
	fmt.Printf("TestArraySetsEqualFail ...\n")
	failed := false

	// case 1
	a := []int{1, 1, 1, 1}
	b := []int{1, 0, 1, 1}
	got := IntArraySetsEqual(a, b)
	expected := false
	if got != expected {
		failed = true
		t.Errorf("got %v but expected %v for equality of %v and %v", got, expected, a, b)
	}

	// case 1
	a = []int{1, 1, 1, 0}
	b = []int{1, 0, 1, 1}
	got = IntArraySetsEqual(a, b)
	expected = true
	if got != expected {
		failed = true
		t.Errorf("got %v but expected %v for equality of %v and %v", got, expected, a, b)
	}

	if failed {
		fmt.Printf("\t... FAILED\n")
	} else {
		fmt.Printf("\t... PASSED\n")
	}
}

func TestGetIntCounts(t *testing.T) {
	fmt.Printf("TestGetIntCountsFail ...\n")
	failed := false

	// case 1
	a := []int{1, 1, 1, 1}
	got := GetIntCounts(a)
	expected := map[int]int{}
	expected[1] = 4
	if !reflect.DeepEqual(expected, got) {
		failed = true
		t.Errorf("got %v but expected %v for original array %v", got, expected, a)
	}

	// case 2
	a = []int{1, 2, 1, 4}
	got = GetIntCounts(a)
	expected = map[int]int{}
	expected[1] = 2
	expected[2] = 1
	expected[4] = 1
	if !reflect.DeepEqual(expected, got) {
		failed = true
		t.Errorf("got %v but expected %v for original array %v", got, expected, a)
	}

	if failed {
		fmt.Printf("\t... FAILED\n")
	} else {
		fmt.Printf("\t... PASSED\n")
	}
}

func TestStringSliceExtendMany(t *testing.T) {
	// partition:
	//      size of slice:      0, 1, >1
	//      replication factor: 1, >1
	t.Log("TestStringSliceExtendMany...")

	// slice size 0
	assertDeepEqual(t, []string{}, StringSliceExtendMany([]string{}, 1))

	// slice size 1, replication factor 1
	assertDeepEqual(t, []string{"a"}, StringSliceExtendMany([]string{"a"}, 1))

	// slice size 1, replication factor >1
	assertDeepEqual(t, []string{"a", "a", "a", "a"}, StringSliceExtendMany([]string{"a"}, 4))

	// slice size >1, replication factor 1
	assertDeepEqual(t, []string{"a", "b", "c"}, StringSliceExtendMany([]string{"a", "b", "c"}, 1))

	// slice size >1, replication factor >1
	assertDeepEqual(t, []string{"1", "2", "3", "4", "1", "2", "3", "4", "1", "2", "3", "4"}, StringSliceExtendMany([]string{"1", "2", "3", "4"}, 3))
}

func assertDeepEqual(t *testing.T, a, b interface{}) {
	if !reflect.DeepEqual(a, b) {
		t.Errorf("Expected %+v and %+v to deep equal", a, b)
	}
}
