'''Comparison and integer sorting algorithms.'''

from typing import List, TypeVar

T = TypeVar('T')


def quick_sort(items: List[T]) -> List[T]:
    '''Return a new list sorted in ascending order using quick sort.

    Time complexity: O(n log n) average, O(n^2) worst case.
    Space complexity: O(n) due to allocated partitions.
    '''
    if len(items) <= 1:
        return list(items)
    pivot = items[len(items) // 2]
    left = [x for x in items if x < pivot]
    middle = [x for x in items if x == pivot]
    right = [x for x in items if x > pivot]
    return quick_sort(left) + middle + quick_sort(right)


def merge_sort(items: List[T]) -> List[T]:
    '''Return a new list sorted in ascending order using stable merge sort.

    Time complexity: O(n log n).
    Space complexity: O(n).
    '''
    if len(items) <= 1:
        return list(items)
    mid = len(items) // 2
    left = merge_sort(items[:mid])
    right = merge_sort(items[mid:])
    return _merge(left, right)


def _merge(left: List[T], right: List[T]) -> List[T]:
    merged: List[T] = []
    i = j = 0
    while i < len(left) and j < len(right):
        if left[i] <= right[j]:
            merged.append(left[i])
            i += 1
        else:
            merged.append(right[j])
            j += 1
    merged.extend(left[i:])
    merged.extend(right[j:])
    return merged


def heap_sort(items: List[T]) -> List[T]:
    '''Return a new list sorted in ascending order using a binary max-heap.

    Time complexity: O(n log n).
    Space complexity: O(n) for the copy.
    '''
    arr = list(items)
    n = len(arr)
    if n <= 1:
        return arr

    for i in range(n // 2 - 1, -1, -1):
        _sift_down(arr, n, i)

    for end in range(n - 1, 0, -1):
        arr[0], arr[end] = arr[end], arr[0]
        _sift_down(arr, end, 0)
    return arr


def _sift_down(arr: List[T], n: int, i: int) -> None:
    largest = i
    left = 2 * i + 1
    right = 2 * i + 2

    if left < n and arr[left] > arr[largest]:
        largest = left
    if right < n and arr[right] > arr[largest]:
        largest = right

    if largest != i:
        arr[i], arr[largest] = arr[largest], arr[i]
        _sift_down(arr, n, largest)


def radix_sort(items: List[int]) -> List[int]:
    '''Return a new list of non-negative integers sorted using LSD radix sort.

    Time complexity: O(d * n) where d is the number of digits.
    Space complexity: O(n).

    Raises:
        ValueError: If any item is negative.
    '''
    if not items:
        return []
    if any(x < 0 for x in items):
        raise ValueError('radix_sort only supports non-negative integers')

    arr = list(items)
    max_val = max(arr)
    exp = 1
    while max_val // exp > 0:
        _counting_sort_by_digit(arr, exp)
        exp *= 10
    return arr


def _counting_sort_by_digit(arr: List[int], exp: int) -> None:
    n = len(arr)
    output = [0] * n
    count = [0] * 10

    for num in arr:
        index = (num // exp) % 10
        count[index] += 1

    for i in range(1, 10):
        count[i] += count[i - 1]

    for i in range(n - 1, -1, -1):
        index = (arr[i] // exp) % 10
        output[count[index] - 1] = arr[i]
        count[index] -= 1

    for i in range(n):
        arr[i] = output[i]


def tim_sort(items: List[T]) -> List[T]:
    '''Return a new list sorted in ascending order using Python\'s Timsort.

    This is a thin wrapper around the built-in ``sorted`` to expose Timsort
    by name for reference and benchmarking.
    '''
    return sorted(items)
