#!/bin/bash

# Session Start Hook: Auto-compact history and trigger Claude to load context
# This hook runs when Claude Code starts a new session or resumes an existing one

# Get current git branch
BRANCH=$(git branch --show-current 2>/dev/null)

# If not in a git repo, exit silently
if [ -z "$BRANCH" ]; then
    exit 0
fi

# Sanitize branch name for filename (replace / with -)
BRANCH_FILE=$(echo "$BRANCH" | tr '/' '-')

# Ensure .claude/branches directory exists
mkdir -p ".claude/branches"

# Check for metadata file in the repo
METADATA_FILE=".claude/branches/${BRANCH_FILE}.json"

if [ ! -f "$METADATA_FILE" ]; then
    # Create initial metadata file if it doesn't exist
    TIMESTAMP=$(date -u +"%Y-%m-%dT%H:%M:%S.000Z")
    COMMIT_COUNT=$(git rev-list --count HEAD ^origin/main 2>/dev/null || echo "0")
    LAST_COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "")

    cat > "$METADATA_FILE" <<EOF
{
  "branch": "$BRANCH",
  "type": "",
  "issue": "",
  "issueUrl": "",
  "description": "",
  "started": "$TIMESTAMP",
  "lastWorked": "$TIMESTAMP",
  "commitCount": $COMMIT_COUNT,
  "lastCommit": "$LAST_COMMIT",
  "conversationHistory": [],
  "notes": []
}
EOF
    echo "📝 Created initial metadata file: $METADATA_FILE"
    echo ""
fi

if [ -f "$METADATA_FILE" ]; then
    # Run auto-compaction first (before Claude loads the context)
    if [ -f ".claude/scripts/compact-history.sh" ]; then
        bash ".claude/scripts/compact-history.sh" "$METADATA_FILE" 2>/dev/null
    fi

    echo "📎 Branch context loaded from $METADATA_FILE"
    echo ""
    echo "🤖 Branch Context:"
    cat "$METADATA_FILE"
    echo ""
    echo ""
    echo "Please present the PR context (branch, type, description, issue) and summarize the conversation history and notes if available."
    echo ""
    echo "Available: /pr-describe, /pr-status"
fi
