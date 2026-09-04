'''Dynamic programming solutions to classic problems.'''

from typing import List, Optional, Sequence, Tuple, TypeVar

from .exceptions import InvalidInputError

T = TypeVar('T')


def knapsack_01(weights: Sequence[int], values: Sequence[int], capacity: int) -> int:
    '''Return the maximum value achievable for the 0/1 knapsack problem.

    ``weights`` and ``values`` must be equal-length sequences of integers.
    ``capacity`` must be a non-negative integer.
    '''
    if not isinstance(capacity, int) or isinstance(capacity, bool):
        raise InvalidInputError('capacity must be a non-negative integer')
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


def longest_common_subsequence_reconstruction(
    a: Sequence[T], b: Sequence[T]
) -> Tuple[int, List[T]]:
    '''Return the length and one longest common subsequence of ``a`` and ``b``.'''
    m, n = len(a), len(b)
    dp = [[0] * (n + 1) for _ in range(m + 1)]
    for i in range(1, m + 1):
        for j in range(1, n + 1):
            if a[i - 1] == b[j - 1]:
                dp[i][j] = dp[i - 1][j - 1] + 1
            else:
                dp[i][j] = max(dp[i - 1][j], dp[i][j - 1])
    result: List[T] = []
    i, j = m, n
    while i > 0 and j > 0:
        if a[i - 1] == b[j - 1]:
            result.append(a[i - 1])
            i -= 1
            j -= 1
        elif dp[i - 1][j] >= dp[i][j - 1]:
            i -= 1
        else:
            j -= 1
    result.reverse()
    return dp[m][n], result


def edit_distance_reconstruction(
    a: Sequence[T], b: Sequence[T]
) -> Tuple[int, List[Tuple[str, Optional[T], Optional[T]]]]:
    '''Return the Levenshtein distance and one optimal edit script.'''
    m, n = len(a), len(b)
    dp = [[0] * (n + 1) for _ in range(m + 1)]
    for i in range(m + 1):
        dp[i][0] = i
    for j in range(n + 1):
        dp[0][j] = j
    for i in range(1, m + 1):
        for j in range(1, n + 1):
            if a[i - 1] == b[j - 1]:
                dp[i][j] = dp[i - 1][j - 1]
            else:
                dp[i][j] = 1 + min(
                    dp[i - 1][j],
                    dp[i][j - 1],
                    dp[i - 1][j - 1],
                )
    script: List[Tuple[str, Optional[T], Optional[T]]] = []
    i, j = m, n
    while i > 0 or j > 0:
        if i == 0:
            script.append(('insert', None, b[j - 1]))
            j -= 1
        elif j == 0:
            script.append(('delete', a[i - 1], None))
            i -= 1
        elif a[i - 1] == b[j - 1]:
            script.append(('match', a[i - 1], b[j - 1]))
            i -= 1
            j -= 1
        else:
            best = dp[i][j]
            if dp[i - 1][j - 1] + 1 == best:
                script.append(('substitute', a[i - 1], b[j - 1]))
                i -= 1
                j -= 1
            elif dp[i][j - 1] + 1 == best:
                script.append(('insert', None, b[j - 1]))
                j -= 1
            else:
                script.append(('delete', a[i - 1], None))
                i -= 1
    script.reverse()
    return dp[m][n], script
