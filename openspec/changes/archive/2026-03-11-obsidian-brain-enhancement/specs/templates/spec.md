# Template Assets Specification

## Purpose

Defines the template markdown files, folder structure, and asset layout for each role pack (Core, Developer, PM/Tech Lead) under `GentlemanNvim/obsidian-brain/`. Templates are plain Obsidian-compatible markdown files with YAML frontmatter.

## Requirements

### Requirement: Directory Structure

The system MUST create the following directory structure under `GentlemanNvim/obsidian-brain/`:

```
GentlemanNvim/obsidian-brain/
  core/
    templates/
      braindump.md
      resource-capture.md
      consolidation.md
      daily-note.md
  developer/
    templates/
      adr.md
      coding-session.md
      tech-debt.md
      debug-journal.md
      sdd-feedback.md
  pm-lead/
    templates/
      meeting-notes.md
      sprint-review.md
      stakeholder-update.md
      risk-registry.md
      daily-brief.md
      weekly-brief.md
      team-intelligence.md
```

#### Scenario: Directory exists in repo

- GIVEN the repository is checked out
- WHEN a user lists `GentlemanNvim/obsidian-brain/`
- THEN the `core/`, `developer/`, and `pm-lead/` subdirectories exist
- AND each contains a `templates/` subdirectory with the expected `.md` files

### Requirement: Template Format

Every template file MUST start with YAML frontmatter delimited by `---`. The frontmatter MUST include at minimum `title`, `date`, and `tags` fields. The frontmatter MUST NOT use plugin-specific syntax (no Templater `<% %>` blocks, no Dataview inline fields). The body MUST use standard Obsidian markdown (headings, lists, checkboxes, wiki-links `[[]]`).

#### Scenario: Valid YAML frontmatter

- GIVEN any template file (e.g., `braindump.md`)
- WHEN the file is opened in a text editor
- THEN the first line is `---`
- AND it contains `title:`, `date:`, and `tags:` fields
- AND the frontmatter is closed with `---`
- AND the body follows in standard markdown

#### Scenario: Obsidian compatibility

- GIVEN any template file
- WHEN the file is opened in Obsidian 1.5+
- THEN the frontmatter is parsed correctly
- AND the file renders without errors or warnings
- AND no plugin-specific syntax causes rendering artifacts

#### Scenario: No plugin-specific syntax

- GIVEN any template file
- WHEN the content is inspected
- THEN it contains no `<% %>` (Templater), no `[field:: value]` (Dataview inline), and no `{{}}` (core Templates plugin interpolation)
- AND placeholders use markdown comments `<!-- TODO: ... -->` or simple text like `{description}`

### Requirement: Core Templates

The system MUST provide the following Core templates:

1. **braindump.md** — Quick capture template with fields for raw thought and optional tags. MUST include sections: Thought, Context, Related Notes.
2. **resource-capture.md** — Link + summary template. MUST include sections: Source URL, Summary, Key Takeaways, Tags.
3. **consolidation.md** — Weekly knowledge synthesis. MUST include sections: Period, Top Insights, Connections Made, Open Questions, Action Items.
4. **daily-note.md** — Daily note template. MUST include sections: Today's Focus, Notes, Reflections, Tomorrow.

#### Scenario: Braindump template content

- GIVEN the file `core/templates/braindump.md`
- WHEN opened
- THEN it has frontmatter with `title: "Braindump"`, `tags: ["braindump", "inbox"]`
- AND the body contains `## Thought`, `## Context`, `## Related Notes` headings

#### Scenario: Resource capture template content

- GIVEN the file `core/templates/resource-capture.md`
- WHEN opened
- THEN it has frontmatter with `tags: ["resource"]`
- AND the body contains `## Source`, `## Summary`, `## Key Takeaways`

#### Scenario: Consolidation template content

- GIVEN the file `core/templates/consolidation.md`
- WHEN opened
- THEN it has frontmatter with `tags: ["consolidation", "weekly"]`
- AND the body contains `## Period`, `## Top Insights`, `## Connections Made`, `## Open Questions`, `## Action Items`

#### Scenario: Daily note template content

- GIVEN the file `core/templates/daily-note.md`
- WHEN opened
- THEN it has frontmatter with `tags: ["daily"]`
- AND the body contains `## Today's Focus`, `## Notes`, `## Reflections`, `## Tomorrow`

### Requirement: Developer Pack Templates

The system MUST provide the following Developer templates:

1. **adr.md** — Architecture Decision Record. MUST include sections: Status, Context, Decision, Consequences, Alternatives Considered.
2. **coding-session.md** — Session log. MUST include sections: Goal, What I Did, Blockers, Decisions Made, Next Steps.
3. **tech-debt.md** — Tracker. MUST include sections: Area, Description, Impact, Effort Estimate, Priority, Plan.
4. **debug-journal.md** — Debug log. MUST include sections: Symptom, Hypothesis, Investigation, Root Cause, Fix, Lessons Learned.
5. **sdd-feedback.md** — SDD feedback loop. MUST include sections: Change Name, Phase Completed, What Worked, What Didn't, Improvements, Link to Engram (optional).

#### Scenario: ADR template content

- GIVEN the file `developer/templates/adr.md`
- WHEN opened
- THEN it has frontmatter with `tags: ["adr", "architecture"]` and a `status:` field
- AND the body contains `## Context`, `## Decision`, `## Consequences`, `## Alternatives Considered`

#### Scenario: SDD feedback template references Engram

- GIVEN the file `developer/templates/sdd-feedback.md`
- WHEN opened
- THEN the body contains a section for `## Link to Engram` with a placeholder `<!-- Optional: engram:// link -->`
- AND the template does NOT require Engram to be installed

### Requirement: PM/Tech Lead Pack Templates

The system MUST provide the following PM/Tech Lead templates:

1. **meeting-notes.md** — Meeting template. MUST include: Attendees, Agenda, Discussion, Decisions, Action Items.
2. **sprint-review.md** — Sprint review. MUST include: Sprint Goal, Completed, Not Completed, Metrics, Retrospective Notes.
3. **stakeholder-update.md** — Status update. MUST include: Summary, Progress, Risks, Next Steps, Ask.
4. **risk-registry.md** — Risk tracker. MUST include: Risk ID, Description, Likelihood, Impact, Mitigation, Owner, Status.
5. **daily-brief.md** — Daily brief. MUST include: Top 3 Priorities, Blockers, Key Decisions Needed, FYI.
6. **weekly-brief.md** — Weekly brief. MUST include: Highlights, Metrics, Risks & Blockers, Next Week Focus, Team Notes.
7. **team-intelligence.md** — Team knowledge capture. MUST include: Observation, Context, Pattern, Recommendation, Shared With.

#### Scenario: Meeting notes template content

- GIVEN the file `pm-lead/templates/meeting-notes.md`
- WHEN opened
- THEN it has frontmatter with `tags: ["meeting"]`
- AND the body contains `## Attendees`, `## Agenda`, `## Discussion`, `## Decisions`, `## Action Items`

#### Scenario: Risk registry template content

- GIVEN the file `pm-lead/templates/risk-registry.md`
- WHEN opened
- THEN it has frontmatter with `tags: ["risk"]`
- AND the body contains a table or structured sections for Risk ID, Description, Likelihood, Impact, Mitigation, Owner, Status

### Requirement: Folder Structure Creation

When the installer copies templates to a project, it MUST create the following folder structure within `{project}/.obsidian-brain/` based on selected packs:

- Core (always): `inbox/`, `resources/`, `knowledge/`, `templates/`
- Developer: `architecture/`, `sessions/`, `debugging/`, `templates/`
- PM/Tech Lead: `meetings/`, `sprints/`, `risks/`, `briefs/`, `templates/`

#### Scenario: Core-only folder structure

- GIVEN the user selected only Core (no optional packs)
- WHEN the installer creates the project vault
- THEN `{project}/.obsidian-brain/` contains `inbox/`, `resources/`, `knowledge/`, `templates/`
- AND the `templates/` folder contains the 4 Core template files

#### Scenario: Core + Developer folder structure

- GIVEN the user selected Core and Developer
- WHEN the installer creates the project vault
- THEN `{project}/.obsidian-brain/` contains Core folders AND `architecture/`, `sessions/`, `debugging/`
- AND `templates/` contains both Core (4) and Developer (5) template files

#### Scenario: All packs folder structure

- GIVEN the user selected Core, Developer, and PM/Tech Lead
- WHEN the installer creates the project vault
- THEN `{project}/.obsidian-brain/` contains all folders from all three packs
- AND `templates/` contains all 16 template files

### Requirement: Template Minimalism

Templates SHOULD be scaffolds, not prescriptive content. Each section SHOULD contain at most a one-line description or placeholder, not multi-paragraph instructions. This allows users to customize after creation.

#### Scenario: Template body length

- GIVEN any template file
- WHEN the total line count (excluding frontmatter) is measured
- THEN it SHOULD be between 15 and 50 lines
- AND each section body contains at most 2 lines of placeholder text
