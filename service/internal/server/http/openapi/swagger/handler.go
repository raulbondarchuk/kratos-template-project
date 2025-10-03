package swagger

import (
	"io/fs"
	stdhttp "net/http"

	openapifs "service/docs"

	kratoshttp "github.com/go-kratos/kratos/v2/transport/http"
)

func AttachEmbeddedSwaggerUI(s *kratoshttp.Server) {
	attachBootstrap(s)

	// redirect /swagger
	s.HandleFunc("/swagger", func(w stdhttp.ResponseWriter, r *stdhttp.Request) {
		stdhttp.Redirect(w, r, "/swagger-ui", stdhttp.StatusSeeOther)
	})

	// logo
	s.HandleFunc("/swagger/logo.png", func(w stdhttp.ResponseWriter, r *stdhttp.Request) {
		data, err := openapifs.FS.ReadFile("logo.png")
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
	s.HandleFunc("/swagger/login", func(w stdhttp.ResponseWriter, r *stdhttp.Request) {
		switch r.Method {
		case stdhttp.MethodGet:
			if isAuthed(r) {
				stdhttp.Redirect(w, r, "/swagger-ui", stdhttp.StatusSeeOther)
				return
			}
			serveLoginPage(w, r, "")

		case stdhttp.MethodPost:
			_ = r.ParseForm()
			username := r.FormValue("username")
			password := r.FormValue("password")

			if token, ok := authenticateWithAPI(username, password); ok {
				setSessionCookie(w, token)                                    // HttpOnly cookie with token
				stdhttp.Redirect(w, r, "/swagger-ui", stdhttp.StatusSeeOther) // without query
				return
			}
			serveLoginPage(w, r, "Error de autorización o autenticación")

		default:
			w.WriteHeader(stdhttp.StatusMethodNotAllowed)
		}
	})

	// logout
	s.HandleFunc("/swagger/logout", func(w stdhttp.ResponseWriter, r *stdhttp.Request) {
		if r.Method != stdhttp.MethodPost && r.Method != stdhttp.MethodGet {
			w.WriteHeader(stdhttp.StatusMethodNotAllowed)
			return
		}
		clearSessionCookie(w)
		stdhttp.Redirect(w, r, "/swagger/login", stdhttp.StatusSeeOther)
	})

	// openapi.yaml
	s.HandleFunc("/swagger/openapi.yaml", authRequired(func(w stdhttp.ResponseWriter, r *stdhttp.Request) {
		data, err := openapifs.FS.ReadFile("openapi.yaml")
		if err != nil {
			httpNotFound(w)
			return
		}
		setOnlyContentType(w, "application/yaml")
		setNoCache(w)
		w.WriteHeader(stdhttp.StatusOK)
		_, _ = w.Write(data)
	}))

	// static /swagger/openapi/*
	if sub, err := fs.Sub(openapifs.FS, "openapi"); err == nil {
		protected := authRequired(func(w stdhttp.ResponseWriter, r *stdhttp.Request) {
			stdhttp.StripPrefix("/swagger/openapi/",
				stdhttp.FileServer(stdhttp.FS(sub)),
			).ServeHTTP(w, r)
		})
		s.Handle("/swagger/openapi/", protected)
	}

	// UI
	s.HandleFunc("/swagger-ui", authRequired(func(w stdhttp.ResponseWriter, r *stdhttp.Request) {
		setOnlyContentType(w, "text/html; charset=utf-8")
		setNoCache(w)
		w.WriteHeader(stdhttp.StatusOK)
		_, _ = w.Write([]byte(uiHTML))
	}))

}
