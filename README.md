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
❯ go run .\main.go
Built-in Map Microseconds  ::: LINEAR SET ::: 241689
Built-in Map Microseconds  ::: LINEAR GET ::: 92691
Checksum: 2199022206976
DAM Microseconds      ::: LINEAR SET ::: 92462
DAM Microseconds      ::: LINEAR GET ::: 72458
Checksum: 2199022206976

Built-in Map Microseconds  ::: RANDOM GET PER OP ::: 0.044669
DAM Microseconds      ::: RANDOM GET PER OP ::: 0.088233

Builtin Map Microseconds   ::: DELETE ::: 138280
DAM Microseconds      ::: DELETE ::: 179705

Memory Used (Built-in):  41,064 bytes
Memory Used (DAM):      85,976 bytes
```

```text
❯ 92691/72458
1.27923762731513
```

That is 1.2 times faster than the built-in `map` for getting values.

And it only consumes roughly 2x more memory than the built-in `map`.
A worth-wile trade-off for the speed.

> NOTE: The memory usage is calculated using a smaller size of map compared to the performance tests.

But you might have noticed the random get per operation tells a different story...
This is because we're using the slowest performance profile.

**When we change to the fastest performance profile, we get:**

```text
❯ go run .\main.go
Built-in Map Microseconds  ::: LINEAR SET ::: 243346
Built-in Map Microseconds  ::: LINEAR GET ::: 92668
Checksum: 2199022206976
DAM Microseconds      ::: LINEAR SET ::: 15104
DAM Microseconds      ::: LINEAR GET ::: 3632
Checksum: 2199022206976

Built-in Map Microseconds  ::: RANDOM GET PER OP ::: 0.044980
DAM Microseconds      ::: RANDOM GET PER OP ::: 0.019487

Builtin Map Microseconds   ::: DELETE ::: 139271
DAM Microseconds      ::: DELETE ::: 112769

Memory Used (Built-in):  41,064 bytes
Memory Used (DAM):      85,976 bytes
```

```text
❯ 92668/3632
25.5143171806167
```

That is 25 times faster than the built-in `map` for getting values in our linear test.

Let us divide the time taken in random get from the built-in `map` and our DAM respectively:

```text
❯ 0.044980/0.019487
2.30820547031354
```

This gives us the number to beat.
The performance of our DAM can only get better from here!

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
