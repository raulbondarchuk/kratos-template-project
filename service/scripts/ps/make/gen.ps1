param(
    [string]$BufGen = "buf.gen.yaml"
)

. "$PSScriptRoot/utils.ps1"

Show-Step "Generating protobuf code from $BufGen"

if (-not (Test-Path $BufGen)) {
    Show-ErrorAndExit "Buf template not found: $BufGen"
}

# If default template in root — use simply `buf generate`
$leaf = Split-Path -Leaf $BufGen
$dir  = Split-Path -Parent $BufGen
$cwd  = (Get-Location).Path

if ($leaf -ieq "buf.gen.yaml" -and ((-not $dir) -or ((Resolve-Path $dir).Path -ieq (Resolve-Path $cwd).Path))) {
    Show-Info "Using default template (auto-detect)"
    buf generate
    if ($LASTEXITCODE -ne 0) { Show-ErrorAndExit "buf generate failed" }
    Show-OK "Protobuf code generated"
    return
}

# Otherwise: read file and pass as inline YAML (reliable on Windows)
Show-Info "Using inline template from file: $((Resolve-Path $BufGen).Path)"
$templateData = Get-Content -Raw -Encoding UTF8 $BufGen
buf generate --template $templateData
if ($LASTEXITCODE -ne 0) { Show-ErrorAndExit "buf generate failed" }

Show-OK "Protobuf code generated"