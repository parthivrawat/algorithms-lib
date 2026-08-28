'''Divide-and-conquer utilities.'''

from typing import List, Tuple, TypeVar, Union

from .exceptions import EmptyInputError

T = TypeVar('T', int, float)


def max_subarray(arr: List[T]) -> T:
    '''Return the maximum sum of any contiguous subarray.

    Uses a divide-and-conquer implementation (Kadane is also available
    elsewhere but this one demonstrates the D&C pattern).
    '''
    if not arr:
        raise EmptyInputError('max_subarray requires a non-empty array')
    return _max_subarray_dc(arr, 0, len(arr) - 1)


def _max_subarray_dc(arr: List[T], left: int, right: int) -> T:
    if left == right:
        return arr[left]

    mid = (left + right) // 2
    left_sum = _max_subarray_dc(arr, left, mid)
    right_sum = _max_subarray_dc(arr, mid + 1, right)
    cross_sum = _max_crossing_sum(arr, left, mid, right)
    return max(left_sum, right_sum, cross_sum)


def _max_crossing_sum(arr: List[T], left: int, mid: int, right: int) -> T:
    total = 0
    left_max = float('-inf')
    for i in range(mid, left - 1, -1):
        total += arr[i]
        left_max = max(left_max, total)

    total = 0
    right_max = float('-inf')
    for i in range(mid + 1, right + 1):
        total += arr[i]
        right_max = max(right_max, total)

    return left_max + right_max


def count_inversions(arr: List[T]) -> int:
    '''Return the number of inversions in the array using merge sort.

    Time complexity: O(n log n).
    '''
    _, count = _sort_and_count(list(arr))
    return count


def _sort_and_count(arr: List[T]) -> Tuple[List[T], int]:  # type: ignore
    n = len(arr)
    if n <= 1:
        return arr, 0

    mid = n // 2
    left, lcount = _sort_and_count(arr[:mid])
    right, rcount = _sort_and_count(arr[mid:])
    merged, scount = _merge_and_count(left, right)
    return merged, lcount + rcount + scount


def _merge_and_count(left: List[T], right: List[T]) -> Tuple[List[T], int]:  # type: ignore
    merged: List[T] = []
    count = 0
    i = j = 0
    while i < len(left) and j < len(right):
        if left[i] <= right[j]:
            merged.append(left[i])
            i += 1
        else:
            merged.append(right[j])
            j += 1
            count += len(left) - i
    merged.extend(left[i:])
    merged.extend(right[j:])
    return merged, count


def fast_power(base: Union[int, float], exponent: int) -> Union[int, float]:
    '''Return ``base`` raised to ``exponent`` using exponentiation by squaring.

    Time complexity: O(log exponent).
    '''
    if exponent < 0:
        return 1 / fast_power(base, -exponent)
    if exponent == 0:
        return 1
    if exponent == 1:
        return base

    half = fast_power(base, exponent // 2)
    if exponent % 2 == 0:
        return half * half
    return base * half * half
