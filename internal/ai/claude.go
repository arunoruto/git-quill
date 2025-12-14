package ai

import (
	"os/exec"
	"strings"
)

type Claude struct{}

func (p Claude) Name() string {
	return "Claude"
}

func (p Claude) IsAvailable() bool {
	_, err := exec.LookPath("claude")
	return err == nil
}

func (p Claude) ListModels() ([]string, error) {
	return []string{
		"(default)",
		"opus",
		"haiku",
		"sonnet",
	}, nil
}

func (p Claude) Generate(req Request) (string, error) {
	fullPrompt := GeneratePrompt(req)

	args := []string{"-p"}

	if req.Model != "" && req.Model != "(default)" {
		args = append(args, "--model", req.Model)
	}

	cmd := exec.Command("claude", args...)
	cmd.Stdin = strings.NewReader(fullPrompt)

	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(out), nil

}
