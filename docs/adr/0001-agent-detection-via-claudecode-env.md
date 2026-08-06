# 1. Agent detection via CLAUDECODE environment variable

Date: 2026-07-18

## Status

Accepted

## Context

swallow must decide automatically whether it is called from an agentic code
agent (suppress output) or by a human (tee output). Alternatives considered:
always suppress, TTY detection, a best-effort list of agent environment
variables, or a swallow-specific variable.

## Decision

Agent mode is active if and only if `CLAUDECODE=1` — the environment marker
Claude Code sets for every process it spawns. Exact match; any other value
means human mode.

Extended 2026-08-06: agent mode is also active when `OPENCODE=1` — the marker
OpenCode sets for every process it spawns. The strict `==1` rule stays
symmetric for both markers; the title and file name stay.

## Consequences

- Deterministic and trivially testable; no TTY heuristics that misfire in
  pipes or CI.
- Other agents are not detected today. Support is added by extending this
  spec first, not by ad-hoc code changes.
- Humans can force agent behavior for a run with `CLAUDECODE=1 swallow …`.
