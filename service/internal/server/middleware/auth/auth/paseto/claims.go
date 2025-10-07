package paseto

import (
	"errors"
	"strings"
	"time"
)

// Claims
type Claims struct {
	Type      string `json:"type"`
	CompanyID uint   `json:"company_id"`
	Username  string `json:"username"`
	CliUser   string `json:"cliuser"`
	Owner     string `json:"owner"`
	Roles     string `json:"roles"`
	Expired   string `json:"expired"`
	Exp       int64  `json:"exp"`
	Iat       int64  `json:"iat"`
}

func (c *Claims) valid() bool {
	return time.Now().Unix() < c.Exp
}

func SkipBearer(token string) string {
	token = strings.TrimSpace(token)
	if strings.HasPrefix(strings.ToLower(token), "bearer ") {
		return strings.TrimSpace(token[7:])
	}
	return token
}

var (
	ErrInvalidSignature = errors.New("invalid signature")
	ErrNotAnAccessToken = errors.New("not an access token")
	ErrTokenExpired     = errors.New("token expired")
)
