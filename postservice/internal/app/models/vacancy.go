package models

type Vacancy struct {
	ID          uint64 `json:"id,omitempty" db:"id"`
	CustomerID  uint64 `json:"customer_id,omitempty" db:"customer_id"`
	ExecutorID  uint64 `json:"executor_id,omitempty" db:"executor_id"`
	Category    string `json:"category,omitempty" db:"category"`
	VacancyName string `json:"vacancy_name,omitempty" db:"vacancy_name"`
	Description string `json:"description,omitempty" db:"description"`
	Salary      uint64 `json:"salary,omitempty" db:"salary"`
	Login       string `json:"user_login,omitempty"`
	Img         string `json:"user_img,omitempty"`
	IsArchived  bool   `json:"is_archived" db:"is_archived"`
}

type VacancySearch struct {
	Keyword string `json:"keyword"`
}

type SuggestVacancyTittle struct {
	Title string `json:"title" db:"order_name"`
}
