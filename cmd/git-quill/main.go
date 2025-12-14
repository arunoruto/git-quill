package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
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
		runTag(os.Args[2:])
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
	cmd.Parse(args)

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

	sageMsg := strings.ReplaceAll(msg, "\"", "\\\"")
	fmt.Println("\n" + msg)
	fmt.Println("\n------------------------------------------------")
	fmt.Println("Create this tag now:")
	fmt.Printf("git tag -a %s -m \"%s\"\n", tagName, sageMsg)
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
