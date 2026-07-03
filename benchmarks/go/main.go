// SPDX-License-Identifier: BSD-3-Clause
//
// Library-level workload for go-ruby-shellwords/shellwords, mirrored
// byte-for-byte by ruby/shellwords.rb. Each sub-benchmark exercises one
// representative shellwords operation; check() emits the canonical output digest
// so run.sh can prove the pure-Go result equals MRI's before the numbers are
// trusted.
package main

import (
	"strings"

	"github.com/go-ruby-shellwords/shellwords"
)

// Inputs — byte-identical to the Ruby side (see ruby/shellwords.rb for the
// annotated originals).

// splitLine: a representative command line with bare words, --flag=value, a
// double-quoted argument with a space, a single-quoted argument, a
// backslash-escaped space, an absolute path, and quote-concatenation.
const splitLine = `prog --name="John Doe" -x 'a b c' plain\ arg --path=/usr/local/bin file"1".txt`

// escapeStr: a string peppered with shell metacharacters that all require escaping.
const escapeStr = `foo bar & baz | qux; echo $HOME > /tmp/f 'x' "y" (z)`

// joinArgs: an argument vector for Join, including a space, an embedded single
// quote, a semicolon, a path with a space, a $-variable, and an EMPTY string
// (which escapes to '' — the empty-argument special case).
var joinArgs = []string{"cmd", "arg with space", "quote'd", "semi;colon", "/path/to file", "$VAR", ""}

// --- operations (each returns its canonical output for the equality check) ---

// opSplit tokenises splitLine; tokens are joined with "\n" for the digest (no
// token contains a newline, so the join is a faithful canonical form).
func opSplit() string {
	toks, err := shellwords.Split(splitLine)
	if err != nil {
		panic(err)
	}
	return strings.Join(toks, "\n")
}

func opEscape() string {
	return shellwords.Escape(escapeStr)
}

func opJoin() string {
	return shellwords.Join(joinArgs)
}

func main() {
	// Verify digests (outside timing): consumed by run.sh.
	check("split", opSplit())
	check("escape", opEscape())
	check("join", opJoin())

	// Timed passes.
	bench("split", 2000, func() { sink = opSplit() })
	bench("escape", 2000, func() { sink = opEscape() })
	bench("join", 2000, func() { sink = opJoin() })
}
