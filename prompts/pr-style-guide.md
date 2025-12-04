# PR Description Writing Style

## Core Principles
- **Expressive**: Use clear, meaningful language that explains WHY, not just WHAT
- **Dense**: Pack information into fewer words, no fluff
- **Minimal**: Short but complete
- **Technical**: Assume the reader knows the codebase

## Writing Rules

1. **Avoid hollow words** like "identifies", "entities", "coordinates", "implements" without context
   - These words sound technical but say nothing about what actually happens

2. **Explain function calls in terms of what they DO**, not just what they're called
   - ❌ Bad: "findAffectedTables() identifies which tables reference the changed entities"
   - ✅ Good: "findAffectedTables() loops through all tables checking if they use any changed ShardGroups"

3. **Connect the flow**: "X calls Y, which does Z"
   - ❌ Bad: "collectTopologyDiff() computes changes"
   - ✅ Good: "notifyAffectedTables() calls collectTopologyDiff() to compare old vs new topologies"

4. **Use concrete examples** in parentheses
   - Example: (e.g., xxhash → range, disabling a global index)
   - Example: (e.g., switching from xxhash to range partitioning)

5. **Reference code locations** as: `file.go:123` or `file.go:123-456`
   - Not as links, just as text references
   - Keep them short and readable

6. **Keep Background and Effect sections dense but readable** (2-3 paragraphs each)
   - Background: What existed, why change is needed, what problem it solves
   - Effect: What changes, concrete impact, how it works

7. **Implementation details should be bullet points** with file:line references
   - Each bullet connects to the next
   - Show the flow of execution
   - Explain what happens and why, not just function names

## Bad vs Good Examples

### Example 1: Hollow vs Expressive

**❌ Bad (hollow):**
"The manager identifies affected entities and coordinates invalidation"

**✅ Good (expressive):**
"The manager diffs old vs new topologies, finds tables using changed ShardGroups, then invokes the callback for each one"

### Example 2: Too technical vs Clear

**❌ Bad (too technical):**
"collectTopologyDiff() implements a differential analysis algorithm"

**✅ Good (clear):**
"collectTopologyDiff() loops through ShardGroups comparing them with Equals(), tracking which ones changed"

### Example 3: Vague vs Specific

**❌ Bad (vague):**
"This improves performance by optimizing the query planner"

**✅ Good (specific):**
"Cached plans now stay valid across topology rebuilds unless their specific ShardGroup changes, avoiding unnecessary replanning"

### Example 4: Implementation Details

**❌ Bad:**
- `manager.go:100` - Callback triggered
- `datatopology_diff.go:21-29` - Coordinates invalidation
- `datatopology_diff.go:44-78` - Computes changes

**✅ Good:**
- `manager.go:100` - After rebuilding the topology snapshot, we call `notifyAffectedTables(oldDt, dt)` to trigger invalidation
- `datatopology_diff.go:21-29` - `notifyAffectedTables()` computes which tables are affected, then loops through them invoking the callback for each one
- `datatopology_diff.go:44-78` - `collectTopologyDiff()` compares old vs new topologies by checking each ShardGroup and SecondaryIndex for changes using their `Equals()` methods

## Summary

Write like you're explaining to a colleague who knows the codebase:
- Skip the buzzwords
- Explain what actually happens
- Use concrete examples
- Show the connections between steps
- Be precise but not verbose
