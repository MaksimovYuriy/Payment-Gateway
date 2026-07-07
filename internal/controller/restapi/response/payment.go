package response

import (
	"payment_gateway/internal/entity"
	"time"
)

type Payment struct {
	Id            int64     `json:"id"`
	MerchantId    int64     `json:"merchant_id"`
	SuccessBankId *int64    `json:"success_bank_id"`
	OrderId       string    `json:"order_id"`
	Amount        string    `json:"amount"`
	Currency      string    `json:"currency"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type ListPayment struct {
	Payments []Payment `json:"payments"`
}

func NewPayment(payment entity.Payment) Payment {
	return Payment{
		Id:            payment.Id,
		MerchantId:    payment.MerchantId,
		SuccessBankId: payment.SuccessBankId,
		OrderId:       payment.OrderId,
		Amount:        payment.Amount,
		Currency:      payment.Currency,
		Status:        string(payment.Status),
		CreatedAt:     payment.CreatedAt,
		UpdatedAt:     payment.UpdatedAt,
	}
}

func NewListPayment(payments []*entity.Payment) ListPayment {
	responseList := make([]Payment, 0, len(payments))
	for _, payment := range payments {
		responseList = append(responseList, NewPayment(*payment))
	}
	return ListPayment{Payments: responseList}
}
