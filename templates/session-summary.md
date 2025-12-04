# Session Summary Template

Summarize this coding session comprehensively. Capture everything needed to resume work seamlessly:

1. **What was built/changed** - Concrete code changes, features added, bugs fixed
2. **Key decisions made** - Technical choices, architectural decisions, trade-offs
3. **User instructions/preferences** - Important directives, constraints, coding style preferences, things user wants/doesn't want
4. **Outstanding questions/blockers** - Unresolved issues, things to revisit
5. **Next steps** - What should happen in the next session

## Format

Write a dense, technical summary. Include all relevant details. Use past tense. Example:

"Implemented conversation persistence system for Claude hooks. Added conversationHistory and compactedHistory fields to branch metadata schema. Created session-summary.md template and enhanced session-end.sh to collect summaries using jq. Built auto-compaction script that keeps last 5 sessions detailed and compacts older entries. User instruction: Summaries must include all important user instructions/preferences - these are critical for context. User wants auto-compaction to happen automatically without manual intervention. Decision: Keep compaction threshold at 5 sessions, run compaction in session-end hook. Outstanding: Need to test if jq can handle multiline summary formatting. Next: Enhance session-end.sh to prompt for summary and append to metadata."

## Critical: Capture User Context

The most important thing is preserving **what the user told you**:
- Explicit requirements ("I want X to do Y")
- Constraints ("Don't do Z", "Keep it minimal")
- Preferences ("I prefer A over B")
- Domain knowledge they shared ("This is how our system works")
- Future intentions ("Eventually we'll need Q")

Without this, the next session loses critical context.
