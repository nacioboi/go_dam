/*/
 ** This software is covered by the MIT License.
 ** See: `./LICENSE`.
/*/

package main

import (
	"fmt"
	"log"
	"os"
	"runtime/pprof"
	"strings"
	"time"
)

//go:noescape
//go:nosplit
func sse_find_idx_uint32_4aat(query rune, size uint32, p *rune) (uint8, bool)

// func caller(a uint32, b uint32) uint32

// func main() {
// 	data := [8]rune{'u', 'a', 'b', 'd', 'e', 'f', 'g', 'h'}

// 	x1, ok1 := sse_find_idx_uint32_4aat('u', 8, &data[0])
// 	fmt.Printf("Index: %v, Ok: %v\n", x1, ok1)

// 	x2, ok2 := sse_find_idx_uint32_4aat('a', 8, &data[0])
// 	fmt.Printf("Index: %v, Ok: %v\n", x2, ok2)

// 	x3, ok3 := sse_find_idx_uint32_4aat('b', 8, &data[0])
// 	fmt.Printf("Index: %v, Ok: %v\n", x3, ok3)

// 	x4, ok4 := sse_find_idx_uint32_4aat('d', 8, &data[0])
// 	fmt.Printf("Index: %v, Ok: %v\n", x4, ok4)

// 	x5, ok5 := sse_find_idx_uint32_4aat('e', 8, &data[0])
// 	fmt.Printf("Index: %v, Ok: %v\n", x5, ok5)

// 	x6, ok6 := sse_find_idx_uint32_4aat('f', 8, &data[0])
// 	fmt.Printf("Index: %v, Ok: %v\n", x6, ok6)

// 	x7, ok7 := sse_find_idx_uint32_4aat('g', 8, &data[0])
// 	fmt.Printf("Index: %v, Ok: %v\n", x7, ok7)

// 	x8, ok8 := sse_find_idx_uint32_4aat('h', 8, &data[0])
// 	fmt.Printf("Index: %v, Ok: %v\n", x8, ok8)

// 	x9, ok9 := sse_find_idx_uint32_4aat('?', 8, &data[0])
// 	fmt.Printf("Index: %v, Ok: %v\n", x9, ok9)

// 	// fmt.Printf("Caller: %v\n", caller(2, 4))

// }

type Huffman__Node struct {
	Char        rune
	Freq        uint32
	Left_Index  uint16
	Right_Index uint16
}

type Huffman__Tree struct {
	Nodes []Huffman__Node
}

type Huffman__Min_Heap struct {
	indices []uint16
	size    uint16
}

// Huffman__New_Huffman_Min_Heap creates a new heap
func Huffman__New_Huffman_Min_Heap(size uint16) *Huffman__Min_Heap {
	return &Huffman__Min_Heap{
		indices: make([]uint16, 0, size),
		size:    0,
	}
}

// Update Push_Top_Index to Push_Index and accept an index
func (h *Huffman__Min_Heap) Push_Index(index uint16, tree Huffman__Tree) {
	h.indices = append(h.indices, index)
	h.size++
	h.heapify_up(h.size-1, tree)
}

func (h *Huffman__Min_Heap) Pop_Index(tree Huffman__Tree) (uint16, bool) {
	if h.size == 0 {
		return 0, false
	}

	root := h.indices[0]
	h.size--
	h.indices[0] = h.indices[h.size]
	h.indices = h.indices[:h.size]

	h.heapify_down(0, tree)
	return root, true
}

func (h *Huffman__Min_Heap) heapify_up(index uint16, tree Huffman__Tree) {
	for index > 0 {
		parent := (index - 1) / 2
		if tree.Nodes[h.indices[index]].Freq >= tree.Nodes[h.indices[parent]].Freq {
			break
		}
		h.indices[index], h.indices[parent] = h.indices[parent], h.indices[index]
		index = parent
	}
}

func (h *Huffman__Min_Heap) heapify_down(index uint16, tree Huffman__Tree) {
	for {
		left := 2*index + 1
		right := 2*index + 2
		smallest := index

		if left < h.size && tree.Nodes[h.indices[left]].Freq < tree.Nodes[h.indices[smallest]].Freq {
			smallest = left
		}
		if right < h.size && tree.Nodes[h.indices[right]].Freq < tree.Nodes[h.indices[smallest]].Freq {
			smallest = right
		}

		if smallest == index {
			break
		}

		h.indices[index], h.indices[smallest] = h.indices[smallest], h.indices[index]
		index = smallest
	}
}

const c_HUFFMAN__NO_CHILD_INDEX uint16 = 0xFFFF

func Huffman__Build_Tree(frequencies Huffman__Frequency_Table) Huffman__Tree {
	tree := Huffman__Tree{
		Nodes: make([]Huffman__Node, 0, len(frequencies.X)*2),
	}
	h := Huffman__New_Huffman_Min_Heap(uint16(len(frequencies.X)))

	// Push all nodes into the heap
	for i := 0; i < len(frequencies.X); i++ {
		node := Huffman__Node{
			Char:        frequencies.X[i],
			Freq:        frequencies.Y[i],
			Left_Index:  c_HUFFMAN__NO_CHILD_INDEX,
			Right_Index: c_HUFFMAN__NO_CHILD_INDEX,
		}
		tree.Nodes = append(tree.Nodes, node)
		index := uint16(len(tree.Nodes) - 1)
		h.Push_Index(index, tree)
	}

	// Build the tree
	for h.size > 1 {
		left, success := h.Pop_Index(tree)
		if !success {
			panic("Failed to pop left node")
		}
		right, success := h.Pop_Index(tree)
		if !success {
			panic("Failed to pop right node")
		}

		parent := Huffman__Node{
			Char:        0,
			Freq:        tree.Nodes[left].Freq + tree.Nodes[right].Freq,
			Left_Index:  left,
			Right_Index: right,
		}
		tree.Nodes = append(tree.Nodes, parent)
		parentIndex := uint16(len(tree.Nodes) - 1)
		h.Push_Index(parentIndex, tree)
	}

	return tree
}

type Huffman__Code_Table struct {
	X []rune
	Y []string
}

func Huffman__Build_Code_Table(tree Huffman__Tree, prefix string) Huffman__Code_Table {
	root_index := uint16(len(tree.Nodes) - 1)
	codes := Huffman__Code_Table{X: make([]rune, 0), Y: make([]string, 0)}
	return huffman__inner__Build_Code_Table(tree, root_index, prefix, codes)
}

func huffman__inner__Build_Code_Table(
	tree Huffman__Tree,
	node_index uint16,
	prefix string,
	codes Huffman__Code_Table,
) Huffman__Code_Table {
	node := &tree.Nodes[node_index]

	if node.Left_Index == c_HUFFMAN__NO_CHILD_INDEX && node.Right_Index == c_HUFFMAN__NO_CHILD_INDEX {
		codes.X = append(codes.X, node.Char)
		codes.Y = append(codes.Y, prefix)
		return codes
	}

	codes = huffman__inner__Build_Code_Table(tree, node.Left_Index, prefix+"0", codes)
	codes = huffman__inner__Build_Code_Table(tree, node.Right_Index, prefix+"1", codes)
	return codes
}

type Huffman__Frequency_Table struct {
	X []rune
	Y []uint32
}

var (
	l__index_of__padded []rune
)

func index_of(slice []rune, char rune) (uint8, bool) {
	if len(slice) == 0 {
		return 0, false
	}

	p := slice
	for len(p)%8 != 0 {

		if l__index_of__padded == nil {
			l__index_of__padded = make([]rune, 0)
		}
		target_size := (len(p)) + (8 - len(p)%8)
		for len(l__index_of__padded) < target_size {
			l__index_of__padded = append(l__index_of__padded, 0)
		}

		copy(l__index_of__padded, p)
		p = l__index_of__padded

		c := rune(0)
		found_unique := false
		for !found_unique {
			if c == char {
				c++
				continue
			}
			for j := 0; j < len(p); j++ {
				if c == p[j] {
					c++
					break
				} else {
					found_unique = true
					break
				}
			}
		}

		for i := len(slice); i < len(p); i++ {
			p[i] = c
		}
	}

	//for i := 0; i < len(slice)/4; i++ {

	x, ok := sse_find_idx_uint32_4aat(char, uint32(len(p)), &p[0])
	if ok {
		if p[x] != char {
			log.Fatalf("Expected %d, got %d\n", char, p[x])
		}
		return x, true
	}

	// var x uint8
	// var ok bool
	// for i := 0; i < len(p); i++ {
	// 	if p[i] == char {
	// 		x = uint8(i)
	// 		ok = true
	// 		break
	// 	}
	// }
	// if ok {
	// 	return x, true
	// }

	//}
	return 0, false
}

func Huffman__Build_Frequency_Table(text string) Huffman__Frequency_Table {
	frequencies := Huffman__Frequency_Table{X: make([]rune, 0), Y: make([]uint32, 0)}

	for _, char := range text {
		if i, ok := index_of(frequencies.X, char); ok {
			frequencies.Y[i]++
		} else {
			frequencies.X = append(frequencies.X, char)
			frequencies.Y = append(frequencies.Y, 1)
		}
	}

	return frequencies
}

func Encode(text string, codes Huffman__Code_Table) string {
	var encoded strings.Builder
	for _, char := range text {
		i, ok := index_of(codes.X, char)
		if !ok {
			panic("Character not found in code table")
		}
		encoded.WriteString(codes.Y[i])
	}
	return encoded.String()
}

// TODO: The solution is to write a string deserializer that takes in a instructions generated by a tool
type Huffman__Compile_Time_Parameters struct {
	frequencies Huffman__Frequency_Table
	root        Huffman__Node
	codes       Huffman__Code_Table
}

func main() {
	var start_t time.Time

	var build_frequency_table_t time.Duration
	// var build_tree_t time.Duration
	// var build_code_table_t time.Duration

	const D = 1024

	// Open file for cpu profiling
	f, err := os.Create("cpu.prof")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	pprof.StartCPUProfile(f)

	start_t = time.Now()
	for i := 0; i < D; i++ {
		Huffman__Build_Frequency_Table(c_THE_ENTIRE_FREE_BIBLE)
	}
	build_frequency_table_t = time.Since(start_t)

	// _frequencies := Huffman__Build_Frequency_Table(c_THE_ENTIRE_FREE_BIBLE)
	// start_t = time.Now()
	// for i := 0; i < D; i++ {
	// 	Huffman__Build_Tree(_frequencies)
	// }
	// build_tree_t = time.Since(start_t)

	// _tree := Huffman__Build_Tree(_frequencies)
	// start_t = time.Now()
	// for i := 0; i < D; i++ {
	// 	Huffman__Build_Code_Table(_tree, "")
	// }
	// build_code_table_t = time.Since(start_t)

	pprof.StopCPUProfile()

	// Passing everything via the stack gives time:
	// Build Frequency Table: 1178393
	// Build Tree: 2018
	//Build Code Table: 3033

	fmt.Printf("Build Frequency Table /op: %v\n", build_frequency_table_t.Microseconds()/D)
	// fmt.Printf("Build Tree: %v\n", build_tree_t.Microseconds())
	// fmt.Printf("Build Code Table: %v\n", build_code_table_t.Microseconds())

	// // Construct Huffman parameters based on all printable ASCII characters...
	// frequencies := Huffman__Build_Frequency_Table(c_THE_ENTIRE_FREE_BIBLE)
	// root := Huffman__Build_Tree(frequencies)
	// codes := Huffman__Build_Code_Table(root, "")

	// // Print Huffman Codes...
	// fmt.Println("Huffman Codes:")
	// for i := 0; i < len(codes.X); i++ {
	// 	fmt.Printf("%c: %s\n", codes.X[i], codes.Y[i])
	// }

	// // Encode the text...
	// encoded := Encode("Hello, world!", codes)
	// fmt.Printf("\nEncoded: %s\n", encoded)
}
