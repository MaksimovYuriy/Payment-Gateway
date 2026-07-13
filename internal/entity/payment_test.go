package entity_test

import (
	"errors"
	"testing"

	"payment_gateway/internal/entity"
	"payment_gateway/internal/lib/apperr"
)

func TestPaymentValidate(t *testing.T) {
	tests := []struct {
		name    string
		payment entity.Payment
		wantErr bool
	}{
		{
			name: "valid payment",
			payment: entity.Payment{
				MerchantId: 1,
				OrderId:    "order-1",
				Amount:     "10.50",
				Currency:   "USD",
			},
			wantErr: false,
		},
		{
			name: "invalid amount",
			payment: entity.Payment{
				MerchantId: 1,
				OrderId:    "order-1",
				Amount:     "0",
				Currency:   "USD",
			},
			wantErr: true,
		},
		{
			name: "empty order id",
			payment: entity.Payment{
				MerchantId: 1,
				OrderId:    " ",
				Amount:     "10.50",
				Currency:   "USD",
			},
			wantErr: true,
		},
		{
			name: "invalid currency",
			payment: entity.Payment{
				MerchantId: 1,
				OrderId:    "order-1",
				Amount:     "10.50",
				Currency:   "US",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.payment.Validate()
			if (err != nil) != tt.wantErr {
				t.Fatalf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil {
				var appErr *apperr.Error
				if !errors.As(err, &appErr) {
					t.Fatalf("Validate() error type = %T, want *apperr.Error", err)
				}
				if appErr.Code != apperr.CodeInvalidInput {
					t.Fatalf("Validate() error code = %s, want %s", appErr.Code, apperr.CodeInvalidInput)
				}
			}
		})
	}
}
