param(
    [string]$AppName = "service",
    [string]$CmdDir  = "cmd/service",
    [string]$BinDir  = "bin"
)

. "$PSScriptRoot/make/utils.ps1"

# === Target folders ===
$stamp   = Get-Date -Format 'dd-MM-yyyy__HH-mm-ss'
$target  = Join-Path $BinDir $stamp
$binPath = Join-Path $target ($AppName + ".linux")

Show-Step "Build initiated"

Show-Info "Preparing build target folder: $target"
if (-not (Test-Path $target)) {
    New-Item -ItemType Directory -Force -Path $target | Out-Null
    Show-Info "Created target directory $target"
}

# === Version ===
Show-Info "Resolving version from git"
$version = git describe --tags --always 2> $null
if (-not $version) {
    $version = "0.0.0-local"
    Show-Info "Git tags not found, fallback version $version"
} else {
    Show-Info "Version detected: $version"
}

# === Cross-compilation ===
Show-Info "Configuring cross-compilation environment"
$env:GOOS        = "linux"
$env:GOARCH      = "amd64"
$env:CGO_ENABLED = "0"
Show-Info "GOOS=$env:GOOS, GOARCH=$env:GOARCH, CGO_ENABLED=$env:CGO_ENABLED"

# === Build ===
Show-Info "Building $AppName for Linux"
go build -ldflags ("-X main.Version=" + $version) -o $binPath "./$CmdDir"
if ($LASTEXITCODE -ne 0) {
    Show-ErrorAndExit "go build failed (exit $LASTEXITCODE)"
}
Show-Info "Build complete -> $binPath"

# === Copy config.yaml ===
$cfg = Join-Path "configs" "config.yaml"
$cfgTarget = Join-Path $target "config.yaml"
Show-Info "Copying config.yaml"
if (Test-Path $cfg) {
    Copy-Item -Force $cfg $cfgTarget
    Show-Info "Copied:  $cfg -> $target"

    # === Update app.version ===
    Show-Info "Updating app.version in config.yaml -> $version"
    $yaml = Get-Content $cfgTarget -Raw -Encoding UTF8

    if ($yaml -match '(?m)^app:') {
        if ($yaml -match '(?m)^\s*version:\s*.+$') {
            # заменяем существующую строку version:
            $yaml = $yaml -replace '(?m)^\s*version:\s*.+$', "  version: $version"
        } else {
            # если version: нет, добавляем сразу после app:
            $yaml = $yaml -replace '(?m)^app:\s*$', "app:`r`n  version: $version"
        }
    } else {
        # если секции app: вообще нет, добавляем в конец
        $yaml += "`r`napp:`r`n  version: $version`r`n"
    }

    [System.IO.File]::WriteAllText($cfgTarget, $yaml, (New-Object System.Text.UTF8Encoding($false)))
    Show-Info "Updated app.version to $version in $cfgTarget"

} else {
    Show-ErrorAndExit "configs/config.yaml not found"
}

# === Copy .env ===
$envFile = ".env"
Show-Info "Copying .env (if exists)"
if (Test-Path $envFile) {
    Copy-Item -Force $envFile $target
    Show-Info "Copied:  $envFile -> $target"
} else {
    Show-Info "⚠️  .env file not found, skipping"
}

Show-OK "Linux build finished successfully"
