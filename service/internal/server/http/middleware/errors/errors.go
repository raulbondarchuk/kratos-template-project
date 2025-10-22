package http_errors

import (
	"fmt"
	"net/http"
	dberr "service/internal/data/helpers"
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

// FromStatusAndError builds a Kratos HTTP error from a status code and an error.
// - If err is nil, it returns nil.
// - If err is already a *kerrors.Error, it attaches fields (if any) and returns it as-is.
// - Otherwise it maps the status code to the corresponding constructor (BadRequest, Conflict, etc).
func FromStatusAndError(status int, reason string, err error, f Fields) error {
	if err == nil {
		return nil
	}
	// Preserve existing Kratos errors to avoid double-wrapping.
	var ke *kerrors.Error
	if errorsAs(err, &ke) { // tiny wrapper below to avoid importing "errors" here
		return withFields(ke, f)
	}

	msg := err.Error()
	switch status {
	case http.StatusBadRequest:
		return BadRequest(reason, msg, f)
	case http.StatusUnauthorized:
		return Unauthorized(reason, msg, f)
	case http.StatusForbidden:
		return Forbidden(reason, msg, f)
	case http.StatusNotFound:
		return NotFound(reason, msg, f)
	case http.StatusConflict:
		return Conflict(reason, msg, f)
	case http.StatusUnprocessableEntity:
		return Unprocessable(reason, msg, f)
	case http.StatusServiceUnavailable:
		return Unavailable(reason, msg, f)
	case http.StatusRequestTimeout, http.StatusGatewayTimeout:
		// Map timeouts to 503 by default (tweak if you want exact passthrough).
		return Unavailable(reason, msg, f)
	default:
		return Internal(reason, msg, f)
	}
}

// FromDBError inspects a DB error and builds a proper HTTP error using your dberr mapper.
// Example mapping:
//   - duplicate key -> 409
//   - FK missing parent -> 400
//   - record not found -> 404
//   - deadlock -> 409 (or 503 by your policy)
//   - unknown -> 500
func FromDBError(reason string, err error, f Fields) error {
	if err == nil {
		return nil
	}
	status := dberr.HTTPStatusFromDBErr(err)
	return FromStatusAndError(status, reason, err, f)
}

// errorsAs is a tiny local wrapper to avoid importing "errors" at each call site.
func errorsAs(err error, target interface{}) bool {
	type asser interface{ As(interface{}) bool }
	if e, ok := any(err).(asser); ok {
		return e.As(target)
	}
	// fallback to stdlib if needed:
	// return errors.As(err, target)
	return false
}
