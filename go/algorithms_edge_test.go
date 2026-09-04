package algorithms

import (
	"errors"
	"math"
	"reflect"
	"testing"
)

// TestEdgeEmptyAndSingleInputs covers empty and single-element inputs across
// sorts, searches, and several DP / divide-and-conquer algorithms.
func TestEdgeEmptyAndSingleInputs(t *testing.T) {
	sorters := []struct {
		name string
		f    func([]int) []int
	}{
		{"QuickSort", QuickSort[int]},
		{"MergeSort", MergeSort[int]},
		{"HeapSort", HeapSort[int]},
		{"NativeSort", NativeSort[int]},
	}
	for _, s := range sorters {
		if got := s.f([]int{}); !reflect.DeepEqual(got, []int{}) {
			t.Errorf("%s(empty) = %v, want []", s.name, got)
		}
		if got := s.f([]int{42}); !reflect.DeepEqual(got, []int{42}) {
			t.Errorf("%s([42]) = %v, want [42]", s.name, got)
		}
	}

	if got, err := RadixSort([]int{}); err != nil || !reflect.DeepEqual(got, []int{}) {
		t.Errorf("RadixSort(empty) = %v, %v, want [], nil", got, err)
	}
	if got, err := RadixSort([]int{42}); err != nil || !reflect.DeepEqual(got, []int{42}) {
		t.Errorf("RadixSort([42]) = %v, %v, want [42], nil", got, err)
	}

	if _, err := BinarySearch([]int{}, 5); !errors.Is(err, ErrNotFound) {
		t.Errorf("BinarySearch(empty) error = %v, want ErrNotFound", err)
	}
	if idx, err := BinarySearch([]int{5}, 5); err != nil || idx != 0 {
		t.Errorf("BinarySearch([5],5) = %d, %v, want 0, nil", idx, err)
	}

	if _, err := JumpSearch([]int{}, 5); !errors.Is(err, ErrNotFound) {
		t.Errorf("JumpSearch(empty) error = %v, want ErrNotFound", err)
	}
	if idx, err := JumpSearch([]int{5}, 5); err != nil || idx != 0 {
		t.Errorf("JumpSearch([5],5) = %d, %v, want 0, nil", idx, err)
	}

	if _, err := InterpolationSearch([]int{}, 5); !errors.Is(err, ErrNotFound) {
		t.Errorf("InterpolationSearch(empty) error = %v, want ErrNotFound", err)
	}
	if idx, err := InterpolationSearch([]int{5}, 5); err != nil || idx != 0 {
		t.Errorf("InterpolationSearch([5],5) = %d, %v, want 0, nil", idx, err)
	}

	if got := LongestCommonSubsequence([]byte{}, []byte("abc")); got != 0 {
		t.Errorf("LCS(empty, abc) = %d, want 0", got)
	}
	if got := EditDistance([]byte{}, []byte("abc")); got != 3 {
		t.Errorf("EditDistance(empty, abc) = %d, want 3", got)
	}
	if got := LongestCommonSubsequence([]byte("x"), []byte("x")); got != 1 {
		t.Errorf("LCS(x,x) = %d, want 1", got)
	}

	if got, err := MaxSubarray([]int{7}); err != nil || got != 7 {
		t.Errorf("MaxSubarray([7]) = %d, %v, want 7, nil", got, err)
	}
	if _, err := MaxSubarray([]int{}); !errors.Is(err, ErrInvalidInput) {
		t.Errorf("MaxSubarray(empty) error = %v, want ErrInvalidInput", err)
	}

	if got := CountInversions([]int{}); got != 0 {
		t.Errorf("CountInversions(empty) = %d, want 0", got)
	}
	if got := CountInversions([]int{7}); got != 0 {
		t.Errorf("CountInversions([7]) = %d, want 0", got)
	}
}

// TestEdgeDuplicates covers inputs containing duplicate values.
func TestEdgeDuplicates(t *testing.T) {
	in := []int{3, 1, 4, 1, 5, 9, 2, 6, 5}
	want := []int{1, 1, 2, 3, 4, 5, 5, 6, 9}

	sorters := []struct {
		name string
		f    func([]int) []int
	}{
		{"QuickSort", QuickSort[int]},
		{"MergeSort", MergeSort[int]},
		{"HeapSort", HeapSort[int]},
		{"NativeSort", NativeSort[int]},
	}
	for _, s := range sorters {
		if got := s.f(in); !reflect.DeepEqual(got, want) {
			t.Errorf("%s(%v) = %v, want %v", s.name, in, got, want)
		}
	}

	if got, err := RadixSort(in); err != nil || !reflect.DeepEqual(got, want) {
		t.Errorf("RadixSort(%v) = %v, %v, want %v, nil", in, got, err, want)
	}

	if got := CountInversions([]int{2, 2, 2}); got != 0 {
		t.Errorf("CountInversions(all equal) = %d, want 0", got)
	}
	if got := CountInversions([]int{3, 1, 3, 1}); got != 3 {
		t.Errorf("CountInversions([3,1,3,1]) = %d, want 3", got)
	}
}

// TestEdgeOverlappingStringMatches checks overlapping occurrences for all three
// exact string-search implementations.
func TestEdgeOverlappingStringMatches(t *testing.T) {
	cases := []struct {
		text, pattern string
		want          []int
	}{
		{"aaaaa", "aaaa", []int{0, 1}},
		{"aaaa", "aa", []int{0, 1, 2}},
		{"abcabcabc", "abcabc", []int{0, 3}},
	}
	searchers := []struct {
		name string
		f    func(string, string) ([]int, error)
	}{
		{"KMPSearch", KMPSearch},
		{"RabinKarpSearch", func(text, pattern string) ([]int, error) {
			return RabinKarpSearch(text, pattern, 256, 1_000_000_007)
		}},
		{"BoyerMooreSearch", BoyerMooreSearch},
	}
	for _, s := range searchers {
		for _, c := range cases {
			got, err := s.f(c.text, c.pattern)
			if err != nil {
				t.Fatalf("%s(%q,%q) error: %v", s.name, c.text, c.pattern, err)
			}
			if !reflect.DeepEqual(got, c.want) {
				t.Errorf("%s(%q,%q) = %v, want %v", s.name, c.text, c.pattern, got, c.want)
			}
		}
	}
}

// TestEdgeDisconnectedGraphs verifies behavior on disconnected unweighted and
// weighted graphs.
func TestEdgeDisconnectedGraphs(t *testing.T) {
	unweighted := map[string][]string{
		"a": {"c"},
		"c": {},
		"b": {"d"},
		"d": {},
	}
	if got, err := BFS(unweighted, "a"); err != nil {
		t.Fatalf("BFS disconnected error: %v", err)
	} else if want := []string{"a", "c"}; !reflect.DeepEqual(got, want) {
		t.Errorf("BFS disconnected = %v, want %v", got, want)
	}
	if got, err := DFS(unweighted, "a"); err != nil {
		t.Fatalf("DFS disconnected error: %v", err)
	} else if want := []string{"a", "c"}; !reflect.DeepEqual(got, want) {
		t.Errorf("DFS disconnected = %v, want %v", got, want)
	}

	weighted := map[string][]Edge{
		"a": {},
		"b": {{To: "c", Weight: 1}},
		"c": {},
	}
	distances, _, err := Dijkstra(weighted, "a")
	if err != nil {
		t.Fatalf("Dijkstra disconnected error: %v", err)
	}
	if distances["a"] != 0 {
		t.Errorf("Dijkstra start distance = %v, want 0", distances["a"])
	}
	if !math.IsInf(distances["b"], 1) || !math.IsInf(distances["c"], 1) {
		t.Errorf("Dijkstra unreachable distances = %v, want +Inf", distances)
	}

	if _, _, err := AStar(weighted, "a", "b", func(_, _ string) float64 { return 0 }); !errors.Is(err, ErrNotFound) {
		t.Errorf("AStar disconnected error = %v, want ErrNotFound", err)
	}

	distBF, _, err := BellmanFord(weighted, "a")
	if err != nil {
		t.Fatalf("BellmanFord disconnected error: %v", err)
	}
	if len(distBF) != 1 || distBF["a"] != 0 {
		t.Errorf("BellmanFord disconnected = %v, want map[a:0]", distBF)
	}
}

// TestEdgeStartEqualsGoal covers the start == goal case for A* and Dijkstra.
func TestEdgeStartEqualsGoal(t *testing.T) {
	graph := map[string][]Edge{
		"a": {{To: "b", Weight: 2}},
		"b": {},
	}
	path, cost, err := AStar(graph, "a", "a", func(_, _ string) float64 { return 0 })
	if err != nil {
		t.Fatalf("AStar start==goal error: %v", err)
	}
	if !reflect.DeepEqual(path, []string{"a"}) || cost != 0 {
		t.Errorf("AStar start==goal = (%v, %v), want ([a], 0)", path, cost)
	}

	distances, _, err := Dijkstra(graph, "a")
	if err != nil {
		t.Fatalf("Dijkstra start==goal error: %v", err)
	}
	if distances["a"] != 0 {
		t.Errorf("Dijkstra start==goal distance = %v, want 0", distances["a"])
	}
}

// TestEdgeZeroCapacityKnapsack checks the 0/1 and fractional knapsacks when the
// backpack has no usable capacity.
func TestEdgeZeroCapacityKnapsack(t *testing.T) {
	got, err := Knapsack01([]int{1, 2, 3}, []int{10, 20, 30}, 0)
	if err != nil || got != 0 {
		t.Errorf("Knapsack01 capacity 0 = %d, %v, want 0, nil", got, err)
	}

	// Fractional knapsack with capacity 0 and a zero-weight, positive-value item.
	gotF, err := FractionalKnapsack([]float64{0, 2}, []float64{50, 20}, 0)
	if err != nil || math.Abs(gotF-50) > 1e-9 {
		t.Errorf("FractionalKnapsack zero-weight = %v, %v, want 50, nil", gotF, err)
	}

	gotF2, err := FractionalKnapsack([]float64{2, 3}, []float64{10, 30}, 0)
	if err != nil || gotF2 != 0 {
		t.Errorf("FractionalKnapsack capacity 0 = %v, %v, want 0, nil", gotF2, err)
	}
}

// TestEdgeHuffmanSingleSymbol verifies that a single-symbol input still gets a
// valid code table.
func TestEdgeHuffmanSingleSymbol(t *testing.T) {
	codes, err := HuffmanCoding(map[string]int{"x": 3})
	if err != nil {
		t.Fatalf("HuffmanCoding single symbol error: %v", err)
	}
	if len(codes) != 1 || codes["x"] != "0" {
		t.Errorf("HuffmanCoding single symbol = %v, want map[x:0]", codes)
	}
}

// TestEdgeFastPowerExtremes covers 0 raised to a negative power and a
// math.MinInt-sized exponent.
func TestEdgeFastPowerExtremes(t *testing.T) {
	for _, exp := range []int{-1, -3, math.MinInt} {
		_, err := FastPower(0, exp)
		if !errors.Is(err, ErrInvalidInput) {
			t.Errorf("FastPower(0,%d) error = %v, want ErrInvalidInput", exp, err)
		}
	}
	if got, err := FastPower(2, math.MinInt); err != nil || got != 0 {
		t.Errorf("FastPower(2,MinInt) = %v, %v, want 0, nil", got, err)
	}
	if got, err := FastPower(-2, math.MinInt); err != nil || got != 0 {
		t.Errorf("FastPower(-2,MinInt) = %v, %v, want 0, nil", got, err)
	}
}

// TestEdgeRadixSortExtremes covers empty / single / duplicate inputs, negative
// rejection, and the largest representable positive int.
func TestEdgeRadixSortExtremes(t *testing.T) {
	cases := []struct {
		in   []int
		want []int
	}{
		{[]int{}, []int{}},
		{[]int{0}, []int{0}},
		{[]int{math.MaxInt, 0, math.MaxInt}, []int{0, math.MaxInt, math.MaxInt}},
		{[]int{5, 5, 5}, []int{5, 5, 5}},
	}
	for _, c := range cases {
		got, err := RadixSort(c.in)
		if err != nil || !reflect.DeepEqual(got, c.want) {
			t.Errorf("RadixSort(%v) = %v, %v, want %v, nil", c.in, got, err, c.want)
		}
	}

	if _, err := RadixSort([]int{-1}); !errors.Is(err, ErrInvalidInput) {
		t.Errorf("RadixSort(negative) error = %v, want ErrInvalidInput", err)
	}
}
