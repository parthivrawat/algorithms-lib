'''Custom exceptions for algorithms-lib.'''


class AlgorithmsError(Exception):
    '''Base exception for all algorithms-lib errors.'''


class EmptyInputError(AlgorithmsError, ValueError):
    '''Raised when an algorithm receives an empty input that requires values.'''


class ValueNotFoundError(AlgorithmsError, ValueError):
    '''Raised when a searched value or path cannot be found.'''


class InvalidGraphError(AlgorithmsError, ValueError):
    '''Raised when a graph algorithm receives an invalid graph or start node.'''


class NegativeCycleError(AlgorithmsError, ValueError):
    '''Raised when a shortest-path algorithm detects a negative-weight cycle.'''
