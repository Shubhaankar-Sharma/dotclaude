Generate a PR description for the current branch using git history and metadata.

Follow these steps:

1. Get the current branch: `git branch --show-current`
2. Read the metadata file: `.claude/branches/<branch>.json`
   - If it doesn't exist, suggest running `/pr-start` first
3. Read the PR description template: `.claude/templates/pr-description.md`
4. Read the style guide: `.claude/prompts/pr-style-guide.md`
5. Analyze the branch commits:
   - Get all commits: `git log origin/main..HEAD --oneline`
   - Get the full diff: `git diff origin/main...HEAD`
   - If there's a specific commit SHA in metadata, you can also analyze: `git show <commit-sha>`
6. If there's an issue number, fetch issue details: `gh issue view <issue-number>`
7. Generate the PR description following these rules:
   - Use the exact template format from `pr-description.md`
   - Follow the writing style from `pr-style-guide.md`
   - Be expressive, dense, and minimal
   - No hollow technical words
   - Include concrete examples
   - Reference code locations as `file.go:123` or `file.go:123-456`
   - Background section: Explain what existed, why this change is needed, what problem it solves
   - Effect section: Dense paragraphs explaining the impact and how it works
   - Implementation details: Bullet points with file:line references and clear explanations

8. Output the complete PR description as markdown for the user to copy-paste

DO NOT create a PR or push anything. Just output the description.
