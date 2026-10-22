# overview

`overview` prints a repository's files as a tree. Beside each file it prints the first sentence of that file's leading comment. The map comes from the source itself, so nobody has to keep a hand-written layout in a doc.

```
├── cmd/
│   └── root.go  ← Package cmd holds the overview command line.
├── internal/
│   ├── header/
│   │   └── header.go  ← Package header reads the one-sentence summary from a file's leading comment.
│   └── tree/
│       └── tree.go  ← Package tree draws file paths and their summaries as a box-drawing tree.
└── main.go  ← Command overview prints a repository's file tree with each file's header summary.
```

## Usage

Run it inside a git work tree. It lists what `git ls-files` lists, so ignored files never appear.

```sh
overview                 # the whole tree
overview src/ui          # one subtree
overview --depth 2       # fold deeper directories into a file count
overview check           # exit 1 when a source file has no header comment
overview check --exclude '*.test.ts' --exclude vendor/
```

`overview check` is the CI half: it keeps the map complete. A generated file (`Code generated ... DO NOT EDIT.`) is never flagged.

## Install

```sh
go install github.com/wow-look-at-my/overview@latest
```

## Build

`go-toolchain` builds and tests it, and runs the end-to-end suite in `dats/`.

MIT licensed.
