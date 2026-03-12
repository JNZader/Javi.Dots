# Javi.Dots

> 🍴 **Fork de [Gentleman.Dots](https://github.com/Gentleman-Programming/Gentleman.Dots)** por [Gentleman Programming](https://youtube.com/@GentlemanProgramming). Las actualizaciones del upstream se mergean regularmente.

> ℹ️ **Actualización (enero 2026)**: OpenCode ahora soporta suscripciones **Claude Max/Pro** mediante el plugin `opencode-anthropic-auth` (incluido en esta configuración).
> Tanto **Claude Code** como **OpenCode** funcionan con tu suscripción de Claude.
> *Nota: este workaround es estable por ahora, pero Anthropic podría bloquearlo en el futuro.*

📄 Leer en: [English](README.md) | **Español**

## Tabla de Contenidos

* [¿Qué es esto?](#qué-es-esto)
* [Inicio rápido](#inicio-rápido)
* [Plataformas soportadas](#plataformas-soportadas)
* [🤖 Herramientas IA y Framework](#-herramientas-ia-y-framework)
* [🎮 Entrenador de Maestría en Vim](#-entrenador-de-maestría-en-vim)
* [Documentación](#documentación)
* [Resumen de herramientas](#resumen-de-herramientas)
* [Estructura del proyecto](#estructura-del-proyecto)
* [Soporte](#soporte)

---

## Vista previa

```
                  ██╗███╗   ██╗███████╗
                  ██║████╗  ██║╚══███╔╝
                  ██║██╔██╗ ██║  ███╔╝
             ██   ██║██║╚██╗██║ ███╔╝
             ╚█████╔╝██║ ╚████║███████╗
              ╚════╝ ╚═╝  ╚═══╝╚══════╝

    ██╗ █████╗ ██╗   ██╗██╗    ██████╗  ██████╗ ████████╗███████╗
    ██║██╔══██╗██║   ██║██║    ██╔══██╗██╔═══██╗╚══██╔══╝██╔════╝
    ██║███████║██║   ██║██║    ██║  ██║██║   ██║   ██║   ███████╗
██  ██║██╔══██║╚██╗ ██╔╝██║    ██║  ██║██║   ██║   ██║   ╚════██║
╚████╔╝██║  ██║ ╚████╔╝ ██║ ██╗██████╔╝╚██████╔╝   ██║   ███████║
 ╚═══╝ ╚═╝  ╚═╝  ╚═══╝  ╚═╝ ╚═╝╚═════╝  ╚═════╝    ╚═╝   ╚══════╝
```

---

## ¿Qué es esto?

Una configuración completa de entorno de desarrollo que incluye:

* **Neovim** con LSP, autocompletado y asistentes de IA (Claude Code, Gemini, OpenCode)
* **Zed** editor con modo Vim y soporte para agentes IA
* **Herramientas IA**: Claude Code, OpenCode, Gemini CLI, GitHub Copilot, Codex CLI, Qwen Code con configs, skills y temas
* **Framework IA**: 199 módulos (72 agentes, 85 skills, 10 hooks, 20 comandos, 10 servidores MCP) + 6 orquestadores de dominio + 36 skills curados, con selección por preset o personalizada
* **Shells**: Fish, Zsh, Nushell
* **Multiplexores de terminal**: Tmux, Zellij
* **Emuladores de terminal**: Alacritty, WezTerm, Kitty, Ghostty

---

## Inicio rápido

### Opción 1: Homebrew (Recomendado)

```bash
brew install JNZader/tap/javi-dots
javi-dots
```

### Opción 2: Descarga directa

```bash
# macOS Apple Silicon
curl -fsSL https://github.com/JNZader/Javi.Dots/releases/latest/download/gentleman-installer-darwin-arm64 -o gentleman.dots

# macOS Intel
curl -fsSL https://github.com/JNZader/Javi.Dots/releases/latest/download/gentleman-installer-darwin-amd64 -o gentleman.dots

# Linux x86_64
curl -fsSL https://github.com/JNZader/Javi.Dots/releases/latest/download/gentleman-installer-linux-amd64 -o gentleman.dots

# Linux ARM64 (Raspberry Pi, etc.)
curl -fsSL https://github.com/JNZader/Javi.Dots/releases/latest/download/gentleman-installer-linux-arm64 -o gentleman.dots

# Luego ejecutar
chmod +x gentleman.dots
./gentleman.dots
```

### Opción 3: Termux (Android)

Termux requiere compilar el instalador localmente (la cross-compilación de Go a Android tiene limitaciones).

```bash
# 1. Instalar dependencias
pkg update && pkg upgrade
pkg install git golang

# 2. Clonar el repositorio
git clone https://github.com/JNZader/Javi.Dots.git
cd Javi.Dots/installer

# 3. Compilar y ejecutar
go build -o ~/gentleman-installer ./cmd/gentleman-installer
cd ~
./gentleman-installer
```

| Soporte en Termux                 | Estado                                               |
| --------------------------------- | ---------------------------------------------------- |
| Shells (Fish, Zsh, Nushell)       | ✅ Disponible                                         |
| Multiplexores (Tmux, Zellij)      | ✅ Disponible                                         |
| Neovim con configuración completa | ✅ Disponible                                         |
| Nerd Fonts                        | ✅ Instaladas automáticamente en `~/.termux/font.ttf` |
| Emuladores de terminal            | ❌ No aplica                                          |
| Homebrew                          | ❌ Usa `pkg`                                          |

> **Tip:** Después de la instalación, reiniciá Termux para aplicar la fuente y luego ejecutá `tmux` o `zellij` para iniciar el entorno configurado.

El instalador TUI te guía para seleccionar tus herramientas preferidas y maneja toda la configuración automáticamente.

> **Usuarios de Windows:** primero debés configurar WSL. Ver la [Guía de instalación manual](docs/manual-installation.md#windows-wsl).

---

## Plataformas soportadas

| Plataforma            | Arquitectura          | Método de instalación       | Gestor de paquetes |
| --------------------- | --------------------- | --------------------------- | ------------------ |
| macOS                 | Apple Silicon (ARM64) | Homebrew, descarga directa  | Homebrew           |
| macOS                 | Intel (x86_64)        | Homebrew, descarga directa  | Homebrew           |
| Linux (Ubuntu/Debian) | x86_64, ARM64         | Homebrew, descarga directa  | Homebrew           |
| Linux (Fedora/RHEL)   | x86_64, ARM64         | Descarga directa            | dnf                |
| Linux (Arch)          | x86_64                | Homebrew, descarga directa  | Homebrew           |
| Windows               | WSL                   | Descarga directa (ver docs) | Homebrew           |
| Android               | Termux (ARM64)        | Compilación local           | pkg                |

---

## 🤖 Herramientas IA y Framework

El instalador incluye un sistema completo de integración con IA (Pasos 8-9):

### Herramientas IA (Paso 8)

Selección múltiple de 6 herramientas de IA (con botón Seleccionar Todo):

| Herramienta | Qué se instala |
|-------------|---------------|
| **Claude Code** | Binario + CLAUDE.md + persona Gentleman + 10+ skills + tema Kanagawa |
| **OpenCode** | Binario + agente Gentleman + 6 orquestadores de dominio + orquestador SDD + tema |
| **Gemini CLI** | CLI vía npm |
| **GitHub Copilot** | Extensión gh |
| **Codex CLI** | Binario vía npm + config AGENTS.md |
| **Qwen Code** | Binario vía npm + QWEN.md + settings.json |

### Orquestadores de Dominio

OpenCode incluye **6 orquestadores de dominio** que reemplazan el selector plano de 72 agentes por una experiencia manejable de 9 agentes:

| Orquestador | Agentes | Ejemplos |
|------------|:-------:|---------|
| `development-orchestrator` | 22 | React Pro, Go Pro, Backend Architect |
| `quality-orchestrator` | 11 | Code Reviewer, Security Auditor, E2E Specialist |
| `infrastructure-orchestrator` | 7 | Cloud Architect, Kubernetes Expert, DevOps |
| `data-ai-orchestrator` | 6 | AI Engineer, Data Scientist, MLOps |
| `business-orchestrator` | 8 | Project Manager, API Designer, Product Strategist |
| `workflow-orchestrator` | 16 | Plan Executor, Wave Executor, Code Migrator |

Elegí un orquestador de dominio → este rutea al especialista correcto. No más scrollear entre 72+ agentes.

### Framework IA (Paso 9)

Elegí un preset o personalizá entre **199 módulos** en 6 categorías:

| Categoría | Módulos | Ejemplos |
|-----------|--------:|---------|
| 🪝 Hooks | 10 | Secret Scanner, Commit Guard, Model Router |
| ⚡ Comandos | 20 | Git Commit, PR Review, TDD, Refactoring |
| 🤖 Agentes | 72 | React Pro, DevOps Engineer, Security Auditor |
| 🎯 Skills | 85 | FastAPI, Spring Boot 4, Kubernetes, PyTorch |
| 📐 SDD | 2 | OpenSpec, Agent Teams Lite |
| 🔌 MCP | 10 | Context7, Engram, Jira, Atlassian, Figma, Notion, Brave Search, Sentry, Cloudflare, VoiceMode |

**Presets**: Minimal, Frontend, Backend, Fullstack, Data, Complete

**Elección SDD**: Instalá [OpenSpec](https://github.com/JNZader/project-starter-framework) (SDD basado en archivos), [Agent Teams Lite](https://github.com/Gentleman-Programming/agent-teams-lite) (SDD liviano con 9 sub-agentes), o ambos.

**Scroll con viewport**: Las listas largas (Skills: 85, Agents: 72) scrollean dentro de la terminal con indicadores `▲`/`▼`.

---

## 🎮 Entrenador de Maestría en Vim

¡Aprendé Vim de forma divertida! El instalador incluye un entrenador interactivo estilo RPG con:

| Módulo                   | Teclas cubiertas                         |
| ------------------------ | ---------------------------------------- |
| 🔤 Movimiento horizontal | `w`, `e`, `b`, `f`, `t`, `0`, `$`, `^`   |
| ↕️ Movimiento vertical   | `j`, `k`, `G`, `gg`, `{`, `}`            |
| 📦 Objetos de texto      | `iw`, `aw`, `i"`, `a(`, `it`, `at`       |
| ✂️ Cambiar y repetir     | `d`, `c`, `dd`, `cc`, `D`, `C`, `x`      |
| 🔄 Sustitución           | `r`, `R`, `s`, `S`, `~`, `gu`, `gU`, `J` |
| 🎬 Macros y registros    | `qa`, `@a`, `@@`, `"ay`, `"+p`           |
| 🔍 Regex / Búsqueda      | `/`, `?`, `n`, `N`, `*`, `#`, `\v`       |

Cada módulo incluye 15 lecciones progresivas, modo práctica con selección inteligente de ejercicios, jefes finales y seguimiento de XP.

Podés iniciarlo desde el menú principal: **Vim Mastery Trainer**

---

## 📦 Inicialización de Proyectos

Bootstrapeá cualquier proyecto con soporte de framework IA:

```bash
# Interactivo
javi-dots  # → Menú Principal → Initialize Project

# No interactivo
javi-dots --non-interactive --init-project \
  --project-path=/ruta/al/proyecto \
  --project-memory=obsidian-brain \
  --project-ci=github --project-engram
```

**Módulos de memoria**: Obsidian Brain, VibeKanban, Engram, Simple, None
**Proveedores CI**: GitHub Actions, GitLab CI, Woodpecker, None

---

## 🎯 Gestor de Skills

Navegá, instalá y eliminá skills de agentes IA del catálogo Gentleman-Skills:

```bash
# Interactivo
javi-dots  # → Menú Principal → Skill Manager

# No interactivo
javi-dots --non-interactive --skill-install=react-19,typescript,tailwind-4
javi-dots --non-interactive --skill-remove=react-19
```

Los skills se organizan por categoría (curated, community, plugin) y se enlazan a `~/.claude/skills/`.

---

## 🔀 Soporte para Forks

Sobrescribí la URL de clone y el directorio para apuntar a tu propio fork:

```bash
# Vía variables de entorno
REPO_URL=https://github.com/TuUsuario/TuFork.git REPO_DIR=TuFork javi-dots

# Vía flags CLI
javi-dots --repo-url=https://github.com/TuUsuario/TuFork.git --repo-dir=TuFork
```

---

## Documentación

| Documento                                                          | Descripción                                                  |
| ------------------------------------------------------------------ | ------------------------------------------------------------ |
| [Guía del instalador TUI](docs/tui-installer.md)                   | Funciones interactivas, navegación, backup y restore        |
| [Herramientas IA y Framework](docs/ai-tools-integration.md)        | Selección de IA, presets, drill-down por categoría, flags CLI |
| [Módulos del Framework IA](docs/ai-framework-modules.md)           | Referencia completa de los 199 módulos en 6 categorías      |
| [Agent Teams Lite](docs/agent-teams-lite.md)                       | Framework SDD liviano con 9 sub-agentes                     |
| [Configuración de IA](docs/ai-configuration.md)                    | Claude Code, OpenCode, Copilot y más                        |
| [Instalación manual](docs/manual-installation.md)                  | Configuración paso a paso para todas las plataformas        |
| [Keymaps de Neovim](docs/neovim-keymaps.md)                        | Referencia completa de atajos                               |
| [Especificación del entrenador Vim](docs/vim-trainer-spec.md)      | Detalles técnicos del entrenador                            |
| [Testing con Docker](docs/docker-testing.md)                       | Tests E2E con contenedores                                  |
| [Contribuir](docs/contributing.md)                                 | Setup de desarrollo, sistema de skills y releases           |

---

## Resumen de herramientas

### Emuladores de terminal

| Herramienta   | Descripción                                  |
| ------------- | -------------------------------------------- |
| **Ghostty**   | Acelerado por GPU, nativo y ultra rápido     |
| **Kitty**     | Rico en funcionalidades, renderizado por GPU |
| **WezTerm**   | Configurable con Lua, multiplataforma        |
| **Alacritty** | Minimalista, escrito en Rust                 |

### Shells

| Herramienta | Descripción                                |
| ----------- | ------------------------------------------ |
| **Nushell** | Datos estructurados y pipelines modernos   |
| **Fish**    | Amigable y con excelentes defaults         |
| **Zsh**     | Altamente personalizable, compatible POSIX |

### Multiplexores

| Herramienta | Descripción                           |
| ----------- | ------------------------------------- |
| **Tmux**    | Probado en batalla, ampliamente usado |
| **Zellij**  | Moderno, plugins WebAssembly          |

### Editores

| Herramienta | Descripción                             |
| ----------- | --------------------------------------- |
| **Neovim**  | Config LazyVim con LSP, completado e IA |
| **Zed**     | Editor de alto rendimiento con modo Vim y soporte IA |

### Prompts

| Herramienta  | Descripción                            |
| ------------ | -------------------------------------- |
| **Starship** | Prompt multi-shell con integración Git |

---

## Estructura del proyecto

```
Javi.Dots/
├── installer/               # Instalador TUI en Go
│   ├── cmd/                 # Punto de entrada
│   ├── internal/            # TUI, sistema y entrenador
│   └── e2e/                 # Tests E2E con Docker
├── docs/                    # Documentación
├── openspec/                # Artefactos de Spec-Driven Development
├── skills/                  # Skills de agentes IA
│
├── GentlemanNvim/           # Configuración Neovim
├── GentlemanClaude/         # Config Claude Code + skills
├── GentlemanOpenCode/       # Config OpenCode
├── GentlemanQwen/           # Config Qwen Code
├── GentlemanZed/            # Config Zed (modo Vim + IA)
│
├── GentlemanFish/
├── GentlemanZsh/
├── GentlemanNushell/
├── GentlemanTmux/
├── GentlemanZellij/
│
├── GentlemanGhostty/
├── GentlemanKitty/
├── alacritty.toml
├── .wezterm.lua
│
└── starship.toml
```

---

## Soporte

* **Issues**: [GitHub Issues](https://github.com/JNZader/Javi.Dots/issues)
* **Upstream**: [Gentleman.Dots](https://github.com/Gentleman-Programming/Gentleman.Dots) por Gentleman Programming
* **Discord**: [Gentleman Programming Community](https://discord.gg/gentleman-programming)
* **YouTube**: [@GentlemanProgramming](https://youtube.com/@GentlemanProgramming)

---

## Licencia

Licencia MIT — libre de usar, modificar y compartir.

**¡Feliz coding!** 🎩
