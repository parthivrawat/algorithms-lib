# Algorithms Library for Go

A comprehensive, zero-dependency collection of common algorithms for Go. The package is designed for production use and is ready for indexing on [pkg.go.dev](https://pkg.go.dev).

## Installation

```bash
go get github.com/parthivrawat/algorithms-lib/go
```

## Quick Start

```go
package main

import (
    "fmt"
    "github.com/parthivrawat/algorithms-lib/go"
)

func main() {
    sorted := algorithms.QuickSort([]int{3, 1, 4, 1, 5, 9, 2, 6})
    fmt.Println(sorted) // [1 1 2 3 4 5 6 9]

    idx, _ := algorithms.BinarySearch(sorted, 4)
    fmt.Println(idx) // 3

    graph := map[string][]algorithms.Edge{
        "a": {{To: "b", Weight: 1}, {To: "c", Weight: 4}},
        "b": {{To: "c", Weight: 2}, {To: "d", Weight: 5}},
        "c": {{To: "d", Weight: 1}},
        "d": {},
    }
    distances, _, _ := algorithms.Dijkstra(graph, "a")
    fmt.Println(distances["d"]) // 4
}
```

## Features

- **Zero runtime dependencies**: Uses only the Go standard library
- **Comprehensive coverage**: Sorting, searching, graphs, dynamic programming, strings, greedy, and divide-and-conquer
- **Generic APIs**: Type-safe with Go 1.21+ generics
- **Well tested**: Full `go test` coverage
- **Idiomatic errors**: Methods return `error` values where appropriate

## Supported Algorithms

### Sorting

- `QuickSort`: quick sort
- `MergeSort`: stable merge sort
- `HeapSort`: binary max-heap sort
- `RadixSort`: LSD radix sort for non-negative integers
- `TimSort`: standard-library sort wrapper

### Searching

- `BinarySearch`: classic binary search
- `InterpolationSearch`: interpolation search for numeric data
- `JumpSearch`: jump search with block size `sqrt(n)`

### Graphs

- `BFS`: breadth-first search
- `DFS`: depth-first search
- `Dijkstra`: Dijkstra shortest paths (non-negative weights)
- `AStar`: A* shortest path with heuristic
- `BellmanFord`: shortest paths with negative-weight cycle detection

### Dynamic Programming

- `Knapsack01`: 0/1 knapsack maximum value
- `LongestCommonSubsequence`: LCS length
- `EditDistance`: Levenshtein distance

### String Algorithms

- `KMPSearch`: Knuth-Morris-Pratt pattern matching
- `RabinKarpSearch`: rolling-hash pattern matching
- `BoyerMooreSearch`: Boyer-Moore with bad-character rule

### Greedy

- `ActivitySelection`: maximum compatible activities
- `FractionalKnapsack`: fractional knapsack maximum value
- `HuffmanCoding`: optimal prefix codes

### Divide and Conquer

- `MaxSubarray`: maximum subarray sum
- `CountInversions`: inversion count via merge sort
- `FastPower`: exponentiation by squaring

## Development

```bash
go test ./...
go build
```

## License

MIT License. See [LICENSE](./LICENSE) for details.
