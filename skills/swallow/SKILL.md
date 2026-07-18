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

Output while the command runs is a single start line —
`running: <command>, swallowing output` — nothing else until it finishes.

- Success: `` done: exit code 0, read logs: `swallow --read <log-file>` `` —
  the hinted command is directly runnable.
- Failure: `done: exit code <n>, full output:` followed by the complete
  replayed output; swallow exits with the command's exit code.

## Notes

- Need output from a successful run? Run the hinted
  `swallow --read <log-file>` command verbatim (works only from the same
  working directory).
- Skip swallow when the output is the answer (`git diff`, `cat`) or for
  interactive commands.
- If `swallow` is not on PATH, run the command directly.
