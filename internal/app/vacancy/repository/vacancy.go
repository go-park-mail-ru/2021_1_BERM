package repository

import (
	"ff/internal/app/models"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	ffErrors "ff/internal/app/server/errors"
)

type VacancyRepository struct {
	db *sqlx.DB
}

const (
	insertVacancy = `INSERT INTO ff.vacancy (
	                  category, 
	                  vacancy_name,
	                  description, 
	                  salary,
                      user_id
	                  )
	       VALUES ($1, $2, $3,$4, $5) RETURNING id`

	selectVacancyByID = "SELECT * FROM ff.vacancy WHERE id=$1"
)

func (v *VacancyRepository) Create(vacancy models.Vacancy) (uint64, error) {
	var vacancyID uint64
	err := v.db.QueryRow(
		insertVacancy,
		vacancy.Category,
		vacancy.VacancyName,
		vacancy.Description,
		vacancy.Salary,
		vacancy.UserID).Scan(&vacancyID)
	if err != nil {
		pqErr := &pq.Error{}
		if errors.As(err, &pqErr) {
			if pqErr.Code == ffErrors.DuplicateErrorCode {
				return 0, errors.Wrap(&ffErrors.DuplicateSourceErr{
					Err: err,
				}, ffErrors.SqlDbSourceError)
			}
		}
		return 0, errors.Wrap(err, ffErrors.SqlDbSourceError)
	}
	return vacancyID, nil
}

func (v *VacancyRepository) FindByID(id uint64) (*models.Vacancy, error) {
	vacancy := models.Vacancy{}
	err := v.db.Get(&vacancy, selectVacancyByID, id)
	if err != nil {
		return nil, errors.Wrap(err, ffErrors.SqlDbSourceError)
	}
	return &vacancy, nil
}

