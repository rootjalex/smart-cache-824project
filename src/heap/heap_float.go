package heap

type HeapItemFloat struct {
	label	string
	key	    float64
}

type MinHeapFloat struct {
	items		[]HeapItemFloat
	labels  	map[string]int
	Size		int
}

func (h *MinHeapFloat) Init() {
	h.items = make([]HeapItemFloat, 0)
	h.labels = make(map[string]int)
}

func (h *MinHeapFloat) MinHeapifyUp(c int) {
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

func (h *MinHeapFloat) MinHeapifyDown(p int) {
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

func (h *MinHeapFloat) Insert(label string, key float64) {
	var i HeapItemFloat
	i.label = label
	i.key = key
	h.items = append(h.items, i)
	h.labels[label] = h.Size
	h.Size++
	h.MinHeapifyUp(h.labels[label])
}

func (h *MinHeapFloat) ExtractMin() string {
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

func (h *MinHeapFloat) ChangeKey(label string, key float64) {
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

func (h *MinHeapFloat) Swap(i int, j int) {
	temp := h.items[i]
	h.items[i] = h.items[j]
	h.items[j] = temp
}

func (h *MinHeapFloat) Contains(key string) bool {
	_, ok := h.labels[key]
	return ok
}

func (h *MinHeapFloat) GetKeyList() []string {
	li := make([]string, len(h.items))
	for i, v := range h.items {
		li[i] = v.label
	}
	return li
}

func (h *MinHeapFloat) GetKey(name string) float64 {
	index, ok := h.labels[name]
	if ok {
		return h.items[index].key
	} else {
		return 0.0
	}
}