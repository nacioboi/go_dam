package dam

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
	c__BITMASK__SIGNS_INFOS uint64 = 0b0111000000000000000000000000000000000000000000000000000000000000
	////////////////////////////////// 432109876543210
	///////////////////////////////////////////////// 5432109876543210
	///////////////////////////////////////////////////////////////// 5432109876543210
	///////////////////////////////////////////////////////////////////////////////// 5432109876543210
	c__BITMASK__MULTIPLIERS uint64 = 0b0000111111111111000000000000000000000000000000000000000000000000
	c__BITMASK__OUR_BITS    uint64 = 0xFFFF

	// The amount of right shift required to bring the extras to the start of the 64-bit integer.
	c__SHIFT_WIDTH__EXTRAS uint64 = 48
	// The amount of right shift, on-top of right shift by `c__SHIFT_WIDTH__EXTRAS`,
	// required to bring the sign infos to the start of the 64-bit integer.
	c__NUM_BITS_FOR_MULTIPLIERS uint64 = 12

	c__NUM_BITS_PER_VALUE       uint64 = 16
	c__NUM_BITS_PER_MULTIPLIER  uint64 = 4
	c__MAX_MULTIPLIER           uint64 = 16 // 1-based instead of zero-based
	c__MAX_NUM_VALUES_PER_ENTRY uint64 = 3
)

// The compressed array stores information in two types of entries:
//
// - Data: Contains multiple (up to 16) values. They are the difference between said value and the checkpoint that precedes.
//
// - Checkpoint: Contains an optimized base value.
type Compressed_Array struct {
	D []uint64
}

func New_Compressed_Array() *Compressed_Array {
	return &Compressed_Array{
		D: make([]uint64, 0),
	}
}

func (ca *Compressed_Array) is_checkpoint(idx uint64) bool {
	v := ca.D[idx]
	return v&c__BITMASK__64th_B == c__BITMASK__64th_B
}

func (ca *Compressed_Array) get_checkpoint_value(idx uint64) uint64 {
	if !ca.is_checkpoint(idx) {
		panic("not a checkpoint")
	}
	return ca.D[idx] & c__BITMASK__BUT_64
}

func (ca *Compressed_Array) get_num_attached_values(checkpoint_idx uint64) uint64 {
	num_attached_values := uint64(0)

	for i := checkpoint_idx + 1; i < uint64(len(ca.D)); i++ {
		if ca.is_checkpoint(i) {
			break
		}
		num_attached_values++
	}

	return num_attached_values
}

func (ca *Compressed_Array) needs_new_checkpoint(checkpoint_idx uint64, value uint64) bool {
	diff, _, _ := ca.calculate_difference(ca.get_checkpoint_value(checkpoint_idx), value)

	if diff >= (c__BITMASK__OUR_BITS * c__MAX_MULTIPLIER) {
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
	checkpoint := c__BITMASK__64th_B | recommended_value
	ca.D = append(ca.D, checkpoint)
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
	if diff >= int64(c__BITMASK__OUR_BITS*c__MAX_MULTIPLIER) {
		return 0, true, is_negative
	}
	return uint64(diff), false, is_negative
}

func (ca *Compressed_Array) apply_sign_infos(idx uint64, j uint64, is_negative bool) {
	sign_infos := ca.D[idx] & c__BITMASK__SIGNS_INFOS
	shifted := sign_infos >> (c__SHIFT_WIDTH__EXTRAS + (c__NUM_BITS_FOR_MULTIPLIERS - 1))
	m := uint64(1) << (3 - j)
	if is_negative {
		shifted |= m
	} else {
		shifted &= (^m & 0b1111)
	}
	mask := shifted << (c__SHIFT_WIDTH__EXTRAS + (c__NUM_BITS_FOR_MULTIPLIERS - 1))
	ca.D[idx] |= mask
}

func (ca *Compressed_Array) get_is_negative(idx, j uint64) bool {
	sign_infos := ca.D[idx] & c__BITMASK__SIGNS_INFOS
	shifted := sign_infos >> (c__SHIFT_WIDTH__EXTRAS + (c__NUM_BITS_FOR_MULTIPLIERS - 1))
	m := uint64(1) << (3 - j)
	return shifted&m != 0
}

func (ca *Compressed_Array) set_multiplier(idx uint64, j uint64, multiplier uint64) {
	multiplier_infos := ca.D[idx] & c__BITMASK__MULTIPLIERS

	m1 := (multiplier >> 0) & 0b1
	m2 := (multiplier >> 1) & 0b1
	m3 := (multiplier >> 2) & 0b1
	m4 := (multiplier >> 3) & 0b1

	shifted := multiplier_infos >> (c__SHIFT_WIDTH__EXTRAS)
	shifted &= 0b111111111111

	base_shifter := c__NUM_BITS_FOR_MULTIPLIERS - ((1 + j) * 4)
	m1_shifter := base_shifter + 3
	m2_shifter := base_shifter + 2
	m3_shifter := base_shifter + 1
	m4_shifter := base_shifter + 0

	if m1 == 1 { // m1
		shifted |= m1 << m1_shifter
	} else {
		shifted &= ^(m1 << m1_shifter)
	}

	if m2 == 1 { // m2
		shifted |= m2 << m2_shifter
	} else {
		shifted &= ^(m2 << m2_shifter)
	}

	if m3 == 1 { // m3
		shifted |= m3 << m3_shifter
	} else {
		shifted &= ^(m3 << m3_shifter)
	}

	if m4 == 1 { // m4
		shifted |= m4 << m4_shifter
	} else {
		shifted &= ^(m4 << m4_shifter)
	}

	mask := shifted << (c__SHIFT_WIDTH__EXTRAS)
	ca.D[idx] |= mask
}

func (ca *Compressed_Array) get_multiplier(idx uint64, j uint64) uint64 {
	multiplier_infos := ca.D[idx] & c__BITMASK__MULTIPLIERS

	shifted := multiplier_infos >> (c__SHIFT_WIDTH__EXTRAS)
	shifted &= 0b111111111111

	base_shifter := c__NUM_BITS_FOR_MULTIPLIERS - ((1 + j) * 4)
	m1_shifter := base_shifter + 3
	m2_shifter := base_shifter + 2
	m3_shifter := base_shifter + 1
	m4_shifter := base_shifter + 0

	var (
		x1 uint64
		x2 uint64
		x3 uint64
		x4 uint64
	)

	x1 = (shifted >> m1_shifter) & 0b1
	x2 = (shifted >> m2_shifter) & 0b1
	x3 = (shifted >> m3_shifter) & 0b1
	x4 = (shifted >> m4_shifter) & 0b1

	x := (x1 << 0) | (x2 << 1) | (x3 << 2) | (x4 << 3)

	return x
}

func (ca *Compressed_Array) find_appropriate_multiplier(diff uint64) (uint64, uint64) {
	// Calculate the multiplier
	appropriate_multiplier := uint64(1)
	did_find_multiplier := false
	for appropriate_multiplier < c__MAX_MULTIPLIER {
		if (diff / appropriate_multiplier) < c__BITMASK__OUR_BITS {
			if diff%appropriate_multiplier == 0 {
				did_find_multiplier = true
				break
			}
		}
		appropriate_multiplier++
	}
	if !did_find_multiplier {
		panic("could not find appropriate multiplier")
	}
	if appropriate_multiplier != 1 {
	}
	return appropriate_multiplier, diff / appropriate_multiplier
}

func (ca *Compressed_Array) Append(value uint64) {
	//fmt.Printf("data: %v\n", ca.data)

	// Get the previous checkpoint index...
	// TODO: This is just an example since we don't always want to get the last checkpoint.
	// E.g. If the diff is too large, we want to search all checkpoint and check if one exists that fits
	//        before adding a new checkpoint.
	if len(ca.D) == 0 {
		ca.create_checkpoint(value)
		return
	}

	// Get some values...
	checkpoint_idx := ca.get_previous_checkpoint_idx(uint64(len(ca.D) - 1))
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

	for i := checkpoint_idx + 1; i < uint64(len(ca.D)); i++ {
		if ca.is_checkpoint(i) {
			break // Stop searching since we have reached the next checkpoint.
		}
		for j := uint64(0); j < c__MAX_NUM_VALUES_PER_ENTRY; j++ {
			if ca.D[i]&(c__BITMASK__OUR_BITS<<(j*c__NUM_BITS_PER_VALUE)) != 0 {
				continue
			}
			diff, does_need_new_checkpoint, is_negative := ca.calculate_difference(checkpoint_value, value)
			if !does_need_new_checkpoint {
				appropriate_multiplier, new_diff := ca.find_appropriate_multiplier(diff)
				ca.D[i] |= uint64(new_diff) << (j * c__NUM_BITS_PER_VALUE)
				ca.apply_sign_infos(i, j, is_negative)
				ca.set_multiplier(i, j, appropriate_multiplier)
			} else {
				// We need to create a new checkpoint.
				ca.create_checkpoint(value)
			}
			return
		}
	}

	diff, does_need_new_checkpoint, is_negative := ca.calculate_difference(checkpoint_value, value)
	if !does_need_new_checkpoint {
		appropriate_multiplier, new_diff := ca.find_appropriate_multiplier(diff)
		ca.D = append(ca.D, uint64(new_diff))
		ca.apply_sign_infos(uint64(len(ca.D)-1), 0, is_negative)
		ca.set_multiplier(uint64(len(ca.D)-1), 0, appropriate_multiplier)
	} else {
		// We need to create a new checkpoint.
		ca.create_checkpoint(value)
	}
}

func (ca *Compressed_Array) Get(idx uint64) uint64 {
	var (
		current_checkpoint_value uint64
		current_logical_index    uint64
	)

	for i := uint64(0); i < uint64(len(ca.D)); i++ {
		if ca.is_checkpoint(i) {
			current_checkpoint_value = ca.get_checkpoint_value(i)
			if current_logical_index == idx {
				return current_checkpoint_value
			}
			current_logical_index++
		} else {
			for j := uint64(0); j < c__MAX_NUM_VALUES_PER_ENTRY; j++ {
				v := (ca.D[i] >> (j * c__NUM_BITS_PER_VALUE)) & c__BITMASK__OUR_BITS
				if v == 0 {
					continue
				}
				if current_logical_index == idx {
					multiplier := ca.get_multiplier(i, j)
					is_negative := ca.get_is_negative(i, j)
					if is_negative {
						return current_checkpoint_value - (v * multiplier)
					} else {
						return current_checkpoint_value + (v * multiplier)
					}
				}
				current_logical_index++
			}
		}
	}

	panic("Index out of bounds")
}
