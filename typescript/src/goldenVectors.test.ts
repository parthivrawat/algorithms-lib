import { describe, expect, it } from 'vitest';
import { readFileSync } from 'node:fs';
import {
  quickSort,
  mergeSort,
  heapSort,
  radixSort,
  nativeSort,
  binarySearch,
  interpolationSearch,
  jumpSearch,
  kmpSearch,
  rabinKarpSearch,
  boyerMooreSearch,
  longestCommonSubsequence,
  editDistance,
  knapsack01,
  fractionalKnapsack,
  activitySelection,
  maxSubarray,
  countInversions,
  fastPower,
} from './index';

interface Vector {
  function: string;
  name: string;
  input: Record<string, any>;
  expected: any;
}

const vectors: { vectors: Vector[] } = JSON.parse(
  readFileSync(
    new URL('../../tests/golden_vectors.json', import.meta.url),
    'utf8'
  )
);

const dispatch: Record<string, (input: Record<string, any>) => any> = {
  quick_sort: ({ items }) => quickSort(items),
  merge_sort: ({ items }) => mergeSort(items),
  heap_sort: ({ items }) => heapSort(items),
  radix_sort: ({ items }) => radixSort(items),
  native_sort: ({ items }) => nativeSort(items),
  binary_search: ({ arr, target }) => binarySearch(arr, target),
  interpolation_search: ({ arr, target }) => interpolationSearch(arr, target),
  jump_search: ({ arr, target }) => jumpSearch(arr, target),
  kmp_search: ({ text, pattern }) => kmpSearch(text, pattern),
  rabin_karp_search: ({ text, pattern, base, mod }) =>
    rabinKarpSearch(text, pattern, base, mod),
  boyer_moore_search: ({ text, pattern }) => boyerMooreSearch(text, pattern),
  longest_common_subsequence: ({ a, b }) => longestCommonSubsequence(a, b),
  edit_distance: ({ a, b }) => editDistance(a, b),
  knapsack_01: ({ weights, values, capacity }) =>
    knapsack01(weights, values, capacity),
  fractional_knapsack: ({ weights, values, capacity }) =>
    fractionalKnapsack(weights, values, capacity),
  activity_selection: ({ activities }) => activitySelection(activities),
  max_subarray: ({ arr }) => maxSubarray(arr),
  count_inversions: ({ arr }) => countInversions(arr),
  fast_power: ({ base, exponent }) => fastPower(base, exponent),
};

describe('golden vectors', () => {
  for (const v of vectors.vectors) {
    it(v.name, () => {
      const fn = dispatch[v.function];
      if (!fn) {
        throw new Error(`No TypeScript dispatch for ${v.function}`);
      }
      expect(fn(v.input)).toEqual(v.expected);
    });
  }
});
