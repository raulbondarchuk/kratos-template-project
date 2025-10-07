package swagger

import (
	"io/fs"
	stdhttp "net/http"

	openapifs "service/docs"

	kratoshttp "github.com/go-kratos/kratos/v2/transport/http"
)

func AttachEmbeddedSwaggerUI(s *kratoshttp.Server) {
	AttachEmbeddedSwaggerUIWithConfig(s, Config{
		Base:   "",
		DocsFS: openapifs.FS,
	})
}

// New: multiple instances with different prefixes
func AttachEmbeddedSwaggerUIWithConfig(s *kratoshttp.Server, cfg Config) {
	cfg.normalize()
	attachBootstrap(s, &cfg)

	// dynamic base for redirects/templates
	p := func(r *stdhttp.Request, path string) string { return baseForReq(r, &cfg) + path }

	// redirect <base>/docs → <base>/docs/ui (double registration)
	reg(s, cfg.Base, "/docs", func(w stdhttp.ResponseWriter, r *stdhttp.Request) {
		stdhttp.Redirect(w, r, p(r, "/docs/ui"), stdhttp.StatusSeeOther)
	})

	// legacy: <base>/swagger → <base>/docs/ui
	reg(s, cfg.Base, "/swagger", func(w stdhttp.ResponseWriter, r *stdhttp.Request) {
		stdhttp.Redirect(w, r, p(r, "/docs/ui"), stdhttp.StatusSeeOther)
	})

	// logo
	reg(s, cfg.Base, "/docs/logo.png", func(w stdhttp.ResponseWriter, r *stdhttp.Request) {
		data, err := fs.ReadFile(cfg.DocsFS, "logo.png")
		if err != nil {
			httpNotFound(w)
			return
		}
		setOnlyContentType(w, "image/png")
		setNoCache(w)
		w.WriteHeader(stdhttp.StatusOK)
		_, _ = w.Write(data)
	})

	// login
	reg(s, cfg.Base, "/docs/login", func(w stdhttp.ResponseWriter, r *stdhttp.Request) {
		switch r.Method {
		case stdhttp.MethodGet:
			if isAuthed(r, &cfg) {
				stdhttp.Redirect(w, r, p(r, "/docs/ui"), stdhttp.StatusSeeOther)
				return
			}
			serveLoginPage(w, r, "", &cfg)

		case stdhttp.MethodPost:
			_ = r.ParseForm()
			username := r.FormValue("username")
			password := r.FormValue("password")

			if isInternalDocsUser(username) {
				if verifyInternalDocsPassword(password) {
					setSessionCookieForReq(w, r, &cfg, username)
					stdhttp.Redirect(w, r, p(r, "/docs/ui"), stdhttp.StatusSeeOther)
					return
				}
				serveLoginPage(w, r, "Error de autorización o autenticación (internal)", &cfg)
				return
			}

			if token, ok := authenticateWithAPI(username, password, &cfg); ok {
				setSessionCookieForReq(w, r, &cfg, token)
				stdhttp.Redirect(w, r, p(r, "/docs/ui"), stdhttp.StatusSeeOther)
				return
			}
			serveLoginPage(w, r, "Error de autorización o autenticación", &cfg)

		default:
			w.WriteHeader(stdhttp.StatusMethodNotAllowed)
		}
	})

	// logout
	reg(s, cfg.Base, "/docs/logout", func(w stdhttp.ResponseWriter, r *stdhttp.Request) {
		if r.Method != stdhttp.MethodPost && r.Method != stdhttp.MethodGet {
			w.WriteHeader(stdhttp.StatusMethodNotAllowed)
			return
		}
		clearSessionCookieForReq(w, r, &cfg)
		stdhttp.Redirect(w, r, p(r, "/docs/login"), stdhttp.StatusSeeOther)
	})

	// openapi.yaml
	reg(s, cfg.Base, "/docs/openapi.yaml", authRequired(&cfg, func(w stdhttp.ResponseWriter, r *stdhttp.Request) {
		data, err := fs.ReadFile(cfg.DocsFS, "openapi.yaml")
		if err != nil {
			httpNotFound(w)
			return
		}
		setOnlyContentType(w, "application/yaml")
		setNoCache(w)
		w.WriteHeader(stdhttp.StatusOK)
		_, _ = w.Write(data)
	}))

	// static <base>/docs/openapi/*
	if sub, err := fs.Sub(cfg.DocsFS, "openapi"); err == nil {
		protected := authRequired(&cfg, func(w stdhttp.ResponseWriter, r *stdhttp.Request) {
			stdhttp.StripPrefix(baseForReq(r, &cfg)+"/docs/openapi/",
				stdhttp.FileServer(stdhttp.FS(sub)),
			).ServeHTTP(w, r)
		})
		regFS(s, cfg.Base, "/docs/openapi/", protected)
	}

	// handler.go (fragnemt UI)
	reg(s, cfg.Base, "/docs/ui", authRequired(&cfg, func(w stdhttp.ResponseWriter, r *stdhttp.Request) {
		setOnlyContentType(w, "text/html; charset=utf-8")
		setNoCache(w)
		w.WriteHeader(stdhttp.StatusOK)

		_ = uiTpl.Execute(w, struct {
			Base        string
			ServiceName string
		}{
			Base:        baseForReq(r, &cfg),
			ServiceName: cfg.ServiceName,
		})
	}))
}
