# 6. Binary end-to-end tests with Ginkgo and Gomega

Date: 2026-07-18

## Status

Accepted

## Context

The behavior of swallow is only real at the process boundary: environment
detection, stream wiring, signals, exit codes. Unit tests of internals would
assert the implementation, not the spec.

## Decision

- One Ginkgo v2 + Gomega suite in `test/`, written as behavior specs.
- `gexec.Build` compiles the actual `cmd/swallow` binary once per suite run;
  `gexec.Start` runs it as a subprocess; assertions use `gbytes` and
  `gexec.Exit`.
- No unit tests unless a behavior cannot be reached through the binary.
- No bash scripts for testing or verification — everything runs inside the
  suite via `go test ./...`.
- The child environment is always built from scratch so the developer's own
  environment (notably `CLAUDECODE=1` when developing inside Claude Code)
  cannot leak into specs.

## Consequences

- Specs double as the executable form of `docs/spec/`.
- The suite needs a Go toolchain at test time (it compiles the binary) —
  given, since tests run via `go test`.
- Slightly slower than unit tests; acceptable at this suite size.
