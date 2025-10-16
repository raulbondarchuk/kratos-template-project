package paseto

import "time"

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

func (c *Claims) valid() bool { return time.Now().Unix() < c.Exp }

// legacy-helper structure (only inside the package)
type LegacyClaims struct {
	Username      string    `json:"username"`
	CompanyID     int       `json:"companyId"`
	CompanyName   string    `json:"companyName"`
	Roles         string    `json:"roles"`
	OwnerUsername string    `json:"ownerUsername"`
	Cliuser       *string   `json:"cliuser,omitempty"`
	IssuedAt      time.Time `json:"iat"`
	ExpiresAt     time.Time `json:"exp"`
}
