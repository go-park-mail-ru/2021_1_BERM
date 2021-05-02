package models

type Review struct {
	ID              uint64 `json:"id,omitempty" db:"id"`
	UserId          uint64 `json:"user_id" db:"user_id"`
	ToUserId        uint64 `json:"to_user_id" db:"to_user_id"`
	OrderId         uint64 `json:"order_id" db:"order_id"`
	Description     string `json:"text" db:"description"`
	Score           uint8  `json:"score" db:"score"`
	OrderName       string `json:"order_name,omitempty"`
	UserLogin       string `json:"user_login,omitempty"`
	UserNameSurname string `json:"user_name_surname,omitempty"`
}
