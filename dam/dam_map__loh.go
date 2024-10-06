/*/
 ** This software is covered by the MIT License.
 ** See: `./LICENSE`.
/*/

package dam

type t_bucket_loh[KT I_Positive_Integer, VT any] struct {
	entries   [256]t_bucket_entry[KT, VT]
	overflows map[KT]VT
}

// Super-Fast Direct-Access Map.
type DAM_LOH[KT I_Positive_Integer, VT any] struct {
	buckets        []t_bucket_loh[KT, VT]
	num_buckets_m1 KT

	overflows_enabled bool
}

// Creates a new DAM (Direct-Access Map) that tries to balance speed and memory usage, strongly prefers memory savings.
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
func New_LOH_DAM[KT I_Positive_Integer, VT any](
	expected_num_inputs KT,
) *DAM_LOH[KT, VT] {
	return _inner_New_LOH_DAM[KT, VT](uint64(expected_num_inputs), true)
}

func _inner_New_LOH_DAM[KT I_Positive_Integer, VT any](
	expected_num_inputs uint64,
	enable_overflows bool,
) *DAM_LOH[KT, VT] {
	expected_num_inputs_runtime := max(256, uint64(next_power_of_two(expected_num_inputs)))
	num_buckets_runtime := max(2, expected_num_inputs_runtime/256)

	if num_buckets_runtime%2 != 0 {
		panic("numBuckets should be a multiple of 2.")
	}

	// Allocate buckets...
	buckets := make([]t_bucket_loh[KT, VT], num_buckets_runtime)

	// Instantiate...
	inst := DAM_LOH[KT, VT]{
		buckets:           buckets,
		num_buckets_m1:    KT(num_buckets_runtime - 1),
		overflows_enabled: enable_overflows,
	}

	if enable_overflows {
		for i := 0; i < len(inst.buckets); i++ {
			inst.buckets[i].overflows = make(map[KT]VT)
		}
	}

	return &inst
}

func (m *DAM_LOH[KT, VT]) Enquire_Number_Of_Buckets() KT {
	return KT(len(m.buckets))
}

// Set a key-value pair in the map.
// Will panic if something goes wrong.
//
// - WARNING: This function is NOT thread-safe.
//
//go:inline
func (m *DAM_LOH[KT, VT]) Set(key KT, value VT) {
	if key == 0 {
		panic("Key cannot be 0.")
	}

	index := key & m.num_buckets_m1
	for KT(len(m.buckets)) <= index {
		x := t_bucket_loh[KT, VT]{}
		m.buckets = append(m.buckets, x)
	}
	buck := &m.buckets[index]

	for i := 0; i < 256; i++ {
		if buck.entries[i].key == key {
			buck.entries[i].value = value
			return
		}
	}

	for i := 0; i < 256; i++ {
		if buck.entries[i].key == 0 {
			buck.entries[i].key = key
			buck.entries[i].value = value
			return
		}
	}

	if m.overflows_enabled {
		buck.overflows[key] = value
		return
	}

	panic("No space left in the bucket.")
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
	index := key & m.num_buckets_m1
	if KT(len(m.buckets)) <= index {
		var zero VT
		return zero, false
	}
	buck := m.buckets[index]

	for i := 0; i < 256; i += 1 {
		if buck.entries[i].key == key {
			return buck.entries[i].value, true
		}
	}

	if m.overflows_enabled {
		res, found := buck.overflows[key]
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
