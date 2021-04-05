package model

type Vacancy struct{
	Id          uint64 `json:"id,omitempty" db:"id"`
	Category    string `json:"category,omitempty" db:"category"`
	VacancyName string `json:"vacancy_name,omitempty" db:"vacancy_name"`
	Description string `json:"description,omitempty" db:"description"`
	Salary      uint64 `json:"salary,omitempty" db:"salary"`
}
