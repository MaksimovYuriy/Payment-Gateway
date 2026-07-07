package entity

import "time"

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
