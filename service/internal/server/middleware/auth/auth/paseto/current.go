package paseto

import "encoding/json"

func (v *Validator) tryNew(raw string) (*Claims, error) {
	var claims Claims
	var footer string

	if err := v.newV2.Decrypt(raw, v.secret, &claims, &footer); err != nil {
		return nil, err
	}
	if !claims.valid() {
		return nil, ErrTokenExpired
	}
	return &claims, nil
}

func (v *Validator) verifyAccessTokenMap(token string) (map[string]interface{}, error) {
	c, err := v.verifyAccessTokenClaims(token)
	if err != nil {
		return nil, err
	}
	var m map[string]interface{}
	b, _ := json.Marshal(c)
	_ = json.Unmarshal(b, &m)
	return m, nil
}

func (v *Validator) verifyAccessTokenRoles(token string) (string, error) {
	c, err := v.verifyAccessTokenClaims(token)
	if err != nil {
		return "", err
	}
	return c.Roles, nil
}
