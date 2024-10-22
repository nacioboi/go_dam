package dam

import "fmt"

const (
	////////////////////////////// 0. 123456789
	/////////////////////////////////////// 1. 0123456789
	///////////////////////////////////////////////// 2. 0123456789
	/////////////////////////////////////////////////////////// 3. 0123456789
	///////////////////////////////////////////////////////////////////// 4. 0123456789
	/////////////////////////////////////////////////////////////////////////////// 5. 0123456789
	///////////////////////////////////////////////////////////////////////////////////////// 6. 01234
	c__BITMASK__64th_B uint64 = 0b1000000000000000000000000000000000000000000000000000000000000000
	c__BITMASK__BUT_64 uint64 = 0b0111111111111111111111111111111111111111111111111111111111111111
	/////////////////////////////////// 123456789012345 // 15 bits metadata.
	c__BITMASK__EXTRA_INFOS uint64 = 0b0111000000000000000000000000000000000000000000000000000000000000
	/////////////////////////////////// 123 // a sign bit for each difference value.
	c__BITMASK__SIGNS_INFOS uint64 = 0b0110000000000000000000000000000000000000000000000000000000000000
	// Bitmask for a single difference value.
	c__BITMASK__OUR_BITS uint64 = (1<<30 - 1)

	// The amount of right shift required to bring the extras to the start of the 64-bit integer.
	c__SHIFT_WIDTH__EXTRAS uint64 = 60

	c__NUM_BITS_PER_VALUE       uint64 = 30
	c__MAX_NUM_VALUES_PER_ENTRY uint64 = 2
)

// The compressed array stores information in two types of entries:
//
// - Data: Contains multiple (up to 16) values. They are the difference between said value and the checkpoint that precedes.
//
// - Checkpoint: Contains an optimized base value.
type Compressed_Array struct {
	D           *[]uint64
	spare       []uint64
	remappings  map[uint64]uint32
	safety_flag bool
}

func New_Compressed_Array() *Compressed_Array {
	data := make([]uint64, 0)
	return &Compressed_Array{
		D:          &data,
		spare:      make([]uint64, 0),
		remappings: make(map[uint64]uint32),
	}
}

func (ca *Compressed_Array) is_checkpoint(idx uint64) bool {
	v := (*ca.D)[idx]
	return v&c__BITMASK__64th_B == c__BITMASK__64th_B
}

func (ca *Compressed_Array) get_checkpoint_value(idx uint64) uint64 {
	if !ca.is_checkpoint(idx) {
		panic("not a checkpoint")
	}
	return (*ca.D)[idx] & c__BITMASK__BUT_64
}

func (ca *Compressed_Array) get_num_attached_values(checkpoint_idx uint64) uint64 {
	num_attached_values := uint64(0)

	for i := checkpoint_idx + 1; i < uint64(len(*ca.D)); i++ {
		if ca.is_checkpoint(i) {
			break
		}
		num_attached_values++
	}

	return num_attached_values
}

func (ca *Compressed_Array) needs_new_checkpoint(checkpoint_idx uint64, value uint64) bool {
	diff, _, _ := ca.calculate_difference(ca.get_checkpoint_value(checkpoint_idx), value)

	if diff >= (c__BITMASK__OUR_BITS) {
		return true
	}

	return false
}

func (ca *Compressed_Array) get_previous_checkpoint_idx(starting_idx uint64) uint64 {
	for i := uint64(starting_idx + 1); i > 0; i-- {
		if ca.is_checkpoint(i - 1) {
			return i - 1
		}
	}
	panic("something went horribly wrong")
}

func (ca *Compressed_Array) create_checkpoint(recommended_value uint64) {
	// Create a new checkpoint with the recommended value
	fmt.Printf("create_checkpoint: %d\n", recommended_value)
	checkpoint := c__BITMASK__64th_B | recommended_value
	(*ca.D) = append((*ca.D), checkpoint)
}

// Calculates value - checkpoint_value, checks if it fits in 16 bits.
// If it doesn't fit, it creates a new checkpoint.
//
// # Returns a (uint64, bool, bool) where
//
// - the first value is the difference between the value and the checkpoint_value,
//
// - the second value is a boolean that indicates if a new checkpoint should be created.
//
// - the third value is a boolean that indicates if the difference is negative.
//
// Note that if a checkpoint is created, the difference must be ignored
func (ca *Compressed_Array) calculate_difference(checkpoint_value, value uint64) (uint64, bool, bool) {
	var is_negative bool
	diff := int64(value - checkpoint_value)
	if diff < 0 {
		diff = -diff
		is_negative = true
	}
	if diff >= int64(c__BITMASK__OUR_BITS) {
		return 0, true, is_negative
	}
	return uint64(diff), false, is_negative
}

func (ca *Compressed_Array) apply_sign_infos(idx uint64, j uint64, is_negative bool) {
	sign_infos := (*ca.D)[idx] & c__BITMASK__SIGNS_INFOS
	shifted := sign_infos >> (c__SHIFT_WIDTH__EXTRAS)
	m := uint64(1) << (c__MAX_NUM_VALUES_PER_ENTRY - j)
	if is_negative {
		shifted |= m
	} else {
		shifted &= (^m & 0b1111)
	}
	mask := shifted << (c__SHIFT_WIDTH__EXTRAS)
	(*ca.D)[idx] |= mask
}

func (ca *Compressed_Array) get_is_negative(idx, j uint64) bool {
	sign_infos := (*ca.D)[idx] & c__BITMASK__SIGNS_INFOS
	shifted := sign_infos >> (c__SHIFT_WIDTH__EXTRAS)
	m := uint64(1) << (c__MAX_NUM_VALUES_PER_ENTRY - j)
	return shifted&m != 0
}

func (ca *Compressed_Array) search_for_matching_checkpoint(value uint64) (uint64, bool) {
	fmt.Printf("search_for_matching_checkpoint: %d\n", value)
	// Search for a checkpoint that fits the value
	for i := uint64(0); i < uint64(len(*ca.D)); i++ {
		if !ca.is_checkpoint(i) {
			continue
		}
		checkpoint_value := ca.get_checkpoint_value(i)
		_, needs_new_checkpoint, _ := ca.calculate_difference(checkpoint_value, value)
		if !needs_new_checkpoint {
			return i, true
		}
	}
	return 0, false
}

func (ca *Compressed_Array) handle_restructuring(matching_checkpoint uint64, value uint64) {
	fmt.Printf("matching_checkpoint: %d\n", matching_checkpoint)
	if ca.safety_flag {
		panic("safety flag is set, recursion is not allowed")
	}

	// A restructure occurs when we need to insert an item.
	ca.spare = make([]uint64, 0)
	the_rest := make([]uint64, 0)
	i := uint64(0)
	for j, v := range *ca.D {
		if j == int(matching_checkpoint) {
			break
		}
		ca.spare = append(ca.spare, v)
		i++
	}
	for j := matching_checkpoint + 1; j < uint64(len(*ca.D)); j++ {
		the_rest = append(the_rest, (*ca.D)[j])
	}

	// Insert the new value
	ca.D = &ca.spare

	// Simply use the append function to add the rest of the values.
	ca.safety_flag = true
	defer func() { ca.safety_flag = false }()

	ca.Append(value)

	for _, v := range the_rest {
		ca.Append(v)
	}
}

func (ca *Compressed_Array) Append(value uint64) {
	//fmt.Printf("data: %v\n", ca.data)

	// Get the previous checkpoint index...
	// TODO: This is just an example since we don't always want to get the last checkpoint.
	// E.g. If the diff is too large, we want to search all checkpoint and check if one exists that fits
	//        before adding a new checkpoint.
	if len(*ca.D) == 0 {
		ca.create_checkpoint(value)
		return
	}

	// Get some values...
	checkpoint_idx := ca.get_previous_checkpoint_idx(uint64(len(*ca.D) - 1))
	checkpoint_value := ca.get_checkpoint_value(checkpoint_idx)

	// Again, example implementation for now: if we cannot fit the value, create a new checkpoint...
	if ca.needs_new_checkpoint(checkpoint_idx, value) {
		matching_checkpoint, did_find_matching_checkpoint := ca.search_for_matching_checkpoint(value)
		if did_find_matching_checkpoint {
			ca.handle_restructuring(matching_checkpoint, value)
		} else {
			ca.create_checkpoint(value)
		}
		return
	}

	// We use a zero bit string to indicate that the bit string is free.
	// For this reason we cannot add a diff that is zero...
	_diff := int64(value - checkpoint_value)
	if _diff < 0 {
		_diff = -_diff
	}
	diff := uint64(_diff)
	if diff == 0 {
		matching_checkpoint, did_find_matching_checkpoint := ca.search_for_matching_checkpoint(value)
		if did_find_matching_checkpoint {
			ca.handle_restructuring(matching_checkpoint, value)
		} else {
			ca.create_checkpoint(value)
		}
		return
	}

	// Otherwise we need to add the value to a normal data entry by means of checking for a free bit string.

	for i := checkpoint_idx + 1; i < uint64(len(*ca.D)); i++ {
		if ca.is_checkpoint(i) {
			break // Stop searching since we have reached the next checkpoint.
		}
		for j := uint64(0); j < c__MAX_NUM_VALUES_PER_ENTRY; j++ {
			if (*ca.D)[i]&(c__BITMASK__OUR_BITS<<(j*c__NUM_BITS_PER_VALUE)) != 0 {
				continue
			}
			diff, does_need_new_checkpoint, is_negative := ca.calculate_difference(checkpoint_value, value)
			if !does_need_new_checkpoint {
				(*ca.D)[i] |= uint64(diff) << (j * c__NUM_BITS_PER_VALUE)
				ca.apply_sign_infos(i, j, is_negative)
			} else {
				// We need to create a new checkpoint.
				matching_checkpoint, did_find_matching_checkpoint := ca.search_for_matching_checkpoint(value)
				if did_find_matching_checkpoint {
					ca.handle_restructuring(matching_checkpoint, value)
				} else {
					ca.create_checkpoint(value)
				}
			}
			return
		}
	}

	diff, does_need_new_checkpoint, is_negative := ca.calculate_difference(checkpoint_value, value)
	if !does_need_new_checkpoint {
		(*ca.D) = append((*ca.D), uint64(diff))
		ca.apply_sign_infos(uint64(len(*ca.D)-1), 0, is_negative)
	} else {
		// We need to create a new checkpoint.
		matching_checkpoint, did_find_matching_checkpoint := ca.search_for_matching_checkpoint(value)
		if did_find_matching_checkpoint {
			ca.handle_restructuring(matching_checkpoint, value)
		} else {
			ca.create_checkpoint(value)
		}
	}
}

func (ca *Compressed_Array) Get(idx uint64) uint64 {
	var (
		current_checkpoint_value uint64
		current_logical_index    uint64
	)

	for i := uint64(0); i < uint64(len(*ca.D)); i++ {
		if ca.is_checkpoint(i) {
			current_checkpoint_value = ca.get_checkpoint_value(i)
			if current_logical_index == idx {
				return current_checkpoint_value
			}
			current_logical_index++
		} else {
			for j := uint64(0); j < c__MAX_NUM_VALUES_PER_ENTRY; j++ {
				v := ((*ca.D)[i] >> (j * c__NUM_BITS_PER_VALUE)) & c__BITMASK__OUR_BITS
				if v == 0 {
					continue
				}
				if current_logical_index == idx {
					is_negative := ca.get_is_negative(i, j)
					if is_negative {
						return current_checkpoint_value - v
					} else {
						return current_checkpoint_value + v
					}
				}
				current_logical_index++
			}
		}
	}

	panic("Index out of bounds")
}
