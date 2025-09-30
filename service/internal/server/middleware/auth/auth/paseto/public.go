package paseto

func VerifyAccessToken(token string) (*Claims, error) {
	return NewValidator().verifyAccessTokenClaims(token)
}

func VerifyAccessTokenMap(token string) (map[string]interface{}, error) {
	return NewValidator().verifyAccessTokenMap(token)
}

func VerifyAccessTokenRoles(token string) (string, error) {
	return NewValidator().verifyAccessTokenRoles(token)
}
