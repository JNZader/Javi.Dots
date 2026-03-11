# Proposal: Markmap Mind Map Integration into Obsidian Brain

## Intent

The Obsidian Brain ecosystem produces structured knowledge artifacts (consolidations, ADRs) that are inherently hierarchical — insights branch into connections, decisions branch into alternatives and consequences. Yet there is no visualization layer to surface this structure. Users must mentally parse nested markdown to see the knowledge graph.

Markmap converts markdown headings/lists into interactive HTML mind maps. Integrating it at three levels — Neovim plugin, templates, and AI skill — gives users visual synthesis of their knowledge without leaving their editor or Obsidian workflow.

## Scope

### In Scope

- **Neovim plugin**: Add `markmap.nvim` (Zeioth) to `GentlemanNvim/nvim/lua/plugins/markmap.lua` — lazy-loaded on `:MarkmapOpen`/`:MarkmapSave`/`:MarkmapWatch`/`:MarkmapWatchStop` commands
- **Consolidation template**: Add an optional `## Mind Map` section with a `markmap` code block to `core/templates/consolidation.md`
- **ADR template**: Add an optional `## Decision Map` section with a `markmap` code block to `developer/templates/adr.md`
- **Consolidation skill**: Update `obsidian-consolidation/SKILL.md` (v1.1 -> v1.2) to instruct the AI to generate a `## Mind Map` section with a markmap code block visualizing the week's key insights and connections

### Out of Scope

- markmap-cli installation (user responsibility; lazy.nvim `build` step handles it if yarn is present)
- markmap-mcp-server (immature, unclear value-add over inline code block generation)
- braindump/resource-capture skill updates (contradicts speed-first philosophy)
- VS Code extension or web-based rendering
- TUI installer changes (no new install steps needed)
- Templates beyond consolidation and ADR (low/no value per exploration analysis)

## Approach

1. **Plugin file**: Create `markmap.lua` following the exact lazy.nvim return-table pattern used by `markdown.lua` and other plugins in the repo. Use `cmd` lazy-loading so the plugin loads only when the user explicitly invokes a markmap command. The `build` step runs `yarn global add markmap-cli`.

2. **Template updates**: Append a clearly labeled optional section at the end of each template (after all existing sections) so the core template flow is unaffected. The markmap code block uses markdown heading syntax that Obsidian Mindmap NextGen renders inline and markmap.nvim renders via `:MarkmapOpen`.

3. **Skill update**: Add a `## Mind Map` generation step to the consolidation skill's process flow and output format. The AI generates a concise markmap (2-3 levels deep) using the already-extracted entities and connections as tree nodes/leaves. Bump version to 1.2.

## Affected Areas

| Area | Impact | Description |
|------|--------|-------------|
| `GentlemanNvim/nvim/lua/plugins/markmap.lua` | New | markmap.nvim lazy.nvim plugin spec |
| `GentlemanNvim/obsidian-brain/core/templates/consolidation.md` | Modified | Add optional `## Mind Map` section at end |
| `GentlemanNvim/obsidian-brain/developer/templates/adr.md` | Modified | Add optional `## Decision Map` section at end |
| `GentlemanClaude/skills/obsidian-consolidation/SKILL.md` | Modified | Add Mind Map generation to process flow, bump to v1.2 |

## Risks

| Risk | Likelihood | Mitigation |
|------|------------|------------|
| yarn not installed — `build` step fails silently | Medium | `cmd` lazy-loading means plugin never errors on startup; only fails when user runs `:MarkmapOpen`. Document yarn as prerequisite in plugin comment. |
| Raw `markmap` code blocks visible to users without renderer | Low | Section is clearly labeled "Mind Map (optional)" and placed at the end of templates. Users without Obsidian Mindmap NextGen see valid markdown (just headings/lists). |
| Template size increase | Low | Each markmap block adds ~8-12 lines. Templates stay under 50 lines (spec limit). |
| Skill output token increase | Low | Markmap block is concise (2-3 levels, entities as leaves). Adds ~100-150 tokens to consolidation output. |

## Rollback Plan

1. Delete `GentlemanNvim/nvim/lua/plugins/markmap.lua`
2. Remove the `## Mind Map` section from `consolidation.md` (restore to current 25-line version)
3. Remove the `## Decision Map` section from `adr.md` (restore to current 22-line version)
4. Revert `obsidian-consolidation/SKILL.md` to v1.1 (remove Mind Map process step and output section)

All changes are additive — no existing behavior is modified. Rollback is a clean removal of added content.

## Dependencies

- `markmap-cli` must be installed globally via `yarn global add markmap-cli` for the Neovim plugin to function (handled by lazy.nvim `build` step if yarn is available)
- No Lua plugin dependencies (markmap.nvim is self-contained)
- No Obsidian plugin dependencies for template rendering (markmap code blocks are valid markdown regardless)

## Success Criteria

- [ ] `:MarkmapOpen` on any markdown buffer in Neovim opens an interactive HTML mind map in the browser
- [ ] `consolidation.md` template contains a valid `markmap` code block that renders in Obsidian with Mindmap NextGen plugin
- [ ] `adr.md` template contains a valid `markmap` code block that renders in Obsidian with Mindmap NextGen plugin
- [ ] The consolidation AI skill generates a `## Mind Map` section with a syntactically valid markmap code block
- [ ] All existing template sections remain unchanged (no regression)
- [ ] The markmap.nvim plugin loads lazily (does not increase Neovim startup time)
- [ ] Templates stay under 50 lines total (per template minimalism spec)
