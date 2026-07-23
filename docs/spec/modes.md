# Modes

swallow runs in exactly one of two modes, decided per run.

## Detection

Agent mode is active if and only if the environment variable `CLAUDECODE`
equals `1` (the marker Claude Code sets for processes it spawns). Any other
value — including empty or `true` — means human mode.

## Agent mode

Every line swallow prints itself carries the `swallow: ` prefix, so wrapped
command output replayed in between can never be mistaken for swallow's own
lines.

- Once the command has started, swallow prints a start line to stdout:
  ``swallow: running <command>, log: `swallow --read <log file name>` ``.
  `<command>` is the basename of argv[0] only — arguments are never echoed,
  because the shell substitutes variables before swallow sees them, so
  echoed arguments could leak secrets into the caller's context (see ADR
  0009). The read hint is part of the start line and appears exactly once
  per run: every agent mode line costs the caller tokens, and the hint —
  dominated by the log file name — is the longest element, so it is never
  repeated. `<log file name>` is the log's bare file name, so the hinted
  command works verbatim from the same working directory (see [reading
  logs](read.md)). A command that cannot be started produces no start line,
  so every start line is followed by a summary line.
- The wrapped command's output is fully suppressed while it runs; it is only
  written to the [log](logging.md). Nothing is printed while the command
  runs.
- Exit code `0`: swallow prints a single summary line to stdout —
  `swallow: done, exit code 0, <c> log lines` — where `<c>` is the log's
  line count, so the caller can judge whether reading the log is
  worthwhile. Exits `0`.
- Exit code `!= 0`: swallow prints a summary line to stderr, replays the log
  — restoring `out|` lines to stdout and `err|` lines to stderr — prints an
  end marker to stderr, and exits with the command's exit code.
  - The replay is capped at the last 100 log lines: failures with huge
    output would otherwise flood the caller's context, defeating swallow's
    purpose; the last lines win because that is where the error usually is.
    The full output is always available via the read hint.
  - Within the cap the summary line is `swallow: done, exit code <n>, full
    output (<c> lines):`; when truncated it is `swallow: done, exit code
    <n>, last 100 of <c> lines:`.
  - The end marker — `swallow: end of output, exit code <n>` — proves the
    replay is complete and repeats the verdict, so the exit code can be
    read from either end of the output. The read hint is not repeated; it
    is on the start line.

## Human mode

- The wrapped command's output is teed live: stdout to stdout, stderr to
  stderr, and both into the log.
- No summary line and no replay; the exit code is propagated unchanged.
