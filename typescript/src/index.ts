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

/* ---------------- Shared types ---------------- */

export type Comparable = number | string;
export type AdjacencyList = Record<string, string[]>;
export type WeightedAdjacencyList = Record<string, [string, number][]>;

/* ---------------- Sorting ---------------- */

export function quickSort<T extends Comparable>(items: T[]): T[] {
  if (items.length <= 1) return [...items];
  const pivot = items[Math.floor(items.length / 2)];
  const left = items.filter((x) => x < pivot);
  const middle = items.filter((x) => x === pivot);
  const right = items.filter((x) => x > pivot);
  return [...quickSort(left), ...middle, ...quickSort(right)];
}

export function mergeSort<T extends Comparable>(items: T[]): T[] {
  if (items.length <= 1) return [...items];
  const mid = Math.floor(items.length / 2);
  const left = mergeSort(items.slice(0, mid));
  const right = mergeSort(items.slice(mid));
  return merge(left, right);
}

function merge<T extends Comparable>(left: T[], right: T[]): T[] {
  const merged: T[] = [];
  let i = 0;
  let j = 0;
  while (i < left.length && j < right.length) {
    if (left[i] <= right[j]) {
      merged.push(left[i]);
      i += 1;
    } else {
      merged.push(right[j]);
      j += 1;
    }
  }
  return merged.concat(left.slice(i), right.slice(j));
}

export function heapSort<T extends Comparable>(items: T[]): T[] {
  const arr = [...items];
  const n = arr.length;
  if (n <= 1) return arr;

  for (let i = Math.floor(n / 2) - 1; i >= 0; i -= 1) {
    siftDown(arr, n, i);
  }

  for (let end = n - 1; end > 0; end -= 1) {
    [arr[0], arr[end]] = [arr[end], arr[0]];
    siftDown(arr, end, 0);
  }
  return arr;
}

function siftDown<T extends Comparable>(arr: T[], n: number, i: number): void {
  let largest = i;
  const left = 2 * i + 1;
  const right = 2 * i + 2;

  if (left < n && arr[left] > arr[largest]) largest = left;
  if (right < n && arr[right] > arr[largest]) largest = right;

  if (largest !== i) {
    [arr[i], arr[largest]] = [arr[largest], arr[i]];
    siftDown(arr, n, largest);
  }
}

export function radixSort(items: number[]): number[] {
  if (items.length === 0) return [];
  if (items.some((x) => x < 0)) {
    throw new Error('radixSort only supports non-negative integers');
  }

  const arr = [...items];
  let maxVal = Math.max(...arr);
  for (let exp = 1; Math.floor(maxVal / exp) > 0; exp *= 10) {
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

export function timSort<T extends Comparable>(items: T[]): T[] {
  return [...items].sort((a, b) => (a < b ? -1 : a > b ? 1 : 0));
}

/* ---------------- Searching ---------------- */

function requireSorted<T extends Comparable>(arr: T[], name: string): void {
  for (let i = 1; i < arr.length; i += 1) {
    if (arr[i] < arr[i - 1]) {
      throw new Error(`${name} requires a list sorted in ascending order`);
    }
  }
}

export function binarySearch<T extends Comparable>(
  arr: T[],
  target: T
): number {
  if (arr.length === 0) return -1;
  requireSorted(arr, 'binarySearch');

  let lo = 0;
  let hi = arr.length - 1;
  while (lo <= hi) {
    const mid = Math.floor((lo + hi) / 2);
    if (arr[mid] === target) return mid;
    if (arr[mid] < target) lo = mid + 1;
    else hi = mid - 1;
  }
  return -1;
}

export function interpolationSearch(arr: number[], target: number): number {
  if (arr.length === 0) return -1;
  requireSorted(arr, 'interpolationSearch');

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

export function jumpSearch<T extends Comparable>(arr: T[], target: T): number {
  if (arr.length === 0) return -1;
  requireSorted(arr, 'jumpSearch');

  const n = arr.length;
  let step = Math.floor(Math.sqrt(n));
  let prev = 0;

  while (arr[Math.min(step, n) - 1] < target) {
    prev = step;
    step += Math.floor(Math.sqrt(n));
    if (prev >= n) return -1;
  }

  while (arr[prev] < target) {
    prev += 1;
    if (prev === Math.min(step, n)) return -1;
  }

  return arr[prev] === target ? prev : -1;
}

/* ---------------- Graphs ---------------- */

class PriorityQueue<T> {
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

export function bfs(graph: AdjacencyList, start: string): string[] {
  if (!(start in graph)) {
    throw new InvalidGraphError(`Start node ${start} not present in graph`);
  }

  const visited = new Set<string>([start]);
  const order: string[] = [];
  const queue: string[] = [start];

  while (queue.length > 0) {
    const node = queue.shift()!;
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

  const distances: Record<string, number> = {};
  for (const node of Object.keys(graph)) distances[node] = Infinity;
  distances[start] = 0;

  const predecessors: Record<string, string | undefined> = {};
  const pq = new PriorityQueue<string>();
  pq.push(start, 0);

  while (!pq.isEmpty()) {
    const current = pq.pop();
    if (!current) break;
    const { item: node, priority: dist } = current;
    if (dist !== distances[node]) continue;

    const neighbors = graph[node] ?? [];
    for (const [nbr, weight] of neighbors) {
      if (weight < 0) {
        throw new InvalidGraphError('dijkstra does not support negative edge weights');
      }
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

  const gScore: Record<string, number> = { [start]: 0 };
  const cameFrom: Record<string, string> = {};
  const pq = new PriorityQueue<string>();
  pq.push(start, heuristic(start, goal));

  while (!pq.isEmpty()) {
    const current = pq.pop();
    if (!current) break;
    const node = current.item;

    if (node === goal) {
      const path = reconstructPath(cameFrom, node);
      return { path, cost: gScore[goal] };
    }

    const neighbors = graph[node] ?? [];
    for (const [nbr, weight] of neighbors) {
      const tentative = gScore[node] + weight;
      if (!(nbr in gScore) || tentative < gScore[nbr]) {
        cameFrom[nbr] = node;
        gScore[nbr] = tentative;
        const fScore = tentative + heuristic(nbr, goal);
        pq.push(nbr, fScore);
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

  const distances: Record<string, number> = { [start]: 0 };
  const predecessors: Record<string, string | undefined> = {};

  for (let i = 0; i < vertices.size - 1; i += 1) {
    let updated = false;
    for (const [u, v, w] of edges) {
      if (u in distances && distances[u] + w < (distances[v] ?? Infinity)) {
        distances[v] = distances[u] + w;
        predecessors[v] = u;
        updated = true;
      }
    }
    if (!updated) break;
  }

  for (const [u, v, w] of edges) {
    if (u in distances && distances[u] + w < (distances[v] ?? Infinity)) {
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
    throw new Error('weights and values must have the same length');
  }
  if (capacity < 0) {
    throw new Error('capacity must be non-negative');
  }

  const dp = new Array<number>(capacity + 1).fill(0);
  for (let i = 0; i < weights.length; i += 1) {
    const w = weights[i];
    const v = values[i];
    if (w < 0) throw new Error('weights must be non-negative');
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

/* ---------------- String algorithms ---------------- */

export function kmpSearch(text: string, pattern: string): number[] {
  if (pattern.length === 0) {
    throw new Error('pattern must not be empty');
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
    throw new Error('pattern must not be empty');
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
    throw new Error('pattern must not be empty');
  }
  const n = text.length;
  const m = pattern.length;
  if (m > n) return [];

  const badChar: Record<string, number> = {};
  for (let i = 0; i < m; i += 1) {
    badChar[pattern[i]] = i;
  }

  const matches: number[] = [];
  let i = 0;
  while (i <= n - m) {
    let j = m - 1;
    while (j >= 0 && text[i + j] === pattern[j]) {
      j -= 1;
    }
    if (j < 0) {
      matches.push(i);
      i += m;
    } else {
      const shift = j - (badChar[text[i + j]] ?? -1);
      i += Math.max(1, shift);
    }
  }
  return matches;
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
      throw new Error(`Activity ${item.index} has start after end`);
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
    throw new Error('weights and values must have the same length');
  }
  if (capacity < 0) {
    throw new Error('capacity must be non-negative');
  }

  const items = weights
    .map((w, i) => ({
      ratio: w > 0 ? values[i] / w : 0,
      weight: w,
      value: values[i],
    }))
    .filter((item) => item.weight > 0)
    .sort((a, b) => b.ratio - a.ratio);

  let total = 0;
  let remaining = capacity;
  for (const { ratio, weight } of items) {
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

  const pq = new PriorityQueue<HuffmanNode>();
  for (const [char, freq] of entries) {
    pq.push({ freq, char }, freq);
  }

  if (pq.isEmpty()) return {};
  const first = pq.pop();
  if (!first) return {};
  if (pq.isEmpty()) {
    return { [first.item.char ?? '']: '0' };
  }
  pq.push(first.item, first.priority);

  while (pq.size() > 1) {
    const left = pq.pop()!;
    const right = pq.pop()!;
    const parent: HuffmanNode = {
      freq: left.priority + right.priority,
      left: left.item,
      right: right.item,
    };
    pq.push(parent, parent.freq);
  }

  const root = pq.pop()!.item;
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

export function countInversions<T extends Comparable>(arr: T[]): number {
  const result = sortAndCount([...arr]);
  return result.count;
}

function sortAndCount<T extends Comparable>(arr: T[]): { sorted: T[]; count: number } {
  const n = arr.length;
  if (n <= 1) return { sorted: arr, count: 0 };

  const mid = Math.floor(n / 2);
  const left = sortAndCount(arr.slice(0, mid));
  const right = sortAndCount(arr.slice(mid));
  const merged = mergeAndCount(left.sorted, right.sorted);
  return {
    sorted: merged.sorted,
    count: left.count + right.count + merged.count,
  };
}

function mergeAndCount<T extends Comparable>(
  left: T[],
  right: T[]
): { sorted: T[]; count: number } {
  const merged: T[] = [];
  let count = 0;
  let i = 0;
  let j = 0;
  while (i < left.length && j < right.length) {
    if (left[i] <= right[j]) {
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
  if (exponent < 0) {
    return 1 / fastPower(base, -exponent);
  }
  if (exponent === 0) return 1;
  if (exponent === 1) return base;

  const half = fastPower(base, Math.floor(exponent / 2));
  if (exponent % 2 === 0) {
    return half * half;
  }
  return base * half * half;
}
