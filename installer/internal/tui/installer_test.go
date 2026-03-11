package tui

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCopyRolePackTemplates(t *testing.T) {
	t.Run("creates core vault structure", func(t *testing.T) {
		repoDir := t.TempDir()
		projectDir := t.TempDir()

		// Set up mock repo structure with core templates
		coreTemplatesDir := filepath.Join(repoDir, "GentlemanNvim", "obsidian-brain", "core", "templates")
		os.MkdirAll(coreTemplatesDir, 0o755)
		os.WriteFile(filepath.Join(coreTemplatesDir, "braindump.md"), []byte("# Braindump"), 0o644)
		os.WriteFile(filepath.Join(coreTemplatesDir, "daily-note.md"), []byte("# Daily"), 0o644)

		err := copyRolePackTemplates(repoDir, projectDir, []string{"core"})
		if err != nil {
			t.Fatalf("copyRolePackTemplates failed: %v", err)
		}

		// Verify core directories
		vaultDir := filepath.Join(projectDir, ".obsidian-brain")
		for _, dir := range []string{"inbox", "resources", "knowledge", "templates", ".obsidian"} {
			path := filepath.Join(vaultDir, dir)
			info, err := os.Stat(path)
			if err != nil {
				t.Errorf("expected directory %s to exist: %v", dir, err)
				continue
			}
			if !info.IsDir() {
				t.Errorf("expected %s to be a directory", dir)
			}
		}

		// Verify templates were copied
		for _, tmpl := range []string{"braindump.md", "daily-note.md"} {
			path := filepath.Join(vaultDir, "templates", tmpl)
			if _, err := os.Stat(path); err != nil {
				t.Errorf("expected template %s to exist: %v", tmpl, err)
			}
		}
	})

	t.Run("creates developer-specific directories", func(t *testing.T) {
		repoDir := t.TempDir()
		projectDir := t.TempDir()

		// Set up mock repo structure with developer templates
		devTemplatesDir := filepath.Join(repoDir, "GentlemanNvim", "obsidian-brain", "developer", "templates")
		os.MkdirAll(devTemplatesDir, 0o755)
		os.WriteFile(filepath.Join(devTemplatesDir, "adr.md"), []byte("# ADR"), 0o644)

		coreTemplatesDir := filepath.Join(repoDir, "GentlemanNvim", "obsidian-brain", "core", "templates")
		os.MkdirAll(coreTemplatesDir, 0o755)

		err := copyRolePackTemplates(repoDir, projectDir, []string{"core", "developer"})
		if err != nil {
			t.Fatalf("copyRolePackTemplates failed: %v", err)
		}

		vaultDir := filepath.Join(projectDir, ".obsidian-brain")
		for _, dir := range []string{"architecture", "sessions", "debugging"} {
			path := filepath.Join(vaultDir, dir)
			if _, err := os.Stat(path); err != nil {
				t.Errorf("expected developer directory %s to exist: %v", dir, err)
			}
		}

		// Verify developer template was copied
		adrPath := filepath.Join(vaultDir, "templates", "adr.md")
		if _, err := os.Stat(adrPath); err != nil {
			t.Errorf("expected adr.md template to exist: %v", err)
		}
	})

	t.Run("creates pm-lead-specific directories", func(t *testing.T) {
		repoDir := t.TempDir()
		projectDir := t.TempDir()

		// Set up mock repo structure with pm-lead templates
		pmTemplatesDir := filepath.Join(repoDir, "GentlemanNvim", "obsidian-brain", "pm-lead", "templates")
		os.MkdirAll(pmTemplatesDir, 0o755)
		os.WriteFile(filepath.Join(pmTemplatesDir, "meeting-notes.md"), []byte("# Meeting"), 0o644)

		coreTemplatesDir := filepath.Join(repoDir, "GentlemanNvim", "obsidian-brain", "core", "templates")
		os.MkdirAll(coreTemplatesDir, 0o755)

		err := copyRolePackTemplates(repoDir, projectDir, []string{"core", "pm-lead"})
		if err != nil {
			t.Fatalf("copyRolePackTemplates failed: %v", err)
		}

		vaultDir := filepath.Join(projectDir, ".obsidian-brain")
		for _, dir := range []string{"meetings", "sprints", "risks", "briefs"} {
			path := filepath.Join(vaultDir, dir)
			if _, err := os.Stat(path); err != nil {
				t.Errorf("expected pm-lead directory %s to exist: %v", dir, err)
			}
		}

		// Verify pm-lead template was copied
		meetingPath := filepath.Join(vaultDir, "templates", "meeting-notes.md")
		if _, err := os.Stat(meetingPath); err != nil {
			t.Errorf("expected meeting-notes.md template to exist: %v", err)
		}
	})

	t.Run("all role packs creates all directories", func(t *testing.T) {
		repoDir := t.TempDir()
		projectDir := t.TempDir()

		// Set up mock repo structure
		for _, pack := range []string{"core", "developer", "pm-lead"} {
			dir := filepath.Join(repoDir, "GentlemanNvim", "obsidian-brain", pack, "templates")
			os.MkdirAll(dir, 0o755)
			os.WriteFile(filepath.Join(dir, pack+"-test.md"), []byte("# "+pack), 0o644)
		}

		err := copyRolePackTemplates(repoDir, projectDir, []string{"core", "developer", "pm-lead"})
		if err != nil {
			t.Fatalf("copyRolePackTemplates failed: %v", err)
		}

		vaultDir := filepath.Join(projectDir, ".obsidian-brain")
		// Core dirs
		for _, dir := range []string{"inbox", "resources", "knowledge"} {
			if _, err := os.Stat(filepath.Join(vaultDir, dir)); err != nil {
				t.Errorf("expected core directory %s: %v", dir, err)
			}
		}
		// Developer dirs
		for _, dir := range []string{"architecture", "sessions", "debugging"} {
			if _, err := os.Stat(filepath.Join(vaultDir, dir)); err != nil {
				t.Errorf("expected developer directory %s: %v", dir, err)
			}
		}
		// PM dirs
		for _, dir := range []string{"meetings", "sprints", "risks", "briefs"} {
			if _, err := os.Stat(filepath.Join(vaultDir, dir)); err != nil {
				t.Errorf("expected pm-lead directory %s: %v", dir, err)
			}
		}
		// .obsidian marker
		if _, err := os.Stat(filepath.Join(vaultDir, ".obsidian")); err != nil {
			t.Errorf("expected .obsidian marker directory: %v", err)
		}

		// All templates should be in the templates dir
		templatesDir := filepath.Join(vaultDir, "templates")
		for _, tmpl := range []string{"core-test.md", "developer-test.md", "pm-lead-test.md"} {
			if _, err := os.Stat(filepath.Join(templatesDir, tmpl)); err != nil {
				t.Errorf("expected template %s: %v", tmpl, err)
			}
		}
	})

	t.Run("skips gracefully when pack source dir missing", func(t *testing.T) {
		repoDir := t.TempDir()
		projectDir := t.TempDir()

		// Don't create any template directories — should not error
		err := copyRolePackTemplates(repoDir, projectDir, []string{"core", "developer"})
		if err != nil {
			t.Fatalf("expected no error when source dirs are missing, got: %v", err)
		}

		// Core directories should still be created (they're always created)
		vaultDir := filepath.Join(projectDir, ".obsidian-brain")
		if _, err := os.Stat(filepath.Join(vaultDir, "inbox")); err != nil {
			t.Errorf("expected inbox directory even without templates: %v", err)
		}
	})

	t.Run("template content is preserved", func(t *testing.T) {
		repoDir := t.TempDir()
		projectDir := t.TempDir()

		coreTemplatesDir := filepath.Join(repoDir, "GentlemanNvim", "obsidian-brain", "core", "templates")
		os.MkdirAll(coreTemplatesDir, 0o755)
		content := "---\ntitle: Test\ntags: [braindump]\n---\n\n## Thought\n"
		os.WriteFile(filepath.Join(coreTemplatesDir, "braindump.md"), []byte(content), 0o644)

		err := copyRolePackTemplates(repoDir, projectDir, []string{"core"})
		if err != nil {
			t.Fatalf("copyRolePackTemplates failed: %v", err)
		}

		dst := filepath.Join(projectDir, ".obsidian-brain", "templates", "braindump.md")
		data, err := os.ReadFile(dst)
		if err != nil {
			t.Fatalf("failed to read copied template: %v", err)
		}
		if string(data) != content {
			t.Errorf("expected template content to match.\nExpected: %q\nGot: %q", content, string(data))
		}
	})
}
