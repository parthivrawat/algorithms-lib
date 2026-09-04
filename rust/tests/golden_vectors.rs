use algorithms_lib::*;
use serde_json::Value;
use std::path::PathBuf;

#[test]
fn golden_vectors() -> Result<(), Box<dyn std::error::Error>> {
    let manifest = PathBuf::from(env!("CARGO_MANIFEST_DIR"));
    let path = manifest.join("../tests/golden_vectors.json");
    let text = std::fs::read_to_string(path)?;
    let data: Value = serde_json::from_str(&text)?;

    for v in data["vectors"].as_array().unwrap() {
        let function = v["function"].as_str().unwrap();
        let input = v["input"].as_object().unwrap();
        let expected = &v["expected"];
        let actual = run(function, input)?;
        assert_eq!(
            actual,
            *expected,
            "{} ({}) mismatch",
            function,
            v["name"].as_str().unwrap()
        );
    }
    Ok(())
}

fn run(
    function: &str,
    input: &serde_json::Map<String, Value>,
) -> Result<Value, Box<dyn std::error::Error>> {
    match function {
        "quick_sort" => {
            let items: Vec<i64> = serde_json::from_value(input["items"].clone())?;
            Ok(serde_json::to_value(quick_sort(&items))?)
        }
        "merge_sort" => {
            let items: Vec<i64> = serde_json::from_value(input["items"].clone())?;
            Ok(serde_json::to_value(merge_sort(&items))?)
        }
        "heap_sort" => {
            let items: Vec<i64> = serde_json::from_value(input["items"].clone())?;
            Ok(serde_json::to_value(heap_sort(&items))?)
        }
        "radix_sort" => {
            let items: Vec<i64> = serde_json::from_value(input["items"].clone())?;
            Ok(serde_json::to_value(radix_sort(&items)?)?)
        }
        "native_sort" => {
            let items: Vec<i64> = serde_json::from_value(input["items"].clone())?;
            Ok(serde_json::to_value(native_sort(&items))?)
        }
        "binary_search" => {
            let arr: Vec<i64> = serde_json::from_value(input["arr"].clone())?;
            let target: i64 = serde_json::from_value(input["target"].clone())?;
            Ok(serde_json::to_value(binary_search(&arr, &target)?)?)
        }
        "interpolation_search" => {
            let arr: Vec<i64> = serde_json::from_value(input["arr"].clone())?;
            let target: i64 = serde_json::from_value(input["target"].clone())?;
            Ok(serde_json::to_value(interpolation_search(&arr, target)?)?)
        }
        "jump_search" => {
            let arr: Vec<i64> = serde_json::from_value(input["arr"].clone())?;
            let target: i64 = serde_json::from_value(input["target"].clone())?;
            Ok(serde_json::to_value(jump_search(&arr, &target)?)?)
        }
        "kmp_search" => {
            let text: String = serde_json::from_value(input["text"].clone())?;
            let pattern: String = serde_json::from_value(input["pattern"].clone())?;
            Ok(serde_json::to_value(kmp_search(&text, &pattern)?)?)
        }
        "rabin_karp_search" => {
            let text: String = serde_json::from_value(input["text"].clone())?;
            let pattern: String = serde_json::from_value(input["pattern"].clone())?;
            let base: i64 = serde_json::from_value(input["base"].clone())?;
            let modulus: i64 = serde_json::from_value(input["mod"].clone())?;
            Ok(serde_json::to_value(rabin_karp_search(
                &text, &pattern, base, modulus,
            )?)?)
        }
        "boyer_moore_search" => {
            let text: String = serde_json::from_value(input["text"].clone())?;
            let pattern: String = serde_json::from_value(input["pattern"].clone())?;
            Ok(serde_json::to_value(boyer_moore_search(&text, &pattern)?)?)
        }
        "longest_common_subsequence" => {
            let a: String = serde_json::from_value(input["a"].clone())?;
            let b: String = serde_json::from_value(input["b"].clone())?;
            let a: Vec<char> = a.chars().collect();
            let b: Vec<char> = b.chars().collect();
            Ok(serde_json::to_value(longest_common_subsequence(&a, &b))?)
        }
        "edit_distance" => {
            let a: String = serde_json::from_value(input["a"].clone())?;
            let b: String = serde_json::from_value(input["b"].clone())?;
            let a: Vec<char> = a.chars().collect();
            let b: Vec<char> = b.chars().collect();
            Ok(serde_json::to_value(edit_distance(&a, &b))?)
        }
        "knapsack_01" => {
            let weights: Vec<i64> = serde_json::from_value(input["weights"].clone())?;
            let values: Vec<i64> = serde_json::from_value(input["values"].clone())?;
            let capacity: i64 = serde_json::from_value(input["capacity"].clone())?;
            Ok(serde_json::to_value(knapsack_01(
                &weights, &values, capacity,
            )?)?)
        }
        "fractional_knapsack" => {
            let weights: Vec<f64> = serde_json::from_value(input["weights"].clone())?;
            let values: Vec<f64> = serde_json::from_value(input["values"].clone())?;
            let capacity: f64 = serde_json::from_value(input["capacity"].clone())?;
            Ok(serde_json::to_value(fractional_knapsack(
                &weights, &values, capacity,
            )?)?)
        }
        "activity_selection" => {
            let activities: Vec<(i64, i64)> = serde_json::from_value(input["activities"].clone())?;
            Ok(serde_json::to_value(activity_selection(&activities)?)?)
        }
        "max_subarray" => {
            let arr: Vec<i64> = serde_json::from_value(input["arr"].clone())?;
            Ok(serde_json::to_value(max_subarray(&arr)?)?)
        }
        "count_inversions" => {
            let arr: Vec<i64> = serde_json::from_value(input["arr"].clone())?;
            Ok(serde_json::to_value(count_inversions(&arr))?)
        }
        "fast_power" => {
            let base: f64 = serde_json::from_value(input["base"].clone())?;
            let exponent: i64 = serde_json::from_value(input["exponent"].clone())?;
            Ok(serde_json::to_value(fast_power(base, exponent)?)?)
        }
        _ => panic!("unknown golden vector function: {}", function),
    }
}
