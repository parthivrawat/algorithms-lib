# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [2.0.0] - 2026-09-04

### Added

- New algorithms and data structures in every package: union-find, topological sort, strongly connected components, minimum spanning tree, Floyd-Warshall, binary search on answer, quickselect, modular fast power, gcd/lcm, sieve of Eratosthenes, Huffman encode/decode, edit distance / LCS reconstruction, priority queue, and graph cycle detection.

### Changed

- **Breaking**: `tim_sort` renamed to `native_sort` across all four packages and documented as a built-in sort wrapper.
- **Breaking**: Go `FastPower` and Rust `fast_power` now return an `error` / `Result` and validate edge cases including `0` raised to a negative power and `math.MinInt` / `i64::MIN` exponents.
- **Breaking**: Go `CountInversions` now returns `int64`.
- **Breaking**: Rust `binary_search`, `interpolation_search`, and `jump_search` use `debug_assert!` for sortedness; release builds treat unsorted input as undefined behavior to preserve `O(log n)`.
- TypeScript sort and search APIs now accept an optional `Comparator<T>` for custom ordering.

### Fixed

- Cross-language correctness: `boyer_moore_search` overlapping matches, `fast_power` edge cases, `fractional_knapsack` zero / negative weight handling, `radix_sort` non-integer / negative / overflow validation, `interpolation_search` precision and overflow, `a_star` non-negative heuristic and stale heap entries, `rabin_karp_search` base / mod validation, and `dijkstra` negative-weight pre-scanning.
- Go: `JumpSearch` empty / nil slice panic, `HuffmanCoding` empty-string key, `BinarySearch` / `MaxSubarray` midpoint overflow, and `requireSorted` error message.
- Rust: `kmp_search` `usize` underflow, `knapsack_01` `usize` truncation, and `NaN` / `inf` graph-weight rejection.
- TypeScript: `bfs` `O(n)` `shift()` queue, `bellmanFord` unreachable-vertex initialization, and generic `Error` usage replaced with `AlgorithmsError` subclasses.
- Python: `quick_sort` space-complexity docstring, `knapsack_01` float capacity, `huffman_coding` non-positive frequencies, and `bfs` / `dfs` missing-neighbor validation.

## [1.0.0] - 2026-09-04

### Added

- Polyglot implementations of common algorithms in Python, TypeScript, Go, and Rust.
- Shared golden test vectors under `tests/golden_vectors.json`.
- GitHub Actions CI matrix covering all four language test suites.
- Root `LICENSE`, `CONTRIBUTING.md`, and `.gitignore`.
- TypeScript `sideEffects: false` for improved tree-shaking.

### Changed

- Migrated Python packaging metadata to PEP 621 (`pyproject.toml`) and raised `requires-python` to `>=3.9`.
- Bumped TypeScript `engines` to `node >=18` and added an ESLint configuration.
- Rust `Cargo.toml` now declares `rust-version` and an explicit `include` allowlist.
