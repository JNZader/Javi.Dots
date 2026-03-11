# Exploration: Markmap Integration into Javi.Dots Ecosystem

## Current State

### Obsidian Brain Templates
16 templates across 3 role packs (`core/`, `developer/`, `pm-lead/`), all using YAML frontmatter + `[[wikilinks]]`. Templates are minimal markdown scaffolds — no code blocks, no embedded visualizations. The consolidation template is 25 lines, braindump is 15 lines, ADR is 22 lines. They are intentionally lean, with AI skills adding role-aware sections at generation time.

### Neovim Plugin Ecosystem
- `obsidian.lua` — obsidian.nvim with multi-workspace detection, lazy=false, depends on plenary.nvim
- `markdown.lua` — render-markdown.nvim for in-buffer heading/bullet rendering, depends on treesitter + mini.nvim
- All plugins follow lazy.nvim spec pattern: `return { "author/plugin", opts = {...} }`
- No existing mind-map or visualization plugins installed
- render-markdown.nvim already handles markdown beautification but not interactive diagrams

### AI Skills
3 Obsidian skills at v1.1: braindump, consolidation, resource-capture. All include entity extraction with `[[wikilinks]]`, role-aware section injection, and `## Entities` sections. Skills generate pure markdown — no embedded code blocks, no diagram syntax. Skills reference templates in `GentlemanNvim/obsidian-brain/`.

### MCP Configuration
MCP servers defined in `GentlemanClaude/mcp-servers.template.json` (8 servers: context7, engram, atlassian, figma, notion, brave-search, sentry, cloudflare). Also mirrored in `GentlemanOpenCode/opencode.json`. No markmap or diagram-related MCP servers present.

## Affected Areas

- `GentlemanNvim/nvim/lua/plugins/` — New plugin file `markmap.lua` for markmap.nvim
- `GentlemanNvim/obsidian-brain/core/templates/consolidation.md` — Candidate for optional markmap block
- `GentlemanNvim/obsidian-brain/developer/templates/adr.md` — Candidate for decision tree markmap
- `GentlemanClaude/skills/obsidian-consolidation/SKILL.md` — Update to optionally generate markmap blocks
- `GentlemanClaude/skills/obsidian-braindump/SKILL.md` — Low-priority update (braindumps are speed-focused)
- `GentlemanClaude/mcp-servers.template.json` — Potential markmap-mcp-server addition
- `GentlemanOpenCode/opencode.json` — Mirror MCP config if added

## Investigation Findings

### 1. markmap.nvim (Zeioth/markmap.nvim)

**What it does**: Neovim plugin that converts any markdown buffer to an interactive HTML mindmap. Commands: `:MarkmapOpen`, `:MarkmapSave`, `:MarkmapWatch`, `:MarkmapWatchStop`. Outputs to `/tmp/markmap.html` by default.

**Dependencies**: Requires `yarn` globally installed + `markmap-cli` (installed via `yarn global add markmap-cli`). No Lua plugin dependencies (no plenary, no treesitter).

**Conflict analysis**: No conflict with obsidian.nvim — markmap.nvim operates on the buffer level and generates external HTML. No conflict with render-markdown.nvim either — they serve different purposes (in-buffer rendering vs. external visualization).

**lazy.nvim integration**: Uses `cmd` lazy-loading (loads only when `:MarkmapOpen` etc. is called). Minimal overhead.

```lua
-- Exact lazy.nvim spec from the README
{
  "Zeioth/markmap.nvim",
  build = "yarn global add markmap-cli",
  cmd = { "MarkmapOpen", "MarkmapSave", "MarkmapWatch", "MarkmapWatchStop" },
  opts = {
    html_output = "/tmp/markmap.html",
    hide_toolbar = false,
    grace_period = 3600000,
  },
  config = function(_, opts) require("markmap").setup(opts) end
},
```

**Cross-platform**: Compatible with Linux, macOS, Windows, and Android Termux (matches Javi.Dots target platforms).

### 2. Obsidian Mindmap NextGen (Obsidian plugin, NOT Neovim)

**What it does**: Obsidian plugin that renders `markmap` code blocks inline inside Obsidian app. Supports per-block and per-file settings via YAML frontmatter.

**Relevance**: This is the Obsidian-app-side counterpart. If templates contain `markmap` code blocks, Obsidian users with this plugin get inline mind maps. Neovim users get them via markmap.nvim's `:MarkmapOpen`. Same markdown, two rendering paths.

**Code block syntax**:
````markdown
```markmap
# Root Topic
## Branch A
- Sub-item 1
- Sub-item 2
## Branch B
### Deep item
```
````

**Settings inheritance**: Global → File frontmatter → Code block (three-tier merge). Supports `coloring`, `height`, `initialExpandLevel`, custom colors.

### 3. Template Analysis — Which Templates Benefit?

| Template | Markmap Value | Rationale |
|----------|--------------|-----------|
| **consolidation** | **HIGH** | Connections between topics, patterns across notes — perfect tree structure |
| **adr** | **HIGH** | Decision tree with alternatives, consequences branching from the decision |
| **braindump** | **LOW** | Speed-focused capture, adding markmap adds friction |
| **resource-capture** | **LOW** | Linear content (source → summary → takeaways), not hierarchical |
| **sprint-review** | **MEDIUM** | Completed vs. not-completed with sub-items, retrospective connections |
| **weekly-brief** | **MEDIUM** | Highlights → Risks → Focus areas, moderately hierarchical |
| **tech-debt** | **MEDIUM** | Area → Impact → Plan hierarchy, useful for complex debt items |
| **daily-note** | **NONE** | Too temporal, no meaningful hierarchy |
| **coding-session** | **LOW** | Mostly linear narrative |
| **debug-journal** | **LOW** | Linear investigation flow |
| **meeting-notes** | **MEDIUM** | Topic-driven with action items branching |
| **sdd-feedback** | **LOW** | Structured feedback, not hierarchical |
| **stakeholder-update** | **LOW** | Mostly flat communication |
| **daily-brief** | **NONE** | Too temporal |
| **risk-registry** | **MEDIUM** | Risk categories with mitigations branching |
| **team-intelligence** | **MEDIUM** | Skills/gaps mapping is inherently tree-like |

### 4. AI Skill Updates — Value Assessment

**consolidation skill** (HIGH value): The consolidation workflow produces `## Top Insights`, `## Connections Made`, `## Technical Patterns` — all inherently hierarchical. A markmap at the end would visualize the knowledge graph structure the consolidation reveals. The skill already does entity extraction, so the markmap could use entities as leaves.

**braindump skill** (LOW value, skip): Braindumps prioritize speed. Adding a markmap code block adds visual noise to what should be a fast capture. Contradicts the "speed over perfection" rule in the skill.

**resource-capture skill** (LOW value, skip): Resources are linear (source → summary → takeaways). A markmap adds little insight.

**Recommendation**: Only update `obsidian-consolidation` skill. Add an optional `## Mind Map` section that generates a markmap code block from the consolidated insights and connections.

### 5. markmap-mcp-server Analysis

**What it would provide**: An MCP server that lets AI agents generate markmap HTML or markdown code blocks programmatically. Could be useful during SDD workflows to visualize architecture, dependency graphs, or task breakdowns.

**Current MCP ecosystem**: 8 servers focused on data/context (context7, engram), collaboration (atlassian, notion, figma), monitoring (sentry, cloudflare), and search (brave). No visualization/diagramming servers.

**Assessment**: Low priority for initial integration. The AI can already generate `markmap` code blocks as plain text — no MCP server needed for that. An MCP server would only add value if it rendered to images or interacted with a running markmap instance. The current `markmap-mcp-server` npm package is early-stage and its added value over inline generation is unclear.

**Recommendation**: OUT of scope for this change. Revisit if the package matures or if there's a concrete workflow that benefits from server-side rendering.

## Approaches

### 1. **Minimal: Neovim Plugin Only**
Add markmap.nvim to the plugin config. No template changes, no skill updates.

- Pros: Zero risk, immediately useful for any markdown file, no breaking changes
- Cons: Doesn't leverage the Obsidian Brain ecosystem, users must manually run `:MarkmapOpen`
- Effort: **Low** (single file addition)

### 2. **Integrated: Plugin + Templates + Consolidation Skill** (RECOMMENDED)
Add markmap.nvim plugin, add optional `markmap` code blocks to consolidation and ADR templates, update the consolidation skill to generate markmap blocks.

- Pros: Full ecosystem integration, Obsidian users with Mindmap NextGen see inline maps, Neovim users use `:MarkmapOpen`, AI generates useful visualizations during consolidation
- Cons: Slightly more complex, templates get longer, users without markmap rendering see raw code blocks
- Effort: **Medium** (plugin file + 2 template updates + 1 skill update)

### 3. **Full: Plugin + All Templates + All Skills + MCP**
Everything in Approach 2 plus markmap blocks in every template that scores MEDIUM+, update all 3 skills, add markmap-mcp-server.

- Pros: Maximum coverage
- Cons: Over-engineered, adds noise to braindumps (contradicts design philosophy), MCP server is immature, too many templates touched for marginal value
- Effort: **High** (many files, ongoing maintenance)

## Recommendation

**Approach 2: Integrated (Plugin + Templates + Consolidation Skill)**

This is the right scope. Here's why:

1. **markmap.nvim** fits the existing plugin pattern perfectly — lazy-loaded, zero conflict, cross-platform.
2. **Consolidation template** is the highest-value target because it synthesizes knowledge across notes — exactly where visual maps shine. ADR is second because decision trees are inherently hierarchical.
3. **Consolidation skill** already does entity extraction and pattern recognition — generating a markmap from that analysis is a natural extension, not forced.
4. **Braindumps stay fast** — the braindump skill's "speed over perfection" philosophy is preserved.
5. **MCP server deferred** — AI can generate markmap code blocks as plain text. No server needed until there's a concrete rendering/interaction use case.

### Scope Summary

**IN SCOPE:**
- New file: `GentlemanNvim/nvim/lua/plugins/markmap.lua`
- Update: `GentlemanNvim/obsidian-brain/core/templates/consolidation.md` — add optional markmap section
- Update: `GentlemanNvim/obsidian-brain/developer/templates/adr.md` — add optional markmap section
- Update: `GentlemanClaude/skills/obsidian-consolidation/SKILL.md` — generate markmap code block in output

**OUT OF SCOPE:**
- markmap-mcp-server (immature, unclear value-add over inline generation)
- braindump/resource-capture skill updates (contradicts speed philosophy)
- VS Code extension, web app, CLI install steps
- Installer TUI changes (markmap-cli is a yarn dependency handled by lazy.nvim build step)
- Templates that score LOW or NONE in the value assessment

## Risks

1. **yarn dependency** — markmap.nvim requires `yarn global add markmap-cli`. If yarn isn't installed, the plugin build step fails silently. Mitigation: the plugin uses `cmd` lazy-loading, so it won't error on startup — only when the user tries to use `:MarkmapOpen`.

2. **Raw code blocks for non-markmap users** — Users without Obsidian Mindmap NextGen plugin will see raw `markmap` code blocks in their notes. Mitigation: make the markmap section clearly labeled as optional/collapsible, and document that it requires the Obsidian plugin for rendering.

3. **Template bloat** — Adding markmap blocks increases template size. Current templates are 15-25 lines. Mitigation: keep markmap blocks as short commented-out examples or in a separate "visualization" section at the end, so they don't interfere with the core template flow.

4. **Skill output length** — The consolidation skill already produces substantial output. Adding a markmap block increases token usage. Mitigation: generate concise markmap (2-3 levels deep, entities as leaves), not a full reproduction of the consolidation content.

## Ready for Proposal

**Yes** — The exploration covers all 5 investigation areas with concrete evidence from the codebase. The recommended approach (Approach 2) has clear scope boundaries, identified risks with mitigations, and touches a manageable number of files (1 new + 3 updates). The orchestrator can proceed with `/sdd:new markmap-obsidian-brain` to create the proposal.
