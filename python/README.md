# Algorithms Library

A comprehensive, zero-dependency collection of common algorithms for Python. The package is designed to be production-ready, strongly typed, and easy to use.

## Installation

```bash
pip install algorithms-lib
```

## Quick Start

```python
from algorithms_lib import quick_sort, merge_sort, binary_search, bfs, dijkstra

# Sorting
print(quick_sort([3, 1, 4, 1, 5, 9, 2, 6]))
# [1, 1, 2, 3, 4, 5, 6, 9]

# Searching
print(binary_search([1, 3, 5, 7, 9], 5))  # 2

# Graph traversal
graph = {'a': ['b', 'c'], 'b': ['d'], 'c': [], 'd': []}
print(bfs(graph, 'a'))  # ['a', 'b', 'c', 'd']

# Shortest path
weighted = {
    'a': [('b', 1), ('c', 4)],
    'b': [('c', 2), ('d', 5)],
    'c': [('d', 1)],
    'd': [],
}
distances, _ = dijkstra(weighted, 'a')
print(distances['d'])  # 4
```

## Features

- **Zero runtime dependencies**: No external packages required
- **Comprehensive coverage**: Sorting, searching, graphs, dynamic programming, strings, greedy, and divide-and-conquer
- **Type safe**: Includes `py.typed` marker for type checkers
- **Well tested**: Full test coverage for normal use, edge cases, and invalid operations
- **Clear error messages**: Explicit, helpful exceptions

## Supported Algorithms

### Sorting

- `quick_sort`: in-place-memory quick sort
- `merge_sort`: stable merge sort
- `heap_sort`: binary max-heap sort
- `radix_sort`: LSD radix sort for non-negative integers
- `tim_sort`: Python\'s Timsort wrapper

### Searching

- `binary_search`: classic binary search
- `interpolation_search`: interpolation search for numeric data
- `jump_search`: jump search with block size `sqrt(n)`

### Graphs

- `bfs`: breadth-first search
- `dfs`: iterative depth-first search
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
- `boyer_moore_search`: Boyer-Moore with bad-character rule

### Greedy

- `activity_selection`: maximum compatible activities
- `fractional_knapsack`: fractional knapsack maximum value
- `huffman_coding`: optimal prefix codes

### Divide and Conquer

- `max_subarray`: maximum subarray sum
- `count_inversions`: inversion count via merge sort
- `fast_power`: exponentiation by squaring

## Development

```bash
pip install -e ".[dev]"
pytest test_algorithms_lib.py -v
```

## License

MIT License. See [LICENSE](./LICENSE) for details.
