# Glossary

- **Run** — one invocation of `swallow <command>`: execute the command, stream
  its output to exactly one log, exit with the command's exit code.
- **Agent Mode** — the mode active when the environment variable `CLAUDECODE`
  equals `1`: output is suppressed and only summarized/replayed.
- **Human Mode** — the mode active otherwise: output is teed live to the
  terminal in addition to the log.
- **Origin** — the working directory swallow was invoked from; determines the
  log subdirectory via its slug.
- **Swallow Dir** — the root directory for logs: `$SWALLOW_DIR` if set,
  otherwise `~/.swallow`.
- **Log** — the single file per run capturing stdout and stderr of the wrapped
  command, line-tagged by stream.
- **Stream Tag** — the `out|` / `err|` prefix on each log line identifying the
  originating stream.
- **Start Line** — the `swallow: running <command>, swallowing output` line
  printed in agent mode once the command has started; never echoes arguments.
- **Summary Line** — the `swallow: done, exit code <n>, …` line of every
  agent mode run: on success it reports no output, a pass-through, or the
  log line count with the read hint; on failure it precedes the replay.
- **Pass-Through** — printing a successful run's output directly after the
  summary line when it is at most 10 log lines: withholding a short output
  costs more than showing it, and it makes a needless wrap harmless.
- **End Marker** — the `swallow: end of output, exit code <n>, …` line
  closing a failure replay: proves the replay is complete and repeats the
  verdict; carries the read hint when the replay was truncated.
- **Read Hint** — the runnable `` `swallow --read <log file name>` `` snippet
  in the long-success summary line and the truncated-failure end marker;
  resolves via bare-name resolution.
- **Replay** — streaming the log back after a failed run in agent mode,
  restoring stdout lines to stdout and stderr lines to stderr.
- **Read** — printing a stored log verbatim via `swallow --read`, permitted
  only for logs of the current origin (the same-origin gate).
- **Retention** — the two-hour lifetime of logs; older logs are pruned at the
  start of every run.
