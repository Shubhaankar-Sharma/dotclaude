You are starting a PR session for the current git branch.

Follow these steps:

1. Get the current git branch name using `git branch --show-current`
2. Check if a metadata file already exists at `.claude/branches/<branch>.json`
   - If it exists, load it and show the user the existing context
   - Ask if they want to update it or keep it as is
3. If no metadata exists or user wants to update, ask the user:
   - **PR type**: feature, bugfix, refactor, docs, or test
   - **Issue link**: GitHub issue URL or number (optional)
   - **Brief description**: One-line description of what this PR does
4. Check if there's an associated GitHub issue by searching with `gh issue list --search "<branch-name>"`
5. Create the metadata file at `.claude/branches/<branch>.json` with this structure:
   ```json
   {
     "branch": "<branch-name>",
     "type": "<type>",
     "issue": "<issue-number or empty>",
     "issueUrl": "<issue-url or empty>",
     "description": "<description>",
     "started": "<ISO 8601 timestamp>",
     "lastWorked": "<ISO 8601 timestamp>",
     "commitCount": <number>,
     "lastCommit": "<commit-sha>",
     "notes": []
   }
   ```
6. Get the commit count with: `git rev-list --count HEAD ^origin/main` (or appropriate base branch)
7. Get the last commit SHA with: `git rev-parse --short HEAD`
8. Output a summary:
   ```
   📎 PR Session Started

   Branch: <branch>
   Type: <type>
   Issue: <issue-link>
   Description: <description>

   💡 Tip: Rename this conversation to:
   [PR] <branch>: <description>

   Context saved to .claude/branches/<branch>.json
   Use /pr-describe when ready to generate the PR description.
   ```

From now on, all work in this conversation is scoped to this branch and PR context.
