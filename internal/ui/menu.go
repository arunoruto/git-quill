package ui

import (
	// "fmt"
	// "os"
	// "os/exec"
	// "strings"
	"fmt"
	"os"

	"github.com/charmbracelet/huh"
)

func GumChoose(title string, options []string) string {
	var selected string

	opts := make([]huh.Option[string], len(options))
	for i, opt := range options {
		opts[i] = huh.NewOption(opt, opt)
	}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title(title).
				Options(opts...).
				Value(&selected),
		),
	)

	err := form.Run()
	if err != nil {
		fmt.Println("Operation cancelled.")
		os.Exit(0)
	}

	return selected
}

func GumFilter(title string, options []string) string {
	var selected []string

	opts := make([]huh.Option[string], len(options))
	for i, opt := range options {
		opts[i] = huh.NewOption(opt, opt)
	}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewMultiSelect[string]().
				Title(title).
				Options(opts...).
				Filterable(true).
				Height(10).
				Value(&selected),
		),
	)

	err := form.Run()
	if err != nil {
		fmt.Println("Operation cancelled.")
		os.Exit(0)
	}

	if len(selected) > 0 {
		return selected[0]
	}

	return ""
}

// func GumChoose(header string, options []string) string {
// 	cmd := exec.Command("gum", "choose", "--header", header)
// 	cmd.Stderr = os.Stderr

// 	stdin, err := cmd.StdinPipe()
// 	if err != nil {
// 		fmt.Println("Error creating stdin pipe:", err)
// 		os.Exit(1)
// 	}

// 	go func() {
// 		defer stdin.Close()
// 		for _, opt := range options {
// 			fmt.Fprintln(stdin, opt)
// 		}
// 	}()

// 	output, err := cmd.Output()
// 	if err != nil {
// 		fmt.Println("Selection cancelled.")
// 		os.Exit(0)
// 	}
// 	return string(output[:len(output)-1])
// }

// func GumFilter(placeholder string, options []string) string {
// 	if len(options) == 0 {
// 		return ""
// 	}

// 	cmd := exec.Command("gum", "filter", "--placeholder", placeholder)
// 	cmd.Stderr = os.Stderr
// 	stdin, _ := cmd.StdinPipe()

// 	go func() {
// 		defer stdin.Close()
// 		for _, opt := range options {
// 			fmt.Fprintln(stdin, opt)
// 		}
// 	}()

// 	output, err := cmd.Output()
// 	if err != nil {
// 		fmt.Println("Selection canceled.")
// 		os.Exit(0)
// 	}

// 	return strings.TrimSpace(string(output))
// }
