package swagger

import (
	"bytes"
	"encoding/json"
	"fmt"
	stdhttp "net/http"
	"service/internal/server/middleware/auth/auth/paseto"
	"service/pkg/utils"
	"strings"
	"time"

	kratoshttp "github.com/go-kratos/kratos/v2/transport/http"
)

func isAuthed(r *stdhttp.Request, cfg *Config) bool {
	c, err := r.Cookie(cfg.CookieName)
	return err == nil && strings.TrimSpace(c.Value) != ""
}

func setSessionCookieForReq(w stdhttp.ResponseWriter, r *stdhttp.Request, cfg *Config, value string) {
	stdhttp.SetCookie(w, &stdhttp.Cookie{
		Name:     cfg.CookieName,
		Value:    value,
		Path:     cookiePathForReq(r, cfg),
		HttpOnly: true,
		SameSite: stdhttp.SameSiteLaxMode,
		MaxAge:   3600,
		// Secure: true,
	})
}

func clearSessionCookieForReq(w stdhttp.ResponseWriter, r *stdhttp.Request, cfg *Config) {
	stdhttp.SetCookie(w, &stdhttp.Cookie{
		Name:     cfg.CookieName,
		Value:    "",
		Path:     cookiePathForReq(r, cfg),
		HttpOnly: true,
		SameSite: stdhttp.SameSiteLaxMode,
		MaxAge:   -1,
		Expires:  time.Unix(0, 0),
	})
}

func authRequired(cfg *Config, h stdhttp.HandlerFunc) stdhttp.HandlerFunc {
	return func(w stdhttp.ResponseWriter, r *stdhttp.Request) {
		if !isAuthed(r, cfg) {
			serveLoginPage(w, r, "", cfg)
			return
		}
		h(w, r)
	}
}

func authenticateWithAPI(username, password string, cfg *Config) (string, bool) {
	payload := map[string]string{"username": username, "password": password}
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return "", false
	}

	req, err := stdhttp.NewRequest(stdhttp.MethodPost, cfg.LoginURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", false
	}
	req.Header.Set("Content-Type", "application/json")

	client := &stdhttp.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", false
	}
	defer resp.Body.Close()

	if resp.StatusCode != stdhttp.StatusOK {
		return "", false
	}

	raw := strings.TrimSpace(resp.Header.Get("Authorization"))
	token := extractBearerToken(raw)
	if token == "" {
		return "", false
	}

	rolesStr, err := paseto.VerifyAccessTokenRoles(token)
	if err != nil {
		return "", false
	}
	roles := parseRoles(rolesStr)

	if !hasAnyRole(roles, cfg.RequiredRoles) {
		return "", false
	}
	return token, true
}

func extractBearerToken(h string) string {
	h = strings.TrimSpace(h)
	if h == "" {
		return ""
	}
	if strings.HasPrefix(strings.ToLower(h), "bearer ") {
		return strings.TrimSpace(h[7:])
	}
	return h
}

func parseRoles(s string) []string {
	if s == "" {
		return nil
	}
	parts := strings.FieldsFunc(s, func(r rune) bool {
		return r == ',' || r == ';' || r == '|' || r == ' '
	})
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}

func hasAnyRole(userRoles, required []string) bool {
	set := make(map[string]struct{}, len(userRoles))
	for _, r := range userRoles {
		set[r] = struct{}{}
	}
	for _, req := range required {
		if _, ok := set[req]; ok {
			return true
		}
	}
	return false
}

func attachBootstrap(s *kratoshttp.Server, cfg *Config) {
	h := authRequired(cfg, func(w stdhttp.ResponseWriter, r *stdhttp.Request) {
		c, _ := r.Cookie(cfg.CookieName)
		token := ""
		if c != nil {
			token = strings.TrimSpace(c.Value)
		}

		setOnlyContentType(w, "application/javascript; charset=utf-8")
		setNoCache(w)
		w.WriteHeader(stdhttp.StatusOK)

		js := fmt.Sprintf(`(function(){
  try{
    var t=%q;
    var K="swagger_access_token";
    var K_LAST="swagger_last_seed";
    var K_OVR="swagger_user_override";

    var s = localStorage.getItem(K) || "";
    var last = localStorage.getItem(K_LAST) || "";
    var hasOverride = !!localStorage.getItem(K_OVR);

    if (t) {
      if (t !== last) {
        localStorage.setItem(K, t);
        localStorage.setItem(K_LAST, t);
        localStorage.removeItem(K_OVR);
        s = t;
      } else if (!s) {
        localStorage.setItem(K, t);
        s = t;
      }
    }

    var i=document.getElementById('authToken');
    if(i){ i.value = localStorage.getItem(K) || s || ""; }
  }catch(e){}
})();`, token)

		_, _ = w.Write([]byte(js))
	})

	reg(s, cfg.Base, "/docs/bootstrap.js", h)
}

func isInternalDocsUser(username string) bool {
	wantUser := utils.EnvFirst("SW_LOGIN")
	if wantUser == "" {
		wantUser = "docs"
	}
	return username == wantUser
}

func verifyInternalDocsPassword(got string) bool {
	want := utils.EnvFirst("SW_PASS")
	return got == want
}
