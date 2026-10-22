# overview

`overview` lists a repository's files, one full path per line. Beside each file it prints the first sentence of that file's leading comment. The map comes from the source itself, so nobody has to keep a hand-written layout in a doc.

EXAMPLE_PLACEHOLDER

Full paths, not a tree: each line stands alone. A reader can pass it straight to a tool, and `grep` or `head` never cuts a line off from its parent directory. A tree costs as many tokens, because the box-drawing characters cost what the repeated prefixes save.

## Usage

Run it inside a git work tree. It lists what `git ls-files` lists, so ignored files never appear.

```sh
overview                 # every file
overview src/ui          # the files under one directory
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
