package model

import (
	validation "github.com/go-ozzo/ozzo-validation"
)

const (
	MinOrderNameLength int = 5
	MaxOrderNameLength int = 300
)

type Order struct {
	ID          uint64   `json:"id,omitempty"`
	OrderName   string   `json:"order_name"`
	CustomerID  uint64   `json:"customer_id"`
	Budget      uint64   `json:"budget"`
	Deadline    uint64   `json:"deadline"`
	Description string   `json:"description"`
	Specializes []string `json:"specializes"`
}

func (o *Order) Validate() error {
	return validation.ValidateStruct(
		o,
		validation.Field(&o.OrderName, validation.Required, validation.Length(MinOrderNameLength, MaxOrderNameLength)),
		validation.Field(&o.Description, validation.Required),
		validation.Field(&o.Specializes, validation.Required),
	)
}
