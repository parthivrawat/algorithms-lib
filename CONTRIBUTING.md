# Contributing to Algorithms Library

Thanks for your interest in contributing.

## Getting Started

This is a polyglot repository. Each language implementation lives in its own directory and can be built/tested independently:

- **Python**: `pip install -e ".[dev]"` then `pytest`
- **TypeScript**: `npm install` then `npm run lint && npm test && npm run build`
- **Go**: `go test ./...` and `go build`
- **Rust**: `cargo test` and `cargo build`

## CI

Pull requests must pass the full GitHub Actions matrix. The workflow runs `pytest`, `vitest`, `go test`, and `cargo test` plus `gofmt`, `cargo fmt`, `cargo clippy`, and ESLint checks.

## Code Style

- Follow the existing conventions in each language directory.
- Keep changes focused and atomic.
- Add or update tests for bug fixes and new functionality.
- Update `README.md`, `CHANGELOG.md`, and language-specific docs when behavior changes.

## Reporting Issues

Please open an issue with:

- The language/package affected
- A minimal input that reproduces the problem
- The expected and actual output

## License

By contributing, you agree that your contributions will be licensed under the MIT License.
