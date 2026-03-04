# /council — Council of Agents

Spawn parallel agents with different perspectives. Collect outputs. Synthesize best ideas.

## Usage

```
/council <task>                    — Default: auto-select perspectives
/council --mode=review <task>      — Review-oriented perspectives
/council --mode=plan <task>        — Planning-oriented perspectives
/council --agents=3 <task>         — Control agent count (2-5)
```

## How It Works

1. Parse the task and mode
2. Create output directory: `mktemp -d /tmp/council-XXXXXX`
3. Spawn agents in parallel — each gets the same task + a unique perspective
4. Each agent writes findings to `<session-dir>/<agent-name>.md`
5. Wait for all agents to complete
6. Read all output files
7. Synthesize: take best ideas from each, resolve conflicts, deduplicate, produce unified output

## Project Name Resolution

To determine the project name (for docs lookup, etc.), use:
```bash
_gc=$(git rev-parse --git-common-dir 2>/dev/null); if [ -n "$_gc" ] && [ "$_gc" != ".git" ]; then basename "$(cd "$_gc/.." && pwd)"; else basename "$(pwd)"; fi
```
This resolves to the real repo name even inside a worktree.

## Agent Mix

| Agent | CLI | How |
|-------|-----|-----|
| claude-1 | Claude sub-agent | `Agent` tool with perspective prompt |
| claude-2 | Claude sub-agent | `Agent` tool with different perspective |
| claude-3 | Claude sub-agent | `Agent` tool with third perspective (if 5 agents) |
| codex | Codex CLI | `codex exec --sandbox read-only -C <dir> <prompt>` (run in background via Bash) |

Claude sub-agents write output via the Agent tool response. Codex writes to stdout, redirected to file.

### Pre-flight check

Before spawning codex, verify it's available:
```bash
which codex && codex exec --help >/dev/null 2>&1
```
If codex is unavailable, replace with an additional Claude sub-agent.

## Perspectives by Mode

### Review mode
1. **Correctness** — bugs, logic errors, invariant violations, type safety
2. **Failure modes** — nil paths, error handling gaps, crash scenarios, resource leaks
3. **API design** — interface clarity, abstraction leaks, naming, contracts
4. **Adversarial** — what breaks under load, concurrency, malicious input, stale state
5. **Completeness** — missing tests, unhandled edge cases, diff minimality

### Plan mode
1. **Most robust** — handles all edge cases, most testable, failure modes covered
2. **Most elegant** — best architecture, cleanest abstractions, minimal surface area
3. **Adversarial** — what breaks, what's missing, what assumption is wrong

### Default mode
Auto-selects perspectives based on task content.

## CLI Commands

```bash
# Session directory (collision-resistant)
SID=$(mktemp -d /tmp/council-XXXXXX)

# Codex (non-interactive, read-only for review, workspace-write for impl)
# Run in background via Bash tool with run_in_background=true
codex exec --sandbox read-only -C "$PROJECT_DIR" "$(cat <<'PROMPT'
<perspective directive>
<task>

Categorize findings as Red (bugs/correctness), Yellow (risks/assumptions), Green (nits/style).
PROMPT
)" > $SID/codex.md 2>&1
```

## Synthesis Protocol

After collecting all outputs:

1. **Read** all `<session-dir>/*.md` files
2. **Deduplicate** — same finding from multiple agents = higher confidence
3. **Conflict resolution** — when agents disagree, present both views with reasoning
4. **Rank** — findings agreed on by 3+ agents ranked highest
5. **Output** — unified result with attribution (which agents flagged what)

## Rules

- Codex runs as a background Bash task
- Claude sub-agents run via Agent tool (parallel)
- ALL agents are read-only unless mode explicitly requires writes. For Claude sub-agents, include in prompt: "You are in READ-ONLY mode. Do NOT use Edit, Write, or any file-modifying tools."
- Orchestrator NEVER produces findings of its own — only synthesizes
- If codex fails (auth, crash), note the gap and proceed with Claude agents
- Always report which agents contributed to each finding
