use algorithms_lib::Error;
use algorithms_lib::*;
use std::collections::HashMap;

#[test]
fn test_quick_sort() {
    assert_eq!(
        quick_sort(&[3, 1, 4, 1, 5, 9, 2, 6]),
        vec![1, 1, 2, 3, 4, 5, 6, 9]
    );
    assert_eq!(
        quick_sort(&["cherry", "apple", "banana"]),
        vec!["apple", "banana", "cherry"]
    );
}

#[test]
fn test_merge_sort() {
    assert_eq!(
        merge_sort(&[3, 1, 4, 1, 5, 9, 2, 6]),
        vec![1, 1, 2, 3, 4, 5, 6, 9]
    );
}

#[test]
fn test_heap_sort() {
    assert_eq!(
        heap_sort(&[3, 1, 4, 1, 5, 9, 2, 6]),
        vec![1, 1, 2, 3, 4, 5, 6, 9]
    );
}

#[test]
fn test_radix_sort() {
    assert_eq!(
        radix_sort(&[170, 45, 75, 90, 2, 802, 24, 66]).unwrap(),
        vec![2, 24, 45, 66, 75, 90, 170, 802]
    );
    assert_eq!(radix_sort(&[1, -5, 3]), Err(Error::InvalidInput));
}

#[test]
fn test_tim_sort() {
    assert_eq!(
        tim_sort(&[3, 1, 4, 1, 5, 9, 2, 6]),
        vec![1, 1, 2, 3, 4, 5, 6, 9]
    );
}

#[test]
fn test_binary_search() {
    let arr = [2, 4, 6, 8, 10, 12];
    assert_eq!(binary_search(&arr, &8).unwrap(), 3);
    assert_eq!(binary_search(&arr, &7), Err(Error::NotFound));
    assert_eq!(binary_search(&[3, 1, 2], &2), Err(Error::InvalidInput));
}

#[test]
fn test_interpolation_search() {
    let arr = [10, 20, 30, 40, 50, 60, 70, 80, 90, 100];
    assert_eq!(interpolation_search(&arr, 50).unwrap(), 4);
    assert_eq!(interpolation_search(&arr, 55), Err(Error::NotFound));
}

#[test]
fn test_jump_search() {
    let arr = [2, 4, 6, 8, 10, 12];
    assert_eq!(jump_search(&arr, &10).unwrap(), 4);
    assert_eq!(jump_search(&arr, &7), Err(Error::NotFound));
}

#[test]
fn test_bfs() {
    let mut graph = HashMap::new();
    graph.insert("a".to_string(), vec!["b".to_string(), "c".to_string()]);
    graph.insert("b".to_string(), vec!["d".to_string()]);
    graph.insert("c".to_string(), Vec::new());
    graph.insert("d".to_string(), Vec::new());
    assert_eq!(
        bfs(&graph, "a").unwrap(),
        vec!["a", "b", "c", "d"]
    );
    assert_eq!(bfs(&graph, "z"), Err(Error::InvalidGraph));
}

#[test]
fn test_dfs() {
    let mut graph = HashMap::new();
    graph.insert("a".to_string(), vec!["b".to_string(), "c".to_string()]);
    graph.insert("b".to_string(), vec!["d".to_string()]);
    graph.insert("c".to_string(), Vec::new());
    graph.insert("d".to_string(), Vec::new());
    assert_eq!(
        dfs(&graph, "a").unwrap(),
        vec!["a", "b", "d", "c"]
    );
}

#[test]
fn test_dijkstra() {
    let mut graph = HashMap::new();
    graph.insert(
        "a".to_string(),
        vec![Edge::new("b", 1.0), Edge::new("c", 4.0)],
    );
    graph.insert(
        "b".to_string(),
        vec![Edge::new("c", 2.0), Edge::new("d", 5.0)],
    );
    graph.insert("c".to_string(), vec![Edge::new("d", 1.0)]);
    graph.insert("d".to_string(), Vec::new());
    let (distances, _) = dijkstra(&graph, "a").unwrap();
    assert_eq!(distances["d"], 4.0);

    let mut negative = HashMap::new();
    negative.insert("a".to_string(), vec![Edge::new("b", -1.0)]);
    negative.insert("b".to_string(), Vec::new());
    assert_eq!(dijkstra(&negative, "a"), Err(Error::InvalidGraph));
}

#[test]
fn test_a_star() {
    let mut graph = HashMap::new();
    graph.insert(
        "a".to_string(),
        vec![Edge::new("b", 1.0), Edge::new("c", 4.0)],
    );
    graph.insert(
        "b".to_string(),
        vec![Edge::new("c", 2.0), Edge::new("d", 5.0)],
    );
    graph.insert("c".to_string(), vec![Edge::new("d", 1.0)]);
    graph.insert("d".to_string(), Vec::new());
    let (path, cost) = a_star(&graph, "a", "d", |_, _| 0.0).unwrap();
    assert_eq!(path, vec!["a", "b", "c", "d"]);
    assert_eq!(cost, 4.0);

    let mut empty = HashMap::new();
    empty.insert("a".to_string(), Vec::new());
    empty.insert("b".to_string(), Vec::new());
    assert_eq!(a_star(&empty, "a", "b", |_, _| 0.0), Err(Error::NotFound));
}

#[test]
fn test_bellman_ford() {
    let mut graph = HashMap::new();
    graph.insert("a".to_string(), vec![Edge::new("b", -1.0)]);
    graph.insert("b".to_string(), vec![Edge::new("c", -2.0)]);
    graph.insert("c".to_string(), vec![Edge::new("d", 1.0)]);
    graph.insert("d".to_string(), Vec::new());
    let (distances, _) = bellman_ford(&graph, "a").unwrap();
    assert_eq!(distances["d"], -2.0);

    let mut negative = HashMap::new();
    negative.insert("a".to_string(), vec![Edge::new("b", 1.0)]);
    negative.insert("b".to_string(), vec![Edge::new("c", -1.0)]);
    negative.insert("c".to_string(), vec![Edge::new("b", -1.0)]);
    assert_eq!(bellman_ford(&negative, "a"), Err(Error::NegativeCycle));
}

#[test]
fn test_knapsack_01() {
    assert_eq!(
        knapsack_01(&[1, 2, 3], &[6, 10, 12], 5).unwrap(),
        22
    );
}

#[test]
fn test_longest_common_subsequence() {
    assert_eq!(
        longest_common_subsequence(b"ABCDE", b"ACE"),
        3
    );
}

#[test]
fn test_edit_distance() {
    assert_eq!(
        edit_distance(b"kitten", b"sitting"),
        3
    );
}

#[test]
fn test_kmp_search() {
    let text = "ABABDABACDABABCABAB";
    let pattern = "ABABCABAB";
    assert_eq!(kmp_search(text, pattern).unwrap(), vec![10]);
}

#[test]
fn test_rabin_karp_search() {
    let text = "ABABDABACDABABCABAB";
    let pattern = "ABABCABAB";
    assert_eq!(
        rabin_karp_search(text, pattern, 256, 1_000_000_007).unwrap(),
        vec![10]
    );
}

#[test]
fn test_boyer_moore_search() {
    let text = "ABABDABACDABABCABAB";
    let pattern = "ABABCABAB";
    assert_eq!(boyer_moore_search(text, pattern).unwrap(), vec![10]);
}

#[test]
fn test_activity_selection() {
    let activities = [
        (1, 3),
        (2, 5),
        (4, 7),
        (1, 8),
        (5, 9),
        (8, 10),
        (9, 11),
        (11, 14),
        (13, 16),
    ];
    let selected = activity_selection(&activities).unwrap();
    assert!(selected.len() >= 4);
}

#[test]
fn test_fractional_knapsack() {
    let total = fractional_knapsack(&[10.0, 20.0, 30.0], &[60.0, 100.0, 120.0], 50.0).unwrap();
    assert!((total - 240.0).abs() < 1e-9);
}

#[test]
fn test_huffman_coding() {
    let mut frequencies = HashMap::new();
    frequencies.insert("a".to_string(), 5);
    frequencies.insert("b".to_string(), 9);
    frequencies.insert("c".to_string(), 12);
    frequencies.insert("d".to_string(), 13);
    frequencies.insert("e".to_string(), 16);
    frequencies.insert("f".to_string(), 45);
    let codes = huffman_coding(&frequencies).unwrap();
    assert_eq!(codes.len(), 6);
    for (_, code) in codes {
        assert!(!code.is_empty());
    }
}

#[test]
fn test_max_subarray() {
    assert_eq!(
        max_subarray(&[-2, 1, -3, 4, -1, 2, 1, -5, 4]).unwrap(),
        6
    );
    assert_eq!(max_subarray(&[-2, -1]).unwrap(), -1);
}

#[test]
fn test_count_inversions() {
    assert_eq!(count_inversions(&[1, 3, 5, 2, 4, 6]), 3);
}

#[test]
fn test_fast_power() {
    assert!((fast_power(2.0, 10) - 1024.0).abs() < 1e-9);
    assert!((fast_power(5.0, 0) - 1.0).abs() < 1e-9);
    assert!((fast_power(2.0, -2) - 0.25).abs() < 1e-9);
}
