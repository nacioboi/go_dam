package dam

////////////////////////////// 0. 123456789
//
/////////////////////////////////////// 1. 0123456789
//
///////////////////////////////////////////////// 2. 0123456789
//
/////////////////////////////////////////////////////////// 3. 0123456789
//
///////////////////////////////////////////////////////////////////// 4. 0123456789
//
/////////////////////////////////////////////////////////////////////////////// 5. 0123456789
//
///////////////////////////////////////////////////////////////////////////////////////// 6. 01234
const SIXTY_FORTH_BIT_BITMASK = 0b1000000000000000000000000000000000000000000000000000000000000000
const ALL_BITS_BUT_64_BITMASK = 0b0111111111111111111111111111111111111111111111111111111111111111

const MAX_NUM_VALUES_PER_ENTRY = 3
const NUM_BITS_PER_VALUE = (64 / MAX_NUM_VALUES_PER_ENTRY) - 1

const ALL_16_BITS_BITMASK = 0xFFFF

// The compressed array stores information in two types of entries:
//
// - Data: Contains multiple (up to 16) values. They are the difference between said value and the checkpoint that precedes.
//
// - Checkpoint: Contains an optimized base value.
type Compressed_Array struct {
	data []uint64
}

func New_Compressed_Array() *Compressed_Array {
	return &Compressed_Array{
		data: make([]uint64, 0),
	}
}

func (ca *Compressed_Array) is_checkpoint(idx uint64) bool {
	return ca.data[idx]&SIXTY_FORTH_BIT_BITMASK == SIXTY_FORTH_BIT_BITMASK
}

func (ca *Compressed_Array) get_checkpoint_value(idx uint64) uint64 {
	if !ca.is_checkpoint(idx) {
		panic("not a checkpoint")
	}
	return ca.data[idx] & ALL_BITS_BUT_64_BITMASK
}

func (ca *Compressed_Array) get_num_attached_values(checkpoint_idx uint64) uint64 {
	num_attached_values := uint64(0)

	for i := checkpoint_idx + 1; i < uint64(len(ca.data)); i++ {
		if ca.is_checkpoint(i) {
			break
		}
		num_attached_values++
	}

	return num_attached_values
}

func (ca *Compressed_Array) needs_new_checkpoint(checkpoint_idx uint64, value uint64) bool {
	checkpoint_value := ca.get_checkpoint_value(checkpoint_idx)
	diff := int64(value - checkpoint_value)

	if diff < 0 {
		panic("not implemented")
	}

	if diff > 0xFFFF {
		return false
	}

	return true
}

func (ca *Compressed_Array) get_previous_checkpoint_idx(starting_idx uint64) uint64 {
	for i := uint64(starting_idx + 1); i > 0; i-- {
		if ca.is_checkpoint(i - 1) {
			return i - 1
		}
	}
	panic("something went horribly wrong")
}

// // Calculate the optimal checkpoint value using the values of all the data attached to it.
// // A checkpoint value that is un-optimized could be really bad for the compression ratio.
// func (ca *Compressed_Array) find_optimal_checkpoint_value(checkpoint_idx uint64) uint64 {
// 	// For now we can just use the minimum value.

// 	optimal_value := uint64(0)
// 	has_set_first_value := false

// 	num_attached_values := ca.get_num_attached_values(checkpoint_idx)

// 	if num_attached_values == 0 {
// 		return ca.get_checkpoint_value(checkpoint_idx)
// 	}

// 	for i := uint64(checkpoint_idx); i < checkpoint_idx+num_attached_values; i++ {
// 		if !has_set_first_value {
// 			optimal_value = ca.data[i]
// 			has_set_first_value = true
// 			continue
// 		}
// 		if ca.data[i] < optimal_value {
// 			optimal_value = ca.data[i]
// 		}
// 	}

// 	return optimal_value
// }

func (ca *Compressed_Array) create_checkpoint(recommended_value uint64) {
	// Create a new checkpoint with the recommended value
	checkpoint := SIXTY_FORTH_BIT_BITMASK | recommended_value
	ca.data = append(ca.data, checkpoint)
}

func (ca *Compressed_Array) Append(value uint64) {
	//fmt.Printf("data: %v\n", ca.data)

	if len(ca.data) == 0 {
		ca.create_checkpoint(value)
		return
	}

	checkpoint_idx := ca.get_previous_checkpoint_idx(uint64(len(ca.data) - 1))

	if !ca.needs_new_checkpoint(checkpoint_idx, value) {
		ca.create_checkpoint(value)

		return
	}

	//checkpoint_value := ca.find_optimal_checkpoint_value(checkpoint_idx)
	checkpoint_value := ca.get_checkpoint_value(checkpoint_idx)

	// Find the first entry that has space for the new value.
	for i := checkpoint_idx + 1; i < uint64(len(ca.data)); i++ {
		if ca.is_checkpoint(i) {
			diff := int64(value - checkpoint_value)
			if diff < 0 {
				panic("not implemented")
			}
			if diff > 0xFFFF {
				// We need to create a new checkpoint.
				ca.create_checkpoint(value)
				return
			}
			ca.data = append(ca.data, uint64(diff))
			return
		}
		for j := uint64(0); j < MAX_NUM_VALUES_PER_ENTRY; j++ {
			if ca.data[i]&(ALL_16_BITS_BITMASK<<(j*16)) != 0 {
				continue
			}
			diff := int64(value - checkpoint_value)
			if diff < 0 {
				panic("not implemented")
			}
			if diff > 0xFFFF {
				// We need to create a new checkpoint.
				ca.create_checkpoint(value)
				return
			}
			//fmt.Printf("diff: %v\n", diff)
			//fmt.Printf("value: %v\n", value)
			ca.data[i] |= uint64(diff) << (j * 16)
			return
		}
	}

	diff := int64(value - checkpoint_value)
	//fmt.Printf("diff: %v\n", diff)
	//fmt.Printf("value: %v\n", value)
	if diff < 0 {
		panic("not implemented")
		// normalized_diff := uint64(-diff)
		// ca.data = append(ca.data, normalized_diff)
		// ca.data[len(ca.data)-1] |= (1 << (63 - 1))
		// return
	}
	if diff > 0xFFFF {
		ca.create_checkpoint(value)
		return
	}
	ca.data = append(ca.data, uint64(diff))
}

func (ca *Compressed_Array) Get(idx uint64) uint64 {
	var (
		current_checkpoint_value uint64
		current_logical_index    uint64
	)

	for i := uint64(0); i < uint64(len(ca.data)); i++ {
		if ca.is_checkpoint(i) {
			current_checkpoint_value = ca.get_checkpoint_value(i)
			current_logical_index++
			if (current_logical_index - 1) == idx {
				return current_checkpoint_value
			}
		} else {
			for j := uint64(0); j < MAX_NUM_VALUES_PER_ENTRY; j++ {
				v := (ca.data[i] >> (j * 16)) & 0xFFFF
				if v != 0 {
					current_logical_index++
				}
				if (current_logical_index - 1) == idx {
					return current_checkpoint_value + v
				}
			}
		}
	}

	// If we reach here, the index is out of bounds
	panic("Index out of bounds")
}
