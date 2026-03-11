# Proposal: Multi-Perspective Explore for SDD

## Intent

The current `/sdd:explore` flow runs a single sub-agent with a single generalist perspective. For complex, cross-cutting changes this produces shallow analysis — one agent trying to cover architecture, risk, testing, and DX in a single pass lacks the depth that independent, focused perspectives provide. We need an opt-in mechanism that fans out N parallel explore sub-agents (each with a different analytical lens), then synthesizes their findings into one comprehensive exploration.

This follows the same fan-out/fan-in pattern already proven by `adversarial-review` and `multi-round-synthesis` skills, applied to the exploration phase of SDD.

## Scope

### In Scope
- Add optional `perspective` input parameter to `sdd-explore` skill
- Add multi-perspective dispatch logic to the orchestrator sections in `GentlemanClaude/CLAUDE.md` and `AGENTS.md`
- Define trigger keywords for multi-perspective mode (explicit user opt-in)
- Define default perspective set (architecture, risk, testing, dx)
- Add synthesis sub-agent prompt template to orchestrator
- Document optional `explore` config section in openspec convention

### Out of Scope
- Changes to other SDD skills (propose, spec, design, tasks, apply, verify, archive)
- Automatic complexity detection to trigger multi-perspective (future enhancement)
- MCP server or external tool integration
- TUI installer changes
- New standalone skills (no `sdd-explore-deep` skill)
- Multi-round iterative synthesis (single synthesis pass is sufficient for MVP)

## Approach

**Orchestrator-side fan-out with config-driven perspectives** (Approaches 1+4 from exploration).

The orchestrator gains a conditional dispatch path for `/sdd:explore`:

1. **Trigger detection**: User says "explore deeply", "multi-perspective", "analizar a fondo", or project config has `explore.mode: deep`
2. **Perspective resolution**: Read perspectives from config (with defaults: architecture, risk, testing, dx; max 4)
3. **Fan-out**: Launch N parallel `sdd-explore` sub-agents via Task tool, each receiving `perspective: X` in their prompt
4. **Synthesis**: Launch a synthesis sub-agent that reads all N exploration outputs and merges them into one `exploration.md` following the standard format but with a `### Perspectives` section
5. **Standard path unchanged**: Without trigger keywords or config, `/sdd:explore` dispatches exactly as today (single agent, no perspective)

The `sdd-explore` skill gets a minor addition: when `perspective` is present in the prompt, the agent focuses its analysis through that lens. When absent, behavior is identical to current.

## Affected Areas

| Area | Impact | Description |
|------|--------|-------------|
| `GentlemanClaude/skills/sdd-explore/SKILL.md` | Modified | Add optional `perspective` input parameter to "What You Receive"; adjust output format to note perspective when present |
| `GentlemanClaude/CLAUDE.md` | Modified | Add `### Multi-Perspective Explore` section to SDD orchestrator with trigger detection, fan-out logic, and synthesis prompt template |
| `AGENTS.md` | Modified | Mirror the Multi-Perspective Explore section from CLAUDE.md |

## Risks

| Risk | Likelihood | Mitigation |
|------|------------|------------|
| Context window limits — synthesis agent receives N exploration reports (potentially 500+ lines each) | Medium | Cap at 4 perspectives by default; instruct perspective agents to keep output concise (<300 lines) |
| Token cost — multi-perspective is 3-5x the cost of single explore | Low | Strictly opt-in via explicit keywords or config; standard single-agent explore remains default |
| Orchestrator complexity creep — adding fan-out logic adds weight to the always-loaded orchestrator prompt | Low | Keep additions to ~30 lines; synthesis is delegated to sub-agent, not done inline |
| Perspective quality — shallow analysis if perspective prompts aren't well-crafted | Medium | Use adversarial-review's perspective prompt patterns as quality reference; iterate on prompts |
| Upstream divergence — multi-perspective is Javi.Dots-specific, not in upstream `agent-teams-lite` | Low | Document clearly as project-specific addition; keep changes isolated to clearly marked sections |

## Rollback Plan

All changes are to markdown prompt files (no code logic). Rollback is:
1. Revert `sdd-explore/SKILL.md` to remove `perspective` input (or leave it — it's backward-compatible since it's optional)
2. Remove the `### Multi-Perspective Explore` section from `GentlemanClaude/CLAUDE.md`
3. Remove the same section from `AGENTS.md`

Since these are prompt-only changes with no runtime code, rollback is a simple git revert with zero risk of breaking existing behavior.

## Dependencies

- None. All affected files already exist and are under our control. No external dependencies, no new packages, no infrastructure changes.

## Success Criteria

- [ ] Running `/sdd:explore <topic>` without trigger keywords produces the same single-agent exploration as before (backward compatibility)
- [ ] Running `/sdd:explore <topic>` with "explore deeply" or "multi-perspective" triggers parallel fan-out of N explore sub-agents
- [ ] Each perspective sub-agent produces a focused exploration through its assigned lens
- [ ] A synthesis sub-agent merges N perspective explorations into one `exploration.md`
- [ ] The synthesized exploration identifies agreements, conflicts, and blind spots across perspectives
- [ ] Default perspectives (architecture, risk, testing, dx) are used when no config overrides exist
- [ ] Perspectives can be customized via project config (when config section exists)
- [ ] The orchestrator stays lightweight — fan-out is dispatch logic only, synthesis is delegated
