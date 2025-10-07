// docs/openapi_embed.go
package openapifs

import "embed"

//go:embed openapi.yaml openapi/* *.png
var FS embed.FS
