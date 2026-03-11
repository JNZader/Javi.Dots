# Delta for Apply

## ADDED Requirements

### Requirement: Parallel Apply Trigger Detection

The orchestrator MUST detect parallel apply mode from two sources, evaluated in priority order:

1. **User keywords** (highest): message contains "parallel apply", "aplica en paralelo", "apply in parallel", or "parallel"
2. **Config flag** (medium): `openspec/config.yaml` contains `apply.parallel: true`

If neither trigger is present, the orchestrator MUST use sequential apply (current behavior, unchanged).

#### Scenario: User keyword triggers parallel mode

- GIVEN a change "my-feature" with `tasks.md` containing 3 independent tasks
- WHEN the user says "/sdd:apply my-feature aplica en paralelo"
- THEN the orchestrator activates parallel apply mode
- AND creates one worktree per task

#### Scenario: Config flag triggers parallel mode

- GIVEN `openspec/config.yaml` contains `apply: { parallel: true }`
- AND no parallel keyword in the user message
- WHEN the user says "/sdd:apply my-feature"
- THEN the orchestrator activates parallel apply mode

#### Scenario: Default remains sequential

- GIVEN no config flag and no parallel keyword
- WHEN the user says "/sdd:apply my-feature"
- THEN the orchestrator uses sequential batch mode (current behavior)
- AND no worktrees are created

### Requirement: Pre-Flight Validation

Before creating worktrees, the orchestrator MUST verify:

1. Git working tree is clean (no uncommitted changes)
2. No leftover `.worktrees/` directory from a previous run
3. Current branch is not in a merge/rebase state

If any check fails, the orchestrator MUST report the issue and abort parallel apply.

#### Scenario: Clean state passes pre-flight

- GIVEN a clean git working tree with no leftover worktrees
- WHEN parallel apply pre-flight runs
- THEN pre-flight passes
- AND worktree creation proceeds

#### Scenario: Dirty working tree aborts

- GIVEN uncommitted changes exist in the working tree
- WHEN parallel apply pre-flight runs
- THEN pre-flight fails with message indicating uncommitted changes
- AND no worktrees are created

#### Scenario: Leftover worktrees detected

- GIVEN `.worktrees/` directory exists from a previous interrupted run
- WHEN parallel apply pre-flight runs
- THEN pre-flight fails with message listing leftover worktrees
- AND suggests cleanup commands (`git worktree remove`, `git worktree prune`)

### Requirement: Worktree Lifecycle Management

The orchestrator MUST manage the full worktree lifecycle:

1. **Create**: `git worktree add .worktrees/sdd-{change}-task-{id} -b sdd/{change}/task-{id}` for each task
2. **Dispatch**: Launch one `sdd-apply` sub-agent per worktree, passing `workdir` parameter
3. **Merge**: After all sub-agents complete, merge branches one at a time into current branch using `--no-ff`
4. **Cleanup**: `git worktree remove .worktrees/sdd-{change}-task-{id}` + `git branch -d sdd/{change}/task-{id}` for each, then `git worktree prune`

#### Scenario: Successful parallel lifecycle

- GIVEN 3 tasks assigned for parallel apply on change "add-widgets"
- WHEN all 3 sub-agents complete successfully
- THEN 3 worktrees are created at `.worktrees/sdd-add-widgets-task-1.1`, `.worktrees/sdd-add-widgets-task-1.2`, `.worktrees/sdd-add-widgets-task-1.3`
- AND 3 branches exist: `sdd/add-widgets/task-1.1`, `sdd/add-widgets/task-1.2`, `sdd/add-widgets/task-1.3`
- AND branches are merged sequentially into the current branch
- AND all worktrees are removed after merge
- AND all task branches are deleted after merge
- AND `git worktree prune` is run

#### Scenario: Branch naming convention

- GIVEN a change named "worktree-parallel-apply" with task 2.1
- WHEN a worktree is created for that task
- THEN the branch is named `sdd/worktree-parallel-apply/task-2.1`
- AND the worktree directory is `.worktrees/sdd-worktree-parallel-apply-task-2.1`

### Requirement: Parallel Dispatch

The orchestrator MUST launch all parallel sub-agents in a single message (multiple Task calls in one response). Each sub-agent MUST receive:

1. The change name
2. The specific task(s) to implement
3. The `workdir` path pointing to its worktree
4. Artifact store mode
5. Instruction to skip `tasks.md` updates (orchestrator handles centrally)

#### Scenario: Parallel dispatch in single message

- GIVEN 3 tasks assigned for parallel apply
- WHEN the orchestrator dispatches sub-agents
- THEN exactly 3 Task calls are made in a single message (not sequential)
- AND each Task call includes a unique `workdir` path

#### Scenario: Sub-agent receives workdir

- GIVEN a sub-agent is launched for task 1.2 in worktree mode
- WHEN the sub-agent reads its prompt context
- THEN `workdir` is set to `.worktrees/sdd-{change}-task-1.2`
- AND the sub-agent operates in that directory for all file operations

### Requirement: Max Worktree Cap

The orchestrator MUST NOT create more than `max_worktrees` worktrees simultaneously. The default cap is 4. The cap MAY be configured via `openspec/config.yaml`:

```yaml
apply:
  parallel: true
  max_worktrees: 4
```

If more tasks than `max_worktrees` are assigned, the orchestrator MUST batch them: run the first N in parallel, wait for completion and merge, then run the next N.

#### Scenario: Tasks within cap

- GIVEN 3 tasks and `max_worktrees: 4`
- WHEN parallel apply runs
- THEN all 3 tasks run in parallel in a single wave

#### Scenario: Tasks exceed cap

- GIVEN 6 tasks and `max_worktrees: 4`
- WHEN parallel apply runs
- THEN 4 tasks run in the first wave
- AND after the first wave merges, the remaining 2 tasks run in a second wave

#### Scenario: Custom cap from config

- GIVEN `openspec/config.yaml` contains `apply: { parallel: true, max_worktrees: 2 }`
- WHEN 4 tasks are assigned for parallel apply
- THEN 2 tasks run in the first wave
- AND 2 tasks run in the second wave

### Requirement: Sequential Merge with Conflict Detection

After all sub-agents in a wave complete, the orchestrator MUST merge branches one at a time. On merge conflict, the orchestrator MUST:

1. Stop merging immediately
2. Report which branch/task caused the conflict
3. Report which files are in conflict
4. Leave the remaining unmerged branches intact for user resolution
5. NOT attempt automatic conflict resolution

#### Scenario: All merges succeed

- GIVEN 3 branches complete without overlapping file changes
- WHEN sequential merge runs
- THEN all 3 branches merge successfully with `--no-ff`
- AND orchestrator proceeds to cleanup

#### Scenario: Merge conflict on second branch

- GIVEN 3 branches where branch 2 conflicts with branch 1's changes
- WHEN sequential merge runs
- THEN branch 1 merges successfully
- AND branch 2 merge fails with conflict
- AND orchestrator reports conflicting files and stops
- AND branch 3 is NOT merged (left intact for later)
- AND the user is told which branches remain unmerged

#### Scenario: Partial failure recovery

- GIVEN 3 sub-agents where sub-agent 2 fails but 1 and 3 succeed
- WHEN the orchestrator collects results
- THEN successful branches (1 and 3) are merged
- AND the failed task is reported to the user
- AND the user decides whether to retry or skip

### Requirement: Central tasks.md Update

In parallel apply mode, the orchestrator MUST update `tasks.md` centrally after all merges complete. Sub-agents MUST NOT update `tasks.md` in worktree mode (to avoid merge conflicts on the tasks file itself).

The orchestrator marks tasks `[x]` based on sub-agent completion reports.

#### Scenario: Orchestrator marks tasks after merge

- GIVEN 3 sub-agents each report their task as complete
- WHEN all merges succeed
- THEN the orchestrator updates `tasks.md` in the main working directory
- AND tasks 1.1, 1.2, and 1.3 are marked `[x]`

#### Scenario: Sub-agent skips tasks.md in worktree mode

- GIVEN a sub-agent receives `workdir` (indicating worktree mode)
- WHEN the sub-agent completes its task
- THEN the sub-agent does NOT modify `tasks.md`
- AND the sub-agent reports completion status in its return summary

### Requirement: Sub-Agent workdir Parameter

The `sdd-apply` skill MUST accept an optional `workdir` parameter in its prompt context. When present:

1. The sub-agent operates in the specified directory for all file operations
2. The sub-agent MUST NOT update `tasks.md` (orchestrator handles this)
3. The sub-agent reads specs/design/tasks from the worktree's `openspec/` directory (relative path)

When `workdir` is absent, behavior is identical to current (backward compatible).

#### Scenario: workdir present - sub-agent uses it

- GIVEN `workdir` is set to `.worktrees/sdd-my-change-task-1.1`
- WHEN the sub-agent implements its task
- THEN all file reads and writes happen relative to that directory
- AND `tasks.md` is NOT updated by the sub-agent

#### Scenario: workdir absent - backward compatible

- GIVEN `workdir` is NOT present in the sub-agent prompt
- WHEN the sub-agent implements its task
- THEN the sub-agent works in the project root (current behavior)
- AND `tasks.md` IS updated by the sub-agent (current behavior)

## MODIFIED Requirements

### Requirement: sdd-apply "What You Receive" Section

(Previously: sub-agent receives change name, task(s), and artifact store mode only.)

The sub-agent now also receives an optional `workdir` parameter. The "What You Receive" section of the skill MUST list:

- Change name
- Specific task(s) to implement
- Artifact store mode
- **workdir** (optional) — when present, the sub-agent operates in this directory and skips `tasks.md` updates

#### Scenario: Updated skill documents workdir

- GIVEN the `sdd-apply/SKILL.md` file
- WHEN the "What You Receive" section is read
- THEN it lists `workdir` as an optional parameter
- AND explains that when present, `tasks.md` updates are skipped
