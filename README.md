# tmux-worktree

A CLI tool for managing git worktrees with tmux integration.

## Features

- 🚀 Create tmux sessions with windows for all git worktrees
- 📦 Create new worktrees in new tmux windows
- 🗑️  Remove worktrees and close tmux windows
- 🎨 Automatic IDE layout setup (configurable pane splits)

## Installation

### From Source

```bash
# Clone the repository
git clone https://github.com/Buzzvil/tmux-worktree
cd tmux-worktree

# Build and install
go build -o tmux-worktree ./cmd/tmux-worktree
sudo mv tmux-worktree /usr/local/bin/

# Or use the install script
./install.sh
```

### Using Go Install

```bash
go install github.com/Buzzvil/tmux-worktree/cmd/tmux-worktree@latest
```

## Usage

### Create Sessions

Create a tmux session with a window for each git worktree:

```bash
tmux-worktree create-sessions
```

This will:
1. Find all git worktrees in the current repository
2. Create a new tmux session named `pointhub-worktrees`
3. Create a window for each worktree
4. Set up an IDE layout in each window (optional)

Options:
- `--no-ide`: Skip IDE layout setup

### Create Window

Create a new worktree in a new tmux window (must be run inside tmux):

```bash
tmux-worktree create-window <worktree-name> <branch-name> [base-branch]
```

Examples:
```bash
# Create worktree from current HEAD
tmux-worktree create-window feature-foo feature/foo

# Create worktree from main branch
tmux-worktree create-window feature-foo feature/foo main
```

Options:
- `--no-ide`: Skip IDE layout setup

### Remove Window

Remove the current worktree and close the tmux window (must be run inside tmux from a worktree):

```bash
tmux-worktree remove-window
```

This will:
1. Check if current directory is a worktree
2. Warn if there are uncommitted changes
3. Ask for confirmation
4. Remove the worktree
5. Close the current tmux window

Options:
- `-f, --force`: Force removal even with uncommitted changes

## IDE Layout

The default IDE layout creates 4 panes:

```
┌─────────────────┬──────┐
│                 │      │
│   pane 0        │pane 3│
│   (70%)         │(30%) │
│                 │claude│
├────────┬────────┤      │
│ pane 1 │ pane 2 │      │
│  (35%) │  (35%) │      │
└────────┴────────┴──────┘
```

- **Pane 0 (top-left)**: Main editor area (75% height)
- **Pane 1 (bottom-left)**: Terminal (25% height, 50% width)
- **Pane 2 (bottom-right)**: Terminal (25% height, 50% width)
- **Pane 3 (right)**: Claude AI (30% width, full height)

## Requirements

- Go 1.21 or higher
- tmux
- git

## Development

```bash
# Clone the repository
git clone https://github.com/Buzzvil/tmux-worktree
cd tmux-worktree

# Install dependencies
go mod download

# Build
go build -o tmux-worktree ./cmd/tmux-worktree

# Run
./tmux-worktree --help
```

## License

MIT
