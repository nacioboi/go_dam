package dam

// huffman__node represents a node in the Huffman tree
type huffman__node struct {
	symbol byte
	freq   int
	left   *huffman__node
	right  *huffman__node
}

// huffman_min_heap represents a min-huffman_min_heap of Nodes
type huffman_min_heap struct {
	nodes []*huffman__node
	size  int
}

// huffman_min_heap__heapify maintains the min-heap property
func (h *huffman_min_heap) huffman_min_heap__heapify(i int) {
	smallest := i
	left := 2*i + 1
	right := 2*i + 2

	if left < h.size && h.nodes[left].freq < h.nodes[smallest].freq {
		smallest = left
	}
	if right < h.size && h.nodes[right].freq < h.nodes[smallest].freq {
		smallest = right
	}
	if smallest != i {
		h.nodes[i], h.nodes[smallest] = h.nodes[smallest], h.nodes[i]
		h.huffman_min_heap__heapify(smallest)
	}
}

// huffman_min_heap__insert adds a new node to the heap
func (h *huffman_min_heap) huffman_min_heap__insert(node *huffman__node) {
	h.nodes = append(h.nodes, node)
	h.size++
	i := h.size - 1

	for i != 0 {
		parent := (i - 1) / 2
		if h.nodes[parent].freq > h.nodes[i].freq {
			h.nodes[i], h.nodes[parent] = h.nodes[parent], h.nodes[i]
			i = parent
		} else {
			break
		}
	}
}

// huffman_min_heap__extract_min removes and returns the node with the minimum frequency
func (h *huffman_min_heap) huffman_min_heap__extract_min() *huffman__node {
	if h.size == 0 {
		return nil
	}
	min := h.nodes[0]
	h.size--
	if h.size > 0 {
		h.nodes[0] = h.nodes[h.size]
		h.huffman_min_heap__heapify(0)
	}
	h.nodes = h.nodes[:h.size]
	return min
}

// huffman_min_heap__build constructs a min-heap from a slice of nodes
func huffman_min_heap__build(nodes []*huffman__node) *huffman_min_heap {
	h := &huffman_min_heap{
		nodes: nodes,
		size:  len(nodes),
	}
	for i := h.size/2 - 1; i >= 0; i-- {
		h.huffman_min_heap__heapify(i)
	}
	return h
}

// huffman__symbol_frequencies holds frequencies of symbols (byte values 0-255)
type huffman__symbol_frequencies struct {
	frequencies [256]int
}

// huffman__count_frequencies counts the frequency of each symbol in the data
func huffman__count_frequencies(data []byte) huffman__symbol_frequencies {
	var sf huffman__symbol_frequencies
	for _, b := range data {
		sf.frequencies[b]++
	}
	return sf
}

// huffman__build_tree builds the Huffman tree and returns the root node
func huffman__build_tree(sf huffman__symbol_frequencies) *huffman__node {
	var nodes []*huffman__node
	// Initialize nodes with leaf nodes for each symbol
	for i := 0; i < 256; i++ {
		if sf.frequencies[i] > 0 {
			node := &huffman__node{
				symbol: byte(i),
				freq:   sf.frequencies[i],
				left:   nil,
				right:  nil,
			}
			nodes = append(nodes, node)
		}
	}
	// Build the min-heap
	heap := huffman_min_heap__build(nodes)
	// Build the Huffman tree
	for heap.size > 1 {
		// Extract two nodes with the lowest frequency
		node1 := heap.huffman_min_heap__extract_min()
		node2 := heap.huffman_min_heap__extract_min()
		// Create a new internal node with these two nodes as children
		merged := &huffman__node{
			freq:  node1.freq + node2.freq,
			left:  node1,
			right: node2,
		}
		// Insert the new node into the heap
		heap.huffman_min_heap__insert(merged)
	}
	// The remaining node is the root of the Huffman tree
	return heap.huffman_min_heap__extract_min()
}

// huffman__symbol_codes holds Huffman codes for symbols (byte values 0-255)
type huffman__symbol_codes struct {
	codes  [256]uint32 // Store the code bits
	length [256]uint8  // Store the length of each code
}

// huffman__generate_codes generates Huffman codes for each symbol
func huffman__generate_codes(node *huffman__node, code uint32, length uint8, sc *huffman__symbol_codes) {
	if node == nil {
		return
	}
	// If it's a leaf node
	if node.left == nil && node.right == nil {
		sc.codes[node.symbol] = code
		sc.length[node.symbol] = length
		return
	}
	// Traverse left with '0' added to the code
	huffman__generate_codes(node.left, code<<1, length+1, sc)
	// Traverse right with '1' added to the code
	huffman__generate_codes(node.right, (code<<1)|1, length+1, sc)
}

// huffman__encode encodes the data using the Huffman codes
func huffman__encode(data []byte, sc huffman__symbol_codes) []byte {
	var encoded []byte
	var currentByte byte
	var bitsFilled uint8

	for _, b := range data {
		code := sc.codes[b]
		codeLen := sc.length[b]
		for codeLen > 0 {
			// Determine the number of bits we can write in the current byte
			bitsToWrite := 8 - bitsFilled
			if bitsToWrite > codeLen {
				bitsToWrite = codeLen
			}
			// Calculate the bits to write
			bits := (code >> (codeLen - bitsToWrite)) & ((1 << bitsToWrite) - 1)
			// Shift bits into position
			currentByte <<= bitsToWrite
			currentByte |= byte(bits)
			bitsFilled += bitsToWrite
			codeLen -= bitsToWrite
			// If the current byte is filled, append it to the encoded data
			if bitsFilled == 8 {
				encoded = append(encoded, currentByte)
				currentByte = 0
				bitsFilled = 0
			}
		}
	}
	// Handle the last byte if it's not fully filled
	if bitsFilled > 0 {
		currentByte <<= (8 - bitsFilled)
		encoded = append(encoded, currentByte)
	}
	return encoded
}

// huffman__decode decodes the encoded bytes using the Huffman tree
func huffman__decode(encoded []byte, root *huffman__node, totalBits int) []byte {
	var decoded []byte
	current := root
	bitIndex := 0

	for i := 0; i < totalBits; i++ {
		// Get the current byte and bit
		byteIndex := bitIndex / 8
		bitOffset := 7 - uint(bitIndex%8)
		bit := (encoded[byteIndex] >> bitOffset) & 1
		bitIndex++
		// Traverse the Huffman tree
		if bit == 0 {
			current = current.left
		} else {
			current = current.right
		}
		// If it's a leaf node
		if current.left == nil && current.right == nil {
			decoded = append(decoded, current.symbol)
			current = root
		}
	}
	return decoded
}
