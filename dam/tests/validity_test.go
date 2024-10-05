package dam_tests

import (
	"testing"

	"github.com/nacioboi/go_dam/dam"
)

func Test_Consistency(t *testing.T) {
	const num_inputs = 1025 // 1 more than a power of 2.
	dam_map := dam.New(
		uint64(num_inputs), dam.With_Performance_Profile[uint64, uint64](dam.PERFORMANCE_PROFILE__NORMAL),
	)

	for i := uint64(0); i < uint64(num_inputs); i++ {
		dam_map.Set(i+1, i)
	}

	for i := uint64(0); i < uint64(num_inputs); i++ {
		x, ok := dam_map.Get(i + 1)
		if !ok {
			t.Error("Key not found.")
		}
		if x != i {
			t.Error("Value mismatch.")
		}
	}
}
