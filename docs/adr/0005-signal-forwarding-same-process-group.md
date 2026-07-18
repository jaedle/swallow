# 5. Signal forwarding, child in the same process group

Date: 2026-07-18

## Status

Accepted

## Context

Agents and process managers signal the swallow process directly; interactive
terminals send SIGINT to the whole foreground process group. The wrapped
command must terminate in both cases, and a TTY-reading child must keep
working.

## Decision

- The child stays in swallow's process group (no `Setpgid`).
- swallow forwards received `SIGINT`/`SIGTERM` to the child unconditionally
  and never exits on a signal itself; it waits and propagates the child's
  exit code, `128 + n` if the child died by signal `n`.

## Consequences

- Interactive Ctrl-C reaches the child twice (once via the terminal's group
  delivery, once forwarded) — harmless for virtually all programs and the
  standard wrapper behavior. Do not "fix" this with `Setpgid`: that would
  detach the child from the foreground group, breaking direct Ctrl-C delivery
  and TTY stdin (`SIGTTIN`).
- A directly-targeted `kill -INT <swallow-pid>` reliably reaches the child.
