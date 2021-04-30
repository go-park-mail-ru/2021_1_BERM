package repository

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"post/internal/app/models"
	"post/pkg/Error"
	"post/pkg/postgresql"
)

const (
	insertVacancy = `INSERT INTO post.vacancy (
						  category, 
						  vacancy_name,
						  description, 
						  salary,
						  customer_id
	                  )
	       VALUES ($1, $2, $3,$4, $5) RETURNING id`

	selectVacancyByID = "SELECT * FROM post.vacancy WHERE id=$1"

	updateVacancy = `UPDATE post.vacancy SET
						vacancy_name =:vacancy_name,
						category =:category,
						customer_id =:customer_id,
						executor_id =:executor_id,
						salary =:salary,
						description =:description
						WHERE id =:id`

	deleteVacancy = `DELETE from post.vacancy WHERE id=$1`

	selectVacanciesByExecutorID = "SELECT * FROM post.vacancy WHERE executor_id=$1"

	selectVacanciesByCustomerID = "SELECT * FROM post.vacancy WHERE customer_id=$1"

	updateExecutor = `UPDATE post.vacancy SET 
                 executor_id =:executor_id
				 WHERE id = :id`
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
		vacancy.CustomerID).Scan(&vacancyID)
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

func (r *Repository) Change(vacancy models.Vacancy) error {
	tx, err := r.db.Begin()
	if err != nil {
		return postgresql.WrapPostgreError(err)
	}
	_, err = tx.Exec(updateVacancy, vacancy)
	if err != nil {
		return &Error.Error{
			Err:              err,
			InternalError:    true,
			ErrorDescription: Error.InternalServerErrorDescription,
		}
	}
	if err = tx.Commit(); err != nil {
		return postgresql.WrapPostgreError(err)
	}
	return nil
}

func (r *Repository) DeleteVacancy(id uint64) error {
	_, err := r.db.Queryx(deleteVacancy, id)
	if err != nil {
		return postgresql.WrapPostgreError(err)
	}
	return nil
}

func (r *Repository) FindByExecutorID(executorID uint64) ([]models.Vacancy, error) {
	var vacancies []models.Vacancy
	err := r.db.Select(&vacancies, selectVacanciesByExecutorID, executorID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, postgresql.WrapPostgreError(err)
	}
	return vacancies, nil
}

func (r *Repository) FindByCustomerID(customerID uint64) ([]models.Vacancy, error) {
	var vacancies []models.Vacancy
	err := r.db.Select(&vacancies, selectVacanciesByCustomerID, customerID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, postgresql.WrapPostgreError(err)
	}
	return vacancies, nil
}

func (r *Repository) UpdateExecutor(vacancy models.Vacancy) error {
	tx := r.db.MustBegin()
	_, err := tx.NamedExec(updateExecutor, &vacancy)
	if err != nil {
		return postgresql.WrapPostgreError(err)
	}
	if err := tx.Commit(); err != nil {
		return postgresql.WrapPostgreError(err)
	}
	return nil
}
