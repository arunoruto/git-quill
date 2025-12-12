package git

import (
	"fmt"
	"os/exec"
	"strings"
)

func GetStagedDiff(limit int) (string, error) {
	ignoredPatterns := []string{
		"*.lock",
		"**/package-lock.json",
		"**/yarn.lock",
		"**/pnpm-lock.yaml",
		"**/go.sum",
		"**/devenv.lock",
		"**/devenv.yaml",
		"*.svg",
		"*.min.js",
		"*.map",
	}

	args := []string{"diff", "--cached", "."}
	for _, pattern := range ignoredPatterns {
		args = append(args, fmt.Sprintf(":(exclude)%s", pattern))
	}

	cmd := exec.Command("git", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("git error: %v", err)
	}

	diff := string(output)
	if strings.TrimSpace(diff) == "" {
		return "", fmt.Errorf("no staged changes found (or all ignored)")
	}

	if limit > 0 && len(diff) > limit {
		diff = diff[:limit] + "\n... [Diff Truncated]"
	}

	return diff, nil
}

func GetStagedFiles() (string, error) {
	cmd := exec.Command("git", "diff", "--cached", "--name-status")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(output)), nil
}
