# scripts/ps/workflow/module-wire.ps1
[CmdletBinding()]
param(
  [Parameter(Mandatory = $true)] [string]$Name
)

Set-StrictMode -Version Latest
$ErrorActionPreference = 'Stop'

$FeatureRoot  = "./internal/feature"
$RegistersAgg = "./cmd/service/registers_agg.go"
$MainWire     = "./cmd/service/wire.go"
$utf8NoBom    = New-Object System.Text.UTF8Encoding($false)

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

$base   = ConvertTo-LowerCase $Name
$pascal = ConvertTo-PascalCase $Name

Show-Step "Updating wire for module '$base'"
Show-Info "Normalized: base='$base', pascal='$pascal'"

# ---- detect latest feature vN ----
$featDir = Join-Path $FeatureRoot $base
if (-not (Test-Path -LiteralPath $featDir)) {
  Show-ErrorAndExit "Feature directory not found: $featDir. Generate the feature first."
}

$V = 0
$existing = @()
Get-ChildItem -LiteralPath $featDir -Directory -ea SilentlyContinue | ForEach-Object {
  if ($_.Name -match '^v(\d+)$') { $n=[int]$Matches[1]; $existing += "v$($n)"; if ($n -gt $V){$V=$n} }
}
if ($V -le 0) { Show-ErrorAndExit "No versioned folder under $featDir (expected v1/v2/...)."}
Show-Info ("Feature versions found: " + ($existing -join ", "))
Show-OK  "Using latest: v$V"

$alias      = "{0}v{1}" -f $base, $V
$importPath = "service/internal/feature/$base/v$V"
Show-Info "Alias: $alias"
Show-Info "Import path: $importPath"

# ---------- helpers ----------
function Add-ImportLine {
  param([string[]]$Lines,[string]$Line)
  if (($Lines -join "`n") -match ('^\s*' + [regex]::Escape($Line) + '\s*$')) { return ,$Lines }
  $iStart = ($Lines | Select-String -Pattern '^\s*import\s*\($' | Select-Object -First 1).LineNumber
  if ($iStart) {
    $iEnd = ($Lines | Select-String -Pattern '^\s*\)$' | Where-Object { $_.LineNumber -gt $iStart } | Select-Object -First 1).LineNumber
    if ($iEnd) {
      $before = $Lines[0..($iEnd-2)]
      $after  = $Lines[($iEnd-1)..($Lines.Count-1)]
      return ,($before + @("`t$Line") + $after)
    }
  } else {
    $pkgIdx = ($Lines | Select-String -Pattern '^\s*package\s+\w+\s*$' | Select-Object -First 1).LineNumber
    if ($pkgIdx) {
      $ins = @("", "import (", "`t$Line", ")", "")
      return ,($Lines[0..$pkgIdx] + $ins + $Lines[($pkgIdx+1)..($Lines.Count-1)])
    }
  }
  ,$Lines
}

function Add-ParamAfterAnchor {
  param([string[]]$Lines,[string]$Anchor,[string]$ParamLine)
  if (($Lines -join "`n") -match [regex]::Escape($ParamLine)) { return ,$Lines }
  for ($i=0; $i -lt $Lines.Count; $i++){
    if ($Lines[$i] -match $Anchor){
      $indent = ($Lines[$i] -replace '^(\s*).*','$1')
      $left = $Lines[0..$i]
      $right = $Lines[($i+1)..($Lines.Count-1)]
      return ,($left + @("$indent$ParamLine") + $right)
    }
  }
  ,$Lines
}

function Add-InSlice {
  param([string[]]$Lines,[string]$HeaderRegex,[string]$ItemLine)
  if (($Lines -join "`n") -match [regex]::Escape($ItemLine)) { return ,$Lines }
  for ($i=0; $i -lt $Lines.Count; $i++){
    if ($Lines[$i] -match $HeaderRegex){
      for ($j=$i+1; $j -lt $Lines.Count; $j++){
        if ($Lines[$j] -match '^\s*\},\s*$'){
          $indent = ($Lines[$i] -replace '^(\s*).*','$1') + "`t"
          $before = $Lines[0..($j-1)]
          $after  = $Lines[$j..($Lines.Count-1)]
          return ,($before + @("$indent$ItemLine") + $after)
        }
      }
      break
    }
  }
  ,$Lines
}

# ========== 1) registers_agg.go ==========
Show-Step "Updating $RegistersAgg"
if (-not (Test-Path -LiteralPath $RegistersAgg)) { Show-ErrorAndExit "File not found: $RegistersAgg" }
$raBefore = Get-Content -LiteralPath $RegistersAgg -Raw -Encoding UTF8
$raL = $raBefore -split "`r?`n"

# 1a) import alias
$raL = Add-ImportLine -Lines $raL -Line "$alias `"$importPath`""
if (($raL -join "`n") -ne $raBefore) { Show-OK "Added import alias to registrars_agg.go" } else { Show-Info "Import alias already present" }

# 1b) typed params
$raTmp = $raL -join "`n"
$raL = Add-ParamAfterAnchor -Lines $raL -Anchor '^[ \t]*// add other HTTP-registrers for modules here:.*$' -ParamLine "$alias`HTTP $alias.HTTPRegister,"
if (($raL -join "`n") -ne $raTmp) { Show-OK "Added HTTP param to BuildAllRegistrars" } else { Show-Info "HTTP param already present" }

$raTmp = $raL -join "`n"
$raL = Add-ParamAfterAnchor -Lines $raL -Anchor '^[ \t]*// add other gRPC-registrers for modules here:.*$'  -ParamLine "$alias`GRPC $alias.GRPCRegister,"
if (($raL -join "`n") -ne $raTmp) { Show-OK "Added GRPC param to BuildAllRegistrars" } else { Show-Info "GRPC param already present" }

# 1c) slice items
$raTmp = $raL -join "`n"
$raL = Add-InSlice -Lines $raL -HeaderRegex 'HTTP:\s*\[\]server_http\.HTTPRegister\s*\{' -ItemLine "server_http.HTTPRegister($alias`HTTP),"
if (($raL -join "`n") -ne $raTmp) { Show-OK "Inserted item into HTTP slice" } else { Show-Info "HTTP slice already contains item" }

$raTmp = $raL -join "`n"
$raL = Add-InSlice -Lines $raL -HeaderRegex 'GRPC:\s*\[\]server_grpc\.GRPCRegister\s*\{' -ItemLine "server_grpc.GRPCRegister($alias`GRPC),"
if (($raL -join "`n") -ne $raTmp) { Show-OK "Inserted item into GRPC slice" } else { Show-Info "GRPC slice already contains item" }

[IO.File]::WriteAllLines($RegistersAgg, $raL, $utf8NoBom)
Show-OK "registrars_agg.go updated"

# ========== 2) cmd/service/wire.go ==========
Show-Step "Updating $MainWire"
if (-not (Test-Path -LiteralPath $MainWire)) { Show-ErrorAndExit "File not found: $MainWire" }
$mwBefore = Get-Content -LiteralPath $MainWire -Raw -Encoding UTF8
$mwL = $mwBefore -split "`r?`n"

# 2a) import alias
$mwL = Add-ImportLine -Lines $mwL -Line "$alias `"$importPath`""
if (($mwL -join "`n") -ne $mwBefore) { Show-OK "Added import alias to wire.go" } else { Show-Info "Import alias already present in wire.go" }

# 2b) ProviderSet in // modules with proper indent
$prov = "$alias.ProviderSet,"

if (($mwL -join "`n") -notmatch [regex]::Escape($prov)) {
  $modulesMatch = ($mwL | Select-String -Pattern '^\s*//\s*modules\s*$' | Select-Object -First 1)
  if ($modulesMatch) {
    $mIdx = $modulesMatch.LineNumber - 1

    # взять отступ первой непустой строки после // modules (например, template.ProviderSet,)
    $indent = $null
    for ($k = $mIdx + 1; $k -lt $mwL.Count; $k++) {
      if ($mwL[$k] -match '^\s*$') { continue }                       # пустые
      $indent = ($mwL[$k] -replace '^(\s*).*','$1')                   # РОВНО как у существующей строки
      break
    }
    if (-not $indent) {
      # fallback: ровно как у комментария // modules (без добавления \t)
      $indent = ($mwL[$mIdx] -replace '^(\s*).*','$1')
    }

    $mwL = $mwL[0..$mIdx] + @("$indent$prov") + $mwL[($mIdx+1)..($mwL.Count-1)]
    Show-OK "Inserted ProviderSet under // modules with indent='$indent'"
  } else {
    # fallback: перед newApp, внутри wire.Build(...)
    $buildStart = ($mwL | Select-String -Pattern 'wire\.Build\(' | Select-Object -First 1).LineNumber
    $newAppIdx  = ($mwL | Select-String -Pattern '^\s*newApp,\s*$' | Where-Object { $_.LineNumber -gt $buildStart } | Select-Object -First 1).LineNumber
    if ($buildStart -and $newAppIdx) {
      $indent = ($mwL[$newAppIdx-1] -replace '^(\s*).*','$1')
      $mwL = $mwL[0..($newAppIdx-2)] + @("$indent$prov") + $mwL[($newAppIdx-1)..($mwL.Count-1)]
      Show-Warn "Inserted ProviderSet before newApp (no // modules anchor found)"
    } else {
      Show-ErrorAndExit "Could not find insertion point for ProviderSet in wire.go"
    }
  }
} else {
  Show-Info "ProviderSet already present in wire.go"
}

[IO.File]::WriteAllLines($MainWire, $mwL, $utf8NoBom)
Show-OK ("wire updated: module '{0}' -> {1}, alias {2}" -f $base, $importPath, $alias)
