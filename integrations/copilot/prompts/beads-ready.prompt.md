---
mode: 'agent'
description: 'Find unblocked beads issues and claim one to work on'
---

Use the terminal to run beads CLI commands to find and claim ready work.

**Step 1 — List ready issues:**

Run in terminal:
```bash
bd ready --json
```

Parse the JSON output and present the issues showing:
- Issue ID
- Priority (P0=critical, P1=high, P2=medium, P3=low, P4=backlog)
- Type (bug / feature / task / epic / chore)
- Title

**Step 2 — Claim work:**

If there are ready issues, ask the user which one to work on. When they choose one:

1. Get full details by running in terminal:
   ```bash
   bd show <id> --json
   ```

2. Claim it by running in terminal:
   ```bash
   bd update <id> --claim --json
   ```

3. Show a confirmation and the issue description so the user has full context.

If there are no ready issues, suggest:
- Creating a new issue: "Shall I create a new issue?"
- Checking stale issues: run `bd stale --days 30 --json`
