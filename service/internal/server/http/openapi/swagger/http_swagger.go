package swagger

import (
	stdhttp "net/http"

	openapifs "service/docs" // embed with openapi.yaml

	kratoshttp "github.com/go-kratos/kratos/v2/transport/http"
)

func AttachEmbeddedSwaggerUI(s *kratoshttp.Server) {
	// Exact handler for /swagger/openapi.yaml
	s.HandleFunc("/swagger/openapi.yaml", func(w stdhttp.ResponseWriter, r *stdhttp.Request) {
		data, err := openapifs.FS.ReadFile("openapi/openapi.yaml")
		if err != nil {
			httpNotFound(w)
			return
		}
		w.Header().Set("Content-Type", "application/yaml; charset=utf-8")
		w.WriteHeader(stdhttp.StatusOK)
		_, _ = w.Write(data)
	})

	// All files from docs/openapi/* are available at /swagger/openapi/*
	s.Handle("/swagger/openapi/",
		stdhttp.StripPrefix("/swagger/openapi/",
			stdhttp.FileServer(stdhttp.FS(openapifs.FS))))

	// Simple UI page
	s.HandleFunc("/swagger-ui", func(w stdhttp.ResponseWriter, r *stdhttp.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write([]byte(`<!doctype html>
<html>
<head>
  <meta charset="utf-8"/>
  <title>Swagger UI</title>
  <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist/swagger-ui.css"/>
</head>
<body>
  <div id="swagger-ui"></div>
  <script src="https://unpkg.com/swagger-ui-dist/swagger-ui-bundle.js"></script>
  <script>
    window.ui = SwaggerUIBundle({
      url: '/swagger/openapi.yaml',
      dom_id: '#swagger-ui',
      presets: [SwaggerUIBundle.presets.apis],
    });
  </script>
</body>
</html>`))
	})
}

func httpNotFound(w stdhttp.ResponseWriter) {
	httpStatus(w, stdhttp.StatusNotFound, "not found")
}

func httpStatus(w stdhttp.ResponseWriter, code int, msg string) {
	w.WriteHeader(code)
	_, _ = w.Write([]byte(msg))
}
