# Tasks: Obsidian Brain Enhancement with Additive Role Packs

## Dependency Graph

```
Phase 1 (Template Assets)     ──┐
Phase 2 (TUI Model Layer)     ──┤
                                ├──→ Phase 3 (TUI Update Logic) ──→ Phase 4 (TUI View Layer)
                                │         │
Phase 7 (Neovim Multi-workspace)         │
Phase 8 (AI Skills)                      │
                                         ▼
                              Phase 5 (Installer Logic) ──→ Phase 6 (CLI Support)
                                                                    │
                                                                    ▼
                                                         Phase 9 (Testing)
                                                                    │
                                                                    ▼
                                                         Phase 10 (Documentation)
```

**Parallelizable**: Phases 1, 2, 7, and 8 have NO interdependencies and can be implemented in any order or in parallel. Phase 3 depends on Phase 2. Phase 4 depends on Phase 3. Phase 5 depends on Phases 1 and 3. Phase 6 depends on Phase 5. Phase 9 depends on Phases 3, 5, and 6. Phase 10 depends on Phase 9.

---

## Phase 1: Template Assets (no code dependencies)

> Creates the Obsidian vault template markdown files and directory structure.
> No Go code dependencies — pure file creation.
> **Depends on**: Nothing
> **Blocked by**: Nothing

### Task 1.1: [x] Create Core templates directory structure

- **File**: `GentlemanNvim/obsidian-brain/core/templates/` (create directory)
- **Action**: Create nested directories `GentlemanNvim/obsidian-brain/core/templates/`
- **Spec ref**: `specs/templates/spec.md` — Requirement: Directory Structure

### Task 1.2: [x] Create `braindump.md` template

- **File**: `GentlemanNvim/obsidian-brain/core/templates/braindump.md` (create)
- **Action**: Create template with YAML frontmatter (`title`, `date`, `tags: ["braindump", "inbox"]`) and sections: `## Thought`, `## Context`, `## Related Notes`
- **Format**: Standard YAML frontmatter, no `<% %>` Templater syntax, no `{{}}` interpolation
- **Spec ref**: `specs/templates/spec.md` — Requirement: Core Templates, Scenario: Braindump template content

### Task 1.3: [x] Create `resource-capture.md` template

- **File**: `GentlemanNvim/obsidian-brain/core/templates/resource-capture.md` (create)
- **Action**: Create template with frontmatter `tags: ["resource"]` and sections: `## Source`, `## Summary`, `## Key Takeaways`
- **Spec ref**: `specs/templates/spec.md` — Scenario: Resource capture template content

### Task 1.4: [x] Create `consolidation.md` template

- **File**: `GentlemanNvim/obsidian-brain/core/templates/consolidation.md` (create)
- **Action**: Create template with frontmatter `tags: ["consolidation", "weekly"]` and sections: `## Period`, `## Top Insights`, `## Connections Made`, `## Open Questions`, `## Action Items`
- **Spec ref**: `specs/templates/spec.md` — Scenario: Consolidation template content

### Task 1.5: [x] Create `daily-note.md` template

- **File**: `GentlemanNvim/obsidian-brain/core/templates/daily-note.md` (create)
- **Action**: Create template with frontmatter `tags: ["daily"]` and sections: `## Today's Focus`, `## Notes`, `## Reflections`, `## Tomorrow`
- **Spec ref**: `specs/templates/spec.md` — Scenario: Daily note template content

### Task 1.6: [x] Create Developer templates directory structure

- **File**: `GentlemanNvim/obsidian-brain/developer/templates/` (create directory)
- **Action**: Create nested directories

### Task 1.7: [x] Create `adr.md` template

- **File**: `GentlemanNvim/obsidian-brain/developer/templates/adr.md` (create)
- **Action**: Create template with frontmatter `tags: ["adr", "architecture"]`, `status:` field, and sections: `## Context`, `## Decision`, `## Consequences`, `## Alternatives Considered`
- **Spec ref**: `specs/templates/spec.md` — Scenario: ADR template content

### Task 1.8: [x] Create `coding-session.md` template

- **File**: `GentlemanNvim/obsidian-brain/developer/templates/coding-session.md` (create)
- **Action**: Create template with sections: `## Goal`, `## What I Did`, `## Blockers`, `## Decisions Made`, `## Next Steps`
- **Spec ref**: `specs/templates/spec.md` — Requirement: Developer Pack Templates

### Task 1.9: [x] Create `tech-debt.md` template

- **File**: `GentlemanNvim/obsidian-brain/developer/templates/tech-debt.md` (create)
- **Action**: Create template with sections: `## Area`, `## Description`, `## Impact`, `## Effort Estimate`, `## Priority`, `## Plan`
- **Spec ref**: `specs/templates/spec.md` — Requirement: Developer Pack Templates

### Task 1.10: [x] Create `debug-journal.md` template

- **File**: `GentlemanNvim/obsidian-brain/developer/templates/debug-journal.md` (create)
- **Action**: Create template with sections: `## Symptom`, `## Hypothesis`, `## Investigation`, `## Root Cause`, `## Fix`, `## Lessons Learned`
- **Spec ref**: `specs/templates/spec.md` — Requirement: Developer Pack Templates

### Task 1.11: [x] Create `sdd-feedback.md` template

- **File**: `GentlemanNvim/obsidian-brain/developer/templates/sdd-feedback.md` (create)
- **Action**: Create template with sections: `## Change Name`, `## Phase Completed`, `## What Worked`, `## What Didn't`, `## Improvements`, `## Link to Engram` (with placeholder `<!-- Optional: engram:// link -->`)
- **Spec ref**: `specs/templates/spec.md` — Scenario: SDD feedback template references Engram

### Task 1.12: [x] Create PM/Tech Lead templates directory structure

- **File**: `GentlemanNvim/obsidian-brain/pm-lead/templates/` (create directory)
- **Action**: Create nested directories

### Task 1.13: [x] Create `meeting-notes.md` template

- **File**: `GentlemanNvim/obsidian-brain/pm-lead/templates/meeting-notes.md` (create)
- **Action**: Create template with frontmatter `tags: ["meeting"]` and sections: `## Attendees`, `## Agenda`, `## Discussion`, `## Decisions`, `## Action Items`
- **Spec ref**: `specs/templates/spec.md` — Scenario: Meeting notes template content

### Task 1.14: [x] Create `sprint-review.md` template

- **File**: `GentlemanNvim/obsidian-brain/pm-lead/templates/sprint-review.md` (create)
- **Action**: Create template with sections: `## Sprint Goal`, `## Completed`, `## Not Completed`, `## Metrics`, `## Retrospective Notes`
- **Spec ref**: `specs/templates/spec.md` — Requirement: PM/Tech Lead Pack Templates

### Task 1.15: [x] Create `stakeholder-update.md` template

- **File**: `GentlemanNvim/obsidian-brain/pm-lead/templates/stakeholder-update.md` (create)
- **Action**: Create template with sections: `## Summary`, `## Progress`, `## Risks`, `## Next Steps`, `## Ask`
- **Spec ref**: `specs/templates/spec.md` — Requirement: PM/Tech Lead Pack Templates

### Task 1.16: [x] Create `risk-registry.md` template

- **File**: `GentlemanNvim/obsidian-brain/pm-lead/templates/risk-registry.md` (create)
- **Action**: Create template with frontmatter `tags: ["risk"]` and structured table/sections for: Risk ID, Description, Likelihood, Impact, Mitigation, Owner, Status
- **Spec ref**: `specs/templates/spec.md` — Scenario: Risk registry template content

### Task 1.17: [x] Create `daily-brief.md` template

- **File**: `GentlemanNvim/obsidian-brain/pm-lead/templates/daily-brief.md` (create)
- **Action**: Create template with sections: `## Top 3 Priorities`, `## Blockers`, `## Key Decisions Needed`, `## FYI`
- **Spec ref**: `specs/templates/spec.md` — Requirement: PM/Tech Lead Pack Templates

### Task 1.18: [x] Create `weekly-brief.md` template

- **File**: `GentlemanNvim/obsidian-brain/pm-lead/templates/weekly-brief.md` (create)
- **Action**: Create template with sections: `## Highlights`, `## Metrics`, `## Risks & Blockers`, `## Next Week Focus`, `## Team Notes`
- **Spec ref**: `specs/templates/spec.md` — Requirement: PM/Tech Lead Pack Templates

### Task 1.19: [x] Create `team-intelligence.md` template

- **File**: `GentlemanNvim/obsidian-brain/pm-lead/templates/team-intelligence.md` (create)
- **Action**: Create template with sections: `## Observation`, `## Context`, `## Pattern`, `## Recommendation`, `## Shared With`
- **Spec ref**: `specs/templates/spec.md` — Requirement: PM/Tech Lead Pack Templates

---

## Phase 2: TUI Model Layer (foundation for TUI changes)

> Adds constants, struct fields, and option/title/description methods.
> **Depends on**: Nothing
> **Blocked by**: Nothing

### Task 2.1: Add `ScreenProjectRolePack` constant to iota block [x]

- **File**: `installer/internal/tui/model.go:80` (modify)
- **Action**: Insert `ScreenProjectRolePack` between `ScreenProjectEngram` (line 79) and `ScreenProjectCI` (currently line 80). The new constant takes the value that was `ScreenProjectCI`; all subsequent constants shift +1 automatically via `iota`.
- **Exact insertion point**: After line 79 (`ScreenProjectEngram`), before current line 80 (`ScreenProjectCI`)
- **Result**:
  ```go
  ScreenProjectEngram          // Yes/No: add Engram alongside Obsidian Brain
  ScreenProjectRolePack        // Multi-select: role packs for Obsidian Brain
  ScreenProjectCI              // Single-select: CI provider
  ```
- **Iota safety**: No hardcoded numeric `Screen` values exist (verified via grep — all references use constant names)
- **Spec ref**: `specs/tui/spec.md` — Requirement: Screen Constant Placement

### Task 2.2: Add `ProjectRolePacks []string` to `UserChoices` struct [x]

- **File**: `installer/internal/tui/model.go:143` (modify)
- **Action**: Add `ProjectRolePacks []string` field to `UserChoices` after `ProjectEngram bool` (line 143) and before `InstallObsidian bool` (line 144)
- **Exact position**: Between `ProjectEngram    bool` and `InstallObsidian  bool`
- **Spec ref**: `specs/tui/spec.md` — Requirement: State Fields; `design.md` — Interfaces: New UserChoices field

### Task 2.3: Add `ProjectRolePacks []string` to `Model` struct [x]

- **File**: `installer/internal/tui/model.go:216` (modify)
- **Action**: Add `ProjectRolePacks []string` field after `ProjectCI string` (line 216). This stores the collected role pack IDs after confirm.
- **Exact position**: After `ProjectCI        string` and before `ProjectLogLines  []string`
- **Spec ref**: `design.md` — Interfaces: New Model fields

### Task 2.4: Add `RolePackSelected []bool` to `Model` struct [x]

- **File**: `installer/internal/tui/model.go:216` (modify)
- **Action**: Add `RolePackSelected []bool` field after `ProjectRolePacks []string` (added in Task 2.3). This tracks toggle state for each role pack in `ScreenProjectRolePack` (index 0 = Developer, index 1 = PM/Tech Lead; Core is implicit/always-on).
- **Pattern**: Follows `AIToolSelected []bool` at line 203
- **Spec ref**: `specs/tui/spec.md` — Requirement: State Fields; `design.md` — `RolePackSelected []bool`

### Task 2.5: Add `GetCurrentOptions()` case for `ScreenProjectRolePack` [x]

- **File**: `installer/internal/tui/model.go` (modify, within `GetCurrentOptions()` at line 354)
- **Action**: Add a `case ScreenProjectRolePack:` in the switch statement, between `ScreenProjectEngram` and `ScreenProjectCI` cases. Return dynamic options with checkbox prefixes based on `RolePackSelected`:
  ```go
  case ScreenProjectRolePack:
      coreLabel := "[x] Core (always included)"
      devLabel := "[ ] Developer Pack"
      pmLabel := "[ ] PM/Tech Lead Pack"
      if m.RolePackSelected != nil && len(m.RolePackSelected) > 0 && m.RolePackSelected[0] {
          devLabel = "[x] Developer Pack"
      }
      if m.RolePackSelected != nil && len(m.RolePackSelected) > 1 && m.RolePackSelected[1] {
          pmLabel = "[x] PM/Tech Lead Pack"
      }
      return []string{coreLabel, devLabel, pmLabel, "─────────────", "✅ Confirm selection"}
  ```
- **Exact location**: After `case ScreenProjectEngram:` return (line 539-540) and before `case ScreenProjectCI:` (line 541)
- **Spec ref**: `specs/tui/spec.md` — Requirement: Options List

### Task 2.6: Add `GetScreenTitle()` case for `ScreenProjectRolePack` [x]

- **File**: `installer/internal/tui/model.go` (modify, within `GetScreenTitle()` at line 560)
- **Action**: Add `case ScreenProjectRolePack:` returning `"📦 Initialize Project — Role Packs"` between the `ScreenProjectEngram` case (line 677) and `ScreenProjectCI` case (line 678)
- **Spec ref**: `specs/tui/spec.md` — Requirement: Screen Title and Description

### Task 2.7: Add `GetScreenDescription()` case for `ScreenProjectRolePack` [x]

- **File**: `installer/internal/tui/model.go` (modify, within `GetScreenDescription()` at line 705)
- **Action**: Add `case ScreenProjectRolePack:` returning `"Select role packs for your Obsidian Brain vault"` between the `ScreenProjectEngram` case (line 755) and `ScreenProjectCI` case (line 757)
- **Spec ref**: `specs/tui/spec.md` — Requirement: Screen Title and Description

### Task 2.8: Initialize `ProjectRolePacks` and `RolePackSelected` in `NewModel()` [x]

- **File**: `installer/internal/tui/model.go` (modify, within `NewModel()` at line 239)
- **Action**: Add initialization in the project init section (after line 289 `ProjectCI: ""`):
  ```go
  ProjectRolePacks:   nil,
  RolePackSelected:   nil,
  ```
- **Note**: Zero-value nil is fine — these are only populated when entering the role pack screen

---

## Phase 3: TUI Update Logic (requires Phase 2)

> Implements navigation, multi-select toggle, and screen flow changes.
> **Depends on**: Phase 2 (screen constant, struct fields, options)
> **Blocked by**: Phase 2

### Task 3.1: Add `rolePackIDMap` variable [x]

- **File**: `installer/internal/tui/update.go` (modify)
- **Action**: Add a package-level variable near `aiToolIDMap` (line 1655):
  ```go
  // rolePackIDMap maps role pack option index to pack ID (0=developer, 1=pm-lead)
  var rolePackIDMap = []string{"developer", "pm-lead"}
  ```
- **Spec ref**: `design.md` — Role pack ID map

### Task 3.2: Add `handleRolePackKeys()` function [x]

- **File**: `installer/internal/tui/update.go` (modify, add new function)
- **Action**: Create `handleRolePackKeys(key string) (tea.Model, tea.Cmd)` following the `handleAIToolsKeys()` pattern at lines 2084-2135. Logic:
  - `up/k`: move cursor up, skip separators
  - `down/j`: move cursor down, skip separators
  - `enter` / `space`:
    - Cursor 0 (Core): no-op (always selected)
    - Cursor 1 (Developer): toggle `m.RolePackSelected[0]`
    - Cursor 2 (PM/Tech Lead): toggle `m.RolePackSelected[1]`
    - Cursor on separator: no-op
    - Cursor 4 (Confirm): collect selected packs into `m.ProjectRolePacks` (`["core"]` + selected), set `m.Screen = ScreenProjectCI`, reset `m.Cursor = 0`
  - `esc/backspace`: call `m.goBackInstallStep()`
- **Spec ref**: `specs/tui/spec.md` — Requirement: Multi-Select Behavior

### Task 3.3: Register `ScreenProjectRolePack` in `handleKeyPress()` switch [x]

- **File**: `installer/internal/tui/update.go:876-888` (modify)
- **Action**: Add `ScreenProjectRolePack` to the screen-specific key dispatch in `handleKeyPress()`. It should NOT be in the `handleSelectionKeys` list (line 887-888) — it needs its own handler. Add a new case:
  ```go
  case ScreenProjectRolePack:
      return m.handleRolePackKeys(key)
  ```
  Insert after the `ScreenProjectEngram` group (in the `handleSelectionKeys` list at line 888) — but note: `ScreenProjectRolePack` must be REMOVED from that list if it was included, and given its own `case` block, like `ScreenAIToolsSelect` at line 891.
- **Spec ref**: `specs/tui/spec.md` — Requirement: HandleKeyPress Registration

### Task 3.4: Add `ScreenProjectRolePack` to space-key exclusion list [x]

- **File**: `installer/internal/tui/update.go:844-867` (modify)
- **Action**: In the `if key == " "` block, add `ScreenProjectRolePack` to the switch cases where space should NOT activate leader mode (should pass through to the handler). Add it alongside `ScreenSkillInstall, ScreenSkillRemove`:
  ```go
  case ScreenSkillInstall, ScreenSkillRemove, ScreenProjectRolePack:
      // Multi-select screens: space toggles selection, pass through
  ```
  Alternatively, handle space in `handleRolePackKeys` — either approach works, but matching the AI tools pattern (which uses `"enter", " "` in the handler) requires space to reach the handler.
- **Spec ref**: `specs/tui/spec.md` — Requirement: HandleKeyPress Registration, Scenario: Space toggles selection

### Task 3.5: Modify `handleSelection()` for `ScreenProjectEngram` — forward to RolePack [x]

- **File**: `installer/internal/tui/update.go:1517-1520` (modify)
- **Action**: Change the `case ScreenProjectEngram:` handler to go to `ScreenProjectRolePack` instead of `ScreenProjectCI`. Current code (lines 1517-1520):
  ```go
  case ScreenProjectEngram:
      m.ProjectEngram = m.Cursor == 0
      m.Screen = ScreenProjectCI
      m.Cursor = 0
  ```
  Change to:
  ```go
  case ScreenProjectEngram:
      m.ProjectEngram = m.Cursor == 0
      m.Screen = ScreenProjectRolePack
      m.Cursor = 0
      m.RolePackSelected = make([]bool, len(rolePackIDMap))
  ```
- **Spec ref**: `specs/tui/spec.md` — Requirement: Forward Navigation, Scenario: Engram screen to RolePack transition

### Task 3.6: Update `goBackInstallStep()` — add `ScreenProjectRolePack` case [x]

- **File**: `installer/internal/tui/update.go:1319-1325` (modify)
- **Action**: Add a new case for `ScreenProjectRolePack` in `goBackInstallStep()`, inserted before the `ScreenProjectCI` case:
  ```go
  case ScreenProjectRolePack:
      m.Screen = ScreenProjectEngram
      m.Cursor = 0
      m.RolePackSelected = nil
      m.ProjectRolePacks = nil
  ```
- **Spec ref**: `specs/tui/spec.md` — Requirement: Backward Navigation, Scenario: ESC from RolePack screen

### Task 3.7: Update `goBackInstallStep()` for `ScreenProjectCI` — back to RolePack when obsidian-brain [x]

- **File**: `installer/internal/tui/update.go:1319-1325` (modify)
- **Action**: Change the `ScreenProjectCI` case. Current code (lines 1319-1325):
  ```go
  case ScreenProjectCI:
      if m.ProjectMemory == "obsidian-brain" {
          m.Screen = ScreenProjectEngram
      } else {
          m.Screen = ScreenProjectMemory
      }
      m.Cursor = 0
  ```
  Change to:
  ```go
  case ScreenProjectCI:
      if m.ProjectMemory == "obsidian-brain" {
          m.Screen = ScreenProjectRolePack
      } else {
          m.Screen = ScreenProjectMemory
      }
      m.Cursor = 0
  ```
- **Spec ref**: `specs/tui/spec.md` — Requirement: Backward Navigation, Scenario: ESC from CI screen with obsidian-brain

### Task 3.8: Add `ScreenProjectRolePack` to `handleEscape()` switch [x]

- **File**: `installer/internal/tui/update.go:1085-1097` (modify)
- **Action**: Add `ScreenProjectRolePack` to the project init screens ESC handling. In `handleEscape()`, the project init screens (lines 1086-1099) use `goBackInstallStep()` implicitly via the selection handler's `esc` case. However, `ScreenProjectRolePack` is handled by `handleRolePackKeys()`, not `handleSelectionKeys()`, so the `esc` case in `handleRolePackKeys()` already calls `goBackInstallStep()`. No additional change is needed in `handleEscape()` IF the handler catches `esc`. 
  
  BUT — looking at the escape handler structure: `handleEscape()` is called before screen-specific handlers (line 871-873). So `esc` is intercepted by `handleEscape()` first. We need to add `ScreenProjectRolePack` to the `handleEscape()` switch to route it to `goBackInstallStep()`. The appropriate place is alongside the other install wizard screens at line 1012:
  ```go
  case ScreenOSSelect, ScreenTerminalSelect, ScreenFontSelect, ScreenShellSelect,
       ScreenWMSelect, ScreenNvimSelect, ScreenZedSelect, ScreenAIToolsSelect,
       ScreenAIFrameworkConfirm, ScreenAIFrameworkPreset, ScreenAIFrameworkCategories,
       ScreenAIFrameworkCategoryItems, ScreenProjectRolePack:
      return m.goBackInstallStep()
  ```
  Note: `ScreenProjectRolePack` is NOT in the project init ESC block (lines 1085-1099) since those screens use `handleSelectionKeys` which has its own ESC. But since `handleRolePackKeys` does NOT catch `esc` via `handleKeyPress` (because `handleEscape()` runs first on `esc`), we need it in the `goBackInstallStep` list at line 1012.
- **Spec ref**: `specs/tui/spec.md` — Requirement: HandleEscape Registration

### Task 3.9: Update `handleSelection()` for `ScreenProjectConfirm` — pass role packs to choices [x]

- **File**: `installer/internal/tui/update.go:1530-1544` (modify)
- **Action**: In the `case ScreenProjectConfirm:` handler (line 1530), add `m.Choices.ProjectRolePacks = m.ProjectRolePacks` alongside the other choice assignments. Current code (lines 1531-1538):
  ```go
  if m.Cursor == 0 { // Confirm
      m.Choices.InitProject = true
      m.Choices.ProjectPath = m.ProjectPathInput
      m.Choices.ProjectStack = m.ProjectStack
      m.Choices.ProjectMemory = m.ProjectMemory
      m.Choices.ProjectCI = m.ProjectCI
      m.Choices.ProjectEngram = m.ProjectEngram
  ```
  Add after `m.Choices.ProjectEngram`:
  ```go
      m.Choices.ProjectRolePacks = m.ProjectRolePacks
  ```
- **Spec ref**: `specs/tui/spec.md` — Requirement: State Fields, Scenario: UserChoices populated before installation

### Task 3.10: Update `runProjectInit()` to pass role packs [x]

- **File**: `installer/internal/tui/update.go:293-302` (modify)
- **Action**: Capture `m.ProjectRolePacks` in `runProjectInit()` and pass it to the installer function. Current code:
  ```go
  func (m Model) runProjectInit() tea.Cmd {
      path := expandPath(m.ProjectPathInput)
      memory := m.ProjectMemory
      ci := m.ProjectCI
      engram := m.ProjectEngram
      return func() tea.Msg {
          err := runProjectInitScript(path, memory, ci, engram)
  ```
  Add `rolePacks := m.ProjectRolePacks` and update function call. The function signature change is handled in Phase 5.
- **Note**: This task is paired with Task 5.2 which extends `runProjectInitScript` signature

---

## Phase 4: TUI View Layer (requires Phase 3)

> Adds rendering for the role pack screen and updates confirmation display.
> **Depends on**: Phase 3 (screen flow must work before rendering)
> **Blocked by**: Phase 3

### Task 4.1: Add `ScreenProjectRolePack` case to `View()` switch [x]

- **File**: `installer/internal/tui/view.go:120` (modify)
- **Action**: Add `ScreenProjectRolePack` as a View case. Currently line 120:
  ```go
  case ScreenProjectStack, ScreenProjectMemory, ScreenProjectObsidianInstall, ScreenProjectEngram, ScreenProjectCI:
      s.WriteString(m.renderSelection())
  ```
  Since the role pack screen uses multi-select with checkboxes (like `renderAIToolSelection()`), it needs its own render function. Change to:
  ```go
  case ScreenProjectStack, ScreenProjectMemory, ScreenProjectObsidianInstall, ScreenProjectEngram, ScreenProjectCI:
      s.WriteString(m.renderSelection())
  case ScreenProjectRolePack:
      s.WriteString(m.renderRolePackSelection())
  ```
- **Spec ref**: `specs/tui/spec.md` — Requirement: View Rendering

### Task 4.2: Create `renderRolePackSelection()` function [x]

- **File**: `installer/internal/tui/view.go` (modify, add new function)
- **Action**: Create `renderRolePackSelection()` following the `renderAIToolSelection()` pattern (lines 296-345). The function should:
  1. Render title via `m.GetScreenTitle()`
  2. Render description via `m.GetScreenDescription()`
  3. Render options with checkbox prefixes (`[x]`/`[ ]`):
     - Core always shows `[x]` (hardcoded)
     - Developer/PM show checkbox based on `m.RolePackSelected`
     - Separator rendered in `MutedStyle`
     - "Confirm selection" has no checkbox
  4. Render help text: `"↑/k up • ↓/j down • [Enter] toggle/confirm • [Esc] back"`
  5. Use cursor indicator `"▸ "` for selected row, `"  "` for others
  6. Use `SelectedStyle` / `UnselectedStyle` for row highlighting
- **Spec ref**: `specs/tui/spec.md` — Requirement: View Rendering, Scenario: Render with checkboxes

### Task 4.3: Update `renderProjectConfirm()` to display role packs [x]

- **File**: `installer/internal/tui/view.go:2345-2382` (modify)
- **Action**: Add a "Packs:" line to the confirmation summary. Insert after the Engram line (line 2362) and before the CI line (line 2363). Only show when `m.ProjectMemory == "obsidian-brain"`:
  ```go
  if m.ProjectMemory == "obsidian-brain" {
      engram := "No"
      if m.ProjectEngram {
          engram = "Yes"
      }
      s.WriteString(fmt.Sprintf("    Engram:  %s\n", engram))
      if len(m.ProjectRolePacks) > 0 {
          s.WriteString(fmt.Sprintf("    Packs:   %s\n", strings.Join(m.ProjectRolePacks, ", ")))
      }
  }
  ```
- **Spec ref**: `specs/tui/spec.md` — Requirement: Confirmation Screen Display, Scenario: Role packs shown in confirmation

---

## Phase 5: Installer Logic (requires Phases 1 and 3)

> Implements template copying and vault folder creation.
> **Depends on**: Phase 1 (templates must exist to copy), Phase 3 (role packs passed from TUI)
> **Blocked by**: Phases 1, 3

### Task 5.1: [x] Add `copyRolePackTemplates()` function

- **File**: `installer/internal/tui/installer.go` (modify, add new function)
- **Action**: Create `copyRolePackTemplates(repoDir, projectPath string, rolePacks []string) error` following the design in `design.md` — copyRolePackTemplates function. Logic:
  1. Create vault base dir: `{projectPath}/.obsidian-brain/`
  2. Create Core folders: `inbox/`, `resources/`, `knowledge/`, `templates/`
  3. Create `.obsidian/` marker dir inside `.obsidian-brain/` (for plugin detection)
  4. If "developer" in rolePacks: create `architecture/`, `sessions/`, `debugging/`
  5. If "pm-lead" in rolePacks: create `meetings/`, `sprints/`, `risks/`, `briefs/`
  6. For each pack in rolePacks: copy template files from `{repoDir}/GentlemanNvim/obsidian-brain/{pack}/templates/` to `{projectPath}/.obsidian-brain/templates/`
  7. Use `system.EnsureDir()` for directory creation (existing pattern at line 1021)
  8. Use `os.ReadFile` / `os.WriteFile` for template copying (design spec shows this pattern)
  9. Skip gracefully if source pack dir doesn't exist
- **Spec ref**: `specs/templates/spec.md` — Requirement: Folder Structure Creation; `design.md` — copyRolePackTemplates function

### Task 5.2: [x] Extend `runProjectInitScript()` to accept and use role packs

- **File**: `installer/internal/tui/installer.go:1700-1759` (modify)
- **Action**: Change function signature from `runProjectInitScript(projectPath, memory, ci string, engram bool)` to `runProjectInitScript(projectPath, memory, ci string, engram bool, rolePacks []string)`. After the existing `init-project.sh` execution, call `copyRolePackTemplates()` when `memory == "obsidian-brain"` and `len(rolePacks) > 0`:
  ```go
  // After existing script execution...
  if memory == "obsidian-brain" && len(rolePacks) > 0 {
      // Find repo dir for template source
      repoDir := findRepoDir()  // or pass it through
      if err := copyRolePackTemplates(repoDir, projectPath, rolePacks); err != nil {
          return fmt.Errorf("failed to copy role pack templates: %w", err)
      }
  }
  ```
  Also update `RunProjectInitScript()` (line 1757) public wrapper to match.
- **Note**: Must also update `runProjectInit()` in `update.go` (Task 3.10) to pass role packs

### Task 5.3: [x] Update `runProjectInit()` in `update.go` to pass role packs

- **File**: `installer/internal/tui/update.go:293-302` (modify)
- **Action**: Pass `m.ProjectRolePacks` to `runProjectInitScript()`:
  ```go
  func (m Model) runProjectInit() tea.Cmd {
      path := expandPath(m.ProjectPathInput)
      memory := m.ProjectMemory
      ci := m.ProjectCI
      engram := m.ProjectEngram
      rolePacks := m.ProjectRolePacks
      return func() tea.Msg {
          err := runProjectInitScript(path, memory, ci, engram, rolePacks)
          return projectInstallCompleteMsg{err: err}
      }
  }
  ```

---

## Phase 6: CLI Support (requires Phase 5)

> Adds `--project-role-pack` flag for non-interactive mode.
> **Depends on**: Phase 5 (installer function must accept role packs)
> **Blocked by**: Phase 5

### Task 6.1: [x] Add `projectRolePack` field to `cliFlags` struct

- **File**: `installer/cmd/gentleman-installer/main.go:17-44` (modify)
- **Action**: Add `projectRolePack string` to the `cliFlags` struct after `projectEngram bool` (line 39):
  ```go
  projectEngram    bool
  projectRolePack  string  // comma-separated: "developer,pm-lead"
  ```
- **Spec ref**: `specs/cli/spec.md` — Requirement: Flag Definition

### Task 6.2: [x] Register `--project-role-pack` flag in `parseFlags()`

- **File**: `installer/cmd/gentleman-installer/main.go:46-81` (modify)
- **Action**: Add `flag.StringVar` registration after the `--project-engram` registration (line 73):
  ```go
  flag.StringVar(&flags.projectRolePack, "project-role-pack", "",
      "Role packs for Obsidian Brain: developer,pm-lead (comma-separated)")
  ```
- **Spec ref**: `specs/cli/spec.md` — Requirement: Flag Definition

### Task 6.3: [x] Add validation and parsing in `runNonInteractive()`

- **File**: `installer/cmd/gentleman-installer/main.go:144-197` (modify)
- **Action**: After the engram validation block (line 178-180), add role pack validation:
  ```go
  // Validate role pack requires obsidian-brain
  if flags.projectRolePack != "" && memory != "obsidian-brain" {
      return fmt.Errorf("--project-role-pack requires --project-memory=obsidian-brain")
  }
  
  // Parse and validate role packs
  var rolePacks []string
  if flags.projectRolePack != "" {
      validPacks := map[string]bool{"developer": true, "pm-lead": true}
      for _, pack := range strings.Split(flags.projectRolePack, ",") {
          pack = strings.TrimSpace(strings.ToLower(pack))
          if pack == "core" {
              continue // core is implicit
          }
          if !validPacks[pack] {
              return fmt.Errorf("invalid role pack: %s (valid: developer, pm-lead)", pack)
          }
          rolePacks = append(rolePacks, pack)
      }
  }
  // Core is always included when obsidian-brain
  if memory == "obsidian-brain" {
      rolePacks = append([]string{"core"}, rolePacks...)
  }
  ```
- **Spec ref**: `specs/cli/spec.md` — Requirements: Validation, Accepted Values, Core Always Included

### Task 6.4: [x] Pass role packs to `RunProjectInitScript()` in non-interactive flow

- **File**: `installer/cmd/gentleman-installer/main.go:192` (modify)
- **Action**: Update the call to `tui.RunProjectInitScript()` to pass role packs:
  ```go
  if err := tui.RunProjectInitScript(absPath, memory, ci, flags.projectEngram, rolePacks); err != nil {
  ```
- **Note**: Requires `RunProjectInitScript()` public API to be updated in Task 5.2

### Task 6.5: [x] Display role packs in non-interactive summary

- **File**: `installer/cmd/gentleman-installer/main.go:182-189` (modify)
- **Action**: After the Engram summary line (line 187-188), add:
  ```go
  if len(rolePacks) > 0 {
      fmt.Printf("  Packs:    %s\n", strings.Join(rolePacks, ", "))
  }
  ```
- **Spec ref**: `specs/cli/spec.md` — Requirement: Non-interactive Output

### Task 6.6: [x] Add `--project-role-pack` to help text in `printHelp()`

- **File**: `installer/cmd/gentleman-installer/main.go:424-509` (modify)
- **Action**: Add a line in the "Project Init Options" section (after line 467):
  ```
  --project-role-pack=<p>  Role packs: developer,pm-lead (comma-separated, core always included)
  ```
  Also add an example:
  ```
  # Initialize with Obsidian Brain + developer pack
  gentleman.dots --non-interactive --init-project --project-path=/path --project-memory=obsidian-brain --project-role-pack=developer
  ```
- **Spec ref**: `specs/cli/spec.md` — Requirement: Help Text

---

## Phase 7: Neovim Multi-workspace (no Go dependencies)

> Modifies obsidian.lua for project vault detection.
> **Depends on**: Nothing (Lua config, independent of Go code)
> **Blocked by**: Nothing

### Task 7.1: Add dynamic workspace detection to `obsidian.lua` [x]

- **File**: `GentlemanNvim/nvim/lua/plugins/obsidian.lua` (modify)
- **Action**: Add project vault detection using `vim.fn.finddir()`. The `opts` table should be computed dynamically. Replace the static `workspaces` table with a function that:
  1. Always includes `GentlemanNotes` workspace (existing personal vault at `~/.config/obsidian`)
  2. Checks `vim.fn.finddir('.obsidian-brain', vim.fn.getcwd() .. ';')` for upward directory search
  3. If found, adds a second workspace entry with `name` derived from the parent directory name and `path` pointing to the `.obsidian-brain/` directory
  4. Falls back silently to personal vault only if detection fails or `.obsidian-brain/` not found
  
  Implementation approach — compute `workspaces` in the `opts` function:
  ```lua
  opts = function()
    local workspaces = {
      {
        name = "GentlemanNotes",
        path = os.getenv("HOME") .. "/.config/obsidian",
      },
    }
    -- Detect project vault
    local brain_dir = vim.fn.finddir('.obsidian-brain', vim.fn.getcwd() .. ';')
    if brain_dir ~= '' then
      local abs_path = vim.fn.fnamemodify(brain_dir, ':p')
      local project_name = vim.fn.fnamemodify(abs_path, ':h:t')
      table.insert(workspaces, {
        name = project_name,
        path = abs_path,
      })
    end
    return {
      legacy_commands = false,
      workspaces = workspaces,
      -- ... rest of existing opts unchanged ...
    }
  end,
  ```
- **Key constraint**: Must NOT break existing personal vault behavior. No errors if `.obsidian-brain/` doesn't exist.
- **Spec ref**: `specs/neovim/spec.md` — All requirements

---

## Phase 8: AI Skills (no Go dependencies)

> Creates three new SKILL.md files and registers them in AGENTS.md.
> **Depends on**: Nothing (pure file creation)
> **Blocked by**: Nothing

### Task 8.1: Create `obsidian-braindump/SKILL.md`

- [x] **File**: `GentlemanClaude/skills/obsidian-braindump/SKILL.md` (create)
- **Action**: Create skill file with YAML frontmatter matching existing format (e.g., `react-19/SKILL.md`):
  - `name: obsidian-braindump`
  - `description:` with trigger keywords ("braindump", "capture thought", "quick note")
  - `license: Apache-2.0`
  - `metadata.author: gentleman-programming`, `metadata.version: "1.0"`
  - Body: workflow description, output format spec (braindump template sections), example, role-awareness section
- **Spec ref**: `specs/skills/spec.md` — Requirement: obsidian-braindump Skill

### Task 8.2: Create `obsidian-consolidation/SKILL.md`

- [x] **File**: `GentlemanClaude/skills/obsidian-consolidation/SKILL.md` (create)
- **Action**: Create skill file with:
  - `name: obsidian-consolidation`
  - `description:` with trigger keywords ("consolidate", "weekly synthesis", "knowledge review")
  - Body: weekly synthesis workflow, consolidation template format, wiki-link usage examples, role-awareness section
- **Spec ref**: `specs/skills/spec.md` — Requirement: obsidian-consolidation Skill

### Task 8.3: Create `obsidian-resource-capture/SKILL.md`

- [x] **File**: `GentlemanClaude/skills/obsidian-resource-capture/SKILL.md` (create)
- **Action**: Create skill file with:
  - `name: obsidian-resource-capture`
  - `description:` with trigger keywords ("capture resource", "save link", "bookmark")
  - Body: resource capture workflow, template format with URL/summary/takeaways, example, role-awareness section
- **Spec ref**: `specs/skills/spec.md` — Requirement: obsidian-resource-capture Skill

### Task 8.4: Register 3 new skills in `AGENTS.md`

- [x] **File**: `AGENTS.md` (modify)
- **Action**: Add 3 rows to the "Generic Skills" table (after existing entries like `zustand-5`):
  ```markdown
  | `obsidian-braindump` | Braindump capture workflow for Obsidian Brain vaults | [SKILL.md](GentlemanClaude/skills/obsidian-braindump/SKILL.md) |
  | `obsidian-consolidation` | Weekly knowledge consolidation for Obsidian Brain | [SKILL.md](GentlemanClaude/skills/obsidian-consolidation/SKILL.md) |
  | `obsidian-resource-capture` | Resource capture and annotation for Obsidian Brain | [SKILL.md](GentlemanClaude/skills/obsidian-resource-capture/SKILL.md) |
  ```
- **Spec ref**: `specs/skills/spec.md` — Requirement: AGENTS.md Registration

---

## Phase 9: Testing (requires Phases 3, 5, 6)

> Adds unit tests for the new functionality.
> **Depends on**: Phases 3, 5, 6 (code must exist to test)
> **Blocked by**: Phases 3, 5, 6

### Task 9.1: [x] Add `TestRolePackScreenOptions` test

- **File**: `installer/internal/tui/project_screens_test.go` (modify)
- **Action**: Add test verifying `GetCurrentOptions()` returns exactly 5 items for `ScreenProjectRolePack`. Follow the pattern of `TestObsidianInstallScreenOptions` (lines 311-339):
  ```go
  func TestRolePackScreenOptions(t *testing.T) {
      t.Run("has 5 options", func(t *testing.T) {
          m := NewModel()
          m.Screen = ScreenProjectRolePack
          m.RolePackSelected = make([]bool, 2)
          opts := m.GetCurrentOptions()
          if len(opts) != 5 {
              t.Errorf("expected 5 options, got %d: %v", len(opts), opts)
          }
      })
      t.Run("title is non-empty", func(t *testing.T) {
          m := NewModel()
          m.Screen = ScreenProjectRolePack
          title := m.GetScreenTitle()
          if title == "" {
              t.Error("expected non-empty title")
          }
      })
      t.Run("description mentions role", func(t *testing.T) {
          m := NewModel()
          m.Screen = ScreenProjectRolePack
          desc := m.GetScreenDescription()
          if !strings.Contains(strings.ToLower(desc), "role") {
              t.Errorf("expected description to mention role, got %q", desc)
          }
      })
  }
  ```
- **Spec ref**: `specs/tui/spec.md` — Requirement: Options List

### Task 9.2: [x] Add `TestRolePackCoreAlwaysOn` test

- **File**: `installer/internal/tui/project_screens_test.go` (modify)
- **Action**: Test that pressing Enter on Core (cursor 0) does not toggle anything:
  ```go
  func TestRolePackCoreAlwaysOn(t *testing.T) {
      m := NewModel()
      m.Screen = ScreenProjectRolePack
      m.RolePackSelected = make([]bool, 2)
      m.Cursor = 0 // Core
      result, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
      nm := result.(Model)
      // Should stay on same screen, no toggle effect
      if nm.Screen != ScreenProjectRolePack {
          t.Errorf("expected ScreenProjectRolePack, got %d", nm.Screen)
      }
  }
  ```
- **Spec ref**: `specs/tui/spec.md` — Scenario: Core cannot be toggled

### Task 9.3: [x] Add `TestRolePackToggle` test

- **File**: `installer/internal/tui/project_screens_test.go` (modify)
- **Action**: Test that pressing Enter on Developer (cursor 1) toggles `RolePackSelected[0]`:
  ```go
  func TestRolePackToggle(t *testing.T) {
      m := NewModel()
      m.Screen = ScreenProjectRolePack
      m.RolePackSelected = make([]bool, 2)
      m.Cursor = 1 // Developer Pack
      result, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
      nm := result.(Model)
      if !nm.RolePackSelected[0] {
          t.Error("expected Developer Pack to be toggled on")
      }
      // Toggle off
      nm.Cursor = 1
      result2, _ := nm.Update(tea.KeyMsg{Type: tea.KeyEnter})
      nm2 := result2.(Model)
      if nm2.RolePackSelected[0] {
          t.Error("expected Developer Pack to be toggled off")
      }
  }
  ```
- **Spec ref**: `specs/tui/spec.md` — Scenario: Toggle Developer pack on/off

### Task 9.4: [x] Add `TestRolePackConfirm` test

- **File**: `installer/internal/tui/project_screens_test.go` (modify)
- **Action**: Test that confirming with selections advances to `ScreenProjectCI` and sets `ProjectRolePacks`:
  ```go
  func TestRolePackConfirm(t *testing.T) {
      m := NewModel()
      m.Screen = ScreenProjectRolePack
      m.RolePackSelected = []bool{true, false} // Developer on, PM off
      m.Cursor = 4 // Confirm
      result, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
      nm := result.(Model)
      if nm.Screen != ScreenProjectCI {
          t.Errorf("expected ScreenProjectCI, got %d", nm.Screen)
      }
      // Should have core + developer
      expected := []string{"core", "developer"}
      // ... verify nm.ProjectRolePacks matches expected
  }
  ```
- **Spec ref**: `specs/tui/spec.md` — Scenario: Confirm with selections

### Task 9.5: [x] Add `TestRolePackBackNav` test

- **File**: `installer/internal/tui/project_screens_test.go` (modify)
- **Action**: Test that backspace from `ScreenProjectRolePack` goes to `ScreenProjectEngram`:
  ```go
  func TestRolePackBackNav(t *testing.T) {
      m := NewModel()
      m.Screen = ScreenProjectRolePack
      m.RolePackSelected = make([]bool, 2)
      result, _ := m.Update(tea.KeyMsg{Type: tea.KeyEsc})
      nm := result.(Model)
      if nm.Screen != ScreenProjectEngram {
          t.Errorf("expected ScreenProjectEngram, got %d", nm.Screen)
      }
      if nm.RolePackSelected != nil {
          t.Error("expected RolePackSelected to be nil after back nav")
      }
  }
  ```
- **Spec ref**: `specs/tui/spec.md` — Scenario: ESC from RolePack screen

### Task 9.6: [x] Add `TestCIBackNavWithRolePack` test

- **File**: `installer/internal/tui/project_screens_test.go` (modify)
- **Action**: Test that backspace from `ScreenProjectCI` with obsidian-brain goes to `ScreenProjectRolePack`:
  ```go
  func TestCIBackNavWithRolePack(t *testing.T) {
      m := NewModel()
      m.Screen = ScreenProjectCI
      m.ProjectMemory = "obsidian-brain"
      result, _ := m.Update(tea.KeyMsg{Type: tea.KeyBackspace})
      nm := result.(Model)
      if nm.Screen != ScreenProjectRolePack {
          t.Errorf("expected ScreenProjectRolePack, got %d", nm.Screen)
      }
  }
  ```
- **Spec ref**: `specs/tui/spec.md` — Scenario: ESC from CI screen with obsidian-brain

### Task 9.7: [x] Add `TestEngramForwardToRolePack` test

- **File**: `installer/internal/tui/project_screens_test.go` (modify)
- **Action**: Test that selecting on `ScreenProjectEngram` goes to `ScreenProjectRolePack` (not `ScreenProjectCI`):
  ```go
  func TestEngramForwardToRolePack(t *testing.T) {
      m := NewModel()
      m.Screen = ScreenProjectEngram
      m.Cursor = 0 // "Yes, add Engram too"
      result, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
      nm := result.(Model)
      if nm.Screen != ScreenProjectRolePack {
          t.Errorf("expected ScreenProjectRolePack, got %d", nm.Screen)
      }
      if nm.RolePackSelected == nil || len(nm.RolePackSelected) != 2 {
          t.Error("expected RolePackSelected to be initialized with 2 elements")
      }
  }
  ```
- **Spec ref**: `specs/tui/spec.md` — Scenario: Engram screen to RolePack transition

### Task 9.8: [x] Add `TestCLIRolePackValidation` test

- **File**: `installer/cmd/gentleman-installer/main_test.go` (create or modify if exists)
- **Action**: Test that `--project-role-pack` validation catches invalid packs and requires obsidian-brain. If no test file exists for main.go, add validation tests for the parsing logic. Alternatively, test via the public `RunProjectInitScript` function behavior.
- **Spec ref**: `specs/cli/spec.md` — Requirement: Validation

### Task 9.9: [x] Add `TestCopyRolePackTemplates` test

- **File**: `installer/internal/tui/installer_test.go` (create or modify if exists)
- **Action**: Test that `copyRolePackTemplates()` creates expected directory structure and copies template files. Use `t.TempDir()` for isolated testing:
  1. Set up a mock repo dir with template files
  2. Call `copyRolePackTemplates(repoDir, projectDir, []string{"core", "developer"})`
  3. Verify `{projectDir}/.obsidian-brain/inbox/` exists
  4. Verify `{projectDir}/.obsidian-brain/templates/braindump.md` exists
  5. Verify `{projectDir}/.obsidian-brain/.obsidian/` marker exists
- **Spec ref**: `specs/templates/spec.md` — Requirement: Folder Structure Creation

---

## Phase 10: Documentation (requires Phase 9)

> Updates documentation to reflect new features.
> **Depends on**: Phase 9 (everything should be tested first)
> **Blocked by**: Phase 9

### Task 10.1: [x] Update `docs/ai-framework-modules.md` with role pack info

- **File**: `docs/ai-framework-modules.md` (modify, if exists)
- **Action**: Add a section documenting the role packs: Core, Developer, PM/Tech Lead. List the templates in each pack and how they integrate with Obsidian Brain.
- **Note**: This file may not exist; check first. If it doesn't, this task is skipped.

### Task 10.2: [x] Update `AGENTS.md` with skill auto-invoke triggers

- **File**: `AGENTS.md` (modify)
- **Action**: In the "Auto-invoke Skills" table, add entries for the 3 new Obsidian skills if they should be auto-invoked. Since they are workflow-guidance skills (not codebase-specific), they likely don't need auto-invoke entries — just the skill table registration done in Task 8.4.
- **Note**: May be a no-op if auto-invoke is not appropriate for these skills

---

## Summary

| Phase | Tasks | Files Modified | Files Created | Depends On |
|-------|-------|---------------|---------------|------------|
| 1. Template Assets | 19 | 0 | 16 template files + 3 dirs | Nothing |
| 2. TUI Model | 8 | 1 (`model.go`) | 0 | Nothing |
| 3. TUI Update | 10 | 1 (`update.go`) | 0 | Phase 2 |
| 4. TUI View | 3 | 1 (`view.go`) | 0 | Phase 3 |
| 5. Installer Logic | 3 | 2 (`installer.go`, `update.go`) | 0 | Phases 1, 3 |
| 6. CLI Support | 6 | 1 (`main.go`) | 0 | Phase 5 |
| 7. Neovim Multi-workspace | 1 | 1 (`obsidian.lua`) | 0 | Nothing |
| 8. AI Skills | 4 | 1 (`AGENTS.md`) | 3 SKILL.md files | Nothing |
| 9. Testing | 9 | 1-2 test files | 0-1 | Phases 3, 5, 6 |
| 10. Documentation | 2 | 1-2 doc files | 0 | Phase 9 |
| **Total** | **65** | **~10** | **~22** | — |

## Recommended Implementation Order

1. **Wave 1** (parallel): Phases 1, 2, 7, 8 — no interdependencies
2. **Wave 2** (sequential after Wave 1): Phase 3 — depends on Phase 2
3. **Wave 3** (sequential after Wave 2): Phase 4 — depends on Phase 3
4. **Wave 4** (sequential after Waves 1+2): Phase 5 — depends on Phases 1 and 3
5. **Wave 5** (sequential after Wave 4): Phase 6 — depends on Phase 5
6. **Wave 6** (sequential after Waves 2+4+5): Phase 9 — depends on Phases 3, 5, 6
7. **Wave 7** (final): Phase 10 — depends on Phase 9
