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

function Get-PascalCase { param([string]$s)
  $parts = ($s -replace '[^A-Za-z0-9]+',' ') -split '\s+' | Where-Object { $_ }
  ($parts | ForEach-Object { $_.Substring(0,1).ToUpper() + $_.Substring(1).ToLower() }) -join ''
}
function Get-LowerCase { param([string]$s) $s.ToLower() }

# -------- find latest feature version vN --------
$base    = Get-LowerCase $Name           # "prueba"
$featDir = Join-Path $FeatureRoot $base
if (-not (Test-Path -LiteralPath $featDir)) {
  throw "Feature directory not found: $featDir. Generate feature first."
}
$V = 0
Get-ChildItem -LiteralPath $featDir -Directory | ForEach-Object {
  if ($_.Name -match '^v(\d+)$') {
    $n = [int]$Matches[1]; if ($n -gt $V) { $V = $n }
  }
}
if ($V -le 0) { throw "No versioned feature folder under $featDir (expected v1/v2/...)."}
$alias      = "{0}v{1}" -f $base, $V                # pruebav2
$importPath = "service/internal/feature/$base/v$V"  # service/internal/feature/prueba/v2

# small helpers (line-based)
function Ensure-Line-AfterAnchor {
  param(
    [string[]]$Lines,
    [string]  $AnchorRegex,
    [string]  $LineToInsert
  )
  # already present?
  if ($Lines -join "`n" -match [regex]::Escape($LineToInsert)) { return ,$Lines }
  for ($i=0; $i -lt $Lines.Count; $i++) {
    if ($Lines[$i] -match $AnchorRegex) {
      $indent = ($Lines[$i] -replace '^(\s*).*','$1')
      $Lines = $Lines[0..$i] + @("$indent$LineToInsert") + $Lines[($i+1)..($Lines.Count-1)]
      break
    }
  }
  ,$Lines
}

function Ensure-In-SliceBlock {
  param(
    [string[]]$Lines,
    [string]  $HeaderRegex,    # e.g. 'HTTP:\s*\[\]server_http\.HTTPRegister\s*\{'
    [string]  $ItemToInsert    # e.g. 'pruebav2HTTP,'
  )
  if ($Lines -join "`n" -match [regex]::Escape($ItemToInsert)) { return ,$Lines }

  for ($i=0; $i -lt $Lines.Count; $i++) {
    if ($Lines[$i] -match $HeaderRegex) {
      $headerIndent = ($Lines[$i] -replace '^(\s*).*','$1')
      # find end '},' of this slice
      $end = $null
      for ($j=$i+1; $j -lt $Lines.Count; $j++) {
        if ($Lines[$j] -match '^\s*\},\s*$') { $end = $j; break }
      }
      if ($null -eq $end) { break } # malformed, skip quietly
      $itemIndent = $headerIndent + "`t"
      $Lines = $Lines[0..($end-1)] + @("$itemIndent$ItemToInsert") + $Lines[$end..($Lines.Count-1)]
      break
    }
  }
  ,$Lines
}

# ========== 1) Update registers_agg.go ==========
if (-not (Test-Path -LiteralPath $RegistersAgg)) {
  throw "File not found: $RegistersAgg"
}
$raRaw   = Get-Content -LiteralPath $RegistersAgg -Raw -Encoding UTF8
$raLines = $raRaw -split "`r?`n"

# 1a) add params into BuildAllRegistrars signature via comments anchors
$httpParam = "$alias`HTTP server_http.HTTPRegister,"
$grpcParam = "$alias`GRPC server_grpc.GRPCRegister,"

$raLines = (Ensure-Line-AfterAnchor -Lines $raLines -AnchorRegex '^[ \t]*// add other HTTP-registrers for modules here:.*$' -LineToInsert $httpParam)
$raLines = (Ensure-Line-AfterAnchor -Lines $raLines -AnchorRegex '^[ \t]*// add other gRPC-registrers for modules here:.*$'  -LineToInsert $grpcParam)

# 1b) add items into slices HTTP:{...} and GRPC:{...}
$raLines = (Ensure-In-SliceBlock -Lines $raLines -HeaderRegex 'HTTP:\s*\[\]server_http\.HTTPRegister\s*\{' -ItemToInsert "$alias`HTTP,")
$raLines = (Ensure-In-SliceBlock -Lines $raLines -HeaderRegex 'GRPC:\s*\[\]server_grpc\.GRPCRegister\s*\{' -ItemToInsert "$alias`GRPC,")

[IO.File]::WriteAllLines($RegistersAgg, $raLines, $utf8NoBom)

# ========== 2) Update cmd/service/wire.go ==========
if (-not (Test-Path -LiteralPath $MainWire)) {
  throw "File not found: $MainWire"
}
$mwRaw   = Get-Content -LiteralPath $MainWire -Raw -Encoding UTF8
$mwLines = $mwRaw -split "`r?`n"

# 2a) ensure import alias in import (...)
$importFound = $false
foreach ($ln in $mwLines) {
  if ($ln -match '^\s*' + [regex]::Escape($alias) + '\s+"'+[regex]::Escape($importPath)+'"$') { $importFound = $true; break }
}
if (-not $importFound) {
  $importStart = $mwLines | Select-String -Pattern '^\s*import\s*\($' | Select-Object -First 1
  if ($importStart) {
    $idxStart = $importStart.LineNumber - 1
    # find closing ')'
    $idxEnd = ($mwLines | Select-String -Pattern '^\s*\)$' | Where-Object { $_.LineNumber -gt $idxStart } | Select-Object -First 1).LineNumber - 1
    if ($idxEnd -ge 0) {
      $mwLines = $mwLines[0..($idxEnd-1)] + @("`t$alias `"$importPath`"") + $mwLines[$idxEnd..($mwLines.Count-1)]
    }
  } else {
    # no import block -> create right after package line
    $pkgIdx = ($mwLines | Select-String -Pattern '^\s*package\s+\w+\s*$' | Select-Object -First 1).LineNumber - 1
    if ($pkgIdx -ge 0) {
      $insert = @(
        ""
        "import ("
        "`t$alias `"$importPath`""
        ")"
        ""
      )
      $mwLines = $mwLines[0..$pkgIdx] + $insert + $mwLines[($pkgIdx+1)..($mwLines.Count-1)]
    }
  }
}

# 2b) ensure <alias>.ProviderSet in // modules section of wire.Build
$provLine = "$alias.ProviderSet,"
$haveProv = ($mwLines -join "`n") -match [regex]::Escape($provLine)
if (-not $haveProv) {
  # try to inject after line with // modules
  $modulesIdx = ($mwLines | Select-String -Pattern '^\s*//\s*modules\s*$' | Select-Object -First 1).LineNumber - 1
  if ($modulesIdx -ge 0) {
    # insert after the comment, keeping same indent + one tab
    $indent = ($mwLines[$modulesIdx] -replace '^(\s*).*','$1')
    $mwLines = $mwLines[0..$modulesIdx] + @("$indent`t$provLine") + $mwLines[($modulesIdx+1)..($mwLines.Count-1)]
  } else {
    # fallback: inside wire.Build(...) before 'newApp,'
    $buildStart = ($mwLines | Select-String -Pattern 'wire\.Build\(' | Select-Object -First 1).LineNumber - 1
    $newAppIdx  = ($mwLines | Select-String -Pattern '^\s*newApp,\s*$' | Where-Object { $_.LineNumber -gt $buildStart } | Select-Object -First 1).LineNumber - 1
    if ($buildStart -ge 0 -and $newAppIdx -ge 0) {
      $indent = ($mwLines[$newAppIdx] -replace '^(\s*).*','$1')
      $mwLines = $mwLines[0..($newAppIdx-1)] + @("$indent$provLine") + $mwLines[$newAppIdx..($mwLines.Count-1)]
    }
  }
}

[IO.File]::WriteAllLines($MainWire, $mwLines, $utf8NoBom)

Write-Host ("wire updated for module '{0}' (v{1}): alias {2}" -f $base, $V, $alias) -ForegroundColor Green
