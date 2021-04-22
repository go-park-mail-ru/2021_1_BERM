package models

type Response struct {
	ID              uint64 `json:"id,omitempty" db:"id"`
	Time            uint64 `json:"time,omitempty" db:"time"`
	PostID          uint64 `json:"post_id,omitempty" db:"order_id"`
	UserID          uint64 `json:"user_id,omitempty" db:"user_id"`
	Rate            uint64 `json:"rate,omitempty" db:"rate"`
	UserLogin       string `json:"user_login,omitempty" db:"user_login"`
	UserImg         string `json:"user_img,omitempty" db:"user_img"`
	OrderResponse   bool   `json:"order_response,omitempty"`
	VacancyResponse bool   `json:"vacancy_response,omitempty"`
}
