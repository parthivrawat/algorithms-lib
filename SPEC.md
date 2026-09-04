# Algorithms Library — Cross-Language API Contract

This document pins the public contract for every algorithm in the four language packages (Python, TypeScript, Go, Rust). It exists so that parity drift between packages can be treated as a bug and so that users can port code between languages with predictable behavior.

## Scope

The 25 user-facing algorithms covered by this contract are:

- **Sorting**: `quick_sort`, `merge_sort`, `heap_sort`, `radix_sort`, `native_sort`
- **Searching**: `binary_search`, `interpolation_search`, `jump_search`
- **Graphs**: `bfs`, `dfs`, `dijkstra`, `a_star`, `bellman_ford`
- **Dynamic programming**: `knapsack_01`, `longest_common_subsequence`, `edit_distance`
- **Strings**: `kmp_search`, `rabin_karp_search`, `boyer_moore_search`
- **Greedy**: `activity_selection`, `fractional_knapsack`, `huffman_coding`
- **Divide and conquer**: `max_subarray`, `count_inversions`, `fast_power`

## Conventions

### Naming

- Python: `snake_case`.
- TypeScript: `camelCase`.
- Go: `PascalCase`.
- Rust: `snake_case`.

Canonical names in this document are `snake_case`. Implementations must expose the same algorithm under the idiom for their language.

### Not-found policy

| Language | Search miss | Graph/path miss | Comment |
|---|---|---|---|
| Python | return `-1` | raise `ValueNotFoundError` | Searchers keep the `-1` sentinel for backwards compatibility. |
| TypeScript | return `-1` | throw `ValueNotFoundError` | Same as Python. |
| Go | return `ErrNotFound` | return `ErrNotFound` | Error-based. |
| Rust | return `Err(Error::NotFound)` | return `Err(Error::NotFound)` | Result-based. |

A missing key in an adjacency list is an invalid graph, not a not-found condition, and must raise the invalid-graph error.

### Error taxonomy

| Concept | Python | TypeScript | Go | Rust |
|---|---|---|---|---|
| Base | `AlgorithmsError` | `AlgorithmsError` | n/a | n/a |
| Not found | `ValueNotFoundError` | `ValueNotFoundError` | `ErrNotFound` | `Error::NotFound` |
| Invalid input | `InvalidInputError` | `InvalidInputError` | `ErrInvalidInput` | `Error::InvalidInput` |
| Invalid graph | `InvalidGraphError` | `InvalidGraphError` | `ErrInvalidGraph` | `Error::InvalidGraph` |
| Negative cycle | `NegativeCycleError` | `NegativeCycleError` | `ErrNegativeCycle` | `Error::NegativeCycle` |

Each error must be distinguishable from the language's generic runtime errors so that callers can catch the library's own error types.

### Sortedness validation

- Python, TypeScript, and Go must explicitly validate sortedness where an algorithm requires a sorted input and raise the invalid-input error if the input is not sorted.
- Rust may use `debug_assert!` with `is_sorted` to keep release builds at `O(log n)` for searches; release builds treat unsorted input as undefined behavior per the caller contract.
- All packages must document the sortedness precondition.

### Sorts: mutating vs non-mutating

- `quick_sort`, `merge_sort`, `heap_sort`, `radix_sort`, and `native_sort` must return a **new** sorted collection and leave the original input untouched.
- Rust additionally exposes `quick_sort_in_place`, `merge_sort_in_place`, `heap_sort_in_place`, and `native_sort_in_place` that take `&mut [T]` and sort in place.
- Go and Python in-place variants are not required by this contract but may be added as separate functions.

### Comparator and key support

Every sort and every search that relies on ordering must accept an optional comparator so that the library works on structured data.

| Language | Parameter | Comparator signature |
|---|---|---|
| Python | `key` (keyword, after `sorted` convention) | `key: Callable[[T], Any]` |
| TypeScript | `compareFn` (positional, optional) | `(a: T, b: T) => number` |
| Go | `cmp` (positional, optional, after `sort` package convention) | `func(a, b T) int` |
| Rust | `compare` or `*_by` variant | `F: Fn(&T, &T) -> Ordering` |

`quick_sort`, `merge_sort`, `heap_sort`, `native_sort`, `binary_search`, `jump_search`, and `count_inversions` are in scope for comparator support.

### Graph edge representation

| Language | Unweighted graph | Weighted graph |
|---|---|---|
| Python | `Mapping[Any, Iterable[Any]]` | `Mapping[Any, Iterable[Tuple[Any, float]]]` |
| TypeScript | `Record<string, string[]>` | `Record<string, [string, number][]>` |
| Go | `map[string][]string` | `map[string][]Edge` where `Edge` is `{To string; Weight float64}` |
| Rust | `HashMap<String, Vec<String>>` | `HashMap<String, Vec<Edge>>` where `Edge { to: String, weight: f64 }` |

Python permits generic node keys; the other three packages use `string` node identifiers. The Python implementation must preserve generic keys while the other packages use string keys consistently.

### Return shapes

| Function family | Canonical return shape | Notes |
|---|---|---|
| Sorts | `sorted_collection` | New collection in ascending order by default or by comparator. |
| Searches | `index` or not-found sentinel/error | Byte/char offsets for string searchers. |
| BFS/DFS | `traversal_order` | List of nodes in visit order. |
| Dijkstra | `(distances, predecessors)` | Distances are numeric; predecessors are node identifiers. |
| Bellman-Ford | `(distances, predecessors)` | Same as Dijkstra; also reports reachable negative-weight cycles. |
| A* | `(path, cost)` | Path is a list of node identifiers; cost is numeric. |

## Function-by-function contract

| Function | Python | TypeScript | Go | Rust | Return semantics |
|---|---|---|---|---|---|
| `quick_sort` | `quick_sort(items: List[T]) -> List[T]` | `quickSort<T>(items: T[], compareFn?: Comparator<T>) -> T[]` | `QuickSort[T Ordered](items []T) []T` | `quick_sort<T: Ord + Clone>(items: &[T]) -> Vec<T>` | New ascending-sorted list. |
| `merge_sort` | `merge_sort(items: List[T]) -> List[T]` | `mergeSort<T>(items: T[], compareFn?: Comparator<T>) -> T[]` | `MergeSort[T Ordered](items []T) []T` | `merge_sort<T: Ord + Clone>(items: &[T]) -> Vec<T>` | New stable ascending-sorted list. |
| `heap_sort` | `heap_sort(items: List[T]) -> List[T]` | `heapSort<T>(items: T[], compareFn?: Comparator<T>) -> T[]` | `HeapSort[T Ordered](items []T) []T` | `heap_sort<T: Ord + Clone>(items: &[T]) -> Vec<T>` | New ascending-sorted list. |
| `radix_sort` | `radix_sort(items: List[int]) -> List[int]` | `radixSort(items: number[]) -> number[]` | `RadixSort(items []int) ([]int, error)` | `radix_sort(items: &[i64]) -> Result<Vec<i64>, Error>` | Sorted non-negative integers; invalid inputs raise/error. |
| `native_sort` | `native_sort(items: List[T]) -> List[T]` | `nativeSort<T>(items: T[], compareFn?: Comparator<T>) -> T[]` | `NativeSort[T Ordered](items []T) []T` | `native_sort<T: Ord + Clone>(items: &[T]) -> Vec<T>` | New list sorted via built-in sort; stability not guaranteed. |
| `binary_search` | `binary_search(arr: List[T], target: T) -> int` | `binarySearch<T>(arr: T[], target: T, compareFn?: Comparator<T>) -> number` | `BinarySearch[T Ordered](arr []T, target T) (int, error)` | `binary_search<T: Ord>(arr: &[T], target: &T) -> Result<usize, Error>` | Index on success; `-1` or not-found error when absent or empty; invalid input error on unsorted input. |
| `interpolation_search` | `interpolation_search(arr: List[int], target: int) -> int` | `interpolationSearch(arr: number[], target: number) -> number` | `InterpolationSearch(arr []int, target int) (int, error)` | `interpolation_search(arr: &[i64], target: i64) -> Result<usize, Error>` | Index on success; `-1` or not-found error when absent or empty; invalid input error on unsorted input. |
| `jump_search` | `jump_search(arr: List[T], target: T) -> int` | `jumpSearch<T>(arr: T[], target: T, compareFn?: Comparator<T>) -> number` | `JumpSearch[T Ordered](arr []T, target T) (int, error)` | `jump_search<T: Ord>(arr: &[T], target: &T) -> Result<usize, Error>` | Index on success; `-1` or not-found error when absent or empty; invalid input error on unsorted input. |
| `bfs` | `bfs(graph, start: Any) -> List[Any]` | `bfs(graph: AdjacencyList, start: string) -> string[]` | `BFS(graph map[string][]string, start string) ([]string, error)` | `bfs(graph, start: &str) -> Result<Vec<String>, Error>` | BFS traversal order; invalid graph error on missing start or neighbor. |
| `dfs` | `dfs(graph, start: Any) -> List[Any]` | `dfs(graph: AdjacencyList, start: string) -> string[]` | `DFS(graph map[string][]string, start string) ([]string, error)` | `dfs(graph, start: &str) -> Result<Vec<String>, Error>` | DFS traversal order; invalid graph error on missing start or neighbor. |
| `dijkstra` | `dijkstra(graph, start: Any) -> Tuple[dict, dict]` | `dijkstra(graph, start: string) -> { distances, predecessors }` | `Dijkstra(graph, start string) (map, map, error)` | `dijkstra(graph, start: &str) -> Result<(HashMap, HashMap), Error>` | Shortest distances and predecessor map; invalid graph error on negative weights. |
| `a_star` | `a_star(graph, start, goal, heuristic) -> Tuple[List, float]` | `aStar(graph, start, goal, heuristic) -> { path, cost }` | `AStar(graph, start, goal string, heuristic func) ([]string, float64, error)` | `a_star(graph, start, goal, heuristic) -> Result<(Vec<String>, f64), Error>` | Path and cost; invalid input error on negative heuristics; invalid graph error on negative weights. |
| `bellman_ford` | `bellman_ford(graph, start: Any) -> Tuple[dict, dict]` | `bellmanFord(graph, start: string) -> { distances, predecessors }` | `BellmanFord(graph, start string) (map, map, error)` | `bellman_ford(graph, start: &str) -> Result<(HashMap, HashMap), Error>` | Shortest distances and predecessor map; negative-cycle error if a reachable cycle exists. |
| `knapsack_01` | `knapsack_01(weights, values, capacity) -> int` | `knapsack01(weights, values, capacity) -> number` | `Knapsack01(weights, values []int, capacity int) (int, error)` | `knapsack_01(weights, values, capacity) -> Result<i64, Error>` | Maximum 0/1 knapsack value; invalid input error on length/negative/overflow. |
| `longest_common_subsequence` | `longest_common_subsequence(a, b) -> int` | `longestCommonSubsequence<T>(a, b) -> number` | `LongestCommonSubsequence[T comparable](a, b []T) int` | `longest_common_subsequence<T: PartialEq>(a, b) -> usize` | Length of LCS. |
| `edit_distance` | `edit_distance(a, b) -> int` | `editDistance<T>(a, b) -> number` | `EditDistance[T comparable](a, b []T) int` | `edit_distance<T: PartialEq>(a, b) -> usize` | Levenshtein distance. |
| `kmp_search` | `kmp_search(text, pattern) -> List[int]` | `kmpSearch(text, pattern) -> number[]` | `KMPSearch(text, pattern string) ([]int, error)` | `kmp_search(text, pattern) -> Result<Vec<usize>, Error>` | All starting byte offsets; invalid input error on empty pattern. |
| `rabin_karp_search` | `rabin_karp_search(text, pattern, base, mod) -> List[int]` | `rabinKarpSearch(text, pattern, base, mod) -> number[]` | `RabinKarpSearch(text, pattern string, base int, mod int64) ([]int, error)` | `rabin_karp_search(text, pattern, base, modulus) -> Result<Vec<usize>, Error>` | All starting byte offsets; invalid input error on empty pattern or non-positive base/mod. |
| `boyer_moore_search` | `boyer_moore_search(text, pattern) -> List[int]` | `boyerMooreSearch(text, pattern) -> number[]` | `BoyerMooreSearch(text, pattern string) ([]int, error)` | `boyer_moore_search(text, pattern) -> Result<Vec<usize>, Error>` | All starting byte offsets, including overlaps; invalid input error on empty pattern. |
| `activity_selection` | `activity_selection(activities) -> List[int]` | `activitySelection(activities) -> number[]` | `ActivitySelection(activities [][2]int) ([]int, error)` | `activity_selection(activities) -> Result<Vec<usize>, Error>` | Indices of a maximum compatible set; invalid input error if any start > end. |
| `fractional_knapsack` | `fractional_knapsack(weights, values, capacity) -> float` | `fractionalKnapsack(weights, values, capacity) -> number` | `FractionalKnapsack(weights, values []float64, capacity float64) (float64, error)` | `fractional_knapsack(weights, values, capacity) -> Result<f64, Error>` | Maximum fractional value; invalid input error on negative capacity or invalid weights/values. |
| `huffman_coding` | `huffman_coding(frequencies) -> Dict[str, str]` | `huffmanCoding(frequencies) -> Record<string, string>` | `HuffmanCoding(frequencies map[string]int) (map[string]string, error)` | `huffman_coding(frequencies) -> Result<HashMap<String, String>, Error>` | Prefix-free code table; invalid input error on empty or non-positive frequencies. |
| `max_subarray` | `max_subarray(arr) -> T` | `maxSubarray(arr) -> number` | `MaxSubarray(arr []int) (int, error)` | `max_subarray(arr) -> Result<i64, Error>` | Maximum contiguous subarray sum; invalid input error on empty input. |
| `count_inversions` | `count_inversions(arr) -> int` | `countInversions<T>(arr, compareFn?: Comparator<T>) -> number` | `CountInversions[T Ordered](arr []T) int64` | `count_inversions<T: Ord + Clone>(arr) -> usize` | Number of inversions. |
| `fast_power` | `fast_power(base, exponent) -> Union[int, float]` | `fastPower(base, exponent) -> number` | `FastPower(base float64, exponent int) (float64, error)` | `fast_power(base, exponent) -> Result<f64, Error>` | `base ** exponent`; invalid input error on `0` raised to a negative power; `0 ** 0 == 1`. |

## Known deviations

- **Comparator support**: TypeScript already accepts `compareFn` for sorts, `binarySearch`, `jumpSearch`, and `countInversions`. Python, Go, and Rust do not yet support custom comparators or keys. Bringing these packages in line with the contract above is the next parity task.
- **Graph node types**: Python supports arbitrary hashable keys; the other packages use `string` keys. The contract allows this deviation because it is idiomatic for Python and explicit in the other packages.
- **Search not-found**: Python and TypeScript return `-1`; Go and Rust return an error. This is a deliberate cross-language divergence because error handling is idiomatic in Go and Rust.

## Golden test vectors

The companion file `tests/golden_vectors.json` in this repository contains canonical inputs and expected outputs for the subset of algorithms that have a language-agnostic result. Each language's test suite should consume this file and verify that its implementation returns the canonical result (after unwrapping any language-specific `Result`/`error`). See `tests/README.md` for consumption examples.
