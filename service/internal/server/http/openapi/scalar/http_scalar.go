package scalar

import (
	stdhttp "net/http"

	openapifs "service/docs" // embed with openapi.yaml

	kratoshttp "github.com/go-kratos/kratos/v2/transport/http"
)

// Raises:
//   - /scalar/openapi.yaml — specification OpenAPI (yaml)
//   - /docs               — HTML with Scalar API Reference
func AttachScalarDocs(s *kratoshttp.Server) {
	// Specification under a separate path (does not conflict with Swagger UI)
	s.HandleFunc("/scalar/openapi.yaml", func(w stdhttp.ResponseWriter, r *stdhttp.Request) {
		data, err := openapifs.FS.ReadFile("openapi.yaml")
		if err != nil {
			httpNotFound(w)
			return
		}
		w.Header().Set("Content-Type", "application/yaml; charset=utf-8")
		w.WriteHeader(stdhttp.StatusOK)
		_, _ = w.Write(data)
	})

	// (optional) serve /scalar/openapi/* if you have nested resources
	s.Handle("/scalar/openapi/",
		stdhttp.StripPrefix("/scalar/openapi/",
			stdhttp.FileServer(stdhttp.FS(openapifs.FS))))

	// Scalar page
	s.HandleFunc("/docs", func(w stdhttp.ResponseWriter, r *stdhttp.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write([]byte(`<!doctype html>
<html lang="es">
<head>
  <meta charset="utf-8"/>
  <meta name="viewport" content="width=device-width,initial-scale=1"/>
  <title>API Reference</title>
  <style>
    html,body,#app{height:100%;margin:0}
    body{font-family:system-ui,-apple-system,Segoe UI,Roboto,Ubuntu,"Helvetica Neue",Arial}
  </style>
</head>
<body>
  <div id="app"></div>

  <!-- Глобальный объект Scalar доступен после этого скрипта -->
  <script src="https://cdn.jsdelivr.net/npm/@scalar/api-reference@latest"></script>
  <script>
    // Можно переопределить ?spec=/swagger/openapi.yaml
    const qs = new URLSearchParams(location.search);
    const specUrl = qs.get('spec') || '/scalar/openapi.yaml';

    // Верный вызов для CDN-версии:
    Scalar.createApiReference('#app', {
      url: specUrl,
      // proxyUrl можно не указывать на том же домене, но оставить строку если вынесешь спеки на другой origin:
      // proxyUrl: 'https://proxy.scalar.com',
      // Пара приятных опций (по желанию):
      // theme: 'default',       // 'default' | 'alternate' | 'dark'
      // layout: 'modern',       // 'modern' | 'classic'
      // metaData: { title: 'API Docs' },
    });
  </script>
</body>
</html>`))
	})
}

func httpNotFound(w stdhttp.ResponseWriter) {
	w.WriteHeader(stdhttp.StatusNotFound)
	_, _ = w.Write([]byte("not found"))
}
