[CmdletBinding()]
param(
  [string]$AppName       = "service",
  [string]$CmdDir        = "cmd/service",
  [string]$Bin           = "bin/service.exe",
  [string]$BufGen        = "buf.gen.yaml",
  [string]$ConfigPath    = "./configs/config.yaml",
  [string]$ReleaseScript = "./scripts/ps/git-release.ps1"
)

$text = @"
=== $AppName - ayuda rapida de [make <cmd>] ===

=== RECOMENDACIONES ====================================================================

Flujo de comandos (RECOMENDADO)
  1.0)   make all     - (init -> gen (deps + gen) -> wire -> run)
  1.1)   make init    - instalar herramientas
  1.2)   make deps    - refrescar buf.lock si hace falta, se puede omitir y utilizar make gen
  1.3)   make gen     - generar codigo protobuf (plantilla: $BufGen) (deps + gen)
  1.4)   make wire    - generar inyeccion (wire) en $CmdDir
  2.0)   make build   - compilar binario ($Bin)
  3.0)   make run     - ejecutar con kratos (o 'make gorun')

Workflow de modulos (RECOMENDADO)
   make module name="foo" [ops="get,upsert,delete"]
    - Crea módulo completo (proto, feature, repo, biz, service, wire). Si ya existe, crea una nueva versión.
    - Si NO se pasan ops → SIN endpoints por defecto; en el .proto se añade un mock GET /vN/foo/mock   // TODO: Mock de endpoint
  make module-delete name="foo" - eliminar modulo foo (todas versiones)
  make module-delete name="foo" version="v2" - eliminar solo version v2 (puede ser cualquier versión existente)

Ejemplos (ops)
  make module name="city"                        -> sin endpoints (se añade mock /mock en .proto)
  make module name="city" ops="get,upsert"       -> con GET + UPSERT
  make module name="city" ops="delete"           -> solo DELETE
  make module-proto name="city" ops="get,upsert" -> solo .proto con GET + UPSERT

OPS (endpoints por defecto)
  - OPS es un parámetro opcional que se puede pasar a make module, make module-proto, make module-feature, make module-service, make module-repo, make module-biz, make module-wire
  
  1) get (alias: find, list, read)
      - Genera RPC Find* en .proto y archivos: s_find.go (service) y r_list.go (repo)

  2) upsert (alias: create, update)
      - Genera RPC Upsert* en .proto y archivos: s_upsert.go (service) y r_upsert.go (repo)

  3) delete (alias: del, remove)
      - Genera RPC Delete*ById en .proto y archivos: s_delete_by_id.go (service) y r_delete_by_id.go (repo)

  4) Si ops está vacío:
      - No se generan RPCs de negocio (GET/UPSERT/DELETE).
      - En .proto se crea un único endpoint de mock: GET /vN/<base>/mock (// TODO: Mock de endpoint) para registrar rutas HTTP.

  5) Cualquier token desconocido en ops se ignora; si todos son desconocidos, se trata como vacío y se genera el mock.

Commit + version automatica
  make commit t="Titulo" d="Descripcion"   - commit con version automatica (etiqueta desde app.version en $ConfigPath)

=== COMANDOS ====================================================================

Comandos principales
  make help     - esta ayuda
  make init     - instalar/actualizar generadores + go mod tidy
  make deps     - actualizar buf.lock si buf.yaml cambio
  make gen      - buf generate
  make wire     - ejecutar wire en $CmdDir
  make build    - compilar -> $Bin
  make run      - kratos run (hot reload)
  make gorun    - go run ./$CmdDir -conf ./configs
  make clean    - limpiar bin/, wire_gen.go, *.pb.go, swagger/openapi
  make docs     - regenerar solo documentacion (docs/, docs/openapi)

=== CONFIG ====================================================================
  
Config (Makefile)
  CMD_DIR=$CmdDir
  BIN=$Bin
  BUF_GEN=$BufGen
  CONFIG_PATH=$ConfigPath

==============================================================================
"@


function Write-WithYellowMake {
  param([string]$Line)
  $rx = [regex]'make\s+\w+'
  $makeMatches = $rx.Matches($Line)
  if ($makeMatches.Count -eq 0) { Write-Host $Line; return }

  $pos = 0
  foreach ($m in $makeMatches) {
    if ($m.Index -gt $pos) {
      Write-Host -NoNewline ($Line.Substring($pos, $m.Index - $pos))
    }
    Write-Host -NoNewline $m.Value -ForegroundColor Yellow
    $pos = $m.Index + $m.Length
  }
  if ($pos -lt $Line.Length) {
    Write-Host ($Line.Substring($pos))
  } else {
    Write-Host ""
  }
}

# salida con color suave para encabezados
foreach ($line in $text -split "`r?`n") {
  if ($line -match '^(===|Flujo de comandos|Recomendaciones|Comandos principales|Comandos de modulo|Workflow de modulos|Commit \+ version automatica|Config|COMANDOS =+|RECOMENDACIONES =+|CONFIG =+)') {
    Write-Host $line -ForegroundColor Cyan
  } else {
    Write-WithYellowMake $line
  }
}
