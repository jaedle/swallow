# AGENTS.md

Guidance for coding agents working on this repository.

## What this is

`swallow` is a CLI wrapper that suppresses the output of a wrapped command when
called from an agentic code agent, to reduce LLM token usage. Output is streamed
to a log file on disk instead.

## Source of truth

`docs/SPEC.md` is the index to the behavior specification in `docs/spec/`.
Change the spec first, then make the code follow. The test suite asserts the
spec, not the implementation. Decisions are recorded in `docs/adr/`, domain
terms in `docs/GLOSSARY.md` — use the glossary terms consistently.

## Layout

| Path | Responsibility | Spec |
| --- | --- | --- |
| `cmd/swallow/` | entrypoint, argument handling, version | `docs/spec/cli.md` |
| `internal/swallow/run.go` | command execution, modes, signals, exit codes | `docs/spec/cli.md`, `docs/spec/modes.md` |
| `internal/swallow/logfile.go` | log directory resolution, naming | `docs/spec/logging.md` |
| `internal/swallow/retention.go` | log pruning | `docs/spec/retention.md` |
| `test/` | behavior specs against the compiled binary | all |
| `ci/config.yaml` | CI configuration for jaedle/pipeline-service | — |

## Verification

`task ci` = fmt + lint + test + build. Run it before considering any work done.

Tests are a Ginkgo v2 + Gomega suite in `test/` that builds the real binary
(`gexec.Build`) and exercises it as a subprocess (`gexec.Start`, `gbytes`).
Do not write unit tests unless a behavior cannot be reached through the binary.
Never write bash scripts for testing or verification.

Test conventions:
- Specs read as behavior: `Describe(<area>)` / `It(<behavior>)`.
- Build the child environment from scratch — the developer's own `CLAUDECODE=1`
  must never leak into a spec.
- Arrange/act/assert as blank-line-separated blocks, bodies short, helpers named.
- Every timeout is a named constant with a rationale.

## Git

- Work directly on `main` (owner's instruction for this repo).
- Commit per vertical slice; never commit when `task ci` fails.
- Commits via the `/commit` skill: emoji + conventional commit format.
- Releases are cut by semantic-release from commits on `main`; the release
  commit (`🔖 chore(release): …`) is created by the pipeline, never by hand.

## Dependencies

`task manage-dependencies` = verify + tidy + vendor. Dependencies are vendored.
When adding a dependency, use the latest version.

## CI / Release

CI is jaedle/pipeline-service reading `ci/config.yaml` (must stay at `ci/`):
verify runs `task ci`, release runs `task release`. Release assets are four
static tarballs (linux/darwin × amd64/arm64) uploaded to the GitHub release,
named for `mise`/`ubi` compatibility (`swallow-v<version>-<os>-<arch>.tar.gz`).
