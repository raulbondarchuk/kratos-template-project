# scripts/ps/workflow/biz.ps1
[CmdletBinding()]
param([Parameter(Mandatory = $true)] [string]$Name)

Set-StrictMode -Version Latest
$ErrorActionPreference = 'Stop'

$ApiRoot     = "./api"
$FeatureRoot = "./internal/feature"
$utf8NoBom   = New-Object System.Text.UTF8Encoding($false)

function ConvertTo-PascalCase { param([string]$s)
  $parts = ($s -replace '[^A-Za-z0-9]+',' ') -split '\s+' | Where-Object { $_ }
  ($parts | ForEach-Object { $_.Substring(0,1).ToUpper() + $_.Substring(1).ToLower() }) -join ''
}
function ConvertTo-LowerCase { param([string]$s) $s.ToLower() }
function ConvertTo-Plural    { param([string]$s) if ($s.ToLower().EndsWith('s')){$s}else{"$s"+"s"} }

$base      = ConvertTo-LowerCase $Name
$pkgBase   = ($base -replace '[^a-z0-9]','_')
$pascal    = ConvertTo-PascalCase $Name
$pluralPascal = ConvertTo-PascalCase (ConvertTo-Plural $base)

# latest api vN
$apiBaseDir = Join-Path $ApiRoot $base
$apiVersion = 1
if (Test-Path $apiBaseDir) {
  $max=0; Get-ChildItem $apiBaseDir -Directory | ForEach-Object { if ($_.Name -match '^v(\d+)$'){ $n=[int]$Matches[1]; if($n -gt $max){$max=$n} } }
  if ($max -gt 0){ $apiVersion=$max }
}
$featureRootV = Join-Path (Join-Path $FeatureRoot $base) "v$apiVersion"
$bizDir = Join-Path $featureRootV "biz"
$null = New-Item -ItemType Directory -Force -Path $bizDir

# biz.go
$p = Join-Path $bizDir "biz.go"
if (-not (Test-Path $p)) {
$txt = @"
package ${pkgBase}_biz

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
)

type ${pascal}Repo interface {
	Find${pluralPascal}(ctx context.Context, id *uint, name *string) ([]${pascal}, error)
	Upsert${pascal}(ctx context.Context, in *${pascal}) (*${pascal}, error)
	Delete${pascal}ById(ctx context.Context, id uint) error
}

type ${pascal}Usecase struct {
	repo ${pascal}Repo
	log  *log.Helper
}

func New${pascal}Usecase(repo ${pascal}Repo, logger log.Logger) *${pascal}Usecase {
	return &${pascal}Usecase{repo: repo, log: log.NewHelper(logger)}
}
"@
[IO.File]::WriteAllText($p, $txt, $utf8NoBom)
}

# entity.go
$p = Join-Path $bizDir "entity.go"
if (-not (Test-Path $p)) {
$txt = @"
package ${pkgBase}_biz

import "time"

type ${pascal} struct {
	ID        uint
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}
"@
[IO.File]::WriteAllText($p, $txt, $utf8NoBom)
}

# usecase.go
$p = Join-Path $bizDir "usecase.go"
if (-not (Test-Path $p)) {
$txt = @"
package ${pkgBase}_biz

import "context"

func (uc *${pascal}Usecase) Find${pluralPascal}(ctx context.Context, id *uint, name *string) ([]${pascal}, error) {
	return uc.repo.Find${pluralPascal}(ctx, id, name)
}

func (uc *${pascal}Usecase) Upsert${pascal}(ctx context.Context, in *${pascal}) (*${pascal}, error) {
	return uc.repo.Upsert${pascal}(ctx, in)
}

func (uc *${pascal}Usecase) Delete${pascal}ById(ctx context.Context, id uint) error {
	return uc.repo.Delete${pascal}ById(ctx, id)
}
"@
[IO.File]::WriteAllText($p, $txt, $utf8NoBom)
}

Write-Host ("biz: {0}/v{1}" -f $base, $apiVersion) -ForegroundColor Green