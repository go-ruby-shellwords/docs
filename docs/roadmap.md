# Roadmap

`go-ruby-shellwords/shellwords` is grown **test-first**, each capability differential-tested against MRI
rather than built in isolation. Ruby's Shellwords — the
deterministic, interpreter-independent slice extracted from rbgo's internals — is
**complete**.

| Stage | What | Status |
| --- | --- | --- |
| shellsplit | POSIX `sh` word splitting honouring single quotes, double quotes and backslash escapes, exactly as Ruby's `Shellwords.shellsplit` does. | **Done** |
| shellescape | Escaping a single string so it survives one round of shell parsing, byte-for-byte identical to MRI (including the empty-string `''` case). | **Done** |
| shelljoin | Joining an array of words into a single shell-safe command line, the inverse of `shellsplit`, matching reference Ruby. | **Done** |
| Error on malformed input | An unmatched quote raises the same `ArgumentError` Ruby raises, at the same point in the input. | **Done** |
| Round-trip fidelity | `shellsplit(shelljoin(words)) == words` over arbitrary words, the property reference Ruby guarantees. | **Done** |
| Differential oracle & coverage | A wide word corpus split, escaped and joined both here and by the system `ruby`/`shellwords`, compared byte-for-byte; 100% coverage, gofmt + go vet clean, green across all six 64-bit Go arches and three OSes. | **Done** |

## Documented out-of-scope boundaries

These are **deliberate**, recorded so the module's surface is unambiguous:

- **No interpreter.** The library implements the deterministic algorithm; it
  never runs arbitrary Ruby. Anything that needs a live binding or evaluation is
  the consumer's job — that is why `rbgo` binds this module rather than the
  reverse.
- **Reference is reference Ruby (MRI).** Byte-for-byte conformance targets MRI's
  behaviour; differences across MRI releases are matched to the reference used by
  the differential oracle.
- **Standalone & reusable.** The module has no dependency on the Ruby runtime;
  the dependency runs the other way.

See [Usage & API](api.md) for the surface and [Why pure Go](why.md) for the
deterministic/interpreter split.
