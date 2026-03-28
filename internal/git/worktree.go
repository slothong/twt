package git

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Worktree represents a git worktree
type Worktree struct {
	Path   string
	Branch string
	IsMain bool
}

// ListWorktrees returns all git worktrees in the current repository
func ListWorktrees() ([]Worktree, error) {
	cmd := exec.Command("git", "worktree", "list", "--porcelain")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to list worktrees: %w", err)
	}

	return parseWorktrees(string(output)), nil
}

// GetRepoRoot returns the root directory of the git repository
func GetRepoRoot() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("not in a git repository: %w", err)
	}

	return strings.TrimSpace(string(output)), nil
}

// CreateWorktree creates a new git worktree with a new branch
func CreateWorktree(path, branch, baseBranch string) error {
	// Check if branch already exists
	if branchExists(branch) {
		return fmt.Errorf("branch '%s' already exists", branch)
	}

	// Check if directory already exists
	if _, err := os.Stat(path); err == nil {
		return fmt.Errorf("directory '%s' already exists", path)
	}

	// Create worktree
	cmd := exec.Command("git", "worktree", "add", "-b", branch, path, baseBranch)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to create worktree: %s", string(output))
	}

	return nil
}

// RemoveWorktree removes a git worktree
func RemoveWorktree(path string, force bool) error {
	args := []string{"worktree", "remove"}
	if force {
		args = append(args, "--force")
	}
	args = append(args, path)

	cmd := exec.Command("git", args...)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to remove worktree: %s", string(output))
	}

	return nil
}

// IsWorktree checks if the current directory is a git worktree
func IsWorktree() (bool, error) {
	if _, err := exec.Command("git", "rev-parse", "--git-dir").Output(); err != nil {
		return false, fmt.Errorf("not in a git repository")
	}

	worktrees, err := ListWorktrees()
	if err != nil {
		return false, err
	}

	cwd, err := os.Getwd()
	if err != nil {
		return false, err
	}

	for _, wt := range worktrees {
		if wt.Path == cwd {
			return true, nil
		}
	}

	return false, nil
}

// GetCurrentWorktree returns the current worktree information
func GetCurrentWorktree() (*Worktree, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	worktrees, err := ListWorktrees()
	if err != nil {
		return nil, err
	}

	for _, wt := range worktrees {
		if wt.Path == cwd {
			return &wt, nil
		}
	}

	return nil, fmt.Errorf("current directory is not a worktree")
}

// HasUncommittedChanges checks if there are uncommitted changes
func HasUncommittedChanges() (bool, error) {
	// Check working tree
	cmd := exec.Command("git", "diff", "--quiet")
	if err := cmd.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 1 {
			return true, nil
		}
		return false, err
	}

	// Check staged changes
	cmd = exec.Command("git", "diff", "--cached", "--quiet")
	if err := cmd.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 1 {
			return true, nil
		}
		return false, err
	}

	return false, nil
}

// parseWorktrees parses the output of 'git worktree list --porcelain'
func parseWorktrees(output string) []Worktree {
	var worktrees []Worktree
	var current Worktree
	isFirst := true

	for _, line := range strings.Split(output, "\n") {
		if line == "" {
			if current.Path != "" {
				current.IsMain = isFirst
				worktrees = append(worktrees, current)
				current = Worktree{}
				isFirst = false
			}
			continue
		}

		if strings.HasPrefix(line, "worktree ") {
			current.Path = strings.TrimPrefix(line, "worktree ")
		} else if strings.HasPrefix(line, "branch ") {
			current.Branch = strings.TrimPrefix(line, "branch refs/heads/")
		}
	}

	// Add last worktree if exists
	if current.Path != "" {
		current.IsMain = isFirst
		worktrees = append(worktrees, current)
	}

	return worktrees
}

// branchExists checks if a git branch exists
func branchExists(branch string) bool {
	cmd := exec.Command("git", "rev-parse", "--verify", branch)
	return cmd.Run() == nil
}

// GetWorktreeName returns the name of a worktree (basename of path)
func GetWorktreeName(path string) string {
	return filepath.Base(path)
}
