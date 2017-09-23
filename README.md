# Keyvast

Kevast is a complex key value store back by block chain and uses complex machine learning for efficient reads and writes.
Just kidding Kevast is just a simple key value store with a REPL client.
This code is by no means production quality, just a quick solution for a code test.

## Benchmark

- OS: Mac OS 10.13 Beta (17A362a)
- CPU: 2.9 GHz Intel Core i7
- MEMORY: 16 GB 2133 MHz LPDDR3

```
BenchmarkKevast_NewDB/Basic-8           200000     106 ns/op      48 B/op     1 allocs/op
BenchmarkKevast_Write/1-8              1000000    18.8 ns/op       0 B/op     0 allocs/op
BenchmarkKevast_Write/10-8             1000000    16.3 ns/op       0 B/op     0 allocs/op
BenchmarkKevast_Write/100-8            1000000    18.1 ns/op       0 B/op     0 allocs/op
BenchmarkKevast_Read/1-8               1000000    19.7 ns/op       0 B/op     0 allocs/op
BenchmarkKevast_Read/10-8              1000000    21.8 ns/op       0 B/op     0 allocs/op
BenchmarkKevast_Read/100-8             1000000    18.8 ns/op       0 B/op     0 allocs/op
BenchmarkKevast_Start/Basic-8           200000     143 ns/op      92 B/op     1 allocs/op
BenchmarkKevast_Abort/1x1-8           20000000    1.90 ns/op       0 B/op     0 allocs/op
BenchmarkKevast_Commit/1x1-8            200000    86.9 ns/op       0 B/op     0 allocs/op
BenchmarkKevast_Commit/1x10-8            20000     579 ns/op       0 B/op     0 allocs/op
BenchmarkKevast_Commit/1x100-8            3000    6064 ns/op       0 B/op     0 allocs/op
```

## License
MIT
