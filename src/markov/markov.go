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
}

// Find highest probabilities from source
// assumes source is in m.nodes
// CAN predict source as likely to be fetched again
func (m *MarkovChain) longPaths(source string, n int) []string {
	// set up min weights
	prob_log := make(map[string]float64)

	// store removed nodes so we can fetch the closest values and check if something has been removed
	var removed heap.MinHeapFloat
	removed.Init()

	// store current guesses
	var heap heap.MinHeapFloat
	heap.Init()

	// relax all edges from source
	src_node := m.nodes[source] 
	
	for _, neighbor := range src_node.adjacencies {
		weight := -math.Log((float64(neighbor.count) / float64(src_node.size)))
		// not seen before, set probability estimate and insert into heap
		prob_log[neighbor.name] = weight
		heap.Insert(neighbor.name, weight)
	}

	// now run Dijkstra's
	for heap.Size > 0 && removed.Size < n {
		name := heap.ExtractMin()
		node := m.nodes[name]
		estimate := prob_log[name]
		removed.Insert(name, estimate)
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
	closest := removed.GetKeyList()
	return closest
}