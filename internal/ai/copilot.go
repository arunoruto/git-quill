package ai

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"
)

type Copilot struct{}

func (p Copilot) Name() string {
	return "Copilot"
}

func (p Copilot) IsAvailable() bool {
	_, err := exec.LookPath("copilot")
	return err == nil
}

func (p Copilot) ListModels() ([]string, error) {
	cmd := exec.Command("copilot", "--help")
	outBytes, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to run copilot --help: %w", err)
	}
	out := string(outBytes)

	re := regexp.MustCompile(`(?s)--model.*?choices:\s*([^)]+)`)
	matches := re.FindStringSubmatch(out)

	if len(matches) < 2 {
		return nil, fmt.Errorf("could not parse models from help text")
	}

	rawList := matches[1]

	models := []string{"(default)"}
	parts := strings.Split(rawList, ",")

	for _, part := range parts {
		clean := strings.TrimSpace(part)
		clean = strings.Trim(clean, "\"")
		if clean != "" {
			models = append(models, clean)
		}
	}

	return models, nil

	// return []string{"(default)"}, nil
}

func (p Copilot) Generate(req Request) (string, error) {
	fullPrompt := BuildPrompt(req)

	args := []string{"-s"}

	if req.Model != "" && req.Model != "(default)" {
		args = append(args, "--model", req.Model)
	}

	cmd := exec.Command("copilot", args...)
	cmd.Stdin = strings.NewReader(fullPrompt)

	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(out), nil

}
