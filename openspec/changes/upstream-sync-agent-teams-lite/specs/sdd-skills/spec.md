# Delta for SDD Skills

## ADDED Requirements

### Requirement: Hybrid Persistence Mode

The persistence contract MUST support a `hybrid` mode where artifacts are persisted to BOTH engram and openspec simultaneously.

#### Scenario: Hybrid mode writes to both backends

- GIVEN `artifact_store.mode` is `hybrid`
- WHEN a sub-agent persists an artifact (e.g., proposal)
- THEN the artifact is saved to engram via `mem_save` with deterministic naming
- AND the artifact is written to `openspec/changes/{change-name}/` filesystem path
- AND both copies contain identical content

#### Scenario: Hybrid mode reads prefer engram

- GIVEN `artifact_store.mode` is `hybrid`
- WHEN a sub-agent needs to read a dependency artifact
- THEN it reads from engram first (two-step recovery)
- AND falls back to filesystem if engram retrieval fails

#### Scenario: Auto-mode resolves to hybrid when applicable

- GIVEN `artifact_store.mode` is `auto` (Javi.Dots extension)
- WHEN engram is available AND user explicitly requested file artifacts
- THEN `auto` resolves to `hybrid`

### Requirement: Inline Engram Instructions

Each SDD skill MUST contain inline engram instructions directly in its SKILL.md, not solely defer to `_shared/engram-convention.md` via "Read and follow" directives.

#### Scenario: Skill contains inline mem_save example

- GIVEN any SDD skill SKILL.md (sdd-init through sdd-archive)
- WHEN the skill's engram mode section is read
- THEN it contains inline `mem_save(...)` and `mem_search(...)` code blocks with the exact parameters for that skill's artifact type
- AND it references `_shared/engram-convention.md` as supplementary, not primary

#### Scenario: Installed skill works without _shared/ directory

- GIVEN a skill is installed to `~/.claude/skills/sdd-propose/SKILL.md`
- AND `~/.claude/skills/_shared/` does NOT exist
- WHEN an AI agent reads the skill
- THEN it has sufficient inline instructions to persist to engram correctly

### Requirement: Mandatory Persist Step

Each SDD skill MUST include an explicit "Persist Artifact" step before the "Return Summary" step.

#### Scenario: Persist step precedes return

- GIVEN any SDD skill SKILL.md
- WHEN the execution steps are read in order
- THEN there is a step titled "Persist Artifact" (or equivalent)
- AND this step comes AFTER the content creation step
- AND BEFORE the "Return Summary" step

#### Scenario: Persist step is mode-aware

- GIVEN the "Persist Artifact" step in any SDD skill
- WHEN the mode is `engram`
- THEN the step calls `mem_save` with deterministic naming
- WHEN the mode is `openspec`
- THEN the step writes to the filesystem path
- WHEN the mode is `hybrid`
- THEN the step does both
- WHEN the mode is `none`
- THEN the step is skipped

### Requirement: Skill Registry Loading

Each SDD skill (except sdd-init which BUILDS the registry) MUST include a "Step 1: Load Skill Registry" that resolves the project's available skills before executing.

#### Scenario: Registry loaded as first step

- GIVEN any SDD skill except sdd-init
- WHEN the execution steps are read
- THEN the first numbered step is "Load Skill Registry" or equivalent

#### Scenario: sdd-init builds the registry

- GIVEN the sdd-init skill
- WHEN its steps are read
- THEN it contains a step to "Build Skill Registry" that scans the project for available skills
- AND persists the registry for other skills to consume

### Requirement: Skill Registry Skill

The system MUST include a new `skill-registry/SKILL.md` skill that provides the registry-building and resolution logic.

#### Scenario: Skill file exists

- GIVEN the repository is checked out
- WHEN listing `GentlemanClaude/skills/skill-registry/`
- THEN a `SKILL.md` file exists
- AND its frontmatter includes `name: skill-registry`

#### Scenario: Registered in AGENTS.md

- GIVEN `AGENTS.md` is read
- WHEN the "Generic Skills" table is examined
- THEN it contains a row for `skill-registry` with description and link

### Requirement: Sub-Agent Context Rules

The persistence contract MUST define which agents read and write which artifact types.

#### Scenario: Context rules table exists

- GIVEN `persistence-contract.md` is read
- WHEN the "Sub-Agent Context Rules" section is examined
- THEN it contains a table mapping each SDD skill to its read dependencies and write targets

#### Scenario: Skills follow context rules

- GIVEN the sub-agent context rules define that sdd-tasks reads `proposal`, `spec`, and `design`
- WHEN sdd-tasks executes
- THEN it retrieves exactly those three artifacts as dependencies

### Requirement: State Persistence

The openspec convention MUST support a `state.yaml` file per change that tracks artifact completion status.

#### Scenario: state.yaml path defined

- GIVEN `openspec-convention.md` is read
- WHEN the directory structure is examined
- THEN `state.yaml` appears under `openspec/changes/{change-name}/`

#### Scenario: Missing state.yaml handled gracefully

- GIVEN a change directory exists WITHOUT `state.yaml`
- WHEN any SDD skill accesses the change
- THEN it treats the state as initial (no artifacts completed)
- AND does NOT error or refuse to proceed

### Requirement: Orchestrator Delegation Rules

The SDD orchestrator sections in `CLAUDE.md` and `AGENTS.md` MUST include explicit delegation rules.

#### Scenario: Delegation rules section exists

- GIVEN `GentlemanClaude/CLAUDE.md` is read
- WHEN the SDD Orchestrator section is examined
- THEN it contains a "Delegation Rules" subsection
- AND the rules specify when to delegate vs handle inline

#### Scenario: Anti-patterns section exists

- GIVEN `GentlemanClaude/CLAUDE.md` is read
- WHEN the SDD Orchestrator section is examined
- THEN it contains an "Anti-patterns" subsection listing behaviors the orchestrator must avoid

#### Scenario: Task Escalation table exists

- GIVEN `GentlemanClaude/CLAUDE.md` is read
- WHEN the SDD Orchestrator section is examined
- THEN it contains a "Task Escalation" table mapping task complexity to handling strategy

### Requirement: Orchestrator Merges Preserve Javi.Dots Sections

When updating `CLAUDE.md` and `AGENTS.md`, the system MUST preserve all Javi.Dots-specific content.

#### Scenario: CLAUDE.md personality preserved

- GIVEN `GentlemanClaude/CLAUDE.md` is updated with upstream orchestrator content
- WHEN the file is read
- THEN the Personality, Language, Tone, Philosophy, Behavior sections remain unchanged
- AND the Domain Routing section remains unchanged
- AND the Framework/Library Detection table remains unchanged
- AND the Plugin Detection section remains unchanged

#### Scenario: AGENTS.md project sections preserved

- GIVEN `AGENTS.md` is updated with upstream orchestrator content
- WHEN the file is read
- THEN the repository skills table remains unchanged
- AND the auto-invoke skills section remains unchanged
- AND the project overview section remains unchanged

#### Scenario: AGENTS.md auto-mode includes hybrid

- GIVEN `AGENTS.md` is updated
- WHEN the Artifact Store Policy auto-resolution chain is read
- THEN `hybrid` appears as a resolution option (when engram is available AND user wants file artifacts)

### Requirement: Shared Directory Sync

The `skills/setup.sh` MUST sync the `_shared/` directory to the user's config alongside skill directories.

#### Scenario: _shared/ copied to user config

- GIVEN `skills/setup.sh` is executed
- WHEN the sync completes
- THEN `~/.claude/skills/_shared/` exists
- AND it contains `persistence-contract.md`, `engram-convention.md`, and `openspec-convention.md`

#### Scenario: Existing _shared/ is overwritten on sync

- GIVEN `~/.claude/skills/_shared/` already exists with older content
- WHEN `skills/setup.sh` is executed
- THEN the files are overwritten with current versions

## MODIFIED Requirements

### Requirement: SDD Skill Version

(Previously: All SDD skills at version "2.0")

All SDD skill SKILL.md files MUST have `metadata.version: "3.3"` in their frontmatter.

#### Scenario: Version updated

- GIVEN any SDD skill SKILL.md
- WHEN the frontmatter is parsed
- THEN `metadata.version` is `"3.3"`

### Requirement: Persistence Contract Mode Table

(Previously: Table had 3 rows — engram, openspec, none)

The persistence contract mode table MUST include a `hybrid` row.

#### Scenario: Four modes documented

- GIVEN `persistence-contract.md` is read
- WHEN the "Behavior Per Mode" table is examined
- THEN it has 4 rows: `engram`, `openspec`, `hybrid`, `none`

### Requirement: Engram Convention Reference Role

(Previously: `_shared/engram-convention.md` was the primary reference, skills said "Read and follow")

The `engram-convention.md` MUST include a NOTE stating that inline instructions in each skill are the primary reference. The convention file is supplementary context for completeness.

#### Scenario: Convention file has inline-primary note

- GIVEN `engram-convention.md` is read
- WHEN the top of the file is examined
- THEN it contains a NOTE or callout indicating that skill-inline instructions take precedence

## REMOVED Requirements

(None — this is a purely additive sync.)
