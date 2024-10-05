/*/
 ** This software is covered by the MIT License.
 ** See: `./LICENSE`.
/*/

package dam

type I_Positive_Integer interface {
	uint8 | uint16 | uint32 | uint64
}

type t_bucket_entry[KT I_Positive_Integer, VT any] struct {
	key   KT
	value VT
}

type bucket[KT I_Positive_Integer, VT any] struct {
	entries []t_bucket_entry[KT, VT]
}

// Super-Fast Direct-Access Map.
type DAM[KT I_Positive_Integer, VT any] struct {
	buckets        []bucket[KT, VT]
	num_buckets_m1 KT

	users_chosen_hash_func func(KT) uint64
	using_users_hash_func  bool

	profile T_Performance_Profile
}

func New[KT I_Positive_Integer, VT any](
	expected_num_inputs KT,
	options ...T_Option[KT, VT],
) *DAM[KT, VT] {
	expected_num_inputs = next_power_of_two(expected_num_inputs)

	profile := PERFORMANCE_PROFILE__SAVE_MEMORY
	for _, opt := range options {
		if opt.t == OPTION_TYPE__WITH_PERFORMANCE_PROFILE {
			profile = opt.other.(T_Performance_Profile)
		}
	}

	var num_buckets KT
	switch profile {
	case PERFORMANCE_PROFILE__FAST:
		num_buckets = expected_num_inputs / 2
	case PERFORMANCE_PROFILE__NORMAL:
		num_buckets = expected_num_inputs / 4
	case PERFORMANCE_PROFILE__SAVE_MEMORY:
		num_buckets = expected_num_inputs / 8
	default:
		panic("Invalid performance profile.")
	}

	if num_buckets%2 != 0 {
		panic("numBuckets should be a multiple of 2.")
	}

	// Allocate buckets...
	num_buckets_runtime := any(num_buckets).(uint64)
	buckets := make([]bucket[KT, VT], num_buckets_runtime)
	estimated_num_entries_per_bucket := expected_num_inputs / num_buckets
	for i := uint64(0); i < num_buckets_runtime; i++ {
		b := bucket[KT, VT]{
			entries: make([]t_bucket_entry[KT, VT], 0, estimated_num_entries_per_bucket),
		}
		buckets[i] = b
	}

	// Instantiate...
	inst := DAM[KT, VT]{
		buckets:        buckets,
		num_buckets_m1: num_buckets - 1,
		profile:        profile,
	}

	// Apply options...
	for _, opt := range options {
		if opt.t != OPTION_TYPE__WITH_PERFORMANCE_PROFILE {
			opt.f(&inst)
		}
	}

	return &inst
}

func (m *DAM[KT, VT]) Enquire_Number_Of_Buckets() KT {
	return m.num_buckets_m1 + 1
}

// Set a key-value pair in the map.
// Will panic if something goes wrong.
//
// - WARNING: This function is NOT thread-safe.
//
//go:inline
func (m *DAM[KT, VT]) Set(key KT, value VT) {
	if key == 0 {
		panic("Key cannot be 0.")
	}

	index := key & m.num_buckets_m1
	buck := &m.buckets[index]

	for i := 0; i < len(buck.entries); i++ {
		if buck.entries[i].key == key {
			buck.entries[i].value = value
			return
		}
	}

	buck.entries = append(buck.entries, t_bucket_entry[KT, VT]{key: key, value: value})
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

// Returns the value and a boolean indicating whether the value was found.
//
// - WARNING: This function is NOT thread-safe.
//
// - NOTE: Remember that keys cannot be 0.
//
// - NOTE: This function will not check if the key is 0.
//
//go:inline
func (m *DAM[KT, VT]) Get(key KT) (VT, bool) {
	// NOTE: Keeping value type here improves performance since we do not modify the value.
	buck := m.buckets[key&m.num_buckets_m1]

	for i := 0; i < len(buck.entries); i += 1 {
		if buck.entries[i].key == key {
			return buck.entries[i].value, true
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
func (m *DAM[KT, VT]) Delete(key KT) bool {
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
