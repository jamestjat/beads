---
mode: 'ask'
description: 'Show the beads workflow guide for AI-supervised coding'
---

# Beads Workflow Guide

**beads (bd)** is a Dolt-powered issue tracker built for AI-supervised coding workflows.

## 1. Find Ready Work

Ask me: *"What issues are ready to work on?"*
Or run: `bd ready --json`

## 2. Claim Your Task

Claim work atomically (sets assignee + in-progress in one step):
```bash
bd update bd-42 --claim --json
```

## 3. Work on It

Implement, test, and document the feature or fix.

## 4. Discover New Work

As you work, file related issues and link them:
```bash
bd create "Found bug" -t bug -p 1 --deps discovered-from:bd-42 --json
```

## 5. Complete the Task

```bash
bd close bd-42 --reason "Fixed timeout in auth middleware" --json
```

## 6. Sync to Remote

At end of session:
```bash
bd dolt push
```

---

## Quick CLI Reference

```bash
bd ready                              # Unblocked issues
bd create "Title" -t task -p 2       # New issue
bd update <id> --claim               # Claim work
bd show <id>                         # Full details
bd close <id> --reason "Done"        # Complete
bd dep add <issue> <blocks>          # Add dependency
bd dolt push                         # Sync to remote
bd prime                             # AI workflow context (~2k tokens)
```

## Priorities

| Value | Meaning |
|-------|---------|
| 0 | Critical — security, data loss, broken builds |
| 1 | High — major features, important bugs |
| 2 | Medium — default |
| 3 | Low — polish, optimization |
| 4 | Backlog — future ideas |

For full documentation run `bd --help` or see the project's COPILOT_INTEGRATION.md.

## Copilot CLI Tips

- Use `/plan` for complex multi-file changes, then convert to beads issues
- Use `/fleet` to parallelize independent ready issues across subagents
- Use `/delegate` for tangential tasks (docs, refactoring)
- Use `/undo` to revert the last turn if an approach fails
- Use `/compact` when context grows large in long sessions
- Use `/session plan` to review the current implementation plan
- Pre-approve bd commands: `/allow-all` or `copilot --allow-tool='shell(bd:*)'`
