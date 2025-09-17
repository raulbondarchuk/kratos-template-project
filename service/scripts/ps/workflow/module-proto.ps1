# scripts/ps/workflow/module-proto.ps1
[CmdletBinding()]
param(
  [Parameter(Mandatory = $true)] [string]$Name
)

Set-StrictMode -Version Latest
$ErrorActionPreference = 'Stop'

# --- utils (logs) ---
try {
  . (Join-Path $PSScriptRoot 'utils.ps1')
} catch {
  function Show-Step { param([string]$Message) Write-Host "`n==> $Message" -ForegroundColor Cyan }
  function Show-Info { param([string]$Message) Write-Host "$Message" -ForegroundColor DarkGray }
  function Show-OK   { param([string]$Message) Write-Host "  [OK] $Message" -ForegroundColor Green }
  function Show-Warn { param([string]$Message) Write-Host "  [WARN] $Message" -ForegroundColor Yellow }
  function Show-ErrorAndExit { param([string]$Message) Write-Host "  [ERROR] $Message" -ForegroundColor Red; exit 1 }
}

# consts
$ApiRoot = "./api"
$utf8NoBom = New-Object System.Text.UTF8Encoding($false)

# --- helpers ---
function ConvertTo-PascalCase {
  [CmdletBinding()]
  param(
    [Parameter(Mandatory=$true, Position=0)]
    [ValidateNotNullOrEmpty()]
    [string]$InputString
  )
  $parts = ($InputString -replace '[^A-Za-z0-9]+',' ') -split '\s+' | Where-Object { $_ }
  ($parts | ForEach-Object { $_.Substring(0,1).ToUpper() + $_.Substring(1).ToLower() }) -join ''
}

function ConvertTo-LowerCase {
  [CmdletBinding()]
  param(
    [Parameter(Mandatory=$true, Position=0)]
    [ValidateNotNullOrEmpty()]
    [string]$InputString
  )
  $InputString.ToLower()
}

function ConvertTo-Plural {
  [CmdletBinding()]
  param(
    [Parameter(Mandatory=$true, Position=0)]
    [ValidateNotNullOrEmpty()]
    [string]$InputNoun
  )
  if ($InputNoun.ToLower().EndsWith('s')) { 
    $InputNoun 
  } else { 
    "$InputNoun" + "s" 
  }
}

Show-Step "Generating .proto module"
Show-Info "Input name: $Name"

# --- normalize names ---
$base         = ConvertTo-LowerCase $Name
$pascal       = ConvertTo-PascalCase $Name
$pluralBase   = ConvertTo-LowerCase (ConvertTo-Plural $base)
$pluralPascal = ConvertTo-PascalCase $pluralBase

Show-Info "Normalized: base='$base', pascal='$pascal'"

# --- detect next available version: v1, v2, v3... (first gap) ---
Show-Step "Detecting next free API version"
$baseDir  = Join-Path -Path $ApiRoot -ChildPath $base
$versions = @()
if (Test-Path -LiteralPath $baseDir) {
  Get-ChildItem -LiteralPath $baseDir -Directory -ErrorAction SilentlyContinue | ForEach-Object {
    if ($_.Name -match '^v(\d+)$') { $versions += [int]$Matches[1] }
  }
    if ($versions.Count -gt 0) {
      $versionList = $versions | Sort-Object | ForEach-Object { "v$_" }
      Show-Info ("Existing versions: " + ($versionList -join ", "))
    } else {
      Show-Info "Existing versions: none"
    }
} else {
  Show-Warn "API base dir not found: $baseDir (will create)"
}
$Version = 1
while ($versions -contains $Version) { $Version++ }
Show-Info "Chosen version: v$Version"

# --- paths & meta for chosen version ---
$pkgDir     = Join-Path -Path $baseDir -ChildPath "v$Version"
$protoFile  = Join-Path $pkgDir "$base.proto"
$errorsFile = Join-Path $pkgDir "errors.proto"

$package    = "api.$base.v$Version"
$goPkg      = "service/api/$base;$base"
$javaOuter  = "$($pascal)ProtoV$Version"
$javaPkg    = "dev.kratos.api.$base.$base"

# Service name must include version suffix for Swagger isolation
$serviceName = "${pascal}v${Version}Service"

# IMPORTANT: versioned route (singular)
$route        = "/v$Version/$base"
$errorsImport = "api/$base/v$Version/errors.proto"

Show-Step "Preparing output"
Show-Info "Package: $package"
Show-Info "Service: $serviceName"
Show-Info "Route:   $route"
Show-Info "Out dir: $pkgDir"

# --- create folders ---
if (-not (Test-Path -LiteralPath $pkgDir)) {
  New-Item -ItemType Directory -Path $pkgDir -Force | Out-Null
  Show-Info "Created directory: $pkgDir"
} else {
  Show-Info "Directory exists: $pkgDir"
}

# --- safety: should not hit because we pick a free version ---
if ((Test-Path $protoFile) -or (Test-Path $errorsFile)) {
  Show-ErrorAndExit "Files already exist for '$base' v$Version. Aborting."
}

# --- errors.proto content ---
$errorsContent = @"
syntax = "proto3";

package $package;

option go_package = "$goPkg";

// Minimum set of codes. Zero is mandatory in proto3.
enum ResponseCode {
  RESPONSE_CODE_UNSPECIFIED = 0;
  RESPONSE_CODE_OK = 200;                    // success
  RESPONSE_CODE_BAD_REQUEST = 400;           // bad request
  RESPONSE_CODE_UNAUTHORIZED = 401;          // unauthorized
  RESPONSE_CODE_FORBIDDEN = 403;             // forbidden
  RESPONSE_CODE_NOT_FOUND = 404;             // not found
  RESPONSE_CODE_METHOD_NOT_ALLOWED = 405;    // method not allowed
  RESPONSE_CODE_INTERNAL_SERVER_ERROR = 500; // server error
  RESPONSE_CODE_NOT_IMPLEMENTED = 501;       // not implemented
  RESPONSE_CODE_SERVICE_UNAVAILABLE = 503;   // service unavailable
}

// Your meta-object: only code and message.
message MetaResponse {
  ResponseCode code = 1; // RESPONSE_CODE_OK or other codes
  string message = 2;    // e.g. "ok", "ip blocked", "db error", "user not found"
}
"@

# --- <name>.proto content (minimal comments) ---
$protoContent = @"
syntax = "proto3";

package $package;

import "$errorsImport";
import "google/api/annotations.proto";
import "google/api/field_behavior.proto";
import "google/protobuf/timestamp.proto";

option go_package = "$goPkg";
option java_multiple_files = true;
option java_outer_classname = "$javaOuter";
option java_package = "$javaPkg";

// ${serviceName}: find, upsert, delete
service $serviceName {
  // GET ${route} - list or search by filters
  rpc Find$pluralPascal(Find${pluralPascal}Request) returns (Find${pluralPascal}Response) {
    option (google.api.http) = { get: "${route}" };
  }

  // POST ${route} - create or update (id=0 create, >0 update)
  rpc Upsert$pascal(Upsert${pascal}Request) returns (Upsert${pascal}Response) {
    option (google.api.http) = {
      post: "${route}"
      body: "*"
    };
  }

  // DELETE ${route}?id=123 - delete by id
  rpc Delete${pascal}ById(Delete${pascal}ByIdRequest) returns (Delete${pascal}ByIdResponse) {
    option (google.api.http) = { delete: "${route}" };
  }
}

message Find${pluralPascal}Request {
  uint32 id   = 1 [(google.api.field_behavior) = OPTIONAL]; // optional
  string name = 2 [(google.api.field_behavior) = OPTIONAL]; // optional
}

message Find${pluralPascal}Response {
  repeated $pascal items = 1;
  MetaResponse meta = 2;
}

message Upsert${pascal}Request {
  uint32 id   = 1 [(google.api.field_behavior) = OPTIONAL]; // 0 -> create
  string name = 2 [(google.api.field_behavior) = REQUIRED]; // required
}

message Upsert${pascal}Response {
  $pascal item = 1;
  MetaResponse meta = 2;
}

message Delete${pascal}ByIdRequest {
  uint32 id = 1 [(google.api.field_behavior) = REQUIRED]; // required
}

message Delete${pascal}ByIdResponse {
  MetaResponse meta = 1;
}

message $pascal {
  uint32 id = 1;
  string name = 2;
  google.protobuf.Timestamp created_at = 3;
  google.protobuf.Timestamp updated_at = 4;
}
"@

# --- write files (UTF-8 no BOM) ---
Show-Step "Writing files"
[IO.File]::WriteAllText($errorsFile, $errorsContent, $utf8NoBom)
Show-OK "errors.proto -> $errorsFile"

[IO.File]::WriteAllText($protoFile,  $protoContent,  $utf8NoBom)
Show-OK "$base.proto   -> $protoFile"

Show-Step "Done"
Show-OK ("Created {0}/v{1}" -f $base, $Version)
Show-Info "Package:  $package"
Show-Info "Service:  $serviceName"
Show-Info "Route:    $route"
