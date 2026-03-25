# GitHub Copilot Integration Guide

This guide explains how to use beads with GitHub Copilot in VS Code.

## Overview

Beads integrates with GitHub Copilot through the `bd` CLI. Copilot Chat in agent mode has built-in terminal access — it runs `bd` commands directly and parses the JSON output. No external MCP server or Python dependency is required.

**Important:** Beads is a system-wide CLI tool. You install it once and use it in any project. Do NOT clone the beads repository into your project.

## Prerequisites

- VS Code 1.96+ with GitHub Copilot extension
- GitHub Copilot subscription (Individual, Business, or Enterprise)
- beads CLI installed (`brew install beads` or `npm install -g @beads/bd`)

## Quick Setup

### Step 1: Install Copilot integration

```bash
bd setup copilot
```

This creates `.github/copilot-instructions.md` with beads workflow context that Copilot reads automatically.

To also install reusable prompt files (optional):

```bash
bd setup copilot --prompts
```

This copies `.prompt.md` files to `.github/prompts/` for slash-command-style shortcuts in Copilot Chat.

**For Copilot CLI users** (recommended):

```bash
bd setup copilot --cli
```

This installs Copilot CLI-specific features:
- **Hooks** (`.github/hooks/beads.json`) — Auto-primes beads context on session start, auto-pushes on session end
- **Skill** (`.github/skills/beads/SKILL.md`) — Beads workflow skill with auto-approved `bd` commands
- **Agent** (`.github/agents/beads.agent.md`) — Dedicated beads issue tracking agent

You can combine flags: `bd setup copilot --cli --prompts`

### Step 2: Initialize beads in your project

```bash
cd your-project
bd init --quiet
```

This creates a `.beads/` directory with the issue database. The init wizard will ask about git hooks—these are optional and you can skip them if unfamiliar.

### Step 3: Restart VS Code

Reload the VS Code window for the instructions to take effect.

## Using Beads with Copilot

### Natural Language Commands

With the integration installed, ask Copilot Chat:

| You say | Copilot runs |
|---------|-------------|
| "What issues are ready to work on?" | `bd ready --json` |
| "Create a bug for the login timeout" | `bd create "Login timeout" -t bug --json` |
| "Show me issue bd-42" | `bd show bd-42 --json` |
| "Mark bd-42 as complete" | `bd close bd-42 --json` |
| "Claim bd-42" | `bd update bd-42 --claim --json` |

### CLI Reference

| Command | Description | Example |
|---------|-------------|---------|
| `bd ready` | List unblocked issues | "What can I work on?" |
| `bd list` | List issues with filters | "Show all open bugs" |
| `bd create` | Create new issue | "Create a task for refactoring" |
| `bd show` | Show issue details | "Show bd-42 details" |
| `bd update` | Update issue fields | "Set bd-42 to in progress" |
| `bd close` | Close an issue | "Complete bd-42" |
| `bd dolt push` | Push changes to remote | "Push my changes" |
| `bd dep add` | Add dependency | "bd-99 blocks bd-42" |
| `bd dep tree` | Show dependency tree | "What depends on bd-42?" |

### Example Workflow

```
You: What issues are ready to work on?

Copilot: [runs bd ready --json]
There are 3 issues ready:
1. [P1] bd-42: Fix authentication timeout
2. [P2] bd-99: Add password reset flow
3. [P3] bd-17: Update API documentation

You: Let me work on bd-42. Claim it.

Copilot: [runs bd update bd-42 --claim --json]
Claimed bd-42 and started work.

You: [... work on the code ...]

You: I found a related bug - the session token isn't being refreshed.
     Create a bug for that, linked to bd-42.

Copilot: [runs bd create "Session token not refreshed" -t bug -p 1 --deps discovered-from:bd-42 --json]
Created bd-103: Session token not refreshed
Linked as discovered-from bd-42.

You: Done with bd-42. Close it with reason "Fixed timeout handling"

Copilot: [runs bd close bd-42 --reason "Fixed timeout handling" --json]
Closed bd-42: Fixed timeout handling

You: Push my changes to the remote

Copilot: [runs bd dolt push]
Pushed: 2 issues updated, synced to Dolt remote.
```

## Reusable Prompts

When installed with `--prompts`, these files appear in `.github/prompts/`:

| Prompt file | Mode | What it does |
|-------------|------|-------------|
| `beads-ready.prompt.md` | agent | Find and claim ready work via terminal |
| `beads-create.prompt.md` | agent | Create issues interactively via terminal |
| `beads-workflow.prompt.md` | ask | Show the full beads workflow guide |
| `plan-to-beads.prompt.md` | agent | Convert a plan document into an epic + tasks |

Invoke them in Copilot Chat by selecting from the prompt picker or by referencing the file.

## Verify Installation

```bash
bd setup copilot --check
```

## Remove Integration

```bash
bd setup copilot --remove                    # Remove instructions only
bd setup copilot --remove --prompts          # Also remove prompt files
bd setup copilot --remove --cli              # Also remove hooks, skills, agents
bd setup copilot --remove --prompts --cli    # Remove everything
```

## Troubleshooting

### Copilot doesn't understand bd commands

1. **Check integration is installed** — `bd setup copilot --check`
2. **Reload VS Code** — Instructions require window reload
3. **Check bd is on PATH:**
   ```bash
   which bd    # macOS/Linux
   where bd    # Windows
   ```

### "bd: command not found"

Install beads:

```bash
# macOS
brew install beads

# npm (any platform)
npm install -g @beads/bd

# Or download from GitHub releases
```

### Changes not persisting

Push changes at end of session:

```bash
bd dolt push
```

Or ask Copilot: "Push my beads changes to the remote"

### No beads database found

Initialize beads in your project:

```bash
cd your-project
bd init --quiet
```

## FAQ

### Do I need to clone the beads repository?

**No.** Beads is a system-wide CLI tool. You install it once (via Homebrew, npm, or pip) and use it in any project. The `.beads/` directory in your project only contains the issue database, not beads itself.

### What are the git hooks and are they safe?

When you run `bd init`, beads can install git hooks that:
- **post-merge**: Import issues when you pull
- **pre-push**: Sync issues before you push

These hooks are safe—they only read/write the `.beads/` directory and never modify your code. You can opt out with `bd init --no-hooks` or skip them during the interactive setup.

### Can I use beads without Copilot?

Yes! Beads works with:
- Terminal (direct CLI)
- Claude Code
- Cursor
- Aider
- Any editor with shell access

### Does this work with Copilot in other editors?

This guide is for VS Code. For other editors:
- **JetBrains IDEs**: Use CLI integration with terminal
- **Neovim**: Use CLI integration directly

## Using with GitHub Copilot CLI

GitHub Copilot CLI (`copilot`) is a terminal-native AI coding assistant that reads the same instruction files as VS Code Copilot Chat. The beads integration works out of the box with Copilot CLI.

### Instruction Discovery

Copilot CLI automatically reads these files (in order):

| Location | Scope |
|----------|-------|
| `~/.copilot/copilot-instructions.md` | Global (all sessions) |
| `.github/copilot-instructions.md` | Repository |
| `.github/instructions/**/*.instructions.md` | Repository (modular) |
| `AGENTS.md` | Repository |

`bd setup copilot` writes to `.github/copilot-instructions.md`, which Copilot CLI reads automatically. If your repo also has an `AGENTS.md` (created by `bd init`), Copilot CLI reads that too — no extra configuration needed.

### Tool Permissions

Pre-approve `bd` commands so Copilot CLI doesn't prompt each time:

```bash
copilot --allow-tool='shell(bd:*)'
```

Common permission patterns for beads workflows:

```bash
# Allow all bd commands
copilot --allow-tool='shell(bd:*)'

# Allow bd and git commands
copilot --allow-tool='shell(bd:*)' --allow-tool='shell(git:*)'

# Allow bd but block destructive git operations
copilot --allow-tool='shell(bd:*)' --allow-tool='shell(git:*)' --deny-tool='shell(git push:*)'
```

### Plan Mode

Use Copilot CLI's `/plan` command for complex tasks, then track with beads:

```
/plan Implement OAuth2 authentication with Google and GitHub providers
```

After Copilot generates a plan, convert it to tracked issues:

```
Convert this plan into a beads epic with tasks
```

This works especially well with the `plan-to-beads.prompt.md` prompt file (installed with `--prompts`).

### Delegate

Use `/delegate` to offload tangential tasks to Copilot coding agent while you keep working:

```
/delegate Update the API documentation for the new endpoints
```

Good candidates for delegation:
- Documentation updates
- Refactoring separate modules
- Creating test scaffolding

Keep core issue work (claim, close, dependency changes) in your local session.

### Multi-Repository Workflows

For projects spanning multiple repos, start Copilot CLI from a parent directory:

```bash
cd ~/projects
copilot
```

Or add directories mid-session:

```
/add-dir /path/to/backend-service
/add-dir /path/to/frontend
```

Each repo can have its own `.beads/` database. Use `bd` commands from the appropriate directory.

### Recommended Workflow

The Copilot CLI best practice workflow maps naturally to beads:

| Copilot CLI Step | Beads Command |
|-----------------|---------------|
| **Explore** — understand codebase | `bd ready --json` — find available work |
| **Plan** — `/plan` the implementation | `plan-to-beads` — convert plan to tracked issues |
| **Implement** — write code | `bd update <id> --claim` — claim and track |
| **Verify** — run tests | Quality gates before closing |
| **Commit** — commit changes | `bd close <id>` + `bd dolt push` |

## See Also

- [Copilot CLI Best Practices](https://docs.github.com/en/copilot/how-tos/copilot-cli/cli-best-practices)
- [CLI Reference](QUICKSTART.md)
- [Installation Guide](INSTALLING.md)
- [Agent Instructions](../AGENT_INSTRUCTIONS.md)
