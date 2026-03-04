# /plan — Council-Informed Planning

Generate implementation plans using /council with different perspectives. Synthesize best ideas into a unified plan.

## Usage

```
/plan <task description>
```

## Behavior

1. Resolve project name (works in worktrees and normal repos):
   ```bash
   _gc=$(git rev-parse --git-common-dir 2>/dev/null)
   if [ -n "$_gc" ] && [ "$_gc" != ".git" ]; then
     PROJECT=$(basename "$(cd "$_gc/.." && pwd)")
   else
     PROJECT=$(basename "$(pwd)")
   fi
   ```
2. Read `/tmp/claude-docs/<project>/arch.md` + `index.md` for codebase context
3. If docs don't exist, spawn a sub-agent to generate them first (via /docs)
4. Invoke `/council --mode=plan` with the task + codebase context

## Perspectives (passed to /council)

### Agent 1: Most Robust (claude sub-agent)
- Prioritize correctness and testability
- Cover all edge cases and failure modes
- Design for crash recovery and idempotency
- Propose comprehensive test strategy
- Ask: "What breaks if any assumption is wrong?"

### Agent 2: Most Elegant (codex exec)
- Prioritize clean architecture and minimal surface area
- Find the simplest abstraction that covers all cases
- Minimize files changed and blast radius
- Reuse existing patterns in the codebase
- Ask: "What's the smallest change that solves this completely?"

### Agent 3: Adversarial (claude sub-agent)
- Actively try to break the other approaches
- Identify assumptions that might not hold
- Find edge cases the other agents missed
- Question whether the task is scoped correctly
- Ask: "What's wrong with this plan? What's missing?"

## Each Agent Outputs

```markdown
# Plan: <perspective name>

## Understanding
What the task requires (verify alignment)

## Approach
Numbered steps with rationale

## Files to Modify
- path/to/file.go — what changes and why

## Risks
- Risk 1: description + mitigation

## Test Strategy
- What tests to write
- What edge cases to cover

## Open Questions
- Anything unclear or that needs user input
```

## Synthesis

Handled by /council's synthesis protocol, plus:

1. **Alignment check** — do all agents agree on what the task requires? If not, flag the disagreement.
2. **Approach merge** — robust agent's test strategy + elegant agent's implementation + adversarial agent's risk mitigations.
3. **File list union** — combine all files-to-modify, flag any that only one agent identified.
4. **Risk union** — all risks from all agents, deduplicated.

## Output Format

```markdown
# Plan: <task>

## Context
<what project docs say about the relevant packages>

## Approach
<synthesized steps from all agents>

## Files
| File | Change | Source |
|------|--------|--------|
| path.go | description | robust + elegant |

## Risks
1. Risk — mitigation (flagged by: adversarial)

## Test Plan
<from robust agent, validated by adversarial>

## Open Questions
<unresolved items needing user input>
```

## Rules

- Delegates to /council for agent spawning and synthesis
- Agents are READ-ONLY — no code changes during planning
- After synthesis, WAIT for user approval before any implementation
- If agents fundamentally disagree on approach, present both with tradeoffs and let user decide
