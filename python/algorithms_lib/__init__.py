'''Algorithms Library for Python.

A zero-dependency collection of common sorting, searching, graph,
dynamic-programming, string, greedy, and divide-and-conquer algorithms.
'''

from .divide_conquer import count_inversions, fast_power, max_subarray
from .dynamic_programming import (
    edit_distance,
    knapsack_01,
    longest_common_subsequence,
)
from .exceptions import (
    AlgorithmsError,
    EmptyInputError,
    InvalidGraphError,
    NegativeCycleError,
    ValueNotFoundError,
)
from .graphs import a_star, bellman_ford, bfs, dfs, dijkstra
from .greedy import activity_selection, fractional_knapsack, huffman_coding
from .searching import binary_search, interpolation_search, jump_search
from .sorting import heap_sort, merge_sort, quick_sort, radix_sort, tim_sort
from .strings import boyer_moore_search, kmp_search, rabin_karp_search

__version__ = '1.0.0'

__all__ = [
    'AlgorithmsError',
    'EmptyInputError',
    'InvalidGraphError',
    'NegativeCycleError',
    'ValueNotFoundError',
    'a_star',
    'activity_selection',
    'bellman_ford',
    'binary_search',
    'boyer_moore_search',
    'bfs',
    'count_inversions',
    'dfs',
    'dijkstra',
    'edit_distance',
    'fast_power',
    'fractional_knapsack',
    'heap_sort',
    'huffman_coding',
    'interpolation_search',
    'jump_search',
    'kmp_search',
    'knapsack_01',
    'longest_common_subsequence',
    'max_subarray',
    'merge_sort',
    'quick_sort',
    'radix_sort',
    'rabin_karp_search',
    'tim_sort',
]
