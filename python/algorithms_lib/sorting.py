'''Comparison and integer sorting algorithms.'''

import heapq
from typing import Generic, Iterable, List, TypeVar

from .exceptions import EmptyInputError

T = TypeVar('T')


def quick_sort(items: List[T]) -> List[T]:
    '''Return a new list sorted in ascending order using quick sort.

    Time complexity: O(n log n) average, O(n^2) worst case.
    Space complexity: O(n log n) average due to allocated partitions,
    O(n^2) worst case. This is not an in-place implementation.
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
        TypeError: If any item is not an integer.
        ValueError: If any item is negative.
    '''
    if not items:
        return []
    for x in items:
        if not isinstance(x, int) or isinstance(x, bool):
            raise TypeError('radix_sort only supports non-negative integers')
        if x < 0:
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


def native_sort(items: List[T]) -> List[T]:
    '''Return a new list sorted in ascending order using the built-in sort.

    This is a thin wrapper around the built-in ``sorted`` for reference and
    benchmarking.
    '''
    return sorted(items)


def quickselect(items: List[T], k: int) -> T:
    '''Return the k-th smallest element of ``items`` (0-indexed).

    ``items`` is not required to be sorted. ``k`` must satisfy
    ``0 <= k < len(items)``.

    Raises:
        TypeError: If ``k`` is not an integer.
        IndexError: If ``k`` is out of range.
    '''
    if not isinstance(k, int) or isinstance(k, bool):
        raise TypeError('k must be an integer')
    if not 0 <= k < len(items):
        raise IndexError('k is out of range')

    arr = list(items)
    while True:
        if len(arr) == 1:
            return arr[0]
        pivot = arr[len(arr) // 2]
        lows = [x for x in arr if x < pivot]
        pivots = [x for x in arr if x == pivot]
        highs = [x for x in arr if x > pivot]
        if k < len(lows):
            arr = lows
        elif k < len(lows) + len(pivots):
            return pivot
        else:
            k -= len(lows) + len(pivots)
            arr = highs


class PriorityQueue(Generic[T]):
    '''A min-priority queue backed by a binary heap.

    Items must support the ``<`` operator.
    '''

    def __init__(self, items: Iterable[T] = ()):
        self._heap: List[T] = list(items)
        heapq.heapify(self._heap)

    def push(self, item: T) -> None:
        '''Add an item to the queue.'''
        heapq.heappush(self._heap, item)

    def pop(self) -> T:
        '''Remove and return the smallest item.'''
        if not self._heap:
            raise EmptyInputError('pop from an empty priority queue')
        return heapq.heappop(self._heap)

    def peek(self) -> T:
        '''Return the smallest item without removing it.'''
        if not self._heap:
            raise EmptyInputError('peek from an empty priority queue')
        return self._heap[0]

    def __len__(self) -> int:
        return len(self._heap)

    def __bool__(self) -> bool:
        return bool(self._heap)
