/*/
 ** This software is covered by the MIT License.
 ** See: `./LICENSE`.
/*/

package dam_tests

import (
	"math/rand"
	"testing"

	"github.com/nacioboi/go_dam_map/dam"
)

func Benchmark_Linear_Builtin_Map_Set(b *testing.B) {
	builtin_map := make(map[int]int)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		builtin_map[i+1] = i
	}
}

func Benchmark_Linear_Builtin_Map_Get(b *testing.B) {
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
}

func Benchmark_Random_Builtin_Map_Set(b *testing.B) {
	builtin_map := make(map[int]int)
	keys := generate_random_keys(b.N)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		builtin_map[keys[i]] = i
	}
}

func Benchmark_Random_Builtin_Map_Get(b *testing.B) {
	builtin_map := make(map[int]int)

	for i := 0; i < b.N; i++ {
		builtin_map[i+1] = i
	}

	keys := generate_random_keys(b.N)

	var t uint64
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, key := range keys {
			x, ok := builtin_map[key]
			if ok {
				t += uint64(x)
			} else {
				panic("Key not found.")
			}
		}
	}
}

func Benchmark_Linear_DAM_Set(b *testing.B) {
	dam_map := dam.New(
		uint64(b.N), dam.With_Performance_Profile[uint64, uint64](dam.PERFORMANCE_PROFILE__NORMAL),
	)
	b.ResetTimer()
	for i := uint64(0); i < uint64(b.N); i++ {
		dam_map.Set(i+1, i)
	}
}

func Benchmark_Linear_FAST_DAM_Get(b *testing.B) {
	dam_map := dam.New(
		uint64(b.N), dam.With_Performance_Profile[uint64, uint64](dam.PERFORMANCE_PROFILE__NORMAL),
	)

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
}

func Benchmark_Random_DAM_Set(b *testing.B) {
	dam_map := dam.New(
		uint64(b.N), dam.With_Performance_Profile[uint64, uint64](dam.PERFORMANCE_PROFILE__NORMAL),
	)
	keys := generate_random_keys(b.N)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		dam_map.Set(uint64(keys[i]), uint64(i))
	}
}

func Benchmark_Random_DAM_Get(b *testing.B) {
	dam_map := dam.New(
		uint64(b.N), dam.With_Performance_Profile[uint64, uint64](dam.PERFORMANCE_PROFILE__NORMAL),
	)

	for i := 0; i < b.N; i++ {
		dam_map.Set(uint64(i+1), uint64(i))
	}

	keys := generate_random_keys(b.N)

	var t uint64
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, key := range keys {
			x, ok := dam_map.Get(uint64(key))
			if ok {
				t += x
			} else {
				panic("Key not found.")
			}
		}
	}
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
