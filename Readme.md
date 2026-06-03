# kommit

> AI-powered git commit messages that don't suck.

![CI](https://github.com/mujib77/kommit/actions/workflows/ci.yml/badge.svg)
![Version](https://img.shields.io/badge/version-v0.2.0-cyan)
![License](https://img.shields.io/badge/license-MIT-blue)

---

## The problem

You staged 4 files. Two are auth changes. One is a README fix. One is a config update. You type `git commit -m "fix stuff"` and move on.

Your git history is now useless.

## The solution

```bash
kommit
```

kommit reads your diff, detects unrelated changes, splits them into logical groups, and generates meaningful commit messages for each one.

---

## Features

**Atomic Commit Splitter**

Detects unrelated changes and splits them into separate logical commits automatically. No other AI commit tool does this.

```
◆ KOMMIT — analyzing your changes...
  staged: 4 files  +47  -12

◆ detected 2 logical groups

  Group 1 — 2 files
    internal/auth/jwt.go
    internal/auth/middleware.go

  1. feat(auth): add JWT token validation
  2. feat(auth): implement middleware authentication  
  3. chore(auth): add auth layer

  [1/2/3] pick  [e] edit  [q] quit

  ✓ committed group 1: feat(auth): add JWT token validation

  Group 2 — 1 file
    README.md

  1. docs: update README with setup instructions
  2. docs: add installation guide
  3. chore: update documentation

  [1/2/3] pick  [e] edit  [q] quit

  ✓ committed group 2: docs: update README with setup instructions
```

**3 message options per commit**

Pick the one that fits, edit it, or regenerate.

**4 AI providers**

Gemini (free), OpenAI, Anthropic, Ollama (local, free).

**Conventional Commits**

`feat:` `fix:` `docs:` `chore:` `refactor:` `test:` out of the box.

---

## Install

```bash
go install github.com/mujib77/kommit@latest
```

---

## Setup

Create `~/.kommit/config.yaml`:

```yaml
provider: gemini        # gemini | openai | anthropic | ollama
api_key: your-key-here
model: gemini-2.5-flash
style: conventional
```

**Free options:**

| Provider | Free | Key needed | Get key |
|----------|------|------------|---------|
| Gemini | ✅ | ✅ | [aistudio.google.com](https://aistudio.google.com) |
| Ollama | ✅ | ❌ | [ollama.com](https://ollama.com) |
| OpenAI | ❌ | ✅ | [platform.openai.com](https://platform.openai.com) |
| Anthropic | ❌ | ✅ | [console.anthropic.com](https://console.anthropic.com) |

**Ollama setup (completely free, runs locally):**

```bash
ollama pull llama3.2
ollama serve
```

Then set `provider: ollama` in config — no API key needed.

---

## Usage

```bash
# stage your changes
git add .

# run kommit
kommit
```

---

## Roadmap

```
v0.1.0  ✅  3 message options, 4 providers, interactive TUI
v0.2.0  ✅  Atomic commit splitter
v0.3.0  →   Repo style learner — learns your team's commit style
v0.4.0  →   PR description generator
v1.0.0  →   brew install, binary releases for all platforms
```

---

## Built With

- [Go](https://golang.org)
- [Cobra](https://github.com/spf13/cobra) — CLI framework
- [Bubbletea](https://github.com/charmbracelet/bubbletea) — TUI
- [Lipgloss](https://github.com/charmbracelet/lipgloss) — terminal styling

---

## License

MIT