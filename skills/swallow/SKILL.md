---
name: swallow
description: >-
  Wrap noisy shell commands (tests, builds, linters, installs) with
  `swallow` to keep their output out of context. Use when only the outcome
  matters — a failure still replays its last 100 lines, so nothing is lost
  by wrapping. Never wrap a command whose output you need to read — git
  status/log/diff, grep and other searches, ls, cat, file creation, version
  checks — run those directly.
---

# swallow

`swallow` runs a command and streams its full output to a log file instead of
your context, saving tokens.

## Usage

Prefix the command:

```sh
swallow go test ./...
swallow npm install
```

Output while the command runs is a single start line —
`` swallow: running <command>, log: `swallow --read <log-file>` `` —
nothing else until it finishes. Lines swallow prints itself always start
with `swallow: `; anything else is output of the wrapped command.

- Success: `swallow: done, exit code 0, <n> log lines` — done, move on. Do
  not read the log after a success unless you genuinely need its content;
  the exit code and line count are the answer.
- Failure: `swallow: done, exit code <n>, …` followed by the replayed
  output (capped at the last 100 lines) and a closing `swallow: end of
  output, exit code <n>` marker; swallow exits with the command's exit
  code. Use the start line's read hint when the replay was truncated.

## Notes

- Need output from a successful run? Run the start line's
  `swallow --read <log-file>` command verbatim (works only from the same
  working directory).
- Skip swallow when the output is the answer (`git diff`, `cat`) or for
  interactive commands.
- If `swallow` is not on PATH, run the command directly.
