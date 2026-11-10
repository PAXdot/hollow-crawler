package crawler

import (
	"fmt"
)

type (
	URLFrontierError struct {
		cause error
		frontierID string
	}
)

func NewURLFrontierError(err error, frontierID string) *URLFrontierError {
	return &URLFrontierError{
		cause:      err,
		frontierID: frontierID,
	}
}

func (e *URLFrontierError) Error() string {
	return fmt.Sprintf("URL frontier (%s) encountered an error: %v", e.frontierID, e.cause)
}

func IsURLFrontierError(err error) bool {
	_, ok := err.(*URLFrontierError)
	return ok
}
