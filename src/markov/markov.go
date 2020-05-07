package markov

import (
	"sync"
	"log"
	"math"
	"../heap"
	// "../utils"
)

type MarkovChain struct {
	nodes			map[string]*Node  // filename -> Node (with adjacencies)
	lastAccess		string
	mu				sync.Mutex
}


// creates empty node for the given name
func MakeMarkovChain() *MarkovChain {
	markov := &MarkovChain{lastAccess: "", nodes: make(map[string]*Node)}
	// set default MC for first call to markov.Access()
	markov.nodes[""] = MakeNode("")
	return markov
}

func (m *MarkovChain) Access(filename string) {
	// utils.DPrintf("Entering Access")
	// defer utils.DPrintf("Leaving Access")
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
	src_node, ok := m.nodes[source]
	if !ok {
		log.Fatalf("THIS SHOULDNEVER HAPPEN %v -> %v", source, m.nodes)
	}
	
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

// returns accessCount, totalCount
func (m *MarkovChain) GetTransProb(filename string, next string) Transition {
	m.mu.Lock()
	defer m.mu.Unlock()
	if len(m.nodes) == 0 {
		return Transition{value:0, total:0}
	} else {
		node, ok := m.nodes[filename]
		if ok {
			return node.GetTransProb(next)
		} else {
			return Transition{value:0, total:0}
		}
	}
}

// add m1 and m2 -> return m1 + m2
func ChainAdd(m1 *MarkovChain, m2 *MarkovChain) *MarkovChain {
	m1.mu.Lock()
	m2.mu.Lock()
	defer m1.mu.Unlock()
	defer m2.mu.Unlock()
	mSum := MakeMarkovChain()


	for key, value := range m1.nodes {
		mSum.nodes[key] = value.Copy()
	}

	for key, value := range m2.nodes {
		node, ok := mSum.nodes[key]
		if ok {
			mSum.nodes[key] = NodeAdd(node, value)
		} else {
			mSum.nodes[key] = value.Copy()
		}
	}
	return mSum
}

// subtract m2 from m1 -> return m1 - m2
// IMPORTANT: assumes all keys in m2 are in m1
// if m2 contains a key not in m1, does NOTHING for that node
// ASSUMES m2 is a prefix of m1
func ChainSub(m1 *MarkovChain, m2 *MarkovChain) *MarkovChain {
	m1.mu.Lock()
	m2.mu.Lock()
	defer m1.mu.Unlock()
	defer m2.mu.Unlock()
	mDiff := MakeMarkovChain()

	for key, value := range m1.nodes {
		mDiff.nodes[key] = value.Copy()
	}

	for key, value := range m2.nodes {
		node, ok := mDiff.nodes[key]
		if ok {
			mDiff.nodes[key] = NodeSub(node, value)
		}
	}

	return mDiff
}

func (m *MarkovChain) Copy() *MarkovChain {
	m.mu.Lock()
	defer m.mu.Unlock()

	// empty chain
	cpy := MakeMarkovChain()

	// deep copy nodes and transitions
	for filename, node := range m.nodes {
		cpy.nodes[filename] = node.Copy()
	}

	// copy last access
	cpy.lastAccess = m.lastAccess

	return cpy
}
