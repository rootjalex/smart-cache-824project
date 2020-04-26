package markov

import (
	"sync"
	"log"
)

type MarkovChain struct {
	nodes 			map[string]*Node  // filename -> Node (with adjacencies)
	lastAccess		string
	mu 				sync.Mutex
}


// creates empty node for the given name
func MakeMarkovChain() *MarkovChain {
	markov := &MarkovChain{lastAccess: "", nodes: make(map[string]*Node)}
	// set default MC for first call to markov.Access()
	markov.nodes[""] = MakeNode("")
	return markov
}

func (m *MarkovChain) Access(filename string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.nodes[m.lastAccess].MakeAccess(filename)

	// check if file has own chain
	if _, ok := m.nodes[filename]; !ok {
		m.nodes[filename] = MakeNode(filename)
	}
	m.lastAccess = filename
}

// predict the next n files after filename is accessed
func (m *MarkovChain) Predict(filename string, n int) []string {
	if n != 1 {
		log.Fatalf("AJ hasn't implemented fetching %d files yet :/", n)
	} 
	if node, ok := m.nodes[filename]; ok {
		next, _ := node.GetMaxNeighbor() 
		return append(make([]string, 0), next)
	} else {
		return make([]string, 0)
	}
}