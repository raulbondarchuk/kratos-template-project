# scripts/ps/make/wire-check.ps1
[CmdletBinding()]
param()

Set-StrictMode -Version Latest
$ErrorActionPreference = 'Stop'

. "$PSScriptRoot/utils.ps1"

$FeatureRoot = "./internal/feature"
$WireFile = "./cmd/service/wire.go"
$utf8NoBom = New-Object System.Text.UTF8Encoding($false)

Show-Step "Checking features and wire.go configuration"

# Count feature modules (exclude common directories like 'routes.go')
$featureCount = 0
if (Test-Path -LiteralPath $FeatureRoot) {
  Get-ChildItem -LiteralPath $FeatureRoot -Directory -ea SilentlyContinue | ForEach-Object {
    if ($_.Name -ne 'common' -and (Test-Path -LiteralPath (Join-Path $_.FullName 'v1'))) {
      $featureCount++
    }
  }
}

Show-Info "Found $featureCount feature module(s)"

# Read wire.go
if (-not (Test-Path -LiteralPath $WireFile)) {
  Show-ErrorAndExit "Wire file not found: $WireFile"
}

$wireContent = Get-Content -LiteralPath $WireFile -Raw -Encoding UTF8

# Check if data.ProviderSet and its import are commented
$isProviderCommented = $wireContent -match '(?m)^\s*//\s*data\.ProviderSet\s*,'
$isImportCommented = $wireContent -match '(?m)^\s*//\s*"service/internal/data"'

if ($featureCount -eq 0) {
  # No features - should be commented
  $needsUpdate = $false
  
  if (-not $isProviderCommented) {
    Show-Info "No features found - commenting out data.ProviderSet"
    $wireContent = $wireContent -replace '(?m)(^\s*)data\.ProviderSet\s*,', '$1// data.ProviderSet,'
    $needsUpdate = $true
  }
  
  if (-not $isImportCommented) {
    Show-Info "No features found - commenting out data import"
    $wireContent = $wireContent -replace '(?m)(^\s*)"service/internal/data"', '$1// "service/internal/data"'
    $needsUpdate = $true
  }
  
  if ($needsUpdate) {
    [IO.File]::WriteAllText($WireFile, $wireContent, $utf8NoBom)
    Show-OK "Updated wire.go: data components commented out"
  } else {
    Show-Info "Wire.go already configured correctly (data components commented out)"
  }
} else {
  # Has features - should be uncommented
  $needsUpdate = $false
  
  if ($isProviderCommented) {
    Show-Info "Features found - uncommenting data.ProviderSet"
    $wireContent = $wireContent -replace '(?m)^\s*//\s*(data\.ProviderSet\s*,)', "`t`t`$1"
    $needsUpdate = $true
  }
  
  if ($isImportCommented) {
    Show-Info "Features found - uncommenting data import"
    $wireContent = $wireContent -replace '(?m)^\s*//\s*("service/internal/data")', "`t`$1"
    $needsUpdate = $true
  }
  
  if ($needsUpdate) {
    [IO.File]::WriteAllText($WireFile, $wireContent, $utf8NoBom)
    Show-OK "Updated wire.go: data components uncommented"
  } else {
    Show-Info "Wire.go already configured correctly (data components active)"
  }
}

Show-OK "Wire configuration check complete"
