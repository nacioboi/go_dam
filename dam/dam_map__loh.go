/*/
 ** This software is covered by the MIT License.
 ** See: `./LICENSE`.
/*/

package dam

const _LOH_DAM__NUM_ITEMS_PER_BUCKET = 128

// Super-Fast Direct-Access Map.
type DAM_LOH[KT I_Large_Positive_Integer, VT any] struct {
	keys   *compressed_array
	values *compressed_array
}

// Creates a new DAM (Direct-Access Map) that tries to balance speed and memory usage, slight preference for memory savings.
//
// Since this is a DAM, we need to know the expected number of inputs in advance.
// This leaves us with the following options:
//
// - `New_Fast_DAM`: Super-fast DAM, sacrifices memory usage for speed.
//
// - `New_Standard_DAM`: Slightly slower DAM, gives up some speed for memory usage.
//
// - `New_MOH_DAM`: (Medium-OverHead DAM), gives up even more speed for memory usage.
//
// - `New_LOH_DAM`: (Low-OverHead DAM), sacrifices speed for memory savings.
//
func New_LOH_DAM[KT I_Large_Positive_Integer, VT any](
	expected_num_inputs KT,
) *DAM_LOH[KT, VT] {
	expected_num_inputs = max(128, next_power_of_two(expected_num_inputs))
	num_buckets := max(2, expected_num_inputs/_LOH_DAM__NUM_ITEMS_PER_BUCKET)

	if num_buckets%2 != 0 {
		panic("numBuckets should be a multiple of 2.")
	}

	// Allocate buckets...
	//num_buckets_runtime := uint64(num_buckets)

	// Instantiate...
	inst := DAM_LOH[KT, VT]{
		// keys:   make([][_LOH_DAM__NUM_ITEMS_PER_BUCKET]uint64, num_buckets_runtime),
		// values: make([][_LOH_DAM__NUM_ITEMS_PER_BUCKET]VT, num_buckets_runtime),
	}

	return &inst
}

func (m *DAM_LOH[KT, VT]) Enquire_Number_Of_Buckets() KT {
	//return KT(len(m.keys))
	panic("Not implemented.")
}

// Set a key-value pair in the map.
// Will panic if something goes wrong.
//
// - WARNING: This function is NOT thread-safe.
//
//go:inline
func (m *DAM_LOH[KT, VT]) Set(key KT, value VT) {
	// if key == 0 {
	// 	panic("Key cannot be 0.")
	// }

	// index := key & KT(len(m.keys)-1)

	// for i := 0; i < _LOH_DAM__NUM_ITEMS_PER_BUCKET; i++ {
	// 	if m.keys[index][i] == uint64(key) {
	// 		m.values[index][i] = value
	// 		return
	// 	}
	// }

	// for i := 0; i < _LOH_DAM__NUM_ITEMS_PER_BUCKET; i++ {
	// 	if m.keys[index][i] == 0 {
	// 		m.keys[index][i] = uint64(key)
	// 		m.values[index][i] = value
	// 		return
	// 	}
	// }

	// panic("no space left in bucket")
}

// Returns the value and a boolean indicating whether the value was found.
//
// - WARNING: This function is NOT thread-safe.
//
// - NOTE: Remember that keys cannot be 0.
//
// - NOTE: This function will not check if the key is 0.
//
//go:inline
func (m *DAM_LOH[KT, VT]) Get(key KT) (VT, bool) {
	// NOTE: Keeping value type here improves performance since we do not modify the value.
	// index := key & KT(len(m.keys)-1)

	// v, ok := avx512_find_idx_64i(uint64(key), &m.keys[index][0])
	// if ok {
	// 	return m.values[index][v], true
	// }

	// // if m.overflows_enabled {
	// // 	// Fetch from overflow...
	// // 	for i := 0; i < len(m.overflow_keys[index]); i++ {
	// // 		if m.overflow_keys[index][i] == uint64(key) {
	// // 			return m.overflow_values[index][i], true
	// // 		}
	// // 	}
	// // }

	var zero VT
	return zero, false
}

// Delete an entry from the map and return a boolean indicating whether the entry was found.
//
// - WARNING: This function is NOT thread-safe.
//
// - NOTE: Remember that keys cannot be 0.
//
// - NOTE: This function will not check if the key is 0.
//
//go:inline
func (m *DAM_LOH[KT, VT]) Delete(key KT) bool {
	// index := key & m.num_buckets_m1
	// buck := &m.buckets[index]

	// loc := -1

	// for i := 0; i < len(buck.entries); i++ {
	// 	if buck.entries[i].key == uint64(key) {
	// 		loc = i
	// 		break
	// 	}
	// }

	// // Rearrange the entire slice...
	// if loc == -1 {
	// 	return false
	// }

	// buck.entries = append(buck.entries[:loc], buck.entries[loc+1:]...)
	// return true
	panic("Not implemented.")
}

// package dam

// const _LOH_DAM__NUM_ITEMS_PER_BUCKET = 64

// // Super-Fast Direct-Access Map.
// type DAM_LOH[KT I_Large_Positive_Integer, VT any] struct {
// 	compressed [][_LOH_DAM__NUM_ITEMS_PER_BUCKET][]byte

// 	huffman_symbol_codes       huffman__symbol_codes
// 	huffman_symbol_frequencies huffman__symbol_frequencies
// 	huffman_head               *huffman__node

// 	update_codes_count uint8

// 	num_buckets_m1 KT

// 	encode_f func(KT, VT) []byte
// 	decode_f func([]byte) (KT, VT)
// }

// // Creates a new DAM (Direct-Access Map) that tries to balance speed and memory usage, strongly prefers memory savings.
// //
// // Since this is a DAM, we need to know the expected number of inputs in advance.
// // This leaves us with the following options:
// //
// // - `New_Fast_DAM`: Super-fast DAM, sacrifices memory usage for speed.
// //
// // - `New_Standard_DAM`: Slightly slower DAM, gives up some speed for memory usage.
// //
// // - `New_MOH_DAM`: (Medium-OverHead DAM), gives up even more speed for memory usage.
// //
// // - `New_LOH_DAM`: (Low-OverHead DAM), sacrifices speed for memory savings.
// func New_LOH_DAM[KT I_Large_Positive_Integer, VT any](
// 	expected_num_inputs uint64,
// 	encode_f func(KT, VT) []byte,
// 	decode_f func([]byte) (KT, VT),
// ) *DAM_LOH[KT, VT] {
// 	expected_num_inputs = max(128, next_power_of_two(expected_num_inputs))
// 	num_buckets := max(2, expected_num_inputs/_LOH_DAM__NUM_ITEMS_PER_BUCKET)

// 	if num_buckets%2 != 0 {
// 		panic("numBuckets should be a multiple of 2.")
// 	}

// 	// Allocate buckets...
// 	num_buckets_runtime := uint64(num_buckets)

// 	// Instantiate...
// 	inst := DAM_LOH[KT, VT]{
// 		compressed:     make([][_LOH_DAM__NUM_ITEMS_PER_BUCKET][]byte, 0, num_buckets_runtime),
// 		num_buckets_m1: KT(num_buckets - 1),
// 		encode_f:       encode_f,
// 		decode_f:       decode_f,
// 	}

// 	return &inst
// }

// func (m *DAM_LOH[KT, VT]) Enquire_Number_Of_Buckets() KT {
// 	return m.num_buckets_m1 + 1
// }

// // Set a key-value pair in the map.
// // Will panic if something goes wrong.
// //
// // - WARNING: This function is NOT thread-safe.
// //
// //go:inline
// func (m *DAM_LOH[KT, VT]) Set(key KT, value VT) {
// 	if key == 0 {
// 		panic("Key cannot be 0.")
// 	}

// 	index := key & m.num_buckets_m1
// 	if KT(len(m.compressed)) <= index {
// 		m.compressed = m.compressed[:index+1]
// 		//m.has_been_compressed = m.has_been_compressed[:index+1]
// 		m.compressed[index] = [_LOH_DAM__NUM_ITEMS_PER_BUCKET][]byte{}
// 		//m.has_been_compressed[index] = false
// 	}

// 	//if m.has_been_compressed[index] {
// 	// for i := 0; i < _LOH_DAM__NUM_ITEMS_PER_BUCKET; i++ {
// 	// 	if m.keys[index][i] == uint64(key) {
// 	// 		m.values[index][i] = m.encode_value_f(value)
// 	// 		return
// 	// 	}
// 	// }

// 	// for i := 0; i < _LOH_DAM__NUM_ITEMS_PER_BUCKET; i++ {
// 	// 	if m.keys[index][i] == 0 {
// 	// 		m.keys[index][i] = uint64(key)
// 	// 		m.values[index][i] = m.encode_value_f(value)
// 	// 		return
// 	// 	}
// 	// }
// 	//}

// 	if m.update_codes_count == 0 {
// 		collated_bytes := make([]byte, 0)
// 		for i := 0; i < len(m.compressed); i++ {
// 			for j := 0; j < len(m.compressed[i]); j++ {
// 				collated_bytes = append(collated_bytes, m.compressed[i][j]...)
// 			}
// 		}
// 		m.huffman_symbol_frequencies = huffman__count_frequencies(collated_bytes)
// 		m.huffman_head = huffman__build_tree(m.huffman_symbol_frequencies)
// 		huffman__generate_codes(m.huffman_head, 0, 0, &m.huffman_symbol_codes)
// 		m.update_codes_count = 255
// 	}
// 	m.update_codes_count--

// 	encoded := huffman__encode(
// 		m.encode_f(key, value),
// 		m.huffman_symbol_codes,
// 	)
// 	for i := 0; i < len(m.compressed[index]); i++ {
// 		x := 0
// 		for j := 0; j < len(m.compressed[index][i]); j++ {
// 			x += int(m.compressed[index][i][j])
// 		}
// 		if x == 0 {
// 			m.compressed[index][i] = encoded
// 			//m.has_been_compressed[index] = true
// 			return
// 		}
// 	}

// 	panic("No space left in the bucket.")
// }

// // Returns the value and a boolean indicating whether the value was found.
// //
// // - WARNING: This function is NOT thread-safe.
// //
// // - NOTE: Remember that keys cannot be 0.
// //
// // - NOTE: This function will not check if the key is 0.
// //
// //go:inline
// func (m *DAM_LOH[KT, VT]) Get(key KT) (VT, bool) {
// 	// NOTE: Keeping value type here improves performance since we do not modify the value.
// 	index := key & m.num_buckets_m1
// 	if KT(len(m.compressed)) <= index {
// 		var zero VT
// 		return zero, false
// 	}

// 	// Fill in here chatgpt...

// 	var zero VT
// 	return zero, false
// }

// // Delete an entry from the map and return a boolean indicating whether the entry was found.
// //
// // - WARNING: This function is NOT thread-safe.
// //
// // - NOTE: Remember that keys cannot be 0.
// //
// // - NOTE: This function will not check if the key is 0.
// //
// //go:inline
// func (m *DAM_LOH[KT, VT]) Delete(key KT) bool {
// 	// index := key & m.num_buckets_m1
// 	// buck := &m.buckets[index]

// 	// loc := -1

// 	// for i := 0; i < len(buck.entries); i++ {
// 	// 	if buck.entries[i].key == uint64(key) {
// 	// 		loc = i
// 	// 		break
// 	// 	}
// 	// }

// 	// // Rearrange the entire slice...
// 	// if loc == -1 {
// 	// 	return false
// 	// }

// 	// buck.entries = append(buck.entries[:loc], buck.entries[loc+1:]...)
// 	// return true
// 	panic("Not implemented.")
// }
