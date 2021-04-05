package postgresstore

import "FL_2/model"

type VacancyRepository struct {
	store *Store
}



func (v *VacancyRepository)Create(vacancy model.Vacancy) (uint64, error){
	var vacansyID uint64
	vacansyID = 0
		err := v.store.db.QueryRow(
			`INSERT INTO vacancy (
	                  category, 
	                  vacancy_name 
	                  description, 
	                  salary
	                  )
	       VALUES ($1, $2, $3,$4) RETURNING id`,
			vacancy.Category,
			vacancy.VacancyName,
			vacancy.Description,
			vacancy.Salary).Scan(&vacansyID)
	return vacansyID, err
}

func (v *VacancyRepository)FindByID(id uint64) (*model.Vacancy, error){
	vacancy := &model.Vacancy{}
	err := v.store.db.Get(&vacancy, "SELECT * FROM vacancy WHERE id=$1", id)
	return vacancy, err
}
