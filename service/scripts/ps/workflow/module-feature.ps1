# scripts/ps/workflow/feature.ps1
[CmdletBinding()]
param(
  [Parameter(Mandatory = $true)] [string]$Name
)

Set-StrictMode -Version Latest
$ErrorActionPreference = 'Stop'

# --- consts ---
$ApiRoot      = "./api"
$FeatureRoot  = "./internal/feature"

# --- helpers ---
function ConvertTo-PascalCase { param([string]$InputString)
  $parts = ($InputString -replace '[^A-Za-z0-9]+',' ') -split '\s+' | Where-Object { $_ }
  ($parts | ForEach-Object { $_.Substring(0,1).ToUpper() + $_.Substring(1).ToLower() }) -join ''
}
function ConvertTo-LowerCase { param([string]$InputString) $InputString.ToLower() }
function ConvertTo-Plural    { param([string]$InputNoun) if ($InputNoun.ToLower().EndsWith('s')) { $InputNoun } else { "$InputNoun" + "s" } }
function ConvertTo-ImportAlias { param([string]$s) ($s.ToLower() -replace '[^a-z0-9]','_') }

$utf8NoBom = New-Object System.Text.UTF8Encoding($false)

# --- names ---
$base         = ConvertTo-LowerCase $Name                # e.g. "prueba"
$pkgBase      = ($base -replace '[^a-z0-9]','_')         # safe for go package names
$alias        = ConvertTo-ImportAlias $base               # api import alias: api_<alias>
$pascal       = ConvertTo-PascalCase $Name               # "Prueba"
$pluralBase   = ConvertTo-LowerCase (ConvertTo-Plural $base)   # "pruebas"
$pluralPascal = ConvertTo-PascalCase $pluralBase               # "Pruebas"

# --- detect latest api version ---
$apiBaseDir = Join-Path $ApiRoot $base
$apiVersion = 1
if (Test-Path -LiteralPath $apiBaseDir) {
  $max = 0
  Get-ChildItem -LiteralPath $apiBaseDir -Directory -ErrorAction SilentlyContinue | ForEach-Object {
    if ($_.Name -match '^v(\d+)$') {
      $n = [int]$Matches[1]
      if ($n -gt $max) { $max = $n }
    }
  }
  if ($max -gt 0) { $apiVersion = $max } else { $apiVersion = 1 }
} else {
  Write-Warning "API dir not found: $apiBaseDir ; using v1. Run your proto generator first."
}
$apiImport = "service/api/$base/v$apiVersion"

# --- paths ---
$rootDir   = Join-Path $FeatureRoot $base
$bizDir    = Join-Path $rootDir "biz"
$repoDir   = Join-Path $rootDir "repo"
$svcDir    = Join-Path $rootDir "service"

# create dirs
$null = New-Item -ItemType Directory -Force -Path $bizDir, $repoDir, $svcDir

# =========================
# registrars.go
# =========================
$registrarsContent = @"
package $pkgBase

import (
	api_$alias "$apiImport"
	${pkgBase}_service "service/internal/feature/$base/service"
	server_grpc "service/internal/server/grpc"
	server_http "service/internal/server/http"

	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
)

var _ api_$alias.${pascal}ServiceHTTPServer = (*${pkgBase}_service.${pascal}Service)(nil)
var _ api_$alias.${pascal}ServiceServer     = (*${pkgBase}_service.${pascal}Service)(nil)

func New${pascal}HTTPRegistrer(s api_$alias.${pascal}ServiceHTTPServer) server_http.HTTPRegister {
	return func(srv *http.Server) {
		api_$alias.Register${pascal}ServiceHTTPServer(srv, s)
	}
}

func New${pascal}GRPCRegistrer(s api_$alias.${pascal}ServiceServer) server_grpc.GRPCRegister {
	return func(srv *grpc.Server) {
		api_$alias.Register${pascal}ServiceServer(srv, s)
	}
}
"@
[IO.File]::WriteAllText((Join-Path $rootDir "registrars.go"), $registrarsContent, $utf8NoBom)

# =========================
# wire.go
# =========================
$wireContent = @"
package $pkgBase

import (
	api_$alias "$apiImport"
	${pkgBase}_biz "service/internal/feature/$base/biz"
	${pkgBase}_repo "service/internal/feature/$base/repo"
	${pkgBase}_service "service/internal/feature/$base/service"

	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	${pkgBase}_repo.New${pascal}Repo,
	${pkgBase}_biz.New${pascal}Usecase,
	${pkgBase}_service.New${pascal}Service,

	wire.Bind(new(api_$alias.${pascal}ServiceHTTPServer), new(*${pkgBase}_service.${pascal}Service)),
	wire.Bind(new(api_$alias.${pascal}ServiceServer),     new(*${pkgBase}_service.${pascal}Service)),

	New${pascal}HTTPRegistrer,
	New${pascal}GRPCRegistrer,
)
"@
[IO.File]::WriteAllText((Join-Path $rootDir "wire.go"), $wireContent, $utf8NoBom)

# =========================
# biz/biz.go
# =========================
$bizGo = @"
package ${pkgBase}_biz

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
)

// ${pascal}Repo describes repository interface
type ${pascal}Repo interface {
	Find${pluralPascal}(ctx context.Context, id *uint, name *string) ([]${pascal}, error)
	Upsert${pascal}(ctx context.Context, in *${pascal}) (*${pascal}, error)
	Delete${pascal}ById(ctx context.Context, id uint) error
}

// ${pascal}Usecase business logic
type ${pascal}Usecase struct {
	repo ${pascal}Repo
	log  *log.Helper
}

func New${pascal}Usecase(repo ${pascal}Repo, logger log.Logger) *${pascal}Usecase {
	return &${pascal}Usecase{
		repo: repo,
		log:  log.NewHelper(logger),
	}
}
"@
[IO.File]::WriteAllText((Join-Path $bizDir "biz.go"), $bizGo, $utf8NoBom)

# =========================
# biz/entity.go
# =========================
$entityGo = @"
package ${pkgBase}_biz

import "time"

// ${pascal} business entity
type ${pascal} struct {
	ID        uint
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}
"@
[IO.File]::WriteAllText((Join-Path $bizDir "entity.go"), $entityGo, $utf8NoBom)

# =========================
# biz/usecase.go
# =========================
$usecaseGo = @"
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
[IO.File]::WriteAllText((Join-Path $bizDir "usecase.go"), $usecaseGo, $utf8NoBom)

# =========================
# repo/repo.go
# =========================
$repoGo = @"
package ${pkgBase}_repo

import (
	\"service/internal/data\"
	${pkgBase}_biz \"service/internal/feature/$base/biz\"

	\"github.com/go-kratos/kratos/v2/log\"
)

// ${base}Repo concrete repository
type ${base}Repo struct {
	data *data.Data
	log  *log.Helper
}

func New${pascal}Repo(data *data.Data, logger log.Logger) ${pkgBase}_biz.${pascal}Repo {
	return &${base}Repo{
		data: data,
		log:  log.NewHelper(logger),
	}
}
"@
[IO.File]::WriteAllText((Join-Path $repoDir "repo.go"), $repoGo, $utf8NoBom)

# =========================
# repo/r_list.go
# =========================
$rList = @"
package ${pkgBase}_repo

import (
	\"context\"
	${pkgBase}_biz \"service/internal/feature/$base/biz\"
)

// Find${pluralPascal} TODO: implement filters
func (r *${base}Repo) Find${pluralPascal}(ctx context.Context, id *uint, name *string) ([]${pkgBase}_biz.${pascal}, error) {
	return []${pkgBase}_biz.${pascal}{}, nil
}
"@
[IO.File]::WriteAllText((Join-Path $repoDir "r_list.go"), $rList, $utf8NoBom)

# =========================
# repo/r_upsert.go
# =========================
$rUpsert = @"
package ${pkgBase}_repo

import (
	\"context\"
	${pkgBase}_biz \"service/internal/feature/$base/biz\"
)

// Upsert${pascal} TODO: implement upsert
func (r *${base}Repo) Upsert${pascal}(ctx context.Context, in *${pkgBase}_biz.${pascal}) (*${pkgBase}_biz.${pascal}, error) {
	return in, nil
}
"@
[IO.File]::WriteAllText((Join-Path $repoDir "r_upsert.go"), $rUpsert, $utf8NoBom)

# =========================
# repo/r_delete_by_id.go
# =========================
$rDelete = @"
package ${pkgBase}_repo

import \"context\"

// Delete${pascal}ById TODO: implement delete
func (r *${base}Repo) Delete${pascal}ById(ctx context.Context, id uint) error {
	return nil
}
"@
[IO.File]::WriteAllText((Join-Path $repoDir "r_delete_by_id.go"), $rDelete, $utf8NoBom)

# =========================
# service/service.go
# =========================
$svcGo = @"
package ${pkgBase}_service

import (
	api_$alias \"$apiImport\"
	${pkgBase}_biz \"service/internal/feature/$base/biz\"
)

// ${pascal}Service implements api_$alias.${pascal}Service
type ${pascal}Service struct {
	api_$alias.Unimplemented${pascal}ServiceServer
	uc *${pkgBase}_biz.${pascal}Usecase
}

func New${pascal}Service(uc *${pkgBase}_biz.${pascal}Usecase) *${pascal}Service {
	return &${pascal}Service{uc: uc}
}
"@
[IO.File]::WriteAllText((Join-Path $svcDir "service.go"), $svcGo, $utf8NoBom)

# =========================
# service/s_find.go
# =========================
$sFind = @"
package ${pkgBase}_service

import (
	\"context\"

	api_$alias \"$apiImport\"
	${pkgBase}_biz \"service/internal/feature/$base/biz\"
	\"service/pkg/converter\"
	\"service/pkg/generic\"
)

func (s *${pascal}Service) Find${pluralPascal}(ctx context.Context, req *api_$alias.Find${pluralPascal}Request) (*api_$alias.Find${pluralPascal}Response, error) {
	var idPtr *uint
	var namePtr *string
	if req.Id != 0 {
		tmp := uint(req.Id)
		idPtr = &tmp
	}
	if req.Name != \"\" {
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
				Message: \"no items\",
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
			Message: \"OK\",
		},
	}, nil
}
"@
[IO.File]::WriteAllText((Join-Path $svcDir "s_find.go"), $sFind, $utf8NoBom)

# =========================
# service/s_upsert.go
# =========================
$sUpsert = @"
package ${pkgBase}_service

import (
	\"context\"

	api_$alias \"$apiImport\"
	${pkgBase}_biz \"service/internal/feature/$base/biz\"
	\"service/pkg/converter\"
	\"service/pkg/generic\"
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
			Message: \"OK\",
		},
	}, nil
}
"@
[IO.File]::WriteAllText((Join-Path $svcDir "s_upsert.go"), $sUpsert, $utf8NoBom)

# =========================
# service/s_delete_by_id.go
# =========================
$sDelete = @"
package ${pkgBase}_service

import (
	\"context\"

	api_$alias \"$apiImport\"
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
			Message: \"OK\",
		},
	}, nil
}
"@
[IO.File]::WriteAllText((Join-Path $svcDir "s_delete_by_id.go"), $sDelete, $utf8NoBom)

Write-Host ("Feature created: internal/feature/{0} (api v{1})" -f $base, $apiVersion) -ForegroundColor Green
