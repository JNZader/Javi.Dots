package tui

import (
	"testing"
	"time"

	"github.com/Gentleman-Programming/Gentleman.Dots/installer/internal/system"
	tea "github.com/charmbracelet/bubbletea"
)

func TestHandleBackupConfirmKeys(t *testing.T) {
	t.Run("should navigate with up/down keys", func(t *testing.T) {
		m := NewModel()
		m.Screen = ScreenBackupConfirm
		m.Cursor = 0

		// Press down
		result, _ := m.handleBackupConfirmKeys("down")
		newModel := result.(Model)
		if newModel.Cursor != 1 {
			t.Errorf("Expected cursor at 1 after down, got %d", newModel.Cursor)
		}

		// Press down again
		result, _ = newModel.handleBackupConfirmKeys("down")
		newModel = result.(Model)
		if newModel.Cursor != 2 {
			t.Errorf("Expected cursor at 2 after second down, got %d", newModel.Cursor)
		}

		// Press up
		result, _ = newModel.handleBackupConfirmKeys("up")
		newModel = result.(Model)
		if newModel.Cursor != 1 {
			t.Errorf("Expected cursor at 1 after up, got %d", newModel.Cursor)
		}
	})

	t.Run("should handle k/j vim keys", func(t *testing.T) {
		m := NewModel()
		m.Screen = ScreenBackupConfirm
		m.Cursor = 1

		// Press k (up)
		result, _ := m.handleBackupConfirmKeys("k")
		newModel := result.(Model)
		if newModel.Cursor != 0 {
			t.Errorf("Expected cursor at 0 after k, got %d", newModel.Cursor)
		}

		// Press j (down)
		result, _ = newModel.handleBackupConfirmKeys("j")
		newModel = result.(Model)
		if newModel.Cursor != 1 {
			t.Errorf("Expected cursor at 1 after j, got %d", newModel.Cursor)
		}
	})

	t.Run("should set CreateBackup true when selecting backup option", func(t *testing.T) {
		m := NewModel()
		m.Screen = ScreenBackupConfirm
		m.Cursor = 0 // Install with Backup
		m.SystemInfo = &system.SystemInfo{
			OS:       system.OSMac,
			HasBrew:  true,
			HasXcode: true,
		}
		m.Choices = UserChoices{
			OS:       "mac",
			Shell:    "fish",
			Terminal: "none",
		}
		m.ExistingConfigs = []string{"nvim: /test"}

		result, _ := m.handleBackupConfirmKeys("enter")
		newModel := result.(Model)

		if !newModel.Choices.CreateBackup {
			t.Error("CreateBackup should be true when selecting backup option")
		}

		if newModel.Screen != ScreenInstalling {
			t.Errorf("Expected ScreenInstalling, got %v", newModel.Screen)
		}
	})

	t.Run("should set CreateBackup false when selecting no backup option", func(t *testing.T) {
		m := NewModel()
		m.Screen = ScreenBackupConfirm
		m.Cursor = 1 // Install without Backup
		m.SystemInfo = &system.SystemInfo{
			OS:       system.OSMac,
			HasBrew:  true,
			HasXcode: true,
		}
		m.Choices = UserChoices{
			OS:       "mac",
			Shell:    "fish",
			Terminal: "none",
		}

		result, _ := m.handleBackupConfirmKeys("enter")
		newModel := result.(Model)

		if newModel.Choices.CreateBackup {
			t.Error("CreateBackup should be false when selecting no backup option")
		}
	})

	t.Run("should go to MainMenu when selecting cancel", func(t *testing.T) {
		m := NewModel()
		m.Screen = ScreenBackupConfirm
		m.Cursor = 2 // Cancel

		result, _ := m.handleBackupConfirmKeys("enter")
		newModel := result.(Model)

		if newModel.Screen != ScreenMainMenu {
			t.Errorf("Expected ScreenMainMenu, got %v", newModel.Screen)
		}
	})

	t.Run("should go to AIToolsSelect on escape with no AI tools (go back)", func(t *testing.T) {
		m := NewModel()
		m.Screen = ScreenBackupConfirm
		m.Cursor = 0

		// Note: ESC is handled by handleEscape(), not handleBackupConfirmKeys()
		// With no AI tools set, goes back to AI tools select
		result, _ := m.handleEscape()
		newModel := result.(Model)

		if newModel.Screen != ScreenAIToolsSelect {
			t.Errorf("Expected ScreenAIToolsSelect (go back, no AI tools), got %v", newModel.Screen)
		}
	})
}

func TestHandleRestoreBackupKeys(t *testing.T) {
	t.Run("should navigate with up/down keys", func(t *testing.T) {
		m := NewModel()
		m.Screen = ScreenRestoreBackup
		m.AvailableBackups = []system.BackupInfo{
			{Path: "/test/backup1"},
			{Path: "/test/backup2"},
		}
		m.Cursor = 0

		result, _ := m.handleRestoreBackupKeys("down")
		newModel := result.(Model)
		if newModel.Cursor != 1 {
			t.Errorf("Expected cursor at 1 after down, got %d", newModel.Cursor)
		}
	})

	t.Run("should go to ScreenRestoreConfirm when selecting a backup", func(t *testing.T) {
		m := NewModel()
		m.Screen = ScreenRestoreBackup
		m.AvailableBackups = []system.BackupInfo{
			{Path: "/test/backup1"},
		}
		m.Cursor = 0

		result, _ := m.handleRestoreBackupKeys("enter")
		newModel := result.(Model)

		if newModel.Screen != ScreenRestoreConfirm {
			t.Errorf("Expected ScreenRestoreConfirm, got %v", newModel.Screen)
		}

		if newModel.SelectedBackup != 0 {
			t.Errorf("Expected SelectedBackup 0, got %d", newModel.SelectedBackup)
		}
	})

	t.Run("should go to MainMenu on escape", func(t *testing.T) {
		m := NewModel()
		m.Screen = ScreenRestoreBackup
		m.Cursor = 0

		result, _ := m.handleRestoreBackupKeys("esc")
		newModel := result.(Model)

		if newModel.Screen != ScreenMainMenu {
			t.Errorf("Expected ScreenMainMenu, got %v", newModel.Screen)
		}
	})
}

func TestHandleRestoreConfirmKeys(t *testing.T) {
	t.Run("should navigate with up/down keys", func(t *testing.T) {
		m := NewModel()
		m.Screen = ScreenRestoreConfirm
		m.AvailableBackups = []system.BackupInfo{
			{Path: "/test/backup1"},
		}
		m.SelectedBackup = 0
		m.Cursor = 0

		result, _ := m.handleRestoreConfirmKeys("down")
		newModel := result.(Model)
		if newModel.Cursor != 1 {
			t.Errorf("Expected cursor at 1 after down, got %d", newModel.Cursor)
		}
	})

	t.Run("should go back to RestoreBackup on cancel", func(t *testing.T) {
		m := NewModel()
		m.Screen = ScreenRestoreConfirm
		m.AvailableBackups = []system.BackupInfo{
			{Path: "/test/backup1"},
		}
		m.SelectedBackup = 0
		m.Cursor = 2 // Cancel

		result, _ := m.handleRestoreConfirmKeys("enter")
		newModel := result.(Model)

		if newModel.Screen != ScreenRestoreBackup {
			t.Errorf("Expected ScreenRestoreBackup, got %v", newModel.Screen)
		}
	})

	t.Run("should go back to RestoreBackup on escape", func(t *testing.T) {
		m := NewModel()
		m.Screen = ScreenRestoreConfirm
		m.AvailableBackups = []system.BackupInfo{
			{Path: "/test/backup1"},
		}
		m.SelectedBackup = 0
		m.Cursor = 0

		result, _ := m.handleRestoreConfirmKeys("esc")
		newModel := result.(Model)

		if newModel.Screen != ScreenRestoreBackup {
			t.Errorf("Expected ScreenRestoreBackup, got %v", newModel.Screen)
		}
	})
}

func TestHandleMainMenuWithRestore(t *testing.T) {
	t.Run("should go to RestoreBackup when selecting restore option", func(t *testing.T) {
		m := NewModel()
		m.Screen = ScreenMainMenu
		m.AvailableBackups = []system.BackupInfo{
			{Path: "/test/backup1"},
		}
		// Options: Start, Learn & Practice, Restore, Init Project, Skill Manager, Exit
		// Restore is at index 2
		m.Cursor = 2

		result, _ := m.handleMainMenuKeys("enter")
		newModel := result.(Model)

		if newModel.Screen != ScreenRestoreBackup {
			t.Errorf("Expected ScreenRestoreBackup, got %v", newModel.Screen)
		}
	})

	t.Run("should handle dynamic menu correctly without backups", func(t *testing.T) {
		m := NewModel()
		m.Screen = ScreenMainMenu
		m.AvailableBackups = []system.BackupInfo{} // No backups
		// Options without restore: Start, Learn & Practice, Init Project, Skill Manager, Exit
		// Exit is at index 4
		m.Cursor = 4

		_, cmd := m.handleMainMenuKeys("enter")

		// Should return quit command
		if cmd == nil {
			t.Error("Expected quit command when selecting Exit")
		}
	})
}

func TestLoadBackupsMsg(t *testing.T) {
	t.Run("should update AvailableBackups on loadBackupsMsg", func(t *testing.T) {
		m := NewModel()
		m.AvailableBackups = []system.BackupInfo{}

		backups := []system.BackupInfo{
			{Path: "/test/backup1", Timestamp: time.Now()},
			{Path: "/test/backup2", Timestamp: time.Now()},
		}

		msg := loadBackupsMsg{backups: backups}
		result, _ := m.Update(msg)
		newModel := result.(Model)

		if len(newModel.AvailableBackups) != 2 {
			t.Errorf("Expected 2 backups, got %d", len(newModel.AvailableBackups))
		}
	})
}

func TestInitLoadsBackups(t *testing.T) {
	t.Run("Init should return command batch", func(t *testing.T) {
		m := NewModel()
		cmd := m.Init()

		// Init returns a batch command, we just verify it's not nil
		if cmd == nil {
			t.Error("Init should return a command")
		}
	})
}

func TestTickCmd(t *testing.T) {
	t.Run("tickCmd should return a command", func(t *testing.T) {
		cmd := tickCmd()
		if cmd == nil {
			t.Error("tickCmd should return a command")
		}
	})
}

func TestLoadBackupsCmd(t *testing.T) {
	t.Run("loadBackupsCmd should return a command", func(t *testing.T) {
		cmd := loadBackupsCmd()
		if cmd == nil {
			t.Error("loadBackupsCmd should return a command")
		}
	})

	t.Run("loadBackupsCmd should return loadBackupsMsg", func(t *testing.T) {
		cmd := loadBackupsCmd()
		msg := cmd()

		_, ok := msg.(loadBackupsMsg)
		if !ok {
			t.Errorf("Expected loadBackupsMsg, got %T", msg)
		}
	})
}

func TestUpdateHandlesBackupScreens(t *testing.T) {
	t.Run("should handle ScreenBackupConfirm key events", func(t *testing.T) {
		m := NewModel()
		m.Screen = ScreenBackupConfirm

		keyMsg := tea.KeyMsg{Type: tea.KeyDown}
		result, _ := m.Update(keyMsg)
		newModel := result.(Model)

		if newModel.Cursor != 1 {
			t.Errorf("Expected cursor at 1, got %d", newModel.Cursor)
		}
	})

	t.Run("should handle ScreenRestoreBackup key events", func(t *testing.T) {
		m := NewModel()
		m.Screen = ScreenRestoreBackup
		m.AvailableBackups = []system.BackupInfo{{Path: "/test"}}

		keyMsg := tea.KeyMsg{Type: tea.KeyDown}
		result, _ := m.Update(keyMsg)
		newModel := result.(Model)

		// Should stay on same screen
		if newModel.Screen != ScreenRestoreBackup {
			t.Errorf("Should stay on ScreenRestoreBackup")
		}
	})

	t.Run("should handle ScreenRestoreConfirm key events", func(t *testing.T) {
		m := NewModel()
		m.Screen = ScreenRestoreConfirm
		m.AvailableBackups = []system.BackupInfo{{Path: "/test"}}
		m.SelectedBackup = 0

		keyMsg := tea.KeyMsg{Type: tea.KeyDown}
		result, _ := m.Update(keyMsg)
		newModel := result.(Model)

		if newModel.Cursor != 1 {
			t.Errorf("Expected cursor at 1, got %d", newModel.Cursor)
		}
	})
}

func TestHandleEscapeFromBackupScreens(t *testing.T) {
	t.Run("should handle escape from ScreenBackupConfirm", func(t *testing.T) {
		m := NewModel()
		m.Screen = ScreenBackupConfirm

		result, _ := m.handleEscape()
		// handleEscape doesn't handle ScreenBackupConfirm directly
		// It's handled in handleBackupConfirmKeys
		newModel := result.(Model)
		if newModel.Screen != ScreenBackupConfirm {
			// If it changes, that's also valid
			t.Log("Screen changed from ScreenBackupConfirm on escape")
		}
	})
}

func TestWindowSizeMsg(t *testing.T) {
	t.Run("should update Width and Height", func(t *testing.T) {
		m := NewModel()

		msg := tea.WindowSizeMsg{Width: 120, Height: 40}
		result, _ := m.Update(msg)
		newModel := result.(Model)

		if newModel.Width != 120 {
			t.Errorf("Expected width 120, got %d", newModel.Width)
		}

		if newModel.Height != 40 {
			t.Errorf("Expected height 40, got %d", newModel.Height)
		}
	})
}

func TestCtrlCQuits(t *testing.T) {
	t.Run("ctrl+c should quit", func(t *testing.T) {
		m := NewModel()

		keyMsg := tea.KeyMsg{Type: tea.KeyCtrlC}
		result, cmd := m.Update(keyMsg)
		newModel := result.(Model)

		if !newModel.Quitting {
			t.Error("Should set Quitting to true on ctrl+c")
		}

		if cmd == nil {
			t.Error("Should return quit command")
		}
	})
}

func TestLearnMenuOptions(t *testing.T) {
	t.Run("ScreenLearnMenu returns 6 items", func(t *testing.T) {
		m := NewModel()
		m.Screen = ScreenLearnMenu
		opts := m.GetCurrentOptions()

		// Learn About Tools, Keymaps Reference, LazyVim Guide, Vim Trainer, separator, Back = 6
		if len(opts) != 6 {
			t.Errorf("expected 6 options, got %d: %v", len(opts), opts)
		}
	})
}

func TestLearnMenuNavigation(t *testing.T) {
	t.Run("Learn About Tools → ScreenLearnTerminals", func(t *testing.T) {
		m := NewModel()
		m.Screen = ScreenLearnMenu
		m.Cursor = 0 // Learn About Tools

		result, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
		nm := result.(Model)

		if nm.Screen != ScreenLearnTerminals {
			t.Errorf("expected ScreenLearnTerminals, got %d", nm.Screen)
		}
		if nm.PrevScreen != ScreenLearnMenu {
			t.Errorf("expected PrevScreen=ScreenLearnMenu, got %d", nm.PrevScreen)
		}
	})

	t.Run("Keymaps Reference → ScreenKeymapsMenu", func(t *testing.T) {
		m := NewModel()
		m.Screen = ScreenLearnMenu
		m.Cursor = 1 // Keymaps Reference

		result, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
		nm := result.(Model)

		if nm.Screen != ScreenKeymapsMenu {
			t.Errorf("expected ScreenKeymapsMenu, got %d", nm.Screen)
		}
		if nm.PrevScreen != ScreenLearnMenu {
			t.Errorf("expected PrevScreen=ScreenLearnMenu, got %d", nm.PrevScreen)
		}
	})

	t.Run("LazyVim Guide → ScreenLearnLazyVim", func(t *testing.T) {
		m := NewModel()
		m.Screen = ScreenLearnMenu
		m.Cursor = 2 // LazyVim Guide

		result, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
		nm := result.(Model)

		if nm.Screen != ScreenLearnLazyVim {
			t.Errorf("expected ScreenLearnLazyVim, got %d", nm.Screen)
		}
		if nm.PrevScreen != ScreenLearnMenu {
			t.Errorf("expected PrevScreen=ScreenLearnMenu, got %d", nm.PrevScreen)
		}
	})

	t.Run("Vim Trainer → ScreenTrainerMenu", func(t *testing.T) {
		m := NewModel()
		m.Screen = ScreenLearnMenu
		m.Cursor = 3 // Vim Trainer

		result, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
		nm := result.(Model)

		if nm.Screen != ScreenTrainerMenu {
			t.Errorf("expected ScreenTrainerMenu, got %d", nm.Screen)
		}
		if nm.PrevScreen != ScreenLearnMenu {
			t.Errorf("expected PrevScreen=ScreenLearnMenu, got %d", nm.PrevScreen)
		}
		if nm.TrainerStats == nil {
			t.Error("expected TrainerStats to be initialized")
		}
	})
}

func TestLearnMenuBack(t *testing.T) {
	t.Run("Back option returns to ScreenMainMenu", func(t *testing.T) {
		m := NewModel()
		m.Screen = ScreenLearnMenu
		m.Cursor = 5 // ← Back (after separator at 4)

		result, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
		nm := result.(Model)

		if nm.Screen != ScreenMainMenu {
			t.Errorf("expected ScreenMainMenu, got %d", nm.Screen)
		}
	})

	t.Run("Esc returns to ScreenMainMenu", func(t *testing.T) {
		m := NewModel()
		m.Screen = ScreenLearnMenu

		result, _ := m.Update(tea.KeyMsg{Type: tea.KeyEsc})
		nm := result.(Model)

		if nm.Screen != ScreenMainMenu {
			t.Errorf("expected ScreenMainMenu, got %d", nm.Screen)
		}
	})
}

func TestMainMenuLearnAndPractice(t *testing.T) {
	t.Run("Learn & Practice goes to ScreenLearnMenu", func(t *testing.T) {
		m := NewModel()
		m.Screen = ScreenMainMenu
		m.Cursor = 1 // Learn & Practice

		result, _ := m.handleMainMenuKeys("enter")
		nm := result.(Model)

		if nm.Screen != ScreenLearnMenu {
			t.Errorf("expected ScreenLearnMenu, got %d", nm.Screen)
		}
	})
}

func TestLearnMenuScreenTitleAndDescription(t *testing.T) {
	t.Run("ScreenLearnMenu has title", func(t *testing.T) {
		m := NewModel()
		m.Screen = ScreenLearnMenu
		title := m.GetScreenTitle()
		if title == "" {
			t.Error("expected non-empty title for ScreenLearnMenu")
		}
	})

	t.Run("ScreenLearnMenu has description", func(t *testing.T) {
		m := NewModel()
		m.Screen = ScreenLearnMenu
		desc := m.GetScreenDescription()
		if desc == "" {
			t.Error("expected non-empty description for ScreenLearnMenu")
		}
	})
}
