# =========================
# Windows PowerShell–only Makefile (Kratos + buf)
# =========================

SHELL := powershell.exe
.SHELLFLAGS := -NoProfile -ExecutionPolicy Bypass -Command

# --- Project layout ---
APP_NAME := service
APP_NAME_TO_BUILD := kratos-template
CMD_DIR  := cmd/$(APP_NAME)
BIN_DIR  := bin

BUF_GEN ?= buf.gen.yaml
BUF_YAML ?= buf.yaml
BUF_LOCK ?= buf.lock

RELEASE_SCRIPT ?= ./scripts/ps/git-release.ps1
CONFIG_PATH    ?= ./configs/config.yaml

# compatibility: if t/d are not set, use TITLE/DESC (if set)
TITLE_EFF = $(if $(strip $(t)),$(t),$(TITLE))
DESC_EFF  = $(if $(strip $(d)),$(d),$(DESC))

# ---------------------------------
# Help - print all commands and their descriptions
# ---------------------------------

HELP_PS := powershell -NoProfile -ExecutionPolicy Bypass -File "./scripts/ps/help.ps1"

.PHONY: help
help:
	@$(HELP_PS) -AppName "$(APP_NAME)" -CmdDir "$(CMD_DIR)" -Bin "$(BIN_DIR)" -BufGen "$(BUF_GEN)" -ConfigPath "$(CONFIG_PATH)" -ReleaseScript "$(RELEASE_SCRIPT)"

# PHONY targets
.PHONY: all init deps gen wire run gorun clean commit release 

# All targets
all: init gen wire run

# Tools setup
init: 
	@powershell -NoProfile -ExecutionPolicy Bypass -File ./scripts/ps/make/init.ps1 -AppName "$(APP_NAME)"

# Generate dependencies (buf.lock)
deps:
	@powershell -NoProfile -ExecutionPolicy Bypass -File ./scripts/ps/make/deps.ps1 -BufYaml "$(BUF_YAML)" -BufLock "$(BUF_LOCK)"

# Generate protobuf code
gen: deps
	@powershell -NoProfile -ExecutionPolicy Bypass -File ./scripts/ps/make/gen.ps1 -BufGen "$(BUF_GEN)"

# Wire (wire_gen.go)
wire:
	@powershell -NoProfile -ExecutionPolicy Bypass -File ./scripts/ps/make/wire.ps1 -CmdDir "$(CMD_DIR)"

# Build
build:
	@powershell -NoProfile -ExecutionPolicy Bypass -File ./scripts/ps/build-linux.ps1 -AppName '$(APP_NAME_TO_BUILD)' -CmdDir '$(CMD_DIR)' -BinDir '$(BIN_DIR)'

# Kratos run
run:
	@powershell -NoProfile -ExecutionPolicy Bypass -File ./scripts/ps/make/run.ps1 -Mode kratos

# Go run
gorun:
	@powershell -NoProfile -ExecutionPolicy Bypass -File ./scripts/ps/make/run.ps1 -Mode go -CmdDir "$(CMD_DIR)" -ConfigPath "$(CONFIG_PATH)"

# Clean
clean clear cls:
	@powershell -NoProfile -ExecutionPolicy Bypass -File ./scripts/ps/make/clean.ps1 -BinDir "$(BIN_DIR)"

# Commit with auto version bump
# Usage: make commit t="..." d="..."
commit:
	@if (-not '$(TITLE_EFF)' -or -not '$(DESC_EFF)') { Write-Host 'Uso: make commit t="..." d="..."'; exit 2 } ; if (-not (Test-Path '$(RELEASE_SCRIPT)')) { Write-Error "Script not found: $(RELEASE_SCRIPT)"; exit 1 } ; & '$(RELEASE_SCRIPT)' -Title '$(TITLE_EFF)' -Desc '$(DESC_EFF)' -ConfigPath '$(CONFIG_PATH)'


# ==============================================================================
# Workflow =====================================================================
# ==============================================================================

PATH_SCRIPTS_WORKFLOW_MODULE := ./scripts/ps/workflow/module.ps1

# Generate a new endpoint skeleton
.PHONY: module

# Create a new endpoint (only .proto scaffolding)
# Example:
#   make module name="prueba"
module:
	@if (-not '$(name)') { Write-Host 'Usage: make module name="<name>"'; exit 2 }
	@powershell -NoProfile -ExecutionPolicy Bypass -File "$(PATH_SCRIPTS_WORKFLOW_MODULE)" -Name "$(name)"
