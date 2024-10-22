/*/
 ** This software is covered by the MIT License.
 ** See: `./LICENSE`.
/*/

package main

import (
	"fmt"
	"log"
	"math/rand"
	"runtime"
	"time"

	"github.com/nacioboi/go_dam/dam"
)

func main() {
	runtime.GC()

	const n = 1024 //<< 20 // 1M

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
			fmt.Printf("CA is incorrect: %d != %d", ca.Get(i), i+1)
			break
		}
	}

	// Clear memory
	ca = nil
	runtime.GC()

	// Create a new random source and generator
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Benchmark normal - RANDOM
	runtime.ReadMemStats(&start_mem_info)
	na = make([]uint64, 0)
	for i := uint64(0); i < n; i++ {
		na = append(na, uint64(r.Int()))
	}
	runtime.ReadMemStats(&end_mem_info)
	println("NA memory usage - RANDOM: ", end_mem_info.Alloc-start_mem_info.Alloc)

	// Clear memory
	na = nil
	runtime.GC()

	// Benchmark compressed - RANDOM
	ca_rv := make([]uint64, 0) // ca -> random values
	runtime.ReadMemStats(&start_mem_info)
	ca = dam.New_Compressed_Array()
	for i := uint64(0); i < n; i++ {
		x := r.Uint64() % 65535 * 4
		ca.Append(x)
		ca_rv = append(ca_rv, x)
	}
	runtime.ReadMemStats(&end_mem_info)
	println("CA memory usage - RANDOM: ", end_mem_info.Alloc-start_mem_info.Alloc)

	// Confirm that the compressed array is correct
	for i := uint64(0); i < n; i++ {
		if ca.Get(i) != ca_rv[i] {
			log.Fatalf("CA is incorrect: %d != %d", ca.Get(i), ca_rv[i])
			break
		}
	}

	// Clear memory
	ca = nil
	runtime.GC()

}
