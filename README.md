# twt

A CLI tool for managing git worktrees with tmux integration.

## Features

- 🚀 Create tmux sessions with windows for all git worktrees
- 📦 Create new worktrees in new tmux windows
- 🗑️  Remove worktrees and close tmux windows
- 🎨 Automatic IDE layout setup (configurable pane splits)

## Installation

### Binary (recommended)

[Releases 페이지](https://github.com/slothong/twt/releases)에서 바이너리를 다운로드하거나, curl로 설치:

```bash
# macOS (Apple Silicon)
curl -L https://github.com/slothong/twt/releases/latest/download/twt_darwin_arm64.tar.gz | tar xz
sudo mv twt /usr/local/bin/

# macOS (Intel)
curl -L https://github.com/slothong/twt/releases/latest/download/twt_darwin_amd64.tar.gz | tar xz
sudo mv twt /usr/local/bin/

# Linux (amd64)
curl -L https://github.com/slothong/twt/releases/latest/download/twt_linux_amd64.tar.gz | tar xz
sudo mv twt /usr/local/bin/
```

### Using Go Install

```bash
go install github.com/slothong/twt/cmd/twt@latest
```

### From Source

```bash
git clone https://github.com/slothong/twt
cd twt
make install
```

## Usage

### Create Sessions

Create a tmux session with a window for each git worktree:

```bash
twt create-sessions
# or use alias
twt cs
```

This will:
1. Find all git worktrees in the current repository
2. Create a new tmux session named `pointhub-worktrees`
3. Create a window for each worktree

### Create Window

Create a new worktree in a new tmux window (must be run inside tmux):

```bash
twt create-window <branch-name> [worktree-name] [base-branch]
# or use aliases
twt cw <branch-name> [worktree-name] [base-branch]
twt new <branch-name> [worktree-name] [base-branch]
```

If `worktree-name` is not provided, it will be auto-generated from `branch-name` by replacing slashes with hyphens.

Examples:
```bash
# Auto-generate worktree name from branch (feature/foo -> feature-foo)
twt cw feature/foo

# Auto-generate with custom base branch
twt new feature/foo main

# Use custom worktree name
twt cw feature/foo custom-name

# Use custom worktree name with base branch
twt new feature/foo custom-name main
```

### Remove Window

Remove the current worktree and close the tmux window (must be run inside tmux from a worktree):

```bash
twt remove-window
# or use aliases
twt rw
twt rm
```

This will:
1. Check if current directory is a worktree
2. Warn if there are uncommitted changes
3. Ask for confirmation
4. Remove the worktree
5. Close the current tmux window

Options:
- `-f, --force`: Force removal even with uncommitted changes

Examples:
```bash
# Remove with confirmation
twt rm

# Force remove without confirmation
twt rm -f
```

## Requirements

- Go 1.21 or higher
- tmux
- git

## Development

```bash
# Clone the repository
git clone https://github.com/slothong/twt
cd twt

# Install dependencies
go mod download

# Build
go build -o twt ./cmd/twt

# Run
./twt --help
```

## License

MIT
