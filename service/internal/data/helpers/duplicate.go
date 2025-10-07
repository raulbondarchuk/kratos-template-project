package data_helpers

import "strings"

func IsDuplicateError(err error) bool {
	return strings.Contains(err.Error(), "duplicate") ||
		strings.Contains(err.Error(), "Duplicate") ||
		strings.Contains(err.Error(), "Error 1062")
}
