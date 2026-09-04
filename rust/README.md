# Algorithms Library for Rust

A comprehensive, zero-dependency collection of common algorithms for Rust. The crate is published on [crates.io](https://crates.io/crates/algorithms-lib) and documented on [docs.rs](https://docs.rs/algorithms-lib).

## Installation

Add the crate to your `Cargo.toml`:

```toml
[dependencies]
algorithms-lib = "2.0.0"
```

## Quick Start

```rust
use algorithms_lib::{quick_sort, binary_search, bfs, dijkstra, Edge};
use std::collections::HashMap;

let sorted = quick_sort(&[3, 1, 4, 1, 5, 9, 2, 6]);
assert_eq!(sorted, vec![1, 1, 2, 3, 4, 5, 6, 9]);

let idx = binary_search(&sorted, &4).unwrap();
assert_eq!(idx, 3);

let mut graph: HashMap<String, Vec<String>> = HashMap::new();
graph.insert("a".to_string(), vec!["b".to_string(), "c".to_string()]);
graph.insert("b".to_string(), vec!["d".to_string()]);
graph.insert("c".to_string(), vec![]);
graph.insert("d".to_string(), vec![]);
assert_eq!(bfs(&graph, "a").unwrap(), vec!["a", "b", "c", "d"]);

let mut weighted: HashMap<String, Vec<Edge>> = HashMap::new();
weighted.insert("a".to_string(), vec![Edge::new("b", 1.0), Edge::new("c", 4.0)]);
weighted.insert("b".to_string(), vec![Edge::new("c", 2.0), Edge::new("d", 5.0)]);
weighted.insert("c".to_string(), vec![Edge::new("d", 1.0)]);
weighted.insert("d".to_string(), vec![]);
let (distances, _) = dijkstra(&weighted, "a").unwrap();
assert_eq!(distances["d"], 4.0);
```

## Features

- **Zero runtime dependencies**: Uses only the Rust standard library
- **Comprehensive coverage**: Sorting, searching, graphs, dynamic programming, strings, greedy, and divide-and-conquer
- **Generic APIs**: Strongly-typed with trait bounds
- **Well tested**: Integration tests for common and edge cases
- **Clear errors**: Methods return `Result<T, Error>` where appropriate

## Supported Algorithms

### Sorting

- `quick_sort`: quick sort
- `merge_sort`: stable merge sort
- `heap_sort`: binary max-heap sort
- `radix_sort`: LSD radix sort for non-negative integers
- `native_sort`: standard library sort wrapper

### Searching

- `binary_search`: classic binary search
- `interpolation_search`: interpolation search for numeric data
- `jump_search`: jump search with block size `sqrt(n)`

### Graphs

- `bfs`: breadth-first search
- `dfs`: depth-first search
- `dijkstra`: Dijkstra shortest paths (non-negative weights)
- `a_star`: A* shortest path with heuristic
- `bellman_ford`: shortest paths with negative-weight cycle detection

### Dynamic Programming

- `knapsack_01`: 0/1 knapsack maximum value
- `longest_common_subsequence`: LCS length
- `edit_distance`: Levenshtein distance

### String Algorithms

- `kmp_search`: Knuth-Morris-Pratt pattern matching
- `rabin_karp_search`: rolling-hash pattern matching
- `boyer_moore_search`: Boyer-Moore with bad-character and good-suffix rules; reports overlapping matches (byte offsets)

### Greedy

- `activity_selection`: maximum compatible activities
- `fractional_knapsack`: fractional knapsack maximum value
- `huffman_coding`: optimal prefix codes

### Divide and Conquer

- `max_subarray`: maximum subarray sum
- `count_inversions`: inversion count via merge sort
- `fast_power`: exponentiation by squaring (returns `Err` for `0` to a negative power)

## Development

```bash
cargo test
cargo build
```

## License

MIT License. See [LICENSE](./LICENSE) for details.
