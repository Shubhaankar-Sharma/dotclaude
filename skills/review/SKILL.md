# /review — Deep Multi-Dimensional Code Review

Comprehensive code review using /council. Every review is deep — structural analysis, behavioral traces, failure modes, code sanity, completeness. No light mode.

## Usage

```
/review                         — Review main..HEAD
/review <base>..<head>          — Review specific ref range
/review --focus=<pkg>           — Narrow to a package/directory
```

## Phases

### Phase 1: Scope

Before spawning agents, scope the changeset:

```bash
git diff --stat <base>..<head>
git log --oneline <base>..<head>
```

Identify: packages touched, files changed, line counts, number of logical changes.

### Phase 2: Invoke /council --mode=review --agents=5

Use /council with `--agents=5` to spawn 5 parallel agents. Each receives the scope summary + the full diff + their dimension-specific instructions below. The 5 dimensions below override /council's default review perspectives.

#### Dimension 1: Structural (claude sub-agent)

**Must produce before findings:**
- Call graph slice for every changed function (fan-in, fan-out, boundary classification)
- Signature delta (additions, removals, modifications with `[+]` `[-]` `[~]` markers)
- Type flow delta (where new/changed types are constructed, passed, asserted, consumed)

**Then analyze:**
- Blast radius of signature changes (high fan-in = high risk)
- New type assertions without ok checks
- Boundary violations (I/O in logic, wrong layer dependencies)
- Ownership shifts (logic moved between packages)
- Abstraction leaks (internal details in public APIs/errors)

#### Dimension 2: Behavioral (claude sub-agent)

**Must produce before findings:**
- Interaction flow traces for each significant changed code path:
  - Happy path, error path, edge case (empty/nil/single-element)
  - At each decision point: which branch taken and why, state mutations tracked
- State machine diagram (mermaid) if the change involves stateful transitions

**Then analyze:**
- State mutation ordering dependencies
- Untaken branches (dead code or missing coverage)
- Implicit pre-conditions ("works because earlier code guarantees it")
- Invariant preservation (does the change maintain documented and structural invariants)

#### Dimension 3: Failure Modes (codex exec)

Uses `codex exec --sandbox read-only` with custom prompt:

**Must enumerate:**
- Error path completeness: every new error return — who handles it, who doesn't
- Nil/zero value analysis: every new pointer/slice/map/interface — when nil, where dereferenced
- Partial failure: multi-step operations — what happens if step N fails mid-way
- Resource leaks: new goroutines, connections, file handles — termination conditions
- Concurrency: new shared state, channel ops, lock ordering

#### Dimension 4: Code Sanity (claude sub-agent)

**Must check:**
- Idiomatic Go (or relevant language) patterns
- Unnecessary complexity — can this be simpler
- Dead code introduced by the change
- Naming clarity — do names communicate intent
- Error messages — are they actionable, not generic
- Comments — do they explain WHY, not WHAT

#### Dimension 5: Completeness (claude sub-agent)

**Must check:**
- Missing test coverage for new code paths
- Missing error handling for new failure modes
- Diff minimality — is every hunk necessary for the stated purpose
- Diff coherence — does this belong in one changeset or should it be split
- Cross-cutting concerns — logging, metrics, documentation updates needed

### Phase 3: Categorize

Each agent categorizes every finding as:

- **Red** — Bugs, correctness issues, data corruption risk, invariant violations, panics, resource leaks with no mitigation. These block merge.
- **Yellow** — Semantic issues that could go wrong: implicit assumptions, missing edge cases, untested error paths, stale state risks, concurrency concerns. Should address before merge.
- **Green** — Nits, style, readability, naming, unnecessary complexity, minor improvements. Nice to have.

### Phase 4: Synthesize

/council handles synthesis. Orchestrator outputs the unified report.

## Output Format

```markdown
# Review: <changeset description>

## Scope
Files: N files in M packages
Lines: +added / -removed
Packages: [list]

## Red (Blocks Merge)
1. [structural, failure-modes] Description — file:line
   Evidence: ...

## Yellow (Address Before Merge)
1. [behavioral] Description — file:line
   Evidence: ...

## Green (Nits)
1. [code-sanity] Description — file:line

## Agent Agreement
- Finding X flagged by: structural, failure-modes, completeness (3/5)
- Finding Y flagged by: behavioral only (1/5)
```

## Rules

- This skill delegates to /council for agent spawning and synthesis
- Agents are READ-ONLY — include in all agent prompts: "Do NOT use Edit, Write, or any file-modifying tools."
- Meta-thinking is mandatory — agents must produce structural artifacts (call graphs, traces, diagrams) BEFORE producing findings
- Every finding must reference specific `file:line`
