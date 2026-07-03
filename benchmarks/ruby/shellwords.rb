# frozen_string_literal: true
# SPDX-License-Identifier: BSD-3-Clause
#
# Reference Shellwords workload, mirrored byte-for-byte by ../go/main.go.
# Ruby's `shellwords` is a PURE-RUBY stdlib module (lib/shellwords.rb), so the
# MRI / YJIT / JRuby / TruffleRuby columns measure each interpreter running that
# same Ruby source; the go column is the pure-Go go-ruby-shellwords library.
# Parity is therefore reachable head-on: same algorithm, different runtime.
require "shellwords"
require_relative "_harness"

# A representative command line exercising bare words, --flag=value, a
# double-quoted argument with a space, a single-quoted argument, a
# backslash-escaped space, an absolute path, and quote-concatenation.
SPLIT_LINE = %q{prog --name="John Doe" -x 'a b c' plain\ arg --path=/usr/local/bin file"1".txt}

# A string peppered with shell metacharacters that all require escaping.
ESCAPE_STR = %q{foo bar & baz | qux; echo $HOME > /tmp/f 'x' "y" (z)}

# An argument vector for shelljoin, including a space, an embedded single quote,
# a semicolon, a path with a space, a $-variable, and an EMPTY string (which
# escapes to '' — the empty-argument special case).
JOIN_ARGS = ["cmd", "arg with space", "quote'd", "semi;colon", "/path/to file", "$VAR", ""].freeze

# --- operations (each returns its canonical output for the equality check) ---

# op_split tokenises SPLIT_LINE; tokens are joined with "\n" for the digest (no
# token contains a newline, so the join is a faithful canonical form).
def op_split
  Shellwords.split(SPLIT_LINE).join("\n")
end

def op_escape
  Shellwords.escape(ESCAPE_STR)
end

def op_join
  Shellwords.join(JOIN_ARGS)
end

# Verify digests (outside timing): consumed by run.sh.
check("split",  op_split)
check("escape", op_escape)
check("join",   op_join)

# Timed passes.
bench("split",  2000) { op_split }
bench("escape", 2000) { op_escape }
bench("join",   2000) { op_join }
