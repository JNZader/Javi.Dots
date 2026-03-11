# Tasks: Upstream Sync with agent-teams-lite v3.3.2

## Phase 1: Shared Conventions (3 files)

- [ ] 1.1 Fetch upstream `skills/_shared/persistence-contract.md` from `raw.githubusercontent.com/Gentleman-Programming/agent-teams-lite/main/skills/_shared/persistence-contract.md` and replace `GentlemanClaude/skills/_shared/persistence-contract.md` wholesale
- [ ] 1.2 Fetch upstream `skills/_shared/engram-convention.md` from `raw.githubusercontent.com/Gentleman-Programming/agent-teams-lite/main/skills/_shared/engram-convention.md` and replace `GentlemanClaude/skills/_shared/engram-convention.md` wholesale
- [ ] 1.3 Fetch upstream `skills/_shared/openspec-convention.md` from `raw.githubusercontent.com/Gentleman-Programming/agent-teams-lite/main/skills/_shared/openspec-convention.md` and replace `GentlemanClaude/skills/_shared/openspec-convention.md` wholesale

## Phase 2: SDD Skills Wholesale Replacement (9 files) + New Skill Registry (1 file)

- [ ] 2.1 Fetch upstream `skills/sdd-init/SKILL.md` from `raw.githubusercontent.com/Gentleman-Programming/agent-teams-lite/main/skills/sdd-init/SKILL.md` and replace `GentlemanClaude/skills/sdd-init/SKILL.md`
- [ ] 2.2 Fetch upstream `skills/sdd-explore/SKILL.md` and replace `GentlemanClaude/skills/sdd-explore/SKILL.md`
- [ ] 2.3 Fetch upstream `skills/sdd-propose/SKILL.md` and replace `GentlemanClaude/skills/sdd-propose/SKILL.md`
- [ ] 2.4 Fetch upstream `skills/sdd-spec/SKILL.md` and replace `GentlemanClaude/skills/sdd-spec/SKILL.md`
- [ ] 2.5 Fetch upstream `skills/sdd-design/SKILL.md` and replace `GentlemanClaude/skills/sdd-design/SKILL.md`
- [ ] 2.6 Fetch upstream `skills/sdd-tasks/SKILL.md` and replace `GentlemanClaude/skills/sdd-tasks/SKILL.md`
- [ ] 2.7 Fetch upstream `skills/sdd-apply/SKILL.md` and replace `GentlemanClaude/skills/sdd-apply/SKILL.md`
- [ ] 2.8 Fetch upstream `skills/sdd-verify/SKILL.md` and replace `GentlemanClaude/skills/sdd-verify/SKILL.md`
- [ ] 2.9 Fetch upstream `skills/sdd-archive/SKILL.md` and replace `GentlemanClaude/skills/sdd-archive/SKILL.md`
- [ ] 2.10 Create directory `GentlemanClaude/skills/skill-registry/` and fetch upstream `skills/skill-registry/SKILL.md` from `raw.githubusercontent.com/Gentleman-Programming/agent-teams-lite/main/skills/skill-registry/SKILL.md`

## Phase 3: Orchestrator Merges (2 files)

- [ ] 3.1 Read `GentlemanClaude/CLAUDE.md` and identify the exact boundaries of the "SDD Orchestrator" section (from `## Spec-Driven Development (SDD) Orchestrator` to end of file or next `##` section)
- [ ] 3.2 Fetch upstream `CLAUDE.md` SDD orchestrator section from `raw.githubusercontent.com/Gentleman-Programming/agent-teams-lite/main/CLAUDE.md`
- [ ] 3.3 Merge into `GentlemanClaude/CLAUDE.md`: replace ONLY the SDD Orchestrator section with upstream content. Preserve: Identity Inheritance subsection (Javi.Dots personality reference), all content above the SDD section (personality, domain routing, skill tables, plugin detection, SDD triggers)
- [ ] 3.4 Verify `GentlemanClaude/CLAUDE.md` post-merge: confirm personality section intact, domain routing intact, skill tables intact, new delegation rules present, anti-patterns present, task escalation table present, hybrid mode in Artifact Store Policy
- [ ] 3.5 Read `AGENTS.md` and identify the exact boundaries of the SDD Orchestrator section
- [ ] 3.6 Merge into `AGENTS.md`: replace ONLY the SDD Orchestrator section content with upstream patterns (delegation rules, anti-patterns, task escalation). Preserve: repository skills table, project overview, auto-invoke tables, all Javi.Dots-specific content
- [ ] 3.7 Add `skill-registry` row to the "Generic Skills" table in `AGENTS.md`: `| skill-registry | Skill dependency resolution and registry loading | [SKILL.md](GentlemanClaude/skills/skill-registry/SKILL.md) |`
- [ ] 3.8 Update `AGENTS.md` Artifact Store Policy `auto` resolution chain: add `hybrid` as option — "If engram available AND user wants file artifacts, use `hybrid`" between engram-only and openspec steps
- [ ] 3.9 Verify `AGENTS.md` post-merge: confirm repository skills table intact, auto-invoke section intact, new delegation rules present, skill-registry in table, hybrid in auto-resolution

## Phase 4: Infrastructure (1 file)

- [ ] 4.1 Read `skills/setup.sh` and identify the skill sync loop that copies directories to `~/.claude/skills/`
- [ ] 4.2 Add a block to `skills/setup.sh` that copies `GentlemanClaude/skills/_shared/` to `~/.claude/skills/_shared/` (before or after the main skill loop). Use `cp -r` or `rsync` matching existing script patterns
- [ ] 4.3 Verify `skills/setup.sh` by dry-run or reading: confirm `_shared/` sync logic exists, confirm existing skill copy loop unchanged

## Phase 5: Validation

- [ ] 5.1 Run `git diff --stat` to confirm exactly 16 files changed (3 shared + 9 skills + 1 new + 2 orchestrator + 1 setup.sh)
- [ ] 5.2 Verify no Javi.Dots-specific content was lost: grep for "Senior Architect", "Rioplatense", "Domain Routing", "Iron Man", "auto-invoke" in CLAUDE.md and AGENTS.md
- [ ] 5.3 Verify all 9 SDD skills have `metadata.version: "3.3"` (or matching upstream version) in frontmatter
- [ ] 5.4 Verify `skill-registry/SKILL.md` exists and has valid frontmatter
- [ ] 5.5 Verify `_shared/` files contain hybrid mode, sub-agent context rules, state.yaml references
- [ ] 5.6 Verify existing openspec changes and archives are untouched
