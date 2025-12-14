package ai

import (
	"os/exec"
	"strings"
)

type Opencode struct{}

func (p Opencode) Name() string {
	return "Opencode"
}

func (p Opencode) IsAvailable() bool {
	_, err := exec.LookPath("opencode")
	return err == nil
}

func (p Opencode) ListModels() ([]string, error) {
	cmd := exec.Command("opencode", "models")
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	var models []string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) > 0 {
			models = append(models, fields[0])
		}
	}

	return models, nil
}

func (p Opencode) Generate(req Request) (string, error) {
	fullPrompt := GeneratePrompt(req)

	model := req.Model
	if model == "" {
		model = "ollama/gemma3:4b"
	}

	cmd := exec.Command("opencode", "run", model)
	cmd.Stdin = strings.NewReader(fullPrompt)

	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(out), nil
}
