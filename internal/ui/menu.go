package ui

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func GumChoose(header string, options []string) string {
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

func GumFilter(placeholder string, options []string) string {
	if len(options) == 0 {
		return ""
	}

	cmd := exec.Command("gum", "filter", "--placeholder", placeholder)
	cmd.Stderr = os.Stderr
	stdin, _ := cmd.StdinPipe()

	go func() {
		defer stdin.Close()
		for _, opt := range options {
			fmt.Fprintln(stdin, opt)
		}
	}()

	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Selection canceled.")
		os.Exit(0)
	}

	return strings.TrimSpace(string(output))
}
