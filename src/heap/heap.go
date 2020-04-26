package heap

type HeapItem struct {
	label	string
	key	    int64
}

type MinHeap struct {
	items		[]HeapItem
	labels  	map[string]int
	Size		int
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
	if p >= h.Size {
		return
	}

	// check children
	l := 2 * p + 1
	r := 2 * p + 2
	if l >= h.Size {
		l = p
	}
	if r >= h.Size {
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
	h.labels[label] = h.Size
	h.Size++
	h.MinHeapifyUp(h.labels[label])
}

func (h *MinHeap) ExtractMin() string {
	// swap first and last terms
	h.Swap(0, h.Size - 1)
	h.labels[h.items[0].label] = 0
	delete(h.labels, h.items[h.Size - 1].label)
	label := h.items[h.Size - 1].label
	h.items = h.items[:(h.Size - 1)]
	h.Size--
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
	temp := h.items[i]
	h.items[i] = h.items[j]
	h.items[j] = temp
}

func (h *MinHeap) Contains(key string) bool {
	_, ok := h.labels[key]
	return ok
}

func (h *MinHeap) GetKeyList() []string {
	li := make([]string, len(h.items))
	for i, v := range h.items {
		li[i] = v.label
	}
	return li
}

func (h *MinHeap) GetKey(name string) int64 {
	index, ok := h.labels[name]
	if ok {
		return h.items[index].key
	} else {
		return 0
	}
}