'''Algorithms Library for Python.

A zero-dependency collection of common sorting, searching, graph,
dynamic-programming, string, greedy, and divide-and-conquer algorithms.
'''

from .divide_conquer import (
    count_inversions,
    fast_power,
    gcd,
    lcm,
    max_subarray,
    modular_fast_power,
    sieve_of_eratosthenes,
)
from .dynamic_programming import (
    edit_distance,
    edit_distance_reconstruction,
    knapsack_01,
    longest_common_subsequence,
    longest_common_subsequence_reconstruction,
)
from .exceptions import (
    AlgorithmsError,
    EmptyInputError,
    InvalidGraphError,
    InvalidInputError,
    NegativeCycleError,
    ValueNotFoundError,
)
from .graphs import (
    UnionFind,
    a_star,
    bellman_ford,
    bfs,
    dfs,
    dijkstra,
    has_cycle,
    floyd_warshall,
    minimum_spanning_tree,
    strongly_connected_components,
    topological_sort,
)
from .greedy import (
    activity_selection,
    fractional_knapsack,
    huffman_coding,
    huffman_decode,
    huffman_encode,
)
from .searching import (
    binary_search,
    binary_search_on_answer,
    interpolation_search,
    jump_search,
)
from .sorting import PriorityQueue, heap_sort, merge_sort, native_sort, quick_sort, quickselect, radix_sort
from .strings import boyer_moore_search, kmp_search, rabin_karp_search

__version__ = '2.0.0'

__all__ = [
    '__version__',
    'AlgorithmsError',
    'EmptyInputError',
    'InvalidGraphError',
    'InvalidInputError',
    'NegativeCycleError',
    'PriorityQueue',
    'UnionFind',
    'ValueNotFoundError',
    'a_star',
    'activity_selection',
    'bellman_ford',
    'binary_search',
    'binary_search_on_answer',
    'boyer_moore_search',
    'bfs',
    'count_inversions',
    'dfs',
    'dijkstra',
    'edit_distance',
    'edit_distance_reconstruction',
    'fast_power',
    'floyd_warshall',
    'fractional_knapsack',
    'gcd',
    'has_cycle',
    'heap_sort',
    'huffman_coding',
    'huffman_decode',
    'huffman_encode',
    'interpolation_search',
    'jump_search',
    'kmp_search',
    'knapsack_01',
    'lcm',
    'longest_common_subsequence',
    'longest_common_subsequence_reconstruction',
    'max_subarray',
    'minimum_spanning_tree',
    'merge_sort',
    'modular_fast_power',
    'quick_sort',
    'quickselect',
    'native_sort',
    'radix_sort',
    'rabin_karp_search',
    'sieve_of_eratosthenes',
    'strongly_connected_components',
    'topological_sort',
]
