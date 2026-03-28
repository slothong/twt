package main

import (
	"fmt"

	"github.com/Buzzvil/tmux-worktree/internal/git"
	"github.com/Buzzvil/tmux-worktree/internal/tmux"
	"github.com/Buzzvil/tmux-worktree/internal/ui"
	"github.com/spf13/cobra"
)

func newRemoveWindowCmd() *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "remove-window",
		Short: "Remove current worktree and close tmux window",
		Long: `Removes the current git worktree and closes the current tmux window.
Must be run inside a tmux session from a worktree directory.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runRemoveWindow(force)
		},
	}

	cmd.Flags().BoolVarP(&force, "force", "f", false, "Force removal even with uncommitted changes")

	return cmd
}

func runRemoveWindow(force bool) error {
	// Check if running inside tmux
	if !tmux.IsInsideTmux() {
		ui.Error("This command must be run inside a tmux session")
		return fmt.Errorf("not in tmux")
	}

	// Get current worktree
	currentWT, err := git.GetCurrentWorktree()
	if err != nil {
		ui.Error("Current directory is not a git worktree")
		return err
	}

	// Check if it's the main worktree
	if currentWT.IsMain {
		ui.Error("Cannot remove the main worktree")
		ui.Info(fmt.Sprintf("Main worktree: %s", currentWT.Path))
		return fmt.Errorf("cannot remove main worktree")
	}

	// Get current window name
	windowName, err := tmux.GetCurrentWindowName()
	if err != nil {
		return err
	}

	ui.Info("This will remove the worktree and close the current window:")
	ui.Info(fmt.Sprintf("  Path: %s", currentWT.Path))
	ui.Info(fmt.Sprintf("  Window: %s", windowName))
	fmt.Println()

	// Check for uncommitted changes
	if !force {
		hasChanges, err := git.HasUncommittedChanges()
		if err != nil {
			return err
		}

		if hasChanges {
			ui.Warning("You have uncommitted changes!")
			fmt.Println()
			if !ui.Confirm("Continue anyway?") {
				ui.Info("Aborted")
				return nil
			}
			force = true // Set force flag for removal
		}
	}

	// Confirm removal
	if !ui.Confirm("Are you sure?") {
		ui.Info("Aborted")
		return nil
	}

	fmt.Println()
	ui.Info("Removing worktree...")

	// Switch to window 0 before closing current window
	if err := tmux.SelectWindow(":0"); err != nil {
		ui.Warning(fmt.Sprintf("Failed to switch window: %v", err))
	}

	// Remove worktree
	if err := git.RemoveWorktree(currentWT.Path, force); err != nil {
		return err
	}

	// Kill current window
	if err := tmux.KillWindow(windowName); err != nil {
		ui.Warning(fmt.Sprintf("Failed to close window: %v", err))
	}

	ui.Success("Worktree removed and window closed!")
	ui.Info(fmt.Sprintf("  Removed: %s", currentWT.Path))

	return nil
}
