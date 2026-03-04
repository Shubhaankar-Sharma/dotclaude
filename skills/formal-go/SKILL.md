# Formal Go

Meta-thinking framework for writing Go where bugs cause real damage.
Not a library. Not domain docs. A reasoning process that makes correct code
the only code that compiles and survives tests.

## The One Idea

Every system has invariants. Most are implicit. Implicit invariants are invisible
to AI and are the #1 source of subtle bugs in generated code. The entire framework
reduces to: **find the invariants, write them down, enforce them mechanically.**

## Phase 1 — Separate System from Environment

Before touching any logic, draw the boundary between what the system *does*
(protocol, business logic, state transitions) and what the system *talks to*
(network, storage, caches, clocks, other services).

**The system** is pure logic: given this state and this input, produce this
output and this new state. It makes no network calls, no disk writes, no
RPCs. It is deterministic and trivially testable.

**The environment** is everything the system interacts with. Each environmental
dependency gets a narrow interface — read, write, and critically: abort/rollback
semantics.

```go
type Resource interface {
    Read(ctx context.Context, key Key) (Value, error)
    Write(ctx context.Context, key Key, val Value) error
    PreCommit(ctx context.Context) error
    Commit(ctx context.Context) error
    Abort(ctx context.Context)
}
```

**Heuristic:** If you see `http.Get` or `sql.Query` inside a state transition
function, the boundary is violated. Extract it behind a Resource interface.

## Phase 2 — Find the Atomicity Boundaries

Ask: "If another goroutine/process/node observes the world between step A and
step B of this operation, does something break?" If yes, A and B must be in
the same critical section.

Each critical section becomes:
1. A unit of rollback — if any step fails, ALL steps are undone
2. A unit of serialization — no other critical section on the same state can interleave
3. A unit of crash recovery — a crash mid-section means the section never happened

```go
func (s *System) HandleRequest(iface ResourceSet) error {
    msg, err := iface.Read(networkResource, selfID)
    if err != nil { return err }
    // ... pure protocol logic ...
    if err := iface.Write(networkResource, destID, response); err != nil { return err }
    if err := iface.PreCommit(); err != nil {
        iface.Abort()
        return err
    }
    iface.Commit()
    return nil
}
```

**Heuristic:** If a function does two writes to different resources and there's
no atomicity story (2PC, saga, or explicit "this is best-effort"), it's a bug
waiting for a crash to expose it.

## Phase 3 — Discover the State Machine

Ask: "What are the possible states of this thing?" Then: "Which transitions
between states are legal?" Then: "What guards must be true for each transition?"

```go
var transitions = map[State]map[State]GuardFunc{ ... }

func (e *Entity) Transition(next State) error {
    guards, ok := transitions[e.state]
    if !ok { return fmt.Errorf("terminal state: %v", e.state) }
    guard, ok := guards[next]
    if !ok { return fmt.Errorf("illegal: %v → %v", e.state, next) }
    if err := guard(e); err != nil {
        return fmt.Errorf("%v → %v blocked: %w", e.state, next, err)
    }
    e.state = next
    return nil
}
```

The transition map IS the spec.

## Phase 4 — Name the Invariants

Invariants are predicates that must be true in EVERY reachable state.

Common shapes:
- **Uniqueness**: Exactly one X in a scope (one primary per shard)
- **Conservation**: Quantity that must not change
- **Ordering**: A before B, never after
- **Mutual exclusion**: X and Y never both true
- **Consistency**: Two representations agree
- **Boundedness**: Value stays within range

```go
func (e *Entity) CheckInvariants() error {
    var errs []string
    if violatesUniqueness(e) { errs = append(errs, "...") }
    if violatesConservation(e) { errs = append(errs, "...") }
    return combineErrors(errs)
}

// mutation gateway — no mutations outside Apply
func (e *Entity) Apply(op func(*Entity) error) error {
    if err := op(e); err != nil { return err }
    return e.CheckInvariants()
}
```

## Phase 5 — Map the Non-Determinism

Look for: `select` statements, retry logic, timeout-driven branching,
order-dependent processing, "pick one" decisions.

Each non-deterministic choice should be:
1. **Enumerable** — you can list the possible outcomes
2. **Safe under all outcomes** — invariants hold regardless of which path is taken
3. **Testable** — property tests explore random choices, not just the happy path

## Phase 6 — Identify the Failure Modes

### For operations in loops (reconcile, retry, watch)
- **Is it idempotent?** Calling twice with same input = same result.
- **Is it level-triggered?** Reacts to current state, not transitions.
- **Does it assume it saw every event?** It won't. Must converge from ANY current state.

### For operations that span steps (workflows, migrations)
- **What if we crash between step N and N+1?** Next restart must resume without re-doing N.
- **What gets cleaned up on cancel?** Every side effect must be reversible or harmless on abort.

### For anything using cached/stale state
- **What if the cache is wrong?** Must degrade safely (error) not silently act on stale data.

### For anything coordinating across nodes
- **What if the lock expires mid-operation?** Must verify lock is still held before mutating.
- **What if both sides think they're the leader?** Fencing must be mechanical, not logical.
- **Are positions/versions compared correctly?** Typed comparison, not string.

## Phase 7 — Write Properties, Not Examples

Use `pgregory.net/rapid`.

```go
rapid.Check(t, func(t *rapid.T) {
    state := generateRandomValidState(t)
    ops := generateRandomOperationSequence(t)
    for i, op := range ops {
        _ = state.Apply(op)
        require.NoError(t, state.CheckInvariants(),
            "invariant violated after op %d: %s", i, op)
    }
})
```

Crash-recovery:
```go
for _, crashAt := range allPhases {
    wf := setupWorkflow()
    wf.RunUntil(crashAt)
    recovered := recoverFromPersistedState(wf.PersistedState())
    require.NoError(t, recovered.Run())
    require.NoError(t, recovered.CheckInvariants())
}
```

## Phase 8 — Reviewing AI-Generated Code

In priority order. Stop and fix before continuing.

1. **Is system separated from environment?** Protocol logic calling I/O directly = boundary violation.
2. **Are atomicity boundaries explicit?** All-or-nothing operations in a single critical section with rollback.
3. **Is there an implicit state machine?** Status strings, phase booleans = missing explicit machine.
4. **Are invariants written and enforced?** Must be code, not comments.
5. **Is it idempotent?** Anything in a loop. Call it twice mentally.
6. **What's the crash story?** Process dies mid-function. Can next run recover?
7. **What's stale?** Any data from cache/watch/remote may be stale.
8. **Are all choices safe?** Non-deterministic paths must all preserve invariants.
9. **Would a wrong implementation pass these tests?** If all example-based, add properties.

## Decision Triggers

| You see...                                  | You should...                                     |
|---------------------------------------------|---------------------------------------------------|
| I/O mixed into protocol logic               | Separate system from environment via Resource      |
| Two writes with no atomicity story           | Define critical section with rollback              |
| String/bool encoding state                  | Extract explicit state machine with transition map |
| Mutable struct with no invariant checker     | Write `CheckInvariants() error`                   |
| Operation in a loop without idempotency test | Add idempotency property test                     |
| Multi-step process without crash test        | Test crash at each step boundary                  |
| Cached data used without staleness test      | Add stale-cache degradation test                  |
| Lock acquired without re-verification        | Add verify step before guarded mutation           |
| Positions/versions compared as strings       | Use typed comparison                              |
| Tests only cover happy-path choices          | Add non-determinism property tests                |
