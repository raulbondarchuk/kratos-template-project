# scripts/ps/workflow/feature.ps1
[CmdletBinding()]
param([Parameter(Mandatory = $true)] [string]$Name)

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

# names
$base      = ConvertTo-LowerCase $Name
$pkgBase   = ($base -replace '[^a-z0-9]','_')
$pascal    = ConvertTo-PascalCase $Name
$alias     = ConvertTo-ImportAlias $base

# detect latest api vN
$apiBaseDir = Join-Path $ApiRoot $base
$apiVersion = 1
if (Test-Path -LiteralPath $apiBaseDir) {
  $max = 0
  Get-ChildItem -LiteralPath $apiBaseDir -Directory | ForEach-Object {
    if ($_.Name -match '^v(\d+)$') { $n=[int]$Matches[1]; if ($n -gt $max){$max=$n} }
  }
  if ($max -gt 0) { $apiVersion = $max }
}
$apiImport     = "service/api/$base/v$apiVersion"
$featureRootV  = Join-Path (Join-Path $FeatureRoot $base) "v$apiVersion"
$svcDirImport  = "service/internal/feature/$base/v$apiVersion/service"
$bizDirImport  = "service/internal/feature/$base/v$apiVersion/biz"
$repoDirImport = "service/internal/feature/$base/v$apiVersion/repo"

# ensure dirs
$null = New-Item -ItemType Directory -Force -Path $featureRootV

# register.go
$registerPath = Join-Path $featureRootV "register.go"
if (-not (Test-Path $registerPath)) {
$register = @"
package $pkgBase

import (
	api_$alias "$apiImport"
	${pkgBase}_service "$svcDirImport"
	server_grpc "service/internal/server/grpc"
	server_http "service/internal/server/http"

	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
)

var _ api_$alias.${pascal}ServiceHTTPServer = (*${pkgBase}_service.${pascal}Service)(nil)
var _ api_$alias.${pascal}ServiceServer     = (*${pkgBase}_service.${pascal}Service)(nil)

func New${pascal}HTTPRegister(s api_$alias.${pascal}ServiceHTTPServer) server_http.HTTPRegister {
	return func(srv *http.Server) {
		api_$alias.Register${pascal}ServiceHTTPServer(srv, s)
	}
}

func New${pascal}GRPCRegister(s api_$alias.${pascal}ServiceServer) server_grpc.GRPCRegister {
	return func(srv *grpc.Server) {
		api_$alias.Register${pascal}ServiceServer(srv, s)
	}
}
"@
[IO.File]::WriteAllText($registerPath, $register, $utf8NoBom)
}

# wire.go
$wirePath = Join-Path $featureRootV "wire.go"
if (-not (Test-Path $wirePath)) {
$wire = @"
package $pkgBase

import (
	api_$alias "$apiImport"
	${pkgBase}_biz "$bizDirImport"
	${pkgBase}_repo "$repoDirImport"
	${pkgBase}_service "$svcDirImport"

	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	${pkgBase}_repo.New${pascal}Repo,
	${pkgBase}_biz.New${pascal}Usecase,
	${pkgBase}_service.New${pascal}Service,

	wire.Bind(new(api_$alias.${pascal}ServiceHTTPServer), new(*${pkgBase}_service.${pascal}Service)),
	wire.Bind(new(api_$alias.${pascal}ServiceServer),     new(*${pkgBase}_service.${pascal}Service)),

	New${pascal}HTTPRegister,
	New${pascal}GRPCRegister,
)
"@
[IO.File]::WriteAllText($wirePath, $wire, $utf8NoBom)
}

Write-Host ("feature: {0}/v{1} (register.go, wire.go)" -f $base, $apiVersion) -ForegroundColor Green