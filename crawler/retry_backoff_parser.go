package crawler

import (
	"errors"
	"time"
	"strconv"
	"regexp"
)

// RetryBackoffParser parses a string representation of a retry/backoff configuration.
type RetryBackoffParser struct{}

// Parse parses the given input string into a RetryBackoffConfig.
func (p *RetryBackoffParser) Parse(input string) (RetryBackoffConfig, error) {
	re := regexp.MustCompile(`^(\d+)(\.\d+)?(ms|s)$`)
	match := re.FindStringSubmatch(input)

	if len(match) == 0 {
		return RetryBackoffConfig{}, errors.New("invalid retry/backoff config format")
	}

	initialTime, err := time.ParseDuration(match[1] + match[2])
	if err != nil {
		return RetryBackoffConfig{}, err
	}

	multiplier, err := strconv.ParseFloat(match[2], 64)
	if err != nil || multiplier < 1 {
		multiplier = 1
	}

	return RetryBackoffConfig{
		Initial:    initialTime,
		Multiplier: multiplier,
	}, nil
}
