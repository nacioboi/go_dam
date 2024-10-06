# DAM - a Direct Access Map??

Yes, you read that right. We gain massive performance improvements by using positive integers as keys instead of just any type.

## Introduction

This project aims to provide a hash map that is stupid fast when compared to the built-in `map`.

## Goals (Top priority first)

- [x] Stupid fast access time compared to the built-in `map` type.
- [x] Implement a hash map with O(1) access time.
- [ ] Implement a hash map with O(1) insert time.
- [ ] Clean and readable code.
- [ ] Easy to use.
- [ ] Configurable.

## See for yourself:

```text
goos: windows
goarch: amd64
pkg: github.com/nacioboi/go_dam/dam/tests
cpu: AMD Ryzen 5 7600 6-Core Processor
Benchmark__Linear_FAST_DAM__Set__-12            410969740                2.901 ns/op           0 B/op          0 allocs/op
Benchmark__Linear_FAST_DAM__Get__-12            793708011                1.515 ns/op           0 B/op          0 allocs/op
Benchmark__Random_FAST_DAM__Set__-12            72179148                18.79 ns/op            0 B/op          0 allocs/op
Benchmark__Random_FAST_DAM__Get__-12            100000000               11.64 ns/op            0 B/op          0 allocs/op
Benchmark__Linear_NORMAL_DAM__Set__-12          377219346                3.180 ns/op           0 B/op          0 allocs/op
Benchmark__Linear_NORMAL_DAM__Get__-12          638553973                1.903 ns/op           0 B/op          0 allocs/op
Benchmark__Random_NORMAL_DAM__Set__-12          61178605                31.09 ns/op            0 B/op          0 allocs/op
Benchmark__Random_NORMAL_DAM__Get__-12          46552225                27.42 ns/op            0 B/op          0 allocs/op
Benchmark__Linear_SAVE_MEMORY_DAM__Set__-12     171420146                7.243 ns/op           0 B/op          0 allocs/op
Benchmark__Linear_SAVE_MEMORY_DAM__Get__-12     254328781                4.440 ns/op           0 B/op          0 allocs/op
Benchmark__Random_SAVE_MEMORY_DAM__Set__-12     44422232                37.70 ns/op            0 B/op          0 allocs/op
Benchmark__Random_SAVE_MEMORY_DAM__Get__-12     37469788                36.44 ns/op            0 B/op          0 allocs/op
Benchmark__Linear_Builtin_Map__Set__-12         13001238               106.4 ns/op            57 B/op          0 allocs/op
Benchmark__Linear_Builtin_Map__Get__-12         35397393                41.38 ns/op            0 B/op          0 allocs/op
Benchmark__Random_Builtin_Map__Set__-12         13061778               116.3 ns/op            57 B/op          0 allocs/op
Benchmark__Random_Builtin_Map__Get__-12         32381409                46.74 ns/op            0 B/op          0 allocs/op
PASS
ok      github.com/nacioboi/go_dam/dam/tests    81.222s
```

Getting the speed improvement for random get:

```powershell
❯  46.74/11.64
4.01546391752577
```

And for linear get:

```powershell
❯ 41.38/1.515
27.3135313531353
````

### Your welcome golang maintainers, i will be submitting a PR to replace the built-in `map` with this package.

jk

## Quick Start

- First, you need to get the package:

```bash
go get github.com/nacioboi/go_dam/dam
```

- Code sample:

```go
package main

import (
	"fmt"
	"time"

	"github.com/nacioboi/go_dam/dam"
)

var _t uint64
var _start time.Time

func bench_linear_DAM_set(dam *dam.DAM[uint64, uint64], n uint64, do_print bool) {
	_start = time.Now()
	for i := uint64(0); i < n; i++ {
		dam.Set(i+1, i)
	}
	since := time.Since(_start)
	if do_print {
		fmt.Println("DAM Microseconds      ::: LINEAR SET :::", since.Microseconds())
	}
}

func bench_linear_DAM_get(dam *dam.DAM[uint64, uint64], n uint64, do_print bool) {
	_t = 0
	_start = time.Now()
	for i := uint64(0); i < n; i++ {
		res := dam.Get(i + 1)
		_t += res.Value
	}
	since := time.Since(_start)
	if do_print {
		fmt.Println("DAM Microseconds      ::: LINEAR GET :::", since.Microseconds())
		fmt.Println("Checksum:", _t)
	}
}

func main() {
	const n = 1024 * 1024

	// Create a new DAM...
	m := dam_map.New_DAM[uint64, uint64](n)

	// Benchmark...
	bench_linear_DAM_set(m, n, true)
	bench_linear_DAM_get(m, n, true)
}
```

## How to contribute

1. Fork the repository.
2. Clone the repository.
3. Make your changes.
4. Create a simple pull request.

Or alternatively, if you spot something simple, just create an issue.

**Help is always appreciated!**

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
