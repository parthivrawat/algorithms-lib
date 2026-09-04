package algorithms

import (
	"errors"
	"math"
	"reflect"
	"testing"
)

func TestQuickSort(t *testing.T) {
	cases := [][]int{
		{},
		{1},
		{3, 1, 4, 1, 5, 9, 2, 6},
		{9, 8, 7, 6, 5, 4, 3, 2, 1},
	}
	for _, c := range cases {
		got := QuickSort(c)
		want := make([]int, len(c))
		copy(want, c)
		sortInts(want)
		if !reflect.DeepEqual(got, want) {
			t.Errorf("QuickSort(%v) = %v, want %v", c, got, want)
		}
	}
	strs := []string{"cherry", "apple", "banana"}
	got := QuickSort(strs)
	want := []string{"apple", "banana", "cherry"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("QuickSort(strings) = %v, want %v", got, want)
	}
}

func TestMergeSort(t *testing.T) {
	in := []int{3, 1, 4, 1, 5, 9, 2, 6}
	got := MergeSort(in)
	want := []int{1, 1, 2, 3, 4, 5, 6, 9}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("MergeSort(%v) = %v, want %v", in, got, want)
	}
}

func TestHeapSort(t *testing.T) {
	in := []int{3, 1, 4, 1, 5, 9, 2, 6}
	got := HeapSort(in)
	want := []int{1, 1, 2, 3, 4, 5, 6, 9}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("HeapSort(%v) = %v, want %v", in, got, want)
	}
}

func TestRadixSort(t *testing.T) {
	in := []int{170, 45, 75, 90, 2, 802, 24, 66}
	got, err := RadixSort(in)
	if err != nil {
		t.Fatalf("RadixSort error: %v", err)
	}
	want := []int{2, 24, 45, 66, 75, 90, 170, 802}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("RadixSort(%v) = %v, want %v", in, got, want)
	}
	if _, err := RadixSort([]int{1, -5, 3}); !errors.Is(err, ErrInvalidInput) {
		t.Errorf("RadixSort negative input error = %v, want ErrInvalidInput", err)
	}
	// exp*10 must not overflow int before the digit loop terminates.
	got, err = RadixSort([]int{math.MaxInt, 1, 0, math.MaxInt - 1})
	if err != nil {
		t.Fatalf("RadixSort large values error: %v", err)
	}
	want = []int{0, 1, math.MaxInt - 1, math.MaxInt}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("RadixSort large values = %v, want %v", got, want)
	}
}

func TestNativeSort(t *testing.T) {
	in := []int{3, 1, 4, 1, 5, 9, 2, 6}
	got := NativeSort(in)
	want := []int{1, 1, 2, 3, 4, 5, 6, 9}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("NativeSort(%v) = %v, want %v", in, got, want)
	}
}

func sortInts(a []int) {
	for i := 0; i < len(a); i++ {
		for j := i + 1; j < len(a); j++ {
			if a[i] > a[j] {
				a[i], a[j] = a[j], a[i]
			}
		}
	}
}

func TestBinarySearch(t *testing.T) {
	arr := []int{2, 4, 6, 8, 10, 12}
	idx, err := BinarySearch(arr, 8)
	if err != nil || idx != 3 {
		t.Errorf("BinarySearch found = %d, err = %v, want 3 nil", idx, err)
	}
	_, err = BinarySearch(arr, 7)
	if !errors.Is(err, ErrNotFound) {
		t.Errorf("BinarySearch missing error = %v, want ErrNotFound", err)
	}
	_, err = BinarySearch([]int{3, 1, 2}, 2)
	if !errors.Is(err, ErrInvalidInput) {
		t.Errorf("BinarySearch unsorted error = %v, want ErrInvalidInput", err)
	}
}

func TestInterpolationSearch(t *testing.T) {
	arr := []int{10, 20, 30, 40, 50, 60, 70, 80, 90, 100}
	idx, err := InterpolationSearch(arr, 50)
	if err != nil || idx != 4 {
		t.Errorf("InterpolationSearch found = %d, err = %v", idx, err)
	}
	_, err = InterpolationSearch(arr, 55)
	if !errors.Is(err, ErrNotFound) {
		t.Errorf("InterpolationSearch missing error = %v, want ErrNotFound", err)
	}

	// Values near math.MaxInt expose float64 precision loss and signed-overflow
	// hazards in the naive (target - arr[lo]) * (hi - lo) / (arr[hi] - arr[lo])
	// computation.
	arr = []int{math.MaxInt - 4, math.MaxInt - 3, math.MaxInt - 2, math.MaxInt - 1, math.MaxInt}
	idx, err = InterpolationSearch(arr, math.MaxInt-2)
	if err != nil || idx != 2 {
		t.Errorf("InterpolationSearch near MaxInt = %d, err = %v", idx, err)
	}

	// Extreme negative-to-positive range would overflow signed int subtraction.
	arr = []int{math.MinInt, -1, 0, 1, math.MaxInt}
	idx, err = InterpolationSearch(arr, math.MaxInt)
	if err != nil || idx != 4 {
		t.Errorf("InterpolationSearch extreme range = %d, err = %v", idx, err)
	}
}

func TestJumpSearch(t *testing.T) {
	arr := []int{2, 4, 6, 8, 10, 12}
	idx, err := JumpSearch(arr, 10)
	if err != nil || idx != 4 {
		t.Errorf("JumpSearch found = %d, err = %v", idx, err)
	}
	_, err = JumpSearch(arr, 7)
	if !errors.Is(err, ErrNotFound) {
		t.Errorf("JumpSearch missing error = %v, want ErrNotFound", err)
	}
	// Empty and nil slices must not panic (arr[min(step,n)-1] = arr[-1]).
	for _, empty := range [][]int{nil, {}} {
		if _, err := JumpSearch(empty, 1); !errors.Is(err, ErrNotFound) {
			t.Errorf("JumpSearch empty slice error = %v, want ErrNotFound", err)
		}
	}
}

func TestBFS(t *testing.T) {
	graph := map[string][]string{
		"a": {"b", "c"},
		"b": {"d"},
		"c": {},
		"d": {},
	}
	got, err := BFS(graph, "a")
	if err != nil {
		t.Fatalf("BFS error: %v", err)
	}
	want := []string{"a", "b", "c", "d"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("BFS = %v, want %v", got, want)
	}
	_, err = BFS(graph, "z")
	if !errors.Is(err, ErrInvalidGraph) {
		t.Errorf("BFS missing start error = %v, want ErrInvalidGraph", err)
	}

	// Neighbors that are not keys of the graph must be rejected.
	bad := map[string][]string{"a": {"b", "zz"}, "b": {}}
	if _, err := BFS(bad, "a"); !errors.Is(err, ErrInvalidGraph) {
		t.Errorf("BFS missing neighbor error = %v, want ErrInvalidGraph", err)
	}
	if _, err := DFS(bad, "a"); !errors.Is(err, ErrInvalidGraph) {
		t.Errorf("DFS missing neighbor error = %v, want ErrInvalidGraph", err)
	}
}

func TestDFS(t *testing.T) {
	graph := map[string][]string{
		"a": {"b", "c"},
		"b": {"d"},
		"c": {},
		"d": {},
	}
	got, err := DFS(graph, "a")
	if err != nil {
		t.Fatalf("DFS error: %v", err)
	}
	want := []string{"a", "b", "d", "c"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("DFS = %v, want %v", got, want)
	}
}

func TestDijkstra(t *testing.T) {
	graph := map[string][]Edge{
		"a": {{To: "b", Weight: 1}, {To: "c", Weight: 4}},
		"b": {{To: "c", Weight: 2}, {To: "d", Weight: 5}},
		"c": {{To: "d", Weight: 1}},
		"d": {},
	}
	distances, _, err := Dijkstra(graph, "a")
	if err != nil {
		t.Fatalf("Dijkstra error: %v", err)
	}
	if distances["d"] != 4 {
		t.Errorf("Dijkstra distance to d = %v, want 4", distances["d"])
	}
	negative := map[string][]Edge{
		"a": {{To: "b", Weight: -1}},
		"b": {},
	}
	_, _, err = Dijkstra(negative, "a")
	if !errors.Is(err, ErrInvalidGraph) {
		t.Errorf("Dijkstra negative error = %v, want ErrInvalidGraph", err)
	}

	unreachable := map[string][]Edge{
		"a": {},
		"b": {{To: "c", Weight: -1}},
		"c": {},
	}
	_, _, err = Dijkstra(unreachable, "a")
	if !errors.Is(err, ErrInvalidGraph) {
		t.Errorf("Dijkstra unreachable negative error = %v, want ErrInvalidGraph", err)
	}
}

func TestAStar(t *testing.T) {
	graph := map[string][]Edge{
		"a": {{To: "b", Weight: 1}, {To: "c", Weight: 4}},
		"b": {{To: "c", Weight: 2}, {To: "d", Weight: 5}},
		"c": {{To: "d", Weight: 1}},
		"d": {},
	}
	path, cost, err := AStar(graph, "a", "d", func(_, _ string) float64 { return 0 })
	if err != nil {
		t.Fatalf("AStar error: %v", err)
	}
	want := []string{"a", "b", "c", "d"}
	if !reflect.DeepEqual(path, want) || cost != 4 {
		t.Errorf("AStar = (%v, %v), want (%v, 4)", path, cost, want)
	}
	_, _, err = AStar(map[string][]Edge{"a": {}, "b": {}}, "a", "b", func(_, _ string) float64 { return 0 })
	if !errors.Is(err, ErrNotFound) {
		t.Errorf("AStar no path error = %v, want ErrNotFound", err)
	}

	// Negative edge weights are not supported by A*.
	_, _, err = AStar(map[string][]Edge{"a": {{To: "b", Weight: -1}}, "b": {}}, "a", "b", func(_, _ string) float64 { return 0 })
	if !errors.Is(err, ErrInvalidGraph) {
		t.Errorf("AStar negative weight error = %v, want ErrInvalidGraph", err)
	}

	// Negative heuristic values are rejected immediately.
	_, _, err = AStar(map[string][]Edge{"a": {}, "b": {}}, "a", "b", func(_, _ string) float64 { return -1 })
	if !errors.Is(err, ErrInvalidInput) {
		t.Errorf("AStar negative heuristic error = %v, want ErrInvalidInput", err)
	}
}

func TestBellmanFord(t *testing.T) {
	graph := map[string][]Edge{
		"a": {{To: "b", Weight: -1}},
		"b": {{To: "c", Weight: -2}},
		"c": {{To: "d", Weight: 1}},
		"d": {},
	}
	distances, _, err := BellmanFord(graph, "a")
	if err != nil {
		t.Fatalf("BellmanFord error: %v", err)
	}
	if distances["d"] != -2 {
		t.Errorf("BellmanFord distance to d = %v, want -2", distances["d"])
	}
	negative := map[string][]Edge{
		"a": {{To: "b", Weight: 1}},
		"b": {{To: "c", Weight: -1}},
		"c": {{To: "b", Weight: -1}},
	}
	_, _, err = BellmanFord(negative, "a")
	if !errors.Is(err, ErrNegativeCycle) {
		t.Errorf("BellmanFord negative cycle error = %v, want ErrNegativeCycle", err)
	}
}

func TestKnapsack01(t *testing.T) {
	got, err := Knapsack01([]int{1, 2, 3}, []int{6, 10, 12}, 5)
	if err != nil || got != 22 {
		t.Errorf("Knapsack01 = %d, err = %v, want 22", got, err)
	}
}

func TestLongestCommonSubsequence(t *testing.T) {
	got := LongestCommonSubsequence([]byte("ABCDE"), []byte("ACE"))
	if got != 3 {
		t.Errorf("LCS = %d, want 3", got)
	}
}

func TestEditDistance(t *testing.T) {
	got := EditDistance([]byte("kitten"), []byte("sitting"))
	if got != 3 {
		t.Errorf("EditDistance = %d, want 3", got)
	}
}

func TestEditDistanceReconstruction(t *testing.T) {
	dist, script, err := EditDistanceReconstruction([]byte("kitten"), []byte("sitting"))
	if err != nil {
		t.Fatalf("EditDistanceReconstruction error: %v", err)
	}
	if dist != 3 {
		t.Errorf("dist = %d, want 3", dist)
	}
	if len(script) < 3 {
		t.Errorf("len(script) = %d, want >= 3", len(script))
	}
}

func TestLongestCommonSubsequenceReconstruction(t *testing.T) {
	length, lcs, err := LongestCommonSubsequenceReconstruction([]byte("ABCDE"), []byte("ACE"))
	if err != nil {
		t.Fatalf("LongestCommonSubsequenceReconstruction error: %v", err)
	}
	if length != 3 {
		t.Errorf("length = %d, want 3", length)
	}
	if string(lcs) != "ACE" {
		t.Errorf("lcs = %s, want ACE", string(lcs))
	}
}

func TestKMPSearch(t *testing.T) {
	text := "ABABDABACDABABCABAB"
	pattern := "ABABCABAB"
	got, err := KMPSearch(text, pattern)
	if err != nil {
		t.Fatalf("KMPSearch error: %v", err)
	}
	want := []int{10}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("KMPSearch = %v, want %v", got, want)
	}
}

func TestRabinKarpSearch(t *testing.T) {
	text := "ABABDABACDABABCABAB"
	pattern := "ABABCABAB"
	got, err := RabinKarpSearch(text, pattern, 256, 1_000_000_007)
	if err != nil {
		t.Fatalf("RabinKarpSearch error: %v", err)
	}
	want := []int{10}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("RabinKarpSearch = %v, want %v", got, want)
	}

	if _, err := RabinKarpSearch(text, pattern, 256, 0); !errors.Is(err, ErrInvalidInput) {
		t.Errorf("RabinKarpSearch mod=0 error = %v, want ErrInvalidInput", err)
	}
	if _, err := RabinKarpSearch(text, pattern, -1, 1_000_000_007); !errors.Is(err, ErrInvalidInput) {
		t.Errorf("RabinKarpSearch base=-1 error = %v, want ErrInvalidInput", err)
	}

	// Overlapping matches and a large modulus near int64 max.
	overlap, err := RabinKarpSearch("aaaaa", "aaaa", 256, 1_000_000_007)
	if err != nil {
		t.Fatalf("RabinKarpSearch overlap error: %v", err)
	}
	if !reflect.DeepEqual(overlap, []int{0, 1}) {
		t.Errorf("RabinKarpSearch overlap = %v, want [0 1]", overlap)
	}
}

func TestBoyerMooreSearch(t *testing.T) {
	text := "ABABDABACDABABCABAB"
	pattern := "ABABCABAB"
	got, err := BoyerMooreSearch(text, pattern)
	if err != nil {
		t.Fatalf("BoyerMooreSearch error: %v", err)
	}
	want := []int{10}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("BoyerMooreSearch = %v, want %v", got, want)
	}
}

func TestBoyerMooreSearchOverlapping(t *testing.T) {
	cases := []struct {
		text, pattern string
		want          []int
	}{
		{"aaaa", "aa", []int{0, 1, 2}},
		{"aaaaa", "aaaa", []int{0, 1}},
		{"abcabcabc", "abcabc", []int{0, 3}},
		{"abababab", "abab", []int{0, 2, 4}},
	}
	for _, tc := range cases {
		got, err := BoyerMooreSearch(tc.text, tc.pattern)
		if err != nil {
			t.Fatalf("BoyerMooreSearch(%q, %q) error: %v", tc.text, tc.pattern, err)
		}
		if !reflect.DeepEqual(got, tc.want) {
			t.Errorf("BoyerMooreSearch(%q, %q) = %v, want %v", tc.text, tc.pattern, got, tc.want)
		}
	}
}

func TestBoyerMooreSearchByteOffsets(t *testing.T) {
	// Multi-byte characters: indices must be byte offsets, matching KMPSearch.
	got, err := BoyerMooreSearch("héllo héllo", "éllo")
	if err != nil {
		t.Fatalf("BoyerMooreSearch error: %v", err)
	}
	want, err := KMPSearch("héllo héllo", "éllo")
	if err != nil {
		t.Fatalf("KMPSearch error: %v", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("BoyerMooreSearch byte offsets = %v, want %v", got, want)
	}
}

func TestActivitySelection(t *testing.T) {
	activities := [][2]int{
		{1, 3}, {2, 5}, {4, 7}, {1, 8}, {5, 9},
		{8, 10}, {9, 11}, {11, 14}, {13, 16},
	}
	got, err := ActivitySelection(activities)
	if err != nil {
		t.Fatalf("ActivitySelection error: %v", err)
	}
	if len(got) < 4 {
		t.Errorf("ActivitySelection length = %d, want >= 4", len(got))
	}
}

func TestFractionalKnapsack(t *testing.T) {
	got, err := FractionalKnapsack([]float64{10, 20, 30}, []float64{60, 100, 120}, 50)
	if err != nil || math.Abs(got-240) > 1e-9 {
		t.Errorf("FractionalKnapsack = %v, err = %v, want 240", got, err)
	}
}

func TestFractionalKnapsackEdgeCases(t *testing.T) {
	// Zero-weight, positive-value items are taken in full for free.
	if got, err := FractionalKnapsack([]float64{0, 10}, []float64{50, 60}, 5); err != nil || math.Abs(got-80) > 1e-9 {
		t.Errorf("FractionalKnapsack zero-weight = %v, err = %v, want 80", got, err)
	}
	// Zero-weight items with non-positive value are never taken.
	if got, err := FractionalKnapsack([]float64{0, 10}, []float64{-50, 60}, 50); err != nil || got != 60 {
		t.Errorf("FractionalKnapsack zero-weight negative value = %v, err = %v, want 60", got, err)
	}
	// Negative weights are rejected.
	if _, err := FractionalKnapsack([]float64{-5, 10}, []float64{10, 60}, 50); !errors.Is(err, ErrInvalidInput) {
		t.Errorf("FractionalKnapsack negative weight error = %v, want ErrInvalidInput", err)
	}
	// Negative-value items are never taken, so the total never goes below zero.
	if got, err := FractionalKnapsack([]float64{10}, []float64{-5}, 10); err != nil || got != 0 {
		t.Errorf("FractionalKnapsack negative value = %v, err = %v, want 0", got, err)
	}
}

func TestHuffmanCoding(t *testing.T) {
	frequencies := map[string]int{
		"a": 5, "b": 9, "c": 12, "d": 13, "e": 16, "f": 45,
	}
	codes, err := HuffmanCoding(frequencies)
	if err != nil {
		t.Fatalf("HuffmanCoding error: %v", err)
	}
	if len(codes) != 6 {
		t.Errorf("HuffmanCoding codes = %v, want 6", codes)
	}
	for k, v := range codes {
		if v == "" {
			t.Errorf("HuffmanCoding empty code for %s", k)
		}
	}
}

func TestHuffmanCodingEdgeCases(t *testing.T) {
	// The empty string is a valid symbol key and must appear in the code table.
	codes, err := HuffmanCoding(map[string]int{"": 2, "a": 3, "b": 1})
	if err != nil {
		t.Fatalf("HuffmanCoding empty-key error: %v", err)
	}
	if _, ok := codes[""]; !ok || len(codes) != 3 {
		t.Errorf("HuffmanCoding codes = %v, want 3 entries including \"\"", codes)
	}

	// A single symbol gets code "0".
	single, err := HuffmanCoding(map[string]int{"x": 3})
	if err != nil || single["x"] != "0" {
		t.Errorf("HuffmanCoding single symbol = %v, %v, want x->0", single, err)
	}

	// Zero and negative frequencies are rejected.
	for _, freq := range []int{0, -1} {
		if _, err := HuffmanCoding(map[string]int{"a": freq}); !errors.Is(err, ErrInvalidInput) {
			t.Errorf("HuffmanCoding freq=%d error = %v, want ErrInvalidInput", freq, err)
		}
	}
}

func TestHuffmanEncodeDecode(t *testing.T) {
	frequencies := map[string]int{"a": 5, "b": 9, "c": 12, "d": 13, "e": 16, "f": 45}
	codes, err := HuffmanCoding(frequencies)
	if err != nil {
		t.Fatalf("HuffmanCoding error: %v", err)
	}
	symbols := []string{"a", "b", "c", "d", "e", "f"}
	encoded, err := HuffmanEncode(symbols, codes)
	if err != nil {
		t.Fatalf("HuffmanEncode error: %v", err)
	}
	decoded, err := HuffmanDecode(encoded, codes)
	if err != nil {
		t.Fatalf("HuffmanDecode error: %v", err)
	}
	if len(decoded) != len(symbols) {
		t.Errorf("decoded length = %d, want %d", len(decoded), len(symbols))
	}
	for i := range symbols {
		if decoded[i] != symbols[i] {
			t.Errorf("decoded[%d] = %s, want %s", i, decoded[i], symbols[i])
		}
	}
}

func TestMaxSubarray(t *testing.T) {
	arr := []int{-2, 1, -3, 4, -1, 2, 1, -5, 4}
	got, err := MaxSubarray(arr)
	if err != nil || got != 6 {
		t.Errorf("MaxSubarray = %d, err = %v, want 6", got, err)
	}
	got, _ = MaxSubarray([]int{-2, -1})
	if got != -1 {
		t.Errorf("MaxSubarray all negative = %d, want -1", got)
	}
}

func TestCountInversions(t *testing.T) {
	got := CountInversions([]int{1, 3, 5, 2, 4, 6})
	if got != 3 {
		t.Errorf("CountInversions = %d, want 3", got)
	}
	// A fully reversed slice of n elements has n*(n-1)/2 inversions.
	const n = 1000
	rev := make([]int, n)
	for i := range rev {
		rev[i] = n - i
	}
	if got := CountInversions(rev); got != int64(n)*(n-1)/2 {
		t.Errorf("CountInversions(reversed %d) = %d, want %d", n, got, int64(n)*(n-1)/2)
	}
}

func TestFastPower(t *testing.T) {
	if got, err := FastPower(2, 10); err != nil || got != 1024 {
		t.Errorf("FastPower(2,10) = %v, %v; want 1024, nil", got, err)
	}
	if got, err := FastPower(5, 0); err != nil || got != 1 {
		t.Errorf("FastPower(5,0) = %v, %v; want 1, nil", got, err)
	}
	if got, err := FastPower(2, -2); err != nil || got != 0.25 {
		t.Errorf("FastPower(2,-2) = %v, %v; want 0.25, nil", got, err)
	}
}

func TestFastPowerEdgeCases(t *testing.T) {
	if _, err := FastPower(0, -1); !errors.Is(err, ErrInvalidInput) {
		t.Errorf("FastPower(0,-1) error = %v, want ErrInvalidInput", err)
	}
	if got, err := FastPower(0, 0); err != nil || got != 1 {
		t.Errorf("FastPower(0,0) = %v, %v; want 1, nil", got, err)
	}
	// math.MinInt cannot be negated in int; 2^MinInt must not hang or crash.
	if got, err := FastPower(2, math.MinInt); err != nil || got != 0 {
		t.Errorf("FastPower(2,MinInt) = %v, %v; want 0, nil", got, err)
	}
	if got, err := FastPower(-2, 3); err != nil || got != -8 {
		t.Errorf("FastPower(-2,3) = %v, %v; want -8, nil", got, err)
	}
	if got, err := FastPower(-2, 4); err != nil || got != 16 {
		t.Errorf("FastPower(-2,4) = %v, %v; want 16, nil", got, err)
	}
}

func TestGcd(t *testing.T) {
	if got, _ := Gcd(48, 18); got != 6 {
		t.Errorf("Gcd(48,18) = %v, want 6", got)
	}
	if got, _ := Gcd(-48, 18); got != 6 {
		t.Errorf("Gcd(-48,18) = %v, want 6", got)
	}
	if got, _ := Gcd(0, 5); got != 5 {
		t.Errorf("Gcd(0,5) = %v, want 5", got)
	}
	if got, _ := Gcd(0, 0); got != 0 {
		t.Errorf("Gcd(0,0) = %v, want 0", got)
	}
}

func TestLcm(t *testing.T) {
	if got, _ := Lcm(4, 6); got != 12 {
		t.Errorf("Lcm(4,6) = %v, want 12", got)
	}
	if got, _ := Lcm(-4, 6); got != 12 {
		t.Errorf("Lcm(-4,6) = %v, want 12", got)
	}
	if got, _ := Lcm(0, 5); got != 0 {
		t.Errorf("Lcm(0,5) = %v, want 0", got)
	}
	if _, err := Lcm(0, 0); !errors.Is(err, ErrInvalidInput) {
		t.Errorf("Lcm(0,0) error = %v, want ErrInvalidInput", err)
	}
}

func TestModularFastPower(t *testing.T) {
	if got, err := ModularFastPower(2, 10, 1000); err != nil || got != 24 {
		t.Errorf("ModularFastPower(2,10,1000) = %v, %v; want 24, nil", got, err)
	}
	if got, err := ModularFastPower(3, 0, 7); err != nil || got != 1 {
		t.Errorf("ModularFastPower(3,0,7) = %v, %v; want 1, nil", got, err)
	}
	if got, err := ModularFastPower(0, 0, 7); err != nil || got != 1 {
		t.Errorf("ModularFastPower(0,0,7) = %v, %v; want 1, nil", got, err)
	}
	if got, err := ModularFastPower(0, 5, 7); err != nil || got != 0 {
		t.Errorf("ModularFastPower(0,5,7) = %v, %v; want 0, nil", got, err)
	}
	if got, err := ModularFastPower(-2, 3, 1000); err != nil || got != 992 {
		t.Errorf("ModularFastPower(-2,3,1000) = %v, %v; want 992, nil", got, err)
	}
	if _, err := ModularFastPower(2, -1, 7); !errors.Is(err, ErrInvalidInput) {
		t.Errorf("ModularFastPower(2,-1,7) error = %v, want ErrInvalidInput", err)
	}
	if _, err := ModularFastPower(2, 2, 0); !errors.Is(err, ErrInvalidInput) {
		t.Errorf("ModularFastPower(2,2,0) error = %v, want ErrInvalidInput", err)
	}
}

func TestQuickSelect(t *testing.T) {
	if got, _ := QuickSelect([]int{3, 2, 1, 5, 4}, 0); got != 1 {
		t.Errorf("QuickSelect(...,0) = %v, want 1", got)
	}
	if got, _ := QuickSelect([]int{3, 2, 1, 5, 4}, 2); got != 3 {
		t.Errorf("QuickSelect(...,2) = %v, want 3", got)
	}
	if got, _ := QuickSelect([]int{3, 2, 1, 5, 4}, 4); got != 5 {
		t.Errorf("QuickSelect(...,4) = %v, want 5", got)
	}
	strs := []string{"cherry", "apple", "banana"}
	if got, _ := QuickSelect(strs, 1); got != "banana" {
		t.Errorf("QuickSelect(strings,1) = %v, want banana", got)
	}
	if _, err := QuickSelect([]int{1, 2, 3}, 5); !errors.Is(err, ErrInvalidInput) {
		t.Errorf("QuickSelect(...,5) error = %v, want ErrInvalidInput", err)
	}
}

func TestSieveOfEratosthenes(t *testing.T) {
	got, _ := SieveOfEratosthenes(10)
	want := []int{2, 3, 5, 7}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("SieveOfEratosthenes(10) = %v, want %v", got, want)
	}
	if got, _ := SieveOfEratosthenes(1); len(got) != 0 {
		t.Errorf("SieveOfEratosthenes(1) = %v, want []", got)
	}
	if _, err := SieveOfEratosthenes(-1); !errors.Is(err, ErrInvalidInput) {
		t.Errorf("SieveOfEratosthenes(-1) error = %v, want ErrInvalidInput", err)
	}
}

func TestTopologicalSort(t *testing.T) {
	graph := map[string][]string{
		"a": {"b", "c"},
		"b": {"d"},
		"c": {"d"},
		"d": {},
	}
	order, err := TopologicalSort(graph)
	if err != nil {
		t.Fatalf("TopologicalSort error = %v", err)
	}
	idx := make(map[string]int)
	for i, v := range order {
		idx[v] = i
	}
	if !(idx["a"] < idx["b"] && idx["a"] < idx["c"] &&
		idx["b"] < idx["d"] && idx["c"] < idx["d"]) {
		t.Errorf("TopologicalSort order %v is not valid", order)
	}

	cyclic := map[string][]string{
		"a": {"b"},
		"b": {"c"},
		"c": {"a"},
	}
	if _, err := TopologicalSort(cyclic); !errors.Is(err, ErrInvalidGraph) {
		t.Errorf("TopologicalSort(cyclic) error = %v, want ErrInvalidGraph", err)
	}
}

func TestUnionFind(t *testing.T) {
	uf := NewUnionFind([]string{"a", "b", "c", "d", "e"})
	if connected, _ := uf.Connected("a", "b"); connected {
		t.Errorf("a and b should not be connected")
	}
	if merged, _ := uf.Union("a", "b"); !merged {
		t.Errorf("Union(a,b) should merge")
	}
	if connected, _ := uf.Connected("a", "b"); !connected {
		t.Errorf("a and b should be connected")
	}
	uf.Union("b", "c")
	if connected, _ := uf.Connected("a", "c"); !connected {
		t.Errorf("a and c should be connected")
	}
	if connected, _ := uf.Connected("a", "d"); connected {
		t.Errorf("a and d should not be connected")
	}
	if uf.Count() != 3 {
		t.Errorf("Count() = %d, want 3", uf.Count())
	}
	if _, err := uf.Find("z"); !errors.Is(err, ErrInvalidInput) {
		t.Errorf("Find(z) error = %v, want ErrInvalidInput", err)
	}
}

func TestPriorityQueue(t *testing.T) {
	pq := NewPriorityQueue([]int{3, 1, 4, 1, 5})
	if got, _ := pq.Peek(); got != 1 {
		t.Errorf("Peek() = %v, want 1", got)
	}
	if pq.Len() != 5 {
		t.Errorf("Len() = %d, want 5", pq.Len())
	}
	if got, _ := pq.Pop(); got != 1 {
		t.Errorf("Pop() = %v, want 1", got)
	}
	if got, _ := pq.Pop(); got != 1 {
		t.Errorf("Pop() = %v, want 1", got)
	}
	if got, _ := pq.Pop(); got != 3 {
		t.Errorf("Pop() = %v, want 3", got)
	}
	pq.Push(0)
	if got, _ := pq.Pop(); got != 0 {
		t.Errorf("Pop() after Push(0) = %v, want 0", got)
	}
	empty := NewPriorityQueue([]int{})
	if _, err := empty.Pop(); !errors.Is(err, ErrInvalidInput) {
		t.Errorf("Pop() on empty error = %v, want ErrInvalidInput", err)
	}
}

func TestMinimumSpanningTree(t *testing.T) {
	edges := []MSTEdge{
		{U: "a", V: "b", Weight: 1},
		{U: "b", V: "c", Weight: 2},
		{U: "a", V: "c", Weight: 3},
		{U: "c", V: "d", Weight: 4},
	}
	total, mst, err := MinimumSpanningTree(edges)
	if err != nil {
		t.Fatalf("MinimumSpanningTree error = %v", err)
	}
	if total != 7 {
		t.Errorf("total = %v, want 7", total)
	}
	if len(mst) != 3 {
		t.Errorf("len(mst) = %d, want 3", len(mst))
	}
}

func TestFloydWarshall(t *testing.T) {
	graph := map[string][]Edge{
		"a": {{To: "b", Weight: 1}, {To: "c", Weight: 4}},
		"b": {{To: "c", Weight: 2}, {To: "d", Weight: 5}},
		"c": {{To: "d", Weight: 1}},
		"d": {},
	}
	dist, err := FloydWarshall(graph)
	if err != nil {
		t.Fatalf("FloydWarshall error = %v", err)
	}
	if dist["a"]["d"] != 4 {
		t.Errorf("a->d = %v, want 4", dist["a"]["d"])
	}
	if dist["a"]["a"] != 0 {
		t.Errorf("a->a = %v, want 0", dist["a"]["a"])
	}
	if dist["b"]["d"] != 3 {
		t.Errorf("b->d = %v, want 3", dist["b"]["d"])
	}
}

func TestFloydWarshallNegativeCycle(t *testing.T) {
	graph := map[string][]Edge{
		"a": {{To: "b", Weight: 1}},
		"b": {{To: "a", Weight: -2}},
	}
	if _, err := FloydWarshall(graph); !errors.Is(err, ErrNegativeCycle) {
		t.Errorf("FloydWarshall error = %v, want ErrNegativeCycle", err)
	}
}

func TestStronglyConnectedComponents(t *testing.T) {
	graph := map[string][]string{
		"a": {"b"},
		"b": {"c", "d"},
		"c": {"a"},
		"d": {"e"},
		"e": {},
	}
	components, err := StronglyConnectedComponents(graph)
	if err != nil {
		t.Fatalf("StronglyConnectedComponents error = %v", err)
	}
	if len(components) != 3 {
		t.Errorf("len(components) = %d, want 3", len(components))
	}
}

func TestHasCycle(t *testing.T) {
	acyclic := map[string][]string{"a": {"b"}, "b": {}, "c": {"d"}, "d": {}}
	if got, _ := HasCycle(acyclic); got {
		t.Errorf("HasCycle(acyclic) = %v, want false", got)
	}
	cyclic := map[string][]string{"a": {"b"}, "b": {"c"}, "c": {"a"}}
	if got, _ := HasCycle(cyclic); !got {
		t.Errorf("HasCycle(cyclic) = %v, want true", got)
	}
}

func TestBinarySearchOnAnswer(t *testing.T) {
	got, err := BinarySearchOnAnswer(1, 10, func(x int) bool { return x >= 6 }, "minimum")
	if err != nil || got != 6 {
		t.Errorf("BinarySearchOnAnswer = %d, err = %v, want 6", got, err)
	}
	got, err = BinarySearchOnAnswer(1, 10, func(x int) bool { return x <= 4 }, "maximum")
	if err != nil || got != 4 {
		t.Errorf("BinarySearchOnAnswer = %d, err = %v, want 4", got, err)
	}
}
