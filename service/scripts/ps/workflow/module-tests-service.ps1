# scripts/ps/workflow/module-tests.ps1
[CmdletBinding()]
param(
  [Parameter(Mandatory=$true)] [string]$Name,
  [string]$Version = "",
  [switch]$Force
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

# --- helpers ---
$utf8NoBom = New-Object System.Text.UTF8Encoding($false)

function ConvertTo-PascalCase { param([string]$s)
  $parts = ($s -replace '[^A-Za-z0-9]+',' ') -split '\s+' | Where-Object { $_ }
  ($parts | ForEach-Object { $_.Substring(0,1).ToUpper() + $_.Substring(1).ToLower() }) -join ''
}
function ConvertTo-LowerCase { param([string]$s) $s.ToLower() }
function ConvertTo-PluralForm { param([string]$s) if ($s.ToLower().EndsWith('s')) { $s } else { "$s" + "s" } }
function ConvertTo-AliasName { param([string]$s) ($s.ToLower() -replace '[^a-z0-9]','_') }

# --- normalize names ---
$base         = ConvertTo-LowerCase $Name
$pascal       = ConvertTo-PascalCase $Name
$pluralBase   = ConvertTo-LowerCase (ConvertTo-PluralForm $base)
$pluralPascal = ConvertTo-PascalCase $pluralBase
$alias        = ConvertTo-AliasName $base
$pkgBase      = ($base -replace '[^a-z0-9]','_')

$ApiRoot     = "./api/$base"
$FeatureRoot = "./internal/feature/$base"

Show-Step "Generating unit tests from .proto + biz interface"
Show-Info "Module: base='$base', pascal='$pascal'"

# --- resolve version ---
if (-not (Test-Path $ApiRoot)) { Show-ErrorAndExit "API dir not found: $ApiRoot" }

if ([string]::IsNullOrWhiteSpace($Version)) {
  $max = 0
  Get-ChildItem $ApiRoot -Directory -ErrorAction SilentlyContinue | ForEach-Object {
    if ($_.Name -match '^v(\d+)$') { $n = [int]$Matches[1]; if ($n -gt $max) { $max = $n } }
  }
  if ($max -eq 0) { Show-ErrorAndExit "No API versions in $ApiRoot" }
  $Version = "v$max"
} else {
  if ($Version -match '^v?(\d+)$') { $Version = "v$($Matches[1])" } else { Show-ErrorAndExit "Invalid version format: '$Version' (use 'vN' or 'N')" }
}
$Route = "/$Version/$base"
Show-Info "Using version: $Version"

$ProtoPath = Join-Path $ApiRoot "$Version/$base.proto"
if (-not (Test-Path $ProtoPath)) { Show-ErrorAndExit "Proto file not found: $ProtoPath" }

# --- imports ---
$imports = @(
  '"context"',
  '"os"',
  '"testing"',
  '"github.com/go-kratos/kratos/v2/log"',
  ('api_{0} "service/api/{1}/{2}"' -f $alias, $base, $Version),
  ('{0}_biz "service/internal/feature/{1}/{2}/biz"' -f $pkgBase, $base, $Version)
)
$importBlock = "import (" + [Environment]::NewLine + "`t" + ($imports -join ([Environment]::NewLine + "`t")) + [Environment]::NewLine + ")"

# --- parse proto RPCs ---
$proto     = Get-Content -LiteralPath $ProtoPath -Raw
$hasList   = [regex]::IsMatch($proto, "\brpc\s+List$pluralPascal\s*\(")
$hasFind   = [regex]::IsMatch($proto, "\brpc\s+Find$pluralPascal\s*\(")
$hasUpsert = [regex]::IsMatch($proto, "\brpc\s+Upsert$pascal\s*\(")
$hasDelete = [regex]::IsMatch($proto, "\brpc\s+Delete${pascal}ById\s*\(")
$hasMock   = [regex]::IsMatch($proto, "\brpc\s+Mock\s*\(")
Show-Info ("RPCs (from proto): List={0} Find={1} Upsert={2} Delete={3} Mock={4}" -f $hasList,$hasFind,$hasUpsert,$hasDelete,$hasMock)

# --- destination paths ---
$svcDir   = Join-Path $FeatureRoot "$Version/service"
$testFile = Join-Path $svcDir "${base}_service_test.go"
if (-not (Test-Path $svcDir)) { New-Item -ItemType Directory -Force -Path $svcDir | Out-Null }
if ((Test-Path $testFile) -and -not $Force) { Show-Warn "Test file exists, use -Force to overwrite: $testFile"; exit 0 }

# --- fake repo ---
$fakeRepoFields = @()
$fakeRepoFields += "List   func(ctx context.Context) ([]${pkgBase}_biz.${pascal}, error)"
if ($hasFind)   { $fakeRepoFields += "Get    func(ctx context.Context, id *uint, name *string) ([]${pkgBase}_biz.${pascal}, error)" }
if ($hasUpsert) { $fakeRepoFields += "Upsert func(ctx context.Context, in *${pkgBase}_biz.${pascal}) (*${pkgBase}_biz.${pascal}, error)" }
if ($hasDelete) { $fakeRepoFields += "Delete func(ctx context.Context, id uint) error" }

$fakeRepoStruct = "type fakeRepo struct {`n`t" + ($fakeRepoFields -join "`n`t") + "`n}"

$fakeRepoMethods = @()
$fakeRepoMethods += @"
func (f *fakeRepo) List${pluralPascal}(ctx context.Context) ([]${pkgBase}_biz.${pascal}, error) {
	if f.List != nil {
		return f.List(ctx)
	}
	return nil, nil
}
"@

if ($hasFind) {
  $fakeRepoMethods += @"
func (f *fakeRepo) Find${pluralPascal}(ctx context.Context, id *uint, name *string) ([]${pkgBase}_biz.${pascal}, error) {
	if f.Get != nil {
		return f.Get(ctx, id, name)
	}
	return nil, nil
}
"@
}

if ($hasUpsert) {
  $fakeRepoMethods += @"
func (f *fakeRepo) Upsert${pascal}(ctx context.Context, in *${pkgBase}_biz.${pascal}) (*${pkgBase}_biz.${pascal}, error) {
	if f.Upsert != nil {
		return f.Upsert(ctx, in)
	}
	return in, nil
}
"@
}

if ($hasDelete) {
  $fakeRepoMethods += @"
func (f *fakeRepo) Delete${pascal}ById(ctx context.Context, id uint) error {
	if f.Delete != nil {
		return f.Delete(ctx, id)
	}
	return nil
}
"@
}

# --- tests ---
$tests = @()

if ($hasList) {
  $tests += @"
func Test${pascal}Service_List${pluralPascal}_OK(t *testing.T) {
	logger := log.NewStdLogger(os.Stdout)
	repo := &fakeRepo{}
	uc := ${pkgBase}_biz.New${pascal}Usecase(repo, logger)
	svc := New${pascal}Service(uc)

	resp, err := svc.List${pluralPascal}(context.Background(), &api_${alias}.List${pluralPascal}Request{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.GetMeta().GetCode() != api_${alias}.ResponseCode_RESPONSE_CODE_OK {
		t.Fatalf("meta.code = %v, want OK", resp.GetMeta().GetCode())
	}
}
"@
}

if ($hasFind) {
  $tests += @"
func Test${pascal}Service_Find${pluralPascal}_OK(t *testing.T) {
    logger := log.NewStdLogger(os.Stdout)
    repo := &fakeRepo{}
    uc := ${pkgBase}_biz.New${pascal}Usecase(repo, logger)
    svc := New${pascal}Service(uc)

    resp, err := svc.Find${pluralPascal}(context.Background(), &api_${alias}.Find${pluralPascal}Request{})
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if resp.GetMeta().GetCode() != api_${alias}.ResponseCode_RESPONSE_CODE_OK {
        t.Fatalf("meta.code = %v, want OK", resp.GetMeta().GetCode())
    }
}
"@
}

if ($hasUpsert) {
  $tests += @"
func Test${pascal}Service_Upsert${pascal}_OK(t *testing.T) {
    logger := log.NewStdLogger(os.Stdout)
    repo := &fakeRepo{}
    uc := ${pkgBase}_biz.New${pascal}Usecase(repo, logger)
    svc := New${pascal}Service(uc)

    req := &api_${alias}.Upsert${pascal}Request{ Id: 0, Name: "${pascal}X" }
    resp, err := svc.Upsert${pascal}(context.Background(), req)
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if resp.GetMeta().GetCode() != api_${alias}.ResponseCode_RESPONSE_CODE_OK {
        t.Fatalf("meta.code = %v, want OK", resp.GetMeta().GetCode())
    }
}
"@
}

if ($hasDelete) {
  $tests += @"
func Test${pascal}Service_Delete${pascal}ById_OK(t *testing.T) {
    logger := log.NewStdLogger(os.Stdout)
    repo := &fakeRepo{}
    uc := ${pkgBase}_biz.New${pascal}Usecase(repo, logger)
    svc := New${pascal}Service(uc)

    resp, err := svc.Delete${pascal}ById(context.Background(), &api_${alias}.Delete${pascal}ByIdRequest{ Id: 7 })
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if resp.GetMeta().GetCode() != api_${alias}.ResponseCode_RESPONSE_CODE_OK {
        t.Fatalf("meta.code = %v, want OK", resp.GetMeta().GetCode())
    }
}
"@
}

if (-not $hasFind -and -not $hasUpsert -and -not $hasDelete) {
  $tests += @"
func Test${pascal}Service_NoBusinessRPCs(t *testing.T) {
    t.Skip("no business RPCs generated in proto (ops empty). Mock endpoint exists at ${Route}/mock")
}
"@
}

# --- write file ---
$header = @"
// Code generated by module-tests.ps1; DO NOT EDIT.
package ${pkgBase}_service

$importBlock

$fakeRepoStruct

$($fakeRepoMethods -join "`n")

$($tests -join "`n")
"@

[IO.File]::WriteAllText($testFile, $header, $utf8NoBom)
Show-OK "Created test: $testFile"