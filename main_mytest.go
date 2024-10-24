package main

import (
	"log"
	"math/rand/v2"
	"os"
	"runtime/pprof"

	"github.com/nacioboi/go_dam/dam"
)

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

func main_mytest() {
	// Open the file for cpu profiling
	f, err := os.Create("cpu.prof")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// encode_value_f := func(key uint64, value uint64) []byte {
	// 	return []byte{byte(key), byte(value)}
	// }
	// decode_value_f := func(encoded []byte) (uint64, uint64) {
	// 	key := uint64(encoded[0])
	// 	value := uint64(encoded[1])
	// 	return key, value
	// }

	n := uint64(1024 * 1024) // * 512)
	dam := dam.DAM_MOH[uint64, uint64]{}.New(n)
	dam.Set(n+2, n+1)

	var x uint64
	var ok bool
	var t uint64

	random_keys := Generate_Random_Keys(n)

	pprof.StartCPUProfile(f)
	for i := uint64(0); i < n; i++ {
		x, ok = dam.Get(random_keys[i])
		if !ok {
			panic("Key not found.")
		}
		t += x
	}
	pprof.StopCPUProfile()
}
