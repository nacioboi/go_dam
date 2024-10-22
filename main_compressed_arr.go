package main

// import (
// 	"fmt"
// 	"log"
// 	"math"
// 	"math/rand"
// 	"runtime"
// 	"runtime/debug"
// 	"time"

// 	"github.com/jlabath/bitarray"
// )

// const (
// 	c_COMPRESSED_ARRAY__MAX_ENTRIES_PER_CHECKPOINT = 65535
// )

// type t_checkpoint_info struct {
// 	value_idx uint8
// }

// type t_differences_info struct {
// 	bit_idx_addition uint8
// }

// type I_Signed_Integer interface {
// 	int8 | int16 | int32 | int64
// }

// // TODO: Use a bit array to store the differences instead of a slice of integers.
// // TODO: It's the only way to ensure that the differences are packed as tightly as possible.
// type Compressed_Array[T I_Signed_Integer] struct {
// 	checkpoints       []uint64             // Stores the checkpoint values
// 	differences       *bitarray.BitArray   // Stores the differences between values and their respective checkpoints
// 	checkpoint_map    []t_checkpoint_info  // Maps each checkpoint to its starting index in the differences slice
// 	differences_map   []t_differences_info // Maps each difference to its starting index in the differences slice
// 	num_values        uint16               // Total number of values stored
// 	bitarray_num_bits uint16
// 	rc_MAX_DIFF       T // Runtime Constant
// 	rc_MIN_DIFF       T // Runtime Constant
// }

// // Creates a new CompressedArray
// func New_Compressed_Array[T I_Signed_Integer]() *Compressed_Array[T] {

// 	// Find the minimum for T...
// 	var rc_MIN_DIFF T
// 	last := rc_MIN_DIFF
// 	for {
// 		if rc_MIN_DIFF > last {
// 			break
// 		}
// 		last = rc_MIN_DIFF
// 		rc_MIN_DIFF--
// 	}
// 	rc_MIN_DIFF = last

// 	// Find the maximum for T...
// 	var rc_MAX_DIFF T
// 	last = rc_MAX_DIFF
// 	for {
// 		if rc_MAX_DIFF < last {
// 			break
// 		}
// 		last = rc_MAX_DIFF
// 		rc_MAX_DIFF++
// 	}
// 	rc_MAX_DIFF = last

// 	bitarray_num_bits := uint16(16)
// 	res := &Compressed_Array[T]{
// 		checkpoints:       make([]uint64, 0),
// 		differences:       bitarray.New(int(bitarray_num_bits)),
// 		checkpoint_map:    make([]t_checkpoint_info, 0),
// 		num_values:        0,
// 		bitarray_num_bits: bitarray_num_bits,
// 		rc_MAX_DIFF:       rc_MAX_DIFF,
// 		rc_MIN_DIFF:       rc_MIN_DIFF,
// 	}

// 	return res
// }

// // Appends a value to the compressed array
// func (ca *Compressed_Array[T]) Append(value uint64) {
// 	var diff int64
// 	//var is_new_checkpoint bool

// 	// Decide whether to create a new checkpoint
// 	if ca.num_values == 0 || (ca.num_values%c_COMPRESSED_ARRAY__MAX_ENTRIES_PER_CHECKPOINT) == 0 {
// 		previous_value_idx := 0
// 		if ca.num_values > 0 {
// 			previous_value_idx = int(ca.checkpoint_map[len(ca.checkpoint_map)-1].value_idx)
// 		}
// 		// Create a new checkpoint
// 		ca.checkpoints = append(ca.checkpoints, value)
// 		ca.checkpoint_map = append(ca.checkpoint_map, t_checkpoint_info{
// 			value_idx: uint8(ca.num_values - uint16(previous_value_idx)),
// 		})
// 		diff = 0
// 		//is_new_checkpoint = true
// 	} else {
// 		// Calculate difference from the last checkpoint
// 		checkpoint_idx := len(ca.checkpoints) - 1
// 		base_value := ca.checkpoints[checkpoint_idx]
// 		diff = int64(value) - int64(base_value)

// 		// If difference is out of range, create a new checkpoint
// 		if diff < int64(ca.rc_MIN_DIFF) || diff > int64(ca.rc_MAX_DIFF) {
// 			prev_value_idx := ca.checkpoint_map[len(ca.checkpoint_map)-1].value_idx
// 			ca.checkpoints = append(ca.checkpoints, value)
// 			ca.checkpoint_map = append(ca.checkpoint_map, t_checkpoint_info{
// 				value_idx: uint8(ca.num_values - uint16(prev_value_idx)),
// 			})
// 			diff = 0
// 			//is_new_checkpoint = true
// 		}
// 	}

// 	// Find the appropriate amount of bits
// 	var num_bits uint16 = 32
// 	var mask uint32 = (1 << num_bits) | ((1 << num_bits) - 1)
// 	if diff > 0 {
// 		for {
// 			if diff&int64(mask) != diff {
// 				if num_bits == 32 {
// 					panic("Difference is too large")
// 				}
// 				num_bits++
// 				break
// 			}
// 			num_bits = uint16(math.Floor(math.Log2(float64(diff))))
// 			mask = (1 << num_bits) - 1
// 		}
// 	} else {
// 		num_bits = 0
// 		mask = 0
// 	}
// 	fmt.Printf("num_bits: %d\n", num_bits)

// 	// Prepare the difference for storage
// 	var last_bit_idx uint16
// 	if len(ca.differences_map) == 0 {
// 		last_bit_idx = 0
// 		ca.differences_map = append(ca.differences_map, t_differences_info{
// 			bit_idx_addition: 0,
// 		})
// 	} else {
// 		last_bit_idx += uint16(ca.differences_map[len(ca.differences_map)-1].bit_idx_addition)
// 		ca.differences_map = append(ca.differences_map, t_differences_info{
// 			bit_idx_addition: uint8(last_bit_idx) + uint8(num_bits),
// 		})
// 	}
// 	fmt.Printf("last_bit_idx: %d\n\n", last_bit_idx)

// 	// Reallocate bit array if necessary
// 	if last_bit_idx+num_bits > ca.bitarray_num_bits {
// 		fmt.Printf("Reallocating bit array to %d bits\n", ca.bitarray_num_bits+1024)
// 		new_differences := bitarray.New(int(ca.bitarray_num_bits + 1024))
// 		for i := 0; i < int(ca.bitarray_num_bits); i++ {
// 			if ca.differences.IsSet(i) {
// 				new_differences.Set(i)
// 			} else {
// 				new_differences.Unset(i)
// 			}
// 		}
// 		ca.bitarray_num_bits += 1024
// 		ca.differences = new_differences
// 	}

// 	// Append the difference to the bit array
// 	for i := 0; i < int(num_bits); i++ {
// 		bit := (diff >> uint(i)) & 1
// 		if bit == 1 {
// 			ca.differences.Set(int(last_bit_idx + uint16(i)))
// 		} else {
// 			ca.differences.Unset(int(last_bit_idx + uint16(i)))
// 		}
// 	}

// 	ca.num_values++

// 	// For debugging
// 	// if is_new_checkpoint {
// 	// 	fmt.Printf("Created checkpoint at index %d with value %d\n", ca.num_values-1, value)
// 	// } else {
// 	// 	fmt.Printf("Appended difference %d at index %d\n", diff, ca.num_values-1)
// 	// }
// }

// // Retrieves a value by its index
// func (ca *Compressed_Array[T]) Get(idx uint16) uint64 {
// 	return 0
// }

// // Returns the length of the compressed array
// func (ca *Compressed_Array[T]) Len() uint16 {
// 	return ca.num_values
// }

// // Testing the CompressedArray
// func mainCompressedTest() {
// 	debug.SetGCPercent(-1) // Disable garbage collection for accurate memory measurements

// 	const n = 64
// 	fmt.Println("Starting test")

// 	var start_mem_info runtime.MemStats
// 	var end_mem_info runtime.MemStats

// 	// Proof that bitarray is working correctly
// 	runtime.ReadMemStats(&start_mem_info)
// 	ba := bitarray.New(10240)
// 	for i := 0; i < 10240; i++ {
// 		ba.Set(i)
// 	}
// 	runtime.ReadMemStats(&end_mem_info)
// 	measured_size_of_bitarray := end_mem_info.Alloc - start_mem_info.Alloc
// 	fmt.Printf("Measured size of bitarray:                      %d bytes\n", measured_size_of_bitarray)
// 	start_mem_info = end_mem_info

// 	runtime.ReadMemStats(&start_mem_info)

// 	r := rand.New(rand.NewSource(time.Now().UnixNano()))
// 	raw_data := make([]uint64, n)
// 	base := uint64(1_000_000)

// 	// Generate random data
// 	for i := 0; i < n; i++ {
// 		offset := int64(r.Intn(8)) // Random number between -20000 and 20000
// 		if i == 0 {
// 			offset = 0
// 		}
// 		raw_data[i] = base + uint64(offset)
// 		//fmt.Printf("Generated value %d for index %d\n", raw_data[i], i)
// 	}

// 	runtime.ReadMemStats(&end_mem_info)
// 	measured_size_of_raw_data := end_mem_info.Alloc - start_mem_info.Alloc
// 	start_mem_info = end_mem_info

// 	// Create and populate the compressed array
// 	ca := New_Compressed_Array[int16]()
// 	for _, v := range raw_data {
// 		ca.Append(v)
// 	}

// 	runtime.ReadMemStats(&end_mem_info)
// 	measured_size_of_compressed_array := end_mem_info.Alloc - start_mem_info.Alloc
// 	fmt.Printf("Measured size of raw data:                       %d bytes\n", measured_size_of_raw_data)
// 	fmt.Printf("Measured size of compressed array:               %d bytes\n", measured_size_of_compressed_array)
// 	fmt.Printf(
// 		"Measured Compression ratio (lower is better):    %.2f%%\n",
// 		float64(measured_size_of_compressed_array)/float64(measured_size_of_raw_data)*100,
// 	)
// 	fmt.Printf("Number of checkpoints:                           %d\n", len(ca.checkpoints))
// 	fmt.Printf("Number of bits used for differences:             %d\n", ca.bitarray_num_bits)

// 	// Verify that all values can be retrieved correctly
// 	for i := uint16(0); i < uint16(ca.Len()); i++ {
// 		res := ca.Get(i)
// 		if res != raw_data[i] {
// 			log.Fatalf("Mismatch at index %d: Expected %d but got %d\n", i, raw_data[i], res)
// 		} else {
// 			//fmt.Printf("Value at index %d matches: %d\n", i, res)
// 		}
// 	}

// 	fmt.Println("All values retrieved and matched successfully!")
// }

// // package main

// // import (
// // 	"fmt"
// // )

// // // Constants for bit masks and sizes
// // const (
// // 	SIXTY_FOURTH_BIT_BITMASK       uint64 = 1 << 63
// // 	ALL_BITS_BUT_64_BITMASK        uint64 = ^SIXTY_FOURTH_BIT_BITMASK
// // 	MAX_NUM_ENTRIES_PER_CHECKPOINT        = 16
// // 	DATA_ENTRIES_PER_UINT64               = 4 // Since we're packing 4 * 16-bit entries
// // )

// // // The compressed array stores information in two types of entries:
// // //
// // // - Checkpoint: A uint64 with the highest bit set to 1, containing the base value.
// // // - Data: uint64 entries containing up to 4 packed 16-bit differences.
// // type compressed_array struct {
// // 	data       []uint64 // Stores both checkpoints and data entries
// // 	num_values int      // Total number of values stored
// // 	counts     []int    // Number of valid differences in each data entry
// // }

// // func new_compressed_array() *compressed_array {
// // 	return &compressed_array{
// // 		data:   make([]uint64, 0),
// // 		counts: make([]int, 0),
// // 	}
// // }

// // // Checks if the entry at idx is a checkpoint
// // func (ca *compressed_array) is_checkpoint(idx int) bool {
// // 	return ca.data[idx]&SIXTY_FOURTH_BIT_BITMASK != 0
// // }

// // // Retrieves the checkpoint value without the highest bit
// // func (ca *compressed_array) get_checkpoint_value(idx int) uint64 {
// // 	return ca.data[idx] & ALL_BITS_BUT_64_BITMASK
// // }

// // // Appends a value to the compressed array
// // func (ca *compressed_array) Append(value uint64) {
// // 	// If there are no checkpoints, create one
// // 	if len(ca.data) == 0 {
// // 		ca.create_checkpoint(value)
// // 	}

// // 	// Get the last checkpoint index
// // 	last_checkpoint_idx := ca.get_last_checkpoint_idx()

// // 	// Get the current base value from the last checkpoint
// // 	base_value := ca.get_checkpoint_value(last_checkpoint_idx)

// // 	// Calculate the difference
// // 	diff := int64(value) - int64(base_value)

// // 	// Check if the difference fits in 16 bits
// // 	if diff < -32768 || diff > 32767 {
// // 		// Difference doesn't fit, create a new checkpoint
// // 		ca.create_checkpoint(value)
// // 		last_checkpoint_idx = len(ca.data) - 1
// // 		base_value = ca.get_checkpoint_value(last_checkpoint_idx)
// // 		diff = 0
// // 	}

// // 	// Append the difference
// // 	ca.append_difference(int16(diff))

// // 	ca.num_values++
// // }

// // // Creates a new checkpoint with the given base value
// // func (ca *compressed_array) create_checkpoint(value uint64) {
// // 	// Set the highest bit to 1 to mark as checkpoint
// // 	checkpoint_value := value | SIXTY_FOURTH_BIT_BITMASK
// // 	ca.data = append(ca.data, checkpoint_value)
// // }

// // // Gets the index of the last checkpoint
// // func (ca *compressed_array) get_last_checkpoint_idx() int {
// // 	for i := len(ca.data) - 1; i >= 0; i-- {
// // 		if ca.is_checkpoint(i) {
// // 			return i
// // 		}
// // 	}
// // 	panic("No checkpoint found")
// // }

// // // Appends a 16-bit difference, packing it into uint64 entries
// // func (ca *compressed_array) append_difference(diff int16) {
// // 	// Check if we need to start a new data entry
// // 	if ca.needs_new_data_entry() {
// // 		ca.data = append(ca.data, 0) // Initialize new data entry
// // 		ca.counts = append(ca.counts, 0)
// // 	}

// // 	// Get the last data entry index
// // 	last_data_idx := len(ca.data) - 1

// // 	// Get the position within the 64-bit data entry (0 to 3)
// // 	position := ca.counts[len(ca.counts)-1]

// // 	// Pack the 16-bit difference into the correct position
// // 	shift := uint(position * 16)
// // 	ca.data[last_data_idx] |= (uint64(uint16(diff)) << shift)

// // 	// Increment the count for this data entry
// // 	ca.counts[len(ca.counts)-1]++
// // }

// // // Determines if a new data entry is needed
// // func (ca *compressed_array) needs_new_data_entry() bool {
// // 	if len(ca.data) == 0 {
// // 		return true
// // 	}

// // 	// Get the last data entry index
// // 	last_idx := len(ca.data) - 1

// // 	// If the last entry is a checkpoint, we need a new data entry
// // 	if ca.is_checkpoint(last_idx) {
// // 		return true
// // 	}

// // 	// Check if the last data entry is full
// // 	if ca.counts[len(ca.counts)-1] >= DATA_ENTRIES_PER_UINT64 {
// // 		return true
// // 	}

// // 	return false
// // }

// // // Retrieves the value at the specified index
// // func (ca *compressed_array) Get(idx int) uint64 {
// // 	if idx < 0 || idx >= ca.num_values {
// // 		panic("Index out of range")
// // 	}

// // 	// Variables to keep track of the current position
// // 	value_count := 0
// // 	var base_value uint64
// // 	count_idx := 0 // Index into counts slice
// // 	for i := 0; i < len(ca.data); i++ {
// // 		if ca.is_checkpoint(i) {
// // 			// Update the base value
// // 			base_value = ca.get_checkpoint_value(i)
// // 			count_idx = 0 // Reset count index for new checkpoint
// // 		} else {
// // 			// Unpack differences from the data entry
// // 			count := ca.counts[count_idx]
// // 			count_idx++

// // 			diffs := ca.unpack_differences(ca.data[i], count)
// // 			for _, diff := range diffs {
// // 				if value_count == idx {
// // 					// Reconstruct the original value
// // 					return base_value + uint64(diff)
// // 				}
// // 				value_count++
// // 				if value_count >= ca.num_values {
// // 					break
// // 				}
// // 			}
// // 		}
// // 	}
// // 	panic("Index not found")
// // }

// // // Unpacks the specified number of differences from a data entry
// // func (ca *compressed_array) unpack_differences(data uint64, count int) []int16 {
// // 	diffs := make([]int16, 0, count)
// // 	for i := 0; i < count; i++ {
// // 		shift := uint(i * 16)
// // 		diff := int16((data >> shift) & 0xFFFF)
// // 		diffs = append(diffs, diff)
// // 	}
// // 	return diffs
// // }

// // func

// // func main_compressed_test() {
// // 	ca := new_compressed_array()
// // 	values := []uint64{
// // 		100000, 100005, 100010, 100015,
// // 		99990, 99995, 100020, 100025,
// // 		200000, 200005, 200010, 200015,
// // 	}

// // 	for _, v := range values {
// // 		ca.Append(v)
// // 	}

// // 	for i := 0; i < ca.num_values; i++ {
// // 		v := ca.Get(i)
// // 		fmt.Printf("Value at index %d: %d\n", i, v)
// // 	}

// // 	fmt.Printf("len(values) = %d\n", len(values))
// // 	fmt.Printf("len(ca.data) = %d\n", len(ca.data))
// // }
