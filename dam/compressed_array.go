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
	c__BITMASK__64th_B      uint64 = 0b1000000000000000000000000000000000000000000000000000000000000000
	c__BITMASK__BUT_64      uint64 = 0b0111111111111111111111111111111111111111111111111111111111111111
	c__BITMASK__EXTRA_INFOS uint64 = 0b0111111111111111000000000000000000000000000000000000000000000000
	c__BITMASK__SIGNS_INFOS uint64 = 0b0111100000000000000000000000000000000000000000000000000000000000
	////////////////////////////////// 432109876543210
	///////////////////////////////////////////////// 5432109876543210
	///////////////////////////////////////////////////////////////// 5432109876543210
	///////////////////////////////////////////////////////////////////////////////// 5432109876543210
	// We will use some of the bits for multipliers for our values.
	//
	// 3 bits can store from 0-7.
	//
	// We are calculating `(x+1)*v`:
	//
	// - where v is the value store in the bits to the right,
	//
	// - and x is the value stored in a 3-bit multiplier.
	//
	// This means that the multiplier 0b000 is 1*v, and 0b111 is 8*v.
	// This effectively gives us 8 times the range as storing 16 bits alone.
	//
	// MAX = 65535 * 8 = 524280
	c__BITMASK__MULTIPLIERS uint64 = 0b0000011111111100000000000000000000000000000000000000000000000000
	c__BITMASK__16_BITS     uint64 = 0xFFFF

	// The amount of right shift required to bring the extras to the start of the 64-bit integer.
	c__SHIFT_WIDTH__EXTRAS uint64 = 48
	// The amount of right shift, on-top of right shift by `c__SHIFT_WIDTH__EXTRAS`,
	// required to bring the sign infos to the start of the 64-bit integer.
	c__SHIFT_WIDTH__AND_SIGN_INFOS uint64 = 11
	c__MAX_NUM_VALUES_PER_ENTRY    uint64 = 3
)

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
	v := ca.data[idx]
	return v&c__BITMASK__64th_B == c__BITMASK__64th_B
}

func (ca *Compressed_Array) get_checkpoint_value(idx uint64) uint64 {
	if !ca.is_checkpoint(idx) {
		panic("not a checkpoint")
	}
	return ca.data[idx] & c__BITMASK__BUT_64
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
	diff, _, _ := ca.calculate_difference(ca.get_checkpoint_value(checkpoint_idx), value)

	if diff >= 0xFFFF {
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
	checkpoint := c__BITMASK__64th_B | recommended_value
	ca.data = append(ca.data, checkpoint)
}

// Calculates value - checkpoint_value, checks if it fits in 16 bits.
// If it doesn't fit, it creates a new checkpoint.
//
// # Returns a (uint64, bool, bool) where
//
// - the first value is the difference between the value and the checkpoint_value,
//
// - the second value is a boolean that indicates if a new checkpoint was created.
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
	if diff > 0xFFFF {
		// We need to create a new checkpoint.
		ca.create_checkpoint(value)
		return 0, true, is_negative
	}
	return uint64(diff), false, is_negative
}

func (ca *Compressed_Array) apply_sign_infos(idx uint64, j uint64, is_negative bool) {
	sign_infos := ca.data[idx] & c__BITMASK__SIGNS_INFOS
	shifted := sign_infos >> (c__SHIFT_WIDTH__EXTRAS + c__SHIFT_WIDTH__AND_SIGN_INFOS)
	m := uint64(1) << (4 - j)
	if is_negative {
		shifted |= m
	} else {
		shifted &= (^m & 0b1111)
	}
	mask := shifted << (c__SHIFT_WIDTH__EXTRAS + c__SHIFT_WIDTH__AND_SIGN_INFOS)
	ca.data[idx] |= mask
}

func (ca *Compressed_Array) get_is_negative(idx, j uint64) bool {
	sign_infos := ca.data[idx] & c__BITMASK__SIGNS_INFOS
	shifted := sign_infos >> (c__SHIFT_WIDTH__EXTRAS + c__SHIFT_WIDTH__AND_SIGN_INFOS)
	m := uint64(1) << (4 - j)
	return shifted&m != 0
}

func (ca *Compressed_Array) Append(value uint64) {
	//fmt.Printf("data: %v\n", ca.data)

	// Get the previous checkpoint index...
	// TODO: This is just an example since we don't always want to get the last checkpoint.
	// E.g. If the diff is too large, we want to search all checkpoint and check if one exists that fits
	//        before adding a new checkpoint.
	if len(ca.data) == 0 {
		ca.create_checkpoint(value)
		return
	}

	// Get some values...
	checkpoint_idx := ca.get_previous_checkpoint_idx(uint64(len(ca.data) - 1))
	checkpoint_value := ca.get_checkpoint_value(checkpoint_idx)

	// Again, example implementation for now: if we cannot fit the value, create a new checkpoint...
	if ca.needs_new_checkpoint(checkpoint_idx, value) {
		ca.create_checkpoint(value)
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
		ca.create_checkpoint(value)
		return
	}

	// Otherwise we need to add the value to a normal data entry by means of checking for a free bit string.

	for i := checkpoint_idx + 1; i < uint64(len(ca.data)); i++ {
		if ca.is_checkpoint(i) {
			break // Stop searching since we have reached the next checkpoint.
		}
		for j := uint64(0); j < c__MAX_NUM_VALUES_PER_ENTRY; j++ {
			if ca.data[i]&(c__BITMASK__16_BITS<<(j*16)) != 0 {
				continue
			}
			diff, did_create_new_checkpoint, is_negative := ca.calculate_difference(checkpoint_value, value)
			if !did_create_new_checkpoint {
				ca.data[i] |= uint64(diff) << (j * 16)
				ca.apply_sign_infos(i, j, is_negative)
			}
			return
		}
	}

	diff, did_create_new_checkpoint, is_negative := ca.calculate_difference(checkpoint_value, value)
	if !did_create_new_checkpoint {
		ca.data = append(ca.data, uint64(diff))
		ca.apply_sign_infos(uint64(len(ca.data)-1), 0, is_negative)
	}
}

func (ca *Compressed_Array) Get(idx uint64) uint64 {
	var (
		current_checkpoint_value uint64
		current_logical_index    uint64
	)

	for i := uint64(0); i < uint64(len(ca.data)); i++ {
		if ca.is_checkpoint(i) {
			current_checkpoint_value = ca.get_checkpoint_value(i)
			if current_logical_index == idx {
				return current_checkpoint_value
			}
			current_logical_index++
		} else {
			for j := uint64(0); j < c__MAX_NUM_VALUES_PER_ENTRY; j++ {
				v := (ca.data[i] >> (j * 16)) & 0xFFFF
				if v == 0 {
					continue
				}
				if current_logical_index == idx {
					is_negative := ca.get_is_negative(i, j)
					if is_negative {
						fmt.Printf("v: %v\n", v)
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
