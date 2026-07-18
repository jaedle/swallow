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

## Consequences

- Deterministic and trivially testable; no TTY heuristics that misfire in
  pipes or CI.
- Other agents are not detected today. Support is added by extending this
  spec first (new ADR), not by ad-hoc code changes.
- Humans can force agent behavior for a run with `CLAUDECODE=1 swallow …`.
