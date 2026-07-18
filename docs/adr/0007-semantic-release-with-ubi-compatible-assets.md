# 7. semantic-release with ubi-compatible GitHub release assets

Date: 2026-07-18

## Status

Accepted

## Context

swallow must be installable via mise-en-place on macOS and Linux without a
package registry. mise's `ubi` backend installs directly from GitHub release
assets, selecting by OS/arch substrings in the asset name and extracting an
executable named like the project. Versioning should be automated from
commits, like the other jaedle repos (semantic-release, emoji conventional
commits, pipeline-service running `task release` on main).

## Decision

- semantic-release derives the version from emoji conventional commits on
  `main`, maintains `CHANGELOG.md`, creates the GitHub release.
- The `@semantic-release/exec` plugin runs `task release:publish` in the
  **prepare** phase (not publish, unlike the docker-based reference repo), so
  the assets exist before `@semantic-release/github` uploads them.
- Assets: `swallow-v<version>-<os>-<arch>.tar.gz` for linux/darwin ×
  amd64/arm64, each containing the static (`CGO_ENABLED=0`) binary `swallow`
  at the archive root.
- Baseline tag `v0.0.0` on the initial commit makes the first release
  `v0.1.0`.
- Install: `mise use -g ubi:jaedle/swallow`.

## Consequences

- No goreleaser, no extra packaging config — four `go build` + `tar`
  invocations in the Taskfile.
- `linux`/`darwin` and `amd64`/`arm64` are exactly the substrings ubi
  matches; renaming assets breaks installation.
- The release job needs the `github-token` secret exposed as `GITHUB_TOKEN`
  with push rights to `main` (CHANGELOG commit, `[skip ci]`).
