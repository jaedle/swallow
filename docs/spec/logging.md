# Logging

Every run captures the wrapped command's stdout and stderr into exactly one
log file.

## Location

- Root: the swallow dir — `$SWALLOW_DIR` if set and non-empty, otherwise
  `~/.swallow`.
- Per origin: logs live in a subdirectory named by the slug of the working
  directory swallow was invoked from.
- Slug: every run of non-alphanumeric characters becomes a single `-`,
  leading/trailing dashes are trimmed. `/` slugs to `root`.
  Example: `/home/jaedle/code/github.com/jaedle/swallow` →
  `home-jaedle-code-github-com-jaedle-swallow`.

## Naming

`<timestamp>-<command>-<suffix>.log`

- Timestamp: `YYYY-MM-DDTHH-MM-SS`, local time — lexically sortable so old
  logs can be recognized for [retention](retention.md).
- Command: slug of the basename of the wrapped command's argv[0] — slugging
  keeps the log name shell-safe so the read hint printed after a run is
  runnable verbatim.
- Suffix: 6 hex characters from a cryptographic random source — reruns and
  concurrent runs always produce distinct logs.

Example: `~/.swallow/home-jaedle-code-github-com-jaedle-swallow/2026-07-18T10-15-30-go-a1b2c3.log`

## Format

One file per run, both streams interleaved line-wise in arrival order. Every
line carries a stream tag:

```
out|a line written to stdout
err|a line written to stderr
```

A final line without a trailing newline is terminated with one in the log.

## Streaming guarantee

Output streams to disk as it is produced. swallow never holds the accumulated
output in memory; memory usage is bounded by the longest single line. Each log
write is one complete tagged line, so concurrent appends from both streams
never corrupt the line framing.
