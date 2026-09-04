/**
 * Algorithms Library for TypeScript.
 *
 * A zero-dependency collection of common algorithms.
 *
 * @author Parthiv Rawat
 * @license MIT
 */

/* ---------------- Exceptions ---------------- */

export class AlgorithmsError extends Error {
  constructor(message: string) {
    super(message);
    this.name = this.constructor.name;
  }
}

export class EmptyInputError extends AlgorithmsError {}
export class ValueNotFoundError extends AlgorithmsError {}
export class InvalidGraphError extends AlgorithmsError {}
export class NegativeCycleError extends AlgorithmsError {}
export class InvalidInputError extends AlgorithmsError {}

/* ---------------- Shared types ---------------- */

export type Comparable = number | string;
export type AdjacencyList = Record<string, string[]>;
export type WeightedAdjacencyList = Record<string, [string, number][]>;

/**
 * Three-way comparator: negative if `a < b`, zero if equal, positive if `a > b`.
 * Pass a custom comparator to any sort/search to order arbitrary element types.
 */
export type Comparator<T> = (a: T, b: T) => number;

const defaultCompare = <T>(a: T, b: T): number => {
  const x = a as unknown as Comparable;
  const y = b as unknown as Comparable;
  return x < y ? -1 : x > y ? 1 : 0;
};

/* ---------------- Sorting ---------------- */

export function quickSort<T>(
  items: T[],
  compareFn: Comparator<T> = defaultCompare
): T[] {
  if (items.length <= 1) return [...items];
  const pivot = items[Math.floor(items.length / 2)];
  const left = items.filter((x) => compareFn(x, pivot) < 0);
  const middle = items.filter((x) => compareFn(x, pivot) === 0);
  const right = items.filter((x) => compareFn(x, pivot) > 0);
  return [
    ...quickSort(left, compareFn),
    ...middle,
    ...quickSort(right, compareFn),
  ];
}

export function mergeSort<T>(
  items: T[],
  compareFn: Comparator<T> = defaultCompare
): T[] {
  if (items.length <= 1) return [...items];
  const mid = Math.floor(items.length / 2);
  const left = mergeSort(items.slice(0, mid), compareFn);
  const right = mergeSort(items.slice(mid), compareFn);
  return merge(left, right, compareFn);
}

function merge<T>(left: T[], right: T[], compareFn: Comparator<T>): T[] {
  const merged: T[] = [];
  let i = 0;
  let j = 0;
  while (i < left.length && j < right.length) {
    if (compareFn(left[i], right[j]) <= 0) {
      merged.push(left[i]);
      i += 1;
    } else {
      merged.push(right[j]);
      j += 1;
    }
  }
  return merged.concat(left.slice(i), right.slice(j));
}

export function heapSort<T>(
  items: T[],
  compareFn: Comparator<T> = defaultCompare
): T[] {
  const arr = [...items];
  const n = arr.length;
  if (n <= 1) return arr;

  for (let i = Math.floor(n / 2) - 1; i >= 0; i -= 1) {
    siftDown(arr, n, i, compareFn);
  }

  for (let end = n - 1; end > 0; end -= 1) {
    [arr[0], arr[end]] = [arr[end], arr[0]];
    siftDown(arr, end, 0, compareFn);
  }
  return arr;
}

function siftDown<T>(
  arr: T[],
  n: number,
  i: number,
  compareFn: Comparator<T>
): void {
  let largest = i;
  const left = 2 * i + 1;
  const right = 2 * i + 2;

  if (left < n && compareFn(arr[left], arr[largest]) > 0) largest = left;
  if (right < n && compareFn(arr[right], arr[largest]) > 0) largest = right;

  if (largest !== i) {
    [arr[i], arr[largest]] = [arr[largest], arr[i]];
    siftDown(arr, n, largest, compareFn);
  }
}

export function radixSort(items: number[]): number[] {
  let maxVal = 0;
  for (const x of items) {
    if (!Number.isSafeInteger(x) || x < 0) {
      throw new InvalidInputError(
        'radixSort only supports non-negative integers'
      );
    }
    if (x > maxVal) maxVal = x;
  }

  const arr = [...items];
  for (let exp = 1; maxVal / exp >= 1; exp *= 10) {
    countingSortByDigit(arr, exp);
  }
  return arr;
}

function countingSortByDigit(arr: number[], exp: number): void {
  const n = arr.length;
  const output = new Array<number>(n).fill(0);
  const count = new Array<number>(10).fill(0);

  for (const num of arr) {
    count[Math.floor(num / exp) % 10] += 1;
  }
  for (let i = 1; i < 10; i += 1) {
    count[i] += count[i - 1];
  }
  for (let i = n - 1; i >= 0; i -= 1) {
    const index = Math.floor(arr[i] / exp) % 10;
    output[count[index] - 1] = arr[i];
    count[index] -= 1;
  }
  for (let i = 0; i < n; i += 1) {
    arr[i] = output[i];
  }
}

export function nativeSort<T>(
  items: T[],
  compareFn: Comparator<T> = defaultCompare
): T[] {
  return [...items].sort(compareFn);
}

/* ---------------- Searching ---------------- */

function requireSorted<T>(
  arr: T[],
  name: string,
  compareFn: Comparator<T>
): void {
  for (let i = 1; i < arr.length; i += 1) {
    if (compareFn(arr[i], arr[i - 1]) < 0) {
      throw new InvalidInputError(
        `${name} requires a list sorted in ascending order`
      );
    }
  }
}

export function binarySearch<T>(
  arr: T[],
  target: T,
  compareFn: Comparator<T> = defaultCompare
): number {
  if (arr.length === 0) return -1;
  requireSorted(arr, 'binarySearch', compareFn);

  let lo = 0;
  let hi = arr.length - 1;
  while (lo <= hi) {
    const mid = lo + Math.floor((hi - lo) / 2);
    const cmp = compareFn(arr[mid], target);
    if (cmp === 0) return mid;
    if (cmp < 0) lo = mid + 1;
    else hi = mid - 1;
  }
  return -1;
}

export function interpolationSearch(arr: number[], target: number): number {
  if (arr.length === 0) return -1;
  requireSorted(arr, 'interpolationSearch', defaultCompare);

  let lo = 0;
  let hi = arr.length - 1;
  while (
    lo <= hi &&
    target >= arr[lo] &&
    target <= arr[hi]
  ) {
    if (arr[hi] === arr[lo]) {
      if (arr[lo] === target) return lo;
      break;
    }

    const pos =
      lo +
      Math.floor(
        ((target - arr[lo]) / (arr[hi] - arr[lo])) * (hi - lo)
      );
    if (pos < lo || pos > hi) break;

    if (arr[pos] === target) return pos;
    if (arr[pos] < target) lo = pos + 1;
    else hi = pos - 1;
  }
  return -1;
}

export function jumpSearch<T>(
  arr: T[],
  target: T,
  compareFn: Comparator<T> = defaultCompare
): number {
  if (arr.length === 0) return -1;
  requireSorted(arr, 'jumpSearch', compareFn);

  const n = arr.length;
  let step = Math.floor(Math.sqrt(n));
  let prev = 0;

  while (compareFn(arr[Math.min(step, n) - 1], target) < 0) {
    prev = step;
    step += Math.floor(Math.sqrt(n));
    if (prev >= n) return -1;
  }

  while (compareFn(arr[prev], target) < 0) {
    prev += 1;
    if (prev === Math.min(step, n)) return -1;
  }

  return compareFn(arr[prev], target) === 0 ? prev : -1;
}

/* ---------------- Graphs ---------------- */

class _PriorityQueue<T> {
  private heap: { item: T; priority: number }[] = [];

  push(item: T, priority: number): void {
    this.heap.push({ item, priority });
    this.siftUp(this.heap.length - 1);
  }

  pop(): { item: T; priority: number } | undefined {
    if (this.heap.length === 0) return undefined;
    const min = this.heap[0];
    const end = this.heap.pop();
    if (this.heap.length > 0 && end) {
      this.heap[0] = end;
      this.siftDown(0);
    }
    return min;
  }

  isEmpty(): boolean {
    return this.heap.length === 0;
  }

  size(): number {
    return this.heap.length;
  }

  private siftUp(i: number): void {
    while (i > 0) {
      const parent = Math.floor((i - 1) / 2);
      if (this.heap[parent].priority <= this.heap[i].priority) break;
      [this.heap[parent], this.heap[i]] = [this.heap[i], this.heap[parent]];
      i = parent;
    }
  }

  private siftDown(i: number): void {
    const n = this.heap.length;
    while (true) {
      let smallest = i;
      const left = 2 * i + 1;
      const right = 2 * i + 2;
      if (left < n && this.heap[left].priority < this.heap[smallest].priority) {
        smallest = left;
      }
      if (
        right < n &&
        this.heap[right].priority < this.heap[smallest].priority
      ) {
        smallest = right;
      }
      if (smallest === i) break;
      [this.heap[i], this.heap[smallest]] = [this.heap[smallest], this.heap[i]];
      i = smallest;
    }
  }
}

function validateAdjacencyList(graph: AdjacencyList, name: string): void {
  for (const node of Object.keys(graph)) {
    for (const nbr of graph[node] ?? []) {
      if (!(nbr in graph)) {
        throw new InvalidGraphError(
          `${name}: neighbor ${nbr} of ${node} is not a key in the graph`
        );
      }
    }
  }
}

export function bfs(graph: AdjacencyList, start: string): string[] {
  if (!(start in graph)) {
    throw new InvalidGraphError(`Start node ${start} not present in graph`);
  }
  validateAdjacencyList(graph, 'bfs');

  const visited = new Set<string>([start]);
  const order: string[] = [];
  const queue: string[] = [start];
  let head = 0;

  while (head < queue.length) {
    const node = queue[head];
    head += 1;
    order.push(node);
    const neighbors = graph[node] ?? [];
    for (const nbr of neighbors) {
      if (!visited.has(nbr)) {
        visited.add(nbr);
        queue.push(nbr);
      }
    }
  }
  return order;
}

export function dfs(graph: AdjacencyList, start: string): string[] {
  if (!(start in graph)) {
    throw new InvalidGraphError(`Start node ${start} not present in graph`);
  }
  validateAdjacencyList(graph, 'dfs');

  const visited = new Set<string>([start]);
  const order: string[] = [];
  const stack: string[] = [start];

  while (stack.length > 0) {
    const node = stack.pop()!;
    order.push(node);
    const neighbors = graph[node] ?? [];
    for (let i = neighbors.length - 1; i >= 0; i -= 1) {
      const nbr = neighbors[i];
      if (!visited.has(nbr)) {
        visited.add(nbr);
        stack.push(nbr);
      }
    }
  }
  return order;
}

export function dijkstra(
  graph: WeightedAdjacencyList,
  start: string
): { distances: Record<string, number>; predecessors: Record<string, string | undefined> } {
  if (!(start in graph)) {
    throw new InvalidGraphError(`Start node ${start} not present in graph`);
  }

  // Pre-scan all edges so negative weights in unreachable components are
  // detected, and collect every vertex referenced in the graph.
  const allNodes = new Set(Object.keys(graph));
  for (const node of Object.keys(graph)) {
    for (const [nbr, weight] of graph[node] ?? []) {
      if (weight < 0) {
        throw new InvalidGraphError(
          'dijkstra does not support negative edge weights'
        );
      }
      allNodes.add(nbr);
    }
  }

  const distances: Record<string, number> = {};
  for (const node of allNodes) distances[node] = Infinity;
  distances[start] = 0;

  const predecessors: Record<string, string | undefined> = {};
  const pq = new _PriorityQueue<string>();
  pq.push(start, 0);

  while (!pq.isEmpty()) {
    const current = pq.pop();
    if (!current) break;
    const { item: node, priority: dist } = current;
    if (dist > distances[node]) continue;

    const neighbors = graph[node] ?? [];
    for (const [nbr, weight] of neighbors) {
      if (!(nbr in distances)) distances[nbr] = Infinity;
      const nd = dist + weight;
      if (nd < distances[nbr]) {
        distances[nbr] = nd;
        predecessors[nbr] = node;
        pq.push(nbr, nd);
      }
    }
  }

  return { distances, predecessors };
}

export function aStar(
  graph: WeightedAdjacencyList,
  start: string,
  goal: string,
  heuristic: (node: string, goal: string) => number
): { path: string[]; cost: number } {
  if (!(start in graph)) {
    throw new InvalidGraphError(`Start node ${start} not present in graph`);
  }
  if (!(goal in graph)) {
    throw new InvalidGraphError(`Goal node ${goal} not present in graph`);
  }

  const h = heuristic(start, goal);
  if (h < 0) {
    throw new InvalidInputError('heuristic must be non-negative');
  }

  const bestF: Record<string, number> = { [start]: h };
  const gScore: Record<string, number> = { [start]: 0 };
  const cameFrom: Record<string, string> = {};
  const pq = new _PriorityQueue<string>();
  pq.push(start, h);

  while (!pq.isEmpty()) {
    const current = pq.pop();
    if (!current) break;
    const { item: node, priority: f } = current;
    if (f > bestF[node]) continue;

    if (node === goal) {
      const path = reconstructPath(cameFrom, node);
      return { path, cost: gScore[goal] };
    }

    const neighbors = graph[node] ?? [];
    for (const [nbr, weight] of neighbors) {
      if (weight < 0) {
        throw new InvalidGraphError(
          'aStar does not support negative edge weights'
        );
      }
      const hNbr = heuristic(nbr, goal);
      if (hNbr < 0) {
        throw new InvalidInputError('heuristic must be non-negative');
      }
      const tentative = gScore[node] + weight;
      if (!(nbr in gScore) || tentative < gScore[nbr]) {
        cameFrom[nbr] = node;
        gScore[nbr] = tentative;
        const fScore = tentative + hNbr;
        if (!(nbr in bestF) || fScore < bestF[nbr]) {
          bestF[nbr] = fScore;
          pq.push(nbr, fScore);
        }
      }
    }
  }

  throw new ValueNotFoundError(`No path from ${start} to ${goal}`);
}

function reconstructPath(cameFrom: Record<string, string>, current: string): string[] {
  const path: string[] = [current];
  while (current in cameFrom) {
    current = cameFrom[current];
    path.push(current);
  }
  return path.reverse();
}

export function bellmanFord(
  graph: WeightedAdjacencyList,
  start: string
): { distances: Record<string, number>; predecessors: Record<string, string | undefined> } {
  if (!(start in graph)) {
    throw new InvalidGraphError(`Start node ${start} not present in graph`);
  }

  const vertices = new Set<string>(Object.keys(graph));
  const edges: [string, string, number][] = [];
  for (const u of Object.keys(graph)) {
    for (const [v, w] of graph[u]) {
      vertices.add(v);
      edges.push([u, v, w]);
    }
  }

  // Initialize every vertex up front so unreachable nodes appear in the
  // result as Infinity/undefined, consistent with dijkstra.
  const distances: Record<string, number> = {};
  const predecessors: Record<string, string | undefined> = {};
  for (const v of vertices) {
    distances[v] = Infinity;
    predecessors[v] = undefined;
  }
  distances[start] = 0;

  for (let i = 0; i < vertices.size - 1; i += 1) {
    let updated = false;
    for (const [u, v, w] of edges) {
      if (distances[u] + w < distances[v]) {
        distances[v] = distances[u] + w;
        predecessors[v] = u;
        updated = true;
      }
    }
    if (!updated) break;
  }

  for (const [u, v, w] of edges) {
    if (distances[u] + w < distances[v]) {
      throw new NegativeCycleError(
        'Graph contains a negative-weight cycle reachable from the start node'
      );
    }
  }

  return { distances, predecessors };
}

/* ---------------- Dynamic Programming ---------------- */

export function knapsack01(
  weights: number[],
  values: number[],
  capacity: number
): number {
  if (weights.length !== values.length) {
    throw new InvalidInputError(
      'weights and values must have the same length'
    );
  }
  if (!Number.isSafeInteger(capacity) || capacity < 0) {
    throw new InvalidInputError(
      'capacity must be a non-negative integer'
    );
  }

  const dp = new Array<number>(capacity + 1).fill(0);
  for (let i = 0; i < weights.length; i += 1) {
    const w = weights[i];
    const v = values[i];
    if (!Number.isSafeInteger(w) || w < 0) {
      throw new InvalidInputError('weights must be non-negative integers');
    }
    for (let c = capacity; c >= w; c -= 1) {
      dp[c] = Math.max(dp[c], dp[c - w] + v);
    }
  }
  return dp[capacity];
}

export function longestCommonSubsequence<T>(a: T[], b: T[]): number {
  const m = a.length;
  const n = b.length;
  let prev = new Array<number>(n + 1).fill(0);
  for (let i = 1; i <= m; i += 1) {
    const curr = new Array<number>(n + 1).fill(0);
    for (let j = 1; j <= n; j += 1) {
      if (a[i - 1] === b[j - 1]) {
        curr[j] = prev[j - 1] + 1;
      } else {
        curr[j] = Math.max(prev[j], curr[j - 1]);
      }
    }
    prev = curr;
  }
  return prev[n];
}

export function editDistance<T>(a: T[], b: T[]): number {
  const m = a.length;
  const n = b.length;
  let prev = Array.from({ length: n + 1 }, (_, i) => i);
  for (let i = 1; i <= m; i += 1) {
    const curr = new Array<number>(n + 1).fill(0);
    curr[0] = i;
    for (let j = 1; j <= n; j += 1) {
      const cost = a[i - 1] === b[j - 1] ? 0 : 1;
      curr[j] = Math.min(
        curr[j - 1] + 1,      // insertion
        prev[j] + 1,          // deletion
        prev[j - 1] + cost    // substitution
      );
    }
    prev = curr;
  }
  return prev[n];
}

export function longestCommonSubsequenceReconstruction<T>(a: T[], b: T[]): [number, T[]] {
  const m = a.length;
  const n = b.length;
  const dp = Array.from({ length: m + 1 }, () => new Array<number>(n + 1).fill(0));
  for (let i = 1; i <= m; i += 1) {
    for (let j = 1; j <= n; j += 1) {
      if (a[i - 1] === b[j - 1]) {
        dp[i][j] = dp[i - 1][j - 1] + 1;
      } else {
        dp[i][j] = Math.max(dp[i - 1][j], dp[i][j - 1]);
      }
    }
  }
  const result: T[] = [];
  let i = m;
  let j = n;
  while (i > 0 && j > 0) {
    if (a[i - 1] === b[j - 1]) {
      result.push(a[i - 1]);
      i -= 1;
      j -= 1;
    } else if (dp[i - 1][j] >= dp[i][j - 1]) {
      i -= 1;
    } else {
      j -= 1;
    }
  }
  result.reverse();
  return [dp[m][n], result];
}

export type EditOperation<T> =
  | { action: 'match'; from: T; to: T }
  | { action: 'insert'; to: T }
  | { action: 'delete'; from: T }
  | { action: 'substitute'; from: T; to: T };

export function editDistanceReconstruction<T>(a: T[], b: T[]): [number, EditOperation<T>[]] {
  const m = a.length;
  const n = b.length;
  const dp: number[][] = [];
  for (let i = 0; i <= m; i += 1) {
    dp[i] = new Array<number>(n + 1).fill(i);
  }
  for (let j = 1; j <= n; j += 1) {
    dp[0][j] = j;
  }
  for (let i = 1; i <= m; i += 1) {
    for (let j = 1; j <= n; j += 1) {
      if (a[i - 1] === b[j - 1]) {
        dp[i][j] = dp[i - 1][j - 1];
      } else {
        dp[i][j] = 1 + Math.min(
          dp[i - 1][j],
          dp[i][j - 1],
          dp[i - 1][j - 1]
        );
      }
    }
  }
  const script: EditOperation<T>[] = [];
  let i = m;
  let j = n;
  while (i > 0 || j > 0) {
    if (i === 0) {
      script.push({ action: 'insert', to: b[j - 1] });
      j -= 1;
    } else if (j === 0) {
      script.push({ action: 'delete', from: a[i - 1] });
      i -= 1;
    } else if (a[i - 1] === b[j - 1]) {
      script.push({ action: 'match', from: a[i - 1], to: b[j - 1] });
      i -= 1;
      j -= 1;
    } else {
      const best = dp[i][j];
      if (dp[i - 1][j - 1] + 1 === best) {
        script.push({ action: 'substitute', from: a[i - 1], to: b[j - 1] });
        i -= 1;
        j -= 1;
      } else if (dp[i][j - 1] + 1 === best) {
        script.push({ action: 'insert', to: b[j - 1] });
        j -= 1;
      } else {
        script.push({ action: 'delete', from: a[i - 1] });
        i -= 1;
      }
    }
  }
  script.reverse();
  return [dp[m][n], script];
}

/* ---------------- String algorithms ---------------- */

export function kmpSearch(text: string, pattern: string): number[] {
  if (pattern.length === 0) {
    throw new InvalidInputError('pattern must not be empty');
  }

  const failure = computeKmpFailure(pattern);
  const matches: number[] = [];
  let j = 0;
  for (let i = 0; i < text.length; i += 1) {
    while (j > 0 && text[i] !== pattern[j]) {
      j = failure[j - 1];
    }
    if (text[i] === pattern[j]) {
      j += 1;
    }
    if (j === pattern.length) {
      matches.push(i - pattern.length + 1);
      j = failure[j - 1];
    }
  }
  return matches;
}

function computeKmpFailure(pattern: string): number[] {
  const failure = new Array<number>(pattern.length).fill(0);
  let j = 0;
  for (let i = 1; i < pattern.length; i += 1) {
    while (j > 0 && pattern[i] !== pattern[j]) {
      j = failure[j - 1];
    }
    if (pattern[i] === pattern[j]) {
      j += 1;
      failure[i] = j;
    }
  }
  return failure;
}

export function rabinKarpSearch(
  text: string,
  pattern: string,
  base = 256,
  mod = 1_000_000_007
): number[] {
  if (pattern.length === 0) {
    throw new InvalidInputError('pattern must not be empty');
  }
  if (!Number.isSafeInteger(base) || base <= 0) {
    throw new InvalidInputError('base must be a positive integer');
  }
  if (!Number.isSafeInteger(mod) || mod <= 0) {
    throw new InvalidInputError('mod must be a positive integer');
  }
  const n = text.length;
  const m = pattern.length;
  if (m > n) return [];

  let h = 1;
  for (let i = 0; i < m - 1; i += 1) {
    h = (h * base) % mod;
  }

  let patternHash = 0;
  let textHash = 0;
  for (let i = 0; i < m; i += 1) {
    patternHash = (patternHash * base + pattern.charCodeAt(i)) % mod;
    textHash = (textHash * base + text.charCodeAt(i)) % mod;
  }

  const matches: number[] = [];
  for (let i = 0; i <= n - m; i += 1) {
    if (textHash === patternHash && text.slice(i, i + m) === pattern) {
      matches.push(i);
    }
    if (i < n - m) {
      textHash =
        (textHash - text.charCodeAt(i) * h) % mod;
      textHash = (textHash * base + text.charCodeAt(i + m)) % mod;
      textHash = (textHash % mod + mod) % mod;
    }
  }
  return matches;
}

export function boyerMooreSearch(text: string, pattern: string): number[] {
  if (pattern.length === 0) {
    throw new InvalidInputError('pattern must not be empty');
  }
  const n = text.length;
  const m = pattern.length;
  if (m > n) return [];

  const badChar: Record<string, number> = {};
  for (let i = 0; i < m; i += 1) {
    badChar[pattern[i]] = i;
  }
  const goodSuffix = computeGoodSuffixShifts(pattern);

  const matches: number[] = [];
  let i = 0;
  while (i <= n - m) {
    let j = m - 1;
    while (j >= 0 && text[i + j] === pattern[j]) {
      j -= 1;
    }
    if (j < 0) {
      matches.push(i);
      i += goodSuffix[0];
    } else {
      const badCharShift = j - (badChar[text[i + j]] ?? -1);
      i += Math.max(goodSuffix[j + 1], badCharShift, 1);
    }
  }
  return matches;
}

function computeGoodSuffixShifts(pattern: string): number[] {
  const m = pattern.length;
  const shift = new Array<number>(m + 1).fill(0);
  const borderPos = new Array<number>(m + 1).fill(0);
  let i = m;
  let j = m + 1;
  borderPos[i] = j;
  while (i > 0) {
    while (j <= m && pattern[i - 1] !== pattern[j - 1]) {
      if (shift[j] === 0) {
        shift[j] = j - i;
      }
      j = borderPos[j];
    }
    i -= 1;
    j -= 1;
    borderPos[i] = j;
  }
  j = borderPos[0];
  for (let k = 0; k <= m; k += 1) {
    if (shift[k] === 0) {
      shift[k] = j;
    }
    if (k === j) {
      j = borderPos[j];
    }
  }
  return shift;
}

/* ---------------- Greedy ---------------- */

export function activitySelection(activities: [number, number][]): number[] {
  if (activities.length === 0) return [];

  const indexed = activities.map((activity, index) => ({
    index,
    start: activity[0],
    end: activity[1],
  }));

  for (const item of indexed) {
    if (item.start > item.end) {
      throw new InvalidInputError(
        `Activity ${item.index} has start after end`
      );
    }
  }

  indexed.sort((a, b) => a.end - b.end);
  const selected: number[] = [indexed[0].index];
  let lastEnd = indexed[0].end;

  for (let i = 1; i < indexed.length; i += 1) {
    const { index, start, end } = indexed[i];
    if (start >= lastEnd) {
      selected.push(index);
      lastEnd = end;
    }
  }
  return selected;
}

export function fractionalKnapsack(
  weights: number[],
  values: number[],
  capacity: number
): number {
  if (weights.length !== values.length) {
    throw new InvalidInputError(
      'weights and values must have the same length'
    );
  }
  if (capacity < 0 || Number.isNaN(capacity)) {
    throw new InvalidInputError('capacity must be non-negative');
  }

  let total = 0;
  const items: { ratio: number; weight: number }[] = [];
  for (let i = 0; i < weights.length; i += 1) {
    if (
      weights[i] < 0 ||
      Number.isNaN(weights[i]) ||
      Number.isNaN(values[i])
    ) {
      throw new InvalidInputError(
        'weights and values must be valid numbers, weights non-negative'
      );
    }
    if (weights[i] === 0) {
      // Zero-weight items consume no capacity: take them if they add value.
      total += Math.max(0, values[i]);
    } else {
      items.push({ ratio: values[i] / weights[i], weight: weights[i] });
    }
  }
  items.sort((a, b) => b.ratio - a.ratio);

  let remaining = capacity;
  for (const { ratio, weight } of items) {
    if (ratio <= 0) break; // taking a non-positive-value item can only lower the total
    const take = Math.min(weight, remaining);
    total += take * ratio;
    remaining -= take;
    if (remaining <= 0) break;
  }
  return total;
}

interface HuffmanNode {
  freq: number;
  char?: string;
  left?: HuffmanNode;
  right?: HuffmanNode;
}

export function huffmanCoding(
  frequencies: Record<string, number>
): Record<string, string> {
  const entries = Object.entries(frequencies);
  if (entries.length === 0) {
    throw new EmptyInputError('frequencies must not be empty');
  }

  const pq = new _PriorityQueue<HuffmanNode>();
  for (const [char, freq] of entries) {
    if (!Number.isFinite(freq) || freq <= 0) {
      throw new InvalidInputError(
        'huffmanCoding frequencies must be positive finite numbers'
      );
    }
    pq.push({ freq, char }, freq);
  }

  const first = pq.pop();
  if (!first) return {};
  if (pq.isEmpty()) {
    return { [first.item.char ?? '']: '0' };
  }
  pq.push(first.item, first.priority);

  while (pq.size() > 1) {
    const left = pq.pop();
    const right = pq.pop();
    if (!left || !right) break;
    const parent: HuffmanNode = {
      freq: left.priority + right.priority,
      left: left.item,
      right: right.item,
    };
    pq.push(parent, parent.freq);
  }

  const rootEntry = pq.pop();
  if (!rootEntry) return {};
  const root = rootEntry.item;
  const codes: Record<string, string> = {};
  assignHuffmanCodes(root, '', codes);
  return codes;
}

function assignHuffmanCodes(
  node: HuffmanNode,
  prefix: string,
  codes: Record<string, string>
): void {
  if (node.char !== undefined) {
    codes[node.char] = prefix || '0';
    return;
  }
  if (node.left) assignHuffmanCodes(node.left, prefix + '0', codes);
  if (node.right) assignHuffmanCodes(node.right, prefix + '1', codes);
}

export function huffmanEncode(
  symbols: string[],
  codeTable: Record<string, string>
): string {
  const parts: string[] = [];
  for (const s of symbols) {
    const code = codeTable[s];
    if (code === undefined) {
      throw new InvalidInputError(`huffmanEncode: symbol ${s} not in code table`);
    }
    parts.push(code);
  }
  return parts.join('');
}

export function huffmanDecode(
  encoded: string,
  codeTable: Record<string, string>
): string[] {
  const reverse: Record<string, string> = {};
  for (const [s, code] of Object.entries(codeTable)) {
    reverse[code] = s;
  }
  const result: string[] = [];
  let current = '';
  for (const bit of encoded) {
    current += bit;
    if (current in reverse) {
      result.push(reverse[current]);
      current = '';
    }
  }
  if (current.length > 0) {
    throw new InvalidInputError('huffmanDecode: incomplete or invalid bit string');
  }
  return result;
}

/* ---------------- Divide and conquer ---------------- */

export function maxSubarray(arr: number[]): number {
  if (arr.length === 0) {
    throw new EmptyInputError('maxSubarray requires a non-empty array');
  }
  return maxSubarrayDc(arr, 0, arr.length - 1);
}

function maxSubarrayDc(arr: number[], left: number, right: number): number {
  if (left === right) return arr[left];

  const mid = Math.floor((left + right) / 2);
  const leftSum = maxSubarrayDc(arr, left, mid);
  const rightSum = maxSubarrayDc(arr, mid + 1, right);
  const crossSum = maxCrossingSum(arr, left, mid, right);
  return Math.max(leftSum, rightSum, crossSum);
}

function maxCrossingSum(arr: number[], left: number, mid: number, right: number): number {
  let sum = 0;
  let leftMax = -Infinity;
  for (let i = mid; i >= left; i -= 1) {
    sum += arr[i];
    leftMax = Math.max(leftMax, sum);
  }

  sum = 0;
  let rightMax = -Infinity;
  for (let i = mid + 1; i <= right; i += 1) {
    sum += arr[i];
    rightMax = Math.max(rightMax, sum);
  }

  return leftMax + rightMax;
}

export function countInversions<T>(
  arr: T[],
  compareFn: Comparator<T> = defaultCompare
): number {
  const result = sortAndCount([...arr], compareFn);
  return result.count;
}

function sortAndCount<T>(
  arr: T[],
  compareFn: Comparator<T>
): { sorted: T[]; count: number } {
  const n = arr.length;
  if (n <= 1) return { sorted: arr, count: 0 };

  const mid = Math.floor(n / 2);
  const left = sortAndCount(arr.slice(0, mid), compareFn);
  const right = sortAndCount(arr.slice(mid), compareFn);
  const merged = mergeAndCount(left.sorted, right.sorted, compareFn);
  return {
    sorted: merged.sorted,
    count: left.count + right.count + merged.count,
  };
}

function mergeAndCount<T>(
  left: T[],
  right: T[],
  compareFn: Comparator<T>
): { sorted: T[]; count: number } {
  const merged: T[] = [];
  let count = 0;
  let i = 0;
  let j = 0;
  while (i < left.length && j < right.length) {
    if (compareFn(left[i], right[j]) <= 0) {
      merged.push(left[i]);
      i += 1;
    } else {
      merged.push(right[j]);
      j += 1;
      count += left.length - i;
    }
  }
  return {
    sorted: merged.concat(left.slice(i), right.slice(j)),
    count,
  };
}

export function fastPower(base: number, exponent: number): number {
  if (!Number.isSafeInteger(exponent)) {
    throw new InvalidInputError('fastPower exponent must be a safe integer');
  }
  if (base === 0 && exponent < 0) {
    throw new InvalidInputError('0 cannot be raised to a negative power');
  }

  let exp = Math.abs(exponent);
  let result = 1;
  let b = base;
  while (exp > 0) {
    if (exp % 2 === 1) {
      result *= b;
    }
    b *= b;
    exp = Math.floor(exp / 2);
  }
  // By convention, 0^0 returns 1 (matching Math.pow and Python's **).
  return exponent < 0 ? 1 / result : result;
}

/* ---------------- Number theory and selection ---------------- */

export function gcd(a: number, b: number): number {
  if (!Number.isSafeInteger(a) || !Number.isSafeInteger(b)) {
    throw new InvalidInputError('gcd arguments must be safe integers');
  }
  a = Math.abs(a);
  b = Math.abs(b);
  while (b !== 0) {
    const t = a % b;
    a = b;
    b = t;
  }
  return a;
}

export function lcm(a: number, b: number): number {
  if (!Number.isSafeInteger(a) || !Number.isSafeInteger(b)) {
    throw new InvalidInputError('lcm arguments must be safe integers');
  }
  if (a === 0 && b === 0) {
    throw new InvalidInputError('lcm of (0, 0) is undefined');
  }
  const g = gcd(a, b);
  return Math.abs(Math.trunc(a / g) * b);
}

export function modularFastPower(
  base: number,
  exponent: number,
  mod: number
): number {
  if (
    !Number.isSafeInteger(base) ||
    !Number.isSafeInteger(exponent) ||
    !Number.isSafeInteger(mod)
  ) {
    throw new InvalidInputError('modularFastPower arguments must be safe integers');
  }
  if (mod <= 0) {
    throw new InvalidInputError('mod must be positive');
  }
  if (exponent < 0) {
    throw new InvalidInputError('exponent must be non-negative');
  }
  if (mod === 1) {
    return 0;
  }
  let result = 1;
  let b = ((base % mod) + mod) % mod;
  let exp = exponent;
  while (exp > 0) {
    if (exp % 2 === 1) {
      result = (result * b) % mod;
    }
    b = (b * b) % mod;
    exp = Math.floor(exp / 2);
  }
  return result;
}

export function quickSelect<T>(
  items: readonly T[],
  k: number,
  compareFn: Comparator<T> = defaultCompare
): T {
  if (!Number.isSafeInteger(k) || k < 0 || k >= items.length) {
    throw new InvalidInputError('k is out of range');
  }
  let arr = [...items];
  while (true) {
    if (arr.length === 1) {
      return arr[0];
    }
    const pivot = arr[Math.floor(arr.length / 2)];
    const lows = arr.filter((x) => compareFn(x, pivot) < 0);
    const pivots = arr.filter((x) => compareFn(x, pivot) === 0);
    const highs = arr.filter((x) => compareFn(x, pivot) > 0);
    if (k < lows.length) {
      arr = lows;
    } else if (k < lows.length + pivots.length) {
      return pivot;
    } else {
      k -= lows.length + pivots.length;
      arr = highs;
    }
  }
}

/* ---------------- Number theory / graph extras ---------------- */

export function sieveOfEratosthenes(n: number): number[] {
  if (!Number.isSafeInteger(n) || n < 0) {
    throw new InvalidInputError('n must be a non-negative safe integer');
  }
  if (n < 2) {
    return [];
  }
  const sieve = new Array<boolean>(n + 1).fill(true);
  sieve[0] = sieve[1] = false;
  for (let p = 2; p * p <= n; p++) {
    if (sieve[p]) {
      for (let multiple = p * p; multiple <= n; multiple += p) {
        sieve[multiple] = false;
      }
    }
  }
  const primes: number[] = [];
  for (let i = 2; i <= n; i++) {
    if (sieve[i]) {
      primes.push(i);
    }
  }
  return primes;
}

export function topologicalSort(graph: AdjacencyList): string[] {
  validateAdjacencyList(graph, 'topologicalSort');
  const nodes = Object.keys(graph);
  const inDegree: Record<string, number> = {};
  for (const node of nodes) {
    inDegree[node] = 0;
  }
  for (const node of nodes) {
    for (const nbr of graph[node]) {
      inDegree[nbr] = (inDegree[nbr] || 0) + 1;
    }
  }

  const queue = nodes.filter((node) => inDegree[node] === 0);
  const order: string[] = [];
  let head = 0;
  while (head < queue.length) {
    const node = queue[head++];
    order.push(node);
    for (const nbr of graph[node]) {
      inDegree[nbr] -= 1;
      if (inDegree[nbr] === 0) {
        queue.push(nbr);
      }
    }
  }

  if (order.length !== nodes.length) {
    throw new InvalidGraphError('topologicalSort: graph contains a cycle');
  }
  return order;
}

/* ---------------- Data structures ---------------- */

export class UnionFind<T = string> {
  private parent = new Map<T, T>();
  private rank = new Map<T, number>();

  constructor(elements: Iterable<T> = []) {
    for (const x of elements) {
      this.parent.set(x, x);
      this.rank.set(x, 0);
    }
  }

  find(x: T): T {
    if (!this.parent.has(x)) {
      throw new InvalidInputError(`Element ${String(x)} is not in the union-find set`);
    }
    let root = x;
    while (this.parent.get(root) !== root) {
      root = this.parent.get(root)!;
    }
    let node = x;
    while (this.parent.get(node) !== root) {
      const parent = this.parent.get(node)!;
      this.parent.set(node, root);
      node = parent;
    }
    return root;
  }

  union(x: T, y: T): boolean {
    const rootX = this.find(x);
    const rootY = this.find(y);
    if (rootX === rootY) return false;
    const rankX = this.rank.get(rootX)!;
    const rankY = this.rank.get(rootY)!;
    if (rankX < rankY) {
      this.parent.set(rootX, rootY);
    } else if (rankX > rankY) {
      this.parent.set(rootY, rootX);
    } else {
      this.parent.set(rootY, rootX);
      this.rank.set(rootX, rankX + 1);
    }
    return true;
  }

  connected(x: T, y: T): boolean {
    return this.find(x) === this.find(y);
  }

  count(): number {
    const roots = new Set<T>();
    for (const x of this.parent.keys()) {
      roots.add(this.find(x));
    }
    return roots.size;
  }
}

export class PriorityQueue<T> {
  private heap: T[] = [];

  constructor(
    items: Iterable<T> = [],
    private compareFn: Comparator<T> = defaultCompare
  ) {
    this.heap = [...items];
    for (let i = Math.floor(this.heap.length / 2) - 1; i >= 0; i--) {
      this.siftDown(i);
    }
  }

  push(item: T): void {
    this.heap.push(item);
    this.siftUp(this.heap.length - 1);
  }

  pop(): T {
    if (this.heap.length === 0) {
      throw new EmptyInputError('pop from an empty priority queue');
    }
    const top = this.heap[0];
    const last = this.heap.pop()!;
    if (this.heap.length > 0) {
      this.heap[0] = last;
      this.siftDown(0);
    }
    return top;
  }

  peek(): T {
    if (this.heap.length === 0) {
      throw new EmptyInputError('peek from an empty priority queue');
    }
    return this.heap[0];
  }

  get length(): number {
    return this.heap.length;
  }

  private siftUp(i: number): void {
    while (i > 0) {
      const parent = Math.floor((i - 1) / 2);
      if (this.compareFn(this.heap[i], this.heap[parent]) < 0) {
        [this.heap[i], this.heap[parent]] = [this.heap[parent], this.heap[i]];
        i = parent;
      } else {
        break;
      }
    }
  }

  private siftDown(i: number): void {
    const n = this.heap.length;
    while (true) {
      let smallest = i;
      const left = 2 * i + 1;
      const right = 2 * i + 2;
      if (left < n && this.compareFn(this.heap[left], this.heap[smallest]) < 0) {
        smallest = left;
      }
      if (right < n && this.compareFn(this.heap[right], this.heap[smallest]) < 0) {
        smallest = right;
      }
      if (smallest === i) break;
      [this.heap[i], this.heap[smallest]] = [this.heap[smallest], this.heap[i]];
      i = smallest;
    }
  }
}

export function minimumSpanningTree(
  edges: [string, string, number][]
): { total: number; mst: [string, string, number][] } {
  const vertices = new Set<string>();
  for (const [u, v] of edges) {
    vertices.add(u);
    vertices.add(v);
  }
  for (const [, , w] of edges) {
    if (!Number.isFinite(w)) {
      throw new InvalidInputError('edge weights must be finite numbers');
    }
  }

  const sorted = [...edges].sort((a, b) => a[2] - b[2]);
  const uf = new UnionFind(vertices);
  const mst: [string, string, number][] = [];
  let total = 0;
  for (const [u, v, w] of sorted) {
    if (u === v) continue;
    if (uf.union(u, v)) {
      total += w;
      mst.push([u, v, w]);
    }
  }
  return { total, mst };
}

export function floydWarshall(
  graph: WeightedAdjacencyList
): Record<string, Record<string, number>> {
  const allNodes = new Set<string>(Object.keys(graph));
  for (const node of Object.keys(graph)) {
    for (const [nbr] of graph[node]) {
      allNodes.add(nbr);
    }
  }
  for (const node of Object.keys(graph)) {
    for (const [nbr] of graph[node]) {
      if (!allNodes.has(nbr)) {
        throw new InvalidGraphError(
          `floydWarshall: neighbor ${nbr} of node ${node} not present in graph`
        );
      }
    }
  }

  const dist: Record<string, Record<string, number>> = {};
  for (const u of allNodes) {
    dist[u] = {};
    for (const v of allNodes) {
      dist[u][v] = u === v ? 0 : Infinity;
    }
  }
  for (const u of Object.keys(graph)) {
    for (const [nbr, weight] of graph[u]) {
      if (weight < dist[u][nbr]) {
        dist[u][nbr] = weight;
      }
    }
  }

  for (const k of Object.keys(dist)) {
    for (const i of Object.keys(dist)) {
      const dik = dist[i][k];
      if (!Number.isFinite(dik)) continue;
      for (const j of Object.keys(dist)) {
        const nd = dik + dist[k][j];
        if (nd < dist[i][j]) {
          dist[i][j] = nd;
        }
      }
    }
  }

  for (const u of Object.keys(dist)) {
    if (dist[u][u] < 0) {
      throw new NegativeCycleError(
        'floydWarshall: graph contains a negative-weight cycle'
      );
    }
  }
  return dist;
}

export function stronglyConnectedComponents(
  graph: AdjacencyList
): string[][] {
  validateAdjacencyList(graph, 'stronglyConnectedComponents');

  const visited = new Set<string>();
  const order: string[] = [];

  function dfs1(start: string): void {
    const stack: [string, boolean][] = [[start, false]];
    while (stack.length > 0) {
      const [node, processed] = stack.pop()!;
      if (processed) {
        order.push(node);
        continue;
      }
      if (visited.has(node)) continue;
      visited.add(node);
      stack.push([node, true]);
      const neighbors = graph[node];
      for (let i = neighbors.length - 1; i >= 0; i--) {
        const nbr = neighbors[i];
        if (!visited.has(nbr)) {
          stack.push([nbr, false]);
        }
      }
    }
  }

  for (const node of Object.keys(graph)) {
    if (!visited.has(node)) {
      dfs1(node);
    }
  }

  const reverse: Record<string, string[]> = {};
  for (const node of Object.keys(graph)) {
    reverse[node] = [];
  }
  for (const [node, neighbors] of Object.entries(graph)) {
    for (const nbr of neighbors) {
      reverse[nbr].push(node);
    }
  }

  visited.clear();
  const components: string[][] = [];

  function dfs2(start: string): string[] {
    const stack = [start];
    const component: string[] = [];
    while (stack.length > 0) {
      const node = stack.pop()!;
      if (visited.has(node)) continue;
      visited.add(node);
      component.push(node);
      for (const nbr of reverse[node]) {
        if (!visited.has(nbr)) {
          stack.push(nbr);
        }
      }
    }
    return component;
  }

  for (let i = order.length - 1; i >= 0; i--) {
    const node = order[i];
    if (!visited.has(node)) {
      components.push(dfs2(node));
    }
  }
  return components;
}

export function hasCycle(graph: AdjacencyList): boolean {
  validateAdjacencyList(graph, 'hasCycle');
  const state = new Map<string, number>();
  for (const u of Object.keys(graph)) {
    state.set(u, 0);
  }
  for (const start of Object.keys(graph)) {
    if (state.get(start) !== 0) continue;
    const stack: [string, number][] = [[start, 0]];
    while (stack.length > 0) {
      const [node, idx] = stack[stack.length - 1];
      if (state.get(node) === 2) {
        stack.pop();
        continue;
      }
      if (state.get(node) === 0) {
        state.set(node, 1);
      }
      const neighbors = graph[node];
      if (idx < neighbors.length) {
        const nbr = neighbors[idx];
        stack[stack.length - 1] = [node, idx + 1];
        if (state.get(nbr) === 1) return true;
        if (state.get(nbr) === 0) {
          stack.push([nbr, 0]);
        }
      } else {
        state.set(node, 2);
        stack.pop();
      }
    }
  }
  return false;
}

export function binarySearchOnAnswer(
  low: number,
  high: number,
  predicate: (x: number) => boolean,
  find: 'minimum' | 'maximum' = 'minimum'
): number | null {
  if (low > high) return null;
  let left = low;
  let right = high;
  let answer: number | null = null;
  if (find === 'minimum') {
    while (left <= right) {
      const mid = Math.floor((left + right) / 2);
      if (predicate(mid)) {
        answer = mid;
        right = mid - 1;
      } else {
        left = mid + 1;
      }
    }
  } else {
    while (left <= right) {
      const mid = Math.floor((left + right) / 2);
      if (predicate(mid)) {
        answer = mid;
        left = mid + 1;
      } else {
        right = mid - 1;
      }
    }
  }
  return answer;
}
