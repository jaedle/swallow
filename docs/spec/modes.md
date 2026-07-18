# Modes

swallow runs in exactly one of two modes, decided per run.

## Detection

Agent mode is active if and only if the environment variable `CLAUDECODE`
equals `1` (the marker Claude Code sets for processes it spawns). Any other
value — including empty or `true` — means human mode.

## Agent mode

- The wrapped command's output is fully suppressed while it runs; it is only
  written to the [log](logging.md).
- Exit code `0`: swallow prints a single line to stdout —
  `everything went fine (log: <path>)` — and exits `0`.
- Exit code `!= 0`: swallow replays the complete log, restoring `out|` lines
  to stdout and `err|` lines to stderr, prints
  `swallow: command failed with exit code <n> (log: <path>)` to stderr, and
  exits with the command's exit code.

## Human mode

- The wrapped command's output is teed live: stdout to stdout, stderr to
  stderr, and both into the log.
- No summary line and no replay; the exit code is propagated unchanged.
