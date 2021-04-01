package model

import (
	validation "github.com/go-ozzo/ozzo-validation"
)

type Order struct {
	Id            uint64   	`json:"id"`
	OrderName     string   	`json:"order_name"`
	CustomerId    uint64   	`json:"customer_id"`
	Description   string    `json:"description"`
	Specializes  []string	`json:"specializes"`
}


func (o *Order) Validate() error {
	return validation.ValidateStruct(
		o,
		validation.Field(&o.OrderName, validation.Required,  validation.Length(5, 300)),
		validation.Field(&o.Description, validation.Required,),
		validation.Field(&o.Specializes, validation.Required),
	)
}