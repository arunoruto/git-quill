package ai

import (
	"os/exec"
	"strings"
)

type Apfel struct{}

func (p Apfel) Name() string {
	return "Apfel"
}

func (p Apfel) IsAvailable() bool {
	_, err := exec.LookPath("apfel")
	return err == nil
}

func (p Apfel) ListModels() ([]string, error) {
	return []string{"(default)"}, nil
}

func (p Apfel) Generate(req Request) (string, error) {
	fullPrompt := GeneratePrompt(req)

	cmd := exec.Command("apfel")
	cmd.Stdin = strings.NewReader(fullPrompt)

	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(out), nil
}
