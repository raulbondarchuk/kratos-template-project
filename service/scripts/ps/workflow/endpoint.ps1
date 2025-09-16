# scripts/ps/workflow/endpoint.ps1
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
  // TODO: add your codes here
}

// Your meta-object: only code and message.
message MetaResponse {
  ResponseCode code = 1; // RESPONSE_CODE_OK or other codes
  string message = 2;    // e.g. "ok", "ip blocked", "db error", "user not found"
}
"@

# --- <name>.proto content (CRUD: upsert, find(+filters), delete) ---
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

// ======================================================
// ${pascal}Service
// ======================================================

/** API para gestion de $pluralBase. */
service ${pascal}Service {
  /** Lista o busca $pluralBase.
   *
   * Descripcion
   * Devuelve lista de $pluralBase. Si hay filtros, se aplican.
   *
   * Parametros (query)
   * - 'id' (opcional)
   * - 'name' (opcional)
   *
   * Respuestas
   * - 200 OK: 'items' + 'meta.code = RESPONSE_CODE_OK'
   * - 200 con error logico: 'meta.code != RESPONSE_CODE_OK' y 'meta.message'
   *
   * Ejemplos
   * GET ${route}
   * GET ${route}?id=1
   */
  rpc Find$pluralPascal(Find${pluralPascal}Request) returns (Find${pluralPascal}Response) {
    option (google.api.http) = {get: "${route}"};
  }

  /** Crea o actualiza (upsert) $base.
   *
   * Descripcion
   * Si 'id' es 0 o no se envia -> crea; si 'id' > 0 -> actualiza.
   *
   * Parametros (body JSON)
   * - 'id' (opcional, 0=create)
   * - 'name' (obligatorio)
   *
   * Respuestas
   * - 200 OK: 'item' resultante + 'meta.code = RESPONSE_CODE_OK'
   * - 200 con error logico: validacion o duplicidad
   *
   * Ejemplos
   * POST ${route}
   * Body: { "name":"$pascal 1" }
   */
  rpc Upsert$pascal(Upsert${pascal}Request) returns (Upsert${pascal}Response) {
    option (google.api.http) = {
      post: "${route}"
      body: "*"
    };
  }

  /** Elimina $base por ID.
   *
   * Parametros (query)
   * - 'id' (obligatorio): ID interno
   *
   * Respuestas
   * - 200 OK: 'meta.code = RESPONSE_CODE_OK'
   * - 200 con error logico: 'meta.code != RESPONSE_CODE_OK'
   *
   * Ejemplos
   * DELETE ${route}?id=123
   */
  rpc Delete${pascal}ById(Delete${pascal}ByIdRequest) returns (Delete${pascal}ByIdResponse) {
    option (google.api.http) = {delete: "${route}"};
  }
}

// ======================================================
// Find$pluralPascal (GET ${route})
// ======================================================

/** Filtros de busqueda. Si vacio -> listado completo. */
message Find${pluralPascal}Request {
  uint32 id = 1 [(google.api.field_behavior) = OPTIONAL];   // filtrar por ID
  string name = 2 [(google.api.field_behavior) = OPTIONAL]; // filtrar por nombre
}

/** Respuesta con lista. */
message Find${pluralPascal}Response {
  repeated $pascal items = 1; // coleccion de resultados
  MetaResponse meta = 2;      // estado de la operacion
}

// ======================================================
// Upsert$pascal (POST ${route})
// ======================================================

/** Cuerpo para crear o actualizar. */
message Upsert${pascal}Request {
  uint32 id = 1 [(google.api.field_behavior) = OPTIONAL];  // 0 -> crear
  string name = 2 [(google.api.field_behavior) = REQUIRED];
}

/** Resultado del upsert. */
message Upsert${pascal}Response {
  $pascal item = 1;  // entidad creada o actualizada
  MetaResponse meta = 2;
}

// ======================================================
// Delete${pascal}ById (DELETE ${route}?id=123)
// ======================================================

/** Peticion para eliminar por ID. */
message Delete${pascal}ByIdRequest {
  uint32 id = 1 [(google.api.field_behavior) = REQUIRED]; // ID interno
}

/** Respuesta de eliminacion. */
message Delete${pascal}ByIdResponse {
  MetaResponse meta = 1;
}

// ======================================================
// Common models
// ======================================================

/** Modelo $base. Ajusta campos al dominio. */
message $pascal {
  uint32 id = 1;
  string name = 2; // ejemplo: "$pascal 1"
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
  $cwd   = (Resolve-Path ".").Path
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
