package main

import (
	"flag"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	// Parse command line arguments
	sessionID := flag.String("session-id", "", "Review session ID")
	baseCommit := flag.String("base-commit", "", "Base commit for diff")
	headCommit := flag.String("head-commit", "", "Head commit for diff")
	metadataPath := flag.String("metadata-path", "", "Path to branch metadata JSON")
	flag.Parse()

	// Validate required arguments
	if *sessionID == "" || *baseCommit == "" || *headCommit == "" || *metadataPath == "" {
		fmt.Fprintln(os.Stderr, "Error: Missing required arguments")
		fmt.Fprintln(os.Stderr, "Usage: diff-review --session-id <id> --base-commit <sha> --head-commit <sha> --metadata-path <path>")
		os.Exit(1)
	}

	// Get diff and line mappings
	diffContent, lineMappings, err := getDiffWithMappings(*baseCommit, *headCommit)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting diff: %v\n", err)
		os.Exit(1)
	}

	if len(lineMappings) == 0 {
		fmt.Println("No changes to review.")
		os.Exit(0)
	}

	// Initialize model
	m := initialModel(*sessionID, *baseCommit, *headCommit, *metadataPath, diffContent, lineMappings)

	// Start Bubble Tea program
	p := tea.NewProgram(m, tea.WithAltScreen(), tea.WithMouseCellMotion())

	finalModel, err := p.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error running TUI: %v\n", err)
		os.Exit(1)
	}

	// Handle exit codes based on final state
	if m, ok := finalModel.(model); ok {
		switch m.exitMode {
		case exitNormal:
			os.Exit(0)
		case exitError:
			os.Exit(1)
		case exitQuit:
			os.Exit(2)
		case exitQuestion:
			os.Exit(3)
		default:
			os.Exit(0)
		}
	}
}
