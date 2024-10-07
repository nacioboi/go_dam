/*/
 ** This software is covered by the MIT License.
 ** See: `./LICENSE`.
/*/

package dam

type t_bucket_fast[KT I_Positive_Integer, VT any] struct {
	first_key KT
}

// Super-Fast Direct-Access Map.
type DAM_FAST[KT I_Positive_Integer, VT any] struct {
	buckets        []t_bucket_fast[KT, VT]
	values         []VT
	overflows      map[KT]VT
	num_buckets_m1 KT

	users_chosen_hash_func func(KT) uint64
	using_users_hash_func  bool
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
func New_Fast_DAM[KT I_Positive_Integer, VT any](
	expected_num_inputs KT,
) *DAM_FAST[KT, VT] {
	expected_num_inputs = max(128, next_power_of_two(expected_num_inputs))
	num_buckets := max(2, expected_num_inputs)

	if num_buckets%2 != 0 {
		panic("numBuckets should be a multiple of 2.")
	}

	// Allocate buckets...
	num_buckets_runtime := uint64(num_buckets)
	buckets := make([]t_bucket_fast[KT, VT], num_buckets_runtime)

	// Instantiate...
	inst := DAM_FAST[KT, VT]{
		buckets:        buckets,
		overflows:      make(map[KT]VT),
		values:         make([]VT, num_buckets_runtime),
		num_buckets_m1: num_buckets - 1,
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
	buck := &m.buckets[index]

	if buck.first_key == key {
		m.values[index] = value
		return
	}

	if buck.first_key == 0 {
		buck.first_key = key
		m.values[index] = value
		return
	}

	//second_hash := (index & (m.num_buckets_m1 >> 1)) + 1
	m.overflows[key] = value
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
	buck := m.buckets[index]

	if buck.first_key == key {
		return m.values[index], true
	}

	//second_hash := (index & (m.num_buckets_m1 >> 1)) + 1
	res, found := m.overflows[key]
	if found {
		return res, true
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
func (m *DAM_FAST[KT, VT]) Delete(key KT) bool {
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
