/*/
 ** This software is covered by the MIT License.
 ** See: `./LICENSE`.
/*/

package main

import (
	"runtime"

	"github.com/nacioboi/go_dam/dam"
)

func main() {
	runtime.GC()

	const n = 10240

	var start_mem_info runtime.MemStats
	var end_mem_info runtime.MemStats

	// Benchmark normal - LINEAR
	runtime.ReadMemStats(&start_mem_info)
	na := make([]uint64, 0)
	for i := uint64(0); i < n; i++ {
		na = append(na, i)
	}
	runtime.ReadMemStats(&end_mem_info)
	println("NA memory usage - LINEAR: ", end_mem_info.Alloc-start_mem_info.Alloc)

	// Clear memory
	na = nil
	runtime.GC()

	// Benchmark compressed - LINEAR
	runtime.ReadMemStats(&start_mem_info)
	ca := dam.New_Compressed_Array()
	for i := uint64(0); i < n; i++ {
		ca.Append(i + 1)
	}
	runtime.ReadMemStats(&end_mem_info)
	println("CA memory usage - LINEAR: ", end_mem_info.Alloc-start_mem_info.Alloc)
	// Confirm that the compressed array is correct
	for i := uint64(0); i < n; i++ {
		if ca.Get(i) != i+1 {
			println("CA is incorrect")
			break
		}
	}

	// Clear memory
	ca = nil
	runtime.GC()

	// // Create a new random source and generator
	// r := rand.New(rand.NewSource(time.Now().UnixNano()))

	// // Benchmark normal - RANDOM
	// runtime.ReadMemStats(&start_mem_info)
	// na = make([]uint64, 0)
	// for i := uint64(0); i < n; i++ {
	// 	na = append(na, uint64(r.Int()))
	// }
	// runtime.ReadMemStats(&end_mem_info)
	// println("NA memory usage - RANDOM: ", end_mem_info.Alloc-start_mem_info.Alloc)

	// // Clear memory
	// na = nil
	// runtime.GC()

	// // Benchmark compressed - RANDOM
	// runtime.ReadMemStats(&start_mem_info)
	// ca = dam.New_Compressed_Array()
	// for i := uint64(0); i < n; i++ {
	// 	ca.Append(uint64(r.Int()))
	// }
	// runtime.ReadMemStats(&end_mem_info)
	// println("CA memory usage - RANDOM: ", end_mem_info.Alloc-start_mem_info.Alloc)

	// // Clear memory
	// ca = nil
	// runtime.GC()

}
