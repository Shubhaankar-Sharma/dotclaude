# Claude PR Workflow

A Claude Code configuration system that links git branches to PRs and generates consistent PR descriptions without the slop.

## The Problem

I have been managing a lot of Claude Code sessions across branches linked to different PRs and also in separate projects. I have been using Claude to generate PR descriptions but it would generate them in a weird format with a lot of emoji slop.

One option was to always resume conversations, but this eventually got difficult. For other work I would always have to prompt it how to generate the description.

## The Solution

After reading [this blog from Anthropic](https://www.anthropic.com/engineering/claude-code-best-practices), I created a simple set of commands and templates to handle this better.

## How It Works

### Starting a PR

When I start working on a PR, I call the `/pr-start` command.

Claude asks a few questions (PR type, issue link, description) and can fetch relevant info from the GitHub issue link you paste.

It saves all this information in a JSON file: `.claude/branches/example-branch.json`

This file contains all the information I provided. I wrote a hook that loads all this information onto Claude's context on session start.

So whenever I start a new session, Claude already knows what I have been working on and where we left off from previous conversations.

### Ending a Session

When I end conversations, I call `/session-end`. This triggers Claude to add information to the branch metadata that it might need next time.

### Generating PR Descriptions

The `/pr-describe` command uses a predefined template of how I want my descriptions to look and uses the branch metadata to generate it.

Example: https://github.com/planetscale/neki/pull/966

## Installation

### Prerequisites

Install jq (required for history compaction):

```bash
brew install jq  # macOS
sudo apt-get install jq  # Linux
```

### Setup

Clone this repo:
```bash
git clone https://github.com/shubhaankar-sharma/dotclaude.git
cd dotclaude
```

Copy the contents to your project's `.claude` directory:
```bash
# From your project root
cp -r /path/to/dotclaude/* .claude/
```

Create the branches directory (git-ignored):
```bash
mkdir -p .claude/branches
```

Make scripts executable:
```bash
chmod +x hooks/session-start.sh
chmod +x scripts/compact-history.sh
chmod +x scripts/build-tui.sh
```

Build the diff-review TUI (requires Go 1.21+):
```bash
cd scripts/diff-review-tui
go mod download
go build -o ../diff-review -ldflags="-s -w" .
chmod +x ../diff-review
cd ../..
```

The `settings.json` file is already configured with the session-start hook and necessary permissions.

## Usage

### Commands

- `/pr-start` - Initialize PR tracking for current branch
- `/pr-describe` - Generate PR description from git history and metadata
- `/pr-status` - Show current PR context and branch status
- `/session-end` - Save session summary before ending work
- `/diff-review` - Launch interactive visual diff review TUI
- `/apply-review` - Apply pending review comments from completed sessions

### Starting Work on a Branch

```bash
git checkout -b your-branch-name
# Open Claude Code and run:
/pr-start
```

Answer the prompts (PR type, issue, description).

### Working Normally

Open any Claude conversation. The session-start hook automatically:
- Compacts old conversation history (keeps last 10 sessions)
- Prompts Claude to load and present your branch context
- Shows what you're working on and next steps

Claude will immediately present your PR context when starting a new session.

### Before Ending Work

```bash
/session-end
```

Saves a comprehensive session summary that gets added to the branch metadata.

### Generating PR Description

```bash
/pr-describe
```

Copy output, paste into GitHub.

### Visual Diff Review

```bash
/diff-review
```

Opens an interactive terminal UI for reviewing your changes with Claude. Features:
- Navigate through diffs with syntax highlighting
- Select lines and ask Claude questions about the code
- Add review comments for batch processing
- Request edits during conversation
- Apply all feedback at once

See the [Visual Diff Review](#visual-diff-review-system) section for details.

## Visual Diff Review System

An interactive terminal-based code review system that combines syntax-highlighted diff viewing with Claude's conversational capabilities.

### Features

- **Interactive TUI**: Navigate diffs with vim-style keybindings (j/k) in a rich terminal interface
- **Syntax Highlighting**: Powered by Chroma - supports 150+ languages with beautiful color schemes
- **Visual Selection**: Select multiple lines (like vim's visual mode) to discuss code blocks
- **Conversational Review**: Press 'q' to ask Claude questions about specific code sections
- **Context-Aware**: Claude sees ±20 lines around your selection plus full codebase access
- **Request Edits**: Ask Claude to make changes during review, then refresh the diff
- **Batch Comments**: Add review comments and apply them all at once
- **Persistent Sessions**: All review data saved in branch metadata

### Setup

The TUI is built in Go for performance and portability. Build it once:

```bash
cd scripts/diff-review-tui
go mod download
go build -o ../diff-review -ldflags="-s -w" .
chmod +x ../diff-review
```

Or use the one-liner:
```bash
cd scripts/diff-review-tui && go mod download && go build -o ../diff-review -ldflags="-s -w" . && chmod +x ../diff-review
```

### Basic Workflow

1. **Start review**:
   ```bash
   /diff-review
   ```

2. **Navigate**: Use `j`/`k` or arrow keys to move through the diff

3. **Ask questions**:
   - Navigate to interesting line
   - Press `Q` (Shift+q) to ask about current line
   - Or press `v` to enter visual mode, select lines with `j`/`k`, then press `q`
   - Claude shows context and answers your question
   - Request edits if needed: "Change this to use X instead"
   - Press Enter to return to TUI

4. **Add comments for later**:
   - Navigate to a line
   - Press `c` to add a comment
   - Type your feedback and press Enter
   - Continue reviewing

5. **Refresh after edits**:
   - Press `r` to reload the diff after Claude makes changes
   - See your requested changes immediately

6. **Apply pending comments**:
   - Press `a` in TUI to save and exit
   - Or run `/apply-review` from command line
   - Claude batch-processes all comments

### Keybindings

| Key | Action |
|-----|--------|
| `j` / `↓` | Move down one line |
| `k` / `↑` | Move up one line |
| `g` | Go to top |
| `G` | Go to bottom |
| `PgDn` | Page down |
| `PgUp` | Page up |
| `v` | Enter visual mode (line selection) |
| `q` | Ask question about selected lines |
| `Q` | Quick question about current line |
| `c` | Add comment for batch processing |
| `s` | Show summary of comments/conversations |
| `a` | Apply pending comments (save & exit) |
| `r` | Refresh diff (after Claude edits) |
| `Esc` | Cancel current mode |
| `Ctrl+C` | Quit review |

### Example Session

```bash
# Make changes on feature branch
git checkout -b feature/new-auth
# ... make code changes ...

# Start interactive review
/diff-review
```

**In TUI:**
1. Navigate to suspicious timeout value on line 45
2. Press `Q` → Ask: "Why was this originally 10s?"
3. Claude explains and shows related code
4. Say: "Change it to 15s and add a comment explaining why"
5. Claude makes the edit
6. Press Enter to return to TUI
7. Press `r` to refresh and see the change
8. Continue reviewing other changes
9. Press `v` on line 67, `j j j` to select 4 lines
10. Press `q` → Ask: "Is this error handling correct?"
11. Claude spots a bug: "You're returning before checking the error"
12. Say: "Fix it"
13. Claude fixes the code
14. Return to TUI, press `r` to refresh
15. Review looks good, press `Ctrl+C` to quit

```bash
# Verify all changes
git diff

# Commit everything (your original changes + Claude's fixes)
git add .
git commit -m "Add authentication with review fixes"
```

### Conversational Review Patterns

**Pattern 1: Quick Question**
- Navigate → Press `Q` → Ask → Return

**Pattern 2: Discuss and Edit**
- Select lines with `v` + `j`/`k` → Press `q` → Discuss → Request edit → Claude edits → Return → Press `r`

**Pattern 3: Batch Comments**
- Navigate → Press `c` → Add comment → Continue → Press `a` to apply all

**Pattern 4: Explore Context**
- Select section → Press `q` → "Show me where this is called" → Claude finds callers → Discuss architecture

### Why It's Powerful

Unlike traditional review tools, Claude has **full codebase access** during questions:
- Find related code across the project
- Explain why code was written a certain way
- Show git blame history
- Make edits based on discussion
- Search for similar patterns
- Understand broader architectural context

### Review Data

All review sessions are stored in branch metadata (`.claude/branches/<branch>.json`):

```json
{
  "reviewSessions": [
    {
      "sessionId": "review-20240115-143000-abc123",
      "timestamp": "2024-01-15T14:30:00.000Z",
      "baseCommit": "abc123",
      "headCommit": "def456",
      "status": "completed",
      "comments": [
        {
          "file": "src/handler.go",
          "line": 45,
          "type": "change",
          "content": "Add nil check before dereferencing",
          "context": "surrounding code...",
          "addedAt": "2024-01-15T14:32:00.000Z"
        }
      ],
      "conversations": [
        {
          "file": "src/handler.go",
          "line": 67,
          "question": "Why buffered channel?",
          "answer": "Prevents goroutine blocking...",
          "timestamp": "2024-01-15T14:35:00.000Z"
        }
      ]
    }
  ]
}
```

### Requirements

- **Go 1.21+** (for building the TUI binary)
- **Git** (for diff generation)
- **jq** (for JSON manipulation in bash scripts)

After building once, the binary has no runtime dependencies.

## What Gets Tracked

Each branch gets a JSON file at `.claude/branches/<branch-name>.json`:

```json
{
  "branch": "feature/new-thing",
  "type": "feature",
  "issue": "123",
  "issueUrl": "https://github.com/org/repo/issues/123",
  "description": "Add support for new thing",
  "started": "2024-01-15T10:30:00.000Z",
  "lastWorked": "2024-01-15T15:45:00.000Z",
  "commitCount": 5,
  "lastCommit": "abc123",
  "conversationHistory": [
    {
      "timestamp": "2024-01-15T15:45:00.000Z",
      "summary": "Implemented core feature logic..."
    }
  ],
  "compactedHistory": [
    {
      "compactedAt": "2024-01-20T10:00:00.000Z",
      "periodStart": "2024-01-15T10:30:00.000Z",
      "periodEnd": "2024-01-18T16:20:00.000Z",
      "sessionCount": 8,
      "summary": "- Initial setup\n- Built API endpoints\n- Added tests"
    }
  ],
  "notes": []
}
```

## Files

### Core System
- `settings.json` - Hook configuration and permissions
- `hooks/session-start.sh` - Auto-context loader with review session checks
- `scripts/compact-history.sh` - History compaction
- `.claude/branches/<branch>.json` - Branch metadata (git-ignored)

### Commands
- `commands/pr-start.md` - Initialize PR tracking
- `commands/pr-describe.md` - Generate PR description
- `commands/pr-status.md` - Show PR context and review sessions
- `commands/session-end.md` - Save session summary
- `commands/diff-review.md` - Launch visual diff review TUI
- `commands/apply-review.md` - Apply pending review comments

### Templates & Prompts
- `prompts/pr-style-guide.md` - PR writing style guide
- `prompts/apply-review-comments.md` - Claude prompt for applying review feedback
- `templates/pr-description.md` - PR description template
- `templates/session-summary.md` - Session summary template
- `templates/review-summary.md` - Review session summary template

### Diff Review TUI
- `scripts/diff-review-tui/main.go` - TUI entry point
- `scripts/diff-review-tui/model.go` - Application state
- `scripts/diff-review-tui/update.go` - Message handling
- `scripts/diff-review-tui/view.go` - Rendering logic
- `scripts/diff-review-tui/diff.go` - Diff parsing and syntax highlighting
- `scripts/diff-review-tui/metadata.go` - JSON persistence
- `scripts/diff-review-tui/go.mod` - Go dependencies
- `scripts/diff-review` - Compiled binary (after build)
- `scripts/build-tui.sh` - Build script

## Customization

### PR Style Guide

Edit `.claude/prompts/pr-style-guide.md` to customize how PR descriptions are written. The default style avoids:
- Hollow technical words
- Excessive formatting
- Generic statements

And emphasizes:
- Expressive, dense, minimal writing
- Concrete examples
- Clear connections between implementation steps
- File:line references for code locations

### PR Description Template

Edit `.claude/templates/pr-description.md` to change the structure of generated PR descriptions.

### History Compaction

Edit `.claude/scripts/compact-history.sh` and change `KEEP_RECENT=10` to keep more or fewer recent sessions.

## Team Usage

In most organizations, most of the `.claude` files are committed to the repo so that other people can also use it.

If you want to do this, remove `.claude/branches/` from your `.gitignore` or keep it excluded if you prefer private branch metadata.

## Requirements

- **Claude Code CLI** - The AI-powered coding assistant
- **Git** - For branch detection and commit history
- **Go 1.21+** - Required to build the diff-review TUI
- **jq** - Required for JSON manipulation and history compaction
- **GitHub CLI (`gh`)** - Optional, for issue linking in PR descriptions

Note: After building the diff-review TUI binary once, it has no runtime dependencies on Go.

## Troubleshooting

### Hook not running

Verify hook is configured in `.claude/settings.json` and the script is executable.

### History not compacting

Check if `jq` is installed:
```bash
which jq
```

### Permissions errors

Ensure scripts are executable:
```bash
chmod +x hooks/session-start.sh
chmod +x scripts/compact-history.sh
chmod +x scripts/build-tui.sh
```

### Diff review TUI not launching

Build the binary:
```bash
cd scripts/diff-review-tui
go mod download
go build -o ../diff-review -ldflags="-s -w" .
chmod +x ../diff-review
```

Verify it works:
```bash
scripts/diff-review --help
```

### Go dependencies failing

Update Go to 1.21 or later:
```bash
go version  # Should show 1.21 or higher
```

On macOS:
```bash
brew install go
# or
brew upgrade go
```

## License

MIT
