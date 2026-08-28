//! Algorithms Library for Rust
//!
//! A comprehensive, zero-dependency collection of common algorithms.

use std::cmp::Ordering;
use std::cmp::Reverse;
use std::collections::{BinaryHeap, HashMap, HashSet, VecDeque};
use std::error::Error as StdError;
use std::fmt;

/// Errors that can be returned by algorithms in this crate.
#[derive(Debug, Clone, PartialEq)]
pub enum Error {
    NotFound,
    InvalidInput,
    InvalidGraph,
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

/// Quick sort.
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

/// LSD radix sort for non-negative integers.
pub fn radix_sort(items: &[i64]) -> Result<Vec<i64>, Error> {
    if items.is_empty() {
        return Ok(Vec::new());
    }
    for &v in items {
        if v < 0 {
            return Err(Error::InvalidInput);
        }
    }
    let mut arr = items.to_vec();
    let mut max_val = arr[0];
    for &v in &arr {
        if v > max_val {
            max_val = v;
        }
    }
    let mut exp = 1;
    while max_val / exp > 0 {
        counting_sort_by_digit(&mut arr, exp);
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

/// Standard-library sort (Timsort-like).
pub fn tim_sort<T: Ord + Clone>(items: &[T]) -> Vec<T> {
    let mut arr = items.to_vec();
    arr.sort();
    arr
}

// ---------------- Searching ----------------

fn check_sorted<T: Ord>(arr: &[T]) -> Result<(), Error> {
    for i in 1..arr.len() {
        if arr[i] < arr[i - 1] {
            return Err(Error::InvalidInput);
        }
    }
    Ok(())
}

/// Binary search.
pub fn binary_search<T: Ord>(arr: &[T], target: &T) -> Result<usize, Error> {
    check_sorted(arr)?;
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

/// Interpolation search for numeric data.
pub fn interpolation_search(arr: &[i64], target: i64) -> Result<usize, Error> {
    check_sorted(arr)?;
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
        let pos = lo_usize
            + ((target - arr[lo_usize]) as f64 / (arr[hi_usize] - arr[lo_usize]) as f64
                * (hi_usize - lo_usize) as f64) as usize;
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
pub fn jump_search<T: Ord>(arr: &[T], target: &T) -> Result<usize, Error> {
    check_sorted(arr)?;
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
    pub to: String,
    pub weight: f64,
}

impl Edge {
    pub fn new(to: &str, weight: f64) -> Self {
        Edge {
            to: to.to_string(),
            weight,
        }
    }
}

/// Breadth-first search.
pub fn bfs(graph: &HashMap<String, Vec<String>>, start: &str) -> Result<Vec<String>, Error> {
    if !graph.contains_key(start) {
        return Err(Error::InvalidGraph);
    }
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
pub fn dfs(graph: &HashMap<String, Vec<String>>, start: &str) -> Result<Vec<String>, Error> {
    if !graph.contains_key(start) {
        return Err(Error::InvalidGraph);
    }
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
    let mut distances: HashMap<String, f64> =
        graph.keys().map(|k| (k.clone(), f64::INFINITY)).collect();
    distances.insert(start.to_string(), 0.0);
    let mut predecessors = HashMap::new();

    let mut pq = BinaryHeap::new();
    pq.push(State {
        cost: 0.0,
        node: start.to_string(),
    });

    while let Some(State { cost, node }) = pq.pop() {
        if cost != distances[&node] {
            continue;
        }
        for e in graph.get(&node).unwrap_or(&Vec::new()) {
            if e.weight < 0.0 {
                return Err(Error::InvalidGraph);
            }
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
    Ok((distances, predecessors))
}

/// A* shortest path.
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
    let mut g_score = HashMap::new();
    g_score.insert(start.to_string(), 0.0);
    let mut came_from = HashMap::new();
    let mut pq = BinaryHeap::new();
    pq.push(State {
        cost: heuristic(start, goal),
        node: start.to_string(),
    });

    while let Some(State { cost: _, node }) = pq.pop() {
        if node == goal {
            return Ok((reconstruct_path(&came_from, &node), g_score[goal]));
        }
        for e in graph.get(&node).unwrap_or(&Vec::new()) {
            let tentative = g_score[&node] + e.weight;
            let current = g_score.get(&e.to).copied().unwrap_or(f64::INFINITY);
            if tentative < current {
                came_from.insert(e.to.clone(), node.clone());
                g_score.insert(e.to.clone(), tentative);
                pq.push(State {
                    cost: tentative + heuristic(&e.to, goal),
                    node: e.to.clone(),
                });
            }
        }
    }
    Err(Error::NotFound)
}

fn reconstruct_path(came_from: &HashMap<String, String>, start: &str) -> Vec<String> {
    let mut current = start.to_string();
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
pub fn knapsack_01(weights: &[i64], values: &[i64], capacity: i64) -> Result<i64, Error> {
    if weights.len() != values.len() {
        return Err(Error::InvalidInput);
    }
    if capacity < 0 {
        return Err(Error::InvalidInput);
    }
    let cap = capacity as usize;
    let mut dp = vec![0; cap + 1];
    for i in 0..weights.len() {
        let w = weights[i] as usize;
        if weights[i] < 0 {
            return Err(Error::InvalidInput);
        }
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

// ---------------- String algorithms ----------------

/// Knuth-Morris-Pratt pattern matching.
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
            matches.push(i - p.len() + 1);
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
pub fn rabin_karp_search(
    text: &str,
    pattern: &str,
    base: i64,
    modulus: i64,
) -> Result<Vec<usize>, Error> {
    if pattern.is_empty() {
        return Err(Error::InvalidInput);
    }
    let t = text.as_bytes();
    let p = pattern.as_bytes();
    let (n, m) = (t.len(), p.len());
    if m > n {
        return Ok(Vec::new());
    }
    let mut h = 1i64;
    for _ in 0..m - 1 {
        h = (h * base) % modulus;
    }
    let mut pattern_hash = 0i64;
    let mut text_hash = 0i64;
    for i in 0..m {
        pattern_hash = (pattern_hash * base + p[i] as i64) % modulus;
        text_hash = (text_hash * base + t[i] as i64) % modulus;
    }
    let mut matches = Vec::new();
    for i in 0..=n - m {
        if text_hash == pattern_hash && &t[i..i + m] == p {
            matches.push(i);
        }
        if i < n - m {
            text_hash = (text_hash - t[i] as i64 * h) % modulus;
            text_hash = (text_hash * base + t[i + m] as i64) % modulus;
            text_hash = ((text_hash % modulus) + modulus) % modulus;
        }
    }
    Ok(matches)
}

/// Boyer-Moore pattern matching with the bad-character rule.
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
    let mut matches = Vec::new();
    let mut i = 0usize;
    while i <= n - m {
        let mut j = (m - 1) as isize;
        while j >= 0 && t[i + j as usize] == p[j as usize] {
            j -= 1;
        }
        if j < 0 {
            matches.push(i);
            i += m;
        } else {
            let ju = j as usize;
            let last = bad_char.get(&t[i + ju]).copied().unwrap_or(usize::MAX);
            let mut shift = if ju > last { ju - last } else { 1 };
            if shift < 1 {
                shift = 1;
            }
            i += shift;
        }
    }
    Ok(matches)
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
        list.push(Indexed { index: i, start, end });
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
pub fn fractional_knapsack(
    weights: &[f64],
    values: &[f64],
    capacity: f64,
) -> Result<f64, Error> {
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
    for i in 0..weights.len() {
        if weights[i] > 0.0 {
            items.push(Item {
                ratio: values[i] / weights[i],
                weight: weights[i],
            });
        }
    }
    items.sort_by(|a, b| {
        b.ratio
            .partial_cmp(&a.ratio)
            .unwrap_or(Ordering::Equal)
    });
    let mut total = 0.0;
    let mut remaining = capacity;
    for it in items {
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
pub fn huffman_coding(frequencies: &HashMap<String, i64>) -> Result<HashMap<String, String>, Error> {
    if frequencies.is_empty() {
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
    if !node.char.is_empty() {
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

/// Exponentiation by squaring.
pub fn fast_power(base: f64, exponent: i64) -> f64 {
    if exponent < 0 {
        return 1.0 / fast_power(base, -exponent);
    }
    if exponent == 0 {
        return 1.0;
    }
    if exponent == 1 {
        return base;
    }
    let half = fast_power(base, exponent / 2);
    if exponent % 2 == 0 {
        half * half
    } else {
        base * half * half
    }
}
