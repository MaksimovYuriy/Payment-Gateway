package entity

import "time"

type Bank struct {
	Id        int64     `validate:""`
	Code      string    `validate:"required"`
	Name      string    `validate:"required"`
	IsActive  bool      `validate:""`
	CreatedAt time.Time `validate:""`
	UpdatedAt time.Time `validate:""`
}
