---
mode: 'agent'
description: 'Convert a plan document into a beads epic with dependency-linked tasks'
---

Convert a plan document into a beads epic and tasks for cross-session tracking.

## Steps

1. **Identify the plan**
   - If the user provided a file path or pasted content, use that.
   - Otherwise ask: "Please paste the plan or provide the file path."

2. **Parse the plan structure**
   - Title: First `# Plan:` or `#` heading
   - Summary: Content under `## Summary` or `## Overview`
   - Phases/tasks: Each `### Phase N` or `### N.` section

3. **Create the epic** by running in terminal:
   ```bash
   bd create "[Plan Title]" -t epic -p 1 --description="[summary]. Plan converted from [filename]." --json
   ```
   Note the returned epic ID from the JSON output.

4. **Create tasks from each phase** (keep descriptions concise — first paragraph only):
   ```bash
   bd create "[Phase title]" -t task -p 2 --description="[first paragraph of phase content]" --json
   ```

5. **Link tasks to the epic**:
   ```bash
   bd dep add <task-id> parent:<epic-id>
   ```

6. **Add sequential phase dependencies** (later phases depend on earlier ones):
   ```bash
   bd dep add <phase2-id> blocks:<phase1-id>
   ```

7. **Return a concise summary** (not raw JSON):
   ```
   Created from: [filename or "pasted plan"]

   Epic: [title] ([epic-id])
     +-- [Phase 1 title] ([id]) -- ready
     +-- [Phase 2 title] ([id]) -- blocked by [phase1-id]
     +-- [Phase 3 title] ([id]) -- blocked by [phase2-id]

   Total: [N] tasks
   Run `bd ready` to start the first phase.
   ```

## Notes

- Original plan content is preserved — this only creates tracking issues.
- Task descriptions use the first paragraph only (keeps them scannable).
- Sequential phases get automatic `blocks` dependencies so they unlock one at a time.
