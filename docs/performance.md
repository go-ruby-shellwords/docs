# Performance

`go-ruby-shellwords/shellwords` is the pure-Go library that
[`rbgo`](https://github.com/go-embedded-ruby/ruby) binds for Ruby's `Shellwords`.
This page records a **comparative benchmark** of that module against the reference
Ruby runtimes, part of the ecosystem-wide per-module parity suite.

## What is measured

The **same** workload — the three representative `Shellwords` operations —
is run through the pure-Go library's Go API and through each reference runtime's
own `shellwords`:

- **split** — tokenise a representative command line
  (`prog --name="John Doe" -x 'a b c' plain\ arg --path=/usr/local/bin file"1".txt`)
  the way the Bourne shell does, honouring `'single'` / `"double"` quotes,
  backslash escapes and quote-concatenation
  (`Shellwords.split` / `String#shellsplit`).
- **escape** — escape a string full of shell metacharacters
  (`foo bar & baz | qux; echo $HOME > /tmp/f 'x' "y" (z)`) so the shell parses it
  back verbatim (`Shellwords.escape` / `String#shellescape`).
- **join** — escape-and-join an argument vector, including a space, an embedded
  quote, a semicolon, a `$`-variable and an **empty** element (which becomes `''`)
  (`Shellwords.join` / `Array#shelljoin`).

Unlike some ecosystem modules, Ruby's `shellwords` is a **pure-Ruby stdlib module**
(`lib/shellwords.rb`) — there is no C extension. So the MRI / YJIT / JRuby /
TruffleRuby columns measure each interpreter executing the same Ruby source, and
the go column measures **this pure-Go library** doing the equivalent work. Parity
is reachable head-on: same algorithm, different runtime. Every operation's output
is checked **byte-identical to MRI** (SHA-256 of the tokens / the escaped string)
before any timing is trusted — the harness aborts on a mismatch.

- **Method:** each process runs 3 untimed warm-up passes, then 25 timed passes of
  a fixed inner loop, timed with a monotonic clock; the **best** pass is reported
  as **ns/op** (lower is better). `vs MRI` < 1.00× means *faster than MRI*.
  Interpreter start-up is outside the timed region, so these are operation costs,
  not `ruby file.rb` process costs.
- **Runtimes:** `ruby` (MRI, the oracle) and `ruby --yjit`; `jruby` (on the JVM);
  `truffleruby` (GraalVM CE Native).

!!! note "rbgo end-to-end row"
    The whole-interpreter `rbgo`-vs-MRI row (single-shot `ruby file.rb` wall time,
    the format used by some other module pages) has not yet been captured for
    `Shellwords` on a controlled host, so it is not shown here rather than printed
    as a fabricated figure. The library-level section below is the real, measured
    parity result for this module.

## Library-level benchmark (Go API vs runtimes) — 2026-07-03

This section measures the **pure-Go library directly, through its Go API** — not
the `rbgo` interpreter path. It isolates the library primitive from
Ruby-interpreter dispatch, answering the parity question head-on: *is the pure-Go
implementation as fast as the reference runtime's own `Shellwords`?* The **same
workload, same inputs, same iteration counts** run through the Go library and
through each reference runtime's stdlib; outputs were checked identical to MRI
before any timing.

- **Host:** Apple M4 Max (`Mac16,5`, arm64), macOS 26.5.1 — **date 2026-07-03**.
- **Runtimes:** Go 1.26.4 · MRI `ruby 4.0.5 +PRISM` · MRI + YJIT · JRuby 10.1.0.0
  (OpenJDK 25) · TruffleRuby 34.0.1 (GraalVM CE Native).
- **Inputs:** the annotated command line, escape string and argument vector above,
  identical on both sides (see [`benchmarks/`](https://github.com/go-ruby-shellwords/docs/tree/main/benchmarks)).
- **Method:** 3 untimed warm-up passes, then 25 timed passes of a fixed 2000-op
  inner loop, monotonic clock, **best** pass reported as **ns/op**.

#### split

| Runtime | ns/op | vs MRI |
| --- | ---: | ---: |
| **go-ruby (pure Go)** | 333.4 | 0.06× |
| MRI | 5338.0 | 1.00× |
| MRI + YJIT | 4909.5 | 0.92× |
| JRuby | 4474.1 | 0.84× |
| TruffleRuby | 7845.5 | 1.47× |

#### escape

| Runtime | ns/op | vs MRI |
| --- | ---: | ---: |
| **go-ruby (pure Go)** | 156.9 | 0.03× |
| MRI | 5899.0 | 1.00× |
| MRI + YJIT | 5627.5 | 0.95× |
| JRuby | 2397.3 | 0.41× |
| TruffleRuby | 2135.7 | 0.36× |

#### join

| Runtime | ns/op | vs MRI |
| --- | ---: | ---: |
| **go-ruby (pure Go)** | 230.9 | 0.04× |
| MRI | 5417.5 | 1.00× |
| MRI + YJIT | 4835.5 | 0.89× |
| JRuby | 1624.6 | 0.30× |
| TruffleRuby | 2527.1 | 0.47× |

Across all three operations the pure-Go library is **dramatically faster than MRI**
— `split` **~16×** (0.06×), `escape` **~38×** (0.03×), `join` **~23×** (0.04×).
This is the expected shape when the reference is *interpreted pure Ruby*: MRI's
`shellsplit` drives a `String#scan` over a multi-alternative regex and `shellescape`
runs two `gsub` passes, where the Go implementation is a single tight byte-scanner
with no per-token regex machine or intermediate allocations.

**go vs YJIT — go wins every operation.** Even against MRI + YJIT the pure-Go
library is far ahead: `split` **~14.7×** faster (333.4 vs 4909.5 ns/op),
`escape` **~35.9×** faster (156.9 vs 5627.5), `join` **~20.9×** faster
(230.9 vs 4835.5). YJIT trims only a few percent off MRI here (0.89×–0.95×):
these string-scanning workloads are dominated by regex/`gsub` C runtime calls that
the JIT cannot inline away.

It also beats the JVM/GraalVM runtimes on every operation: JRuby (0.84× / 0.41× /
0.30× the *MRI* time, i.e. slower than the go column on all three) and TruffleRuby
(1.47× / 0.36× / 0.47×) — TruffleRuby is actually **slower than MRI on `split`** at
this warm-up budget. Output stays byte-identical to MRI throughout (the harness
verifies the SHA-256 of every op before timing).

!!! note "Reproduce"
    The harness is committed under
    [`benchmarks/`](https://github.com/go-ruby-shellwords/docs/tree/main/benchmarks):
    a self-contained Go driver (`go/`, pins the published library by
    pseudo-version in `go.mod`, no `replace`), the equivalent `ruby/shellwords.rb`
    workload, and `run.sh`. Run `bash benchmarks/run.sh`; it first verifies the Go
    output is byte-identical to MRI (SHA-256 per op) and aborts on mismatch, then
    times. Env `OUTER`/`WARM` tune the pass budget and `RUBY`/`JRUBY`/`TRUFFLERUBY`
    select the runtime binaries.

!!! warning "Warm-up budget & noise — framing"
    Numbers reflect a **fixed warm-process budget** (3 warm-up + 25 timed passes
    in one process). The JVM/GraalVM JITs (JRuby, TruffleRuby) may need a larger
    warm-up to reach steady state, so their columns can **understate** peak
    throughput. The go column's ops are sub-microsecond (156–333 ns/op), so treat
    per-row differences under ~10% as noise. Every number here is a **real measured
    value** from the dated run above — nothing is fabricated, estimated, or
    cherry-picked. The go-ruby column is the pure-Go library; every other column is
    that interpreter's own pure-Ruby `shellwords` doing the equivalent work.
