/*/
 ** This software is covered by the MIT License.
 ** See: `./LICENSE`.
/*/

package dam

// Super-Fast Direct-Access Map.
type DAM_FAST[KT I_Positive_Integer, VT any] struct {
	first_keys      []KT
	first_values    []VT
	overflow_keys   [][]KT
	overflow_values [][]VT

	num_buckets_m1 KT
}

// Creates a new DAM (Direct-Access Map) that sacrifices memory usage for speed.
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
func (_ DAM_FAST[KT, VT]) New(
	expected_num_inputs KT,
) *DAM_FAST[KT, VT] {
	expected_num_inputs = max(128, next_power_of_two(expected_num_inputs))
	num_buckets := max(2, expected_num_inputs)

	if num_buckets%2 != 0 {
		panic("numBuckets should be a multiple of 2.")
	}

	// Allocate buckets...
	num_buckets_runtime := uint64(num_buckets)

	// Instantiate...
	inst := DAM_FAST[KT, VT]{
		first_keys:      make([]KT, num_buckets_runtime),
		overflow_keys:   make([][]KT, num_buckets_runtime),
		overflow_values: make([][]VT, num_buckets_runtime),
		first_values:    make([]VT, num_buckets_runtime),
		num_buckets_m1:  num_buckets - 1,
	}

	return &inst
}

func (m *DAM_FAST[KT, VT]) Enquire_Number_Of_Buckets() KT {
	return m.num_buckets_m1 + 1
}

// Set a key-value pair in the map.
// Will panic if something goes wrong.
//
// - WARNING: This function is NOT thread-safe.
//
//go:inline
func (m *DAM_FAST[KT, VT]) Set(key KT, value VT) {
	if key == 0 {
		panic("Key cannot be 0.")
	}

	index := key & m.num_buckets_m1

	if m.first_keys[index] == key {
		m.first_values[index] = value
		return
	}

	if m.first_keys[index] == 0 {
		m.first_keys[index] = key
		m.first_values[index] = value
		return
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
func (m *DAM_FAST[KT, VT]) Get(key KT) (VT, bool) {
	// NOTE: Keeping value type here improves performance since we do not modify the value.
	index := key & m.num_buckets_m1

	if m.first_keys[index] == key {
		return m.first_values[index], true
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

//go:inline
func (m *DAM_FAST[KT, VT]) Delete(key KT) bool {

	panic("Not implemented.")
}
