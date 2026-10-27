package domain

import (
	"errors"
	"strings"
	"testing"
)

func TestValidateEvent(t *testing.T) {
	tests := []struct {
		name      string
		eventType string
		payload   []byte
		wantErr   error
	}{
		{
			name:      "type ok",
			eventType: "order.created",
		},
		{
			name:      "type space",
			eventType: "order created",
			wantErr:   ErrInvalidEventType,
		},
		{
			name:      "type capital",
			eventType: "Order created",
			wantErr:   ErrInvalidEventType,
		},
		{
			name:      "type empty",
			eventType: "",
			wantErr:   ErrInvalidEventType,
		},
		{
			name:      "type len exactly limit",
			eventType: strings.Repeat("s", MaxEventTypeLen),
		},
		{
			name:      "type len out of limit",
			eventType: strings.Repeat("s", MaxEventTypeLen+1),
			wantErr:   ErrInvalidEventType,
		},
		{
			name:      "type empty segment",
			eventType: "order..created",
			wantErr:   ErrInvalidEventType,
		},
		{
			name:      "type trailing dot",
			eventType: "order.",
			wantErr:   ErrInvalidEventType,
		},
		{
			name:      "type non-latin",
			eventType: "заказ.создан",
			wantErr:   ErrInvalidEventType,
		},
		{
			name:      "type underscore ok",
			eventType: "invoice_payment.failed",
		},
		{
			name:      "payload ok",
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
				err := ValidateEvent(tt.eventType, tt.payload)

				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("want %v, got %v", tt.wantErr, err)
				}

			},
		)
	}
}
