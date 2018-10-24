package augit

import (
	"encoding/json"
	"fmt"
)

// APIError is the wrapper type for ErrorDetails, mostly for JSON reasons
type APIError struct {
	Error ErrorDetails `json:"error"`
}

// ErrorDetails contains an error code int and message string
type ErrorDetails struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// LogAndFormatError takes an error code and error message and returns a
// JSON byte slice containing both. Also prints a readable version of the
// error details for logging.
// Note that this is an app-specific error code, not an HTTP status code
// Can be used just as a logger by discarding the returned value
func LogAndFormatError(code int, message string) []byte {
	fmt.Printf("[ERROR] %d - %s\n", code, message)
	errBytes, err := json.Marshal(&APIError{
		Error: ErrorDetails{
			Code:    code,
			Message: message,
		},
	})
	if err != nil {
		return []byte{}
	}
	return errBytes
}
