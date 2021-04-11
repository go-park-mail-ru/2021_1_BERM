package model

type Order struct {
	ID          uint64 `json:"id,omitempty" db:"id"`
	OrderName   string `json:"order_name" db:"order_name"`
	CustomerID  uint64 `json:"customer_id,omitempty"`
	ExecutorID  uint64 `json:"executor_id,omitempty"`
	Budget      uint64 `json:"budget"`
	Deadline    uint64 `json:"deadline"`
	Description string `json:"description"`
	Category    string `json:"category"`
	Login       string `json:"login,omitempty"`
	Img         string `json:"img,omitempty"`
}
