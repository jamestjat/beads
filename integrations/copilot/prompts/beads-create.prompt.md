---
mode: 'agent'
description: 'Create a new beads issue interactively'
---

Create a new beads issue using the terminal. The user may have already provided details in their message — use what they supplied and ask only for missing required fields.

Gather:
1. **Title** (required) — a short, action-oriented summary
2. **Type** (default: `task`) — one of: `bug`, `feature`, `task`, `epic`, `chore`
3. **Priority** (default: `2`) — `0`=critical, `1`=high, `2`=medium, `3`=low, `4`=backlog
4. **Description** (recommended) — context that will help future work

Optional:
- **Parent issue** — for child tasks of an epic
- **Discovered from** — if this came up while working on another issue (enter the parent ID)
- **Blocks** — if this must be done before another issue

Run in terminal:
```bash
bd create "Title" --description="Detailed context" -t <type> -p <priority> --json
```

If a "discovered-from" relationship was given, also run:
```bash
bd dep add <new-id> discovered-from:<parent-id>
```

Show the created issue ID, title, and a brief confirmation.
