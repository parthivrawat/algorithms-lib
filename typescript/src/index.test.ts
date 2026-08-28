import { describe, it, expect } from 'vitest';
import {
  aStar,
  activitySelection,
  bellmanFord,
  binarySearch,
  boyerMooreSearch,
  bfs,
  countInversions,
  dfs,
  dijkstra,
  editDistance,
  fastPower,
  fractionalKnapsack,
  heapSort,
  huffmanCoding,
  interpolationSearch,
  InvalidGraphError,
  jumpSearch,
  kmpSearch,
  knapsack01,
  longestCommonSubsequence,
  maxSubarray,
  mergeSort,
  NegativeCycleError,
  quickSort,
  radixSort,
  rabinKarpSearch,
  timSort,
  ValueNotFoundError,
} from './index';

const sortingDatasets: (number[] | string[])[] = [
  [],
  [1],
  [3, 1, 4, 1, 5, 9, 2, 6],
  [9, 8, 7, 6, 5, 4, 3, 2, 1],
  ['cherry', 'apple', 'banana'],
];
const sortingCases = sortingDatasets.map((data) => [data]);

describe('Sorting', () => {
  it.each(sortingCases)('quickSort sorts %j', (data) => {
    const sorted = [...data].sort((a, b) =>
      a < b ? -1 : a > b ? 1 : 0
    ) as typeof data;
    expect(quickSort(data as never)).toEqual(sorted);
  });

  it.each(sortingCases)('mergeSort sorts %j', (data) => {
    const sorted = [...data].sort((a, b) =>
      a < b ? -1 : a > b ? 1 : 0
    ) as typeof data;
    expect(mergeSort(data as never)).toEqual(sorted);
  });

  it.each(sortingCases)('heapSort sorts %j', (data) => {
    const sorted = [...data].sort((a, b) =>
      a < b ? -1 : a > b ? 1 : 0
    ) as typeof data;
    expect(heapSort(data as never)).toEqual(sorted);
  });

  it.each(sortingCases)('timSort sorts %j', (data) => {
    const sorted = [...data].sort((a, b) =>
      a < b ? -1 : a > b ? 1 : 0
    ) as typeof data;
    expect(timSort(data as never)).toEqual(sorted);
  });

  it('radixSort sorts non-negative integers', () => {
    expect(radixSort([170, 45, 75, 90, 2, 802, 24, 66])).toEqual([
      2, 24, 45, 66, 75, 90, 170, 802,
    ]);
  });

  it('radixSort rejects negative integers', () => {
    expect(() => radixSort([1, -5, 3])).toThrow();
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
    expect(() => binarySearch([3, 1, 2], 2)).toThrow();
  });

  it('interpolationSearch finds existing value', () => {
    expect(interpolationSearch(Array.from({ length: 10 }, (_, i) => (i + 1) * 10), 50)).toBe(4);
  });

  it('jumpSearch finds existing value', () => {
    expect(jumpSearch(sortedNums, 10)).toBe(4);
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
});

describe('Dynamic programming', () => {
  it('knapsack01 returns maximum value', () => {
    expect(knapsack01([1, 2, 3], [6, 10, 12], 5)).toBe(22);
  });

  it('longestCommonSubsequence returns length', () => {
    expect(longestCommonSubsequence('ABCDE'.split(''), 'ACE'.split(''))).toBe(3);
  });

  it('editDistance returns Levenshtein distance', () => {
    expect(editDistance('kitten'.split(''), 'sitting'.split(''))).toBe(3);
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

  it('huffmanCoding returns a code for every symbol', () => {
    const frequencies = { a: 5, b: 9, c: 12, d: 13, e: 16, f: 45 };
    const codes = huffmanCoding(frequencies);
    expect(Object.keys(codes).length).toBe(6);
    for (const code of Object.values(codes)) {
      expect(typeof code).toBe('string');
      expect(code.length).toBeGreaterThan(0);
    }
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
});
