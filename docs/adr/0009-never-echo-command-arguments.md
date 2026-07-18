# 9. Never echo command arguments

Date: 2026-07-18

## Status

Accepted

## Context

Agent mode prints a start line announcing the wrapped command. Echoing the
full argv seems helpful, but the shell expands variables before swallow sees
them: `swallow curl -H "Authorization: Bearer $TOKEN" …` reaches swallow
with the secret already substituted, and swallow cannot tell an expanded
variable from a literal. Anything echoed lands in the calling agent's
context — the very surface swallow exists to keep clean. Truncation does not
help; secrets routinely sit at the front of a command line.

The echo also carries no information: the caller composed the command and
already has it.

## Decision

swallow never prints command arguments. The start line contains argv[0]
only; log file names likewise derive from argv[0] (ADR 0003).

## Consequences

- Secrets passed as arguments cannot leak through swallow's own output.
  Secrets the command itself prints still end up in the log and the failure
  replay — identical to running the command without swallow.
- Two runs of the same program with different arguments produce
  indistinguishable start lines; the caller disambiguates by context.
- Future output changes must preserve this rule rather than "improve" the
  start line with full argv.
