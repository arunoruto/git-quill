package ai

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func GenerateMessage(provider, model, stagedFiles, diff string) (string, error) {
	rules := `
		You are a git commit message generator.
		Follow the Conventional Commits specification.
		Rules:
		- Types: fix, feat, build, chore, ci, docs, style, refactor, perf, test.
		- Use present tense.
		- Max title length: 50 chars.
		- No markdown code blocks.
		- No conversational text.`

	data := fmt.Sprintf("Files changed:\n%s\n\nDiff:\n```diff\n%s\n```", stagedFiles, diff)

	trigger := `Based on the diff above, generate the commit message now.
	Output raw text only.`

	var cmd *exec.Cmd
	switch provider {
	case "ollama":
		fullPrompt := fmt.Sprintf(`
			%s

			---

			%s

			---

			%s
		`, rules, data, trigger)

		if model == "" {
			model = "gemma3:4b"
		}
		cmd = exec.Command("ollama", "run", model)
		cmd.Stdin = strings.NewReader(fullPrompt)
	case "opencode":
		args := []string{"run", "--agent", "commit"}
		if model != "" {
			args = append(args, "-m", model)
		}
		prompt := fmt.Sprintf("%s\n%s", data, trigger)
		args = append(args, prompt)

		cmd = exec.Command("opencode", args...)
		cmd.Stderr = os.Stderr
	default:
		return "", fmt.Errorf("provider %s not implemented yet", provider)
	}

	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return cleanOutput(string(output)), nil
}

func cleanOutput(raw string) string {
	s := strings.TrimSpace(raw)
	s = strings.ReplaceAll(s, "```", "")
	return s
}
