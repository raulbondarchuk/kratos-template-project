param(
    [string]$BufYaml = "buf.yaml",
    [string]$BufLock = "buf.lock"
)

. "$PSScriptRoot/utils.ps1"

Show-Step "Checking buf.lock freshness"

if (-not (Test-Path $BufLock) -or ((Get-Item $BufYaml).LastWriteTime -gt (Get-Item $BufLock).LastWriteTime)) {
    Show-Info "buf dep update required"
    buf dep update
    if ($LASTEXITCODE -ne 0) { Show-ErrorAndExit "buf dep update failed" }
    Show-OK "Dependencies updated"
} else {
    Show-Info "skip buf dep update (lock is fresh)"
    Show-OK "Dependencies already fresh"
}