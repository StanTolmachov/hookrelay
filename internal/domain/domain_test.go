package domain

import (
	"bytes"
	"errors"
	"strings"
	"testing"
	"time"
	"uuid"
)

func TestNewEvent(t *testing.T) {
	tests := []struct {
		name      string
		eventType string
		payload   []byte
		wantErr   error
	}{
		{
			name:      "ok",
			eventType: "order.created",
			payload:   []byte(strings.Repeat("s", MaxPayloadBytes)),
		},
		{
			name:      "payload out of limit",
			eventType: "order.created",
			payload:   []byte(strings.Repeat("s", MaxPayloadBytes+1)),
			wantErr:   ErrPayloadTooLarge,
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				before := time.Now()
				event, err := NewEvent(tt.eventType, tt.payload)
				after := time.Now()

				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("want %v, got %v", tt.wantErr, err)
				}

				if tt.wantErr != nil {
					if event != nil {
						t.Fatalf("want nil event, got %v", event)
					}
					return
				}

				if event.ID == uuid.Nil() {
					t.Fatalf("ID is nil")
				}
				if v := event.ID.String()[14]; v != '7' {
					t.Fatalf("ID version is %c, want 7", v)
				}
				if event.CreatedAt.Location() != time.UTC {
					t.Fatalf("CreatedAt location = %v, want UTC", event.CreatedAt.Location())
				}
				if event.CreatedAt.Before(before) || event.CreatedAt.After(after) {
					t.Fatalf("CreatedAt = %v, want between %v and %v", event.CreatedAt, before, after)
				}

				if event.Type != tt.eventType {
					t.Fatalf("event.Type = %v, want = %v", event.Type, tt.eventType)
				}
				if !bytes.Equal(event.Payload, tt.payload) {
					t.Fatalf("payload len = %d, want = %d", len(event.Payload), len(tt.payload))
				}
			},
		)
	}
}
