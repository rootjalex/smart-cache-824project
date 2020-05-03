package utils

import (
	"testing"
	"fmt"
    "reflect"
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