'''String searching algorithms.'''

from typing import List


def kmp_search(text: str, pattern: str) -> List[int]:
    '''Return all starting indices of ``pattern`` in ``text`` using KMP.

    Time complexity: O(n + m).
    '''
    if not pattern:
        raise ValueError('pattern must not be empty')

    failure = _kmp_failure(pattern)
    matches: List[int] = []
    j = 0
    for i in range(len(text)):
        while j > 0 and text[i] != pattern[j]:
            j = failure[j - 1]
        if text[i] == pattern[j]:
            j += 1
        if j == len(pattern):
            matches.append(i - len(pattern) + 1)
            j = failure[j - 1]
    return matches


def _kmp_failure(pattern: str) -> List[int]:
    failure = [0] * len(pattern)
    j = 0
    for i in range(1, len(pattern)):
        while j > 0 and pattern[i] != pattern[j]:
            j = failure[j - 1]
        if pattern[i] == pattern[j]:
            j += 1
            failure[i] = j
    return failure


def rabin_karp_search(
    text: str, pattern: str, base: int = 256, mod: int = 1_000_000_007
) -> List[int]:
    '''Return all starting indices of ``pattern`` in ``text`` using rolling hash.

    Time complexity: O(n + m) average, O(n * m) worst-case due to hash
    collisions (mitigated by a full character comparison).
    '''
    if not pattern:
        raise ValueError('pattern must not be empty')
    n, m = len(text), len(pattern)
    if m > n:
        return []

    h = pow(base, m - 1, mod)
    pattern_hash = 0
    text_hash = 0
    for i in range(m):
        pattern_hash = (pattern_hash * base + ord(pattern[i])) % mod
        text_hash = (text_hash * base + ord(text[i])) % mod

    matches: List[int] = []
    for i in range(n - m + 1):
        if text_hash == pattern_hash and text[i:i + m] == pattern:
            matches.append(i)
        if i < n - m:
            text_hash = (text_hash - ord(text[i]) * h) % mod
            text_hash = (text_hash * base + ord(text[i + m])) % mod
            text_hash %= mod
    return matches


def boyer_moore_search(text: str, pattern: str) -> List[int]:
    '''Return all starting indices of ``pattern`` in ``text`` using Boyer-Moore
    with the bad-character rule.

    Time complexity: O(n * m) worst case, but often much better in practice.
    '''
    if not pattern:
        raise ValueError('pattern must not be empty')

    n, m = len(text), len(pattern)
    if m > n:
        return []

    bad_char = {}
    for i in range(m):
        bad_char[pattern[i]] = i

    matches: List[int] = []
    i = 0
    while i <= n - m:
        j = m - 1
        while j >= 0 and text[i + j] == pattern[j]:
            j -= 1
        if j < 0:
            matches.append(i)
            i += m
        else:
            shift = j - bad_char.get(text[i + j], -1)
            i += max(1, shift)
    return matches
