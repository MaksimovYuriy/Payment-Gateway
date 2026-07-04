package entity

import "time"

type Merchant struct {
	Id                 int64
	Name               string
	ApiKeyHash         string
	Domain             string
	WebhookUrl         string
	SuccessRedirectUrl string
	FailureRedirectUrl string
	IsActive           bool
	CreatedAt          time.Time
	UpdatedAt          time.Time
}
