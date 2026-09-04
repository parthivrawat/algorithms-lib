# Algorithms Library

A comprehensive, zero-dependency collection of common algorithms for TypeScript. The package is built for both CommonJS and ESM, with full declaration files, and is ready for publication on npm.

## Installation

```bash
npm install algos-lib
```

## Quick Start

```typescript
import {
  quickSort,
  binarySearch,
  bfs,
  dijkstra,
  type AdjacencyList,
  type WeightedAdjacencyList,
} from 'algos-lib';

const sorted = quickSort([3, 1, 4, 1, 5, 9, 2, 6]);
console.log(binarySearch(sorted, 4)); // 3

const graph: AdjacencyList = { a: ['b', 'c'], b: ['d'], c: [], d: [] };
console.log(bfs(graph, 'a')); // ['a', 'b', 'c', 'd']

const weighted: WeightedAdjacencyList = {
  a: [['b', 1], ['c', 4]],
  b: [['c', 2], ['d', 5]],
  c: [['d', 1]],
  d: [],
};
const { distances } = dijkstra(weighted, 'a');
console.log(distances['d']); // 4

// Custom comparators let the generic sorts/searches order arbitrary types:
const people = [{ name: 'amy', age: 30 }, { name: 'bob', age: 25 }];
quickSort(people, (a, b) => a.age - b.age); // sorted by age
```

## Features

- **Zero runtime dependencies**: No external packages required
- **Comprehensive coverage**: Sorting, searching, graphs, dynamic programming, strings, greedy, and divide-and-conquer
- **Strongly typed**: Generic APIs with full TypeScript declarations
- **Custom comparators**: every generic sort/search accepts an optional `compareFn` for arbitrary element types
- **Typed error taxonomy**: `AlgorithmsError` subclasses (`InvalidInputError`, `InvalidGraphError`, `NegativeCycleError`, `ValueNotFoundError`, `EmptyInputError`) for precise error handling
- **Dual format**: CommonJS (`dist/index.js`) and ESM (`dist/index.mjs`) with `dist/index.d.ts`
- **Tested**: Vitest test suite covering common cases and edge cases

## Supported Algorithms

### Sorting

- `quickSort`: quick sort
- `mergeSort`: stable merge sort
- `heapSort`: binary max-heap sort
- `radixSort`: LSD radix sort for non-negative integers
- `nativeSort`: array sort using the built-in sort

### Searching

- `binarySearch`: classic binary search
- `interpolationSearch`: interpolation search for numeric data
- `jumpSearch`: jump search with block size `sqrt(n)`

### Graphs

- `bfs`: breadth-first search
- `dfs`: iterative depth-first search
- `dijkstra`: Dijkstra shortest paths (non-negative weights)
- `aStar`: A* shortest path with heuristic
- `bellmanFord`: shortest paths with negative-weight cycle detection

### Dynamic Programming

- `knapsack01`: 0/1 knapsack maximum value
- `longestCommonSubsequence`: LCS length
- `editDistance`: Levenshtein distance

### String Algorithms

- `kmpSearch`: Knuth-Morris-Pratt pattern matching
- `rabinKarpSearch`: rolling-hash pattern matching
- `boyerMooreSearch`: Boyer-Moore with bad-character and good-suffix rules; reports overlapping matches

### Greedy

- `activitySelection`: maximum compatible activities
- `fractionalKnapsack`: fractional knapsack maximum value
- `huffmanCoding`: optimal prefix codes

### Divide and Conquer

- `maxSubarray`: maximum subarray sum
- `countInversions`: inversion count via merge sort
- `fastPower`: exponentiation by squaring

## Development

```bash
npm install
npm test
npm run build
```

## License

MIT License. See [LICENSE](./LICENSE) for details.
