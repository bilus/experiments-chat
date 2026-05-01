# Agent Environment Guide

Use this guide before running local commands.

## Devbox And Direnv

Use devbox and direnv for all project commands.

Preferred command shape:

```bash
direnv exec . go test ./...
```

Do not rely on a system Go toolchain being available outside devbox.

## Fresh Setup

If setting up the repo from scratch:

```bash
devbox init
devbox add go_1_25@latest
devbox generate direnv
direnv allow
```

The project may use a newer generated Go version in `go.mod`; always trust the repo and devbox state over global machine assumptions.

## Command Discipline

Use `direnv exec .` for Go commands, formatting, and tests.

Examples:

```bash
direnv exec . gofmt -w internal/chat/agent.go
direnv exec . go test ./...
```

For this project, the default verification command is:

```bash
direnv exec . go test ./...
```

