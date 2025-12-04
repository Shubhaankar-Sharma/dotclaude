Save the current session summary and update branch metadata before ending your Claude session.

Follow these steps:

1. Get the current branch: `git branch --show-current`
2. Sanitize branch name: replace `/` with `-` for filename
3. Check if metadata exists: `.claude/branches/<sanitized-branch>.json`
4. If metadata exists:
   - Read the session-summary template: `.claude/templates/session-summary.md`
   - Create a comprehensive session summary including:
     - What was built/changed (code, files, features)
     - Key decisions and trade-offs made
     - User instructions, preferences, and feedback given
     - Any blockers or issues encountered
     - Next steps and what to work on next
   - Update the metadata file:
     - Append new entry to `conversationHistory` array with `{timestamp, summary}`
     - Update `lastWorked` to current ISO timestamp
     - Update `commitCount` from `git rev-list --count HEAD ^origin/main`
     - Update `lastCommit` from `git rev-parse --short HEAD`
   - Confirm save: "💾 Session summary saved to .claude/branches/<branch>.json"
5. If no metadata exists:
   ```
   No PR context found for branch '<branch>'.

   Nothing to save. Use /pr-start if you want to track this branch.
   ```
