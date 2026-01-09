#!/bin/bash

# Verification script for Visual Diff Review System

echo "Visual Diff Review System - Installation Verification"
echo "======================================================="
echo ""

# Check for required commands
check_command() {
    if command -v "$1" &> /dev/null; then
        echo "✓ $1 is installed"
        return 0
    else
        echo "✗ $1 is NOT installed"
        return 1
    fi
}

# Check for files
check_file() {
    if [ -f "$1" ]; then
        echo "✓ $1 exists"
        return 0
    else
        echo "✗ $1 is MISSING"
        return 1
    fi
}

# Check for executable
check_executable() {
    if [ -x "$1" ]; then
        echo "✓ $1 is executable"
        return 0
    else
        echo "✗ $1 is NOT executable (run: chmod +x $1)"
        return 1
    fi
}

echo "1. Checking Prerequisites:"
echo "--------------------------"
check_command git
check_command go
check_command jq
echo ""

echo "2. Checking Command Files:"
echo "--------------------------"
check_file "commands/diff-review.md"
check_file "commands/apply-review.md"
check_file "prompts/apply-review-comments.md"
check_file "templates/review-summary.md"
echo ""

echo "3. Checking Modified Files:"
echo "--------------------------"
check_file "hooks/session-start.sh"
check_file "commands/pr-status.md"
check_file "settings.json"
echo ""

echo "4. Checking Go Source Files:"
echo "--------------------------"
check_file "scripts/diff-review-tui/main.go"
check_file "scripts/diff-review-tui/model.go"
check_file "scripts/diff-review-tui/update.go"
check_file "scripts/diff-review-tui/view.go"
check_file "scripts/diff-review-tui/diff.go"
check_file "scripts/diff-review-tui/metadata.go"
check_file "scripts/diff-review-tui/go.mod"
echo ""

echo "5. Checking TUI Binary:"
echo "--------------------------"
if [ -f "scripts/diff-review" ]; then
    check_executable "scripts/diff-review"
    echo ""
    echo "Testing binary:"
    scripts/diff-review --help 2>&1 | head -3
else
    echo "✗ scripts/diff-review binary NOT BUILT"
    echo ""
    echo "To build the binary, run:"
    echo "  cd scripts/diff-review-tui"
    echo "  go mod download"
    echo "  go build -o ../diff-review -ldflags=\"-s -w\" ."
    echo "  chmod +x ../diff-review"
fi
echo ""

echo "6. Checking Script Permissions:"
echo "--------------------------"
check_executable "hooks/session-start.sh"
check_executable "scripts/compact-history.sh"
echo ""

echo "======================================================="
echo "Verification Complete!"
echo ""
echo "Next Steps:"
echo "1. Build the TUI binary (if not already built)"
echo "2. Try: /diff-review"
echo "3. Read the README.md for full documentation"
