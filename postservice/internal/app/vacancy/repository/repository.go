package repository

import (
	"FL_2/model"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	"post/internal/app/models"
	customErr "post/internal/app/errors"
)

const (
	duplicateErrorCode = "23505"
	sqlDbSourceError   = "SQL sb source error"
)

const (
	insertVacancy = `INSERT INTO post.vacancy (
						  category, 
						  vacancy_name,
						  description, 
						  salary,
						  user_id
	                  )
	       VALUES ($1, $2, $3,$4, $5) RETURNING id`

	selectVacancyByID = "SELECT * FROM post.vacancy WHERE id=$1"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepo(db *sqlx.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) Create(vacancy models.Vacancy) (uint64, error) {
	var vacancyID uint64
	err := r.db.QueryRow(
		insertVacancy,
		vacancy.Category,
		vacancy.VacancyName,
		vacancy.Description,
		vacancy.Salary,
		vacancy.UserID).Scan(&vacancyID)
	if err != nil {
		pqErr := &pq.Error{}
		if errors.As(err, &pqErr) {
			if pqErr.Code == duplicateErrorCode {
				return 0, errors.Wrap(&customErr.DuplicateSourceErr{
					Err: err,
				}, sqlDbSourceError)
			}
		}
		return 0, errors.Wrap(err, sqlDbSourceError)
	}
	return vacancyID, nil
}

func (r *Repository) FindByID(id uint64) (*models.Vacancy, error) {
	vacancy := models.Vacancy{}
	err := r.db.Get(&vacancy, selectVacancyByID, id)
	if err != nil {
		return nil, errors.Wrap(err, sqlDbSourceError)
	}
	return &vacancy, nil
}



