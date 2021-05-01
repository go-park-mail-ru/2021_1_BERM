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
	insertResponse = `INSERT INTO post.responses (
                            post_id, 
                            user_id, 
                            rate, 
                            time,
                            order_response,
                            vacancy_response,
                            text
		)
        VALUES (
                $1, 
                $2, 
                $3,
				$4,
                $5,
                $6,
                $7
                ) RETURNING id`

	selectOrderResponseByPostID = "SELECT * FROM post.responses WHERE post_id = $1 AND order_response = true"

	selectVacancyResponseByPostID = "SELECT * FROM post.responses WHERE post_id = $1 AND vacancy_response = true"

	updateOrderResponse = `UPDATE post.responses SET 
                          rate=:rate,
                          time=:time,
                          text=:text
                          WHERE user_id=:user_id AND post_id=:post_id AND order_response = true`

	updateVacancyResponse = `UPDATE post.responses SET 
                          rate=:rate,
                          time=:time,
                          text=:text
							WHERE user_id=:user_id AND post_id=:post_id AND vacancy_response = true`

	deleteOrderResponse = `DELETE FROM post.responses 
				 WHERE user_id=:user_id AND post_id=:post_id AND order_response = true`

	deleteVacancyResponse = `DELETE FROM post.responses 
				 WHERE user_id=:user_id AND post_id=:post_id AND vacancy_response = true`
)

type Repository struct {
	db *sqlx.DB
}

func NewRepo(db *sqlx.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) Create(response models.Response) (uint64, error) {
	var responseID uint64
	err := r.db.QueryRow(
		insertResponse,
		response.PostID,
		response.UserID,
		response.Rate,
		response.Time,
		response.OrderResponse,
		response.VacancyResponse,
		response.Text).Scan(&responseID)
	if err != nil {
		return 0, postgresql.WrapPostgreError(err)
	}

	return responseID, nil
}

func (r *Repository) FindByOrderPostID(id uint64) ([]models.Response, error) {
	var responses []models.Response
	if err := r.db.Select(&responses, selectOrderResponseByPostID, id); err != nil {
		return nil, postgresql.WrapPostgreError(err)
	}
	return responses, nil
}

func (r *Repository) FindByVacancyPostID(id uint64) ([]models.Response, error) {
	var responses []models.Response
	if err := r.db.Select(&responses, selectVacancyResponseByPostID, id); err != nil {
		return nil, postgresql.WrapPostgreError(err)
	}
	return responses, nil
}

func (r *Repository) ChangeOrderResponse(response models.Response) (*models.Response, error) {
	tx := r.db.MustBegin()
	_, err := tx.NamedExec(updateOrderResponse, &response)
	if err != nil {
		return nil, postgresql.WrapPostgreError(err)
	}
	if err := tx.Commit(); err != nil {
		return nil, postgresql.WrapPostgreError(err)
	}
	return &response, nil
}

func (r *Repository) ChangeVacancyResponse(response models.Response) (*models.Response, error) {
	tx := r.db.MustBegin()
	_, err := tx.NamedExec(updateVacancyResponse, &response)
	if err != nil {
		return nil, postgresql.WrapPostgreError(err)
	}
	if err := tx.Commit(); err != nil {
		return nil, postgresql.WrapPostgreError(err)
	}
	return &response, nil
}

func (r *Repository) DeleteOrderResponse(response models.Response) error {
	tx := r.db.MustBegin()
	_, err := tx.NamedExec(deleteOrderResponse, &response)
	if err != nil {
		return postgresql.WrapPostgreError(err)
	}
	if err = tx.Commit(); err != nil {
		return postgresql.WrapPostgreError(err)
	}
	return nil
}

func (r *Repository) DeleteVacancyResponse(response models.Response) error {
	tx := r.db.MustBegin()
	_, err := tx.NamedExec(deleteVacancyResponse, &response)
	if err != nil {
		return postgresql.WrapPostgreError(err)
	}
	if err = tx.Commit(); err != nil {
		return postgresql.WrapPostgreError(err)
	}
	return nil
}
