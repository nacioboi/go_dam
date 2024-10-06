/*/
 ** This software is covered by the MIT License.
 ** See: `./LICENSE`.
/*/

package dam

type t_bucket_moh[KT I_Positive_Integer, VT any] struct {
	keys   [64]uint64
	values [64]VT
}

// Super-Fast Direct-Access Map.
type DAM_MOH[KT I_Positive_Integer, VT any] struct {
	buckets        []t_bucket_moh[KT, VT]
	overflows      map[KT]VT
	num_buckets_m1 KT

	overflows_enabled bool
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
func New_MOH_DAM[KT I_Positive_Integer, VT any](
	expected_num_inputs KT,
) *DAM_MOH[KT, VT] {
	return _inner_New_MOH_DAM[KT, VT](expected_num_inputs, true)
}

func _inner_New_MOH_DAM[KT I_Positive_Integer, VT any](
	expected_num_inputs KT,
	enable_overflows bool,
) *DAM_MOH[KT, VT] {
	expected_num_inputs = max(64, next_power_of_two(expected_num_inputs))
	num_buckets := max(2, expected_num_inputs/64)

	if num_buckets%2 != 0 {
		panic("numBuckets should be a multiple of 2.")
	}

	// Allocate buckets...
	num_buckets_runtime := uint64(num_buckets)
	buckets := make([]t_bucket_moh[KT, VT], num_buckets_runtime)

	// Instantiate...
	inst := DAM_MOH[KT, VT]{
		buckets:           buckets,
		num_buckets_m1:    num_buckets - 1,
		overflows_enabled: enable_overflows,
	}

	if enable_overflows {
		inst.overflows = make(map[KT]VT)
	}

	return &inst
}

func (m *DAM_MOH[KT, VT]) Enquire_Number_Of_Buckets() KT {
	return m.num_buckets_m1 + 1
}

// Set a key-value pair in the map.
// Will panic if something goes wrong.
//
// - WARNING: This function is NOT thread-safe.
//
//go:inline
func (m *DAM_MOH[KT, VT]) Set(key KT, value VT) {
	if key == 0 {
		panic("Key cannot be 0.")
	}

	index := key & m.num_buckets_m1
	buck := &m.buckets[index]

	for i := 0; i < 64; i++ {
		if buck.keys[i] == uint64(key) {
			buck.values[i] = value
			return
		}
	}

	for i := 0; i < 64; i++ {
		if buck.keys[i] == 0 {
			buck.keys[i] = uint64(key)
			buck.values[i] = value
			return
		}
	}

	if m.overflows_enabled {
		m.overflows[key] = value
		return
	}

	panic("no space left in bucket")
}

// The runtime overhead was too much:
// //go:noescape
// //go:nosplit
// //go:nobounds
// func simd_find_idx(ptr *uint64, pad uint64, len uint64, key uint64) uint8
//
// - WARNING: This function is NOT thread-safe.
//
// - NOTE: Remember that keys cannot be 0.
//
// - NOTE: This function will not check if the key is 0.
//
//go:inline
// func (m *DAM[KT, VT]) Find(key KT) uint8 {
// 	// NOTE: Keeping value type here improves performance since we do not modify the value.
// 	buck := m.buckets[key&m.num_buckets_m1]
//
// 	starting_point := uint8(m.extras[key/8]>>int(key%8)) & 1
//
// 	var my_var uint64
// 	i := uint8(0)
// 	for ; i < PREALLOC_SPACE_PER_BUCKET/2; i++ {
// 		my_var = uint64(buck.keys[starting_point+i*2])
// 		m.candidates[i] = my_var
// 	}
//
//	if len(candidates)%4 != 0 {
//		padding := 4 - len(candidates)%4
//		for i := 0; i < padding; i++ {
//			candidates = append(candidates, 0)
//		}
//	}
//
// 	//runtime.LockOSThread()
// 	v := simd_find_idx(&m.candidates[0], 0, PREALLOC_SPACE_PER_BUCKET, uint64(key))
// 	//runtime.UnlockOSThread()
// 	x := starting_point + (v * 2) - 2
//
// 	//if buck.keys[x] == key {
// 	return x + 1
// 	//}
//
// 	panic("Not implemented.")
// }

// Warning: This function must be called with keys satisfying the following conditions:
// 1. len(keys) % 4 == 0
// 2. no duplicate keys.
//
//go:noescape
//go:nosplit
func simd_find_idx(ptr *uint64, length uint64, key uint64) uint8

// Returns the value and a boolean indicating whether the value was found.
//
// - WARNING: This function is NOT thread-safe.
//
// - NOTE: Remember that keys cannot be 0.
//
// - NOTE: This function will not check if the key is 0.
//
//go:inline
func (m *DAM_MOH[KT, VT]) Get(key KT) (VT, bool) {
	// NOTE: Keeping value type here improves performance since we do not modify the value.
	buck := m.buckets[key&m.num_buckets_m1]

	// Took me an entire day coding assembly for a 5ns improvement. :facepalm:
	// NOTE: From what i can tell, the below optimization is only faster sometimes, other times it is slower.
	v := simd_find_idx(&buck.keys[0], 64, uint64(key))
	if v != 0 {
		return buck.values[v-1], true
	}

	if m.overflows_enabled {
		res, found := m.overflows[key]
		if found {
			return res, true
		}
	}

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
func (m *DAM_MOH[KT, VT]) Delete(key KT) bool {
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
