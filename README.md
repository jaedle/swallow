# swallow

Wraps a command and swallows its output when called from an agentic code
agent — the agent sees one line instead of thousands, saving LLM tokens. The
full output always lands in a log file.

```
$ CLAUDECODE=1 swallow go test ./...
running: go, swallowing output
done: exit code 0, read logs: `swallow --read 2026-07-18T10-15-30-go-a1b2c3.log`
```

## Behavior

- **Agent mode** (`CLAUDECODE=1`, set by Claude Code): output is suppressed.
  On success swallow prints the two lines above — the `--read` hint is
  directly runnable. On failure it prints `done: exit code <n>, full output:`
  and replays the complete output — stdout to stdout, stderr to stderr — and
  exits with the command's exit code. Command arguments are never echoed
  (they may contain shell-expanded secrets).
- **Human mode** (otherwise): output passes through live, and is still
  captured in the log.
- Logs: one file per run under `~/.swallow/<origin>/` (override the root with
  `$SWALLOW_DIR`), named with timestamp, command and a unique suffix. Logs
  older than two hours are pruned on every run.
- stdin, `SIGINT`/`SIGTERM` and the exit code pass through to the wrapped
  command.
- **Reading logs**: `swallow --read <log-file>` prints a captured log
  verbatim; a bare file name resolves against the current origin's log
  directory. Only logs of the current origin can be read — anything else is
  refused.

Full specification: [docs/SPEC.md](docs/SPEC.md).

## Install

With [mise-en-place](https://mise.jdx.dev):

```sh
mise use -g ubi:jaedle/swallow
```

Or grab a static binary for linux/darwin (amd64/arm64) from the
[releases](https://github.com/jaedle/swallow/releases).

To teach your coding agent to use swallow, install the
[agent skill](skills/swallow/SKILL.md):

```sh
npx skills add jaedle/swallow
```

## Usage

```
swallow [--] <command> [args...]
swallow --read <log-file>
swallow --version
swallow --help
```

The command is executed directly, without a shell.

## Development

```sh
mise install
task ci
```

`task ci` = fmt + lint + test + build. Tests are a Ginkgo/Gomega suite
exercising the compiled binary. See [AGENTS.md](AGENTS.md) for conventions.
