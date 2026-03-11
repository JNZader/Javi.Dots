# AI Skills Specification

## Purpose

Defines the three new AI skills (`obsidian-braindump`, `obsidian-consolidation`, `obsidian-resource-capture`) to be created in `GentlemanClaude/skills/`. These skills guide AI coding assistants in generating Obsidian-compatible notes using the templates from the corresponding role packs.

## Requirements

### Requirement: Skill Format

Each skill MUST follow the existing SKILL.md frontmatter format as used by other skills in the repository (e.g., `react-19/SKILL.md`). The frontmatter MUST include `name`, `description`, `license`, and `metadata` (author, version) fields delimited by `---`.

#### Scenario: Valid frontmatter

- GIVEN any of the 3 new skill files
- WHEN the frontmatter is parsed
- THEN it contains:
  - `name:` matching the skill directory name
  - `description:` with a `>` multi-line scalar containing a trigger description
  - `license: Apache-2.0`
  - `metadata.author: gentleman-programming`
  - `metadata.version: "1.0"`

#### Scenario: Consistent with existing skills

- GIVEN the file `GentlemanClaude/skills/obsidian-braindump/SKILL.md`
- WHEN compared structurally to `GentlemanClaude/skills/react-19/SKILL.md`
- THEN the frontmatter structure is identical (same keys, same format)

### Requirement: obsidian-braindump Skill

The system MUST create `GentlemanClaude/skills/obsidian-braindump/SKILL.md`. This skill MUST instruct the AI to:

1. Generate a braindump note following the `braindump.md` template format
2. Include YAML frontmatter with title, date, and tags
3. Capture the user's raw thought in the `## Thought` section
4. Suggest relevant tags based on content
5. Optionally suggest related notes using `[[wiki-link]]` syntax

#### Scenario: Braindump skill trigger

- GIVEN an AI assistant has loaded the braindump skill
- WHEN the user says "brain dump this idea" or "capture this thought"
- THEN the AI generates a markdown note matching the braindump template structure

#### Scenario: Braindump output format

- GIVEN the AI generates a braindump note
- WHEN the output is examined
- THEN it starts with `---` YAML frontmatter
- AND contains `## Thought`, `## Context`, `## Related Notes` sections
- AND related notes use `[[note-name]]` wiki-link syntax

### Requirement: obsidian-consolidation Skill

The system MUST create `GentlemanClaude/skills/obsidian-consolidation/SKILL.md`. This skill MUST instruct the AI to:

1. Generate a weekly consolidation note following the `consolidation.md` template format
2. Synthesize insights from recent notes (the user provides context)
3. Identify connections between notes
4. Surface open questions and action items
5. Use `[[wiki-link]]` syntax to reference source notes

#### Scenario: Consolidation skill trigger

- GIVEN an AI assistant has loaded the consolidation skill
- WHEN the user says "consolidate this week's notes" or "weekly synthesis"
- THEN the AI generates a consolidation note matching the template structure

#### Scenario: Consolidation references sources

- GIVEN the AI generates a consolidation note
- WHEN the output is examined
- THEN the `## Top Insights` and `## Connections Made` sections reference source notes using `[[wiki-link]]` syntax

### Requirement: obsidian-resource-capture Skill

The system MUST create `GentlemanClaude/skills/obsidian-resource-capture/SKILL.md`. This skill MUST instruct the AI to:

1. Generate a resource capture note following the `resource-capture.md` template format
2. Accept a URL or reference and generate a summary
3. Extract key takeaways
4. Suggest relevant tags
5. Format the source as a markdown link

#### Scenario: Resource capture skill trigger

- GIVEN an AI assistant has loaded the resource capture skill
- WHEN the user says "capture this resource" or "save this link" and provides a URL
- THEN the AI generates a resource capture note matching the template structure

#### Scenario: Resource capture output format

- GIVEN the AI generates a resource capture note
- WHEN the output is examined
- THEN it contains the source URL in `## Source` as a markdown link `[title](url)`
- AND `## Summary` contains a 2-4 sentence summary
- AND `## Key Takeaways` contains a bulleted list

### Requirement: Role-Awareness

Each skill SHOULD include a section that instructs the AI to adapt its output based on the active role packs. If the Developer pack is active, braindumps MAY include code-related fields. If the PM pack is active, consolidation notes MAY include team-facing summaries.

#### Scenario: Developer-aware braindump

- GIVEN the braindump skill instructions include a role-awareness section
- WHEN the user's vault contains developer pack templates (detected by context)
- THEN the AI MAY add a `## Code Reference` section to the braindump output

#### Scenario: PM-aware consolidation

- GIVEN the consolidation skill instructions include a role-awareness section
- WHEN the user's vault contains PM pack templates (detected by context)
- THEN the AI MAY add a `## Team Summary` section suitable for sharing

#### Scenario: No role packs detected

- GIVEN only Core templates exist in the vault
- WHEN any skill generates output
- THEN the output matches the base Core template format without role-specific additions

### Requirement: Skill File Location

Each skill MUST be placed in its own directory under `GentlemanClaude/skills/`:

- `GentlemanClaude/skills/obsidian-braindump/SKILL.md`
- `GentlemanClaude/skills/obsidian-consolidation/SKILL.md`
- `GentlemanClaude/skills/obsidian-resource-capture/SKILL.md`

#### Scenario: Skill directories exist

- GIVEN the repository is checked out
- WHEN listing `GentlemanClaude/skills/`
- THEN `obsidian-braindump/`, `obsidian-consolidation/`, and `obsidian-resource-capture/` directories exist
- AND each contains a single `SKILL.md` file

### Requirement: AGENTS.md Registration

The system MUST register the 3 new skills in the `AGENTS.md` skill table under the "Generic Skills" section.

#### Scenario: Skills appear in AGENTS.md

- GIVEN `AGENTS.md` is read
- WHEN the "Generic Skills" table is examined
- THEN it contains rows for `obsidian-braindump`, `obsidian-consolidation`, and `obsidian-resource-capture`
- AND each row includes the skill name, description, and link to the SKILL.md file

### Requirement: Skill Content Structure

Each skill body (below the frontmatter) MUST follow this structure:

1. A brief description of when to use the skill (1-2 sentences)
2. The output format specification (template sections and their content guidelines)
3. Examples of generated output (at least one)
4. Role-awareness section (optional behavior based on detected packs)

#### Scenario: Skill body structure

- GIVEN any of the 3 new skill files
- WHEN the body (below `---` closing the frontmatter) is read
- THEN it contains:
  - A heading describing the workflow
  - A code block showing the expected output format
  - At least one example
  - A section on role-aware behavior

### Requirement: No External Dependencies

The skills MUST NOT require any external tools, APIs, or plugins to function. They are pure instruction sets for AI assistants that generate markdown text.

#### Scenario: Standalone operation

- GIVEN an AI assistant with only the skill file loaded
- WHEN the user triggers the skill
- THEN the AI can generate the expected output without any external calls, APIs, or file system access beyond writing the note
