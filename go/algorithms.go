// Package algorithms provides a comprehensive, zero-dependency collection
// of common algorithms for Go.
package algorithms

import (
	"cmp"
	"container/heap"
	"errors"
	"fmt"
	"math"
	"math/bits"
	"sort"
	"strings"
)

// Ordered is a constraint for types that support ordering operators.
// It is an alias for the standard library's cmp.Ordered (Go 1.21+).
type Ordered = cmp.Ordered

// Common errors returned by this package.
var (
	ErrNotFound      = errors.New("not found")
	ErrInvalidInput  = errors.New("invalid input")
	ErrInvalidGraph  = errors.New("invalid graph")
	ErrNegativeCycle = errors.New("negative-weight cycle detected")
)

// ---------------- Sorting ----------------

// QuickSort returns a new sorted slice using quick sort.
func QuickSort[T Ordered](items []T) []T {
	if len(items) <= 1 {
		out := make([]T, len(items))
		copy(out, items)
		return out
	}
	pivot := items[len(items)/2]
	var left, middle, right []T
	for _, x := range items {
		if x < pivot {
			left = append(left, x)
		} else if x > pivot {
			right = append(right, x)
		} else {
			middle = append(middle, x)
		}
	}
	return append(append(QuickSort(left), middle...), QuickSort(right)...)
}

// MergeSort returns a new sorted slice using stable merge sort.
func MergeSort[T Ordered](items []T) []T {
	if len(items) <= 1 {
		out := make([]T, len(items))
		copy(out, items)
		return out
	}
	mid := len(items) / 2
	left := MergeSort(items[:mid])
	right := MergeSort(items[mid:])
	return merge(left, right)
}

func merge[T Ordered](left, right []T) []T {
	merged := make([]T, 0, len(left)+len(right))
	i, j := 0, 0
	for i < len(left) && j < len(right) {
		if left[i] <= right[j] {
			merged = append(merged, left[i])
			i++
		} else {
			merged = append(merged, right[j])
			j++
		}
	}
	merged = append(merged, left[i:]...)
	merged = append(merged, right[j:]...)
	return merged
}

// HeapSort returns a new sorted slice using a binary max-heap.
func HeapSort[T Ordered](items []T) []T {
	arr := make([]T, len(items))
	copy(arr, items)
	n := len(arr)
	if n <= 1 {
		return arr
	}
	for i := n/2 - 1; i >= 0; i-- {
		siftDown(arr, n, i)
	}
	for end := n - 1; end > 0; end-- {
		arr[0], arr[end] = arr[end], arr[0]
		siftDown(arr, end, 0)
	}
	return arr
}

func siftDown[T Ordered](arr []T, n, i int) {
	for {
		largest := i
		left := 2*i + 1
		right := 2*i + 2
		if left < n && arr[left] > arr[largest] {
			largest = left
		}
		if right < n && arr[right] > arr[largest] {
			largest = right
		}
		if largest == i {
			break
		}
		arr[i], arr[largest] = arr[largest], arr[i]
		i = largest
	}
}

// RadixSort returns a new sorted slice of non-negative integers using LSD radix sort.
func RadixSort(items []int) ([]int, error) {
	if len(items) == 0 {
		return []int{}, nil
	}
	maxVal := 0
	for _, v := range items {
		if v < 0 {
			return nil, ErrInvalidInput
		}
		if v > maxVal {
			maxVal = v
		}
	}
	arr := make([]int, len(items))
	copy(arr, items)
	for exp := 1; maxVal/exp > 0; {
		countingSortByDigit(arr, exp)
		// Stop before exp*10 can overflow int. No further digit pass is needed
		// anyway: a value that large has no higher digit than exp's place.
		if exp > math.MaxInt/10 {
			break
		}
		exp *= 10
	}
	return arr, nil
}

func countingSortByDigit(arr []int, exp int) {
	n := len(arr)
	output := make([]int, n)
	count := make([]int, 10)
	for _, num := range arr {
		count[(num/exp)%10]++
	}
	for i := 1; i < 10; i++ {
		count[i] += count[i-1]
	}
	for i := n - 1; i >= 0; i-- {
		index := (arr[i] / exp) % 10
		output[count[index]-1] = arr[i]
		count[index]--
	}
	copy(arr, output)
}

// NativeSort returns a new sorted slice using the standard library sort.
func NativeSort[T Ordered](items []T) []T {
	arr := make([]T, len(items))
	copy(arr, items)
	sort.Slice(arr, func(i, j int) bool { return arr[i] < arr[j] })
	return arr
}

// ---------------- Searching ----------------

func requireSorted[T Ordered](arr []T, name string) error {
	for i := 1; i < len(arr); i++ {
		if arr[i] < arr[i-1] {
			return fmt.Errorf("%s: %w", name, ErrInvalidInput)
		}
	}
	return nil
}

// BinarySearch returns the index of the target in a sorted slice.
// Returns ErrNotFound if the target is not present and ErrInvalidInput if the slice is not sorted.
func BinarySearch[T Ordered](arr []T, target T) (int, error) {
	if err := requireSorted(arr, "BinarySearch"); err != nil {
		return -1, err
	}
	lo, hi := 0, len(arr)-1
	for lo <= hi {
		mid := lo + (hi-lo)/2
		if arr[mid] == target {
			return mid, nil
		}
		if arr[mid] < target {
			lo = mid + 1
		} else {
			hi = mid - 1
		}
	}
	return -1, ErrNotFound
}

// InterpolationSearch returns the index of the target in a sorted numeric slice.
func InterpolationSearch(arr []int, target int) (int, error) {
	if err := requireSorted(arr, "InterpolationSearch"); err != nil {
		return -1, err
	}
	lo, hi := 0, len(arr)-1
	for lo <= hi && target >= arr[lo] && target <= arr[hi] {
		if arr[hi] == arr[lo] {
			if arr[lo] == target {
				return lo, nil
			}
			break
		}
		// Use 128-bit unsigned arithmetic so large signed int ranges (e.g.
		// math.MinInt to math.MaxInt) don't overflow before the division.
		diffTarget := uint64(target) - uint64(arr[lo])
		diffPos := uint64(hi - lo)
		diffVal := uint64(arr[hi]) - uint64(arr[lo])
		hiPart, loPart := bits.Mul64(diffTarget, diffPos)
		offset, _ := bits.Div64(hiPart, loPart, diffVal)
		pos := lo + int(offset)
		if pos < lo || pos > hi {
			break
		}
		if arr[pos] == target {
			return pos, nil
		}
		if arr[pos] < target {
			lo = pos + 1
		} else {
			hi = pos - 1
		}
	}
	return -1, ErrNotFound
}

// JumpSearch returns the index of the target in a sorted slice.
func JumpSearch[T Ordered](arr []T, target T) (int, error) {
	if err := requireSorted(arr, "JumpSearch"); err != nil {
		return -1, err
	}
	n := len(arr)
	if n == 0 {
		return -1, ErrNotFound
	}
	step := int(math.Sqrt(float64(n)))
	prev := 0
	for arr[min(step, n)-1] < target {
		prev = step
		step += int(math.Sqrt(float64(n)))
		if prev >= n {
			return -1, ErrNotFound
		}
	}
	for arr[prev] < target {
		prev++
		if prev == min(step, n) {
			return -1, ErrNotFound
		}
	}
	if arr[prev] == target {
		return prev, nil
	}
	return -1, ErrNotFound
}

// BinarySearchOnAnswer finds the smallest or largest x in [low, high] that
// satisfies a monotone predicate.
func BinarySearchOnAnswer(low, high int, predicate func(int) bool, find string) (int, error) {
	if find != "minimum" && find != "maximum" {
		return 0, fmt.Errorf("BinarySearchOnAnswer: find must be 'minimum' or 'maximum': %w", ErrInvalidInput)
	}
	if low > high {
		return 0, ErrNotFound
	}
	left, right := low, high
	answer := low - 1
	if find == "minimum" {
		for left <= right {
			mid := left + (right-left)/2
			if predicate(mid) {
				answer = mid
				right = mid - 1
			} else {
				left = mid + 1
			}
		}
	} else {
		for left <= right {
			mid := left + (right-left)/2
			if predicate(mid) {
				answer = mid
				left = mid + 1
			} else {
				right = mid - 1
			}
		}
	}
	if answer < low {
		return 0, ErrNotFound
	}
	return answer, nil
}

// ---------------- Graphs ----------------

// Edge represents a weighted directed edge in a graph.
type Edge struct {
	To     string
	Weight float64
}

// validateAdjacencyList returns ErrInvalidGraph if any neighbor referenced in
// the adjacency list is not itself a key of the graph.
func validateAdjacencyList(graph map[string][]string, name string) error {
	for node, neighbors := range graph {
		for _, nbr := range neighbors {
			if _, ok := graph[nbr]; !ok {
				return fmt.Errorf("%s: neighbor %q of %q is not a key in the graph: %w",
					name, nbr, node, ErrInvalidGraph)
			}
		}
	}
	return nil
}

// BFS returns the vertices in breadth-first order.
// Every neighbor listed in the adjacency list must also be a key of the graph.
func BFS(graph map[string][]string, start string) ([]string, error) {
	if _, ok := graph[start]; !ok {
		return nil, ErrInvalidGraph
	}
	if err := validateAdjacencyList(graph, "BFS"); err != nil {
		return nil, err
	}
	visited := make(map[string]bool)
	visited[start] = true
	queue := []string{start}
	result := []string{}
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		result = append(result, current)
		for _, nbr := range graph[current] {
			if !visited[nbr] {
				visited[nbr] = true
				queue = append(queue, nbr)
			}
		}
	}
	return result, nil
}

// DFS returns the vertices in depth-first order.
// Every neighbor listed in the adjacency list must also be a key of the graph.
func DFS(graph map[string][]string, start string) ([]string, error) {
	if _, ok := graph[start]; !ok {
		return nil, ErrInvalidGraph
	}
	if err := validateAdjacencyList(graph, "DFS"); err != nil {
		return nil, err
	}
	visited := make(map[string]bool)
	visited[start] = true
	stack := []string{start}
	result := []string{}
	for len(stack) > 0 {
		current := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		result = append(result, current)
		neighbors := graph[current]
		for i := len(neighbors) - 1; i >= 0; i-- {
			nbr := neighbors[i]
			if !visited[nbr] {
				visited[nbr] = true
				stack = append(stack, nbr)
			}
		}
	}
	return result, nil
}

type pqItem struct {
	node     string
	priority float64
}

type graphPQ []*pqItem

func (pq graphPQ) Len() int { return len(pq) }
func (pq graphPQ) Less(i, j int) bool {
	if pq[i].priority == pq[j].priority {
		return pq[i].node < pq[j].node
	}
	return pq[i].priority < pq[j].priority
}
func (pq graphPQ) Swap(i, j int) { pq[i], pq[j] = pq[j], pq[i] }
func (pq *graphPQ) Push(x interface{}) {
	*pq = append(*pq, x.(*pqItem))
}
func (pq *graphPQ) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil // release the popped slot so the item can be GC'd
	*pq = old[:n-1]
	return item
}

// Dijkstra computes shortest distances from the start node.
// Edge weights must be non-negative.
func Dijkstra(graph map[string][]Edge, start string) (map[string]float64, map[string]string, error) {
	if _, ok := graph[start]; !ok {
		return nil, nil, ErrInvalidGraph
	}

	// Pre-scan all edges so negative weights in unreachable components are
	// detected, and collect every vertex referenced in the graph.
	allNodes := make(map[string]struct{})
	for k := range graph {
		allNodes[k] = struct{}{}
	}
	for _, edges := range graph {
		for _, e := range edges {
			if e.Weight < 0 {
				return nil, nil, ErrInvalidGraph
			}
			allNodes[e.To] = struct{}{}
		}
	}

	distances := make(map[string]float64)
	for k := range allNodes {
		distances[k] = math.Inf(1)
	}
	distances[start] = 0
	predecessors := make(map[string]string)

	pq := &graphPQ{}
	heap.Init(pq)
	heap.Push(pq, &pqItem{node: start, priority: 0})

	for pq.Len() > 0 {
		item := heap.Pop(pq).(*pqItem)
		node := item.node
		dist := item.priority
		if dist > distances[node] {
			continue
		}
		for _, e := range graph[node] {
			nd := dist + e.Weight
			if nd < distances[e.To] {
				distances[e.To] = nd
				predecessors[e.To] = node
				heap.Push(pq, &pqItem{node: e.To, priority: nd})
			}
		}
	}
	return distances, predecessors, nil
}

// AStar finds a shortest path from start to goal using a heuristic.
//
// The heuristic must be non-negative. For an optimal first result it should be
// admissible and consistent; inconsistent but admissible heuristics still
// return an optimal path, but may cause nodes to be re-expanded.
//
// All edge weights must be non-negative.
func AStar(graph map[string][]Edge, start, goal string, heuristic func(string, string) float64) ([]string, float64, error) {
	if _, ok := graph[start]; !ok {
		return nil, 0, ErrInvalidGraph
	}
	if _, ok := graph[goal]; !ok {
		return nil, 0, ErrInvalidGraph
	}

	h := heuristic(start, goal)
	if h < 0 {
		return nil, 0, ErrInvalidInput
	}

	gScore := map[string]float64{start: 0}
	bestF := map[string]float64{start: h}
	cameFrom := make(map[string]string)
	pq := &graphPQ{}
	heap.Init(pq)
	heap.Push(pq, &pqItem{node: start, priority: h})

	for pq.Len() > 0 {
		item := heap.Pop(pq).(*pqItem)
		node := item.node
		if f, ok := bestF[node]; ok && item.priority > f {
			continue
		}
		if node == goal {
			return reconstructPath(cameFrom, node), gScore[goal], nil
		}
		for _, e := range graph[node] {
			if e.Weight < 0 {
				return nil, 0, ErrInvalidGraph
			}
			hNbr := heuristic(e.To, goal)
			if hNbr < 0 {
				return nil, 0, ErrInvalidInput
			}
			tentative := gScore[node] + e.Weight
			if g, ok := gScore[e.To]; !ok || tentative < g {
				cameFrom[e.To] = node
				gScore[e.To] = tentative
				f := tentative + hNbr
				if best, ok := bestF[e.To]; !ok || f < best {
					bestF[e.To] = f
					heap.Push(pq, &pqItem{node: e.To, priority: f})
				}
			}
		}
	}
	return nil, 0, ErrNotFound
}

func reconstructPath(cameFrom map[string]string, current string) []string {
	path := []string{current}
	for {
		prev, ok := cameFrom[current]
		if !ok {
			break
		}
		current = prev
		path = append(path, current)
	}
	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}
	return path
}

type edge3 struct {
	from   string
	to     string
	weight float64
}

// BellmanFord computes shortest distances and detects negative cycles.
func BellmanFord(graph map[string][]Edge, start string) (map[string]float64, map[string]string, error) {
	if _, ok := graph[start]; !ok {
		return nil, nil, ErrInvalidGraph
	}
	vertices := make(map[string]struct{})
	var edges []edge3
	for u, list := range graph {
		vertices[u] = struct{}{}
		for _, e := range list {
			vertices[e.To] = struct{}{}
			edges = append(edges, edge3{from: u, to: e.To, weight: e.Weight})
		}
	}
	distances := map[string]float64{start: 0}
	predecessors := make(map[string]string)

	for i := 0; i < len(vertices)-1; i++ {
		updated := false
		for _, e := range edges {
			du, ok := distances[e.from]
			if !ok {
				continue
			}
			dv := math.Inf(1)
			if v, ok := distances[e.to]; ok {
				dv = v
			}
			if du+e.weight < dv {
				distances[e.to] = du + e.weight
				predecessors[e.to] = e.from
				updated = true
			}
		}
		if !updated {
			break
		}
	}

	for _, e := range edges {
		du, ok := distances[e.from]
		if !ok {
			continue
		}
		dv := math.Inf(1)
		if v, ok := distances[e.to]; ok {
			dv = v
		}
		if du+e.weight < dv {
			return nil, nil, ErrNegativeCycle
		}
	}
	return distances, predecessors, nil
}

// ---------------- Dynamic programming ----------------

// Knapsack01 returns the maximum value for the 0/1 knapsack problem.
func Knapsack01(weights, values []int, capacity int) (int, error) {
	if len(weights) != len(values) {
		return 0, ErrInvalidInput
	}
	if capacity < 0 {
		return 0, ErrInvalidInput
	}
	dp := make([]int, capacity+1)
	for i := 0; i < len(weights); i++ {
		w := weights[i]
		v := values[i]
		if w < 0 {
			return 0, ErrInvalidInput
		}
		for c := capacity; c >= w; c-- {
			if dp[c-w]+v > dp[c] {
				dp[c] = dp[c-w] + v
			}
		}
	}
	return dp[capacity], nil
}

// LongestCommonSubsequence returns the length of the LCS.
func LongestCommonSubsequence[T comparable](a, b []T) int {
	m, n := len(a), len(b)
	prev := make([]int, n+1)
	for i := 1; i <= m; i++ {
		curr := make([]int, n+1)
		for j := 1; j <= n; j++ {
			if a[i-1] == b[j-1] {
				curr[j] = prev[j-1] + 1
			} else {
				if prev[j] > curr[j-1] {
					curr[j] = prev[j]
				} else {
					curr[j] = curr[j-1]
				}
			}
		}
		prev = curr
	}
	return prev[n]
}

// EditDistance returns the Levenshtein distance between two sequences.
func EditDistance[T comparable](a, b []T) int {
	m, n := len(a), len(b)
	prev := make([]int, n+1)
	for i := 0; i <= n; i++ {
		prev[i] = i
	}
	for i := 1; i <= m; i++ {
		curr := make([]int, n+1)
		curr[0] = i
		for j := 1; j <= n; j++ {
			cost := 1
			if a[i-1] == b[j-1] {
				cost = 0
			}
			insertion := curr[j-1] + 1
			deletion := prev[j] + 1
			substitution := prev[j-1] + cost
			curr[j] = minOfThree(insertion, deletion, substitution)
		}
		prev = curr
	}
	return prev[n]
}

func minOfThree(a, b, c int) int {
	if a < b {
		if a < c {
			return a
		}
		return c
	}
	if b < c {
		return b
	}
	return c
}

// LongestCommonSubsequenceReconstruction returns the length and one LCS.
func LongestCommonSubsequenceReconstruction[T comparable](a, b []T) (int, []T, error) {
	m, n := len(a), len(b)
	dp := make([][]int, m+1)
	for i := range dp {
		dp[i] = make([]int, n+1)
	}
	for i := 1; i <= m; i++ {
		for j := 1; j <= n; j++ {
			if a[i-1] == b[j-1] {
				dp[i][j] = dp[i-1][j-1] + 1
			} else {
				if dp[i-1][j] > dp[i][j-1] {
					dp[i][j] = dp[i-1][j]
				} else {
					dp[i][j] = dp[i][j-1]
				}
			}
		}
	}
	result := make([]T, 0, dp[m][n])
	i, j := m, n
	for i > 0 && j > 0 {
		if a[i-1] == b[j-1] {
			result = append(result, a[i-1])
			i--
			j--
		} else if dp[i-1][j] >= dp[i][j-1] {
			i--
		} else {
			j--
		}
	}
	for k, l := 0, len(result)-1; k < l; k, l = k+1, l-1 {
		result[k], result[l] = result[l], result[k]
	}
	return dp[m][n], result, nil
}

// EditOp represents one step in an edit script.
type EditOp[T comparable] struct {
	Action string
	From   *T
	To     *T
}

// EditDistanceReconstruction returns the Levenshtein distance and an edit script.
func EditDistanceReconstruction[T comparable](a, b []T) (int, []EditOp[T], error) {
	m, n := len(a), len(b)
	dp := make([][]int, m+1)
	for i := range dp {
		dp[i] = make([]int, n+1)
		dp[i][0] = i
	}
	for j := 0; j <= n; j++ {
		dp[0][j] = j
	}
	for i := 1; i <= m; i++ {
		for j := 1; j <= n; j++ {
			if a[i-1] == b[j-1] {
				dp[i][j] = dp[i-1][j-1]
			} else {
				dp[i][j] = 1 + minOfThree(dp[i-1][j], dp[i][j-1], dp[i-1][j-1])
			}
		}
	}
	var script []EditOp[T]
	i, j := m, n
	for i > 0 || j > 0 {
		if i == 0 {
			x := b[j-1]
			script = append(script, EditOp[T]{Action: "insert", To: &x})
			j--
		} else if j == 0 {
			x := a[i-1]
			script = append(script, EditOp[T]{Action: "delete", From: &x})
			i--
		} else if a[i-1] == b[j-1] {
			x, y := a[i-1], b[j-1]
			script = append(script, EditOp[T]{Action: "match", From: &x, To: &y})
			i--
			j--
		} else {
			best := dp[i][j]
			if dp[i-1][j-1]+1 == best {
				x, y := a[i-1], b[j-1]
				script = append(script, EditOp[T]{Action: "substitute", From: &x, To: &y})
				i--
				j--
			} else if dp[i][j-1]+1 == best {
				x := b[j-1]
				script = append(script, EditOp[T]{Action: "insert", To: &x})
				j--
			} else {
				x := a[i-1]
				script = append(script, EditOp[T]{Action: "delete", From: &x})
				i--
			}
		}
	}
	for k, l := 0, len(script)-1; k < l; k, l = k+1, l-1 {
		script[k], script[l] = script[l], script[k]
	}
	return dp[m][n], script, nil
}

// ---------------- String algorithms ----------------

// KMPSearch returns all starting indices of the pattern in the text.
func KMPSearch(text, pattern string) ([]int, error) {
	if pattern == "" {
		return nil, ErrInvalidInput
	}
	failure := computeKmpFailure(pattern)
	matches := []int{}
	j := 0
	for i := 0; i < len(text); i++ {
		for j > 0 && text[i] != pattern[j] {
			j = failure[j-1]
		}
		if text[i] == pattern[j] {
			j++
		}
		if j == len(pattern) {
			matches = append(matches, i-len(pattern)+1)
			j = failure[j-1]
		}
	}
	return matches, nil
}

func computeKmpFailure(pattern string) []int {
	failure := make([]int, len(pattern))
	j := 0
	for i := 1; i < len(pattern); i++ {
		for j > 0 && pattern[i] != pattern[j] {
			j = failure[j-1]
		}
		if pattern[i] == pattern[j] {
			j++
			failure[i] = j
		}
	}
	return failure
}

// mulMod returns (a*b) % mod using 128-bit unsigned arithmetic to avoid
// overflow. mod must be non-zero; a and b must fit in uint64.
func mulMod(a, b, mod uint64) uint64 {
	hi, lo := bits.Mul64(a, b)
	_, rem := bits.Div64(hi, lo, mod)
	return rem
}

// RabinKarpSearch returns all starting indices of the pattern using rolling hash.
//
// base and mod must be positive. mod is treated as an unsigned modulus to
// avoid int64 overflow during the rolling updates.
func RabinKarpSearch(text, pattern string, base int, mod int64) ([]int, error) {
	if pattern == "" {
		return nil, ErrInvalidInput
	}
	if mod <= 0 {
		return nil, ErrInvalidInput
	}
	if base <= 0 {
		return nil, ErrInvalidInput
	}

	n := len(text)
	m := len(pattern)
	if m > n {
		return []int{}, nil
	}

	baseU := uint64(base)
	modU := uint64(mod)

	h := uint64(1)
	for i := 0; i < m-1; i++ {
		h = mulMod(h, baseU, modU)
	}

	var patternHash, textHash uint64
	for i := 0; i < m; i++ {
		patternHash = mulMod(patternHash, baseU, modU) + uint64(pattern[i])
		if patternHash >= modU {
			patternHash -= modU
		}
		textHash = mulMod(textHash, baseU, modU) + uint64(text[i])
		if textHash >= modU {
			textHash -= modU
		}
	}

	matches := []int{}
	for i := 0; i <= n-m; i++ {
		if textHash == patternHash && text[i:i+m] == pattern {
			matches = append(matches, i)
		}
		if i < n-m {
			sub := mulMod(uint64(text[i]), h, modU)
			if textHash >= sub {
				textHash -= sub
			} else {
				textHash += modU - sub
			}
			textHash = mulMod(textHash, baseU, modU)
			textHash += uint64(text[i+m])
			if textHash >= modU {
				textHash -= modU
			}
		}
	}
	return matches, nil
}

// BoyerMooreSearch returns all starting byte indices of the pattern in the
// text, including overlapping occurrences, using the bad-character and
// good-suffix heuristics. Like KMPSearch and RabinKarpSearch, indices are
// byte offsets, not rune offsets.
func BoyerMooreSearch(text, pattern string) ([]int, error) {
	if pattern == "" {
		return nil, ErrInvalidInput
	}
	n := len(text)
	m := len(pattern)
	if m > n {
		return []int{}, nil
	}
	badChar := make(map[byte]int)
	for i := 0; i < m; i++ {
		badChar[pattern[i]] = i
	}
	goodSuffix := goodSuffixShifts(pattern)
	matches := []int{}
	i := 0
	for i <= n-m {
		j := m - 1
		for j >= 0 && text[i+j] == pattern[j] {
			j--
		}
		if j < 0 {
			matches = append(matches, i)
			i += goodSuffix[0]
		} else {
			last, ok := badChar[text[i+j]]
			if !ok {
				last = -1
			}
			shift := goodSuffix[j+1]
			if bc := j - last; bc > shift {
				shift = bc
			}
			i += shift
		}
	}
	return matches, nil
}

// goodSuffixShifts builds the Boyer-Moore good-suffix table: shift[k] is the
// shift applied after a mismatch at pattern index k-1, and shift[0] is the
// border-based shift applied after a full match so that overlapping matches
// are still reported.
func goodSuffixShifts(pattern string) []int {
	m := len(pattern)
	shift := make([]int, m+1)
	borderPos := make([]int, m+1)
	i, j := m, m+1
	borderPos[i] = j
	for i > 0 {
		for j <= m && pattern[i-1] != pattern[j-1] {
			if shift[j] == 0 {
				shift[j] = j - i
			}
			j = borderPos[j]
		}
		i--
		j--
		borderPos[i] = j
	}
	j = borderPos[0]
	for i := 0; i <= m; i++ {
		if shift[i] == 0 {
			shift[i] = j
		}
		if i == j {
			j = borderPos[j]
		}
	}
	return shift
}

// ---------------- Greedy ----------------

// ActivitySelection returns the indices of a maximum set of compatible activities.
func ActivitySelection(activities [][2]int) ([]int, error) {
	if len(activities) == 0 {
		return []int{}, nil
	}
	type indexed struct {
		index int
		start int
		end   int
	}
	indexedList := make([]indexed, len(activities))
	for i, a := range activities {
		if a[0] > a[1] {
			return nil, ErrInvalidInput
		}
		indexedList[i] = indexed{index: i, start: a[0], end: a[1]}
	}
	sort.Slice(indexedList, func(i, j int) bool {
		return indexedList[i].end < indexedList[j].end
	})
	selected := []int{indexedList[0].index}
	lastEnd := indexedList[0].end
	for i := 1; i < len(indexedList); i++ {
		if indexedList[i].start >= lastEnd {
			selected = append(selected, indexedList[i].index)
			lastEnd = indexedList[i].end
		}
	}
	return selected, nil
}

// FractionalKnapsack returns the maximum value for the fractional knapsack
// problem. Negative or NaN weights and NaN values are rejected with
// ErrInvalidInput; a zero-weight item contributes its full value when that
// value is positive, since it consumes no capacity. Items with a
// non-positive value-to-weight ratio are never taken.
func FractionalKnapsack(weights, values []float64, capacity float64) (float64, error) {
	if len(weights) != len(values) {
		return 0, ErrInvalidInput
	}
	if capacity < 0 {
		return 0, ErrInvalidInput
	}
	type item struct {
		ratio  float64
		weight float64
	}
	items := make([]item, 0, len(weights))
	total := 0.0
	for i := 0; i < len(weights); i++ {
		if weights[i] < 0 || math.IsNaN(weights[i]) || math.IsNaN(values[i]) {
			return 0, ErrInvalidInput
		}
		if weights[i] == 0 {
			// Zero-weight items consume no capacity: take them if they add value.
			total += math.Max(0, values[i])
			continue
		}
		items = append(items, item{
			ratio:  values[i] / weights[i],
			weight: weights[i],
		})
	}
	sort.Slice(items, func(i, j int) bool { return items[i].ratio > items[j].ratio })
	remaining := capacity
	for _, it := range items {
		if it.ratio <= 0 {
			break // taking a non-positive-value item can only lower the total
		}
		take := math.Min(it.weight, remaining)
		total += take * it.ratio
		remaining -= take
		if remaining <= 0 {
			break
		}
	}
	return total, nil
}

type huffmanNode struct {
	freq  int
	char  string
	leaf  bool
	left  *huffmanNode
	right *huffmanNode
}

type huffmanPQ []*huffmanNode

func (pq huffmanPQ) Len() int { return len(pq) }
func (pq huffmanPQ) Less(i, j int) bool {
	if pq[i].freq == pq[j].freq {
		return pq[i].char < pq[j].char
	}
	return pq[i].freq < pq[j].freq
}
func (pq huffmanPQ) Swap(i, j int) { pq[i], pq[j] = pq[j], pq[i] }
func (pq *huffmanPQ) Push(x interface{}) {
	*pq = append(*pq, x.(*huffmanNode))
}
func (pq *huffmanPQ) Pop() interface{} {
	old := *pq
	n := len(old)
	node := old[n-1]
	old[n-1] = nil // release the popped slot so the node can be GC'd
	*pq = old[:n-1]
	return node
}

// HuffmanCoding returns a prefix-free code table for the given frequencies.
// Every frequency must be positive. Symbol keys may be any string, including
// the empty string.
func HuffmanCoding(frequencies map[string]int) (map[string]string, error) {
	if len(frequencies) == 0 {
		return nil, ErrInvalidInput
	}
	for char, freq := range frequencies {
		if freq <= 0 {
			return nil, fmt.Errorf(
				"HuffmanCoding: frequency for %q must be positive: %w",
				char, ErrInvalidInput)
		}
	}
	pq := &huffmanPQ{}
	heap.Init(pq)
	for char, freq := range frequencies {
		heap.Push(pq, &huffmanNode{freq: freq, char: char, leaf: true})
	}
	if pq.Len() == 1 {
		node := heap.Pop(pq).(*huffmanNode)
		return map[string]string{node.char: "0"}, nil
	}
	for pq.Len() > 1 {
		left := heap.Pop(pq).(*huffmanNode)
		right := heap.Pop(pq).(*huffmanNode)
		parent := &huffmanNode{
			freq:  left.freq + right.freq,
			left:  left,
			right: right,
		}
		heap.Push(pq, parent)
	}
	root := heap.Pop(pq).(*huffmanNode)
	codes := make(map[string]string)
	assignHuffmanCodes(root, "", codes)
	return codes, nil
}

func assignHuffmanCodes(node *huffmanNode, prefix string, codes map[string]string) {
	if node == nil {
		return
	}
	if node.leaf {
		if prefix == "" {
			codes[node.char] = "0"
		} else {
			codes[node.char] = prefix
		}
		return
	}
	assignHuffmanCodes(node.left, prefix+"0", codes)
	assignHuffmanCodes(node.right, prefix+"1", codes)
}

// HuffmanEncode encodes a slice of symbols using a Huffman code table.
func HuffmanEncode(symbols []string, codeTable map[string]string) (string, error) {
	var parts []string
	for _, s := range symbols {
		code, ok := codeTable[s]
		if !ok {
			return "", fmt.Errorf("HuffmanEncode: symbol %q not in code table: %w", s, ErrInvalidInput)
		}
		parts = append(parts, code)
	}
	return strings.Join(parts, ""), nil
}

// HuffmanDecode decodes a Huffman bit string using a code table.
func HuffmanDecode(encoded string, codeTable map[string]string) ([]string, error) {
	reverse := make(map[string]string, len(codeTable))
	for s, code := range codeTable {
		reverse[code] = s
	}
	var result []string
	current := ""
	for _, bit := range encoded {
		current += string(bit)
		if s, ok := reverse[current]; ok {
			result = append(result, s)
			current = ""
		}
	}
	if current != "" {
		return nil, fmt.Errorf("HuffmanDecode: incomplete or invalid bit string: %w", ErrInvalidInput)
	}
	return result, nil
}

// ---------------- Divide and conquer ----------------

// MaxSubarray returns the maximum sum of any contiguous subarray.
func MaxSubarray(arr []int) (int, error) {
	if len(arr) == 0 {
		return 0, ErrInvalidInput
	}
	return maxSubarrayDc(arr, 0, len(arr)-1), nil
}

func maxSubarrayDc(arr []int, left, right int) int {
	if left == right {
		return arr[left]
	}
	mid := left + (right-left)/2
	leftSum := maxSubarrayDc(arr, left, mid)
	rightSum := maxSubarrayDc(arr, mid+1, right)
	crossSum := maxCrossingSum(arr, left, mid, right)
	if leftSum > rightSum {
		if leftSum > crossSum {
			return leftSum
		}
		return crossSum
	}
	if rightSum > crossSum {
		return rightSum
	}
	return crossSum
}

func maxCrossingSum(arr []int, left, mid, right int) int {
	sum := arr[mid]
	leftMax := arr[mid]
	for i := mid - 1; i >= left; i-- {
		sum += arr[i]
		if sum > leftMax {
			leftMax = sum
		}
	}
	sum = arr[mid+1]
	rightMax := arr[mid+1]
	for i := mid + 2; i <= right; i++ {
		sum += arr[i]
		if sum > rightMax {
			rightMax = sum
		}
	}
	return leftMax + rightMax
}

// CountInversions returns the number of inversions in the slice.
// The count is int64 so it cannot overflow for large inputs on 32-bit
// platforms: a slice of n elements can have n*(n-1)/2 inversions.
func CountInversions[T Ordered](arr []T) int64 {
	_, count := sortAndCount(arr)
	return count
}

func sortAndCount[T Ordered](arr []T) ([]T, int64) {
	n := len(arr)
	if n <= 1 {
		out := make([]T, n)
		copy(out, arr)
		return out, 0
	}
	mid := n / 2
	left, lcount := sortAndCount(arr[:mid])
	right, rcount := sortAndCount(arr[mid:])
	merged, scount := mergeAndCount(left, right)
	return merged, lcount + rcount + scount
}

func mergeAndCount[T Ordered](left, right []T) ([]T, int64) {
	merged := make([]T, 0, len(left)+len(right))
	count := int64(0)
	i, j := 0, 0
	for i < len(left) && j < len(right) {
		if left[i] <= right[j] {
			merged = append(merged, left[i])
			i++
		} else {
			merged = append(merged, right[j])
			j++
			count += int64(len(left) - i)
		}
	}
	merged = append(merged, left[i:]...)
	merged = append(merged, right[j:]...)
	return merged, count
}

// FastPower returns base raised to the exponent using exponentiation by
// squaring. It returns ErrInvalidInput when base is zero and the exponent is
// negative. By convention, 0^0 returns 1.
func FastPower(base float64, exponent int) (float64, error) {
	if base == 0 && exponent < 0 {
		return 0, ErrInvalidInput
	}
	// Use the unsigned magnitude of the exponent so that math.MinInt negates
	// correctly: -uint64(e) wraps to |e| for every negative e.
	exp := uint64(exponent)
	if exponent < 0 {
		exp = -exp
	}
	result := 1.0
	b := base
	for exp > 0 {
		if exp&1 == 1 {
			result *= b
		}
		b *= b
		exp >>= 1
	}
	if exponent < 0 {
		result = 1 / result
	}
	return result, nil
}

// Gcd returns the greatest common divisor of a and b.
// The result is always non-negative. Gcd(0, 0) returns 0.
func Gcd(a, b int) (int, error) {
	if a < 0 {
		a = -a
	}
	if b < 0 {
		b = -b
	}
	for b != 0 {
		a, b = b, a%b
	}
	return a, nil
}

// Lcm returns the least common multiple of a and b.
// It returns ErrInvalidInput if both inputs are zero or if the result overflows int.
func Lcm(a, b int) (int, error) {
	if a == 0 && b == 0 {
		return 0, ErrInvalidInput
	}
	g, err := Gcd(a, b)
	if err != nil {
		return 0, err
	}
	// Compute (a/g) * b in int64 to detect overflow.
	prod := int64(a/g) * int64(b)
	if prod < 0 {
		prod = -prod
	}
	if prod > int64(math.MaxInt) {
		return 0, ErrInvalidInput
	}
	return int(prod), nil
}

// ModularFastPower returns (base^exponent) % mod using exponentiation by squaring.
// It returns ErrInvalidInput if mod is not positive or exponent is negative.
func ModularFastPower(base, exponent, mod int) (int, error) {
	if mod <= 0 {
		return 0, ErrInvalidInput
	}
	if exponent < 0 {
		return 0, ErrInvalidInput
	}
	if mod == 1 {
		return 0, nil
	}
	m := int64(mod)
	b := int64(base % mod)
	if b < 0 {
		b += m
	}
	var result int64 = 1
	e := int64(exponent)
	for e > 0 {
		if e&1 == 1 {
			result = (result * b) % m
		}
		b = (b * b) % m
		e >>= 1
	}
	return int(result), nil
}

// QuickSelect returns the k-th smallest element (0-indexed) of items.
// It returns ErrInvalidInput if k is out of range.
func QuickSelect[T Ordered](items []T, k int) (T, error) {
	var zero T
	if k < 0 || k >= len(items) {
		return zero, ErrInvalidInput
	}
	arr := make([]T, len(items))
	copy(arr, items)
	for {
		if len(arr) == 1 {
			return arr[0], nil
		}
		pivot := arr[len(arr)/2]
		var lows, pivots, highs []T
		for _, x := range arr {
			if x < pivot {
				lows = append(lows, x)
			} else if x > pivot {
				highs = append(highs, x)
			} else {
				pivots = append(pivots, x)
			}
		}
		if k < len(lows) {
			arr = lows
		} else if k < len(lows)+len(pivots) {
			return pivot, nil
		} else {
			k -= len(lows) + len(pivots)
			arr = highs
		}
	}
}

// SieveOfEratosthenes returns all prime numbers less than or equal to n.
func SieveOfEratosthenes(n int) ([]int, error) {
	if n < 0 {
		return nil, ErrInvalidInput
	}
	if n < 2 {
		return []int{}, nil
	}
	sieve := make([]bool, n+1)
	for i := 2; i <= n; i++ {
		sieve[i] = true
	}
	for p := 2; p*p <= n; p++ {
		if sieve[p] {
			for multiple := p * p; multiple <= n; multiple += p {
				sieve[multiple] = false
			}
		}
	}
	primes := make([]int, 0)
	for i := 2; i <= n; i++ {
		if sieve[i] {
			primes = append(primes, i)
		}
	}
	return primes, nil
}

// TopologicalSort returns a topological ordering of the DAG.
func TopologicalSort(graph map[string][]string) ([]string, error) {
	if err := validateAdjacencyList(graph, "TopologicalSort"); err != nil {
		return nil, err
	}
	inDegree := make(map[string]int)
	for node := range graph {
		inDegree[node] = 0
	}
	for _, neighbors := range graph {
		for _, nbr := range neighbors {
			inDegree[nbr]++
		}
	}

	var queue []string
	for node, degree := range inDegree {
		if degree == 0 {
			queue = append(queue, node)
		}
	}

	order := make([]string, 0, len(inDegree))
	head := 0
	for head < len(queue) {
		node := queue[head]
		head++
		order = append(order, node)
		for _, nbr := range graph[node] {
			inDegree[nbr]--
			if inDegree[nbr] == 0 {
				queue = append(queue, nbr)
			}
		}
	}

	if len(order) != len(inDegree) {
		return nil, ErrInvalidGraph
	}
	return order, nil
}

// UnionFind is a disjoint-set union–find data structure.
type UnionFind[T comparable] struct {
	parent map[T]T
	rank   map[T]int
}

// NewUnionFind creates a union–find structure for the given elements.
func NewUnionFind[T comparable](elements []T) *UnionFind[T] {
	uf := &UnionFind[T]{
		parent: make(map[T]T),
		rank:   make(map[T]int),
	}
	for _, x := range elements {
		uf.parent[x] = x
		uf.rank[x] = 0
	}
	return uf
}

// Find returns the representative of the set containing x.
func (uf *UnionFind[T]) Find(x T) (T, error) {
	var zero T
	if _, ok := uf.parent[x]; !ok {
		return zero, ErrInvalidInput
	}
	root := x
	for uf.parent[root] != root {
		root = uf.parent[root]
	}
	node := x
	for uf.parent[node] != root {
		parent := uf.parent[node]
		uf.parent[node] = root
		node = parent
	}
	return root, nil
}

// Union merges the sets containing x and y.
func (uf *UnionFind[T]) Union(x, y T) (bool, error) {
	rootX, err := uf.Find(x)
	if err != nil {
		return false, err
	}
	rootY, err := uf.Find(y)
	if err != nil {
		return false, err
	}
	if rootX == rootY {
		return false, nil
	}
	rankX := uf.rank[rootX]
	rankY := uf.rank[rootY]
	if rankX < rankY {
		uf.parent[rootX] = rootY
	} else if rankX > rankY {
		uf.parent[rootY] = rootX
	} else {
		uf.parent[rootY] = rootX
		uf.rank[rootX] = rankX + 1
	}
	return true, nil
}

// Connected reports whether x and y are in the same set.
func (uf *UnionFind[T]) Connected(x, y T) (bool, error) {
	rootX, err := uf.Find(x)
	if err != nil {
		return false, err
	}
	rootY, err := uf.Find(y)
	if err != nil {
		return false, err
	}
	return rootX == rootY, nil
}

// Count returns the number of disjoint sets.
func (uf *UnionFind[T]) Count() int {
	roots := make(map[T]struct{})
	for x := range uf.parent {
		r, _ := uf.Find(x)
		roots[r] = struct{}{}
	}
	return len(roots)
}

// PriorityQueue is a min-priority queue.
type PriorityQueue[T Ordered] struct {
	items []T
}

// NewPriorityQueue creates a priority queue with the given items.
func NewPriorityQueue[T Ordered](items []T) *PriorityQueue[T] {
	pq := &PriorityQueue[T]{items: append([]T(nil), items...)}
	for i := len(pq.items)/2 - 1; i >= 0; i-- {
		pq.siftDown(i)
	}
	return pq
}

// Push adds an item.
func (pq *PriorityQueue[T]) Push(item T) {
	pq.items = append(pq.items, item)
	pq.siftUp(len(pq.items) - 1)
}

// Pop removes and returns the smallest item.
func (pq *PriorityQueue[T]) Pop() (T, error) {
	var zero T
	if len(pq.items) == 0 {
		return zero, ErrInvalidInput
	}
	top := pq.items[0]
	last := pq.items[len(pq.items)-1]
	pq.items = pq.items[:len(pq.items)-1]
	if len(pq.items) > 0 {
		pq.items[0] = last
		pq.siftDown(0)
	}
	return top, nil
}

// Peek returns the smallest item without removing it.
func (pq *PriorityQueue[T]) Peek() (T, error) {
	var zero T
	if len(pq.items) == 0 {
		return zero, ErrInvalidInput
	}
	return pq.items[0], nil
}

// Len returns the number of items.
func (pq *PriorityQueue[T]) Len() int {
	return len(pq.items)
}

func (pq *PriorityQueue[T]) siftUp(i int) {
	for i > 0 {
		parent := (i - 1) / 2
		if pq.items[i] < pq.items[parent] {
			pq.items[i], pq.items[parent] = pq.items[parent], pq.items[i]
			i = parent
		} else {
			break
		}
	}
}

func (pq *PriorityQueue[T]) siftDown(i int) {
	for {
		smallest := i
		left := 2*i + 1
		right := 2*i + 2
		if left < len(pq.items) && pq.items[left] < pq.items[smallest] {
			smallest = left
		}
		if right < len(pq.items) && pq.items[right] < pq.items[smallest] {
			smallest = right
		}
		if smallest == i {
			break
		}
		pq.items[i], pq.items[smallest] = pq.items[smallest], pq.items[i]
		i = smallest
	}
}

// MSTEdge is an undirected weighted edge for minimum spanning tree algorithms.
type MSTEdge struct {
	U      string
	V      string
	Weight float64
}

// MinimumSpanningTree uses Kruskal's algorithm to build a minimum spanning forest.
func MinimumSpanningTree(edges []MSTEdge) (float64, []MSTEdge, error) {
	vertices := make(map[string]struct{})
	for _, e := range edges {
		vertices[e.U] = struct{}{}
		vertices[e.V] = struct{}{}
		if math.IsNaN(e.Weight) || math.IsInf(e.Weight, 0) {
			return 0, nil, ErrInvalidInput
		}
	}

	sorted := make([]MSTEdge, len(edges))
	copy(sorted, edges)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Weight < sorted[j].Weight
	})

	uf := NewUnionFind([]string{})
	for v := range vertices {
		uf.parent[v] = v
		uf.rank[v] = 0
	}

	var total float64
	var mst []MSTEdge
	for _, e := range sorted {
		if e.U == e.V {
			continue
		}
		merged, _ := uf.Union(e.U, e.V)
		if merged {
			total += e.Weight
			mst = append(mst, e)
		}
	}
	return total, mst, nil
}

// FloydWarshall returns the all-pairs shortest-path distance matrix.
func FloydWarshall(graph map[string][]Edge) (map[string]map[string]float64, error) {
	allNodes := make(map[string]struct{})
	for u := range graph {
		allNodes[u] = struct{}{}
	}
	for _, edges := range graph {
		for _, e := range edges {
			allNodes[e.To] = struct{}{}
		}
	}
	for _, edges := range graph {
		for _, e := range edges {
			if _, ok := allNodes[e.To]; !ok {
				return nil, ErrInvalidGraph
			}
			if math.IsNaN(e.Weight) || math.IsInf(e.Weight, 0) {
				return nil, ErrInvalidInput
			}
		}
	}

	dist := make(map[string]map[string]float64)
	for u := range allNodes {
		dist[u] = make(map[string]float64)
		for v := range allNodes {
			if u == v {
				dist[u][v] = 0
			} else {
				dist[u][v] = math.Inf(1)
			}
		}
	}
	for u, edges := range graph {
		for _, e := range edges {
			if e.Weight < dist[u][e.To] {
				dist[u][e.To] = e.Weight
			}
		}
	}

	for k := range allNodes {
		for i := range allNodes {
			dik := dist[i][k]
			if math.IsInf(dik, 1) {
				continue
			}
			for j := range allNodes {
				nd := dik + dist[k][j]
				if nd < dist[i][j] {
					dist[i][j] = nd
				}
			}
		}
	}

	for u := range allNodes {
		if dist[u][u] < 0 {
			return nil, ErrNegativeCycle
		}
	}
	return dist, nil
}

// StronglyConnectedComponents returns the SCCs of the graph using Kosaraju's algorithm.
func StronglyConnectedComponents(graph map[string][]string) ([][]string, error) {
	if err := validateAdjacencyList(graph, "StronglyConnectedComponents"); err != nil {
		return nil, err
	}

	visited := make(map[string]bool)
	order := make([]string, 0)

	type frame struct {
		node    string
		visited bool
	}

	for u := range graph {
		if !visited[u] {
			stack := []frame{{u, false}}
			for len(stack) > 0 {
				f := stack[len(stack)-1]
				stack = stack[:len(stack)-1]
				if f.visited {
					order = append(order, f.node)
					continue
				}
				if visited[f.node] {
					continue
				}
				visited[f.node] = true
				stack = append(stack, frame{f.node, true})
				for i := len(graph[f.node]) - 1; i >= 0; i-- {
					nbr := graph[f.node][i]
					if !visited[nbr] {
						stack = append(stack, frame{nbr, false})
					}
				}
			}
		}
	}

	reverse := make(map[string][]string)
	for u := range graph {
		reverse[u] = nil
	}
	for u, neighbors := range graph {
		for _, v := range neighbors {
			reverse[v] = append(reverse[v], u)
		}
	}

	visited = make(map[string]bool)
	var components [][]string

	for i := len(order) - 1; i >= 0; i-- {
		if !visited[order[i]] {
			component := make([]string, 0)
			stack := []string{order[i]}
			for len(stack) > 0 {
				node := stack[len(stack)-1]
				stack = stack[:len(stack)-1]
				if visited[node] {
					continue
				}
				visited[node] = true
				component = append(component, node)
				for _, nbr := range reverse[node] {
					if !visited[nbr] {
						stack = append(stack, nbr)
					}
				}
			}
			components = append(components, component)
		}
	}
	return components, nil
}

// HasCycle returns whether the directed graph contains a cycle.
func HasCycle(graph map[string][]string) (bool, error) {
	if err := validateAdjacencyList(graph, "HasCycle"); err != nil {
		return false, err
	}
	state := make(map[string]int)
	for u := range graph {
		state[u] = 0
	}
	for u := range graph {
		if state[u] != 0 {
			continue
		}
		type frame struct {
			node string
			idx  int
		}
		stack := []frame{{u, 0}}
		for len(stack) > 0 {
			f := &stack[len(stack)-1]
			if state[f.node] == 2 {
				stack = stack[:len(stack)-1]
				continue
			}
			if state[f.node] == 0 {
				state[f.node] = 1
			}
			neighbors := graph[f.node]
			if f.idx < len(neighbors) {
				nbr := neighbors[f.idx]
				f.idx++
				if state[nbr] == 1 {
					return true, nil
				}
				if state[nbr] == 0 {
					stack = append(stack, frame{nbr, 0})
				}
			} else {
				state[f.node] = 2
				stack = stack[:len(stack)-1]
			}
		}
	}
	return false, nil
}
