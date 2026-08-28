'''Tests for the algorithms library.'''

import pytest

from algorithms_lib import (
    AlgorithmsError,
    NegativeCycleError,
    InvalidGraphError,
    ValueNotFoundError,
    a_star,
    activity_selection,
    bellman_ford,
    binary_search,
    boyer_moore_search,
    bfs,
    count_inversions,
    dfs,
    dijkstra,
    edit_distance,
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
    quick_sort,
    radix_sort,
    rabin_karp_search,
    tim_sort,
)


class TestSorting:
    datasets = [
        [],
        [1],
        [3, 1, 4, 1, 5, 9, 2, 6],
        [9, 8, 7, 6, 5, 4, 3, 2, 1],
        ['cherry', 'apple', 'banana'],
    ]

    @pytest.mark.parametrize('data', datasets)
    def test_quick_sort(self, data):
        assert quick_sort(data) == sorted(data)

    @pytest.mark.parametrize('data', datasets)
    def test_merge_sort(self, data):
        assert merge_sort(data) == sorted(data)

    @pytest.mark.parametrize('data', datasets)
    def test_heap_sort(self, data):
        assert heap_sort(data) == sorted(data)

    @pytest.mark.parametrize('data', datasets)
    def test_tim_sort(self, data):
        assert tim_sort(data) == sorted(data)

    def test_radix_sort(self):
        assert radix_sort([170, 45, 75, 90, 2, 802, 24, 66]) == [
            2, 24, 45, 66, 75, 90, 170, 802
        ]

    def test_radix_sort_empty(self):
        assert radix_sort([]) == []

    def test_radix_sort_negative_raises(self):
        with pytest.raises(ValueError):
            radix_sort([1, -5, 3])


class TestSearching:
    sorted_nums = [2, 4, 6, 8, 10, 12]

    def test_binary_search_found(self):
        assert binary_search(self.sorted_nums, 8) == 3

    def test_binary_search_not_found(self):
        assert binary_search(self.sorted_nums, 7) == -1

    def test_binary_search_empty(self):
        assert binary_search([], 1) == -1

    def test_binary_search_unsorted_raises(self):
        with pytest.raises(ValueError):
            binary_search([3, 1, 2], 2)

    def test_interpolation_search_found(self):
        assert interpolation_search(list(range(10, 110, 10)), 50) == 4

    def test_interpolation_search_not_found(self):
        assert interpolation_search(list(range(10, 110, 10)), 55) == -1

    def test_jump_search_found(self):
        assert jump_search(self.sorted_nums, 10) == 4

    def test_jump_search_not_found(self):
        assert jump_search(self.sorted_nums, 7) == -1


class TestGraphs:
    def test_bfs_basic(self):
        graph = {'a': ['b', 'c'], 'b': ['d'], 'c': [], 'd': []}
        assert bfs(graph, 'a') == ['a', 'b', 'c', 'd']

    def test_bfs_missing_start_raises(self):
        with pytest.raises(InvalidGraphError):
            bfs({'a': []}, 'z')

    def test_dfs_basic(self):
        graph = {'a': ['b', 'c'], 'b': ['d'], 'c': [], 'd': []}
        assert dfs(graph, 'a') == ['a', 'b', 'd', 'c']

    def test_dijkstra_basic(self):
        graph = {
            'a': [('b', 1), ('c', 4)],
            'b': [('c', 2), ('d', 5)],
            'c': [('d', 1)],
            'd': [],
        }
        distances, _ = dijkstra(graph, 'a')
        assert distances['d'] == 4

    def test_dijkstra_negative_weight_raises(self):
        graph = {'a': [('b', -1)], 'b': []}
        with pytest.raises(InvalidGraphError):
            dijkstra(graph, 'a')

    def test_a_star_basic(self):
        graph = {
            'a': [('b', 1), ('c', 4)],
            'b': [('c', 2), ('d', 5)],
            'c': [('d', 1)],
            'd': [],
        }
        path, cost = a_star(graph, 'a', 'd', lambda node, goal: 0)
        assert path == ['a', 'b', 'c', 'd']
        assert cost == 4

    def test_a_star_no_path_raises(self):
        graph = {'a': [], 'b': []}
        with pytest.raises(ValueError):
            a_star(graph, 'a', 'b', lambda node, goal: 0)

    def test_bellman_ford_negative_weights(self):
        graph = {
            'a': [('b', -1)],
            'b': [('c', -2)],
            'c': [('d', 1)],
            'd': [],
        }
        distances, _ = bellman_ford(graph, 'a')
        assert distances['d'] == -2

    def test_bellman_ford_negative_cycle_raises(self):
        graph = {
            'a': [('b', 1)],
            'b': [('c', -1)],
            'c': [('b', -1)],
        }
        with pytest.raises(NegativeCycleError):
            bellman_ford(graph, 'a')


class TestDynamicProgramming:
    def test_knapsack_01(self):
        weights = [1, 2, 3]
        values = [6, 10, 12]
        assert knapsack_01(weights, values, 5) == 22

    def test_knapsack_01_mismatched_input_raises(self):
        with pytest.raises(ValueError):
            knapsack_01([1, 2], [1], 3)

    def test_longest_common_subsequence(self):
        assert longest_common_subsequence('ABCDE', 'ACE') == 3

    def test_edit_distance(self):
        assert edit_distance('kitten', 'sitting') == 3


class TestStrings:
    text = 'ABABDABACDABABCABAB'
    pattern = 'ABABCABAB'
    expected_index = 10

    def test_kmp_search(self):
        assert kmp_search(self.text, self.pattern) == [self.expected_index]

    def test_kmp_search_overlapping(self):
        assert kmp_search('aaaaa', 'aaaa') == [0, 1]

    def test_kmp_search_empty_pattern_raises(self):
        with pytest.raises(ValueError):
            kmp_search('hello', '')

    def test_rabin_karp_search(self):
        assert rabin_karp_search(self.text, self.pattern) == [self.expected_index]

    def test_boyer_moore_search(self):
        assert boyer_moore_search(self.text, self.pattern) == [self.expected_index]


class TestGreedy:
    def test_activity_selection(self):
        activities = [
            (1, 3),
            (2, 5),
            (4, 7),
            (1, 8),
            (5, 9),
            (8, 10),
            (9, 11),
            (11, 14),
            (13, 16),
        ]
        selected = activity_selection(activities)
        assert len(selected) >= 4

    def test_activity_selection_invalid_activity_raises(self):
        with pytest.raises(ValueError):
            activity_selection([(2, 1)])

    def test_fractional_knapsack(self):
        weights = [10, 20, 30]
        values = [60, 100, 120]
        assert fractional_knapsack(weights, values, 50) == 240.0

    def test_huffman_coding(self):
        frequencies = {'a': 5, 'b': 9, 'c': 12, 'd': 13, 'e': 16, 'f': 45}
        codes = huffman_coding(frequencies)
        assert len(codes) == 6
        for code in codes.values():
            assert isinstance(code, str)
            assert code


class TestDivideAndConquer:
    def test_max_subarray(self):
        arr = [-2, 1, -3, 4, -1, 2, 1, -5, 4]
        assert max_subarray(arr) == 6

    def test_max_subarray_all_negative(self):
        assert max_subarray([-2, -1]) == -1

    def test_count_inversions(self):
        assert count_inversions([1, 3, 5, 2, 4, 6]) == 3

    def test_fast_power(self):
        assert fast_power(2, 10) == 1024
        assert fast_power(5, 0) == 1
        assert fast_power(2, -2) == 0.25
