package paseto

import (
	"sync"

	"github.com/o1egl/paseto"
)

type Validator struct {
	secret       []byte
	legacySecret string

	newV2 *paseto.V2
}

var (
	validatorInstance *Validator
	onceValidator     sync.Once
)

func NewValidator() *Validator {
	onceValidator.Do(func() {

		currentSK := mustLoadEnv("SK_PASETO")
		legacySK := mustLoadEnv("SK_PASETO")

		secret := []byte(currentSK)
		if len(secret) != 32 {
			secret = sha256KeyFromString(currentSK)
		}

		validatorInstance = &Validator{
			secret:       secret,
			newV2:        paseto.NewV2(),
			legacySecret: legacySK,
		}
	})
	return validatorInstance
}

func (v *Validator) verifyAccessTokenClaims(token string) (*Claims, error) {
	raw := SkipBearer(token)

	if c, err := v.tryNew(raw); err == nil {
		return c, nil
	}

	if c, err := v.tryLegacy(raw); err == nil {
		return c, nil
	}
	return nil, ErrInvalidSignature
}
