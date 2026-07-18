# 8. Log reading gated to the current origin

Date: 2026-07-18

## Status

Accepted

## Context

After a successful run in agent mode the agent only sees the log path. To
inspect the output it needs a way to read the log — but an agent working in
one project should not browse logs of other projects, which may contain
unrelated or sensitive output.

## Decision

`swallow --read <log-file>` prints a log verbatim, gated to the current
origin: the lexically cleaned path (relative paths resolved against the
working directory) must lie directly in `<swallow-dir>/<slug(cwd)>/`. Every
other path is refused with exit `1`, without disclosing whether it exists.

The gate compares paths lexically. It scopes an agent to its own project's
logs; it is not a security boundary against a hostile local user, who could
read the files directly anyway.

## Consequences

- An agent can follow up `everything went fine (log: <path>)` with
  `swallow --read <path>` from the same working directory.
- Reading a log requires being in the originating directory (or one with the
  same slug — the known slug-collision trade-off of ADR 0003).
- Symlink tricks are not defended against, consistent with the lexical
  comparison and the non-goal of defending against local users.
