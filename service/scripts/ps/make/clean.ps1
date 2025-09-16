param(
    [string]$BinDir = "bin"
)

. "$PSScriptRoot/utils.ps1"

Show-Step "Cleaning project"

# --- bin/
if (Test-Path $BinDir) {
    Remove-Item -Recurse -Force $BinDir
    Show-Info "Removed $BinDir"
} else {
    Show-Info "$BinDir not found, skipping"
}

# --- wire_gen.go
$wireFiles = Get-ChildItem -Recurse -Filter 'wire_gen.go' -ErrorAction SilentlyContinue
if ($wireFiles) {
    $wireFiles | Remove-Item -Force -ErrorAction SilentlyContinue
    Show-Info "Removed wire_gen.go files"
} else {
    Show-Info "No wire_gen.go found"
}

# --- *.pb.go
$pbFiles = Get-ChildItem -Recurse -Filter '*.pb.go' -ErrorAction SilentlyContinue
if ($pbFiles) {
    $pbFiles | Remove-Item -Force -ErrorAction SilentlyContinue
    Show-Info "Removed *.pb.go files"
} else {
    Show-Info "No *.pb.go found"
}

# --- openapi.yaml
$openapiFiles = Get-ChildItem -Recurse -Filter 'openapi.yaml' -ErrorAction SilentlyContinue
if ($openapiFiles) {
    $openapiFiles | Remove-Item -Force -ErrorAction SilentlyContinue
    Show-Info "Removed openapi.yaml files"
} else {
    Show-Info "No openapi.yaml found"
}

# --- service.swagger.json
$swaggerFiles = Get-ChildItem -Recurse -Filter 'service.swagger.json' -ErrorAction SilentlyContinue
if ($swaggerFiles) {
    $swaggerFiles | Remove-Item -Force -ErrorAction SilentlyContinue
    Show-Info "Removed service.swagger.json files"
} else {
    Show-Info "No service.swagger.json found"
}

Show-OK "Project cleaned successfully"
