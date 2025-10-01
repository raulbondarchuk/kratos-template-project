[CmdletBinding()]
param()  # without parameters — always full rebuild

Set-StrictMode -Version Latest
$ErrorActionPreference = 'Stop'

. "$PSScriptRoot/utils.ps1"

# Check if there are any .proto files in api directory
$protoFiles = @(Get-ChildItem -Path "api" -Recurse -Filter "*.proto" -ErrorAction SilentlyContinue)
$protoCount = $protoFiles.Count

if ($protoCount -eq 0) {
  Show-Info "No .proto files found in api directory. Nothing to generate."
  exit 0
}

if (-not (Get-Command buf -ErrorAction SilentlyContinue)) {
  Show-ErrorAndExit "buf tool not found in PATH"
}

Show-Step "Rebuilding ALL API documentation"
Show-Info ("Found {0} .proto files in api directory" -f $protoCount)

# Select template: first buf.gen.docs.yaml, then buf.gen.yaml
$tplCandidates = @("buf.gen.docs.yaml", "buf.gen.yaml")
$tpl = $tplCandidates | Where-Object { Test-Path $_ } | Select-Object -First 1
if (-not $tpl) {
  Show-ErrorAndExit "Template not found: neither buf.gen.docs.yaml, nor buf.gen.yaml"
}
$absTpl = (Resolve-Path $tpl).Path
Show-Info "Using template: $absTpl"

# --- Clean docs/openapi (can be deleted completely) ---
if (Test-Path "docs/openapi") {
  Show-Info "Cleaning folder: docs/openapi"
  Remove-Item -Recurse -Force "docs/openapi"
  Show-Info   "Removed: docs/openapi"
}

# --- Careful cleaning docs (preserve openapi_embed.go and others) ---
if (Test-Path "docs") {
  $preserve = @("openapi_embed.go", ".gitkeep")
  $extsGenerated = @(".json", ".yaml", ".yml")

  Show-Info "Cleaning generated files in docs (preserve: $($preserve -join ', '))"

  $toDelete = Get-ChildItem "docs" -Recurse -File -ErrorAction SilentlyContinue | Where-Object {
    ($preserve -notcontains $_.Name) -and ($extsGenerated -contains $_.Extension.ToLower())
  }

  $deleted = 0
  foreach ($f in $toDelete) {
    Remove-Item -Force $f.FullName
    $deleted++
  }
  Show-Info "Removed files in docs: $deleted"

  # remove empty folders inside docs (docs itself is not touched)
  $dirsRemoved = 0
  Get-ChildItem "docs" -Recurse -Directory -ErrorAction SilentlyContinue |
    Sort-Object FullName -Descending | ForEach-Object {
      $hasFiles = (Get-ChildItem $_.FullName -Recurse -File -ErrorAction SilentlyContinue | Measure-Object).Count -gt 0
      if (-not $hasFiles) {
        Remove-Item -Recurse -Force $_.FullName
        $dirsRemoved++
      }
    }
  if ($dirsRemoved -gt 0) { Show-Info "Removed empty directories under docs: $dirsRemoved" }
} else {
  # if folder docs is not found — it's not a problem
  Show-Info "Folder 'docs' not found, nothing to clean there."
}

# --- Generation (inline-template — more reliable on Windows) ---
$templateData = Get-Content -Raw -Encoding UTF8 $tpl
buf generate --template $templateData
if ($LASTEXITCODE -ne 0) { Show-ErrorAndExit "buf generate failed" }

# --- Report ---
if (Test-Path "docs/openapi") {
  $files = Get-ChildItem "docs/openapi" -Recurse -File -ErrorAction SilentlyContinue | Select-Object -Expand FullName | Sort-Object
  if ($files) { Show-Info ("OpenAPI v2 files:`n  " + ($files -join "`n  ")) } else { Show-Warn "docs/openapi is empty" }
}
if (Test-Path "docs") {
  $files = Get-ChildItem "docs" -Recurse -File -ErrorAction SilentlyContinue |
           Where-Object { $_.FullName -notmatch '\\docs\\openapi(\\|$)' } |
           Select-Object -Expand FullName | Sort-Object
  if ($files) { Show-Info ("OpenAPI files:`n  " + ($files -join "`n  ")) }
}

Show-OK "Documentation generation complete"
