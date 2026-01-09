# Build Instructions for Diff Review TUI

## Prerequisites

- Go 1.21 or later
- Git

## Build Steps

```bash
cd scripts/diff-review-tui

# Download dependencies
go mod download

# Build the binary
go build -o ../diff-review -ldflags="-s -w" .

# Make it executable
chmod +x ../diff-review
```

## Verify

```bash
../diff-review --help
```

## One-liner

```bash
cd scripts/diff-review-tui && go mod download && go build -o ../diff-review -ldflags="-s -w" . && chmod +x ../diff-review
```
