package paseto

import (
	"strconv"
	"time"

	o1 "github.com/o1egl/paseto"
)

func (v *Validator) tryLegacy(raw string) (*Claims, error) {
	var jt o1.JSONToken
	var footer string

	key := sha256KeyFromString(v.legacySecret)
	if err := o1.NewV2().Decrypt(raw, key, &jt, &footer); err != nil {
		return nil, err
	}

	lc := &LegacyClaims{
		Username:  jt.Subject,
		IssuedAt:  jt.IssuedAt,
		ExpiresAt: jt.Expiration,
		Roles:     jt.Get("roles"),
	}
	if v := jt.Get("ownerUsername"); v != "" {
		lc.OwnerUsername = v
	}
	if v := jt.Get("cliuser"); v != "" {
		lc.Cliuser = &v
	}
	if v := jt.Get("companyId"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			lc.CompanyID = n
		}
	}

	// map to new format Claims
	if time.Now().After(lc.ExpiresAt) {
		return nil, ErrTokenExpired
	}
	if lc.Cliuser == nil {
		lc.Cliuser = new(string)
	}
	return &Claims{
		Type:      "access",
		CompanyID: uint(lc.CompanyID),
		Username:  lc.Username,
		CliUser:   *lc.Cliuser,
		Owner:     lc.OwnerUsername,
		Roles:     lc.Roles,
		Exp:       lc.ExpiresAt.Unix(),
		Iat:       lc.IssuedAt.Unix(),
		Expired:   lc.ExpiresAt.Format(time.RFC3339),
	}, nil
}
