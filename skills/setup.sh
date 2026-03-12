#!/bin/bash
# ============================================================================
# Gentleman.Dots AI Skills Setup Script
# ============================================================================
# This script synchronizes AGENTS.md to tool-specific instruction files.
# AGENTS.md is the single source of truth - edits propagate to all tools.
#
# Usage:
#   ./skills/setup.sh              # Interactive menu
#   ./skills/setup.sh --all        # Generate all formats
#   ./skills/setup.sh --claude     # Generate CLAUDE.md only
#   ./skills/setup.sh --gemini     # Generate GEMINI.md only
#   ./skills/setup.sh --copilot    # Generate .github/copilot-instructions.md
#   ./skills/setup.sh --codex      # Generate CODEX.md only
#   ./skills/setup.sh --opencode   # Sync to OpenCode config
#   ./skills/setup.sh --install-opencode-agents  # Install orchestrator agents
#   ./skills/setup.sh --install-engram  # Install Engram MCP server
#
# ============================================================================

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m'
BOLD='\033[1m'

# Script directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

# ============================================================================
# Utility Functions
# ============================================================================

log_info() {
    printf "${BLUE}[INFO]${NC} %s\n" "$1"
}

log_success() {
    printf "${GREEN}[SUCCESS]${NC} %s\n" "$1"
}

log_warning() {
    printf "${YELLOW}[WARNING]${NC} %s\n" "$1"
}

log_error() {
    printf "${RED}[ERROR]${NC} %s\n" "$1"
}

log_header() {
    printf "\n${CYAN}${BOLD}════════════════════════════════════════${NC}\n"
    printf "${CYAN}${BOLD}  %s${NC}\n" "$1"
    printf "${CYAN}${BOLD}════════════════════════════════════════${NC}\n\n"
}

# ============================================================================
# Generation Functions
# ============================================================================

# Find all AGENTS.md files in the repository
find_agents_files() {
    find "$REPO_ROOT" -name "AGENTS.md" -type f 2>/dev/null
}

# Generate CLAUDE.md from AGENTS.md
generate_claude() {
    local agents_file="$1"
    local dir=$(dirname "$agents_file")
    local claude_file="$dir/CLAUDE.md"

    log_info "Generating CLAUDE.md from $agents_file"

    # Add Claude-specific header
    cat > "$claude_file" << 'EOF'
# Claude Code Instructions

> **Auto-generated from AGENTS.md** - Do not edit directly.
> Run `./skills/setup.sh --claude` to regenerate.

EOF

    # Append AGENTS.md content
    cat "$agents_file" >> "$claude_file"

    log_success "Created $claude_file"
}

# Generate GEMINI.md from AGENTS.md
generate_gemini() {
    local agents_file="$1"
    local dir=$(dirname "$agents_file")
    local gemini_file="$dir/GEMINI.md"

    log_info "Generating GEMINI.md from $agents_file"

    cat > "$gemini_file" << 'EOF'
# Gemini CLI Instructions

> **Auto-generated from AGENTS.md** - Do not edit directly.
> Run `./skills/setup.sh --gemini` to regenerate.

EOF

    cat "$agents_file" >> "$gemini_file"

    log_success "Created $gemini_file"
}

# Generate .github/copilot-instructions.md from AGENTS.md
generate_copilot() {
    local agents_file="$1"
    local dir=$(dirname "$agents_file")
    local copilot_dir="$dir/.github"
    local copilot_file="$copilot_dir/copilot-instructions.md"

    log_info "Generating copilot-instructions.md from $agents_file"

    mkdir -p "$copilot_dir"

    cat > "$copilot_file" << 'EOF'
# GitHub Copilot Instructions

> **Auto-generated from AGENTS.md** - Do not edit directly.
> Run `./skills/setup.sh --copilot` to regenerate.

EOF

    cat "$agents_file" >> "$copilot_file"

    log_success "Created $copilot_file"
}

# Generate CODEX.md from AGENTS.md
generate_codex() {
    local agents_file="$1"
    local dir=$(dirname "$agents_file")
    local codex_file="$dir/CODEX.md"

    log_info "Generating CODEX.md from $agents_file"

    cat > "$codex_file" << 'EOF'
# OpenAI Codex Instructions

> **Auto-generated from AGENTS.md** - Do not edit directly.
> Run `./skills/setup.sh --codex` to regenerate.

EOF

    cat "$agents_file" >> "$codex_file"

    log_success "Created $codex_file"
}

# Sync skills to OpenCode config directory
sync_opencode() {
    local opencode_dir="$HOME/.config/opencode/skill"

    log_info "Syncing skills to OpenCode config..."

    if [ ! -d "$opencode_dir" ]; then
        log_warning "OpenCode skill directory not found: $opencode_dir"
        log_info "Creating directory..."
        mkdir -p "$opencode_dir"
    fi

    # Copy AGENTS.md as the main instruction file
    if [ -f "$REPO_ROOT/AGENTS.md" ]; then
        cp "$REPO_ROOT/AGENTS.md" "$opencode_dir/AGENTS.md"
        log_success "Copied AGENTS.md to OpenCode"
    fi

    # Copy individual skills
    local skills_dir="$REPO_ROOT/GentlemanClaude/skills"
    if [ -d "$skills_dir" ]; then
        for skill_dir in "$skills_dir"/*/; do
            if [ -f "$skill_dir/SKILL.md" ]; then
                skill_name=$(basename "$skill_dir")
                mkdir -p "$opencode_dir/$skill_name"
                cp "$skill_dir/SKILL.md" "$opencode_dir/$skill_name/SKILL.md"
                log_info "  → Copied $skill_name"
            fi
        done
        log_success "Synced all skills to OpenCode"
    fi

    # Copy _shared/ directory (convention files referenced by SDD skills)
    local shared_dir="$REPO_ROOT/GentlemanClaude/skills/_shared"
    if [ -d "$shared_dir" ]; then
        mkdir -p "$opencode_dir/_shared"
        cp "$shared_dir"/*.md "$opencode_dir/_shared/" 2>/dev/null || true
        log_info "  → Copied _shared/"
    fi
}

# Install orchestrator agents to OpenCode
install_opencode_agents() {
    local agents_dir="$HOME/.config/opencode/agents"
    local specialists_dir="$HOME/.config/opencode/specialists"
    local repo_agents_dir="$REPO_ROOT/opencode-agents"

    log_header "Installing OpenCode Orchestrator Agents"

    # Create directories
    mkdir -p "$agents_dir"
    mkdir -p "$specialists_dir"/{development,infrastructure,quality,data-ai,business,domains,planner,docs}

    # Backup existing agents if any
    if [ -f "$agents_dir/gentleman.md" ] || [ -f "$agents_dir/master-orchestrator.md" ]; then
        local backup_dir="$agents_dir/backup-$(date +%Y%m%d-%H%M%S)"
        mkdir -p "$backup_dir"
        # Backup old agent files (not our orchestrators)
        find "$agents_dir" -maxdepth 1 -name "*.md" -not -name "gentleman.md" \
            -not -name "*orchestrator.md" -exec cp {} "$backup_dir/" \; 2>/dev/null || true
        log_info "Backed up existing agents to $backup_dir"
    fi

    # Copy orchestrator agents from repo
    if [ -d "$repo_agents_dir" ]; then
        log_info "Installing orchestrator agents..."
        for agent_file in "$repo_agents_dir"/*.md; do
            if [ -f "$agent_file" ]; then
                cp "$agent_file" "$agents_dir/"
                log_info "  → Installed $(basename "$agent_file")"
            fi
        done
        log_success "Installed 12 orchestrator agents"
    else
        log_warning "Agent templates not found: $repo_agents_dir"
        log_info "Run from repo root with ./skills/setup.sh"
    fi

    # Note about specialists
    log_info ""
    log_info "Specialist agents are stored in: $specialists_dir"
    log_info "These are hidden from OpenCode UI and accessed via orchestrators"
    log_info ""
    log_success "OpenCode agent setup complete!"
    log_info ""
    log_info "Available orchestrators in OpenCode:"
    log_info "  • gentleman - Javi.Dots expert"
    log_info "  • sdd-orchestrator - SDD workflow"
    log_info "  • master-orchestrator - Delegates to specialized orchestrators"
    log_info "  • architect-orchestrator - System design"
    log_info "  • frontend-orchestrator - React, Vue, Angular"
    log_info "  • backend-orchestrator - Python, Go, Java, Node"
    log_info "  • devops-orchestrator - Docker, K8s, CI/CD"
    log_info "  • data-ai-orchestrator - ML, AI, Data"
    log_info "  • quality-orchestrator - Testing, Security"
    log_info "  • business-orchestrator - PM, Analysis"
    log_info "  • specialized-orchestrator - Blockchain, Games"
    log_info "  • planner-orchestrator - Project planning"
}

# Install unified orchestrator instructions for Copilot
install_copilot_orchestrator() {
    local instructions_file="$REPO_ROOT/.github/copilot-instructions.md"
    local unified_dir="$REPO_ROOT/unified-instructions"

    log_header "Installing Copilot Orchestrator Instructions"

    if [ ! -f "$unified_dir/orchestrator.md" ]; then
        log_warning "Unified instructions not found: $unified_dir/orchestrator.md"
        log_info "Skipping Copilot orchestrator installation"
        return 1
    fi

    # Backup existing if present
    if [ -f "$instructions_file" ]; then
        local backup_file="$REPO_ROOT/.github/copilot-instructions.md.backup-$(date +%Y%m%d-%H%M%S)"
        cp "$instructions_file" "$backup_file"
        log_info "Backed up existing copilot-instructions.md"
    fi

    # Install unified orchestrator
    cp "$unified_dir/orchestrator.md" "$instructions_file"
    log_success "Installed unified orchestrator for GitHub Copilot"
    log_info ""
    log_info "Copilot now uses context-aware orchestration:"
    log_info "  • Detects task domain automatically"
    log_info "  • Applies appropriate patterns from unified instructions"
    log_info "  • Single instruction file with all domains"
}

# Install unified orchestrator instructions for Gemini
install_gemini_orchestrator() {
    local instructions_file="$REPO_ROOT/GEMINI.md"
    local unified_dir="$REPO_ROOT/unified-instructions"

    log_header "Installing Gemini Orchestrator Instructions"

    if [ ! -f "$unified_dir/orchestrator.md" ]; then
        log_warning "Unified instructions not found: $unified_dir/orchestrator.md"
        log_info "Skipping Gemini orchestrator installation"
        return 1
    fi

    # Backup existing if present
    if [ -f "$instructions_file" ]; then
        local backup_file="$REPO_ROOT/GEMINI.md.backup-$(date +%Y%m%d-%H%M%S)"
        cp "$instructions_file" "$backup_file"
        log_info "Backed up existing GEMINI.md"
    fi

    # Install unified orchestrator with Gemini header
    cat > "$instructions_file" << 'HEADER'
# Gemini CLI Instructions

> Auto-generated orchestrator for Gemini CLI
> Run `./skills/setup.sh --install-gemini-orchestrator` to regenerate

HEADER

    cat "$unified_dir/orchestrator.md" >> "$instructions_file"
    log_success "Installed unified orchestrator for Gemini CLI"
    log_info ""
    log_info "Gemini now uses context-aware orchestration"
}

# Install unified orchestrator instructions for Codex
install_codex_orchestrator() {
    local instructions_file="$REPO_ROOT/CODEX.md"
    local unified_dir="$REPO_ROOT/unified-instructions"

    log_header "Installing Codex Orchestrator Instructions"

    if [ ! -f "$unified_dir/orchestrator.md" ]; then
        log_warning "Unified instructions not found: $unified_dir/orchestrator.md"
        log_info "Skipping Codex orchestrator installation"
        return 1
    fi

    # Backup existing if present
    if [ -f "$instructions_file" ]; then
        local backup_file="$REPO_ROOT/CODEX.md.backup-$(date +%Y%m%d-%H%M%S)"
        cp "$instructions_file" "$backup_file"
        log_info "Backed up existing CODEX.md"
    fi

    # Install unified orchestrator with Codex header
    cat > "$instructions_file" << 'HEADER'
# OpenAI Codex Instructions

> Auto-generated orchestrator for OpenAI Codex
> Run `./skills/setup.sh --install-codex-orchestrator` to regenerate

HEADER

    cat "$unified_dir/orchestrator.md" >> "$instructions_file"
    log_success "Installed unified orchestrator for OpenAI Codex"
    log_info ""
    log_info "Codex now uses context-aware orchestration"
}

# Install all orchestrators
install_all_orchestrators() {
    install_opencode_agents
    install_copilot_orchestrator
    install_gemini_orchestrator
    install_codex_orchestrator
}

# Sync skills to Claude Code config directory
sync_claude_config() {
    local claude_dir="$HOME/.claude/skills"

    log_info "Syncing skills to Claude Code config..."

    if [ ! -d "$claude_dir" ]; then
        log_warning "Claude skills directory not found: $claude_dir"
        log_info "Creating directory..."
        mkdir -p "$claude_dir"
    fi

    # Copy individual skills
    local skills_dir="$REPO_ROOT/GentlemanClaude/skills"
    if [ -d "$skills_dir" ]; then
        for skill_dir in "$skills_dir"/*/; do
            if [ -f "$skill_dir/SKILL.md" ]; then
                skill_name=$(basename "$skill_dir")
                mkdir -p "$claude_dir/$skill_name"

                # Remove existing file if read-only, then copy
                local dest_file="$claude_dir/$skill_name/SKILL.md"
                if [ -f "$dest_file" ]; then
                    chmod u+w "$dest_file" 2>/dev/null || true
                fi
                cp -f "$skill_dir/SKILL.md" "$dest_file"

                # Copy assets if they exist
                if [ -d "$skill_dir/assets" ]; then
                    chmod -R u+w "$claude_dir/$skill_name/assets" 2>/dev/null || true
                    cp -rf "$skill_dir/assets" "$claude_dir/$skill_name/"
                fi

                log_info "  → Copied $skill_name"
            fi
        done
        log_success "Synced all skills to ~/.claude/skills/"
    fi

    # Copy _shared/ directory (convention files referenced by SDD skills)
    local shared_dir="$REPO_ROOT/GentlemanClaude/skills/_shared"
    if [ -d "$shared_dir" ]; then
        mkdir -p "$claude_dir/_shared"
        chmod -R u+w "$claude_dir/_shared" 2>/dev/null || true
        cp -f "$shared_dir"/*.md "$claude_dir/_shared/"
        log_info "  → Copied _shared/"
    fi
}

# Sync hook scripts to Claude Code config directory
sync_hooks() {
    local hooks_src="$REPO_ROOT/GentlemanClaude/hooks"
    local hooks_dest="$HOME/.claude/hooks"

    log_info "Syncing hooks to ~/.claude/hooks/..."

    if [ ! -d "$hooks_src" ]; then
        log_warning "No hooks directory found at $hooks_src"
        return
    fi

    mkdir -p "$hooks_dest"

    for hook_file in "$hooks_src"/*.sh; do
        [ -f "$hook_file" ] || continue
        local hook_name=$(basename "$hook_file")
        local dest_file="$hooks_dest/$hook_name"

        # No-clobber: only copy if destination does NOT already exist
        if [ ! -f "$dest_file" ]; then
            cp "$hook_file" "$dest_file"
            chmod +x "$dest_file"
            log_info "  → Copied $hook_name"
        else
            log_info "  → Skipped $hook_name (already exists)"
        fi
    done

    log_success "Synced hooks to ~/.claude/hooks/"
}

# Install Engram MCP server for persistent memory
install_engram() {
    log_header "Installing Engram (Persistent Memory)"

    # Check if engram is already installed
    if command -v engram &> /dev/null; then
        log_info "Engram already installed: $(engram --version 2>/dev/null || echo 'version unknown')"
        setup_engram_for_opencode
        return 0
    fi

    # Install via Homebrew
    if command -v brew &> /dev/null; then
        log_info "Installing engram via Homebrew..."
        if brew install gentleman-programming/tap/engram; then
            log_success "Engram installed successfully"
            setup_engram_for_opencode
        else
            log_error "Failed to install engram via Homebrew"
            log_info "You can manually install from: https://github.com/Gentleman-Programming/engram"
            return 1
        fi
    else
        log_warning "Homebrew not found. Cannot install engram automatically."
        log_info "Please install Homebrew first: https://brew.sh"
        log_info "Or manually install engram from: https://github.com/Gentleman-Programming/engram"
        return 1
    fi
}

# Setup engram for OpenCode
setup_engram_for_opencode() {
    log_info "Setting up engram for OpenCode..."

    if engram setup opencode; then
        log_success "Engram configured for OpenCode"
        log_info ""
        log_info "Next steps:"
        log_info "  1. Restart OpenCode to load the engram plugin"
        log_info "  2. Run 'engram serve &' to start the HTTP server for session tracking"
        log_info ""
        log_info "Engram provides persistent memory across sessions:"
        log_info "  • mem_save - Save observations"
        log_info "  • mem_search - Search memories"
        log_info "  • mem_context - Recover session context"
        log_info "  • mem_session_summary - Summarize sessions"
    else
        log_warning "Failed to configure engram for OpenCode"
        return 1
    fi
}

# Generate all formats for a single AGENTS.md
generate_all_for_file() {
    local agents_file="$1"

    generate_claude "$agents_file"
    generate_gemini "$agents_file"
    generate_copilot "$agents_file"
    generate_codex "$agents_file"
}

# Generate all formats for all AGENTS.md files
generate_all() {
    log_header "Generating All Formats"

    local agents_files=$(find_agents_files)

    if [ -z "$agents_files" ]; then
        log_error "No AGENTS.md files found in repository"
        exit 1
    fi

    for agents_file in $agents_files; do
        log_info "Processing: $agents_file"
        generate_all_for_file "$agents_file"
        echo ""
    done

    log_success "All formats generated!"
}

# ============================================================================
# Interactive Menu
# ============================================================================

show_menu() {
    log_header "Gentleman.Dots AI Skills Setup"

    echo "This script synchronizes AGENTS.md to tool-specific formats."
    echo "AGENTS.md is the single source of truth for all AI assistants."
    echo ""
    echo "Select which assistants to configure:"
    echo ""
    echo "  ${CYAN}1)${NC} Claude Code      (CLAUDE.md)"
    echo "  ${CYAN}2)${NC} Gemini CLI       (GEMINI.md)"
    echo "  ${CYAN}3)${NC} GitHub Copilot   (.github/copilot-instructions.md)"
    echo "  ${CYAN}4)${NC} OpenAI Codex     (CODEX.md)"
    echo "  ${CYAN}5)${NC} All of the above"
    echo ""
    echo "  ${CYAN}6)${NC} Sync to ~/.claude/skills/"
    echo "  ${CYAN}7)${NC} Sync to OpenCode config"
    echo "  ${CYAN}8)${NC} Sync to all user configs"
    echo "  ${CYAN}9)${NC} Sync hooks to ~/.claude/hooks/"
    echo "  ${CYAN}10)${NC} Install OpenCode orchestrator agents"
    echo "  ${CYAN}11)${NC} Install Copilot orchestrator"
    echo "  ${CYAN}12)${NC} Install Gemini orchestrator"
    echo "  ${CYAN}13)${NC} Install Codex orchestrator"
    echo "  ${CYAN}14)${NC} Install ALL orchestrators"
    echo "  ${CYAN}15)${NC} Install Engram (persistent memory)"
    echo ""
    echo "  ${CYAN}0)${NC} Exit"
    echo ""
    printf "Enter choice [0-15]: "
}

handle_menu_choice() {
    local choice="$1"
    local agents_file="$REPO_ROOT/AGENTS.md"

    if [ ! -f "$agents_file" ]; then
        log_error "AGENTS.md not found at $agents_file"
        exit 1
    fi

    case $choice in
        1)
            generate_claude "$agents_file"
            ;;
        2)
            generate_gemini "$agents_file"
            ;;
        3)
            generate_copilot "$agents_file"
            ;;
        4)
            generate_codex "$agents_file"
            ;;
        5)
            generate_all_for_file "$agents_file"
            ;;
        6)
            sync_claude_config
            ;;
        7)
            sync_opencode
            ;;
        8)
            sync_claude_config
            sync_opencode
            sync_hooks
            ;;
        9)
            sync_hooks
            ;;
        10)
            install_opencode_agents
            ;;
        11)
            install_copilot_orchestrator
            ;;
        12)
            install_gemini_orchestrator
            ;;
        13)
            install_codex_orchestrator
            ;;
        14)
            install_all_orchestrators
            ;;
        15)
            install_engram
            ;;
        0)
            log_info "Exiting..."
            exit 0
            ;;
        *)
            log_error "Invalid choice: $choice"
            exit 1
            ;;
    esac
}

interactive_menu() {
    show_menu
    read -r choice
    handle_menu_choice "$choice"
}

# ============================================================================
# CLI Argument Parsing
# ============================================================================

show_help() {
    cat << EOF
Gentleman.Dots AI Skills Setup

Usage: ./skills/setup.sh [OPTIONS]

Options:
  --claude      Generate CLAUDE.md from AGENTS.md
  --gemini      Generate GEMINI.md from AGENTS.md
  --copilot     Generate .github/copilot-instructions.md
  --codex       Generate CODEX.md from AGENTS.md
  --all         Generate all format-specific files
  --sync-claude Sync skills + hooks to ~/.claude/
  --sync-opencode Sync skills to OpenCode config
  --sync-hooks  Sync hooks to ~/.claude/hooks/
  --sync-all    Sync skills + hooks to all user config directories
  --install-opencode-agents  Install orchestrator agents to OpenCode
  --install-copilot-orchestrator  Install unified orchestrator for Copilot
  --install-gemini-orchestrator   Install unified orchestrator for Gemini
  --install-codex-orchestrator    Install unified orchestrator for Codex
  --install-all-orchestrators     Install all orchestrators
  --install-engram                Install Engram MCP server for persistent memory
  --help        Show this help message

Examples:
  ./skills/setup.sh              # Interactive menu
  ./skills/setup.sh --all        # Generate all formats
  ./skills/setup.sh --claude     # Claude Code only
  ./skills/setup.sh --sync-all   # Sync to user configs
  ./skills/setup.sh --install-opencode-agents  # Install OpenCode agents
  ./skills/setup.sh --install-all-orchestrators  # Install all orchestrators
  ./skills/setup.sh --install-engram  # Install Engram persistent memory
EOF
}

parse_args() {
    local agents_file="$REPO_ROOT/AGENTS.md"

    if [ ! -f "$agents_file" ]; then
        log_error "AGENTS.md not found at $agents_file"
        exit 1
    fi

    case "$1" in
        --claude)
            generate_claude "$agents_file"
            ;;
        --gemini)
            generate_gemini "$agents_file"
            ;;
        --copilot)
            generate_copilot "$agents_file"
            ;;
        --codex)
            generate_codex "$agents_file"
            ;;
        --all)
            generate_all_for_file "$agents_file"
            ;;
        --sync-claude)
            sync_claude_config
            sync_hooks
            ;;
        --sync-opencode)
            sync_opencode
            ;;
        --sync-hooks)
            sync_hooks
            ;;
        --sync-all)
            sync_claude_config
            sync_opencode
            sync_hooks
            ;;
        --install-opencode-agents)
            install_opencode_agents
            ;;
        --install-copilot-orchestrator)
            install_copilot_orchestrator
            ;;
        --install-gemini-orchestrator)
            install_gemini_orchestrator
            ;;
        --install-codex-orchestrator)
            install_codex_orchestrator
            ;;
        --install-all-orchestrators)
            install_all_orchestrators
            ;;
        --install-engram)
            install_engram
            ;;
        --help|-h)
            show_help
            exit 0
            ;;
        *)
            log_error "Unknown option: $1"
            show_help
            exit 1
            ;;
    esac
}

# ============================================================================
# Main
# ============================================================================

main() {
    cd "$REPO_ROOT"

    if [ $# -eq 0 ]; then
        interactive_menu
    else
        parse_args "$@"
    fi

    echo ""
    log_success "Done!"
}

main "$@"
