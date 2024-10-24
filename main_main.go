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
	"unsafe"
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

func (_ Huffman__Min_Heap) New(size uint16) *Huffman__Min_Heap {
	return &Huffman__Min_Heap{
		indices: make([]uint16, 0, size),
		size:    0,
	}
}

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
	h := Huffman__Min_Heap{}.New(uint16(len(frequencies.X)))

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

func ensure_padded(slice []rune, by int, and_avoid rune) []rune {
	p := slice
	for len(p)%by != 0 {

		if l__index_of__padded == nil {
			l__index_of__padded = make([]rune, 0)
		}
		target_size := (len(p)) + (by - len(p)%by)
		for len(l__index_of__padded) < target_size {
			l__index_of__padded = append(l__index_of__padded, 0)
		}

		copy(l__index_of__padded, p)
		p = l__index_of__padded

		c := rune(0)
		found_unique := false
		for !found_unique {
			if c == and_avoid {
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
	return p
}

func index_of(slice []rune, char rune) (uint8, bool) {
	if len(slice) == 0 {
		return 0, false
	}

	x, ok := sse_find_idx_uint32_4aat(char, uint32(len(slice)), &slice[0])
	if ok {
		if slice[x] != char {
			log.Fatalf("Expected %d, got %d\n", char, slice[x])
		}
		return x, true
	}

	return 0, false
}

func Huffman__Build_Frequency_Table(text string) Huffman__Frequency_Table {
	frequencies := Huffman__Frequency_Table{X: make([]rune, 0), Y: make([]uint32, 0)}

	var has_changed bool = true
	var tmp_arr []rune
	var cache_k [255]rune
	var cache_v [255]uint8
	for _, char := range text {
		if has_changed {
			tmp_arr = ensure_padded(frequencies.X, 8, char)
		}

		if c := cache_k[char]; c != 0 {
			if c != char {
				panic("Cache key mismatch")
			}
			frequencies.Y[cache_v[char]]++
			continue
		}
		if i, ok := index_of(tmp_arr, char); ok {
			cache_k[char] = char
			cache_v[char] = i
			frequencies.Y[i]++
			continue
		} else {
			frequencies.X = append(frequencies.X, char)
			frequencies.Y = append(frequencies.Y, 1)
			has_changed = true
			for k := range cache_k {
				cache_k[k] = 0
			}
		}
	}

	return frequencies
}

func RV_Encode(text string, codes Huffman__Code_Table) []byte {
	var encoded []byte
	var currentByte uint8
	var bitIndex uint8 = 0

	for _, char := range text {
		for i := 0; i < len(codes.X); i++ {
			if codes.X[i] == char {
				code := codes.Y[i]
				// Convert each bit from the code into packed bytes
				for _, bit := range code {
					if bit == '1' {
						currentByte |= (1 << (7 - bitIndex)) // Set the bit
					}
					bitIndex++
					if bitIndex == 8 {
						encoded = append(encoded, currentByte)
						currentByte = 0
						bitIndex = 0
					}
				}
				break
			}
		}
	}

	return encoded
}

func RV_Decode(encoded []byte, codes Huffman__Code_Table) string {
	var decoded strings.Builder
	var buffer string // Accumulates bits as a string

	for _, byteVal := range encoded {
		for bitIndex := 0; bitIndex < 8; bitIndex++ {
			// Extract bits from the byte and append to buffer
			if byteVal&(1<<(7-bitIndex)) != 0 {
				buffer += "1"
			} else {
				buffer += "0"
			}

			// Check if buffer matches any Huffman code
			for i := 0; i < len(codes.X); i++ {
				code := codes.Y[i]
				if buffer == code {
					decoded.WriteRune(codes.X[i])
					buffer = "" // Reset buffer after match
					break
				}
			}
		}
	}

	return decoded.String()
}

type Huffman__Compile_Time_Codes_Registrar struct{}

type Huffman__Compile_Time_Codes struct {
	X   unsafe.Pointer
	Y   unsafe.Pointer
	Len uint32
}

var (
	l__global_hctc__x []string
	l__global_hctc__y []Huffman__Compile_Time_Codes
)

func (_ Huffman__Compile_Time_Codes_Registrar) register(name string, x []rune, y []string) {
	inst := Huffman__Compile_Time_Codes{
		X:   unsafe.Pointer(&x[0]),
		Y:   unsafe.Pointer(&y[0]),
		Len: uint32(len(x)),
	}

	if l__global_hctc__x == nil {
		l__global_hctc__x = make([]string, 0)
		l__global_hctc__y = make([]Huffman__Compile_Time_Codes, 0)
	}

	l__global_hctc__x = append(l__global_hctc__x, name)
	l__global_hctc__y = append(l__global_hctc__y, inst)
}

func (_ Huffman__Compile_Time_Codes_Registrar) Get(name string) Huffman__Compile_Time_Codes {
	for i := 0; i < len(l__global_hctc__x); i++ {
		if l__global_hctc__x[i] == name {
			return l__global_hctc__y[i]
		}
	}
	panic("Code table not found")
}

func CV_Encode(text string, codes Huffman__Compile_Time_Codes) []byte {
	var encoded []byte
	var currentByte uint8
	var bitIndex uint8 = 0

	for _, char := range text {
		for i := 0; i < int(codes.Len); i++ {
			if *(*rune)(unsafe.Add(codes.X, i<<2)) == char {
				code := *(*string)(unsafe.Add(codes.Y, i<<4))
				// Convert each bit from the code into packed bytes
				for _, bit := range code {
					if bit == '1' {
						currentByte |= (1 << (7 - bitIndex)) // Set the bit
					}
					bitIndex++
					if bitIndex == 8 {
						encoded = append(encoded, currentByte)
						currentByte = 0
						bitIndex = 0
					}
				}
				break
			}
		}
	}

	// Append any remaining bits as a final byte
	if bitIndex > 0 {
		encoded = append(encoded, currentByte)
	}

	return encoded
}

func CV_Decode(encoded []byte, codes Huffman__Compile_Time_Codes) string {
	var decoded strings.Builder
	var buffer string // Accumulates bits as a string

	for _, byteVal := range encoded {
		for bitIndex := 0; bitIndex < 8; bitIndex++ {
			// Extract bits from the byte and append to buffer
			if byteVal&(1<<(7-bitIndex)) != 0 {
				buffer += "1"
			} else {
				buffer += "0"
			}

			// Check if buffer matches any Huffman code
			for i := 0; i < int(codes.Len); i++ {
				code := *(*string)(unsafe.Add(codes.Y, i<<4))
				if buffer == code {
					decoded.WriteRune(*(*rune)(unsafe.Add(codes.X, i<<2)))
					buffer = "" // Reset buffer after match
					break
				}
			}
		}
	}

	return decoded.String()
}

func init() {
	Huffman__Compile_Time_Codes_Registrar{}.register("TFB-ASCIIV", _TFB_ASCIIV_CV_X, _TFB_ASCIIV_CV_Y)
}

func main() {
	var start_t time.Time

	var rv__build_frequency_table_t time.Duration
	var rv__build_tree_t time.Duration
	var rv__build_code_table_t time.Duration

	const N = 1024 //* 256

	// Open file for cpu profiling
	f, err := os.Create("cpu.prof")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	pprof.StartCPUProfile(f)

	// Separate...
	fmt.Println()

	// Benchmark the RV setup functions...
	start_t = time.Now()
	for i := 0; i < N/4; i++ {
		Huffman__Build_Frequency_Table(c__THE_ENTIRE_TFB__ASCIIV)
	}
	rv__build_frequency_table_t = time.Since(start_t)

	_rv__frequencies := Huffman__Build_Frequency_Table(c__THE_ENTIRE_TFB__ASCIIV)
	start_t = time.Now()
	for i := 0; i < N/4; i++ {
		Huffman__Build_Tree(_rv__frequencies)
	}
	rv__build_tree_t = time.Since(start_t)

	_rv__tree := Huffman__Build_Tree(_rv__frequencies)
	start_t = time.Now()
	for i := 0; i < N/4; i++ {
		Huffman__Build_Code_Table(_rv__tree, "")
	}
	rv__build_code_table_t = time.Since(start_t)

	pprof.StopCPUProfile()

	fmt.Printf("Build Frequency Table ns/op:  %v\n", rv__build_frequency_table_t.Nanoseconds()/N/4)
	fmt.Printf("Build Tree ns/op:             %v\n", rv__build_tree_t.Nanoseconds()/N/4)
	fmt.Printf("Build Code Table ns/op:       %v\n", rv__build_code_table_t.Nanoseconds()/N/4)

	// Construct Huffman parameters based on all printable ASCII characters...
	rv__frequencies := Huffman__Build_Frequency_Table(c__THE_ENTIRE_TFB__ASCIIV)
	rv__root := Huffman__Build_Tree(rv__frequencies)
	rv__codes := Huffman__Build_Code_Table(rv__root, "")

	// Separate...
	fmt.Println()

	// Print codes just for fun...
	fmt.Printf("var _TFB_ASCIIV_CV_X []rune = []rune{")
	for i, char := range rv__codes.X {
		if i%8 == 0 {
			fmt.Println()
		}
		fmt.Printf("0x%x, ", char)
	}
	fmt.Println("}")
	fmt.Printf("var _TFB_ASCIIV_CV_Y []string = []string{")
	for i, code := range rv__codes.Y {
		if i%8 == 0 {
			fmt.Println()
		}
		fmt.Printf("\"%s\", ", code)
	}
	fmt.Println("\n}")

	cv__codes := Huffman__Compile_Time_Codes_Registrar{}.Get("TFB-ASCIIV")

	// Separate...
	fmt.Println()

	// Benchmark the two versions of encode, runtime (RV) and compile-time (CV) ...
	start_t = time.Now()
	for j := 0; j < N; j++ {
		for _, char := range rv__codes.X {
			index_of(rv__codes.X, char)
		}
	}
	fmt.Printf("RV ns/op:  %v\n", time.Since(start_t).Nanoseconds()/N)

	start_t = time.Now()
	for j := 0; j < N; j++ {
		for i := 0; i < int(cv__codes.Len); i++ {
			_ = *(*rune)(unsafe.Add(cv__codes.X, i<<2))
			_ = *(*string)(unsafe.Add(cv__codes.Y, i<<4))
		}
	}
	fmt.Printf("CV ns/op:  %v\n", time.Since(start_t).Nanoseconds()/N)

	// Separate...
	fmt.Println()

	// Show results...
	rv__encoded := RV_Encode("Hello, world!", rv__codes)
	fmt.Printf("RV Encoded:                %x\n", rv__encoded)
	fmt.Printf("RV Encoded Length:         %d\n", len(rv__encoded))
	fmt.Printf("RV Decoded (ignore []s):  [%s]\n", RV_Decode(rv__encoded, rv__codes)) // Not implemented yet...

	cv__encoded := CV_Encode("Hello, world!", cv__codes)
	fmt.Printf("CV Encoded:                %x\n", cv__encoded)
	fmt.Printf("CV Encoded Length:         %d\n", len(cv__encoded))
	fmt.Printf("CV Decoded (ignore []s):  [%s]\n", CV_Decode(cv__encoded, cv__codes))
}
