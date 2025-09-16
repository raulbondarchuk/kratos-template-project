// internal/server/http/registrar.go
package server_http

import "github.com/go-kratos/kratos/v2/transport/http"

// HTTPRegistrar is a function that registers routes on the server.
type HTTPRegister func(*http.Server)
