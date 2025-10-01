# scripts/ps/workflow/module-delete.ps1
[CmdletBinding()]
param(
  [Parameter(Mandatory = $true)] [string]$Name,
  [string]$Version  # optional: v2, v3...
)

Set-StrictMode -Version Latest
$ErrorActionPreference = 'Stop'

$ApiRoot      = "./api"
$FeatureRoot  = "./internal/feature"
$RegistersAgg = "./cmd/service/registers_agg.go"
$MainWire     = "./cmd/service/wire.go"
$RoutesFile   = "./internal/feature/routes.go"
$utf8NoBom    = New-Object System.Text.UTF8Encoding($false)

# --- logs ---
try { . (Join-Path $PSScriptRoot 'utils.ps1') } catch {
  function Show-Step { param([string]$Message) Write-Host "`n==> $Message" -ForegroundColor Cyan }
  function Show-Info { param([string]$Message) Write-Host "$Message" -ForegroundColor DarkGray }
  function Show-OK   { param([string]$Message) Write-Host "  [OK] $Message" -ForegroundColor Green }
  function Show-Warn { param([string]$Message) Write-Host "  [WARN] $Message" -ForegroundColor Yellow }
  function Show-ErrorAndExit { param([string]$Message) Write-Host "  [ERROR] $Message" -ForegroundColor Red; exit 1 }
}

function ConvertTo-LowerCase { param([string]$s) $s.ToLower() }
$base = ConvertTo-LowerCase $Name

# === guard: forbid 'common' module deletion ===
if ($base -eq 'common') {
  Show-ErrorAndExit "Module 'common' is reserved and cannot be deleted."
}

# --- discover versions present ---
$featBase = Join-Path $FeatureRoot $base
$apiBase  = Join-Path $ApiRoot $base

$featVersions = @()
if (Test-Path -LiteralPath $featBase) {
  Get-ChildItem -LiteralPath $featBase -Directory -ea SilentlyContinue | ForEach-Object {
    if ($_.Name -match '^v(\d+)$') { $featVersions += [int]$Matches[1] }
  }
}

$apiVersions = @()
if (Test-Path -LiteralPath $apiBase) {
  Get-ChildItem -LiteralPath $apiBase -Directory -ea SilentlyContinue | ForEach-Object {
    if ($_.Name -match '^v(\d+)$') { $apiVersions += [int]$Matches[1] }
  }
}

# --- choose versions to delete ---
$versionsToDelete = @()
if ([string]::IsNullOrWhiteSpace($Version)) {
  Show-Step "Delete ALL versions of module '$base'"

  $versionsToDelete = @()
  $versionsToDelete += $featVersions
  $versionsToDelete += $apiVersions
  $versionsToDelete = $versionsToDelete | Sort-Object -Unique

  if (-not $versionsToDelete -or @($versionsToDelete).Count -eq 0) {
    Show-Warn "No versions found under feature/api (still cleaning code links)."
  }
} else {
  if ($Version -notmatch '^v(\d+)$') { Show-ErrorAndExit "Invalid -Version '$Version'. Use vN (e.g. v2)." }
  $vNum = [int]$Matches[1]

  $exists = ($featVersions -contains $vNum) -or ($apiVersions -contains $vNum)
  if (-not $exists) {
    $availList = @()
    $availList += $featVersions
    $availList += $apiVersions
    $availList  = $availList | Sort-Object -Unique
    $avail      = ($availList | ForEach-Object { "v$_" }) -join ", "
    if ([string]::IsNullOrWhiteSpace($avail)) { $avail = "<none>" }
    Show-ErrorAndExit "Version $Version of '$base' not found. Available: $avail"
  }
  Show-Step "Delete ONLY version $Version of module '$base'"
  $versionsToDelete = @($vNum)
}

# --- helper: remove directories if exist ---
function Remove-DirIfExists([string]$Path) {
  if (Test-Path -LiteralPath $Path) {
    Show-Info "Deleting: $Path"
    Remove-Item -LiteralPath $Path -Recurse -Force
    Show-OK "Removed: $Path"
  } else {
    Show-Info "Skip (not found): $Path"
  }
}

# --- 1) delete feature dirs ---
if (@($versionsToDelete).Count -gt 0) {
  foreach ($v in $versionsToDelete) { Remove-DirIfExists (Join-Path $featBase ("v$v")) }
  if ((Test-Path -LiteralPath $featBase) -and ((Get-ChildItem -LiteralPath $featBase -Force | Measure-Object).Count -eq 0)) {
    Remove-DirIfExists $featBase
  }
} else {
  Remove-DirIfExists $featBase
}

# --- 2) delete api dirs ---
if (@($versionsToDelete).Count -gt 0) {
  foreach ($v in $versionsToDelete) { Remove-DirIfExists (Join-Path $apiBase ("v$v")) }
  if ((Test-Path -LiteralPath $apiBase) -and ((Get-ChildItem -LiteralPath $apiBase -Force | Measure-Object).Count -eq 0)) {
    Remove-DirIfExists $apiBase
  }
} else {
  Remove-DirIfExists $apiBase
}

# --- build aliases to purge from code ---
$aliases = @()
if (@($versionsToDelete).Count -gt 0) {
  foreach ($v in $versionsToDelete) { $aliases += ("{0}v{1}" -f $base, $v) }
}

# --- helper: multiline regex replace with count ---
function Update-MultilineText {
  param([string]$Text, [string]$Pattern, [string]$With, [string]$What)
  # use inline-flag (?m) inside $Pattern
  $re   = New-Object System.Text.RegularExpressions.Regex($Pattern)
  $cnt  = ($re.Matches($Text)).Count
  if ($cnt -gt 0) { Show-OK ("Removed {0} occurrence(s): {1}" -f $cnt, $What) }
  else { Show-Info ("No matches to remove: {0}" -f $What) }
  return $re.Replace($Text, $With)
}

# --- 3) clean cmd/service/registers_agg.go ---
if (-not (Test-Path -LiteralPath $RegistersAgg)) {
  Show-Warn "File not found: $RegistersAgg (skip)."
} else {
  Show-Step "Cleaning links in $RegistersAgg"
  $txt = Get-Content -LiteralPath $RegistersAgg -Raw -Encoding UTF8

  if ($aliases.Count -gt 0) {
    foreach ($a in $aliases) {
      # import line: \tpruebav2 "…"
      $patImport    = '(?m)^\s*' + [regex]::Escape($a) + '\s+"[^"]+"\s*(//.*)?\s*$'
      # params: pruebav2HTTP pruebav2.HTTPRegister,
      $patParamHTTP = '(?m)^\s*' + [regex]::Escape($a) + 'HTTP\s+' + [regex]::Escape($a) + '\.HTTPRegister\s*,\s*$'
      $patParamGRPC = '(?m)^\s*' + [regex]::Escape($a) + 'GRPC\s+' + [regex]::Escape($a) + '\.GRPCRegister\s*,\s*$'
      # slice items
      $patItemHTTP  = '(?m)^\s*server_http\.HTTPRegister\(\s*' + [regex]::Escape($a) + 'HTTP\s*\)\s*,\s*$'
      $patItemGRPC  = '(?m)^\s*server_grpc\.GRPCRegister\(\s*' + [regex]::Escape($a) + 'GRPC\s*\)\s*,\s*$'

      $txt = Update-MultilineText -Text $txt -Pattern $patImport    -With '' -What "import: $a"
      $txt = Update-MultilineText -Text $txt -Pattern $patParamHTTP -With '' -What "param: ${a}HTTP"
      $txt = Update-MultilineText -Text $txt -Pattern $patParamGRPC -With '' -What "param: ${a}GRPC"
      $txt = Update-MultilineText -Text $txt -Pattern $patItemHTTP  -With '' -What "slice item: ${a}HTTP"
      $txt = Update-MultilineText -Text $txt -Pattern $patItemGRPC  -With '' -What "slice item: ${a}GRPC"
    }
  } else {
    # remove ALL versions basev\d+
    $patImport    = '(?m)^\s*' + [regex]::Escape($base) + 'v\d+\s+"[^"]+"\s*(//.*)?\s*$'
    $patParamHTTP = '(?m)^\s*' + [regex]::Escape($base) + 'v\d+HTTP\s+' + [regex]::Escape($base) + 'v\d+\.HTTPRegister\s*,\s*$'
    $patParamGRPC = '(?m)^\s*' + [regex]::Escape($base) + 'v\d+GRPC\s+' + [regex]::Escape($base) + 'v\d+\.GRPCRegister\s*,\s*$'
    $patItemHTTP  = '(?m)^\s*server_http\.HTTPRegister\(\s*' + [regex]::Escape($base) + 'v\d+HTTP\s*\)\s*,\s*$'
    $patItemGRPC  = '(?m)^\s*server_grpc\.GRPCRegister\(\s*' + [regex]::Escape($base) + 'v\d+GRPC\s*\)\s*,\s*$'

    $txt = Update-MultilineText -Text $txt -Pattern $patImport    -With '' -What "imports: ${base}v*"
    $txt = Update-MultilineText -Text $txt -Pattern $patParamHTTP -With '' -What "params: ${base}v* HTTP"
    $txt = Update-MultilineText -Text $txt -Pattern $patParamGRPC -With '' -What "params: ${base}v* GRPC"
    $txt = Update-MultilineText -Text $txt -Pattern $patItemHTTP  -With '' -What "slice items: ${base}v* HTTP"
    $txt = Update-MultilineText -Text $txt -Pattern $patItemGRPC  -With '' -What "slice items: ${base}v* GRPC"
  }

  # collapse extra empty lines
  $txt = [regex]::Replace($txt, "(\r?\n){3,}", "`r`n`r`n")
  [IO.File]::WriteAllText($RegistersAgg, $txt, $utf8NoBom)
  Show-OK "Updated: $RegistersAgg"
}

# --- 4) clean cmd/service/wire.go ---
if (-not (Test-Path -LiteralPath $MainWire)) {
  Show-Warn "File not found: $MainWire (skip)."
} else {
  Show-Step "Cleaning links in $MainWire"
  $txt = Get-Content -LiteralPath $MainWire -Raw -Encoding UTF8

  if ($aliases.Count -gt 0) {
    foreach ($a in $aliases) {
      $patImp  = '(?m)^\s*' + [regex]::Escape($a) + '\s+"[^"]+"\s*$'
      $patProv = '(?m)^\s*' + [regex]::Escape($a) + '\.ProviderSet\s*,\s*$'
      $txt = Update-MultilineText -Text $txt -Pattern $patImp  -With '' -What "wire import: $a"
      $txt = Update-MultilineText -Text $txt -Pattern $patProv -With '' -What "wire provider: $a"
    }
  } else {
    $patImp  = '(?m)^\s*' + [regex]::Escape($base) + 'v\d+\s+"[^"]+"\s*$'
    $patProv = '(?m)^\s*' + [regex]::Escape($base) + 'v\d+\.ProviderSet\s*,\s*$'
    $txt = Update-MultilineText -Text $txt -Pattern $patImp  -With '' -What "wire imports: ${base}v*"
    $txt = Update-MultilineText -Text $txt -Pattern $patProv -With '' -What "wire providers: ${base}v*"
  }

  $txt = [regex]::Replace($txt, "(\r?\n){3,}", "`r`n`r`n")
  [IO.File]::WriteAllText($MainWire, $txt, $utf8NoBom)
  Show-OK "Updated: $MainWire"
}

# --- 5) clean internal/feature/routes.go ---
if (-not (Test-Path -LiteralPath $RoutesFile)) {
  Show-Warn "File not found: $RoutesFile (skip)."
} else {
  Show-Step "Cleaning routes in $RoutesFile"
  $txt = Get-Content -LiteralPath $RoutesFile -Raw -Encoding UTF8

  # Patterns for matching routes entries
  if ($Version -and $Version -match '^v(\d+)$') {
    $v = $Matches[1]
    # For specific version
    $patImport = '(?m)^\s*' + [regex]::Escape($base) + '_v' + $v + '\s+"service/internal/feature/' + [regex]::Escape($base) + '/v' + $v + '"\s*$'
    $patSvcImport = '(?m)^\s*' + [regex]::Escape($base) + '_v' + $v + '_service\s+"service/internal/feature/' + [regex]::Escape($base) + '/v' + $v + '/service"\s*$'
    $patParam = '(?m)^\s*' + [regex]::Escape($base) + 'V' + $v + 'Svc\s+\*' + [regex]::Escape($base) + '_v' + $v + '_service\.[A-Za-z0-9]+Service\s*,\s*$'
    $patGroup = '(?m)^\s*' + [regex]::Escape($base) + '_v' + $v + '\.GetServiceEndpoints\(' + [regex]::Escape($base) + 'V' + $v + 'Svc\)\s*,?\s*$'
  } else {
    # For all versions
    $patImport = '(?m)^\s*' + [regex]::Escape($base) + '_v\d+\s+"service/internal/feature/' + [regex]::Escape($base) + '/v\d+"\s*$'
    $patSvcImport = '(?m)^\s*' + [regex]::Escape($base) + '_v\d+_service\s+"service/internal/feature/' + [regex]::Escape($base) + '/v\d+/service"\s*$'
    $patParam = '(?m)^\s*' + [regex]::Escape($base) + 'V\d+Svc\s+\*' + [regex]::Escape($base) + '_v\d+_service\.[A-Za-z0-9]+Service\s*,\s*$'
    $patGroup = '(?m)^\s*' + [regex]::Escape($base) + '_v\d+\.GetServiceEndpoints\(' + [regex]::Escape($base) + 'V\d+Svc\)\s*,?\s*$'
  }

  # Remove all matches
  $txt = Update-MultilineText -Text $txt -Pattern $patImport    -With '' -What "routes import: $base"
  $txt = Update-MultilineText -Text $txt -Pattern $patSvcImport -With '' -What "routes service import: ${base}_service"
  $txt = Update-MultilineText -Text $txt -Pattern $patParam     -With '' -What "routes param: ${base}Svc"
  $txt = Update-MultilineText -Text $txt -Pattern $patGroup     -With '' -What "routes group: ${base}.GetServiceEndpoints"

  # Fix formatting after removals
  # First clean up any empty lines between items
  $txt = $txt -replace '(?m)(\t\t[^\n\r]+),\s*\n\s*\n\s*(\t\t)', "`$1,`n`$2"
  
  # Then handle the endpoint list formatting
  if ($txt -match '(?m)return\s+\[\]endpoint\.ServiceGroup\s*\{([^\}]+)\}') {
    $body = $Matches[1]
    # If only comments remain
    if ($body -match '^\s*(?://[^\n]*\n\s*)*$') {
      $txt = $txt -replace '(?m)return\s+\[\]endpoint\.ServiceGroup\s*\{[^\}]+\}', 'return []endpoint.ServiceGroup{}'
    } else {
      # Clean up the body: ensure comma after each item except the last one
      $items = @()
      $comments = @()
      foreach ($line in ($body -split "`n")) {
        if ($line -match '^\s*//') {
          $comments += $line
        } elseif ($line -match '\S') {
          # Remove any existing comma
          $line = $line -replace ',\s*$', ''
          $items += $line
        }
      }
      
      # Rebuild the body with proper formatting
      $newBody = "`n"
      if ($comments.Count -gt 0) {
        $newBody += ($comments -join "`n") + "`n"
      }
      for ($i = 0; $i -lt $items.Count; $i++) {
        $item = $items[$i].TrimEnd()
        # Always add comma after each item
        $newBody += "${item},`n"
      }
      $newBody += "`t"
      
      $txt = $txt -replace '(?m)return\s+\[\]endpoint\.ServiceGroup\s*\{[^\}]+\}', "return []endpoint.ServiceGroup{$newBody}"
    }
  }

  # collapse extra empty lines
  $txt = [regex]::Replace($txt, "(\r?\n){3,}", "`r`n`r`n")
  [IO.File]::WriteAllText($RoutesFile, $txt, $utf8NoBom)
  Show-OK "Updated: $RoutesFile"
}

# --- done ---
if ($Version) {
  Show-OK ("Module '{0}' version {1} removed (feature/api + code links cleaned)" -f $base, $Version)
} else {
  Show-OK ("Module '{0}' fully removed (all versions; feature/api + code links cleaned)" -f $base)
}
