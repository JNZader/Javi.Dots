package main

import (
	"strings"
	"testing"
)

func TestParseRolePacks(t *testing.T) {
	t.Run("role pack requires obsidian-brain", func(t *testing.T) {
		_, err := parseRolePacks("developer", "simple")
		if err == nil {
			t.Fatal("expected error when using role pack without obsidian-brain")
		}
		if !strings.Contains(err.Error(), "obsidian-brain") {
			t.Errorf("error should mention obsidian-brain, got: %s", err.Error())
		}
	})

	t.Run("valid role pack values accepted", func(t *testing.T) {
		packs, err := parseRolePacks("developer,pm-lead", "obsidian-brain")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		expected := []string{"core", "developer", "pm-lead"}
		if len(packs) != len(expected) {
			t.Fatalf("expected %d packs, got %d: %v", len(expected), len(packs), packs)
		}
		for i, want := range expected {
			if packs[i] != want {
				t.Errorf("packs[%d] = %q, want %q", i, packs[i], want)
			}
		}
	})

	t.Run("invalid role pack value rejected", func(t *testing.T) {
		_, err := parseRolePacks("invalid-pack", "obsidian-brain")
		if err == nil {
			t.Fatal("expected error for invalid role pack")
		}
		if !strings.Contains(err.Error(), "developer") || !strings.Contains(err.Error(), "pm-lead") {
			t.Errorf("error should list valid options, got: %s", err.Error())
		}
	})

	t.Run("empty role pack with obsidian-brain defaults to core only", func(t *testing.T) {
		packs, err := parseRolePacks("", "obsidian-brain")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(packs) != 1 || packs[0] != "core" {
			t.Errorf("expected [core], got %v", packs)
		}
	})

	t.Run("core is auto-prepended and deduplicated", func(t *testing.T) {
		packs, err := parseRolePacks("core,developer", "obsidian-brain")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		expected := []string{"core", "developer"}
		if len(packs) != len(expected) {
			t.Fatalf("expected %d packs (no duplicate core), got %d: %v", len(expected), len(packs), packs)
		}
		for i, want := range expected {
			if packs[i] != want {
				t.Errorf("packs[%d] = %q, want %q", i, packs[i], want)
			}
		}
	})

	t.Run("whitespace is trimmed and case is normalized", func(t *testing.T) {
		packs, err := parseRolePacks("  Developer , PM-Lead  ", "obsidian-brain")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		expected := []string{"core", "developer", "pm-lead"}
		if len(packs) != len(expected) {
			t.Fatalf("expected %d packs, got %d: %v", len(expected), len(packs), packs)
		}
		for i, want := range expected {
			if packs[i] != want {
				t.Errorf("packs[%d] = %q, want %q", i, packs[i], want)
			}
		}
	})

	t.Run("non-obsidian-brain memory returns nil packs when empty", func(t *testing.T) {
		packs, err := parseRolePacks("", "simple")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(packs) != 0 {
			t.Errorf("expected empty packs for non-obsidian-brain, got %v", packs)
		}
	})

	t.Run("role pack rejected with none memory", func(t *testing.T) {
		_, err := parseRolePacks("developer", "none")
		if err == nil {
			t.Fatal("expected error when using role pack with memory=none")
		}
		if !strings.Contains(err.Error(), "obsidian-brain") {
			t.Errorf("error should mention obsidian-brain, got: %s", err.Error())
		}
	})
}
