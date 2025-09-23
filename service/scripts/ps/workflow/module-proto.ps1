# scripts/ps/workflow/module-proto.ps1
[CmdletBinding()]
param(
  [Parameter(Mandatory = $true)] [string]$Name,
  [string]$Ops = ""
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
$ApiRoot   = "./api"
$utf8NoBom = New-Object System.Text.UTF8Encoding($false)

# --- helpers ---
function ConvertTo-PascalCase {
  param([Parameter(Mandatory=$true)][string]$InputString)
  $parts = ($InputString -replace '[^A-Za-z0-9]+',' ') -split '\s+' | Where-Object { $_ }
  ($parts | ForEach-Object { $_.Substring(0,1).ToUpper() + $_.Substring(1).ToLower() }) -join ''
}
function ConvertTo-LowerCase { param([Parameter(Mandatory=$true)][string]$InputString) $InputString.ToLower() }
function ConvertTo-Plural    { param([Parameter(Mandatory=$true)][string]$InputNoun) if ($InputNoun.ToLower().EndsWith('s')){$InputNoun}else{"$InputNoun"+"s"} }

# --- parse ops -> flags ---
$opsList = @()
if ($Ops) { $opsList = $Ops.ToLower().Split(',') | ForEach-Object { $_.Trim() } | Where-Object { $_ } }
$HasGet = $false; $HasUpsert = $false; $HasDelete = $false
foreach ($op in $opsList) {
  switch ($op) {
    'get'    { $HasGet = $true }
    'find'   { $HasGet = $true }
    'list'   { $HasGet = $true }
    'read'   { $HasGet = $true }
    'upsert' { $HasUpsert = $true }
    'create' { $HasUpsert = $true }
    'update' { $HasUpsert = $true }
    'delete' { $HasDelete = $true }
    'del'    { $HasDelete = $true }
    'remove' { $HasDelete = $true }
    default  { Show-Warn "Unknown op '$op' ignored" }
  }
}
$AnyOps       = $HasGet -or $HasUpsert -or $HasDelete
$GenerateMock = -not $AnyOps   # если нет ни одной операции — добавляем мок

Show-Step "Generating .proto module"
Show-Info "Input name: $Name ; ops=[$Ops]"

# --- normalize names ---
$base         = ConvertTo-LowerCase $Name
$pascal       = ConvertTo-PascalCase $Name
$pluralBase   = ConvertTo-LowerCase (ConvertTo-Plural $base)
$pluralPascal = ConvertTo-PascalCase $pluralBase

# --- detect next free version ---
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

# --- paths & meta ---
$pkgDir      = Join-Path -Path $baseDir -ChildPath "v$Version"
$protoFile   = Join-Path $pkgDir "$base.proto"

$package     = "api.$base.v$Version"
$goPkg       = "service/api/$base;$base"
$javaOuter   = "$($pascal)ProtoV$Version"
$javaPkg     = "dev.kratos.api.$base.$base"

$serviceName = "${pascal}v${Version}Service"
$route       = "/v$Version/$base"

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

# --- safety ---
if (Test-Path $protoFile) {
  Show-ErrorAndExit "File already exists for '$base' v$Version. Aborting."
}

# --- dynamic imports ---
$importLines = @()
$importLines += 'import "google/protobuf/timestamp.proto";'
$importLines += 'import "api/common/v1/common.proto";'        # <<< общий Meta
if ($AnyOps -or $GenerateMock) {
  $importLines += 'import "google/api/annotations.proto";'
  $importLines += 'import "google/api/field_behavior.proto";'
}
$importsBlock = ($importLines -join "`n")

# fq тип Meta
$MetaFQ = "api.common.v1.MetaResponse"

# --- service methods & messages ---
$serviceMethods = @()
$messages = @()

if ($HasGet) {
  $serviceMethods += @"
  // GET ${route} - list or search by filters (query: id OR name)
  rpc Find${pluralPascal}(Find${pluralPascal}Request) returns (Find${pluralPascal}Response) {
    option (google.api.http) = { get: "${route}" };
  }
"@

  $messages += @"
message Find${pluralPascal}Request {
  // oneof можно добавить при желании строгости; пока просто опциональные поля
  optional uint32 id   = 1 [(google.api.field_behavior) = OPTIONAL];
  optional string name = 2 [(google.api.field_behavior) = OPTIONAL];

  // optional пагинация (раскомментируй при нужде)
  // optional uint32 limit  = 10;
  // optional uint32 offset = 11;
}
message Find${pluralPascal}Response {
  repeated $pascal items = 1;
  $MetaFQ meta = 2;
}
"@
}

if ($HasUpsert) {
  $serviceMethods += @"
  // POST ${route} - create or update (id empty/0 => create, >0 => update)
  rpc Upsert${pascal}(Upsert${pascal}Request) returns (Upsert${pascal}Response) {
    option (google.api.http) = {
      post: "${route}"
      body: "*"
    };
  }
"@

  $messages += @"
message Upsert${pascal}Request {
  optional uint32 id   = 1; // 0 or unset => create; >0 => update
  string name          = 2 [(google.api.field_behavior) = REQUIRED];
}
message Upsert${pascal}Response {
  $pascal item = 1;
  $MetaFQ meta = 2;
}
"@
}

if ($HasDelete) {
  $serviceMethods += @"
  // DELETE ${route}?id=123 - delete by id (query param)
  rpc Delete${pascal}ById(Delete${pascal}ByIdRequest) returns (Delete${pascal}ByIdResponse) {
    option (google.api.http) = { delete: "${route}" };
  }
"@

  $messages += @"
message Delete${pascal}ByIdRequest {
  uint32 id = 1 [(google.api.field_behavior) = REQUIRED];
}
message Delete${pascal}ByIdResponse {
  $MetaFQ meta = 1;
}
"@
}

# --- mock endpoint when no ops selected ---
if ($GenerateMock) {
  $serviceMethods += @"
// Mock endpoint (no ops selected)
  rpc Mock(MockRequest) returns (MockResponse) {
    option (google.api.http) = { get: "${route}/mock" };
  }
"@

  $messages += @"
message MockRequest {}
message MockResponse {
  string message = 1;      // e.g. ""pong""
  $MetaFQ meta   = 2;
}
"@
}

# --- service block ---
$serviceBlock = ""
if ($serviceMethods.Count -gt 0) {
  $serviceBlock = @"
// ${serviceName}: generated with ops=[$Ops]
service $serviceName {
$($serviceMethods -join "`n")
}
"@
} else {
  $serviceBlock = @"
// ${serviceName}: no RPCs
service $serviceName {}
"@
}

$entityMessage = @"
message $pascal {
  uint32 id = 1;
  string name = 2;
  google.protobuf.Timestamp created_at = 3;
  google.protobuf.Timestamp updated_at = 4;
}
"@

$messagesText = ""
if ($messages.Count -gt 0) {
  $messagesText = ($messages -join "`n")
}

# --- final proto content ---
$protoContent = @"
syntax = "proto3";

package $package;

$importsBlock

option go_package = "$goPkg";
option java_multiple_files = true;
option java_outer_classname = "$javaOuter";
option java_package = "$javaPkg";

$serviceBlock

$messagesText

$entityMessage
"@

# --- write file ---
Show-Step "Writing files"
[IO.File]::WriteAllText($protoFile,  $protoContent,  $utf8NoBom)
Show-OK "$base.proto -> $protoFile"

Show-Step "Done"
Show-OK ("Created {0}/v{1}" -f $base, $Version)
Show-Info "Package:  $package"
Show-Info "Service:  $serviceName"
Show-Info "Route:    $route"
