function Show-Step {
    param([string]$Message)
    Write-Host ""
    Write-Host "==> $Message" -ForegroundColor Cyan
}

function Show-Info {
    param([string]$Message)
    Write-Host "$Message" -ForegroundColor DarkGray
}

function Show-OK {
    param([string]$Message)
    Write-Host "  [OK] $Message" -ForegroundColor Green
}

function Show-Warn {
    param([string]$Message)
    Write-Host "  [WARN] $Message" -ForegroundColor Yellow
}

function Show-ErrorAndExit {
    param([string]$Message)
    Write-Host "  [ERROR] $Message" -ForegroundColor Red
    exit 1
}