package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"slices"
	"strings"
	"time"

	"github.com/arunoruto/git-quill/internal/ai"
	"github.com/arunoruto/git-quill/internal/git"
	"github.com/arunoruto/git-quill/internal/ui"
	"github.com/lmittmann/tint"
)

type Config struct {
	Provider string
	Model    string
	Verbose  bool
}

func main() {
	if len(os.Args) < 2 {
		printMainUsage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "commit":
		runCommit(os.Args[2:])
	case "tag":
		fmt.Println("Tag generation feature coming soon!")
	case "-h", "--help", "help":
		printMainUsage()
	default:
		fmt.Printf("Unknown command: %s\n", os.Args[1])
		printMainUsage()
		os.Exit(1)
	}
}

func setupLogger(verbose bool) {
	level := slog.LevelInfo
	if verbose {
		level = slog.LevelDebug
	}

	// opts := &slog.HandlerOptions{Level: level}
	// logger := slog.New(slog.NewTextHandler(os.Stderr, opts))
	opts := &tint.Options{
		Level:      level,
		TimeFormat: time.DateTime,
		NoColor:    false,
	}
	logger := slog.New(tint.NewHandler(os.Stderr, opts))
	slog.SetDefault(logger)
}

func printMainUsage() {
	fmt.Println("Usage: git quill <command> [options]")
	fmt.Println("\nCommands:")
	fmt.Println("\t commit Generate an AI commit message")
	fmt.Println("\t tag    Generate an AI tag message")
}

func registerSharedFlags(fs *flag.FlagSet) *Config {
	cfg := &Config{}
	fs.StringVar(&cfg.Provider, "p", "", "AI Provider")
	fs.StringVar(&cfg.Model, "m", "", "Model Name")
	fs.BoolVar(&cfg.Verbose, "v", false, "Verbose logging")
	return cfg
}

func runCommit(args []string) {
	cmd := flag.NewFlagSet("commit", flag.ExitOnError)
	config := registerSharedFlags(cmd)
	var brief bool
	cmd.BoolVar(&brief, "b", false, "Brief summary")
	cmd.Parse(args)
	setupLogger(config.Verbose)

	diff, err := git.GetStagedDiff(20000)
	if err != nil {
		fmt.Printf("Error reading git: %v\n", err)
		os.Exit(1)
	}

	files, _ := git.GetStagedFiles()

	providers := findAvailableProviders()
	slog.Debug("Possible providers", strings.Join(providers, ", "), "")
	config.Provider = strings.ToLower(config.Provider)
	if config.Provider == "" {
		config.Provider = ui.GumChoose("Select AI Provider", providers)
	} else if !slices.Contains(providers, config.Provider) {
		fmt.Printf("The given provider %s is not in the list of providers (%s)", config.Provider, strings.Join(providers, ", "))
		os.Exit(0)
	}

	slog.Debug("Fetching models...", "provider", config.Provider)
	models, err := ai.ListModels(config.Provider)
	if err == nil || len(models) > 0 {
		config.Model = ui.GumFilter("Search Model...", models)
	} else {
		switch config.Provider {
		case "ollama":
			config.Model = "gemma3:4b"
		default:
			config.Model = "(default)"
		}
	}

	slog.Debug("Running Commit", "provider", config.Provider, "model", config.Model, "brief", brief)
	slog.Debug("generating commit message...")
	msg, err := ai.GenerateMessage(config.Provider, config.Model, files, diff)
	if err != nil {
		fmt.Printf("AI Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(msg)
}

func runTag(args []string) {
	cmd := flag.NewFlagSet("tag", flag.ExitOnError)
	config := registerSharedFlags(cmd)
	cmd.Parse(args)
	setupLogger(config.Verbose)
	slog.Debug("Running Tag", "provider", config.Provider)
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
