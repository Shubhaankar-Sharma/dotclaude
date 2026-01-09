package main

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"

	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/formatters"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
	"github.com/bluekeyes/go-gitdiff/gitdiff"
)

// getDiffWithMappings gets the diff between two commits and builds line mappings
func getDiffWithMappings(baseCommit, headCommit string) (string, []LineMapping, error) {
	// Get raw diff
	cmd := exec.Command("git", "diff", "--unified=3", "--color=never",
		fmt.Sprintf("%s..%s", baseCommit, headCommit))

	rawDiff, err := cmd.Output()
	if err != nil {
		return "", nil, fmt.Errorf("failed to get diff: %w", err)
	}

	diffText := string(rawDiff)

	// Parse with go-gitdiff to build line mappings
	files, _, err := gitdiff.Parse(strings.NewReader(diffText))
	if err != nil {
		return "", nil, fmt.Errorf("failed to parse diff: %w", err)
	}

	// Build line mappings
	mappings := buildLineMappings(files)

	// Syntax highlight the diff
	highlighted := highlightDiff(diffText)

	return highlighted, mappings, nil
}

// buildLineMappings creates a mapping from diff line numbers to file locations
func buildLineMappings(files []*gitdiff.File) []LineMapping {
	var mappings []LineMapping
	diffLineNum := 0

	for _, file := range files {
		// Skip binary files
		if file.IsBinary || file.IsNew || file.IsDelete {
			// Add a simple mapping for the file header
			mappings = append(mappings, LineMapping{
				DiffLine: diffLineNum,
				File:     getFileName(file),
				OldLine:  0,
				NewLine:  0,
				Content:  fmt.Sprintf("diff --git a/%s b/%s", file.OldName, file.NewName),
				Op:       "header",
			})
			diffLineNum++
			continue
		}

		// Process text fragments
		for _, fragment := range file.TextFragments {
			oldLine := int(fragment.OldPosition)
			newLine := int(fragment.NewPosition)

			for _, line := range fragment.Lines {
				mapping := LineMapping{
					DiffLine: diffLineNum,
					File:     getFileName(file),
					Content:  string(line.Line),
				}

				switch line.Op {
				case gitdiff.OpAdd:
					mapping.NewLine = newLine
					mapping.Op = "add"
					newLine++
				case gitdiff.OpDelete:
					mapping.OldLine = oldLine
					mapping.Op = "delete"
					oldLine++
				case gitdiff.OpContext:
					mapping.OldLine = oldLine
					mapping.NewLine = newLine
					mapping.Op = "context"
					oldLine++
					newLine++
				}

				mappings = append(mappings, mapping)
				diffLineNum++
			}
		}
	}

	return mappings
}

// getFileName returns the appropriate filename from a File object
func getFileName(file *gitdiff.File) string {
	if file.NewName != "" && file.NewName != "/dev/null" {
		return file.NewName
	}
	if file.OldName != "" && file.OldName != "/dev/null" {
		return file.OldName
	}
	return "unknown"
}

// highlightDiff applies syntax highlighting to diff text
func highlightDiff(diffText string) string {
	// Get the diff lexer
	lexer := lexers.Get("diff")
	if lexer == nil {
		lexer = lexers.Fallback
	}

	// Use a style that works well in terminals
	style := styles.Get("monokai")
	if style == nil {
		style = styles.Fallback
	}

	// Use terminal256 formatter for better color support
	formatter := formatters.Get("terminal256")
	if formatter == nil {
		formatter = formatters.Fallback
	}

	// Tokenize and format
	var buf bytes.Buffer
	iterator, err := lexer.Tokenise(nil, diffText)
	if err != nil {
		// If highlighting fails, return plain text
		return diffText
	}

	err = formatter.Format(&buf, style, iterator)
	if err != nil {
		// If formatting fails, return plain text
		return diffText
	}

	return buf.String()
}

// highlightCode applies syntax highlighting to a code snippet based on file extension
func highlightCode(code, filename string) string {
	// Get lexer based on filename
	lexer := lexers.Match(filename)
	if lexer == nil {
		lexer = lexers.Fallback
	}

	// Ensure we get a non-nil lexer
	lexer = chroma.Coalesce(lexer)

	style := styles.Get("monokai")
	if style == nil {
		style = styles.Fallback
	}

	formatter := formatters.Get("terminal256")
	if formatter == nil {
		formatter = formatters.Fallback
	}

	var buf bytes.Buffer
	iterator, err := lexer.Tokenise(nil, code)
	if err != nil {
		return code
	}

	err = formatter.Format(&buf, style, iterator)
	if err != nil {
		return code
	}

	return buf.String()
}
