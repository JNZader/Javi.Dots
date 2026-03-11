# Tasks: Worktree-Based Parallel Apply for SDD

## Phase 1: Foundation (sdd-apply skill update + .gitignore)

- [x] 1.1 Update `GentlemanClaude/skills/sdd-apply/SKILL.md` "What You Receive" section: add `workdir` (optional) parameter with description — "When present, the sub-agent operates in this directory instead of the project root and MUST NOT update `tasks.md`."
- [x] 1.2 Update `GentlemanClaude/skills/sdd-apply/SKILL.md` Step 4 ("Mark Tasks Complete"): add conditional — "If `workdir` was provided (worktree mode), skip `tasks.md` updates. Report completed tasks in your return summary instead. The orchestrator will update `tasks.md` centrally."
- [x] 1.3 Update `GentlemanClaude/skills/sdd-apply/SKILL.md` metadata version from `"2.0"` to `"2.1"`
- [x] 1.4 Add `.worktrees/` entry to root `.gitignore`

## Phase 2: Orchestrator — Parallel Apply Section (CLAUDE.md)

- [x] 2.1 Add `### Parallel Apply with Worktrees` section to `GentlemanClaude/CLAUDE.md` after the existing `### Apply Strategy` section. Include: trigger detection (user keywords + config flag), pre-flight checks (clean git state, no leftover worktrees), worktree creation command pattern (`git worktree add .worktrees/sdd-{change}-task-{id} -b sdd/{change}/task-{id}`), parallel dispatch pattern (N Task calls in single message with `workdir` parameter), sequential merge strategy (`--no-ff`, stop on first conflict), central `tasks.md` update after merge, cleanup commands (`git worktree remove` + `git branch -d` + `git worktree prune`), and max worktree cap (default 4, configurable via `apply.max_worktrees`).

## Phase 3: Orchestrator Mirror (AGENTS.md)

- [x] 3.1 Add the same `### Parallel Apply with Worktrees` section to `AGENTS.md` under the SDD orchestrator section, mirroring the content added to `GentlemanClaude/CLAUDE.md` in task 2.1. Keep both files in sync.

## Phase 4: Verification

- [ ] 4.1 Verify `sdd-apply/SKILL.md` backward compatibility: confirm that when `workdir` is absent, the skill instructions produce identical behavior to the pre-change version (sequential apply, sub-agent updates `tasks.md`)
- [ ] 4.2 Verify trigger detection: confirm the orchestrator section documents both keyword triggers ("parallel apply", "aplica en paralelo", "parallel") and config trigger (`apply.parallel: true` in `openspec/config.yaml`)
- [ ] 4.3 Verify merge strategy documentation: confirm the orchestrator section specifies sequential merge with `--no-ff`, stop-on-first-conflict behavior, and partial failure handling
- [ ] 4.4 Verify `AGENTS.md` and `CLAUDE.md` parallel apply sections are identical in content
- [ ] 4.5 Verify `.gitignore` contains `.worktrees/` entry
