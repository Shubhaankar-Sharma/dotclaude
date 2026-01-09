package main

import (
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

type mode int

const (
	viewMode mode = iota
	visualMode
	commentMode
	summaryMode
)

type exitMode int

const (
	exitNormal exitMode = iota
	exitError
	exitQuit
	exitQuestion
)

// LineMapping maps diff lines to file locations
type LineMapping struct {
	DiffLine int
	File     string
	OldLine  int
	NewLine  int
	Content  string
	Op       string // "add", "delete", "context"
}

// Comment represents a review comment
type Comment struct {
	ID      string `json:"id"`
	File    string `json:"file"`
	Line    int    `json:"line"`
	Type    string `json:"type"`
	Content string `json:"content"`
	Context string `json:"context"`
	AddedAt string `json:"addedAt"`
}

// Conversation represents a Q&A session
type Conversation struct {
	ID        string `json:"id"`
	File      string `json:"file"`
	Line      int    `json:"line"`
	Question  string `json:"question"`
	Answer    string `json:"answer"`
	Timestamp string `json:"timestamp"`
}

// model is the Bubble Tea model
type model struct {
	// Bubble components
	viewport  viewport.Model
	textInput textinput.Model

	// Application state
	mode     mode
	exitMode exitMode
	ready    bool

	// Diff data
	diffContent  string
	lineMappings []LineMapping
	currentLine  int

	// Visual selection
	visualStart int
	visualEnd   int
	inVisual    bool

	// Review data
	comments      []Comment
	conversations []Conversation
	sessionID     string
	metadataPath  string

	// Git data
	baseCommit string
	headCommit string

	// Window dimensions
	width  int
	height int

	// Error state
	err error
}

func initialModel(sessionID, baseCommit, headCommit, metadataPath, diffContent string, lineMappings []LineMapping) model {
	ti := textinput.New()
	ti.Placeholder = "Enter your comment..."
	ti.CharLimit = 500
	ti.Width = 80

	vp := viewport.New(80, 20)
	vp.SetContent(diffContent)

	return model{
		viewport:      vp,
		textInput:     ti,
		mode:          viewMode,
		exitMode:      exitNormal,
		ready:         false,
		diffContent:   diffContent,
		lineMappings:  lineMappings,
		currentLine:   0,
		visualStart:   0,
		visualEnd:     0,
		inVisual:      false,
		comments:      []Comment{},
		conversations: []Conversation{},
		sessionID:     sessionID,
		metadataPath:  metadataPath,
		baseCommit:    baseCommit,
		headCommit:    headCommit,
		width:         80,
		height:        24,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}
