package ai

import (
	"os/exec"
	"strings"
)

type Ollama struct{}

func (p Ollama) Name() string {
	return "ollama"
}

func (p Ollama) IsAvailable() bool {
	_, err := exec.LookPath("ollama")
	return err == nil
}

func (p Ollama) ListModels() ([]string, error) {
	cmd := exec.Command("ollama", "list")
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	var models []string

	for i, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if i == 0 && strings.HasPrefix(line, "NAME") {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) > 0 {
			models = append(models, fields[0])
		}
	}

	return models, nil
}

func (p Ollama) Generate(req Request) (string, error) {
	fullPrompt := GeneratePrompt(req)

	model := req.Model
	if model == "" {
		model = "gemma3:4b"
	}

	cmd := exec.Command("ollama", "run", model)
	cmd.Stdin = strings.NewReader(fullPrompt)

	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(out), nil

}
