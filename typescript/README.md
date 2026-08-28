# Algorithms Library

A comprehensive, zero-dependency collection of common algorithms for TypeScript. The package is built for both CommonJS and ESM, with full declaration files, and is ready for publication on npm.

## Installation

```bash
npm install algos-lib
```

## Quick Start

```typescript
import { quickSort, binarySearch, bfs, dijkstra } from 'algos-lib';

const sorted = quickSort([3, 1, 4, 1, 5, 9, 2, 6]);
console.log(binarySearch(sorted, 4)); // 3

const graph = { a: ['b', 'c'], b: ['d'], c: [], d: [] };
console.log(bfs(graph, 'a')); // ['a', 'b', 'c', 'd']

const weighted = {
  a: [['b', 1] as [string, number], ['c', 4] as [string, number]],
  b: [['c', 2] as [string, number], ['d', 5] as [string, number]],
  c: [['d', 1] as [string, number]],
  d: [],
};
const { distances } = dijkstra(weighted, 'a');
console.log(distances['d']); // 4
```

## Features

- **Zero runtime dependencies**: No external packages required
- **Comprehensive coverage**: Sorting, searching, graphs, dynamic programming, strings, greedy, and divide-and-conquer
- **Strongly typed**: Generic APIs with full TypeScript declarations
- **Dual format**: CommonJS (`dist/index.js`) and ESM (`dist/index.mjs`) with `dist/index.d.ts`
- **Well tested**: Vitest test suite covering common and edge cases

## Supported Algorithms

### Sorting

- `quickSort`: quick sort
- `mergeSort`: stable merge sort
- `heapSort`: binary max-heap sort
- `radixSort`: LSD radix sort for non-negative integers
- `timSort`: array sort using JavaScript's built-in Timsort

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
- `boyerMooreSearch`: Boyer-Moore with bad-character rule

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
