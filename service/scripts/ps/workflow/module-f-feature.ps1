# scripts/ps/workflow/module-feature.ps1
[CmdletBinding()]
param(
  [Parameter(Mandatory = $true)] [string]$Name
)

Set-StrictMode -Version Latest
$ErrorActionPreference = 'Stop'

$ApiRoot     = "./api"
$FeatureRoot = "./internal/feature"
$utf8NoBom   = New-Object System.Text.UTF8Encoding($false)

# --- utils (logs) ---
try { . (Join-Path $PSScriptRoot 'utils.ps1') } catch {
  function Show-Step { param([string]$Message) Write-Host "`n==> $Message" -ForegroundColor Cyan }
  function Show-Info { param([string]$Message) Write-Host "$Message" -ForegroundColor DarkGray }
  function Show-OK   { param([string]$Message) Write-Host "  [OK] $Message" -ForegroundColor Green }
  function Show-Warn { param([string]$Message) Write-Host "  [WARN] $Message" -ForegroundColor Yellow }
  function Show-ErrorAndExit { param([string]$Message) Write-Host "  [ERROR] $Message" -ForegroundColor Red; exit 1 }
}

function ConvertTo-PascalCase { param([string]$s)
  $parts = ($s -replace '[^A-Za-z0-9]+',' ') -split '\s+' | Where-Object { $_ }
  ($parts | ForEach-Object { $_.Substring(0,1).ToUpper() + $_.Substring(1).ToLower() }) -join ''
}
function ConvertTo-LowerCase { param([string]$s) $s.ToLower() }
function ConvertTo-ImportAlias { param([string]$s) ($s.ToLower() -replace '[^a-z0-9]','_') }

Show-Step "Generating feature module (register.go, wire.go)"

$base   = ConvertTo-LowerCase $Name            # e.g. "prueba"
$pascal = ConvertTo-PascalCase $Name           # "Prueba"
$pkg    = ($base -replace '[^a-z0-9]','_')
$alias  = ConvertTo-ImportAlias $base
Show-Info "Module: base='$base', pascal='$pascal', pkg='$pkg', alias='$alias'"

# --- latest API version (vN) ---
$apiBaseDir = Join-Path $ApiRoot $base
$apiV = 1
Show-Info "Detecting latest API version in '$apiBaseDir'..."
if (Test-Path $apiBaseDir) {
  $max = 0
  Get-ChildItem $apiBaseDir -Directory -ea SilentlyContinue | ForEach-Object {
    if ($_.Name -match '^v(\d+)$') { $n=[int]$Matches[1]; if ($n -gt $max){$max=$n} }
  }
  if ($max -gt 0) { $apiV = $max }
  Show-Info "Using API version: v$apiV"
} else {
  Show-Warn "API dir not found: $apiBaseDir ; using v1."
}

# feature path uses same vN as API
$featRootV = Join-Path (Join-Path $FeatureRoot $base) "v$apiV"
$svcDir  = Join-Path $featRootV "service"
$bizDir  = Join-Path $featRootV "biz"
$repoDir = Join-Path $featRootV "repo"

Show-Info "Ensuring feature directories:"
$null = New-Item -ItemType Directory -Force -Path $svcDir, $bizDir, $repoDir
Show-OK "Created/exists: $featRootV (service, biz, repo)"

$apiImport  = "service/api/$base/v$apiV"
$svcImport  = "service/internal/feature/$base/v$apiV/service"
$bizImport  = "service/internal/feature/$base/v$apiV/biz"
$repoImport = "service/internal/feature/$base/v$apiV/repo"

# -------------------------
# register.go (module-local types + registrars) — with version in service name
# -------------------------
$registerPath = Join-Path $featRootV "register.go"
if (-not (Test-Path $registerPath)) {
  Show-Info "Writing: $registerPath"
$registerGo = @"
package $pkg

import (
	api_$alias "$apiImport"
	${pkg}_service "$svcImport"

	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
)

// Module-local types to avoid wire type collisions
type HTTPRegister func(*http.Server)
type GRPCRegister func(*grpc.Server)

// NOTE: versioned service interfaces from proto: ${pascal}v${apiV}Service...
var _ api_$alias.${pascal}v${apiV}ServiceHTTPServer = (*${pkg}_service.${pascal}Service)(nil)
var _ api_$alias.${pascal}v${apiV}ServiceServer     = (*${pkg}_service.${pascal}Service)(nil)

func New${pascal}HTTPRegistrer(s api_$alias.${pascal}v${apiV}ServiceHTTPServer) HTTPRegister {
	return func(srv *http.Server) {
		api_$alias.Register${pascal}v${apiV}ServiceHTTPServer(srv, s)
	}
}

func New${pascal}GRPCRegistrer(s api_$alias.${pascal}v${apiV}ServiceServer) GRPCRegister {
	return func(srv *grpc.Server) {
		api_$alias.Register${pascal}v${apiV}ServiceServer(srv, s)
	}
}
"@
  [IO.File]::WriteAllText($registerPath, $registerGo, $utf8NoBom)
  Show-OK "Created: register.go"
} else {
  Show-Info "Skip (exists): $registerPath"
}

# -------------------------
# wire.go (ProviderSet + binds) — with version in service name
# -------------------------
$wirePath = Join-Path $featRootV "wire.go"
if (-not (Test-Path $wirePath)) {
  Show-Info "Writing: $wirePath"
$wireGo = @"
package $pkg

import (
	api_$alias "$apiImport"
	${pkg}_biz "$bizImport"
	${pkg}_repo "$repoImport"
	${pkg}_service "$svcImport"

	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	${pkg}_repo.New${pascal}Repo,
	${pkg}_biz.New${pascal}Usecase,
	${pkg}_service.New${pascal}Service,

	// map generated service interfaces (versioned) to our implementation
	wire.Bind(new(api_$alias.${pascal}v${apiV}ServiceHTTPServer), new(*${pkg}_service.${pascal}Service)),
	wire.Bind(new(api_$alias.${pascal}v${apiV}ServiceServer),     new(*${pkg}_service.${pascal}Service)),

	// module-local registrars
	New${pascal}HTTPRegistrer,
	New${pascal}GRPCRegistrer,
)
"@
  [IO.File]::WriteAllText($wirePath, $wireGo, $utf8NoBom)
  Show-OK "Created: wire.go"
} else {
  Show-Info "Skip (exists): $wirePath"
}

Show-OK ("feature module generated: {0}/v{1}  (register.go, wire.go)" -f $base, $apiV)
