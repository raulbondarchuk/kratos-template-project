param(
  [string]$AppName = "service",
  [string]$CmdDir  = "cmd/service",
  [string]$Bin     = "bin/service.exe",
  [string]$BufGen  = "buf.gen.yaml",

  # coherentes con tu Makefile
  [string]$ConfigPath = "./configs/config.yaml",
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
  init         Install/upgrade toolchain, then vendor deps (offline):
               - Installs: buf, protoc-gen-go, protoc-gen-go-grpc, protoc-gen-go-http,
                 protoc-gen-openapiv2, protoc-gen-openapi (gnostic).
               - Calls: make third_party

  third_party  Vendor minimal proto deps via sparse clone (offline-friendly):
               - third_party/googleapis: google/api, google/rpc, google/type
               - third_party/grpc-gateway: protoc-gen-openapiv2/options
               (Add paths via sparse-checkout if you import more .proto)

  gen          Generate code with local plugins (no BSR calls):
               - buf generate (uses $BufGen)

  wire         Generate DI with Google Wire:
               - runs Wire in $CmdDir (produces wire_gen.go)

  build        Build binary -> $Bin
  run          kratos run (hot reload)
  gorun        go run ./$CmdDir -conf ./configs
  clean        Remove bin/ and all wire_gen.go

Commit & Auto Version Bump:
  Driven by ${ReleaseScript}:
    - Reads the base major ONLY from $ConfigPath (app.version like v1, v2). You never pass versions.
    - Tags are vX.N; next N is computed from origin (shared counter across devs).
    - Tag message equals commit message (Title + Desc).
    - Branch policy:
        * Forbidden on 'main' and 'master'.
        * Forbidden in detached HEAD.
        * If no upstream, first push sets it (-u origin HEAD:<branch>).
    - Push flow: push branch, then push tag (auto-retries with next vX.N on collision).
    - If no changes, exits cleanly.

Usage:
  make commit t="Your title" d="Your description"
  (compat: TITLE/DESC still accepted if t/d are not provided)
  make release t="..." d="..."   # alias to commit

Examples:
  make commit t="feat: ingest LTA" d="parse \$MSG:11"
  make build APP_NAME=myapp

Config:
  APP_NAME = $AppName
  CMD_DIR  = $CmdDir
  BIN      = $Bin
  BUF_GEN  = $BufGen
  CONFIG_PATH    = $ConfigPath
  RELEASE_SCRIPT = $ReleaseScript

Tips:
  - Offline setup: ensure buf.yaml lists local modules:
        version: v2
        modules:
          - path: .
          - path: third_party/googleapis
          - path: third_party/grpc-gateway
  - Local plugins only: buf.gen.yaml should use 'local:' plugins.
  - PATH: make sure $(go env GOPATH)\bin is in PATH for this session.
  - Old Git? Run 'git -C third_party/... sparse-checkout init --cone' before 'set'.
"@
} else {
  $text = @"
Uso: make <target> [VAR=VALOR]

Targets:
  init         Instala/actualiza toolchain y luego hace vendor (offline):
               - Instala: buf, protoc-gen-go, protoc-gen-go-grpc, protoc-gen-go-http,
                 protoc-gen-openapiv2, protoc-gen-openapi (gnostic).
               - Llama: make third_party

  third_party  Vendor mínimo de dependencias proto con sparse clone (offline):
               - third_party/googleapis: google/api, google/rpc, google/type
               - third_party/grpc-gateway: protoc-gen-openapiv2/options
               (Añade rutas en sparse-checkout si importas más .proto)

  gen          Generación con plugins locales (sin llamadas a BSR):
               - buf generate (usa $BufGen)

  wire         Genera DI con Google Wire:
               - ejecuta Wire en $CmdDir (crea wire_gen.go)

  build        Compila binario -> $Bin
  run          kratos run (hot reload)
  gorun        go run ./$CmdDir -conf ./configs
  clean        Elimina bin/ y todos los wire_gen.go

Commit & Auto Version Bump:
  Orquestado por ${ReleaseScript}:
    - Lee la versión base SOLO desde $ConfigPath (app.version tipo v1, v2). No pasas versiones.
    - Etiquetas vX.N; el siguiente N se calcula mirando origin (contador compartido).
    - El mensaje del tag = mensaje del commit (Title + Desc).
    - Política de ramas:
        * Prohibido en 'main' y 'master'.
        * Prohibido en detached HEAD.
        * Si no hay upstream, el primer push lo crea (-u origin HEAD:<rama>).
    - Flujo: push de la rama y luego push del tag (reintenta con siguiente vX.N si hay colisión).
    - Si no hay cambios, sale sin error.

Uso:
  make commit t="Tu título" d="Tu descripción"
  (compat: TITLE/DESC siguen funcionando si no pasas t/d)
  make release t="..." d="..."   # alias de commit

Ejemplos:
  make commit t="feat: ingest LTA" d="parsear \$MSG:11"
  make build APP_NAME=myapp

Config:
  APP_NAME = $AppName
  CMD_DIR  = $CmdDir
  BIN      = $Bin
  BUF_GEN  = $BufGen
  CONFIG_PATH    = $ConfigPath
  RELEASE_SCRIPT = $ReleaseScript

Tips:
  - Modo offline: asegúrate de que buf.yaml tenga módulos locales:
        version: v2
        modules:
          - path: .
          - path: third_party/googleapis
          - path: third_party/grpc-gateway
  - Solo plugins locales: buf.gen.yaml debe usar 'local:'.
  - PATH: añade $(go env GOPATH)\bin al PATH de la sesión.
  - Git antiguo: ejecuta 'git -C third_party/... sparse-checkout init --cone' antes de 'set'.
"@
}

Write-Host $text
