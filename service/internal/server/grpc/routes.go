package server_grpc

import (
	api_template "service/api/template"

	"github.com/go-kratos/kratos/v2/transport/grpc"
)

func LoadRoutes(srv *grpc.Server,
	template api_template.TemplatesServer,
) {
	api_template.RegisterTemplatesServer(srv, template)
}
