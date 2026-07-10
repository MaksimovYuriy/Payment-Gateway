package entity

import (
	"fmt"
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
	Id            int64
	MerchantId    int64
	SuccessBankId *int64
	OrderId       string
	Amount        string
	Currency      string
	Status        PaymentStatus
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func (p *Payment) Validate() error {
	amount, err := strconv.ParseFloat(p.Amount, 64)
	if err != nil || amount <= 0 {
		return fmt.Errorf("Incorrect amount")
	}
	if strings.TrimSpace(p.OrderId) == "" {
		return fmt.Errorf("Incorrect order id")
	}
	if len(p.Currency) != 3 {
		return fmt.Errorf("Incorrect currency")
	}
	return nil
}

func (p *Payment) SetStatusProcessing() {
	p.Status = PaymentStatusProcessing
}
