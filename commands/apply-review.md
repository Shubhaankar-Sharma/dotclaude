# Apply Review Command

Apply pending review comments from the most recent review session.

## Overview

Reads review comments from branch metadata and applies them as code changes. Uses Claude to intelligently interpret and implement each comment while maintaining code style and context.

## Process

1. **Find Pending Reviews**
   - Read current branch metadata
   - Find most recent "completed" review session
   - Check if comments exist

2. **Load Context**
   - Read all comments from session
   - Group by file for efficient processing
   - Prepare context for Claude

3. **Apply Comments**
   - Use Claude to process each comment
   - Read relevant files
   - Make changes according to feedback
   - Preserve code style and formatting

4. **Update Metadata**
   - Mark session as "applied"
   - Record resolution timestamp
   - Save updated metadata

## Implementation

```bash
#!/bin/bash

set -e

# Get current branch
BRANCH=$(git branch --show-current)
if [ -z "$BRANCH" ]; then
    echo "Error: Not on a branch"
    exit 1
fi

# Get branch metadata path
METADATA_DIR="$HOME/.claude/branches"
METADATA_FILE="$METADATA_DIR/${BRANCH}.json"

if [ ! -f "$METADATA_FILE" ]; then
    echo "No metadata found for branch: $BRANCH"
    exit 1
fi

# Check if jq is available
if ! command -v jq &> /dev/null; then
    echo "Error: jq is required. Install with 'brew install jq'"
    exit 1
fi

# Find most recent completed review session
SESSION=$(jq -r '.reviewSessions | map(select(.status == "completed")) | last' "$METADATA_FILE")

if [ "$SESSION" = "null" ] || [ -z "$SESSION" ]; then
    echo "No completed review sessions found."
    echo "Run /diff-review first to create a review session."
    exit 0
fi

SESSION_ID=$(echo "$SESSION" | jq -r '.sessionId')
COMMENT_COUNT=$(echo "$SESSION" | jq -r '.comments | length')

if [ "$COMMENT_COUNT" = "0" ]; then
    echo "No comments to apply in session: $SESSION_ID"
    exit 0
fi

echo "Applying $COMMENT_COUNT comment(s) from review session: $SESSION_ID"
echo ""

# Extract comments to temp file
COMMENTS_FILE="/tmp/claude-review-apply-${SESSION_ID}.json"
echo "$SESSION" | jq '.comments' > "$COMMENTS_FILE"

# Show summary
echo "Comments to apply:"
jq -r '.[] | "  - \(.file):\(.line) - \(.content)"' "$COMMENTS_FILE"
echo ""

read -p "Proceed with applying these comments? (y/n): " CONFIRM
if [ "$CONFIRM" != "y" ] && [ "$CONFIRM" != "Y" ]; then
    echo "Cancelled."
    rm -f "$COMMENTS_FILE"
    exit 0
fi

echo ""
echo "Applying comments..."
echo ""

# Call Claude with the apply-review-comments prompt
# Claude will read the comments file and make the necessary changes
echo "Reading comments from: $COMMENTS_FILE"
echo ""
echo "Instructions for Claude:"
cat "$(dirname "$0")/../prompts/apply-review-comments.md"
echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "Please apply the comments from the JSON file above."
echo "Use the Edit tool to make changes to each file."
echo ""

# After Claude processes (this is where the conversation happens)
# User confirms when done
echo ""
read -p "Have all comments been applied? (y/n): " APPLIED

if [ "$APPLIED" = "y" ] || [ "$APPLIED" = "Y" ]; then
    # Update metadata to mark session as applied
    TIMESTAMP=$(date -u +%Y-%m-%dT%H:%M:%S.000Z)

    jq --arg sid "$SESSION_ID" \
       --arg ts "$TIMESTAMP" \
       '(.reviewSessions[] | select(.sessionId == $sid)) |= (
           .status = "applied" |
           .appliedAt = $ts
       )' "$METADATA_FILE" > "$METADATA_FILE.tmp" && mv "$METADATA_FILE.tmp" "$METADATA_FILE"

    echo ""
    echo "✓ Review comments applied successfully!"
    echo "  Session: $SESSION_ID"
    echo "  Comments: $COMMENT_COUNT"
    echo ""
    echo "Next steps:"
    echo "  1. Review changes: git diff"
    echo "  2. Run tests if needed"
    echo "  3. Commit changes: git add . && git commit"
else
    echo ""
    echo "Comments not fully applied. Session remains in 'completed' state."
    echo "Run /apply-review again when ready."
fi

# Clean up
rm -f "$COMMENTS_FILE"
```

## Example Usage

```bash
# After completing a /diff-review session with comments
$ /apply-review

Applying 3 comment(s) from review session: review-20240115-143000-abc123

Comments to apply:
  - src/handler.go:45 - Add nil check before dereferencing
  - src/handler.go:67 - Change to use buffered channel
  - src/api.go:23 - Update error message to be more specific

Proceed with applying these comments? (y/n): y

Applying comments...
[Claude makes the changes]

Have all comments been applied? (y/n): y

✓ Review comments applied successfully!
  Session: review-20240115-143000-abc123
  Comments: 3

Next steps:
  1. Review changes: git diff
  2. Run tests if needed
  3. Commit changes: git add . && git commit
```

## Notes

- Only applies comments from "completed" sessions
- Requires jq for JSON processing
- Claude interprets comments intelligently
- Preserves code style and context
- Session marked as "applied" after successful completion
- Can be run multiple times if needed
