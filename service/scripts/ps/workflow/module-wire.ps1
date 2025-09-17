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

function ConvertTo-PascalCase { param([string]$s)
  $parts = ($s -replace '[^A-Za-z0-9]+',' ') -split '\s+' | Where-Object { $_ }
  ($parts | ForEach-Object { $_.Substring(0,1).ToUpper() + $_.Substring(1).ToLower() }) -join ''
}
function ConvertTo-LowerCase { param([string]$s) $s.ToLower() }

$base    = ConvertTo-LowerCase $Name     # "prueba"
# Используем $pascal в сообщении в конце скрипта
$pascal  = ConvertTo-PascalCase $Name    # "Prueba"

# ---- detect latest feature vN ----
$featDir = Join-Path $FeatureRoot $base
if (-not (Test-Path -LiteralPath $featDir)) {
  throw "Feature directory not found: $featDir. Generate the feature first."
}
$V = 0
Get-ChildItem -LiteralPath $featDir -Directory -ea SilentlyContinue | ForEach-Object {
  if ($_.Name -match '^v(\d+)$') { $n=[int]$Matches[1]; if ($n -gt $V){$V=$n} }
}
if ($V -le 0) { throw "No versioned folder under $featDir (expected v1/v2/...)."}
$alias      = "{0}v{1}" -f $base, $V
$importPath = "service/internal/feature/$base/v$V"

# ---------- helpers ----------
function Add-ImportLine {
  param([string[]]$Lines,[string]$Line)
  if (($Lines -join "`n") -match [regex]::Escape($Line)) { return ,$Lines }
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
      # find closing brace of this slice (line with only '},')
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
if (-not (Test-Path -LiteralPath $RegistersAgg)) { throw "File not found: $RegistersAgg" }
$ra = Get-Content -LiteralPath $RegistersAgg -Raw -Encoding UTF8
$raL = $ra -split "`r?`n"

# 1a) ensure import of module alias
$raL = Add-ImportLine -Lines $raL -Line "$alias `"$importPath`""

# 1b) ensure typed params in BuildAllRegistrars
$raL = Add-ParamAfterAnchor -Lines $raL -Anchor '^[ \t]*// add other HTTP-registrers for modules here:.*$' -ParamLine "$alias`HTTP $alias.HTTPRegister,"
$raL = Add-ParamAfterAnchor -Lines $raL -Anchor '^[ \t]*// add other gRPC-registrers for modules here:.*$'  -ParamLine "$alias`GRPC $alias.GRPCRegister,"

# 1c) ensure slice items with cast to common types
$raL = Add-InSlice -Lines $raL -HeaderRegex 'HTTP:\s*\[\]server_http\.HTTPRegister\s*\{' -ItemLine "server_http.HTTPRegister($alias`HTTP),"
$raL = Add-InSlice -Lines $raL -HeaderRegex 'GRPC:\s*\[\]server_grpc\.GRPCRegister\s*\{' -ItemLine "server_grpc.GRPCRegister($alias`GRPC),"

[IO.File]::WriteAllLines($RegistersAgg, $raL, $utf8NoBom)

# ========== 2) cmd/service/wire.go ==========
if (-not (Test-Path -LiteralPath $MainWire)) { throw "File not found: $MainWire" }
$mw = Get-Content -LiteralPath $MainWire -Raw -Encoding UTF8
$mwL = $mw -split "`r?`n"

# 2a) import alias
$mwL = Add-ImportLine -Lines $mwL -Line "$alias `"$importPath`""

# 2b) ProviderSet in // modules section (or before newApp,)
$prov = "$alias.ProviderSet,"
if (($mwL -join "`n") -notmatch [regex]::Escape($prov)) {
  $modulesIdx = ($mwL | Select-String -Pattern '^\s*//\s*modules\s*$' | Select-Object -First 1).LineNumber
  if ($modulesIdx) {
    $idx = $modulesIdx
    $indent = ($mwL[$idx] -replace '^(\s*).*','$1') + "`t"
    $mwL = $mwL[0..$idx] + @("$indent$prov") + $mwL[($idx+1)..($mwL.Count-1)]
  } else {
    $buildStart = ($mwL | Select-String -Pattern 'wire\.Build\(' | Select-Object -First 1).LineNumber
    $newAppIdx  = ($mwL | Select-String -Pattern '^\s*newApp,\s*$' | Where-Object { $_.LineNumber -gt $buildStart } | Select-Object -First 1).LineNumber
    if ($buildStart -and $newAppIdx) {
      $indent = ($mwL[$newAppIdx-1] -replace '^(\s*).*','$1')
      $mwL = $mwL[0..($newAppIdx-2)] + @("$indent$prov") + $mwL[($newAppIdx-1)..($mwL.Count-1)]
    }
  }
}

[IO.File]::WriteAllLines($MainWire, $mwL, $utf8NoBom)

Write-Host ("wire updated: module '{0}' -> {1}, alias {2}, type {3}" -f $base, $importPath, $alias, $pascal) -ForegroundColor Green
