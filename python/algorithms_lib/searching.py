'''Searching algorithms for sorted sequences.'''

from typing import Callable, List, Optional, TypeVar

from .exceptions import InvalidInputError

T = TypeVar('T')


def _require_sorted(arr: List[T], name: str) -> None:
    for i in range(1, len(arr)):
        if arr[i] < arr[i - 1]:
            raise InvalidInputError(
                f'{name} requires a list sorted in ascending order'
            )


def binary_search(arr: List[T], target: T) -> int:
    '''Return the index of ``target`` in a sorted list, or -1 if not found.

    Time complexity: O(log n).
    '''
    if not arr:
        return -1
    _require_sorted(arr, 'binary_search')

    lo, hi = 0, len(arr) - 1
    while lo <= hi:
        mid = (lo + hi) // 2
        if arr[mid] == target:
            return mid
        if arr[mid] < target:
            lo = mid + 1
        else:
            hi = mid - 1
    return -1


def interpolation_search(arr: List[int], target: int) -> int:
    '''Return the index of ``target`` in a sorted numeric list, or -1.

    Time complexity: O(log log n) average for uniformly distributed data,
    O(n) worst case.
    '''
    if not arr:
        return -1
    _require_sorted(arr, 'interpolation_search')

    lo, hi = 0, len(arr) - 1
    while lo <= hi and target >= arr[lo] and target <= arr[hi]:
        if arr[hi] == arr[lo]:
            if arr[lo] == target:
                return lo
            break

        pos = lo + (target - arr[lo]) * (hi - lo) // (arr[hi] - arr[lo])
        if pos < lo or pos > hi:
            break

        if arr[pos] == target:
            return pos
        if arr[pos] < target:
            lo = pos + 1
        else:
            hi = pos - 1
    return -1


def jump_search(arr: List[T], target: T) -> int:
    '''Return the index of ``target`` in a sorted list, or -1.

    Time complexity: O(sqrt(n)).
    '''
    if not arr:
        return -1
    _require_sorted(arr, 'jump_search')

    n = len(arr)
    step = int(n ** 0.5)
    prev = 0

    while arr[min(step, n) - 1] < target:
        prev = step
        step += int(n ** 0.5)
        if prev >= n:
            return -1

    while arr[prev] < target:
        prev += 1
        if prev == min(step, n):
            return -1

    if arr[prev] == target:
        return prev
    return -1


def binary_search_on_answer(
    low: int,
    high: int,
    predicate: Callable[[int], bool],
    find: str = 'minimum'
) -> Optional[int]:
    '''Find the smallest or largest integer ``x`` in ``[low, high]`` for which
    ``predicate(x)`` is true, assuming a monotone predicate.

    ``find`` must be ``'minimum'`` or ``'maximum'``. Returns ``None`` if no
    value in the range satisfies the predicate.
    '''
    if find not in ('minimum', 'maximum'):
        raise InvalidInputError(
            "binary_search_on_answer: find must be 'minimum' or 'maximum'"
        )
    if low > high:
        return None
    if find == 'minimum':
        left, right = low, high
        answer = None
        while left <= right:
            mid = (left + right) // 2
            if predicate(mid):
                answer = mid
                right = mid - 1
            else:
                left = mid + 1
        return answer
    left, right = low, high
    answer = None
    while left <= right:
        mid = (left + right) // 2
        if predicate(mid):
            answer = mid
            left = mid + 1
        else:
            right = mid - 1
    return answer
