package server_http

import (
	api_template "service/api/template"

	"github.com/go-kratos/kratos/v2/transport/http"
)

// HTTPRegister is a function that registers routes on the server.
type HTTPRegister func(*http.Server)

func LoadRoutes(srv *http.Server,
	template api_template.TemplatesHTTPServer,
) {
	api_template.RegisterTemplatesHTTPServer(srv, template)
}
