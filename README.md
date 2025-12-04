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
chmod +x .claude/hooks/session-start.sh
chmod +x .claude/scripts/compact-history.sh
```

The `settings.json` file is already configured with the session-start hook and necessary permissions.

## Usage

### Commands

- `/pr-start` - Initialize PR tracking for current branch
- `/pr-describe` - Generate PR description from git history and metadata
- `/pr-status` - Show current PR context and branch status
- `/session-end` - Save session summary before ending work

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

- `.claude/settings.json` - Hook configuration
- `.claude/hooks/session-start.sh` - Auto-context loader
- `.claude/scripts/compact-history.sh` - History compaction
- `.claude/branches/<branch>.json` - Branch metadata (git-ignored)
- `.claude/commands/*.md` - Slash commands
- `.claude/prompts/pr-style-guide.md` - PR writing style
- `.claude/templates/pr-description.md` - PR description template
- `.claude/templates/session-summary.md` - Session summary template

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

- Claude Code CLI
- Git (for branch detection and commit history)
- GitHub CLI (`gh`) - Optional, for issue linking in PR descriptions
- jq - Required for history compaction

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
chmod +x .claude/hooks/session-start.sh
chmod +x .claude/scripts/compact-history.sh
```

## License

MIT
