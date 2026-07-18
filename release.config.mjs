// Commits follow gitmoji + conventional commits (e.g. "✨ feat: ..."), so the
// header patterns allow an optional leading emoji token before the type.
const headerPattern = /^(?:[^\w\s]+\s)?(\w+)(?:\(([^)]*)\))?!?: (.+)$/;
const breakingHeaderPattern = /^(?:[^\w\s]+\s)?(\w+)(?:\(([^)]*)\))?!: (.+)$/;
const parserOpts = { headerPattern, breakingHeaderPattern };

export default {
  branches: ["main"],
  plugins: [
    ["@semantic-release/commit-analyzer", { preset: "conventionalcommits", parserOpts }],
    ["@semantic-release/release-notes-generator", { preset: "conventionalcommits", parserOpts }],
    ["@semantic-release/changelog", { changelogFile: "CHANGELOG.md" }],
    ["@semantic-release/exec", { prepareCmd: "task release:publish VERSION=${nextRelease.version}" }],
    ["@semantic-release/github", { assets: [{ path: "out/swallow-v*.tar.gz" }] }],
    [
      "@semantic-release/git",
      {
        assets: ["CHANGELOG.md"],
        message: "🔖 chore(release): ${nextRelease.version} [skip ci]\n\n${nextRelease.notes}",
      },
    ],
  ],
};
