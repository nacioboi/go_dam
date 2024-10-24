/*/
 ** This software is covered by the MIT License.
 ** See: `./LICENSE`.
/*/

package dam_tests

import (
	"math/rand/v2"
	"runtime"
	"testing"

	"github.com/nacioboi/go_dam/dam"
)

func Benchmark__FAST_DAM__Memory_Use_W_Overflows__(b *testing.B) {
	defer runtime.GC()
	runtime.GC()

	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	before := m.TotalAlloc

	x := dam.DAM_FAST[uint64, uint64]{}.New(uint64(4096))
	for i := uint64(0); i < 4096; i++ {
		x.Set(i+1, i)
		x.Set(i+1+4096, i+4096)
	}

	runtime.ReadMemStats(&m)
	after := m.TotalAlloc
	b.Logf("Benchmark__FAST_DAM__Memory_Use_W_Overflows__: %d", after-before)
}

func Benchmark__STD_DAM__Memory_Use_W_Overflows__(b *testing.B) {
	defer runtime.GC()
	runtime.GC()

	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	before := m.TotalAlloc

	x := dam.DAM_STD[uint64, uint64]{}.New(uint64(4096))
	for i := uint64(0); i < 4096; i++ {
		x.Set(i+1, i)
		x.Set(i+1+4096, i+4096)
	}

	runtime.ReadMemStats(&m)
	after := m.TotalAlloc
	b.Logf("Benchmark__STD_DAM__Memory_Use_W_Overflows__: %d", after-before)
}

func Benchmark__MOH_DAM__Memory_Use_W_Overflows__(b *testing.B) {
	defer runtime.GC()
	runtime.GC()

	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	before := m.TotalAlloc

	x := dam.DAM_MOH[uint64, uint64]{}.New(uint64(4096))
	for i := uint64(0); i < 4096; i++ {
		x.Set(i+1, i)
		x.Set(i+1+4096, i+4096)
	}

	runtime.ReadMemStats(&m)
	after := m.TotalAlloc
	b.Logf("Benchmark__MOH_DAM__Memory_Use_W_Overflows__: %d", after-before)
}

func Benchmark__FAST_DAM__Memory_Usage__(b *testing.B) {
	defer runtime.GC()
	runtime.GC()

	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	before := m.TotalAlloc

	x := dam.DAM_FAST[uint64, uint64]{}.New(uint64(4096))
	for i := uint64(0); i < 4096; i++ {
		x.Set(i+1, i)
	}

	runtime.ReadMemStats(&m)
	after := m.TotalAlloc
	b.Logf("Benchmark__FAST_DAM__Memory_Usage__: %d", after-before)
}

func Benchmark__STD_DAM__Memory_Usage__(b *testing.B) {
	defer runtime.GC()
	runtime.GC()

	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	before := m.TotalAlloc

	x := dam.DAM_STD[uint64, uint64]{}.New(uint64(4096))
	for i := uint64(0); i < 4096; i++ {
		x.Set(i+1, i)
	}

	runtime.ReadMemStats(&m)
	after := m.TotalAlloc
	b.Logf("Benchmark__STD_DAM__Memory_Usage__: %d", after-before)
}

func Benchmark__MOH_DAM__Memory_Usage__(b *testing.B) {
	defer runtime.GC()
	runtime.GC()

	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	before := m.TotalAlloc

	x := dam.DAM_MOH[uint64, uint64]{}.New(uint64(4096))
	for i := uint64(0); i < 4096; i++ {
		x.Set(i+1, i)
	}

	runtime.ReadMemStats(&m)
	after := m.TotalAlloc
	b.Logf("Benchmark__MOH_DAM__Memory_Usage__: %d", after-before)
}

func Benchmark__LOH_DAM__Memory_Usage__(b *testing.B) {
	defer runtime.GC()
	runtime.GC()

	// encode_value_f := func(key uint64, value uint64) []byte {
	// 	return []byte{byte(key), byte(value)}
	// }
	// decode_value_f := func(encoded []byte) (uint64, uint64) {
	// 	key := uint64(encoded[0])
	// 	value := uint64(encoded[1])
	// 	return key, value
	// }

	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	before := m.TotalAlloc

	x := dam.DAM_LOH[uint64, uint64]{}.New(uint64(4096))
	for i := uint64(0); i < 4096; i++ {
		x.Set(i+1, i)
	}

	runtime.ReadMemStats(&m)
	after := m.TotalAlloc
	b.Logf("Benchmark__LOH_DAM__Memory_Usage__: %d", after-before)
}

func Benchmark__Builtin_Map__Memory_Usage__(b *testing.B) {
	defer runtime.GC()
	runtime.GC()

	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	before := m.TotalAlloc

	x := make(map[int]int)
	for i := 0; i < 4096; i++ {
		x[i] = i
	}

	runtime.ReadMemStats(&m)
	after := m.TotalAlloc
	b.Logf("Benchmark__Builtin_Map__Memory_Usage__: %d", after-before)
}

func Benchmark__Linear_FAST_DAM__Set__(b *testing.B) {
	defer runtime.GC()

	dam_map := dam.DAM_FAST[uint64, uint64]{}.New(uint64(b.N))

	b.ResetTimer()
	for i := uint64(0); i < uint64(b.N); i++ {
		dam_map.Set(i+1, i)
	}

	dam_map = nil
}

func Benchmark__Linear_FAST_DAM__Get__(b *testing.B) {
	defer runtime.GC()

	dam_map := dam.DAM_FAST[uint64, uint64]{}.New(uint64(b.N))

	for i := uint64(0); i < uint64(b.N); i++ {
		dam_map.Set(i+1, i)
	}

	var t uint64
	b.ResetTimer()
	for i := uint64(0); i < uint64(b.N); i++ {
		x, ok := dam_map.Get(i + 1)
		if ok {
			t += x
		} else {
			panic("Key not found.")
		}
	}

	dam_map = nil
}

func Benchmark__Random_FAST_DAM__Set__(b *testing.B) {
	defer runtime.GC()

	dam_map := dam.DAM_FAST[uint64, uint64]{}.New(uint64(b.N))

	keys := generate_random_keys(b.N)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		dam_map.Set(uint64(keys[i]), uint64(i))
	}

	dam_map = nil
}

func Benchmark__Random_FAST_DAM__Get__(b *testing.B) {
	defer runtime.GC()

	dam_map := dam.DAM_FAST[uint64, uint64]{}.New(uint64(b.N))

	for i := 0; i < b.N; i++ {
		dam_map.Set(uint64(i+1), uint64(i))
	}

	keys := generate_random_keys(b.N)

	var t uint64
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		x, ok := dam_map.Get(uint64(keys[i]))
		if ok {
			t += x
		} else {
			panic("Key not found.")
		}
	}

	dam_map = nil
}

func Benchmark__Linear_STD_DAM__Set__(b *testing.B) {
	defer runtime.GC()

	dam_map := dam.DAM_STD[uint64, uint64]{}.New(uint64(b.N))
	b.ResetTimer()
	for i := uint64(0); i < uint64(b.N); i++ {
		dam_map.Set(i+1, i)
	}

	dam_map = nil
}

func Benchmark__Linear_STD_DAM__Get__(b *testing.B) {
	defer runtime.GC()

	dam_map := dam.DAM_STD[uint64, uint64]{}.New(uint64(b.N))

	for i := uint64(0); i < uint64(b.N); i++ {
		dam_map.Set(i+1, i)
	}

	var t uint64
	b.ResetTimer()
	for i := uint64(0); i < uint64(b.N); i++ {
		x, ok := dam_map.Get(i + 1)
		if ok {
			t += x
		} else {
			panic("Key not found.")
		}
	}

	dam_map = nil
}

func Benchmark__Random_STD_DAM__Set__(b *testing.B) {
	defer runtime.GC()

	dam_map := dam.DAM_STD[uint64, uint64]{}.New(uint64(b.N))

	keys := generate_random_keys(b.N)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		dam_map.Set(uint64(keys[i]), uint64(i))
	}

	dam_map = nil
}

func Benchmark__Random_STD_DAM__Get__(b *testing.B) {
	defer runtime.GC()

	dam_map := dam.DAM_STD[uint64, uint64]{}.New(uint64(b.N))

	for i := 0; i < b.N; i++ {
		dam_map.Set(uint64(i+1), uint64(i))
	}

	keys := generate_random_keys(b.N)

	var t uint64
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		x, ok := dam_map.Get(uint64(keys[i]))
		if ok {
			t += x
		} else {
			panic("Key not found.")
		}
	}

	dam_map = nil
}

func Benchmark__Linear_MOH_DAM__Set__(b *testing.B) {
	defer runtime.GC()

	dam_map := dam.DAM_MOH[uint64, uint64]{}.New(uint64(b.N))

	b.ResetTimer()
	for i := uint64(0); i < uint64(b.N); i++ {
		dam_map.Set(i+1, i)
	}

	dam_map = nil
}

func Benchmark__Linear_MOH_DAM__Get__(b *testing.B) {
	defer runtime.GC()

	dam_map := dam.DAM_MOH[uint64, uint64]{}.New(uint64(b.N))

	for i := uint64(0); i < uint64(b.N); i++ {
		dam_map.Set(i+1, i)
	}

	var t uint64
	b.ResetTimer()
	for i := uint64(0); i < uint64(b.N); i++ {
		x, ok := dam_map.Get(i + 1)
		if ok {
			t += x
		} else {
			panic("Key not found.")
		}
	}

	dam_map = nil
}

func Benchmark__Random_MOH_DAM__Set__(b *testing.B) {
	defer runtime.GC()

	dam_map := dam.DAM_MOH[uint64, uint64]{}.New(uint64(b.N))

	keys := generate_random_keys(b.N)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		dam_map.Set(uint64(keys[i]), uint64(i))
	}

	dam_map = nil
}

func Benchmark__Random_MOH_DAM__Get__(b *testing.B) {
	defer runtime.GC()

	dam_map := dam.DAM_MOH[uint64, uint64]{}.New(uint64(b.N))

	for i := 0; i < b.N; i++ {
		dam_map.Set(uint64(i+1), uint64(i))
	}

	keys := generate_random_keys(b.N)

	var t uint64
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		x, ok := dam_map.Get(uint64(keys[i]))
		if ok {
			t += x
		} else {
			panic("Key not found.")
		}
	}

	dam_map = nil
}

// func Benchmark__Linear_LOH_DAM__Set__(b *testing.B) {
// 	defer runtime.GC()

// 	encode_value := func(v uint64) []byte {
// 		return []byte{byte(v)}
// 	}
// 	decode_value := func(b []byte) uint64 {
// 		return uint64(b[0])
// 	}

// 	dam_map := dam.New_LOH_DAM[uint64, uint64](uint64(b.N), encode_value, decode_value)

// 	b.ResetTimer()
// 	for i := uint64(0); i < uint64(b.N); i++ {
// 		dam_map.Set(i+1, i)
// 	}

// 	dam_map = nil
// }

// func Benchmark__Linear_LOH_DAM__Get__(b *testing.B) {
// 	defer runtime.GC()

// 	encode_value := func(v uint64) []byte {
// 		return []byte{byte(v)}
// 	}
// 	decode_value := func(b []byte) uint64 {
// 		return uint64(b[0])
// 	}

// 	dam_map := dam.New_LOH_DAM[uint64, uint64](uint64(b.N), encode_value, decode_value)

// 	for i := uint64(0); i < uint64(b.N); i++ {
// 		dam_map.Set(i+1, i)
// 	}

// 	var t uint64
// 	b.ResetTimer()
// 	for i := uint64(0); i < uint64(b.N); i++ {
// 		x, ok := dam_map.Get(i + 1)
// 		if ok {
// 			t += x
// 		} else {
// 			b.Fatalf("Key not found: %d", i+1)
// 		}
// 	}

// 	dam_map = nil
// }

// func Benchmark__Random_LOH_DAM__Set__(b *testing.B) {
// 	defer runtime.GC()

// 	encode_value := func(v uint64) []byte {
// 		return []byte{byte(v)}
// 	}
// 	decode_value := func(b []byte) uint64 {
// 		return uint64(b[0])
// 	}

// 	dam_map := dam.New_LOH_DAM[uint64, uint64](uint64(b.N), encode_value, decode_value)

// 	keys := generate_random_keys(b.N)

// 	b.ResetTimer()
// 	for i := 0; i < b.N; i++ {
// 		dam_map.Set(uint64(keys[i]), uint64(i))
// 	}

// 	dam_map = nil
// }

// func Benchmark__Random_LOH_DAM__Get__(b *testing.B) {
// 	defer runtime.GC()

// 	encode_value := func(v uint64) []byte {
// 		return []byte{byte(v)}
// 	}
// 	decode_value := func(b []byte) uint64 {
// 		return uint64(b[0])
// 	}

// 	dam_map := dam.New_LOH_DAM[uint64, uint64](uint64(b.N), encode_value, decode_value)

// 	for i := 0; i < b.N; i++ {
// 		dam_map.Set(uint64(i+1), uint64(i))
// 	}

// 	keys := generate_random_keys(b.N)

// 	var t uint64
// 	b.ResetTimer()
// 	for i := 0; i < b.N; i++ {
// 		x, ok := dam_map.Get(uint64(keys[i]))
// 		if ok {
// 			t += x
// 		} else {
// 			panic("Key not found.")
// 		}
// 	}

// 	dam_map = nil
// }

func Benchmark__Linear_Builtin_Map__Set__(b *testing.B) {
	defer runtime.GC()

	builtin_map := make(map[int]int)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		builtin_map[i+1] = i
	}

	builtin_map = nil
}

func Benchmark__Linear_Builtin_Map__Get__(b *testing.B) {
	defer runtime.GC()

	builtin_map := make(map[int]int)

	for i := 0; i < b.N; i++ {
		builtin_map[i+1] = i
	}

	var t int
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		x, ok := builtin_map[i+1]
		if ok {
			t += x
		} else {
			panic("Key not found.")
		}
	}

	builtin_map = nil
}

func Benchmark__Random_Builtin_Map__Set__(b *testing.B) {
	defer runtime.GC()

	builtin_map := make(map[int]int)

	keys := generate_random_keys(b.N)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		builtin_map[keys[i]] = i
	}

	builtin_map = nil
}

func Benchmark__Random_Builtin_Map__Get__(b *testing.B) {
	defer runtime.GC()

	builtin_map := make(map[int]int)

	for i := 0; i < b.N; i++ {
		builtin_map[i+1] = i
	}

	keys := generate_random_keys(b.N)

	var t uint64
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		x, ok := builtin_map[keys[i]]
		if ok {
			t += uint64(x)
		} else {
			panic("Key not found.")
		}
	}

	builtin_map = nil
}

func Benchmark__FAST_DAM__Get__W_Overflows__(b *testing.B) {
	defer runtime.GC()

	dam_map := dam.DAM_FAST[uint64, uint64]{}.New(uint64(b.N))

	for i := 0; i < b.N; i++ {
		dam_map.Set(uint64(i+1), uint64(i))
		dam_map.Set(uint64((i+1)*2), uint64(i*2))
	}

	var t uint64
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		x, ok := dam_map.Get(uint64((i + 1) * 2))
		if ok {
			t += x
		} else {
			panic("Key not found.")
		}
	}

	dam_map = nil
}

func Benchmark__STD_DAM__Get__W_Overflows__(b *testing.B) {
	defer runtime.GC()

	dam_map := dam.DAM_STD[uint64, uint64]{}.New(uint64(b.N))

	for i := 0; i < b.N; i++ {
		dam_map.Set(uint64(i+1), uint64(i))
		dam_map.Set(uint64((i+1)*2), uint64(i*2))
	}

	var t uint64
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		x, ok := dam_map.Get(uint64((i + 1) * 2))
		if ok {
			t += x
		} else {
			panic("Key not found.")
		}
	}

	dam_map = nil
}

func Benchmark__MOH_DAM__Get__W_Overflows__(b *testing.B) {
	defer runtime.GC()

	dam_map := dam.DAM_MOH[uint64, uint64]{}.New(uint64(b.N))

	for i := 0; i < b.N; i++ {
		dam_map.Set(uint64(i+1), uint64(i))
		dam_map.Set(uint64((i+1)*2), uint64(i*2))
	}

	var t uint64
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		x, ok := dam_map.Get(uint64((i + 1) * 2))
		if ok {
			t += x
		} else {
			panic("Key not found.")
		}
	}

	dam_map = nil
}

func generate_random_keys(n int) []int {
	keys := make([]int, n)
	for i := 0; i < n; i++ {
		keys[i] = i + 1
	}
	rand.Shuffle(int(n), func(i, j int) {
		keys[i], keys[j] = keys[j], keys[i]
	})
	return keys
}
