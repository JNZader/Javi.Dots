# Design: Obsidian Brain Enhancement with Additive Role Packs

## Technical Approach

This change follows the exact patterns already established in the codebase for screen flow, multi-select, installer logic, and CLI flags. Key reference implementations:

- **Screen constant + iota enum**: `model.go:22-91` — new `ScreenProjectRolePack` inserted between `ScreenProjectEngram` (line 79) and `ScreenProjectCI` (line 80)
- **Multi-select pattern**: `ScreenAIToolsSelect` in `update.go:2084-2135` — `AIToolSelected []bool` toggle state + confirm action that collects selected IDs
- **Conditional screen insertion**: `ScreenProjectObsidianInstall` + `ScreenProjectEngram` — only shown when `obsidian-brain` is selected (`update.go:1501-1510`). Role pack follows the same conditional pattern.
- **Installer file operations**: `installer.go:1018-1022` — `system.EnsureDir()` + `system.CopyDir()` for directory creation and config copying
- **CLI flags**: `main.go:17-44` — `cliFlags` struct with `flag.StringVar` registration
- **Multi-workspace obsidian.nvim**: Plugin already supports `workspaces` table (`obsidian.lua:15-19`)

## Architecture Decisions

### Decision 1: Template asset location — `GentlemanNvim/obsidian-brain/`

**Chosen**: `GentlemanNvim/obsidian-brain/` with subdirectories `core/`, `developer/`, `pm-lead/`

**Rationale**: The existing Obsidian integration lives entirely in `GentlemanNvim/` — the plugin config (`nvim/lua/plugins/obsidian.lua`), and the installer creates obsidian directories in `stepInstallNvim()` at `installer.go:1018-1022`. Templates are Obsidian vault content used by the `obsidian.nvim` plugin, so they belong alongside the Neovim-related assets. `GentlemanClaude/` is exclusively for AI agent skills (SKILL.md files), not vault content.

**Rejected alternatives**:
- `GentlemanClaude/obsidian-brain/` — Wrong scope; Claude directory is for AI skills, not Obsidian vault templates
- `installer/assets/` — No precedent for this directory; installer currently copies from `GentlemanNvim/` directly

### Decision 2: Multi-select for role packs — follow `ScreenAIToolsSelect` pattern

**Chosen**: Multi-select with toggle checkboxes, matching `handleAIToolsKeys()` at `update.go:2084-2135`

**Pattern**:
- `RolePackSelected []bool` on `Model` (like `AIToolSelected []bool` at `model.go:203`)
- Options list: `"Core (always included)"`, `"Developer Pack"`, `"PM/Tech Lead Pack"`, separator, `"Confirm selection"`
- Core is always checked and cannot be unchecked (hardcoded `true` at index 0)
- Enter on toggleable items flips `RolePackSelected[cursor]`
- Enter on "Confirm" collects selected pack IDs into `UserChoices.ProjectRolePacks`

**Rejected**: Single-select (would limit users to one pack — additive is the whole point)

### Decision 3: Template file format — standard YAML frontmatter

**Chosen**: Standard Obsidian-compatible YAML frontmatter (`---` delimited) with section headers. No Templater-specific syntax (e.g., `<% tp.* %>`) to avoid requiring a community plugin.

**Format**:
```markdown
---
type: braindump
created: "{{date}}"
tags: []
---

# Braindump

## Capture
...
```

The `{{date}}` syntax is supported by Obsidian's built-in Templates core plugin (Settings > Core Plugins > Templates). No external plugins required.

**Rationale**: The existing `obsidian.lua` config at line 50-55 already sets `templates.subdir = "templates"` and `templates.date_format = "%Y-%m-%d-%a"`. Standard YAML frontmatter works with both the core Templates plugin and Templater if users install it later.

### Decision 4: Multi-workspace detection — static + conditional config

**Chosen**: Static config with a conditional workspace registration. The `obsidian.lua` plugin config will check for an environment variable or a directory marker to add project workspaces.

**Implementation**:
```lua
workspaces = {
  {
    name = "GentlemanNotes",
    path = os.getenv("HOME") .. "/.config/obsidian",
  },
},
-- Dynamic: detect project vault if cwd contains .obsidian-brain/
detect_cwd = true,
```

The `obsidian.nvim` plugin (v3+) natively supports `detect_cwd = true` which automatically registers any directory containing a `.obsidian/` folder as a workspace. We use `.obsidian-brain/` as the vault root and create a `.obsidian/` marker inside it so the plugin detects it.

**Alternative considered**: Dynamic Lua function scanning — adds complexity and fragility. The plugin's built-in detection is simpler and officially supported.

### Decision 5: Screen constant placement — between `ScreenProjectEngram` and `ScreenProjectCI`

**Exact insertion point**: After `ScreenProjectEngram` (currently line 79) and before `ScreenProjectCI` (currently line 80).

```go
ScreenProjectEngram          // Yes/No: add Engram alongside Obsidian Brain   (line 79)
ScreenProjectRolePack        // NEW: Multi-select role packs                  (line 80)
ScreenProjectCI              // Single-select: CI provider                    (line 81)
```

**Iota safety**: All Screen constants use `iota` (line 22-91). Inserting a new constant shifts all subsequent values by +1. This is safe because:
1. No hardcoded numeric Screen values exist (verified via grep — all references use the constant names)
2. The `Screen` type is only used in switch statements that match on constant names
3. Tests reference constants, not numeric values (verified in `project_screens_test.go`)

## Data Flow

```
┌─────────────────────────────────────────────────────────────┐
│                        TUI FLOW                              │
│                                                              │
│  ScreenProjectMemory ──"obsidian-brain"──→ ScreenProjectObsidianInstall │
│          │                                        │          │
│          │                                        ▼          │
│          │                              ScreenProjectEngram  │
│          │                                        │          │
│          │                                        ▼          │
│          │                          ScreenProjectRolePack    │
│          │                           (multi-select, NEW)     │
│          │                                        │          │
│          │                                        ▼          │
│          └──"other memory"──→ ScreenProjectCI ◀───┘         │
│                                        │                     │
│                                        ▼                     │
│                              ScreenProjectConfirm            │
│                              (shows role packs)              │
│                                        │                     │
│                                        ▼                     │
│                           ScreenProjectInstalling            │
└─────────────────────────────────────────────────────────────┘
                                        │
                                        ▼
┌─────────────────────────────────────────────────────────────┐
│                    INSTALLER LOGIC                            │
│                                                              │
│  runProjectInit()                                            │
│    │                                                         │
│    ├─ runProjectInitScript(path, memory, ci, engram)         │
│    │   └─ passes --role-pack=core,developer to init-project.sh│
│    │                                                         │
│    └─ copyRolePackTemplates()  ← NEW function                │
│        ├─ src: {repoDir}/GentlemanNvim/obsidian-brain/{pack}/│
│        └─ dst: {projectPath}/.obsidian-brain/templates/      │
│            ├─ creates folder structure (inbox/, resources/,   │
│            │   knowledge/, templates/)                        │
│            └─ creates .obsidian/ marker for plugin detection  │
└─────────────────────────────────────────────────────────────┘
```

## File Changes

| File | Action | Changes |
|------|--------|---------|
| `installer/internal/tui/model.go` | Modify | Add `ScreenProjectRolePack` constant (line 80). Add `ProjectRolePacks []string` to `UserChoices` (after `ProjectEngram`). Add `RolePackSelected []bool` to `Model`. Add case to `GetCurrentOptions()`, `GetScreenTitle()`, `GetScreenDescription()`. |
| `installer/internal/tui/update.go` | Modify | Add `ScreenProjectRolePack` to `handleKeyPress()` switch in `ScreenAIToolsSelect`-style handler list. Add `handleRolePackKeys()` function (multi-select toggle). Modify `handleSelection()` for `ScreenProjectEngram` → next screen logic. Modify `goBackInstallStep()` for `ScreenProjectRolePack` and `ScreenProjectCI` back-nav. Modify `runProjectInit()` to pass role packs. Add `ScreenProjectRolePack` to `handleEscape()`. |
| `installer/internal/tui/view.go` | Modify | Add `ScreenProjectRolePack` case to `View()` switch — calls `renderRolePackSelection()`. Add `renderRolePackSelection()` function (checkbox rendering like `renderAIToolSelection()`). Modify `renderProjectConfirm()` to display selected role packs. |
| `installer/internal/tui/installer.go` | Modify | Extend obsidian directory creation section (~line 1018-1022) to copy role pack templates. Add `copyRolePackTemplates()` function. |
| `installer/cmd/gentleman-installer/main.go` | Modify | Add `projectRolePack string` to `cliFlags` struct. Add `flag.StringVar` registration. Add validation in `runNonInteractive()`. |
| `GentlemanNvim/nvim/lua/plugins/obsidian.lua` | Modify | Add `detect_cwd = true` to opts. Keep existing `GentlemanNotes` workspace untouched. |
| `AGENTS.md` | Modify | Add 3 new skills to the "Generic Skills" table. |
| `installer/internal/tui/project_screens_test.go` | Modify | Add tests for `ScreenProjectRolePack` screen (options, title, description, navigation, selection). |
| `GentlemanNvim/obsidian-brain/core/templates/braindump.md` | Create | Quick capture template |
| `GentlemanNvim/obsidian-brain/core/templates/resource-capture.md` | Create | Link + summary + tags template |
| `GentlemanNvim/obsidian-brain/core/templates/consolidation.md` | Create | Weekly knowledge synthesis template |
| `GentlemanNvim/obsidian-brain/core/templates/daily-note.md` | Create | Daily note template |
| `GentlemanNvim/obsidian-brain/developer/templates/adr.md` | Create | Architecture Decision Record template |
| `GentlemanNvim/obsidian-brain/developer/templates/coding-session.md` | Create | Coding session log template |
| `GentlemanNvim/obsidian-brain/developer/templates/tech-debt.md` | Create | Tech debt tracker template |
| `GentlemanNvim/obsidian-brain/developer/templates/debug-journal.md` | Create | Debug journal template |
| `GentlemanNvim/obsidian-brain/developer/templates/sdd-feedback.md` | Create | SDD feedback loop template |
| `GentlemanNvim/obsidian-brain/pm-lead/templates/meeting-notes.md` | Create | Meeting notes template |
| `GentlemanNvim/obsidian-brain/pm-lead/templates/sprint-review.md` | Create | Sprint review template |
| `GentlemanNvim/obsidian-brain/pm-lead/templates/stakeholder-update.md` | Create | Stakeholder update template |
| `GentlemanNvim/obsidian-brain/pm-lead/templates/risk-registry.md` | Create | Risk registry template |
| `GentlemanNvim/obsidian-brain/pm-lead/templates/daily-brief.md` | Create | Daily brief template |
| `GentlemanNvim/obsidian-brain/pm-lead/templates/weekly-brief.md` | Create | Weekly brief template |
| `GentlemanNvim/obsidian-brain/pm-lead/templates/team-intelligence.md` | Create | Team intelligence template |
| `GentlemanClaude/skills/obsidian-braindump/SKILL.md` | Create | Braindump capture workflow skill |
| `GentlemanClaude/skills/obsidian-consolidation/SKILL.md` | Create | Weekly knowledge consolidation skill |
| `GentlemanClaude/skills/obsidian-resource-capture/SKILL.md` | Create | Resource capture and annotation skill |

## Interfaces / Contracts

### New UserChoices field

```go
// UserChoices stores all user selections (model.go)
type UserChoices struct {
	// ... existing fields ...
	// Project init
	InitProject   bool
	ProjectPath   string
	ProjectStack  string
	ProjectMemory string
	ProjectCI        string
	ProjectEngram    bool
	ProjectRolePacks []string // NEW: selected role packs, e.g. ["core", "developer"]
	InstallObsidian  bool
}
```

### New Model fields

```go
// Model is the main application state (model.go)
type Model struct {
	// ... existing fields ...
	// Project init
	ProjectPathInput string
	ProjectPathError string
	ProjectStack     string
	ProjectMemory    string
	ProjectEngram    bool
	ProjectCI        string
	ProjectLogLines  []string
	// Role pack multi-select (NEW)
	RolePackSelected []bool // Toggle state for each role pack in ScreenProjectRolePack
	// ... rest of fields ...
}
```

### New Screen constant

```go
const (
	// ... existing screens ...
	// Project Init screens
	ScreenProjectPath       // Text input: project directory
	ScreenProjectStack      // Single-select: detected stack confirmation/override
	ScreenProjectMemory          // Single-select: memory module
	ScreenProjectObsidianInstall // Offer to install Obsidian app if not detected
	ScreenProjectEngram          // Yes/No: add Engram alongside Obsidian Brain
	ScreenProjectRolePack        // NEW: Multi-select role packs (Core always on)
	ScreenProjectCI              // Single-select: CI provider
	ScreenProjectConfirm         // Summary before execution
	ScreenProjectInstalling      // Progress log
	ScreenProjectResult          // Success/error
	// ... rest of screens ...
)
```

### Role pack ID map (in update.go)

```go
// rolePackIDMap maps role pack option index to pack ID
var rolePackIDMap = []string{"core", "developer", "pm-lead"}
```

### GetCurrentOptions for ScreenProjectRolePack

```go
case ScreenProjectRolePack:
	return []string{
		"Core (always included)",
		"Developer Pack",
		"PM/Tech Lead Pack",
		"─────────────",
		"✅ Confirm selection",
	}
```

### GetScreenTitle / GetScreenDescription

```go
case ScreenProjectRolePack:
	return "📦 Initialize Project — Role Packs"

case ScreenProjectRolePack:
	return "Select role packs. Core is always included."
```

### New CLI flag

```go
// In cliFlags struct (main.go)
type cliFlags struct {
	// ... existing fields ...
	projectRolePack string // comma-separated: "developer,pm-lead"
}

// Registration
flag.StringVar(&flags.projectRolePack, "project-role-pack", "",
	"Role packs: developer,pm-lead (comma-separated, core is always included)")
```

### Validation in runNonInteractive

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
		if !validPacks[pack] {
			return fmt.Errorf("invalid role pack: %s (valid: developer, pm-lead)", pack)
		}
		rolePacks = append(rolePacks, pack)
	}
}
// Core is always included
rolePacks = append([]string{"core"}, rolePacks...)
```

### obsidian.lua workspace config

```lua
return {
  "obsidian-nvim/obsidian.nvim",
  version = "*",
  lazy = false,
  enabled = function()
    return not vim.g.disable_obsidian
  end,
  dependencies = {
    "nvim-lua/plenary.nvim",
  },
  opts = {
    legacy_commands = false,
    workspaces = {
      {
        name = "GentlemanNotes",
        path = os.getenv("HOME") .. "/.config/obsidian",
      },
    },
    -- Detect project vaults automatically when cwd contains .obsidian/
    detect_cwd = true,
    completion = {
      cmp = true,
    },
    picker = {
      name = "snacks.pick",
    },
    callbacks = {
      -- ... existing callbacks unchanged ...
    },
    templates = {
      subdir = "templates",
      date_format = "%Y-%m-%d-%a",
      gtime_format = "%H:%M",
      tags = "",
    },
  },
}
```

### copyRolePackTemplates function (installer.go)

```go
// copyRolePackTemplates copies selected role pack templates into the project vault
func copyRolePackTemplates(repoDir, projectPath string, rolePacks []string) error {
	vaultDir := filepath.Join(projectPath, ".obsidian-brain")
	templatesDir := filepath.Join(vaultDir, "templates")

	// Create vault folder structure
	for _, dir := range []string{
		vaultDir,
		filepath.Join(vaultDir, "inbox"),
		filepath.Join(vaultDir, "resources"),
		filepath.Join(vaultDir, "knowledge"),
		templatesDir,
		filepath.Join(vaultDir, ".obsidian"), // marker for plugin detection
	} {
		if err := system.EnsureDir(dir); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	// Copy templates from each selected role pack
	for _, pack := range rolePacks {
		srcDir := filepath.Join(repoDir, "GentlemanNvim", "obsidian-brain", pack, "templates")
		if _, err := os.Stat(srcDir); os.IsNotExist(err) {
			continue // pack doesn't exist yet, skip gracefully
		}
		entries, err := os.ReadDir(srcDir)
		if err != nil {
			return fmt.Errorf("failed to read pack %s templates: %w", pack, err)
		}
		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}
			src := filepath.Join(srcDir, entry.Name())
			dst := filepath.Join(templatesDir, entry.Name())
			data, err := os.ReadFile(src)
			if err != nil {
				return fmt.Errorf("failed to read template %s: %w", src, err)
			}
			if err := os.WriteFile(dst, data, 0644); err != nil {
				return fmt.Errorf("failed to write template %s: %w", dst, err)
			}
		}
	}

	return nil
}
```

### Screen flow modifications in handleSelection (update.go)

```go
// After ScreenProjectEngram selection:
case ScreenProjectEngram:
	m.ProjectEngram = m.Cursor == 0
	// NEW: show role pack selection
	m.Screen = ScreenProjectRolePack
	m.Cursor = 0
	// Initialize role pack selection: Core always on
	m.RolePackSelected = []bool{true, false, false}
```

### Back navigation in goBackInstallStep (update.go)

```go
case ScreenProjectRolePack:
	m.Screen = ScreenProjectEngram
	m.Cursor = 0

case ScreenProjectCI:
	if m.ProjectMemory == "obsidian-brain" {
		m.Screen = ScreenProjectRolePack  // Changed from ScreenProjectEngram
	} else {
		m.Screen = ScreenProjectMemory
	}
	m.Cursor = 0
```

### AI skill frontmatter format (matching react-19/SKILL.md)

```yaml
---
name: obsidian-braindump
description: >
  Braindump capture workflow for Obsidian Brain vaults.
  Trigger: When capturing quick ideas, notes, or unstructured thoughts.
license: Apache-2.0
metadata:
  author: gentleman-programming
  version: "1.0"
---
```

## Testing Strategy

### Existing test patterns (from `project_screens_test.go`)

The project screens already have comprehensive tests (1323 lines) covering:
- Screen navigation (forward and back)
- Option counts per screen
- Title/description non-empty checks
- Conditional screen flow (obsidian-brain → ObsidianInstall → Engram)
- Path input validation, cursor movement, tab completion, file browser

### New tests to add (in `project_screens_test.go`)

| Test | Pattern | Validates |
|------|---------|-----------|
| `TestRolePackScreenOptions` | `TestGetCurrentOptionsProjectScreens` | 5 options (3 packs + separator + confirm) |
| `TestRolePackScreenTitle` | `TestGetScreenTitleProjectScreens` | Non-empty title |
| `TestRolePackCoreAlwaysOn` | Custom | Core cannot be deselected |
| `TestRolePackToggle` | `TestAIToolsSelectToggle` (if exists) | Developer/PM toggling works |
| `TestRolePackConfirm` | Custom | Confirm collects selected packs into `Choices.ProjectRolePacks` |
| `TestRolePackBackNav` | `TestProjectEscapeBackNavigation` | Backspace → ScreenProjectEngram |
| `TestCIBackNavWithRolePack` | Extends existing test | CI backspace → ScreenProjectRolePack (not Engram) |
| `TestRolePackOnlyForObsidian` | Custom | Role pack screen skipped when memory != obsidian-brain |
| `TestProjectConfirmShowsRolePacks` | Custom | renderProjectConfirm() includes pack names |

### E2E test coverage

The E2E test flag should be: `--project-memory=obsidian-brain --project-role-pack=developer --project-engram`

## Migration / Rollout

No migration needed. All changes are purely additive:
- New iota constant doesn't break existing constants (no hardcoded values)
- New `UserChoices` and `Model` fields have zero values (empty slice/nil) — existing code unaffected
- Role pack screen only appears in the obsidian-brain flow path
- `obsidian.lua` change (`detect_cwd = true`) is backward compatible — existing personal vault continues working
- Templates in `GentlemanNvim/obsidian-brain/` are new files with no dependencies
- CLI flag is optional and defaults to empty string

## Open Questions

1. **`init-project.sh` flag support**: The `runProjectInitScript()` at `installer.go:1729` currently passes `--memory`, `--ci`, `--engram` to the external `project-starter-framework` script. Adding `--role-pack` requires the external script to understand this flag. Two options:
   - (a) Handle template copying entirely in Go (in `installer.go`), independent of `init-project.sh` — this is the safer approach and what the design above assumes
   - (b) Pass `--role-pack` to `init-project.sh` and let it handle copying — requires coordination with the external repo

   **Recommendation**: Option (a) — handle in Go. The template source files live in `GentlemanNvim/obsidian-brain/` within this repo, so copying them is a local operation.

2. **`detect_cwd` plugin behavior**: Need to verify that `obsidian.nvim` v3+ supports `detect_cwd = true` as documented. If not available, the fallback is to check for `vim.fn.isdirectory(vim.fn.getcwd() .. '/.obsidian-brain/.obsidian')` in a custom workspace override function. This should be verified during implementation.

3. **Template naming conflicts**: If Core and Developer packs both have a template with the same filename, the later copy overwrites the earlier. The design ensures unique filenames across packs (braindump.md, adr.md, etc. are all distinct), but this should be enforced as a convention.
