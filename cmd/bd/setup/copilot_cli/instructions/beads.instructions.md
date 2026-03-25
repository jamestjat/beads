# Beads Issue Tracking (CLI)

This project uses **bd (beads)** for all issue tracking. Do NOT use markdown TODOs.

## Hooks (Pre-configured)

- **sessionStart**: `bd prime` runs automatically to load context
- **sessionEnd**: `bd dolt push` runs automatically to sync

You do not need to manually prime or push — hooks handle it.

## Core Workflow

1. **Find work**: `bd ready --json`
2. **Claim**: `bd update <id> --claim --json`
3. **Work**: Implement, test, document
4. **Discover**: `bd create "Title" --description="Context" -t bug -p 1 --deps discovered-from:<parent-id> --json`
5. **Complete**: `bd close <id> --reason "Done" --json`

## Parallel Work with /fleet

Use `/fleet` to parallelize independent ready issues:

1. Run `bd ready --json` to list unblocked work
2. Check `bd dep tree` to find issues with no mutual dependencies
3. Use `/fleet` to distribute one issue per subagent
4. **Important**: Only the orchestrator session should run `bd dolt push` — subagents must NOT push

## Delegation

Use `/delegate` for tangential work discovered during your main task:
- Documentation updates
- Refactoring separate modules
- Filing and implementing discovered bugs

Keep core issue work (claim, close, dependency changes) in your local session.

## Proactive Issue Discovery

While working on any task, **immediately** file issues when you find:
- Bugs, broken behavior, error handling gaps
- Missing tests, documentation gaps, code quality issues
- Performance problems or tech debt

```bash
bd create "<title>" --description="Found while working on <current-task>: <details>" -t bug -p <priority> --deps discovered-from:<current-id> --json
```

File and continue working — do not stop your current task.

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
bd prime              # Full workflow context
```

## Session Management

- `/new` — Start fresh conversation (triggers sessionEnd hook, pushes changes)
- `/clear` — Abandon session (no push)
- `/undo` — Revert last turn and file changes
- `/compact` — Manually compress context when sessions grow large
- `/session plan` — Review current implementation plan

## Rules

- Always use `--json` for programmatic output
- Always include `--description` when creating issues
- Link discovered work with `discovered-from` dependencies
- Use `/allow-all` or `copilot --allow-tool='shell(bd:*)'` to pre-approve bd commands
