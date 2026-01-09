#!/bin/bash

# Build script for diff-review TUI
set -e

cd "$(dirname "$0")/diff-review-tui"

echo "Initializing Go module..."
go mod init diff-review-tui 2>/dev/null || true

echo "Installing dependencies..."
go get github.com/charmbracelet/bubbletea@latest
go get github.com/charmbracelet/bubbles/viewport@latest
go get github.com/charmbracelet/bubbles/textinput@latest
go get github.com/charmbracelet/lipgloss@latest
go get github.com/alecthomas/chroma/v2@latest
go get github.com/bluekeyes/go-gitdiff/gitdiff@latest

echo "Building binary..."
go build -o ../diff-review -ldflags="-s -w" .

echo "Making binary executable..."
chmod +x ../diff-review

echo ""
echo "✓ Build complete!"
echo "Binary created at: $(dirname "$0")/diff-review"
