# Modes

swallow runs in exactly one of two modes, decided per run.

## Detection

Agent mode is active if and only if the environment variable `CLAUDECODE`
equals `1` (the marker Claude Code sets for processes it spawns). Any other
value — including empty or `true` — means human mode.

## Agent mode

- Before the command runs, swallow prints a start line to stdout:
  `running: <command>, swallowing output`. `<command>` is the command name
  (argv[0]) only — arguments are never echoed, because the shell substitutes
  variables before swallow sees them, so echoed arguments could leak secrets
  into the caller's context (see ADR 0009).
- The wrapped command's output is fully suppressed while it runs; it is only
  written to the [log](logging.md). Nothing is printed while the command
  runs.
- Exit code `0`: swallow prints a single summary line to stdout —
  ``done: exit code 0, read logs: `swallow --read <log file name>` `` —
  where `<log file name>` is the log's bare file name, so the hinted command
  works verbatim from the same working directory (see [reading
  logs](read.md)). Exits `0`.
- Exit code `!= 0`: swallow prints `done: exit code <n>, full output:` to
  stderr, then replays the complete log, restoring `out|` lines to stdout
  and `err|` lines to stderr, and exits with the command's exit code. No
  `--read` hint — the output was just replayed.

## Human mode

- The wrapped command's output is teed live: stdout to stdout, stderr to
  stderr, and both into the log.
- No summary line and no replay; the exit code is propagated unchanged.
