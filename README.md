# Algorithms Library

A polyglot, production-ready collection of common algorithms implemented for Python, TypeScript, Go, and Rust.

## Packages

| Language | Package | Registry |
|---|---|---|
| Python | `algorithms-lib` | [PyPI](https://pypi.org/project/algorithms-lib) |
| TypeScript | `algos-lib` | [npm](https://www.npmjs.com/package/algos-lib) |
| Go | `github.com/parthivrawat/algorithms-lib/go` | [pkg.go.dev](https://pkg.go.dev/github.com/parthivrawat/algorithms-lib/go) |
| Rust | `algorithms-lib` | [crates.io](https://crates.io/crates/algorithms-lib) |

## Overview

This repository provides consistent, minimal implementations of common algorithms across four major languages. Each implementation is zero-dependency, fully tested, and packaged for its respective registry.

## Implemented Algorithms

- **Sorting**: `quick_sort`, `merge_sort`, `heap_sort`, `radix_sort`, `tim_sort`
- **Searching**: `binary_search`, `interpolation_search`, `jump_search`
- **Graphs**: `bfs`, `dfs`, `dijkstra`, `a_star`, `bellman_ford`
- **Dynamic Programming**: `knapsack_01`, `longest_common_subsequence`, `edit_distance`
- **Strings**: `kmp_search`, `rabin_karp_search`, `boyer_moore_search`
- **Greedy**: `activity_selection`, `fractional_knapsack`, `huffman_coding`
- **Divide and Conquer**: `max_subarray`, `count_inversions`, `fast_power`

## Repository Layout

```
algorithms-lib/
├── python/     # PyPI package
├── typescript/ # npm package
├── go/         # Go module
├── rust/       # crates.io package
└── README.md   # This file
```

## Development

Each language directory contains its own build and test commands:

- **Python**: `pip install -e ".[dev]"` then `pytest`
- **TypeScript**: `npm install` then `npm test` and `npm run build`
- **Go**: `go test ./...` and `go build`
- **Rust**: `cargo test` and `cargo build`

## License

MIT License. See each language directory's `LICENSE` file for details.
