# Proposal: Obsidian Brain Enhancement with Additive Role Packs

## Intent

The current Obsidian integration in Javi.Dots is minimal: `stepInstallNvim()` in `installer/internal/tui/installer.go:1019-1022` creates `~/.config/obsidian/` and `~/.config/obsidian/templates/` as empty directories. The `obsidian.nvim` plugin (`GentlemanNvim/nvim/lua/plugins/obsidian.lua`) configures a single workspace called "GentlemanNotes" pointing at `~/.config/obsidian` with an empty `templates/` subdir. The `workflow-obsidian-brain` skill exists in the AI framework module registry (`docs/ai-framework-modules.md:439`) but is installed via the external `project-starter-framework` repo, not shipped with Javi.Dots itself.

Users who select "Obsidian Brain" as their memory module during project init (`ScreenProjectMemory`) get a vault with zero templates, zero structure, and zero guidance. The skill referenced in `ai-framework-modules.md` only provides AI agent instructions, not actual vault templates or Obsidian configuration.

Adding role-based packs (Core, Developer, PM/Tech Lead) inspired by the COG Second Brain methodology would make Obsidian Brain actually useful from day one by shipping templates, folder structure, and role-aware AI skills directly with Javi.Dots.

## Scope

### In Scope

1. **Core vault templates** shipped in `GentlemanNvim/obsidian-brain/core/` (or similar path):
   - Braindump template (quick capture)
   - Resource capture template (link + summary + tags)
   - Knowledge consolidation template (weekly synthesis)
   - Daily note template
   - Folder structure: `inbox/`, `resources/`, `knowledge/`, `templates/`

2. **Developer role pack** in `GentlemanNvim/obsidian-brain/developer/`:
   - ADR (Architecture Decision Record) template
   - Coding session log template
   - Tech debt tracker template
   - Debug journal template
   - SDD feedback loop template (links to Engram when active)

3. **PM/Tech Lead role pack** in `GentlemanNvim/obsidian-brain/pm-lead/`:
   - Meeting notes template
   - Sprint review template
   - Stakeholder update template
   - Risk registry template
   - Daily/weekly brief templates
   - Team intelligence template

4. **New TUI screen** (`ScreenProjectRolePack`) for role pack selection:
   - Inserted into the project init flow: after `ScreenProjectEngram`, before `ScreenProjectCI`
   - Only shown when `ProjectMemory == "obsidian-brain"`
   - Multi-select: Core (always on), Developer, PM/Tech Lead
   - Options: `[x] Core (included)`, `[ ] Developer Pack`, `[ ] PM/Tech Lead Pack`

5. **New AI skills** in `GentlemanClaude/skills/`:
   - `obsidian-braindump/SKILL.md` - Braindump capture workflow skill
   - `obsidian-consolidation/SKILL.md` - Weekly knowledge consolidation skill
   - `obsidian-resource-capture/SKILL.md` - Resource capture and annotation skill
   - Skills adapt their output format based on which role packs are active

6. **CLI flag** `--project-role-pack=developer,pm-lead` for non-interactive mode

7. **Multi-workspace support** in `obsidian.nvim` config:
   - Personal vault: `~/.config/obsidian` (existing)
   - Project vault: `{project_path}/.obsidian-brain/` (new, when obsidian-brain selected)
   - Dynamic workspace registration based on project detection

8. **Model/state changes**:
   - New field `ProjectRolePacks []string` in `UserChoices` struct
   - New field `ProjectRolePacks []string` in `Model` struct
   - New Screen constant `ScreenProjectRolePack`

### Out of Scope

- COG-specific features (daily-brief, weekly-checkin, team-brief, competitive intelligence analysis)
- Obsidian desktop app plugin management (user installs community plugins manually)
- Modifying the `project-starter-framework` repo (this is Javi.Dots-only)
- Obsidian mobile/iCloud sync configuration
- Changes to Engram itself
- Dataview plugin queries (user-installed plugin, we only provide templates that are Dataview-ready)
- Obsidian community theme installation

## Approach

### Phase 1: Template Assets
Create the template markdown files and folder structure under a new `GentlemanNvim/obsidian-brain/` directory. Each role pack is a subdirectory. Templates use standard Obsidian frontmatter (YAML) and are compatible with Templater and core Templates plugins.

### Phase 2: TUI Integration
1. Add `ScreenProjectRolePack` constant to `model.go` (between `ScreenProjectEngram` and `ScreenProjectCI`)
2. Add `ProjectRolePacks []string` to both `UserChoices` and `Model` structs
3. Add multi-select handler in `update.go` for the new screen (similar to `ScreenAIToolsSelect` pattern)
4. Add rendering in `view.go` following the existing `renderSelection()` pattern
5. Wire the screen into the project init flow in `handleSelection()` — after Engram, before CI
6. Update `goBackInstallStep()` for back navigation
7. Update `renderProjectConfirm()` to display selected role packs

### Phase 3: Installer Logic
1. Extend `runProjectInitScript()` to pass `--role-pack=core,developer,pm-lead` to `init-project.sh`
2. Add template copying logic in `installer.go`: copy selected pack templates into `{project}/.obsidian-brain/templates/`
3. Create vault folder structure based on pack selection
4. Optionally generate `.obsidian/` config (vault settings)

### Phase 4: Obsidian.nvim Enhancement
1. Modify `GentlemanNvim/nvim/lua/plugins/obsidian.lua` to support multiple workspaces
2. Add dynamic workspace detection: if current directory contains `.obsidian-brain/`, register it
3. Keep existing personal vault as fallback workspace

### Phase 5: AI Skills
1. Create 3 new skills in `GentlemanClaude/skills/` following the existing frontmatter format (see `react-19/SKILL.md`)
2. Skills reference templates by name and adapt output based on role pack context
3. Register skills in `AGENTS.md` skill table

### Phase 6: CLI Support
1. Add `--project-role-pack` flag in `main.go` alongside existing `--project-memory` and `--project-engram` flags
2. Validate that role-pack requires `--project-memory=obsidian-brain`
3. Pass through to installer logic

## Affected Areas

### Files to Modify
- `installer/internal/tui/model.go` - Screen constant, UserChoices, Model struct, GetCurrentOptions, GetScreenTitle, GetScreenDescription
- `installer/internal/tui/update.go` - handleSelection(), goBackInstallStep(), handleKeyPress(), handleEscape()
- `installer/internal/tui/view.go` - View() switch case, renderProjectConfirm(), new render function
- `installer/internal/tui/installer.go` - stepInstallNvim() obsidian section, runProjectInitScript()
- `installer/cmd/gentleman-installer/main.go` - cliFlags struct, parseFlags(), runNonInteractive()
- `GentlemanNvim/nvim/lua/plugins/obsidian.lua` - workspaces configuration
- `AGENTS.md` - skill table registration
- `docs/ai-framework-modules.md` - document new skills (optional, since these are Javi.Dots-native)

### Files to Create
- `GentlemanNvim/obsidian-brain/core/templates/braindump.md`
- `GentlemanNvim/obsidian-brain/core/templates/resource-capture.md`
- `GentlemanNvim/obsidian-brain/core/templates/consolidation.md`
- `GentlemanNvim/obsidian-brain/core/templates/daily-note.md`
- `GentlemanNvim/obsidian-brain/developer/templates/adr.md`
- `GentlemanNvim/obsidian-brain/developer/templates/coding-session.md`
- `GentlemanNvim/obsidian-brain/developer/templates/tech-debt.md`
- `GentlemanNvim/obsidian-brain/developer/templates/debug-journal.md`
- `GentlemanNvim/obsidian-brain/developer/templates/sdd-feedback.md`
- `GentlemanNvim/obsidian-brain/pm-lead/templates/meeting-notes.md`
- `GentlemanNvim/obsidian-brain/pm-lead/templates/sprint-review.md`
- `GentlemanNvim/obsidian-brain/pm-lead/templates/stakeholder-update.md`
- `GentlemanNvim/obsidian-brain/pm-lead/templates/risk-registry.md`
- `GentlemanNvim/obsidian-brain/pm-lead/templates/daily-brief.md`
- `GentlemanNvim/obsidian-brain/pm-lead/templates/weekly-brief.md`
- `GentlemanNvim/obsidian-brain/pm-lead/templates/team-intelligence.md`
- `GentlemanClaude/skills/obsidian-braindump/SKILL.md`
- `GentlemanClaude/skills/obsidian-consolidation/SKILL.md`
- `GentlemanClaude/skills/obsidian-resource-capture/SKILL.md`

## Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| Template format incompatible with Obsidian versions | Low | Medium | Use standard YAML frontmatter, avoid plugin-specific syntax. Test with Obsidian 1.5+ |
| Multi-workspace obsidian.nvim breaks existing users | Medium | High | Guard new workspace detection behind a conditional check; keep existing "GentlemanNotes" workspace untouched as default |
| Screen constant ordering breaks enum iota | Medium | High | Insert `ScreenProjectRolePack` carefully between `ScreenProjectEngram` (line 79) and `ScreenProjectCI` (line 80). All subsequent constants auto-increment, but any hardcoded references to screen numbers would break. Verify no hardcoded screen values exist. |
| `runProjectInitScript` interface change | Low | Medium | The script already accepts `--engram` flag; adding `--role-pack` follows the same pattern. Falls back gracefully if script doesn't understand the flag yet. |
| Too many TUI screens in project init flow | Low | Low | Role pack screen only appears when obsidian-brain is selected. Most users won't see it. |
| Template content becomes stale or opinionated | Medium | Low | Templates should be minimal scaffolds (frontmatter + section headers), not prescriptive content. Users customize after. |

## Rollback Plan

1. **TUI changes**: Remove `ScreenProjectRolePack` constant and related handlers. The iota-based enum means removing the constant shifts all subsequent values, so this requires a single clean revert of model.go, update.go, view.go changes.
2. **Templates**: Delete `GentlemanNvim/obsidian-brain/` directory entirely. No other code depends on it until the installer copies files.
3. **obsidian.nvim**: Revert `obsidian.lua` to single-workspace config. The personal vault at `~/.config/obsidian` continues working.
4. **Skills**: Delete the 3 new skill directories from `GentlemanClaude/skills/`. They're standalone files with no dependencies.
5. **CLI flag**: Remove `--project-role-pack` from `main.go`. Existing flags are unaffected.

All changes are additive and self-contained. A full revert is a single `git revert` of the implementation commit(s).

## Dependencies

- **Existing infrastructure**: Obsidian Brain memory module flow already exists in the TUI (ScreenProjectMemory -> ScreenProjectObsidianInstall -> ScreenProjectEngram). This change extends that flow.
- **project-starter-framework**: The `init-project.sh` script may need updates to support `--role-pack` flag, but this can be handled gracefully (ignore unknown flags).
- **No new Go dependencies**: All changes use existing Bubbletea patterns already in the codebase.
- **No new npm/external dependencies**: Templates are plain Markdown files.
- **obsidian.nvim plugin**: Already a dependency (`obsidian-nvim/obsidian.nvim` in the Lua config). Multi-workspace is a built-in feature of the plugin.

## Success Criteria

1. **TUI flow works end-to-end**: User selects Obsidian Brain -> sees role pack selection -> selected packs appear in confirmation summary -> templates are copied to project
2. **Templates are valid Obsidian notes**: Each template opens correctly in Obsidian with proper frontmatter, and works with the core Templates plugin
3. **Multi-workspace detection**: Opening Neovim in a project with `.obsidian-brain/` auto-registers the project workspace in obsidian.nvim without breaking personal vault access
4. **CLI parity**: `--project-role-pack=developer,pm-lead` produces the same result as TUI selection
5. **Back navigation works**: ESC from role pack screen goes back to Engram/Memory screen correctly
6. **Skills load correctly**: The 3 new skills appear in the Skill Manager browse view and follow the existing SKILL.md frontmatter format
7. **No regression**: Existing installation flow (OS -> Terminal -> ... -> AI Tools) is completely unaffected. Project init without obsidian-brain skips the role pack screen entirely.
8. **E2E validation**: Docker-based E2E tests pass with `--project-memory=obsidian-brain --project-role-pack=developer`
