package server_http

import (
	api_template "service/bin/proto/endpoints/template"

	"github.com/go-kratos/kratos/v2/transport/http"
)

func LoadRoutes(srv *http.Server,
	template api_template.TemplatesHTTPServer,
) {
	api_template.RegisterTemplatesHTTPServer(srv, template)
}
