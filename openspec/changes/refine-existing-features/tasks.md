# Tasks: Refine Existing Features

## Phase 1: Multi-Perspective Explore Improvements (#1, #2, #3)

- [x] 1.1 Update `GentlemanClaude/CLAUDE.md` "Config Reading" subsection under Multi-Perspective Explore: add documentation clarifying `explore.perspectives` accepts ANY string names, not just the 4 defaults. Add examples: `performance`, `cost`, `compliance`, `security`.
- [x] 1.2 Update `GentlemanClaude/CLAUDE.md` "Config Reading" subsection: add `explore.rounds` to the config schema example (default: 1, max: 3). Document the round lifecycle: fan-out → synthesis → check for "NEEDS FURTHER ANALYSIS" → repeat if rounds remain.
- [x] 1.3 Update `GentlemanClaude/CLAUDE.md` "Fan-Out Dispatch" subsection: add instructions for round 2+ behavior — each perspective agent receives the prior round's synthesis and is instructed to dig deeper on unresolved conflicts and blind spots.
- [x] 1.4 Update `GentlemanClaude/CLAUDE.md` "Synthesis Dispatch" subsection: replace the free-form agree/conflict bullet points with the Agreement Matrix table format. Include the complete table template with columns: Finding, one column per perspective, Confidence. Add marking legend: checkmark = agrees, X = disagrees, dash = not covered.
- [x] 1.5 Update `GentlemanClaude/CLAUDE.md` "Synthesis Dispatch": add instruction that synthesis agent MUST flag items as "NEEDS FURTHER ANALYSIS" if confidence is Low and rounds remain available.
- [x] 1.6 Mirror ALL changes from tasks 1.1-1.5 in `AGENTS.md` — the SDD section in AGENTS.md has the same Multi-Perspective Explore content and MUST stay in sync.

## Phase 2: Installer & Path-Scoped Rules (#4, #5)

- [x] 2.1 Create `GentlemanClaude/hooks/comment-check.sh`: PostToolUse hook script (5-15 lines) that checks recent tool output for TODO/FIXME/HACK comments without explanatory context and warns. Include shebang, comment header describing purpose/trigger, and the check logic.
- [x] 2.2 Create `GentlemanClaude/hooks/todo-tracker.sh`: Stop hook script (5-15 lines) that outputs a reminder to review and update the todo list before session ends. Include shebang and comment header.
- [x] 2.3 Update `skills/setup.sh`: add a `sync_hooks()` function that copies `GentlemanClaude/hooks/*.sh` to `~/.claude/hooks/`. Only copy if destination file does NOT already exist (no-clobber). Make copied files executable.
- [x] 2.4 Update `skills/setup.sh`: add menu option 9 ("Sync hooks to ~/.claude/hooks/") calling `sync_hooks()`. Update `--sync-claude` and `--sync-all` to also call `sync_hooks()`. Update `show_help()` with new option.
- [x] 2.5 Add `### Path-Scoped Rules` section to `GentlemanClaude/CLAUDE.md` (after the "How to use skills" subsection, before "---" separator preceding SDD). Include 4 mappings: `**/*.go` → Go conventions, `**/*.lua` → Neovim/lazy.nvim conventions, `**/SKILL.md` → skill-creator skill, `**/*.md` in `obsidian-brain/` → Obsidian conventions.

## Phase 3: Obsidian Consolidation Improvement (#6)

- [x] 3.1 Update `GentlemanClaude/skills/obsidian-consolidation/SKILL.md` step 7 (Generate Mind Map): change the instruction to specify that markmap leaf nodes MUST use `[[wikilinked entity]]` syntax instead of plain text. Example: `- [[Go]] [[error-wrapping]]` instead of `- Go error handling`.
- [x] 3.2 Update the Mind Map template example in the "Core Only" section of `obsidian-consolidation/SKILL.md`: replace plain text node labels with `[[wikilink]]` syntax in the markmap code block example.
- [x] 3.3 Update rule 10 at the end of `obsidian-consolidation/SKILL.md`: add explicit statement that mind map leaf nodes MUST use `[[wikilinks]]` matching the entities already extracted in the `## Entities` section.

## Phase 4: Infrastructure Improvements (#7, #8)

- [x] 4.1 Add `### Skill Versions` table to `AGENTS.md` after the "Generic Skills" table. Populate with all 40 skills from `GentlemanClaude/skills/*/SKILL.md`, reading the `version:` frontmatter field. Columns: Skill, Version, Last Updated (use `2026-03-11` as current date for all). Sort alphabetically.
- [x] 4.2 Update `GentlemanClaude/CLAUDE.md` Delegation Rules section: replace "3. **Self-check before every response:** ..." with a concrete markdown checklist:
  ```
  Before responding, verify:
  - [ ] Am I about to read source code? → DELEGATE
  - [ ] Am I about to write/edit code? → DELEGATE
  - [ ] Am I about to analyze architecture? → DELEGATE
  - [ ] Am I about to run tests/builds? → DELEGATE
  - [ ] Am I about to write specs/proposals/design? → DELEGATE
  If ALL unchecked → safe to respond inline.
  ```
- [x] 4.3 Mirror the same self-check checklist change from task 4.2 in `AGENTS.md` Delegation Rules section — MUST stay in sync with CLAUDE.md.
- [ ] 4.4 Verify: confirm all 8 improvements are reflected across CLAUDE.md, AGENTS.md, setup.sh, and skill files. Check that no section references stale content from before the changes.
