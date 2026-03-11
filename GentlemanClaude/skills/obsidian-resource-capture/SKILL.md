---
name: obsidian-resource-capture
description: >
  Capture URLs and resources with auto-extracted insights into Obsidian vault.
  Trigger: When user says capture resource, save link, bookmark, or shares a URL to save with notes.
license: Apache-2.0
metadata:
  author: gentleman-programming
  version: "1.0"
---

## Purpose

Capture external resources (URLs, articles, docs, videos, tools) into the Obsidian Brain vault with structured summaries and key takeaways. This skill transforms a raw URL into a searchable, categorized knowledge artifact using the `resource-capture.md` template.

## When to Invoke

Trigger this skill when the user says any of:
- "capture this resource", "capture this link"
- "save this link", "save this article"
- "bookmark this", "add to resources"
- "I found this useful", "check this out" (with a URL)
- "resource capture", "log this reference"
- "save for later" (with a URL or reference)

## Process Flow

1. **Extract the URL/reference** — Identify the URL or resource reference from the user's message. If no URL is provided, ask for one.
2. **Fetch or infer content** — If WebFetch is available, retrieve the page content. If not, use the URL, title, and any context the user provides to generate the summary.
3. **Generate summary** — Write a 2-4 sentence summary of what the resource covers and why it's valuable.
4. **Extract key takeaways** — Identify 3-5 actionable or notable points from the resource.
5. **Categorize with tags** — Auto-suggest tags based on the content topic, technology, and resource type.
6. **Detect vault location** — Look for `.obsidian-brain/` in the project or `~/.config/obsidian/` as fallback.
7. **Save or output** — Write to `resources/` directory if vault is accessible, otherwise output the full markdown.

## Role-Aware Behavior

The resource capture adapts its categorization and extraction based on active role packs:

### Core Only (no role packs detected)
Use the base template exactly:
```markdown
---
title: "<resource title>"
date: "{{date}} {{time}}"
tags:
  - resource
  - <topic-tag>
---

## Source

[<resource title>](<url>)

## Summary

<2-4 sentence summary of what this resource covers and why it matters>

## Key Takeaways

- <takeaway 1>
- <takeaway 2>
- <takeaway 3>
```

### Developer Pack Active
When developer templates are detected, add a technical analysis lens:

- Add a `## Technical Relevance` section:
  ```markdown
  ## Technical Relevance

  - **Applies to**: <which parts of the codebase or architecture this is relevant to>
  - **Technology**: <language, framework, tool mentioned>
  - **Pattern/Concept**: <design pattern, architecture principle, or technique>
  - **Try it**: <specific action to apply this knowledge — e.g., "Test this approach in the auth module">
  ```
- Add a `## Code Snippets` section if the resource contains code:
  ```markdown
  ## Code Snippets

  ```<language>
  // Key code example from the resource
  ```
  ```
- Suggest tech-specific tags: `#golang`, `#react`, `#architecture`, `#performance`, `#security`, `#api`, `#testing`
- Add related notes links to existing technical notes: `Related: [[adr-jwt-auth]], [[tech-debt-database]]`

### PM/Tech Lead Pack Active
When PM templates are detected, add a strategic analysis lens:

- Add a `## Strategic Relevance` section:
  ```markdown
  ## Strategic Relevance

  - **Applies to**: <which project, initiative, or team concern>
  - **Category**: <competitor analysis | market research | best practice | tool evaluation | process improvement>
  - **Share with**: <team members or stakeholders who should see this>
  - **Action**: <specific next step — e.g., "Discuss in next sprint planning">
  ```
- Add a `## Discussion Points` section:
  ```markdown
  ## Discussion Points

  - <question this resource raises for the team>
  - <how this could affect current priorities>
  - <comparison to current approach>
  ```
- Suggest management-specific tags: `#competitor`, `#market-research`, `#process`, `#tool-evaluation`, `#team-practice`

### Both Packs Active
Merge all sections. A resource about a competitor's technical architecture, for example, would get both `## Technical Relevance` (architecture patterns to learn from) and `## Strategic Relevance` (competitive positioning).

## Template Reference

Uses: `resource-capture.md` from `GentlemanNvim/obsidian-brain/core/templates/resource-capture.md`

Template structure:
```yaml
---
title:
date: "{{date}} {{time}}"
tags:
  - resource
---
```
Sections: `## Source`, `## Summary`, `## Key Takeaways`

## Output Format

The generated file should be saved as:
```
.obsidian-brain/resources/YYYY-MM-DD-<slugified-title>.md
```

Example filename: `2025-01-15-go-error-handling-best-practices.md`

### Complete Example (Developer Pack Active)

```markdown
---
title: "Go Error Handling Best Practices"
date: "2025-01-15 10:45"
tags:
  - resource
  - golang
  - error-handling
  - best-practices
---

## Source

[Go Error Handling Best Practices](https://go.dev/blog/error-handling)

## Summary

Comprehensive guide on idiomatic error handling in Go, covering custom error types, error wrapping with `%w`, sentinel errors, and the `errors.Is`/`errors.As` patterns introduced in Go 1.13. Particularly relevant for teams migrating from exception-based languages.

## Key Takeaways

- Always wrap errors with context using `fmt.Errorf("operation failed: %w", err)` to preserve the error chain
- Use sentinel errors (`var ErrNotFound = errors.New(...)`) for expected error conditions that callers need to check
- Prefer `errors.Is()` over `==` comparison to handle wrapped errors correctly
- Custom error types should implement `Error() string` and optionally `Unwrap() error`
- Avoid `panic()` for expected error conditions — reserve it for truly unrecoverable states

## Technical Relevance

- **Applies to**: Error handling across all Go services, especially `services/auth/` and `services/api-gateway/`
- **Technology**: Go 1.21+
- **Pattern**: Error wrapping chain with sentinel errors at package boundaries
- **Try it**: Refactor the auth service error handling to use wrapped errors instead of string matching

## Code Snippets

```go
// Custom error type with context
type ValidationError struct {
    Field   string
    Message string
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("validation failed on %s: %s", e.Field, e.Message)
}

// Wrapping errors
if err := db.Query(ctx, query); err != nil {
    return fmt.Errorf("user lookup failed: %w", err)
}
```

## Related Notes

- [[tech-debt-error-handling]]
- [[coding-session-auth-refactor]]
```

## Critical Rules

1. **URL is required** — If no URL or resource reference is provided, ask for one before generating the note.
2. **Always format source as markdown link** — The `## Source` section MUST use `[title](url)` format, never bare URLs.
3. **Summary is 2-4 sentences** — Not a paragraph, not a single sentence. Concise but informative.
4. **Key takeaways are actionable** — Each takeaway should be something the reader can DO or APPLY, not just a fact restatement.
5. **No external dependencies for generation** — If WebFetch is unavailable, generate based on URL structure, user-provided context, and your knowledge. Do NOT refuse to generate because you cannot fetch the URL.
6. **Suggest tags based on content** — Always include `resource` tag plus 2-4 topic-specific tags. Do not use generic tags like `#interesting` or `#useful`.
7. **One resource per note** — Each capture is for a single URL/resource. If the user shares multiple URLs, create separate notes.
8. **Role sections are additions** — Role-aware sections supplement the core template, they never replace `## Source`, `## Summary`, or `## Key Takeaways`.
9. **Respect the user's description** — If the user provides their own notes about why the resource matters, incorporate their perspective into the summary and takeaways rather than overriding it.
