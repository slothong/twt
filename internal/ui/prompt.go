package ui

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Confirm prompts the user for a yes/no confirmation
func Confirm(message string) bool {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("%s (y/N): ", message)

	response, err := reader.ReadString('\n')
	if err != nil {
		return false
	}

	response = strings.ToLower(strings.TrimSpace(response))
	return response == "y" || response == "yes"
}

// Success prints a success message
func Success(message string) {
	fmt.Printf("✓ %s\n", message)
}

// Error prints an error message
func Error(message string) {
	fmt.Fprintf(os.Stderr, "Error: %s\n", message)
}

// Info prints an info message
func Info(message string) {
	fmt.Println(message)
}

// Warning prints a warning message
func Warning(message string) {
	fmt.Printf("Warning: %s\n", message)
}
