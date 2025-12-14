package ai

import (
	"fmt"
	"os"
)

func GeneratePrompt(req Request) string {
	switch req.Task {
	case "tag":
		return generateTagPrompt(req)
	case "commit":
		return generateCommitPrompt(req)
	default:
		fmt.Printf("The prompt type '%s' isn't implemented yet!\n", req.Task)
		os.Exit(1)
	}

	return ""
}

func buildCommitSystemPrompt(req Request) string {
	rules := `You are a git commit message generator.
Follow the Conventional Commits specification.

Format:
<type>[optional scope]: <description>

[optional body]

Rules:
- Types: fix, feat, build, chore, ci, docs, style, refactor, perf, test.
- Use present tense.
- Max title length: 50 chars.
- No markdown code blocks.
- No conversational text.`

	if req.IsBrief {
		rules += "\nConstraint: I prefer a very short, one-sentence summary."
	}

	if req.UseEmoji {
		rules += "\nConstraint: Use GitMojis (e.g. 🐛 fix:)."
	}

	return rules
}

func buildCommitDiffPayload(req Request) string {
	return fmt.Sprintf(`
	Files changed:
	%s

	Diff:
	`+"```diff"+`
	%s
	`+"```", req.StagedFiles, req.Diff)
}

func buildCommitTrigger() string {
	return "Based on the diff above, generate the commit message now. Output raw text only."
}

func generateCommitPrompt(req Request) string {
	rules := buildCommitSystemPrompt(req)
	data := buildCommitDiffPayload(req)
	trigger := buildCommitTrigger()

	return fmt.Sprintf("%s\n\n---\n\n%s\n\n---\n\n%s", rules, data, trigger)
}

func buildTagSystemPrompt(req Request) string {
	rules := `
	You are a release notes generator for a git repository.
Your task is to create a creative title and a summary for a new release based on a list of commit messages.

The output should be in two parts:
1. A creative and concise title for the release on the first line.
2. The body of the release notes, starting from the third line (leave a blank line after the title).

Rules for the body:
- Group changes by their type. Use titles for sections, like "Features:", "Bug Fixes:", "Miscellaneous:".
- Do not use markdown headings.
- IMPORTANT: Do not start any lines with the hash symbol (#).
- Each item in the list should be a brief, clear summary of the change.
- Omit the commit hashes from the output.
- The tone should be professional and user-friendly.
- Do not include conversational text or any text outside of the release notes themselves.`

	if req.UseEmoji {
		rules += "\nConstraint: Use fun emojis for section headers."
	}

	return rules
}

func buildTagDiffPayload(req Request) string {
	return fmt.Sprintf("Commits since last release:\n%s", req.Diff)
}

func buildTagTrigger() string {
	return "Generate the release notes now."
}

func generateTagPrompt(req Request) string {
	rules := buildTagSystemPrompt(req)
	data := buildTagDiffPayload(req)
	trigger := buildTagTrigger()

	return fmt.Sprintf("%s\n\n---\n\n%s\n\n---\n\n%s", rules, data, trigger)
}
