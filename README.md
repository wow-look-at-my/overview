# overview

`overview` lists a repository's files, one full path per line. Beside each file it prints the first sentence of that file's leading comment. The map comes from the source itself, so nobody has to keep a hand-written layout in a doc.

```
$ overview cmd internal/header/header.go internal/listing/listing.go main.go
cmd/files.go  ← File listing through git, and the header read of each file.
cmd/root.go  ← Package cmd holds the overview command line.
internal/header/header.go  ← Package header reads the one-sentence summary from a file's leading comment.
internal/listing/listing.go  ← Package listing writes file paths and their summaries, one full path per line.
main.go  ← Command overview lists a repository's files, each with its header summary.
```

Full paths, not a tree: each line stands alone. A reader can pass it straight to a tool, and `grep` or `head` never cuts a line off from its parent directory. A tree costs as many tokens, because the box-drawing characters cost what the repeated prefixes save.

## Usage

Run it inside a git work tree. It lists what `git ls-files` lists, so ignored files never appear.

```sh
overview                 # every file
overview src/ui          # the files under one directory
```

overview only reads. [slopfix](https://github.com/wow-look-at-my/slopfix) checks that every source file carries a header comment. The map stays complete.

## Install

```sh
go install github.com/wow-look-at-my/overview@latest
```

## Build

`go-toolchain` builds and tests it, and runs the end-to-end suite in `dats/`.

MIT licensed.
