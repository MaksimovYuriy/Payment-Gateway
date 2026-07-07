package bankprocessor

import (
	"context"
	"payment_gateway/internal/entity"
	"time"
)

type BankProcessor interface {
	Process(ctx context.Context, payment *entity.Payment, bankId int64) Result
}

type Result struct {
	Success           bool
	ExternalPaymentId string
	ErrorCode         string
	ErrorMessage      string
}

type MockProcessor struct {
	delay time.Duration
}

func NewMockProcessor(delay time.Duration) *MockProcessor {
	return &MockProcessor{delay: delay}
}

func (p *MockProcessor) Process(ctx context.Context, payment *entity.Payment, bankId int64) Result {
	select {
	case <-time.After(p.delay):
		return Result{
			Success:           true,
			ExternalPaymentId: "mock-payment-id",
		}
	case <-ctx.Done():
		return Result{
			Success:      false,
			ErrorCode:    "context_cancelled",
			ErrorMessage: ctx.Err().Error(),
		}
	}
}
