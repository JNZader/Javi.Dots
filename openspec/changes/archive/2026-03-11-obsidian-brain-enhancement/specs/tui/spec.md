# TUI Flow Specification

## Purpose

Defines the behavior of the new `ScreenProjectRolePack` screen and its integration into the existing project initialization flow (model, update, view). Covers screen constant placement, state fields, navigation, multi-select interaction, and confirmation display.

## Requirements

### Requirement: Screen Constant Placement

The system MUST define `ScreenProjectRolePack` as a new `Screen` constant in `model.go`, inserted between `ScreenProjectEngram` (currently line 79) and `ScreenProjectCI` (currently line 80) in the `iota` block.

#### Scenario: Correct iota ordering

- GIVEN the Screen constant block uses `iota` starting at `ScreenWelcome`
- WHEN `ScreenProjectRolePack` is inserted between `ScreenProjectEngram` and `ScreenProjectCI`
- THEN `ScreenProjectRolePack` receives the value previously held by `ScreenProjectCI`, and all subsequent constants auto-increment by 1
- AND no existing code references screen constants by hardcoded integer value (verified: all references use symbolic names)

#### Scenario: No hardcoded screen number references

- GIVEN the codebase uses symbolic `Screen` constants everywhere (switch cases, assignments)
- WHEN a new constant is inserted into the iota block
- THEN no existing behavior changes because no code compares `Screen` values to literal integers

### Requirement: State Fields

The system MUST add a `ProjectRolePacks []string` field to both the `UserChoices` struct and the `Model` struct in `model.go`. The system MUST add a `RolePackSelected []bool` field to `Model` for tracking toggle state.

#### Scenario: Field initialization

- GIVEN a new `Model` is created via `NewModel()`
- WHEN the model is initialized
- THEN `ProjectRolePacks` is `nil` (empty slice) and `RolePackSelected` is `nil`

#### Scenario: UserChoices populated before installation

- GIVEN the user confirms project initialization on `ScreenProjectConfirm`
- WHEN the confirm handler runs (Cursor == 0)
- THEN `m.Choices.ProjectRolePacks` is set from `m.ProjectRolePacks`

### Requirement: Options List

The system MUST provide the following options for `ScreenProjectRolePack` in `GetCurrentOptions()`:

1. `"[x] Core (always included)"` (non-toggleable, always shown as selected)
2. `"[ ] Developer Pack"` (toggleable)
3. `"[ ] PM/Tech Lead Pack"` (toggleable)
4. `"─────────────"` (separator)
5. `"✅ Confirm selection"` (action)

#### Scenario: Default option display

- GIVEN the user navigates to `ScreenProjectRolePack`
- WHEN `GetCurrentOptions()` is called
- THEN it returns exactly 5 items: Core (pre-selected), Developer, PM/Tech Lead, separator, Confirm

#### Scenario: Toggled option display

- GIVEN `RolePackSelected[1]` is `true` (Developer toggled on)
- WHEN `GetCurrentOptions()` is called
- THEN the Developer option shows `"[x] Developer Pack"` and PM/Tech Lead shows `"[ ] PM/Tech Lead Pack"`

### Requirement: Screen Title and Description

The system MUST return `"📦 Initialize Project — Role Packs"` for `GetScreenTitle()` and `"Select role packs for your Obsidian Brain vault"` for `GetScreenDescription()` when `Screen == ScreenProjectRolePack`.

#### Scenario: Title displayed

- GIVEN `m.Screen == ScreenProjectRolePack`
- WHEN `GetScreenTitle()` is called
- THEN it returns `"📦 Initialize Project — Role Packs"`

### Requirement: Multi-Select Behavior

The system MUST use a toggle-on-enter interaction for `ScreenProjectRolePack`, following the pattern established by `ScreenAIToolsSelect`. Pressing Enter on Developer or PM/Tech Lead toggles their selection state. Pressing Enter on Core does nothing (always selected). Pressing Enter on "Confirm selection" advances to the next screen.

#### Scenario: Toggle Developer pack on

- GIVEN `RolePackSelected` is `[false, false]` (indices 0=Developer, 1=PM/Tech Lead; Core is implicit)
- WHEN the user presses Enter on index 1 (Developer Pack option, which is options list index 1)
- THEN `RolePackSelected[0]` becomes `true`
- AND `ProjectRolePacks` is updated to include `"developer"`

#### Scenario: Toggle Developer pack off

- GIVEN `RolePackSelected[0]` is `true`
- WHEN the user presses Enter on the Developer Pack option again
- THEN `RolePackSelected[0]` becomes `false`
- AND `"developer"` is removed from `ProjectRolePacks`

#### Scenario: Core cannot be toggled

- GIVEN the cursor is on the Core option (index 0)
- WHEN the user presses Enter
- THEN nothing changes — Core is always included and the toggle is a no-op

#### Scenario: Confirm with selections

- GIVEN `RolePackSelected` is `[true, true]` (both Developer and PM/Tech Lead selected)
- WHEN the user presses Enter on "Confirm selection"
- THEN `m.ProjectRolePacks` is set to `["core", "developer", "pm-lead"]`
- AND `m.Screen` transitions to `ScreenProjectCI`
- AND `m.Cursor` resets to 0

#### Scenario: Confirm with no optional packs

- GIVEN `RolePackSelected` is `[false, false]`
- WHEN the user presses Enter on "Confirm selection"
- THEN `m.ProjectRolePacks` is set to `["core"]` (Core is always included)
- AND `m.Screen` transitions to `ScreenProjectCI`

### Requirement: Forward Navigation

The system MUST insert `ScreenProjectRolePack` into the project init flow after `ScreenProjectEngram` and before `ScreenProjectCI`. The screen MUST only appear when `ProjectMemory == "obsidian-brain"`.

#### Scenario: Obsidian Brain selected, screen shown

- GIVEN the user selected `ProjectMemory == "obsidian-brain"`
- WHEN the user confirms Engram selection on `ScreenProjectEngram`
- THEN `m.Screen` transitions to `ScreenProjectRolePack` (not directly to `ScreenProjectCI`)
- AND `m.RolePackSelected` is initialized to `make([]bool, 2)` (Developer, PM/Tech Lead)

#### Scenario: Non-obsidian-brain memory, screen skipped

- GIVEN the user selected `ProjectMemory == "vibekanban"` (or any non-obsidian-brain)
- WHEN the user reaches the point where `ScreenProjectRolePack` would appear
- THEN `m.Screen` transitions directly to `ScreenProjectCI`, skipping role packs entirely

#### Scenario: Engram screen to RolePack transition

- GIVEN `m.Screen == ScreenProjectEngram` and `m.ProjectMemory == "obsidian-brain"`
- WHEN the user selects "Yes, add Engram too" or "No, just Obsidian Brain"
- THEN `m.Screen` becomes `ScreenProjectRolePack`

### Requirement: Backward Navigation

The system MUST handle ESC/backspace from `ScreenProjectRolePack` by navigating back to `ScreenProjectEngram`. The system MUST update `goBackInstallStep()` for `ScreenProjectCI` to go back to `ScreenProjectRolePack` when `ProjectMemory == "obsidian-brain"`, otherwise back to `ScreenProjectMemory`.

#### Scenario: ESC from RolePack screen

- GIVEN `m.Screen == ScreenProjectRolePack`
- WHEN the user presses ESC or backspace
- THEN `m.Screen` transitions to `ScreenProjectEngram`
- AND `m.Cursor` resets to 0
- AND `m.RolePackSelected` is reset to `nil` (selections cleared)
- AND `m.ProjectRolePacks` is reset to `nil`

#### Scenario: ESC from CI screen with obsidian-brain

- GIVEN `m.Screen == ScreenProjectCI` and `m.ProjectMemory == "obsidian-brain"`
- WHEN the user presses ESC
- THEN `m.Screen` transitions to `ScreenProjectRolePack` (not `ScreenProjectEngram`)

#### Scenario: ESC from CI screen without obsidian-brain

- GIVEN `m.Screen == ScreenProjectCI` and `m.ProjectMemory == "vibekanban"`
- WHEN the user presses ESC
- THEN `m.Screen` transitions to `ScreenProjectMemory` (existing behavior unchanged)

### Requirement: Confirmation Screen Display

The system MUST display selected role packs in `renderProjectConfirm()` when `ProjectMemory == "obsidian-brain"`.

#### Scenario: Role packs shown in confirmation

- GIVEN `m.ProjectMemory == "obsidian-brain"` and `m.ProjectRolePacks == ["core", "developer"]`
- WHEN `renderProjectConfirm()` is called
- THEN the output includes a line: `"    Packs:   core, developer"`
- AND this line appears after the Engram line and before the CI line

#### Scenario: No role packs line for non-obsidian memory

- GIVEN `m.ProjectMemory == "vibekanban"`
- WHEN `renderProjectConfirm()` is called
- THEN no "Packs:" line appears in the output

### Requirement: View Rendering

The system MUST render `ScreenProjectRolePack` using a dedicated render function (or via `renderSelection()` with checkbox support), following the `renderAIToolSelection()` pattern for multi-select with checkboxes.

#### Scenario: Render with checkboxes

- GIVEN `m.Screen == ScreenProjectRolePack`
- WHEN `View()` is called
- THEN the output shows options with `[x]`/`[ ]` prefixes
- AND Core always shows `[x]`
- AND the title, description, and help text are rendered

### Requirement: HandleKeyPress Registration

The system MUST register `ScreenProjectRolePack` in `handleKeyPress()` so it routes to a handler function (either reusing `handleAIToolsKeys` pattern or a dedicated `handleRolePackKeys` function). The space key MUST toggle selections (not activate leader mode).

#### Scenario: Space toggles selection

- GIVEN `m.Screen == ScreenProjectRolePack` and cursor is on Developer Pack
- WHEN the user presses space
- THEN the Developer Pack selection is toggled
- AND leader mode is NOT activated

#### Scenario: HandleKeyPress routing

- GIVEN `m.Screen == ScreenProjectRolePack`
- WHEN any key is pressed
- THEN the key is routed to the role pack handler (not to `handleSelectionKeys`)

### Requirement: HandleEscape Registration

The system MUST add `ScreenProjectRolePack` to the `handleEscape()` function for proper ESC handling.

#### Scenario: ESC navigates back

- GIVEN `m.Screen == ScreenProjectRolePack`
- WHEN ESC is pressed
- THEN `goBackInstallStep()` is called
- AND the screen transitions to `ScreenProjectEngram`
