package entity

import "time"

type PAttemptStatus string

const (
	PAttemptStatusCreated    PAttemptStatus = "created"
	PAttemptStatusProcessing PAttemptStatus = "processing"
	PAttemptStatusSucceeded  PAttemptStatus = "succeeded"
	PAttemptStatusFailed     PAttemptStatus = "failed"
	PAttemptStatusTimeout    PAttemptStatus = "timeout"
	PAttemptStatusCancelled  PAttemptStatus = "cancelled"
)

type PaymentAttempt struct {
	Id                int64
	PaymentId         int64
	BankId            int64
	Status            PAttemptStatus
	ExternalPaymentId *string
	ErrorMessage      *string
	ErrorCode         *string
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

func (pa *PaymentAttempt) SetStatusProcessing() {
	pa.Status = PAttemptStatusProcessing
}
