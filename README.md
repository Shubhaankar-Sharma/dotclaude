# dotclaude

Claude Code configuration for multi-session, multi-agent Go infrastructure work.

## What This Does

- **Worktree isolation** — Every Claude session auto-enters a git worktree. Parallel sessions never step on each other.
- **Agent cascade** — The main Claude session is an orchestrator. It never writes code directly. All work is delegated to focused sub-agents with minimal context.
- **Council of agents** — Spawn parallel agents (Claude sub-agents + Codex CLI) with different perspectives. Synthesize the best ideas.
- **Deep code review** — 5-agent parallel review: structural, behavioral, failure modes, code sanity, completeness. Red/Yellow/Green.
- **Council-informed planning** — 3 perspectives (robust, elegant, adversarial) into one plan.
- **Project docs in /tmp** — Lightweight arch/index/package docs to `/tmp/claude-docs/<project>/` for fast ramp-up.
- **AeroSpace + SketchyBar integration** — Sessions register with workspace indicators.

## Setup

```bash
git clone https://github.com/Shubhaankar-Sharma/dotclaude.git
cd dotclaude
./install.sh
```

The install script symlinks skills and hooks into `~/.claude/` and copies settings/CLAUDE.md.

## Structure

```
skills/
├── docs/SKILL.md       # /docs — Generate project docs to /tmp
├── council/SKILL.md    # /council — Multi-agent council
├── review/SKILL.md     # /review — Deep 5-agent code review
├── plan/SKILL.md       # /plan — Council-informed planning
└── formal-go/SKILL.md  # /formal-go — Invariant-first Go reasoning
hooks/
├── register-session.sh # Register session for SketchyBar
├── focus-workspace.sh  # Focus AeroSpace workspace
├── session-active.sh   # Mark active (SketchyBar)
└── session-waiting.sh  # Mark waiting (SketchyBar)
settings.json           # Permissions, hooks, plugins
CLAUDE.md               # Orchestrator instructions + guardrails
```

## Skills

### /docs
Generates codebase docs to `/tmp/claude-docs/<project>/`: `arch.md` (500 tokens), `index.md` (300 tokens), `pkg/<name>.md` (1000 tokens). Agents produce mermaid diagrams and call graphs before writing.

### /council
Spawns parallel agents with different perspectives (Claude sub-agents + Codex CLI). Each writes to `/tmp/council-XXXXXX/`. Orchestrator synthesizes: deduplicates, ranks by agreement, resolves conflicts.

### /review
5-agent parallel review via /council: structural (call graphs, type flow), behavioral (traces, invariants), failure modes (nil/error/concurrency), code sanity (idiomatic, complexity), completeness (tests, diff minimality). Red/Yellow/Green categorization.

### /plan
3-perspective planning via /council: robust (correctness), elegant (minimal), adversarial (what breaks). Synthesizes into unified plan with files, risks, test strategy.

### /formal-go
Meta-thinking for Go: separate system/environment, find atomicity boundaries, discover state machines, name invariants, map non-determinism, identify failure modes, write property tests, review checklist.

## Agent Cascade

Default behavior (not a skill). Main session is orchestrator, never writes code:

```
User → Main (orchestrator)
         ├→ reads /tmp/claude-docs/<project>/
         ├→ asks: which agent pair? (Claude+Claude, Claude+Codex, Codex+Claude)
         ├→ Pkg Agent → Impl Agent (TDD) → Review Agent
         └→ synthesizes → User
```

## Hooks (AeroSpace + SketchyBar)

Optional. Requires [AeroSpace](https://github.com/nikitabobko/AeroSpace) and [SketchyBar](https://github.com/FelixKratz/SketchyBar).

| Event | Hook | What |
|-------|------|------|
| SessionStart | register-session.sh | Register in `/tmp/claude-agents/<ws>.json` |
| Notification | focus-workspace.sh | Focus AeroSpace workspace |
| Stop | focus-workspace.sh + session-waiting.sh | Mark waiting |
| PreToolUse / UserPromptSubmit | session-active.sh | Mark active |

## Prerequisites

- [Claude Code CLI](https://docs.anthropic.com/en/docs/claude-code)
- [Codex CLI](https://github.com/openai/codex) (for council agents)
- [AeroSpace](https://github.com/nikitabobko/AeroSpace) + [SketchyBar](https://github.com/FelixKratz/SketchyBar) (optional)

## License

MIT
