package request

type RegistrationMerchant struct {
	Name               string `json:"name" validate:"required,max=255"`
	Domain             string `json:"domain" validate:"required,max=255"`
	WebhookUrl         string `json:"webhook_url" validate:"required"`
	SuccessRedirectUrl string `json:"success_redirect_url" validate:"required"`
	FailureRedirectUrl string `json:"failure_redirect_url" validate:"required"`
	IsActive           bool   `json:"is_active"`
}
