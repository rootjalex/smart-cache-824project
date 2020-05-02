package markov

import (
	"testing"
	"fmt"
)


func TestNodeAddSimple(t *testing.T) {
	fmt.Printf("TestNodeAddSimple ...\n")
	failed := false

	node1 := MakeNode("test")
	node1.MakeAccess("A.txt")
	node1.MakeAccess("B.txt")
	node1.MakeAccess("C.txt")
	node1.MakeAccess("A.txt")
	node1.MakeAccess("B.txt")
	node1.MakeAccess("C.txt")
	node1.MakeAccess("E.txt")
	// node1: A: 2, B: 2, C:2, E:1

	node2 := MakeNode("test")
	node2.MakeAccess("D.txt")
	node2.MakeAccess("B.txt")
	node2.MakeAccess("A.txt")
	node2.MakeAccess("C.txt")
	node2.MakeAccess("B.txt")
	// node2: A: 1, B: 2, C:1, D:1

	node := NodeAdd(node1, node2)
	// should be: A: 3, B: 4, C:3, D:1, E:1

	// expected value
	size := 12
	expected := make(map[string]int)
	expected["A.txt"] = 3
	expected["B.txt"] = 4
	expected["C.txt"] = 3
	expected["D.txt"] = 1
	expected["E.txt"] = 1

	for key, value := range expected {
		received, total := node.GetTransProb(key)
		if received != value {
			t.Errorf("Expected VALUE: %v, got: %v, for filename %v", value, received, key)
		}
		if total != size {
			t.Errorf("Expected TOTAL: %v, got: %v, for filename %v", size, total, key)
		}
	}
	if failed {
		fmt.Printf("\t... FAILED\n")
	} else {
		fmt.Printf("\t... PASSED\n")
	}
}

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
