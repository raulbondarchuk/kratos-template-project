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
BIN      := $(BIN_DIR)/$(APP_NAME).exe
BUF_GEN  := buf.gen.yaml

RELEASE_SCRIPT ?= ./scripts/ps/git-release.ps1
CONFIG_PATH    ?= ./configs/config.yaml

# ---------------------------------
# Help - print all commands and their descriptions
# ---------------------------------

HELP_PS := powershell -NoProfile -ExecutionPolicy Bypass -File "./scripts/ps/help.ps1"

.PHONY: help
help:
	@$(HELP_PS) -AppName "$(APP_NAME)" -CmdDir "$(CMD_DIR)" -Bin "$(BIN)" -BufGen "$(BUF_GEN)" -ConfigPath "$(CONFIG_PATH)" -ReleaseScript "$(RELEASE_SCRIPT)" -Lang "$(LANG)"

# ---------------------------------
# Tools setup
# ---------------------------------
.PHONY: init
init:
	@echo ">> installing tools"
	go install github.com/bufbuild/buf/cmd/buf@latest
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/go-kratos/kratos/cmd/protoc-gen-go-http/v2@latest
	go install github.com/google/gnostic/cmd/protoc-gen-openapi@latest
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest
	go install github.com/google/wire/cmd/wire@latest
	go install github.com/go-kratos/kratos/cmd/kratos/v2@latest
	go mod tidy

# ---------------------------------
# Go mod tidy
# ---------------------------------
.PHONY: tidy
tidy:
	@echo ">> go mod tidy"
	go mod tidy

# ---------------------------------
# Generate (proto + wire)
# ---------------------------------
.PHONY: gen
gen:
	@echo ">> buf generate"
	buf dep update
	buf generate --template $(BUF_GEN)
	@echo ">> wire generate"
	Set-Location '$(CMD_DIR)'; wire

# ---------------------------------
# Build
# ---------------------------------

.PHONY: build
build:
	@powershell -NoProfile -ExecutionPolicy Bypass -File ./scripts/ps/build-linux.ps1 -AppName '$(APP_NAME_TO_BUILD)' -CmdDir '$(CMD_DIR)' -BinDir '$(BIN_DIR)'


# ---------------------------------
# Run (go run, kratos run)
# ---------------------------------
.PHONY: run krun
run:
	@echo ">> running with go run"
	go run ./$(CMD_DIR) -conf ./configs

krun:
	@echo ">> running with kratos run"
	kratos run

# ---------------------------------
# Clean
# ---------------------------------
.PHONY: clean
clean:
	@echo ">> clean"
	if (Test-Path '$(BIN_DIR)') { Remove-Item -Recurse -Force '$(BIN_DIR)' }
	Get-ChildItem -Recurse -Filter 'wire_gen.go' | Remove-Item -Force -ErrorAction SilentlyContinue

# ---------------------------------
# Commit with auto version bump
# ---------------------------------

# compatibility: if t/d are not set, use TITLE/DESC (if set)
TITLE_EFF = $(if $(strip $(t)),$(t),$(TITLE))
DESC_EFF  = $(if $(strip $(d)),$(d),$(DESC))

.PHONY: commit
# uso: make commit t="..." d="..."
commit:
	@if (-not '$(TITLE_EFF)' -or -not '$(DESC_EFF)') { Write-Host 'Uso: make commit t="..." d="..."'; exit 2 } ; if (-not (Test-Path '$(RELEASE_SCRIPT)')) { Write-Error "Script not found: $(RELEASE_SCRIPT)"; exit 1 } ; & '$(RELEASE_SCRIPT)' -Title '$(TITLE_EFF)' -Desc '$(DESC_EFF)' -ConfigPath '$(CONFIG_PATH)'
# alias
.PHONY: release
release:
	@$(MAKE) commit t="$(t)" d="$(d)"
%:
	@: