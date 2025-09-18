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

== RECOMENDACIONES ====================================================================

Flujo de comandos (RECOMENDADO)
  1.0)   make all     - (init -> gen (deps + gen) -> wire -> run)
  1.1)   make init    - instalar herramientas
  1.2)   make deps    - refrescar buf.lock si hace falta, se puede omitir y utilizar make gen
  1.3)   make gen     - generar codigo protobuf (plantilla: $BufGen) (deps + gen)
  1.4)   make wire    - generar inyeccion (wire) en $CmdDir
  2.0)   make build   - compilar binario ($Bin)
  3.0)   make run     - ejecutar con kratos (o 'make gorun')

Workflow de modulos (RECOMENDADO)
  make module name="foo" - crear modulo completo (proto, feature, repo, biz, service, wire) En caso de que ya exista, se crea una version nueva.
  make module-delete name="foo" - eliminar modulo foo (todas versiones)
  make module-delete name="foo" version="v2" - eliminar solo version v2 (puede ser cualquier versión existente)

Commit + version automatica
  make commit t="Titulo" d="Descripcion"   - commit con version automatica (etiqueta desde app.version en $ConfigPath)

== COMANDOS ====================================================================

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

== CONFIG ====================================================================
  
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
