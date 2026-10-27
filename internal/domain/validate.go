package domain

import (
	"fmt"
	"regexp"
)

const MaxEventTypeLen = 100

var typeRe = regexp.MustCompile(`^[a-z0-9_]+(\.[a-z0-9_]+)*$`)

func ValidateEvent(eventType string) error {
	if eventType == "" {
		return fmt.Errorf("%w: empty", ErrInvalidEventType)
	}
	if len(eventType) > MaxEventTypeLen {
		return fmt.Errorf("%w: len out of limit", ErrInvalidEventType)
	}
	if !typeRe.MatchString(eventType) {
		return fmt.Errorf("%w: %q", ErrInvalidEventType, eventType)
	}
	return nil
}
