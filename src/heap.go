package cache 

type HeapItem struct {
	label 	string 
	key 	int64
}

type MinHeap struct {
	items 	[]HeapItem
	labels  map[string]int
	n 		int
}

func (h *MinHeap) Init() {
	h.items = make([]HeapItem, 0)
	h.labels = make(map[string]int)
}

func (h *MinHeap) MinHeapifyUp(c int) {
	if c == 0 {
		return
	}
	p := (c - 1) / 2 
	if h.items[p].key > h.items[c].key {
		// swap terms
		h.Swap(p, c)
		h.labels[h.items[p].label] = p
		h.labels[h.items[c].label] = c
		h.MinHeapifyUp(p)
	}
}

func (h *MinHeap) MinHeapifyDown(p int) {
	if p >= h.n {
		return
	}

	// check children
	l := 2 * p + 1
	r := 2 * p + 2
	if l >= h.n {
		l = p 
	}
	if r >= h.n {
		r = p 
	}

	// set child pointer
	var c int
	if h.items[r].key > h.items[l].key {
		c = l 
	} else {
		c = r
	}

	if h.items[p].key > h.items[c].key {
		// swap terms
		h.Swap(p, c)
		h.labels[h.items[p].label] = p
		h.labels[h.items[c].label] = c
		h.MinHeapifyDown(c)
	}
}

func (h *MinHeap) Insert(label string, key int64) {
	var i HeapItem 
	i.label = label 
	i.key = key 
	h.items = append(h.items, i)
	h.labels[label] = h.n
	h.n++
	DPrintf("Inserting item %d, %s", key, label)
	DPrintf("Resultant: %v", h.items)
	h.MinHeapifyUp(h.labels[label])
}

func (h *MinHeap) ExtractMin() string {
	// swap first and last terms
	h.Swap(0, h.n - 1)
	h.labels[h.items[0].label] = 0
	delete(h.labels, h.items[h.n - 1].label)
	DPrintf("Removing item %d, %s", h.items[h.n - 1].key, h.items[h.n - 1].label)
	label := h.items[h.n - 1].label
	h.items = h.items[:(h.n - 1)]
	h.n--
	h.MinHeapifyDown(0)
	return label
}

func (h *MinHeap) ChangeKey(label string, key int64) {
	index, ok := h.labels[label]
	if ok {
		if key < h.items[index].key {
			h.items[index].key = key 
			h.MinHeapifyUp(index)
		} else {
			h.items[index].key = key 
			h.MinHeapifyDown(index)
		}
	}
}

func (h *MinHeap) Swap(i int, j int) {
	DPrintf("Swap %d, %s", h.items[i].key, h.items[i].label)
	DPrintf("With %d, %s", h.items[j].key, h.items[j].label)
	temp := h.items[i]
	h.items[i] = h.items[j]
	h.items[j] = temp 
	DPrintf("Finished: %d, %s", h.items[i].key, h.items[i].label)
	DPrintf("With %d, %s\n", h.items[j].key, h.items[j].label)
}