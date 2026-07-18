---
name: swallow
description: >-
  Wrap noisy shell commands (tests, builds, linters, installs) with `swallow`
  to keep their output out of context. Use whenever running a command expected
  to produce long output.
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

- Success: one line — `everything went fine (log: <path>)`.
- Failure: the complete output is replayed, and swallow exits with the
  command's exit code.

## Notes

- Need output from a successful run? Read or grep the printed log path.
- Skip swallow when the output is the answer (`git diff`, `cat`) or for
  interactive commands.
- If `swallow` is not on PATH, run the command directly.
