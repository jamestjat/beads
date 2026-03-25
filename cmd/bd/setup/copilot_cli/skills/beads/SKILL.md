---
name: beads
description: >-
  Beads issue tracking workflow using the bd CLI. Use this skill whenever the
  user asks about issues, tasks, bugs, work items, what to work on, tracking
  progress, dependencies between tasks, or syncing work — even if they don't
  mention "beads" or "bd" by name. Also use when the user says things like
  "create a ticket", "what's left to do", "mark that done", "what blocks X",
  or "push my changes". This project uses bd for ALL issue tracking — never
  create markdown TODOs or task lists.
allowed-tools: "shell(bd:*)"
---

# Beads Issue Tracking

**beads (bd)** is a Dolt-powered issue tracker built for AI coding workflows. It replaces markdown TODOs with a dependency-aware task graph that syncs via Dolt version control.

The `--json` flag matters because it gives you structured output you can parse reliably instead of human-formatted tables. Use it on every read/write command (not needed for `bd dolt push/pull` or `bd dep add` which don't return structured data).

## Core Workflow

1. **Find work**: `bd ready --json` — shows unblocked issues sorted by priority. Pick the highest-priority item; break ties by oldest first.
2. **Claim**: `bd update <id> --claim --json` — atomically sets assignee + status in one step, preventing double-assignment
3. **Work**: Implement, test, document
4. **Discover**: File anything you find along the way (bugs, tech debt, ideas) — link it back to your current task so context isn't lost:
   ```bash
   bd create "Title" --description="Found while working on <current-task>: <details>" -t bug -p 1 --deps discovered-from:<parent-id> --json
   ```
5. **Complete**: `bd close <id> --reason "what was done" --json` — use the reason to summarize the actual fix, not just "Done"
6. **Sync**: `bd dolt push` (runs automatically via sessionEnd hook)

## Commands Reference

| Command | What it does | Example use |
|---------|-------------|-------------|
| `bd ready --json` | Unblocked issues | "what can I work on?" |
| `bd list --status open --json` | Filter issues | "show all open bugs" |
| `bd create "Title" -t task -p 2 --description="..." --json` | New issue | "create a ticket for X" |
| `bd show <id> --json` | Full details | "tell me about bd-42" |
| `bd update <id> --claim --json` | Claim work | "I'll take bd-42" |
| `bd close <id> --reason "Fixed X" --json` | Complete | "mark bd-42 done" |
| `bd dep add <id> blocks:<other>` | Add dependency | "bd-99 blocks bd-42" |
| `bd dep tree <id>` | Dependency tree | "what depends on bd-42?" |
| `bd stale --days 30 --json` | Forgotten issues | "anything forgotten?" |
| `bd prime` | Full context dump | "reload beads context" |
| `bd dolt push` / `pull` | Sync to remote | "push my changes" |

## Issue Types & Priorities

**Types**: `bug` (broken), `feature` (new), `task` (work item), `epic` (large w/ subtasks), `chore` (maintenance)

**Priorities**: `0` critical (security, data loss) · `1` high (bugs, major features) · `2` medium (default) · `3` low (polish) · `4` backlog

When in doubt: bugs default to `1`, features to `2`, chores to `3`.

## Parallel Work with /fleet

When multiple independent issues are ready, use `/fleet` to work them simultaneously:

1. `bd ready --json` to find unblocked work
2. `bd dep tree` to confirm issues don't block each other
3. `/fleet` to distribute one issue per subagent
4. Each subagent claims → works → closes its issue
5. **Only the orchestrator runs `bd dolt push`** — subagents writing to Dolt concurrently can cause conflicts

## Session Management

- `/new` — Fresh conversation (triggers sessionEnd hook, auto-pushes)
- `/clear` — Abandon session (no push, no hooks — use to discard unwanted work)
- `/undo` — Revert last turn including both code and bd state changes
- `/compact` — Compress context; `bd prime` re-injects beads context afterward
- `/session plan` — Review current plan

## Proactive Discovery

While working on any task, if you notice bugs, missing tests, code smells, or improvement opportunities — file them immediately as linked issues. This is the single most valuable habit in the beads workflow because it captures context at the moment you have it, rather than hoping someone remembers later.

```bash
bd create "<clear title>" --description="Found while working on <current-task>: <details>" -t bug -p <priority> --deps discovered-from:<current-id> --json
```

Don't stop your current task to investigate — just file and keep going.

## Key Conventions

- **`--json` on read/write commands** — structured output is more reliable to parse (not needed for push/pull/dep add)
- **`--description` on every create** — issues without descriptions lose context and become useless
- **`discovered-from` links** — connect found work to the task where you found it
- **Pre-approve bd**: `/allow-all` or `copilot --allow-tool='shell(bd:*)'`
