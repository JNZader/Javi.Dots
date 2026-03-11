# Proposal: Upstream Sync with agent-teams-lite v3.3.2

## Intent

Javi.Dots SDD skills are at v2.0 while upstream agent-teams-lite is at v3.3.2 and engram is at v1.9.4. The gap means Javi.Dots lacks: hybrid persistence mode, inline engram instructions (critical for installed skills where `_shared/` is unreachable), mandatory persist steps, skill registry loading, state.yaml tracking, sub-agent context rules, and delegation anti-patterns. This sync brings full feature parity with upstream while preserving all Javi.Dots-specific customizations (personality, domain routing, skill tables, plugin detection).

## Scope

### In Scope
- Replace 3 shared convention files (`_shared/persistence-contract.md`, `engram-convention.md`, `openspec-convention.md`) with upstream v3.3.2 content
- Replace 9 SDD skill files (`sdd-{init,explore,propose,spec,design,tasks,apply,verify,archive}/SKILL.md`) with upstream v3.3.2 content
- Create new `skill-registry/SKILL.md` from upstream v3.3.1
- Surgically merge orchestrator additions into `GentlemanClaude/CLAUDE.md` (delegation rules, anti-patterns, task escalation, hybrid mode) while preserving all Javi.Dots sections
- Surgically merge orchestrator additions into `AGENTS.md` (same additions + skill-registry table entry + hybrid in auto-mode resolution)
- Fix `skills/setup.sh` to sync `_shared/` directory to user config

### Out of Scope
- Non-SDD sub-agent engram awareness (domain orchestrators saving discoveries) — separate follow-up
- Engram v1.9.4 API changes (confirmed compatible — no breaking changes to mem_save/mem_search/mem_get_observation/mem_update)
- Migrating existing openspec changes to include state.yaml (skills handle missing state.yaml gracefully)
- Updating openspec/config.yaml (no schema changes needed)

## Approach

**Hybrid: Wholesale replace for skills, surgical merge for orchestrator files.**

1. **Shared conventions (3 files)**: Wholesale replace from upstream. These are foundation files that all skills reference. Zero Javi.Dots customizations.
2. **SDD skills (9 files)**: Wholesale replace from upstream. All verified to have zero Javi.Dots-specific content. Fetch from `raw.githubusercontent.com/Gentleman-Programming/agent-teams-lite/main/skills/sdd-*/SKILL.md`.
3. **New skill-registry (1 file)**: Create from upstream. Fetch from `raw.githubusercontent.com/Gentleman-Programming/agent-teams-lite/main/skills/skill-registry/SKILL.md`.
4. **Orchestrator files (2 files)**: Identify the SDD Orchestrator section boundaries in each file. Replace only those sections with upstream content. Preserve personality, domain routing, skill tables, plugin detection, SDD triggers. Add `hybrid` to AGENTS.md auto-mode resolution chain.
5. **Infrastructure (1 file)**: Add `_shared/` sync logic to `skills/setup.sh`.

## Affected Areas

| Area | Impact | Description |
|------|--------|-------------|
| `GentlemanClaude/skills/_shared/persistence-contract.md` | Modified | Add hybrid mode, sub-agent context rules, skill registry, state persistence, detail level |
| `GentlemanClaude/skills/_shared/engram-convention.md` | Modified | Add inline call NOTE, state artifact, mem_update improvements |
| `GentlemanClaude/skills/_shared/openspec-convention.md` | Modified | Add state.yaml, archive structure improvements |
| `GentlemanClaude/skills/sdd-init/SKILL.md` | Modified | Wholesale replace — add Step 4 (Skill Registry), hybrid mode, Step 5 (Persist), state.yaml |
| `GentlemanClaude/skills/sdd-explore/SKILL.md` | Modified | Wholesale replace — add Step 1 (Load Registry), inline engram, hybrid, persist |
| `GentlemanClaude/skills/sdd-propose/SKILL.md` | Modified | Wholesale replace — same pattern |
| `GentlemanClaude/skills/sdd-spec/SKILL.md` | Modified | Wholesale replace — same pattern |
| `GentlemanClaude/skills/sdd-design/SKILL.md` | Modified | Wholesale replace — same pattern |
| `GentlemanClaude/skills/sdd-tasks/SKILL.md` | Modified | Wholesale replace — same pattern |
| `GentlemanClaude/skills/sdd-apply/SKILL.md` | Modified | Wholesale replace — add registry, inline engram, hybrid, persist, TDD from config |
| `GentlemanClaude/skills/sdd-verify/SKILL.md` | Modified | Wholesale replace — add registry, inline engram, hybrid, persist, Spec Compliance Matrix |
| `GentlemanClaude/skills/sdd-archive/SKILL.md` | Modified | Wholesale replace — same pattern |
| `GentlemanClaude/skills/skill-registry/SKILL.md` | New | Create from upstream — skill dependency resolution and registry loading |
| `GentlemanClaude/CLAUDE.md` | Modified | Surgical merge — SDD orchestrator section only |
| `AGENTS.md` | Modified | Surgical merge — SDD orchestrator section + skill table + auto-mode resolution |
| `skills/setup.sh` | Modified | Add `_shared/` directory sync logic |

## Risks

| Risk | Likelihood | Mitigation |
|------|------------|------------|
| Javi.Dots-specific sections in CLAUDE.md/AGENTS.md accidentally overwritten | Medium | Identify exact line boundaries of SDD sections before editing. Diff after merge. |
| Upstream skill paths reference `skills/_shared/` but installed path is `~/.claude/skills/_shared/` | Low | Upstream v3.3.2 inlines engram instructions as primary, `_shared/` is fallback reference. Also fixing setup.sh to sync `_shared/`. |
| Existing openspec changes lack state.yaml | Low | Skills designed to handle missing state.yaml as initial state. No migration needed. |
| `auto` mode resolution in AGENTS.md needs hybrid addition | Low | Add `hybrid` as step 1.5: "If engram available AND user wants file artifacts, use hybrid." Clear logic. |
| setup.sh `_shared/` sync may conflict with future upstream changes to setup patterns | Low | Simple additive change — one extra copy block. Easy to maintain. |

## Rollback Plan

All 16 files are tracked in git. Rollback is:
```bash
git checkout HEAD~1 -- GentlemanClaude/skills/ AGENTS.md GentlemanClaude/CLAUDE.md skills/setup.sh
```
For partial rollback, each phase is independently revertable since shared conventions are backwards-compatible (skills reference `_shared/` as optional context, not hard dependencies).

## Dependencies

- Upstream content must be fetchable from `raw.githubusercontent.com/Gentleman-Programming/agent-teams-lite/main/`
- No external tooling dependencies — all changes are file replacements/edits

## Success Criteria

- [ ] All 9 SDD skills at v3.3.2 content (version in frontmatter updated)
- [ ] All 3 shared conventions match upstream v3.3.2
- [ ] `skill-registry/SKILL.md` exists and matches upstream
- [ ] `CLAUDE.md` preserves personality, domain routing, skill tables, plugin detection
- [ ] `CLAUDE.md` SDD orchestrator section has delegation rules, anti-patterns, task escalation
- [ ] `AGENTS.md` preserves project overview, auto-invoke tables, Javi.Dots-specific content
- [ ] `AGENTS.md` skill table includes skill-registry entry
- [ ] `AGENTS.md` auto-mode resolution includes hybrid
- [ ] `skills/setup.sh` copies `_shared/` to `~/.claude/skills/_shared/`
- [ ] Existing openspec changes and archives remain intact
- [ ] `git diff` shows no unintended deletions of Javi.Dots content
