# CLI

## Invocation

```
swallow [--] <command> [args...]
swallow --read <log-file>
swallow --version
```

- Everything after `swallow` (and an optional `--` separator) is the command,
  executed directly — no shell interpretation.
- Without a command: usage on stderr, exit `2`.
- `--read`: print a captured log, see [reading logs](read.md). A command named
  `--read` can be wrapped via the `--` separator.
- `--version`: print the version, exit `0`.

## stdin

The wrapped command inherits swallow's stdin unchanged.

## Signals

`SIGINT` and `SIGTERM` received by swallow are forwarded to the wrapped
command. swallow never exits on a signal itself; it waits for the command and
propagates its result. The command stays in swallow's process group (see ADR
0005).

swallow finishes once the command has exited *and* its output streams are
closed. Background children of the command that keep stdout/stderr open keep
the run (and its log) alive until they release the streams — this guarantees
the log is complete.

## Exit codes

| Code | Meaning |
| --- | --- |
| exit code of the command | command ran (including `0`) |
| `128 + n` | command was killed by signal `n` |
| `126` | command found but could not be started |
| `127` | command not found (no log is created) |
| `2` | usage error |
| `1` | swallow-internal failure (e.g. log not writable) |
