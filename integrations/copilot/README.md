# GitHub Copilot Integration for Beads

CLI-first integration between beads issue tracking and GitHub Copilot in VS Code.

## Overview

This integration provides:

- **Copilot instructions** — Workflow context injected into `.github/copilot-instructions.md` so Copilot understands the `bd` CLI
- **Reusable prompts** — `.prompt.md` files that work like slash commands in Copilot Chat, driving `bd` CLI commands via the built-in terminal

No external MCP server or Python dependency is required. Copilot uses its built-in terminal capability to run `bd` commands directly.

## Quick Setup

Run all commands from **your project directory** (not the beads repo):

```bash
cd your-project
bd init --quiet        # Initialize beads database
bd setup copilot       # Install Copilot integration
```

This creates / updates one file:

| File | Purpose |
|------|---------|
| `.github/copilot-instructions.md` | Teaches Copilot your bd workflow |

To also install reusable prompt files:

```bash
bd setup copilot --prompts
```

This additionally copies prompt files to `.github/prompts/`:

| Prompt file | What it does |
|-------------|-------------|
| `beads-ready.prompt.md` | Find and claim ready work |
| `beads-create.prompt.md` | Create a new issue interactively |
| `beads-workflow.prompt.md` | Show the full beads workflow guide |
| `plan-to-beads.prompt.md` | Convert a plan document into an epic + tasks |

Then reload VS Code for changes to take effect.

## Prerequisites

- VS Code 1.96+ with GitHub Copilot extension
- GitHub Copilot subscription
- beads CLI installed (`bd` on PATH)

## Verify

```bash
bd setup copilot --check
```

Ask Copilot Chat: *"What beads issues are ready to work on?"*

## Remove

```bash
bd setup copilot --remove             # Remove instructions only
bd setup copilot --remove --prompts   # Also remove prompt files
```

The beads section is surgically removed from `copilot-instructions.md`; any other content you've written there is preserved.

---

## Copilot CLI Support

This integration also works with [GitHub Copilot CLI](https://docs.github.com/en/copilot/how-tos/copilot-cli/cli-best-practices) (`copilot`), the terminal-native AI coding assistant.

Copilot CLI auto-reads `.github/copilot-instructions.md` and `AGENTS.md`, so the beads integration works out of the box.

**Pre-approve bd commands** to avoid per-command prompts:

```bash
copilot --allow-tool='shell(bd:*)'
```

**Use plan mode** for complex work, then convert to tracked issues:

```
/plan Implement the new authentication flow
# After plan is generated:
Convert this plan into a beads epic with tasks
```

See [COPILOT_INTEGRATION.md](../../docs/COPILOT_INTEGRATION.md) for the full Copilot CLI section.

---

## How It Works

Copilot Chat in agent mode has built-in terminal access. The instructions file teaches Copilot the `bd` CLI commands and workflow. When you ask Copilot to manage issues, it runs `bd` commands in the terminal with `--json` output and parses the structured results.

### Example Copilot Chat Session

```
You: What issues are ready to work on?
Copilot: [runs bd ready --json] Three issues ready: bd-42 (P1), bd-99 (P2), bd-17 (P3)

You: Claim bd-42 and start work.
Copilot: [runs bd update bd-42 --claim --json] Claimed bd-42: Fix authentication timeout

You: Found a related bug — session tokens aren't being refreshed.
Copilot: [runs bd create "Session token not refreshed" -t bug -p 1 --deps discovered-from:bd-42 --json]
         Created bd-103 (discovered-from bd-42)

You: Done! Close bd-42 with "Fixed timeout handling in auth middleware".
Copilot: [runs bd close bd-42 --reason "Fixed timeout handling in auth middleware" --json]
         Closed bd-42.

You: Push my changes.
Copilot: [runs bd dolt push] Pushed: 2 issues updated.
```

---

## How it Compares to Claude Code

| Feature | Claude Code | GitHub Copilot |
|---------|-------------|----------------|
| Setup command | `bd setup claude` | `bd setup copilot` / `bd setup copilot --cli` |
| Context file | `CLAUDE.md` | `.github/copilot-instructions.md` |
| Modular instructions | N/A | `.github/instructions/beads.instructions.md` (with `--cli`) |
| Session hooks | `bd prime` on SessionStart / PreCompact | `.github/hooks/beads.json` (with `--cli`) |
| Slash commands | `claude-plugin/commands/` | `.github/prompts/*.prompt.md` |
| Skills | `.claude/skills/` | `.github/skills/beads/SKILL.md` (with `--cli`) |
| Custom agents | N/A | `.github/agents/beads.agent.md` (with `--cli`) |
| Parallel work | N/A | `/fleet` for independent issues |
| CLI access | Direct shell + hooks | Built-in terminal in agent mode |
| Dependencies | None (CLI only) | None (CLI only) |

---

## Troubleshooting

### Copilot doesn't understand bd commands

1. Verify instructions are installed: `bd setup copilot --check`
2. Reload VS Code window
3. Ensure `bd` is on your PATH: `which bd` / `where bd`

### "bd: command not found"

Install beads: `brew install beads` or `npm install -g @beads/bd`

### No database found

```bash
bd init --quiet
```

---

## See Also

- [Detailed integration guide](../../docs/COPILOT_INTEGRATION.md)
- [Claude Code integration](../claude-code/README.md)
