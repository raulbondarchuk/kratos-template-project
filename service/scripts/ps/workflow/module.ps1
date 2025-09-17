[CmdletBinding()]
param(
  [Parameter(Mandatory = $true)] [string]$Name,
  [int]$Version = 1,
  [string]$ApiRoot = "./api",
  [switch]$Force
)

Set-StrictMode -Version Latest
$ErrorActionPreference = 'Stop'

# --- helpers (Approved Verbs) ---
function ConvertTo-PascalCase {
  param([Parameter(Mandatory=$true)][string]$InputString)
  $parts = ($InputString -replace '[^A-Za-z0-9]+',' ') -split '\s+' | Where-Object { $_ }
  ($parts | ForEach-Object { $_.Substring(0,1).ToUpper() + $_.Substring(1).ToLower() }) -join ''
}

function ConvertTo-LowerCase {
  param([Parameter(Mandatory=$true)][string]$InputString)
  $InputString.ToLower()
}

function ConvertTo-Plural {
  param([Parameter(Mandatory=$true)][string]$InputNoun)
  if ($InputNoun.ToLower().EndsWith('s')) { return $InputNoun }
  return "$InputNoun" + "s"
}

$utf8NoBom = New-Object System.Text.UTF8Encoding($false)

# --- derive names/paths ---
$base         = ConvertTo-LowerCase $Name
$pascal       = ConvertTo-PascalCase $Name
$pluralBase   = ConvertTo-LowerCase (ConvertTo-Plural $base)
$pluralPascal = ConvertTo-PascalCase $pluralBase

$pkgDir     = Join-Path -Path $ApiRoot -ChildPath "$base/v$Version"
$protoFile  = Join-Path $pkgDir "$base.proto"
$errorsFile = Join-Path $pkgDir "errors.proto"

$package   = "api.$base.v$Version"
$goPkg     = "service/api/$base;$base"
$javaOuter = "$($pascal)ProtoV$Version"
$javaPkg   = "dev.kratos.api.$base.$base"

$route        = "/$pluralBase"
$errorsImport = "api/$base/v$Version/errors.proto"

# --- create folder ---
if (-not (Test-Path -LiteralPath $pkgDir)) {
  New-Item -ItemType Directory -Path $pkgDir | Out-Null
}

# --- guard overwrite ---
if (-not $Force) {
  if (Test-Path $protoFile)  { throw "File exists: $protoFile (use -Force to overwrite)" }
  if (Test-Path $errorsFile) { throw "File exists: $errorsFile (use -Force to overwrite)" }
}

# --- errors.proto content (keep) ---
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

// $pascal Service: find, upsert, delete
service ${pascal}Service {
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
[IO.File]::WriteAllText($errorsFile, $errorsContent, $utf8NoBom)
[IO.File]::WriteAllText($protoFile,  $protoContent,  $utf8NoBom)

Write-Host "Created:" -ForegroundColor Green
Write-Host "  $errorsFile"
Write-Host "  $protoFile"

# --- buf format (module at "." -> filter with --path) ---
$buf = Get-Command buf -ErrorAction SilentlyContinue
if ($buf) {
  $cwd    = (Resolve-Path ".").Path
  $dirAbs = (Resolve-Path $pkgDir).Path
  $relDir = if ($dirAbs.StartsWith($cwd, [StringComparison]::OrdinalIgnoreCase)) {
    $dirAbs.Substring($cwd.Length).TrimStart('\')
  } else { $pkgDir }
  $relDir = ($relDir -replace '\\','/')

  try {
    & buf format -w --path "$relDir" . 2>$null | Out-Null
  } catch {
    Write-Warning ("buf format skipped: {0}" -f $_.Exception.Message)
  }
}
