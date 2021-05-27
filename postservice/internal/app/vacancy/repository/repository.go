package vacancy

import (
	"context"
	"database/sql"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"post/internal/app/models"
	"post/pkg/error/errortools"
	"post/pkg/types"
	"strings"
)

const (
	ctxParam      types.CtxKey = 4
	insertVacancy              = `INSERT INTO post.vacancy (
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

	selectArchiveVacancies = "SELECT * FROM post.archive_vacancy"

	insertArchiveVacancy = `INSERT INTO post.archive_vacancy (
                          id,
						  category, 
						  vacancy_name,
						  description, 
						  salary,
						  customer_id,
                          executor_id,
                          is_archived
	                  )
	       VALUES ($1, $2, $3,$4, $5, $6, $7, $8) RETURNING id`

	searchVacanciesInTitle = "SELECT * FROM post.vacancy WHERE to_tsvector(vacancy_name) @@ to_tsquery($1)"

	searchVacanciesInText = "SELECT * FROM post.vacancy WHERE to_tsvector(description) @@ to_tsquery($1)"

	getActualVacancy = "SELECT * FROM post.vacancy " +
		"WHERE CASE WHEN $1 != 0 THEN salary >= $1 ELSE true END " +
		"AND CASE WHEN $2 != 0  THEN salary <= $2 ELSE true END " +
		"AND CASE WHEN $3 != '~' THEN to_tsvector(vacancy_name) @@ to_tsquery($3) ELSE true END " +
		"AND CASE WHEN $4 != '~' THEN category = $4 ELSE true END " +
		"ORDER BY salary LIMIT $5 OFFSET $6"

	getActualVacancyDesc = "SELECT * FROM post.vacancy " +
		"WHERE CASE WHEN $1 != 0 THEN salary >= $1 ELSE true END " +
		"AND CASE WHEN $2 != 0  THEN salary <= $2 ELSE true END " +
		"AND CASE WHEN $3 != '~' THEN to_tsvector(vacancy_name) @@ to_tsquery($3) ELSE true END " +
		"AND CASE WHEN $4 != '~' THEN category = $4 ELSE true END " +
		"ORDER BY salary DESC LIMIT $5 OFFSET $6"

	selectArchiveVacancyByID = "SELECT * FROM post.archive_vacancy WHERE id=$1"

	selectArchiveVacanciesByExecutorID = "SELECT * FROM post.archive_vacancy WHERE executor_id=$1"

	selectArchiveVacanciesByCustomerID = "SELECT * FROM post.archive_vacancy WHERE customer_id=$1"

	selectTittle = `SELECT DISTINCT vacancy_name FROM post.vacancy WHERE vacancy_name LIKE $1 LIMIT 5`

	selectAllTittle = `SELECT DISTINCT vacancy_name FROM post.vacancy LIMIT 5`
)

type Repository struct {
	db *sqlx.DB
}

func NewRepo(db *sqlx.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) Create(vacancy models.Vacancy, ctx context.Context) (uint64, error) {
	var vacancyID uint64
	err := r.db.QueryRow(
		insertVacancy,
		vacancy.Category,
		vacancy.VacancyName,
		vacancy.Description,
		vacancy.Salary,
		vacancy.CustomerID).Scan(&vacancyID)
	if err != nil {
		customErr := errortools.SqlErrorChoice(err)
		return 0, errors.Wrap(customErr, err.Error())
	}
	return vacancyID, nil
}

func (r *Repository) FindByID(id uint64, ctx context.Context) (*models.Vacancy, error) {
	vacancy := models.Vacancy{}
	err := r.db.Get(&vacancy, selectVacancyByID, id)
	if err != nil {
		customErr := errortools.SqlErrorChoice(err)
		return nil, errors.Wrap(customErr, err.Error())
	}
	return &vacancy, nil
}

func (r *Repository) GetActualVacancies(ctx context.Context) ([]models.Vacancy, error) {
	var vacancies []models.Vacancy
	param := ctx.Value(ctxParam).(map[string]interface{})
	category := param["category"].(string)
	limit := param["limit"].(int)
	offset := param["offset"].(int)
	desk := param["desc"].(bool)
	salaryFrom := param["from"].(int)
	salaryTo := param["to"].(int)
	searchStr := param["search_str"].(string)
	if searchStr != "~" {
		search := strings.Split(searchStr, " ")
		var res string
		for i, s := range search {
			if i == len(search)-1 {
				res += " " + s
				break
			}
			res += s + " <->"
		}
		searchStr = res
		searchStr += ":*"
	}
	if desk {
		if err := r.db.Select(
			&vacancies,
			getActualVacancy,
			salaryFrom,
			salaryTo,
			searchStr,
			category,
			limit,
			offset); err != nil {
			customErr := errortools.SqlErrorChoice(err)
			return nil, errors.Wrap(customErr, err.Error())
		}
	} else {
		if err := r.db.Select(
			&vacancies,
			getActualVacancyDesc,
			salaryFrom,
			salaryTo,
			searchStr,
			category,
			limit,
			offset); err != nil {
			customErr := errortools.SqlErrorChoice(err)
			return nil, errors.Wrap(customErr, err.Error())
		}
	}
	return vacancies, nil
}

func (r *Repository) GetVacancyNum(ctx context.Context) (uint64, error) {
	var num uint64
	if err := r.db.Get(&num, "SELECT COUNT(id) FROM post.vacancy"); err != nil {
		customErr := errortools.SqlErrorChoice(err)
		return 0, errors.Wrap(customErr, err.Error())
	}
	return num, nil
}

func (r *Repository) Change(vacancy models.Vacancy, ctx context.Context) error {
	tx, err := r.db.Beginx()
	if err != nil {
		customErr := errortools.SqlErrorChoice(err)
		return errors.Wrap(customErr, err.Error())
	}
	_, err = tx.NamedExec(updateVacancy, &vacancy)
	if err != nil {
		customErr := errortools.SqlErrorChoice(err)
		return errors.Wrap(customErr, err.Error())
	}
	if err = tx.Commit(); err != nil {
		customErr := errortools.SqlErrorChoice(err)
		return errors.Wrap(customErr, err.Error())
	}
	return nil
}

func (r *Repository) DeleteVacancy(id uint64, ctx context.Context) error {
	_, err := r.db.Queryx(deleteVacancy, id)
	if err != nil {
		customErr := errortools.SqlErrorChoice(err)
		return errors.Wrap(customErr, err.Error())
	}
	return nil
}

func (r *Repository) FindByExecutorID(executorID uint64, ctx context.Context) ([]models.Vacancy, error) {
	var vacancies []models.Vacancy
	err := r.db.Select(&vacancies, selectVacanciesByExecutorID, executorID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		customErr := errortools.SqlErrorChoice(err)
		return nil, errors.Wrap(customErr, err.Error())
	}
	return vacancies, nil
}

func (r *Repository) FindByCustomerID(customerID uint64, ctx context.Context) ([]models.Vacancy, error) {
	var vacancies []models.Vacancy
	err := r.db.Select(&vacancies, selectVacanciesByCustomerID, customerID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		customErr := errortools.SqlErrorChoice(err)
		return nil, errors.Wrap(customErr, err.Error())
	}
	return vacancies, nil
}

func (r *Repository) UpdateExecutor(vacancy models.Vacancy, ctx context.Context) error {
	tx, err := r.db.Beginx()
	if err != nil {
		customErr := errortools.SqlErrorChoice(err)
		return errors.Wrap(customErr, err.Error())
	}
	_, err = tx.NamedExec(updateExecutor, &vacancy)
	if err != nil {
		customErr := errortools.SqlErrorChoice(err)
		return errors.Wrap(customErr, err.Error())
	}
	if err := tx.Commit(); err != nil {
		customErr := errortools.SqlErrorChoice(err)
		return errors.Wrap(customErr, err.Error())
	}
	return nil
}

func (r *Repository) CreateArchive(vacancy models.Vacancy, ctx context.Context) (uint64, error) {
	var vacancyID uint64
	vacancy.IsArchived = true
	err := r.db.QueryRow(
		insertArchiveVacancy,
		vacancy.ID,
		vacancy.Category,
		vacancy.VacancyName,
		vacancy.Description,
		vacancy.Salary,
		vacancy.CustomerID,
		vacancy.ExecutorID,
		vacancy.IsArchived).Scan(&vacancyID)
	if err != nil {
		customErr := errortools.SqlErrorChoice(err)
		return 0, errors.Wrap(customErr, err.Error())
	}
	return vacancyID, nil
}

func (r *Repository) GetArchiveVacancies(ctx context.Context) ([]models.Vacancy, error) {
	var vacancies []models.Vacancy
	if err := r.db.Select(&vacancies, selectArchiveVacancies); err != nil {
		customErr := errortools.SqlErrorChoice(err)
		return nil, errors.Wrap(customErr, err.Error())
	}
	return vacancies, nil
}

func (r *Repository) SearchVacancy(keyword string, ctx context.Context) ([]models.Vacancy, error) {
	var vacancies []models.Vacancy
	if keyword == "" {
		return nil, nil
	}
	keyword += ":*"
	if err := r.db.Select(&vacancies, searchVacanciesInTitle, keyword); err != nil {
		customErr := errortools.SqlErrorChoice(err)
		return nil, errors.Wrap(customErr, err.Error())
	}
	if len(vacancies) == 0 {
		if err := r.db.Select(&vacancies, searchVacanciesInText, keyword); err != nil {
			customErr := errortools.SqlErrorChoice(err)
			return nil, errors.Wrap(customErr, err.Error())
		}
	}
	return vacancies, nil
}

func (r *Repository) FindArchiveByID(id uint64, ctx context.Context) (*models.Vacancy, error) {
	vacancy := models.Vacancy{}
	if err := r.db.Get(&vacancy, selectArchiveVacancyByID, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		customErr := errortools.SqlErrorChoice(err)
		return nil, errors.Wrap(customErr, err.Error())
	}
	return &vacancy, nil
}

func (r *Repository) GetArchiveVacanciesByExecutorID(executorID uint64, ctx context.Context) ([]models.Vacancy, error) {
	var vacancies []models.Vacancy
	if err := r.db.Select(&vacancies, selectArchiveVacanciesByExecutorID, executorID); err != nil {
		customErr := errortools.SqlErrorChoice(err)
		return nil, errors.Wrap(customErr, err.Error())
	}
	return vacancies, nil
}

func (r *Repository) GetArchiveVacanciesByCustomerID(customerID uint64, ctx context.Context) ([]models.Vacancy, error) {
	var vacancies []models.Vacancy
	if err := r.db.Select(&vacancies, selectArchiveVacanciesByCustomerID, customerID); err != nil {
		customErr := errortools.SqlErrorChoice(err)
		return nil, errors.Wrap(customErr, err.Error())
	}
	return vacancies, nil
}

func (r *Repository) SuggestVacancyTitle(
	suggestWord string,
	ctx context.Context) ([]models.SuggestVacancyTittle, error) {
	var suggestTittles []models.SuggestVacancyTittle
	if suggestWord == "" {
		if err := r.db.Select(&suggestTittles, selectAllTittle); err != nil {
			customErr := errortools.SqlErrorChoice(err)
			return nil, errors.Wrap(customErr, err.Error())
		}
		return suggestTittles, nil
	}
	suggestWord += "%"
	if err := r.db.Select(&suggestTittles, selectTittle, suggestWord); err != nil {
		customErr := errortools.SqlErrorChoice(err)
		return nil, errors.Wrap(customErr, err.Error())
	}
	return suggestTittles, nil
}
