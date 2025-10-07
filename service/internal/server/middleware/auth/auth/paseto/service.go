package paseto

import (
	"os"
	"sync"

	"github.com/o1egl/paseto"
)

type Validator struct {
	secret []byte
	paseto *paseto.V2
}

var (
	validatorInstance *Validator
	onceValidator     sync.Once
)

// NewValidator creates a new validator
func NewValidator() *Validator {
	onceValidator.Do(func() {
		secret := []byte(os.Getenv("SK_PASETO"))
		if len(secret) != 32 {
			panic("SK_PASETO must be exactly 32 bytes for V2.Local")
		}
		validatorInstance = &Validator{
			secret: secret,
			paseto: paseto.NewV2(),
		}
	})
	return validatorInstance
}
