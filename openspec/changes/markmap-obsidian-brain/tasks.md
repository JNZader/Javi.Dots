# Tasks: Markmap Mind Map Integration into Obsidian Brain

## Phase 1: Neovim Plugin (Foundation)

- [x] 1.1 Create `GentlemanNvim/nvim/lua/plugins/markmap.lua` with lazy.nvim spec: `return { "Zeioth/markmap.nvim", build = "yarn global add markmap-cli", cmd = { "MarkmapOpen", "MarkmapSave", "MarkmapWatch", "MarkmapWatchStop" }, opts = { html_output = "/tmp/markmap.html", hide_toolbar = false, grace_period = 3600000 }, config = function(_, opts) require("markmap").setup(opts) end }`. Follow the `return { ... }` pattern from `markdown.lua`.

## Phase 2: Template Updates

- [x] 2.1 Append `## Mind Map` section to `GentlemanNvim/obsidian-brain/core/templates/consolidation.md` after `## Action Items`. Include a `markmap` fenced code block with placeholder tree: root = "Weekly Insights", branches = themes (e.g., "Theme A", "Theme B"), leaves = entities. Verify total file stays under 50 lines.

- [x] 2.2 Append `## Decision Map` section to `GentlemanNvim/obsidian-brain/developer/templates/adr.md` after `## Alternatives Considered`. Include a `markmap` fenced code block with placeholder decision tree: root = "Decision", branches = alternatives, leaves = consequences. Verify total file stays under 50 lines.

## Phase 3: AI Skill Update

- [x] 3.1 Bump `metadata.version` from `"1.1"` to `"1.2"` in `GentlemanClaude/skills/obsidian-consolidation/SKILL.md` frontmatter.

- [x] 3.2 Add step 6.5 (between entity extraction step 6 and save/output step 7) to the "Process Flow" section in `obsidian-consolidation/SKILL.md`: "Generate Mind Map — Build a `markmap` fenced code block from the entities and connections identified in steps 3-5. Use 2-3 levels: period as root, insight themes as branches, individual `[[entity]]` wikilinks as leaves. Keep it concise (max 15 leaf nodes)."

- [x] 3.3 Update the "Template Reference" section in `obsidian-consolidation/SKILL.md` to include `## Mind Map` in the listed template sections.

- [x] 3.4 Add `## Mind Map` section to the "Core Only" role-aware output template in `obsidian-consolidation/SKILL.md`, placed after `## Entities`. Show the expected markmap code block format with wikilink entities as leaves.

- [x] 3.5 Add `## Mind Map` to the "Complete Example (Core Only)" section in `obsidian-consolidation/SKILL.md`. Generate a concrete markmap code block using the example's entities (JWT, connection-pool, api-gateway, onboarding-v2, etc.) grouped under insight themes.

- [x] 3.6 Add a "Mind Map Rules" subsection to "Critical Rules" in `obsidian-consolidation/SKILL.md`: (1) Mind Map section is mandatory — every consolidation MUST include it. (2) Reuse entities from `## Entities` — do NOT introduce new entities in the map. (3) Max 3 levels deep, max 15 leaf nodes. (4) Use `[[wikilinks]]` for entity leaves.

## Phase 4: Verification

- [ ] 4.1 Verify `markmap.lua` loads lazily: check that `:Lazy` shows it as "not loaded" before running any `:Markmap*` command. Verify the return table is syntactically valid Lua (no parse errors on Neovim startup).

- [ ] 4.2 Verify `consolidation.md` retains all original sections (`## Period`, `## Top Insights`, `## Connections Made`, `## Open Questions`, `## Action Items`) unchanged, and the new `## Mind Map` section is appended last. Confirm total line count is under 50.

- [ ] 4.3 Verify `adr.md` retains all original sections (`## Context`, `## Decision`, `## Consequences`, `## Alternatives Considered`) unchanged, and the new `## Decision Map` section is appended last. Confirm total line count is under 50.

- [ ] 4.4 Verify `obsidian-consolidation/SKILL.md` has version `"1.2"`, includes the Mind Map generation step, includes the Mind Map in the output template, includes a concrete example, and lists Mind Map rules in Critical Rules. Confirm no existing sections were removed or reordered.
