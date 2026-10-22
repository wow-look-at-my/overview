# overview

A Go CLI that prints the repository tree with each file's header-comment summary, plus `overview check` to fail on a missing header. Run `overview` for the map of this repo. Never write a hand-kept layout into a doc.

- Build and test with `go-toolchain`, not bare `go`. It also runs `dats/cli.dats` against the built binary.
- `internal/header` owns the comment syntaxes, one style per file type. A new type is one map entry plus a case in `header_test.go`.
- `internal/tree` is pure rendering. `cmd/` holds one cobra command per file, and each command registers itself in `init()`.
