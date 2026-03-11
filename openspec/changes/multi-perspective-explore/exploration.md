# Exploration: Multi-Perspective Explore for SDD

## Current State

### How `/sdd:explore` Works Today

The explore flow is a **single-agent, single-pass** operation:

1. **Orchestrator** (in `GentlemanClaude/CLAUDE.md` or `AGENTS.md`) receives `/sdd:explore <topic>`
2. Orchestrator launches **ONE** sub-agent via Task tool with the `sdd-explore` skill
3. Sub-agent reads `GentlemanClaude/skills/sdd-explore/SKILL.md`, follows its 6-step process:
   - Step 1: Load Skill Registry
   - Step 2: Understand the Request
   - Step 3: Investigate the Codebase (read files, search patterns)
   - Step 4: Analyze Options (compare approaches in a table)
   - Step 5: Persist Artifact (write `exploration.md` to openspec or engram)
   - Step 6: Return Structured Analysis (status, executive_summary, artifacts, risks, next_recommended)
4. Orchestrator receives the result, presents summary to user, asks if they want to proceed to `/sdd:propose`

**Key constraint**: The sub-agent carries a single generalist perspective. It analyzes architecture, risks, testing implications, and DX all at once in one pass. There is no mechanism for the orchestrator to decompose the exploration topic into specialized angles or run parallel investigations.

### Orchestrator Dispatch Logic

From `GentlemanClaude/CLAUDE.md:217-228`:
- `/sdd-explore` maps to the `sdd-explore` skill — ONE skill, ONE sub-agent
- The orchestrator's Sub-Agent Launching Pattern (line 281-302) shows a single Task call per phase
- No provision for parallel Task calls within a single SDD phase
- The orchestrator rules (line 242-252) say "ONLY track state, summarize, ask for approval, launch sub-agents" — but nothing prevents launching MULTIPLE sub-agents for the same phase

From `AGENTS.md:173-191`:
- Same command mapping, same single-dispatch pattern
- The `auto` artifact store mode resolution is the only AGENTS.md-specific addition

### Existing Multi-Agent Patterns in the Codebase

Two skills already implement fan-out/fan-in:

1. **`adversarial-review`** (`~/.config/opencode/skills/quality/adversarial-review/SKILL.md`):
   - 3 fixed perspectives (Security, Quality, Test) run in parallel on the same diff
   - A Synthesizer agent merges findings using a consensus algorithm
   - Supports custom perspectives via config
   - **Key design**: Perspectives are INDEPENDENT (they don't see each other's output during review)

2. **`multi-round-synthesis`** (`~/.config/opencode/skills/workflow/multi-round-synthesis/SKILL.md`):
   - Coordinator delegates to N specialists in parallel (Round 1)
   - Reviews responses, detects conflicts and gaps
   - Sends follow-up questions (Round 2+)
   - Synthesizes final answer when satisfied (max 3 rounds)
   - **Key design**: Iterative refinement through multiple rounds, coordinator decides termination

The **Make It Heavy** (Doriandarko) pattern adds a decomposition step:
   - Takes a query, breaks it into N sub-questions from different analytical angles
   - Fans out N independent agents
   - Fans in with synthesis
   - **Key design**: The decomposition itself is an LLM step (the orchestrator or a decomposer agent decides the angles)

## Affected Areas

- `GentlemanClaude/skills/sdd-explore/SKILL.md` — The explore sub-agent skill definition
- `GentlemanClaude/CLAUDE.md` — Orchestrator SDD section (dispatch logic, sub-agent launching pattern)
- `AGENTS.md` — Mirror of orchestrator SDD section
- `GentlemanClaude/skills/_shared/openspec-convention.md` — Artifact paths (if multi-explore creates multiple files)
- `openspec/config.yaml` — Config rules for explore phase (if adding config options)

## Approaches

### 1. **Orchestrator-Side Fan-Out** — The orchestrator decides when and how to decompose

The orchestrator gains a new dispatch mode for `/sdd:explore`. When triggered (by user keyword or automatic complexity detection), it:
1. Decomposes the topic into N perspective prompts (either hardcoded perspectives or LLM-generated)
2. Launches N parallel `sdd-explore` sub-agents via Task tool, each with a perspective constraint added to their prompt
3. Launches a synthesis sub-agent (or does lightweight synthesis itself) to merge N explorations into one `exploration.md`

**What changes**:
- `CLAUDE.md` + `AGENTS.md`: Add multi-perspective dispatch logic to the orchestrator's `/sdd:explore` handling
- `sdd-explore/SKILL.md`: Optionally add a `perspective` input field so the sub-agent knows it's exploring from a specific angle (minor change)
- No new skills needed

**Pros**:
- Cleanest separation of concerns: orchestrator decides WHEN to fan-out, sub-agents just explore
- The `sdd-explore` skill stays simple — it still does one exploration, just with an optional perspective constraint
- Reuses the existing skill without duplication
- The orchestrator already has the pattern for parallel Task calls (it's allowed by the rules, just not used for explore)
- Mirrors how `adversarial-review` works (orchestrator launches N perspectives in parallel)

**Cons**:
- Orchestrator becomes slightly more complex (but still lightweight — it's prompt logic, not code)
- Synthesis step needs definition: who does it and how?
- Perspective definitions need to live somewhere (in the orchestrator? in config?)

**Effort**: Medium

### 2. **Skill-Internal Fan-Out** — The `sdd-explore` sub-agent does multi-perspective internally

The skill itself gains a "deep mode" where it:
1. Analyzes the topic
2. Decides which perspectives to apply
3. Sequentially (it can't launch sub-agents from within a sub-agent) analyzes from each perspective
4. Self-synthesizes into one exploration document

**What changes**:
- `sdd-explore/SKILL.md`: Major rewrite to add internal multi-perspective logic

**Pros**:
- Self-contained — no orchestrator changes needed
- Single artifact output

**Cons**:
- Sub-agents can't launch sub-sub-agents (Task tool limitation) — so perspectives run SEQUENTIALLY in one context, not in parallel
- Bloats the sub-agent's context window significantly (analyzing from 5 perspectives in one agent)
- Defeats the purpose of fan-out (no parallelism, no independent perspectives)
- Makes the skill much more complex and harder to maintain
- Single agent bias: one agent analyzing "from different perspectives" is really just one agent pretending to have different viewpoints — it lacks genuine independence

**Effort**: High

### 3. **New Skill: `sdd-explore-deep`** — A separate skill for multi-perspective exploration

A new skill that IS the synthesis agent. It:
1. Receives the topic from the orchestrator
2. Defines perspectives
3. Launches N `sdd-explore` sub-agents (if sub-agents can launch sub-agents) OR returns instructions to the orchestrator about what sub-agents to launch

**What changes**:
- New `GentlemanClaude/skills/sdd-explore-deep/SKILL.md`
- `CLAUDE.md` + `AGENTS.md`: Add `/sdd:explore-deep` command or auto-routing logic

**Pros**:
- Clean separation: existing skill untouched
- Could serve as a reusable "deep exploration" pattern

**Cons**:
- Sub-agents launching sub-agents is not a supported pattern — the sub-agent would need to return a "fan-out plan" to the orchestrator, adding a round-trip
- Adds another skill to maintain, register, and sync
- Splits the explore concept across two skills unnecessarily
- Skill proliferation without clear value — it's really just orchestrator logic in a skill costume

**Effort**: Medium-High

### 4. **Configuration-Driven Mode in Existing Skill** — A `mode: deep` config option

Add a configuration flag in `openspec/config.yaml` under `rules.explore`:

```yaml
rules:
  explore:
    mode: deep  # "standard" (default) or "deep" (multi-perspective)
    perspectives:
      - architecture
      - risk
      - testing
      - dx
```

The orchestrator reads this config and dispatches accordingly (Approach 1 with config driving the decision).

**What changes**:
- `openspec/config.yaml` schema: Add `explore` rules section
- `CLAUDE.md` + `AGENTS.md`: Orchestrator reads config to decide dispatch mode
- `sdd-explore/SKILL.md`: Minor — accept optional `perspective` parameter
- `_shared/openspec-convention.md`: Document the new config section

**Pros**:
- Project-customizable: different projects can define different perspectives
- No new skills or commands — just configuration
- Follows the existing pattern where `openspec/config.yaml` has per-phase rules
- User doesn't need to remember a new command

**Cons**:
- Config-always-on might be wasteful for simple explorations
- Still needs orchestrator logic to read config and fan out

**Effort**: Medium

## Recommendation

**Approach 1 (Orchestrator-Side Fan-Out) combined with Approach 4 (Config-Driven)**.

Rationale:

1. **The orchestrator is the right place for fan-out logic.** This is coordination, not exploration. The orchestrator already has rules for launching sub-agents, and nothing prevents it from launching multiple. Adding fan-out here follows the existing architecture pattern — it's the same thing `adversarial-review` does, applied to explore.

2. **The `sdd-explore` skill should stay simple.** It's a focused sub-agent that explores one topic from one angle. Adding an optional `perspective` constraint (a one-line addition to "What You Receive") is all it needs. The skill doesn't need to know it's part of a multi-perspective fan-out.

3. **Config drives perspective definitions, not hardcoding.** Different projects have different concerns. A CLI tool project (like Javi.Dots) cares about cross-platform, DX, and performance. A web app cares about security, scalability, and accessibility. Perspectives should be configurable in `openspec/config.yaml`, not hardcoded in the orchestrator.

4. **Trigger should be opt-in with smart defaults.** Multi-perspective should activate when:
   - User explicitly asks: "explore deeply", "multi-perspective explore", "explore from all angles"
   - Config has `explore.mode: deep`
   - Orchestrator detects complexity (cross-cutting concern, architecture change, >3 domains affected) — but this is a stretch goal, not MVP

5. **Synthesis should be a sub-agent, not the orchestrator.** The orchestrator must stay lightweight. A synthesis sub-agent receives N exploration reports and produces one merged `exploration.md`. This sub-agent doesn't need a new skill — it's just a Task call with a synthesis prompt (like the Synthesizer in adversarial-review).

### Recommended Perspective Set for Dev Tool Projects

| Perspective | Focus | When Most Valuable |
|-------------|-------|-------------------|
| **Architecture** | Patterns, coupling, modularity, code organization, existing conventions | Always (core perspective) |
| **Risk & Feasibility** | Breaking changes, migration risks, backward compatibility, effort estimation | Features touching public API or core modules |
| **Testing & Quality** | Testability, coverage strategy, regression risk, CI impact | Changes that add new behavior |
| **User Experience / DX** | Ergonomics, discoverability, learning curve, documentation needs | User-facing features, CLI changes, config changes |
| **Performance & Scale** | Resource usage, startup time, memory footprint, concurrent operations | Infrastructure changes, data processing |

For Javi.Dots specifically, a good default would be **Architecture + Risk + DX** (3 perspectives). Testing and Performance would be opt-in for relevant changes.

### Synthesis Strategy

The synthesis sub-agent should:
1. Receive all N perspective explorations
2. Identify agreements (all perspectives agree on approach X)
3. Identify conflicts (perspective A says approach X, perspective B says approach Y)
4. Identify blind spots (things only one perspective caught that others missed)
5. Produce a unified `exploration.md` that follows the existing format but with a `### Perspectives` section showing the multi-angle analysis
6. Make a single recommendation with the weight of N perspectives behind it

### Implementation Scope

**Changes needed**:

| File | Change Type | Description |
|------|-------------|-------------|
| `GentlemanClaude/CLAUDE.md` | Modify | Add multi-perspective dispatch logic to `/sdd:explore` handling. Add trigger keywords. Add synthesis sub-agent prompt template. |
| `AGENTS.md` | Modify | Mirror orchestrator changes from CLAUDE.md |
| `GentlemanClaude/skills/sdd-explore/SKILL.md` | Minor modify | Add optional `perspective` to "What You Receive" section. Adjust Step 6 output to note perspective if present. |
| `GentlemanClaude/skills/_shared/openspec-convention.md` | Minor modify | Document `exploration.md` can be the result of multi-perspective synthesis |
| `openspec/config.yaml` | Modify | Add `explore` rules section with mode and perspectives |

**What stays the same**:
- The `sdd-explore` skill's core logic (Steps 1-6)
- All other SDD skills
- The openspec directory structure
- The dependency graph (explore is still optional, before proposal)
- Engram artifact naming convention

## Risks

- **R1: Context window limits** — Running 5 perspectives means 5 sub-agent calls. Each produces an exploration report. The synthesis sub-agent receives all 5 reports. If each is 500 lines, that's 2500 lines of input to the synthesizer. Mitigation: Cap at 3-4 perspectives by default; instruct perspective agents to be concise.

- **R2: Token cost** — Multi-perspective is 3-5x the cost of single explore. Users should be aware this is opt-in and more expensive. Mitigation: Default to standard (single) explore. Only fan-out when explicitly requested or configured.

- **R3: Diminishing returns** — Not every exploration benefits from multiple perspectives. A simple "should we upgrade dependency X?" doesn't need architecture, risk, testing, DX, and performance angles. Mitigation: Orchestrator should use judgment (or user keyword detection) to decide when multi-perspective adds value.

- **R4: Orchestrator complexity creep** — The orchestrator is supposed to stay lightweight. Adding fan-out logic, perspective definitions, synthesis prompts, and config reading adds weight. Mitigation: Keep the orchestrator's new logic to ~30 lines of prompt additions. The synthesis is delegated to a sub-agent, so the orchestrator just does dispatch.

- **R5: Upstream divergence** — The SDD skills are synced from upstream `agent-teams-lite`. Adding multi-perspective logic to the orchestrator sections of `CLAUDE.md` and `AGENTS.md` is Javi.Dots-specific and will need to be preserved during future upstream syncs. Mitigation: Document clearly in the exploration that these are Javi.Dots-specific additions (upstream has no multi-perspective concept).

- **R6: Perspective quality** — An LLM asked to "explore from a risk perspective" might produce shallow analysis if the perspective prompt isn't well-crafted. Mitigation: Invest in perspective prompt engineering. Use the adversarial-review skill's perspective prompts as a quality reference.

## Ready for Proposal

**Yes** — The scope is well-defined and the approach is clear:

- Primary change: Orchestrator dispatch logic in `CLAUDE.md` + `AGENTS.md` (~30 lines each)
- Secondary change: Minor `sdd-explore/SKILL.md` update (add `perspective` to inputs)
- Config addition: `openspec/config.yaml` explore rules
- No new skills needed
- Follows existing patterns (adversarial-review fan-out, multi-round-synthesis coordination)

The orchestrator should tell the user: "Ready to proceed with `/sdd:propose multi-perspective-explore`. The change adds opt-in multi-perspective exploration to the SDD workflow, where the orchestrator fans out 3-5 parallel explore sub-agents with different analytical lenses, then synthesizes the results into one comprehensive exploration document."
