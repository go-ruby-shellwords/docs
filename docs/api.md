# Usage & API

The public API lives at the module root (`github.com/go-ruby-shellwords/shellwords`). It is **Ruby-shaped but Go-idiomatic**: `Split` / `Escape` / `Join` mirror Shellwords' `shellsplit` / `shellescape` / `shelljoin`, while the surface follows Go conventions — an explicit `error`, value types, no global state.

!!! success "Status: implemented"
    The library is built and importable as `github.com/go-ruby-shellwords/shellwords`, bound into
    `rbgo` as a native module; see [Roadmap](roadmap.md).

## Install

```sh
go get github.com/go-ruby-shellwords/shellwords
```

## Worked example

```go
w, _ := shellwords.Split(`a "b c" d\ e`)    // ["a", "b c", "d e"]
e := shellwords.Escape("a b")               // "a\\ b"
s := shellwords.Join([]string{"a", "b c"})  // "a b\\ c"
```

## Shape

```go
// Split performs POSIX sh word splitting (Shellwords.shellsplit).
func Split(s string) ([]string, error)

// Escape escapes a string so it survives one round of shell parsing
// (Shellwords.shellescape).
func Escape(s string) string

// Join joins words into a single shell-safe command line, the inverse
// of Split (Shellwords.shelljoin).
func Join(words []string) string
```

## MRI conformance

Correctness is defined by reference Ruby. A **differential oracle** runs a wide
corpus through both the system `ruby` and this library and compares the results
**byte-for-byte** — not approximated from memory. The oracle tests skip
themselves where `ruby` is not on `PATH` (e.g. the qemu arch lanes), so the
cross-arch builds still validate the library.

## Relationship to Ruby

`go-ruby-shellwords/shellwords` is **standalone and reusable**, and is the backend bound into
[go-embedded-ruby](https://github.com/go-embedded-ruby/ruby) by `rbgo` as a
native module — the same way [go-ruby-regexp](https://github.com/go-ruby-regexp)
and [go-ruby-erb](https://github.com/go-ruby-erb) are bound. The dependency runs
the other way: this library has no dependency on the Ruby runtime.
