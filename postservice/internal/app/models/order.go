package models

//easyjson:json
type OrderList []Order

//easyjson:json
type Order struct {
	ID          uint64 `json:"id,omitempty" db:"id"`
	OrderName   string `json:"order_name,omitempty" db:"order_name"`
	CustomerID  uint64 `json:"customer_id,omitempty" db:"customer_id"`
	ExecutorID  uint64 `json:"executor_id,omitempty" db:"executor_id"`
	Budget      uint64 `json:"budget,omitempty" db:"budget"`
	Deadline    uint64 `json:"deadline,omitempty" db:"deadline"`
	Description string `json:"description,omitempty" db:"description"`
	Category    string `json:"category,omitempty" db:"category"`
	UserLogin   string `json:"user_login,omitempty"`
	UserImg     string `json:"user_img,omitempty"`
	IsArchived  bool   `json:"is_archived" db:"is_archived"`
}

//easyjson:json
type OrderSearch struct {
	Keyword string `json:"keyword"`
}

//easyjson:json
type SuggestOrderTitleList []SuggestOrderTitle

//easyjson:json
type SuggestOrderTitle struct {
	Title string `json:"title" db:"order_name"`
}
