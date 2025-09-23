# scripts/ps/workflow/service.ps1
[CmdletBinding()]
param(
  [Parameter(Mandatory = $true)] [string]$Name,
  [string]$Ops = ""
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
function ConvertTo-Plural    { param([string]$s) if ($s.ToLower().EndsWith('s')){$s}else{"$s"+"s"} }
function ConvertTo-ImportAlias { param([string]$s) ($s.ToLower() -replace '[^a-z0-9]','_') }

# --- parse ops ---
$opsList = @()
if ($Ops) { $opsList = $Ops.ToLower().Split(',') | ForEach-Object { $_.Trim() } | Where-Object { $_ } }
$HasGet = $false; $HasUpsert = $false; $HasDelete = $false
foreach($op in $opsList){
  switch ($op) {
    'get' { $HasGet = $true } 'find' { $HasGet = $true } 'list' { $HasGet = $true } 'read' { $HasGet = $true }
    'upsert' { $HasUpsert = $true } 'create' { $HasUpsert = $true } 'update' { $HasUpsert = $true }
    'delete' { $HasDelete = $true } 'del' { $HasDelete = $true } 'remove' { $HasDelete = $true }
    default { Show-Warn "Unknown op '$op' ignored" }
  }
}

Show-Step "Generating service layer"

$base         = ConvertTo-LowerCase $Name
$pkgBase      = ($base -replace '[^a-z0-9]','_')
$pascal       = ConvertTo-PascalCase $Name
$pluralPascal = ConvertTo-PascalCase (ConvertTo-Plural $base)
$alias        = ConvertTo-ImportAlias $base

Show-Info "Module: base='$base', pascal='$pascal', pkg='$pkgBase', alias='$alias', ops=[$Ops]"

# latest api vN
$apiBaseDir = Join-Path $ApiRoot $base
$apiVersion = 1
Show-Info "Detecting latest API version in '$apiBaseDir'..."
if (Test-Path $apiBaseDir) {
  $max=0; Get-ChildItem $apiBaseDir -Directory -ErrorAction SilentlyContinue | ForEach-Object {
    if ($_.Name -match '^v(\d+)$'){ $n=[int]$Matches[1]; if($n -gt $max){$max=$n} }
  }
  if ($max -gt 0){ $apiVersion=$max }
  Show-Info "Using API version: v$apiVersion"
} else {
  Show-Warn "API dir not found: $apiBaseDir ; using v1."
}

$featureRootV = Join-Path (Join-Path $FeatureRoot $base) "v$apiVersion"
$svcDir = Join-Path $featureRootV "service"
$null = New-Item -ItemType Directory -Force -Path $svcDir
Show-Info "Service dir ready: $svcDir"

$apiImport     = "service/api/$base/v$apiVersion"
$bizImport     = "service/internal/feature/$base/v$apiVersion/biz"
$metaImport    = "service/internal/server/http/meta"
$serviceName   = "${pascal}v${apiVersion}Service"

# ---------- service.go (always)
$p = Join-Path $svcDir "service.go"
if (-not (Test-Path $p)) {
  Show-Info "Writing: $p"
$txt = @"
package ${pkgBase}_service

import (
	api_$alias "$apiImport"
	${pkgBase}_biz "$bizImport"
)

type ${pascal}Service struct {
	api_$alias.Unimplemented${serviceName}Server
	uc *${pkgBase}_biz.${pascal}Usecase
}

func New${pascal}Service(uc *${pkgBase}_biz.${pascal}Usecase) *${pascal}Service {
	return &${pascal}Service{uc: uc}
}
"@
  [IO.File]::WriteAllText($p, $txt, $utf8NoBom)
  Show-OK "Created: service.go"
} else { Show-Info "Skip (exists): $p" }

# ---------- s_find.go (GET)
if ($HasGet) {
  $p = Join-Path $svcDir "s_find.go"
  if (-not (Test-Path $p)) {
    Show-Info "Writing: $p"
$txt = @"
package ${pkgBase}_service

import (
	"context"

	api_$alias   "$apiImport"
	${pkgBase}_biz "$bizImport"
	srvmeta      "$metaImport"
	"service/pkg/converter"
	"service/pkg/generic"
)

func (s *${pascal}Service) Find${pluralPascal}(ctx context.Context, req *api_$alias.Find${pluralPascal}Request) (*api_$alias.Find${pluralPascal}Response, error) {
	// presence-aware (optional fields)
	var idPtr *uint
	var namePtr *string

	if req.Id != nil && *req.Id != 0 {
	  v := uint(*req.Id)
	  idPtr = &v
	}
	if req.Name != nil && *req.Name != "" {
	  v := *req.Name
	  namePtr = &v
	}

	bizRes, err := s.uc.Find${pluralPascal}(ctx, idPtr, namePtr)
	if err != nil {
		return &api_$alias.Find${pluralPascal}Response{
			Meta: srvmeta.WithDetails(srvmeta.MetaInternal("failed to find ${base}"), map[string]string{"error": err.Error()}),
		}, nil
	}

	if len(bizRes) == 0 {
		return &api_$alias.Find${pluralPascal}Response{
			Meta: srvmeta.MetaNoContent("no items"),
		}, nil
	}

	dto, err := generic.ToDTOSliceGeneric[${pkgBase}_biz.${pascal}, api_$alias.${pascal}](bizRes)
	if err != nil {
		return &api_$alias.Find${pluralPascal}Response{
			Meta: srvmeta.WithDetails(srvmeta.MetaInternal("failed to marshal dto"), map[string]string{"error": err.Error()}),
		}, nil
	}

	for i := range bizRes {
		dto[i].CreatedAt = converter.ConvertToGoogleTimestamp(bizRes[i].CreatedAt)
		dto[i].UpdatedAt = converter.ConvertToGoogleTimestamp(bizRes[i].UpdatedAt)
	}

	return &api_$alias.Find${pluralPascal}Response{
		Items: generic.ToPointerSliceGeneric(dto),
		Meta:  srvmeta.MetaOK("OK"),
	}, nil
}
"@
    [IO.File]::WriteAllText($p, $txt, $utf8NoBom)
    Show-OK "Created: s_find.go"
  } else { Show-Info "Skip (exists): $p" }
} else {
  Show-Info "No GET op; skip s_find.go"
}

# ---------- s_upsert.go (UPSERT)
if ($HasUpsert) {
  $p = Join-Path $svcDir "s_upsert.go"
  if (-not (Test-Path $p)) {
    Show-Info "Writing: $p"
$txt = @"
package ${pkgBase}_service

import (
	"context"

	api_$alias   "$apiImport"
	${pkgBase}_biz "$bizImport"
	srvmeta      "$metaImport"
	"service/pkg/converter"
	"service/pkg/generic"
)

func (s *${pascal}Service) Upsert${pascal}(ctx context.Context, req *api_$alias.Upsert${pascal}Request) (*api_$alias.Upsert${pascal}Response, error) {
	in := &${pkgBase}_biz.${pascal}{
		ID:   uint(req.GetId()),
		Name: req.GetName(),
	}
	res, err := s.uc.Upsert${pascal}(ctx, in)
	if err != nil {
		return &api_$alias.Upsert${pascal}Response{
			Meta: srvmeta.WithDetails(srvmeta.MetaInternal("failed to upsert ${base}"), map[string]string{"error": err.Error()}),
		}, nil
	}

	dto, err := generic.ToDTOGeneric[${pkgBase}_biz.${pascal}, api_$alias.${pascal}](*res)
	if err != nil {
		return &api_$alias.Upsert${pascal}Response{
			Meta: srvmeta.WithDetails(srvmeta.MetaInternal("failed to marshal dto"), map[string]string{"error": err.Error()}),
		}, nil
	}
	dto.CreatedAt = converter.ConvertToGoogleTimestamp(res.CreatedAt)
	dto.UpdatedAt = converter.ConvertToGoogleTimestamp(res.UpdatedAt)

	return &api_$alias.Upsert${pascal}Response{
		Item: &dto,
		Meta: srvmeta.MetaOK("OK"),
	}, nil
}
"@
    [IO.File]::WriteAllText($p, $txt, $utf8NoBom)
    Show-OK "Created: s_upsert.go"
  } else { Show-Info "Skip (exists): $p" }
} else {
  Show-Info "No UPSERT op; skip s_upsert.go"
}

# ---------- s_delete_by_id.go (DELETE)
if ($HasDelete) {
  $p = Join-Path $svcDir "s_delete_by_id.go"
  if (-not (Test-Path $p)) {
    Show-Info "Writing: $p"
$txt = @"
package ${pkgBase}_service

import (
	"context"

	api_$alias "$apiImport"
	srvmeta    "$metaImport"
)

func (s *${pascal}Service) Delete${pascal}ById(ctx context.Context, req *api_$alias.Delete${pascal}ByIdRequest) (*api_$alias.Delete${pascal}ByIdResponse, error) {
	if err := s.uc.Delete${pascal}ById(ctx, uint(req.GetId())); err != nil {
		return &api_$alias.Delete${pascal}ByIdResponse{
			Meta: srvmeta.WithDetails(srvmeta.MetaInternal("failed to delete ${base}"), map[string]string{"error": err.Error()}),
		}, nil
	}
	return &api_$alias.Delete${pascal}ByIdResponse{
		Meta: srvmeta.MetaOK("OK"),
	}, nil
}
"@
    [IO.File]::WriteAllText($p, $txt, $utf8NoBom)
    Show-OK "Created: s_delete_by_id.go"
  } else { Show-Info "Skip (exists): $p" }
} else {
  Show-Info "No DELETE op; skip s_delete_by_id.go"
}

Show-OK ("service generated: {0}/v{1}" -f $base, $apiVersion)
