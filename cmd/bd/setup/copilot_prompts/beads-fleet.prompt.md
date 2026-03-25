---
mode: 'agent'
description: 'Parallelize beads work across fleet subagents'
---

Use the terminal to find independent beads issues and distribute them across fleet subagents.

**Step 1 — List ready issues:**

Run in terminal:
```bash
bd ready --json
```

Parse the JSON output and present the issues showing ID, priority, type, and title.

**Step 2 — Check dependencies:**

For each ready issue, check if it blocks or is related to other ready issues:
```bash
bd dep tree <id> --json
```

Identify issues that have **no mutual dependencies** with each other — these can be worked in parallel.

**Step 3 — Compose fleet command:**

Group independent issues and compose a `/fleet` command that assigns one issue per subagent. Each subagent should:

1. Claim the issue: `bd update <id> --claim --json`
2. Read the full details: `bd show <id> --json`
3. Implement the fix or feature
4. Run relevant tests
5. Close the issue: `bd close <id> --reason "Done" --json`

**Important**: Subagents must NOT run `bd dolt push`. Only the orchestrator (this session) pushes after all subagents complete.

**Step 4 — After fleet completes:**

Once all subagents finish:
```bash
bd dolt push
```

Report which issues were completed and any that had problems.
