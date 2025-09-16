param([string]$AppName = "service")
. "$PSScriptRoot/utils.ps1"

Show-Step "Installing tools for $AppName"

function Install-Tool {
    param(
        [string]$Name,
        [string]$Command
    )
    Show-Info "Installing $Name..."
    Invoke-Expression $Command
    if ($LASTEXITCODE -ne 0) {
        Show-ErrorAndExit "$Name install failed"
    }
    Show-OK "$Name installed"
}

# --- Go toolchain check ---
$goVersion = (go version) 2>$null
if (-not $goVersion) {
    Show-ErrorAndExit "Go not found! Please install Go first."
}
Show-Info "Found Go: $goVersion"

# --- Install tools ---
Install-Tool "protoc-gen-go"          'go install google.golang.org/protobuf/cmd/protoc-gen-go@latest'
Install-Tool "protoc-gen-go-grpc"     'go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest'
Install-Tool "protoc-gen-go-http"     'go install github.com/go-kratos/kratos/cmd/protoc-gen-go-http/v2@latest'
Install-Tool "protoc-gen-grpc-gateway"'go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest'
Install-Tool "protoc-gen-openapiv2"   'go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest'
Install-Tool "protoc-gen-openapi"     'go install github.com/google/gnostic/cmd/protoc-gen-openapi@latest'
Install-Tool "wire"                   'go install github.com/google/wire/cmd/wire@latest'

# --- Go mod tidy ---
Show-Step "Tidying dependencies"
go mod tidy
if ($LASTEXITCODE -ne 0) {
    Show-ErrorAndExit "go mod tidy failed"
}
Show-OK "Dependencies tidied"

Show-OK "All tools installed successfully for $AppName!"
