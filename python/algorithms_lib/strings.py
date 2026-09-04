'''String searching algorithms.'''

from typing import List

from .exceptions import InvalidInputError


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

    Raises:
        ValueError: If ``pattern`` is empty, ``mod`` is not positive, or
            ``base`` is not positive.
    '''
    if not pattern:
        raise InvalidInputError('pattern must not be empty')
    if mod <= 0:
        raise InvalidInputError('mod must be positive')
    if base <= 0:
        raise InvalidInputError('base must be positive')

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
    return matches


def boyer_moore_search(text: str, pattern: str) -> List[int]:
    '''Return all starting indices of ``pattern`` in ``text`` using Boyer-Moore
    with both the bad-character and good-suffix rules.

    Time complexity: O(m) preprocessing and sublinear scanning on typical 
    inputs; O(n * m) worst case when many matches are reported.
    '''
    if not pattern:
        raise ValueError('pattern must not be empty')

    n, m = len(text), len(pattern)
    if m > n:
        return []

    bad_char = {}
    for i in range(m):
        bad_char[pattern[i]] = i
    good_suffix = _good_suffix_shifts(pattern)

    matches: List[int] = []
    i = 0
    while i <= n - m:
        j = m - 1
        while j >= 0 and text[i + j] == pattern[j]:
            j -= 1
        if j < 0:
            matches.append(i)
            i += good_suffix[0]
        else:
            bad_char_shift = j - bad_char.get(text[i + j], -1)
            i += max(good_suffix[j + 1], bad_char_shift, 1)
    return matches


def _good_suffix_shifts(pattern: str) -> List[int]:
    '''Build the Boyer-Moore good-suffix shift table.

    ``shift[k]`` is the shift applied when a mismatch occurs at pattern index
    ``k - 1``; ``shift[0]`` is applied after a full match and is derived from
    the pattern's longest border, so overlapping matches are still found.
    '''
    m = len(pattern)
    shift = [0] * (m + 1)
    border_pos = [0] * (m + 1)
    i, j = m, m + 1
    border_pos[i] = j
    while i > 0:
        while j <= m and pattern[i - 1] != pattern[j - 1]:
            if shift[j] == 0:
                shift[j] = j - i
            j = border_pos[j]
        i -= 1
        j -= 1
        border_pos[i] = j
    j = border_pos[0]
    for i in range(m + 1):
        if shift[i] == 0:
            shift[i] = j
        if i == j:
            j = border_pos[j]
    return shift
