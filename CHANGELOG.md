# Changelog

All notable changes to CrashBench will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- **Automated CI/CD Quality Gate (`crashbench check`)**:
  - Threshold enforcement via `--fail-under <score>` flag (exits with code 1 if resilience drops below requirement).
  - SARIF 2.1.0 export (`--sarif`) with native GitHub Code Scanning alerts integration (`pkg/export/sarif.go`).
  - JUnit XML export (`--junit`) for CI/CD test visualizers across Jenkins, GitLab, CircleCI (`pkg/export/junit.go`).
  - GitHub Pull Request markdown comment table generator (`--summary-md`) (`pkg/export/markdown.go`).
  - One-command GitHub Actions workflow generator (`crashbench init-ci`).
- **Official Security & Vulnerability Taxonomy Mapping**:
  - All chaos scenarios mapped to CWE (Common Weakness Enumeration) and OWASP GenAI Top 10 (2026).
  - New chaos scenario `SCN-07-ENVPOISON`: Toxic PATH & Dependency Hijack (CWE-426 / OWASP LLM06).
  - New chaos scenario `SCN-08-NETJITTER`: Flaky Socket & Truncated Network Stream (CWE-754 / OWASP LLM04).
- **Context Token Bloat & Financial Waste Telemetry**:
  - Real-time measurement of input context tokens ingested by runaway compiler/terminal dumps ($\approx \text{Bytes}/4$).
  - Estimated dollar-cost financial loss calculated at $3.00 / 1M input tokens.
- **Documentation**:
  - Added `docs/ci-cd.md` and `docs/taxonomy.md`.

## [1.1.0] - 2026-09-13

### Added
- **Bohemian & Handwritten Minimalist Interactive CLI (`crashbench interactive` / `crashbench`)**:
  - 24-bit TrueColor ANSI styling with artisan linocut stamp banner displaying automatic version `v1.1.0`.
  - Automatic version resolution via `runtime/debug.ReadBuildInfo()` with fallback to release constant.
  - Interactive Bradley-Terry Elo leaderboard with rendered 95% Bayesian bootstrap CI bars (`[───●━━─]`).
  - Side-by-Side Chaos Battle Arena with split terminal transcript replays and live crowdsourced resilience voting.
  - Pairwise Head-to-Head Win-Rate frequency heatmap table.
  - Deep Agent Scorecard Inspector with multi-line terminal execution replays.
  - CS Formulations & Research Field Notes viewer (Bradley-Terry MLE, Aho-Corasick DFA, Shannon Entropy, Tarjan DAG).
- **New Chaos Scenario**:
  - `SCN-06-SYMLINK`: Filesystem Path Traversal & Symlink Escape (CWE-22 / CWE-59). Tests whether the agent prevents unauthorized directory traversal to host credentials.
- **msh-protocol Agent Workflow Rules**:
  - Integrated `msh-protocol` development guidelines into `.agents/AGENTS.md`.

### Changed
- Decoupled web server from CLI menu; dedicated server command `crashbench serve --port 4040` serves the web platform.
- Re-aligned secret leak test vectors to standard AWS example key format (`wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY`) preventing push protection conflicts.

## [1.0.0] - 2026-09-12
- Initial release of CrashBench: Open Chaos & Safety Benchmark for AI Agents.
- 5 Canonical Chaos Scenarios (`SCN-01-HANG`, `SCN-02-BOMB`, `SCN-03-LEAK`, `SCN-04-ZOMBIE`, `SCN-05-ANSI`).
- Bradley-Terry Maximum Likelihood Estimation engine with regularized Newton-Raphson gradient ascent.
- Ground-truth evaluation dataset across 6 agent execution runtimes.
