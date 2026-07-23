---
name: swallow
description: >-
  Wrap shell commands with `swallow` to keep long output out of context:
  tests, builds, linters, installs — any command where you don't need the
  full output of a success. When in doubt, wrap: output of 10 lines or
  fewer passes straight through, and a failure still replays its last 100
  lines, so wrapping never hides what you need. Skip it only when the
  output itself is the answer — git diff, grep, cat.
---

# swallow

`swallow` runs a command and, when its output is long, streams it to a log
file instead of your context, saving tokens. Short output is shown directly.

## Usage

Prefix the command:

```sh
swallow go test ./...
swallow npm install
```

Output while the command runs is a single start line —
`swallow: running <command>, swallowing output` — nothing else until it
finishes. Lines swallow prints itself always start with `swallow: `;
anything else is output of the wrapped command.

- Success with no output: `swallow: done, exit code 0, no output`.
- Success with 10 lines or fewer: `swallow: done, exit code 0, output
  (<n> lines):` followed by the output itself — nothing was withheld.
- Success with longer output: `` swallow: done, exit code 0, <n> log
  lines, read: `swallow --read <log-file>` `` — do not read the log unless
  you genuinely need its content; the exit code is the answer.
- Failure: `swallow: done, exit code <n>, …` followed by the replayed
  output (capped at the last 100 lines) and a closing `swallow: end of
  output, …` marker; when the replay was truncated the marker carries a
  read hint to the full log. swallow exits with the command's exit code.

## Notes

- Need output from a successful run? Run the hinted
  `swallow --read <log-file>` command verbatim (works only from the same
  working directory).
- Skip swallow when the output is the answer (`git diff`, `cat`) or for
  interactive commands.
- If `swallow` is not on PATH, run the command directly.
