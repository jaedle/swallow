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
| `skills/swallow/` | agent skill installed via `npx skills add` | — |
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

## Agent skill

The frontmatter `description` in `skills/swallow/SKILL.md` is the trigger
surface: it alone decides when a coding agent invokes the skill. Change it
only with behavioral evidence, not by rewording:

- A/B the old and new description against freshly spawned agents (several
  Claude models × reasoning efforts) doing small fixture tasks; each agent
  reports the skills it invoked and every command it ran.
- Short-output tasks (`git status`/`log`/`diff`, `ls`, file creation, version
  checks) must not invoke the skill; noisy tasks (test suites, dependency
  installs, builds) must keep wrapping.
- Name offending commands explicitly. Abstract rules measured worse: a
  "~30 lines of output" threshold made models predict output size and skip
  legitimate test/build wraps.
- `task verify-skill` only proves the frontmatter parses; it says nothing
  about trigger quality.

Reference run: [PR #6](https://github.com/jaedle/swallow/pull/6) — false
positives on short commands 4/18 → 0/18, noisy-command coverage 16/18.

Findings from the 2026-07 invocation research
([PR #7](https://github.com/jaedle/swallow/pull/7) — 152 headless one-shot
trials, haiku/sonnet/opus × reasoning efforts, three description variants):

- Fresh agents did not over-invoke at all: 0 false positives across 82
  short-output / output-needed tasks in every variant. The measured failure
  mode is the opposite — noisy commands ran unwrapped in ~60% of trials,
  worst inside multi-step tasks. Perceived over-invocation is more likely
  session dynamics (skill body loaded once, repo instructions priming
  swallow) than the description.
- An "output longer than N lines" threshold again measured no better:
  false positives unchanged at 0, coverage slightly worse. Do not encode
  output-size predictions.
- Claude Code truncates skill listings for small models: haiku sessions saw
  the bare name `- swallow` with no description text, so description tuning
  cannot reach haiku-driven agents at all.

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
