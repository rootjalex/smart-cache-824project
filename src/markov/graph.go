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
	node := &Node{name: name, size: 0, adjacencies:make([]Edge, 0), neighbors: make(map[string]int), best: nil}
	return node
}

// returns the neighbor with highest probability 
func (n *Node) GetMaxNeighbor() (string, float64) {
	n.mu.Lock()
	defer n.mu.Unlock()
	if n.size == 0 {
		return "", 0.0
	} else {
		return n.best.name, (float64(n.best.count) / float64(n.size))
	}
}

func (n *Node) MakeAccess(filename string) {
	n.mu.Lock()
	defer n.mu.Unlock()

	n.size++

	neighbor, ok := n.neighbors[filename]

	if ok {
		n.adjacencies[neighbor].count++
	} else {
		var e Edge 
		e.count = 1
		e.name = filename

		// set index in map and append to end of list
		n.neighbors[filename] = len(n.adjacencies) 
		n.adjacencies = append(n.adjacencies, e)

		neighbor = n.neighbors[filename]
	}

	if n.best == nil || n.adjacencies[neighbor].count > n.best.count {
		n.best = &n.adjacencies[neighbor]
	}
}