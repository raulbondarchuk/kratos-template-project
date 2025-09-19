# scripts/ps/workflow/module-tests-delete.ps1
[CmdletBinding()]
param(
  [Parameter(Mandatory=$true)] [string]$Name,
  [string]$Version = ""
)

Set-StrictMode -Version Latest
$ErrorActionPreference = 'Stop'

# --- utils (logs) ---
try {
  . (Join-Path $PSScriptRoot 'utils.ps1')
} catch {
  function Show-Step { param([string]$Message) Write-Host "`n==> $Message" -ForegroundColor Cyan }
  function Show-Info { param([string]$Message) Write-Host "$Message" -ForegroundColor DarkGray }
  function Show-OK   { param([string]$Message) Write-Host "  [OK] $Message" -ForegroundColor Green }
  function Show-Warn { param([string]$Message) Write-Host "  [WARN] $Message" -ForegroundColor Yellow }
  function Show-ErrorAndExit { param([string]$Message) Write-Host "  [ERROR] $Message" -ForegroundColor Red; exit 1 }
}

# --- helpers ---
function ConvertTo-LowerCase { param([string]$s) $s.ToLower() }

# --- normalize names ---
$base = ConvertTo-LowerCase $Name
$ApiRoot = "./api/$base"
$FeatureRoot = "./internal/feature/$base"

Show-Step "Deleting unit tests"
Show-Info "Module: base='$base'"

# --- resolve version ---
if (-not (Test-Path $ApiRoot)) { Show-ErrorAndExit "API dir not found: $ApiRoot" }

$versions = @()
if ([string]::IsNullOrWhiteSpace($Version)) {
  # Delete all versions
  Get-ChildItem $ApiRoot -Directory -ErrorAction SilentlyContinue | ForEach-Object {
    if ($_.Name -match '^v\d+$') { $versions += $_.Name }
  }
  if ($versions.Count -eq 0) { Show-ErrorAndExit "No API versions in $ApiRoot" }
} else {
  # Delete specific version
  if ($Version -match '^v?(\d+)$') {
    $v = "v$($Matches[1])"
    if (-not (Test-Path (Join-Path $ApiRoot $v))) { Show-ErrorAndExit "Version not found: $v" }
    $versions = @($v)
  } else {
    Show-ErrorAndExit "Invalid version format: '$Version' (use 'vN' or 'N')"
  }
}

# --- delete test files ---
foreach ($v in $versions) {
  $testFile = Join-Path $FeatureRoot "$v/service/${base}_service_test.go"
  if (Test-Path $testFile) {
    Remove-Item $testFile -Force
    Show-OK "Deleted test: $testFile"
  } else {
    Show-Info "Test file not found: $testFile"
  }
}
