package repository

import (
	"github.com/jmoiron/sqlx"
	"post/internal/app/models"
	"post/pkg/postgresql"
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
		return 0, postgresql.WrapPostgreError(err)
	}
	return vacancyID, nil
}

func (r *Repository) FindByID(id uint64) (*models.Vacancy, error) {
	vacancy := models.Vacancy{}
	err := r.db.Get(&vacancy, selectVacancyByID, id)
	if err != nil {
		return nil, postgresql.WrapPostgreError(err)
	}
	return &vacancy, nil
}



