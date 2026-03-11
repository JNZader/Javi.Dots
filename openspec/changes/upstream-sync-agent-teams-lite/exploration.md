# Exploration: Upstream Sync with agent-teams-lite v3.3.2 and engram v1.9.4

## Current State

Javi.Dots SDD skills are at **v2.0**, matching upstream agent-teams-lite ~v2.x. The current setup:

- **9 SDD skills** in `GentlemanClaude/skills/sdd-*/SKILL.md` (init, explore, propose, spec, design, tasks, apply, verify, archive) — all version `"2.0"`
- **3 shared conventions** in `GentlemanClaude/skills/_shared/` (persistence-contract.md, engram-convention.md, openspec-convention.md)
- **No `skill-registry` skill** — doesn't exist at all
- **No `.atl/` directory** — no local skill registry fallback
- **No `hybrid` mode** — persistence-contract.md only knows `engram | openspec | none`
- **No inline engram instructions** in skills — all 9 skills defer to `skills/_shared/engram-convention.md` via "Read and follow" directives
- **No mandatory persist steps** — skills end with "Return Structured Analysis" not "Persist Artifact then Return"
- **No state.yaml** — openspec-convention.md doesn't mention it
- **No sub-agent context rules** — persistence-contract.md has no who-reads/who-writes table
- **No delegation rules or anti-patterns** in CLAUDE.md orchestrator section
- **No Task Escalation table** in CLAUDE.md
- **`_shared/` is NOT synced** to `~/.claude/skills/` by `setup.sh` (the sync loop only copies directories containing `SKILL.md`) — this is actually the ROOT CAUSE of the v3.2.3 inline engram fix

### Javi.Dots-Specific Customizations (things we MUST preserve)

1. **CLAUDE.md personality & tone** — Senior Architect persona, Rioplatense Spanish, Iron Man/Jarvis analogies (lines 1-43)
2. **Domain Routing / Sub-agent orchestration** (lines 45-81) — 6 domain orchestrators (development, infrastructure, data-ai, quality, business, workflow) — NOT in upstream
3. **Framework/Library Detection table** (lines 84-117) — 17+ skill mappings for React, Next.js, TypeScript, etc. — NOT in upstream
4. **Plugin Detection** (lines 119-128) — merge-checks, trim-md, mermaid plugins — NOT in upstream
5. **SDD Triggers section** (lines 168-176) — Javi.Dots specific trigger keywords
6. **AGENTS.md header** (lines 1-117) — Repository skills table, project overview, auto-invoke skills — Javi.Dots specific
7. **AGENTS.md `auto` mode** in Artifact Store Policy (line 136) — upstream CLAUDE.md has `engram | openspec | none`, AGENTS.md adds `auto` with resolution priority
8. **Multi-round-synthesis skill** references Javi.Dots by name (line 558)
9. **Adversarial-review skill** mentions Javi.Dots (line 1321)
10. **openspec/ directory** already has real archived changes and active changes — must not break

## Affected Areas

### Files to CREATE (new):
- `GentlemanClaude/skills/skill-registry/SKILL.md` — New skill from upstream v3.3.1

### Files to UPDATE (SDD skills — version bump 2.0 → 3.3):
- `GentlemanClaude/skills/sdd-init/SKILL.md` — Add Step 4 (Build Skill Registry), hybrid mode, Step 5 (Persist), state.yaml support
- `GentlemanClaude/skills/sdd-explore/SKILL.md` — Add Step 1 (Load Skill Registry), inline engram, hybrid mode, mandatory persist step
- `GentlemanClaude/skills/sdd-propose/SKILL.md` — Same: registry, inline engram, hybrid, persist
- `GentlemanClaude/skills/sdd-spec/SKILL.md` — Same pattern
- `GentlemanClaude/skills/sdd-design/SKILL.md` — Same pattern
- `GentlemanClaude/skills/sdd-tasks/SKILL.md` — Same pattern
- `GentlemanClaude/skills/sdd-apply/SKILL.md` — Registry, inline engram, hybrid, mandatory persist step, TDD workflow from config
- `GentlemanClaude/skills/sdd-verify/SKILL.md` — Registry, inline engram, hybrid, persist, Spec Compliance Matrix
- `GentlemanClaude/skills/sdd-archive/SKILL.md` — Registry, inline engram, hybrid, persist

### Files to UPDATE (shared conventions):
- `GentlemanClaude/skills/_shared/persistence-contract.md` — Add `hybrid` mode row, Sub-Agent Context Rules section, Skill Registry section, State Persistence table
- `GentlemanClaude/skills/_shared/engram-convention.md` — Add NOTE about inline instructions being primary, state artifact section, mem_update section improvements
- `GentlemanClaude/skills/_shared/openspec-convention.md` — Add state.yaml, archive structure improvements

### Files to UPDATE (orchestrator / instructions):
- `GentlemanClaude/CLAUDE.md` — Add Delegation Rules section, Anti-patterns section, Task Escalation table, skill registry loading in Sub-Agent Launching Pattern, Engram Topic Key Format table, hybrid mode to Artifact Store Policy
- `AGENTS.md` — Mirror the same orchestrator changes from CLAUDE.md (this file is the "single source of truth")

### Files to UPDATE (infrastructure):
- `skills/setup.sh` — Must sync `_shared/` directory to `~/.claude/skills/_shared/` (currently skipped because no SKILL.md)

## Approaches

### 1. **Wholesale replacement from upstream** — Pull upstream skill files as-is, then re-add Javi.Dots customizations
   - Pros: Guaranteed feature parity with upstream, easy to track future diffs
   - Cons: Risk of losing Javi.Dots customizations if not careful, CLAUDE.md has significant Javi.Dots-only sections that upstream doesn't have
   - Effort: **Medium** — The 9 SDD skills + 3 shared conventions can be replaced wholesale since they have NO Javi.Dots-specific content. CLAUDE.md and AGENTS.md need surgical merge.

### 2. **Incremental patch** — Apply changes version-by-version (v3.1.0, v3.2.0, v3.2.3, v3.3.0, v3.3.1, v3.3.2)
   - Pros: Granular control, easy to test each version's changes
   - Cons: More work, 6 separate rounds of changes, harder to maintain
   - Effort: **High** — Unnecessary since there are no intermediate consumers

### 3. **Hybrid: Replace SDD skills, surgically merge orchestrator** (RECOMMENDED)
   - Pros: SDD skills are pure upstream with zero Javi.Dots customizations — wholesale replacement is safe. Orchestrator files (CLAUDE.md, AGENTS.md) need surgical merge to preserve personality, domain routing, skill tables, etc.
   - Cons: Need to verify upstream skill files don't reference paths or structures that differ from Javi.Dots
   - Effort: **Medium-Low** — Most work is copy-paste for skills, careful merge for 2 orchestrator files

## Recommendation

**Approach 3: Hybrid (replace skills, merge orchestrator)**.

Rationale:
1. All 9 SDD skill files are 100% upstream content with zero Javi.Dots customization (verified by grep — no mentions of Javi.Dots, no repo-specific paths, no custom steps)
2. All 3 shared convention files are 100% upstream content (verified)
3. The new `skill-registry` skill is a clean addition
4. Only CLAUDE.md and AGENTS.md have Javi.Dots-specific sections that need merge, not replacement
5. `setup.sh` needs a small fix to sync `_shared/` directory

### Implementation Order:
1. Update `_shared/` conventions (3 files) — foundation that skills reference
2. Update all 9 SDD skills (wholesale replace from upstream v3.3.2)
3. Create `skill-registry/SKILL.md` (new skill from upstream v3.3.1)
4. Merge orchestrator additions into CLAUDE.md (preserve Javi.Dots sections)
5. Merge orchestrator additions into AGENTS.md (keep in sync with CLAUDE.md)
6. Fix `setup.sh` to sync `_shared/` directory
7. Update version references and AGENTS.md skill table

## Risks

### R1: `_shared/` sync gap (CRITICAL)
The `setup.sh` loop only copies directories with `SKILL.md`. The `_shared/` directory has `.md` files but no `SKILL.md`. When skills are installed to `~/.claude/skills/`, the `_shared/` conventions are missing. This is exactly why upstream v3.2.3 inlined engram instructions into each skill — sub-agents couldn't reach the shared files. Our update MUST either:
- (a) Add `_shared/` sync to `setup.sh`, OR
- (b) Accept that inline instructions in each skill are the fix (upstream's approach)
- **Recommendation**: Do BOTH — sync `_shared/` for reference AND keep inline instructions as upstream does.

### R2: `auto` mode in AGENTS.md vs upstream's explicit modes
AGENTS.md uses `auto` as default mode with a resolution chain. Upstream doesn't have `auto`. The `hybrid` mode addition doesn't conflict, but we need to update the `auto` resolution chain to include `hybrid` as an option. Specifically: if user has both engram AND wants file artifacts, `auto` could resolve to `hybrid`.

### R3: Existing openspec changes may need state.yaml
We have 4 active changes and 1 archived change in `openspec/`. The upstream v3.3.0 adds `state.yaml` per change. Existing changes don't have this file. The skills should handle missing state.yaml gracefully (treat as initial state).

### R4: Engram v1.9.4 compatibility
The exploration task mentions engram v1.9.4, but the skill changes are about agent-teams-lite. We should verify if engram v1.9.4 has API changes that affect `mem_search`/`mem_save`/`mem_get_observation`/`mem_update` calls in the inlined instructions. If API is unchanged, no risk.

### R5: Non-SDD sub-agents saving discoveries
Upstream v3.3.0 says "Non-SDD sub-agents must save discoveries to engram." This affects the Domain Routing pattern in CLAUDE.md — domain orchestrators (development, infrastructure, etc.) would need engram awareness. This is a broader change that may warrant a separate follow-up.

## Ready for Proposal

**Yes** — The scope is well-defined:
- 15 files to update/create (9 skills + 3 shared + 1 new skill + 2 orchestrator files)
- 1 infrastructure fix (setup.sh)
- Clear approach: wholesale replace for skills, surgical merge for orchestrators
- Low risk of breaking existing functionality since SDD skills are pure upstream content
- The orchestrator should tell the user: "Ready to proceed with /sdd-new upstream-sync-agent-teams-lite. The change involves updating 16 files to sync with upstream v3.3.2, with a focus on hybrid mode support, inline engram instructions, skill registry, and mandatory persist steps."
