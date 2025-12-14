package git

import (
	"fmt"
	"os/exec"
	"strings"
)

func GetLastTag() string {
	cmd := exec.Command("git", "describe", "--tag", "--abbrev=0")
	out, err := cmd.Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

func GetCommitsSince(tag string) (string, error) {
	rangeSpec := "HEAD"
	if tag != "" {
		rangeSpec = fmt.Sprintf("%s..HEAD", tag)
	}

	cmd := exec.Command("git", "log", "--pretty=format:- %h %s", rangeSpec)
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("git log error: %w", err)
	}

	return strings.TrimSpace(string(out)), nil
}
