package domain

import (
	"encoding/json"
	"time"
	"uuid"
)

type DeliveryStatus string

const (
	DeliveryStatusPending    DeliveryStatus = "pending"
	DeliveryStatusDelivering DeliveryStatus = "delivering"
	DeliveryStatusDelivered  DeliveryStatus = "delivered"
	DeliveryStatusFailed     DeliveryStatus = "failed"
)

type Endpoint struct {
	ID         uuid.UUID
	URL        string
	Secret     string
	EventTypes []string
	IsActive   bool
	CreatedAt  time.Time
}

type Event struct {
	ID        uuid.UUID
	Type      string
	Payload   json.RawMessage
	CreatedAt time.Time
}

type Delivery struct {
	ID         uuid.UUID
	EventID    uuid.UUID
	EndpointID uuid.UUID
	Status     DeliveryStatus
	Attempts   int
	NextTryAt  time.Time
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type Attempt struct {
	ID         uuid.UUID
	DeliveryID uuid.UUID
	AttemptNo  int
	StatusCode *int
	DurationMS int
	Error      string
	CreatedAt  time.Time
}

func NewEvent(eventType string, payload []byte) (*Event, error) {
	err := ValidateEvent(eventType, payload)
	if err != nil {
		return nil, err
	}

	return &Event{
		ID:        uuid.NewV7(),
		Type:      eventType,
		Payload:   payload,
		CreatedAt: time.Now().UTC(),
	}, nil
}
