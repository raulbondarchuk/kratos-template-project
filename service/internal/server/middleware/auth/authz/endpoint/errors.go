package endpoint

import "errors"

const (
	ReasonAuthz = "AUTHORIZATION"
)

var (
	ErrMissingAuthorizationHeader = errors.New("missing authorization header")
	ErrInvalidToken               = errors.New("invalid token")
	ErrInsufficientPermissions    = errors.New("insufficient permissions")
)
