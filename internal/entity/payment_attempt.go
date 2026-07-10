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
	Id                int64          `validate:""`
	PaymentId         int64          `validate:"required"`
	BankId            int64          `validate:"required"`
	Status            PAttemptStatus `validate:"required"`
	ExternalPaymentId *string        `validate:"omitempty"`
	ErrorMessage      *string        `validate:"omitempty"`
	ErrorCode         *string        `validate:"omitempty"`
	CreatedAt         time.Time      `validate:""`
	UpdatedAt         time.Time      `validate:""`
}

func (pa *PaymentAttempt) SetStatusProcessing() {
	pa.Status = PAttemptStatusProcessing
}
