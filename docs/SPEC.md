# swallow — Specification

`swallow` wraps a command, streams its output to a log file, and suppresses it
for agentic callers to reduce LLM token usage.

This document is the index. The behavior is specified per aspect:

- [CLI](spec/cli.md) — invocation, arguments, exit codes, stdin, signals
- [Modes](spec/modes.md) — agent mode vs. human mode, suppression, replay
- [Logging](spec/logging.md) — log location, naming, format, streaming
- [Retention](spec/retention.md) — pruning of old logs

Terms are defined in the [glossary](GLOSSARY.md). Decisions are recorded as
[ADRs](adr/).

The Ginkgo suite in `test/` is the executable form of this specification.
