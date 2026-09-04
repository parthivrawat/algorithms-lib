import { describe, it, expect } from 'vitest';
import {
  AdjacencyList,
  AlgorithmsError,
  aStar,
  activitySelection,
  bellmanFord,
  binarySearch,
  binarySearchOnAnswer,
  boyerMooreSearch,
  bfs,
  countInversions,
  dfs,
  dijkstra,
  editDistance,
  editDistanceReconstruction,
  EmptyInputError,
  fastPower,
  floydWarshall,
  fractionalKnapsack,
  gcd,
  hasCycle,
  heapSort,
  huffmanCoding,
  huffmanDecode,
  huffmanEncode,
  interpolationSearch,
  InvalidGraphError,
  InvalidInputError,
  jumpSearch,
  kmpSearch,
  knapsack01,
  lcm,
  longestCommonSubsequence,
  longestCommonSubsequenceReconstruction,
  maxSubarray,
  mergeSort,
  minimumSpanningTree,
  modularFastPower,
  NegativeCycleError,
  quickSelect,
  quickSort,
  nativeSort,
  radixSort,
  rabinKarpSearch,
  sieveOfEratosthenes,
  stronglyConnectedComponents,
  topologicalSort,
  UnionFind,
  ValueNotFoundError,
  PriorityQueue,
  WeightedAdjacencyList,
} from './index';

const sortingDatasets: (number | string)[][] = [
  [],
  [1],
  [3, 1, 4, 1, 5, 9, 2, 6],
  [9, 8, 7, 6, 5, 4, 3, 2, 1],
  ['cherry', 'apple', 'banana'],
];
const sortingCases = sortingDatasets.map((data) => [data]);

const ascending = (a: number | string, b: number | string): number =>
  a < b ? -1 : a > b ? 1 : 0;

describe('Sorting', () => {
  it.each(sortingCases)('quickSort sorts %j', (data) => {
    expect(quickSort(data)).toEqual([...data].sort(ascending));
  });

  it.each(sortingCases)('mergeSort sorts %j', (data) => {
    expect(mergeSort(data)).toEqual([...data].sort(ascending));
  });

  it.each(sortingCases)('heapSort sorts %j', (data) => {
    expect(heapSort(data)).toEqual([...data].sort(ascending));
  });

  it.each(sortingCases)('nativeSort sorts %j', (data) => {
    expect(nativeSort(data)).toEqual([...data].sort(ascending));
  });

  it('sorts support custom comparators on objects', () => {
    const people = [
      { name: 'amy', age: 30 },
      { name: 'bob', age: 25 },
      { name: 'cid', age: 40 },
    ];
    const byAge = (a: { age: number }, b: { age: number }): number =>
      a.age - b.age;
    for (const sort of [quickSort, mergeSort, heapSort, nativeSort]) {
      expect(sort(people, byAge).map((p) => p.name)).toEqual([
        'bob',
        'amy',
        'cid',
      ]);
    }
  });

  it('mergeSort with comparator is stable', () => {
    const items = [
      { k: 1, tag: 'a' },
      { k: 1, tag: 'b' },
      { k: 0, tag: 'c' },
    ];
    const sorted = mergeSort(items, (a, b) => a.k - b.k);
    expect(sorted.map((x) => x.tag)).toEqual(['c', 'a', 'b']);
  });

  it('radixSort sorts non-negative integers', () => {
    expect(radixSort([170, 45, 75, 90, 2, 802, 24, 66])).toEqual([
      2, 24, 45, 66, 75, 90, 170, 802,
    ]);
  });

  it('radixSort rejects negative integers', () => {
    expect(() => radixSort([1, -5, 3])).toThrow(AlgorithmsError);
  });

  it('radixSort rejects non-integer input', () => {
    expect(() => radixSort([1.5, 2, 3])).toThrow(AlgorithmsError);
    expect(() => radixSort([1, Number.NaN])).toThrow(AlgorithmsError);
    expect(() => radixSort([1, Number.MAX_SAFE_INTEGER + 1])).toThrow(AlgorithmsError);
  });

  it('radixSort sorts large values', () => {
    expect(radixSort([Number.MAX_SAFE_INTEGER, 1, 0])).toEqual([
      0,
      1,
      Number.MAX_SAFE_INTEGER,
    ]);
  });
});

describe('Searching', () => {
  const sortedNums = [2, 4, 6, 8, 10, 12];

  it('binarySearch finds existing value', () => {
    expect(binarySearch(sortedNums, 8)).toBe(3);
  });

  it('binarySearch returns -1 for missing value', () => {
    expect(binarySearch(sortedNums, 7)).toBe(-1);
  });

  it('binarySearch throws for unsorted input', () => {
    expect(() => binarySearch([3, 1, 2], 2)).toThrow(InvalidInputError);
  });

  it('binarySearch supports custom comparators', () => {
    const arr = [{ v: 1 }, { v: 3 }, { v: 5 }];
    const cmp = (a: { v: number }, b: { v: number }): number => a.v - b.v;
    expect(binarySearch(arr, { v: 3 }, cmp)).toBe(1);
    expect(binarySearch(arr, { v: 4 }, cmp)).toBe(-1);
  });

  it('interpolationSearch finds existing value', () => {
    expect(interpolationSearch(Array.from({ length: 10 }, (_, i) => (i + 1) * 10), 50)).toBe(4);
  });

  it('jumpSearch finds existing value', () => {
    expect(jumpSearch(sortedNums, 10)).toBe(4);
  });

  it('binarySearchOnAnswer finds boundary', () => {
    expect(binarySearchOnAnswer(1, 10, (x) => x >= 6)).toBe(6);
    expect(binarySearchOnAnswer(1, 10, (x) => x <= 4, 'maximum')).toBe(4);
  });
});

describe('Graphs', () => {
  it('bfs traverses in breadth-first order', () => {
    const graph = { a: ['b', 'c'], b: ['d'], c: [], d: [] };
    expect(bfs(graph, 'a')).toEqual(['a', 'b', 'c', 'd']);
  });

  it('bfs throws for missing start', () => {
    expect(() => bfs({ a: [] }, 'z')).toThrow(InvalidGraphError);
  });

  it('bfs/dfs reject neighbors that are not graph keys', () => {
    const bad: Record<string, string[]> = { a: ['b', 'zz'], b: [] };
    expect(() => bfs(bad, 'a')).toThrow(InvalidGraphError);
    expect(() => dfs(bad, 'a')).toThrow(InvalidGraphError);
  });

  it('dfs traverses in depth-first order', () => {
    const graph = { a: ['b', 'c'], b: ['d'], c: [], d: [] };
    expect(dfs(graph, 'a')).toEqual(['a', 'b', 'd', 'c']);
  });

  it('dijkstra finds shortest distances', () => {
    const graph: Record<string, [string, number][]> = {
      a: [['b', 1], ['c', 4]],
      b: [['c', 2], ['d', 5]],
      c: [['d', 1]],
      d: [],
    };
    const { distances } = dijkstra(graph, 'a');
    expect(distances['d']).toBe(4);
  });

  it('dijkstra rejects negative weights', () => {
    const graph: Record<string, [string, number][]> = {
      a: [['b', -1]],
      b: [],
    };
    expect(() => dijkstra(graph, 'a')).toThrow(InvalidGraphError);
  });

  it('dijkstra rejects negative weights in unreachable components', () => {
    const graph: Record<string, [string, number][]> = {
      a: [],
      b: [['c', -1]],
      c: [],
    };
    expect(() => dijkstra(graph, 'a')).toThrow(InvalidGraphError);
  });

  it('aStar finds a shortest path', () => {
    const graph: Record<string, [string, number][]> = {
      a: [['b', 1], ['c', 4]],
      b: [['c', 2], ['d', 5]],
      c: [['d', 1]],
      d: [],
    };
    const result = aStar(graph, 'a', 'd', () => 0);
    expect(result.path).toEqual(['a', 'b', 'c', 'd']);
    expect(result.cost).toBe(4);
  });

  it('aStar throws when no path exists', () => {
    const graph: Record<string, [string, number][]> = { a: [], b: [] };
    expect(() => aStar(graph, 'a', 'b', () => 0)).toThrow(ValueNotFoundError);
  });

  it('aStar rejects negative edge weights', () => {
    const graph: Record<string, [string, number][]> = { a: [['b', -1]], b: [] };
    expect(() => aStar(graph, 'a', 'b', () => 0)).toThrow(InvalidGraphError);
  });

  it('aStar rejects negative heuristic values', () => {
    const graph: Record<string, [string, number][]> = { a: [], b: [] };
    expect(() => aStar(graph, 'a', 'b', () => -1)).toThrow(AlgorithmsError);
  });

  it('bellmanFord handles negative weights', () => {
    const graph: Record<string, [string, number][]> = {
      a: [['b', -1]],
      b: [['c', -2]],
      c: [['d', 1]],
      d: [],
    };
    const { distances } = bellmanFord(graph, 'a');
    expect(distances['d']).toBe(-2);
  });

  it('bellmanFord detects negative cycles', () => {
    const graph: Record<string, [string, number][]> = {
      a: [['b', 1]],
      b: [['c', -1]],
      c: [['b', -1]],
    };
    expect(() => bellmanFord(graph, 'a')).toThrow(NegativeCycleError);
  });

  it('bellmanFord initializes unreachable vertices to Infinity', () => {
    const graph: Record<string, [string, number][]> = {
      a: [['b', 1]],
      b: [],
      z: [['y', 2]],
      y: [],
    };
    const { distances, predecessors } = bellmanFord(graph, 'a');
    expect(distances['z']).toBe(Infinity);
    expect(distances['y']).toBe(Infinity);
    expect(predecessors['z']).toBeUndefined();
  });
});

describe('Dynamic programming', () => {
  it('knapsack01 returns maximum value', () => {
    expect(knapsack01([1, 2, 3], [6, 10, 12], 5)).toBe(22);
  });

  it('knapsack01 rejects non-integer capacity and weights', () => {
    expect(() => knapsack01([1, 2], [6, 10], 2.5)).toThrow(InvalidInputError);
    expect(() => knapsack01([1.5, 2], [6, 10], 5)).toThrow(InvalidInputError);
    expect(() => knapsack01([1], [6], -1)).toThrow(InvalidInputError);
    expect(() => knapsack01([1, 2], [6], 5)).toThrow(InvalidInputError);
  });

  it('longestCommonSubsequence returns length', () => {
    expect(longestCommonSubsequence('ABCDE'.split(''), 'ACE'.split(''))).toBe(3);
  });

  it('editDistance returns Levenshtein distance', () => {
    expect(editDistance('kitten'.split(''), 'sitting'.split(''))).toBe(3);
  });

  it('editDistanceReconstruction returns distance and script', () => {
    const [distance, script] = editDistanceReconstruction(
      'kitten'.split(''),
      'sitting'.split('')
    );
    expect(distance).toBe(3);
    expect(script.length).toBeGreaterThanOrEqual(3);
  });

  it('longestCommonSubsequenceReconstruction returns length and subsequence', () => {
    const [length, lcs] = longestCommonSubsequenceReconstruction(
      'ABCDE'.split(''),
      'ACE'.split('')
    );
    expect(length).toBe(3);
    expect(lcs).toEqual(['A', 'C', 'E']);
  });
});

describe('String algorithms', () => {
  const text = 'ABABDABACDABABCABAB';
  const pattern = 'ABABCABAB';

  it('kmpSearch finds pattern', () => {
    expect(kmpSearch(text, pattern)).toEqual([10]);
  });

  it('rabinKarpSearch finds pattern', () => {
    expect(rabinKarpSearch(text, pattern)).toEqual([10]);
  });

  it('boyerMooreSearch finds pattern', () => {
    expect(boyerMooreSearch(text, pattern)).toEqual([10]);
  });

  it('boyerMooreSearch finds overlapping matches', () => {
    expect(boyerMooreSearch('aaaa', 'aa')).toEqual([0, 1, 2]);
    expect(boyerMooreSearch('aaaaa', 'aaaa')).toEqual([0, 1]);
    expect(boyerMooreSearch('abcabcabc', 'abcabc')).toEqual([0, 3]);
    expect(boyerMooreSearch('abababab', 'abab')).toEqual([0, 2, 4]);
  });
});

describe('Greedy', () => {
  it('activitySelection selects compatible activities', () => {
    const activities: [number, number][] = [
      [1, 3],
      [2, 5],
      [4, 7],
      [1, 8],
      [5, 9],
      [8, 10],
      [9, 11],
      [11, 14],
      [13, 16],
    ];
    const selected = activitySelection(activities);
    expect(selected.length).toBeGreaterThanOrEqual(4);
  });

  it('fractionalKnapsack returns maximum value', () => {
    expect(fractionalKnapsack([10, 20, 30], [60, 100, 120], 50)).toBe(240);
  });

  it('fractionalKnapsack takes zero-weight items with positive value in full', () => {
    expect(fractionalKnapsack([0, 10], [50, 60], 5)).toBe(80);
    expect(fractionalKnapsack([0], [5], 0)).toBe(5);
  });

  it('fractionalKnapsack rejects negative weights', () => {
    expect(() => fractionalKnapsack([-5, 10], [10, 60], 50)).toThrow();
  });

  it('fractionalKnapsack never takes non-positive-value items', () => {
    expect(fractionalKnapsack([10], [-5], 10)).toBe(0);
    expect(fractionalKnapsack([0, 10], [-50, 60], 50)).toBe(60);
  });

  it('huffmanCoding returns a code for every symbol', () => {
    const frequencies = { a: 5, b: 9, c: 12, d: 13, e: 16, f: 45 };
    const codes = huffmanCoding(frequencies);
    expect(Object.keys(codes).length).toBe(6);
    for (const code of Object.values(codes)) {
      expect(typeof code).toBe('string');
      expect(code.length).toBeGreaterThan(0);
    }
  });

  it('huffmanCoding handles a single symbol', () => {
    expect(huffmanCoding({ x: 3 })).toEqual({ x: '0' });
  });

  it('huffmanCoding rejects non-positive or non-finite frequencies', () => {
    expect(() => huffmanCoding({ a: 0 })).toThrow(InvalidInputError);
    expect(() => huffmanCoding({ a: -1 })).toThrow(InvalidInputError);
    expect(() => huffmanCoding({ a: Number.NaN })).toThrow(InvalidInputError);
    expect(() => huffmanCoding({ a: Infinity })).toThrow(InvalidInputError);
  });

  it('huffmanEncode and huffmanDecode are inverses', () => {
    const frequencies = { a: 5, b: 9, c: 12, d: 13, e: 16, f: 45 };
    const codes = huffmanCoding(frequencies);
    const symbols = ['a', 'b', 'c', 'd', 'e', 'f'];
    const encoded = huffmanEncode(symbols, codes);
    const decoded = huffmanDecode(encoded, codes);
    expect(decoded).toEqual(symbols);
  });

  it('huffmanEncode rejects unknown symbol', () => {
    expect(() => huffmanEncode(['x'], { a: '0' })).toThrow(InvalidInputError);
  });
});

describe('Divide and conquer', () => {
  it('maxSubarray finds maximum sum', () => {
    expect(maxSubarray([-2, 1, -3, 4, -1, 2, 1, -5, 4])).toBe(6);
  });

  it('maxSubarray works for all-negative input', () => {
    expect(maxSubarray([-2, -1])).toBe(-1);
  });

  it('countInversions counts inversions', () => {
    expect(countInversions([1, 3, 5, 2, 4, 6])).toBe(3);
  });

  it('fastPower computes powers', () => {
    expect(fastPower(2, 10)).toBe(1024);
    expect(fastPower(5, 0)).toBe(1);
    expect(fastPower(2, -2)).toBe(0.25);
  });

  it('fastPower handles edge cases', () => {
    expect(() => fastPower(0, -1)).toThrow(AlgorithmsError);
    expect(fastPower(0, 0)).toBe(1);
    expect(() => fastPower(2, 2.5)).toThrow(AlgorithmsError);
    expect(() => fastPower(2, Number.MAX_SAFE_INTEGER + 1)).toThrow(AlgorithmsError);
    expect(fastPower(-2, 3)).toBe(-8);
    expect(fastPower(-2, 4)).toBe(16);
    expect(fastPower(2, -10)).toBeCloseTo(2 ** -10, 15);
  });
});

describe('Number theory and selection', () => {
  it('gcd returns the greatest common divisor', () => {
    expect(gcd(48, 18)).toBe(6);
    expect(gcd(-48, 18)).toBe(6);
    expect(gcd(0, 5)).toBe(5);
    expect(gcd(0, 0)).toBe(0);
    expect(() => gcd(48.5, 18)).toThrow(InvalidInputError);
  });

  it('lcm returns the least common multiple', () => {
    expect(lcm(4, 6)).toBe(12);
    expect(lcm(-4, 6)).toBe(12);
    expect(lcm(0, 5)).toBe(0);
    expect(() => lcm(0, 0)).toThrow(InvalidInputError);
    expect(() => lcm(4.5, 6)).toThrow(InvalidInputError);
  });

  it('modularFastPower computes (base^exponent) % mod', () => {
    expect(modularFastPower(2, 10, 1000)).toBe(24);
    expect(modularFastPower(3, 0, 7)).toBe(1);
    expect(modularFastPower(0, 0, 7)).toBe(1);
    expect(modularFastPower(0, 5, 7)).toBe(0);
    expect(modularFastPower(-2, 3, 1000)).toBe(992);
    expect(() => modularFastPower(2, -1, 7)).toThrow(InvalidInputError);
    expect(() => modularFastPower(2, 2, 0)).toThrow(InvalidInputError);
  });

  it('quickSelect returns the k-th smallest element', () => {
    expect(quickSelect([3, 2, 1, 5, 4], 0)).toBe(1);
    expect(quickSelect([3, 2, 1, 5, 4], 2)).toBe(3);
    expect(quickSelect([3, 2, 1, 5, 4], 4)).toBe(5);
    expect(quickSelect(['cherry', 'apple', 'banana'], 1)).toBe('banana');
    expect(() => quickSelect([1, 2, 3], 5)).toThrow(InvalidInputError);
  });
});

describe('Number theory and graph extras', () => {
  it('sieveOfEratosthenes returns primes up to n', () => {
    expect(sieveOfEratosthenes(10)).toEqual([2, 3, 5, 7]);
    expect(sieveOfEratosthenes(1)).toEqual([]);
    expect(sieveOfEratosthenes(2)).toEqual([2]);
    expect(() => sieveOfEratosthenes(-1)).toThrow(InvalidInputError);
    expect(() => sieveOfEratosthenes(2.5)).toThrow(InvalidInputError);
  });

  it('topologicalSort returns a valid ordering', () => {
    const graph: AdjacencyList = {
      a: ['b', 'c'],
      b: ['d'],
      c: ['d'],
      d: [],
    };
    const order = topologicalSort(graph);
    expect(order.indexOf('a')).toBeLessThan(order.indexOf('b'));
    expect(order.indexOf('a')).toBeLessThan(order.indexOf('c'));
    expect(order.indexOf('b')).toBeLessThan(order.indexOf('d'));
    expect(order.indexOf('c')).toBeLessThan(order.indexOf('d'));
  });

  it('topologicalSort rejects cyclic graphs', () => {
    const graph: AdjacencyList = {
      a: ['b'],
      b: ['c'],
      c: ['a'],
    };
    expect(() => topologicalSort(graph)).toThrow(InvalidGraphError);
  });
});

describe('Data structures', () => {
  it('unionFind performs union and find', () => {
    const uf = new UnionFind(['a', 'b', 'c', 'd', 'e']);
    expect(uf.connected('a', 'b')).toBe(false);
    expect(uf.union('a', 'b')).toBe(true);
    expect(uf.connected('a', 'b')).toBe(true);
    uf.union('b', 'c');
    expect(uf.connected('a', 'c')).toBe(true);
    expect(uf.connected('a', 'd')).toBe(false);
    expect(uf.count()).toBe(3);
    expect(uf.union('a', 'c')).toBe(false);
  });

  it('unionFind rejects unknown elements', () => {
    const uf = new UnionFind(['a', 'b']);
    expect(() => uf.find('z')).toThrow(InvalidInputError);
  });

  it('priorityQueue returns items in order', () => {
    const pq = new PriorityQueue([3, 1, 4, 1, 5]);
    expect(pq.peek()).toBe(1);
    expect(pq.length).toBe(5);
    expect(pq.pop()).toBe(1);
    expect(pq.pop()).toBe(1);
    expect(pq.pop()).toBe(3);
    pq.push(0);
    expect(pq.pop()).toBe(0);
  });

  it('priorityQueue with custom comparator works', () => {
    const pq = new PriorityQueue([3, 1, 4, 1, 5], (a, b) => b - a);
    expect(pq.pop()).toBe(5);
    expect(pq.pop()).toBe(4);
  });

  it('priorityQueue pop on empty throws', () => {
    const pq = new PriorityQueue<number>();
    expect(() => pq.pop()).toThrow(EmptyInputError);
  });
});

describe('Graph algorithms', () => {
  it('minimumSpanningTree computes the MST', () => {
    const edges: [string, string, number][] = [
      ['a', 'b', 1],
      ['b', 'c', 2],
      ['a', 'c', 3],
      ['c', 'd', 4],
    ];
    const { total, mst } = minimumSpanningTree(edges);
    expect(total).toBe(7);
    expect(mst.length).toBe(3);
    expect(mst.reduce((sum, e) => sum + e[2], 0)).toBe(7);
  });

  it('floydWarshall computes all-pairs shortest paths', () => {
    const graph: WeightedAdjacencyList = {
      a: [['b', 1], ['c', 4]],
      b: [['c', 2], ['d', 5]],
      c: [['d', 1]],
      d: [],
    };
    const dist = floydWarshall(graph);
    expect(dist['a']['d']).toBe(4);
    expect(dist['a']['a']).toBe(0);
    expect(dist['b']['d']).toBe(3);
  });

  it('floydWarshall detects negative cycles', () => {
    const graph: WeightedAdjacencyList = {
      a: [['b', 1]],
      b: [['a', -2]],
    };
    expect(() => floydWarshall(graph)).toThrow(NegativeCycleError);
  });

  it('stronglyConnectedComponents groups components', () => {
    const graph: AdjacencyList = {
      a: ['b'],
      b: ['c', 'd'],
      c: ['a'],
      d: ['e'],
      e: [],
    };
    const components = stronglyConnectedComponents(graph);
    const normalized = components
      .map((c) => [...c].sort())
      .sort((a, b) => a.join('').localeCompare(b.join('')));
    expect(normalized).toEqual([['a', 'b', 'c'], ['d'], ['e']]);
  });

  it('hasCycle detects cycles', () => {
    const acyclic: AdjacencyList = { a: ['b'], b: [], c: ['d'], d: [] };
    expect(hasCycle(acyclic)).toBe(false);
    const cyclic: AdjacencyList = { a: ['b'], b: ['c'], c: ['a'] };
    expect(hasCycle(cyclic)).toBe(true);
  });
});
