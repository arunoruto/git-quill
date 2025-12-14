package ai

import (
	"os/exec"
	"strings"
)

type Gemini struct{}

func (p Gemini) Name() string {
	return "Gemini"
}

func (p Gemini) IsAvailable() bool {
	_, err := exec.LookPath("gemini")
	return err == nil
}

func (p Gemini) ListModels() ([]string, error) {
	return []string{
		"(default)",
		"gemini-2.5-pro",
		"gemini-2.5-flash",
		"gemini-2.5-flash-lite",
	}, nil
}

func (p Gemini) Generate(req Request) (string, error) {
	fullPrompt := GeneratePrompt(req)

	var args []string

	if req.Model != "" && req.Model != "(default)" {
		args = append(args, "--model", req.Model)
	}

	cmd := exec.Command("gemini", args...)
	cmd.Stdin = strings.NewReader(fullPrompt)

	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(out), nil

}
