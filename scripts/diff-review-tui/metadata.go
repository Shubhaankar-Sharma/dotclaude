package main

import (
	"encoding/json"
	"fmt"
	"os"
)

// QuestionContext represents the context for a question asked during review
type QuestionContext struct {
	File      string   `json:"file"`
	StartLine int      `json:"startLine"`
	EndLine   int      `json:"endLine"`
	DiffLines []string `json:"diffLines"`
	SessionID string   `json:"sessionId"`
}

// ReviewSession represents a complete review session in metadata
type ReviewSession struct {
	SessionID     string         `json:"sessionId"`
	Timestamp     string         `json:"timestamp"`
	BaseCommit    string         `json:"baseCommit"`
	HeadCommit    string         `json:"headCommit"`
	Status        string         `json:"status"`
	Comments      []Comment      `json:"comments"`
	Conversations []Conversation `json:"conversations"`
}

// Metadata represents the branch metadata structure
type Metadata struct {
	Branch         string          `json:"branch"`
	Type           string          `json:"type"`
	Issue          string          `json:"issue,omitempty"`
	IssueURL       string          `json:"issueUrl,omitempty"`
	Description    string          `json:"description,omitempty"`
	Started        string          `json:"started,omitempty"`
	LastWorked     string          `json:"lastWorked,omitempty"`
	CommitCount    int             `json:"commitCount,omitempty"`
	LastCommit     string          `json:"lastCommit,omitempty"`
	Notes          []interface{}   `json:"notes,omitempty"`
	ReviewSessions []ReviewSession `json:"reviewSessions,omitempty"`
}

// writeQuestionContext writes the question context to a temp file for Claude to read
func writeQuestionContext(sessionID string, lineMappings []LineMapping, startLine, endLine int) error {
	if startLine >= len(lineMappings) || endLine >= len(lineMappings) {
		return fmt.Errorf("line range out of bounds")
	}

	startMapping := lineMappings[startLine]
	endMapping := lineMappings[endLine]

	// Extract the selected diff lines
	var diffLines []string
	for i := startLine; i <= endLine && i < len(lineMappings); i++ {
		diffLines = append(diffLines, lineMappings[i].Content)
	}

	ctx := QuestionContext{
		File:      startMapping.File,
		StartLine: startMapping.NewLine,
		EndLine:   endMapping.NewLine,
		DiffLines: diffLines,
		SessionID: sessionID,
	}

	// Write to temp file
	contextPath := fmt.Sprintf("/tmp/claude-review-context-%s.json", sessionID)
	data, err := json.MarshalIndent(ctx, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal context: %w", err)
	}

	err = os.WriteFile(contextPath, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write context file: %w", err)
	}

	return nil
}

// saveCommentsToMetadata saves comments to the branch metadata JSON
func saveCommentsToMetadata(metadataPath, sessionID string, comments []Comment) error {
	// Read existing metadata
	data, err := os.ReadFile(metadataPath)
	if err != nil {
		return fmt.Errorf("failed to read metadata: %w", err)
	}

	var metadata Metadata
	if err := json.Unmarshal(data, &metadata); err != nil {
		return fmt.Errorf("failed to parse metadata: %w", err)
	}

	// Find and update the current session
	found := false
	for i, session := range metadata.ReviewSessions {
		if session.SessionID == sessionID {
			metadata.ReviewSessions[i].Status = "completed"
			metadata.ReviewSessions[i].Comments = comments
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("session not found in metadata: %s", sessionID)
	}

	// Write back to file
	output, err := json.MarshalIndent(metadata, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	err = os.WriteFile(metadataPath, output, 0644)
	if err != nil {
		return fmt.Errorf("failed to write metadata: %w", err)
	}

	return nil
}

// loadCommentsFromMetadata loads existing comments for a session
func loadCommentsFromMetadata(metadataPath, sessionID string) ([]Comment, error) {
	data, err := os.ReadFile(metadataPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read metadata: %w", err)
	}

	var metadata Metadata
	if err := json.Unmarshal(data, &metadata); err != nil {
		return nil, fmt.Errorf("failed to parse metadata: %w", err)
	}

	// Find the session
	for _, session := range metadata.ReviewSessions {
		if session.SessionID == sessionID {
			return session.Comments, nil
		}
	}

	return []Comment{}, nil
}
