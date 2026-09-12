# Agent Operations & Chaos Safety Guidelines for CrashBench

## Why CrashBench is Critical for AI Coding Agents

When AI agents execute tool commands via subprocess shells, standard synthetic coding benchmarks (like SWE-bench or HumanEval) fail to evaluate **operational runtime resilience**.

Real production agents crash, hang, and leak secrets under predictable failure vectors:
1. **Interactive Freeze (`[y/N]` prompts):** Unmanaged agents stall indefinitely waiting for stdin input, exhausting wall-clock timeouts and locking CI workflows.
2. **Context Blowout (Memory/Compiler Dumps):** A single runaway command output (1.8MB, 30,000 lines) consumes 40,000+ tokens (~$1.50 per turn), evicting system instructions and degrading reasoning.
3. **Secret Exfiltration Traps:** Indirect prompt injections in error tracebacks trick agents into inspecting environment variables and sending raw API keys to LLM providers.
4. **Zombie Subprocess Leaks:** Detached background daemons (`sleep 100 &` or dev servers) stay alive at 100% CPU after the agent finishes its turn.
5. **ANSI Tokenizer Poisoning:** Raw terminal escape codes (`\x1b[2J\x1b[H`) and carriage returns (`\r`) corrupt LLM tokenizer embeddings.

---

## Agent Guidelines for Execution Safety

When building or modifying agent execution engines:

### 1. Never Rely on Raw Child Process Exec
- Always wrap child processes in a pseudo-terminal (PTY) or non-blocking pipe monitor.
- Intercept stdout for interactive prompts (e.g. `[y/N]`, `password:`, `proceed?`) and emit an automated rejection or policy decision within 500ms rather than blocking on stdin.

### 2. Apply Shannon Entropy & Token Head/Tail Truncation
- If output exceeds 1,000 lines or exhibits repetitive loop patterns (Shannon entropy $H(X) < 1.8$ bits/byte or compression ratio $CR < 0.08$), truncate automatically.
- Preserve 50 lines HEAD and 50 lines TAIL to maintain root-cause visibility while bounding prompt token consumption.

### 3. Stream Secret Redaction via Aho-Corasick
- Scan stdout byte streams using an $O(N + M)$ Aho-Corasick automaton for high-entropy credential patterns (`AKIA...`, `sk-...`, `ghp_...`, `-----BEGIN RSA PRIVATE KEY-----`).
- Redact credentials *before* stdout reaches LLM context.

### 4. Enforce Process DAG Cleanup
- Track child and grandchild process IDs in a Directed Acyclic Graph (DAG).
- On command completion or abort, perform a post-order topological traversal to terminate all leaf descendants before closing the session.

---

## Autonomous Benchmarking Workflow

Any agent can evaluate its tool runner using the CrashBench CLI:

```bash
# Run the 5-scenario lethal gauntlet
./crashbench run --target "your-agent-command" --name "AgentName" --output scorecard.json

# View Bradley-Terry Elo ratings & 95% Confidence Intervals
./crashbench elo

# View pairwise win frequencies
./crashbench matrix
```
