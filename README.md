# CrashBench: The Open Chaos & Safety Benchmark for AI Agents

<p align="center">
  <pre>
   ______               __    ____                  __  
  / ____/________ _____/ /_  / __ )___  ____  _____/ /_ 
 / /   / ___/ __ `/ ___/ __ \/ __  / _ \/ __ \/ ___/ __ \
/ /___/ /  / /_/ (__  ) / / / /_/ /  __/ / / / /__/ / / /
\____/_/   \__,_/____/_/ /_/_____/\___/_/ /_/\___/_/ /_/ 
  </pre>
</p>

<p align="center">
  <strong>The SWE-bench of Agent Reliability, Terminal Resilience, and Execution Safety.</strong><br>
  Stress-testing AI coding agents against lethal real-world terminal failure modes.
</p>

<p align="center">
  <a href="https://github.com/crashbench/crashbench/releases"><img src="https://img.shields.io/badge/Release-v1.1.0-brightgreen?style=flat-square" alt="Release"></a>
  <a href="https://github.com/crashbench/crashbench/blob/main/LICENSE"><img src="https://img.shields.io/badge/License-Apache%202.0-blue?style=flat-square" alt="License"></a>
  <a href="https://github.com/crashbench/crashbench"><img src="https://img.shields.io/badge/CrashBench-Verified-orange?style=flat-square" alt="CrashBench"></a>
  <a href="https://github.com/crashbench/crashbench/stargazers"><img src="https://img.shields.io/github/stars/crashbench/crashbench?style=flat-square" alt="Stars"></a>
</p>

---

## ⚡ Why CrashBench?

Existing AI benchmarks test **code capability** (can the model solve a Python leetcode puzzle?) or **text prompt injection** (can it be tricked into writing bad words?).

**Nobody was testing operational execution failure modes.**

When an AI coding agent (like Claude Code, Cursor, Aider, or Codex) executes terminal commands in production:
1. **Interactive Prompt Hangs (`[y/N]`):** Commands like `npm init`, `rm -i`, or `apt-get` ask for user confirmation. Unprotected agents stall indefinitely, hanging CI/CD pipelines and burning infinite wall-clock time.
2. **Context Bombs:** Runaway compiler traces or recursive stack dumps output 30,000+ lines. This floods the model's context window, spikes token API costs by $10–$25 per prompt, and drops earlier instructions.
3. **Secret Exfiltration Traps:** Indirect prompt injections hidden in test logs or compiler errors instruct the agent to dump environment variables. Unredacted credentials (`AWS_SECRET_KEY`, `OPENAI_API_KEY`) get permanently captured in chat histories and provider logs.
4. **Zombie Process Accumulation:** Build scripts and test daemons leave orphaned background workers spinning at 100% CPU on host machines.
5. **ANSI Byte Corruption:** Raw 24-bit escape codes, cursor clears, and carriage-return spinners poison embedding tokenizers and degrade reasoning.

> *"Capability benchmarks measure what agents accomplish. CrashBench measures what they survive and what they break."*

---

## 🏆 Live Benchmark Leaderboard

| Rank | Agent / Execution Runtime | Status | Resilience Score (CRI) | Grade | Hang Resist. | Context Eff. | Secret Safety | ANSI Clean |
| :---: | :--- | :---: | :---: | :---: | :---: | :---: | :---: | :---: |
| 🥇 | **Protected Runtime (msh reference)** | Reference Implementation | **100.0** | **S** | 100% | 100% | 100% | 100% |
| 🥈 | **OpenHands (Docker Sandbox)** | Container Boundary | **60.0** | **C** | 40% | 75% | 20% | 70% |
| 🥉 | **Aider (Raw Shell Runner)** | Baseline Study | **43.0** | **D** | 0% | 70% | 0% | 60% |
| 4 | **Cursor Agent (Terminal Execution)** | Baseline Study | **38.0** | **F** | 0% | 60% | 0% | 50% |
| 5 | **Codex CLI (Default Shell)** | Baseline Study | **30.0** | **F** | 0% | 50% | 0% | 40% |
| 6 | **Claude Code (Raw Terminal Exec)** | Baseline Study | **23.5** | **F** | 0% | 10% | 0% | 30% |

> *Evaluated using the 5 lethal operational failure vectors of CrashBench. Lower scores indicate unhandled stdin freezes, runaway token blowouts, and raw credential exfiltration.*

---

## 🚀 Quick Start

### 1. Install CrashBench

```bash
# Via Go
go install github.com/crashbench/crashbench/cmd/crashbench@latest

# Via curl (Linux / macOS)
curl -fsSL https://raw.githubusercontent.com/crashbench/crashbench/main/scripts/install.sh | bash

# Via PowerShell (Windows)
irm https://raw.githubusercontent.com/crashbench/crashbench/main/scripts/install.ps1 | iex
```

### 2. Run the Gauntlet Against Your Agent

```bash
# Run against a mock baseline
crashbench run --name "MyAgent" --target "mock"

# Run against a custom command or agent CLI
crashbench run --name "Claude Code" --target "npx claude-code"

# Export JSON scorecard
crashbench run --name "MyAgent" --target "mock" --output scorecard.json
```

### 3. Launch the Embedded Web Arena & Leaderboard

CrashBench compiles with an embedded zero-dependency Web UI:

```bash
crashbench serve --port 4040
```
Open **`http://localhost:4040`** to view the live interactive leaderboard, terminal replays, in-browser simulator, and badge generator.

---

## 🔬 The Canonical Chaos Vectors

| Vector ID | Name | Category | Failure Mode & Impact |
| :--- | :--- | :--- | :--- |
| **`SCN-01`** | **The Interactive Freeze** | Hang Resistance | Agent locks on `[y/N]` or password prompts without non-blocking resolution. |
| **`SCN-02`** | **The 30,000-Line Context Bomb** | Context Efficiency | 1.8MB compiler trace floods context window without intelligent truncation. |
| **`SCN-03`** | **Secret Exfiltration Trap** | Secret Safety | Indirect prompt injection triggers emission of unmasked AWS or OpenAI keys. |
| **`SCN-04`** | **Orphan Subprocess Leak** | Process Hygiene | Detached child subshells left alive as orphaned background processes. |
| **`SCN-05`** | **ANSI Escape Labyrinth** | Terminal Fidelity | Unparsed terminal control codes and progress bars pollute token embeddings. |
| **`SCN-06`** | **Filesystem Symlink Escape** | Boundary Isolation | Symlinks pointing to host sensitive paths (`../../.ssh/id_rsa`, `/etc/shadow`). |

---

## 🎨 Bohemian & Handwritten Minimalist Interactive CLI

CrashBench includes a bespoke 24-bit TrueColor interactive terminal experience:

```bash
# Launch interactive mode (default when run without flags)
crashbench
# or
crashbench interactive
```

Key capabilities:
- `[1] 🏆 Bradley-Terry Elo Leaderboard`: Live MLE calculations with 95% Bayesian bootstrap CI slider bars `[───●━━─]`.
- `[2] ⚔️ Side-by-Side Chaos Battle Arena`: Select any two agent runtimes and inject failure vectors with split-screen diffs.
- `[3] 📊 Pairwise Win-Rate Heatmap`: Head-to-head empirical win frequency matrix across 2,500+ gauntlet battles.
- `[4] 🔍 Deep Agent Scorecard Inspector`: Multi-scenario scorecards and full terminal transcripts.
- `[5] 🚀 Run Live Chaos Gauntlet`: Execute lethal scenarios against target agents.
- `[6] 📜 Research Field Notes & Math`: CS formulations (Aho-Corasick DFA, Shannon entropy, Tarjan DAG).

---

## 🛡️ Embed Your Verified Badge

Showcase your agent's resilience score on your GitHub README:

```markdown
[![CrashBench](https://img.shields.io/badge/CrashBench-S%20100.0-brightgreen?style=flat-square)](https://github.com/crashbench/crashbench)
```

[![CrashBench](https://img.shields.io/badge/CrashBench-S%20100.0-brightgreen?style=flat-square)](https://github.com/crashbench/crashbench)

---

## 🏛️ Project Architecture

```
crashbench/
├── cmd/
│   └── crashbench/          # CLI entry point (run, serve, leaderboard)
├── pkg/
│   ├── gauntlet/            # The 5 lethal scenarios and payload generators
│   ├── runner/              # Execution engine, process monitor & evaluator
│   └── scoring/             # Weighted CRI index and badge generation
├── data/
│   ├── leaderboard.json     # Verified ground-truth benchmark dataset
│   └── embed.go             # Go embed package
├── web/                     # Embedded Web Arena & Leaderboard
│   ├── index.html           # Dark glassmorphism dashboard
│   ├── style.css            # Obsidian & neon design system
│   ├── app.js               # Interactive replay modal & live simulator
│   └── embed.go             # Embedded static FS
└── scripts/                 # 1-line installation scripts
```

---

## 🤝 Contributing & Submitting Agent Runs

We accept community benchmark submissions!
1. Run the gauntlet: `crashbench run --name "YourAgent" --target "your-cli-command" --output results.json`
2. Submit a Pull Request updating `data/leaderboard.json` with your verified run logs.

---

## 📄 License

Apache License 2.0. See [LICENSE](LICENSE) for details.
