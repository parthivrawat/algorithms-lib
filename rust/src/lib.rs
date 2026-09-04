//! Algorithms Library for Rust
//!
//! A comprehensive, zero-dependency collection of common algorithms across
//! sorting, searching, graphs, dynamic programming, strings, greedy, and
//! divide-and-conquer.
//!
//! # Conventions
//!
//! - Fallible operations return [`Result<T, Error>`](Error); lookups that miss
//!   return `Err(Error::NotFound)` rather than a sentinel index.
//! - `*_search` functions on strings report **byte offsets**, not char indices.
//! - Most sorts take `&[T]` and return `Vec<T>`; `*_in_place` variants sort
//!   `&mut [T]` without returning an allocation.
//!
//! # Example
//!
//! ```
//! use algorithms_lib::{quick_sort, binary_search};
//!
//! let sorted = quick_sort(&[3, 1, 4, 1, 5]);
//! assert_eq!(sorted, vec![1, 1, 3, 4, 5]);
//! assert_eq!(binary_search(&sorted, &4), Ok(3));
//! ```

#![warn(missing_docs)]
#![allow(
    clippy::for_kv_map,
    clippy::multiple_bound_locations,
    clippy::needless_range_loop,
    clippy::type_complexity,
    clippy::unnecessary_lazy_evaluations,
    clippy::useless_vec
)]

use std::cmp::Ordering;
use std::cmp::Reverse;
use std::collections::{BinaryHeap, HashMap, HashSet, VecDeque};
use std::error::Error as StdError;
use std::fmt;
use std::hash::Hash;

/// Errors that can be returned by algorithms in this crate.
#[derive(Debug, Clone, PartialEq)]
pub enum Error {
    /// The requested value was not present (search miss, unreachable goal).
    NotFound,
    /// An argument violates the function's input contract.
    InvalidInput,
    /// The graph is malformed or violates the algorithm's weight requirements.
    InvalidGraph,
    /// A negative-weight cycle reachable from the start was detected.
    NegativeCycle,
}

impl fmt::Display for Error {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        match self {
            Error::NotFound => write!(f, "not found"),
            Error::InvalidInput => write!(f, "invalid input"),
            Error::InvalidGraph => write!(f, "invalid graph"),
            Error::NegativeCycle => write!(f, "negative-weight cycle detected"),
        }
    }
}

impl StdError for Error {}

// ---------------- Sorting ----------------

/// Quick sort. Returns a new sorted `Vec`; see [`quick_sort_in_place`] to sort
/// a mutable slice without allocating the result.
pub fn quick_sort<T: Ord + Clone>(items: &[T]) -> Vec<T> {
    if items.len() <= 1 {
        return items.to_vec();
    }
    let pivot = items[items.len() / 2].clone();
    let mut left = Vec::new();
    let mut middle = Vec::new();
    let mut right = Vec::new();
    for x in items {
        if *x < pivot {
            left.push(x.clone());
        } else if *x > pivot {
            right.push(x.clone());
        } else {
            middle.push(x.clone());
        }
    }
    let mut result = quick_sort(&left);
    result.extend(middle);
    result.extend(quick_sort(&right));
    result
}

/// Stable merge sort.
pub fn merge_sort<T: Ord + Clone>(items: &[T]) -> Vec<T> {
    if items.len() <= 1 {
        return items.to_vec();
    }
    let mid = items.len() / 2;
    let left = merge_sort(&items[..mid]);
    let right = merge_sort(&items[mid..]);
    merge(&left, &right)
}

fn merge<T: Ord>(left: &[T], right: &[T]) -> Vec<T>
where
    T: Clone,
{
    let mut merged = Vec::with_capacity(left.len() + right.len());
    let mut i = 0;
    let mut j = 0;
    while i < left.len() && j < right.len() {
        if left[i] <= right[j] {
            merged.push(left[i].clone());
            i += 1;
        } else {
            merged.push(right[j].clone());
            j += 1;
        }
    }
    merged.extend_from_slice(&left[i..]);
    merged.extend_from_slice(&right[j..]);
    merged
}

/// Binary max-heap sort.
pub fn heap_sort<T: Ord + Clone>(items: &[T]) -> Vec<T> {
    let mut arr = items.to_vec();
    let n = arr.len();
    if n <= 1 {
        return arr;
    }
    for i in (0..n / 2).rev() {
        sift_down(&mut arr, n, i);
    }
    for end in (1..n).rev() {
        arr.swap(0, end);
        sift_down(&mut arr, end, 0);
    }
    arr
}

fn sift_down<T: Ord>(arr: &mut [T], n: usize, i: usize) {
    let mut largest = i;
    let left = 2 * i + 1;
    let right = 2 * i + 2;
    if left < n && arr[left] > arr[largest] {
        largest = left;
    }
    if right < n && arr[right] > arr[largest] {
        largest = right;
    }
    if largest != i {
        arr.swap(i, largest);
        sift_down(arr, n, largest);
    }
}

/// Sorts `items` in place using quick sort's partitioning result.
///
/// This variant avoids returning an allocated `Vec`; internally it sorts a
/// scratch buffer and copies back, so it does not require `T: Ord` callers to
/// receive a new collection.
pub fn quick_sort_in_place<T: Ord + Clone>(items: &mut [T]) {
    let sorted = quick_sort(items);
    items.clone_from_slice(&sorted);
}

/// Sorts `items` in place using stable merge sort.
///
/// Merge sort inherently needs O(n) auxiliary storage; this variant keeps that
/// scratch internal instead of returning a new `Vec`.
pub fn merge_sort_in_place<T: Ord + Clone>(items: &mut [T]) {
    let sorted = merge_sort(items);
    items.clone_from_slice(&sorted);
}

/// Sorts `items` in place using a binary max-heap. True in-place sort:
/// no auxiliary buffer proportional to `items` is allocated.
pub fn heap_sort_in_place<T: Ord>(items: &mut [T]) {
    let n = items.len();
    if n <= 1 {
        return;
    }
    for i in (0..n / 2).rev() {
        sift_down(items, n, i);
    }
    for end in (1..n).rev() {
        items.swap(0, end);
        sift_down(items, end, 0);
    }
}

/// LSD radix sort for non-negative integers.
pub fn radix_sort(items: &[i64]) -> Result<Vec<i64>, Error> {
    if items.is_empty() {
        return Ok(Vec::new());
    }
    let mut max_val = 0i64;
    for &v in items {
        if v < 0 {
            return Err(Error::InvalidInput);
        }
        if v > max_val {
            max_val = v;
        }
    }
    let mut arr = items.to_vec();
    let mut exp: i64 = 1;
    while max_val / exp > 0 {
        counting_sort_by_digit(&mut arr, exp);
        // Stop before exp * 10 can overflow i64. No further digit pass is
        // needed anyway: no i64 has a higher digit than exp's place.
        if exp > i64::MAX / 10 {
            break;
        }
        exp *= 10;
    }
    Ok(arr)
}

fn counting_sort_by_digit(arr: &mut [i64], exp: i64) {
    let n = arr.len();
    let mut output = vec![0; n];
    let mut count = vec![0; 10];
    for &num in arr.iter() {
        count[((num / exp) % 10) as usize] += 1;
    }
    for i in 1..10 {
        count[i] += count[i - 1];
    }
    for i in (0..n).rev() {
        let index = ((arr[i] / exp) % 10) as usize;
        output[count[index] - 1] = arr[i];
        count[index] -= 1;
    }
    arr.copy_from_slice(&output);
}

/// Standard-library sort wrapper.
///
/// This delegates to the language's built-in `slice::sort`.
pub fn native_sort<T: Ord + Clone>(items: &[T]) -> Vec<T> {
    let mut arr = items.to_vec();
    arr.sort();
    arr
}

/// Sorts `items` in place using the standard library sort (`slice::sort`).
pub fn native_sort_in_place<T: Ord>(items: &mut [T]) {
    items.sort();
}

// ---------------- Searching ----------------

fn is_sorted<T: Ord>(arr: &[T]) -> bool {
    arr.windows(2).all(|w| w[0] <= w[1])
}

/// Binary search.
///
/// `arr` must be sorted in ascending order; this is a caller contract checked
/// with `debug_assert!` in debug builds so searches stay O(log n). On unsorted
/// input in release builds the result is unspecified (typically `NotFound`).
///
/// Returns `Ok(index)` of the target, or `Err(Error::NotFound)` if absent.
///
/// # Examples
/// ```
/// use algorithms_lib::binary_search;
/// assert_eq!(binary_search(&[2, 4, 6, 8], &6), Ok(2));
/// ```
pub fn binary_search<T: Ord>(arr: &[T], target: &T) -> Result<usize, Error> {
    debug_assert!(is_sorted(arr), "binary_search requires a sorted slice");
    let mut lo = 0;
    let mut hi = arr.len() as isize - 1;
    while lo <= hi {
        let mid = (lo + hi) / 2;
        let mid_usize = mid as usize;
        if &arr[mid_usize] == target {
            return Ok(mid_usize);
        }
        if &arr[mid_usize] < target {
            lo = mid + 1;
        } else {
            hi = mid - 1;
        }
    }
    Err(Error::NotFound)
}

/// Binary search on answer for monotone predicates.
///
/// `low` and `high` are inclusive. `find` is `"minimum"` for the smallest
/// satisfying value or `"maximum"` for the largest.
pub fn binary_search_on_answer(
    low: i64,
    high: i64,
    predicate: impl Fn(i64) -> bool,
    find: &str,
) -> Result<i64, Error> {
    if find != "minimum" && find != "maximum" {
        return Err(Error::InvalidInput);
    }
    if low > high {
        return Err(Error::NotFound);
    }
    let mut left = low;
    let mut right = high;
    let mut answer: Option<i64> = None;
    if find == "minimum" {
        while left <= right {
            let mid = left + (right - left) / 2;
            if predicate(mid) {
                answer = Some(mid);
                right = mid - 1;
            } else {
                left = mid + 1;
            }
        }
    } else {
        while left <= right {
            let mid = left + (right - left) / 2;
            if predicate(mid) {
                answer = Some(mid);
                left = mid + 1;
            } else {
                right = mid - 1;
            }
        }
    }
    answer.ok_or(Error::NotFound)
}

/// Interpolation search for numeric data.
///
/// `arr` must be sorted in ascending order (caller contract, checked with
/// `debug_assert!` in debug builds). Returns `Ok(index)` of the target, or
/// `Err(Error::NotFound)` if absent.
pub fn interpolation_search(arr: &[i64], target: i64) -> Result<usize, Error> {
    debug_assert!(
        is_sorted(arr),
        "interpolation_search requires a sorted slice"
    );
    let mut lo = 0;
    let mut hi = arr.len() as isize - 1;
    while lo <= hi && target >= arr[lo as usize] && target <= arr[hi as usize] {
        if arr[hi as usize] == arr[lo as usize] {
            if arr[lo as usize] == target {
                return Ok(lo as usize);
            }
            break;
        }
        let lo_usize = lo as usize;
        let hi_usize = hi as usize;
        // Use i128 for the intermediate products so extreme i64 values cannot
        // overflow before the division.
        let numerator = (target as i128 - arr[lo_usize] as i128) * ((hi_usize - lo_usize) as i128);
        let denominator = arr[hi_usize] as i128 - arr[lo_usize] as i128;
        let pos = lo_usize + (numerator / denominator) as usize;
        if pos < lo_usize || pos > hi_usize {
            break;
        }
        if arr[pos] == target {
            return Ok(pos);
        }
        if arr[pos] < target {
            lo = pos as isize + 1;
        } else {
            hi = pos as isize - 1;
        }
    }
    Err(Error::NotFound)
}

/// Jump search.
///
/// `arr` must be sorted in ascending order (caller contract, checked with
/// `debug_assert!` in debug builds). Returns `Ok(index)` of the target, or
/// `Err(Error::NotFound)` if absent.
pub fn jump_search<T: Ord>(arr: &[T], target: &T) -> Result<usize, Error> {
    debug_assert!(is_sorted(arr), "jump_search requires a sorted slice");
    let n = arr.len();
    if n == 0 {
        return Err(Error::NotFound);
    }
    let mut step = (n as f64).sqrt() as usize;
    if step == 0 {
        step = 1;
    }
    let mut prev = 0;
    while arr[step.min(n) - 1] < *target {
        prev = step;
        step += (n as f64).sqrt() as usize;
        if prev >= n {
            return Err(Error::NotFound);
        }
    }
    while arr[prev] < *target {
        prev += 1;
        if prev == step.min(n) {
            return Err(Error::NotFound);
        }
    }
    if &arr[prev] == target {
        return Ok(prev);
    }
    Err(Error::NotFound)
}

// ---------------- Graphs ----------------

/// Weighted edge for graph algorithms.
#[derive(Debug, Clone, PartialEq)]
pub struct Edge {
    /// Destination vertex.
    pub to: String,
    /// Edge weight. `dijkstra`/`a_star` require `weight >= 0.0` and finite;
    /// `bellman_ford` accepts negative weights.
    pub weight: f64,
}

impl Edge {
    /// Creates an edge to vertex `to` with the given `weight`.
    pub fn new(to: &str, weight: f64) -> Self {
        Edge {
            to: to.to_string(),
            weight,
        }
    }
}

/// Returns `Err(Error::InvalidGraph)` if any neighbor referenced in the
/// adjacency map is not itself a key of the graph.
fn validate_adjacency_list(graph: &HashMap<String, Vec<String>>) -> Result<(), Error> {
    for neighbors in graph.values() {
        for nbr in neighbors {
            if !graph.contains_key(nbr) {
                return Err(Error::InvalidGraph);
            }
        }
    }
    Ok(())
}

/// Breadth-first search.
///
/// Every neighbor listed in the adjacency map must also be a key of the map;
/// missing keys return `Err(Error::InvalidGraph)`.
pub fn bfs(graph: &HashMap<String, Vec<String>>, start: &str) -> Result<Vec<String>, Error> {
    if !graph.contains_key(start) {
        return Err(Error::InvalidGraph);
    }
    validate_adjacency_list(graph)?;
    let mut visited = HashSet::new();
    visited.insert(start.to_string());
    let mut queue = VecDeque::new();
    queue.push_back(start.to_string());
    let mut result = Vec::new();
    while let Some(current) = queue.pop_front() {
        result.push(current.clone());
        if let Some(neighbors) = graph.get(&current) {
            for nbr in neighbors {
                if visited.insert(nbr.clone()) {
                    queue.push_back(nbr.clone());
                }
            }
        }
    }
    Ok(result)
}

/// Depth-first search.
///
/// Every neighbor listed in the adjacency map must also be a key of the map;
/// missing keys return `Err(Error::InvalidGraph)`.
pub fn dfs(graph: &HashMap<String, Vec<String>>, start: &str) -> Result<Vec<String>, Error> {
    if !graph.contains_key(start) {
        return Err(Error::InvalidGraph);
    }
    validate_adjacency_list(graph)?;
    let mut visited = HashSet::new();
    visited.insert(start.to_string());
    let mut stack = vec![start.to_string()];
    let mut result = Vec::new();
    while let Some(current) = stack.pop() {
        result.push(current.clone());
        if let Some(neighbors) = graph.get(&current) {
            for nbr in neighbors.iter().rev() {
                if visited.insert(nbr.clone()) {
                    stack.push(nbr.clone());
                }
            }
        }
    }
    Ok(result)
}

#[derive(Clone)]
struct State {
    cost: f64,
    node: String,
}

impl PartialEq for State {
    fn eq(&self, other: &Self) -> bool {
        self.cost.total_cmp(&other.cost) == Ordering::Equal && self.node == other.node
    }
}

impl Eq for State {}

impl PartialOrd for State {
    fn partial_cmp(&self, other: &Self) -> Option<Ordering> {
        Some(self.cmp(other))
    }
}

impl Ord for State {
    fn cmp(&self, other: &Self) -> Ordering {
        other
            .cost
            .total_cmp(&self.cost)
            .then_with(|| self.node.cmp(&other.node))
    }
}

/// Dijkstra's shortest paths.
pub fn dijkstra(
    graph: &HashMap<String, Vec<Edge>>,
    start: &str,
) -> Result<(HashMap<String, f64>, HashMap<String, String>), Error> {
    if !graph.contains_key(start) {
        return Err(Error::InvalidGraph);
    }

    // Pre-scan all edges so negative weights in unreachable components are
    // detected, and collect every vertex referenced in the graph.
    let mut all_nodes: HashSet<String> = graph.keys().cloned().collect();
    let mut negative = false;
    for edges in graph.values() {
        for e in edges {
            if e.weight < 0.0 || !e.weight.is_finite() {
                negative = true;
            }
            all_nodes.insert(e.to.clone());
        }
    }
    if negative {
        return Err(Error::InvalidGraph);
    }

    let mut distances: HashMap<String, f64> = all_nodes
        .iter()
        .map(|k| (k.clone(), f64::INFINITY))
        .collect();
    distances.insert(start.to_string(), 0.0);
    let mut predecessors = HashMap::new();

    let mut pq = BinaryHeap::new();
    pq.push(State {
        cost: 0.0,
        node: start.to_string(),
    });

    while let Some(State { cost, node }) = pq.pop() {
        if cost > distances[&node] {
            continue;
        }
        if let Some(edges) = graph.get(&node) {
            for e in edges {
                let nd = cost + e.weight;
                let current = distances.get(&e.to).copied().unwrap_or(f64::INFINITY);
                if nd < current {
                    distances.insert(e.to.clone(), nd);
                    predecessors.insert(e.to.clone(), node.clone());
                    pq.push(State {
                        cost: nd,
                        node: e.to.clone(),
                    });
                }
            }
        }
    }
    Ok((distances, predecessors))
}

/// A* shortest path.
///
/// The heuristic must be non-negative. For an optimal first result it should be
/// admissible and consistent; inconsistent but admissible heuristics still
/// return an optimal path, but may cause nodes to be re-expanded.
///
/// All edge weights must be non-negative.
pub fn a_star<F>(
    graph: &HashMap<String, Vec<Edge>>,
    start: &str,
    goal: &str,
    heuristic: F,
) -> Result<(Vec<String>, f64), Error>
where
    F: Fn(&str, &str) -> f64,
{
    if !graph.contains_key(start) || !graph.contains_key(goal) {
        return Err(Error::InvalidGraph);
    }

    let h = heuristic(start, goal);
    if h < 0.0 {
        return Err(Error::InvalidInput);
    }

    let mut g_score = HashMap::new();
    g_score.insert(start.to_string(), 0.0);
    let mut best_f: HashMap<String, f64> = HashMap::new();
    best_f.insert(start.to_string(), h);
    let mut came_from = HashMap::new();
    let mut pq = BinaryHeap::new();
    pq.push(State {
        cost: h,
        node: start.to_string(),
    });

    while let Some(State { cost, node }) = pq.pop() {
        if let Some(&best) = best_f.get(&node) {
            if cost > best {
                continue;
            }
        }
        if node == goal {
            return Ok((reconstruct_path(&came_from, &node), g_score[goal]));
        }
        if let Some(edges) = graph.get(&node) {
            for e in edges {
                if e.weight < 0.0 || !e.weight.is_finite() {
                    return Err(Error::InvalidGraph);
                }
                let h_nbr = heuristic(&e.to, goal);
                if h_nbr < 0.0 {
                    return Err(Error::InvalidInput);
                }
                let tentative = g_score[&node] + e.weight;
                let current = g_score.get(&e.to).copied().unwrap_or(f64::INFINITY);
                if tentative < current {
                    came_from.insert(e.to.clone(), node.clone());
                    g_score.insert(e.to.clone(), tentative);
                    let f = tentative + h_nbr;
                    let best = best_f.get(&e.to).copied().unwrap_or(f64::INFINITY);
                    if f < best {
                        best_f.insert(e.to.clone(), f);
                        pq.push(State {
                            cost: f,
                            node: e.to.clone(),
                        });
                    }
                }
            }
        }
    }
    Err(Error::NotFound)
}

fn reconstruct_path(came_from: &HashMap<String, String>, goal: &str) -> Vec<String> {
    let mut current = goal.to_string();
    let mut path = vec![current.clone()];
    while let Some(prev) = came_from.get(&current) {
        path.push(prev.clone());
        current = prev.clone();
    }
    path.reverse();
    path
}

struct Edge3 {
    from: String,
    to: String,
    weight: f64,
}

/// Bellman-Ford shortest paths with negative cycle detection.
pub fn bellman_ford(
    graph: &HashMap<String, Vec<Edge>>,
    start: &str,
) -> Result<(HashMap<String, f64>, HashMap<String, String>), Error> {
    if !graph.contains_key(start) {
        return Err(Error::InvalidGraph);
    }
    let mut vertices = HashSet::new();
    let mut edges = Vec::new();
    for (u, list) in graph {
        vertices.insert(u.clone());
        for e in list {
            vertices.insert(e.to.clone());
            edges.push(Edge3 {
                from: u.clone(),
                to: e.to.clone(),
                weight: e.weight,
            });
        }
    }
    let mut distances = HashMap::new();
    distances.insert(start.to_string(), 0.0);
    let mut predecessors = HashMap::new();

    for _ in 0..vertices.len().saturating_sub(1) {
        let mut updated = false;
        for e in &edges {
            if let Some(&du) = distances.get(&e.from) {
                let dv = distances.get(&e.to).copied().unwrap_or(f64::INFINITY);
                if du + e.weight < dv {
                    distances.insert(e.to.clone(), du + e.weight);
                    predecessors.insert(e.to.clone(), e.from.clone());
                    updated = true;
                }
            }
        }
        if !updated {
            break;
        }
    }

    for e in &edges {
        if let Some(&du) = distances.get(&e.from) {
            let dv = distances.get(&e.to).copied().unwrap_or(f64::INFINITY);
            if du + e.weight < dv {
                return Err(Error::NegativeCycle);
            }
        }
    }
    Ok((distances, predecessors))
}

// ---------------- Dynamic programming ----------------

/// 0/1 knapsack.
///
/// `capacity` and every weight must be non-negative and fit in `usize`
/// (on 32-bit platforms values above `u32::MAX` are rejected instead of
/// silently truncated).
pub fn knapsack_01(weights: &[i64], values: &[i64], capacity: i64) -> Result<i64, Error> {
    if weights.len() != values.len() {
        return Err(Error::InvalidInput);
    }
    let cap = usize::try_from(capacity).map_err(|_| Error::InvalidInput)?;
    let mut dp = vec![0; cap + 1];
    for i in 0..weights.len() {
        let w = usize::try_from(weights[i]).map_err(|_| Error::InvalidInput)?;
        let v = values[i];
        for c in (w..=cap).rev() {
            if dp[c - w] + v > dp[c] {
                dp[c] = dp[c - w] + v;
            }
        }
    }
    Ok(dp[cap])
}

/// Longest common subsequence length.
pub fn longest_common_subsequence<T: PartialEq>(a: &[T], b: &[T]) -> usize {
    let (m, n) = (a.len(), b.len());
    let mut prev = vec![0; n + 1];
    for i in 1..=m {
        let mut curr = vec![0; n + 1];
        for j in 1..=n {
            if a[i - 1] == b[j - 1] {
                curr[j] = prev[j - 1] + 1;
            } else {
                curr[j] = prev[j].max(curr[j - 1]);
            }
        }
        prev = curr;
    }
    prev[n]
}

/// Levenshtein distance.
pub fn edit_distance<T: PartialEq>(a: &[T], b: &[T]) -> usize {
    let (m, n) = (a.len(), b.len());
    let mut prev: Vec<usize> = (0..=n).collect();
    for i in 1..=m {
        let mut curr = vec![0; n + 1];
        curr[0] = i;
        for j in 1..=n {
            let cost = if a[i - 1] == b[j - 1] { 0 } else { 1 };
            let insertion = curr[j - 1] + 1;
            let deletion = prev[j] + 1;
            let substitution = prev[j - 1] + cost;
            curr[j] = insertion.min(deletion).min(substitution);
        }
        prev = curr;
    }
    prev[n]
}

/// One step in an edit-distance script.
#[derive(Debug, PartialEq, Clone)]
pub enum EditOp<T: Clone + PartialEq> {
    /// Keep the same symbol.
    Match(T),
    /// Delete a symbol from the source.
    Delete(T),
    /// Insert a symbol from the target.
    Insert(T),
    /// Replace one symbol with another.
    Substitute {
        /// Original symbol.
        from: T,
        /// Replacement symbol.
        to: T,
    },
}

/// Levenshtein distance and one optimal edit script.
pub fn edit_distance_reconstruction<T: PartialEq + Clone>(
    a: &[T],
    b: &[T],
) -> (usize, Vec<EditOp<T>>) {
    let (m, n) = (a.len(), b.len());
    let mut dp = vec![vec![0; n + 1]; m + 1];
    for i in 0..=m {
        dp[i][0] = i;
    }
    for j in 0..=n {
        dp[0][j] = j;
    }
    for i in 1..=m {
        for j in 1..=n {
            if a[i - 1] == b[j - 1] {
                dp[i][j] = dp[i - 1][j - 1];
            } else {
                dp[i][j] = 1 + dp[i - 1][j].min(dp[i][j - 1]).min(dp[i - 1][j - 1]);
            }
        }
    }
    let mut script = Vec::new();
    let (mut i, mut j) = (m, n);
    while i > 0 || j > 0 {
        if i == 0 {
            script.push(EditOp::Insert(b[j - 1].clone()));
            j -= 1;
        } else if j == 0 {
            script.push(EditOp::Delete(a[i - 1].clone()));
            i -= 1;
        } else if a[i - 1] == b[j - 1] {
            script.push(EditOp::Match(a[i - 1].clone()));
            i -= 1;
            j -= 1;
        } else {
            let best = dp[i][j];
            if dp[i - 1][j - 1] + 1 == best {
                script.push(EditOp::Substitute {
                    from: a[i - 1].clone(),
                    to: b[j - 1].clone(),
                });
                i -= 1;
                j -= 1;
            } else if dp[i][j - 1] + 1 == best {
                script.push(EditOp::Insert(b[j - 1].clone()));
                j -= 1;
            } else {
                script.push(EditOp::Delete(a[i - 1].clone()));
                i -= 1;
            }
        }
    }
    script.reverse();
    (dp[m][n], script)
}

/// Longest common subsequence length and one actual subsequence.
pub fn longest_common_subsequence_reconstruction<T: PartialEq + Clone>(
    a: &[T],
    b: &[T],
) -> (usize, Vec<T>) {
    let (m, n) = (a.len(), b.len());
    let mut dp = vec![vec![0; n + 1]; m + 1];
    for i in 1..=m {
        for j in 1..=n {
            if a[i - 1] == b[j - 1] {
                dp[i][j] = dp[i - 1][j - 1] + 1;
            } else {
                dp[i][j] = dp[i - 1][j].max(dp[i][j - 1]);
            }
        }
    }
    let mut result = Vec::with_capacity(dp[m][n]);
    let (mut i, mut j) = (m, n);
    while i > 0 && j > 0 {
        if a[i - 1] == b[j - 1] {
            result.push(a[i - 1].clone());
            i -= 1;
            j -= 1;
        } else if dp[i - 1][j] >= dp[i][j - 1] {
            i -= 1;
        } else {
            j -= 1;
        }
    }
    result.reverse();
    (dp[m][n], result)
}

// ---------------- String algorithms ----------------

/// Knuth-Morris-Pratt pattern matching.
///
/// Returns all starting **byte offsets** of `pattern` in `text`, including
/// overlapping occurrences. Returns `Err(Error::InvalidInput)` when `pattern`
/// is empty.
pub fn kmp_search(text: &str, pattern: &str) -> Result<Vec<usize>, Error> {
    if pattern.is_empty() {
        return Err(Error::InvalidInput);
    }
    let t = text.as_bytes();
    let p = pattern.as_bytes();
    let failure = compute_kmp_failure(p);
    let mut matches = Vec::new();
    let mut j = 0;
    for i in 0..t.len() {
        while j > 0 && t[i] != p[j] {
            j = failure[j - 1];
        }
        if t[i] == p[j] {
            j += 1;
        }
        if j == p.len() {
            matches.push(i + 1 - p.len());
            j = failure[j - 1];
        }
    }
    Ok(matches)
}

fn compute_kmp_failure(pattern: &[u8]) -> Vec<usize> {
    let mut failure = vec![0; pattern.len()];
    let mut j = 0;
    for i in 1..pattern.len() {
        while j > 0 && pattern[i] != pattern[j] {
            j = failure[j - 1];
        }
        if pattern[i] == pattern[j] {
            j += 1;
            failure[i] = j;
        }
    }
    failure
}

/// Rolling-hash pattern matching.
///
/// Returns all starting **byte offsets** of `pattern` in `text`, including
/// overlapping occurrences.
///
/// `base` and `modulus` must be positive; `pattern` must not be empty —
/// violations return `Err(Error::InvalidInput)`. The arithmetic is performed
/// in `i128` so large `base`/`modulus` values do not overflow.
pub fn rabin_karp_search(
    text: &str,
    pattern: &str,
    base: i64,
    modulus: i64,
) -> Result<Vec<usize>, Error> {
    if pattern.is_empty() {
        return Err(Error::InvalidInput);
    }
    if modulus <= 0 {
        return Err(Error::InvalidInput);
    }
    if base <= 0 {
        return Err(Error::InvalidInput);
    }

    let t = text.as_bytes();
    let p = pattern.as_bytes();
    let (n, m) = (t.len(), p.len());
    if m > n {
        return Ok(Vec::new());
    }

    let base = base as i128;
    let modulus = modulus as i128;

    let mut h = 1i128;
    for _ in 0..m - 1 {
        h = (h * base) % modulus;
    }
    let mut pattern_hash = 0i128;
    let mut text_hash = 0i128;
    for i in 0..m {
        pattern_hash = (pattern_hash * base + p[i] as i128) % modulus;
        text_hash = (text_hash * base + t[i] as i128) % modulus;
    }
    let mut matches = Vec::new();
    for i in 0..=n - m {
        if text_hash == pattern_hash && &t[i..i + m] == p {
            matches.push(i);
        }
        if i < n - m {
            text_hash = (text_hash - (t[i] as i128) * h).rem_euclid(modulus);
            text_hash = (text_hash * base + t[i + m] as i128).rem_euclid(modulus);
        }
    }
    Ok(matches)
}

/// Boyer-Moore pattern matching with the bad-character and good-suffix rules.
///
/// Returns all starting **byte offsets** of `pattern` in `text`, including
/// overlapping occurrences. Returns `Err(Error::InvalidInput)` when `pattern`
/// is empty.
pub fn boyer_moore_search(text: &str, pattern: &str) -> Result<Vec<usize>, Error> {
    if pattern.is_empty() {
        return Err(Error::InvalidInput);
    }
    let t = text.as_bytes();
    let p = pattern.as_bytes();
    let (n, m) = (t.len(), p.len());
    if m > n {
        return Ok(Vec::new());
    }
    let mut bad_char = HashMap::new();
    for (i, &b) in p.iter().enumerate() {
        bad_char.insert(b, i);
    }
    let good_suffix = good_suffix_shifts(p);
    let mut matches = Vec::new();
    let mut i = 0usize;
    while i <= n - m {
        let mut j = m as isize - 1;
        while j >= 0 && t[i + j as usize] == p[j as usize] {
            j -= 1;
        }
        if j < 0 {
            matches.push(i);
            i += good_suffix[0];
        } else {
            let last = bad_char
                .get(&t[i + j as usize])
                .map(|&l| l as isize)
                .unwrap_or(-1);
            i += (good_suffix[j as usize + 1] as isize).max(j - last).max(1) as usize;
        }
    }
    Ok(matches)
}

/// Builds the Boyer-Moore good-suffix table: `shift[k]` is the shift applied
/// after a mismatch at pattern index `k - 1`, and `shift[0]` is the
/// border-based shift applied after a full match so that overlapping matches
/// are still reported.
fn good_suffix_shifts(pattern: &[u8]) -> Vec<usize> {
    let m = pattern.len();
    let mut shift = vec![0usize; m + 1];
    let mut border_pos = vec![0usize; m + 1];
    let (mut i, mut j) = (m, m + 1);
    border_pos[i] = j;
    while i > 0 {
        while j <= m && pattern[i - 1] != pattern[j - 1] {
            if shift[j] == 0 {
                shift[j] = j - i;
            }
            j = border_pos[j];
        }
        i -= 1;
        j -= 1;
        border_pos[i] = j;
    }
    let mut j = border_pos[0];
    for i in 0..=m {
        if shift[i] == 0 {
            shift[i] = j;
        }
        if i == j {
            j = border_pos[j];
        }
    }
    shift
}

// ---------------- Greedy ----------------

/// Maximum compatible activity indices.
pub fn activity_selection(activities: &[(i64, i64)]) -> Result<Vec<usize>, Error> {
    if activities.is_empty() {
        return Ok(Vec::new());
    }
    #[derive(Clone, Copy)]
    struct Indexed {
        index: usize,
        start: i64,
        end: i64,
    }
    let mut list = Vec::with_capacity(activities.len());
    for (i, &(start, end)) in activities.iter().enumerate() {
        if start > end {
            return Err(Error::InvalidInput);
        }
        list.push(Indexed {
            index: i,
            start,
            end,
        });
    }
    list.sort_by(|a, b| a.end.cmp(&b.end));
    let mut selected = vec![list[0].index];
    let mut last_end = list[0].end;
    for item in list.iter().skip(1) {
        if item.start >= last_end {
            selected.push(item.index);
            last_end = item.end;
        }
    }
    Ok(selected)
}

/// Fractional knapsack maximum value.
///
/// Negative or NaN weights and NaN values are rejected with
/// `Err(Error::InvalidInput)`; a zero-weight item contributes its full value
/// when that value is positive, since it consumes no capacity. Items with a
/// non-positive value-to-weight ratio are never taken.
pub fn fractional_knapsack(weights: &[f64], values: &[f64], capacity: f64) -> Result<f64, Error> {
    if weights.len() != values.len() {
        return Err(Error::InvalidInput);
    }
    if capacity < 0.0 {
        return Err(Error::InvalidInput);
    }
    #[derive(Clone, Copy)]
    struct Item {
        ratio: f64,
        weight: f64,
    }
    let mut items = Vec::with_capacity(weights.len());
    let mut total = 0.0;
    for i in 0..weights.len() {
        if weights[i] < 0.0 || weights[i].is_nan() || values[i].is_nan() {
            return Err(Error::InvalidInput);
        }
        if weights[i] == 0.0 {
            // Zero-weight items consume no capacity: take them if they add value.
            total += values[i].max(0.0);
            continue;
        }
        items.push(Item {
            ratio: values[i] / weights[i],
            weight: weights[i],
        });
    }
    items.sort_by(|a, b| b.ratio.partial_cmp(&a.ratio).unwrap_or(Ordering::Equal));
    let mut remaining = capacity;
    for it in items {
        if it.ratio <= 0.0 {
            break; // taking a non-positive-value item can only lower the total
        }
        let take = it.weight.min(remaining);
        total += take * it.ratio;
        remaining -= take;
        if remaining <= 0.0 {
            break;
        }
    }
    Ok(total)
}

struct HuffmanNode {
    freq: i64,
    char: String,
    left: Option<Box<HuffmanNode>>,
    right: Option<Box<HuffmanNode>>,
}

impl PartialEq for HuffmanNode {
    fn eq(&self, other: &Self) -> bool {
        self.freq == other.freq && self.char == other.char
    }
}

impl Eq for HuffmanNode {}

impl PartialOrd for HuffmanNode {
    fn partial_cmp(&self, other: &Self) -> Option<Ordering> {
        Some(self.cmp(other))
    }
}

impl Ord for HuffmanNode {
    fn cmp(&self, other: &Self) -> Ordering {
        self.freq
            .cmp(&other.freq)
            .then_with(|| self.char.cmp(&other.char))
    }
}

/// Huffman coding.
///
/// Every frequency must be positive; violations return
/// `Err(Error::InvalidInput)`. Symbol keys may be any string, including the
/// empty string.
pub fn huffman_coding(
    frequencies: &HashMap<String, i64>,
) -> Result<HashMap<String, String>, Error> {
    if frequencies.is_empty() {
        return Err(Error::InvalidInput);
    }
    if frequencies.values().any(|&f| f <= 0) {
        return Err(Error::InvalidInput);
    }
    let mut pq = BinaryHeap::new();
    for (ch, &freq) in frequencies {
        pq.push(Reverse(HuffmanNode {
            freq,
            char: ch.clone(),
            left: None,
            right: None,
        }));
    }
    if pq.len() == 1 {
        let node = pq.pop().unwrap().0;
        let mut codes = HashMap::new();
        codes.insert(node.char, "0".to_string());
        return Ok(codes);
    }
    while pq.len() > 1 {
        let left = pq.pop().unwrap().0;
        let right = pq.pop().unwrap().0;
        let parent = HuffmanNode {
            freq: left.freq + right.freq,
            char: String::new(),
            left: Some(Box::new(left)),
            right: Some(Box::new(right)),
        };
        pq.push(Reverse(parent));
    }
    let root = pq.pop().unwrap().0;
    let mut codes = HashMap::new();
    assign_huffman_codes(&root, "", &mut codes);
    Ok(codes)
}

fn assign_huffman_codes(node: &HuffmanNode, prefix: &str, codes: &mut HashMap<String, String>) {
    // Internal nodes always have both children; leaves have neither. This also
    // lets the empty string be used as a symbol key.
    if node.left.is_none() && node.right.is_none() {
        if prefix.is_empty() {
            codes.insert(node.char.clone(), "0".to_string());
        } else {
            codes.insert(node.char.clone(), prefix.to_string());
        }
        return;
    }
    if let Some(ref left) = node.left {
        assign_huffman_codes(left, &format!("{}0", prefix), codes);
    }
    if let Some(ref right) = node.right {
        assign_huffman_codes(right, &format!("{}1", prefix), codes);
    }
}

/// Encode a sequence of symbols using a Huffman code table.
pub fn huffman_encode(
    symbols: &[String],
    code_table: &HashMap<String, String>,
) -> Result<String, Error> {
    let mut parts = Vec::with_capacity(symbols.len());
    for s in symbols {
        let code = code_table.get(s).ok_or_else(|| Error::InvalidInput)?;
        parts.push(code.clone());
    }
    Ok(parts.concat())
}

/// Decode a Huffman bit string using a code table.
pub fn huffman_decode(
    encoded: &str,
    code_table: &HashMap<String, String>,
) -> Result<Vec<String>, Error> {
    let mut reverse = HashMap::with_capacity(code_table.len());
    for (s, code) in code_table {
        reverse.insert(code.clone(), s.clone());
    }
    let mut result = Vec::new();
    let mut current = String::new();
    for bit in encoded.chars() {
        current.push(bit);
        if let Some(s) = reverse.get(&current) {
            result.push(s.clone());
            current.clear();
        }
    }
    if !current.is_empty() {
        return Err(Error::InvalidInput);
    }
    Ok(result)
}

// ---------------- Divide and conquer ----------------

/// Maximum subarray sum.
pub fn max_subarray(arr: &[i64]) -> Result<i64, Error> {
    if arr.is_empty() {
        return Err(Error::InvalidInput);
    }
    Ok(max_subarray_dc(arr, 0, arr.len() - 1))
}

fn max_subarray_dc(arr: &[i64], left: usize, right: usize) -> i64 {
    if left == right {
        return arr[left];
    }
    let mid = (left + right) / 2;
    let left_sum = max_subarray_dc(arr, left, mid);
    let right_sum = max_subarray_dc(arr, mid + 1, right);
    let cross_sum = max_crossing_sum(arr, left, mid, right);
    left_sum.max(right_sum).max(cross_sum)
}

fn max_crossing_sum(arr: &[i64], left: usize, mid: usize, right: usize) -> i64 {
    let mut sum = arr[mid];
    let mut left_max = arr[mid];
    for i in (left..mid).rev() {
        sum += arr[i];
        if sum > left_max {
            left_max = sum;
        }
    }
    sum = arr[mid + 1];
    let mut right_max = arr[mid + 1];
    for i in mid + 2..=right {
        sum += arr[i];
        if sum > right_max {
            right_max = sum;
        }
    }
    left_max + right_max
}

/// Inversion count.
pub fn count_inversions<T: Ord + Clone>(arr: &[T]) -> usize {
    let (_, count) = sort_and_count(arr);
    count
}

fn sort_and_count<T: Ord + Clone>(arr: &[T]) -> (Vec<T>, usize) {
    let n = arr.len();
    if n <= 1 {
        return (arr.to_vec(), 0);
    }
    let mid = n / 2;
    let (left, lcount) = sort_and_count(&arr[..mid]);
    let (right, rcount) = sort_and_count(&arr[mid..]);
    let (merged, scount) = merge_and_count(&left, &right);
    (merged, lcount + rcount + scount)
}

fn merge_and_count<T: Ord + Clone>(left: &[T], right: &[T]) -> (Vec<T>, usize) {
    let mut merged = Vec::with_capacity(left.len() + right.len());
    let mut count = 0;
    let mut i = 0;
    let mut j = 0;
    while i < left.len() && j < right.len() {
        if left[i] <= right[j] {
            merged.push(left[i].clone());
            i += 1;
        } else {
            merged.push(right[j].clone());
            j += 1;
            count += left.len() - i;
        }
    }
    merged.extend_from_slice(&left[i..]);
    merged.extend_from_slice(&right[j..]);
    (merged, count)
}

/// Exponentiation by squaring. Returns `Err(Error::InvalidInput)` when `base`
/// is zero and `exponent` is negative. By convention, `0^0` returns `1.0`.
pub fn fast_power(base: f64, exponent: i64) -> Result<f64, Error> {
    if base == 0.0 && exponent < 0 {
        return Err(Error::InvalidInput);
    }
    // unsigned_abs yields |exponent| as u64 and is correct even for i64::MIN,
    // whose negation cannot be represented in i64.
    let mut exp = exponent.unsigned_abs();
    let mut result = 1.0f64;
    let mut b = base;
    while exp > 0 {
        if exp & 1 == 1 {
            result *= b;
        }
        b *= b;
        exp >>= 1;
    }
    if exponent < 0 {
        result = 1.0 / result;
    }
    Ok(result)
}

/// Greatest common divisor.
///
/// The result is always non-negative. Returns [`Error::InvalidInput`] if the
/// result does not fit in an `i64` (this only happens for `gcd(i64::MIN, 0)`).
pub fn gcd(a: i64, b: i64) -> Result<i64, Error> {
    let mut a = a.unsigned_abs();
    let mut b = b.unsigned_abs();
    while b != 0 {
        (a, b) = (b, a % b);
    }
    if a > i64::MAX as u64 {
        return Err(Error::InvalidInput);
    }
    Ok(a as i64)
}

/// Least common multiple.
///
/// Returns [`Error::InvalidInput`] if both inputs are zero or if the result
/// overflows an `i64`.
pub fn lcm(a: i64, b: i64) -> Result<i64, Error> {
    if a == 0 && b == 0 {
        return Err(Error::InvalidInput);
    }
    let g = gcd(a, b)?;
    let a_div = (a / g) as i128;
    let b_i = b as i128;
    let result = (a_div * b_i).abs();
    if result > i64::MAX as i128 {
        return Err(Error::InvalidInput);
    }
    Ok(result as i64)
}

/// Modular exponentiation by squaring: `(base^exponent) % mod`.
///
/// Returns [`Error::InvalidInput`] if `mod` is not positive or `exponent`
/// is negative.
pub fn modular_fast_power(base: i64, exponent: i64, modulus: i64) -> Result<i64, Error> {
    if modulus <= 0 {
        return Err(Error::InvalidInput);
    }
    if exponent < 0 {
        return Err(Error::InvalidInput);
    }
    if modulus == 1 {
        return Ok(0);
    }
    let m = modulus as i128;
    let mut b = (base % modulus) as i128;
    if b < 0 {
        b += m;
    }
    let mut result: i128 = 1;
    let mut exp = exponent as u64;
    while exp > 0 {
        if exp & 1 == 1 {
            result = (result * b) % m;
        }
        b = (b * b) % m;
        exp >>= 1;
    }
    Ok(result as i64)
}

/// Returns the k-th smallest element (0-indexed) of `items`.
///
/// Returns [`Error::InvalidInput`] if `k` is out of range.
pub fn quick_select<T: Ord + Clone>(items: &[T], k: usize) -> Result<T, Error> {
    if k >= items.len() {
        return Err(Error::InvalidInput);
    }
    let mut arr = items.to_vec();
    let mut k = k;
    loop {
        if arr.len() == 1 {
            return Ok(arr[0].clone());
        }
        let pivot = arr[arr.len() / 2].clone();
        let mut lows = Vec::new();
        let mut pivots = Vec::new();
        let mut highs = Vec::new();
        for x in &arr {
            match x.cmp(&pivot) {
                Ordering::Less => lows.push(x.clone()),
                Ordering::Greater => highs.push(x.clone()),
                Ordering::Equal => pivots.push(x.clone()),
            }
        }
        if k < lows.len() {
            arr = lows;
        } else if k < lows.len() + pivots.len() {
            return Ok(pivot);
        } else {
            k -= lows.len() + pivots.len();
            arr = highs;
        }
    }
}

/// Sieve of Eratosthenes. Returns all prime numbers `<= n`.
pub fn sieve_of_eratosthenes(n: i64) -> Result<Vec<i64>, Error> {
    if n < 0 {
        return Err(Error::InvalidInput);
    }
    let n = n as usize;
    if n < 2 {
        return Ok(Vec::new());
    }
    let mut sieve = vec![true; n + 1];
    sieve[0] = false;
    sieve[1] = false;
    let limit = (n as f64).sqrt() as usize;
    for p in 2..=limit {
        if sieve[p] {
            let mut multiple = p * p;
            while multiple <= n {
                sieve[multiple] = false;
                multiple += p;
            }
        }
    }
    Ok(sieve
        .iter()
        .enumerate()
        .filter(|(_, &is_prime)| is_prime)
        .map(|(i, _)| i as i64)
        .collect())
}

/// Topological sort using Kahn's algorithm.
pub fn topological_sort(graph: &HashMap<String, Vec<String>>) -> Result<Vec<String>, Error> {
    validate_adjacency_list(graph)?;
    let mut in_degree: HashMap<String, usize> = HashMap::new();
    for (node, _) in graph {
        in_degree.insert(node.clone(), 0);
    }
    for (_, neighbors) in graph {
        for nbr in neighbors {
            *in_degree.get_mut(nbr).unwrap() += 1;
        }
    }

    let mut queue: VecDeque<String> = in_degree
        .iter()
        .filter(|(_, &deg)| deg == 0)
        .map(|(node, _)| node.clone())
        .collect();

    let mut order = Vec::new();
    while let Some(node) = queue.pop_front() {
        order.push(node.clone());
        for nbr in graph.get(&node).unwrap_or(&Vec::new()) {
            let deg = in_degree.get_mut(nbr).unwrap();
            *deg -= 1;
            if *deg == 0 {
                queue.push_back(nbr.clone());
            }
        }
    }

    if order.len() != in_degree.len() {
        return Err(Error::InvalidGraph);
    }
    Ok(order)
}

/// Disjoint-set union–find with union by rank and path compression.
pub struct UnionFind<T: Hash + Eq + Clone> {
    parent: HashMap<T, T>,
    rank: HashMap<T, u32>,
}

impl<T: Hash + Eq + Clone> UnionFind<T> {
    /// Create a new union–find structure for the given elements.
    pub fn new(elements: impl IntoIterator<Item = T>) -> Self {
        let mut parent = HashMap::new();
        let mut rank = HashMap::new();
        for x in elements {
            parent.insert(x.clone(), x.clone());
            rank.insert(x, 0);
        }
        Self { parent, rank }
    }

    /// Return the representative of the set containing `x`.
    pub fn find(&mut self, x: &T) -> Result<T, Error> {
        if !self.parent.contains_key(x) {
            return Err(Error::InvalidInput);
        }
        let mut root = x.clone();
        while self.parent[&root] != root {
            root = self.parent[&root].clone();
        }
        let mut node = x.clone();
        while self.parent[&node] != root {
            let parent = self.parent[&node].clone();
            self.parent.insert(node.clone(), root.clone());
            node = parent;
        }
        Ok(root)
    }

    /// Merge the sets containing `x` and `y`.
    pub fn union(&mut self, x: &T, y: &T) -> Result<bool, Error> {
        let root_x = self.find(x)?;
        let root_y = self.find(y)?;
        if root_x == root_y {
            return Ok(false);
        }
        let rank_x = *self.rank.get(&root_x).unwrap_or(&0);
        let rank_y = *self.rank.get(&root_y).unwrap_or(&0);
        if rank_x < rank_y {
            self.parent.insert(root_x, root_y);
        } else if rank_x > rank_y {
            self.parent.insert(root_y, root_x);
        } else {
            self.parent.insert(root_y, root_x.clone());
            self.rank.insert(root_x, rank_x + 1);
        }
        Ok(true)
    }

    /// Return whether `x` and `y` are in the same set.
    pub fn connected(&mut self, x: &T, y: &T) -> Result<bool, Error> {
        Ok(self.find(x)? == self.find(y)?)
    }

    /// Return the current number of disjoint sets.
    pub fn count(&mut self) -> usize {
        let keys: Vec<T> = self.parent.keys().cloned().collect();
        let mut roots = HashSet::new();
        for k in keys {
            if let Ok(root) = self.find(&k) {
                roots.insert(root);
            }
        }
        roots.len()
    }
}

/// A min-priority queue.
pub struct PriorityQueue<T: Ord> {
    heap: BinaryHeap<Reverse<T>>,
}

impl<T: Ord> PriorityQueue<T> {
    /// Create a new priority queue with the given items.
    pub fn new(items: impl IntoIterator<Item = T>) -> Self {
        let heap: BinaryHeap<Reverse<T>> = items.into_iter().map(Reverse).collect();
        Self { heap }
    }

    /// Add an item.
    pub fn push(&mut self, item: T) {
        self.heap.push(Reverse(item));
    }

    /// Remove and return the smallest item.
    pub fn pop(&mut self) -> Result<T, Error> {
        self.heap.pop().map(|r| r.0).ok_or(Error::InvalidInput)
    }

    /// Return the smallest item without removing it.
    pub fn peek(&self) -> Result<&T, Error> {
        self.heap.peek().map(|r| &r.0).ok_or(Error::InvalidInput)
    }

    /// Return the number of items.
    pub fn len(&self) -> usize {
        self.heap.len()
    }

    /// Return whether the queue is empty.
    pub fn is_empty(&self) -> bool {
        self.heap.is_empty()
    }
}

/// An undirected, weighted edge for minimum spanning tree algorithms.
#[derive(Clone)]
pub struct MSTEdge {
    /// One endpoint of the edge.
    pub u: String,
    /// The other endpoint of the edge.
    pub v: String,
    /// Edge weight.
    pub weight: f64,
}

/// Kruskal's minimum spanning tree (or forest) algorithm.
pub fn minimum_spanning_tree(edges: &[MSTEdge]) -> Result<(f64, Vec<MSTEdge>), Error> {
    for e in edges {
        if !e.weight.is_finite() {
            return Err(Error::InvalidInput);
        }
    }

    let mut sorted: Vec<MSTEdge> = edges.to_vec();
    sorted.sort_by(|a, b| a.weight.partial_cmp(&b.weight).unwrap_or(Ordering::Equal));

    let mut vertices = HashSet::new();
    for e in edges {
        vertices.insert(e.u.clone());
        vertices.insert(e.v.clone());
    }
    let mut uf = UnionFind::new(Vec::new());
    for v in vertices {
        uf.parent.insert(v.clone(), v.clone());
        uf.rank.insert(v, 0);
    }

    let mut total = 0.0;
    let mut mst = Vec::new();
    for e in sorted {
        if e.u == e.v {
            continue;
        }
        if uf.union(&e.u, &e.v)? {
            total += e.weight;
            mst.push(e);
        }
    }
    Ok((total, mst))
}

/// Floyd–Warshall all-pairs shortest paths.
pub fn floyd_warshall(
    graph: &HashMap<String, Vec<Edge>>,
) -> Result<HashMap<String, HashMap<String, f64>>, Error> {
    let mut all_nodes = HashSet::new();
    for u in graph.keys() {
        all_nodes.insert(u.clone());
    }
    for edges in graph.values() {
        for e in edges {
            all_nodes.insert(e.to.clone());
        }
    }
    for (_u, edges) in graph {
        for e in edges {
            if !all_nodes.contains(&e.to) {
                return Err(Error::InvalidGraph);
            }
            if !e.weight.is_finite() {
                return Err(Error::InvalidInput);
            }
        }
    }

    let inf = f64::INFINITY;
    let mut dist: HashMap<String, HashMap<String, f64>> = HashMap::new();
    for u in &all_nodes {
        let mut row = HashMap::new();
        for v in &all_nodes {
            row.insert(v.clone(), if u == v { 0.0 } else { inf });
        }
        dist.insert(u.clone(), row);
    }
    for (u, edges) in graph {
        for e in edges {
            if e.weight < dist[u][&e.to] {
                dist.get_mut(u).unwrap().insert(e.to.clone(), e.weight);
            }
        }
    }

    for k in &all_nodes {
        for i in &all_nodes {
            let dik = dist[i][k];
            if dik.is_infinite() {
                continue;
            }
            for j in &all_nodes {
                let nd = dik + dist[k][j];
                if nd < dist[i][j] {
                    dist.get_mut(i).unwrap().insert(j.clone(), nd);
                }
            }
        }
    }

    for u in &all_nodes {
        if dist[u][u] < 0.0 {
            return Err(Error::NegativeCycle);
        }
    }
    Ok(dist)
}

/// Kosaraju's strongly connected components algorithm.
pub fn strongly_connected_components(
    graph: &HashMap<String, Vec<String>>,
) -> Result<Vec<Vec<String>>, Error> {
    validate_adjacency_list(graph)?;

    let mut visited = HashSet::new();
    let mut order: Vec<String> = Vec::new();

    for u in graph.keys() {
        if !visited.contains(u) {
            let mut stack: Vec<(String, bool)> = vec![(u.clone(), false)];
            while let Some((node, processed)) = stack.pop() {
                if processed {
                    order.push(node);
                    continue;
                }
                if visited.contains(&node) {
                    continue;
                }
                visited.insert(node.clone());
                stack.push((node.clone(), true));
                if let Some(neighbors) = graph.get(&node) {
                    for nbr in neighbors.iter().rev() {
                        if !visited.contains(nbr) {
                            stack.push((nbr.clone(), false));
                        }
                    }
                }
            }
        }
    }

    let mut reverse: HashMap<String, Vec<String>> = HashMap::new();
    for u in graph.keys() {
        reverse.insert(u.clone(), Vec::new());
    }
    for (u, neighbors) in graph {
        for v in neighbors {
            reverse.get_mut(v).unwrap().push(u.clone());
        }
    }

    let mut visited = HashSet::new();
    let mut components: Vec<Vec<String>> = Vec::new();

    for u in order.iter().rev() {
        if !visited.contains(u) {
            let mut component = Vec::new();
            let mut stack = vec![u.clone()];
            while let Some(node) = stack.pop() {
                if visited.contains(&node) {
                    continue;
                }
                visited.insert(node.clone());
                component.push(node.clone());
                if let Some(neighbors) = reverse.get(&node) {
                    for nbr in neighbors {
                        if !visited.contains(nbr) {
                            stack.push(nbr.clone());
                        }
                    }
                }
            }
            components.push(component);
        }
    }
    Ok(components)
}

/// Returns whether the directed graph contains a cycle.
pub fn has_cycle(graph: &HashMap<String, Vec<String>>) -> Result<bool, Error> {
    validate_adjacency_list(graph)?;

    let mut state: HashMap<String, i8> = HashMap::new();
    for u in graph.keys() {
        state.insert(u.clone(), 0);
    }

    for u in graph.keys() {
        if state[u] != 0 {
            continue;
        }
        let mut stack: Vec<(String, usize)> = vec![(u.clone(), 0)];
        while let Some((node, idx)) = stack.pop() {
            if *state.get(&node).unwrap() == 2 {
                continue;
            }
            if *state.get(&node).unwrap() == 0 {
                state.insert(node.clone(), 1);
            }
            let neighbors = graph.get(&node).cloned().unwrap_or_default();
            if idx < neighbors.len() {
                let nbr = &neighbors[idx];
                stack.push((node, idx + 1));
                match *state.get(nbr).unwrap_or(&0) {
                    1 => return Ok(true),
                    0 => stack.push((nbr.clone(), 0)),
                    _ => {}
                }
            } else {
                state.insert(node, 2);
            }
        }
    }
    Ok(false)
}
