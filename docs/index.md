# go-ruby-shellwords documentation

**Ruby's Shellwords split, escape & join in pure Go — MRI-compatible, no cgo.**

`go-ruby-shellwords/shellwords` is a faithful, pure-Go (zero cgo) reimplementation of Ruby's Shellwords,
matching reference Ruby (MRI) byte-for-byte. The module path is
`github.com/go-ruby-shellwords/shellwords`.

It was **extracted from rbgo's prelude/internals into a reusable standalone
library**: the module is standalone and importable by any Go program, and it is
the backend bound into [go-embedded-ruby](https://github.com/go-embedded-ruby/ruby)
by `rbgo` as a native module — just like
[go-ruby-regexp](https://github.com/go-ruby-regexp) and
[go-ruby-erb](https://github.com/go-ruby-erb). The dependency runs the other
way: this library has **no dependency on the Ruby runtime**.

!!! success "Status: shellsplit + shellescape + shelljoin complete — MRI byte-exact"
    Faithful port of Ruby's Shellwords: **`shellsplit`** POSIX `sh` word splitting honouring single quotes, double quotes and backslash escapes, **`shellescape`** so a string survives one round of shell parsing (including the empty-string `''` case), and **`shelljoin`** as the inverse — with `shellsplit(shelljoin(words)) == words` guaranteed and the same **`ArgumentError`** on an unmatched quote. Validated by a **differential oracle** against the system `ruby` / `shellwords` — split, escaped and joined output compared byte-for-byte — at 100% coverage, `gofmt` + `go vet` clean, CI green across the six 64-bit Go targets and three OSes.

## Quick taste

```go
w, _ := shellwords.Split(`a "b c" d\ e`)    // ["a", "b c", "d e"]
e := shellwords.Escape("a b")               // "a\\ b"
s := shellwords.Join([]string{"a", "b c"})  // "a b\\ c"
```

## Repositories

| Repo | What it is |
| --- | --- |
| [`shellwords`](https://github.com/go-ruby-shellwords/shellwords) | the library — Ruby's Shellwords in pure Go |
| [`docs`](https://github.com/go-ruby-shellwords/docs) | this documentation site (MkDocs Material, versioned with mike) |
| [`go-ruby-shellwords.github.io`](https://github.com/go-ruby-shellwords/go-ruby-shellwords.github.io) | the organization landing page (Hugo) |
| [`brand`](https://github.com/go-ruby-shellwords/brand) | logo and brand assets |

## Principles

- **Pure Go, `CGO_ENABLED=0`** — trivial cross-compilation, a single static
  binary, no C toolchain.
- **MRI byte-exact.** Output matches reference Ruby exactly, not approximately,
  validated by a differential oracle against the `ruby` binary.
- **Standalone & reusable.** Extracted from rbgo's internals; no dependency on
  the Ruby runtime — the dependency runs the other way.
- **100% test coverage** is the target, enforced as a CI gate, across 6 arches
  and 3 OSes.

## Where to go next

- [Why pure Go](why.md) — why this slice of Ruby is deterministic enough to live
  as a standalone, interpreter-independent Go library.
- [Usage & API](api.md) — the public surface and worked examples.
- [Roadmap](roadmap.md) — what is done and what is downstream by design.

Source lives at [github.com/go-ruby-shellwords/shellwords](https://github.com/go-ruby-shellwords/shellwords).
