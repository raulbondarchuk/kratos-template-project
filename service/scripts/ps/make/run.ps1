param(
    [ValidateSet("kratos", "go")]
    [string]$Mode = "kratos", 
    [string]$CmdDir = "cmd/service",
    [string]$ConfigPath = "configs"
)

. "$PSScriptRoot/utils.ps1"

if ($Mode -eq "kratos") {
    Show-Step "Running with kratos run"
    kratos run
    if ($LASTEXITCODE -ne 0) { Show-ErrorAndExit "kratos run failed" }
    Show-OK "kratos run finished"
}
elseif ($Mode -eq "go") {
    Show-Step "Running with go run"
    go run "./$CmdDir" -conf "./$ConfigPath"
    if ($LASTEXITCODE -ne 0) { Show-ErrorAndExit "go run failed" }
    Show-OK "go run finished"
}
