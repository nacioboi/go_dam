/*/
 ** This software is covered by the MIT License.
 ** See: `./LICENSE`.
/*/

package dam

const _STD_DAM__NUM_ITEMS_PER_BUCKET = 4

// Super-Fast Direct-Access Map.
type DAM_STD[KT I_Positive_Integer, VT any] struct {
	keys            [][_STD_DAM__NUM_ITEMS_PER_BUCKET]KT
	values          [][_STD_DAM__NUM_ITEMS_PER_BUCKET]VT
	overflow_keys   [][]KT
	overflow_values [][]VT

	num_buckets_m1 KT
}

// Creates a new DAM (Direct-Access Map) that tries to balance speed and memory usage, slight preference for speed.
//
// Since this is a DAM, we need to know the expected number of inputs in advance.
// This leaves us with the following options:
//
// - `New_Fast_DAM`: Super-fast DAM, sacrifices memory usage for speed.
//
// - `New_STD_DAM`: Slightly slower DAM, gives up some speed for memory usage.
//
// - `New_MOH_DAM`: (Medium-OverHead DAM), gives up even more speed for memory usage.
//
// - `New_LOH_DAM`: (Low-OverHead DAM), sacrifices speed for memory savings.
//
func New_STD_DAM[KT I_Positive_Integer, VT any](
	expected_num_inputs KT,
) *DAM_STD[KT, VT] {
	expected_num_inputs = max(128, next_power_of_two(expected_num_inputs))
	num_buckets := max(2, expected_num_inputs/4)

	if num_buckets%2 != 0 {
		panic("numBuckets should be a multiple of 2.")
	}

	// Allocate buckets...
	num_buckets_runtime := uint64(num_buckets)

	// Instantiate...
	inst := DAM_STD[KT, VT]{
		keys:            make([][_STD_DAM__NUM_ITEMS_PER_BUCKET]KT, num_buckets_runtime),
		values:          make([][_STD_DAM__NUM_ITEMS_PER_BUCKET]VT, num_buckets_runtime),
		overflow_keys:   make([][]KT, num_buckets_runtime),
		overflow_values: make([][]VT, num_buckets_runtime),
		num_buckets_m1:  num_buckets - 1,
	}

	return &inst
}

func (m *DAM_STD[KT, VT]) Enquire_Number_Of_Buckets() KT {
	return m.num_buckets_m1 + 1
}

// Set a key-value pair in the map.
// Will panic if something goes wrong.
//
// - WARNING: This function is NOT thread-safe.
//
//go:inline
func (m *DAM_STD[KT, VT]) Set(key KT, value VT) {
	if key == 0 {
		panic("Key cannot be 0.")
	}

	index := key & m.num_buckets_m1

	for i := 0; i < _STD_DAM__NUM_ITEMS_PER_BUCKET; i++ {
		if m.keys[index][i] == key {
			m.values[index][i] = value
			return
		}
	}

	for i := 0; i < _STD_DAM__NUM_ITEMS_PER_BUCKET; i++ {
		if m.keys[index][i] == 0 {
			m.keys[index][i] = key
			m.values[index][i] = value
			return
		}
	}

	// Overflow...
	for i := 0; i < len(m.overflow_keys[index]); i++ {
		if m.overflow_keys[index][i] == key {
			m.overflow_values[index][i] = value
			return
		}
	}

	m.overflow_keys[index] = append(m.overflow_keys[index], key)
	m.overflow_values[index] = append(m.overflow_values[index], value)
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
func (m *DAM_STD[KT, VT]) Get(key KT) (VT, bool) {
	// NOTE: Keeping value type here improves performance since we do not modify the value.
	index := key & m.num_buckets_m1

	for i := 0; i < _STD_DAM__NUM_ITEMS_PER_BUCKET; i += 1 {
		if m.keys[index][i] == key {
			return m.values[index][i], true
		}
	}

	// Fetch from overflow...
	for i := 0; i < len(m.overflow_keys[index]); i++ {
		if m.overflow_keys[index][i] == key {
			return m.overflow_values[index][i], true
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
func (m *DAM_STD[KT, VT]) Delete(key KT) bool {
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
