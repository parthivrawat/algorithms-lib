use algorithms_lib::*;
use std::collections::HashMap;

#[test]
fn edge_sorts_empty_single_and_duplicates() {
    let empty: &[i64] = &[];

    // Empty inputs
    assert!(quick_sort(empty).is_empty());
    assert!(merge_sort(empty).is_empty());
    assert!(heap_sort(empty).is_empty());
    assert!(native_sort(empty).is_empty());
    assert_eq!(radix_sort(&[]).unwrap(), Vec::<i64>::new());

    // Single-element inputs
    assert_eq!(quick_sort(&[7]), vec![7]);
    assert_eq!(merge_sort(&[7]), vec![7]);
    assert_eq!(heap_sort(&[7]), vec![7]);
    assert_eq!(native_sort(&[7]), vec![7]);
    assert_eq!(radix_sort(&[7]).unwrap(), vec![7]);

    // Duplicates
    let dups = [3, 1, 2, 1, 3, 2, 1];
    let want = [1, 1, 1, 2, 2, 3, 3];
    assert_eq!(quick_sort(&dups), want);
    assert_eq!(merge_sort(&dups), want);
    assert_eq!(heap_sort(&dups), want);
    assert_eq!(native_sort(&dups), want);
    assert_eq!(radix_sort(&[3, 1, 2, 1, 3, 2, 1]).unwrap(), want);
}

#[test]
fn edge_searches_empty_single_and_duplicates() {
    let empty: &[i64] = &[];

    // Binary / jump / interpolation search on empty and single-element arrays
    assert_eq!(binary_search(empty, &5), Err(Error::NotFound));
    assert_eq!(binary_search(&[5], &5), Ok(0));
    assert_eq!(binary_search(&[5], &3), Err(Error::NotFound));
    assert_eq!(binary_search(&[1, 1, 1, 1], &1).unwrap() < 4, true);

    assert_eq!(jump_search(empty, &5), Err(Error::NotFound));
    assert_eq!(jump_search(&[5], &5), Ok(0));
    assert_eq!(jump_search(&[5, 5, 5, 5], &5).unwrap() < 4, true);

    assert_eq!(interpolation_search(empty, 5), Err(Error::NotFound));
    assert_eq!(interpolation_search(&[5], 5), Ok(0));
    assert_eq!(
        interpolation_search(&[10, 10, 10, 20], 10).unwrap() < 4,
        true
    );
}

#[test]
fn edge_string_search_overlapping() {
    // Overlapping occurrences for all three string searchers
    assert_eq!(kmp_search("aaaa", "aa").unwrap(), vec![0, 1, 2]);
    assert_eq!(kmp_search("aaaaa", "aaaa").unwrap(), vec![0, 1]);

    assert_eq!(boyer_moore_search("aaaa", "aa").unwrap(), vec![0, 1, 2]);
    assert_eq!(boyer_moore_search("aaaaa", "aaaa").unwrap(), vec![0, 1]);

    assert_eq!(
        rabin_karp_search("aaaa", "aa", 256, 1_000_000_007).unwrap(),
        vec![0, 1, 2]
    );
    assert_eq!(
        rabin_karp_search("aaaaa", "aaaa", 256, 1_000_000_007).unwrap(),
        vec![0, 1]
    );

    // Pattern longer than text
    assert_eq!(kmp_search("ab", "abc").unwrap(), Vec::<usize>::new());
    assert_eq!(
        boyer_moore_search("ab", "abc").unwrap(),
        Vec::<usize>::new()
    );
    assert_eq!(
        rabin_karp_search("ab", "abc", 256, 1_000_000_007).unwrap(),
        Vec::<usize>::new()
    );
}

#[test]
fn edge_bfs_dfs_disconnected() {
    let mut graph: HashMap<String, Vec<String>> = HashMap::new();
    graph.insert("a".to_string(), vec!["b".to_string(), "c".to_string()]);
    graph.insert("b".to_string(), Vec::new());
    graph.insert("c".to_string(), Vec::new());
    graph.insert("d".to_string(), vec!["e".to_string()]);
    graph.insert("e".to_string(), Vec::new());

    // BFS / DFS from a should stay within the {a,b,c} component
    assert_eq!(bfs(&graph, "a").unwrap(), vec!["a", "b", "c"]);
    assert_eq!(dfs(&graph, "a").unwrap(), vec!["a", "b", "c"]);

    // BFS / DFS from d should stay within the {d,e} component
    assert_eq!(bfs(&graph, "d").unwrap(), vec!["d", "e"]);
    assert_eq!(dfs(&graph, "d").unwrap(), vec!["d", "e"]);
}

#[test]
fn edge_dijkstra_and_astar_disconnected_and_start_equals_goal() {
    let mut graph: HashMap<String, Vec<Edge>> = HashMap::new();
    graph.insert("a".to_string(), vec![Edge::new("b", 1.0)]);
    graph.insert("b".to_string(), Vec::new());
    graph.insert("c".to_string(), vec![Edge::new("d", 2.0)]);
    graph.insert("d".to_string(), Vec::new());

    let (distances, _pred) = dijkstra(&graph, "a").unwrap();
    assert_eq!(distances["a"], 0.0);
    assert_eq!(distances["b"], 1.0);
    assert!(distances["c"].is_infinite());
    assert!(distances["d"].is_infinite());

    // A* between disconnected components is not reachable
    assert_eq!(a_star(&graph, "a", "c", |_, _| 0.0), Err(Error::NotFound));

    // start == goal
    let mut solo: HashMap<String, Vec<Edge>> = HashMap::new();
    solo.insert("a".to_string(), Vec::new());
    let (distances, _pred) = dijkstra(&solo, "a").unwrap();
    assert_eq!(distances["a"], 0.0);

    let (path, cost) = a_star(&solo, "a", "a", |_, _| 0.0).unwrap();
    assert_eq!(path, vec!["a".to_string()]);
    assert_eq!(cost, 0.0);
}

#[test]
fn edge_knapsack_zero_capacity() {
    // 0/1 knapsack with zero capacity
    assert_eq!(knapsack_01(&[1, 2], &[10, 20], 0).unwrap(), 0);
    assert_eq!(knapsack_01(&[], &[], 0).unwrap(), 0);
    // Zero-weight items can still be taken at zero capacity
    assert_eq!(knapsack_01(&[0, 2], &[50, 20], 0).unwrap(), 50);
    assert_eq!(knapsack_01(&[0], &[-5], 0).unwrap(), 0);
    assert_eq!(knapsack_01(&[1], &[10], -1), Err(Error::InvalidInput));

    // Fractional knapsack with zero capacity
    assert_eq!(
        fractional_knapsack(&[0.0, 10.0], &[50.0, 60.0], 0.0).unwrap(),
        50.0
    );
    assert_eq!(fractional_knapsack(&[], &[], 0.0).unwrap(), 0.0);
    assert_eq!(fractional_knapsack(&[0.0], &[-10.0], 0.0).unwrap(), 0.0);
    assert_eq!(fractional_knapsack(&[10.0], &[60.0], 0.0).unwrap(), 0.0);
    assert_eq!(
        fractional_knapsack(&[10.0], &[60.0], -5.0),
        Err(Error::InvalidInput)
    );
}

#[test]
fn edge_huffman_single_symbol() {
    let mut freq = HashMap::new();
    freq.insert("x".to_string(), 1);
    let codes = huffman_coding(&freq).unwrap();
    assert_eq!(codes.len(), 1);
    assert_eq!(codes["x"], "0");

    // Empty string as the only symbol is still valid
    let mut empty = HashMap::new();
    empty.insert(String::new(), 5);
    let codes = huffman_coding(&empty).unwrap();
    assert_eq!(codes.len(), 1);
    assert_eq!(codes[""], "0");

    // Empty frequency map is rejected
    assert_eq!(
        huffman_coding(&HashMap::<String, i64>::new()),
        Err(Error::InvalidInput)
    );
}

#[test]
fn edge_fast_power_extreme() {
    // 0 raised to a negative exponent is undefined and should fail
    assert_eq!(fast_power(0.0, -1), Err(Error::InvalidInput));
    assert_eq!(fast_power(0.0, -5), Err(Error::InvalidInput));
    assert_eq!(fast_power(0.0, i64::MIN), Err(Error::InvalidInput));

    // 0^0 and 0^positive are defined
    assert_eq!(fast_power(0.0, 0).unwrap(), 1.0);
    assert_eq!(fast_power(0.0, 5).unwrap(), 0.0);

    // i64::MIN exponent must not overflow or hang
    assert_eq!(fast_power(2.0, i64::MIN).unwrap(), 0.0);
    assert_eq!(fast_power(-1.0, i64::MIN).unwrap(), 1.0);
    assert_eq!(fast_power(1.0, i64::MIN).unwrap(), 1.0);
}

#[test]
fn edge_radix_sort_extreme() {
    // Empty, single, and all-duplicates
    assert_eq!(radix_sort(&[]).unwrap(), Vec::<i64>::new());
    assert_eq!(radix_sort(&[5]).unwrap(), vec![5]);
    assert_eq!(radix_sort(&[0, 0, 0]).unwrap(), vec![0, 0, 0]);

    // Largest possible i64 values should not overflow
    assert_eq!(radix_sort(&[i64::MAX]).unwrap(), vec![i64::MAX]);
    assert_eq!(
        radix_sort(&[i64::MAX, i64::MAX, 0]).unwrap(),
        vec![0, i64::MAX, i64::MAX]
    );

    // Negative values are rejected
    assert_eq!(radix_sort(&[i64::MIN]), Err(Error::InvalidInput));
    assert_eq!(radix_sort(&[1, -1, 2]), Err(Error::InvalidInput));
}
