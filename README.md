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

TODO -- Add code sample

## How to contribute

1. Fork the repository.
2. Clone the repository.
3. Make your changes.
4. Create a simple pull request.

Or alternatively, if you spot something simple, just create an issue.

**Help is always appreciated!**

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
