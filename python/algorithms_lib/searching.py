'''Searching algorithms for sorted sequences.'''

from typing import List, TypeVar

from .exceptions import ValueNotFoundError

T = TypeVar('T')


def _require_sorted(arr: List[T], name: str) -> None:
    for i in range(1, len(arr)):
        if arr[i] < arr[i - 1]:
            raise ValueError(f'{name} requires a list sorted in ascending order')


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

        pos = lo + int(
            ((target - arr[lo]) / (arr[hi] - arr[lo])) * (hi - lo)
        )
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
