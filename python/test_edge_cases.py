'''Edge-case tests for the algorithms library.'''

import pytest

from algorithms_lib import (
    a_star,
    binary_search,
    boyer_moore_search,
    bfs,
    count_inversions,
    dfs,
    dijkstra,
    fast_power,
    fractional_knapsack,
    heap_sort,
    huffman_coding,
    interpolation_search,
    jump_search,
    kmp_search,
    knapsack_01,
    longest_common_subsequence,
    max_subarray,
    merge_sort,
    native_sort,
    quick_sort,
    radix_sort,
    rabin_karp_search,
)


class TestEdgeCases:
    def test_empty_single_element_and_duplicate_sorts(self):
        cases = [[], [1], [2, 2, 2], [3, 1, 4, 1, 5, 9, 2, 6]]
        for data in cases:
            assert quick_sort(data) == sorted(data)
            assert merge_sort(data) == sorted(data)
            assert heap_sort(data) == sorted(data)
            assert native_sort(data) == sorted(data)

    def test_empty_single_and_duplicate_dc_dp(self):
        assert max_subarray([5]) == 5
        assert max_subarray([-5]) == -5
        assert count_inversions([]) == 0
        assert count_inversions([1]) == 0
        assert count_inversions([2, 2, 2]) == 0
        assert longest_common_subsequence('', 'abc') == 0
        assert longest_common_subsequence('A', 'A') == 1

    def test_radix_sort_empty_single_duplicates_and_extreme(self):
        assert radix_sort([]) == []
        assert radix_sort([42]) == [42]
        assert radix_sort([5, 5, 5, 1]) == [1, 5, 5, 5]
        assert radix_sort([2 ** 63 - 1, 0, 1]) == [0, 1, 2 ** 63 - 1]
        assert radix_sort([0, 10 ** 18, 10 ** 18, 0]) == [0, 0, 10 ** 18, 10 ** 18]
        with pytest.raises(ValueError):
            radix_sort([-1, 0])
        with pytest.raises(TypeError):
            radix_sort([1.0, 2])
        with pytest.raises(TypeError):
            radix_sort([True, 1])

    def test_searching_empty_single_and_duplicates(self):
        assert binary_search([], 1) == -1
        assert binary_search([5], 5) == 0
        assert binary_search([1, 2, 2, 3], 2) in (1, 2)
        assert binary_search([1, 2, 2, 3], 4) == -1

        assert interpolation_search([], 1) == -1
        assert interpolation_search([7], 7) == 0
        assert interpolation_search([1, 2, 2, 3], 2) in (1, 2)
        assert interpolation_search([1, 2, 3], 4) == -1

        assert jump_search([], 1) == -1
        assert jump_search([5], 5) == 0
        assert jump_search([1, 2, 2, 3], 2) in (1, 2)
        assert jump_search([1, 2, 3], 4) == -1

    def test_overlapping_string_matches(self):
        assert boyer_moore_search('aaaa', 'aa') == [0, 1, 2]
        assert boyer_moore_search('aaaaa', 'aaaa') == [0, 1]
        assert boyer_moore_search('abcabcabc', 'abcabc') == [0, 3]
        assert boyer_moore_search('abababab', 'abab') == [0, 2, 4]
        assert boyer_moore_search('abc', 'abcd') == []

        assert kmp_search('aaaa', 'aa') == [0, 1, 2]
        assert kmp_search('aaaaa', 'aaaa') == [0, 1]
        assert kmp_search('abcabcabc', 'abcabc') == [0, 3]
        assert kmp_search('abc', 'abcd') == []

        assert rabin_karp_search('aaaa', 'aa') == [0, 1, 2]
        assert rabin_karp_search('aaaaa', 'aaaa') == [0, 1]
        assert rabin_karp_search('abcabcabc', 'abcabc') == [0, 3]
        assert rabin_karp_search('abc', 'abcd') == []

    def test_disconnected_graphs(self):
        # Both components are valid adjacency-list keys; traversal must stay in the
        # start component.
        graph = {
            'a': ['b'],
            'b': [],
            'c': ['d'],
            'd': [],
        }
        assert bfs(graph, 'a') == ['a', 'b']
        assert dfs(graph, 'a') == ['a', 'b']

        weighted = {
            'a': [('b', 1)],
            'b': [],
            'c': [('d', 2)],
            'd': [],
        }
        distances, _ = dijkstra(weighted, 'a')
        assert distances['a'] == 0
        assert distances['b'] == 1
        assert distances['c'] == float('inf')
        assert distances['d'] == float('inf')

    def test_start_equals_goal(self):
        graph = {'a': []}
        distances, _ = dijkstra(graph, 'a')
        assert distances['a'] == 0

        path, cost = a_star(graph, 'a', 'a', lambda node, goal: 0)
        assert path == ['a']
        assert cost == 0.0

    def test_zero_capacity_knapsacks(self):
        assert knapsack_01([], [], 0) == 0
        assert knapsack_01([1, 2], [10, 20], 0) == 0
        assert knapsack_01([0, 10], [50, 60], 0) == 50
        assert fractional_knapsack([10, 20], [60, 100], 0) == 0.0
        assert fractional_knapsack([0, 10], [50, 60], 0) == 50.0

    def test_single_symbol_huffman(self):
        assert huffman_coding({'x': 5}) == {'x': '0'}

    def test_fast_power_extreme_exponents(self):
        with pytest.raises(ValueError):
            fast_power(0, -1)
        with pytest.raises(ValueError):
            fast_power(0, -5)

        # 64-bit MinInt-sized exponent should not crash and must respect sign.
        minint = -(2 ** 63)
        assert fast_power(1, minint) == 1.0
        assert fast_power(-1, minint) == 1.0

        assert fast_power(2, 0) == 1
        assert fast_power(0, 0) == 1
