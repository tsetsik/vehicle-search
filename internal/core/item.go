package core

import "github.com/go-playground/validator"

type (
	Item struct {
		Make     string `json:"make" validate:"required"`
		Model    string `json:"model" validate:"required"`
		Year     int    `json:"year" validate:"required,gte=1886,lte=2026"`
		FuelType string `json:"fuel_type" validate:"required,oneof=petrol diesel electric hybrid"`
	}
)

func (i *Item) Validate() error {
	validator := validator.New()
	return validator.Struct(i)
}
