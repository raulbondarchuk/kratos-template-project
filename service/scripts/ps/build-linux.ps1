param(
    [string]$AppName = "service",
    [string]$CmdDir  = "cmd/service",
    [string]$BinDir  = "bin"
)

# Folder with date: dd-MM-yyyy__HH-mm-ss (without colons)
$stamp   = Get-Date -Format 'dd-MM-yyyy__HH-mm-ss'
$target  = Join-Path $BinDir $stamp
$binPath = Join-Path $target ($AppName + ".linux")

if (-not (Test-Path $target)) {
    New-Item -ItemType Directory -Force -Path $target | Out-Null
}

# Version from git (with fallback)
$version = git describe --tags --always 2> $null
if (-not $version) { $version = "0.0.0-local" }

# Cross-compilation for Linux (without CGO)
$env:GOOS        = "linux"
$env:GOARCH      = "amd64"
$env:CGO_ENABLED = "0"

# Build
go build -ldflags ("-X main.Version=" + $version) -o $binPath "./$CmdDir"
if ($LASTEXITCODE -ne 0) {
    throw "go build failed (exit $LASTEXITCODE)"
}

# Copy config.yaml
$cfg = Join-Path "configs" "config.yaml"
if (Test-Path $cfg) {
    Copy-Item -Force $cfg $target
    Write-Host "Copied:  $cfg -> $target"
} else {
    throw "configs/config.yaml not found"
}

# Copy .env if exists
$envFile = ".env"
if (Test-Path $envFile) {
    Copy-Item -Force $envFile $target
    Write-Host "Copied:  $envFile -> $target"
} else {
    Write-Host "⚠️  .env file not found, skipping"
}

Write-Host "Built:   $binPath"
