package markov

import (
	"sync"
	"log"
)

// Individual edge in the markov graph
// count represents frequency
type Edge struct {
	count			int
	name			string
}

// sparse representation of adjacencies. double space for efficient lookups + iteration
type Node struct {
	name			string
	size			int
	adjacencies	    []Edge			// fast iterator 
	neighbors		map[string]int	// filename -> index in adjacencies. fast lookup
	mu				sync.Mutex		// for when the chain should be concurrent
}

// creates empty node for the given name
func MakeNode(name string) *Node {
	node := &Node{name: name, size: 0, adjacencies:make([]Edge, 0), neighbors: make(map[string]int)}
	return node
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
}

func (n *Node) Copy() *Node {
	n.mu.Lock()
	defer n.mu.Unlock()

	// copy edges slice
	edges := make([]Edge, len(n.adjacencies))
	copy(edges, n.adjacencies)

	// copy index map
	neighbors := make(map[string]int)
	for key, value := range  n.neighbors {
		neighbors[key] = value
	}

	return &Node{name: n.name, size: n.size, adjacencies: edges, neighbors: neighbors}
}

func (e *Edge) Copy() Edge {
	return Edge{name:e.name, count: e.count}
}

// add e1 and e2 -> return e1 + e2
func EdgeAdd(e1 *Edge, e2 *Edge) Edge {
	if e1.name != e2.name {
		log.Fatalf("Attempt to add two Edges with different names %v and %v", e1, e2)
	}
	edge := Edge{name: e1.name, count: e1.count + e2.count}
	return edge
}

// add n1 and n2 -> return n1 + n2
func NodeAdd(n1 *Node, n2 *Node) *Node {
	n1.mu.Lock()
	n2.mu.Lock()
	defer n1.mu.Unlock()
	defer n2.mu.Unlock()
	if n1.name != n2.name {
		log.Fatalf("Attempt to add two Nodes with different names %v and %v", n1, n2)
	}
	if (n1.size == 0) {
		return n2.Copy()
	} else if (n2.size == 0) {
		return n1.Copy()
	}

	node := &Node{}
	node.name = n1.name
	node.size = 0
	node.adjacencies = make([]Edge, 0)
	node.neighbors = make(map[string]int)

	index := 0

	total := 0

	for key1, index1 := range n1.neighbors {
		var edge Edge
		edge1 := n1.adjacencies[index1]
		index2, ok := n2.neighbors[key1]
		if ok {
			// add these two edges
			edge2 := n2.adjacencies[index2]
			edge = EdgeAdd(&edge1, &edge2)
		} else {
			edge = edge1.Copy()
		}

		total += edge.count
		node.neighbors[edge.name] = index
		node.adjacencies = append(node.adjacencies, edge)
		index++
	}

	for key2, index2 := range n2.neighbors {
		_, skip := node.neighbors[key2]
		if !skip {
			// this was not found above
			edge2 := n2.adjacencies[index2]
			edge := edge2.Copy()
			node.neighbors[edge.name] = index
			node.adjacencies = append(node.adjacencies, edge)
			index++
			total += edge.count
		}
	}

	node.size = total
	return node
}

// returns accessCount, totalCount
func (n *Node) GetTransProb(filename string) (int, int) {
	n.mu.Lock()
	defer n.mu.Unlock()
	if n.size == 0 {
		return 0, n.size
	} else {
		index, ok := n.neighbors[filename]
		if !ok {
			return 0, n.size
		} else {
			return n.adjacencies[index].count, n.size
		}
	}
}