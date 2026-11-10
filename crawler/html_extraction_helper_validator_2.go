package crawler

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

// HTMLExtractionHelperValidator validates HTML extraction helpers.
type HTMLExtractionHelperValidator struct {
	regexps map[string]*regexp.Regexp
}

// NewHTMLExtractionHelperValidator creates a new HTMLExtractionHelperValidator.
func NewHTMLExtractionHelperValidator() *HTMLExtractionHelperValidator {
	return &HTMLExtractionHelperValidator{
		regexps: map[string]*regexp.Regexp{
			"url":            regexp.MustCompile(`^https?://[^\s,()<>]+$`),
			"extractor-name": regexp.MustCompile(`^[a-zA-Z0-9_-]+$`),
		},
	}
}

// ValidateURL validates the URL format.
func (v *HTMLExtractionHelperValidator) ValidateURL(u string) error {
	if !v.regexps["url"].MatchString(u) {
		return fmt.Errorf("invalid URL format: %s", u)
	}
	return nil
}

// ValidateExtractorName validates the extractor name format.
func (v *HTMLExtractionHelperValidator) ValidateExtractorName(n string) error {
	if !v.regexps["extractor-name"].MatchString(n) {
		return fmt.Errorf("invalid extractor name format: %s", n)
	}
	return nil
}

// ValidateExtractors validates the list of extractors.
func (v *HTMLExtractionHelperValidator) ValidateExtractors(extractors []HTMLExtractionHelper) error {
	for _, e := range extractors {
		if err := v.ValidateURL(e.URL); err != nil {
			return fmt.Errorf("invalid URL: %w", err)
		}
		if err := v.ValidateExtractorName(e.Name); err != nil {
			return fmt.Errorf("invalid extractor name: %w", err)
		}
	}
	return nil
}
