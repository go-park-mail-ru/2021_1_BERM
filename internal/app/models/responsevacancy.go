package models

type ResponseVacancy struct {
	ID        uint64 `json:"id,omitempty" db:"id"`
	Time      uint64 `json:"time,omitempty" db:"time"`
	VacancyID uint64 `json:"vacancy_id,omitempty" db:"vacancy_id"`
	UserID    uint64 `json:"user_id,omitempty" db:"user_id"`
	Rate      uint64 `json:"rate,omitempty" db:"rate"`
	UserLogin string `json:"user_login,omitempty" db:"user_login"`
	UserImg   string `json:"user_img,omitempty" db:"user_img"`
}
