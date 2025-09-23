package http_errors

import (
	"fmt"
	"service/pkg/logger"

	kerrors "github.com/go-kratos/kratos/v2/errors"
)

type Fields map[string]string

func withFields(e *kerrors.Error, f Fields) error {
	if f == nil {
		return e
	}
	// Filter out empty values
	nonEmptyFields := make(Fields)
	for k, v := range f {
		if v != "" {
			nonEmptyFields[k] = v
		}
	}
	if len(nonEmptyFields) == 0 {
		return e
	}
	return e.WithMetadata(map[string]string(nonEmptyFields))
}

// 400 Bad Request
func BadRequest(reason, msg string, f Fields) error {
	logger.Warn("Bad Request", map[string]interface{}{"reason": reason, "msg": msg, "fields": f})
	return withFields(kerrors.BadRequest(reason, msg), f)
}
func BadRequestf(reason, format string, a ...any) error {
	logger.Warn("Bad Request format", map[string]interface{}{"reason": reason, "format": format, "a": a})
	return kerrors.BadRequest(reason, fmt.Sprintf(format, a...))
}

// 401 Unauthorized
func Unauthorized(reason, msg string, f Fields) error {
	logger.Warn("Unauthorized", map[string]interface{}{"reason": reason, "msg": msg, "fields": f})
	return withFields(kerrors.Unauthorized(reason, msg), f)
}

// 403 Forbidden
func Forbidden(reason, msg string, f Fields) error {
	logger.Warn("Forbidden", map[string]interface{}{"reason": reason, "msg": msg, "fields": f})
	return withFields(kerrors.Forbidden(reason, msg), f)
}

// 404 Not Found
func NotFound(reason, msg string, f Fields) error {
	logger.Warn("Not Found", map[string]interface{}{"reason": reason, "msg": msg, "fields": f})
	return withFields(kerrors.NotFound(reason, msg), f)
}

// 409 Conflict
func Conflict(reason, msg string, f Fields) error {
	logger.Warn("Conflict", map[string]interface{}{"reason": reason, "msg": msg, "fields": f})
	return withFields(kerrors.Conflict(reason, msg), f)
}

// 422 Unprocessable Entity
func Unprocessable(reason, msg string, f Fields) error {
	logger.Warn("Unprocessable", map[string]interface{}{"reason": reason, "msg": msg, "fields": f})
	return withFields(kerrors.New(422, reason, msg), f)
}

// 500 Internal Server Error
func Internal(reason, msg string, f Fields) error {
	logger.Error("Internal", map[string]interface{}{"reason": reason, "msg": msg, "fields": f})
	return withFields(kerrors.InternalServer(reason, msg), f)
}

// 503 Service Unavailable
func Unavailable(reason, msg string, f Fields) error {
	logger.Error("Unavailable", map[string]interface{}{"reason": reason, "msg": msg, "fields": f})
	return withFields(kerrors.ServiceUnavailable(reason, msg), f)
}
