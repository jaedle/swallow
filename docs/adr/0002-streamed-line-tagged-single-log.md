# 2. Streamed, line-tagged single log

Date: 2026-07-18

## Status

Accepted

## Context

The wrapped command's output must be captured without holding it in memory
(the binary shall stay small in memory even for huge outputs). On failure in
agent mode the output must be replayed with stdout and stderr restored to
their original streams, so the stream identity must survive capture. A single
log per run is wanted so "have a look at the log" is one file.

## Decision

- One log file per run.
- Both streams are read line-wise through pipes and appended to the log as
  they arrive, each line prefixed with a stream tag (`out|` / `err|`).
- Every log write is one complete tagged line; the file is opened with
  `O_APPEND` so concurrent appends from the two stream readers keep framing
  intact.
- Replay parses the tags and routes each line back to its stream.

## Consequences

- Memory is bounded by the longest single line, never by the output size.
- Interleaving is line-granular in arrival order — equivalent to what a
  terminal shows, not guaranteed byte-exact.
- The log is not byte-identical to the raw output (tags, enforced trailing
  newline); readers and replay must strip tags.
