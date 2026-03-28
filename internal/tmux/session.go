package tmux

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Window represents a tmux window
type Window struct {
	Index int
	Name  string
}

// IsInsideTmux checks if the current process is inside a tmux session
func IsInsideTmux() bool {
	return os.Getenv("TMUX") != ""
}

// HasSession checks if a tmux session exists
func HasSession(name string) bool {
	cmd := exec.Command("tmux", "has-session", "-t", name)
	return cmd.Run() == nil
}

// NewSession creates a new tmux session
func NewSession(name, windowName, path string) error {
	cmd := exec.Command("tmux", "new-session", "-d", "-s", name, "-n", windowName, "-c", path)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to create session: %s", string(output))
	}
	return nil
}

// NewWindow creates a new tmux window
func NewWindow(session, name, path string) error {
	cmd := exec.Command("tmux", "new-window", "-t", session, "-n", name, "-c", path)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to create window: %s", string(output))
	}
	return nil
}

// KillWindow kills a tmux window
func KillWindow(target string) error {
	cmd := exec.Command("tmux", "kill-window", "-t", target)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to kill window: %s", string(output))
	}
	return nil
}

// SelectWindow selects a tmux window
func SelectWindow(target string) error {
	cmd := exec.Command("tmux", "select-window", "-t", target)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to select window: %s", string(output))
	}
	return nil
}

// SelectPane selects a tmux pane
func SelectPane(target string) error {
	cmd := exec.Command("tmux", "select-pane", "-t", target)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to select pane: %s", string(output))
	}
	return nil
}

// SplitWindow splits a tmux window/pane
func SplitWindow(target string, horizontal bool, percentage int, command string) error {
	args := []string{"split-window", "-t", target}

	if horizontal {
		args = append(args, "-h")
	} else {
		args = append(args, "-v")
	}

	args = append(args, "-p", fmt.Sprintf("%d", percentage))

	if command != "" {
		args = append(args, command)
	}

	cmd := exec.Command("tmux", args...)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to split window: %s", string(output))
	}
	return nil
}

// SendKeys sends keys to a tmux pane
func SendKeys(target, keys string) error {
	cmd := exec.Command("tmux", "send-keys", "-t", target, keys, "C-m")
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to send keys: %s", string(output))
	}
	return nil
}

// ListWindows returns all windows in a tmux session
func ListWindows(session string) ([]Window, error) {
	cmd := exec.Command("tmux", "list-windows", "-t", session, "-F", "#{window_index}:#{window_name}")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to list windows: %w", err)
	}

	var windows []Window
	for _, line := range strings.Split(strings.TrimSpace(string(output)), "\n") {
		if line == "" {
			continue
		}

		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}

		var idx int
		fmt.Sscanf(parts[0], "%d", &idx)

		windows = append(windows, Window{
			Index: idx,
			Name:  parts[1],
		})
	}

	return windows, nil
}

// AttachSession attaches to a tmux session
func AttachSession(name string) error {
	cmd := exec.Command("tmux", "attach", "-t", name)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// SwitchClient switches to another tmux session (when inside tmux)
func SwitchClient(name string) error {
	cmd := exec.Command("tmux", "switch-client", "-t", name)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to switch client: %s", string(output))
	}
	return nil
}

// GetCurrentWindowName returns the current window name
func GetCurrentWindowName() (string, error) {
	cmd := exec.Command("tmux", "display-message", "-p", "#W")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get window name: %w", err)
	}
	return strings.TrimSpace(string(output)), nil
}

