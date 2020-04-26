package markov

import (
	"testing"
	"fmt"
)

func TestSimpleMarkovPredict(t *testing.T) {
	fmt.Printf("TestSimpleMarkovPredict ...\n")
	failed := false

	m := MakeMarkovChain()

	m.Access("A.txt")
	m.Access("B.txt")
	m.Access("A.txt")

	got_b := m.Predict("A.txt", 1)
	got_a := m.Predict("B.txt", 1)

	expect_b := "B.txt"
	expect_a := "A.txt"

	if len(got_b) != 1 || got_b[0] != expect_b {
		failed = true
		t.Errorf("Expected prediction: [ %v ], got prediction: %v", expect_b, got_b)
	} 

	if len(got_a) != 1 || got_a[0] != expect_a {
		failed = true
		t.Errorf("Expected prediction: [ %v ], got prediction: %v", expect_a, got_a)
	} 


	if failed {
		fmt.Printf("\t... FAILED\n")
	} else {
		fmt.Printf("\t... PASSED\n")
	}
}


func TestMultiplePathsMarkovPredict(t *testing.T) {
	fmt.Printf("TestMultiplePathsMarkovPredict ...\n")
	failed := false

	m := MakeMarkovChain()

	m.Access("A.txt")
	m.Access("B.txt")
	m.Access("C.txt")

	m.Access("A.txt")
	m.Access("B.txt")
	m.Access("C.txt")

	m.Access("A.txt")
	m.Access("B.txt")
	m.Access("D.txt")

	m.Access("A.txt")
	m.Access("B.txt")
	m.Access("E.txt")

	m.Access("B.txt")
	m.Access("C.txt")

	m.Access("A.txt")
	m.Access("B.txt")
	m.Access("D.txt")

	got_a := m.Predict("A.txt", 3)
	got_b := m.Predict("B.txt", 3)

	// order does not matter, content is important
	expect_a := make(map[string]bool)
	expect_a["B.txt"] = true
	expect_a["C.txt"] = true
	expect_a["A.txt"] = true
	expect_b := make(map[string]bool)
	expect_b["C.txt"] = true
	expect_b["B.txt"] = true
	expect_b["A.txt"] = true

	if len(got_a) == len(expect_a) {
		for _, name := range got_a {
			if _, ok := expect_a[name]; !ok {
				failed = true
				t.Errorf("Expected A prediction: [ %v ], got prediction: %v", expect_a, got_a)
			}
		}
	} else {
		failed = true
		t.Errorf("Expected A prediction: [ %v ], got prediction: %v", expect_a, got_a)
	}

	if len(got_b) == len(expect_b) {
		for _, name := range got_b {
			if _, ok := expect_b[name]; !ok {
				failed = true
				t.Errorf("Expected B prediction: [ %v ], got prediction: %v", expect_b, got_b)
			}
		}
	} else {
		failed = true
		t.Errorf("Expected B prediction: [ %v ], got prediction: %v", expect_b, got_b)
	}

	if failed {
		fmt.Printf("\t... FAILED\n")
	} else {
		fmt.Printf("\t... PASSED\n")
	}
}
