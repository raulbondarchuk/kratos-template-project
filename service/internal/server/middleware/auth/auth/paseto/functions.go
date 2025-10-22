package paseto

import (
	"encoding/json"
)

func (v *Validator) verifyAccessTokenClaims(token string) (*Claims, error) {
	token = SkipBearer(token)

	var claims Claims
	var footer string

	if err := v.paseto.Decrypt(token, v.secret, &claims, &footer); err != nil {
		return nil, ErrInvalidSignature
	}
	if claims.Type != "access" {
		return nil, ErrNotAnAccessToken
	}
	if !claims.valid() {
		return nil, ErrTokenExpired
	}
	return &claims, nil
}

func (v *Validator) verifyAccessTokenMap(token string) (map[string]interface{}, error) {
	claims, err := v.verifyAccessTokenClaims(token)
	if err != nil {
		return nil, err
	}
	var m map[string]interface{}
	data, _ := json.Marshal(claims)
	_ = json.Unmarshal(data, &m)
	return m, nil
}

func (v *Validator) verifyAccessTokenRoles(token string) (string, error) {
	claims, err := v.verifyAccessTokenClaims(token)
	if err != nil {
		return "", err
	}
	return claims.Roles, nil
}
