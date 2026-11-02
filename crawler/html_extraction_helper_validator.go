package crawler

import (
	"bytes"
	"context"
	"io"
	"net/http"

	"golang.org/x/net/html"
)

// HTMLValidationError represents an error during HTML extraction validation.
type HTMLValidationError struct {
	URL    string
	Reason string
}

func (e HTMLValidationError) Error() string {
	return "HTML validation failed: " + e.Reason
}

// ValidateHTML validates the HTML extracted from a URL and returns an error if validation fails.
func ValidateHTML(ctx context.Context, url, htmlContent string) error {
	// Parse HTML to check if it's well-formed.
	_, err := html.Parse(strings.NewReader(htmlContent))
	if err != nil {
		return HTMLValidationError{URL: url, Reason: "HTML is not well-formed"}
	}

	// Add more validation rules as needed, e.g., check if charset is supported.

	return nil
}

// ValidateResponseHTML validates the HTML extracted from an HTTP response.
func ValidateResponseHTML(ctx context.Context, resp *http.Response) error {
	if resp.StatusCode != http.StatusOK {
		return HTMLValidationError{URL: resp.Request.URL.String(), Reason: "Non-200 status code"}
	}
	defer resp.Body.Close()

	htmlContent, err := io.ReadAll(resp.Body)
	if err != nil {
		return HTMLValidationError{URL: resp.Request.URL.String(), Reason: "Failed to read response body"}
	}

	return ValidateHTML(ctx, resp.Request.URL.String(), string(htmlContent))
}
