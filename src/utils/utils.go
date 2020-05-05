package utils

import (
	"fmt"
	"log"
	"reflect"

	"../config"
)

const Debug int = 1

func DPrintf(format string, args ...interface{}) {
	if Debug > 0 {
		log.Printf(format, args...)
	}
}

func CountValues(a []int, i int) int {
	count := 0
	for _, val := range a {
		if val == i {
			count++
		}
	}
	return count
}

// JOIN
func JoinInts(a []int, b []int) []int {
	var res []int
	for _, val := range a {
		res = append(res, val)
	}
	for _, val := range b {
		res = append(res, val)
	}
	return res
}

// Equal tells whether a and b contain the same elements.
// A nil argument is equivalent to an empty slice.
func IntArrayEqual(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

func GetIntCounts(a []int) map[int]int {
	aCounts := map[int]int{}
	for _, v := range a {
		if _, ok := aCounts[v]; !ok {
			aCounts[v] = 0
		}
		aCounts[v]++
	}
	return aCounts
}

func GetStringCounts(a []string) map[string]int {
	aCounts := map[string]int{}
	for _, v := range a {
		if _, ok := aCounts[v]; !ok {
			aCounts[v] = 0
		}
		aCounts[v]++
	}
	return aCounts
}

func GetDataTypeCounts(a []config.DataType) map[config.DataType]int {
	aCounts := map[config.DataType]int{}
	for _, v := range a {
		if _, ok := aCounts[v]; !ok {
			aCounts[v] = 0
		}
		aCounts[v]++
	}
	return aCounts
}

// Checks if two Arrays contain same elements
func IntArraySetsEqual(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	aCounts := GetIntCounts(a)
	bCounts := GetIntCounts(b)
	return reflect.DeepEqual(aCounts, bCounts)
}

// Checks if two Arrays contain same elements
func StringArraySetsEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	aCounts := GetStringCounts(a)
	bCounts := GetStringCounts(b)
	return reflect.DeepEqual(aCounts, bCounts)
}

// Checks if two Arrays contain same elements
func DataTypeArraySetsEqual(a, b []config.DataType) bool {
	if len(a) != len(b) {
		return false
	}
	aCounts := GetDataTypeCounts(a)
	bCounts := GetDataTypeCounts(b)
	return reflect.DeepEqual(aCounts, bCounts)
}

// Returns a new slice with the given slice `s` replicated in it `n` times
func StringSliceExtendMany(s []string, n int) []string {
	// replication factor must be >= 1
	if n < 1 {
		panic(fmt.Sprintf("Can only replicate with replication factor geq to 1 but got %v", n))
	}
	out := make([]string, n*len(s))
	for i := 0; i < n; i++ {
		for j, v := range s {
			offset := i*len(s) + j
			out[offset] = v
		}
	}
	return out
}

func Max(a int, b int) int {
	if a > b {
		return a
	}
	return b
}

func Min(a int, b int) int {
	if a < b {
		return a
	}
	return b
}
