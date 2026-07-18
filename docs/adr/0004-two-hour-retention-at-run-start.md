# 4. Two-hour retention, pruned at run start

Date: 2026-07-18

## Status

Accepted

## Context

Logs exist so an agent (or human) can inspect the last runs; they have no
long-term value and would otherwise accumulate unboundedly. Timestamped names
were chosen so old logs are recognizable (ADR 0003).

## Decision

Every run starts by pruning log files older than two hours (by modification
time) across the whole swallow dir, removing emptied origin directories.
Best-effort, no daemon, no explicit clean command.

## Consequences

- Disk usage is self-limiting for anyone who keeps using swallow; a user who
  stops using it keeps at most two hours of logs.
- Modification time (not the name's timestamp) is the criterion — cheap, and
  it protects logs of still-running commands, which are written continuously.
- A run idle for more than two hours can lose its log to a concurrent run's
  pruning — accepted trade-off.
