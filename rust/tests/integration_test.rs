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
    // exp * 10 must not overflow i64 before the digit loop terminates.
    assert_eq!(
        radix_sort(&[i64::MAX, 1, 0, i64::MAX - 1]).unwrap(),
        vec![0, 1, i64::MAX - 1, i64::MAX]
    );
}

#[test]
fn test_native_sort() {
    assert_eq!(
        native_sort(&[3, 1, 4, 1, 5, 9, 2, 6]),
        vec![1, 1, 2, 3, 4, 5, 6, 9]
    );
}

#[test]
fn test_in_place_sorts() {
    let want = vec![1, 1, 2, 3, 4, 5, 6, 9];

    let mut a = vec![3, 1, 4, 1, 5, 9, 2, 6];
    quick_sort_in_place(&mut a);
    assert_eq!(a, want);

    let mut b = vec![3, 1, 4, 1, 5, 9, 2, 6];
    merge_sort_in_place(&mut b);
    assert_eq!(b, want);

    let mut c = vec![3, 1, 4, 1, 5, 9, 2, 6];
    heap_sort_in_place(&mut c);
    assert_eq!(c, want);

    let mut d = vec![3, 1, 4, 1, 5, 9, 2, 6];
    native_sort_in_place(&mut d);
    assert_eq!(d, want);
}

#[test]
fn test_binary_search() {
    let arr = [2, 4, 6, 8, 10, 12];
    assert_eq!(binary_search(&arr, &8).unwrap(), 3);
    assert_eq!(binary_search(&arr, &7), Err(Error::NotFound));
}

#[test]
#[cfg(debug_assertions)]
#[should_panic(expected = "binary_search requires a sorted slice")]
fn test_binary_search_unsorted_debug_assert() {
    // Sortedness is a caller contract checked with debug_assert! in debug
    // builds; in release builds the result on unsorted input is unspecified.
    let _ = binary_search(&[3, 1, 2], &2);
}

#[test]
fn test_interpolation_search() {
    let arr = [10, 20, 30, 40, 50, 60, 70, 80, 90, 100];
    assert_eq!(interpolation_search(&arr, 50).unwrap(), 4);
    assert_eq!(interpolation_search(&arr, 55), Err(Error::NotFound));

    // Values near i64::MAX expose precision loss and overflow in the naive
    // float/difference computation.
    let arr = [
        i64::MAX - 4,
        i64::MAX - 3,
        i64::MAX - 2,
        i64::MAX - 1,
        i64::MAX,
    ];
    assert_eq!(interpolation_search(&arr, i64::MAX - 2).unwrap(), 2);
    assert_eq!(
        interpolation_search(&arr, i64::MAX - 10),
        Err(Error::NotFound)
    );

    // Extreme negative-to-positive range would overflow i64 subtraction.
    let arr = [i64::MIN, -1, 0, 1, i64::MAX];
    assert_eq!(interpolation_search(&arr, i64::MAX).unwrap(), 4);
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
    assert_eq!(bfs(&graph, "a").unwrap(), vec!["a", "b", "c", "d"]);
    assert_eq!(bfs(&graph, "z"), Err(Error::InvalidGraph));

    // Neighbors that are not keys of the map must be rejected.
    let mut bad: HashMap<String, Vec<String>> = HashMap::new();
    bad.insert("a".to_string(), vec!["zz".to_string()]);
    assert_eq!(bfs(&bad, "a"), Err(Error::InvalidGraph));
    assert_eq!(dfs(&bad, "a"), Err(Error::InvalidGraph));
}

#[test]
fn test_dfs() {
    let mut graph = HashMap::new();
    graph.insert("a".to_string(), vec!["b".to_string(), "c".to_string()]);
    graph.insert("b".to_string(), vec!["d".to_string()]);
    graph.insert("c".to_string(), Vec::new());
    graph.insert("d".to_string(), Vec::new());
    assert_eq!(dfs(&graph, "a").unwrap(), vec!["a", "b", "d", "c"]);
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

    let mut unreachable = HashMap::new();
    unreachable.insert("a".to_string(), Vec::new());
    unreachable.insert("b".to_string(), vec![Edge::new("c", -1.0)]);
    unreachable.insert("c".to_string(), Vec::new());
    assert_eq!(dijkstra(&unreachable, "a"), Err(Error::InvalidGraph));
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

    // A* does not support negative edge weights.
    let mut negative = HashMap::new();
    negative.insert("a".to_string(), vec![Edge::new("b", -1.0)]);
    negative.insert("b".to_string(), Vec::new());
    assert_eq!(
        a_star(&negative, "a", "b", |_, _| 0.0),
        Err(Error::InvalidGraph)
    );

    // Negative heuristic values are rejected immediately.
    assert_eq!(
        a_star(&empty, "a", "b", |_, _| -1.0),
        Err(Error::InvalidInput)
    );
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
    assert_eq!(knapsack_01(&[1, 2, 3], &[6, 10, 12], 5).unwrap(), 22);
    // Negative weights and capacity are rejected, not silently truncated.
    assert_eq!(knapsack_01(&[-1, 2], &[6, 10], 5), Err(Error::InvalidInput));
    assert_eq!(knapsack_01(&[1], &[6], -1), Err(Error::InvalidInput));
}

#[test]
fn test_longest_common_subsequence() {
    assert_eq!(longest_common_subsequence(b"ABCDE", b"ACE"), 3);
}

#[test]
fn test_edit_distance() {
    assert_eq!(edit_distance(b"kitten", b"sitting"), 3);
}

#[test]
fn test_edit_distance_reconstruction() {
    let (dist, script) = edit_distance_reconstruction(b"kitten", b"sitting");
    assert_eq!(dist, 3);
    assert!(script.len() >= 3);
}

#[test]
fn test_longest_common_subsequence_reconstruction() {
    let (length, lcs) = longest_common_subsequence_reconstruction(b"ABCDE", b"ACE");
    assert_eq!(length, 3);
    assert_eq!(lcs, vec![b'A', b'C', b'E']);
}

#[test]
fn test_kmp_search() {
    let text = "ABABDABACDABABCABAB";
    let pattern = "ABABCABAB";
    assert_eq!(kmp_search(text, pattern).unwrap(), vec![10]);
    // Empty patterns are rejected.
    assert_eq!(kmp_search(text, ""), Err(Error::InvalidInput));
    // Offsets are byte offsets, not char indices ('é' is 2 bytes in UTF-8).
    assert_eq!(kmp_search("héllo héllo", "éllo").unwrap(), vec![1, 8]);
}

#[test]
fn test_rabin_karp_search() {
    let text = "ABABDABACDABABCABAB";
    let pattern = "ABABCABAB";
    assert_eq!(
        rabin_karp_search(text, pattern, 256, 1_000_000_007).unwrap(),
        vec![10]
    );

    assert_eq!(
        rabin_karp_search(text, pattern, 256, 0),
        Err(Error::InvalidInput)
    );
    assert_eq!(
        rabin_karp_search(text, pattern, -1, 1_000_000_007),
        Err(Error::InvalidInput)
    );

    // Overlapping matches and a modulus that would overflow i64 arithmetic.
    assert_eq!(
        rabin_karp_search("aaaaa", "aaaa", 256, 1_000_000_007).unwrap(),
        vec![0, 1]
    );
    assert_eq!(
        rabin_karp_search("hello", "ll", 256, i64::MAX).unwrap(),
        vec![2]
    );
}

#[test]
fn test_boyer_moore_search() {
    let text = "ABABDABACDABABCABAB";
    let pattern = "ABABCABAB";
    assert_eq!(boyer_moore_search(text, pattern).unwrap(), vec![10]);
}

#[test]
fn test_boyer_moore_search_overlapping() {
    assert_eq!(boyer_moore_search("aaaa", "aa").unwrap(), vec![0, 1, 2]);
    assert_eq!(boyer_moore_search("aaaaa", "aaaa").unwrap(), vec![0, 1]);
    assert_eq!(
        boyer_moore_search("abcabcabc", "abcabc").unwrap(),
        vec![0, 3]
    );
    assert_eq!(
        boyer_moore_search("abababab", "abab").unwrap(),
        vec![0, 2, 4]
    );
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
fn test_fractional_knapsack_edge_cases() {
    // Zero-weight, positive-value items are taken in full for free.
    assert_eq!(
        fractional_knapsack(&[0.0, 10.0], &[50.0, 60.0], 5.0).unwrap(),
        80.0
    );
    // Zero-weight items with non-positive value are never taken.
    assert_eq!(
        fractional_knapsack(&[0.0, 10.0], &[-50.0, 60.0], 50.0).unwrap(),
        60.0
    );
    // Negative weights are rejected.
    assert_eq!(
        fractional_knapsack(&[-5.0, 10.0], &[10.0, 60.0], 50.0),
        Err(Error::InvalidInput)
    );
    // Negative-value items are never taken, so the total never goes below zero.
    assert_eq!(fractional_knapsack(&[10.0], &[-5.0], 10.0).unwrap(), 0.0);
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
fn test_huffman_coding_edge_cases() {
    // The empty string is a valid symbol key and must appear in the code table.
    let mut freq = HashMap::new();
    freq.insert(String::new(), 2);
    freq.insert("a".to_string(), 3);
    freq.insert("b".to_string(), 1);
    let codes = huffman_coding(&freq).unwrap();
    assert_eq!(codes.len(), 3);
    assert!(codes.contains_key(""));

    // Non-positive frequencies are rejected.
    let mut bad = HashMap::new();
    bad.insert("a".to_string(), 0);
    assert_eq!(huffman_coding(&bad), Err(Error::InvalidInput));
    bad.insert("a".to_string(), -1);
    assert_eq!(huffman_coding(&bad), Err(Error::InvalidInput));
}

#[test]
fn test_huffman_encode_decode() {
    let mut frequencies = HashMap::new();
    frequencies.insert("a".to_string(), 5);
    frequencies.insert("b".to_string(), 9);
    frequencies.insert("c".to_string(), 12);
    frequencies.insert("d".to_string(), 13);
    frequencies.insert("e".to_string(), 16);
    frequencies.insert("f".to_string(), 45);
    let codes = huffman_coding(&frequencies).unwrap();
    let symbols = vec![
        "a".to_string(),
        "b".to_string(),
        "c".to_string(),
        "d".to_string(),
        "e".to_string(),
        "f".to_string(),
    ];
    let encoded = huffman_encode(&symbols, &codes).unwrap();
    let decoded = huffman_decode(&encoded, &codes).unwrap();
    assert_eq!(decoded, symbols);
}

#[test]
fn test_max_subarray() {
    assert_eq!(max_subarray(&[-2, 1, -3, 4, -1, 2, 1, -5, 4]).unwrap(), 6);
    assert_eq!(max_subarray(&[-2, -1]).unwrap(), -1);
}

#[test]
fn test_count_inversions() {
    assert_eq!(count_inversions(&[1, 3, 5, 2, 4, 6]), 3);
}

#[test]
fn test_fast_power() {
    assert!((fast_power(2.0, 10).unwrap() - 1024.0).abs() < 1e-9);
    assert!((fast_power(5.0, 0).unwrap() - 1.0).abs() < 1e-9);
    assert!((fast_power(2.0, -2).unwrap() - 0.25).abs() < 1e-9);
}

#[test]
fn test_fast_power_edge_cases() {
    assert_eq!(fast_power(0.0, -1), Err(Error::InvalidInput));
    assert_eq!(fast_power(0.0, 0).unwrap(), 1.0);
    // i64::MIN cannot be negated in i64; must not hang or panic.
    assert_eq!(fast_power(2.0, i64::MIN).unwrap(), 0.0);
    assert_eq!(fast_power(-2.0, 3).unwrap(), -8.0);
    assert_eq!(fast_power(-2.0, 4).unwrap(), 16.0);
}

#[test]
fn test_gcd() {
    assert_eq!(gcd(48, 18).unwrap(), 6);
    assert_eq!(gcd(-48, 18).unwrap(), 6);
    assert_eq!(gcd(0, 5).unwrap(), 5);
    assert_eq!(gcd(0, 0).unwrap(), 0);
}

#[test]
fn test_lcm() {
    assert_eq!(lcm(4, 6).unwrap(), 12);
    assert_eq!(lcm(-4, 6).unwrap(), 12);
    assert_eq!(lcm(0, 5).unwrap(), 0);
    assert_eq!(lcm(0, 0), Err(Error::InvalidInput));
}

#[test]
fn test_modular_fast_power() {
    assert_eq!(modular_fast_power(2, 10, 1000).unwrap(), 24);
    assert_eq!(modular_fast_power(3, 0, 7).unwrap(), 1);
    assert_eq!(modular_fast_power(0, 0, 7).unwrap(), 1);
    assert_eq!(modular_fast_power(0, 5, 7).unwrap(), 0);
    assert_eq!(modular_fast_power(-2, 3, 1000).unwrap(), 992);
    assert_eq!(modular_fast_power(2, -1, 7), Err(Error::InvalidInput));
    assert_eq!(modular_fast_power(2, 2, 0), Err(Error::InvalidInput));
}

#[test]
fn test_quick_select() {
    assert_eq!(quick_select(&[3, 2, 1, 5, 4], 0).unwrap(), 1);
    assert_eq!(quick_select(&[3, 2, 1, 5, 4], 2).unwrap(), 3);
    assert_eq!(quick_select(&[3, 2, 1, 5, 4], 4).unwrap(), 5);
    assert_eq!(quick_select(&["cherry", "apple", "banana"], 1).unwrap(), "banana");
    assert_eq!(quick_select(&[1, 2, 3], 5), Err(Error::InvalidInput));
}

#[test]
fn test_sieve_of_eratosthenes() {
    assert_eq!(sieve_of_eratosthenes(10).unwrap(), vec![2, 3, 5, 7]);
    assert!(sieve_of_eratosthenes(1).unwrap().is_empty());
    assert_eq!(sieve_of_eratosthenes(-1), Err(Error::InvalidInput));
}

#[test]
fn test_topological_sort() {
    let mut graph: HashMap<String, Vec<String>> = HashMap::new();
    graph.insert("a".to_string(), vec!["b".to_string(), "c".to_string()]);
    graph.insert("b".to_string(), vec!["d".to_string()]);
    graph.insert("c".to_string(), vec!["d".to_string()]);
    graph.insert("d".to_string(), Vec::new());
    let order = topological_sort(&graph).unwrap();
    let pos: HashMap<String, usize> = order
        .iter()
        .enumerate()
        .map(|(i, n)| (n.clone(), i))
        .collect();
    assert!(pos["a"] < pos["b"]);
    assert!(pos["a"] < pos["c"]);
    assert!(pos["b"] < pos["d"]);
    assert!(pos["c"] < pos["d"]);

    let mut cyclic: HashMap<String, Vec<String>> = HashMap::new();
    cyclic.insert("a".to_string(), vec!["b".to_string()]);
    cyclic.insert("b".to_string(), vec!["c".to_string()]);
    cyclic.insert("c".to_string(), vec!["a".to_string()]);
    assert_eq!(topological_sort(&cyclic), Err(Error::InvalidGraph));
}

#[test]
fn test_union_find() {
    let mut uf = UnionFind::new(["a", "b", "c", "d", "e"]);
    assert!(!uf.connected(&"a", &"b").unwrap());
    assert!(uf.union(&"a", &"b").unwrap());
    assert!(uf.connected(&"a", &"b").unwrap());
    assert!(uf.union(&"b", &"c").unwrap());
    assert!(uf.connected(&"a", &"c").unwrap());
    assert!(!uf.connected(&"a", &"d").unwrap());
    assert_eq!(uf.count(), 3);
    assert!(!uf.union(&"a", &"c").unwrap());
    assert_eq!(uf.find(&"z"), Err(Error::InvalidInput));
}

#[test]
fn test_priority_queue() {
    let mut pq = PriorityQueue::new([3, 1, 4, 1, 5]);
    assert_eq!(*pq.peek().unwrap(), 1);
    assert_eq!(pq.len(), 5);
    assert_eq!(pq.pop().unwrap(), 1);
    assert_eq!(pq.pop().unwrap(), 1);
    assert_eq!(pq.pop().unwrap(), 3);
    pq.push(0);
    assert_eq!(pq.pop().unwrap(), 0);
    let mut empty: PriorityQueue<i64> = PriorityQueue::new([]);
    assert_eq!(empty.pop(), Err(Error::InvalidInput));
}

#[test]
fn test_minimum_spanning_tree() {
    let edges = [
        MSTEdge { u: "a".to_string(), v: "b".to_string(), weight: 1.0 },
        MSTEdge { u: "b".to_string(), v: "c".to_string(), weight: 2.0 },
        MSTEdge { u: "a".to_string(), v: "c".to_string(), weight: 3.0 },
        MSTEdge { u: "c".to_string(), v: "d".to_string(), weight: 4.0 },
    ];
    let (total, mst) = minimum_spanning_tree(&edges).unwrap();
    assert!((total - 7.0).abs() < 1e-9);
    assert_eq!(mst.len(), 3);
}

#[test]
fn test_floyd_warshall() {
    let mut graph: HashMap<String, Vec<Edge>> = HashMap::new();
    graph.insert("a".to_string(), vec![
        Edge { to: "b".to_string(), weight: 1.0 },
        Edge { to: "c".to_string(), weight: 4.0 },
    ]);
    graph.insert("b".to_string(), vec![
        Edge { to: "c".to_string(), weight: 2.0 },
        Edge { to: "d".to_string(), weight: 5.0 },
    ]);
    graph.insert("c".to_string(), vec![
        Edge { to: "d".to_string(), weight: 1.0 },
    ]);
    graph.insert("d".to_string(), Vec::new());
    let dist = floyd_warshall(&graph).unwrap();
    assert!((dist["a"]["d"] - 4.0).abs() < 1e-9);
    assert!((dist["a"]["a"]).abs() < 1e-9);
    assert!((dist["b"]["d"] - 3.0).abs() < 1e-9);
}

#[test]
fn test_floyd_warshall_negative_cycle() {
    let mut graph: HashMap<String, Vec<Edge>> = HashMap::new();
    graph.insert("a".to_string(), vec![
        Edge { to: "b".to_string(), weight: 1.0 },
    ]);
    graph.insert("b".to_string(), vec![
        Edge { to: "a".to_string(), weight: -2.0 },
    ]);
    assert_eq!(floyd_warshall(&graph), Err(Error::NegativeCycle));
}

#[test]
fn test_strongly_connected_components() {
    let mut graph: HashMap<String, Vec<String>> = HashMap::new();
    graph.insert("a".to_string(), vec!["b".to_string()]);
    graph.insert("b".to_string(), vec!["c".to_string(), "d".to_string()]);
    graph.insert("c".to_string(), vec!["a".to_string()]);
    graph.insert("d".to_string(), vec!["e".to_string()]);
    graph.insert("e".to_string(), Vec::new());
    let components = strongly_connected_components(&graph).unwrap();
    assert_eq!(components.len(), 3);
}

#[test]
fn test_has_cycle() {
    let mut acyclic: HashMap<String, Vec<String>> = HashMap::new();
    acyclic.insert("a".to_string(), vec!["b".to_string()]);
    acyclic.insert("b".to_string(), Vec::new());
    acyclic.insert("c".to_string(), vec!["d".to_string()]);
    acyclic.insert("d".to_string(), Vec::new());
    assert_eq!(has_cycle(&acyclic).unwrap(), false);

    let mut cyclic: HashMap<String, Vec<String>> = HashMap::new();
    cyclic.insert("a".to_string(), vec!["b".to_string()]);
    cyclic.insert("b".to_string(), vec!["c".to_string()]);
    cyclic.insert("c".to_string(), vec!["a".to_string()]);
    assert_eq!(has_cycle(&cyclic).unwrap(), true);
}

#[test]
fn test_binary_search_on_answer() {
    assert_eq!(binary_search_on_answer(1, 10, |x| x >= 6, "minimum").unwrap(), 6);
    assert_eq!(binary_search_on_answer(1, 10, |x| x <= 4, "maximum").unwrap(), 4);
}
