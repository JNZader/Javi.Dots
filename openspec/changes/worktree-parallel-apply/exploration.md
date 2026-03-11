# Exploration: Worktree-Based Parallel Apply for SDD

## Current State

### How `/sdd:apply` Works Today

The apply flow is **sequential, single-workspace**:

1. **Orchestrator** receives `/sdd:apply <change-name>`
2. Reads `tasks.md` from `openspec/changes/{change-name}/tasks.md`
3. Batches tasks (e.g., "Phase 1, tasks 1.1-1.3") based on the Apply Strategy rule:
   > "For large task lists, batch tasks to sub-agents. Do NOT send all tasks at once."
4. Launches **ONE** `sdd-apply` sub-agent per batch via Task tool
5. Sub-agent reads specs, design, and tasks, then implements the batch **in the main working directory**
6. Sub-agent marks `[x]` on completed tasks in `tasks.md`
7. Orchestrator shows progress, asks user to continue with next batch
8. Repeat until all tasks done

**Key constraint**: Each sub-agent works in the same workspace (`/home/javier/Javi.Dots`). Tasks execute sequentially because the orchestrator waits for each batch to complete before launching the next. There is no mechanism for parallel sub-agent execution during apply.

### What the `sdd-apply` Sub-Agent Receives

From `GentlemanClaude/skills/sdd-apply/SKILL.md`:
- Change name
- Specific task(s) to implement (e.g., "Phase 1, tasks 1.1-1.3")
- Artifact store mode
- Access to: proposal, spec, design, tasks (reads all four before writing code)

The sub-agent's workflow (Step 3b, standard mode):
1. Read task description
2. Read relevant spec scenarios (acceptance criteria)
3. Read design decisions (constraints)
4. Read existing code patterns
5. Write the code
6. Mark task `[x]` in tasks.md
7. Report issues/deviations

**Nothing in the sub-agent skill references a working directory.** The sub-agent implicitly works wherever the Task tool places it (the project root by default).

### Orchestrator Dispatch Logic

From `GentlemanClaude/CLAUDE.md:459-469` (Apply Strategy):
```
For large task lists, batch tasks to sub-agents (e.g., "implement Phase 1, tasks 1.1-1.3").
Do NOT send all tasks at once -- break into manageable batches.
After each batch, show progress to user and ask to continue.
```

The Sub-Agent Launching Pattern (`CLAUDE.md:387-408`) shows a single Task call per dispatch. However, the orchestrator rules do **not prohibit** launching multiple Task calls simultaneously. The multi-perspective explore feature (already implemented) proves the orchestrator can fan out N parallel sub-agents in a single message.

### Existing Worktree Infrastructure

The `worktree-flow` skill (`~/.config/opencode/skills/workflow/worktree-flow/SKILL.md`) provides comprehensive worktree automation:

- **Directory layout**: `.worktrees/{task-id}/` with one branch per task
- **Lifecycle**: CREATE -> WORK -> VALIDATE -> PR -> CLEANUP
- **Wave execution**: Parallel agents, each in their own worktree
- **Shell scripts**: `worktree-create.sh`, `worktree-cleanup.sh`, `worktree-status.sh`, `worktree-pr.sh`
- **Configuration**: `.ai-config/worktree-flow.yaml` with branch prefix, max parallel, validation commands
- **Anti-patterns**: Don't use for sequential tasks, don't share worktrees between agents, don't skip validation

Key `worktree-flow` patterns directly applicable:
```bash
# Create worktree from current HEAD
git worktree add .worktrees/task-N -b feat/task-N origin/main

# Agent operates entirely within the worktree
# All git operations are scoped to this worktree

# After completion: merge or PR
```

The skill already defines integration points with `wave-workflow` and `playbooks`, but **has no integration with SDD apply**.

## Affected Areas

- `GentlemanClaude/CLAUDE.md` -- Orchestrator apply dispatch logic (Add parallel mode detection and worktree management)
- `AGENTS.md` -- Mirror of orchestrator SDD section (Keep in sync with CLAUDE.md changes)
- `GentlemanClaude/skills/sdd-apply/SKILL.md` -- Sub-agent skill (Minor: accept optional `workdir` context)
- `GentlemanClaude/skills/_shared/openspec-convention.md` -- Document `tasks.md` concurrent update considerations
- `.gitignore` -- Ensure `.worktrees/` is ignored (may already be there)

## Approaches

### 1. **Orchestrator-Only: Worktree Management in Orchestrator Prompt** -- The orchestrator manages worktrees before/after launching sub-agents

The orchestrator gains a parallel apply mode. When triggered, it:
1. Detects independent tasks (or trusts the user's opt-in)
2. Runs `git worktree add` for each task before launching sub-agents
3. Launches N parallel `sdd-apply` sub-agents, each with a `workdir` parameter pointing to their worktree
4. Waits for all sub-agents to complete
5. Merges branches back to current branch one at a time (sequential merge)
6. Reports conflicts to user if any
7. Cleans up worktrees

**What changes**:
- `CLAUDE.md` + `AGENTS.md`: Add parallel apply dispatch logic (~40 lines), trigger detection, merge strategy, cleanup
- `sdd-apply/SKILL.md`: Add optional `workdir` to "What You Receive" (sub-agent operates in specified directory instead of project root). One-line addition.
- No new skills needed

**Pros**:
- Minimal skill changes -- the sub-agent doesn't know or care it's in a worktree
- Follows the established pattern: orchestrator coordinates, sub-agents do focused work
- Reuses existing `worktree-flow` patterns without creating a dependency
- The orchestrator can use the Task tool's working directory to scope each sub-agent
- Single responsibility: orchestrator manages git state, sub-agents write code

**Cons**:
- Orchestrator becomes more complex (but it's prompt logic, not code)
- The orchestrator needs to run git commands (worktree add/remove/merge) which stretches its "delegate-only" rule -- but these are coordination commands, not implementation work
- `tasks.md` update becomes tricky: each sub-agent updates their own worktree's copy, but the main copy needs to reflect all completions after merge

**Effort**: Medium

### 2. **New Skill: `sdd-parallel-apply`** -- A dedicated skill for worktree-based parallel execution

Create a new skill that handles the full lifecycle:
1. Reads tasks.md, identifies parallelizable tasks
2. Creates worktrees
3. Returns instructions to orchestrator for parallel dispatch
4. After all sub-agents complete, handles merge and cleanup

**What changes**:
- New `GentlemanClaude/skills/sdd-parallel-apply/SKILL.md`
- `CLAUDE.md` + `AGENTS.md`: Route to new skill when parallel apply triggered
- Register in AGENTS.md skill table

**Pros**:
- Clean separation -- all worktree logic in one skill file
- Testable in isolation
- Can be evolved independently

**Cons**:
- Sub-agents can't launch sub-sub-agents, so this skill would need to return a "dispatch plan" to the orchestrator, adding a round-trip
- Duplicates orchestrator coordination logic in a skill
- Skill proliferation: adds another file to maintain and sync
- The actual parallel dispatch still happens in the orchestrator (only the orchestrator can launch parallel Task calls)

**Effort**: Medium-High

### 3. **Hybrid: Orchestrator Dispatch + Helper Script** -- Orchestrator manages dispatch, shell scripts handle git mechanics

The orchestrator uses the `worktree-flow` skill's script patterns as helpers:
1. Runs a shell script to create all worktrees at once
2. Launches parallel sub-agents (one per worktree)
3. Runs a shell script to merge all branches and cleanup

**What changes**:
- `CLAUDE.md` + `AGENTS.md`: Parallel dispatch logic
- Create actual shell scripts (e.g., `scripts/sdd-worktree-create.sh`, `scripts/sdd-worktree-merge.sh`)
- `sdd-apply/SKILL.md`: Accept workdir

**Pros**:
- Shell scripts are reusable outside SDD
- Git operations are centralized and testable
- Less prompt complexity in orchestrator

**Cons**:
- Introduces runtime dependencies (bash scripts that need to exist on disk)
- Scripts need to be maintained, tested, and shipped
- Adds complexity for what is essentially 5-6 git commands
- Over-engineering for an opt-in feature that most users won't use

**Effort**: High

## Recommendation

**Approach 1: Orchestrator-Only** with patterns borrowed from `worktree-flow`.

Rationale:

1. **The orchestrator IS the right place for this.** Parallel dispatch is a coordination concern, not an implementation concern. The orchestrator already coordinates sequential batches; adding parallel dispatch is a natural extension. This follows the same pattern as multi-perspective explore (orchestrator fans out N sub-agents, collects results).

2. **The `sdd-apply` skill should stay simple.** The sub-agent's job is: read specs, write code, mark tasks done. Whether it's working in the main directory or a worktree is irrelevant to its logic. The only change needed is accepting an optional `workdir` path in its prompt context.

3. **No shell scripts needed.** The orchestrator can run git commands directly via Bash tool calls. The commands are simple and well-understood:
   - `git worktree add .worktrees/task-N -b sdd/{change}/task-N`
   - `git worktree remove .worktrees/task-N`
   - `git merge sdd/{change}/task-N`
   These don't warrant wrapping in scripts for an opt-in feature.

4. **Opt-in is critical.** Default behavior MUST remain sequential. The user explicitly requests parallel mode via keyword or config. This means the new code path only fires when requested -- zero risk to existing users.

### Recommended Architecture

```
User: "/sdd:apply change-name --parallel" (or "aplica en paralelo")
         |
         v
Orchestrator detects parallel trigger
         |
         v
Read tasks.md, select independent tasks
         |
         v
FOR EACH task:
  git worktree add .worktrees/sdd-{change}-task-N -b sdd/{change}/task-N
         |
         v
Launch N parallel sdd-apply sub-agents (single message, N Task calls)
  Each receives: workdir=".worktrees/sdd-{change}-task-N", task assignment
         |
         v
Wait for all sub-agents to complete
         |
         v
FOR EACH completed branch (sequential):
  git merge sdd/{change}/task-N --no-ff
  IF conflict: STOP, report to user
         |
         v
Update tasks.md in main with all [x] marks
         |
         v
Cleanup: git worktree remove + git branch -d for each
         |
         v
Report results to user
```

### Trigger Mechanism

| Trigger | Source | Priority |
|---------|--------|----------|
| User keywords: "parallel", "en paralelo", "parallel apply" | User message | Highest |
| Config: `apply.parallel: true` in `openspec/config.yaml` | Project config | Medium |
| Default | N/A | Sequential (no change) |

### Task Independence

**Trust the user (MVP).** If they opt into parallel mode, they accept the risk of merge conflicts. The orchestrator should:
1. Warn once: "Parallel mode may cause merge conflicts if tasks touch the same files. Continue?"
2. If conflicts occur during merge: stop, report which tasks conflict, ask user how to proceed
3. NOT attempt automatic dependency detection (out of scope, hard to get right, false sense of security)

### tasks.md Concurrency

Each worktree gets its own copy of the codebase, including `openspec/changes/{change}/tasks.md`. Each sub-agent marks `[x]` on its assigned tasks in its own copy. After merge:
- The orchestrator does a final pass to consolidate `[x]` marks from all branches into the main `tasks.md`
- Alternative: the orchestrator updates `tasks.md` centrally after all merges complete, based on sub-agent reports (simpler, avoids merge conflicts on tasks.md itself)

**Recommended**: Orchestrator updates tasks.md centrally after all sub-agents report. Sub-agents should NOT update tasks.md in worktree mode (it will cause merge conflicts). Instead, they report completed tasks in their return summary, and the orchestrator marks them.

### Branch Naming Convention

```
sdd/{change-name}/task-{N}
```

Examples:
- `sdd/worktree-parallel-apply/task-1.1`
- `sdd/worktree-parallel-apply/task-1.2`
- `sdd/worktree-parallel-apply/task-1.3`

Worktree directories:
```
.worktrees/
  sdd-worktree-parallel-apply-task-1.1/
  sdd-worktree-parallel-apply-task-1.2/
  sdd-worktree-parallel-apply-task-1.3/
```

### Merge Strategy

1. Sequential merge: merge one branch at a time into current branch
2. Use `--no-ff` to preserve branch history
3. On conflict: **STOP immediately**. Report conflicting files and which tasks caused it. Ask user to resolve manually.
4. On success: delete branch, remove worktree
5. After all merges: run `git worktree prune`

### .gitignore

Ensure `.worktrees/` is in `.gitignore`. Check if it's already there; if not, add it as part of the change.

## Scope

### IN Scope
- Orchestrator parallel dispatch logic in `CLAUDE.md` and `AGENTS.md`
- Trigger detection (user keywords + config flag)
- Worktree creation before parallel launch
- Parallel sub-agent dispatch (N Task calls in single message)
- Sequential branch merge after completion
- Conflict detection and reporting (stop on first conflict)
- Worktree and branch cleanup
- Central `tasks.md` update after merge
- Minor `sdd-apply/SKILL.md` update (accept optional workdir, skip tasks.md update in worktree mode)
- `.gitignore` update for `.worktrees/`
- Config addition: `apply.parallel` option in `openspec/config.yaml`

### OUT of Scope
- Automatic task dependency detection
- Automatic merge conflict resolution
- TUI changes
- PR creation per task (this is local merge, not PR-based)
- Wave-based execution (that's `worktree-flow` + `wave-workflow` territory)
- Dependency installation in worktrees (Go projects don't typically need this)
- New skills or new shell scripts

## Risks

- **R1: Merge conflicts on shared files** -- If two tasks modify the same file (e.g., both add functions to the same package), merge will conflict. Mitigation: User opts in, accepting this risk. Orchestrator stops on first conflict and reports clearly. User resolves manually.

- **R2: tasks.md conflict** -- All sub-agents try to update tasks.md in their worktree. Mitigation: Sub-agents skip tasks.md updates in worktree mode; orchestrator updates centrally after merge.

- **R3: Disk space** -- Each worktree is a full working tree copy (minus `.git` which is shared). For Javi.Dots (~50MB), this is negligible. For larger projects, N worktrees = N * project_size. Mitigation: Cap max parallel worktrees (default 4, configurable). Cleanup immediately after merge.

- **R4: Orchestrator delegation rule tension** -- The orchestrator runs `git worktree add/remove/merge` commands directly via Bash. This is technically "work" that could be delegated. However, these are coordination commands (creating the environment for sub-agents), not implementation work. The same way the orchestrator can read config files to make dispatch decisions, it can run git commands to set up execution environments. This is a reasonable exception that should be documented.

- **R5: Sub-agent isolation assumptions** -- The `sdd-apply` skill assumes it can read specs/design/tasks from the standard openspec paths. In a worktree, these files are at the worktree root, which is correct (worktree is a full copy). However, if the sub-agent uses absolute paths to the main project, it would read the wrong files. Mitigation: The orchestrator passes the worktree path; the sub-agent should use relative paths (which the skill already does implicitly via `openspec/changes/{change}/...`).

- **R6: Partial failure** -- If 2 of 3 sub-agents succeed and 1 fails, the orchestrator needs to handle mixed results. Mitigation: Merge successful branches first. Report the failed task. User decides whether to retry the failed task sequentially or abort.

- **R7: Git state corruption on crash** -- If the process is interrupted mid-merge, git can be left in a conflicted state with dangling worktrees. Mitigation: Document recovery procedure (`git merge --abort`, `git worktree prune`, manual `.worktrees/` cleanup). The orchestrator should check for leftover worktrees at the start of parallel apply.

## Ready for Proposal

**Yes** -- The scope is well-defined and follows established patterns:

- Primary change: Orchestrator dispatch logic in `CLAUDE.md` + `AGENTS.md` (~50-60 lines of prompt additions)
- Secondary change: Minor `sdd-apply/SKILL.md` update (accept workdir, conditional tasks.md behavior)
- Config addition: `apply.parallel` option in `openspec/config.yaml`
- Gitignore: Add `.worktrees/` if not present
- No new skills, no shell scripts, no code changes to the Go codebase
- Follows the proven pattern of multi-perspective explore fan-out
- Entirely opt-in; zero impact on existing sequential behavior

The orchestrator should tell the user: "Ready to proceed with `/sdd:propose worktree-parallel-apply`. The change adds opt-in parallel execution to `/sdd:apply` using git worktrees. When the user says 'aplica en paralelo' or sets `apply.parallel: true` in config, the orchestrator creates isolated worktrees per task, launches sub-agents in parallel, merges branches back sequentially, and cleans up. Default behavior remains sequential -- nothing changes unless explicitly requested."
