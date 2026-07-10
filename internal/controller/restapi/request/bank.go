package request

type RegistrationBank struct {
	Code     string `json:"code" validate:"required,max=32"`
	Name     string `json:"name" validate:"required,max=255"`
	IsActive bool   `json:"is_active"`
}
