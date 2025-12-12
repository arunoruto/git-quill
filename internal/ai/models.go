package ai

import (
	"os/exec"
	"strings"
)

func ListModels(provider string) ([]string, error) {
	var cmd *exec.Cmd

	switch provider {
	case "ollama":
		cmd = exec.Command("ollama", "list")
	case "opencode":
		cmd = exec.Command("opencode", "models")
	default:
		return nil, nil
	}

	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	return parseModelOutput(provider, string(output)), nil
}

func parseModelOutput(provider, raw string) []string {
	lines := strings.Split(strings.TrimSpace(raw), "\n")
	var models []string

	for i, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		switch provider {
		case "ollama":
			if i == 0 && strings.HasPrefix(line, "NAME") {
				continue
			}
			fields := strings.Fields(line)
			if len(fields) > 0 {
				models = append(models, fields[0])
			}
		default:
			models = append(models, line)
		}
	}

	return models
}
