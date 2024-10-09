package main

import (
	"fmt"
	"math/rand"
	"time"
	"unsafe"
)

// Constants for bit masks and sizes
const (
	SIXTY_FOURTH_BIT_BITMASK       uint64 = 1 << 63
	ALL_BITS_BUT_64_BITMASK        uint64 = ^SIXTY_FOURTH_BIT_BITMASK
	MAX_NUM_ENTRIES_PER_CHECKPOINT        = 16
	DATA_ENTRIES_PER_UINT64               = 4 // Since we're packing 4 * 16-bit entries
)

// The compressed array stores information in two types of entries:
//
// - Checkpoint: A uint64 with the highest bit set to 1, containing the base value.
// - Data: uint64 entries containing up to 4 packed 16-bit differences.
type compressed_array struct {
	data       []uint64         // Stores both checkpoints and data entries
	num_values int              // Total number of values stored
	counts     map[uint64]uint8 // Number of valid differences in each data entry (values from 1 to 4)
}

func new_compressed_array() *compressed_array {
	return &compressed_array{
		data:   make([]uint64, 0),
		counts: make(map[uint64]uint8),
	}
}

func (ca *compressed_array) is_checkpoint(idx uint64) bool {
	return ca.data[idx]&SIXTY_FOURTH_BIT_BITMASK != 0
}

func (ca *compressed_array) get_checkpoint_value(idx uint64) uint64 {
	return ca.data[idx] & ALL_BITS_BUT_64_BITMASK
}

func (ca *compressed_array) get_last_checkpoint_idx() uint64 {
	for i := uint64(len(ca.data) - 1); i >= 0; i-- {
		if ca.is_checkpoint(i) {
			return i
		}
	}
	panic("No checkpoint found")
}

func (ca *compressed_array) find_last_data_idx_of_checkpoint(checkpoint_idx uint64) uint64 {
	for i := checkpoint_idx + 1; i < uint64(len(ca.data)); i++ {
		if ca.is_checkpoint(i) {
			return i - 1
		}
	}
	return uint64(len(ca.data) - 1)
}

func (ca *compressed_array) checkpoint_can_fit_diff(checkpoint_idx uint64, value uint64, base uint64) (uint64, bool) {
	var diff int64 = int64(value) - int64(base)

	// Check if the difference fits in 16 bits
	if diff < -32768 || diff > 32767 {
		// Check if there are any other checkpoints that can fit the difference

		// First, search backwards
		if checkpoint_idx > 0 {
			for i := checkpoint_idx - 1; i >= 0; i-- {
				if ca.is_checkpoint(i) {
					base_value := ca.get_checkpoint_value(i)
					diff = int64(value) - int64(base_value)
					if diff >= -32768 && diff <= 32767 {
						return i, true
					}
				}
			}
		}

		// Next, search forwards
		if checkpoint_idx+1 < uint64(len(ca.data)) {
			for i := checkpoint_idx + 1; i < uint64(len(ca.data)); i++ {
				if ca.is_checkpoint(i) {
					base_value := ca.get_checkpoint_value(i)
					diff = int64(value) - int64(base_value)
					if diff >= -32768 && diff <= 32767 {
						return i, true
					}
				}
			}
		}

		// Failed to find a checkpoint that can fit the difference
		return 0, false
	}

	return checkpoint_idx, true
}

// Determines if a new data entry is needed
func (ca *compressed_array) needs_new_data_entry(checkpoint_idx uint64, last_entry_idx uint64) bool {
	if len(ca.data) == 0 {
		return true
	}

	// If the last entry is a checkpoint, we need a new data entry
	if ca.is_checkpoint(last_entry_idx) {
		return true
	}

	// Check if the last data entry is full
	if ca.counts[checkpoint_idx] >= DATA_ENTRIES_PER_UINT64 {
		return true
	}

	return false
}

func (ca *compressed_array) unpack_differences(data uint64, count int) []int16 {
	diffs := make([]int16, 0, count)
	for i := 0; i < count; i++ {
		shift := uint(i * 16)
		diff := int16((data >> shift) & 0xFFFF)
		diffs = append(diffs, diff)
	}
	return diffs
}

func (ca *compressed_array) create_checkpoint(value uint64) {
	// Set the highest bit to 1 to mark as checkpoint
	checkpoint_value := value | SIXTY_FOURTH_BIT_BITMASK
	ca.data = append(ca.data, checkpoint_value)
}

func (ca *compressed_array) append_difference(checkpoint_idx uint64, diff int16) {
	last_data_idx := ca.find_last_data_idx_of_checkpoint(checkpoint_idx)

	// Check if we need to start a new data entry
	if ca.needs_new_data_entry(checkpoint_idx, last_data_idx) {
		ca.data = append(ca.data, 0) // Initialize new data entry
		ca.counts[checkpoint_idx] = 0
	}

	// Get the position within the 64-bit data entry (0 to 3)
	position := ca.counts[checkpoint_idx]

	// Pack the 16-bit difference into the correct position
	shift := uint(position * 16 % 48)
	ca.data[last_data_idx] |= (uint64(uint16(diff)) << shift)

	// Increment the count for this data entry
	ca.counts[checkpoint_idx]++
}

// Appends a value to the compressed array
func (ca *compressed_array) Append(value uint64) {
	// If there are no checkpoints, create one
	if len(ca.data) == 0 {
		ca.create_checkpoint(value)
	}

	// Get the last checkpoint index
	last_checkpoint_idx := ca.get_last_checkpoint_idx()
	if !ca.is_checkpoint(last_checkpoint_idx) {
		panic("oopsies")
	}

	// Get the current base value from the last checkpoint
	base_value := ca.get_checkpoint_value(last_checkpoint_idx)

	// Calculate the difference
	idx, ok := ca.checkpoint_can_fit_diff(last_checkpoint_idx, value, base_value)
	if !ok {
		// Difference doesn't fit, create a new checkpoint
		ca.create_checkpoint(value)
	} else {

		base_value = ca.get_checkpoint_value(idx)
		diff := int64(value) - int64(base_value)

		// Append the difference
		ca.append_difference(idx, int16(diff))
	}

	ca.num_values++
}

// Retrieves the value at the specified index
func (ca *compressed_array) Get(idx int) uint64 {
	if idx < 0 || idx >= ca.num_values {
		panic("Index out of range")
	}

	// Keep track of the current number of values processed
	value_count := 0
	var base_value uint64

	// Iterate through the data to find the checkpoint and unpack the differences
	for i := uint64(0); i < uint64(len(ca.data)); i++ {
		if ca.is_checkpoint(i) {
			// Update the base value when we encounter a checkpoint
			base_value = ca.get_checkpoint_value(i)

			return base_value
		} else {
			// Unpack the differences from the data entry
			count := int(ca.counts[i])
			diffs := ca.unpack_differences(ca.data[i], count)

			// Iterate over each difference in this entry
			for _, diff := range diffs {
				if value_count == idx {
					// Found the value we're looking for
					return base_value + uint64(diff)
				}
				value_count++
			}

			return base_value + uint64(diffs[idx])
		}
	}

	// If we reach here, something went wrong
	panic("Index not found")
}

func main_compressed_test() {
	// Seed the random number generator
	rand.Seed(time.Now().UnixNano())

	// Generate a large amount of random data
	numValues := 4 //_000 // One million values
	rawData := make([]uint64, numValues)

	// Generate random data with some locality to ensure differences fit in 16 bits
	// For example, generate values around a moving base
	base := uint64(1_000_000)
	for i := 0; i < numValues; i++ {
		// Randomly change the base value occasionally
		// if rand.Intn(1000) == 0 {
		// 	base = uint64(rand.Int63n(1_000_000_000))
		// }

		// Generate a value within +/- 30000 of the base
		offset := int64(rand.Intn(60001)) - 30000 // Range: -30000 to +30000
		value := base + uint64(offset)
		rawData[i] = value
	}

	// Create a compressed array and append the values
	ca := new_compressed_array()
	for _, v := range rawData {
		ca.Append(v)
	}

	// Benchmark the size of raw data (in bytes)
	sizeOfRawData := numValues * int(unsafe.Sizeof(uint64(0))) // Each uint64 is 8 bytes
	fmt.Printf("Size of raw data: %d bytes\n", sizeOfRawData)

	// Calculate the size of the compressed array data
	sizeOfDataSlice := len(ca.data) * int(unsafe.Sizeof(uint64(0)))    // Each uint64 is 8 bytes
	sizeOfCountsSlice := len(ca.counts) * int(unsafe.Sizeof(uint8(0))) // Each uint8 is 1 byte
	sizeOfCompressedArray := sizeOfDataSlice + sizeOfCountsSlice
	fmt.Printf("Calculated size of compressed_array data: %d bytes\n", sizeOfCompressedArray)

	// Calculate compression ratio
	compressionRatio := float64(sizeOfCompressedArray) / float64(sizeOfRawData) * 100
	fmt.Printf("Compression ratio: %.2f%%\n", compressionRatio)

	// Optionally, verify some values
	for i := 0; i < 5; i++ {
		idx := rand.Intn(numValues)
		originalValue := rawData[idx]
		compressedValue := ca.Get(idx)
		if originalValue != compressedValue {
			fmt.Printf("Mismatch at index %d: original=%d, compressed=%d\n", idx, originalValue, compressedValue)
		} else {
			fmt.Printf("Value at index %d matches: %d\n", idx, compressedValue)
		}
	}
}

// package main

// import (
// 	"fmt"
// )

// // Constants for bit masks and sizes
// const (
// 	SIXTY_FOURTH_BIT_BITMASK       uint64 = 1 << 63
// 	ALL_BITS_BUT_64_BITMASK        uint64 = ^SIXTY_FOURTH_BIT_BITMASK
// 	MAX_NUM_ENTRIES_PER_CHECKPOINT        = 16
// 	DATA_ENTRIES_PER_UINT64               = 4 // Since we're packing 4 * 16-bit entries
// )

// // The compressed array stores information in two types of entries:
// //
// // - Checkpoint: A uint64 with the highest bit set to 1, containing the base value.
// // - Data: uint64 entries containing up to 4 packed 16-bit differences.
// type compressed_array struct {
// 	data       []uint64 // Stores both checkpoints and data entries
// 	num_values int      // Total number of values stored
// 	counts     []int    // Number of valid differences in each data entry
// }

// func new_compressed_array() *compressed_array {
// 	return &compressed_array{
// 		data:   make([]uint64, 0),
// 		counts: make([]int, 0),
// 	}
// }

// // Checks if the entry at idx is a checkpoint
// func (ca *compressed_array) is_checkpoint(idx int) bool {
// 	return ca.data[idx]&SIXTY_FOURTH_BIT_BITMASK != 0
// }

// // Retrieves the checkpoint value without the highest bit
// func (ca *compressed_array) get_checkpoint_value(idx int) uint64 {
// 	return ca.data[idx] & ALL_BITS_BUT_64_BITMASK
// }

// // Appends a value to the compressed array
// func (ca *compressed_array) Append(value uint64) {
// 	// If there are no checkpoints, create one
// 	if len(ca.data) == 0 {
// 		ca.create_checkpoint(value)
// 	}

// 	// Get the last checkpoint index
// 	last_checkpoint_idx := ca.get_last_checkpoint_idx()

// 	// Get the current base value from the last checkpoint
// 	base_value := ca.get_checkpoint_value(last_checkpoint_idx)

// 	// Calculate the difference
// 	diff := int64(value) - int64(base_value)

// 	// Check if the difference fits in 16 bits
// 	if diff < -32768 || diff > 32767 {
// 		// Difference doesn't fit, create a new checkpoint
// 		ca.create_checkpoint(value)
// 		last_checkpoint_idx = len(ca.data) - 1
// 		base_value = ca.get_checkpoint_value(last_checkpoint_idx)
// 		diff = 0
// 	}

// 	// Append the difference
// 	ca.append_difference(int16(diff))

// 	ca.num_values++
// }

// // Creates a new checkpoint with the given base value
// func (ca *compressed_array) create_checkpoint(value uint64) {
// 	// Set the highest bit to 1 to mark as checkpoint
// 	checkpoint_value := value | SIXTY_FOURTH_BIT_BITMASK
// 	ca.data = append(ca.data, checkpoint_value)
// }

// // Gets the index of the last checkpoint
// func (ca *compressed_array) get_last_checkpoint_idx() int {
// 	for i := len(ca.data) - 1; i >= 0; i-- {
// 		if ca.is_checkpoint(i) {
// 			return i
// 		}
// 	}
// 	panic("No checkpoint found")
// }

// // Appends a 16-bit difference, packing it into uint64 entries
// func (ca *compressed_array) append_difference(diff int16) {
// 	// Check if we need to start a new data entry
// 	if ca.needs_new_data_entry() {
// 		ca.data = append(ca.data, 0) // Initialize new data entry
// 		ca.counts = append(ca.counts, 0)
// 	}

// 	// Get the last data entry index
// 	last_data_idx := len(ca.data) - 1

// 	// Get the position within the 64-bit data entry (0 to 3)
// 	position := ca.counts[len(ca.counts)-1]

// 	// Pack the 16-bit difference into the correct position
// 	shift := uint(position * 16)
// 	ca.data[last_data_idx] |= (uint64(uint16(diff)) << shift)

// 	// Increment the count for this data entry
// 	ca.counts[len(ca.counts)-1]++
// }

// // Determines if a new data entry is needed
// func (ca *compressed_array) needs_new_data_entry() bool {
// 	if len(ca.data) == 0 {
// 		return true
// 	}

// 	// Get the last data entry index
// 	last_idx := len(ca.data) - 1

// 	// If the last entry is a checkpoint, we need a new data entry
// 	if ca.is_checkpoint(last_idx) {
// 		return true
// 	}

// 	// Check if the last data entry is full
// 	if ca.counts[len(ca.counts)-1] >= DATA_ENTRIES_PER_UINT64 {
// 		return true
// 	}

// 	return false
// }

// // Retrieves the value at the specified index
// func (ca *compressed_array) Get(idx int) uint64 {
// 	if idx < 0 || idx >= ca.num_values {
// 		panic("Index out of range")
// 	}

// 	// Variables to keep track of the current position
// 	value_count := 0
// 	var base_value uint64
// 	count_idx := 0 // Index into counts slice
// 	for i := 0; i < len(ca.data); i++ {
// 		if ca.is_checkpoint(i) {
// 			// Update the base value
// 			base_value = ca.get_checkpoint_value(i)
// 			count_idx = 0 // Reset count index for new checkpoint
// 		} else {
// 			// Unpack differences from the data entry
// 			count := ca.counts[count_idx]
// 			count_idx++

// 			diffs := ca.unpack_differences(ca.data[i], count)
// 			for _, diff := range diffs {
// 				if value_count == idx {
// 					// Reconstruct the original value
// 					return base_value + uint64(diff)
// 				}
// 				value_count++
// 				if value_count >= ca.num_values {
// 					break
// 				}
// 			}
// 		}
// 	}
// 	panic("Index not found")
// }

// // Unpacks the specified number of differences from a data entry
// func (ca *compressed_array) unpack_differences(data uint64, count int) []int16 {
// 	diffs := make([]int16, 0, count)
// 	for i := 0; i < count; i++ {
// 		shift := uint(i * 16)
// 		diff := int16((data >> shift) & 0xFFFF)
// 		diffs = append(diffs, diff)
// 	}
// 	return diffs
// }

// func

// func main_compressed_test() {
// 	ca := new_compressed_array()
// 	values := []uint64{
// 		100000, 100005, 100010, 100015,
// 		99990, 99995, 100020, 100025,
// 		200000, 200005, 200010, 200015,
// 	}

// 	for _, v := range values {
// 		ca.Append(v)
// 	}

// 	for i := 0; i < ca.num_values; i++ {
// 		v := ca.Get(i)
// 		fmt.Printf("Value at index %d: %d\n", i, v)
// 	}

// 	fmt.Printf("len(values) = %d\n", len(values))
// 	fmt.Printf("len(ca.data) = %d\n", len(ca.data))
// }
