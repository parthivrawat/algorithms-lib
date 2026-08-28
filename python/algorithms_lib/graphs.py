'''Graph traversal and shortest-path algorithms.'''

import heapq
from collections import deque
from typing import Any, Callable, Dict, Iterable, List, Mapping, Tuple, Union

from .exceptions import InvalidGraphError, NegativeCycleError, ValueNotFoundError

AdjacencyList = Mapping[Any, Iterable[Any]]
WeightedAdjacencyList = Mapping[Any, Iterable[Tuple[Any, Union[int, float]]]]


def _neighbors(graph: Mapping[Any, Any], node: Any) -> Iterable[Any]:
    return graph.get(node, ())


def bfs(graph: AdjacencyList, start: Any) -> List[Any]:
    '''Return a list of nodes in breadth-first order.

    ``graph`` is an adjacency list mapping a node to an iterable of neighbors.
    '''
    if start not in graph:
        raise InvalidGraphError(f'Start node {start!r} not present in graph')

    visited = {start}
    order: List[Any] = []
    queue: deque = deque([start])

    while queue:
        node = queue.popleft()
        order.append(node)
        for nbr in _neighbors(graph, node):
            if nbr not in visited:
                visited.add(nbr)
                queue.append(nbr)
    return order


def dfs(graph: AdjacencyList, start: Any) -> List[Any]:
    '''Return a list of nodes in depth-first order.

    An iterative implementation is used to avoid recursion-depth limits.
    '''
    if start not in graph:
        raise InvalidGraphError(f'Start node {start!r} not present in graph')

    visited = {start}
    order: List[Any] = []
    stack = [start]

    while stack:
        node = stack.pop()
        order.append(node)
        # Reverse iteration preserves the original neighbor order.
        for nbr in reversed(list(_neighbors(graph, node))):
            if nbr not in visited:
                visited.add(nbr)
                stack.append(nbr)
    return order


def dijkstra(
    graph: WeightedAdjacencyList, start: Any
) -> Tuple[Dict[Any, float], Dict[Any, Any]]:
    '''Compute shortest distances from ``start`` using Dijkstra\'s algorithm.

    Returns a tuple ``(distances, predecessors)``.

    Raises:
        InvalidGraphError: If a negative edge weight is found.
    '''
    if start not in graph:
        raise InvalidGraphError(f'Start node {start!r} not present in graph')

    distances: Dict[Any, float] = {node: float('inf') for node in graph}
    distances[start] = 0
    predecessors: Dict[Any, Any] = {}
    heap = [(0.0, start)]

    while heap:
        dist, node = heapq.heappop(heap)
        if dist != distances[node]:
            continue

        for nbr, weight in _neighbors(graph, node):
            if weight < 0:
                raise InvalidGraphError('Dijkstra does not support negative edge weights')
            if nbr not in distances:
                distances[nbr] = float('inf')
            nd = dist + weight
            if nd < distances[nbr]:
                distances[nbr] = nd
                predecessors[nbr] = node
                heapq.heappush(heap, (nd, nbr))
    return distances, predecessors


def a_star(
    graph: WeightedAdjacencyList,
    start: Any,
    goal: Any,
    heuristic: Callable[[Any, Any], float],
) -> Tuple[List[Any], float]:
    '''Find a shortest path from ``start`` to ``goal`` using A* search.

    ``heuristic(node, goal)`` must return a non-negative estimate of the
    remaining cost. Returns ``(path, cost)``.

    Raises:
        InvalidGraphError: If the start or goal is not in the graph.
        ValueNotFoundError: If no path exists.
    '''
    if start not in graph:
        raise InvalidGraphError(f'Start node {start!r} not present in graph')
    if goal not in graph:
        raise InvalidGraphError(f'Goal node {goal!r} not present in graph')

    open_set = [(0.0, start)]
    g_score: Dict[Any, float] = {start: 0.0}
    came_from: Dict[Any, Any] = {}

    while open_set:
        _, current = heapq.heappop(open_set)
        if current == goal:
            path = _reconstruct_path(came_from, current)
            return path, g_score[goal]

        for nbr, weight in _neighbors(graph, current):
            tentative = g_score[current] + weight
            if nbr not in g_score or tentative < g_score[nbr]:
                came_from[nbr] = current
                g_score[nbr] = tentative
                f_score = tentative + heuristic(nbr, goal)
                heapq.heappush(open_set, (f_score, nbr))

    raise ValueNotFoundError(f'No path from {start!r} to {goal!r}')


def _reconstruct_path(came_from: Dict[Any, Any], current: Any) -> List[Any]:
    path = [current]
    while current in came_from:
        current = came_from[current]
        path.append(current)
    return list(reversed(path))


def bellman_ford(
    graph: WeightedAdjacencyList, start: Any
) -> Tuple[Dict[Any, float], Dict[Any, Any]]:
    '''Compute shortest distances from ``start`` using the Bellman-Ford algorithm.

    Returns ``(distances, predecessors)``. Unlike Dijkstra, Bellman-Ford
    supports negative edge weights and detects negative cycles.

    Raises:
        InvalidGraphError: If the start node is not in the graph.
        NegativeCycleError: If a reachable negative-weight cycle exists.
    '''
    if start not in graph:
        raise InvalidGraphError(f'Start node {start!r} not present in graph')

    vertices = set(graph.keys())
    edges: List[Tuple[Any, Any, Union[int, float]]] = []
    for u in graph:
        for v, w in _neighbors(graph, u):
            vertices.add(v)
            edges.append((u, v, w))

    distances: Dict[Any, float] = {start: 0.0}
    predecessors: Dict[Any, Any] = {}

    for _ in range(len(vertices) - 1):
        updated = False
        for u, v, w in edges:
            if u in distances and distances[u] + w < distances.get(v, float('inf')):
                distances[v] = distances[u] + w
                predecessors[v] = u
                updated = True
        if not updated:
            break

    for u, v, w in edges:
        if u in distances and distances[u] + w < distances.get(v, float('inf')):
            raise NegativeCycleError(
                'Graph contains a negative-weight cycle reachable from the start node'
            )

    return distances, predecessors
