# Delta for Explore

## ADDED Requirements

### Requirement: Perspective-Constrained Exploration

The `sdd-explore` skill MUST accept an optional `perspective` input parameter. When `perspective` is provided, the sub-agent MUST focus its entire analysis through that specific analytical lens. When `perspective` is absent, the skill MUST behave identically to its current general-purpose exploration.

#### Scenario: Exploration with architecture perspective

- GIVEN the orchestrator launches `sdd-explore` with `perspective: architecture`
- WHEN the sub-agent executes its exploration
- THEN the analysis focuses on patterns, coupling, modularity, code organization, and existing conventions
- AND the output's recommendation section is framed in architectural terms
- AND the returned summary includes `perspective: architecture` in its metadata

#### Scenario: Exploration with risk perspective

- GIVEN the orchestrator launches `sdd-explore` with `perspective: risk`
- WHEN the sub-agent executes its exploration
- THEN the analysis focuses on breaking changes, migration risks, backward compatibility, and effort estimation
- AND the Risks section is expanded with detailed risk assessment
- AND the returned summary includes `perspective: risk` in its metadata

#### Scenario: Exploration with testing perspective

- GIVEN the orchestrator launches `sdd-explore` with `perspective: testing`
- WHEN the sub-agent executes its exploration
- THEN the analysis focuses on testability, coverage strategy, regression risk, and CI impact
- AND the Approaches comparison evaluates each option's testing implications
- AND the returned summary includes `perspective: testing` in its metadata

#### Scenario: Exploration with dx perspective

- GIVEN the orchestrator launches `sdd-explore` with `perspective: dx`
- WHEN the sub-agent executes its exploration
- THEN the analysis focuses on ergonomics, discoverability, learning curve, and documentation needs
- AND the Approaches comparison evaluates each option's developer experience impact
- AND the returned summary includes `perspective: dx` in its metadata

#### Scenario: Exploration without perspective (backward compatibility)

- GIVEN the orchestrator launches `sdd-explore` without a `perspective` parameter
- WHEN the sub-agent executes its exploration
- THEN the behavior is identical to the current general-purpose exploration
- AND no perspective metadata is included in the returned summary
- AND the full 6-step process executes unchanged

### Requirement: Multi-Perspective Orchestrator Dispatch

The SDD orchestrator MUST support a multi-perspective dispatch mode for `/sdd:explore` that fans out N parallel explore sub-agents, each with a different perspective, then synthesizes results.

#### Scenario: Trigger via explicit keyword

- GIVEN the user invokes `/sdd:explore <topic>` with keywords "explore deeply", "multi-perspective", "analizar a fondo", or "explore from all angles"
- WHEN the orchestrator processes the command
- THEN the orchestrator activates multi-perspective mode
- AND launches N parallel `sdd-explore` sub-agents (one per perspective)
- AND each sub-agent receives `perspective: <name>` in its prompt

#### Scenario: Trigger via config

- GIVEN the project has an `openspec/config.yaml` with `explore.mode: deep`
- AND the user invokes `/sdd:explore <topic>` without explicit keywords
- WHEN the orchestrator processes the command
- THEN the orchestrator activates multi-perspective mode
- AND uses the perspectives listed in `explore.perspectives` from config

#### Scenario: Standard mode when no trigger present

- GIVEN no explicit multi-perspective keywords are used
- AND no config sets `explore.mode: deep`
- WHEN the user invokes `/sdd:explore <topic>`
- THEN the orchestrator dispatches a single `sdd-explore` sub-agent as it does today
- AND no fan-out or synthesis occurs

#### Scenario: Default perspectives when config has no overrides

- GIVEN multi-perspective mode is triggered
- AND no `explore.perspectives` section exists in config
- WHEN the orchestrator resolves perspectives
- THEN it uses the defaults: architecture, risk, testing, dx (4 perspectives)

#### Scenario: Custom perspectives from config

- GIVEN multi-perspective mode is triggered
- AND `openspec/config.yaml` contains:
  ```yaml
  explore:
    mode: deep
    perspectives:
      - architecture
      - security
      - performance
  ```
- WHEN the orchestrator resolves perspectives
- THEN it uses only the 3 perspectives from config: architecture, security, performance
- AND does NOT include the defaults that are not listed

#### Scenario: Maximum perspective cap

- GIVEN multi-perspective mode is triggered
- AND config lists more than 4 perspectives
- WHEN the orchestrator resolves perspectives
- THEN it SHOULD use only the first 4 perspectives
- AND log/warn that the remaining perspectives were skipped due to the cap

### Requirement: Parallel Fan-Out Execution

The orchestrator MUST launch all perspective sub-agents in parallel (not sequentially) to maximize efficiency and ensure independent analysis.

#### Scenario: Parallel launch of perspective agents

- GIVEN multi-perspective mode is active with 3 perspectives (architecture, risk, dx)
- WHEN the orchestrator dispatches sub-agents
- THEN all 3 `sdd-explore` Task calls are made in a single message (parallel execution)
- AND each Task call includes the topic, change name, artifact store mode, and its assigned perspective

#### Scenario: Independent perspective analysis

- GIVEN 3 perspective sub-agents are running in parallel
- WHEN each sub-agent produces its exploration
- THEN no sub-agent sees the output of any other sub-agent
- AND each exploration is based solely on codebase investigation from its assigned perspective

### Requirement: Synthesis Sub-Agent

After all perspective explorations complete, the orchestrator MUST launch a synthesis sub-agent that merges N perspective explorations into one unified `exploration.md`.

#### Scenario: Synthesis receives all perspective outputs

- GIVEN 3 perspective sub-agents have completed with their exploration reports
- WHEN the orchestrator launches the synthesis sub-agent
- THEN the synthesis prompt includes the full text of all 3 perspective explorations
- AND the synthesis prompt instructs the agent to identify agreements, conflicts, and blind spots

#### Scenario: Synthesis output format

- GIVEN the synthesis sub-agent receives N perspective explorations
- WHEN it produces the merged exploration
- THEN the output follows the standard `exploration.md` format (Current State, Affected Areas, Approaches, Recommendation, Risks, Ready for Proposal)
- AND includes an additional `### Perspectives` section summarizing what each perspective contributed
- AND the Recommendation section reflects the weight of multiple perspectives

#### Scenario: Synthesis identifies conflicts

- GIVEN perspective A recommends Approach X and perspective B recommends Approach Y
- WHEN the synthesis sub-agent merges the explorations
- THEN the `### Perspectives` section explicitly notes the disagreement
- AND the Recommendation section explains how the conflict was resolved (which approach and why)

#### Scenario: Synthesis identifies blind spots

- GIVEN perspective A identifies a concern not mentioned by any other perspective
- WHEN the synthesis sub-agent merges the explorations
- THEN the concern is included in the merged exploration
- AND is marked as a single-perspective finding in the `### Perspectives` section

#### Scenario: Synthesis persists as standard exploration

- GIVEN the synthesis sub-agent produces a merged exploration
- WHEN the artifact is persisted (to openspec or engram)
- THEN it is saved as the standard `exploration.md` for the change
- AND downstream phases (propose, spec) can consume it without knowing it came from multi-perspective

### Requirement: Orchestrator Stays Lightweight

The orchestrator MUST NOT perform synthesis, analysis, or exploration work inline. All heavy work MUST be delegated to sub-agents.

#### Scenario: Orchestrator only dispatches

- GIVEN multi-perspective mode is active
- WHEN the orchestrator handles the full flow
- THEN the orchestrator only: (1) detects trigger, (2) resolves perspectives, (3) launches parallel sub-agents, (4) launches synthesis sub-agent, (5) presents final summary to user
- AND the orchestrator does NOT read source code, analyze approaches, or merge explorations itself

## MODIFIED Requirements

(None — this change adds new behavior without modifying existing requirements.)

## REMOVED Requirements

(None — no existing behavior is removed.)
