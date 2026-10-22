# overview

A Go CLI that lists the repository's files, one full path per line, each with its header summary. Run `overview` for the map of this repo. It only reads: a check of file headers belongs in slopfix, not here. Never write a hand-kept layout into a doc.

- Build and test with `go-toolchain`, not bare `go`. It also runs `dats/cli.dats` against the built binary.
- `internal/header` owns the comment syntaxes, one style per file type. A new type is one map entry plus a case in `header_test.go`.
- `internal/listing` is pure output. The output stays flat on purpose: a tree costs as many tokens and is harder to read back. `cmd/` holds one cobra command per file, and each command registers itself in `init()`.
