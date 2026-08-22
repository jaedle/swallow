# 10. Cursor agent detection via CURSOR_AGENT presence

Date: 2026-08-22

## Status

Accepted

## Context

Cursor's agent sets `CURSOR_AGENT` for every process it spawns. Cursor's own
documentation tests the marker by presence (`[[ -n "$CURSOR_AGENT" ]]`), and
the semantics of its value are undocumented — a strict `==1` rule like the
`CLAUDECODE`/`OPENCODE` markers (see ADR 0001) would risk silent misses
whenever Cursor sets any other value.

## Decision

Agent mode is also active when `CURSOR_AGENT` is set to any non-empty value.
Unlike the existing markers, presence is the test, not a specific value. The
strict `==1` semantics of `CLAUDECODE` and `OPENCODE` stay unchanged.

## Consequences

- Detection matches Cursor's own documented check, so any value Cursor
  chooses today or later is honored.
- Value semantics differ between markers: `CLAUDECODE`/`OPENCODE` are strict
  `==1`, `CURSOR_AGENT` is non-empty — an accepted asymmetry, recorded here.
- A human running `CURSOR_AGENT=x swallow …` forces agent mode, same as the
  existing markers allow.
