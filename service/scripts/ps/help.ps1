param(
  [string]$AppName = "service",
  [string]$CmdDir  = "cmd/service",
  [string]$Bin     = "bin/service.exe",
  [string]$BufGen  = "buf.gen.yaml",

  # coherentes con tu Makefile
  [string]$ConfigPath = "./configs/config.yaml",
  [string]$ReleaseScript = "./scripts/ps/git-release.ps1",

  # idioma: "es" por defecto, "en" para ingles
  [string]$Lang = "es"
)

$Lang = ($Lang.Trim().ToLower())
if ($Lang -ne "en") { $Lang = "es" }

if ($Lang -eq "en") {
  $text = @"
Usage: make <target> [VAR=VALUE]

Targets:
  init         Install/upgrade toolchain and run go mod tidy:
               - buf, protoc-gen-go, protoc-gen-go-grpc, protoc-gen-go-http (Kratos),
                 protoc-gen-openapi (gnostic), grpc-gateway, wire, kratos CLI
  gen          Generate code:
               - buf generate (using $BufGen)
               - wire (in $CmdDir, produces wire_gen.go)
  build        Build binary -> $Bin
  run          go run ./$CmdDir -conf ./configs
  krun         kratos run (hot reload mode)
  tidy         go mod tidy
  clean        Remove bin/ and all wire_gen.go

Commit & Auto Version Bump:
  Driven by ${ReleaseScript}:
    - Reads base version ONLY from $ConfigPath (app.version like v1, v2). You never pass versions.
    - Tags are vX.N (1,2,3...). The next patch number is computed from **origin** (git ls-remote),
      so all developers share the same counter.
    - Tag message equals the commit message (Title + Desc).
    - Branch policy:
        * Forbidden from 'main' and 'master'.
        * Forbidden in detached HEAD.
        * If no upstream, first push sets it (-u origin HEAD:<branch>).
    - Push flow:
        1) push current branch,
        2) push the tag with collision retry (auto-picks next vX.N if name is taken).
    - If there are no changes, exits cleanly.

Usage (full syntax only):
  make commit t="Your title" d="Your description"
  (compat: TITLE/DESC still accepted if t/d are not provided)

Also available:
  make release t="..." d="..."   (alias to commit)

Examples:
  make commit t="prueba commit" d="Esto es una prueba commit"
  make commit t="feat: ingest LTA" d="support `\$MSG:11"     # escape $ with backtick in PowerShell

Config:
  APP_NAME = $AppName
  CMD_DIR  = $CmdDir
  BIN      = $Bin
  BUF_GEN  = $BufGen

  Commit config:
  CONFIG_PATH    = $ConfigPath
  RELEASE_SCRIPT = $ReleaseScript

Tips:
  - If proto imports fail: buf dep update
  - OpenAPI outputs are configured in $BufGen (e.g., docs/…)
  - Override app name:  make build APP_NAME=myapp
  - Extra go flags:     make build GOFLAGS=-trimpath
"@
} else {
  $text = @"
Uso: make <target> [VAR=VALOR]

Targets:
  init         Instala/actualiza toolchain y ejecuta go mod tidy:
               - buf, protoc-gen-go, protoc-gen-go-grpc, protoc-gen-go-http (Kratos),
                 protoc-gen-openapi (gnostic), grpc-gateway, wire, kratos CLI
  gen          Generacion de codigo:
               - buf generate (con $BufGen)
               - wire (en $CmdDir, produce wire_gen.go)
  build        Compila binario -> $Bin
  run          go run ./$CmdDir -conf ./configs
  krun         kratos run (hot reload)
  tidy         go mod tidy
  clean        Borra bin/ y todos los wire_gen.go

Commit & Auto Version Bump:
  Orquestado por ${ReleaseScript}:
    - Lee la version base SOLO desde $ConfigPath (app.version tipo v1, v2). Nunca pasas version.
    - Etiquetas vX.N (1,2,3...). El siguiente numero se calcula mirando **origin** (git ls-remote),
      para que todos usen el mismo contador.
    - El mensaje del tag es el mismo que el mensaje del commit (Title + Desc).
    - Politica de ramas:
        * Prohibido en 'main' y 'master'.
        * Prohibido en detached HEAD.
        * Si no hay upstream, el primer push lo crea (-u origin HEAD:<rama>).
    - Flujo de push:
        1) push de la rama actual,
        2) push del tag con reintento si hay colision (elige automaticamente el siguiente vX.N).
    - Si no hay cambios, sale sin error.

Uso (solo modo completo):
  make commit t="Tu titulo" d="Tu descripcion"
  (compat: TITLE/DESC siguen funcionando si no pasas t/d)

Alias:
  make release t="..." d="..."   (alias de commit)

Ejemplos:
  make commit t="prueba commit" d="Esto es una prueba commit"
  make commit t="feat: ingest LTA" d="support `\$MSG:11"     # escapa $ con acento grave en PowerShell

Config:
  APP_NAME = $AppName
  CMD_DIR  = $CmdDir
  BIN      = $Bin
  BUF_GEN  = $BufGen

  Config de commit:
  CONFIG_PATH    = $ConfigPath
  RELEASE_SCRIPT = $ReleaseScript

Tips:
  - Si fallan imports proto: buf dep update
  - Salidas OpenAPI se configuran en $BufGen (p. ej., docs/…)
  - Cambiar nombre app:  make build APP_NAME=myapp
  - Flags extra go:     make build GOFLAGS=-trimpath
"@
}

Write-Host $text
