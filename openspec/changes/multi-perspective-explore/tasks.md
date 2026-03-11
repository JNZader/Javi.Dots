# Tasks: Multi-Perspective Explore for SDD

## Phase 1: Skill Foundation

- [x] 1.1 Modify `GentlemanClaude/skills/sdd-explore/SKILL.md` — Add `perspective` to the "What You Receive" section: `- Perspective (optional) — analytical lens to focus through (e.g., "architecture", "risk", "testing", "dx"). When present, constrain ALL analysis to this perspective. When absent, explore generally.`
- [x] 1.2 Modify `GentlemanClaude/skills/sdd-explore/SKILL.md` — Update Step 2 (Understand the Request) to include: "If a `perspective` was provided, frame your entire analysis through that lens. Prioritize findings relevant to that perspective. You are a specialist, not a generalist."
- [x] 1.3 Modify `GentlemanClaude/skills/sdd-explore/SKILL.md` — Update Step 6 (Return Structured Analysis) to include `perspective: {name}` in the returned metadata when a perspective was provided. Add a note that the exploration title should read `## Exploration: {topic} (Perspective: {name})` when perspective-constrained.
- [x] 1.4 Modify `GentlemanClaude/skills/sdd-explore/SKILL.md` — Bump `metadata.version` from `"2.0"` to `"2.1"` in the frontmatter.

## Phase 2: Orchestrator Logic

- [x] 2.1 Modify `GentlemanClaude/CLAUDE.md` — Add a `### Multi-Perspective Explore` subsection inside the `## Spec-Driven Development (SDD) Orchestrator` section (after the existing `/sdd-explore` command mapping). Include: trigger keyword list (`"explore deeply"`, `"multi-perspective"`, `"analizar a fondo"`, `"explore from all angles"`, or config `explore.mode: deep`), default perspectives (architecture, risk, testing, dx), and max cap (4).
- [x] 2.2 Modify `GentlemanClaude/CLAUDE.md` — In the same `### Multi-Perspective Explore` section, add the fan-out dispatch pattern: "When multi-perspective is triggered, launch N parallel `sdd-explore` sub-agents in a SINGLE message via Task tool. Each Task call passes `perspective: X` along with the standard context (project, change name, topic, artifact store mode)."
- [x] 2.3 Modify `GentlemanClaude/CLAUDE.md` — In the same section, add the synthesis sub-agent prompt template. The template instructs the synthesis agent to: (1) read all N perspective explorations, (2) identify agreements, conflicts, and blind spots, (3) produce one merged `exploration.md` following standard format plus a `### Perspectives` section, (4) persist the merged artifact via the active store mode.
- [x] 2.4 Modify `GentlemanClaude/CLAUDE.md` — Add config reading instruction: "If `openspec/config.yaml` exists and contains an `explore` section, use `explore.mode` and `explore.perspectives` to determine dispatch mode and perspective list. Config perspectives override defaults."
- [x] 2.5 Modify `AGENTS.md` — Mirror the complete `### Multi-Perspective Explore` section from `GentlemanClaude/CLAUDE.md` (tasks 2.1-2.4) into the SDD Orchestrator section of `AGENTS.md`, keeping content identical.

## Phase 3: Verification

- [ ] 3.1 Verify backward compatibility: Read the updated `sdd-explore/SKILL.md` and confirm that the 6-step process is unchanged when no `perspective` is provided. The skill MUST produce identical output for standard (non-perspective) calls.
- [ ] 3.2 Verify trigger detection: Read the updated `GentlemanClaude/CLAUDE.md` and confirm that the orchestrator's standard `/sdd:explore` path (single agent dispatch) is still the default. Multi-perspective ONLY activates on explicit keywords or config.
- [ ] 3.3 Verify fan-out pattern: Confirm the orchestrator instructions specify parallel Task calls (single message with multiple Task tool invocations), not sequential launches.
- [ ] 3.4 Verify synthesis delegation: Confirm the orchestrator instructions delegate synthesis to a sub-agent via Task tool, NOT inline. The orchestrator must NOT read exploration outputs or merge them itself.
- [ ] 3.5 Verify CLAUDE.md and AGENTS.md consistency: Diff the Multi-Perspective Explore sections in both files and confirm they are equivalent in intent and instruction.
- [ ] 3.6 Verify spec scenarios: Walk through each scenario from `openspec/changes/multi-perspective-explore/specs/explore/spec.md` against the modified files and confirm all are satisfiable.

## Phase 4: Cleanup

- [ ] 4.1 Verify no unintended changes to other SDD skills (propose, spec, design, tasks, apply, verify, archive) by checking their SKILL.md files were not modified.
- [ ] 4.2 Verify the `sdd-explore` skill version bump is present (`"2.0"` → `"2.1"`).
