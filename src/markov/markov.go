package markov

import (
	"sync"
	"log"
	"math"
	"../heap"
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
	m.mu.Lock()
	defer m.mu.Unlock()
	if n < 0 {
		log.Fatalf("AJ hasn't implemented fetching %d files yet :/", n)
	} 
	return m.longPaths(filename, n)
	// if n != 1 {
	// 	log.Fatalf("AJ hasn't implemented fetching %d files yet :/", n)
	// } 
	// if node, ok := m.nodes[filename]; ok {
	// 	next, _ := node.GetMaxNeighbor() 
	// 	return append(make([]string, 0), next)
	// } else {
	// 	return make([]string, 0)
	// }
}

// Find highest probabilities from source
// assumes source is in m.nodes
func (m *MarkovChain) longPaths(source string, n int) []string {
	// set up min weights
	prob_log := make(map[string]float64)
	prob_log[source] = 0.0

	// store removed nodes so we can fetch the closest values and check if something has been removed
	var removed heap.MinHeapFloat
	removed.Init()

	// store current guesses
	var heap heap.MinHeapFloat
	heap.Init()
	heap.Insert(source, 1.0)

	for heap.Size > 0 && removed.Size < n {
		name := heap.ExtractMin()
		node := m.nodes[name]
		estimate := prob_log[name]
		if name != source {
			removed.Insert(name, estimate)
		}
		for _, neighbor := range node.adjacencies {
			if _, ok := prob_log[neighbor.name]; !ok {
				// not seen before, set probability estimate and insert into heap
				prob_log[neighbor.name] = math.Inf(1)
				heap.Insert(neighbor.name, math.Inf(1))
			}
			if !removed.Contains(neighbor.name) {
				// then try to relax weight estimate
				weight := -math.Log((float64(neighbor.count) / float64(node.size)))
				if (weight + estimate) < prob_log[neighbor.name] {
					// then relax this edge
					prob_log[neighbor.name] = (weight + estimate)
					heap.ChangeKey(neighbor.name, (weight + estimate))
				}
			}
		}
	}
	return removed.GetKeyList()
}