// Package algorithms provides a comprehensive, zero-dependency collection
// of common algorithms for Go.
package algorithms

import (
	"container/heap"
	"errors"
	"math"
	"sort"
)

// Ordered is a constraint for types that support ordering operators.
type Ordered interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~float32 | ~float64 |
		~string
}

// Common errors returned by this package.
var (
	ErrNotFound     = errors.New("not found")
	ErrInvalidInput = errors.New("invalid input")
	ErrInvalidGraph = errors.New("invalid graph")
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
	for _, v := range items {
		if v < 0 {
			return nil, ErrInvalidInput
		}
	}
	arr := make([]int, len(items))
	copy(arr, items)
	maxVal := arr[0]
	for _, v := range arr[1:] {
		if v > maxVal {
			maxVal = v
		}
	}
	for exp := 1; maxVal/exp > 0; exp *= 10 {
		countingSortByDigit(arr, exp)
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

// TimSort returns a new sorted slice using the standard library sort.
func TimSort[T Ordered](items []T) []T {
	arr := make([]T, len(items))
	copy(arr, items)
	sort.Slice(arr, func(i, j int) bool { return arr[i] < arr[j] })
	return arr
}

// ---------------- Searching ----------------

func requireSorted[T Ordered](arr []T, name string) error {
	for i := 1; i < len(arr); i++ {
		if arr[i] < arr[i-1] {
			return ErrInvalidInput
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
		mid := (lo + hi) / 2
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
		pos := lo + int(float64(target-arr[lo])/float64(arr[hi]-arr[lo])*float64(hi-lo))
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

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// ---------------- Graphs ----------------

// Edge represents a weighted directed edge in a graph.
type Edge struct {
	To     string
	Weight float64
}

// BFS returns the vertices in breadth-first order.
func BFS(graph map[string][]string, start string) ([]string, error) {
	if _, ok := graph[start]; !ok {
		return nil, ErrInvalidGraph
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
func DFS(graph map[string][]string, start string) ([]string, error) {
	if _, ok := graph[start]; !ok {
		return nil, ErrInvalidGraph
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
	*pq = old[:n-1]
	return item
}

// Dijkstra computes shortest distances from the start node.
// Edge weights must be non-negative.
func Dijkstra(graph map[string][]Edge, start string) (map[string]float64, map[string]string, error) {
	if _, ok := graph[start]; !ok {
		return nil, nil, ErrInvalidGraph
	}
	distances := make(map[string]float64)
	for k := range graph {
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
		if dist != distances[node] {
			continue
		}
		for _, e := range graph[node] {
			if e.Weight < 0 {
				return nil, nil, ErrInvalidGraph
			}
			if _, ok := distances[e.To]; !ok {
				distances[e.To] = math.Inf(1)
			}
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
func AStar(graph map[string][]Edge, start, goal string, heuristic func(string, string) float64) ([]string, float64, error) {
	if _, ok := graph[start]; !ok {
		return nil, 0, ErrInvalidGraph
	}
	if _, ok := graph[goal]; !ok {
		return nil, 0, ErrInvalidGraph
	}
	gScore := map[string]float64{start: 0}
	cameFrom := make(map[string]string)
	pq := &graphPQ{}
	heap.Init(pq)
	heap.Push(pq, &pqItem{node: start, priority: heuristic(start, goal)})

	for pq.Len() > 0 {
		item := heap.Pop(pq).(*pqItem)
		node := item.node
		if node == goal {
			return reconstructPath(cameFrom, node), gScore[goal], nil
		}
		for _, e := range graph[node] {
			tentative := gScore[node] + e.Weight
			if g, ok := gScore[e.To]; !ok || tentative < g {
				cameFrom[e.To] = node
				gScore[e.To] = tentative
				f := tentative + heuristic(e.To, goal)
				heap.Push(pq, &pqItem{node: e.To, priority: f})
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

// RabinKarpSearch returns all starting indices of the pattern using rolling hash.
func RabinKarpSearch(text, pattern string, base int, mod int64) ([]int, error) {
	if pattern == "" {
		return nil, ErrInvalidInput
	}
	n := len(text)
	m := len(pattern)
	if m > n {
		return []int{}, nil
	}
	h := int64(1)
	for i := 0; i < m-1; i++ {
		h = (h * int64(base)) % mod
	}
	var patternHash, textHash int64
	for i := 0; i < m; i++ {
		patternHash = (patternHash*int64(base) + int64(pattern[i])) % mod
		textHash = (textHash*int64(base) + int64(text[i])) % mod
	}
	matches := []int{}
	for i := 0; i <= n-m; i++ {
		if textHash == patternHash && text[i:i+m] == pattern {
			matches = append(matches, i)
		}
		if i < n-m {
			textHash = (textHash - int64(text[i])*h) % mod
			textHash = (textHash*int64(base) + int64(text[i+m])) % mod
			textHash = (textHash%mod + mod) % mod
		}
	}
	return matches, nil
}

// BoyerMooreSearch returns all starting indices of the pattern in the text.
func BoyerMooreSearch(text, pattern string) ([]int, error) {
	if pattern == "" {
		return nil, ErrInvalidInput
	}
	runes := []rune(text)
	patRunes := []rune(pattern)
	n := len(runes)
	m := len(patRunes)
	if m > n {
		return []int{}, nil
	}
	badChar := make(map[rune]int)
	for i, r := range patRunes {
		badChar[r] = i
	}
	matches := []int{}
	i := 0
	for i <= n-m {
		j := m - 1
		for j >= 0 && runes[i+j] == patRunes[j] {
			j--
		}
		if j < 0 {
			matches = append(matches, i)
			i += m
		} else {
			last, ok := badChar[runes[i+j]]
			if !ok {
				last = -1
			}
			shift := j - last
			if shift < 1 {
				shift = 1
			}
			i += shift
		}
	}
	return matches, nil
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

// FractionalKnapsack returns the maximum value for the fractional knapsack problem.
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
		value  float64
	}
	items := make([]item, 0, len(weights))
	for i := 0; i < len(weights); i++ {
		if weights[i] > 0 {
			items = append(items, item{
				ratio:  values[i] / weights[i],
				weight: weights[i],
				value:  values[i],
			})
		}
	}
	sort.Slice(items, func(i, j int) bool { return items[i].ratio > items[j].ratio })
	total := 0.0
	remaining := capacity
	for _, it := range items {
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
	freq       int
	char       string
	left       *huffmanNode
	right      *huffmanNode
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
	*pq = old[:n-1]
	return node
}

// HuffmanCoding returns a prefix-free code table for the given frequencies.
func HuffmanCoding(frequencies map[string]int) (map[string]string, error) {
	if len(frequencies) == 0 {
		return nil, ErrInvalidInput
	}
	pq := &huffmanPQ{}
	heap.Init(pq)
	for char, freq := range frequencies {
		heap.Push(pq, &huffmanNode{freq: freq, char: char})
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
	if node.char != "" {
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
	mid := (left + right) / 2
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
func CountInversions[T Ordered](arr []T) int {
	_, count := sortAndCount(arr)
	return count
}

func sortAndCount[T Ordered](arr []T) ([]T, int) {
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

func mergeAndCount[T Ordered](left, right []T) ([]T, int) {
	merged := make([]T, 0, len(left)+len(right))
	count := 0
	i, j := 0, 0
	for i < len(left) && j < len(right) {
		if left[i] <= right[j] {
			merged = append(merged, left[i])
			i++
		} else {
			merged = append(merged, right[j])
			j++
			count += len(left) - i
		}
	}
	merged = append(merged, left[i:]...)
	merged = append(merged, right[j:]...)
	return merged, count
}

// FastPower returns base raised to the exponent using exponentiation by squaring.
func FastPower(base float64, exponent int) float64 {
	if exponent < 0 {
		return 1 / FastPower(base, -exponent)
	}
	if exponent == 0 {
		return 1
	}
	if exponent == 1 {
		return base
	}
	half := FastPower(base, exponent/2)
	if exponent%2 == 0 {
		return half * half
	}
	return base * half * half
}
