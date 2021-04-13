package postgresstore

import (
	"FL_2/model"
	"github.com/lib/pq"
	"github.com/pkg/errors"
)

type ResponseRepository struct {
	store *Store
}

//id            SERIAL PRIMARY KEY NOT NULL,
//order_id      INTEGER            NOT NULL,
//user_id       INTEGER            NOT NULL,
//rate          INTEGER            NOT NULL,
//user_login VARCHAR            NOT NULL,
//user_img      VARCHAR DEFAULT '',
func (r *ResponseRepository) Create(response model.Response) (uint64, error) {
	var responseID uint64
	err := r.store.db.QueryRow(
		`INSERT INTO responses (
                   order_id, 
                   user_id, 
                   rate, 
                   user_login, 
                   user_img, 
                   time
		)
        VALUES (
                $1, 
                $2, 
                $3,
				$4,
				$5,
                $6
                ) RETURNING id`,
		response.OrderID,
		response.UserID,
		response.Rate,
		response.UserLogin,
		response.UserImg,
		response.Time).Scan(&responseID)
	if err != nil {
		pqErr := &pq.Error{}
		if errors.As(err, &pqErr){
			if pqErr.Code == duplicateErrorCode{
				return 0, errors.Wrap(&DuplicateSourceErr{
					Err: err,
				}, sqlDbSourceError)
			}
		}
		return 0, errors.Wrap(err, sqlDbSourceError)
	}

	return responseID, nil
}

func (r *ResponseRepository) FindById(id uint64) ([]model.Response, error) {
	var responses []model.Response
	if err := r.store.db.Select(&responses, "SELECT * FROM responses WHERE order_id = $1", id); err != nil {
		return nil, errors.Wrap(err, sqlDbSourceError)
	}
	return responses, nil
}

func (r *ResponseRepository) Change(response model.Response) (*model.Response, error) {
	tx := r.store.db.MustBegin()
	_, err := tx.NamedExec(`UPDATE responses SET 
                 rate=:rate,
                 time=:time
				 WHERE user_id=:user_id AND order_id=:order_id`, &response)
	if err != nil {
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, errors.Wrap(err, sqlDbSourceError)
	}
	return &response, nil
}

func (r *ResponseRepository) Delete(response model.Response) error {
	tx := r.store.db.MustBegin()
	_, err := tx.NamedExec(`DELETE FROM responses 
				 WHERE user_id=:user_id AND order_id=:order_id`, &response)
	if err != nil {
		return errors.Wrap(err, sqlDbSourceError)
	}
	if err = tx.Commit(); err != nil {
		return errors.Wrap(err, sqlDbSourceError)
	}
	return nil
}
