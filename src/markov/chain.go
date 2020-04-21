package markov

import (
	"sync"
)
// Individual edge in the markov graph
// count represents frequency
type Edge struct {
	count 			int
	name  			string
}

// sparse representation of adjacencies. double space for efficient lookups + iteration
type Node struct {
	name 			string
	size 			int
	adjacencies 	[]Edge			// fast iterator 
	neighbors		map[string]int 	// filename -> index in adjacencies. fast lookup
	mu 				sync.Mutex 		// for when the chain should be concurrent
	best 			*Edge
}

// creates empty node for the given name
func MakeNode(name string) *Node {
	node := &Node{name: name, size: 0, adjacencies:make([]Edge, 0), neighbors: make(map[string]int)}
	return node
}

// returns the neighbor with highest probability 
func (n *Node) GetMaxNeighbor() (string, float32) {
	if n.size == 0 {
		return "", 0.0
	}

	maxFreq := 0
	maxName := ""

	// find max neighbor
	for _, node := range n.adjacencies {
		if node.count > maxFreq {
			maxFreq = node.count 
			maxName = node.name
		}
	}
	return maxName, (float32(maxFreq) / float32(n.size))
}