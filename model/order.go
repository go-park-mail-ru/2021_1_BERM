package model

import (
	validation "github.com/go-ozzo/ozzo-validation"
)

type Order struct {
	ID          uint64   `json:"id,omitempty" db:"id"`
	OrderName   string   `json:"order_name" db:"order_name"`
	CustomerID  uint64   `json:"customer_id"`
	ExecutorID	uint64	 `json:"executor_id"`
	Budget      uint64   `json:"budget"`
	Deadline    uint64   `json:"deadline"`
	Description string   `json:"description"`
	Category 	string 	 `json:"category"`
}


func (o *Order) Validate() error {
	return validation.ValidateStruct(
		o,
		validation.Field(&o.OrderName, validation.Required,  validation.Length(5, 300)),
		validation.Field(&o.Description, validation.Required,),
		validation.Field(&o.Category, validation.Required),
	)
}