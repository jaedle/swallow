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
- **Replay** — streaming the log back after a failed run in agent mode,
  restoring stdout lines to stdout and stderr lines to stderr.
- **Retention** — the two-hour lifetime of logs; older logs are pruned at the
  start of every run.
