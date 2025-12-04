#!/bin/bash

# Compact History Script: Keeps conversation history manageable
# Moves older session summaries into compactedHistory to prevent unbounded growth

METADATA_FILE="$1"

if [ -z "$METADATA_FILE" ] || [ ! -f "$METADATA_FILE" ]; then
    exit 0
fi

# Check if jq is installed
if ! command -v jq &> /dev/null; then
    exit 0
fi

# Configuration: Keep last N sessions in conversationHistory
KEEP_RECENT=10

# Read the current conversation history length
HISTORY_LENGTH=$(jq '.conversationHistory | length' "$METADATA_FILE" 2>/dev/null || echo "0")

# If history is small enough, no compaction needed
if [ "$HISTORY_LENGTH" -le "$KEEP_RECENT" ]; then
    exit 0
fi

# Calculate how many entries to compact
COMPACT_COUNT=$((HISTORY_LENGTH - KEEP_RECENT))

# Create a timestamp for this compaction
COMPACT_TIMESTAMP=$(date -u +"%Y-%m-%dT%H:%M:%S.000Z")

# Extract the entries to compact (oldest entries)
ENTRIES_TO_COMPACT=$(jq --argjson count "$COMPACT_COUNT" '.conversationHistory[:$count]' "$METADATA_FILE")

# Create a summary of the compacted entries
# Get the timestamp range
FIRST_TIMESTAMP=$(echo "$ENTRIES_TO_COMPACT" | jq -r '.[0].timestamp')
LAST_TIMESTAMP=$(echo "$ENTRIES_TO_COMPACT" | jq -r '.[-1].timestamp')

# Concatenate all summaries
COMBINED_SUMMARY=$(echo "$ENTRIES_TO_COMPACT" | jq -r '.[].summary' | sed 's/^/- /')

# Create the compacted entry
COMPACTED_ENTRY=$(jq -n \
    --arg from "$FIRST_TIMESTAMP" \
    --arg to "$LAST_TIMESTAMP" \
    --arg summary "$COMBINED_SUMMARY" \
    --argjson count "$COMPACT_COUNT" \
    '{
        compactedAt: "'$COMPACT_TIMESTAMP'",
        periodStart: $from,
        periodEnd: $to,
        sessionCount: $count,
        summary: $summary
    }')

# Update the metadata file
jq --argjson entry "$COMPACTED_ENTRY" \
   --argjson keep "$KEEP_RECENT" \
   '
   .compactedHistory = (.compactedHistory // []) + [$entry] |
   .conversationHistory = .conversationHistory[-$keep:]
   ' "$METADATA_FILE" > "${METADATA_FILE}.tmp"

# Replace the original file if update succeeded
if [ $? -eq 0 ]; then
    mv "${METADATA_FILE}.tmp" "$METADATA_FILE"
else
    rm -f "${METADATA_FILE}.tmp"
fi