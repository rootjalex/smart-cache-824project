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
