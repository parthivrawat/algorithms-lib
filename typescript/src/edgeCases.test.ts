import { describe, it, expect } from 'vitest';
import {
  aStar,
  bfs,
  boyerMooreSearch,
  dfs,
  dijkstra,
  fastPower,
  fractionalKnapsack,
  huffmanCoding,
  InvalidInputError,
  jumpSearch,
  kmpSearch,
  knapsack01,
  longestCommonSubsequence,
  editDistance,
  maxSubarray,
  radixSort,
  rabinKarpSearch,
  ValueNotFoundError,
  quickSort,
  mergeSort,
  heapSort,
  nativeSort,
  binarySearch,
  interpolationSearch,
  countInversions,
  AdjacencyList,
  WeightedAdjacencyList,
} from './index';

describe('Edge cases', () => {
  describe('empty and single-element inputs', () => {
    it('sorting functions return empty arrays unchanged', () => {
      const empty: number[] = [];
      expect(quickSort(empty)).toEqual([]);
      expect(mergeSort(empty)).toEqual([]);
      expect(heapSort(empty)).toEqual([]);
      expect(nativeSort(empty)).toEqual([]);
      expect(radixSort([])).toEqual([]);
    });

    it('sorting functions handle single-element arrays', () => {
      const data = [42];
      expect(quickSort(data)).toEqual([42]);
      expect(mergeSort(data)).toEqual([42]);
      expect(heapSort(data)).toEqual([42]);
      expect(nativeSort(data)).toEqual([42]);
      expect(radixSort([42])).toEqual([42]);
    });

    it('searches return -1 on empty input', () => {
      expect(binarySearch([], 5)).toBe(-1);
      expect(interpolationSearch([], 5)).toBe(-1);
      expect(jumpSearch([], 5)).toBe(-1);
    });

    it('bfs and dfs handle a single isolated node', () => {
      const graph: AdjacencyList = { a: [] };
      expect(bfs(graph, 'a')).toEqual(['a']);
      expect(dfs(graph, 'a')).toEqual(['a']);
    });

    it('knapsack functions handle empty inputs and zero capacity', () => {
      expect(knapsack01([], [], 0)).toBe(0);
      expect(knapsack01([1, 2], [10, 20], 0)).toBe(0);
      expect(fractionalKnapsack([], [], 0)).toBe(0);
      expect(fractionalKnapsack([1, 2], [10, 20], 0)).toBe(0);
    });

    it('lcs and editDistance handle empty inputs', () => {
      expect(longestCommonSubsequence([], [])).toBe(0);
      expect(longestCommonSubsequence([1, 2], [])).toBe(0);
      expect(editDistance([], [])).toBe(0);
      expect(editDistance(['a'], [])).toBe(1);
    });

    it('maxSubarray on a single element', () => {
      expect(maxSubarray([7])).toBe(7);
      expect(maxSubarray([-4])).toBe(-4);
    });
  });

  describe('duplicates', () => {
    it('sorting functions preserve duplicate values', () => {
      const data = [3, 1, 3, 2, 1];
      const expected = [1, 1, 2, 3, 3];
      expect(quickSort(data)).toEqual(expected);
      expect(mergeSort(data)).toEqual(expected);
      expect(heapSort(data)).toEqual(expected);
      expect(nativeSort(data)).toEqual(expected);
      expect(radixSort(data)).toEqual(expected);
    });

    it('countInversions handles duplicates correctly', () => {
      expect(countInversions([1, 2, 2, 3])).toBe(0);
      expect(countInversions([2, 2, 1])).toBe(2);
    });
  });

  describe('overlapping string matches', () => {
    it('kmpSearch reports overlapping matches', () => {
      expect(kmpSearch('aaaa', 'aa')).toEqual([0, 1, 2]);
      expect(kmpSearch('abababab', 'abab')).toEqual([0, 2, 4]);
    });

    it('rabinKarpSearch reports overlapping matches', () => {
      expect(rabinKarpSearch('aaaa', 'aa')).toEqual([0, 1, 2]);
      expect(rabinKarpSearch('abcabcabc', 'abcabc')).toEqual([0, 3]);
    });

    it('boyerMooreSearch reports overlapping matches', () => {
      expect(boyerMooreSearch('aaaa', 'aa')).toEqual([0, 1, 2]);
      expect(boyerMooreSearch('abcabcabc', 'abcabc')).toEqual([0, 3]);
    });
  });

  describe('disconnected graphs', () => {
    it('bfs and dfs stay in the start component', () => {
      const graph: AdjacencyList = { a: ['b'], b: [], c: ['d'], d: [] };
      expect(bfs(graph, 'a')).toEqual(['a', 'b']);
      expect(dfs(graph, 'a')).toEqual(['a', 'b']);
    });

    it('dijkstra leaves disconnected vertices at Infinity', () => {
      const graph: WeightedAdjacencyList = { a: [['b', 1]], b: [], c: [['d', 2]], d: [] };
      const { distances } = dijkstra(graph, 'a');
      expect(distances.a).toBe(0);
      expect(distances.b).toBe(1);
      expect(distances.c).toBe(Infinity);
      expect(distances.d).toBe(Infinity);
    });

    it('aStar throws when there is no path between components', () => {
      const graph: WeightedAdjacencyList = { a: [], b: [] };
      expect(() => aStar(graph, 'a', 'b', () => 0)).toThrow(ValueNotFoundError);
    });
  });

  describe('start equals goal', () => {
    it('dijkstra with start == goal', () => {
      const graph: WeightedAdjacencyList = { a: [] };
      const { distances, predecessors } = dijkstra(graph, 'a');
      expect(distances.a).toBe(0);
      expect(predecessors.a).toBeUndefined();
    });

    it('aStar with start == goal', () => {
      const graph: WeightedAdjacencyList = { a: [] };
      const result = aStar(graph, 'a', 'a', () => 0);
      expect(result.path).toEqual(['a']);
      expect(result.cost).toBe(0);
    });
  });

  describe('zero capacity and zero-weight items', () => {
    it('knapsack01 returns 0 for positive-weight items with zero capacity', () => {
      expect(knapsack01([1, 2], [10, 20], 0)).toBe(0);
    });

    it('knapsack01 still takes zero-weight items with zero capacity', () => {
      expect(knapsack01([0], [5], 0)).toBe(5);
    });

    it('fractionalKnapsack with zero capacity', () => {
      expect(fractionalKnapsack([1, 2], [10, 20], 0)).toBe(0);
      expect(fractionalKnapsack([0], [5], 0)).toBe(5);
    });
  });

  describe('single-symbol Huffman', () => {
    it('huffmanCoding returns a single code for one symbol', () => {
      expect(huffmanCoding({ x: 3 })).toEqual({ x: '0' });
    });
  });

  describe('extreme numeric inputs', () => {
    it('fastPower rejects 0 raised to a negative exponent', () => {
      expect(() => fastPower(0, -2)).toThrow(InvalidInputError);
      expect(() => fastPower(0, -1)).toThrow(InvalidInputError);
    });

    it('fastPower with Number.MAX_SAFE_INTEGER exponents', () => {
      expect(fastPower(1, Number.MAX_SAFE_INTEGER)).toBe(1);
      expect(fastPower(2, Number.MAX_SAFE_INTEGER)).toBe(Infinity);
      expect(fastPower(-1, Number.MAX_SAFE_INTEGER)).toBe(-1);
      expect(() => fastPower(2, Number.MAX_SAFE_INTEGER + 1)).toThrow(
        InvalidInputError
      );
    });

    it('radixSort handles and rejects extreme values', () => {
      expect(radixSort([Number.MAX_SAFE_INTEGER, 0, 1])).toEqual([
        0,
        1,
        Number.MAX_SAFE_INTEGER,
      ]);
      expect(() => radixSort([Number.MAX_SAFE_INTEGER + 1])).toThrow(
        InvalidInputError
      );
      expect(() => radixSort([-1])).toThrow(InvalidInputError);
      expect(() => radixSort([1.5])).toThrow(InvalidInputError);
    });
  });
});
