/*/
 ** This software is covered by the MIT License.
 ** See: `./LICENSE`.
/*/

package dam

const _MOH_DAM__NUM_ITEMS_PER_BUCKET = 64

// Super-Fast Direct-Access Map.
type DAM_MOH[KT I_Positive_Integer, VT any] struct {
	keys            [][_MOH_DAM__NUM_ITEMS_PER_BUCKET]uint64
	values          [][_MOH_DAM__NUM_ITEMS_PER_BUCKET]VT
	overflow_keys   [][]uint64
	overflow_values [][]VT
	len             uint64

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
func (_ DAM_MOH[KT, VT]) New(
	expected_num_inputs KT,
) *DAM_MOH[KT, VT] {
	return _inner_New_MOH_DAM[KT, VT](expected_num_inputs, true)
}

func _inner_New_MOH_DAM[KT I_Positive_Integer, VT any](
	expected_num_inputs KT,
	enable_overflows bool,
) *DAM_MOH[KT, VT] {
	expected_num_inputs = max(128, next_power_of_two(expected_num_inputs))
	num_buckets := max(2, expected_num_inputs/_MOH_DAM__NUM_ITEMS_PER_BUCKET)

	if num_buckets%2 != 0 {
		panic("numBuckets should be a multiple of 2.")
	}

	// Allocate buckets...
	num_buckets_runtime := uint64(num_buckets)

	// Instantiate...
	inst := DAM_MOH[KT, VT]{
		keys:              make([][_MOH_DAM__NUM_ITEMS_PER_BUCKET]uint64, num_buckets_runtime),
		values:            make([][_MOH_DAM__NUM_ITEMS_PER_BUCKET]VT, num_buckets_runtime),
		num_buckets_m1:    num_buckets - 1,
		overflows_enabled: enable_overflows,
	}

	if enable_overflows {
		inst.overflow_keys = make([][]uint64, num_buckets_runtime)
		inst.overflow_values = make([][]VT, num_buckets_runtime)
	}

	return &inst
}

func (m *DAM_MOH[KT, VT]) Enquire_Number_Of_Buckets() KT {
	return KT(m.num_buckets_m1 + 1)
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

	for i := 0; i < _MOH_DAM__NUM_ITEMS_PER_BUCKET; i++ {
		if m.keys[index][i] == uint64(key) {
			m.values[index][i] = value
			return
		}
	}

	for i := 0; i < _MOH_DAM__NUM_ITEMS_PER_BUCKET; i++ {
		if m.keys[index][i] == 0 {
			m.keys[index][i] = uint64(key)
			m.values[index][i] = value
			return
		}
	}

	if m.overflows_enabled {
		for i := 0; i < len(m.overflow_keys[index]); i++ {
			if m.overflow_keys[index][i] == uint64(key) {
				m.overflow_values[index][i] = value
			}
		}
		m.overflow_keys[index] = append(m.overflow_keys[index], uint64(key))
		m.overflow_values[index] = append(m.overflow_values[index], value)
		return
	}

	panic("no space left in bucket")
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
func (m *DAM_MOH[KT, VT]) Get(key KT) (VT, bool) {
	// NOTE: Keeping value type here improves performance since we do not modify the value.
	index := key & m.num_buckets_m1

	v, ok := avx512_find_idx_64i(uint64(key), &m.keys[index][0])
	if ok {
		return m.values[index][v], true
	}

	// if m.overflows_enabled {
	// 	// Fetch from overflow...
	// 	for i := 0; i < len(m.overflow_keys[index]); i++ {
	// 		if m.overflow_keys[index][i] == uint64(key) {
	// 			return m.overflow_values[index][i], true
	// 		}
	// 	}
	// }

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
