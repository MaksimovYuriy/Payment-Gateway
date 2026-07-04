package entity

import "time"

type Bank struct {
	Id        int64
	Code      string
	Name      string
	IsActive  bool
	CreatedAt time.Time
	UpdatedAt time.Time
}
