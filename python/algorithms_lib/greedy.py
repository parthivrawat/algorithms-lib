'''Greedy algorithm implementations.'''

import heapq
from typing import Dict, List, Tuple, Union

from .exceptions import EmptyInputError


def activity_selection(activities: List[Tuple[int, int]]) -> List[int]:
    '''Return the indices of a maximum-size set of compatible activities.

    ``activities`` is a list of ``(start, end)`` tuples.
    '''
    if not activities:
        return []

    indexed = list(enumerate(activities))
    for idx, (start, end) in indexed:
        if start > end:
            raise ValueError(f'Activity {idx} has start after end')

    indexed.sort(key=lambda x: x[1][1])
    selected: List[int] = [indexed[0][0]]
    last_end = indexed[0][1][1]

    for idx, (start, end) in indexed[1:]:
        if start >= last_end:
            selected.append(idx)
            last_end = end
    return selected


def fractional_knapsack(
    weights: List[Union[int, float]],
    values: List[Union[int, float]],
    capacity: Union[int, float],
) -> float:
    '''Return the maximum value achievable in the fractional knapsack problem.

    Unlike 0/1 knapsack, items can be taken partially.
    '''
    if len(weights) != len(values):
        raise ValueError('weights and values must have the same length')
    if capacity < 0:
        raise ValueError('capacity must be non-negative')

    items = [
        (values[i] / weights[i], weights[i], values[i])
        for i in range(len(weights))
        if weights[i] > 0
    ]
    items.sort(reverse=True, key=lambda x: x[0])

    total = 0.0
    remaining = float(capacity)
    for ratio, weight, value in items:
        take = min(weight, remaining)
        total += take * ratio
        remaining -= take
        if remaining <= 0:
            break
    return total


def huffman_coding(frequencies: Dict[str, int]) -> Dict[str, str]:
    '''Return a prefix-free Huffman code table for the given symbol frequencies.'''
    if not frequencies:
        raise EmptyInputError('frequencies must not be empty')

    class _Node:
        def __init__(self, freq, char=None, left=None, right=None):
            self.freq = freq
            self.char = char
            self.left = left
            self.right = right

        def __lt__(self, other):
            return self.freq < other.freq

    heap = [_Node(freq, char=char) for char, freq in frequencies.items()]
    heapq.heapify(heap)

    if len(heap) == 1:
        return {heap[0].char: '0'}

    while len(heap) > 1:
        left = heapq.heappop(heap)
        right = heapq.heappop(heap)
        parent = _Node(left.freq + right.freq, left=left, right=right)
        heapq.heappush(heap, parent)

    codes: Dict[str, str] = {}
    _assign_codes(heap[0], '', codes)
    return codes


def _assign_codes(node, prefix: str, codes: Dict[str, str]) -> None:
    if node is None:
        return
    if node.char is not None:
        codes[node.char] = prefix or '0'
        return
    _assign_codes(node.left, prefix + '0', codes)
    _assign_codes(node.right, prefix + '1', codes)
