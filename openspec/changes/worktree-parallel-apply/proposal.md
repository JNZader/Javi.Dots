# Proposal: Worktree-Based Parallel Apply for SDD

## Intent

The SDD `/sdd:apply` phase executes tasks sequentially in a single workspace. When tasks are independent (different files, no shared state), this wastes time — sub-agents wait idle while one finishes. Users working on changes with 4+ independent tasks (e.g., adding separate config sections to separate files) experience unnecessary serial delay.

This change adds **opt-in worktree-based parallel execution** to `/sdd:apply`. The orchestrator creates isolated git worktrees per task, launches sub-agents in parallel (each in its own worktree), then merges branches back sequentially. Default behavior remains sequential — nothing changes unless explicitly requested.

## Scope

### In Scope
- Add parallel apply dispatch logic to the SDD orchestrator (`CLAUDE.md` + `AGENTS.md`)
- Add trigger detection: user keywords ("parallel apply", "aplica en paralelo") and config flag (`apply.parallel: true`)
- Implement worktree lifecycle in orchestrator: create, dispatch, merge, cleanup
- Add optional `workdir` input parameter to `sdd-apply/SKILL.md` (sub-agent works in specified directory instead of project root)
- Add conditional `tasks.md` update behavior: sub-agents skip `tasks.md` marking in worktree mode; orchestrator updates centrally after merge
- Document `apply` config section in `openspec/config.yaml` schema
- Add `.worktrees/` to `.gitignore`

### Out of Scope
- Automatic task independence/dependency detection (user opts in, accepts merge conflict risk)
- Automatic merge conflict resolution (stop and report, user decides)
- TUI installer changes
- New skills or new shell scripts
- Go code changes
- PR creation per task branch (local merge only)
- Wave-based execution (that's `worktree-flow` + `wave-workflow` territory)

## Approach

**Orchestrator-only management** — borrow patterns from the existing `worktree-flow` skill and multi-perspective explore fan-out, but keep all worktree lifecycle logic in the orchestrator prompt (no new skills or scripts).

Flow:
1. Orchestrator detects parallel trigger (keyword or config)
2. Pre-flight: verify clean git state, check for leftover worktrees
3. For each task in the batch: `git worktree add .worktrees/sdd-{change}-task-{N} -b sdd/{change}/task-{N}`
4. Launch N parallel `sdd-apply` sub-agents (single message, N Task calls), each with `workdir` pointing to its worktree
5. Wait for all sub-agents to complete
6. Merge branches sequentially into current branch (`--no-ff`); stop on first conflict
7. Update `tasks.md` centrally based on sub-agent reports
8. Cleanup: `git worktree remove` + `git branch -d` + `git worktree prune`

The sub-agent does not know it's in a worktree. It receives a `workdir` and works there. When `workdir` is absent, behavior is unchanged (backward compatible).

## Affected Areas

| Area | Impact | Description |
|------|--------|-------------|
| `GentlemanClaude/CLAUDE.md` | Modified | Add `### Parallel Apply with Worktrees` section to SDD orchestrator |
| `AGENTS.md` | Modified | Mirror the same parallel apply section (kept in sync) |
| `GentlemanClaude/skills/sdd-apply/SKILL.md` | Modified | Add optional `workdir` input; conditional `tasks.md` update logic |
| `.gitignore` | Modified | Add `.worktrees/` entry |

## Risks

| Risk | Likelihood | Mitigation |
|------|------------|------------|
| Merge conflicts on shared files | Medium | User opts in accepting risk; orchestrator stops on first conflict and reports clearly |
| `tasks.md` concurrent update conflict | High (if not handled) | Sub-agents skip `tasks.md` in worktree mode; orchestrator updates centrally |
| Disk space with multiple worktrees | Low | Cap at 4 worktrees (configurable); cleanup immediately after merge |
| Orchestrator running git commands stretches delegation rule | Low | Document as coordination commands (environment setup), not implementation work |
| Sub-agent using absolute paths to main project | Low | Sub-agent uses relative paths implicitly via `openspec/changes/{change}/...`; worktree contains full copy |
| Partial failure (some sub-agents succeed, some fail) | Low | Merge successful branches first; report failed task; user retries or aborts |
| Git state corruption on crash mid-merge | Low | Document recovery: `git merge --abort`, `git worktree prune`, manual cleanup; pre-flight checks for leftover worktrees |

## Rollback Plan

All changes are to markdown documentation files (`CLAUDE.md`, `AGENTS.md`, `SKILL.md`, `.gitignore`). Rollback is a simple `git revert` of the commit. No runtime code, no database migrations, no binary artifacts. The feature is opt-in, so removing it affects zero existing workflows.

## Dependencies

- None. All changes are to documentation/prompt files. No external tools or libraries required beyond git (which is already a prerequisite for the project).

## Success Criteria

- [ ] Running `/sdd:apply` without parallel trigger behaves identically to current behavior (backward compatible)
- [ ] Running `/sdd:apply` with "parallel" keyword creates worktrees, dispatches parallel sub-agents, merges, and cleans up
- [ ] Setting `apply.parallel: true` in `openspec/config.yaml` enables parallel mode without keywords
- [ ] Merge conflict during sequential merge stops execution and reports conflicting files to user
- [ ] `tasks.md` is updated correctly after all merges complete (no duplicate marks, no missed marks)
- [ ] `.worktrees/` directory is cleaned up after successful parallel apply
- [ ] Max worktree cap (default 4) is respected
- [ ] Pre-flight check detects and warns about leftover worktrees from previous runs
