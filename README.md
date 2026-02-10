# Text to Speech (Go)

![CI](https://github.com/justSteven-lang/text-to-speech/actions/workflows/ci.yml/badge.svg)

A simple **Text-to-Speech (TTS) CLI application** written in Go, using `espeak` to generate audio output.

This project is intentionally built as a **learning playground for CI lifecycle and junior DevOps practices**, including testing, linting, coverage gates, Docker, and GitHub Actions.

---

## ✨ Features

- Convert text input into speech
- Generate audio output as `.wav` file
- Simple CLI interface
- CI pipeline with:
  - Linting (`golangci-lint`)
  - Unit tests
  - Coverage threshold
  - Docker build

---

## 📦 Requirements

### Local Development

- Go **>= 1.25**
- `espeak`

### Ubuntu / Debian

```bash
sudo apt update
sudo apt install -y espeak
```

---

## 🚀 Usage

### Run locally

```bash
go run main.go "Hello from Text to Speech"
```

Output:

```
output.wav
```

---

## 🧪 Testing

Run tests with coverage:

```bash
go test ./... -cover
```

The CI pipeline enforces a **minimum test coverage threshold** to maintain code quality.

---

## 🔁 CI Pipeline Overview

The GitHub Actions CI pipeline performs the following steps:

1. Set up Go environment
2. Install OS dependencies
3. Run linter
4. Run unit tests with coverage check
5. Build Docker image

---

## 🎯 Learning Goals

This repository is designed to practice:

- CI/CD fundamentals
- GitHub Actions workflows
- Go project structure
- DevOps mindset (quality gates, automation, reproducibility)

---

## 📄 License

MIT
