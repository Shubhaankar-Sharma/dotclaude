## Worktrees
- At session start, if inside a git repository, ALWAYS enter a worktree before doing any work
- This isolates parallel Claude sessions from stepping on each other
- After entering the worktree, re-register the session by running: `~/.claude/hooks/register-session.sh`

## Project Name Resolution
- ALWAYS resolve project name via:
  ```bash
  _gc=$(git rev-parse --git-common-dir 2>/dev/null)
  if [ -n "$_gc" ] && [ "$_gc" != ".git" ]; then
    PROJECT=$(basename "$(cd "$_gc/.." && pwd)")
  else
    PROJECT=$(basename "$(pwd)")
  fi
  ```
- This works inside worktrees (resolves to real repo name) and in normal repos (falls back to dirname)
- Use this for `/tmp/claude-docs/<project>/` lookups and any project-name-dependent paths

## Agent Protocol
- You are an orchestrator. NEVER write code or file artifacts directly in main context.
- ALL implementation and doc generation is delegated to focused sub-agents. Less context = more focus.
- For any code question, first check `/tmp/claude-docs/<project>/` — read `arch.md` + `index.md` to ramp up.
- If docs are missing or stale, spawn a sub-agent to generate them (via /docs).
- For planning, use /plan which invokes /council with multiple perspectives.
- For implementation, ask the user which agent pair to use:
  1. Claude + Claude (default, fastest)
  2. Claude + Codex (codex review — codex gives great reviews)
  3. Codex + Claude (sandboxed impl, use `codex apply` to extract changes)

## Implementation Agent Template
When spawning an impl sub-agent, include in its prompt:
- Task description and scope (which files to modify)
- Relevant docs from `/tmp/claude-docs/<project>/pkg/<name>.md`
- TDD requirement: write failing test first, then implement
- File scope constraint: "Only modify these files: [list]"
- "Run `go test ./...` before reporting success"
- "Report back: what you changed, what tests you added, what's still unclear"

When spawning a review sub-agent:
- "You are in READ-ONLY mode. Do NOT use Edit, Write, or any file-modifying tools."
- Give it the diff and the task description
- "Categorize findings as Red (bugs) / Yellow (risks) / Green (nits)"
- Max 3 impl/review iterations. If still failing after 3, escalate to user.

## Review Mode
- When asked to review, use /review — it delegates to /council for 5-agent parallel review
- /review is for code review. Agent-pair selection is for implementation. These are separate flows.
- NEVER implement fixes during review unless explicitly asked
- Categorize findings as Red (bugs) / Yellow (risks) / Green (nits)

## Go Development
- Primary language is Go
- For Go projects, apply /formal-go reasoning when writing or reviewing code that involves state machines, concurrency, or crash recovery
- Use `go mod tidy` to resolve dependency/toolchain issues before suggesting manual upgrades
- Run tests after changes to ensure nothing breaks

## Code Changes
- NEVER re-add code the user has explicitly removed
- NEVER delete files without explicit confirmation
- Minimize blast radius — smallest change that solves the problem
- When making multi-file changes, list all files to be modified and get approval first

## Git Operations
- NEVER push without explicit user approval
- NEVER add co-authors unless asked
- NEVER commit generated/temporary files — check .gitignore first
- When unsure about git push/pull direction, ASK before executing

## Approach
- When user says "review" or "analyze", provide analysis — do NOT implement fixes unless asked
- Do NOT start implementing until the user confirms the plan/approach
- Do NOT start autonomous exploration unless explicitly asked
- Ask before making changes that exceed the scope of what was requested
