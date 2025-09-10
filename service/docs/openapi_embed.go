// docs/openapi_embed.go
package openapifs

import "embed"

//go:embed openapi.yaml openapi/*
var FS embed.FS
