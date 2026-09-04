'''Tests for the algorithms library.'''

import pytest

from algorithms_lib import (
    __version__,
    AlgorithmsError,
    EmptyInputError,
    InvalidGraphError,
    InvalidInputError,
    NegativeCycleError,
    PriorityQueue,
    UnionFind,
    ValueNotFoundError,
    a_star,
    activity_selection,
    bellman_ford,
    binary_search,
    binary_search_on_answer,
    boyer_moore_search,
    bfs,
    count_inversions,
    sieve_of_eratosthenes,
    topological_sort,
    dfs,
    dijkstra,
    has_cycle,
    edit_distance,
    edit_distance_reconstruction,
    fast_power,
    floyd_warshall,
    minimum_spanning_tree,
    strongly_connected_components,
    fractional_knapsack,
    heap_sort,
    huffman_coding,
    huffman_decode,
    huffman_encode,
    interpolation_search,
    jump_search,
    kmp_search,
    knapsack_01,
    longest_common_subsequence,
    longest_common_subsequence_reconstruction,
    gcd,
    lcm,
    max_subarray,
    merge_sort,
    modular_fast_power,
    quick_sort,
    quickselect,
    native_sort,
    radix_sort,
    rabin_karp_search,
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
    def test_native_sort(self, data):
        assert native_sort(data) == sorted(data)

    def test_radix_sort(self):
        assert radix_sort([170, 45, 75, 90, 2, 802, 24, 66]) == [
            2, 24, 45, 66, 75, 90, 170, 802
        ]

    def test_radix_sort_empty(self):
        assert radix_sort([]) == []

    def test_radix_sort_negative_raises(self):
        with pytest.raises(ValueError):
            radix_sort([1, -5, 3])

    def test_radix_sort_non_integer_raises(self):
        with pytest.raises(TypeError):
            radix_sort([1.5, 2, 3])
        with pytest.raises(TypeError):
            radix_sort([1, 'a'])
        with pytest.raises(TypeError):
            radix_sort([1.0, 2.0])

    def test_radix_sort_large_values(self):
        assert radix_sort([10 ** 18, 1, 0, 10 ** 18 - 1]) == [
            0, 1, 10 ** 18 - 1, 10 ** 18
        ]


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

    def test_interpolation_search_large_values(self):
        # Values beyond 2^53 would be rounded by float division.
        base = 2 ** 60
        arr = list(range(base, base + 1000, 100))
        assert interpolation_search(arr, base + 300) == 3
        assert interpolation_search(arr, base + 999) == -1

    def test_interpolation_search_extreme_range(self):
        arr = [-(2 ** 63), -1, 0, 1, 2 ** 63 - 1]
        assert interpolation_search(arr, 2 ** 63 - 1) == 4
        assert interpolation_search(arr, 2) == -1

    def test_jump_search_found(self):
        assert jump_search(self.sorted_nums, 10) == 4

    def test_jump_search_not_found(self):
        assert jump_search(self.sorted_nums, 7) == -1

    def test_binary_search_on_answer(self):
        assert binary_search_on_answer(1, 10, lambda x: x >= 6) == 6
        assert binary_search_on_answer(
            1, 10, lambda x: x <= 4, find='maximum'
        ) == 4


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

    def test_bfs_missing_neighbor_raises(self):
        with pytest.raises(InvalidGraphError):
            bfs({'a': ['b']}, 'a')

    def test_dfs_missing_neighbor_raises(self):
        with pytest.raises(InvalidGraphError):
            dfs({'a': ['b']}, 'a')

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

    def test_dijkstra_unreachable_negative_weight_raises(self):
        graph = {'a': [], 'b': [('c', -1)], 'c': []}
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

    def test_a_star_negative_edge_weight_raises(self):
        graph = {'a': [('b', -1)], 'b': []}
        with pytest.raises(InvalidGraphError):
            a_star(graph, 'a', 'b', lambda node, goal: 0)

    def test_a_star_negative_heuristic_raises(self):
        graph = {'a': [], 'b': []}
        with pytest.raises(ValueError):
            a_star(graph, 'a', 'b', lambda node, goal: -1)

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

    def test_knapsack_01_float_capacity_raises(self):
        with pytest.raises(InvalidInputError):
            knapsack_01([1], [1], 2.5)
        with pytest.raises(InvalidInputError):
            knapsack_01([1], [1], True)

    def test_longest_common_subsequence(self):
        assert longest_common_subsequence('ABCDE', 'ACE') == 3

    def test_edit_distance(self):
        assert edit_distance('kitten', 'sitting') == 3

    def test_edit_distance_reconstruction(self):
        distance, script = edit_distance_reconstruction('kitten', 'sitting')
        assert distance == 3
        assert len(script) >= 3

    def test_longest_common_subsequence_reconstruction(self):
        length, lcs = longest_common_subsequence_reconstruction('ABCDE', 'ACE')
        assert length == 3
        assert lcs == ['A', 'C', 'E']


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

    def test_rabin_karp_search_mod_validation(self):
        with pytest.raises(ValueError):
            rabin_karp_search('hello', 'll', mod=0)
        with pytest.raises(ValueError):
            rabin_karp_search('hello', 'll', base=-1)

    def test_rabin_karp_search_overlapping(self):
        assert rabin_karp_search('aaaaa', 'aaaa') == [0, 1]

    def test_boyer_moore_search(self):
        assert boyer_moore_search(self.text, self.pattern) == [self.expected_index]

    def test_boyer_moore_search_overlapping(self):
        assert boyer_moore_search('aaaa', 'aa') == [0, 1, 2]
        assert boyer_moore_search('aaaaa', 'aaaa') == [0, 1]
        assert boyer_moore_search('abcabcabc', 'abcabc') == [0, 3]
        assert boyer_moore_search('abababab', 'abab') == [0, 2, 4]


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

    def test_fractional_knapsack_zero_weight_items(self):
        # Zero-weight, positive-value items are taken in full for free.
        assert fractional_knapsack([0, 10], [50, 60], 5) == 80.0
        # Zero-weight items with non-positive value are never taken.
        assert fractional_knapsack([0, 10], [-50, 60], 50) == 60.0
        assert fractional_knapsack([0], [5], 0) == 5.0

    def test_fractional_knapsack_negative_weight_raises(self):
        with pytest.raises(ValueError):
            fractional_knapsack([-5, 10], [10, 60], 50)

    def test_fractional_knapsack_negative_values_not_taken(self):
        # Items are optional: taking a negative-value item would lower the total.
        assert fractional_knapsack([10], [-5], 10) == 0.0
        assert fractional_knapsack([10, 5], [-5, 30], 10) == 30.0

    def test_huffman_coding(self):
        frequencies = {'a': 5, 'b': 9, 'c': 12, 'd': 13, 'e': 16, 'f': 45}
        codes = huffman_coding(frequencies)
        assert len(codes) == 6
        for code in codes.values():
            assert isinstance(code, str)
            assert code

    def test_huffman_coding_non_positive_frequency_raises(self):
        with pytest.raises(InvalidInputError):
            huffman_coding({'a': 0})
        with pytest.raises(InvalidInputError):
            huffman_coding({'a': -1})

    def test_huffman_encode_decode(self):
        frequencies = {'a': 5, 'b': 9, 'c': 12, 'd': 13, 'e': 16, 'f': 45}
        codes = huffman_coding(frequencies)
        symbols = ['a', 'b', 'c', 'd', 'e', 'f']
        encoded = huffman_encode(symbols, codes)
        decoded = huffman_decode(encoded, codes)
        assert decoded == symbols

    def test_huffman_encode_unknown_symbol_raises(self):
        with pytest.raises(InvalidInputError):
            huffman_encode(['x'], {'a': '0'})


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

    def test_fast_power_edge_cases(self):
        with pytest.raises(ValueError):
            fast_power(0, -1)
        assert fast_power(0, 0) == 1
        with pytest.raises(TypeError):
            fast_power(2, 2.5)
        assert fast_power(-2, 3) == -8
        assert fast_power(-2, 4) == 16
        assert fast_power(2, -10) == 2 ** -10
        assert fast_power(2, 10) == 2 ** 10 and isinstance(fast_power(2, 10), int)


class TestNumberTheory:
    def test_gcd(self):
        assert gcd(48, 18) == 6
        assert gcd(-48, 18) == 6
        assert gcd(0, 5) == 5
        assert gcd(0, 0) == 0

    def test_lcm(self):
        assert lcm(4, 6) == 12
        assert lcm(-4, 6) == 12
        assert lcm(0, 5) == 0

    def test_lcm_zero_zero_raises(self):
        with pytest.raises(InvalidInputError):
            lcm(0, 0)

    def test_modular_fast_power(self):
        assert modular_fast_power(2, 10, 1000) == 24
        assert modular_fast_power(3, 0, 7) == 1
        assert modular_fast_power(0, 0, 7) == 1
        assert modular_fast_power(0, 5, 7) == 0

    def test_modular_fast_power_invalid_input(self):
        with pytest.raises(ValueError):
            modular_fast_power(2, -1, 7)
        with pytest.raises(ValueError):
            modular_fast_power(2, 2, 0)

    def test_quickselect(self):
        assert quickselect([3, 2, 1, 5, 4], 0) == 1
        assert quickselect([3, 2, 1, 5, 4], 2) == 3
        assert quickselect([3, 2, 1, 5, 4], 4) == 5
        assert quickselect(['cherry', 'apple', 'banana'], 1) == 'banana'

    def test_quickselect_out_of_range(self):
        with pytest.raises(IndexError):
            quickselect([1, 2, 3], 5)
        with pytest.raises(IndexError):
            quickselect([1, 2, 3], -1)

    def test_sieve_of_eratosthenes(self):
        assert sieve_of_eratosthenes(10) == [2, 3, 5, 7]
        assert sieve_of_eratosthenes(1) == []
        assert sieve_of_eratosthenes(2) == [2]
        with pytest.raises(ValueError):
            sieve_of_eratosthenes(-1)


class TestGraph:
    def test_topological_sort(self):
        graph = {
            'a': ['b', 'c'],
            'b': ['d'],
            'c': ['d'],
            'd': [],
        }
        order = topological_sort(graph)
        assert order.index('a') < order.index('b')
        assert order.index('a') < order.index('c')
        assert order.index('b') < order.index('d')
        assert order.index('c') < order.index('d')

    def test_topological_sort_cycle(self):
        graph = {
            'a': ['b'],
            'b': ['c'],
            'c': ['a'],
        }
        with pytest.raises(InvalidGraphError):
            topological_sort(graph)

    def test_minimum_spanning_tree(self):
        edges = [
            ('a', 'b', 1),
            ('b', 'c', 2),
            ('a', 'c', 3),
            ('c', 'd', 4),
        ]
        total, mst = minimum_spanning_tree(edges)
        assert total == 7
        assert len(mst) == 3
        assert sum(w for _, _, w in mst) == 7

    def test_floyd_warshall(self):
        graph = {
            'a': [('b', 1), ('c', 4)],
            'b': [('c', 2), ('d', 5)],
            'c': [('d', 1)],
            'd': [],
        }
        dist = floyd_warshall(graph)
        assert dist['a']['d'] == 4
        assert dist['a']['a'] == 0
        assert dist['b']['d'] == 3

    def test_floyd_warshall_negative_cycle(self):
        graph = {
            'a': [('b', 1)],
            'b': [('a', -2)],
        }
        with pytest.raises(NegativeCycleError):
            floyd_warshall(graph)

    def test_strongly_connected_components(self):
        graph = {
            'a': ['b'],
            'b': ['c', 'd'],
            'c': ['a'],
            'd': ['e'],
            'e': [],
        }
        components = strongly_connected_components(graph)
        assert sorted([sorted(c) for c in components]) == [
            ['a', 'b', 'c'], ['d'], ['e']
        ]

    def test_has_cycle(self):
        assert not has_cycle({'a': ['b'], 'b': [], 'c': ['d'], 'd': []})
        assert has_cycle({'a': ['b'], 'b': ['c'], 'c': ['a']})

    def test_has_cycle_missing_node_raises(self):
        with pytest.raises(InvalidGraphError):
            has_cycle({'a': ['b']})


class TestDataStructures:
    def test_union_find(self):
        uf = UnionFind(['a', 'b', 'c', 'd', 'e'])
        assert not uf.connected('a', 'b')
        uf.union('a', 'b')
        assert uf.connected('a', 'b')
        uf.union('b', 'c')
        assert uf.connected('a', 'c')
        assert not uf.connected('a', 'd')
        assert uf.count() == 3

    def test_union_find_invalid_element(self):
        uf = UnionFind(['a', 'b'])
        with pytest.raises(InvalidInputError):
            uf.find('z')

    def test_priority_queue(self):
        pq = PriorityQueue([3, 1, 4, 1, 5])
        assert pq.peek() == 1
        assert len(pq) == 5
        assert pq.pop() == 1
        assert pq.pop() == 1
        assert pq.pop() == 3
        pq.push(0)
        assert pq.pop() == 0

    def test_priority_queue_empty(self):
        pq = PriorityQueue()
        with pytest.raises(EmptyInputError):
            pq.pop()


class TestExports:
    def test_version_exported(self):
        import algorithms_lib

        assert algorithms_lib.__version__ == '2.0.0'
        assert '__version__' in algorithms_lib.__all__
        assert 'InvalidInputError' in algorithms_lib.__all__
        assert 'gcd' in algorithms_lib.__all__
        assert 'lcm' in algorithms_lib.__all__
        assert 'modular_fast_power' in algorithms_lib.__all__
        assert 'quickselect' in algorithms_lib.__all__
