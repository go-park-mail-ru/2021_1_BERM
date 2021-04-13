package postgresstore

import (
	"FL_2/model"
)

type ResponseOrderRepository struct {
	store *Store
}

func (r *ResponseOrderRepository) Create(response model.ResponseOrder) (uint64, error) {
	var responseID uint64
	err := r.store.db.QueryRow(
		`INSERT INTO order_responses (
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
		return 0, err
	}

	return responseID, nil
}

func (r *ResponseOrderRepository) FindByOrderId(id uint64) ([]model.ResponseOrder, error) {
	var responses []model.ResponseOrder
	if err := r.store.db.Select(&responses, "SELECT * FROM order_responses WHERE order_id = $1", id); err != nil {
		return nil, err
	}
	return responses, nil
}

func (r *ResponseOrderRepository) Change(response model.ResponseOrder) (*model.ResponseOrder, error) {
	tx := r.store.db.MustBegin()
	_, err := tx.NamedExec(`UPDATE order_responses SET 
                 rate=:rate,
                 time=:time
				 WHERE user_id=:user_id AND order_id=:order_id`, &response)
	if err != nil {
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return &response, nil
}

func (r *ResponseOrderRepository) Delete(response model.ResponseOrder) error {
	tx := r.store.db.MustBegin()
	_, err := tx.NamedExec(`DELETE FROM order_responses 
				 WHERE user_id=:user_id AND order_id=:order_id`, &response)
	if err != nil {
		return  err
	}
	if err := tx.Commit(); err != nil {
		return  err
	}
	return nil
}
