# Delta for Markmap Integration

## ADDED Requirements

### Requirement: markmap.nvim Plugin

The system MUST provide a `GentlemanNvim/nvim/lua/plugins/markmap.lua` file containing a lazy.nvim plugin spec for `Zeioth/markmap.nvim`.

The plugin MUST be lazy-loaded using the `cmd` key, activating only on `:MarkmapOpen`, `:MarkmapSave`, `:MarkmapWatch`, and `:MarkmapWatchStop` commands.

The plugin MUST include a `build` step that runs `yarn global add markmap-cli`.

The plugin MUST call `require("markmap").setup(opts)` in its `config` function.

#### Scenario: Plugin file follows lazy.nvim pattern

- GIVEN the file `GentlemanNvim/nvim/lua/plugins/markmap.lua`
- WHEN its structure is compared to `GentlemanNvim/nvim/lua/plugins/markdown.lua`
- THEN it uses the same `return { ... }` table pattern
- AND it includes `cmd`, `build`, `opts`, and `config` keys

#### Scenario: Lazy loading via cmd

- GIVEN Neovim starts with the markmap plugin installed
- WHEN the user has NOT run any `:Markmap*` command
- THEN the plugin is NOT loaded (verified via `:Lazy`)
- AND Neovim startup time is unaffected

#### Scenario: MarkmapOpen generates HTML mind map

- GIVEN a markdown buffer is open in Neovim
- AND `markmap-cli` is installed globally
- WHEN the user runs `:MarkmapOpen`
- THEN an HTML file is generated at the configured `html_output` path
- AND the default browser opens with an interactive mind map

#### Scenario: markmap-cli not installed

- GIVEN `markmap-cli` is NOT installed globally
- AND `yarn` is NOT available
- WHEN Neovim starts
- THEN no error is raised (plugin is lazy-loaded)
- WHEN the user runs `:MarkmapOpen`
- THEN an error message indicates markmap-cli is missing

### Requirement: No Conflict with Existing Plugins

The markmap.nvim plugin MUST NOT conflict with `obsidian.nvim` or `render-markdown.nvim`. It operates on the buffer level and generates external HTML output, independent of in-buffer rendering plugins.

#### Scenario: Coexistence with render-markdown.nvim

- GIVEN both `render-markdown.nvim` and `markmap.nvim` are installed
- WHEN a markdown file is opened
- THEN `render-markdown.nvim` provides in-buffer heading/bullet rendering as before
- AND `:MarkmapOpen` generates an external HTML visualization
- AND neither plugin interferes with the other

#### Scenario: Coexistence with obsidian.nvim

- GIVEN both `obsidian.nvim` and `markmap.nvim` are installed
- WHEN a file in an Obsidian workspace is opened
- THEN obsidian.nvim provides wiki-link completion and navigation as before
- AND `:MarkmapOpen` works on the same buffer without conflict

## MODIFIED Requirements

### Requirement: Consolidation Template Content (from templates/spec.md)

The consolidation template (`core/templates/consolidation.md`) MUST now additionally include a `## Mind Map` section at the end of the template, after `## Action Items`. This section MUST contain a `markmap` fenced code block with placeholder content showing the expected tree structure.

(Previously: The template contained only `## Period`, `## Top Insights`, `## Connections Made`, `## Open Questions`, `## Action Items` — no visualization section.)

The existing sections MUST remain unchanged.

#### Scenario: Consolidation template with Mind Map section

- GIVEN the file `core/templates/consolidation.md`
- WHEN opened
- THEN it contains all original sections: `## Period`, `## Top Insights`, `## Connections Made`, `## Open Questions`, `## Action Items`
- AND after `## Action Items`, a `## Mind Map` section exists
- AND the `## Mind Map` section contains a fenced code block with language identifier `markmap`
- AND the code block contains a placeholder tree structure using markdown headings

#### Scenario: Template stays within size limit

- GIVEN the file `core/templates/consolidation.md` with the new Mind Map section
- WHEN the total line count is measured
- THEN it SHOULD be between 25 and 45 lines
- AND it MUST be under 50 lines (per template minimalism spec)

#### Scenario: Mind Map section is at the end

- GIVEN the file `core/templates/consolidation.md`
- WHEN the section order is examined
- THEN `## Mind Map` is the LAST section in the file
- AND it does not interrupt the flow of the core template sections

### Requirement: ADR Template Content (from templates/spec.md)

The ADR template (`developer/templates/adr.md`) MUST now additionally include a `## Decision Map` section at the end of the template, after `## Alternatives Considered`. This section MUST contain a `markmap` fenced code block with placeholder content showing a decision tree structure (decision at root, alternatives as branches, consequences as leaves).

(Previously: The template contained only `## Context`, `## Decision`, `## Consequences`, `## Alternatives Considered` — no visualization section.)

The existing sections MUST remain unchanged.

#### Scenario: ADR template with Decision Map section

- GIVEN the file `developer/templates/adr.md`
- WHEN opened
- THEN it contains all original sections: `## Context`, `## Decision`, `## Consequences`, `## Alternatives Considered`
- AND after `## Alternatives Considered`, a `## Decision Map` section exists
- AND the `## Decision Map` section contains a fenced code block with language identifier `markmap`
- AND the code block contains a placeholder decision tree using markdown headings

#### Scenario: ADR template stays within size limit

- GIVEN the file `developer/templates/adr.md` with the new Decision Map section
- WHEN the total line count is measured
- THEN it SHOULD be between 22 and 40 lines
- AND it MUST be under 50 lines (per template minimalism spec)

### Requirement: obsidian-consolidation Skill Mind Map Generation (from skills/spec.md)

The consolidation skill (`GentlemanClaude/skills/obsidian-consolidation/SKILL.md`) MUST be updated to version 1.2. The skill MUST now additionally instruct the AI to generate a `## Mind Map` section containing a `markmap` fenced code block as part of its consolidation output.

(Previously: The skill generated `## Top Insights`, `## Connections Made`, `## Open Questions`, `## Action Items`, `## Entities`, and role-aware sections — but no visualization section.)

The Mind Map MUST be derived from the already-extracted entities and connections, not from new analysis. It MUST be concise (2-3 levels deep).

#### Scenario: Skill version bump

- GIVEN the file `GentlemanClaude/skills/obsidian-consolidation/SKILL.md`
- WHEN the frontmatter is read
- THEN `metadata.version` is `"1.2"`

#### Scenario: Mind Map in process flow

- GIVEN the consolidation skill instructions
- WHEN the process flow is examined
- THEN there is a step between entity extraction and save/output that generates the Mind Map
- AND the step instructs the AI to build the markmap from entities and connections already identified

#### Scenario: Mind Map output format

- GIVEN the AI generates a consolidation note using the updated skill
- WHEN the output is examined
- THEN it contains a `## Mind Map` section
- AND the section contains a fenced code block with language identifier `markmap`
- AND the markmap content uses markdown headings for hierarchy:
  - Level 1 (`#`): The consolidation period (root)
  - Level 2 (`##`): Top insight themes / connection clusters
  - Level 3 (`###`): Individual entities as `[[wikilinks]]`

#### Scenario: Mind Map conciseness

- GIVEN the AI generates a Mind Map section
- WHEN the markmap code block is examined
- THEN it has at most 3 levels of depth
- AND it has at most 15 leaf nodes
- AND it reuses entities already extracted (no new analysis)

#### Scenario: Existing output sections unchanged

- GIVEN the updated consolidation skill
- WHEN the AI generates a consolidation note
- THEN all previous sections still appear: `## Period`, `## Top Insights`, `## Connections Made`, `## Open Questions`, `## Action Items`, `## Entities`
- AND role-aware sections (if applicable) remain unchanged
- AND `## Mind Map` appears after `## Entities`

#### Scenario: Mind Map section placement in output

- GIVEN the AI generates a consolidation note
- WHEN the section order is examined
- THEN `## Mind Map` appears after `## Entities` and before any role-aware addendum sections
- AND if no Entities section exists (edge case), the Mind Map still appears as the last core section

### Requirement: Template Reference in Skill

The consolidation skill's "Template Reference" section MUST be updated to reflect the new `## Mind Map` section in the consolidation template.

(Previously: Listed sections as `## Period`, `## Top Insights`, `## Connections Made`, `## Open Questions`, `## Action Items`.)

#### Scenario: Updated template reference

- GIVEN the consolidation skill's "Template Reference" section
- WHEN examined
- THEN it lists all template sections including `## Mind Map`

## REMOVED Requirements

(None — this change is purely additive.)
