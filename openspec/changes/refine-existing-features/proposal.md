# Proposal: Refine Existing Features

## Intent

The SDD orchestrator, multi-perspective explore, installer hooks, consolidation skill, and AGENTS.md infrastructure have grown organically over several changes. Eight concrete refinements have been identified that improve configurability, iteration depth, output quality, developer experience, and maintainability — without introducing new features from scratch. This change bundles them into one coordinated pass because many touch the same files (CLAUDE.md, AGENTS.md, setup.sh).

## Scope

### In Scope

1. **Configurable explore perspectives** — Document that `explore.perspectives` in `openspec/config.yaml` accepts custom perspective names (e.g., "performance", "cost", "compliance"), not just the 4 defaults.
2. **Multi-round iteration for explore** — Add `explore.rounds` config (default 1, max 3) so synthesis can trigger deeper passes when conflicts remain unresolved.
3. **Disagreement table in synthesis** — Replace free-form agree/conflict text with a structured Agreement Matrix table in the synthesis sub-agent prompt template.
4. **Hook scripts for Claude Code** — Create `comment-check.sh` and `todo-tracker.sh` in `GentlemanClaude/hooks/`, add sync step to `skills/setup.sh`.
5. **Path-scoped rules** — Add a `### Path-Scoped Rules` documentation section to `GentlemanClaude/CLAUDE.md` mapping file globs to convention sets.
6. **Wikilinks in mind map nodes** — Update `obsidian-consolidation/SKILL.md` Mind Map instructions to use `[[wikilinked entities]]` as node text.
7. **Skill version table** — Add a `### Skill Versions` table to `AGENTS.md` auto-populated from SKILL.md frontmatter `version:` fields.
8. **Orchestrator self-check checklist** — Replace the plain-text delegation self-check with a concrete inline checklist in both CLAUDE.md and AGENTS.md.

### Out of Scope

- New SDD phases or commands
- Changes to the Go TUI codebase
- Changes to Neovim/LazyVim configuration
- Engram backend changes
- New skills (only existing skill/file modifications)
- Automated version bumping or CI validation of the version table

## Approach

All 8 improvements are documentation/config/script changes — no Go code is modified. The work is organized into 4 phases to minimize merge conflicts:

1. **Phase 1 (Multi-perspective explore)** — Improvements #1, #2, #3 all modify the Multi-Perspective Explore section in CLAUDE.md and AGENTS.md. Do them together.
2. **Phase 2 (Installer + path rules)** — Improvements #4, #5 add new files and a new section. Hook scripts are new files in `GentlemanClaude/hooks/`, setup.sh gets a sync function, CLAUDE.md gets a path-scoped rules section.
3. **Phase 3 (Obsidian)** — Improvement #6 is isolated to one skill file.
4. **Phase 4 (Infrastructure)** — Improvements #7, #8 modify AGENTS.md and CLAUDE.md in non-overlapping sections.

## Affected Areas

| Area | Impact | Description |
|------|--------|-------------|
| `GentlemanClaude/CLAUDE.md` | Modified | Multi-perspective explore section (#1,#2,#3), path-scoped rules (#5), self-check checklist (#8) |
| `AGENTS.md` | Modified | Multi-perspective explore section (#1,#2,#3), skill version table (#7), self-check checklist (#8) |
| `GentlemanClaude/skills/obsidian-consolidation/SKILL.md` | Modified | Mind Map generation instructions (#6) |
| `GentlemanClaude/hooks/comment-check.sh` | New | PostToolUse hook script (#4) |
| `GentlemanClaude/hooks/todo-tracker.sh` | New | Stop hook script (#4) |
| `skills/setup.sh` | Modified | New `sync_hooks()` function and menu option (#4) |

## Risks

| Risk | Likelihood | Mitigation |
|------|------------|------------|
| Multi-round explore adds latency/cost | Medium | Default rounds=1 preserves current behavior; max cap at 3 |
| Hook scripts may conflict with user's existing hooks | Low | Install to dedicated `GentlemanClaude/hooks/` dir; sync copies but doesn't overwrite existing user hooks |
| Agreement Matrix table is harder to produce than free text | Low | Provide explicit template with example; synthesis agent has structured output |
| Version table becomes stale if not maintained | Medium | Document it as a manual step in setup.sh; suggest running setup.sh after version bumps |
| Path-scoped rules section may grow unwieldy | Low | Start with 4 clear rules; document pattern for adding more |

## Rollback Plan

All changes are to markdown files and shell scripts. Rollback is a simple `git revert` of the commit(s). No database migrations, no binary changes, no infrastructure state to unwind.

## Dependencies

- None. All files being modified already exist in the repository.

## Success Criteria

- [ ] `openspec/config.yaml` schema documented with `explore.perspectives` accepting custom names and `explore.rounds` accepting 1-3
- [ ] Multi-round iteration logic is clearly specified in the Multi-Perspective Explore section of both CLAUDE.md and AGENTS.md
- [ ] Synthesis sub-agent prompt template includes the Agreement Matrix table format
- [ ] `GentlemanClaude/hooks/comment-check.sh` exists and is <15 lines
- [ ] `GentlemanClaude/hooks/todo-tracker.sh` exists and is <15 lines
- [ ] `skills/setup.sh` includes a `sync_hooks()` function
- [ ] `GentlemanClaude/CLAUDE.md` contains `### Path-Scoped Rules` section with 4 glob→convention mappings
- [ ] `obsidian-consolidation/SKILL.md` Mind Map instructions specify `[[wikilinks]]` in node text
- [ ] `AGENTS.md` contains `### Skill Versions` table with version data from all 40 skills
- [ ] Both CLAUDE.md and AGENTS.md Delegation Rules include the concrete self-check checklist
