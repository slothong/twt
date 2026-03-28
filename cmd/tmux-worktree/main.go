package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var version = "dev"

func main() {
	rootCmd := &cobra.Command{
		Use:   "tmux-worktree",
		Short: "Manage git worktrees with tmux integration",
		Long: `tmux-worktree is a CLI tool for managing git worktrees with tmux integration.
It allows you to create tmux sessions/windows for worktrees and manage them easily.`,
		Version: version,
	}

	rootCmd.AddCommand(
		newCreateSessionsCmd(),
		newCreateWindowCmd(),
		newRemoveWindowCmd(),
	)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
