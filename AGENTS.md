# AGENTS.md
Repository guidance for coding agents working in `github.com/colt3k/nglog`.

## Scope and hierarchy
- This file applies to the whole repository.
- There are additional scoped guides in:
  - `ng/AGENTS.md`
  - `internal/pkg/AGENTS.md`
- When working inside those directories, follow the more specific guide too.

## Repository shape
- Language: Go
- Module: `github.com/colt3k/nglog`
- Go version from `go.mod`: `go 1.23.0`
- Toolchain from `go.mod`: `go1.23.8`
- Main areas:
  - `ng/` = core logging implementation
  - `ers/` = error helpers
  - `internal/pkg/` = enum, scheduler, utility helpers
  - `test/` = example/integration-style tests

## Editor / assistant rule files
- No `.cursorrules` file was found.
- No `.cursor/rules/` directory was found.
- No `.github/copilot-instructions.md` file was found.
- The only existing agent-specific docs are `ng/AGENTS.md` and `internal/pkg/AGENTS.md`.

## Working norms
- Prefer minimal, compatibility-preserving edits.
- This is a library repo; avoid renaming public types, functions, constants, or packages unless the task explicitly requires it.
- Do not “modernize” unrelated code while fixing a bug.
- Match local file style in older files, then run `gofmt` on edited files.
- Be careful with tests: several are noisy, side-effectful, or environment-dependent.

## Build, test, lint, and verification commands

## Vendoring caveat
- The checked-in `vendor/` directory is inconsistent with `go.mod`.
- Plain `go test ./...` and `go build ./...` fail with an inconsistent vendoring error.
- Prefer `-mod=readonly` for routine agent work unless the task is specifically about vendoring.
- Do not run `go mod vendor` unless the user asks for vendored dependency updates.

## Primary commands
- Build all packages:
  - `go build -mod=readonly ./...`
- Run all tests:
  - `go test -mod=readonly ./...`
- Run vet across the repo:
  - `go vet -mod=readonly ./...`
- Format edited files:
  - `gofmt -w path/to/file.go`

## Single-test and targeted test commands
- Run one package:
  - `go test -mod=readonly ./internal/pkg/sched`
  - `go test -mod=readonly ./ers/bserr`
  - `go test -mod=readonly ./test`
- Run one named test:
  - `go test -mod=readonly ./internal/pkg/sched -run '^TestCronSched_Parse$'`
  - `go test -mod=readonly ./test -run '^TestJSON$'`
  - `go test -mod=readonly ./ers/bserr -run '^TestBSErr_Err$'`
- Run one test verbosely:
  - `go test -mod=readonly ./test -run '^TestJSON$' -v`

## Current verification reality
- `go test -mod=readonly ./internal/pkg/sched -run '^TestCronSched_Parse$'` passes.
- `go test -mod=readonly ./test -run '^TestJSON$'` passes.
- `go build -mod=readonly ./...` and `go vet -mod=readonly ./...` run without the vendoring error.
- `go test -mod=readonly ./...` is not clean today:
  - `ers/bserr` intentionally exercises fatal/error paths,
  - the full run emits a large amount of log output,
  - some tests write files or interact with local services.

## Tests with side effects or environment assumptions
- `test/log_test.go` writes files such as:
  - `test.log`
  - `myfile.log`
  - `output.txt`
  - `logtest/roll_test.log*`
- `TestHTTP` expects an HTTP listener on `localhost:8080`.
- `TestSyslog` and `TestSyslogXML` interact with syslog.
- `TestMail`, `TestTCPSocket`, and `TestErrStacks` are currently skipped.
- Prefer the narrowest possible test scope while iterating.

## Where to look first
- Public root interface: `logger.go`
- Core logger implementation: `ng/logger.go`
- Layout interface and text formatting: `ng/textlayout.go`
- File appender: `ng/appender_file.go`
- HTTP appender: `ng/appender_http.go`
- Error helpers: `ers/bserr/errs.go`
- Scheduler utilities: `internal/pkg/sched/sched.go`
- Representative tests/examples: `test/log_test.go`

## Code style and implementation guidance

## Formatting
- Always run `gofmt` on changed Go files.
- The repo has older formatting habits in some files, but `gofmt` remains the safest baseline.
- Keep diffs small; do not reformat untouched files for cosmetic consistency.

## Imports
- Use normal Go import blocks.
- The repo is mixed: some files separate stdlib and project imports cleanly, some older files are less tidy.
- Preferred approach for new or edited code:
  - stdlib imports first,
  - blank line,
  - third-party / module imports,
  - alias imports only when needed for clarity or collision avoidance.
- The `ng` package is often imported as `log`:
  - `log "github.com/colt3k/nglog/ng"`

## Naming
- Packages are lowercase short names: `ng`, `ers`, `util`, `sched`.
- Exported identifiers use PascalCase.
- Unexported fields and helpers use lowerCamelCase.
- Existing public constants use all-caps enum-like names such as `INFO`, `WARN`, `ERROR`, `DBGL2`, `DBGL3`, `FATALNE`.
- Preserve historical public names even when they are awkward by modern standards.

## Types and APIs
- Prefer existing concrete types and interfaces over inventing parallel abstractions.
- Reuse current interfaces where appropriate:
  - `Logger` in `logger.go`
  - `Layout` in `ng/textlayout.go`
  - appender interfaces in `ng/`
- Avoid unnecessary API churn in exported packages.
- `ng/AGENTS.md` says public types/functions must be documented; follow that rule for new exported APIs.

## Error handling
- Handle errors explicitly and early.
- Existing code commonly returns `error` from constructors and setup functions.
- Common patterns already in the repo:
  - `return nil, err`
  - `return nil, fmt.Errorf(...)`
  - logging through `ng` helpers in error-wrapper code
- Some legacy files print directly with `fmt.Println`/`fmt.Printf` or call fatal logging helpers; do not spread those patterns into new code unless you are preserving behavior in that exact subsystem.
- Prefer returning errors from library code over introducing new process-exit behavior.

## Logging patterns
- The core abstraction is the `ng` logger and appender/layout system.
- When modifying logging internals, preserve behavior before attempting cleanup.
- Common calls include:
  - `log.Modify(...)`
  - `log.Logln(...)`
  - `log.Logf(...)`
  - `log.WithFields(...)`
- Keep level semantics intact; `DBGL2`, `DBGL3`, and `FATALNE` are intentional public surface area.

## Tests
- Tests mostly use the standard `testing` package directly.
- Assertions are lightweight; many tests just execute code paths.
- Keep new tests close to the package they cover.
- Prefer narrow, deterministic tests over new environment-dependent tests.
- If adding tests around I/O appenders, avoid unnecessary network or system dependencies.

## Compatibility guidance
- `internal/pkg/AGENTS.md` notes that internal package changes must maintain backward compatibility.
- Treat exported behavior in `ng/` and `ers/` as stable unless the task explicitly allows breakage.
- Be careful with constructor signatures, exported constants, and formatter output.

## Practical agent checklist
1. Read this file and any deeper `AGENTS.md` in the target directory.
2. Inspect the exact package you are changing before editing.
3. Use `-mod=readonly` for build/test/vet unless vendoring is the task.
4. Prefer package-level or single-test runs during iteration.
5. Expect noisy output and side effects from `test/log_test.go`.
6. Keep public APIs stable.
7. Run `gofmt -w` on edited Go files.
8. Re-run the narrowest relevant verification command before finishing.
