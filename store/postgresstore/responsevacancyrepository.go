package postgresstore

import (
	"FL_2/model"
	"github.com/lib/pq"
	"github.com/pkg/errors"
)

type ResponseVacancyRepository struct {
	store *Store
}

func (r *ResponseVacancyRepository) Create(response model.ResponseVacancy) (uint64, error) {
	var responseID uint64
	err := r.store.db.QueryRow(
		`INSERT INTO vacancy_responses (
                   vacancy_id, 
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
		response.VacancyID,
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

func (r *ResponseVacancyRepository) FindByVacancyID(id uint64) ([]model.ResponseVacancy, error) {
	var responses []model.ResponseVacancy
	if err := r.store.db.Select(&responses, "SELECT * FROM vacancy_responses WHERE vacancy_id = $1", id); err != nil {
		return nil, errors.Wrap(err, sqlDbSourceError)
	}
	return responses, nil
}

func (r *ResponseVacancyRepository) Change(response model.ResponseVacancy) (*model.ResponseVacancy, error) {
	tx := r.store.db.MustBegin()
	_, err := tx.NamedExec(`UPDATE vacancy_responses SET 
                 rate=:rate,
                 time=:time
				 WHERE user_id=:user_id AND vacancy_id=:order_id`, &response)
	if err != nil {
		return nil, errors.Wrap(err, sqlDbSourceError)
	}
	if err := tx.Commit(); err != nil {
		return nil, errors.Wrap(err, sqlDbSourceError)
	}
	return &response, nil
}

func (r *ResponseVacancyRepository) Delete(response model.ResponseVacancy) error {
	tx := r.store.db.MustBegin()
	_, err := tx.NamedExec(`DELETE FROM vacancy_responses 
				 WHERE user_id=:user_id AND vacancy_id=:order_id`, &response)
	if err != nil {
		return  errors.Wrap(err, sqlDbSourceError)
	}
	if err := tx.Commit(); err != nil {
		return  errors.Wrap(err, sqlDbSourceError)
	}
	return nil
}
