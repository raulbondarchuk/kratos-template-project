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

function ConvertTo-PascalCase { param([string]$s)
  $parts = ($s -replace '[^A-Za-z0-9]+',' ') -split '\s+' | Where-Object { $_ }
  ($parts | ForEach-Object { $_.Substring(0,1).ToUpper() + $_.Substring(1).ToLower() }) -join ''
}
function ConvertTo-LowerCase { param([string]$s) $s.ToLower() }
function ConvertTo-ImportAlias { param([string]$s) ($s.ToLower() -replace '[^a-z0-9]','_') }

$base   = ConvertTo-LowerCase $Name            # e.g. "prueba"
$pascal = ConvertTo-PascalCase $Name           # "Prueba"
$pkg    = ($base -replace '[^a-z0-9]','_')
$alias  = ConvertTo-ImportAlias $base

# --- latest API version (vN) ---
$apiBaseDir = Join-Path $ApiRoot $base
$apiV = 1
if (Test-Path $apiBaseDir) {
  $max = 0
  Get-ChildItem $apiBaseDir -Directory -ea SilentlyContinue | ForEach-Object {
    if ($_.Name -match '^v(\d+)$') { $n=[int]$Matches[1]; if ($n -gt $max){$max=$n} }
  }
  if ($max -gt 0) { $apiV = $max }
} else {
  Write-Warning "API dir not found: $apiBaseDir ; using v1."
}

# feature path uses same vN as API
$featRootV = Join-Path (Join-Path $FeatureRoot $base) "v$apiV"
$svcDir = Join-Path $featRootV "service"
$bizDir = Join-Path $featRootV "biz"
$repoDir = Join-Path $featRootV "repo"
$null = New-Item -ItemType Directory -Force -Path $svcDir,$bizDir,$repoDir

$apiImport    = "service/api/$base/v$apiV"
$svcImport    = "service/internal/feature/$base/v$apiV/service"
$bizImport    = "service/internal/feature/$base/v$apiV/biz"
$repoImport   = "service/internal/feature/$base/v$apiV/repo"

# -------------------------
# register.go (module-local types + registrars)
# -------------------------
$registerPath = Join-Path $featRootV "register.go"
if (-not (Test-Path $registerPath)) {
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

var _ api_$alias.${pascal}ServiceHTTPServer = (*${pkg}_service.${pascal}Service)(nil)
var _ api_$alias.${pascal}ServiceServer     = (*${pkg}_service.${pascal}Service)(nil)

func New${pascal}HTTPRegistrer(s api_$alias.${pascal}ServiceHTTPServer) HTTPRegister {
	return func(srv *http.Server) {
		api_$alias.Register${pascal}ServiceHTTPServer(srv, s)
	}
}

func New${pascal}GRPCRegistrer(s api_$alias.${pascal}ServiceServer) GRPCRegister {
	return func(srv *grpc.Server) {
		api_$alias.Register${pascal}ServiceServer(srv, s)
	}
}
"@
[IO.File]::WriteAllText($registerPath, $registerGo, $utf8NoBom)
}

# -------------------------
# wire.go (ProviderSet + binds)
# -------------------------
$wirePath = Join-Path $featRootV "wire.go"
if (-not (Test-Path $wirePath)) {
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

	// map generated service interfaces to our implementation
	wire.Bind(new(api_$alias.${pascal}ServiceHTTPServer), new(*${pkg}_service.${pascal}Service)),
	wire.Bind(new(api_$alias.${pascal}ServiceServer),     new(*${pkg}_service.${pascal}Service)),

	// module-local registrars
	New${pascal}HTTPRegistrer,
	New${pascal}GRPCRegistrer,
)
"@
[IO.File]::WriteAllText($wirePath, $wireGo, $utf8NoBom)
}

Write-Host ("feature module: {0}/v{1}  (register.go, wire.go)" -f $base, $apiV) -ForegroundColor Green
