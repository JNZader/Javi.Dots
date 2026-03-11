# CLI Flag Specification

## Purpose

Defines the `--project-role-pack` CLI flag for non-interactive project initialization, its validation rules, and integration with the existing `--init-project` flow in `main.go`.

## Requirements

### Requirement: Flag Definition

The system MUST add a `--project-role-pack` flag of type `string` to the `cliFlags` struct in `main.go`. The field name MUST be `projectRolePack`. The flag MUST be registered via `flag.StringVar` in `parseFlags()` with a default value of `""` and description `"Role packs for Obsidian Brain: developer,pm-lead (comma-separated)"`.

#### Scenario: Flag parsed from command line

- GIVEN the user runs `gentleman.dots --non-interactive --init-project --project-path=/foo --project-memory=obsidian-brain --project-role-pack=developer`
- WHEN `parseFlags()` executes
- THEN `flags.projectRolePack` is set to `"developer"`

#### Scenario: Flag not provided

- GIVEN the user runs `gentleman.dots --non-interactive --init-project --project-path=/foo --project-memory=obsidian-brain`
- WHEN `parseFlags()` executes
- THEN `flags.projectRolePack` is `""` (empty string)

### Requirement: Validation — Requires obsidian-brain

The system MUST validate that `--project-role-pack` is only accepted when `--project-memory=obsidian-brain`. If `--project-role-pack` is provided with any other memory module, the system MUST return an error.

#### Scenario: Valid combination

- GIVEN `--project-memory=obsidian-brain` and `--project-role-pack=developer,pm-lead`
- WHEN validation runs in `runNonInteractive()`
- THEN no error is returned

#### Scenario: Invalid memory module

- GIVEN `--project-memory=vibekanban` and `--project-role-pack=developer`
- WHEN validation runs in `runNonInteractive()`
- THEN an error is returned: `"--project-role-pack requires --project-memory=obsidian-brain"`

#### Scenario: Engram memory module

- GIVEN `--project-memory=engram` and `--project-role-pack=developer`
- WHEN validation runs
- THEN an error is returned: `"--project-role-pack requires --project-memory=obsidian-brain"`

### Requirement: Accepted Values

The system MUST accept a comma-separated list of role pack IDs. Valid IDs are: `developer`, `pm-lead`. The system MUST reject unknown pack IDs. Whitespace around values MUST be trimmed. Values MUST be case-insensitive.

#### Scenario: Single valid value

- GIVEN `--project-role-pack=developer`
- WHEN parsed
- THEN `rolePacks` is `["developer"]` (plus implicit `"core"`)

#### Scenario: Multiple valid values

- GIVEN `--project-role-pack=developer,pm-lead`
- WHEN parsed
- THEN `rolePacks` is `["developer", "pm-lead"]` (plus implicit `"core"`)

#### Scenario: Invalid value

- GIVEN `--project-role-pack=designer`
- WHEN validation runs
- THEN an error is returned: `"invalid role pack: designer (valid: developer, pm-lead)"`

#### Scenario: Mixed valid and invalid

- GIVEN `--project-role-pack=developer,designer`
- WHEN validation runs
- THEN an error is returned for `"designer"`

#### Scenario: Whitespace trimming

- GIVEN `--project-role-pack= developer , pm-lead `
- WHEN parsed
- THEN `rolePacks` is `["developer", "pm-lead"]` (whitespace trimmed)

#### Scenario: Case insensitivity

- GIVEN `--project-role-pack=Developer,PM-Lead`
- WHEN parsed
- THEN values are normalized to lowercase: `["developer", "pm-lead"]`

### Requirement: Default Behavior

When `--project-role-pack` is not provided but `--project-memory=obsidian-brain` is set, the system MUST default to Core only (no optional packs). The installer MUST still create the Core vault structure and templates.

#### Scenario: Default when flag omitted

- GIVEN `--project-memory=obsidian-brain` and no `--project-role-pack` flag
- WHEN `runNonInteractive()` processes the project init
- THEN the system uses `rolePacks = ["core"]`
- AND Core templates and folder structure are created

### Requirement: Core Always Included

The system MUST always include `"core"` in the role packs list, regardless of what the user specifies. The user SHOULD NOT need to specify `core` explicitly.

#### Scenario: Core not specified but included

- GIVEN `--project-role-pack=developer`
- WHEN the role packs list is finalized
- THEN it contains `["core", "developer"]`

#### Scenario: Core explicitly specified (idempotent)

- GIVEN `--project-role-pack=core,developer`
- WHEN the role packs list is finalized
- THEN it contains `["core", "developer"]` (no duplicate `"core"`)

### Requirement: Pass-through to Installer

The system MUST pass the role packs to `runProjectInitScript()` (or equivalent installer logic) so that templates are copied based on pack selection. The function signature MUST be extended to accept role packs.

#### Scenario: Packs passed to init script

- GIVEN `rolePacks = ["core", "developer"]`
- WHEN `runProjectInitScript()` is called
- THEN it receives the role packs and creates the appropriate vault structure and copies the appropriate templates

### Requirement: Non-interactive Output

The system MUST display selected role packs in the non-interactive summary output when `--project-role-pack` is provided.

#### Scenario: Summary displays packs

- GIVEN `--project-role-pack=developer,pm-lead`
- WHEN the non-interactive summary is printed
- THEN the output includes `"  Packs:    core, developer, pm-lead"`

### Requirement: Help Text

The system MUST include `--project-role-pack` in the help text printed by `printHelp()`, under the "Project Init Options" section.

#### Scenario: Help text includes flag

- GIVEN the user runs `gentleman.dots --help`
- WHEN help text is displayed
- THEN it includes a line describing `--project-role-pack` with valid values and a usage example
