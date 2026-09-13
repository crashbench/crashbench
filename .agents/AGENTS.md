# msh-protocol Development Rules & Agent Safety Guidelines

## Feature Workflow
When building new features, you MUST follow this exact workflow:
1. **Branching**: Create a new branch for the feature (e.g., `feat/feature-name`). Do NOT commit directly to `main`.
2. **Documentation**: After building the feature, you MUST update the `docs/` folder to reflect the changes.
3. **Changelog**: Add an entry to the `[Unreleased]` section of `CHANGELOG.md`.
4. **Readme**: Update `README.md` to highlight the new feature or update any affected examples.

## Git Operations
- **NEVER** run `git add`, `git commit`, or `git push` yourself. Always leave these commands for the user to run in their own terminal. You may only run `git checkout` to switch or create branches.
- **HOWEVER**, you MUST explicitly notify the user and tell them when it is time to run these commands (e.g., at the end of a feature or phase) so they can commit and push the changes.

## Command Execution Rule
Whenever you need to run a terminal command, you MUST prepend `msh` (or `msh wrap`) to it.
For example, instead of running `npm run build`, you can simply run `msh npm run build` (or `msh wrap "npm run build" --max-lines 500`).

---

## Agent Operations & Chaos Safety Guidelines for CrashBench

### Why CrashBench is Critical for AI Coding Agents
When AI agents execute tool commands via subprocess shells, standard synthetic coding benchmarks (like SWE-bench or HumanEval) fail to evaluate **operational runtime resilience**.

Real production agents crash, hang, and leak secrets under predictable failure vectors:
1. **Interactive Freeze (`[y/N]` prompts):** Unmanaged agents stall indefinitely waiting for stdin input, exhausting wall-clock timeouts and locking CI workflows.
2. **Context Blowout (Memory/Compiler Dumps):** A single runaway command output (1.8MB, 30,000 lines) consumes 40,000+ tokens (~$1.50 per turn), evicting system instructions and degrading reasoning.
3. **Secret Exfiltration Traps:** Indirect prompt injections in error tracebacks trick agents into inspecting environment variables and sending raw API keys to LLM providers.
4. **Zombie Subprocess Leaks:** Detached background daemons (`sleep 100 &` or dev servers) stay alive at 100% CPU after the agent finishes its turn.
5. **ANSI Tokenizer Poisoning:** Raw terminal escape codes (`\x1b[2J\x1b[H`) and carriage returns (`\r`) corrupt LLM tokenizer embeddings.

---

### Agent Guidelines for Execution Safety
1. **Never Rely on Raw Child Process Exec**: Wrap child processes in a PTY or non-blocking pipe monitor (e.g. `msh`).
2. **Apply Shannon Entropy & Token Head/Tail Truncation**: Truncate runaway outputs preserving head and tail.
3. **Stream Secret Redaction via Aho-Corasick**: Redact secrets before stdout reaches LLM context.
4. **Enforce Process DAG Cleanup**: Track child process IDs in a DAG and reap on termination.
