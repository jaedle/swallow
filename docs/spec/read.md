# Reading logs

```
swallow --read <log-file>
```

Prints a previously captured log to stdout, verbatim — including the stream
tags, exactly as stored on disk. There is no replay: both streams stay on
stdout, in the tagged form. Reading behaves identically in agent mode and
human mode.

## Same-origin gate

Only logs of the current origin may be read: the requested path must point
directly into the origin directory of the working directory swallow is
invoked from — `<swallow dir>/<slug(cwd)>/`. Logs of other origins, or any
path outside the swallow dir, are refused.

- A relative path is resolved against the working directory.
- The path is lexically cleaned before the check, so `..` segments cannot
  escape the gate.
- On refusal swallow prints `swallow: refusing to read <path>: not a log of
  the current origin` to stderr and exits `1`. The file's existence is not
  disclosed.

## Exit codes

| Code | Meaning |
| --- | --- |
| `0` | log printed |
| `1` | refused by the same-origin gate, or the log cannot be read |
| `2` | usage error (missing operand) |

## Retention

Reading is not a run: it never prunes logs. Note that logs older than the
[retention](retention.md) lifetime are pruned by the next run, so only
recent logs remain readable.
