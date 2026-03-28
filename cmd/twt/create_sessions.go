package main

import (
	"fmt"

	"github.com/slothong/twt/internal/git"
	"github.com/slothong/twt/internal/tmux"
	"github.com/slothong/twt/internal/ui"
	"github.com/spf13/cobra"
)

func newCreateSessionsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-sessions",
		Short: "Create tmux session with windows for all git worktrees",
		Long:  `Creates a tmux session with a window for each git worktree.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCreateSessions()
		},
	}

	return cmd
}

func runCreateSessions() error {
	sessionName := "pointhub-worktrees"

	// Check if session already exists
	if tmux.HasSession(sessionName) {
		ui.Error(fmt.Sprintf("tmux session '%s' already exists", sessionName))
		ui.Info(fmt.Sprintf("Attach to it with: tmux attach -t %s", sessionName))
		ui.Info(fmt.Sprintf("Or kill it with: tmux kill-session -t %s", sessionName))
		return fmt.Errorf("session already exists")
	}

	// Get all worktrees
	worktrees, err := git.ListWorktrees()
	if err != nil {
		return fmt.Errorf("failed to list worktrees: %w", err)
	}

	if len(worktrees) == 0 {
		ui.Error("No git worktrees found")
		return fmt.Errorf("no worktrees")
	}

	ui.Info(fmt.Sprintf("Found %d worktree(s):", len(worktrees)))
	for _, wt := range worktrees {
		ui.Info(fmt.Sprintf("  %s", wt.Path))
	}
	fmt.Println()

	// Create session with first worktree
	first := worktrees[0]
	firstName := git.GetWorktreeName(first.Path)
	ui.Info(fmt.Sprintf("Creating tmux session '%s' with first window '%s'", sessionName, firstName))

	if err := tmux.NewSession(sessionName, firstName, first.Path); err != nil {
		return err
	}

	// Create windows for remaining worktrees
	for i := 1; i < len(worktrees); i++ {
		wt := worktrees[i]
		name := git.GetWorktreeName(wt.Path)
		ui.Info(fmt.Sprintf("Creating window '%s'", name))

		if err := tmux.NewWindow(sessionName, name, wt.Path); err != nil {
			return err
		}
	}

	fmt.Println()
	ui.Success(fmt.Sprintf("tmux session '%s' created successfully!", sessionName))

	// Attach or switch to session
	if tmux.IsInsideTmux() {
		ui.Info(fmt.Sprintf("Switching to session '%s'...", sessionName))
		return tmux.SwitchClient(sessionName)
	} else {
		ui.Info(fmt.Sprintf("Attaching to session '%s'...", sessionName))
		return tmux.AttachSession(sessionName)
	}
}
