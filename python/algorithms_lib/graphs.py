'''Graph traversal and shortest-path algorithms.'''

import heapq
import math
from collections import deque
from typing import Any, Callable, Dict, Generic, Iterable, List, Mapping, Set, Tuple, TypeVar, Union

from .exceptions import (
    InvalidGraphError,
    InvalidInputError,
    NegativeCycleError,
    ValueNotFoundError,
)

AdjacencyList = Mapping[Any, Iterable[Any]]
WeightedAdjacencyList = Mapping[Any, Iterable[Tuple[Any, Union[int, float]]]]


def _neighbors(graph: Mapping[Any, Any], node: Any) -> Iterable[Any]:
    return graph.get(node, ())


def _validate_adjacency_list(graph: Mapping[Any, Any], name: str) -> None:
    for node, neighbors in graph.items():
        for nbr in neighbors:
            if nbr not in graph:
                raise InvalidGraphError(
                    f'{name}: neighbor {nbr!r} of node {node!r} not present in graph'
                )


def bfs(graph: AdjacencyList, start: Any) -> List[Any]:
    '''Return a list of nodes in breadth-first order.

    ``graph`` is an adjacency list mapping a node to an iterable of neighbors.
    Every neighbor referenced must appear as a key in ``graph``.
    '''
    if start not in graph:
        raise InvalidGraphError(f'Start node {start!r} not present in graph')
    _validate_adjacency_list(graph, 'bfs')

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
    Every neighbor referenced must appear as a key in ``graph``.
    '''
    if start not in graph:
        raise InvalidGraphError(f'Start node {start!r} not present in graph')
    _validate_adjacency_list(graph, 'dfs')

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

    # Pre-scan all edges so negative weights in unreachable components are
    # detected, and collect every vertex referenced in the graph.
    all_nodes = set(graph)
    for node, edges in graph.items():
        for nbr, weight in edges:
            if weight < 0:
                raise InvalidGraphError(
                    'Dijkstra does not support negative edge weights'
                )
            all_nodes.add(nbr)

    distances: Dict[Any, float] = {node: float('inf') for node in all_nodes}
    distances[start] = 0
    predecessors: Dict[Any, Any] = {}
    heap = [(0.0, start)]

    while heap:
        dist, node = heapq.heappop(heap)
        if dist != distances[node]:
            continue

        for nbr, weight in _neighbors(graph, node):
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
    remaining cost. For an optimal first result the heuristic should be
    admissible and consistent; inconsistent but admissible heuristics still
    return an optimal path, but may cause nodes to be re-expanded.

    All edge weights must be non-negative.

    Raises:
        InvalidGraphError: If the start or goal is not in the graph, or if an
            edge has a negative weight.
        ValueError: If the heuristic returns a negative value.
        ValueNotFoundError: If no path exists.
    '''
    if start not in graph:
        raise InvalidGraphError(f'Start node {start!r} not present in graph')
    if goal not in graph:
        raise InvalidGraphError(f'Goal node {goal!r} not present in graph')

    h = heuristic(start, goal)
    if h < 0:
        raise InvalidInputError('heuristic must be non-negative')

    open_set = [(h, start)]
    best_f: Dict[Any, float] = {start: h}
    g_score: Dict[Any, float] = {start: 0.0}
    came_from: Dict[Any, Any] = {}

    while open_set:
        f, current = heapq.heappop(open_set)
        if f > best_f[current]:
            continue

        if current == goal:
            path = _reconstruct_path(came_from, current)
            return path, g_score[goal]

        for nbr, weight in _neighbors(graph, current):
            if weight < 0:
                raise InvalidGraphError(
                    'a_star does not support negative edge weights'
                )
            h_nbr = heuristic(nbr, goal)
            if h_nbr < 0:
                raise InvalidInputError('heuristic must be non-negative')

            tentative = g_score[current] + weight
            if nbr not in g_score or tentative < g_score[nbr]:
                came_from[nbr] = current
                g_score[nbr] = tentative
                f_score = tentative + h_nbr
                if nbr not in best_f or f_score < best_f[nbr]:
                    best_f[nbr] = f_score
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
    supports negative edge weights and detects negative cycles that are
    reachable from ``start``. Cycles in components not reachable from
    ``start`` are not reported.

    Raises:
        InvalidGraphError: If the start node is not in the graph.
        NegativeCycleError: If a negative-weight cycle reachable from
            ``start`` exists.
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


def topological_sort(graph: AdjacencyList) -> List[Any]:
    '''Return a topological ordering of the directed acyclic ``graph``.

    ``graph`` is an adjacency list mapping each node to its dependents.
    Returns a list of nodes in which every node appears before any node
    that depends on it.

    Raises:
        InvalidGraphError: If the graph contains a cycle, or if a neighbor
            referenced by a node is not present as a key in the graph.
    '''
    _validate_adjacency_list(graph, 'topological_sort')

    in_degree: Dict[Any, int] = {node: 0 for node in graph}
    for node, neighbors in graph.items():
        for nbr in neighbors:
            in_degree[nbr] += 1

    queue: deque = deque(
        [node for node, degree in in_degree.items() if degree == 0]
    )
    order: List[Any] = []
    while queue:
        node = queue.popleft()
        order.append(node)
        for nbr in _neighbors(graph, node):
            in_degree[nbr] -= 1
            if in_degree[nbr] == 0:
                queue.append(nbr)

    if len(order) != len(in_degree):
        raise InvalidGraphError(
            'topological_sort: graph contains a cycle'
        )
    return order


T = TypeVar('T')


class UnionFind(Generic[T]):
    '''Disjoint-set union–find with union by rank and path compression.

    ``elements`` is the initial set of elements. Use ``find``, ``union``,
    ``connected`` and ``count`` to query or mutate the partition.
    '''

    def __init__(self, elements: Iterable[T] = ()):
        self._parent: Dict[T, T] = {x: x for x in elements}
        self._rank: Dict[T, int] = {x: 0 for x in elements}

    def find(self, x: T) -> T:
        '''Return the representative of the set containing ``x``.'''
        if x not in self._parent:
            raise InvalidInputError(f'{x!r} is not in the union-find set')
        root = x
        while self._parent[root] != root:
            root = self._parent[root]
        node = x
        while self._parent[node] != root:
            parent = self._parent[node]
            self._parent[node] = root
            node = parent
        return root

    def union(self, x: T, y: T) -> bool:
        '''Merge the sets containing ``x`` and ``y``.

        Returns ``True`` if a merge happened, ``False`` if they were
        already in the same set.
        '''
        root_x = self.find(x)
        root_y = self.find(y)
        if root_x == root_y:
            return False
        if self._rank[root_x] < self._rank[root_y]:
            self._parent[root_x] = root_y
        elif self._rank[root_x] > self._rank[root_y]:
            self._parent[root_y] = root_x
        else:
            self._parent[root_y] = root_x
            self._rank[root_x] += 1
        return True

    def connected(self, x: T, y: T) -> bool:
        '''Return whether ``x`` and ``y`` are in the same set.'''
        return self.find(x) == self.find(y)

    def count(self) -> int:
        '''Return the current number of disjoint sets.'''
        roots = set()
        for x in self._parent:
            roots.add(self.find(x))
        return len(roots)


def minimum_spanning_tree(
    edges: List[Tuple[Any, Any, Union[int, float]]]
) -> Tuple[float, List[Tuple[Any, Any, float]]]:
    '''Return the total weight and edges of a minimum spanning forest.

    ``edges`` is a list of ``(u, v, weight)`` tuples. The returned edges
    form a minimum spanning tree for each connected component of the graph.
    Self-loops are ignored.

    Raises:
        InvalidInputError: If any weight is not a finite number.
    '''
    vertices = set()
    for u, v, _ in edges:
        vertices.add(u)
        vertices.add(v)

    for u, v, w in edges:
        if isinstance(w, bool) or not isinstance(w, (int, float)):
            raise InvalidInputError('edge weights must be numbers')
        if math.isnan(w) or math.isinf(w):
            raise InvalidInputError('edge weights must be finite numbers')

    uf = UnionFind(vertices)
    sorted_edges = sorted(edges, key=lambda e: e[2])
    total = 0.0
    mst: List[Tuple[Any, Any, float]] = []
    for u, v, w in sorted_edges:
        if u == v:
            continue
        if uf.union(u, v):
            total += float(w)
            mst.append((u, v, float(w)))
    return total, mst


def floyd_warshall(
    graph: WeightedAdjacencyList
) -> Dict[Any, Dict[Any, float]]:
    '''Return the all-pairs shortest-path distance matrix.

    ``graph`` is a weighted adjacency list. Missing edges have distance
    ``float('inf')``. The diagonal is ``0`` for every vertex.

    Raises:
        NegativeCycleError: If the graph contains a negative-weight cycle
            reachable from some vertex.
    '''
    all_nodes = set(graph)
    for node, edges in graph.items():
        for nbr, _ in edges:
            all_nodes.add(nbr)
    for node, edges in graph.items():
        for nbr, _ in edges:
            if nbr not in all_nodes:
                raise InvalidGraphError(
                    f'floyd_warshall: neighbor {nbr!r} of node {node!r} not present in graph'
                )

    dist: Dict[Any, Dict[Any, float]] = {
        u: {v: float('inf') for v in all_nodes} for u in all_nodes
    }
    for u in all_nodes:
        dist[u][u] = 0.0
    for u, edges in graph.items():
        for nbr, weight in edges:
            if weight < dist[u][nbr]:
                dist[u][nbr] = float(weight)

    for k in all_nodes:
        for i in all_nodes:
            dik = dist[i][k]
            if dik == float('inf'):
                continue
            for j in all_nodes:
                nd = dik + dist[k][j]
                if nd < dist[i][j]:
                    dist[i][j] = nd

    for u in all_nodes:
        if dist[u][u] < 0:
            raise NegativeCycleError(
                'floyd_warshall: graph contains a negative-weight cycle'
            )
    return dist


def strongly_connected_components(
    graph: AdjacencyList
) -> List[List[Any]]:
    '''Return the strongly connected components of ``graph`` using Kosaraju.

    The components are returned in reverse topological order.
    '''
    _validate_adjacency_list(graph, 'strongly_connected_components')

    visited: Set[Any] = set()
    order: List[Any] = []

    def _dfs1(node: Any) -> None:
        stack: List[Tuple[Any, bool]] = [(node, False)]
        while stack:
            current, processed = stack.pop()
            if processed:
                order.append(current)
                continue
            if current in visited:
                continue
            visited.add(current)
            stack.append((current, True))
            for nbr in reversed(list(_neighbors(graph, current))):
                if nbr not in visited:
                    stack.append((nbr, False))

    for node in graph:
        if node not in visited:
            _dfs1(node)

    reverse: Dict[Any, List[Any]] = {node: [] for node in graph}
    for node, neighbors in graph.items():
        for nbr in neighbors:
            reverse[nbr].append(node)

    visited.clear()
    components: List[List[Any]] = []

    def _dfs2(node: Any) -> List[Any]:
        stack = [node]
        component: List[Any] = []
        while stack:
            current = stack.pop()
            if current in visited:
                continue
            visited.add(current)
            component.append(current)
            for nbr in reverse[current]:
                if nbr not in visited:
                    stack.append(nbr)
        return component

    for node in reversed(order):
        if node not in visited:
            components.append(_dfs2(node))
    return components


def has_cycle(graph: AdjacencyList) -> bool:
    '''Return whether the directed ``graph`` contains a cycle.

    Raises:
        InvalidGraphError: If a neighbor referenced in the graph does not
            appear as a key.
    '''
    _validate_adjacency_list(graph, 'has_cycle')
    state: Dict[Any, int] = {node: 0 for node in graph}
    for start in graph:
        if state[start] != 0:
            continue
        stack: List[Tuple[Any, int]] = [(start, 0)]
        while stack:
            node, idx = stack[-1]
            if state[node] == 2:
                stack.pop()
                continue
            if state[node] == 0:
                state[node] = 1
            neighbors = list(_neighbors(graph, node))
            if idx < len(neighbors):
                nbr = neighbors[idx]
                stack[-1] = (node, idx + 1)
                if state[nbr] == 1:
                    return True
                if state[nbr] == 0:
                    stack.append((nbr, 0))
            else:
                state[node] = 2
                stack.pop()
    return False
