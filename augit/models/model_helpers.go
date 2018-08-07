package models

import "strings"

const (
	// PopErrNotFoundFragment is a string found when there are no rows in a result set for a query
	PopErrNotFoundFragment = "sql: no rows in result set"
)

// IsErrRecordNotFound returns true if the error string contains PopNotFoundFragment
func IsErrRecordNotFound(err error) bool {
	return strings.Contains(err.Error(), PopErrNotFoundFragment)
}
