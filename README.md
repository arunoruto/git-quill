# git-quill

`git-quill` is a command-line tool that uses AI to help you write commit messages and tag messages. It's designed to be a simple and interactive way to improve your git workflow.

## Features

- **AI-powered commit messages:** Generate commit messages from your staged changes.
- **AI-powered tag messages:** Generate release notes from the commits since your last tag.
- **Multiple AI providers:** Supports Ollama, Opencode, Gemini, and Copilot.
- **Interactive UI:** Select your preferred AI provider and model through a simple and interactive UI.
- **Brief and detailed summaries:** Choose between a brief or detailed summary for your commit messages.

## Installation

1.  Make sure you have Go installed on your system.
2.  Clone the repository:
    ```bash
    git clone https://github.com/arunoruto/git-quill.git
    ```
3.  Navigate to the project directory:
    ```bash
    cd git-quill
    ```
4.  Build the project:
    ```bash
    go build -o git-quill ./cmd/git-quill
    ```
5.  Move the binary to a directory in your `PATH` (e.g., `/usr/local/bin`):
    ```bash
    mv git-quill /usr/local/bin/
    ```

## Usage

### `commit`

Generate a commit message from your staged changes:

```bash
git-quill commit
```

You can also specify the AI provider and model:

```bash
git-quill commit -p ollama -m gemma3:4b
```

For a brief summary, use the `-b` flag:

```bash
git-quill commit -b
```

### `tag`

Generate a tag message (release notes) from the commits since your last tag:

```bash
git-quill tag
```

You can also specify the AI provider and model:

```bash
git-quill tag -p ollama -m gemma3:4b
```

## Supported AI Providers

-   [Ollama](https://ollama.ai/)
-   [Opencode](https://www.opencode.ai/)
-   [Gemini](https://gemini.google.com/)
-   [Copilot](https://github.com/features/copilot)

## License

This project is licensed under the MIT License.
