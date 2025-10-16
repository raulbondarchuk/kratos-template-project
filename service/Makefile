# =========================
# Windows PowerShell–only Makefile (Kratos + buf)
# =========================

SHELL := powershell.exe
.SHELLFLAGS := -NoProfile -ExecutionPolicy Bypass -Command

# --- Project layout ---
APP_NAME := service
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
.PHONY: all init deps gen wire wire-check run gorun clean commit release wire-ignore

PATH_SCRIPTS_PS := ./scripts/ps
PATH_SCRIPTS_MAKE := ./scripts/ps/make

# All targets
all: init gen wire run

# Tools setup
init: 
	@powershell -NoProfile -ExecutionPolicy Bypass -File "$(PATH_SCRIPTS_MAKE)/init.ps1" -AppName "$(APP_NAME)"

# Generate dependencies (buf.lock)
deps:
	@powershell -NoProfile -ExecutionPolicy Bypass -File "$(PATH_SCRIPTS_MAKE)/deps.ps1" -BufYaml "$(BUF_YAML)" -BufLock "$(BUF_LOCK)"

# Generate protobuf code
gen: deps
	@powershell -NoProfile -ExecutionPolicy Bypass -File "$(PATH_SCRIPTS_MAKE)/gen.ps1" -BufGen "$(BUF_GEN)"

# Wire (wire_gen.go)
wire: 
	make wire-check
	@powershell -NoProfile -ExecutionPolicy Bypass -File "$(PATH_SCRIPTS_MAKE)/wire.ps1" -CmdDir "$(CMD_DIR)"

wire-check:
	@powershell -NoProfile -ExecutionPolicy Bypass -File "$(PATH_SCRIPTS_MAKE)/wire-check.ps1"

# Kratos run
run:
	@powershell -NoProfile -ExecutionPolicy Bypass -File "$(PATH_SCRIPTS_MAKE)/run.ps1" -Mode kratos

# Go run
gorun:
	@powershell -NoProfile -ExecutionPolicy Bypass -File "$(PATH_SCRIPTS_MAKE)/run.ps1" -Mode go -CmdDir "$(CMD_DIR)" -ConfigPath "$(CONFIG_PATH)"

# ==============================================================================
# Docs =========================================================================
# ==============================================================================
PATH_SCRIPTS_MAKE_DOCS := $(PATH_SCRIPTS_MAKE)/docs.ps1
.PHONY: docs
docs:
	@powershell -NoProfile -ExecutionPolicy Bypass -File "$(PATH_SCRIPTS_MAKE_DOCS)"

# ==============================================================================
# Commit =====================================================================
# ==============================================================================

# Commit with auto version bump
# Usage: make commit t="..." d="..."
commit:
	@if (-not '$(TITLE_EFF)' -or -not '$(DESC_EFF)') { Write-Host 'Uso: make commit t="..." d="..."'; exit 2 } ; if (-not (Test-Path '$(RELEASE_SCRIPT)')) { Write-Error "Script not found: $(RELEASE_SCRIPT)"; exit 1 } ; & '$(RELEASE_SCRIPT)' -Title '$(TITLE_EFF)' -Desc '$(DESC_EFF)' -ConfigPath '$(CONFIG_PATH)'


# ==============================================================================
# Workflow =====================================================================
# ==============================================================================

PATH_SCRIPTS_WORKFLOW_MODULE_PROTO := ./scripts/ps/workflow/module-proto.ps1
PATH_SCRIPTS_WORKFLOW_MODULE_WIRE := ./scripts/ps/workflow/module-wire.ps1

PATH_SCRIPTS_WORKFLOW_MODULE_F_FEATURE := ./scripts/ps/workflow/module-f-feature.ps1
PATH_SCRIPTS_WORKFLOW_MODULE_F_SERVICE := ./scripts/ps/workflow/module-f-service.ps1
PATH_SCRIPTS_WORKFLOW_MODULE_F_REPO := ./scripts/ps/workflow/module-f-repo.ps1
PATH_SCRIPTS_WORKFLOW_MODULE_F_BIZ := ./scripts/ps/workflow/module-f-biz.ps1

PATH_SCRIPTS_WORKFLOW_MODULE_DELETE := ./scripts/ps/workflow/module-delete.ps1

# Generate a new module skeleton
.PHONY: module module-proto module-feature module-service module-repo module-biz module-wire module-delete
# Create a new module (only .proto scaffolding)
# Example:
# make module name="city" → sin endpoints
# make module name="city" ops="get,upsert" → con GET+UPSERT
# make module name="city" ops="delete" → solo DELETE
module:
	make module-proto   		name="$(name)" $(if $(strip $(ops)),ops="$(ops)",)
	make gen
	make module-feature 		name="$(name)"
	make module-repo    		name="$(name)" $(if $(strip $(ops)),ops="$(ops)",)
	make module-biz     		name="$(name)" $(if $(strip $(ops)),ops="$(ops)",)
	make module-service 		name="$(name)" $(if $(strip $(ops)),ops="$(ops)",)
	make module-wire    name="$(name)"
	make wire
	
	

# Generate a new module .proto in API skeleton
module-proto:
	@if (-not '$(name)') { Write-Host 'Usage: make module-proto name="<name>" [ops="get,upsert,delete"]'; exit 2 }
	@powershell -NoProfile -ExecutionPolicy Bypass -File "$(PATH_SCRIPTS_WORKFLOW_MODULE_PROTO)" -Name "$(name)" $(if $(strip $(ops)),-Ops "$(ops)",)

# Generate a new module feature in API skeleton
module-feature:
	@if (-not '$(name)') { Write-Host 'Usage: make module-feature name="<name>"'; exit 2 }
	@powershell -NoProfile -ExecutionPolicy Bypass -File "$(PATH_SCRIPTS_WORKFLOW_MODULE_F_FEATURE)" -Name "$(name)"

# Generate a new module service in API skeleton
module-service:
	@if (-not '$(name)') { Write-Host 'Usage: make module-service name="<name>" [ops="get,upsert,delete"]'; exit 2 }
	@powershell -NoProfile -ExecutionPolicy Bypass -File "$(PATH_SCRIPTS_WORKFLOW_MODULE_F_SERVICE)" -Name "$(name)" $(if $(strip $(ops)),-Ops "$(ops)",)

# Generate a new module repo in API skeleton
module-repo:
	@if (-not '$(name)') { Write-Host 'Usage: make module-repo name="<name>" [ops="get,upsert,delete"]'; exit 2 }
	@powershell -NoProfile -ExecutionPolicy Bypass -File "$(PATH_SCRIPTS_WORKFLOW_MODULE_F_REPO)" -Name "$(name)" $(if $(strip $(ops)),-Ops "$(ops)",)


# Generate a new module biz in API skeleton
module-biz:
	@if (-not '$(name)') { Write-Host 'Usage: make module-biz name="<name>" [ops="get,upsert,delete"]'; exit 2 }
	@powershell -NoProfile -ExecutionPolicy Bypass -File "$(PATH_SCRIPTS_WORKFLOW_MODULE_F_BIZ)" -Name "$(name)" $(if $(strip $(ops)),-Ops "$(ops)",)


# Generate a new module wire in API skeleton
module-wire:
	@if (-not '$(name)') { Write-Host 'Usage: make module-wire name="<name>"'; exit 2 }
	@powershell -NoProfile -ExecutionPolicy Bypass -File "$(PATH_SCRIPTS_WORKFLOW_MODULE_WIRE)" -Name "$(name)"

# Delete a module (all versions or a specific version with version="vN")
# Examples:
#   make module-delete name="prueba"
#   make module-delete name="prueba" version="v2"
module-delete:
	@if (-not '$(name)') { Write-Host 'Usage: make module-delete name="<name>" [version="vN"]'; exit 2 }
	@powershell -NoProfile -ExecutionPolicy Bypass -File "$(PATH_SCRIPTS_WORKFLOW_MODULE_DELETE)" -Name "$(name)" $(if $(strip $(version)),-Version $(version),)
	make wire
	make docs

PATH_SCRIPTS_WORKFLOW_MODULE_TESTS_SERVICE := ./scripts/ps/workflow/module-tests-service.ps1

.PHONY: tmodule module-tests-service tmodule-delete module-tests-delete
PATH_SCRIPTS_WORKFLOW_MODULE_TESTS_DELETE := ./scripts/ps/workflow/module-tests-delete.ps1

# make module-tests name="city" version="v2" force=1
# version=v2 (optional) - specific version, otherwise latest
# force=1 (optional) - overwrite if exists
tmodule:
	make module-tests-service name="$(name)" $(if $(strip $(version)),version="$(version)",) $(if $(strip $(force)),force="$(force)",)
	
module-tests-service:
	@if (-not '$(name)') { Write-Host 'Usage: make module-tests name="<name>" [version="v2"] [force=1]'; exit 2 }
	@powershell -NoProfile -ExecutionPolicy Bypass -File "$(PATH_SCRIPTS_WORKFLOW_MODULE_TESTS_SERVICE)" -Name "$(name)" $(if $(strip $(version)),-Version "$(version)",) $(if $(strip $(force)),-Force,)

# make tmodule-delete name="city" version="v2"
# version=v2 (optional) - specific version, otherwise all versions
tmodule-delete:
	make module-tests-delete name="$(name)" $(if $(strip $(version)),version="$(version)",)

module-tests-delete:
	@if (-not '$(name)') { Write-Host 'Usage: make module-tests-delete name="<name>" [version="v2"]'; exit 2 }
	@powershell -NoProfile -ExecutionPolicy Bypass -File "$(PATH_SCRIPTS_WORKFLOW_MODULE_TESTS_DELETE)" -Name "$(name)" $(if $(strip $(version)),-Version "$(version)",)
