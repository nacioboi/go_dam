/*/
 ** This software is covered by the MIT License.
 ** See: `./LICENSE`.
/*/

package main

import (
	"fmt"
	"strings"
)

func format_Number_With_Commas(n int64) string {
	s := fmt.Sprintf("%d", n)
	if n < 0 {
		s = s[1:]
	}
	var result []string
	for len(s) > 3 {
		result = append(result, s[len(s)-3:])
		s = s[:len(s)-3]
	}
	result = append(result, s)
	for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
		result[i], result[j] = result[j], result[i]
	}
	formattedNumber := strings.Join(result, ",")
	if n < 0 {
		formattedNumber = "-" + formattedNumber
	}
	return formattedNumber
}

// func asm_std(keys []uint64, key uint64) (uint8, bool)

// // Warning: This function must be called with keys satisfying the following conditions:
// // 1. len(keys) % 4 == 0
// // 2. no duplicate keys.
// func simd_find_idx(keys []uint64, key uint64) uint8

// var keys = []uint64{
// 	1,
// 	2,
// 	3,
// 	4,
// 	5,
// 	6,
// 	7,
// 	8,
// 	9,
// 	10,
// 	11,
// 	12,
// 	13,
// 	14,
// 	15,
// 	16,
// 	17,
// 	18,
// 	19,
// 	20,
// 	21,
// 	22,
// 	23,
// 	24,
// 	25,
// 	26,
// 	27,
// 	28,
// 	29,
// 	30,
// 	31,
// 	32,
// 	152,
// 	222,
// 	332,
// 	442,
// }
// var start time.Time

// const n = 1_000_000 * 128

// func main() {
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

// 	// Benchmark v2
// 	start = time.Now()
// 	pprof.StartCPUProfile(f)
// 	for i := 0; i < n; i++ {
// 		v = simd_find_idx(keys, 222)
// 		if v == 0 {
// 			panic("not found")
// 		}
// 		if v != 34 {
// 			log.Fatalf("wrong value: %d", v)
// 		}
// 	}
// 	pprof.StopCPUProfile()
// 	fmt.Println("v2:", time.Since(start))

// 	// Benchmark std
// 	start = time.Now()
// 	for i := 0; i < n; i++ {
// 		v, ok = asm_std(keys, 32)
// 		if !ok {
// 			panic("not found")
// 		}
// 		if v != 31 {
// 			log.Fatalf("wrong value: %d", v)
// 		}
// 	}
// 	fmt.Println("std:", time.Since(start))

// }

// Refactor test goals:
// - Easy to read.
// - Good, readable output.
// - Easy to add new tests.
// - Easy to add new benchmarks.

//func main() {
// debug.SetGCPercent(-1)
// defer debug.SetGCPercent(100)
// defer runtime.GC()

// f, err := os.Create("cpu.prof")
// if err != nil {
// 	log.Fatal(err)
// }
// defer f.Close()

// n_normal := uint64(1024 * 1024 * 32)
// n_memory := uint64(1024 * 1024)

// // Create maps...
// bm := make(map[uint64]uint64)
// dam_save := dam_map.New[uint64, uint64](
// 	n_normal,
// 	dam_map.With_Performance_Profile[uint64, uint64](dam_map.PERFORMANCE_PROFILE__SAVE_MEMORY),
// )
// dam_normal := dam_map.New[uint64, uint64](
// 	n_normal,
// 	dam_map.With_Performance_Profile[uint64, uint64](dam_map.PERFORMANCE_PROFILE__NORMAL),
// )
// dam_fast := dam_map.New[uint64, uint64](
// 	n_normal,
// 	dam_map.With_Performance_Profile[uint64, uint64](dam_map.PERFORMANCE_PROFILE__FAST),
// )

// bm_m_f := func() map[uint64]uint64 {
// 	return make(map[uint64]uint64)
// }

// dam_save_m_f := func() *dam_map.DAM[uint64, uint64] {
// 	return dam_map.New[uint64, uint64](
// 		n_memory,
// 		dam_map.With_Performance_Profile[uint64, uint64](dam_map.PERFORMANCE_PROFILE__SAVE_MEMORY),
// 	)
// }
// dam_normal_m_f := func() *dam_map.DAM[uint64, uint64] {
// 	return dam_map.New[uint64, uint64](
// 		n_memory,
// 		dam_map.With_Performance_Profile[uint64, uint64](dam_map.PERFORMANCE_PROFILE__NORMAL),
// 	)
// }
// dam_fast_m_f := func() *dam_map.DAM[uint64, uint64] {
// 	return dam_map.New[uint64, uint64](
// 		n_memory,
// 		dam_map.With_Performance_Profile[uint64, uint64](dam_map.PERFORMANCE_PROFILE__FAST),
// 	)
// }

// var res tests.Test_Result

// // Benchmark Linear Set...
// res = tests.Bench_Linear_Builtin_Map_Set(bm, n_normal)
// fmt.Printf("\nBM       :: LINEAR SET            :: %d\n", res.Elapsed_Time)
// fmt.Printf("BM       :: MICROSECONDS PER OP   :: %f\n", float64(res.Elapsed_Time)/float64(n_normal))
// res = tests.Bench_Linear_DAM_Set(dam_save, n_normal)
// fmt.Printf("DAM SAVE  :: LINEAR SET          :: %d\n", res.Elapsed_Time)
// fmt.Printf("DAM SAVE  :: MICROSECONDS PER OP :: %f\n", float64(res.Elapsed_Time)/float64(n_normal))
// res = tests.Bench_Linear_DAM_Set(dam_normal, n_normal)
// fmt.Printf("DAM NORMAL  :: LINEAR SET        :: %d\n", res.Elapsed_Time)
// fmt.Printf("DAM NORMAL  :: MICROSECONDS PER OP :: %f\n", float64(res.Elapsed_Time)/float64(n_normal))
// res = tests.Bench_Linear_DAM_Set(dam_fast, n_normal)
// fmt.Printf("DAM FAST  :: LINEAR SET          :: %d\n", res.Elapsed_Time)
// fmt.Printf("DAM FAST  :: MICROSECONDS PER OP :: %f\n", float64(res.Elapsed_Time)/float64(n_normal))

// // Benchmark Linear Get...
// res = tests.Bench_Linear_Builtin_Map_Get(bm, n_normal)
// fmt.Printf("\nBM       :: LINEAR GET            :: %d\n", res.Elapsed_Time)
// fmt.Printf("BM       :: MICROSECONDS PER OP   :: %f\n", float64(res.Elapsed_Time)/float64(n_normal))
// bm_checksum := res.Checksum
// res = tests.Bench_Linear_DAM_Get(dam_save, n_normal)
// fmt.Printf("DAM SAVE  :: LINEAR GET          :: %d\n", res.Elapsed_Time)
// fmt.Printf("DAM SAVE  :: MICROSECONDS PER OP :: %f\n", float64(res.Elapsed_Time)/float64(n_normal))
// dam_save_checksum := res.Checksum
// res = tests.Bench_Linear_DAM_Get(dam_normal, n_normal)
// fmt.Printf("DAM NORMAL  :: LINEAR GET        :: %d\n", res.Elapsed_Time)
// fmt.Printf("DAM NORMAL  :: MICROSECONDS PER OP :: %f\n", float64(res.Elapsed_Time)/float64(n_normal))
// dam_normal_checksum := res.Checksum
// res = tests.Bench_Linear_DAM_Get(dam_fast, n_normal)
// fmt.Printf("DAM FAST  :: LINEAR GET          :: %d\n", res.Elapsed_Time)
// fmt.Printf("DAM FAST  :: MICROSECONDS PER OP :: %f\n", float64(res.Elapsed_Time)/float64(n_normal))
// dam_fast_checksum := res.Checksum

// // Benchmark Random get...
// data := tests.Generate_Random_Keys(n_normal)
// res = tests.Bench_Random_Builtin_Map_Get(bm, data)
// fmt.Printf("\nBM       :: RANDOM GET            :: %d\n", res.Elapsed_Time)
// fmt.Printf("BM       :: MICROSECONDS PER OP   :: %f\n", float64(res.Elapsed_Time)/float64(n_normal))
// res = tests.Bench_Random_DAM_Get(dam_save, data)
// fmt.Printf("DAM SAVE  :: RANDOM GET          :: %d\n", res.Elapsed_Time)
// fmt.Printf("DAM SAVE  :: MICROSECONDS PER OP :: %f\n", float64(res.Elapsed_Time)/float64(n_normal))
// res = tests.Bench_Random_DAM_Get(dam_normal, data)
// fmt.Printf("DAM NORMAL  :: RANDOM GET        :: %d\n", res.Elapsed_Time)
// fmt.Printf("DAM NORMAL  :: MICROSECONDS PER OP :: %f\n", float64(res.Elapsed_Time)/float64(n_normal))
// pprof.StartCPUProfile(f)
// res = tests.Bench_Random_DAM_Get(dam_fast, data)
// pprof.StopCPUProfile()
// fmt.Printf("DAM FAST  :: RANDOM GET          :: %d\n", res.Elapsed_Time)
// fmt.Printf("DAM FAST  :: MICROSECONDS PER OP :: %f\n", float64(res.Elapsed_Time)/float64(n_normal))

// // Benchmark memory usage...
// var dam_max_mem uint64
// var dam_save_mem uint64
// var dam_normal_mem uint64
// var dam_fast_mem uint64

// res = tests.Bench_Mem_Usage_Builtin_Map(bm_m_f, n_memory)
// fmt.Printf("\nBM       :: MEMORY USAGE           :: %s\n", format_Number_With_Commas(int64(res.Memory_Usage)))

// res = tests.Bench_Mem_Usage_DAM(dam_save_m_f, n_memory)
// dam_max_mem = max(dam_max_mem, res.Memory_Usage)
// dam_save_mem = res.Memory_Usage
// fmt.Printf("DAM SAVE  :: MEMORY USAGE         :: %s\n", format_Number_With_Commas(int64(res.Memory_Usage)))

// res = tests.Bench_Mem_Usage_DAM(dam_normal_m_f, n_memory)
// dam_max_mem = max(dam_max_mem, res.Memory_Usage)
// dam_normal_mem = res.Memory_Usage
// fmt.Printf("DAM NORMAL  :: MEMORY USAGE       :: %s\n", format_Number_With_Commas(int64(res.Memory_Usage)))

// res = tests.Bench_Mem_Usage_DAM(dam_fast_m_f, n_memory)
// dam_max_mem = max(dam_max_mem, res.Memory_Usage)
// dam_fast_mem = res.Memory_Usage
// fmt.Printf("DAM FAST  :: MEMORY USAGE         :: %s\n", format_Number_With_Commas(int64(res.Memory_Usage)))

// // Print memory percentage of max memory...
// fmt.Printf("\nSFDA SAVE :: MEMORY USAGE PERCENTAGE :: %f\n", float64(dam_save_mem)/float64(dam_max_mem))
// fmt.Printf("DAM NORMAL  :: MEMORY USAGE PERCENTAGE :: %f\n", float64(dam_normal_mem)/float64(dam_max_mem))
// fmt.Printf("DAM FAST  :: MEMORY USAGE PERCENTAGE :: %f\n", float64(dam_fast_mem)/float64(dam_max_mem))

// // Print checksums...
// fmt.Printf("\nBM       :: CHECKSUM               :: %d\n", bm_checksum)
// fmt.Printf("DAM SAVE  :: CHECKSUM             :: %d\n", dam_save_checksum)
// fmt.Printf("DAM NORMAL  :: CHECKSUM           :: %d\n", dam_normal_checksum)
// fmt.Printf("DAM FAST  :: CHECKSUM             :: %d\n", dam_fast_checksum)

// // Assert checksums...
// for _, checksum := range []uint64{dam_save_checksum, dam_normal_checksum, dam_fast_checksum} {
// 	if checksum != bm_checksum {
// 		log.Fatalf("Checksums do not match!")
// 	}
// }
// fmt.Printf("\nChecksums match!\n")

//}

func main() {
}
