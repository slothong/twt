#!/usr/bin/env bash

set -euo pipefail

INSTALL_DIR="${INSTALL_DIR:-$HOME/.local/bin}"
BINARY_NAME="tmux-worktree"

echo "Building $BINARY_NAME..."
go build -o "$BINARY_NAME" ./cmd/tmux-worktree

echo "Installing to $INSTALL_DIR..."
mkdir -p "$INSTALL_DIR"
mv "$BINARY_NAME" "$INSTALL_DIR/"
chmod +x "$INSTALL_DIR/$BINARY_NAME"

echo ""
echo "✓ $BINARY_NAME installed successfully!"
echo "  Location: $INSTALL_DIR/$BINARY_NAME"
echo ""

# Check if INSTALL_DIR is in PATH
if [[ ":$PATH:" != *":$INSTALL_DIR:"* ]]; then
  echo "Warning: $INSTALL_DIR is not in your PATH"
  echo "Add this to your ~/.zshrc or ~/.bashrc:"
  echo "  export PATH=\"\$PATH:$INSTALL_DIR\""
  echo ""
fi

echo "Usage:"
echo "  $BINARY_NAME create-sessions       # Create tmux session with all worktrees"
echo "  $BINARY_NAME create-window <name> <branch>  # Create new worktree in tmux window"
echo "  $BINARY_NAME remove-window          # Remove current worktree and close window"
echo ""
echo "Run '$BINARY_NAME --help' for more information"
