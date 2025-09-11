param(
  [string]$AppName = "service",
  [string]$CmdDir  = "cmd/service",
  [string]$Bin     = "bin/service.exe",
  [string]$BufGen  = "buf.gen.yaml",

  # coherentes con tu Makefile
  [string]$ConfigPath   = "./configs/config.yaml",
  [string]$ReleaseScript = "./scripts/ps/git-release.ps1",

  # idioma: "es" por defecto, "en" para inglés
  [string]$Lang = "es"
)

$Lang = ($Lang.Trim().ToLower())
if ($Lang -ne "en") { $Lang = "es" }

if ($Lang -eq "en") {
  $text = @"
Usage: make <target> [VAR=VALUE]

Targets:
  help         Show this help (multi-language via LANG=en|es)
  init         Install/upgrade codegen tools and run 'go mod tidy':
               - protoc-gen-go, protoc-gen-go-grpc
               - protoc-gen-go-http (Kratos v2)
               - protoc-gen-grpc-gateway, protoc-gen-openapiv2
               - protoc-gen-openapi (gnostic), wire
  deps         Update buf.lock only when needed (when buf.yaml is newer or lock is missing)
  gen          buf generate (LOCAL plugins; template: $BufGen)
  wire         Run 'wire' in $CmdDir (produces wire_gen.go)
  build        Build Linux binary via ./scripts/ps/build-linux.ps1 -> $Bin
  run          Hot reload with 'kratos run'
  gorun        Run app with 'go run ./$CmdDir -conf ./configs'
  clean        Remove bin/, wire_gen.go, *.pb.go, openapi.yaml, service.swagger.json

Commit & Auto Version Bump:
  Driven by ${ReleaseScript}:
    - Reads base version ONLY from $ConfigPath (app.version like v1, v2).
    - Tags use vX.N (shared counter from 'origin').
    - Tag message equals commit message (Title + Desc).
    - Branch policy: forbidden on 'main'/'master' and detached HEAD.
    - First push sets upstream if missing; retries tag on collision.

Usage:
  make commit t="Your title" d="Your description"
  # alias:
  make release t="..." d="..."

Config:
  APP_NAME = $AppName
  CMD_DIR  = $CmdDir
  BIN      = $Bin
  BUF_GEN  = $BufGen
  CONFIG_PATH    = $ConfigPath
  RELEASE_SCRIPT = $ReleaseScript

Notes & Tips:
  - Local BUF plugins (no cloud quotas). Ensure GOBIN/GOPATH\bin is in PATH:
      go, go-grpc, go-http, openapiv2, openapi
  - 'deps' avoids hitting the network on every build (updates lock only when needed).
  - 'run' requires 'kratos' CLI installed separately.
  - If you see undefined 'BindForm' in generated code:
      go get github.com/go-kratos/kratos/v2@<ver>
      go install github.com/go-kratos/kratos/v2/cmd/protoc-gen-go-http@<ver>
"@
} else {
  $text = @"
Uso: make <target> [VAR=VALOR]

Targets:
  help         Muestra esta ayuda (multi-idioma con LANG=en|es)
  init         Instala/actualiza herramientas de generación y ejecuta 'go mod tidy':
               - protoc-gen-go, protoc-gen-go-grpc
               - protoc-gen-go-http (Kratos v2)
               - protoc-gen-grpc-gateway, protoc-gen-openapiv2
               - protoc-gen-openapi (gnostic), wire
  deps         Actualiza buf.lock solo cuando hace falta (si buf.yaml es más reciente o falta el lock)
  gen          buf generate (plugins LOCALES; plantilla: $BufGen)
  wire         Ejecuta 'wire' en $CmdDir (genera wire_gen.go)
  build        Compila binario Linux con ./scripts/ps/build-linux.ps1 -> $Bin
  run          Recarga en caliente con 'kratos run'
  gorun        Ejecuta con 'go run ./$CmdDir -conf ./configs'
  clean        Borra bin/, wire_gen.go, *.pb.go, openapi.yaml, service.swagger.json

Commit & Auto Version Bump:
  Orquestado por ${ReleaseScript}:
    - Lee la versión base SOLO de $ConfigPath (app.version tipo v1, v2).
    - Tags vX.N (contador compartido desde 'origin').
    - Mensaje del tag = mensaje del commit (Title + Desc).
    - Política de ramas: prohibido en 'main'/'master' y en detached HEAD.
    - Primer push crea upstream si falta; reintenta tag si hay colisión.

Uso:
  make commit t="Tu título" d="Tu descripción"
  # alias:
  make release t="..." d="..."

Config:
  APP_NAME = $AppName
  CMD_DIR  = $CmdDir
  BIN      = $Bin
  BUF_GEN  = $BufGen
  CONFIG_PATH    = $ConfigPath
  RELEASE_SCRIPT = $ReleaseScript

Notas y Tips:
  - BUF con plugins locales (sin cuotas en la nube). Asegura GOBIN/GOPATH\bin en PATH:
      go, go-grpc, go-http, openapiv2, openapi
  - 'deps' evita tocar la red en cada build (actualiza el lock solo cuando hace falta).
  - 'run' requiere tener el CLI 'kratos' instalado aparte.
  - Si ves 'undefined BindForm' en el código generado:
      go get github.com/go-kratos/kratos/v2@<ver>
      go install github.com/go-kratos/kratos/v2/cmd/protoc-gen-go-http@<ver>
"@
}

Write-Host $text
