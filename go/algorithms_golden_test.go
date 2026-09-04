package algorithms

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"testing"
)

func cast(from, to interface{}) error {
	b, err := json.Marshal(from)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, to)
}

type goldenCase struct {
	Function string                 `json:"function"`
	Name     string                 `json:"name"`
	Input    map[string]interface{} `json:"input"`
	Expected interface{}            `json:"expected"`
}

type goldenFile struct {
	Vectors []goldenCase `json:"vectors"`
}

func TestGoldenVectors(t *testing.T) {
	f, err := os.Open("../tests/golden_vectors.json")
	if err != nil {
		t.Fatalf("open golden vectors: %v", err)
	}
	defer f.Close()

	var file goldenFile
	if err := json.NewDecoder(f).Decode(&file); err != nil {
		t.Fatalf("decode golden vectors: %v", err)
	}

	for _, c := range file.Vectors {
		c := c
		t.Run(fmt.Sprintf("%s-%s", c.Function, c.Name), func(t *testing.T) {
			var actual interface{}
			var err error

			switch c.Function {
			case "quick_sort":
				var items []int
				if err = cast(c.Input["items"], &items); err != nil {
					break
				}
				actual = QuickSort(items)
			case "merge_sort":
				var items []int
				if err = cast(c.Input["items"], &items); err != nil {
					break
				}
				actual = MergeSort(items)
			case "heap_sort":
				var items []int
				if err = cast(c.Input["items"], &items); err != nil {
					break
				}
				actual = HeapSort(items)
			case "radix_sort":
				var items []int
				if err = cast(c.Input["items"], &items); err != nil {
					break
				}
				actual, err = RadixSort(items)
			case "native_sort":
				var items []int
				if err = cast(c.Input["items"], &items); err != nil {
					break
				}
				actual = NativeSort(items)
			case "binary_search":
				var arr []int
				var target int
				if err = cast(c.Input["arr"], &arr); err != nil {
					break
				}
				if err = cast(c.Input["target"], &target); err != nil {
					break
				}
				actual, err = BinarySearch(arr, target)
			case "interpolation_search":
				var arr []int
				var target int
				if err = cast(c.Input["arr"], &arr); err != nil {
					break
				}
				if err = cast(c.Input["target"], &target); err != nil {
					break
				}
				actual, err = InterpolationSearch(arr, target)
			case "jump_search":
				var arr []int
				var target int
				if err = cast(c.Input["arr"], &arr); err != nil {
					break
				}
				if err = cast(c.Input["target"], &target); err != nil {
					break
				}
				actual, err = JumpSearch(arr, target)
			case "kmp_search":
				var text, pattern string
				if err = cast(c.Input["text"], &text); err != nil {
					break
				}
				if err = cast(c.Input["pattern"], &pattern); err != nil {
					break
				}
				actual, err = KMPSearch(text, pattern)
			case "rabin_karp_search":
				var text, pattern string
				var base int
				var mod int64
				if err = cast(c.Input["text"], &text); err != nil {
					break
				}
				if err = cast(c.Input["pattern"], &pattern); err != nil {
					break
				}
				if err = cast(c.Input["base"], &base); err != nil {
					break
				}
				if err = cast(c.Input["mod"], &mod); err != nil {
					break
				}
				actual, err = RabinKarpSearch(text, pattern, base, mod)
			case "boyer_moore_search":
				var text, pattern string
				if err = cast(c.Input["text"], &text); err != nil {
					break
				}
				if err = cast(c.Input["pattern"], &pattern); err != nil {
					break
				}
				actual, err = BoyerMooreSearch(text, pattern)
			case "longest_common_subsequence":
				var aStr, bStr string
				if err = cast(c.Input["a"], &aStr); err != nil {
					break
				}
				if err = cast(c.Input["b"], &bStr); err != nil {
					break
				}
				a, b := []rune(aStr), []rune(bStr)
				actual = LongestCommonSubsequence(a, b)
			case "edit_distance":
				var aStr, bStr string
				if err = cast(c.Input["a"], &aStr); err != nil {
					break
				}
				if err = cast(c.Input["b"], &bStr); err != nil {
					break
				}
				a, b := []rune(aStr), []rune(bStr)
				actual = EditDistance(a, b)
			case "knapsack_01":
				var weights, values []int
				var capacity int
				if err = cast(c.Input["weights"], &weights); err != nil {
					break
				}
				if err = cast(c.Input["values"], &values); err != nil {
					break
				}
				if err = cast(c.Input["capacity"], &capacity); err != nil {
					break
				}
				actual, err = Knapsack01(weights, values, capacity)
			case "fractional_knapsack":
				var weights, values []float64
				var capacity float64
				if err = cast(c.Input["weights"], &weights); err != nil {
					break
				}
				if err = cast(c.Input["values"], &values); err != nil {
					break
				}
				if err = cast(c.Input["capacity"], &capacity); err != nil {
					break
				}
				actual, err = FractionalKnapsack(weights, values, capacity)
			case "activity_selection":
				var activities [][2]int
				if err = cast(c.Input["activities"], &activities); err != nil {
					break
				}
				actual, err = ActivitySelection(activities)
			case "max_subarray":
				var arr []int
				if err = cast(c.Input["arr"], &arr); err != nil {
					break
				}
				actual, err = MaxSubarray(arr)
			case "count_inversions":
				var arr []int
				if err = cast(c.Input["arr"], &arr); err != nil {
					break
				}
				actual = CountInversions(arr)
			case "fast_power":
				var base float64
				var exponent int
				if err = cast(c.Input["base"], &base); err != nil {
					break
				}
				if err = cast(c.Input["exponent"], &exponent); err != nil {
					break
				}
				actual, err = FastPower(base, exponent)
			default:
				t.Fatalf("unknown golden vector function: %s", c.Function)
			}

			if err != nil {
				t.Fatalf("%s error: %v", c.Function, err)
			}

			gotBytes, err := json.Marshal(actual)
			if err != nil {
				t.Fatalf("marshal actual: %v", err)
			}
			wantBytes, err := json.Marshal(c.Expected)
			if err != nil {
				t.Fatalf("marshal expected: %v", err)
			}
			if !bytes.Equal(gotBytes, wantBytes) {
				t.Errorf("%s = %s, want %s", c.Function, gotBytes, wantBytes)
			}
		})
	}
}
