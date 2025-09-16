param(
    [string]$CmdDir = "cmd/service"
)

. "$PSScriptRoot/utils.ps1"

Show-Step "Running wire in $CmdDir"

if (-not (Test-Path $CmdDir)) {
    Show-ErrorAndExit "CmdDir not found: $CmdDir"
}

Show-Info "Changing directory to $CmdDir"
Push-Location $CmdDir

Show-Info "Executing 'wire'"
wire
if ($LASTEXITCODE -ne 0) {
    Pop-Location
    Show-ErrorAndExit "wire failed"
}

Show-Info "Returning to project root"
Pop-Location

Show-OK "Wire code generated in $CmdDir"
