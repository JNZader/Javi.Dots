# Multi-workspace obsidian.nvim Specification

## Purpose

Defines the changes to `GentlemanNvim/nvim/lua/plugins/obsidian.lua` to support dynamic project workspace detection alongside the existing personal vault. This allows users who initialize a project with Obsidian Brain to edit their project vault notes from within Neovim without manual configuration.

## Requirements

### Requirement: Personal Vault Always Available

The system MUST keep the existing `GentlemanNotes` workspace pointing at `~/.config/obsidian` as the first workspace in the `workspaces` table. This workspace MUST remain functional for all users, including those who never use project vaults.

#### Scenario: Personal vault untouched

- GIVEN a user who has not initialized any project with Obsidian Brain
- WHEN they open Neovim and run `:ObsidianWorkspaces`
- THEN only `GentlemanNotes` appears in the workspace list
- AND it points to `~/.config/obsidian`

#### Scenario: Personal vault alongside project vault

- GIVEN a user who has both `~/.config/obsidian` (personal) and `~/myproject/.obsidian-brain/` (project)
- WHEN they open Neovim in `~/myproject/` and run `:ObsidianWorkspaces`
- THEN both `GentlemanNotes` and the project workspace appear

### Requirement: Dynamic Project Workspace Detection

The system MUST add logic to detect a project vault by checking whether the current working directory (or any parent up to the filesystem root) contains a `.obsidian-brain/` directory. If detected, a second workspace entry MUST be added to the `workspaces` table dynamically.

#### Scenario: Project vault detected

- GIVEN the user opens Neovim in `~/projects/myapp/` and `~/projects/myapp/.obsidian-brain/` exists
- WHEN obsidian.nvim loads
- THEN the workspaces table contains two entries:
  1. `{name = "GentlemanNotes", path = "~/.config/obsidian"}`
  2. `{name = "myapp", path = "~/projects/myapp/.obsidian-brain"}`

#### Scenario: Project vault not detected

- GIVEN the user opens Neovim in `~/random-dir/` and no `.obsidian-brain/` exists in any parent
- WHEN obsidian.nvim loads
- THEN the workspaces table contains only the `GentlemanNotes` entry

#### Scenario: Nested project directory

- GIVEN the user opens Neovim in `~/projects/myapp/src/components/`
- AND `~/projects/myapp/.obsidian-brain/` exists
- WHEN obsidian.nvim loads
- THEN the project workspace is detected (searching up from cwd to find `.obsidian-brain/`)

### Requirement: Workspace Naming

The project workspace name SHOULD be derived from the directory name containing `.obsidian-brain/`. For example, if the vault is at `/home/user/projects/myapp/.obsidian-brain/`, the workspace name SHOULD be `"myapp"`.

#### Scenario: Workspace name from directory

- GIVEN `.obsidian-brain/` is found in `/home/user/cool-project/`
- WHEN the workspace entry is created
- THEN `name` is `"cool-project"`

#### Scenario: Workspace name with special characters

- GIVEN `.obsidian-brain/` is found in `/home/user/my.project-v2/`
- WHEN the workspace entry is created
- THEN `name` is `"my.project-v2"` (preserved as-is)

### Requirement: Templates Subdir

The project workspace MUST set its `templates.subdir` to `"templates"` (matching the folder structure created by the installer). This ensures Obsidian's template insertion commands work within the project vault.

#### Scenario: Templates path resolution

- GIVEN the project workspace path is `~/projects/myapp/.obsidian-brain/`
- WHEN the user runs `:ObsidianTemplate` within the project workspace
- THEN obsidian.nvim looks for templates in `~/projects/myapp/.obsidian-brain/templates/`

### Requirement: No Regression for Users Without Project Vaults

The system MUST NOT error, warn, or behave differently for users who have never used Obsidian Brain in a project. The detection MUST be silent — if no `.obsidian-brain/` is found, the plugin loads with only the personal vault exactly as it does today.

#### Scenario: Clean user experience

- GIVEN a user who never selected Obsidian Brain for any project
- WHEN they open Neovim anywhere
- THEN obsidian.nvim loads exactly as before (single GentlemanNotes workspace)
- AND no errors appear in `:messages`
- AND no warnings appear in the status line

#### Scenario: Directory detection failure graceful

- GIVEN the detection function encounters a permission error on a parent directory
- WHEN obsidian.nvim loads
- THEN the plugin falls back to the personal vault only
- AND no error is raised to the user

### Requirement: Detection Method

The detection SHOULD use `vim.fn.finddir('.obsidian-brain', vim.fn.getcwd() .. ';')` or equivalent upward-search Lua API. The detection MUST happen at plugin load time (in the `opts` function or a `config` callback), not lazily.

#### Scenario: Detection at startup

- GIVEN the user opens Neovim with `nvim .`
- WHEN the plugin loads
- THEN the workspace detection runs before any Obsidian commands are available
- AND the correct workspaces are already configured

### Requirement: Multiple Workspaces Feature

The system MUST rely on the built-in multi-workspace feature of obsidian.nvim (the `workspaces` table in opts). No custom workspace switching logic is needed — obsidian.nvim handles workspace selection based on the current buffer's file path.

#### Scenario: Automatic workspace switching

- GIVEN both GentlemanNotes and a project workspace are configured
- WHEN the user opens a file inside `~/.config/obsidian/`
- THEN obsidian.nvim uses the GentlemanNotes workspace context

- WHEN the user opens a file inside `~/projects/myapp/.obsidian-brain/`
- THEN obsidian.nvim uses the project workspace context
