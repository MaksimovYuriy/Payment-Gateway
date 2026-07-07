package request

type CreatePayment struct {
	MerchantId int64  `json:"merchant_id" validate:"required"`
	BankId     int64  `json:"bank_id" validate:"required"`
	OrderId    string `json:"order_id" validate:"required,max=255"`
	Amount     string `json:"amount" validate:"required"`
	Currency   string `json:"currency" validate:"required,max=3"`
}
