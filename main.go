package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func main() {
	// This is a simple wrapper around the actual implementation in cmd/revengo
	// It allows users to run the tool directly from the repository root

	// Get the path to the cmd/revengo directory
	execPath, err := os.Executable()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	dir := filepath.Dir(execPath)
	cmdPath := filepath.Join(dir, "cmd", "revengo", "main.go")

	// Check if the file exists
	if _, err := os.Stat(cmdPath); os.IsNotExist(err) {
		// Try a relative path from current directory
		cmdPath = filepath.Join("cmd", "revengo", "main.go")
		if _, err := os.Stat(cmdPath); os.IsNotExist(err) {
			fmt.Printf("Error: Could not find the main implementation at %s\n", cmdPath)
			os.Exit(1)
		}
	}

	// Execute the actual implementation
	cmd := exec.Command("go", append([]string{"run", cmdPath}, os.Args[1:]...)...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			os.Exit(exitErr.ExitCode())
		}
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
