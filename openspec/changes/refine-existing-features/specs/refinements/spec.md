# Delta for Refinements

## ADDED Requirements

### Requirement: Custom Explore Perspectives

The orchestrator MUST accept user-defined perspective names in `openspec/config.yaml` under `explore.perspectives`, not limited to the 4 default names.

The system MUST validate that each perspective entry is a non-empty string. The system MUST treat unrecognized names as valid custom perspectives and pass them through to sub-agents without filtering.

#### Scenario: User defines custom perspectives in config

- GIVEN `openspec/config.yaml` contains `explore.perspectives: [architecture, performance, compliance]`
- WHEN multi-perspective explore is triggered
- THEN the orchestrator fans out 3 sub-agents with perspectives `architecture`, `performance`, `compliance`
- AND the default `testing`, `risk`, `dx` perspectives are NOT used

#### Scenario: No perspectives defined in config

- GIVEN `openspec/config.yaml` does NOT contain an `explore.perspectives` key
- WHEN multi-perspective explore is triggered
- THEN the orchestrator uses the 4 defaults: `architecture`, `testing`, `risk`, `dx`

#### Scenario: More than 4 perspectives defined

- GIVEN `openspec/config.yaml` contains `explore.perspectives` with 6 entries
- WHEN multi-perspective explore is triggered
- THEN the orchestrator uses only the first 4 entries
- AND warns the user that 2 perspectives were skipped

---

### Requirement: Multi-Round Explore Iteration

The orchestrator MUST support iterative explore rounds controlled by `explore.rounds` in `openspec/config.yaml`.

The system MUST default `explore.rounds` to 1 (current single-pass behavior). The system MUST cap `explore.rounds` at 3 regardless of config value. The system MUST stop early if the synthesis contains no "NEEDS FURTHER ANALYSIS" items.

#### Scenario: Default single round (backward compatible)

- GIVEN `explore.rounds` is not set OR is set to 1
- WHEN multi-perspective explore completes
- THEN exactly ONE fan-out + synthesis cycle executes
- AND behavior is identical to the current implementation

#### Scenario: Multi-round with convergence in round 2

- GIVEN `explore.rounds` is set to 3
- WHEN round 1 synthesis produces 2 "NEEDS FURTHER ANALYSIS" items
- AND round 2 synthesis produces 0 "NEEDS FURTHER ANALYSIS" items
- THEN round 3 is NOT executed (early convergence)
- AND the final synthesis is from round 2

#### Scenario: Multi-round with max rounds reached

- GIVEN `explore.rounds` is set to 2
- WHEN round 1 synthesis produces 3 "NEEDS FURTHER ANALYSIS" items
- AND round 2 synthesis still produces 1 "NEEDS FURTHER ANALYSIS" item
- THEN exploration stops after round 2 (max reached)
- AND the remaining unresolved item is flagged in the final synthesis

#### Scenario: Round 2+ perspective agents receive prior synthesis

- GIVEN round 1 synthesis has completed
- WHEN round 2 perspective agents are launched
- THEN each agent receives the prior round's synthesis as additional context
- AND each agent is instructed to focus on unresolved conflicts and blind spots identified in that synthesis

#### Scenario: Config value exceeds cap

- GIVEN `explore.rounds` is set to 5
- WHEN multi-perspective explore begins
- THEN the orchestrator caps rounds at 3
- AND warns the user that the configured value exceeds the maximum

---

### Requirement: Agreement Matrix in Synthesis

The synthesis sub-agent MUST produce a structured Agreement Matrix table instead of free-form agree/conflict text.

#### Scenario: All perspectives agree on a finding

- GIVEN 4 perspective agents all identify "Use pattern X" as a recommendation
- WHEN the synthesis agent processes the reports
- THEN the Agreement Matrix shows checkmarks for all 4 perspectives
- AND the Confidence column shows "High"

#### Scenario: Perspectives disagree on an approach

- GIVEN `architecture` and `risk` recommend Approach Y, but `testing` and `dx` oppose it
- WHEN the synthesis agent processes the reports
- THEN the Agreement Matrix row for "Approach Y" shows checkmark for architecture and risk, X for testing and dx
- AND the Confidence column shows "Low — needs resolution"
- AND the synthesis includes a resolution recommendation explaining which view is favored and why

#### Scenario: Finding only raised by one perspective

- GIVEN only the `risk` perspective identifies "Migration risk Z"
- WHEN the synthesis agent processes the reports
- THEN the Agreement Matrix row shows checkmark for risk, dash (—) for others
- AND the Confidence column shows "Single perspective — unvalidated"

---

### Requirement: Hook Scripts for Claude Code

The system MUST provide hook scripts in `GentlemanClaude/hooks/` that enhance the Claude Code developer experience.

Each hook script MUST be a self-contained shell script of 5-15 lines. Each hook script MUST include a comment header describing its purpose and trigger event.

#### Scenario: comment-check hook detects TODO without context

- GIVEN the `comment-check.sh` hook is installed
- WHEN a tool writes code containing `// TODO` without explanatory text after it
- THEN the hook outputs a warning message identifying the file and line

#### Scenario: comment-check hook passes clean code

- GIVEN the `comment-check.sh` hook is installed
- WHEN a tool writes code with no TODO/FIXME/HACK comments
- THEN the hook produces no output (silent pass)

#### Scenario: todo-tracker hook reminds on session end

- GIVEN the `todo-tracker.sh` hook is installed
- WHEN a Stop event fires (session ending)
- THEN the hook outputs a reminder to review and update the todo list

#### Scenario: setup.sh syncs hooks to user config

- GIVEN `GentlemanClaude/hooks/` contains `comment-check.sh` and `todo-tracker.sh`
- WHEN `./skills/setup.sh --sync-claude` runs
- THEN hooks are copied to `~/.claude/hooks/`
- AND existing user hooks are NOT overwritten (copy only if destination doesn't exist)

---

### Requirement: Path-Scoped Rules

`GentlemanClaude/CLAUDE.md` MUST contain a `### Path-Scoped Rules` section that maps file glob patterns to convention sets.

The section MUST be documentation/instructions only — no executable code. The section MUST include at least 4 glob→convention mappings.

#### Scenario: AI edits a Go file

- GIVEN the path-scoped rules section is present in CLAUDE.md
- WHEN the AI is about to edit a `**/*.go` file
- THEN the instructions direct it to apply Go conventions (table-driven tests, error wrapping with `%w`, explicit returns)

#### Scenario: AI edits a Lua file

- GIVEN the path-scoped rules section is present
- WHEN the AI is about to edit a `**/*.lua` file
- THEN the instructions direct it to apply Neovim/lazy.nvim conventions

#### Scenario: AI edits a SKILL.md file

- GIVEN the path-scoped rules section is present
- WHEN the AI is about to edit a `**/SKILL.md` file
- THEN the instructions direct it to load the `skill-creator` skill first

#### Scenario: AI edits an Obsidian markdown file

- GIVEN the path-scoped rules section is present
- WHEN the AI is about to edit a `**/*.md` file inside `obsidian-brain/` or `.obsidian-brain/`
- THEN the instructions direct it to apply Obsidian conventions (wikilinks, frontmatter, templates)

---

### Requirement: Skill Version Table

`AGENTS.md` MUST contain a `### Skill Versions` table listing all skills with their version and last-updated date.

The table MUST be populated by reading the `version:` field from each SKILL.md frontmatter. The table MUST include all skills under `GentlemanClaude/skills/`.

#### Scenario: Table includes all current skills

- GIVEN 40 SKILL.md files exist under `GentlemanClaude/skills/`
- WHEN the `### Skill Versions` table is generated
- THEN the table contains 40 rows, one per skill
- AND each row shows the skill name, version number, and last-updated date

#### Scenario: Table is sorted alphabetically

- GIVEN the skill version table exists
- WHEN a reader scans the table
- THEN skills are listed in alphabetical order by name

---

### Requirement: Orchestrator Self-Check Checklist

The Delegation Rules section in BOTH `CLAUDE.md` and `AGENTS.md` MUST include a concrete inline checklist replacing the plain-text self-check instruction.

The checklist MUST enumerate at least 5 specific categories of work that require delegation.

#### Scenario: Orchestrator uses checklist before responding

- GIVEN the self-check checklist is present in the Delegation Rules
- WHEN the orchestrator prepares a response
- THEN it evaluates each checklist item (read source code?, write code?, analyze architecture?, run tests?, write specs?)
- AND if ANY item applies, delegates to a sub-agent instead of responding inline

#### Scenario: All items unchecked allows inline response

- GIVEN the orchestrator is answering a simple coordination question
- WHEN all 5+ checklist items evaluate to "no"
- THEN the orchestrator responds inline (no delegation needed)

---

## MODIFIED Requirements

### Requirement: Mind Map Generation in Consolidation Skill

(Previously: Mind Map used plain text labels for nodes, separate from wikilinked entities)

The `obsidian-consolidation` SKILL.md Mind Map section MUST instruct the agent to use `[[wikilinked entities]]` as node text in the markmap code block, connecting the visual mind map directly to the Obsidian knowledge graph.

#### Scenario: Mind map nodes use wikilinks

- GIVEN a consolidation is being generated with entities `[[Go]]`, `[[error-wrapping]]`, `[[connection-pool]]`
- WHEN the Mind Map section is produced
- THEN markmap leaf nodes use `[[Go]]`, `[[error-wrapping]]`, `[[connection-pool]]` instead of plain text "Go error handling", "connection pool"

#### Scenario: Mind map reuses extracted entities only

- GIVEN the Entities section has been generated with a specific set of entities
- WHEN the Mind Map is produced
- THEN all node labels are drawn from the already-extracted entities
- AND no new entities are introduced in the map that weren't in the Entities section

## REMOVED Requirements

(None — this change only adds and modifies.)
