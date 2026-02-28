package tui

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestExpandPath(t *testing.T) {
	t.Run("should expand ~/foo to home dir + /foo", func(t *testing.T) {
		home, err := os.UserHomeDir()
		if err != nil {
			t.Fatalf("could not get home dir: %v", err)
		}
		result := expandPath("~/foo")
		expected := filepath.Join(home, "foo")
		if result != expected {
			t.Errorf("expected %q, got %q", expected, result)
		}
	})

	t.Run("should return absolute path unchanged", func(t *testing.T) {
		result := expandPath("/abs/path")
		if result != "/abs/path" {
			t.Errorf("expected /abs/path, got %q", result)
		}
	})

	t.Run("should return relative path unchanged", func(t *testing.T) {
		result := expandPath("relative")
		if result != "relative" {
			t.Errorf("expected 'relative', got %q", result)
		}
	})
}

func TestDetectStack(t *testing.T) {
	tests := []struct {
		name     string
		file     string
		expected string
	}{
		{"go.mod → go", "go.mod", "go"},
		{"package.json → node", "package.json", "node"},
		{"angular.json → angular", "angular.json", "angular"},
		{"Cargo.toml → rust", "Cargo.toml", "rust"},
		{"pyproject.toml → python", "pyproject.toml", "python"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			f, err := os.Create(filepath.Join(dir, tt.file))
			if err != nil {
				t.Fatalf("failed to create indicator file: %v", err)
			}
			f.Close()

			result := detectStack(dir)
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}

	t.Run("empty dir → unknown", func(t *testing.T) {
		dir := t.TempDir()
		result := detectStack(dir)
		if result != "unknown" {
			t.Errorf("expected 'unknown', got %q", result)
		}
	})
}

func TestProjectPathValidation(t *testing.T) {
	t.Run("empty input + Enter sets error, screen stays", func(t *testing.T) {
		m := NewModel()
		m.Screen = ScreenProjectPath
		m.ProjectPathInput = ""

		result, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
		nm := result.(Model)

		if nm.ProjectPathError == "" {
			t.Error("expected ProjectPathError to be set for empty input")
		}
		if nm.Screen != ScreenProjectPath {
			t.Errorf("expected screen to stay at ScreenProjectPath, got %d", nm.Screen)
		}
	})

	t.Run("non-existent path + Enter sets error", func(t *testing.T) {
		m := NewModel()
		m.Screen = ScreenProjectPath
		m.ProjectPathInput = "/tmp/absolutely-does-not-exist-12345"

		result, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
		nm := result.(Model)

		if nm.ProjectPathError == "" {
			t.Error("expected ProjectPathError to be set for non-existent path")
		}
		if nm.Screen != ScreenProjectPath {
			t.Errorf("expected screen to stay at ScreenProjectPath, got %d", nm.Screen)
		}
	})

	t.Run("valid directory path + Enter advances to ScreenProjectStack", func(t *testing.T) {
		dir := t.TempDir()
		// Create a go.mod to test stack detection
		f, _ := os.Create(filepath.Join(dir, "go.mod"))
		f.Close()

		m := NewModel()
		m.Screen = ScreenProjectPath
		m.ProjectPathInput = dir

		result, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
		nm := result.(Model)

		if nm.Screen != ScreenProjectStack {
			t.Errorf("expected ScreenProjectStack, got %d", nm.Screen)
		}
		if nm.ProjectPathError != "" {
			t.Errorf("expected no error, got %q", nm.ProjectPathError)
		}
		if nm.ProjectStack != "go" {
			t.Errorf("expected ProjectStack='go', got %q", nm.ProjectStack)
		}
	})
}

func TestProjectPathBackspace(t *testing.T) {
	t.Run("backspace removes character before cursor", func(t *testing.T) {
		m := NewModel()
		m.Screen = ScreenProjectPath
		m.ProjectPathInput = "hello"
		m.ProjectPathCursor = 5 // cursor at end

		result, _ := m.Update(tea.KeyMsg{Type: tea.KeyBackspace})
		nm := result.(Model)

		if nm.ProjectPathInput != "hell" {
			t.Errorf("expected 'hell', got %q", nm.ProjectPathInput)
		}
		if nm.ProjectPathCursor != 4 {
			t.Errorf("expected cursor at 4, got %d", nm.ProjectPathCursor)
		}
	})

	t.Run("multiple backspaces eventually reach empty string without panic", func(t *testing.T) {
		m := NewModel()
		m.Screen = ScreenProjectPath
		m.ProjectPathInput = "ab"
		m.ProjectPathCursor = 2

		for i := 0; i < 5; i++ {
			result, _ := m.Update(tea.KeyMsg{Type: tea.KeyBackspace})
			m = result.(Model)
		}

		if m.ProjectPathInput != "" {
			t.Errorf("expected empty string, got %q", m.ProjectPathInput)
		}
		if m.ProjectPathCursor != 0 {
			t.Errorf("expected cursor at 0, got %d", m.ProjectPathCursor)
		}
	})
}

func TestProjectPathTyping(t *testing.T) {
	t.Run("character keys accumulate in ProjectPathInput", func(t *testing.T) {
		m := NewModel()
		m.Screen = ScreenProjectPath
		m.ProjectPathInput = ""
		m.ProjectPathCursor = 0

		chars := []rune{'/', 't', 'm', 'p'}
		for _, c := range chars {
			result, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{c}})
			m = result.(Model)
		}

		if m.ProjectPathInput != "/tmp" {
			t.Errorf("expected '/tmp', got %q", m.ProjectPathInput)
		}
		if m.ProjectPathCursor != 4 {
			t.Errorf("expected cursor at 4, got %d", m.ProjectPathCursor)
		}
	})

	t.Run("space is included in input, does not activate leader mode", func(t *testing.T) {
		m := NewModel()
		m.Screen = ScreenProjectPath
		m.ProjectPathInput = "/my"
		m.ProjectPathCursor = 3

		// Send space
		result, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{' '}})
		nm := result.(Model)

		if !strings.Contains(nm.ProjectPathInput, " ") {
			t.Errorf("expected space in input, got %q", nm.ProjectPathInput)
		}
		if nm.LeaderMode {
			t.Error("space in ScreenProjectPath should NOT activate leader mode")
		}
	})
}

func TestProjectMemoryConditionalEngram(t *testing.T) {
	t.Run("obsidian-brain without obsidian installed goes to ScreenProjectObsidianInstall", func(t *testing.T) {
		// In test env, "obsidian" binary is NOT in PATH, so it goes to install screen
		m := NewModel()
		m.Screen = ScreenProjectMemory
		m.Cursor = 0 // obsidian-brain

		result, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
		nm := result.(Model)

		if nm.ProjectMemory != "obsidian-brain" {
			t.Errorf("expected ProjectMemory='obsidian-brain', got %q", nm.ProjectMemory)
		}
		// Obsidian binary not in PATH → goes to install screen
		if nm.Screen != ScreenProjectObsidianInstall {
			t.Errorf("expected ScreenProjectObsidianInstall, got %d", nm.Screen)
		}
	})

	t.Run("vibekanban skips Engram, goes to ScreenProjectCI", func(t *testing.T) {
		m := NewModel()
		m.Screen = ScreenProjectMemory
		m.Cursor = 1 // vibekanban

		result, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
		nm := result.(Model)

		if nm.Screen != ScreenProjectCI {
			t.Errorf("expected ScreenProjectCI, got %d", nm.Screen)
		}
		if nm.ProjectMemory != "vibekanban" {
			t.Errorf("expected ProjectMemory='vibekanban', got %q", nm.ProjectMemory)
		}
	})
}

func TestObsidianInstallSelection(t *testing.T) {
	t.Run("Yes sets InstallObsidian=true and goes to Engram", func(t *testing.T) {
		m := NewModel()
		m.Screen = ScreenProjectObsidianInstall
		m.Cursor = 0 // Yes

		result, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
		nm := result.(Model)

		if !nm.Choices.InstallObsidian {
			t.Error("expected InstallObsidian=true")
		}
		if nm.Screen != ScreenProjectEngram {
			t.Errorf("expected ScreenProjectEngram, got %d", nm.Screen)
		}
	})

	t.Run("No sets InstallObsidian=false and goes to Engram", func(t *testing.T) {
		m := NewModel()
		m.Screen = ScreenProjectObsidianInstall
		m.Cursor = 1 // No

		result, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
		nm := result.(Model)

		if nm.Choices.InstallObsidian {
			t.Error("expected InstallObsidian=false")
		}
		if nm.Screen != ScreenProjectEngram {
			t.Errorf("expected ScreenProjectEngram, got %d", nm.Screen)
		}
	})
}

func TestObsidianInstallBackNav(t *testing.T) {
	t.Run("backspace on ObsidianInstall goes back to Memory", func(t *testing.T) {
		m := NewModel()
		m.Screen = ScreenProjectObsidianInstall

		result, _ := m.Update(tea.KeyMsg{Type: tea.KeyBackspace})
		nm := result.(Model)

		if nm.Screen != ScreenProjectMemory {
			t.Errorf("expected ScreenProjectMemory, got %d", nm.Screen)
		}
	})

	t.Run("backspace on Engram goes to ObsidianInstall when obsidian not in PATH", func(t *testing.T) {
		m := NewModel()
		m.Screen = ScreenProjectEngram

		result, _ := m.Update(tea.KeyMsg{Type: tea.KeyBackspace})
		nm := result.(Model)

		// Obsidian not in PATH → back goes to ObsidianInstall
		if nm.Screen != ScreenProjectObsidianInstall {
			t.Errorf("expected ScreenProjectObsidianInstall, got %d", nm.Screen)
		}
	})
}

func TestObsidianInstallScreenOptions(t *testing.T) {
	t.Run("has 2 options", func(t *testing.T) {
		m := NewModel()
		m.Screen = ScreenProjectObsidianInstall
		opts := m.GetCurrentOptions()

		if len(opts) != 2 {
			t.Errorf("expected 2 options, got %d: %v", len(opts), opts)
		}
	})

	t.Run("title is non-empty", func(t *testing.T) {
		m := NewModel()
		m.Screen = ScreenProjectObsidianInstall
		title := m.GetScreenTitle()
		if title == "" {
			t.Error("expected non-empty title")
		}
	})

	t.Run("description mentions obsidian", func(t *testing.T) {
		m := NewModel()
		m.Screen = ScreenProjectObsidianInstall
		desc := m.GetScreenDescription()
		if !strings.Contains(strings.ToLower(desc), "obsidian") {
			t.Errorf("expected description to mention obsidian, got %q", desc)
		}
	})
}

func TestProjectCIAdvancesToConfirm(t *testing.T) {
	t.Run("selecting CI advances to ScreenProjectConfirm", func(t *testing.T) {
		m := NewModel()
		m.Screen = ScreenProjectCI
		m.Cursor = 0 // GitHub Actions

		result, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
		nm := result.(Model)

		if nm.Screen != ScreenProjectConfirm {
			t.Errorf("expected ScreenProjectConfirm, got %d", nm.Screen)
		}
		if nm.ProjectCI != "github" {
			t.Errorf("expected ProjectCI='github', got %q", nm.ProjectCI)
		}
	})
}

func TestProjectEscapeBackNavigation(t *testing.T) {
	// Selection-based project screens use backspace (via handleSelectionKeys)
	// for back navigation, since ESC is intercepted by handleEscape() first.
	// ScreenProjectPath uses ESC via handleEscape() directly.

	t.Run("ScreenProjectStack → Backspace → ScreenProjectPath", func(t *testing.T) {
		m := NewModel()
		m.Screen = ScreenProjectStack

		result, _ := m.Update(tea.KeyMsg{Type: tea.KeyBackspace})
		nm := result.(Model)

		if nm.Screen != ScreenProjectPath {
			t.Errorf("expected ScreenProjectPath, got %d", nm.Screen)
		}
	})

	t.Run("ScreenProjectMemory → Backspace → ScreenProjectStack", func(t *testing.T) {
		m := NewModel()
		m.Screen = ScreenProjectMemory

		result, _ := m.Update(tea.KeyMsg{Type: tea.KeyBackspace})
		nm := result.(Model)

		if nm.Screen != ScreenProjectStack {
			t.Errorf("expected ScreenProjectStack, got %d", nm.Screen)
		}
	})

	t.Run("ScreenProjectEngram → Backspace → ScreenProjectObsidianInstall (obsidian not in PATH)", func(t *testing.T) {
		m := NewModel()
		m.Screen = ScreenProjectEngram

		result, _ := m.Update(tea.KeyMsg{Type: tea.KeyBackspace})
		nm := result.(Model)

		// Obsidian not in PATH → back goes to ObsidianInstall
		if nm.Screen != ScreenProjectObsidianInstall {
			t.Errorf("expected ScreenProjectObsidianInstall, got %d", nm.Screen)
		}
	})

	t.Run("ScreenProjectCI with memory=obsidian-brain → Backspace → ScreenProjectEngram", func(t *testing.T) {
		m := NewModel()
		m.Screen = ScreenProjectCI
		m.ProjectMemory = "obsidian-brain"

		result, _ := m.Update(tea.KeyMsg{Type: tea.KeyBackspace})
		nm := result.(Model)

		if nm.Screen != ScreenProjectEngram {
			t.Errorf("expected ScreenProjectEngram, got %d", nm.Screen)
		}
	})

	t.Run("ScreenProjectCI with memory=simple → Backspace → ScreenProjectMemory", func(t *testing.T) {
		m := NewModel()
		m.Screen = ScreenProjectCI
		m.ProjectMemory = "simple"

		result, _ := m.Update(tea.KeyMsg{Type: tea.KeyBackspace})
		nm := result.(Model)

		if nm.Screen != ScreenProjectMemory {
			t.Errorf("expected ScreenProjectMemory, got %d", nm.Screen)
		}
	})

	t.Run("ScreenProjectConfirm → Backspace → ScreenProjectCI", func(t *testing.T) {
		m := NewModel()
		m.Screen = ScreenProjectConfirm

		result, _ := m.Update(tea.KeyMsg{Type: tea.KeyBackspace})
		nm := result.(Model)

		if nm.Screen != ScreenProjectCI {
			t.Errorf("expected ScreenProjectCI, got %d", nm.Screen)
		}
	})

	t.Run("ScreenProjectPath typing mode → Esc → ScreenMainMenu", func(t *testing.T) {
		m := NewModel()
		m.Screen = ScreenProjectPath
		m.ProjectPathMode = PathModeTyping

		result, _ := m.Update(tea.KeyMsg{Type: tea.KeyEsc})
		nm := result.(Model)

		if nm.Screen != ScreenMainMenu {
			t.Errorf("expected ScreenMainMenu, got %d", nm.Screen)
		}
	})

	t.Run("ScreenProjectPath completion mode → Esc → stays, returns to typing", func(t *testing.T) {
		m := NewModel()
		m.Screen = ScreenProjectPath
		m.ProjectPathMode = PathModeCompletion
		m.ProjectPathCompletions = []string{"foo", "bar"}

		result, _ := m.Update(tea.KeyMsg{Type: tea.KeyEsc})
		nm := result.(Model)

		if nm.Screen != ScreenProjectPath {
			t.Errorf("expected ScreenProjectPath, got %d", nm.Screen)
		}
		if nm.ProjectPathMode != PathModeTyping {
			t.Errorf("expected PathModeTyping, got %d", nm.ProjectPathMode)
		}
	})

	t.Run("ScreenProjectPath browser mode → Esc → stays, returns to typing", func(t *testing.T) {
		m := NewModel()
		m.Screen = ScreenProjectPath
		m.ProjectPathMode = PathModeBrowser
		m.FileBrowserEntries = []string{"dir1"}

		result, _ := m.Update(tea.KeyMsg{Type: tea.KeyEsc})
		nm := result.(Model)

		if nm.Screen != ScreenProjectPath {
			t.Errorf("expected ScreenProjectPath, got %d", nm.Screen)
		}
		if nm.ProjectPathMode != PathModeTyping {
			t.Errorf("expected PathModeTyping, got %d", nm.ProjectPathMode)
		}
	})
}

func TestProjectConfirmCancel(t *testing.T) {
	t.Run("Cursor=1 (Cancel) → ScreenMainMenu", func(t *testing.T) {
		m := NewModel()
		m.Screen = ScreenProjectConfirm
		m.Cursor = 1 // Cancel

		result, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
		nm := result.(Model)

		if nm.Screen != ScreenMainMenu {
			t.Errorf("expected ScreenMainMenu, got %d", nm.Screen)
		}
	})
}

func TestGetCurrentOptionsProjectScreens(t *testing.T) {
	t.Run("ScreenProjectMemory → 5 options", func(t *testing.T) {
		m := NewModel()
		m.Screen = ScreenProjectMemory
		opts := m.GetCurrentOptions()

		if len(opts) != 5 {
			t.Errorf("expected 5 options, got %d: %v", len(opts), opts)
		}
	})

	t.Run("ScreenProjectEngram → 2 options", func(t *testing.T) {
		m := NewModel()
		m.Screen = ScreenProjectEngram
		opts := m.GetCurrentOptions()

		if len(opts) != 2 {
			t.Errorf("expected 2 options, got %d: %v", len(opts), opts)
		}
	})

	t.Run("ScreenProjectCI → 4 options", func(t *testing.T) {
		m := NewModel()
		m.Screen = ScreenProjectCI
		opts := m.GetCurrentOptions()

		if len(opts) != 4 {
			t.Errorf("expected 4 options, got %d: %v", len(opts), opts)
		}
	})

	t.Run("ScreenProjectConfirm → 2 options", func(t *testing.T) {
		m := NewModel()
		m.Screen = ScreenProjectConfirm
		opts := m.GetCurrentOptions()

		if len(opts) != 2 {
			t.Errorf("expected 2 options, got %d: %v", len(opts), opts)
		}
	})
}

func TestGetScreenTitleProjectScreens(t *testing.T) {
	screens := []Screen{
		ScreenProjectPath,
		ScreenProjectStack,
		ScreenProjectMemory,
		ScreenProjectObsidianInstall,
		ScreenProjectEngram,
		ScreenProjectCI,
		ScreenProjectConfirm,
		ScreenProjectInstalling,
		ScreenProjectResult,
	}

	m := NewModel()
	for _, s := range screens {
		t.Run("screen title non-empty", func(t *testing.T) {
			m.Screen = s
			title := m.GetScreenTitle()
			if title == "" {
				t.Errorf("expected non-empty title for screen %d", s)
			}
		})
	}
}

func TestGetScreenDescriptionProjectStack(t *testing.T) {
	t.Run("with ProjectStack=go, description contains 'go'", func(t *testing.T) {
		m := NewModel()
		m.Screen = ScreenProjectStack
		m.ProjectStack = "go"

		desc := m.GetScreenDescription()
		if !strings.Contains(strings.ToLower(desc), "go") {
			t.Errorf("expected description to contain 'go', got %q", desc)
		}
	})

	t.Run("with ProjectStack empty, description is generic", func(t *testing.T) {
		m := NewModel()
		m.Screen = ScreenProjectStack
		m.ProjectStack = ""

		desc := m.GetScreenDescription()
		if desc == "" {
			t.Error("expected non-empty generic description")
		}
		// Should NOT contain "Auto-detected"
		if strings.Contains(desc, "Auto-detected") {
			t.Errorf("empty stack should give generic description, got %q", desc)
		}
	})
}

func TestProjectResultEnter(t *testing.T) {
	t.Run("Enter on ScreenProjectResult → ScreenMainMenu", func(t *testing.T) {
		m := NewModel()
		m.Screen = ScreenProjectResult

		result, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
		nm := result.(Model)

		if nm.Screen != ScreenMainMenu {
			t.Errorf("expected ScreenMainMenu, got %d", nm.Screen)
		}
	})
}

// --- Cursor navigation tests ---

func TestCursorLeftRight(t *testing.T) {
	t.Run("left moves cursor back", func(t *testing.T) {
		m := NewModel()
		m.Screen = ScreenProjectPath
		m.ProjectPathInput = "/tmp/foo"
		m.ProjectPathCursor = 8

		result, _ := m.Update(tea.KeyMsg{Type: tea.KeyLeft})
		nm := result.(Model)

		if nm.ProjectPathCursor != 7 {
			t.Errorf("expected cursor at 7, got %d", nm.ProjectPathCursor)
		}
	})

	t.Run("right moves cursor forward", func(t *testing.T) {
		m := NewModel()
		m.Screen = ScreenProjectPath
		m.ProjectPathInput = "/tmp/foo"
		m.ProjectPathCursor = 4

		result, _ := m.Update(tea.KeyMsg{Type: tea.KeyRight})
		nm := result.(Model)

		if nm.ProjectPathCursor != 5 {
			t.Errorf("expected cursor at 5, got %d", nm.ProjectPathCursor)
		}
	})

	t.Run("left at start stays at 0", func(t *testing.T) {
		m := NewModel()
		m.Screen = ScreenProjectPath
		m.ProjectPathInput = "/tmp"
		m.ProjectPathCursor = 0

		result, _ := m.Update(tea.KeyMsg{Type: tea.KeyLeft})
		nm := result.(Model)

		if nm.ProjectPathCursor != 0 {
			t.Errorf("expected cursor at 0, got %d", nm.ProjectPathCursor)
		}
	})

	t.Run("right at end stays at end", func(t *testing.T) {
		m := NewModel()
		m.Screen = ScreenProjectPath
		m.ProjectPathInput = "/tmp"
		m.ProjectPathCursor = 4

		result, _ := m.Update(tea.KeyMsg{Type: tea.KeyRight})
		nm := result.(Model)

		if nm.ProjectPathCursor != 4 {
			t.Errorf("expected cursor at 4, got %d", nm.ProjectPathCursor)
		}
	})
}

func TestCursorHomeEnd(t *testing.T) {
	t.Run("ctrl+a moves cursor to start", func(t *testing.T) {
		m := NewModel()
		m.Screen = ScreenProjectPath
		m.ProjectPathInput = "/tmp/foo"
		m.ProjectPathCursor = 5

		result, _ := m.handlePathTypingKeys("ctrl+a")
		nm := result.(Model)

		if nm.ProjectPathCursor != 0 {
			t.Errorf("expected cursor at 0, got %d", nm.ProjectPathCursor)
		}
	})

	t.Run("ctrl+e moves cursor to end", func(t *testing.T) {
		m := NewModel()
		m.Screen = ScreenProjectPath
		m.ProjectPathInput = "/tmp/foo"
		m.ProjectPathCursor = 2

		result, _ := m.handlePathTypingKeys("ctrl+e")
		nm := result.(Model)

		if nm.ProjectPathCursor != 8 {
			t.Errorf("expected cursor at 8, got %d", nm.ProjectPathCursor)
		}
	})

	t.Run("home key moves cursor to start", func(t *testing.T) {
		m := NewModel()
		m.Screen = ScreenProjectPath
		m.ProjectPathInput = "/tmp/foo"
		m.ProjectPathCursor = 5

		result, _ := m.Update(tea.KeyMsg{Type: tea.KeyHome})
		nm := result.(Model)

		if nm.ProjectPathCursor != 0 {
			t.Errorf("expected cursor at 0, got %d", nm.ProjectPathCursor)
		}
	})

	t.Run("end key moves cursor to end", func(t *testing.T) {
		m := NewModel()
		m.Screen = ScreenProjectPath
		m.ProjectPathInput = "/tmp/foo"
		m.ProjectPathCursor = 2

		result, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnd})
		nm := result.(Model)

		if nm.ProjectPathCursor != 8 {
			t.Errorf("expected cursor at 8, got %d", nm.ProjectPathCursor)
		}
	})
}

func TestCursorInsertMiddle(t *testing.T) {
	t.Run("insert char at cursor position", func(t *testing.T) {
		m := NewModel()
		m.Screen = ScreenProjectPath
		m.ProjectPathInput = "/tmpfoo"
		m.ProjectPathCursor = 4

		result, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}})
		nm := result.(Model)

		if nm.ProjectPathInput != "/tmp/foo" {
			t.Errorf("expected '/tmp/foo', got %q", nm.ProjectPathInput)
		}
		if nm.ProjectPathCursor != 5 {
			t.Errorf("expected cursor at 5, got %d", nm.ProjectPathCursor)
		}
	})
}

func TestCursorBackspaceMiddle(t *testing.T) {
	t.Run("backspace at cursor position in middle", func(t *testing.T) {
		m := NewModel()
		m.Screen = ScreenProjectPath
		m.ProjectPathInput = "/tmp/foo"
		m.ProjectPathCursor = 4

		result, _ := m.Update(tea.KeyMsg{Type: tea.KeyBackspace})
		nm := result.(Model)

		if nm.ProjectPathInput != "/tm/foo" {
			t.Errorf("expected '/tm/foo', got %q", nm.ProjectPathInput)
		}
		if nm.ProjectPathCursor != 3 {
			t.Errorf("expected cursor at 3, got %d", nm.ProjectPathCursor)
		}
	})
}

func TestCursorCtrlU(t *testing.T) {
	t.Run("ctrl+u clears input and resets cursor", func(t *testing.T) {
		m := NewModel()
		m.Screen = ScreenProjectPath
		m.ProjectPathInput = "/tmp/foo/bar"
		m.ProjectPathCursor = 8

		result, _ := m.handlePathTypingKeys("ctrl+u")
		nm := result.(Model)

		if nm.ProjectPathInput != "" {
			t.Errorf("expected empty string, got %q", nm.ProjectPathInput)
		}
		if nm.ProjectPathCursor != 0 {
			t.Errorf("expected cursor at 0, got %d", nm.ProjectPathCursor)
		}
	})
}

func TestCursorCtrlW(t *testing.T) {
	t.Run("ctrl+w deletes to previous /", func(t *testing.T) {
		m := NewModel()
		m.Screen = ScreenProjectPath
		m.ProjectPathInput = "/tmp/foo/bar"
		m.ProjectPathCursor = 12

		result, _ := m.handlePathTypingKeys("ctrl+w")
		nm := result.(Model)

		if nm.ProjectPathInput != "/tmp/foo/" {
			t.Errorf("expected '/tmp/foo/', got %q", nm.ProjectPathInput)
		}
	})

	t.Run("ctrl+w at start does nothing", func(t *testing.T) {
		m := NewModel()
		m.Screen = ScreenProjectPath
		m.ProjectPathInput = "/tmp"
		m.ProjectPathCursor = 0

		result, _ := m.handlePathTypingKeys("ctrl+w")
		nm := result.(Model)

		if nm.ProjectPathInput != "/tmp" {
			t.Errorf("expected '/tmp', got %q", nm.ProjectPathInput)
		}
		if nm.ProjectPathCursor != 0 {
			t.Errorf("expected cursor at 0, got %d", nm.ProjectPathCursor)
		}
	})
}

// --- Tab-completion tests ---

func TestTabSingleMatch(t *testing.T) {
	t.Run("single match auto-completes inline", func(t *testing.T) {
		dir := t.TempDir()
		os.Mkdir(filepath.Join(dir, "unique-dir"), 0o755)

		m := NewModel()
		m.Screen = ScreenProjectPath
		m.ProjectPathInput = filepath.Join(dir, "uni")
		m.ProjectPathCursor = len([]rune(m.ProjectPathInput))

		result, _ := m.handlePathTypingKeys("tab")
		nm := result.(Model)

		expected := filepath.Join(dir, "unique-dir") + "/"
		if nm.ProjectPathInput != expected {
			t.Errorf("expected %q, got %q", expected, nm.ProjectPathInput)
		}
		if nm.ProjectPathMode != PathModeTyping {
			t.Errorf("expected PathModeTyping, got %d", nm.ProjectPathMode)
		}
	})
}

func TestTabMultipleMatches(t *testing.T) {
	t.Run("multiple matches show dropdown", func(t *testing.T) {
		dir := t.TempDir()
		os.Mkdir(filepath.Join(dir, "projects"), 0o755)
		os.Mkdir(filepath.Join(dir, "prometheus"), 0o755)

		m := NewModel()
		m.Screen = ScreenProjectPath
		m.ProjectPathInput = filepath.Join(dir, "pro")
		m.ProjectPathCursor = len([]rune(m.ProjectPathInput))

		result, _ := m.handlePathTypingKeys("tab")
		nm := result.(Model)

		if nm.ProjectPathMode != PathModeCompletion {
			t.Errorf("expected PathModeCompletion, got %d", nm.ProjectPathMode)
		}
		if len(nm.ProjectPathCompletions) != 2 {
			t.Errorf("expected 2 completions, got %d", len(nm.ProjectPathCompletions))
		}
		if nm.ProjectPathCompIdx != 0 {
			t.Errorf("expected comp index 0, got %d", nm.ProjectPathCompIdx)
		}
	})
}

func TestTabNoMatch(t *testing.T) {
	t.Run("no matches show error message", func(t *testing.T) {
		dir := t.TempDir()

		m := NewModel()
		m.Screen = ScreenProjectPath
		m.ProjectPathInput = filepath.Join(dir, "nonexistent")
		m.ProjectPathCursor = len([]rune(m.ProjectPathInput))

		result, _ := m.handlePathTypingKeys("tab")
		nm := result.(Model)

		if nm.ProjectPathError == "" {
			t.Error("expected error message for no matches")
		}
		if nm.ProjectPathMode != PathModeTyping {
			t.Errorf("expected PathModeTyping, got %d", nm.ProjectPathMode)
		}
	})
}

func TestTabSelectCompletion(t *testing.T) {
	t.Run("down + enter selects completion", func(t *testing.T) {
		dir := t.TempDir()
		os.Mkdir(filepath.Join(dir, "alpha"), 0o755)
		os.Mkdir(filepath.Join(dir, "beta"), 0o755)

		m := NewModel()
		m.Screen = ScreenProjectPath
		m.ProjectPathMode = PathModeCompletion
		m.ProjectPathInput = dir + "/"
		m.ProjectPathCompletions = []string{"alpha", "beta"}
		m.ProjectPathCompIdx = 0

		// Move down to "beta"
		result, _ := m.handlePathCompletionKeys("down")
		m = result.(Model)
		if m.ProjectPathCompIdx != 1 {
			t.Errorf("expected comp index 1, got %d", m.ProjectPathCompIdx)
		}

		// Select with enter
		result, _ = m.handlePathCompletionKeys("enter")
		nm := result.(Model)

		expected := filepath.Join(dir, "beta") + "/"
		if nm.ProjectPathInput != expected {
			t.Errorf("expected %q, got %q", expected, nm.ProjectPathInput)
		}
		if nm.ProjectPathMode != PathModeTyping {
			t.Errorf("expected PathModeTyping, got %d", nm.ProjectPathMode)
		}
	})
}

func TestTabEscCancels(t *testing.T) {
	t.Run("esc in completion returns to typing", func(t *testing.T) {
		m := NewModel()
		m.Screen = ScreenProjectPath
		m.ProjectPathMode = PathModeCompletion
		m.ProjectPathCompletions = []string{"foo", "bar"}
		m.ProjectPathCompIdx = 0
		m.ProjectPathInput = "/tmp/"

		result, _ := m.handlePathCompletionKeys("esc")
		nm := result.(Model)

		if nm.ProjectPathMode != PathModeTyping {
			t.Errorf("expected PathModeTyping, got %d", nm.ProjectPathMode)
		}
		if nm.ProjectPathCompletions != nil {
			t.Error("expected completions to be cleared")
		}
		// Input should be unchanged
		if nm.ProjectPathInput != "/tmp/" {
			t.Errorf("expected input unchanged, got %q", nm.ProjectPathInput)
		}
	})
}

// --- File browser tests ---

func TestBrowserOpenCtrlB(t *testing.T) {
	t.Run("ctrl+b opens browser mode", func(t *testing.T) {
		dir := t.TempDir()
		os.Mkdir(filepath.Join(dir, "subdir"), 0o755)

		m := NewModel()
		m.Screen = ScreenProjectPath
		m.ProjectPathInput = dir
		m.ProjectPathCursor = len([]rune(dir))
		m.Height = 30

		result, _ := m.handlePathTypingKeys("ctrl+b")
		nm := result.(Model)

		if nm.ProjectPathMode != PathModeBrowser {
			t.Errorf("expected PathModeBrowser, got %d", nm.ProjectPathMode)
		}
		if nm.FileBrowserRoot != dir {
			t.Errorf("expected root %q, got %q", dir, nm.FileBrowserRoot)
		}
		if len(nm.FileBrowserEntries) == 0 {
			t.Error("expected at least one entry")
		}
	})
}

func TestBrowserNavigation(t *testing.T) {
	t.Run("j/k moves cursor in browser", func(t *testing.T) {
		m := NewModel()
		m.Screen = ScreenProjectPath
		m.ProjectPathMode = PathModeBrowser
		m.FileBrowserEntries = []string{"alpha", "beta", "gamma"}
		m.FileBrowserCursor = 0
		m.Height = 30

		// Move down
		result, _ := m.handlePathBrowserKeys("j")
		m = result.(Model)
		if m.FileBrowserCursor != 1 {
			t.Errorf("expected cursor at 1, got %d", m.FileBrowserCursor)
		}

		result, _ = m.handlePathBrowserKeys("j")
		m = result.(Model)
		if m.FileBrowserCursor != 2 {
			t.Errorf("expected cursor at 2, got %d", m.FileBrowserCursor)
		}

		// Move up
		result, _ = m.handlePathBrowserKeys("k")
		m = result.(Model)
		if m.FileBrowserCursor != 1 {
			t.Errorf("expected cursor at 1, got %d", m.FileBrowserCursor)
		}
	})
}

func TestBrowserDrillIn(t *testing.T) {
	t.Run("enter on subdir drills into it", func(t *testing.T) {
		dir := t.TempDir()
		subdir := filepath.Join(dir, "myproject")
		os.Mkdir(subdir, 0o755)
		os.Mkdir(filepath.Join(subdir, "src"), 0o755)

		m := NewModel()
		m.Screen = ScreenProjectPath
		m.ProjectPathMode = PathModeBrowser
		m.FileBrowserRoot = dir
		m.FileBrowserEntries = []string{"myproject"}
		m.FileBrowserCursor = 2 // index 2 = first subdirectory entry
		m.Height = 30

		result, _ := m.handlePathBrowserKeys("enter")
		nm := result.(Model)

		if nm.FileBrowserRoot != subdir {
			t.Errorf("expected root %q, got %q", subdir, nm.FileBrowserRoot)
		}
		if nm.FileBrowserCursor != 0 {
			t.Errorf("expected cursor reset to 0, got %d", nm.FileBrowserCursor)
		}
	})
}

func TestBrowserGoUp(t *testing.T) {
	t.Run("h goes to parent directory", func(t *testing.T) {
		dir := t.TempDir()
		subdir := filepath.Join(dir, "child")
		os.Mkdir(subdir, 0o755)

		m := NewModel()
		m.Screen = ScreenProjectPath
		m.ProjectPathMode = PathModeBrowser
		m.FileBrowserRoot = subdir
		m.FileBrowserEntries = []string{}
		m.FileBrowserCursor = 0
		m.Height = 30

		result, _ := m.handlePathBrowserKeys("h")
		nm := result.(Model)

		if nm.FileBrowserRoot != dir {
			t.Errorf("expected root %q, got %q", dir, nm.FileBrowserRoot)
		}
	})
}

func TestBrowserSelectDir(t *testing.T) {
	t.Run("select this directory sets input path", func(t *testing.T) {
		dir := t.TempDir()

		m := NewModel()
		m.Screen = ScreenProjectPath
		m.ProjectPathMode = PathModeBrowser
		m.FileBrowserRoot = dir
		m.FileBrowserEntries = []string{"subdir"}
		m.FileBrowserCursor = 0 // "Select this directory"
		m.Height = 30

		result, _ := m.handlePathBrowserKeys("enter")
		nm := result.(Model)

		if nm.ProjectPathInput != dir {
			t.Errorf("expected input %q, got %q", dir, nm.ProjectPathInput)
		}
		if nm.ProjectPathMode != PathModeTyping {
			t.Errorf("expected PathModeTyping, got %d", nm.ProjectPathMode)
		}
	})
}

func TestBrowserEscClose(t *testing.T) {
	t.Run("esc closes browser and returns to typing", func(t *testing.T) {
		m := NewModel()
		m.Screen = ScreenProjectPath
		m.ProjectPathMode = PathModeBrowser
		m.FileBrowserEntries = []string{"dir1"}
		m.FileBrowserCursor = 0
		m.Height = 30

		result, _ := m.handlePathBrowserKeys("esc")
		nm := result.(Model)

		if nm.ProjectPathMode != PathModeTyping {
			t.Errorf("expected PathModeTyping, got %d", nm.ProjectPathMode)
		}
		if nm.FileBrowserEntries != nil {
			t.Error("expected entries to be cleared")
		}
	})
}

func TestBrowserToggleHidden(t *testing.T) {
	t.Run("dot key toggles hidden files", func(t *testing.T) {
		dir := t.TempDir()
		os.Mkdir(filepath.Join(dir, ".hidden"), 0o755)
		os.Mkdir(filepath.Join(dir, "visible"), 0o755)

		m := NewModel()
		m.Screen = ScreenProjectPath
		m.ProjectPathMode = PathModeBrowser
		m.FileBrowserRoot = dir
		m.FileBrowserShowHidden = false
		m.FileBrowserEntries = listDirectories(dir, "", false)
		m.Height = 30

		// Initially should not have .hidden
		for _, e := range m.FileBrowserEntries {
			if e == ".hidden" {
				t.Error("hidden dir should not be visible initially")
			}
		}

		// Toggle hidden
		result, _ := m.handlePathBrowserKeys(".")
		nm := result.(Model)

		if !nm.FileBrowserShowHidden {
			t.Error("expected FileBrowserShowHidden=true")
		}
		found := false
		for _, e := range nm.FileBrowserEntries {
			if e == ".hidden" {
				found = true
				break
			}
		}
		if !found {
			t.Error("expected .hidden to be visible after toggle")
		}
	})
}

// --- Pre-fill + helper tests ---

func TestProjectPathPrefilledWithCwd(t *testing.T) {
	t.Run("input starts with cwd when entering from main menu", func(t *testing.T) {
		cwd, err := os.Getwd()
		if err != nil {
			t.Fatalf("could not get cwd: %v", err)
		}

		m := NewModel()
		m.Screen = ScreenMainMenu
		// Find the Initialize Project option index
		opts := m.GetCurrentOptions()
		for i, opt := range opts {
			if strings.Contains(opt, "Initialize Project") {
				m.Cursor = i
				break
			}
		}

		result, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
		nm := result.(Model)

		if nm.Screen != ScreenProjectPath {
			t.Errorf("expected ScreenProjectPath, got %d", nm.Screen)
		}
		if nm.ProjectPathInput != cwd {
			t.Errorf("expected input %q, got %q", cwd, nm.ProjectPathInput)
		}
		if nm.ProjectPathCursor != len([]rune(cwd)) {
			t.Errorf("expected cursor at %d, got %d", len([]rune(cwd)), nm.ProjectPathCursor)
		}
	})
}

func TestContractHome(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("could not get home dir: %v", err)
	}

	t.Run("home dir becomes ~", func(t *testing.T) {
		result := contractHome(home)
		if result != "~" {
			t.Errorf("expected '~', got %q", result)
		}
	})

	t.Run("subdir of home uses ~", func(t *testing.T) {
		result := contractHome(filepath.Join(home, "projects"))
		if result != "~/projects" {
			t.Errorf("expected '~/projects', got %q", result)
		}
	})

	t.Run("non-home path unchanged", func(t *testing.T) {
		result := contractHome("/tmp/foo")
		if result != "/tmp/foo" {
			t.Errorf("expected '/tmp/foo', got %q", result)
		}
	})
}

func TestSplitPathForCompletion(t *testing.T) {
	t.Run("path with partial name", func(t *testing.T) {
		parent, prefix := splitPathForCompletion("/home/user/pro")
		if parent != "/home/user" {
			t.Errorf("expected parent '/home/user', got %q", parent)
		}
		if prefix != "pro" {
			t.Errorf("expected prefix 'pro', got %q", prefix)
		}
	})

	t.Run("path with trailing slash", func(t *testing.T) {
		parent, prefix := splitPathForCompletion("/home/user/")
		if parent != "/home/user" {
			t.Errorf("expected parent '/home/user', got %q", parent)
		}
		if prefix != "" {
			t.Errorf("expected empty prefix, got %q", prefix)
		}
	})

	t.Run("empty path returns home dir", func(t *testing.T) {
		home, err := os.UserHomeDir()
		if err != nil {
			t.Fatalf("could not get home dir: %v", err)
		}
		parent, prefix := splitPathForCompletion("")
		if parent != home {
			t.Errorf("expected parent %q, got %q", home, parent)
		}
		if prefix != "" {
			t.Errorf("expected empty prefix, got %q", prefix)
		}
	})
}

func TestListDirectories(t *testing.T) {
	t.Run("lists subdirectories matching prefix", func(t *testing.T) {
		dir := t.TempDir()
		os.Mkdir(filepath.Join(dir, "projects"), 0o755)
		os.Mkdir(filepath.Join(dir, "prometheus"), 0o755)
		os.Mkdir(filepath.Join(dir, "docs"), 0o755)
		// Create a regular file (should not appear)
		os.WriteFile(filepath.Join(dir, "proFile.txt"), []byte("test"), 0o644)

		dirs := listDirectories(dir, "pro", false)
		if len(dirs) != 2 {
			t.Errorf("expected 2 dirs, got %d: %v", len(dirs), dirs)
		}
	})

	t.Run("hides dotfiles when showHidden=false", func(t *testing.T) {
		dir := t.TempDir()
		os.Mkdir(filepath.Join(dir, ".config"), 0o755)
		os.Mkdir(filepath.Join(dir, "visible"), 0o755)

		dirs := listDirectories(dir, "", false)
		for _, d := range dirs {
			if strings.HasPrefix(d, ".") {
				t.Errorf("should not include hidden dir %q", d)
			}
		}
	})

	t.Run("shows dotfiles when showHidden=true", func(t *testing.T) {
		dir := t.TempDir()
		os.Mkdir(filepath.Join(dir, ".config"), 0o755)
		os.Mkdir(filepath.Join(dir, "visible"), 0o755)

		dirs := listDirectories(dir, "", true)
		found := false
		for _, d := range dirs {
			if d == ".config" {
				found = true
				break
			}
		}
		if !found {
			t.Error("expected .config in results when showHidden=true")
		}
	})

	t.Run("returns sorted results", func(t *testing.T) {
		dir := t.TempDir()
		os.Mkdir(filepath.Join(dir, "zeta"), 0o755)
		os.Mkdir(filepath.Join(dir, "alpha"), 0o755)
		os.Mkdir(filepath.Join(dir, "beta"), 0o755)

		dirs := listDirectories(dir, "", false)
		if len(dirs) < 3 {
			t.Fatalf("expected 3 dirs, got %d", len(dirs))
		}
		if dirs[0] != "alpha" || dirs[1] != "beta" || dirs[2] != "zeta" {
			t.Errorf("expected sorted [alpha beta zeta], got %v", dirs)
		}
	})
}

func TestDeleteKey(t *testing.T) {
	t.Run("delete removes char at cursor", func(t *testing.T) {
		m := NewModel()
		m.Screen = ScreenProjectPath
		m.ProjectPathInput = "/tmp/foo"
		m.ProjectPathCursor = 5 // at 'f'

		result, _ := m.Update(tea.KeyMsg{Type: tea.KeyDelete})
		nm := result.(Model)

		if nm.ProjectPathInput != "/tmp/oo" {
			t.Errorf("expected '/tmp/oo', got %q", nm.ProjectPathInput)
		}
		if nm.ProjectPathCursor != 5 {
			t.Errorf("expected cursor at 5, got %d", nm.ProjectPathCursor)
		}
	})
}
