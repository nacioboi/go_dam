package dam_tests

import (
	"fmt"
	"testing"

	"github.com/nacioboi/go_dam/dam"
)

func Test_Consistency_LOH(t *testing.T) {
	const num_inputs = 1025 // 1 more than a power of 2.
	dam_map := dam.New_LOH_DAM[uint64, uint64](uint64(num_inputs))

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

func Test_Consistency_MOH(t *testing.T) {
	const num_inputs = 1025 // 1 more than a power of 2.
	dam_map := dam.New_MOH_DAM[uint64, uint64](uint64(num_inputs))

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

func Test_Consistency_STD(t *testing.T) {
	const num_inputs = 1025 // 1 more than a power of 2.
	dam_map := dam.New_STD_DAM[uint64, uint64](uint64(num_inputs))

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

func Test_Consistency_Fast(t *testing.T) {
	const num_inputs = 1025 // 1 more than a power of 2.
	dam_map := dam.New_Fast_DAM[uint64, uint64](uint64(num_inputs))

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

func Test_Collision_LOH(t *testing.T) {
	const num_inputs = 384
	dam_map := dam.New_LOH_DAM[uint64, uint64](uint64(num_inputs))

	const v = 2
	for i := uint64(0); i < uint64(num_inputs); i++ {
		dam_map.Set((i+1)*v, i*v)
	}

	for i := uint64(0); i < uint64(num_inputs); i++ {
		x, ok := dam_map.Get((i + 1) * v)
		if !ok {
			t.Error("Key not found.")
		}
		if x != i*v {
			fmt.Printf("Expected: %d, Got: %d\n", i*v, x)
			t.Error("Value mismatch.")
		}
	}
}

func Test_Collision_MOH(t *testing.T) {
	const num_inputs = 1024
	dam_map := dam.New_MOH_DAM[uint64, uint64](uint64(num_inputs))

	const v = 32
	for i := uint64(0); i < uint64(num_inputs*v); i++ {
		dam_map.Set((i+1)*v, i*v)
	}

	for i := uint64(0); i < uint64(num_inputs*v); i++ {
		x, ok := dam_map.Get((i + 1) * v)
		if !ok {
			t.Error("Key not found.")
		}
		if x != i*v {
			t.Error("Value mismatch.")
		}
	}
}

func Test_Collision_STD(t *testing.T) {
	const num_inputs = 1024
	dam_map := dam.New_STD_DAM[uint64, uint64](uint64(num_inputs))

	const v = 32
	for i := uint64(0); i < uint64(num_inputs*v); i++ {
		dam_map.Set((i+1)*v, i*v)
	}

	for i := uint64(0); i < uint64(num_inputs*v); i++ {
		x, ok := dam_map.Get((i + 1) * v)
		if !ok {
			t.Error("Key not found.")
		}
		if x != i*v {
			t.Error("Value mismatch.")
		}
	}
}

func Test_Collision_Fast(t *testing.T) {
	const num_inputs = 128
	dam_map := dam.New_Fast_DAM[uint64, uint64](uint64(num_inputs))

	const v = 2
	for i := uint64(0); i < uint64(num_inputs*v); i++ {
		dam_map.Set((i+1)*v, i*v)
	}

	for i := uint64(0); i < uint64(num_inputs*v); i++ {
		x, ok := dam_map.Get((i + 1) * v)
		if !ok {
			t.Error("Key not found.")
		}
		if x != i*v {
			t.Error("Value mismatch.")
		}
	}
}
