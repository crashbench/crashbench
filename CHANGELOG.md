# Changelog

All notable changes to CrashBench will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

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
