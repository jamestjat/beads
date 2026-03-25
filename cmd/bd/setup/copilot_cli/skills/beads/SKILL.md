---
name: beads
description: >-
  Beads issue tracking workflow. Use this skill when working with bd commands
  for issue tracking, creating issues, claiming work, closing tasks, managing
  dependencies, and syncing with Dolt remotes.
allowed-tools: "shell(bd:*)"
---

# Beads Issue Tracking

**beads (bd)** is a Dolt-powered issue tracker for AI-supervised coding workflows.

## Core Workflow

1. **Find work**: `bd ready --json` — list unblocked issues
2. **Claim**: `bd update <id> --claim --json` — atomically assign and start
3. **Work**: Implement, test, document
4. **Discover**: `bd create "Title" --description="Context" -t bug -p 1 --deps discovered-from:<parent-id> --json`
5. **Complete**: `bd close <id> --reason "Done" --json`
6. **Sync**: `bd dolt push`

## Commands

```
bd ready              # Unblocked issues
bd list --status open # Filter issues
bd create "Title" --description="..." -t task -p 2 --json
bd show <id> --json   # Full details
bd update <id> --claim --json
bd close <id> --reason "Done" --json
bd dep add <id> blocks:<other>
bd dep tree <id>
bd stale --days 30    # Forgotten issues
bd dolt push          # Push to remote
bd dolt pull          # Pull from remote
bd prime              # Full workflow context
```

## Issue Types

- `bug` — Something broken
- `feature` — New functionality
- `task` — Work item (tests, docs, refactoring)
- `epic` — Large feature with subtasks
- `chore` — Maintenance

## Priorities

- `0` Critical (security, data loss)
- `1` High (major features, important bugs)
- `2` Medium (default)
- `3` Low (polish)
- `4` Backlog (future ideas)

## Rules

- Always use `--json` for programmatic output
- Always include `--description` when creating issues
- Link discovered work with `discovered-from` dependencies
- Use `bd dolt push` to sync at end of session

## Auto-Discovery

While working on any task, if you encounter:
- **Bugs**: Broken behavior, error handling gaps, incorrect logic
- **Improvements**: Code smells, missing tests, documentation gaps, performance issues
- **Tech debt**: Deprecated patterns, hardcoded values, missing validation

**Immediately** create a beads issue linked to the current task:

```bash
bd create "<clear title>" --description="Found while working on <current-task>: <details>" -t bug -p <priority> --deps discovered-from:<current-id> --json
```

Do not stop your current task — file the issue and continue working. This ensures nothing is forgotten.
