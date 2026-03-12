package tui

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/Gentleman-Programming/Gentleman.Dots/installer/internal/system"
)

// StepError provides context about which step failed and why
type StepError struct {
	StepID      string
	StepName    string
	Description string
	Cause       error
}

func (e *StepError) Error() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Step '%s' failed\n", e.StepName))
	sb.WriteString(fmt.Sprintf("Description: %s\n", e.Description))
	if e.Cause != nil {
		sb.WriteString(fmt.Sprintf("\nDetails:\n%v", e.Cause))
	}
	return sb.String()
}

func (e *StepError) Unwrap() error {
	return e.Cause
}

// wrapStepError creates a detailed error for a step failure
func wrapStepError(stepID, stepName, description string, cause error) error {
	return &StepError{
		StepID:      stepID,
		StepName:    stepName,
		Description: description,
		Cause:       cause,
	}
}

// executeStep runs the actual installation for a step
func executeStep(stepID string, m *Model) error {
	switch stepID {
	case "backup":
		return stepBackupConfigs(m)
	case "clone":
		return stepCloneRepo(m)
	case "homebrew":
		return stepInstallHomebrew(m)
	case "deps":
		return stepInstallDeps(m)
	case "xcode":
		return stepInstallXcode(m)
	case "terminal":
		return stepInstallTerminal(m)
	case "font":
		return stepInstallFont(m)
	case "shell":
		return stepInstallShell(m)
	case "wm":
		return stepInstallWM(m)
	case "nvim":
		return stepInstallNvim(m)
	case "zed":
		return stepInstallZed(m)
	case "aitools":
		return stepInstallAITools(m)
	case "aiframework":
		return stepInstallAIFramework(m)
	case "engram":
		return stepInstallEngram(m)
	case "cleanup":
		return stepCleanup(m)
	case "setshell":
		return stepSetDefaultShell(m)
	default:
		return fmt.Errorf("unknown step: %s", stepID)
	}
}

func stepBackupConfigs(m *Model) error {
	stepID := "backup"
	if len(m.ExistingConfigs) == 0 {
		SendLog(stepID, "No existing configs to backup")
		return nil
	}

	SendLog(stepID, fmt.Sprintf("Backing up %d existing configs...", len(m.ExistingConfigs)))

	// Extract just the config keys from the ExistingConfigs slice
	configKeys := make([]string, len(m.ExistingConfigs))
	for i, config := range m.ExistingConfigs {
		configKeys[i] = config
		SendLog(stepID, fmt.Sprintf("  → %s", config))
	}

	backupDir, err := system.CreateBackup(configKeys)
	if err != nil {
		return fmt.Errorf("failed to create backup: %w", err)
	}

	m.BackupDir = backupDir
	SendLog(stepID, fmt.Sprintf("✓ Backup created at: %s", backupDir))
	return nil
}

func stepCloneRepo(m *Model) error {
	stepID := "clone"
	repoDir := m.RepoDir

	// Check if already exists
	if _, err := os.Stat(repoDir); err == nil {
		SendLog(stepID, "Removing existing "+repoDir+" directory...")
		result := system.RunWithLogs("rm -rf "+repoDir, nil, func(line string) {
			SendLog(stepID, line)
		})
		if result.Error != nil {
			return wrapStepError("clone", "Clone Repository",
				"Failed to remove existing "+repoDir+" directory",
				result.Error)
		}
	}

	SendLog(stepID, "Cloning repository from GitHub...")
	result := system.RunWithLogs("git clone --progress "+m.RepoURL+" "+repoDir, nil, func(line string) {
		SendLog(stepID, line)
	})
	if result.Error != nil {
		return wrapStepError("clone", "Clone Repository",
			"Failed to clone the repository. Check your internet connection and git installation.",
			result.Error)
	}

	// Verify clone was successful
	if _, err := os.Stat(repoDir); os.IsNotExist(err) {
		return wrapStepError("clone", "Clone Repository",
			"Repository was cloned but directory not found",
			fmt.Errorf("%s directory does not exist after clone", repoDir))
	}

	SendLog(stepID, "✓ Repository cloned successfully")
	return nil
}

func stepInstallHomebrew(m *Model) error {
	stepID := "homebrew"

	// Termux doesn't use Homebrew - it uses pkg
	if m.SystemInfo.IsTermux {
		SendLog(stepID, "Skipping Homebrew (Termux uses pkg package manager)")
		return nil
	}

	if system.CommandExists("brew") {
		SendLog(stepID, "Homebrew already installed, skipping...")
		return nil
	}

	SendLog(stepID, "Installing Homebrew package manager...")
	result := system.RunWithLogs(`/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"`, nil, func(line string) {
		SendLog(stepID, line)
	})
	if result.Error != nil {
		return wrapStepError("homebrew", "Install Homebrew",
			"Failed to install Homebrew package manager. Check your internet connection.",
			result.Error)
	}

	// Add to PATH
	homeDir := os.Getenv("HOME")
	brewPrefix := system.GetBrewPrefix()

	shellConfig := fmt.Sprintf(`eval "$(%s/bin/brew shellenv)"`, brewPrefix)

	SendLog(stepID, "Configuring shell to use Homebrew...")
	// Add to common shell configs
	for _, rcFile := range []string{".bashrc", ".zshrc"} {
		rcPath := filepath.Join(homeDir, rcFile)
		if f, err := os.OpenFile(rcPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); err == nil {
			f.WriteString("\n" + shellConfig + "\n")
			f.Close()
		}
	}

	// Source it now
	system.Run(shellConfig, nil)

	SendLog(stepID, "✓ Homebrew installed successfully")
	return nil
}

func stepInstallDeps(m *Model) error {
	stepID := "deps"

	// Termux: use pkg (no sudo needed)
	// Check both SystemInfo and Choices.OS for redundancy
	isTermux := m.SystemInfo.IsTermux || m.Choices.OS == "termux"
	if isTermux {
		SendLog(stepID, "Updating Termux packages...")
		result := system.RunPkgWithLogs("update", nil, func(line string) {
			SendLog(stepID, line)
		})
		if result.Error != nil {
			return wrapStepError("deps", "Install Dependencies",
				"Failed to update Termux packages",
				result.Error)
		}
		result = system.RunPkgWithLogs("upgrade -y", nil, func(line string) {
			SendLog(stepID, line)
		})
		if result.Error != nil {
			// Upgrade failures are not critical
			SendLog(stepID, "Warning: package upgrade had issues, continuing...")
		}
		SendLog(stepID, "Installing base dependencies...")
		result = system.RunPkgInstall("git curl", nil, func(line string) {
			SendLog(stepID, line)
		})
		if result.Error != nil {
			return wrapStepError("deps", "Install Dependencies",
				"Failed to install base dependencies on Termux",
				result.Error)
		}
		return nil
	}

	// Arch Linux
	if m.SystemInfo.OS == system.OSArch {
		result := system.RunSudo("pacman -Syu --noconfirm", nil)
		if result.Error != nil {
			return wrapStepError("deps", "Install Dependencies",
				"Failed to update Arch Linux packages",
				result.Error)
		}
		result = system.RunSudo("pacman -S --needed --noconfirm base-devel curl file git wget unzip fontconfig", nil)
		if result.Error != nil {
			return wrapStepError("deps", "Install Dependencies",
				"Failed to install base dependencies on Arch Linux",
				result.Error)
		}
		return nil
	}

	// Fedora/RHEL
	if m.SystemInfo.OS == system.OSFedora {
		result := system.RunSudo("dnf check-update || true", nil) // dnf check-update returns 100 if updates available
		result = system.RunSudo("dnf install -y @development-tools curl file git wget unzip fontconfig", nil)
		if result.Error != nil {
			return wrapStepError("deps", "Install Dependencies",
				"Failed to install base dependencies on Fedora/RHEL",
				result.Error)
		}
		return nil
	}

	// Debian/Ubuntu
	result := system.RunSudo("apt-get update", nil)
	if result.Error != nil {
		return wrapStepError("deps", "Install Dependencies",
			"Failed to update apt package list",
			result.Error)
	}
	result = system.RunSudo("apt-get install -y build-essential curl file git unzip fontconfig", nil)
	if result.Error != nil {
		return wrapStepError("deps", "Install Dependencies",
			"Failed to install base dependencies on Debian/Ubuntu",
			result.Error)
	}
	return nil
}

func stepInstallXcode(m *Model) error {
	result := system.Run("xcode-select --install", nil)
	if result.Error != nil {
		// xcode-select returns error if already installed, which is fine
		if result.ExitCode == 1 && strings.Contains(result.Stderr, "already installed") {
			return nil
		}
		return wrapStepError("xcode", "Install Xcode CLI",
			"Failed to install Xcode Command Line Tools. You may need to install them manually from the App Store.",
			result.Error)
	}
	return nil
}

func stepInstallTerminal(m *Model) error {
	terminal := m.Choices.Terminal
	homeDir := os.Getenv("HOME")
	repoDir := m.RepoDir
	stepID := "terminal"

	switch terminal {
	case "alacritty":
		if !system.CommandExists("alacritty") {
			SendLog(stepID, "Installing Alacritty...")
			var result *system.ExecResult
			if m.SystemInfo.OS == system.OSArch {
				result = system.RunSudoWithLogs("pacman -S --noconfirm alacritty", nil, func(line string) {
					SendLog(stepID, line)
				})
			} else if m.SystemInfo.OS == system.OSMac {
				result = system.RunBrewWithLogs("install --cask alacritty", nil, func(line string) {
					SendLog(stepID, line)
				})
			} else if m.SystemInfo.OS == system.OSFedora {
				// Fedora: install from dnf
				result = system.RunSudoWithLogs("dnf install -y alacritty", nil, func(line string) {
					SendLog(stepID, line)
				})
			} else if m.SystemInfo.OS == system.OSDebian || m.SystemInfo.OS == system.OSLinux {
				// Debian/Ubuntu: compile from source (PPAs are unreliable)
				SendLog(stepID, "Building Alacritty from source...")
				SendLog(stepID, "Installing build dependencies...")
				result = system.RunSudoWithLogs("apt-get install -y cmake pkg-config libfreetype6-dev libfontconfig1-dev libxcb-xfixes0-dev libxkbcommon-dev python3 gzip scdoc git curl", nil, func(line string) {
					SendLog(stepID, line)
				})
				if result.Error != nil {
					return wrapStepError("terminal", "Install Alacritty",
						"Failed to install build dependencies",
						result.Error)
				}
				// Install Rust/Cargo only for this build
				cargoPath := filepath.Join(homeDir, ".cargo/bin/cargo")
				if !system.CommandExists("cargo") && !system.CommandExists(cargoPath) {
					SendLog(stepID, "Installing Rust/Cargo toolchain...")
					result = system.RunWithLogs("curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh -s -- -y", nil, func(line string) {
						SendLog(stepID, line)
					})
					if result.Error != nil {
						return wrapStepError("terminal", "Install Alacritty",
							"Failed to install Rust",
							result.Error)
					}
					cargoPath = filepath.Join(homeDir, ".cargo/bin/cargo")
				}
				// Clone and build Alacritty
				SendLog(stepID, "Cloning Alacritty repository...")
				alacrittyDir := filepath.Join(os.TempDir(), "alacritty-build")
				os.RemoveAll(alacrittyDir)
				result = system.RunWithLogs(fmt.Sprintf("git clone https://github.com/alacritty/alacritty.git %s", alacrittyDir), nil, func(line string) {
					SendLog(stepID, line)
				})
				if result.Error != nil {
					return wrapStepError("terminal", "Install Alacritty",
						"Failed to clone Alacritty repository",
						result.Error)
				}
				SendLog(stepID, "Building Alacritty (this may take 5-10 minutes)...")
				if !system.CommandExists("cargo") {
					cargoPath = filepath.Join(homeDir, ".cargo/bin/cargo")
				} else {
					cargoPath = "cargo"
				}
				result = system.RunWithLogs(fmt.Sprintf("%s build --release --manifest-path %s/Cargo.toml", cargoPath, alacrittyDir), nil, func(line string) {
					SendLog(stepID, line)
				})
				if result.Error != nil {
					return wrapStepError("terminal", "Install Alacritty",
						"Failed to build Alacritty",
						result.Error)
				}
				SendLog(stepID, "Installing Alacritty binary...")
				result = system.RunSudoWithLogs(fmt.Sprintf("cp %s/target/release/alacritty /usr/local/bin/alacritty", alacrittyDir), nil, func(line string) {
					SendLog(stepID, line)
				})
				if result.Error != nil {
					return wrapStepError("terminal", "Install Alacritty",
						"Failed to install Alacritty binary",
						result.Error)
				}
				system.RunSudoWithLogs(fmt.Sprintf("cp %s/extra/linux/Alacritty.desktop /usr/share/applications/", alacrittyDir), nil, func(line string) {
					SendLog(stepID, line)
				})
				os.RemoveAll(alacrittyDir)
				SendLog(stepID, "✓ Alacritty built and installed from source")
			} else {
				return wrapStepError("terminal", "Install Alacritty",
					"Unsupported operating system for Alacritty installation",
					fmt.Errorf("OS type: %v", m.SystemInfo.OS))
			}
			if result.Error != nil {
				return wrapStepError("terminal", "Install Alacritty",
					"Failed to install Alacritty terminal emulator",
					result.Error)
			}
		} else {
			SendLog(stepID, "Alacritty already installed")
		}
		SendLog(stepID, "Copying Alacritty configuration...")
		if err := system.EnsureDir(filepath.Join(homeDir, ".config/alacritty")); err != nil {
			return wrapStepError("terminal", "Install Alacritty",
				"Failed to create Alacritty config directory",
				err)
		}
		if err := system.CopyFile(filepath.Join(repoDir, "alacritty.toml"), filepath.Join(homeDir, ".config/alacritty/alacritty.toml")); err != nil {
			return wrapStepError("terminal", "Install Alacritty",
				"Failed to copy Alacritty configuration",
				err)
		}
		SendLog(stepID, "✓ Alacritty configured")

	case "wezterm":
		if !system.CommandExists("wezterm") {
			SendLog(stepID, "Installing WezTerm...")
			var result *system.ExecResult
			if m.SystemInfo.OS == system.OSArch {
				result = system.RunSudoWithLogs("pacman -S --noconfirm wezterm", nil, func(line string) {
					SendLog(stepID, line)
				})
			} else if m.SystemInfo.OS == system.OSFedora {
				// Fedora: enable COPR and install
				system.RunSudo("dnf copr enable -y wezfurlong/wezterm-nightly", nil)
				result = system.RunSudoWithLogs("dnf install -y wezterm", nil, func(line string) {
					SendLog(stepID, line)
				})
			} else if m.SystemInfo.OS == system.OSMac {
				result = system.RunBrewWithLogs("install --cask wezterm", nil, func(line string) {
					SendLog(stepID, line)
				})
			} else {
				system.Run("brew tap wez/wezterm-linuxbrew", nil)
				result = system.RunBrewWithLogs("install wezterm", nil, func(line string) {
					SendLog(stepID, line)
				})
			}
			if result.Error != nil {
				return wrapStepError("terminal", "Install WezTerm",
					"Failed to install WezTerm terminal emulator",
					result.Error)
			}
		} else {
			SendLog(stepID, "WezTerm already installed")
		}
		SendLog(stepID, "Copying WezTerm configuration...")
		if err := system.EnsureDir(filepath.Join(homeDir, ".config/wezterm")); err != nil {
			return wrapStepError("terminal", "Install WezTerm",
				"Failed to create WezTerm config directory",
				err)
		}
		if err := system.CopyFile(filepath.Join(repoDir, ".wezterm.lua"), filepath.Join(homeDir, ".config/wezterm/wezterm.lua")); err != nil {
			return wrapStepError("terminal", "Install WezTerm",
				"Failed to copy WezTerm configuration",
				err)
		}
		SendLog(stepID, "✓ WezTerm configured")

	case "kitty":
		if !system.CommandExists("kitty") && m.SystemInfo.OS == system.OSMac {
			SendLog(stepID, "Installing Kitty...")
			result := system.RunBrewWithLogs("install --cask kitty", nil, func(line string) {
				SendLog(stepID, line)
			})
			if result.Error != nil {
				return wrapStepError("terminal", "Install Kitty",
					"Failed to install Kitty terminal emulator",
					result.Error)
			}
		} else {
			SendLog(stepID, "Kitty already installed")
		}
		SendLog(stepID, "Copying Kitty configuration...")
		if err := system.EnsureDir(filepath.Join(homeDir, ".config/kitty")); err != nil {
			return wrapStepError("terminal", "Install Kitty",
				"Failed to create Kitty config directory",
				err)
		}
		if err := system.CopyDir(filepath.Join(repoDir, "GentlemanKitty"), filepath.Join(homeDir, ".config", "kitty")); err != nil {
			return wrapStepError("terminal", "Install Kitty",
				"Failed to copy Kitty configuration",
				err)
		}
		SendLog(stepID, "✓ Kitty configured")

	case "ghostty":
		if !system.CommandExists("ghostty") {
			SendLog(stepID, "Installing Ghostty...")
			var result *system.ExecResult
			if m.SystemInfo.OS == system.OSArch {
				result = system.RunSudoWithLogs("pacman -S --noconfirm ghostty", nil, func(line string) {
					SendLog(stepID, line)
				})
			} else if m.SystemInfo.OS == system.OSFedora {
				// Fedora: enable COPR and install
				system.RunSudo("dnf copr enable -y pgdev/ghostty", nil)
				result = system.RunSudoWithLogs("dnf install -y ghostty", nil, func(line string) {
					SendLog(stepID, line)
				})
			} else if m.SystemInfo.OS == system.OSMac {
				result = system.RunBrewWithLogs("install --cask ghostty", nil, func(line string) {
					SendLog(stepID, line)
				})
			} else {
				result = system.RunWithLogs(`/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/mkasberg/ghostty-ubuntu/HEAD/install.sh)"`, nil, func(line string) {
					SendLog(stepID, line)
				})
			}
			if result.Error != nil {
				return wrapStepError("terminal", "Install Ghostty",
					"Failed to install Ghostty terminal emulator",
					result.Error)
			}
		} else {
			SendLog(stepID, "Ghostty already installed")
		}
		SendLog(stepID, "Copying Ghostty configuration...")
		if err := system.EnsureDir(filepath.Join(homeDir, ".config/ghostty")); err != nil {
			return wrapStepError("terminal", "Install Ghostty",
				"Failed to create Ghostty config directory",
				err)
		}
		if err := system.CopyDir(filepath.Join(repoDir, "GentlemanGhostty"), filepath.Join(homeDir, ".config", "ghostty")); err != nil {
			return wrapStepError("terminal", "Install Ghostty",
				"Failed to copy Ghostty configuration",
				err)
		}
		SendLog(stepID, "✓ Ghostty configured")
	}

	return nil
}

func stepInstallFont(m *Model) error {
	homeDir := os.Getenv("HOME")
	stepID := "font"

	// Termux: fonts work differently - copy to ~/.termux/font.ttf
	isTermux := m.SystemInfo.IsTermux || m.Choices.OS == "termux"
	if isTermux {
		SendLog(stepID, "Downloading JetBrainsMono Nerd Font for Termux...")
		termuxDir := filepath.Join(homeDir, ".termux")
		if err := system.EnsureDir(termuxDir); err != nil {
			return wrapStepError("font", "Install Nerd Font",
				"Failed to create .termux directory",
				err)
		}

		// Download a single TTF file for Termux
		result := system.RunWithLogs(fmt.Sprintf("curl -fsSL -o %s/font.ttf https://github.com/ryanoasis/nerd-fonts/raw/HEAD/patched-fonts/JetBrainsMono/Ligatures/Regular/JetBrainsMonoNerdFont-Regular.ttf", termuxDir), nil, func(line string) {
			SendLog(stepID, line)
		})
		if result.Error != nil {
			return wrapStepError("font", "Install Nerd Font",
				"Failed to download font. Check your internet connection.",
				result.Error)
		}

		SendLog(stepID, "Reloading Termux settings...")
		system.Run("termux-reload-settings", nil)
		SendLog(stepID, "✓ Font installed - restart Termux to apply")
		return nil
	}

	if m.SystemInfo.OS == system.OSMac {
		SendLog(stepID, "Installing Iosevka Term Nerd Font...")
		result := system.RunBrewWithLogs("install --cask font-iosevka-term-nerd-font", nil, func(line string) {
			SendLog(stepID, line)
		})
		if result.Error != nil {
			return wrapStepError("font", "Install Iosevka Nerd Font",
				"Failed to install font via Homebrew. Try installing manually from https://www.nerdfonts.com/",
				result.Error)
		}
		SendLog(stepID, "✓ Font installed")
		return nil
	}

	// Linux
	fontDir := filepath.Join(homeDir, ".local/share/fonts")
	SendLog(stepID, "Creating fonts directory...")
	if err := system.EnsureDir(fontDir); err != nil {
		return wrapStepError("font", "Install Iosevka Nerd Font",
			"Failed to create fonts directory",
			err)
	}

	SendLog(stepID, "Downloading Iosevka Term Nerd Font...")
	result := system.RunWithLogs(fmt.Sprintf("curl -fsSL -o %s/IosevkaTerm.zip https://github.com/ryanoasis/nerd-fonts/releases/download/v3.3.0/IosevkaTerm.zip", fontDir), nil, func(line string) {
		SendLog(stepID, line)
	})
	if result.Error != nil {
		return wrapStepError("font", "Install Iosevka Nerd Font",
			"Failed to download font. Check your internet connection.",
			result.Error)
	}

	SendLog(stepID, "Extracting font archive...")
	result = system.RunWithLogs(fmt.Sprintf("unzip -o %s/IosevkaTerm.zip -d %s/", fontDir, fontDir), nil, func(line string) {
		SendLog(stepID, line)
	})
	if result.Error != nil {
		return wrapStepError("font", "Install Iosevka Nerd Font",
			"Failed to extract font archive",
			result.Error)
	}

	SendLog(stepID, "Updating font cache...")
	system.RunWithLogs("fc-cache -fv", nil, func(line string) {
		SendLog(stepID, line)
	})
	SendLog(stepID, "✓ Font installed")
	return nil
}

func stepInstallShell(m *Model) error {
	homeDir := os.Getenv("HOME")
	repoDir := m.RepoDir
	shell := m.Choices.Shell
	stepID := "shell"

	// Common dependencies
	SendLog(stepID, "Creating required directories...")
	system.EnsureDir(filepath.Join(homeDir, ".config"))
	system.EnsureDir(filepath.Join(homeDir, ".cache/starship"))
	system.EnsureDir(filepath.Join(homeDir, ".cache/carapace"))
	system.EnsureDir(filepath.Join(homeDir, ".local/share/atuin"))

	switch shell {
	case "fish":
		SendLog(stepID, "Installing Fish shell and plugins...")
		var result *system.ExecResult
		if m.SystemInfo.IsTermux {
			result = system.RunPkgInstall("fish starship zoxide", nil, func(line string) {
				SendLog(stepID, line)
			})
		} else {
			result = system.RunBrewWithLogs("install fish carapace zoxide atuin starship", nil, func(line string) {
				SendLog(stepID, line)
			})
		}
		if result.Error != nil {
			return wrapStepError("shell", "Install Fish",
				"Failed to install Fish shell and dependencies",
				result.Error)
		}
		SendLog(stepID, "Copying Fish configuration...")
		if err := system.CopyFile(filepath.Join(repoDir, "starship.toml"), filepath.Join(homeDir, ".config/starship.toml")); err != nil {
			return wrapStepError("shell", "Install Fish",
				"Failed to copy starship configuration",
				err)
		}
		if err := system.CopyDir(filepath.Join(repoDir, "GentlemanFish", "fish"), filepath.Join(homeDir, ".config", "fish")); err != nil {
			return wrapStepError("shell", "Install Fish",
				"Failed to copy Fish configuration",
				err)
		}
		// Patch config.fish based on WM choice
		SendLog(stepID, "Configuring shell for window manager...")
		if err := system.PatchFishForWM(filepath.Join(homeDir, ".config/fish/config.fish"), m.Choices.WindowMgr, m.Choices.InstallNvim); err != nil {
			return wrapStepError("shell", "Install Fish",
				"Failed to configure config.fish for window manager",
				err)
		}
		// Remove tmux.fish function if not using tmux
		if m.Choices.WindowMgr != "tmux" {
			os.Remove(filepath.Join(homeDir, ".config/fish/functions/tmux.fish"))
		}
		// Termux: Add fish to $PREFIX/etc/shells so tmux doesn't complain
		if m.SystemInfo.IsTermux {
			SendLog(stepID, "Adding fish to Termux shells...")
			prefix := os.Getenv("PREFIX")
			if prefix == "" {
				prefix = "/data/data/com.termux/files/usr"
			}
			shellsFile := filepath.Join(prefix, "etc", "shells")
			system.EnsureDir(filepath.Join(prefix, "etc"))
			f, err := os.OpenFile(shellsFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err == nil {
				f.WriteString(filepath.Join(prefix, "bin", "fish") + "\n")
				f.Close()
			}
		}
		SendLog(stepID, "✓ Fish shell configured")

	case "zsh":
		SendLog(stepID, "Installing Zsh and plugins...")
		var result *system.ExecResult
		if m.SystemInfo.IsTermux {
			// Termux has zsh in pkg, but plugins need to be installed differently
			result = system.RunPkgInstall("zsh starship zoxide", nil, func(line string) {
				SendLog(stepID, line)
			})
		} else {
			result = system.RunBrewWithLogs("install zsh carapace zoxide atuin zsh-autosuggestions zsh-syntax-highlighting zsh-autocomplete powerlevel10k", nil, func(line string) {
				SendLog(stepID, line)
			})
		}
		if result.Error != nil {
			return wrapStepError("shell", "Install Zsh",
				"Failed to install Zsh and plugins",
				result.Error)
		}
		SendLog(stepID, "Copying Zsh configuration...")
		if err := system.CopyFile(filepath.Join(repoDir, "GentlemanZsh/.zshrc"), filepath.Join(homeDir, ".zshrc")); err != nil {
			return wrapStepError("shell", "Install Zsh",
				"Failed to copy .zshrc configuration",
				err)
		}
		// Patch .zshrc based on WM choice
		SendLog(stepID, "Configuring shell for window manager...")
		if err := system.PatchZshForWM(filepath.Join(homeDir, ".zshrc"), m.Choices.WindowMgr, m.Choices.InstallNvim); err != nil {
			return wrapStepError("shell", "Install Zsh",
				"Failed to configure .zshrc for window manager",
				err)
		}
		if err := system.CopyFile(filepath.Join(repoDir, "GentlemanZsh/.p10k.zsh"), filepath.Join(homeDir, ".p10k.zsh")); err != nil {
			return wrapStepError("shell", "Install Zsh",
				"Failed to copy Powerlevel10k configuration",
				err)
		}
		if err := system.CopyDir(filepath.Join(repoDir, "GentlemanZsh", ".oh-my-zsh"), filepath.Join(homeDir, ".oh-my-zsh")); err != nil {
			return wrapStepError("shell", "Install Zsh",
				"Failed to copy Oh-My-Zsh directory",
				err)
		}
		// Termux: Add zsh to $PREFIX/etc/shells so tmux doesn't complain
		if m.SystemInfo.IsTermux {
			SendLog(stepID, "Adding zsh to Termux shells...")
			prefix := os.Getenv("PREFIX")
			if prefix == "" {
				prefix = "/data/data/com.termux/files/usr"
			}
			shellsFile := filepath.Join(prefix, "etc", "shells")
			system.EnsureDir(filepath.Join(prefix, "etc"))
			f, err := os.OpenFile(shellsFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err == nil {
				f.WriteString(filepath.Join(prefix, "bin", "zsh") + "\n")
				f.Close()
			}
		}
		SendLog(stepID, "✓ Zsh configured with Powerlevel10k")

	case "nushell":
		SendLog(stepID, "Installing Nushell and dependencies...")
		var result *system.ExecResult
		if m.SystemInfo.IsTermux {
			result = system.RunPkgInstall("nushell starship zoxide jq", nil, func(line string) {
				SendLog(stepID, line)
			})
		} else {
			result = system.RunBrewWithLogs("install nushell carapace zoxide atuin jq bash starship", nil, func(line string) {
				SendLog(stepID, line)
			})
		}
		if result.Error != nil {
			return wrapStepError("shell", "Install Nushell",
				"Failed to install Nushell and dependencies",
				result.Error)
		}
		SendLog(stepID, "Copying Nushell configuration...")
		if err := system.CopyFile(filepath.Join(repoDir, "starship.toml"), filepath.Join(homeDir, ".config/starship.toml")); err != nil {
			return wrapStepError("shell", "Install Nushell",
				"Failed to copy starship configuration",
				err)
		}
		if err := system.CopyFile(filepath.Join(repoDir, "bash-env-json"), filepath.Join(homeDir, ".config/bash-env-json")); err != nil {
			return wrapStepError("shell", "Install Nushell",
				"Failed to copy bash-env-json",
				err)
		}
		if err := system.CopyFile(filepath.Join(repoDir, "bash-env.nu"), filepath.Join(homeDir, ".config/bash-env.nu")); err != nil {
			return wrapStepError("shell", "Install Nushell",
				"Failed to copy bash-env.nu",
				err)
		}

		var nuDir string
		if runtime.GOOS == "darwin" {
			nuDir = filepath.Join(homeDir, "Library/Application Support/nushell")
		} else {
			nuDir = filepath.Join(homeDir, ".config/nushell")
		}
		if err := system.EnsureDir(nuDir); err != nil {
			return wrapStepError("shell", "Install Nushell",
				"Failed to create Nushell config directory",
				err)
		}
		if err := system.CopyDir(filepath.Join(repoDir, "GentlemanNushell"), nuDir); err != nil {
			return wrapStepError("shell", "Install Nushell",
				"Failed to copy Nushell configuration",
				err)
		}
		// Patch config.nu based on WM choice
		SendLog(stepID, "Configuring shell for window manager...")
		if err := system.PatchNushellForWM(filepath.Join(nuDir, "config.nu"), m.Choices.WindowMgr); err != nil {
			return wrapStepError("shell", "Install Nushell",
				"Failed to configure config.nu for window manager",
				err)
		}
		// Termux: Add nu to $PREFIX/etc/shells so tmux doesn't complain
		if m.SystemInfo.IsTermux {
			SendLog(stepID, "Adding nushell to Termux shells...")
			prefix := os.Getenv("PREFIX")
			if prefix == "" {
				prefix = "/data/data/com.termux/files/usr"
			}
			shellsFile := filepath.Join(prefix, "etc", "shells")
			system.EnsureDir(filepath.Join(prefix, "etc"))
			f, err := os.OpenFile(shellsFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err == nil {
				f.WriteString(filepath.Join(prefix, "bin", "nu") + "\n")
				f.Close()
			}
		}
		SendLog(stepID, "✓ Nushell configured")
	}

	return nil
}

func stepInstallWM(m *Model) error {
	homeDir := os.Getenv("HOME")
	repoDir := m.RepoDir
	wm := m.Choices.WindowMgr
	stepID := "wm"

	switch wm {
	case "tmux":
		if !system.CommandExists("tmux") {
			SendLog(stepID, "Installing Tmux...")
			var result *system.ExecResult
			if m.SystemInfo.IsTermux {
				result = system.RunPkgInstall("tmux", nil, func(line string) {
					SendLog(stepID, line)
				})
			} else {
				result = system.RunBrewWithLogs("install tmux", nil, func(line string) {
					SendLog(stepID, line)
				})
			}
			if result.Error != nil {
				return wrapStepError("wm", "Install Tmux",
					"Failed to install Tmux",
					result.Error)
			}
		} else {
			SendLog(stepID, "Tmux already installed")
		}

		// TPM
		tpmDir := filepath.Join(homeDir, ".tmux/plugins/tpm")
		if _, err := os.Stat(tpmDir); os.IsNotExist(err) {
			SendLog(stepID, "Cloning TPM (Tmux Plugin Manager)...")
			result := system.RunWithLogs(fmt.Sprintf("git clone https://github.com/tmux-plugins/tpm %s", tpmDir), nil, func(line string) {
				SendLog(stepID, line)
			})
			if result.Error != nil {
				return wrapStepError("wm", "Install Tmux",
					"Failed to clone TPM (Tmux Plugin Manager)",
					result.Error)
			}
		}

		SendLog(stepID, "Copying Tmux configuration...")
		if err := system.EnsureDir(filepath.Join(homeDir, ".tmux")); err != nil {
			return wrapStepError("wm", "Install Tmux",
				"Failed to create .tmux directory",
				err)
		}
		if err := system.CopyDir(filepath.Join(repoDir, "GentlemanTmux", "plugins"), filepath.Join(homeDir, ".tmux", "plugins")); err != nil {
			return wrapStepError("wm", "Install Tmux",
				"Failed to copy Tmux plugins",
				err)
		}
		if err := system.CopyFile(filepath.Join(repoDir, "GentlemanTmux/tmux.conf"), filepath.Join(homeDir, ".tmux.conf")); err != nil {
			return wrapStepError("wm", "Install Tmux",
				"Failed to copy tmux.conf",
				err)
		}

		// Configure tmux to use the user's chosen shell
		SendLog(stepID, "Configuring tmux default shell...")
		tmuxConfPath := filepath.Join(homeDir, ".tmux.conf")
		shellName := ""
		switch m.Choices.Shell {
		case "fish":
			shellName = "fish"
		case "zsh":
			shellName = "zsh"
		case "nushell":
			shellName = "nu"
		}
		if shellName != "" {
			// Find the full path to the shell
			shellFullPath := ""
			if m.SystemInfo.IsTermux {
				// In Termux, construct the path directly (which command has issues)
				prefix := os.Getenv("PREFIX")
				if prefix == "" {
					prefix = "/data/data/com.termux/files/usr"
				}
				shellFullPath = filepath.Join(prefix, "bin", shellName)
			} else {
				result := system.Run(fmt.Sprintf("which %s", shellName), nil)
				if result.Error == nil && result.Output != "" {
					shellFullPath = strings.TrimSpace(result.Output)
				}
			}
			if shellFullPath == "" {
				shellFullPath = shellName // Fallback
			}

			// Replace placeholder in tmux.conf with actual shell config
			content, err := os.ReadFile(tmuxConfPath)
			if err == nil {
				shellConfig := fmt.Sprintf("set -g default-command \"%s\"\nset -g default-shell \"%s\"", shellFullPath, shellFullPath)
				newContent := strings.Replace(string(content), "# GENTLEMAN_DEFAULT_SHELL", shellConfig, 1)
				os.WriteFile(tmuxConfPath, []byte(newContent), 0644)
			}
		}

		// Install plugins
		SendLog(stepID, "Installing Tmux plugins...")
		system.RunWithLogs(filepath.Join(homeDir, ".tmux/plugins/tpm/bin/install_plugins"), nil, func(line string) {
			SendLog(stepID, line)
		})
		SendLog(stepID, "✓ Tmux configured")

	case "zellij":
		if !system.CommandExists("zellij") {
			SendLog(stepID, "Installing Zellij...")
			var result *system.ExecResult
			if m.SystemInfo.IsTermux {
				result = system.RunPkgInstall("zellij", nil, func(line string) {
					SendLog(stepID, line)
				})
			} else {
				result = system.RunBrewWithLogs("install zellij", nil, func(line string) {
					SendLog(stepID, line)
				})
			}
			if result.Error != nil {
				return wrapStepError("wm", "Install Zellij",
					"Failed to install Zellij",
					result.Error)
			}
		} else {
			SendLog(stepID, "Zellij already installed")
		}

		SendLog(stepID, "Copying Zellij configuration...")
		zellijDir := filepath.Join(homeDir, ".config/zellij")
		if err := system.EnsureDir(zellijDir); err != nil {
			return wrapStepError("wm", "Install Zellij",
				"Failed to create Zellij config directory",
				err)
		}
		if err := system.CopyDir(filepath.Join(repoDir, "GentlemanZellij", "zellij"), zellijDir); err != nil {
			return wrapStepError("wm", "Install Zellij",
				"Failed to copy Zellij configuration",
				err)
		}

		// Configure zellij to use the user's chosen shell
		SendLog(stepID, "Configuring zellij default shell...")
		zellijConfPath := filepath.Join(zellijDir, "config.kdl")
		shellPath := ""
		switch m.Choices.Shell {
		case "fish":
			shellPath = "fish"
		case "zsh":
			shellPath = "zsh"
		case "nushell":
			shellPath = "nu"
		}
		if shellPath != "" {
			// Append default_shell config to zellij config.kdl
			f, err := os.OpenFile(zellijConfPath, os.O_APPEND|os.O_WRONLY, 0644)
			if err == nil {
				f.WriteString(fmt.Sprintf("\n// Default shell (configured by Gentleman.Dots)\ndefault_shell \"%s\"\n", shellPath))
				f.Close()
			}
		}
		SendLog(stepID, "✓ Zellij configured")
	}

	return nil
}

func stepInstallNvim(m *Model) error {
	homeDir := os.Getenv("HOME")
	repoDir := m.RepoDir
	stepID := "nvim"

	// Obsidian app installation (if user opted in)
	if m.Choices.InstallObsidian {
		SendLog(stepID, "Installing Obsidian app...")
		var obsResult *system.ExecResult
		switch m.SystemInfo.OS {
		case system.OSMac:
			obsResult = system.RunBrewWithLogs("install --cask obsidian", nil, func(line string) {
				SendLog(stepID, line)
			})
		case system.OSArch:
			obsResult = system.RunSudoWithLogs("pacman -S --noconfirm obsidian", nil, func(line string) {
				SendLog(stepID, line)
			})
		case system.OSDebian, system.OSLinux:
			obsResult = system.RunWithLogs("flatpak install -y flathub md.obsidian.Obsidian", nil, func(line string) {
				SendLog(stepID, line)
			})
		case system.OSFedora:
			obsResult = system.RunWithLogs("flatpak install -y flathub md.obsidian.Obsidian", nil, func(line string) {
				SendLog(stepID, line)
			})
		}
		if obsResult != nil && obsResult.Error != nil {
			SendLog(stepID, "⚠ Obsidian install failed: "+obsResult.Error.Error())
			SendLog(stepID, "You can install Obsidian manually from https://obsidian.md")
		} else {
			SendLog(stepID, "✓ Obsidian installed")
		}
	}

	// Obsidian directories for Neovim plugin
	SendLog(stepID, "Creating Obsidian directories...")
	obsidianDir := filepath.Join(homeDir, ".config/obsidian")
	system.EnsureDir(obsidianDir)
	system.EnsureDir(filepath.Join(obsidianDir, "templates"))

	// Check Node.js
	if !system.CommandExists("node") {
		SendLog(stepID, "Installing Node.js...")
		var result *system.ExecResult
		if m.SystemInfo.IsTermux {
			result = system.RunPkgInstall("nodejs", nil, func(line string) {
				SendLog(stepID, line)
			})
		} else {
			result = system.RunBrewWithLogs("install node", nil, func(line string) {
				SendLog(stepID, line)
			})
		}
		if result.Error != nil {
			return wrapStepError("nvim", "Install Neovim",
				"Failed to install Node.js (required for LSP servers)",
				result.Error)
		}
	} else {
		SendLog(stepID, "Node.js already installed")
	}

	// Install dependencies
	SendLog(stepID, "Installing Neovim and dependencies...")
	var result *system.ExecResult
	if m.SystemInfo.IsTermux {
		// Termux package names (neovim instead of nvim, clang instead of gcc)
		result = system.RunPkgInstall("neovim git clang fzf fd ripgrep bat curl lazygit", nil, func(line string) {
			SendLog(stepID, line)
		})
	} else {
		result = system.RunBrewWithLogs("install nvim git gcc fzf fd ripgrep coreutils bat curl lazygit tree-sitter", nil, func(line string) {
			SendLog(stepID, line)
		})
	}
	if result.Error != nil {
		return wrapStepError("nvim", "Install Neovim",
			"Failed to install Neovim and dependencies",
			result.Error)
	}

	// Copy config
	SendLog(stepID, "Copying Neovim configuration...")
	nvimDir := filepath.Join(homeDir, ".config/nvim")
	if err := system.EnsureDir(nvimDir); err != nil {
		return wrapStepError("nvim", "Install Neovim",
			"Failed to create Neovim config directory",
			err)
	}
	// Copy nvim config directory
	srcNvim := filepath.Join(repoDir, "GentlemanNvim", "nvim")
	if err := system.CopyDir(srcNvim, nvimDir); err != nil {
		return wrapStepError("nvim", "Install Neovim",
			"Failed to copy Neovim configuration",
			err)
	}

	SendLog(stepID, "✓ Neovim configured with Gentleman setup")
	return nil
}

func stepInstallZed(m *Model) error {
	homeDir := os.Getenv("HOME")
	repoDir := m.RepoDir
	stepID := "zed"

	// Skip on Termux — Zed requires GUI with Vulkan
	if m.SystemInfo.IsTermux {
		SendLog(stepID, "Skipping Zed on Termux (requires GUI with Vulkan)")
		return nil
	}

	// Install Zed binary
	if !system.CommandExists("zed") {
		SendLog(stepID, "Installing Zed editor...")
		var result *system.ExecResult
		switch m.SystemInfo.OS {
		case system.OSMac:
			result = system.RunBrewWithLogs("install --cask zed", nil, func(line string) {
				SendLog(stepID, line)
			})
		case system.OSArch:
			result = system.RunSudoWithLogs("pacman -S --noconfirm zed", nil, func(line string) {
				SendLog(stepID, line)
			})
		case system.OSDebian, system.OSLinux, system.OSFedora:
			result = system.RunWithLogs("bash -c 'curl -f https://zed.dev/install.sh | sh'", nil, func(line string) {
				SendLog(stepID, line)
			})
		default:
			result = system.RunWithLogs("bash -c 'curl -f https://zed.dev/install.sh | sh'", nil, func(line string) {
				SendLog(stepID, line)
			})
		}
		if result != nil && result.Error != nil {
			SendLog(stepID, "Warning: Zed install failed: "+result.Error.Error())
			SendLog(stepID, "You can install Zed manually from https://zed.dev/download")
		} else {
			SendLog(stepID, "Zed binary installed")
		}
	} else {
		SendLog(stepID, "Zed already installed")
	}

	// Copy config
	SendLog(stepID, "Copying Zed configuration...")
	zedDir := filepath.Join(homeDir, ".config", "zed")
	if err := system.EnsureDir(zedDir); err != nil {
		return wrapStepError("zed", "Install Zed",
			"Failed to create Zed config directory",
			err)
	}

	srcZed := filepath.Join(repoDir, "GentlemanZed")
	if err := system.CopyDir(srcZed, zedDir); err != nil {
		return wrapStepError("zed", "Install Zed",
			"Failed to copy Zed configuration",
			err)
	}

	SendLog(stepID, "✓ Zed configured with Gentleman setup")
	return nil
}

// hasAITool checks if a tool is in the selected AI tools list
func hasAITool(tools []string, name string) bool {
	for _, t := range tools {
		if t == name {
			return true
		}
	}
	return false
}

func stepInstallAITools(m *Model) error {
	homeDir := os.Getenv("HOME")
	repoDir := m.RepoDir
	stepID := "aitools"

	// Install and configure Claude Code
	if hasAITool(m.Choices.AITools, "claude") {
		SendLog(stepID, "Installing Claude Code...")
		system.RunWithLogs(`curl -fsSL https://claude.ai/install.sh | bash`, nil, func(line string) {
			SendLog(stepID, line)
		})

		SendLog(stepID, "Configuring Claude Code...")
		claudeDir := filepath.Join(homeDir, ".claude")
		system.EnsureDir(claudeDir)
		system.EnsureDir(filepath.Join(claudeDir, "output-styles"))
		system.EnsureDir(filepath.Join(claudeDir, "skills"))
		system.EnsureDir(filepath.Join(claudeDir, "plugins"))
		system.CopyFile(filepath.Join(repoDir, "GentlemanClaude/CLAUDE.md"), filepath.Join(claudeDir, "CLAUDE.md"))
		system.CopyFile(filepath.Join(repoDir, "GentlemanClaude/settings.json"), filepath.Join(claudeDir, "settings.json"))
		system.CopyFile(filepath.Join(repoDir, "GentlemanClaude/statusline.sh"), filepath.Join(claudeDir, "statusline.sh"))
		system.Run(fmt.Sprintf("chmod +x %s", filepath.Join(claudeDir, "statusline.sh")), nil)
		system.CopyFile(filepath.Join(repoDir, "GentlemanClaude/output-styles/gentleman.md"), filepath.Join(claudeDir, "output-styles/gentleman.md"))
		system.CopyFile(filepath.Join(repoDir, "GentlemanClaude/mcp-servers.template.json"), filepath.Join(claudeDir, "mcp-servers.template.json"))
		system.CopyFile(filepath.Join(repoDir, "GentlemanClaude/tweakcc-theme.json"), filepath.Join(claudeDir, "tweakcc-theme.json"))
		SendLog(stepID, "⚙️ Copied CLAUDE.md, statusline, output styles, config")

		SendLog(stepID, "Applying tweakcc theme...")
		result := system.Run("npx tweakcc --apply", nil)
		if result.Error == nil {
			SendLog(stepID, "🎨 Applied tweakcc theme")
		} else {
			SendLog(stepID, "⚠️ Could not apply tweakcc theme (run 'npx tweakcc --apply' manually)")
		}
	}

	// Install and configure OpenCode
	if hasAITool(m.Choices.AITools, "opencode") {
		SendLog(stepID, "Installing OpenCode...")
		system.RunWithLogs(`curl -fsSL https://opencode.ai/install | bash`, nil, func(line string) {
			SendLog(stepID, line)
		})

		SendLog(stepID, "Configuring OpenCode...")
		openCodeDir := filepath.Join(homeDir, ".config/opencode")
		system.EnsureDir(openCodeDir)
		system.EnsureDir(filepath.Join(openCodeDir, "themes"))
		system.CopyFile(filepath.Join(repoDir, "GentlemanOpenCode/opencode.json"), filepath.Join(openCodeDir, "opencode.json"))
		system.CopyFile(filepath.Join(repoDir, "GentlemanOpenCode/themes/gentleman.json"), filepath.Join(openCodeDir, "themes/gentleman.json"))
		SendLog(stepID, "🧠 Copied OpenCode config")
	}

	// Install Gemini CLI
	if hasAITool(m.Choices.AITools, "gemini") {
		SendLog(stepID, "Installing Gemini CLI...")
		result := system.RunWithLogs(`npm install -g @google/gemini-cli`, nil, func(line string) {
			SendLog(stepID, line)
		})
		if result.Error != nil {
			SendLog(stepID, "⚠️ Could not install Gemini CLI (run 'npm install -g @google/gemini-cli' manually)")
		} else {
			SendLog(stepID, "✓ Gemini CLI installed")
		}
	}

	// Install and configure OpenAI Codex CLI
	if hasAITool(m.Choices.AITools, "codex") {
		SendLog(stepID, "Installing Codex CLI...")
		result := system.RunWithLogs(`npm install -g @openai/codex`, nil, func(line string) {
			SendLog(stepID, line)
		})
		if result.Error != nil {
			SendLog(stepID, "⚠️ Could not install Codex CLI (run 'npm install -g @openai/codex' manually)")
		} else {
			SendLog(stepID, "✓ Codex CLI installed")
		}

		SendLog(stepID, "Configuring Codex CLI...")
		codexDir := filepath.Join(homeDir, ".codex")
		system.EnsureDir(codexDir)

		// Copy CLAUDE.md as AGENTS.md (Codex reads AGENTS.md for instructions)
		system.CopyFile(filepath.Join(repoDir, "GentlemanClaude/CLAUDE.md"), filepath.Join(codexDir, "AGENTS.md"))
		SendLog(stepID, "⚙️ Copied AGENTS.md to ~/.codex/")
	}

	// Install and configure Qwen Code
	if hasAITool(m.Choices.AITools, "qwen") {
		SendLog(stepID, "Installing Qwen Code...")
		result := system.RunWithLogs(`npm install -g @qwen-code/qwen-code@latest`, nil, func(line string) {
			SendLog(stepID, line)
		})
		if result.Error != nil {
			SendLog(stepID, "⚠️ Could not install Qwen Code (run 'npm install -g @qwen-code/qwen-code@latest' manually)")
		} else {
			SendLog(stepID, "✓ Qwen Code installed")
		}

		SendLog(stepID, "Configuring Qwen Code...")
		qwenDir := filepath.Join(homeDir, ".qwen")
		system.EnsureDir(qwenDir)
		system.EnsureDir(filepath.Join(qwenDir, "skills"))
		system.CopyFile(filepath.Join(repoDir, "GentlemanQwen/QWEN.md"), filepath.Join(qwenDir, "QWEN.md"))
		system.CopyFile(filepath.Join(repoDir, "GentlemanQwen/settings.json"), filepath.Join(qwenDir, "settings.json"))
		SendLog(stepID, "🧠 Copied QWEN.md and settings config to ~/.qwen/")
	}

	// Install GitHub Copilot CLI extension
	if hasAITool(m.Choices.AITools, "copilot") {
		SendLog(stepID, "Installing GitHub Copilot CLI...")
		result := system.RunWithLogs(`gh extension install github/gh-copilot`, nil, func(line string) {
			SendLog(stepID, line)
		})
		if result.Error != nil {
			SendLog(stepID, "⚠️ Could not install GitHub Copilot (run 'gh extension install github/gh-copilot' manually)")
		} else {
			SendLog(stepID, "✓ GitHub Copilot CLI installed")
		}
	}

	// Centralize skills from Gentleman-Skills repo
	if err := setupCentralizedSkills(m); err != nil {
		SendLog(stepID, fmt.Sprintf("⚠️ Centralized skills setup failed: %v", err))
	}

	SendLog(stepID, "✓ AI tools configured")
	return nil
}

// setupCentralizedSkills clones Gentleman-Skills repo to ~/.gentleman/skills/ and
// creates symlinks into each CLI's skill discovery path.
// Central source: ~/.gentleman/skills/ (curated/ + community/)
// Claude:              ~/.claude/skills/<name>  → central/<name>
// OpenCode/Codex/Gemini: ~/.agents/skills/<name> → central/<name>
func setupCentralizedSkills(m *Model) error {
	homeDir := os.Getenv("HOME")
	stepID := "aitools"
	centralDir := filepath.Join(homeDir, ".gentleman", "skills")

	// Determine if any CLI needs skills
	needsClaude := hasAITool(m.Choices.AITools, "claude")
	needsAgents := hasAITool(m.Choices.AITools, "opencode") ||
		hasAITool(m.Choices.AITools, "codex") ||
		hasAITool(m.Choices.AITools, "gemini")

	if !needsClaude && !needsAgents {
		return nil
	}

	SendLog(stepID, "Setting up centralized skills...")

	// Clone or update Gentleman-Skills repo
	needsClone := true
	if info, err := os.Stat(centralDir); err == nil {
		if time.Since(info.ModTime()) < time.Hour {
			needsClone = false
			SendLog(stepID, "Using cached Gentleman-Skills repo")
		} else {
			os.RemoveAll(centralDir)
		}
	}

	if needsClone {
		SendLog(stepID, "Cloning Gentleman-Skills...")
		system.EnsureDir(filepath.Join(homeDir, ".gentleman"))
		result := system.RunWithLogs(
			"git clone --depth 1 https://github.com/Gentleman-Programming/Gentleman-Skills.git "+centralDir,
			nil, func(line string) { SendLog(stepID, line) },
		)
		if result.Error != nil {
			return fmt.Errorf("failed to clone Gentleman-Skills: %w", result.Error)
		}
	}

	// Discover all skill directories (curated/ and community/)
	var skillPaths []string
	for _, subdir := range []string{"curated", "community"} {
		dir := filepath.Join(centralDir, subdir)
		entries, err := os.ReadDir(dir)
		if err != nil {
			continue
		}
		for _, entry := range entries {
			if !entry.IsDir() {
				continue
			}
			// Verify it has a SKILL.md
			skillFile := filepath.Join(dir, entry.Name(), "SKILL.md")
			if _, err := os.Stat(skillFile); err == nil {
				skillPaths = append(skillPaths, filepath.Join(dir, entry.Name()))
			}
		}
	}

	if len(skillPaths) == 0 {
		SendLog(stepID, "⚠️ No skills found in Gentleman-Skills repo")
		return nil
	}

	SendLog(stepID, fmt.Sprintf("Found %d skills in Gentleman-Skills", len(skillPaths)))

	// Create symlinks for Claude (~/.claude/skills/<name>)
	if needsClaude {
		claudeSkillsDir := filepath.Join(homeDir, ".claude", "skills")
		system.EnsureDir(claudeSkillsDir)
		linked := 0
		for _, sp := range skillPaths {
			name := filepath.Base(sp)
			dst := filepath.Join(claudeSkillsDir, name)
			// Remove existing (file, dir, or stale symlink)
			os.RemoveAll(dst)
			if err := os.Symlink(sp, dst); err != nil {
				SendLog(stepID, fmt.Sprintf("⚠️ Could not symlink %s for Claude: %v", name, err))
			} else {
				linked++
			}
		}
		SendLog(stepID, fmt.Sprintf("🔗 Linked %d skills → ~/.claude/skills/", linked))
	}

	// Create symlinks for OpenCode/Codex/Gemini (~/.agents/skills/<name>)
	if needsAgents {
		agentsSkillsDir := filepath.Join(homeDir, ".agents", "skills")
		system.EnsureDir(agentsSkillsDir)
		linked := 0
		for _, sp := range skillPaths {
			name := filepath.Base(sp)
			dst := filepath.Join(agentsSkillsDir, name)
			// Remove existing (file, dir, or stale symlink)
			os.RemoveAll(dst)
			if err := os.Symlink(sp, dst); err != nil {
				SendLog(stepID, fmt.Sprintf("⚠️ Could not symlink %s for agents: %v", name, err))
			} else {
				linked++
			}
		}
		SendLog(stepID, fmt.Sprintf("🔗 Linked %d skills → ~/.agents/skills/", linked))
	}

	return nil
}

func stepInstallAIFramework(m *Model) error {
	stepID := "aiframework"

	// Determine which features to install via setup-global.sh
	var features []string
	if m.Choices.AIFrameworkPreset != "" {
		// Map presets to feature combinations
		presetFeatures := map[string][]string{
			"minimal":   {"hooks", "commands", "sdd"},
			"frontend":  {"hooks", "commands", "skills", "agents", "sdd"},
			"backend":   {"hooks", "commands", "skills", "agents", "sdd"},
			"fullstack": {"hooks", "commands", "skills", "agents", "sdd", "mcp"},
			"data":      {"hooks", "commands", "skills", "agents", "sdd", "mcp"},
			"complete":  {"hooks", "commands", "skills", "agents", "sdd", "mcp"},
		}
		if f, ok := presetFeatures[m.Choices.AIFrameworkPreset]; ok {
			features = f
		} else {
			features = []string{"hooks", "commands", "skills", "agents", "sdd", "mcp"}
		}
	} else if len(m.Choices.AIFrameworkModules) > 0 {
		// Custom selection — already feature-level IDs from collectSelectedFeatures
		features = m.Choices.AIFrameworkModules
	}

	// Run project-starter-framework setup if there are features to install
	if len(features) > 0 {
		// Clean up any leftover clone from a previous failed run
		system.Run("rm -rf /tmp/project-starter-framework-install", nil)

		SendLog(stepID, "Cloning project-starter-framework...")
		result := system.RunWithLogs(
			"git clone --depth 1 https://github.com/JNZader/project-starter-framework.git /tmp/project-starter-framework-install",
			nil, func(line string) { SendLog(stepID, line) },
		)
		if result.Error != nil {
			return wrapStepError("aiframework", "Install AI Framework",
				"Failed to clone project-starter-framework", result.Error)
		}

		// Build the setup-global.sh command
		setupCmd := "/tmp/project-starter-framework-install/scripts/setup-global.sh --auto --skip-install"

		// Determine which CLIs to configure based on selected AI tools
		var clis []string
		if hasAITool(m.Choices.AITools, "claude") {
			clis = append(clis, "claude")
		}
		if hasAITool(m.Choices.AITools, "opencode") {
			clis = append(clis, "opencode")
		}
		if hasAITool(m.Choices.AITools, "gemini") {
			clis = append(clis, "gemini")
		}
		if hasAITool(m.Choices.AITools, "copilot") {
			clis = append(clis, "copilot")
		}
		if hasAITool(m.Choices.AITools, "codex") {
			clis = append(clis, "codex")
		}
		if hasAITool(m.Choices.AITools, "qwen") {
			clis = append(clis, "qwen")
		}
		if len(clis) > 0 {
			setupCmd += " --clis=" + strings.Join(clis, ",")
		}

		setupCmd += " --features=" + strings.Join(features, ",")

		SendLog(stepID, "Running framework setup...")
		SendLog(stepID, fmt.Sprintf("Command: %s", setupCmd))
		result = system.RunWithLogs(setupCmd, nil, func(line string) {
			SendLog(stepID, line)
		})
		if result.Error != nil {
			return wrapStepError("aiframework", "Install AI Framework",
				"Framework setup failed", result.Error)
		}

		// Cleanup cloned framework repo
		system.Run("rm -rf /tmp/project-starter-framework-install", nil)

		SendLog(stepID, "✓ AI framework configured")
	}

	// Install Agent Teams Lite if selected (separate SDD framework)
	if m.Choices.InstallAgentTeamsLite {
		SendLog(stepID, "Installing Agent Teams Lite...")
		if err := installAgentTeamsLite(m); err != nil {
			SendLog(stepID, fmt.Sprintf("⚠️ Agent Teams Lite failed: %v (you can install manually)", err))
		} else {
			SendLog(stepID, "✓ Agent Teams Lite installed")
		}
	}

	return nil
}

// stepInstallEngram installs Engram MCP server and configures auto-start service
func stepInstallEngram(m *Model) error {
	stepID := "engram"
	homeDir := os.Getenv("HOME")

	// Check if engram is already installed
	engramPath, err := exec.LookPath("engram")
	if err == nil {
		SendLog(stepID, fmt.Sprintf("Engram already installed at: %s", engramPath))
		SendLog(stepID, "Checking configuration...")

		// Check if configured for OpenCode
		opencodePluginPath := filepath.Join(homeDir, ".config/opencode/plugins/engram.ts")
		if _, err := os.Stat(opencodePluginPath); err == nil {
			SendLog(stepID, "✓ Engram already configured for OpenCode")

			// Check if service is configured
			if setupEngramService(m) {
				SendLog(stepID, "✓ Engram service configured")
			}

			return nil
		}
	}

	// Install engram via Homebrew
	SendLog(stepID, "Installing Engram via Homebrew...")
	result := system.RunWithLogs("brew install gentleman-programming/tap/engram", nil, func(line string) {
		SendLog(stepID, line)
	})
	if result.Error != nil {
		SendLog(stepID, "⚠️ Could not install Engram via Homebrew")
		SendLog(stepID, "You can manually install from: https://github.com/Gentleman-Programming/engram")
		return nil // Don't fail the entire installation
	}

	SendLog(stepID, "✓ Engram installed")

	// Configure for OpenCode
	SendLog(stepID, "Configuring Engram for OpenCode...")
	result = system.RunWithLogs("engram setup opencode", nil, func(line string) {
		SendLog(stepID, line)
	})
	if result.Error != nil {
		SendLog(stepID, fmt.Sprintf("⚠️ Could not configure Engram: %v", result.Error))
	} else {
		SendLog(stepID, "✓ Engram configured for OpenCode")
	}

	// Setup auto-start service
	if setupEngramService(m) {
		SendLog(stepID, "✓ Engram service configured for auto-start")
	}

	SendLog(stepID, "✓ Engram installation complete")
	SendLog(stepID, "Engram provides persistent memory: mem_save, mem_search, mem_context")

	return nil
}

// setupEngramService configures systemd (Linux) or launchd (macOS) for engram
func setupEngramService(m *Model) bool {
	stepID := "engram"
	homeDir := os.Getenv("HOME")

	switch runtime.GOOS {
	case "linux":
		return setupEngramSystemd(homeDir, stepID)
	case "darwin":
		return setupEngramLaunchd(homeDir, stepID)
	default:
		SendLog(stepID, "⚠️ Auto-start service not supported on this OS")
		return false
	}
}

// setupEngramSystemd creates systemd user service for Linux
func setupEngramSystemd(homeDir, stepID string) bool {
	configDir := filepath.Join(homeDir, ".config/systemd/user")
	serviceFile := filepath.Join(configDir, "engram.service")

	// Ensure directory exists
	if err := os.MkdirAll(configDir, 0755); err != nil {
		SendLog(stepID, fmt.Sprintf("⚠️ Could not create systemd directory: %v", err))
		return false
	}

	// Create service file
	serviceContent := `[Unit]
Description=Engram Persistent Memory Server
After=network.target

[Service]
Type=simple
ExecStart=%s/.local/bin/engram serve
Restart=always
RestartSec=10
Environment="HOME=%s"
Environment="ENGRAM_DATA_DIR=%s/.engram"

[Install]
WantedBy=default.target
`

	// Find engram binary location
	engramPath := filepath.Join(homeDir, ".local/bin/engram")
	if _, err := os.Stat(engramPath); os.IsNotExist(err) {
		// Try to find in PATH
		if path, err := exec.LookPath("engram"); err == nil {
			engramPath = path
		}
	}

	content := fmt.Sprintf(serviceContent, homeDir, homeDir, homeDir)

	if err := os.WriteFile(serviceFile, []byte(content), 0644); err != nil {
		SendLog(stepID, fmt.Sprintf("⚠️ Could not create systemd service: %v", err))
		return false
	}

	// Enable and start service
	result := system.Run("systemctl --user daemon-reload", nil)
	if result.Error != nil {
		SendLog(stepID, fmt.Sprintf("⚠️ Could not reload systemd: %v", result.Error))
		return false
	}

	result = system.Run("systemctl --user enable engram.service", nil)
	if result.Error != nil {
		SendLog(stepID, fmt.Sprintf("⚠️ Could not enable engram service: %v", result.Error))
		return false
	}

	// Try to start, but don't fail if it doesn't (might need logout/login)
	result = system.Run("systemctl --user start engram.service", nil)
	if result.Error != nil {
		SendLog(stepID, "Note: Engram service enabled but not started (will start on next login)")
	} else {
		SendLog(stepID, "✓ Engram service started")
	}

	return true
}

// setupEngramLaunchd creates launchd plist for macOS
func setupEngramLaunchd(homeDir, stepID string) bool {
	launchAgentsDir := filepath.Join(homeDir, "Library/LaunchAgents")
	plistFile := filepath.Join(launchAgentsDir, "com.gentleman.engram.plist")

	// Ensure directory exists
	if err := os.MkdirAll(launchAgentsDir, 0755); err != nil {
		SendLog(stepID, fmt.Sprintf("⚠️ Could not create LaunchAgents directory: %v", err))
		return false
	}

	// Find engram binary
	engramPath := "/opt/homebrew/bin/engram"
	if _, err := os.Stat(engramPath); os.IsNotExist(err) {
		engramPath = "/usr/local/bin/engram"
		if _, err := os.Stat(engramPath); os.IsNotExist(err) {
			// Try to find in PATH
			if path, err := exec.LookPath("engram"); err == nil {
				engramPath = path
			}
		}
	}

	plistContent := `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key>
    <string>com.gentleman.engram</string>
    <key>ProgramArguments</key>
    <array>
        <string>%s</string>
        <string>serve</string>
    </array>
    <key>RunAtLoad</key>
    <true/>
    <key>KeepAlive</key>
    <true/>
    <key>EnvironmentVariables</key>
    <dict>
        <key>HOME</key>
        <string>%s</string>
        <key>ENGRAM_DATA_DIR</key>
        <string>%s/.engram</string>
    </dict>
    <key>StandardOutPath</key>
    <string>%s/.engram/engram.log</string>
    <key>StandardErrorPath</key>
    <string>%s/.engram/engram.error.log</string>
</dict>
</plist>
`

	content := fmt.Sprintf(plistContent, engramPath, homeDir, homeDir, homeDir, homeDir)

	if err := os.WriteFile(plistFile, []byte(content), 0644); err != nil {
		SendLog(stepID, fmt.Sprintf("⚠️ Could not create launchd plist: %v", err))
		return false
	}

	// Load the plist
	result := system.Run(fmt.Sprintf("launchctl load %s", plistFile), nil)
	if result.Error != nil {
		SendLog(stepID, fmt.Sprintf("⚠️ Could not load launchd service: %v", result.Error))
		return false
	}

	// Try to start
	result = system.Run("launchctl start com.gentleman.engram", nil)
	if result.Error != nil {
		SendLog(stepID, "Note: Engram service loaded but not started (will start on next login)")
	} else {
		SendLog(stepID, "✓ Engram service started")
	}

	return true
}

// installAgentTeamsLite clones the agent-teams-lite repo and runs install.sh for each selected AI tool.
func installAgentTeamsLite(m *Model) error {
	const repoURL = "https://github.com/Gentleman-Programming/agent-teams-lite.git"
	const clonePath = "/tmp/agent-teams-lite-install"
	stepID := "aiframework"

	// Cleanup any leftover
	system.Run("rm -rf "+clonePath, nil)

	SendLog(stepID, "Cloning agent-teams-lite...")
	result := system.RunWithLogs(
		"git clone --depth 1 "+repoURL+" "+clonePath,
		nil, func(line string) { SendLog(stepID, line) },
	)
	if result.Error != nil {
		return fmt.Errorf("failed to clone agent-teams-lite: %w", result.Error)
	}

	// Make install script executable
	system.Run("chmod +x "+clonePath+"/scripts/install.sh", nil)

	// Map our AI tool IDs to agent-teams-lite agent names
	agentMap := map[string]string{
		"claude":   "claude-code",
		"opencode": "opencode",
		"gemini":   "gemini-cli",
		"codex":    "codex",
		"qwen":     "qwen-code",
	}

	installed := 0
	for _, tool := range m.Choices.AITools {
		agentName, ok := agentMap[tool]
		if !ok {
			continue
		}
		SendLog(stepID, fmt.Sprintf("Installing Agent Teams Lite for %s...", agentName))
		installCmd := fmt.Sprintf("%s/scripts/install.sh --agent %s", clonePath, agentName)
		result = system.RunWithLogs(installCmd, nil, func(line string) {
			SendLog(stepID, line)
		})
		if result.Error != nil {
			SendLog(stepID, fmt.Sprintf("⚠️ Agent Teams Lite install failed for %s", agentName))
		} else {
			installed++
		}
	}

	// Cleanup
	system.Run("rm -rf "+clonePath, nil)

	if installed == 0 {
		return fmt.Errorf("no AI tools could be configured with Agent Teams Lite")
	}
	return nil
}

func stepCleanup(m *Model) error {
	stepID := "cleanup"
	SendLog(stepID, "Removing temporary files...")
	// Only remove the cloned repo - no sudo needed
	result := system.Run("rm -rf "+m.RepoDir, nil)
	if result.Error != nil {
		// Non-critical error, just log it
		SendLog(stepID, "Warning: Could not remove temporary directory")
		return nil
	}
	SendLog(stepID, "✓ Cleanup complete")
	return nil
}

// stepSetDefaultShell sets the selected shell as the user's default shell
// In non-interactive mode, this handles Termux specially (via .bashrc)
// and attempts to set the shell on other systems if possible
func stepSetDefaultShell(m *Model) error {
	stepID := "setshell"
	shell := m.Choices.Shell
	homeDir := os.Getenv("HOME")

	var shellCmd string
	switch shell {
	case "fish":
		shellCmd = "fish"
	case "zsh":
		shellCmd = "zsh"
	case "nushell":
		shellCmd = "nu"
	default:
		SendLog(stepID, fmt.Sprintf("Unknown shell: %s, skipping", shell))
		return nil
	}

	// Termux: no chsh available, modify .bashrc to auto-start shell
	if m.SystemInfo.IsTermux {
		SendLog(stepID, "Configuring shell auto-start for Termux...")

		// Find the shell path
		shellPath := system.Run(fmt.Sprintf("which %s", shellCmd), nil)
		if shellPath.Error != nil || strings.TrimSpace(shellPath.Output) == "" {
			SendLog(stepID, fmt.Sprintf("Shell '%s' not found in PATH, skipping", shellCmd))
			return nil
		}
		shellPathStr := strings.TrimSpace(shellPath.Output)

		// Read existing .bashrc
		bashrcPath := filepath.Join(homeDir, ".bashrc")
		existingContent := ""
		if data, err := os.ReadFile(bashrcPath); err == nil {
			existingContent = string(data)
		}

		// Check if already configured
		if strings.Contains(existingContent, "# Gentleman.Dots shell auto-start") {
			SendLog(stepID, "Shell auto-start already configured in ~/.bashrc")
			return nil
		}

		// Append auto-start configuration
		autoStartConfig := fmt.Sprintf(`
# Gentleman.Dots shell auto-start
if [ -x "%s" ] && [ -z "$GENTLEMANDOTS_SHELL_STARTED" ]; then
    export GENTLEMANDOTS_SHELL_STARTED=1
    exec %s
fi
`, shellPathStr, shellPathStr)

		f, err := os.OpenFile(bashrcPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return wrapStepError("setshell", "Set Default Shell",
				"Failed to open ~/.bashrc for writing",
				err)
		}
		defer f.Close()

		if _, err := f.WriteString(autoStartConfig); err != nil {
			return wrapStepError("setshell", "Set Default Shell",
				"Failed to write shell auto-start to ~/.bashrc",
				err)
		}

		SendLog(stepID, fmt.Sprintf("✓ Configured %s to auto-start in ~/.bashrc", shell))
		SendLog(stepID, "Close and reopen Termux for changes to take effect")
		return nil
	}

	// Non-Termux: Try to set shell using sudo usermod (works if NOPASSWD configured)
	// Find the shell path first
	shellPath := system.Run(fmt.Sprintf("which %s", shellCmd), nil)
	if shellPath.Error != nil || strings.TrimSpace(shellPath.Output) == "" {
		SendLog(stepID, fmt.Sprintf("Shell '%s' not found in PATH, skipping", shellCmd))
		return nil
	}
	shellPathStr := strings.TrimSpace(shellPath.Output)

	// Get current username
	currentUser := os.Getenv("USER")
	if currentUser == "" {
		currentUser = os.Getenv("LOGNAME")
	}
	if currentUser == "" {
		// Fallback to whoami command (useful in Docker containers)
		whoamiResult := system.Run("whoami", nil)
		if whoamiResult.Error == nil {
			currentUser = strings.TrimSpace(whoamiResult.Output)
		}
	}
	if currentUser == "" {
		SendLog(stepID, "Could not determine current user, skipping shell change")
		return nil
	}

	// First, ensure shell is in /etc/shells
	SendLog(stepID, fmt.Sprintf("Adding %s to /etc/shells if needed...", shellPathStr))
	checkShells := system.Run(fmt.Sprintf("grep -q '^%s$' /etc/shells", shellPathStr), nil)
	if checkShells.Error != nil {
		// Shell not in /etc/shells, try to add it
		addResult := system.RunSudo(fmt.Sprintf("sh -c 'echo \"%s\" >> /etc/shells'", shellPathStr), nil)
		if addResult.Error != nil {
			SendLog(stepID, fmt.Sprintf("Could not add %s to /etc/shells (may need manual setup)", shellPathStr))
		}
	}

	// Try sudo usermod first (more reliable than chsh in scripts)
	SendLog(stepID, fmt.Sprintf("Setting %s as default shell for %s...", shell, currentUser))
	result := system.RunSudo(fmt.Sprintf("usermod -s %s %s", shellPathStr, currentUser), nil)
	if result.Error != nil {
		// usermod failed, try chsh as fallback
		SendLog(stepID, "usermod failed, trying chsh...")
		result = system.RunSudo(fmt.Sprintf("chsh -s %s %s", shellPathStr, currentUser), nil)
		if result.Error != nil {
			// Both failed - not critical, just inform user
			SendLog(stepID, fmt.Sprintf("Could not set default shell automatically"))
			SendLog(stepID, fmt.Sprintf("Run manually: chsh -s %s", shellPathStr))
			return nil
		}
	}

	SendLog(stepID, fmt.Sprintf("✓ Default shell set to %s", shell))
	SendLog(stepID, "Log out and log back in for changes to take effect")
	return nil
}

// runProjectInitScript clones project-starter-framework and runs init-project.sh.
// If memory is "obsidian-brain" and rolePacks is non-empty, it also copies
// role pack templates into the project vault after the script finishes.
func runProjectInitScript(projectPath, memory, ci string, engram bool, rolePacks []string) error {
	cacheDir := filepath.Join(os.TempDir(), "project-starter-framework-install")

	// Check cache freshness (1 hour)
	needsClone := true
	if info, err := os.Stat(cacheDir); err == nil {
		if time.Since(info.ModTime()) < time.Hour {
			needsClone = false
		} else {
			os.RemoveAll(cacheDir)
		}
	}

	if needsClone {
		if globalProgram != nil {
			globalProgram.Send(projectInstallLogMsg{line: "Cloning project-starter-framework..."})
		}
		cmd := exec.Command("git", "clone", "--depth", "1",
			"https://github.com/JNZader/project-starter-framework.git", cacheDir)
		if out, err := cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("failed to clone framework: %s: %w", string(out), err)
		}
	}

	// Build command
	scriptPath := filepath.Join(cacheDir, "init-project.sh")
	os.Chmod(scriptPath, 0755)

	// Map memory choice to numeric ID
	memoryMap := map[string]string{
		"obsidian-brain": "1",
		"vibekanban":     "2",
		"engram":         "3",
		"simple":         "4",
		"none":           "5",
	}
	memoryNum := memoryMap[memory]
	if memoryNum == "" {
		memoryNum = "5"
	}

	// Map CI choice to numeric ID
	ciMap := map[string]string{
		"github":     "1",
		"gitlab":     "2",
		"woodpecker": "3",
		"none":       "4",
	}
	ciNum := ciMap[ci]
	if ciNum == "" {
		ciNum = "4"
	}

	args := []string{scriptPath, "--non-interactive", "--memory=" + memoryNum, "--ci=" + ciNum}
	if engram {
		args = append(args, "--engram")
	}

	if globalProgram != nil {
		globalProgram.Send(projectInstallLogMsg{line: fmt.Sprintf("Running: bash %s", strings.Join(args, " "))})
	}

	cmd := exec.Command("bash", args...)
	cmd.Dir = projectPath

	output, err := cmd.CombinedOutput()
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" && globalProgram != nil {
			globalProgram.Send(projectInstallLogMsg{line: line})
		}
	}

	if err != nil {
		return fmt.Errorf("init-project.sh failed: %w", err)
	}

	// Copy role pack templates if obsidian-brain is selected and packs were chosen
	if memory == "obsidian-brain" && len(rolePacks) > 0 {
		if globalProgram != nil {
			globalProgram.Send(projectInstallLogMsg{line: "Copying Obsidian Brain role pack templates..."})
		}
		repoRoot := findRepoDirForTemplates("")
		if repoRoot == "" {
			if globalProgram != nil {
				globalProgram.Send(projectInstallLogMsg{line: "⚠ Could not locate template assets, skipping role pack templates"})
			}
		} else {
			if err := copyRolePackTemplates(repoRoot, projectPath, rolePacks); err != nil {
				if globalProgram != nil {
					globalProgram.Send(projectInstallLogMsg{line: fmt.Sprintf("⚠ Role pack template copy failed: %v", err)})
				}
			} else {
				if globalProgram != nil {
					globalProgram.Send(projectInstallLogMsg{line: fmt.Sprintf("✓ Copied templates for packs: %s", strings.Join(rolePacks, ", "))})
				}
			}
		}
	}

	return nil
}

// RunProjectInitScript exposes runProjectInitScript for CLI usage
func RunProjectInitScript(projectPath, memory, ci string, engram bool, rolePacks []string) error {
	return runProjectInitScript(projectPath, memory, ci, engram, rolePacks)
}

// findRepoDirForTemplates locates the Javi.Dots repo root so we can find
// GentlemanNvim/obsidian-brain/ template assets. It tries:
// 1. The provided repoDir (used during full installation when repo is cloned)
// 2. Walking up from the executable path
// 3. Walking up from the current working directory
func findRepoDirForTemplates(repoDir string) string {
	marker := filepath.Join("GentlemanNvim", "obsidian-brain")

	// 1. Check the provided repoDir
	if repoDir != "" {
		if _, err := os.Stat(filepath.Join(repoDir, marker)); err == nil {
			return repoDir
		}
	}

	// 2. Walk up from executable path
	if exePath, err := os.Executable(); err == nil {
		dir := filepath.Dir(exePath)
		for i := 0; i < 10; i++ {
			if _, err := os.Stat(filepath.Join(dir, marker)); err == nil {
				return dir
			}
			parent := filepath.Dir(dir)
			if parent == dir {
				break
			}
			dir = parent
		}
	}

	// 3. Walk up from current working directory
	if cwd, err := os.Getwd(); err == nil {
		dir := cwd
		for i := 0; i < 10; i++ {
			if _, err := os.Stat(filepath.Join(dir, marker)); err == nil {
				return dir
			}
			parent := filepath.Dir(dir)
			if parent == dir {
				break
			}
			dir = parent
		}
	}

	return ""
}

// copyRolePackTemplates copies selected role pack templates into the project vault.
// repoDir is the path to the Javi.Dots repo root containing GentlemanNvim/obsidian-brain/.
func copyRolePackTemplates(repoDir, projectPath string, rolePacks []string) error {
	vaultDir := filepath.Join(projectPath, ".obsidian-brain")
	templatesDir := filepath.Join(vaultDir, "templates")

	// Create core vault folder structure
	coreDirs := []string{
		vaultDir,
		filepath.Join(vaultDir, "inbox"),
		filepath.Join(vaultDir, "resources"),
		filepath.Join(vaultDir, "knowledge"),
		templatesDir,
		filepath.Join(vaultDir, ".obsidian"), // marker for obsidian.nvim plugin detection
	}
	for _, dir := range coreDirs {
		if err := system.EnsureDir(dir); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	// Create role-specific directories
	for _, pack := range rolePacks {
		switch pack {
		case "developer":
			for _, dir := range []string{"architecture", "sessions", "debugging"} {
				if err := system.EnsureDir(filepath.Join(vaultDir, dir)); err != nil {
					return fmt.Errorf("failed to create developer directory %s: %w", dir, err)
				}
			}
		case "pm-lead":
			for _, dir := range []string{"meetings", "sprints", "risks", "briefs"} {
				if err := system.EnsureDir(filepath.Join(vaultDir, dir)); err != nil {
					return fmt.Errorf("failed to create pm-lead directory %s: %w", dir, err)
				}
			}
		}
	}

	// Copy templates from each selected role pack
	for _, pack := range rolePacks {
		srcDir := filepath.Join(repoDir, "GentlemanNvim", "obsidian-brain", pack, "templates")
		if _, err := os.Stat(srcDir); os.IsNotExist(err) {
			continue // pack directory doesn't exist yet, skip gracefully
		}
		entries, err := os.ReadDir(srcDir)
		if err != nil {
			return fmt.Errorf("failed to read pack %s templates: %w", pack, err)
		}
		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}
			src := filepath.Join(srcDir, entry.Name())
			dst := filepath.Join(templatesDir, entry.Name())
			data, err := os.ReadFile(src)
			if err != nil {
				return fmt.Errorf("failed to read template %s: %w", src, err)
			}
			if err := os.WriteFile(dst, data, 0644); err != nil {
				return fmt.Errorf("failed to write template %s: %w", dst, err)
			}
		}
	}

	return nil
}
