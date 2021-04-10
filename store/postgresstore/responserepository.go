package postgresstore

import (
	"FL_2/model"
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
func (r *ResponseRepository)Create(response model.Response) (uint64, error){
	var responseID uint64
	err := r.store.db.QueryRow(
		`INSERT INTO response (
                   order_id, 
                   user_id_id, 
                   rare, 
                   user_login, 
                   user_img, 
		)
        VALUES (
                $1, 
                $2, 
                $3,
				$4,
				$5,
                ) RETURNING id`,
                response.OrderID,
                response.UserID,
                response.Rate,
                response.UserLogin,
                response.UserImg, ).Scan(&responseID)
	if err != nil {
		return 0, err
	}

	return responseID, nil
}
func (r *ResponseRepository)FindById(id uint64)  ([]model.Response, error){
	var responses []model.Response
	if err := r.store.db.Select(&responses, "SELECT * FROM response WHERE order_id = $1", id); err != nil {
		return nil, err
	}
	return responses, nil
}
