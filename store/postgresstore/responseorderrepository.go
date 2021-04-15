package postgresstore

import (
	"FL_2/model"
	"github.com/lib/pq"
	"github.com/pkg/errors"
)

type ResponseOrderRepository struct {
	store *Store
}

const (
	insertOrderResponse = `INSERT INTO ff.order_responses (
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
                ) RETURNING id`

	selectOrderResponseByOrderID = "SELECT * FROM ff.order_responses WHERE order_id = $1"

	updateOrderResponse = `UPDATE ff.order_responses SET 
                 rate=:rate,
                 time=:time
				 WHERE user_id=:user_id AND order_id=:order_id`

	deleteOrderResponse = `DELETE FROM ff.order_responses 
				 WHERE user_id=:user_id AND order_id=:order_id`
)

func (r *ResponseOrderRepository) Create(response model.ResponseOrder) (uint64, error) {
	var responseID uint64
	err := r.store.db.QueryRow(
		insertOrderResponse,
		response.OrderID,
		response.UserID,
		response.Rate,
		response.UserLogin,
		response.UserImg,
		response.Time).Scan(&responseID)
	if err != nil {
		pqErr := &pq.Error{}
		if errors.As(err, &pqErr) {
			if pqErr.Code == duplicateErrorCode {
				return 0, errors.Wrap(&DuplicateSourceErr{
					Err: err,
				}, sqlDbSourceError)
			}
		}
		return 0, errors.Wrap(err, sqlDbSourceError)
	}

	return responseID, nil
}

func (r *ResponseOrderRepository) FindByOrderID(id uint64) ([]model.ResponseOrder, error) {
	var responses []model.ResponseOrder
	if err := r.store.db.Select(&responses, selectOrderResponseByOrderID, id); err != nil {
		return nil, errors.Wrap(err, sqlDbSourceError)
	}
	return responses, nil
}

func (r *ResponseOrderRepository) Change(response model.ResponseOrder) (*model.ResponseOrder, error) {
	tx := r.store.db.MustBegin()
	_, err := tx.NamedExec(updateOrderResponse, &response)
	if err != nil {
		return nil, errors.Wrap(err, sqlDbSourceError)
	}
	if err := tx.Commit(); err != nil {
		return nil, errors.Wrap(err, sqlDbSourceError)
	}
	return &response, nil
}

func (r *ResponseOrderRepository) Delete(response model.ResponseOrder) error {
	tx := r.store.db.MustBegin()
	_, err := tx.NamedExec(deleteOrderResponse, &response)
	if err != nil {
		return errors.Wrap(err, sqlDbSourceError)
	}
	if err = tx.Commit(); err != nil {
		return errors.Wrap(err, sqlDbSourceError)
	}
	return nil
}
