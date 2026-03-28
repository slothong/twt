package main

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/slothong/twt/internal/git"
	"github.com/slothong/twt/internal/tmux"
	"github.com/slothong/twt/internal/ui"
	"github.com/spf13/cobra"
)

func newCreateWindowCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create-window <branch-name> [worktree-name] [base-branch]",
		Aliases: []string{"cw", "new"},
		Short:   "Create a new worktree in a new tmux window",
		Long: `Creates a new git worktree and opens it in a new tmux window.
Must be run inside a tmux session.

If worktree-name is not provided, it will be auto-generated from branch-name
by replacing slashes with hyphens (e.g., feature/foo -> feature-foo).`,
		Args: cobra.RangeArgs(1, 3),
		RunE: func(cmd *cobra.Command, args []string) error {
			branchName := args[0]
			worktreeName := ""
			baseBranch := "HEAD"

			// Parse arguments based on count
			if len(args) >= 2 {
				worktreeName = args[1]
			}
			if len(args) >= 3 {
				baseBranch = args[2]
			}

			// Auto-generate worktree name if not provided
			if worktreeName == "" {
				worktreeName = strings.ReplaceAll(branchName, "/", "-")
			}

			return runCreateWindow(worktreeName, branchName, baseBranch)
		},
	}

	return cmd
}

func runCreateWindow(worktreeName, branchName, baseBranch string) error {
	// Check if running inside tmux
	if !tmux.IsInsideTmux() {
		ui.Error("This command must be run inside a tmux session")
		return fmt.Errorf("not in tmux")
	}

	// Get repository root
	repoRoot, err := git.GetRepoRoot()
	if err != nil {
		return err
	}

	// Calculate worktree path
	repoParent := filepath.Dir(repoRoot)
	worktreePath := filepath.Join(repoParent, worktreeName)

	ui.Info(fmt.Sprintf("Creating worktree '%s' with branch '%s' from '%s'...", worktreeName, branchName, baseBranch))

	// Create worktree
	if err := git.CreateWorktree(worktreePath, branchName, baseBranch); err != nil {
		return err
	}

	// Create new tmux window
	ui.Info(fmt.Sprintf("Creating tmux window '%s'...", worktreeName))
	if err := tmux.NewWindow("", worktreeName, worktreePath); err != nil {
		// If window creation fails, try to clean up the worktree
		git.RemoveWorktree(worktreePath, true)
		return err
	}

	ui.Success("Worktree window created!")
	ui.Info(fmt.Sprintf("  Path: %s", worktreePath))
	ui.Info(fmt.Sprintf("  Branch: %s", branchName))

	return nil
}
