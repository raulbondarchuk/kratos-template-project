package dberr

import (
	"context"
	"errors"
	"net/http"

	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

// HTTPStatusFromDBErr maps a database error to the corresponding HTTP status code.
// Returns 500 (Internal Server Error) if the error cannot be classified.
func HTTPStatusFromDBErr(err error) int {
	if err == nil {
		return http.StatusOK
	}

	// Context-related errors
	if errors.Is(err, context.DeadlineExceeded) {
		return http.StatusGatewayTimeout // 504
	}
	if errors.Is(err, context.Canceled) {
		return http.StatusRequestTimeout // 408
	}

	// GORM: record not found
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return http.StatusNotFound // 404
	}

	// MySQL-specific errors
	var me *mysql.MySQLError
	if errors.As(err, &me) {
		switch me.Number {

		// Unique or duplicate constraint violations
		case 1062: // ER_DUP_ENTRY
			return http.StatusConflict // 409

		// Foreign key constraint violations
		case 1452: // ER_NO_REFERENCED_ROW_2: add/update child → missing parent
			return http.StatusBadRequest // 400 (invalid parent reference)
		case 1451, 1217: // ER_ROW_IS_REFERENCED_2 / ER_ROW_IS_REFERENCED: delete/update parent blocked
			return http.StatusConflict // 409

		// Invalid or incomplete data
		case 1048: // ER_BAD_NULL_ERROR: column cannot be null
			return http.StatusBadRequest // 400
		case 1364: // ER_NO_DEFAULT_FOR_FIELD
			return http.StatusBadRequest // 400
		case 1406: // ER_DATA_TOO_LONG
			return http.StatusBadRequest // 400
		case 1292: // ER_TRUNCATED_WRONG_VALUE
			return http.StatusBadRequest // 400

		// Database lock or timeout issues
		case 1205: // ER_LOCK_WAIT_TIMEOUT
			return http.StatusServiceUnavailable // 503 (retry possible)
		case 1213: // ER_LOCK_DEADLOCK
			return http.StatusConflict // 409 (or 503, depending on retry policy)

		// Schema or infrastructure errors
		case 1146: // ER_NO_SUCH_TABLE
			return http.StatusInternalServerError // 500
		case 1054: // ER_BAD_FIELD_ERROR
			return http.StatusInternalServerError // 500
		}
	}

	// Default: internal server error
	return http.StatusInternalServerError
}
