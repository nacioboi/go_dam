/*/
 ** This software is covered by the MIT License.
 ** See: `./LICENSE`.
/*/

package main

import (
	"log"
	"math/rand/v2"
	"os"
	"runtime/pprof"

	"github.com/nacioboi/go_dam/dam"
)

// func format_Number_With_Commas(n int64) string {
// 	s := fmt.Sprintf("%d", n)
// 	if n < 0 {
// 		s = s[1:]
// 	}
// 	var result []string
// 	for len(s) > 3 {
// 		result = append(result, s[len(s)-3:])
// 		s = s[:len(s)-3]
// 	}
// 	result = append(result, s)
// 	for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
// 		result[i], result[j] = result[j], result[i]
// 	}
// 	formattedNumber := strings.Join(result, ",")
// 	if n < 0 {
// 		formattedNumber = "-" + formattedNumber
// 	}
// 	return formattedNumber
// }

// func asm_std(keys []uint64, key uint64) (uint8, bool)

// // Warning: This function must be called with keys satisfying the following conditions:
// // 1. len(keys) % 4 == 0
// // 2. no duplicate keys.
// //
// //go:noescape
// //go:nosplit
// func simd_find_idx(ptr *uint64, length uint64, key uint64) uint8

// var keys = []uint64{
// 	1, 2, 3, 4, 5, 6, 7, 8, 9, 10,
// 	11, 12, 13, 14, 15, 16, 17, 18, 19, 20,
// 	21, 22, 23, 24, 25, 26, 27, 28, 29, 30,
// 	31, 32, 33, 34, 35, 36, 37, 38, 39, 40,
// 	41, 42, 43, 44, 45, 46, 47, 48, 49, 50,
// 	51, 52, 53, 54, 55, 56, 57, 58, 59, 60,
// 	61, 62, 63, 64, 65, 66, 67, 68, 69, 70,
// 	71, 72, 73, 74, 75, 76, 77, 78, 79, 80,
// 	81, 82, 83, 84, 85, 86, 87, 88, 89, 90,
// 	91, 92, 93, 94, 95, 96, 97, 98, 99, 100,
// 	101, 102, 103, 104, 105, 106, 107, 108, 109, 110,
// 	111, 112, 113, 114, 115, 116, 117, 118, 119, 120,
// 	121, 122, 123, 124, 125, 126, 127, 128,
// }
// var start time.Time

// const n = 1024 * 1024 * 256

// func std(keys []uint64, key uint64) (uint8, bool) {
// 	for i, k := range keys {
// 		if k == key {
// 			return uint8(i), true
// 		}
// 	}
// 	return 0, false
// }

// func main() {
// 	if len(keys)%8 != 0 {
// 		panic("len(keys) must be a multiple of 8")
// 	}
// 	if cap(keys) != len(keys) {
// 		panic("cap(keys) must be equal to len(keys)")
// 	}
// 	// Check for duplicates
// 	// m := make(map[uint64]struct{})
// 	// for _, k := range keys {
// 	// 	if _, ok := m[k]; ok {
// 	// 		panic("duplicate key")
// 	// 	}
// 	// 	m[k] = struct{}{}
// 	// }
// 	// m = nil

// 	// Open file for benchmarking
// 	f, err := os.Create("cpu.prof")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer f.Close()

// 	debug.SetGCPercent(-1)
// 	defer debug.SetGCPercent(100)

// 	var v uint8
// 	var ok bool

// 	time_total_std := time.Duration(0)
// 	time_total_v2 := time.Duration(0)

// 	// Benchmark v2
// 	for i := 0; i < n; i++ {
// 		to_find := rand.Uint64() % 128
// 		if to_find == 0 {
// 			to_find = 1
// 		}

// 		start = time.Now()
// 		v = simd_find_idx(&keys[0], uint64(len(keys)), to_find)
// 		if v == 0 {
// 			panic("not found")
// 		}
// 		if v != uint8(to_find) {
// 			log.Fatalf("wrong value: %d", v)
// 		}
// 		time_total_v2 += time.Since(start)

// 		// Benchmark std
// 		start = time.Now()
// 		v, ok = asm_std(keys, to_find)
// 		if !ok {
// 			panic("not found")
// 		}
// 		if v != uint8(to_find-1) {
// 			log.Fatalf("wrong value: %d", v)
// 		}
// 		time_total_std += time.Since(start)
// 	}

// 	fmt.Printf("v2 average time: %s\n", time_total_v2/time.Duration(n))
// 	fmt.Printf("std average time: %s\n", time_total_std/time.Duration(n))

// }

func Generate_Random_Keys(n uint64) []uint64 {
	keys := make([]uint64, n)
	for i := uint64(0); i < n; i++ {
		keys[i] = i + 1
	}
	rand.Shuffle(int(n), func(i, j int) {
		keys[i], keys[j] = keys[j], keys[i]
	})
	return keys
}

func main() {
	// Open the file for cpu profiling
	f, err := os.Create("cpu.prof")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	n := uint64(1024 * 1024 * 32)
	dam := dam.New_MOH_DAM[uint64, uint64](1024)
	dam.Set(n+2, n+1)

	var x uint64
	var t uint64

	random_keys := Generate_Random_Keys(n)

	pprof.StartCPUProfile(f)
	for i := uint64(0); i < n; i++ {
		x, _ = dam.Get(random_keys[i])
		t += x
	}
	pprof.StopCPUProfile()
}
