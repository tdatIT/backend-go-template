package datetime

import (
	"errors"
	"strings"
	"time"
)

type TimeRequest string

func (t TimeRequest) ToTime() (time.Time, error) {
	s := strings.TrimSpace(string(t))
	if s == "" {
		return time.Time{}, errors.New("empty time string")
	}

	formats := []string{
		time.RFC3339,
		"2006-01-02 15:04",
		"2006-01-02",
	}

	for _, f := range formats {
		if parsed, err := time.Parse(f, s); err == nil {
			return parsed, nil
		}
	}
	return time.Time{}, errors.New("invalid time format: " + s)
}
