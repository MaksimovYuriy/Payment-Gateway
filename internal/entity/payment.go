package entity

import (
	"payment_gateway/internal/lib/apperr"
	"strconv"
	"strings"
	"time"
)

type PaymentStatus string

const (
	PaymentStatusCreated    PaymentStatus = "created"
	PaymentStatusProcessing PaymentStatus = "processing"
	PaymentStatusCompleted  PaymentStatus = "completed"
	PaymentStatusFailed     PaymentStatus = "failed"
	PaymentStatusCancelled  PaymentStatus = "cancelled"
)

type Payment struct {
	Id            int64         `validate:""`
	MerchantId    int64         `validate:"required"`
	SuccessBankId *int64        `validate:""`
	OrderId       string        `validate:"required"`
	Amount        string        `validate:"required"`
	Currency      string        `validate:"required"`
	Status        PaymentStatus `validate:""`
	CreatedAt     time.Time     `validate:""`
	UpdatedAt     time.Time     `validate:""`
}

func (p *Payment) Validate() error {
	amount, err := strconv.ParseFloat(p.Amount, 64)
	if err != nil || amount <= 0 {
		return apperr.InvalidInput("incorrect amount")
	}
	if strings.TrimSpace(p.OrderId) == "" {
		return apperr.InvalidInput("incorrect order id")
	}
	if len(p.Currency) != 3 {
		return apperr.InvalidInput("incorrect currency")
	}
	return nil
}

func (p *Payment) SetStatusProcessing() {
	p.Status = PaymentStatusProcessing
}
