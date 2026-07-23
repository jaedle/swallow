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
  `swallow: running <command>, swallowing output`. `<command>` is the
  basename of argv[0] only — arguments are never echoed, because the shell
  substitutes variables before swallow sees them, so echoed arguments could
  leak secrets into the caller's context (see ADR 0009). A command that
  cannot be started produces no start line, so every start line is followed
  by a summary line.
- The wrapped command's output is fully suppressed while it runs; it is only
  written to the [log](logging.md). Nothing is printed while the command
  runs.
- Exit code `0`: swallow prints a summary line to stdout and exits `0`.
  The summary depends on the log's line count `<c>`:
  - `<c>` = 0: `swallow: done, exit code 0, no output`.
  - `<c>` at most 10: `swallow: done, exit code 0, output (<c> lines):`
    followed by the full replay — restoring `out|` lines to stdout and
    `err|` lines to stderr. This pass-through makes a needless wrap
    harmless: withholding a short output costs more than showing it (the
    read hint alone is about as long as five output lines, and a hidden
    short answer usually triggers a read round trip on top), so callers
    can wrap when in doubt.
  - otherwise: ``swallow: done, exit code 0, <c> log lines, read:
    `swallow --read <log file name>` `` — output was withheld, so the
    summary carries the read hint; the line count tells the caller whether
    reading is worthwhile.
- Exit code `!= 0`: swallow prints a summary line to stderr, replays the log
  — restoring `out|` lines to stdout and `err|` lines to stderr — prints an
  end marker to stderr, and exits with the command's exit code.
  - The replay is capped at the last 100 log lines: failures with huge
    output would otherwise flood the caller's context, defeating swallow's
    purpose; the last lines win because that is where the error usually is.
    The full output stays available via the read hint.
  - Within the cap the summary line is `swallow: done, exit code <n>, full
    output (<c> lines):` and the end marker is `swallow: end of output,
    exit code <n>`.
  - When truncated the summary line is `swallow: done, exit code <n>, last
    100 of <c> lines:` and the end marker carries the read hint —
    ``swallow: end of output, exit code <n>, read: `swallow --read <log
    file name>` `` — because only then was output withheld.
  - The end marker proves the replay is complete and repeats the verdict,
    so the exit code can be read from either end of the output.
- The read hint appears at most once per run, and only when output was
  actually withheld — every agent mode line costs the caller tokens, and
  the hint, dominated by the log file name, is the longest element.
  `<log file name>` is the log's bare file name, so the hinted command
  works verbatim from the same working directory (see [reading
  logs](read.md)).

## Human mode

- The wrapped command's output is teed live: stdout to stdout, stderr to
  stderr, and both into the log.
- No summary line and no replay; the exit code is propagated unchanged.
