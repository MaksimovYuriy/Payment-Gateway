package response

import (
	"time"

	"payment_gateway/internal/entity"
)

type Bank struct {
	Id        int64     `json:"id"`
	Code      string    `json:"code"`
	Name      string    `json:"name"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ListBank struct {
	Banks []Bank `json:"banks"`
}

func NewBank(bank entity.Bank) Bank {
	return Bank{
		Id:        bank.Id,
		Code:      bank.Code,
		Name:      bank.Name,
		IsActive:  bank.IsActive,
		CreatedAt: bank.CreatedAt,
		UpdatedAt: bank.UpdatedAt,
	}
}

func NewListBank(banks []*entity.Bank) ListBank {
	list := make([]Bank, 0, len(banks))
	for _, bank := range banks {
		list = append(list, NewBank(*bank))
	}

	return ListBank{Banks: list}
}
