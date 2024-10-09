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

// We must set an upper limit since having only one checkpoint for the entire array would be bad for the compression ratio.
const MAX_NUM_ENTRIES_PER_CHECKPOINT = 16

// The compressed array stores information in two types of entries:
//
// - Data: Contains multiple (up to 16) values. They are the difference between said value and the checkpoint that precedes.
//
// - Checkpoint: Contains an optimized base value.
type compressed_array struct {
	data        []uint64
	split_infos []uint8
}

func new_compressed_array() *compressed_array {
	return &compressed_array{
		data: make([]uint64, 0),
	}
}

func (ca *compressed_array) is_checkpoint(idx uint64) bool {
	return ca.data[idx]&SIXTY_FORTH_BIT_BITMASK == SIXTY_FORTH_BIT_BITMASK
}

func (ca *compressed_array) get_checkpoint_value(idx uint64) uint64 {
	return ca.data[idx] & ALL_BITS_BUT_64_BITMASK
}

func (ca *compressed_array) can_checkpoint_fit_more_values(checkpoint_idx uint64) bool {
	num_attached_values := uint64(0)

	for i := checkpoint_idx + 1; i < uint64(len(ca.data)); i++ {
		if ca.is_checkpoint(i) {
			break
		}
		num_attached_values++
	}

	if num_attached_values >= MAX_NUM_ENTRIES_PER_CHECKPOINT {
		return false
	}
	return true
}

func (ca *compressed_array) get_previous_checkpoint_idx(starting_idx uint64) uint64 {
	for i := uint64(starting_idx); i >= 0; i-- {
		if ca.is_checkpoint(i) {
			return i
		}
	}
	panic("something went horribly wrong")
}

// mathematically calculate the optimal checkpoint value using the values of all the data attached to it.
// The goal is to squeeze as much information as possible in the 64 bits of the checkpoint.
// So a checkpoint value that is un-optimized could be really bad for the compression ratio.
func (ca *compressed_array) find_optimal_checkpoint_value(checkpoint_idx uint64) uint64 {
	panic("not implemented")
}

func (ca *compressed_array) create_checkpoint() {
	if len(ca.data) == 0 {
		ca.data = append(ca.data, SIXTY_FORTH_BIT_BITMASK)
		return
	}

	panic("not implemented")
}

func (ca *compressed_array) Append(value uint64) {

}
