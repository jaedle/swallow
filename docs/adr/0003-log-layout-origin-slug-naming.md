# 3. Log layout: origin slug directories, timestamped unique names

Date: 2026-07-18

## Status

Accepted

## Context

Logs from different projects should be distinguishable, rerunning a command
must never overwrite an earlier log, and retention needs to recognize old
logs cheaply.

## Decision

`<swallow-dir>/<origin-slug>/<timestamp>-<command>-<random>.log`

- The full originating working directory is slugged (non-alphanumeric runs →
  `-`) and used as a subdirectory, grouping logs per project.
- The filename carries a sortable local timestamp, the wrapped command's
  basename, and 6 random hex characters for uniqueness.

## Consequences

- Two different directories slugging to the same value share a subdirectory —
  harmless, logs remain unique by name.
- Retention can prune by file modification time and remove empty origin
  directories.
- Slugs of deep paths are long but stay well below filesystem name limits.
