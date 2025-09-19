# scripts/ps/workflow/repo.ps1
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

Show-Step "Generating repo layer"

$base         = ConvertTo-LowerCase $Name
$pkgBase      = ($base -replace '[^a-z0-9]','_')
$pascal       = ConvertTo-PascalCase $Name
$pluralPascal = ConvertTo-PascalCase (ConvertTo-Plural $base)

Show-Info "Module: base='$base', pascal='$pascal', pkg='$pkgBase', ops=[$Ops]"

# latest api vN -> feature vN dir
$apiBaseDir = Join-Path $ApiRoot $base
$apiVersion = 1
Show-Info "Detecting latest API version in '$apiBaseDir'..."
if (Test-Path $apiBaseDir) {
  $max=0
  Get-ChildItem $apiBaseDir -Directory -ErrorAction SilentlyContinue | ForEach-Object {
    if ($_.Name -match '^v(\d+)$'){ $n=[int]$Matches[1]; if($n -gt $max){$max=$n} }
  }
  if ($max -gt 0){ $apiVersion=$max }
  Show-Info "Using API version: v$apiVersion"
} else {
  Show-Warn "API dir not found: $apiBaseDir ; using v1."
}

$featureRootV = Join-Path (Join-Path $FeatureRoot $base) "v$apiVersion"
$repoDir = Join-Path $featureRootV "repo"
$null = New-Item -ItemType Directory -Force -Path $repoDir
Show-Info "Repo dir ready: $repoDir"

# repo.go (always)
$p = Join-Path $repoDir "repo.go"
if (-not (Test-Path $p)) {
  Show-Info "Writing: $p"
$txt = @"
package ${pkgBase}_repo

import (
	"service/internal/data"
	${pkgBase}_biz "service/internal/feature/$base/v$apiVersion/biz"

	"github.com/go-kratos/kratos/v2/log"
)

type ${base}Repo struct {
	data *data.Data
	log  *log.Helper
}

func New${pascal}Repo(data *data.Data, logger log.Logger) ${pkgBase}_biz.${pascal}Repo {
	return &${base}Repo{data: data, log: log.NewHelper(logger)}
}
"@
  [IO.File]::WriteAllText($p, $txt, $utf8NoBom)
  Show-OK "Created: repo.go"
} else { Show-Info "Skip (exists): $p" }

# r_list.go (GET)
if ($HasGet) {
  $p = Join-Path $repoDir "r_list.go"
  if (-not (Test-Path $p)) {
    Show-Info "Writing: $p"
$txt = @"
package ${pkgBase}_repo

import (
	"context"
	${pkgBase}_biz "service/internal/feature/$base/v$apiVersion/biz"
)

func (r *${base}Repo) Find${pluralPascal}(ctx context.Context, id *uint, name *string) ([]${pkgBase}_biz.${pascal}, error) {
	return []${pkgBase}_biz.${pascal}{}, nil
}
"@
    [IO.File]::WriteAllText($p, $txt, $utf8NoBom)
    Show-OK "Created: r_list.go"
  } else { Show-Info "Skip (exists): $p" }
} else {
  Show-Info "No GET op; skip r_list.go"
}

# r_upsert.go (UPSERT)
if ($HasUpsert) {
  $p = Join-Path $repoDir "r_upsert.go"
  if (-not (Test-Path $p)) {
    Show-Info "Writing: $p"
$txt = @"
package ${pkgBase}_repo

import (
	"context"
	${pkgBase}_biz "service/internal/feature/$base/v$apiVersion/biz"
)

func (r *${base}Repo) Upsert${pascal}(ctx context.Context, in *${pkgBase}_biz.${pascal}) (*${pkgBase}_biz.${pascal}, error) {
	return in, nil
}
"@
    [IO.File]::WriteAllText($p, $txt, $utf8NoBom)
    Show-OK "Created: r_upsert.go"
  } else { Show-Info "Skip (exists): $p" }
} else {
  Show-Info "No UPSERT op; skip r_upsert.go"
}

# r_delete_by_id.go (DELETE)
if ($HasDelete) {
  $p = Join-Path $repoDir "r_delete_by_id.go"
  if (-not (Test-Path $p)) {
    Show-Info "Writing: $p"
$txt = @"
package ${pkgBase}_repo

import "context"

func (r *${base}Repo) Delete${pascal}ById(ctx context.Context, id uint) error {
	return nil
}
"@
    [IO.File]::WriteAllText($p, $txt, $utf8NoBom)
    Show-OK "Created: r_delete_by_id.go"
  } else { Show-Info "Skip (exists): $p" }
} else {
  Show-Info "No DELETE op; skip r_delete_by_id.go"
}

Show-OK ("repo generated: {0}/v{1}" -f $base, $apiVersion)
