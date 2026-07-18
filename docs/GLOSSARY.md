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
- **Start Line** — the `running: <command>, swallowing output` line printed
  in agent mode once the command has started; never echoes arguments.
- **Summary Line** — the `done: exit code <n>, …` line ending every agent
  mode run: on success it carries the read hint, on failure it precedes the
  replay.
- **Read Hint** — the runnable `` `swallow --read <log file name>` `` snippet
  in the success summary line; resolves via bare-name resolution.
- **Replay** — streaming the log back after a failed run in agent mode,
  restoring stdout lines to stdout and stderr lines to stderr.
- **Read** — printing a stored log verbatim via `swallow --read`, permitted
  only for logs of the current origin (the same-origin gate).
- **Retention** — the two-hour lifetime of logs; older logs are pruned at the
  start of every run.
