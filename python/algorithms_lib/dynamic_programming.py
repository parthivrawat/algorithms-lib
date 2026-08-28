'''Dynamic programming solutions to classic problems.'''

from typing import Sequence, TypeVar

from .exceptions import EmptyInputError

T = TypeVar('T')


def knapsack_01(weights: Sequence[int], values: Sequence[int], capacity: int) -> int:
    '''Return the maximum value achievable for the 0/1 knapsack problem.

    ``weights`` and ``values`` must be equal-length sequences. ``capacity``
    must be non-negative.
    '''
    if len(weights) != len(values):
        raise ValueError('weights and values must have the same length')
    if capacity < 0:
        raise ValueError('capacity must be non-negative')

    n = len(weights)
    dp = [0] * (capacity + 1)
    for i in range(n):
        w = weights[i]
        v = values[i]
        if w < 0:
            raise ValueError('weights must be non-negative')
        for c in range(capacity, w - 1, -1):
            dp[c] = max(dp[c], dp[c - w] + v)
    return dp[capacity]


def longest_common_subsequence(a: Sequence[T], b: Sequence[T]) -> int:
    '''Return the length of the longest common subsequence of ``a`` and ``b``.'''
    m, n = len(a), len(b)
    prev = [0] * (n + 1)
    for i in range(1, m + 1):
        curr = [0] * (n + 1)
        for j in range(1, n + 1):
            if a[i - 1] == b[j - 1]:
                curr[j] = prev[j - 1] + 1
            else:
                curr[j] = max(prev[j], curr[j - 1])
        prev = curr
    return prev[n]


def edit_distance(a: Sequence[T], b: Sequence[T]) -> int:
    '''Return the Levenshtein distance between sequences ``a`` and ``b``.'''
    m, n = len(a), len(b)
    prev = list(range(n + 1))
    for i in range(1, m + 1):
        curr = [i] + [0] * n
        for j in range(1, n + 1):
            cost = 0 if a[i - 1] == b[j - 1] else 1
            curr[j] = min(
                curr[j - 1] + 1,      # insertion
                prev[j] + 1,          # deletion
                prev[j - 1] + cost,   # substitution
            )
        prev = curr
    return prev[n]
