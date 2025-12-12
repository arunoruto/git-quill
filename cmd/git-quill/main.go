package main

import (
	"fmt"
	"os"
	"os/exec"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "commit":
		runCommit()
	case "tag":
		fmt.Println("Tag generation feature coming soon!")
	default:
		fmt.Printf("Unknown command: %s\n", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Usage: git quill <command> [options]")
	fmt.Println("\nCommands:")
	fmt.Println("\t commit Generate an AI commit message")
	fmt.Println("\t tag    Generate an AI tag message")
}

func runCommit() {
	if !hasStagedChanges() {
		fmt.Println("Error: No staged changes to commit.")
		os.Exit(1)
	}

	providers := findAvailableProviders()
	if len(providers) == 0 {
		fmt.Println("Error: No AI providers found (ollama, opencode, etc)")
		os.Exit(1)
	}

	selectedProvider := gumChoose("Select AI Provider", providers)
	fmt.Printf("You selected: %s\n", selectedProvider)
}

func hasStagedChanges() bool {
	cmd := exec.Command("git", "diff", "--cached", "--quiet")
	err := cmd.Run()
	return err != nil
}

func findAvailableProviders() []string {
	candidates := []string{"ollama", "opencode", "copilot", "gemini"}
	var found []string
	for _, p := range candidates {
		if _, err := exec.LookPath(p); err == nil {
			found = append(found, p)
		}
	}
	return found
}

func gumChoose(header string, options []string) string {
	cmd := exec.Command("gum", "choose", "--header", header)
	cmd.Stderr = os.Stderr

	stdin, err := cmd.StdinPipe()
	if err != nil {
		fmt.Println("Error creating stdin pipe:", err)
		os.Exit(1)
	}

	go func() {
		defer stdin.Close()
		for _, opt := range options {
			fmt.Fprintln(stdin, opt)
		}
	}()

	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Selection cancelled.")
		os.Exit(0)
	}
	return string(output[:len(output)-1])
}
