#!/bin/bash
set -e

REPO_DIR="$(cd "$(dirname "$0")" && pwd)"
CLAUDE_DIR="$HOME/.claude"

echo "dotclaude installer"
echo "==================="

# Ensure ~/.claude directories exist
mkdir -p "$CLAUDE_DIR/skills" "$CLAUDE_DIR/hooks"

# Symlink skills (each skill dir gets symlinked)
echo ""
echo "Linking skills..."
for skill_dir in "$REPO_DIR"/skills/*/; do
  skill_name=$(basename "$skill_dir")
  target="$CLAUDE_DIR/skills/$skill_name"
  if [ -L "$target" ]; then
    unlink "$target"
  elif [ -d "$target" ]; then
    echo "  WARNING: $target is a real directory, skipping (back up and remove manually)"
    continue
  fi
  ln -sv "$skill_dir" "$target"
done

# Symlink hooks
echo ""
echo "Linking hooks..."
for hook in "$REPO_DIR"/hooks/*.sh; do
  hook_name=$(basename "$hook")
  target="$CLAUDE_DIR/hooks/$hook_name"
  if [ -L "$target" ]; then
    unlink "$target"
  elif [ -f "$target" ]; then
    echo "  WARNING: $target exists as a regular file, backing up to ${target}.bak"
    mv "$target" "${target}.bak"
  fi
  ln -sv "$hook" "$target"
done
chmod +x "$CLAUDE_DIR/hooks/"*.sh 2>/dev/null || true

# Copy settings.json (don't symlink — user may have local overrides)
echo ""
if [ -f "$CLAUDE_DIR/settings.json" ]; then
  echo "settings.json already exists — not overwriting."
  echo "  Compare with: diff $CLAUDE_DIR/settings.json $REPO_DIR/settings.json"
else
  cp -v "$REPO_DIR/settings.json" "$CLAUDE_DIR/settings.json"
fi

# Copy CLAUDE.md to default project
echo ""
DEFAULT_PROJECT="$CLAUDE_DIR/projects/-Users-$(whoami)"
mkdir -p "$DEFAULT_PROJECT"
if [ -f "$DEFAULT_PROJECT/CLAUDE.md" ]; then
  echo "CLAUDE.md already exists — not overwriting."
  echo "  Compare with: diff $DEFAULT_PROJECT/CLAUDE.md $REPO_DIR/CLAUDE.md"
else
  cp -v "$REPO_DIR/CLAUDE.md" "$DEFAULT_PROJECT/CLAUDE.md"
fi

echo ""
echo "Done. Skills and hooks are symlinked — edits in this repo update ~/.claude/ automatically."
echo ""
echo "Next steps:"
echo "  1. Review/merge settings.json if it wasn't copied"
echo "  2. Add to ~/.zshrc (for AeroSpace hooks):"
echo '     if [[ $- == *i* ]] && [[ -z "$AEROSPACE_WS" ]]; then'
echo '       export AEROSPACE_WS=$(aerospace list-workspaces --focused 2>/dev/null)'
echo '     fi'
