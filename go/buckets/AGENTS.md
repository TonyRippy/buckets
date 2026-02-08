# Repository Guidelines

## Project Structure & Module Organization
This repository is a single Go module (`go.mod`, module name `buckets`) focused on bucketing strategies.

- Core interfaces and shared types: `buckets.go` (`BucketingStrategy`, `Range`, bounds logic).
- Strategy implementation: `floor.go` (current `floor` bucketer and parser registration).
- Parser/registry: `parse.go` (`RegisterParser`, `Parse`).
- Tests: `*_test.go` files alongside implementation (`buckets_test.go`, `floor_test.go`, `parse_test.go`).

Keep new strategy code in its own `<strategy>.go` file with matching `<strategy>_test.go`.

## Build, Test, and Development Commands
- `go test ./...` runs the full test suite.
- `go test -run TestFloorBucketer ./...` runs targeted tests while iterating.
- `go test -cover ./...` checks package-level coverage.
- `go fmt ./...` formats all Go files.
- `go vet ./...` catches common static issues.

Run `go fmt` and `go test ./...` before opening a PR.

## Coding Style & Naming Conventions
- Follow standard Go formatting (`gofmt`), tabs for indentation, and idiomatic Go naming.
- Public API names use `CamelCase` (for example `FloorBucketer`, `RegisterParser`).
- Unexported internals use `camelCase` (for example `floorBucketer`, `assertRangeEquals`).
- Keep parser keys lowercase (current code normalizes names/args with `strings.ToLower`).
- Prefer clear error messages with context, e.g. `fmt.Errorf("invalid width %g", width)`.

## Testing Guidelines
- Use Go’s built-in `testing` package with table-driven tests.
- Name tests `TestXxx` and subtests with `t.Run(...)` using descriptive case labels.
- Cover both success and error paths (`Parse` unknown/empty input, invalid widths, boundary ranges).
- Keep helper assertions in test files and mark helpers with `t.Helper()`.

## Commit & Pull Request Guidelines
Git history is minimal; follow the observed style: short, imperative, sentence-case summaries (for example, `Add parser validation for duplicate keys`).

- Keep commits focused to one logical change.
- PRs should include:
  - What changed and why.
  - Test evidence (`go test ./...` output summary).
  - Any API or behavior changes with examples (input spec and expected bucket/range).
