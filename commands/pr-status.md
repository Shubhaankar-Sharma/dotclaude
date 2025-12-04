Show the current PR context and branch status.

Follow these steps:

1. Get the current branch: `git branch --show-current`
2. Check if metadata exists: `.claude/branches/<branch>.json`
3. If metadata exists, display:
   ```
   📎 Current PR Context

   Branch: <branch>
   Type: <type>
   Issue: <issue-link>
   Description: <description>
   Started: <date>
   Last worked: <date>
   Commits: <count>
   Last commit: <sha>
   ```
4. Show recent commits: `git log --oneline -5`
5. Show git status: `git status --short`
6. If no metadata exists, output:
   ```
   No PR context found for branch '<branch>'.

   Run /pr-start to initialize PR tracking.
   ```
