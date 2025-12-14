package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
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
	mainCmd := flag.NewFlagSet("git-quill", flag.ExitOnError)

	var listProviders bool
	mainCmd.BoolVar(&listProviders, "list-providers", false, "List available AI providers")

	var listModelsProvider string
	mainCmd.StringVar(&listModelsProvider, "list-models", "", "List models for specific provider")

	var help bool
	mainCmd.BoolVar(&help, "h", false, "Show help")
	mainCmd.BoolVar(&help, "help", false, "Show help")

	if err := mainCmd.Parse(os.Args[1:]); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	if help {
		printMainUsage()
		os.Exit(0)
	} else if listProviders {
		printProviders()
		os.Exit(0)
	} else if listModelsProvider != "" {
		printModels(listModelsProvider)
		os.Exit(0)
	}

	args := mainCmd.Args()

	if len(args) < 1 {
		printMainUsage()
		os.Exit(1)
	}

	switch args[0] {
	case "commit":
		runCommit(args[1:])
	case "tag":
		runTag(args[1:])
	default:
		fmt.Printf("Unknown command: %s\n", args[0])
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

func printProviders() {
	providers := ai.GetAvailableProviders()
	for _, p := range providers {
		fmt.Println(strings.ToLower(p.Name()))
	}
}

func printModels(providerName string) {
	p, err := ai.GetProviderByName(providerName)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	models, err := p.ListModels()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	for _, m := range models {
		fmt.Println(m)
	}
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
	if err := cmd.Parse(args); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	setupLogger(config.Verbose)

	diff, err := git.GetStagedDiff(20000)
	if err != nil {
		fmt.Printf("Error reading git: %v\n", err)
		os.Exit(1)
	}

	files, _ := git.GetStagedFiles()

	provider, err := resolveAI(config)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	req := ai.Request{
		Task:        "commit",
		Diff:        diff,
		StagedFiles: files,
		Model:       config.Model,
		IsBrief:     brief,
	}

	slog.Debug("Running Commit", "provider", config.Provider, "model", config.Model, "brief", brief)
	slog.Debug("generating commit message...")
	msg, err := provider.Generate(req)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	fmt.Println(msg)
}

func runTag(args []string) {
	cmd := flag.NewFlagSet("tag", flag.ExitOnError)
	config := registerSharedFlags(cmd)

	var raw bool
	cmd.BoolVar(&raw, "raw", false, "Output raw markdown only (no instructions)")

	if err := cmd.Parse(args); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	setupLogger(config.Verbose)
	slog.Debug("Running Tag", "provider", config.Provider)

	tagName := "vX.Y.Z"
	if len(cmd.Args()) > 0 {
		tagName = cmd.Args()[0]
	}

	lastTag := git.GetLastTag()
	slog.Debug("Found last tag", "tag", lastTag)

	commits, err := git.GetCommitsSince(lastTag)
	if err != nil {
		fmt.Printf("Error reading git log: %v\n", err)
		os.Exit(1)
	}

	if strings.TrimSpace(commits) == "" {
		fmt.Printf("No new commits found since the last tag '%s'.\n", lastTag)
		os.Exit(0)
	}

	provider, err := resolveAI(config)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	req := ai.Request{
		Task:  "tag",
		Diff:  commits,
		Model: config.Model,
	}

	slog.Debug("Running tag", "provider", config.Provider, "model", config.Model)
	slog.Debug("generating release notes since ", lastTag, "...")
	msg, err := provider.Generate(req)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	if raw {
		fmt.Print(msg)
	} else {
		sageMsg := strings.ReplaceAll(msg, "\"", "\\\"")
		fmt.Println("\n" + msg)
		fmt.Println("\n------------------------------------------------")
		fmt.Println("Create this tag now:")
		fmt.Printf("git tag -a %s -m \"%s\"\n", tagName, sageMsg)
	}
}

func resolveAI(cfg *Config) (ai.Provider, error) {
	available := ai.GetAvailableProviders()
	if len(available) == 0 {
		return nil, fmt.Errorf("no AI providers found (install ollama, copilot, etc)")
	}

	var provider ai.Provider
	var err error

	if cfg.Provider == "" {
		names := make([]string, len(available))
		for i, p := range available {
			names[i] = p.Name()
		}
		choice := ui.Select("Select AI Provider", names)
		if choice == "" {
			return nil, fmt.Errorf("selection cancelled")
		}
		provider, _ = ai.GetProviderByName(choice)
		cfg.Provider = provider.Name()
	} else {
		provider, err = ai.GetProviderByName(cfg.Provider)
		if err != nil {
			return nil, err
		}
	}

	if cfg.Model == "" {
		models, err := provider.ListModels()
		if err == nil && len(models) > 0 {
			cfg.Model = ui.Select("Select Model", models)
			if cfg.Model == "" {
				return nil, fmt.Errorf("selection cancelled")
			}
		}
	}

	return provider, nil
}
