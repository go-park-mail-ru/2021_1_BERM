package postgresstore

import (
	"FL_2/model"
	"github.com/lib/pq"
	"github.com/pkg/errors"
)

type VacancyRepository struct {
	store *Store
}

const (
	insertVacancy = `INSERT INTO ff.vacancy (
	                  category, 
	                  vacancy_name,
	                  description, 
	                  salary,
                      user_id
	                  )
	       VALUES ($1, $2, $3,$4, $5) RETURNING id`

	selectVacancyByID = "SELECT * FROM ff.vacancy WHERE id=$1"
)

func (v *VacancyRepository) Create(vacancy model.Vacancy) (uint64, error) {
	var vacancyID uint64
	err := v.store.Db.QueryRow(
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
				return 0, errors.Wrap(&DuplicateSourceErr{
					Err: err,
				}, sqlDbSourceError)
			}
		}
		return 0, errors.Wrap(err, sqlDbSourceError)
	}
	return vacancyID, nil
}

func (v *VacancyRepository) FindByID(id uint64) (*model.Vacancy, error) {
	vacancy := model.Vacancy{}
	err := v.store.Db.Get(&vacancy, selectVacancyByID, id)
	if err != nil {
		return nil, errors.Wrap(err, sqlDbSourceError)
	}
	return &vacancy, nil
}
