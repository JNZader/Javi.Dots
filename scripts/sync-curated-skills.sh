#!/bin/bash
# sync-curated-skills.sh - Sync curated skills between GentlemanClaude/skills and GentlemanOpenCode/skill
#
# The canonical source is GentlemanClaude/skills/. This script copies all SKILL.md
# files from there to GentlemanOpenCode/skill/, ensuring both directories stay in sync.
#
# Usage: ./scripts/sync-curated-skills.sh [--check]
#   --check   Dry-run mode: report differences without modifying files

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(dirname "$SCRIPT_DIR")"

SOURCE="$REPO_ROOT/GentlemanClaude/skills"
TARGET="$REPO_ROOT/GentlemanOpenCode/skill"

CHECK_ONLY=false
[[ "${1:-}" == "--check" ]] && CHECK_ONLY=true

if [[ ! -d "$SOURCE" ]]; then
  echo "ERROR: Source directory not found: $SOURCE" >&2
  exit 1
fi

if [[ ! -d "$TARGET" ]]; then
  echo "ERROR: Target directory not found: $TARGET" >&2
  exit 1
fi

diffs=0
synced=0
total=0

for skill_dir in "$SOURCE"/*/; do
  skill_name="$(basename "$skill_dir")"
  src_file="$skill_dir/SKILL.md"
  dst_file="$TARGET/$skill_name/SKILL.md"
  total=$((total + 1))

  if [[ ! -f "$src_file" ]]; then
    echo "WARN: No SKILL.md in $skill_dir" >&2
    continue
  fi

  if [[ ! -f "$dst_file" ]]; then
    diffs=$((diffs + 1))
    if $CHECK_ONLY; then
      echo "MISSING: $skill_name (not in OpenCode)"
    else
      mkdir -p "$TARGET/$skill_name"
      cp "$src_file" "$dst_file"
      synced=$((synced + 1))
      echo "ADDED: $skill_name"
    fi
  elif ! diff -q "$src_file" "$dst_file" >/dev/null 2>&1; then
    diffs=$((diffs + 1))
    if $CHECK_ONLY; then
      echo "CHANGED: $skill_name"
      diff --color=auto "$src_file" "$dst_file" | head -20 || true
      echo "---"
    else
      cp "$src_file" "$dst_file"
      synced=$((synced + 1))
      echo "UPDATED: $skill_name"
    fi
  fi
done

# Check for orphans in target (skills that exist in OpenCode but not in Claude)
for target_dir in "$TARGET"/*/; do
  target_name="$(basename "$target_dir")"
  if [[ ! -d "$SOURCE/$target_name" ]]; then
    echo "ORPHAN: $target_name (in OpenCode but not in Claude)"
    diffs=$((diffs + 1))
  fi
done

echo ""
if $CHECK_ONLY; then
  if [[ $diffs -eq 0 ]]; then
    echo "OK: All $total skills are in sync"
    exit 0
  else
    echo "DRIFT: $diffs skill(s) differ between GentlemanClaude/skills and GentlemanOpenCode/skill"
    exit 1
  fi
else
  if [[ $synced -eq 0 ]]; then
    echo "OK: All $total skills already in sync"
  else
    echo "SYNCED: $synced of $total skills updated (Claude -> OpenCode)"
  fi
fi
