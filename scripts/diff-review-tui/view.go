package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("39"))

	fileStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240"))

	commentCountStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("10"))

	visualStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("62")).
			Foreground(lipgloss.Color("230")).
			Bold(true)

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241"))

	selectedLineStyle = lipgloss.NewStyle().
				Background(lipgloss.Color("237")).
				Foreground(lipgloss.Color("15"))

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("9")).
			Bold(true)
)

func (m model) View() string {
	if !m.ready {
		return "Loading diff..."
	}

	if m.err != nil {
		return errorStyle.Render(fmt.Sprintf("Error: %v", m.err))
	}

	var header string
	var footer string
	var content string

	switch m.mode {
	case viewMode:
		header = m.viewModeHeader()
		footer = m.viewModeFooter()
		content = m.viewport.View()

	case visualMode:
		header = m.visualModeHeader()
		footer = m.visualModeFooter()
		content = m.renderWithSelection()

	case commentMode:
		header = m.commentModeHeader()
		footer = m.commentModeFooter()
		content = m.viewport.View()

	case summaryMode:
		return m.summaryView()
	}

	return lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		content,
		footer,
	)
}

func (m model) viewModeHeader() string {
	var location string
	if m.currentLine < len(m.lineMappings) {
		mapping := m.lineMappings[m.currentLine]
		if mapping.NewLine > 0 {
			location = fmt.Sprintf("%s:%d", mapping.File, mapping.NewLine)
		} else if mapping.OldLine > 0 {
			location = fmt.Sprintf("%s:%d (deleted)", mapping.File, mapping.OldLine)
		} else {
			location = mapping.File
		}
	}

	return fmt.Sprintf("%s | %s | Comments: %s\n",
		titleStyle.Render(fmt.Sprintf("Review: %s", m.sessionID[:12])),
		fileStyle.Render(location),
		commentCountStyle.Render(fmt.Sprintf("%d", len(m.comments))),
	)
}

func (m model) viewModeFooter() string {
	help := "[v] Visual [q] Question [Q] Quick Question [c] Comment [s] Summary [a] Apply [r] Refresh [Ctrl+C] Quit"
	return helpStyle.Render(help)
}

func (m model) visualModeHeader() string {
	start := min(m.visualStart, m.visualEnd)
	end := max(m.visualStart, m.visualEnd)
	lineCount := end - start + 1

	var location string
	if start < len(m.lineMappings) && end < len(m.lineMappings) {
		startMapping := m.lineMappings[start]
		endMapping := m.lineMappings[end]
		location = fmt.Sprintf("%s:%d-%d",
			startMapping.File,
			startMapping.NewLine,
			endMapping.NewLine,
		)
	}

	return fmt.Sprintf("%s | %s | %s\n",
		titleStyle.Render(fmt.Sprintf("Review: %s", m.sessionID[:12])),
		fileStyle.Render(location),
		visualStyle.Render(fmt.Sprintf(" VISUAL (%d lines) ", lineCount)),
	)
}

func (m model) visualModeFooter() string {
	help := "[q] Question [c] Comment [Esc] Cancel [j/k] Extend Selection"
	return helpStyle.Render(help)
}

func (m model) commentModeHeader() string {
	var location string
	if m.inVisual {
		start := min(m.visualStart, m.visualEnd)
		end := max(m.visualStart, m.visualEnd)
		lineCount := end - start + 1
		if start < len(m.lineMappings) {
			location = fmt.Sprintf("%s (selecting %d lines)", m.lineMappings[start].File, lineCount)
		}
	} else if m.currentLine < len(m.lineMappings) {
		mapping := m.lineMappings[m.currentLine]
		location = fmt.Sprintf("%s:%d", mapping.File, mapping.NewLine)
	}

	return fmt.Sprintf("%s | %s\n%s\n",
		titleStyle.Render("Add Comment"),
		fileStyle.Render(location),
		m.textInput.View(),
	)
}

func (m model) commentModeFooter() string {
	help := "[Enter] Save [Esc] Cancel"
	return helpStyle.Render(help)
}

func (m model) renderWithSelection() string {
	if !m.inVisual || len(m.lineMappings) == 0 {
		return m.viewport.View()
	}

	lines := strings.Split(m.diffContent, "\n")
	start := min(m.visualStart, m.visualEnd)
	end := max(m.visualStart, m.visualEnd)

	var result []string
	for i, line := range lines {
		if i >= start && i <= end {
			result = append(result, selectedLineStyle.Render(line))
		} else {
			result = append(result, line)
		}
	}

	content := strings.Join(result, "\n")

	// Create a new viewport with the highlighted content
	vp := m.viewport
	vp.SetContent(content)

	return vp.View()
}

func (m model) summaryView() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("Review Session Summary"))
	b.WriteString("\n\n")

	b.WriteString(fmt.Sprintf("Session ID: %s\n", m.sessionID))
	b.WriteString(fmt.Sprintf("Base Commit: %s\n", m.baseCommit[:8]))
	b.WriteString(fmt.Sprintf("Head Commit: %s\n", m.headCommit[:8]))
	b.WriteString("\n")

	b.WriteString(fmt.Sprintf("Comments: %d\n", len(m.comments)))
	b.WriteString(fmt.Sprintf("Conversations: %d\n", len(m.conversations)))
	b.WriteString("\n")

	if len(m.comments) > 0 {
		b.WriteString(lipgloss.NewStyle().Bold(true).Render("Comments:"))
		b.WriteString("\n\n")

		for i, comment := range m.comments {
			b.WriteString(fmt.Sprintf("%d. %s:%d\n", i+1, comment.File, comment.Line))
			b.WriteString(fmt.Sprintf("   Type: %s\n", comment.Type))
			b.WriteString(fmt.Sprintf("   %s\n", comment.Content))
			b.WriteString("\n")
		}
	} else {
		b.WriteString("No comments added yet.\n\n")
	}

	if len(m.conversations) > 0 {
		b.WriteString(lipgloss.NewStyle().Bold(true).Render("Conversations:"))
		b.WriteString("\n\n")

		for i, conv := range m.conversations {
			b.WriteString(fmt.Sprintf("%d. %s:%d\n", i+1, conv.File, conv.Line))
			b.WriteString(fmt.Sprintf("   Q: %s\n", conv.Question))
			b.WriteString(fmt.Sprintf("   A: %s\n", conv.Answer))
			b.WriteString("\n")
		}
	}

	b.WriteString("\n")
	b.WriteString(helpStyle.Render("[Esc] Back to Review"))

	return b.String()
}
