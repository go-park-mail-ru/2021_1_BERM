package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	customErr "post/internal/app/errors"
	"post/internal/app/models"
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
                   time
		)
        VALUES (
                $1, 
                $2, 
                $3,
				$4
                ) RETURNING id`

	selectResponseByPostID = "SELECT * FROM post.responses WHERE post_id = $1"

	updateResponse = `UPDATE post.responses SET 
                 rate=:rate,
                 time=:time
				 WHERE user_id=:user_id AND post_id=:post_id`

	deleteResponse = `DELETE FROM post.responses 
				 WHERE user_id=:user_id AND post_id=:post_id`
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
		response.UserLogin,
		response.UserImg,
		response.Time).Scan(&responseID)
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

	return responseID, nil
}

func (r *Repository) FindByPostID(id uint64) ([]models.Response, error) {
	var responses []models.Response
	if err := r.db.Select(&responses, selectResponseByPostID, id); err != nil {
		return nil, errors.Wrap(err, sqlDbSourceError)
	}
	return responses, nil
}

func (r *Repository) Change(response models.Response) (*models.Response, error) {
	tx := r.db.MustBegin()
	_, err := tx.NamedExec(updateResponse, &response)
	if err != nil {
		return nil, errors.Wrap(err, sqlDbSourceError)
	}
	if err := tx.Commit(); err != nil {
		return nil, errors.Wrap(err, sqlDbSourceError)
	}
	return &response, nil
}

func (r *Repository) Delete(response models.Response) error {
	tx := r.db.MustBegin()
	_, err := tx.NamedExec(deleteResponse, &response)
	if err != nil {
		return errors.Wrap(err, sqlDbSourceError)
	}
	if err = tx.Commit(); err != nil {
		return errors.Wrap(err, sqlDbSourceError)
	}
	return nil
}
