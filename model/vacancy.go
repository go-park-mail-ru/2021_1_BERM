package model

type Vacancy struct {
	Id          uint64 `json:"id,omitempty" db:"id"`
	UserId      uint64 `json:"user_id,omitempty" db:"user_id"`
	Category    string `json:"category,omitempty" db:"category"`
	VacancyName string `json:"vacancy_name,omitempty" db:"vacancy_name"`
	Description string `json:"description,omitempty" db:"description"`
	Salary      uint64 `json:"salary,omitempty" db:"salary"`
	Login       string `json:"login,omitempty"`
	Img         string `json:"img,omitempty"`
}
