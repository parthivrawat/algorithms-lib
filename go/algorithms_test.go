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
}

func TestTimSort(t *testing.T) {
	in := []int{3, 1, 4, 1, 5, 9, 2, 6}
	got := TimSort(in)
	want := []int{1, 1, 2, 3, 4, 5, 6, 9}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("TimSort(%v) = %v, want %v", in, got, want)
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
}

func TestFastPower(t *testing.T) {
	if FastPower(2, 10) != 1024 {
		t.Errorf("FastPower(2,10) = %v, want 1024", FastPower(2, 10))
	}
	if FastPower(5, 0) != 1 {
		t.Errorf("FastPower(5,0) = %v, want 1", FastPower(5, 0))
	}
	if FastPower(2, -2) != 0.25 {
		t.Errorf("FastPower(2,-2) = %v, want 0.25", FastPower(2, -2))
	}
}
