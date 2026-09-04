'''Divide-and-conquer utilities.'''

from typing import List, Tuple, TypeVar, Union

from .exceptions import EmptyInputError, InvalidInputError

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

    ``exponent`` must be an integer. By convention ``0 ** 0`` returns ``1``,
    matching Python's built-in ``**`` operator.

    Raises:
        TypeError: If ``exponent`` is not an integer.
        ValueError: If ``base`` is zero and ``exponent`` is negative.

    Time complexity: O(log exponent).
    '''
    if not isinstance(exponent, int) or isinstance(exponent, bool):
        raise TypeError('exponent must be an integer')
    if base == 0 and exponent < 0:
        raise ValueError('0 cannot be raised to a negative power')

    result: Union[int, float] = 1
    b: Union[int, float] = base
    exp = abs(exponent)
    while exp > 0:
        if exp & 1:
            result *= b
        b *= b
        exp >>= 1
    if exponent < 0:
        result = 1 / result
    return result


def gcd(a: int, b: int) -> int:
    '''Return the greatest common divisor of ``a`` and ``b``.

    ``a`` and ``b`` must be integers. The result is always non-negative.
    By convention, ``gcd(0, 0)`` returns ``0``.

    Raises:
        TypeError: If ``a`` or ``b`` is not an integer.
    '''
    if not isinstance(a, int) or isinstance(a, bool):
        raise TypeError('gcd arguments must be integers')
    if not isinstance(b, int) or isinstance(b, bool):
        raise TypeError('gcd arguments must be integers')
    a, b = abs(a), abs(b)
    while b:
        a, b = b, a % b
    return a


def lcm(a: int, b: int) -> int:
    '''Return the least common multiple of ``a`` and ``b``.

    ``a`` and ``b`` must be integers. The result is always non-negative.

    Raises:
        TypeError: If ``a`` or ``b`` is not an integer.
        InvalidInputError: If both ``a`` and ``b`` are zero.
    '''
    if not isinstance(a, int) or isinstance(a, bool):
        raise TypeError('lcm arguments must be integers')
    if not isinstance(b, int) or isinstance(b, bool):
        raise TypeError('lcm arguments must be integers')
    if a == 0 and b == 0:
        raise InvalidInputError('lcm of (0, 0) is undefined')
    g = gcd(a, b)
    return abs(a // g * b)


def modular_fast_power(base: int, exponent: int, mod: int) -> int:
    '''Return ``(base ** exponent) % mod`` using exponentiation by squaring.

    ``exponent`` must be a non-negative integer and ``mod`` must be a
    positive integer. By convention, ``0 ** 0`` returns ``1``.

    Raises:
        TypeError: If any argument is not an integer.
        ValueError: If ``exponent`` is negative or ``mod`` is not positive.
    '''
    if not isinstance(base, int) or isinstance(base, bool):
        raise TypeError('modular_fast_power arguments must be integers')
    if not isinstance(exponent, int) or isinstance(exponent, bool):
        raise TypeError('modular_fast_power arguments must be integers')
    if not isinstance(mod, int) or isinstance(mod, bool):
        raise TypeError('modular_fast_power arguments must be integers')
    if mod <= 0:
        raise ValueError('mod must be positive')
    if exponent < 0:
        raise ValueError('exponent must be non-negative')
    if mod == 1:
        return 0

    base = base % mod
    result = 1
    while exponent > 0:
        if exponent & 1:
            result = (result * base) % mod
        base = (base * base) % mod
        exponent >>= 1
    return result


def sieve_of_eratosthenes(n: int) -> List[int]:
    '''Return the list of prime numbers less than or equal to ``n``.

    Raises:
        TypeError: If ``n`` is not an integer.
        ValueError: If ``n`` is negative.
    '''
    if not isinstance(n, int) or isinstance(n, bool):
        raise TypeError('n must be an integer')
    if n < 0:
        raise ValueError('n must be non-negative')
    if n < 2:
        return []

    sieve = [True] * (n + 1)
    sieve[0] = sieve[1] = False
    for p in range(2, int(n ** 0.5) + 1):
        if sieve[p]:
            for multiple in range(p * p, n + 1, p):
                sieve[multiple] = False
    return [i for i, is_prime in enumerate(sieve) if is_prime]
