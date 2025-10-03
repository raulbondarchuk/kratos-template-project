package swagger

import (
	"bytes"
	"encoding/json"
	"fmt"
	stdhttp "net/http"
	"service/internal/server/middleware/auth/auth/paseto"
	"strings"
	"time"

	kratoshttp "github.com/go-kratos/kratos/v2/transport/http"
)

func isAuthed(r *stdhttp.Request) bool {
	c, err := r.Cookie(cookieName)
	return err == nil && strings.TrimSpace(c.Value) != ""
}

func setSessionCookie(w stdhttp.ResponseWriter, value string) {
	stdhttp.SetCookie(w, &stdhttp.Cookie{
		Name:     cookieName,
		Value:    value, // the token
		Path:     cookiePath,
		HttpOnly: true,
		SameSite: stdhttp.SameSiteLaxMode,
		MaxAge:   3600,
		// Secure: true, // enable for HTTPS
	})
}

func clearSessionCookie(w stdhttp.ResponseWriter) {
	stdhttp.SetCookie(w, &stdhttp.Cookie{
		Name:     cookieName,
		Value:    "",
		Path:     cookiePath,
		HttpOnly: true,
		SameSite: stdhttp.SameSiteLaxMode,
		MaxAge:   -1,
		Expires:  time.Unix(0, 0),
	})
}

func authRequired(h stdhttp.HandlerFunc) stdhttp.HandlerFunc {
	return func(w stdhttp.ResponseWriter, r *stdhttp.Request) {
		if !isAuthed(r) {
			serveLoginPage(w, r, "")
			return
		}
		h(w, r)
	}
}

// authenticateWithAPI do a request to loginURL, extract the token from Authorization,
// validate the roles and return the token + success flag.
func authenticateWithAPI(username, password string) (string, bool) {
	payload := map[string]string{"username": username, "password": password}
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return "", false
	}

	req, err := stdhttp.NewRequest(stdhttp.MethodPost, loginURL, bytes.NewBuffer(jsonData))
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

	// 1) get the token from Authorization
	raw := strings.TrimSpace(resp.Header.Get("Authorization"))
	token := extractBearerToken(raw)
	if token == "" {
		return "", false
	}

	// 2) validate the roles, which are stored in the token
	rolesStr, err := paseto.VerifyAccessTokenRoles(token) // expect only the token, not the whole header
	if err != nil {
		return "", false
	}
	roles := parseRoles(rolesStr) // normalize to []string

	// 3) require at least one role
	if !hasAnyRole(roles, requiredRoles) {
		return "", false
	}

	return token, true
}

// Trim "Bearer " if it exists
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

// Parse roles like "admin, user" / "admin user" / "admin;user"
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

// at least one of the requiredRoles is present
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

func attachBootstrap(s *kratoshttp.Server) {
	s.HandleFunc("/swagger/bootstrap.js", authRequired(func(w stdhttp.ResponseWriter, r *stdhttp.Request) {
		c, _ := r.Cookie(cookieName)
		token := ""
		if c != nil {
			token = strings.TrimSpace(c.Value)
		}

		setOnlyContentType(w, "application/javascript; charset=utf-8")
		setNoCache(w)
		w.WriteHeader(stdhttp.StatusOK)

		// t — token from HttpOnly-cookie (session). Work with localStorage keys:
		// K        = token, which is used by Swagger UI
		// K_LAST   = last "seeded" token from the server (for recognizing new authorization)
		// K_OVR    = flag, that the user manually changed the token
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
        // New authorization → rewrite localStorage with a fresh token
        localStorage.setItem(K, t);
        localStorage.setItem(K_LAST, t);
        localStorage.removeItem(K_OVR); // reset the flag, because this is a new session
        s = t;
      } else if (!s) {
        // Old session, but localStorage is empty → seed it
        localStorage.setItem(K, t);
        s = t;
      }
      // If t === last and there is an override — do nothing (respect the user's value)
    }

    // Set the current value from localStorage to the input
    var i=document.getElementById('authToken'); 
    if(i){ i.value = localStorage.getItem(K) || s || ""; }
  }catch(e){}
})();`, token)

		_, _ = w.Write([]byte(js))
	}))
}
