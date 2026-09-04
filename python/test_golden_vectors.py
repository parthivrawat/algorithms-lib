'''Consume the shared golden vectors and verify the Python implementation.'''

import json
import math
from pathlib import Path
from typing import Any

import pytest

import algorithms_lib


GOLDEN_VECTORS = (
    Path(__file__).resolve().parents[1] / 'tests' / 'golden_vectors.json'
)


def _approx(expected: Any, actual: Any) -> bool:
    '''Compare expected and actual, allowing a small float tolerance.'''
    if isinstance(expected, float):
        return math.isclose(actual, expected, rel_tol=1e-9, abs_tol=1e-9)
    if isinstance(expected, list):
        if not isinstance(actual, list) or len(actual) != len(expected):
            return False
        return all(_approx(e, a) for e, a in zip(expected, actual))
    return actual == expected


def _load_vectors() -> list[dict[str, Any]]:
    with GOLDEN_VECTORS.open('r', encoding='utf-8') as f:
        data = json.load(f)
    return data['vectors']


@pytest.mark.parametrize(
    'vector',
    _load_vectors(),
    ids=lambda v: f"{v['function']} - {v.get('name', v['function'])}",
)
def test_golden_vector(vector: dict[str, Any]) -> None:
    fn = getattr(algorithms_lib, vector['function'])
    actual = fn(**vector['input'])
    assert _approx(vector['expected'], actual), (
        f"{vector['function']} returned {actual!r}, expected {vector['expected']!r}"
    )
