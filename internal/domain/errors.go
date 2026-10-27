package domain

import "errors"

var (
	ErrInvalidEventType = errors.New("invalid event type")
	ErrPayloadTooLarge  = errors.New("payload too large")
	ErrEndpointNotFound = errors.New("endpoint not found")
)
