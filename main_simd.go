package main

import (
	"fmt"
	"log"
	"math/rand/v2"
	"os"
	"runtime/debug"
	"runtime/pprof"
	"strings"
	"time"
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

func golang_find_idx(key uint64, keys [64]uint64) (uint8, bool) {
	for i, k := range keys {
		if k == key {
			return uint8(i), true
		}
	}
	return 0, false
}

// Warning: This function must be called with keys satisfying the following conditions:
// 1. len(keys) % 4 == 0
// 2. no duplicate keys.
//
//go:noescape
//go:nosplit
func avx512_find_idx_64(key uint64, arr [64]uint64) (uint8, bool)

func simd_find_idx(ptr *uint64, n int, key uint64) uint16

func std(keys []uint64, key uint64) (uint8, bool) {
	for i, k := range keys {
		if k == key {
			return uint8(i), true
		}
	}
	return 0, false
}

var keys = [64]uint64{
	1, 2, 3, 4, 5, 6, 7, 8, 9, 10,
	11, 12, 13, 14, 15, 16, 17, 18, 19, 20,
	21, 22, 23, 24, 25, 26, 27, 28, 29, 30,
	31, 32, 33, 34, 35, 36, 37, 38, 39, 40,
	41, 42, 43, 44, 45, 46, 47, 48, 49, 50,
	51, 52, 53, 54, 55, 56, 57, 58, 59, 60,
	61, 62, 63, 64,
}

var start time.Time

const n = 1024 * 1024 * 64 //* 768

func main_simd() {
	if len(keys)%8 != 0 {
		panic("len(keys) must be a multiple of 8")
	}
	if cap(keys) != len(keys) {
		panic("cap(keys) must be equal to len(keys)")
	}
	// Check for duplicates
	// m := make(map[uint64]struct{})
	// for _, k := range keys {
	// 	if _, ok := m[k]; ok {
	// 		panic("duplicate key")
	// 	}
	// 	m[k] = struct{}{}
	// }
	// m = nil

	// Open file for benchmarking
	f, err := os.Create("cpu.prof")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	debug.SetGCPercent(-1)
	defer debug.SetGCPercent(100)

	var v uint8
	var ok bool

	time_total_std := time.Duration(0)
	time_total_v2 := time.Duration(0)

	// Benchmark v2
	pprof.StartCPUProfile(f)
	for i := 0; i < n; i++ {
		to_find := rand.Uint64() % 64
		if to_find == 0 {
			to_find = 64
		}

		start = time.Now()
		v, ok = avx512_find_idx_64(to_find, keys)
		if !ok {
			panic("not found")
		}
		if v != uint8(to_find-1) {
			log.Fatalf("wrong value: %d", v)
		}
		time_total_v2 += time.Since(start)

		// Benchmark std
		start = time.Now()
		v, ok = golang_find_idx(to_find, keys)
		if !ok {
			panic("not found")
		}
		if v != uint8(to_find-1) {
			log.Fatalf("wrong value: %d", v)
		}
		time_total_std += time.Since(start)
	}
	pprof.StopCPUProfile()

	fmt.Printf("v2 average time: %s\n", time_total_v2/time.Duration(n))
	fmt.Printf("std average time: %s\n", time_total_std/time.Duration(n))

}
