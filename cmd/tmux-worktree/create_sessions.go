package main

import (
	"fmt"
	"time"

	"github.com/Buzzvil/tmux-worktree/internal/git"
	"github.com/Buzzvil/tmux-worktree/internal/tmux"
	"github.com/Buzzvil/tmux-worktree/internal/ui"
	"github.com/spf13/cobra"
)

func newCreateSessionsCmd() *cobra.Command {
	var noIDE bool

	cmd := &cobra.Command{
		Use:   "create-sessions",
		Short: "Create tmux session with windows for all git worktrees",
		Long: `Creates a tmux session with a window for each git worktree.
Each window will have the IDE layout set up automatically.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCreateSessions(noIDE)
		},
	}

	cmd.Flags().BoolVar(&noIDE, "no-ide", false, "Don't set up IDE layout in windows")

	return cmd
}

func runCreateSessions(noIDE bool) error {
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

	// Set up IDE layout for each window if requested
	if !noIDE {
		fmt.Println()
		ui.Info("Setting up IDE layout in each window...")

		windows, err := tmux.ListWindows(sessionName)
		if err != nil {
			return err
		}

		for _, win := range windows {
			ui.Info(fmt.Sprintf("  -> Setting up window '%s' (index: %d)", win.Name, win.Index))
			target := fmt.Sprintf("%s:%d", sessionName, win.Index)

			if err := tmux.SetupIDELayout(target); err != nil {
				ui.Warning(fmt.Sprintf("Failed to set up IDE layout for window '%s': %v", win.Name, err))
				continue
			}

			time.Sleep(500 * time.Millisecond)
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
