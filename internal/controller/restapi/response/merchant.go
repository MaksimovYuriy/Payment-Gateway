package response

import (
	"time"

	"payment_gateway/internal/entity"
)

type Merchant struct {
	Id                 int64     `json:"id"`
	Name               string    `json:"name"`
	Domain             string    `json:"domain"`
	WebhookUrl         string    `json:"webhook_url"`
	SuccessRedirectUrl string    `json:"success_redirect_url"`
	FailureRedirectUrl string    `json:"failure_redirect_url"`
	IsActive           bool      `json:"is_active"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

type ListMerchant struct {
	Merchants []Merchant `json:"merchants"`
}

type MerchantRegistration struct {
	Merchant Merchant `json:"merchant"`
	ApiKey   string   `json:"api_key"`
}

func NewMerchant(merchant entity.Merchant) Merchant {
	return Merchant{
		Id:                 merchant.Id,
		Name:               merchant.Name,
		Domain:             merchant.Domain,
		WebhookUrl:         merchant.WebhookUrl,
		SuccessRedirectUrl: merchant.SuccessRedirectUrl,
		FailureRedirectUrl: merchant.FailureRedirectUrl,
		IsActive:           merchant.IsActive,
		CreatedAt:          merchant.CreatedAt,
		UpdatedAt:          merchant.UpdatedAt,
	}
}

func NewListMerchant(merchants []*entity.Merchant) ListMerchant {
	list := make([]Merchant, 0, len(merchants))
	for _, merchant := range merchants {
		list = append(list, NewMerchant(*merchant))
	}

	return ListMerchant{Merchants: list}
}

func NewMerchantRegistration(merchant entity.Merchant, apiKey string) MerchantRegistration {
	return MerchantRegistration{
		Merchant: NewMerchant(merchant),
		ApiKey:   apiKey,
	}
}
