# Diff Review Command

Launch interactive visual diff review TUI for code review with Claude.

## Overview

Opens a terminal-based interactive diff viewer that allows you to:
- Navigate through git diffs with syntax highlighting
- Ask Claude questions about specific code sections (press 'q')
- Add review comments for batch processing (press 'c')
- Request edits during conversation
- Apply all collected feedback at once

## Workflow

1. **Initialize Review Session**
   - Get current branch metadata path
   - Create new review session with unique ID
   - Store base and head commits
   - Save session to metadata JSON

2. **Launch TUI**
   - Execute Go binary: `scripts/diff-review`
   - Pass session ID and commit range as arguments
   - TUI handles all user interaction

3. **Handle TUI Exit Codes**
   - Exit 0: Normal quit, no action
   - Exit 1: Error occurred
   - Exit 2: Quit without applying
   - Exit 3: Question mode - Claude conversation needed

4. **Question Mode** (Exit code 3)
   - Read question context from `/tmp/claude-review-context-<session-id>.json`
   - Parse: file path, line range, selected diff lines
   - Read ACTUAL source file (not just diff)
   - Show ±20 lines of context around selection
   - Allow user to ask questions or request edits
   - After conversation, ask if user wants to return to TUI
   - If yes: relaunch TUI with same session
   - If no: end session

5. **Apply Mode** (Triggered by 'a' in TUI)
   - Comments are already saved to metadata
   - Run `/apply-review` command to process

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
mkdir -p "$METADATA_DIR"
METADATA_FILE="$METADATA_DIR/${BRANCH}.json"

# Get git commits
BASE_COMMIT=$(git merge-base HEAD main 2>/dev/null || git rev-parse HEAD~1)
HEAD_COMMIT=$(git rev-parse HEAD)

# Generate session ID
SESSION_ID="review-$(date +%Y%m%d-%H%M%S)-$(head -c 4 /dev/urandom | xxd -p)"

# Initialize metadata if doesn't exist
if [ ! -f "$METADATA_FILE" ]; then
    cat > "$METADATA_FILE" <<EOF
{
  "branch": "$BRANCH",
  "type": "feature",
  "reviewSessions": []
}
EOF
fi

# Add new review session to metadata using jq
if command -v jq &> /dev/null; then
    jq --arg sid "$SESSION_ID" \
       --arg base "$BASE_COMMIT" \
       --arg head "$HEAD_COMMIT" \
       --arg ts "$(date -u +%Y-%m-%dT%H:%M:%S.000Z)" \
       '.reviewSessions += [{
           "sessionId": $sid,
           "timestamp": $ts,
           "baseCommit": $base,
           "headCommit": $head,
           "status": "in_progress",
           "comments": [],
           "conversations": []
       }]' "$METADATA_FILE" > "$METADATA_FILE.tmp" && mv "$METADATA_FILE.tmp" "$METADATA_FILE"
else
    echo "Warning: jq not found. Install with 'brew install jq' for metadata tracking."
fi

# Launch TUI
TUI_BINARY="$(dirname "$0")/../scripts/diff-review"

if [ ! -x "$TUI_BINARY" ]; then
    echo "Error: TUI binary not found or not executable at $TUI_BINARY"
    echo "Build it with: cd scripts/diff-review-tui && go build -o ../diff-review ."
    exit 1
fi

# Loop to handle question mode
while true; do
    "$TUI_BINARY" --session-id "$SESSION_ID" \
                  --base-commit "$BASE_COMMIT" \
                  --head-commit "$HEAD_COMMIT" \
                  --metadata-path "$METADATA_FILE"

    EXIT_CODE=$?

    case $EXIT_CODE in
        0)
            # Normal quit
            echo "Review session ended."
            break
            ;;
        1)
            # Error
            echo "Error occurred during review."
            exit 1
            ;;
        2)
            # Quit without applying
            echo "Review quit without applying changes."
            break
            ;;
        3)
            # Question mode - handle conversation
            CONTEXT_FILE="/tmp/claude-review-context-${SESSION_ID}.json"

            if [ ! -f "$CONTEXT_FILE" ]; then
                echo "Error: Question context file not found"
                exit 1
            fi

            # Parse context file
            FILE=$(jq -r '.file' "$CONTEXT_FILE")
            START_LINE=$(jq -r '.startLine' "$CONTEXT_FILE")
            END_LINE=$(jq -r '.endLine' "$CONTEXT_FILE")

            # Show context from actual file
            echo ""
            echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
            echo "Question about: $FILE"
            echo "Lines: $START_LINE-$END_LINE"
            echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
            echo ""

            # Show ±20 lines of context
            CONTEXT_START=$((START_LINE - 20))
            CONTEXT_END=$((END_LINE + 20))
            [ $CONTEXT_START -lt 1 ] && CONTEXT_START=1

            echo "Context from source file:"
            echo ""
            sed -n "${CONTEXT_START},${CONTEXT_END}p" "$FILE" | nl -v $CONTEXT_START -w 4 -s ": " | \
                awk -v start=$START_LINE -v end=$END_LINE '{
                    line_num = substr($0, 1, 4);
                    gsub(/^[ \t]+/, "", line_num);
                    line_num = int(line_num);
                    if (line_num >= start && line_num <= end) {
                        print $0 " [SELECTED]"
                    } else {
                        print $0
                    }
                }'

            echo ""
            echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
            echo ""
            echo "What would you like to know about this code?"
            echo "(You can ask questions, request edits, or explore related code)"
            echo ""

            # User can now interact with Claude
            # Claude has full codebase access and can make edits

            # After conversation, ask about returning to TUI
            echo ""
            read -p "Return to diff review? (y/n): " RETURN_TO_TUI

            if [ "$RETURN_TO_TUI" != "y" ] && [ "$RETURN_TO_TUI" != "Y" ]; then
                echo "Ending review session."
                break
            fi

            echo ""
            echo "Returning to TUI... (press 'r' to refresh if Claude made changes)"
            sleep 1
            # Loop continues, TUI relaunches
            ;;
        *)
            echo "Unknown exit code: $EXIT_CODE"
            exit 1
            ;;
    esac
done

# Clean up context file
rm -f "/tmp/claude-review-context-${SESSION_ID}.json"

echo ""
echo "Review session complete!"
echo "Session ID: $SESSION_ID"
echo ""
echo "To apply pending comments, run: /apply-review"
```

## Keybindings in TUI

- `j/k` or `↑/↓`: Navigate through diff
- `v`: Enter visual mode (select multiple lines)
- `q`: Ask question about current line (exit to Claude)
- `Q`: Ask question in visual mode (about selected lines)
- `c`: Add comment for later batch processing
- `s`: Show summary of all comments
- `a`: Apply all pending comments (saves and exits)
- `r`: Refresh diff (after Claude makes edits)
- `Esc`: Cancel current mode
- `Ctrl+C`: Quit review

## Notes

- TUI uses Chroma for syntax highlighting
- All review data persists in branch metadata JSON
- Questions give Claude full codebase access
- Edits made during conversation show up on refresh
- Binary file: `scripts/diff-review` (built from Go source)
