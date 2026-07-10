package entity

import "time"

type Merchant struct {
	Id                 int64     `validate:""`
	Name               string    `validate:"required"`
	ApiKeyHash         string    `validate:"required"`
	Domain             string    `validate:"required"`
	WebhookUrl         string    `validate:"required"`
	SuccessRedirectUrl string    `validate:"required"`
	FailureRedirectUrl string    `validate:"required"`
	IsActive           bool      `validate:""`
	CreatedAt          time.Time `validate:""`
	UpdatedAt          time.Time `validate:""`
}
