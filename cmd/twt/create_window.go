package main

import (
	"fmt"
	"path/filepath"

	"github.com/slothong/twt/internal/git"
	"github.com/slothong/twt/internal/tmux"
	"github.com/slothong/twt/internal/ui"
	"github.com/spf13/cobra"
)

func newCreateWindowCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-window <worktree-name> <branch-name> [base-branch]",
		Short: "Create a new worktree in a new tmux window",
		Long: `Creates a new git worktree and opens it in a new tmux window.
Must be run inside a tmux session.`,
		Args: cobra.RangeArgs(2, 3),
		RunE: func(cmd *cobra.Command, args []string) error {
			worktreeName := args[0]
			branchName := args[1]
			baseBranch := "HEAD"
			if len(args) > 2 {
				baseBranch = args[2]
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
