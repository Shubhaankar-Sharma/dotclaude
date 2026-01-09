package main

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type refreshMsg struct {
	diffContent  string
	lineMappings []LineMapping
}

type errMsg struct {
	err error
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		headerHeight := 2
		footerHeight := 2
		verticalMargins := headerHeight + footerHeight

		m.viewport.Width = msg.Width
		m.viewport.Height = msg.Height - verticalMargins

		if !m.ready {
			m.viewport.SetContent(m.diffContent)
			m.ready = true
		}

		return m, nil

	case tea.KeyMsg:
		// Route based on current mode
		switch m.mode {
		case viewMode:
			return m.handleViewMode(msg)
		case visualMode:
			return m.handleVisualMode(msg)
		case commentMode:
			return m.handleCommentMode(msg)
		case summaryMode:
			return m.handleSummaryMode(msg)
		}

	case refreshMsg:
		m.diffContent = msg.diffContent
		m.lineMappings = msg.lineMappings
		m.viewport.SetContent(msg.diffContent)
		m.currentLine = 0
		m.viewport.GotoTop()
		return m, nil

	case errMsg:
		m.err = msg.err
		return m, nil
	}

	// Always update viewport for scrolling
	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
}

func (m model) handleViewMode(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c":
		m.exitMode = exitQuit
		return m, tea.Quit

	case "q":
		// Normal quit
		m.exitMode = exitQuit
		return m, tea.Quit

	case "v":
		// Enter visual mode for line selection
		m.mode = visualMode
		m.visualStart = m.currentLine
		m.visualEnd = m.currentLine
		m.inVisual = true
		return m, nil

	case "Q":
		// Question mode on current line (capital Q)
		return m, tea.Sequence(
			m.saveQuestionContext(m.currentLine, m.currentLine),
			tea.Quit,
		)

	case "c":
		// Enter comment mode
		m.mode = commentMode
		m.textInput.Focus()
		return m, textinput.Blink

	case "s":
		// Show summary
		m.mode = summaryMode
		return m, nil

	case "a":
		// Apply pending comments (save and exit normally)
		if len(m.comments) > 0 {
			return m, tea.Sequence(
				m.saveComments(),
				tea.Quit,
			)
		}
		return m, nil

	case "r":
		// Refresh diff
		return m, m.refreshDiff()

	case "j", "down":
		if m.currentLine < len(m.lineMappings)-1 {
			m.currentLine++
			m.viewport.LineDown(1)
		}
		return m, nil

	case "k", "up":
		if m.currentLine > 0 {
			m.currentLine--
			m.viewport.LineUp(1)
		}
		return m, nil

	case "g":
		// Go to top
		m.currentLine = 0
		m.viewport.GotoTop()
		return m, nil

	case "G":
		// Go to bottom
		m.currentLine = len(m.lineMappings) - 1
		m.viewport.GotoBottom()
		return m, nil

	case "pgdown":
		m.viewport.ViewDown()
		m.currentLine = min(m.currentLine+m.viewport.Height/2, len(m.lineMappings)-1)
		return m, nil

	case "pgup":
		m.viewport.ViewUp()
		m.currentLine = max(m.currentLine-m.viewport.Height/2, 0)
		return m, nil
	}

	return m, nil
}

func (m model) handleVisualMode(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		// Cancel visual selection
		m.mode = viewMode
		m.inVisual = false
		return m, nil

	case "q":
		// Ask question about selected range
		start := min(m.visualStart, m.visualEnd)
		end := max(m.visualStart, m.visualEnd)
		return m, tea.Sequence(
			m.saveQuestionContext(start, end),
			tea.Quit,
		)

	case "c":
		// Add comment on selected range
		m.mode = commentMode
		m.textInput.Focus()
		return m, textinput.Blink

	case "j", "down":
		// Extend selection downward
		if m.currentLine < len(m.lineMappings)-1 {
			m.currentLine++
			m.visualEnd = m.currentLine
			m.viewport.LineDown(1)
		}
		return m, nil

	case "k", "up":
		// Extend selection upward
		if m.currentLine > 0 {
			m.currentLine--
			m.visualEnd = m.currentLine
			m.viewport.LineUp(1)
		}
		return m, nil

	case "ctrl+c":
		m.exitMode = exitQuit
		return m, tea.Quit
	}

	return m, nil
}

func (m model) handleCommentMode(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg.String() {
	case "esc":
		// Cancel comment
		m.mode = viewMode
		m.textInput.Blur()
		m.textInput.SetValue("")
		return m, nil

	case "enter":
		// Save comment
		if m.textInput.Value() != "" {
			// Determine range
			var startLine, endLine int
			var file string
			var line int

			if m.inVisual {
				startLine = min(m.visualStart, m.visualEnd)
				endLine = max(m.visualStart, m.visualEnd)
				file = m.lineMappings[startLine].File
				line = m.lineMappings[startLine].NewLine
			} else {
				startLine = m.currentLine
				endLine = m.currentLine
				mapping := m.lineMappings[m.currentLine]
				file = mapping.File
				line = mapping.NewLine
			}

			// Extract context
			context := ""
			for i := startLine; i <= endLine && i < len(m.lineMappings); i++ {
				context += m.lineMappings[i].Content + "\n"
			}

			comment := Comment{
				ID:      fmt.Sprintf("comment-%d", len(m.comments)+1),
				File:    file,
				Line:    line,
				Type:    "change",
				Content: m.textInput.Value(),
				Context: context,
				AddedAt: time.Now().UTC().Format(time.RFC3339),
			}

			m.comments = append(m.comments, comment)

			// Exit visual mode if active
			m.inVisual = false
		}

		m.mode = viewMode
		m.textInput.Blur()
		m.textInput.SetValue("")
		return m, nil

	case "ctrl+c":
		m.exitMode = exitQuit
		return m, tea.Quit
	}

	// Pass message to text input
	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m model) handleSummaryMode(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc", "q":
		// Return to view mode
		m.mode = viewMode
		return m, nil

	case "ctrl+c":
		m.exitMode = exitQuit
		return m, tea.Quit
	}

	return m, nil
}

func (m model) refreshDiff() tea.Cmd {
	return func() tea.Msg {
		diffContent, lineMappings, err := getDiffWithMappings(m.baseCommit, m.headCommit)
		if err != nil {
			return errMsg{err}
		}
		return refreshMsg{
			diffContent:  diffContent,
			lineMappings: lineMappings,
		}
	}
}

func (m model) saveQuestionContext(startLine, endLine int) tea.Cmd {
	return func() tea.Msg {
		m.exitMode = exitQuestion
		if err := writeQuestionContext(m.sessionID, m.lineMappings, startLine, endLine); err != nil {
			return errMsg{err}
		}
		return tea.Quit()
	}
}

func (m model) saveComments() tea.Cmd {
	return func() tea.Msg {
		if err := saveCommentsToMetadata(m.metadataPath, m.sessionID, m.comments); err != nil {
			return errMsg{err}
		}
		return tea.Quit()
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
