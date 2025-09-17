# scripts/ps/workflow/service.ps1
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
function ConvertTo-ImportAlias { param([string]$s) ($s.ToLower() -replace '[^a-z0-9]','_') }

$base         = ConvertTo-LowerCase $Name
$pkgBase      = ($base -replace '[^a-z0-9]','_')
$pascal       = ConvertTo-PascalCase $Name
$pluralPascal = ConvertTo-PascalCase (ConvertTo-Plural $base)
$alias        = ConvertTo-ImportAlias $base

# latest api vN
$apiBaseDir = Join-Path $ApiRoot $base
$apiVersion = 1
if (Test-Path $apiBaseDir) {
  $max=0; Get-ChildItem $apiBaseDir -Directory | ForEach-Object { if ($_.Name -match '^v(\d+)$'){ $n=[int]$Matches[1]; if($n -gt $max){$max=$n} } }
  if ($max -gt 0){ $apiVersion=$max }
}
$featureRootV = Join-Path (Join-Path $FeatureRoot $base) "v$apiVersion"
$svcDir = Join-Path $featureRootV "service"
$null = New-Item -ItemType Directory -Force -Path $svcDir

$apiImport = "service/api/$base/v$apiVersion"
$bizImport = "service/internal/feature/$base/v$apiVersion/biz"

# service.go
$p = Join-Path $svcDir "service.go"
if (-not (Test-Path $p)) {
$txt = @"
package ${pkgBase}_service

import (
	api_$alias "$apiImport"
	${pkgBase}_biz "$bizImport"
)

type ${pascal}Service struct {
	api_$alias.Unimplemented${pascal}ServiceServer
	uc *${pkgBase}_biz.${pascal}Usecase
}

func New${pascal}Service(uc *${pkgBase}_biz.${pascal}Usecase) *${pascal}Service {
	return &${pascal}Service{uc: uc}
}
"@
[IO.File]::WriteAllText($p, $txt, $utf8NoBom)
}

# s_find.go
$p = Join-Path $svcDir "s_find.go"
if (-not (Test-Path $p)) {
$txt = @"
package ${pkgBase}_service

import (
	"context"

	api_$alias "$apiImport"
	${pkgBase}_biz "$bizImport"
	"service/pkg/converter"
	"service/pkg/generic"
)

func (s *${pascal}Service) Find${pluralPascal}(ctx context.Context, req *api_$alias.Find${pluralPascal}Request) (*api_$alias.Find${pluralPascal}Response, error) {
	var idPtr *uint
	var namePtr *string
	if req.Id != 0 {
		tmp := uint(req.Id)
		idPtr = &tmp
	}
	if req.Name != "" {
		tmp := req.Name
		namePtr = &tmp
	}

	bizRes, err := s.uc.Find${pluralPascal}(ctx, idPtr, namePtr)
	if err != nil {
		return &api_$alias.Find${pluralPascal}Response{
			Meta: &api_$alias.MetaResponse{
				Code:    api_$alias.ResponseCode_RESPONSE_CODE_INTERNAL_SERVER_ERROR,
				Message: err.Error(),
			},
		}, nil
	}
	if len(bizRes) == 0 {
		return &api_$alias.Find${pluralPascal}Response{
			Meta: &api_$alias.MetaResponse{
				Code:    api_$alias.ResponseCode_RESPONSE_CODE_OK,
				Message: "no items",
			},
		}, nil
	}

	dto, err := generic.ToDTOSliceGeneric[${pkgBase}_biz.${pascal}, api_$alias.${pascal}](bizRes)
	if err != nil {
		return &api_$alias.Find${pluralPascal}Response{
			Meta: &api_$alias.MetaResponse{
				Code:    api_$alias.ResponseCode_RESPONSE_CODE_INTERNAL_SERVER_ERROR,
				Message: err.Error(),
			},
		}, nil
	}

	for i := range bizRes {
		dto[i].CreatedAt = converter.ConvertToGoogleTimestamp(bizRes[i].CreatedAt)
		dto[i].UpdatedAt = converter.ConvertToGoogleTimestamp(bizRes[i].UpdatedAt)
	}

	return &api_$alias.Find${pluralPascal}Response{
		Items: generic.ToPointerSliceGeneric(dto),
		Meta: &api_$alias.MetaResponse{
			Code:    api_$alias.ResponseCode_RESPONSE_CODE_OK,
			Message: "OK",
		},
	}, nil
}
"@
[IO.File]::WriteAllText($p, $txt, $utf8NoBom)
}

# s_upsert.go
$p = Join-Path $svcDir "s_upsert.go"
if (-not (Test-Path $p)) {
$txt = @"
package ${pkgBase}_service

import (
	"context"

	api_$alias "$apiImport"
	${pkgBase}_biz "$bizImport"
	"service/pkg/converter"
	"service/pkg/generic"
)

func (s *${pascal}Service) Upsert${pascal}(ctx context.Context, req *api_$alias.Upsert${pascal}Request) (*api_$alias.Upsert${pascal}Response, error) {
	in := &${pkgBase}_biz.${pascal}{
		ID:   uint(req.Id),
		Name: req.Name,
	}
	res, err := s.uc.Upsert${pascal}(ctx, in)
	if err != nil {
		return &api_$alias.Upsert${pascal}Response{
			Meta: &api_$alias.MetaResponse{
				Code:    api_$alias.ResponseCode_RESPONSE_CODE_INTERNAL_SERVER_ERROR,
				Message: err.Error(),
			},
		}, nil
	}

	dto, err := generic.ToDTOGeneric[${pkgBase}_biz.${pascal}, api_$alias.${pascal}](*res)
	if err != nil {
		return &api_$alias.Upsert${pascal}Response{
			Meta: &api_$alias.MetaResponse{
				Code:    api_$alias.ResponseCode_RESPONSE_CODE_INTERNAL_SERVER_ERROR,
				Message: err.Error(),
			},
		}, nil
	}
	dto.CreatedAt = converter.ConvertToGoogleTimestamp(res.CreatedAt)
	dto.UpdatedAt = converter.ConvertToGoogleTimestamp(res.UpdatedAt)

	return &api_$alias.Upsert${pascal}Response{
		Item: &dto,
		Meta: &api_$alias.MetaResponse{
			Code:    api_$alias.ResponseCode_RESPONSE_CODE_OK,
			Message: "OK",
		},
	}, nil
}
"@
[IO.File]::WriteAllText($p, $txt, $utf8NoBom)
}

# s_delete_by_id.go
$p = Join-Path $svcDir "s_delete_by_id.go"
if (-not (Test-Path $p)) {
$txt = @"
package ${pkgBase}_service

import (
	"context"

	api_$alias "$apiImport"
)

func (s *${pascal}Service) Delete${pascal}ById(ctx context.Context, req *api_$alias.Delete${pascal}ByIdRequest) (*api_$alias.Delete${pascal}ByIdResponse, error) {
	if err := s.uc.Delete${pascal}ById(ctx, uint(req.Id)); err != nil {
		return &api_$alias.Delete${pascal}ByIdResponse{
			Meta: &api_$alias.MetaResponse{
				Code:    api_$alias.ResponseCode_RESPONSE_CODE_INTERNAL_SERVER_ERROR,
				Message: err.Error(),
			},
		}, nil
	}
	return &api_$alias.Delete${pascal}ByIdResponse{
		Meta: &api_$alias.MetaResponse{
			Code:    api_$alias.ResponseCode_RESPONSE_CODE_OK,
			Message: "OK",
		},
	}, nil
}
"@
[IO.File]::WriteAllText($p, $txt, $utf8NoBom)
}

Write-Host ("service: {0}/v{1}" -f $base, $apiVersion) -ForegroundColor Green
