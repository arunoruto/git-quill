package ai

import "fmt"

func buildSystemPrompt(req Request) string {
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

func buildDiffPayload(req Request) string {
	return fmt.Sprintf(`
	Files changed:
	%s

	Diff:
	`+"```diff"+`
	%s
	`+"```", req.StagedFiles, req.Diff)
}

func buildTrigger() string {
	return "Based on the diff above, generate the commit message now. Output raw text only."
}

func BuildPrompt(req Request) string {
	rules := buildSystemPrompt(req)
	data := buildDiffPayload(req)
	trigger := buildTrigger()

	return fmt.Sprintf("%s\n\n---\n\n%s\n\n---\n\n%s", rules, data, trigger)
}
