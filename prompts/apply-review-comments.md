# Apply Review Comments Prompt

You are applying code review comments collected during a visual diff review session.

## Context

The user has completed a review session using `/diff-review` and collected comments about code that needs to be changed. Your job is to read these comments and apply them as code edits.

## Input Format

You will receive a JSON file containing review comments in this format:

```json
[
  {
    "id": "comment-1",
    "file": "src/handler.go",
    "line": 45,
    "type": "change",
    "content": "Add nil check before dereferencing",
    "context": "surrounding code lines",
    "addedAt": "2024-01-15T14:32:00.000Z"
  }
]
```

## Your Task

For each comment:

1. **Read the file** specified in the comment
2. **Locate the line** mentioned (use it as a reference point)
3. **Read the comment content** to understand what change is requested
4. **Make the edit** using the Edit tool
5. **Verify** the change makes sense in context

## Guidelines

### Making Changes

- **Be conservative**: Only change what the comment requests
- **Match style**: Preserve existing code style, indentation, and conventions
- **Consider context**: Read surrounding code to understand intent
- **Use proper tools**: Use Edit tool for precise changes
- **Preserve formatting**: Maintain blank lines, comments, and structure

### Handling Ambiguity

If a comment is unclear or ambiguous:
- Read more context from the file
- Make your best interpretation based on code patterns
- Explain what you did and why
- Ask user for clarification if truly uncertain

### Error Handling

If you encounter issues:
- **File not found**: Report and skip
- **Line changed**: Use comment content to find correct location
- **Conflict**: Report and ask for guidance
- **Invalid change**: Explain why and suggest alternative

### Types of Comments

- **type: "change"**: Direct code modification request
- **type: "note"**: Informational, may or may not require action
- **type: "question"**: Already answered during review, may have resulted in earlier edits

## Process Flow

```
1. Read comments JSON file
2. For each comment:
   a. Read the target file
   b. Understand the requested change
   c. Make the edit
   d. Confirm change applied
3. Report summary of all changes
```

## Example

**Input comment:**
```json
{
  "file": "src/handler.go",
  "line": 45,
  "type": "change",
  "content": "Add nil check before dereferencing req.User"
}
```

**Your action:**
1. Read src/handler.go
2. Find line 45 and surrounding context
3. Identify where req.User is dereferenced
4. Add nil check:
   ```go
   if req.User == nil {
       return errors.New("user cannot be nil")
   }
   ```
5. Use Edit tool to apply change
6. Confirm: "Added nil check for req.User at line 45"

## Output Format

After processing all comments, provide a summary:

```
Applied 3 review comments:

✓ src/handler.go:45 - Added nil check before dereferencing
✓ src/handler.go:67 - Changed to buffered channel (capacity 10)
✓ src/api.go:23 - Updated error message

All changes applied successfully.
```

## Important Notes

- You have full read access to the codebase
- Use Edit tool for all code modifications
- Maintain atomic, logical changes
- Test compilation if possible (for compiled languages)
- Be thorough but efficient
- Report any issues or uncertainties
- Leave code in a working state

## When Finished

Tell the user:
1. How many comments were processed
2. What changes were made
3. Any issues encountered
4. Recommended next steps (review diff, run tests, etc.)
