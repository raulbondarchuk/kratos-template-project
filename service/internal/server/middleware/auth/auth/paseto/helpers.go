package paseto

import (
	"crypto/sha256"
	"errors"
	"os"
	"strings"
)

var (
	ErrInvalidSignature = errors.New("invalid signature")
	ErrNotAnAccessToken = errors.New("not an access token")
	ErrTokenExpired     = errors.New("token expired")
)

func SkipBearer(token string) string {
	token = strings.TrimSpace(token)
	if strings.HasPrefix(strings.ToLower(token), "bearer ") {
		return strings.TrimSpace(token[7:])
	}
	return token
}

func mustLoadEnv(name string) string {
	v := os.Getenv(name)
	if v == "" {
		panic("missing env: " + name)
	}
	return v
}

func sha256KeyFromString(base string) []byte {
	sum := sha256.Sum256([]byte(base))
	return sum[:]
}
