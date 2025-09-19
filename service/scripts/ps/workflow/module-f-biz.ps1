# scripts/ps/workflow/biz.ps1
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

# --- parse ops -> flags ---
$opsList = @()
if ($Ops) { $opsList = $Ops.ToLower().Split(',') | ForEach-Object { $_.Trim() } | Where-Object { $_ } }
$HasGet = $false; $HasUpsert = $false; $HasDelete = $false
foreach($op in $opsList){
  switch ($op) {
    'get' { $HasGet = $true }
    'find' { $HasGet = $true }
    'list' { $HasGet = $true }
    'read' { $HasGet = $true }
    'upsert' { $HasUpsert = $true }
    'create' { $HasUpsert = $true }
    'update' { $HasUpsert = $true }
    'delete' { $HasDelete = $true }
    'del' { $HasDelete = $true }
    'remove' { $HasDelete = $true }
    default { Show-Warn "Unknown op '$op' ignored" }
  }
}
$AnyOps = $HasGet -or $HasUpsert -or $HasDelete

Show-Step "Generating biz layer"
$base         = ConvertTo-LowerCase $Name
$pkgBase      = ($base -replace '[^a-z0-9]','_')
$pascal       = ConvertTo-PascalCase $Name
$pluralPascal = ConvertTo-PascalCase (ConvertTo-Plural $base)
Show-Info "Module: base='$base', pascal='$pascal', plural='$pluralPascal', ops=[$Ops]"

# --- detect latest API vN ---
Show-Info "Detecting latest API version in '$ApiRoot/$base'..."
$apiBaseDir = Join-Path $ApiRoot $base
$apiVersion = 1
if (Test-Path $apiBaseDir) {
  $max = 0
  Get-ChildItem $apiBaseDir -Directory | ForEach-Object {
    if ($_.Name -match '^v(\d+)$') { $n=[int]$Matches[1]; if ($n -gt $max) { $max=$n } }
  }
  if ($max -gt 0) { $apiVersion=$max }
  Show-Info "Using API version: v$apiVersion"
} else {
  Show-Warn "API dir not found: $apiBaseDir. Using v1."
}

$featureRootV = Join-Path (Join-Path $FeatureRoot $base) "v$apiVersion"
$bizDir = Join-Path $featureRootV "biz"

Show-Info "Ensuring directories:"
$null = New-Item -ItemType Directory -Force -Path $bizDir
Show-Info "Created/exists: $bizDir"

# --- biz.go (dynamic imports + interface) ---
$p = Join-Path $bizDir "biz.go"
if (-not (Test-Path $p)) {
  Show-Info "Writing: $p"

  $imports = @()
  if ($AnyOps) { $imports += '"context"' }
  $imports += '"github.com/go-kratos/kratos/v2/log"'
  $importsBlock = "import (`n`t" + ($imports -join "`n`t") + "`n)"

  $repoMethods = @()
  if ($HasGet)    { $repoMethods += "Find${pluralPascal}(ctx context.Context, id *uint, name *string) ([]${pascal}, error)" }
  if ($HasUpsert) { $repoMethods += "Upsert${pascal}(ctx context.Context, in *${pascal}) (*${pascal}, error)" }
  if ($HasDelete) { $repoMethods += "Delete${pascal}ById(ctx context.Context, id uint) error" }
  $repoMethodsText = ""
  if ($repoMethods.Count -gt 0) { $repoMethodsText = "`n`t" + ($repoMethods -join "`n`t") + "`n" }

  $txt = @"
package ${pkgBase}_biz

$importsBlock

type ${pascal}Repo interface {
$repoMethodsText}
type ${pascal}Usecase struct {
	repo ${pascal}Repo
	log  *log.Helper
}

func New${pascal}Usecase(repo ${pascal}Repo, logger log.Logger) *${pascal}Usecase {
	return &${pascal}Usecase{repo: repo, log: log.NewHelper(logger)}
}
"@
  [IO.File]::WriteAllText($p, $txt, $utf8NoBom)
  Show-OK "Created: biz.go"
} else {
  Show-Info "Skip (exists): $p"
}

# --- entity.go (always) ---
$p = Join-Path $bizDir "entity.go"
if (-not (Test-Path $p)) {
  Show-Info "Writing: $p"
  $txt = @"
package ${pkgBase}_biz

import "time"

type ${pascal} struct {
	ID        uint
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}
"@
  [IO.File]::WriteAllText($p, $txt, $utf8NoBom)
  Show-OK "Created: entity.go"
} else {
  Show-Info "Skip (exists): $p"
}

# --- usecase.go (only if any ops) ---
$p = Join-Path $bizDir "usecase.go"
if ($AnyOps -and -not (Test-Path $p)) {
  Show-Info "Writing: $p"

  $ucMethods = @()
  if ($HasGet) {
    $ucMethods += @"
func (uc *${pascal}Usecase) Find${pluralPascal}(ctx context.Context, id *uint, name *string) ([]${pascal}, error) {
	return uc.repo.Find${pluralPascal}(ctx, id, name)
}
"@
  }
  if ($HasUpsert) {
    $ucMethods += @"
func (uc *${pascal}Usecase) Upsert${pascal}(ctx context.Context, in *${pascal}) (*${pascal}, error) {
	return uc.repo.Upsert${pascal}(ctx, in)
}
"@
  }
  if ($HasDelete) {
    $ucMethods += @"
func (uc *${pascal}Usecase) Delete${pascal}ById(ctx context.Context, id uint) error {
	return uc.repo.Delete${pascal}ById(ctx, id)
}
"@
  }
  $methodsText = ($ucMethods -join "`n")

  $txt = @"
package ${pkgBase}_biz

import "context"

$methodsText
"@
  [IO.File]::WriteAllText($p, $txt, $utf8NoBom)
  Show-OK "Created: usecase.go"
} elseif (-not $AnyOps) {
  Show-Info "No ops requested; skip usecase.go"
} else {
  Show-Info "Skip (exists): $p"
}

Show-OK ("biz generated: {0}/v{1}" -f $base, $apiVersion)
