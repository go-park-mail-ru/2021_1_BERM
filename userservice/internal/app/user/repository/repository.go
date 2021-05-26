package user

import (
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"user/internal/app/models"
	"user/pkg/error/errortools"
)

type Repository struct {
	Db *sqlx.DB
}

func NewRepo(db *sqlx.DB) *Repository {
	return &Repository{
		db,
	}
}

func (r *Repository) Create(user models.NewUser, ctx context.Context) (uint64, error) {
	var ID uint64
	err := r.Db.QueryRow(
		CreateUserRequest,
		user.Email,
		user.EncryptPassword,
		user.Login,
		user.NameSurname,
		user.About,
		user.Executor).Scan(&ID)
	if err != nil {
		customErr := errortools.SqlErrorChoice(err)
		return 0, errors.Wrap(customErr, err.Error())
	}
	return ID, nil
}

func (r *Repository) Change(user models.ChangeUser, ctx context.Context) error {
	tx := r.Db.MustBegin()
	_, err := tx.NamedExec(UpdateUser, user)
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

func (r *Repository) FindUserByID(ID uint64, ctx context.Context) (*models.UserInfo, error) {
	user := models.UserInfo{}
	err := r.Db.Get(&user, SelectUserByID, ID)
	if err != nil {
		customErr := errortools.SqlErrorChoice(err)
		return nil, errors.Wrap(customErr, err.Error())
	}
	return &user, nil
}

func (r *Repository) FindUserByEmail(email string, ctx context.Context) (*models.UserInfo, error) {
	user := models.UserInfo{}
	err := r.Db.Get(&user, SelectUserByEmail, email)
	if err != nil {
		customErr := errortools.SqlErrorChoice(err)
		return nil, errors.Wrap(customErr, err.Error())
	}
	return &user, nil
}

func (r *Repository) SetUserImg(ID uint64, img string, ctx context.Context) error {
	err := r.Db.QueryRow(UpdateUserImg, img, ID).Err()
	if err != nil {
		customErr := errortools.SqlErrorChoice(err)
		return errors.Wrap(customErr, err.Error())
	}
	return err
}
const(
	getUsersRating = `SELECT users.id, email, password, login, name_surname, about, executor, img, AVG(score) AS rating
		FROM userservice.users AS users
		INNER JOIN userservice.reviews
		 ON users.id = reviews.to_user_id
		WHERE CASE WHEN $1 != 0 THEN (SELECT AVG(score) FROM userservice.reviews WHERE to_user_id = users.id) >= $1 ELSE true END
		AND CASE WHEN $2 != 0 THEN (SELECT AVG(score) FROM userservice.reviews WHERE to_user_id = users.id) <= $2 ELSE true END
		AND CASE WHEN $3 != '~' THEN to_tsvector(name_surname) @@ to_tsquery($3) ELSE true END
		GROUP BY users.id
		ORDER BY rating LIMIT $4 OFFSET $5`
	getUsersRatingDesc = `SELECT users.id, email, password, login, name_surname, about, executor, img, AVG(score) AS rating
		FROM userservice.users AS users
		INNER JOIN userservice.reviews
		 ON users.id = reviews.to_user_id
		WHERE CASE WHEN $1 != 0 THEN (SELECT AVG(score) FROM userservice.reviews WHERE to_user_id = users.id) >= $1 ELSE true END
		AND CASE WHEN $2 != 0 THEN (SELECT AVG(score) FROM userservice.reviews WHERE to_user_id = users.id) <= $2 ELSE true END
		AND CASE WHEN $3 != '~' THEN to_tsvector(name_surname) @@ to_tsquery($3) ELSE true END
		GROUP BY users.id
		ORDER BY rating DESC LIMIT $4 OFFSET $5`

	getUsersNick = `SELECT users.id, email, password, login, name_surname, about, executor, img, AVG(score) AS rating
		FROM userservice.users AS users
		INNER JOIN userservice.reviews
		 ON users.id = reviews.to_user_id
		WHERE CASE WHEN $1 != 0 THEN (SELECT AVG(score) FROM userservice.reviews WHERE to_user_id = users.id) >= $1 ELSE true END
		AND CASE WHEN $2 != 0 THEN (SELECT AVG(score) FROM userservice.reviews WHERE to_user_id = users.id) <= $2 ELSE true END
		AND CASE WHEN $3 != '~' THEN to_tsvector(name_surname) @@ to_tsquery($3) ELSE true END
		GROUP BY users.id
		ORDER BY name_surname DESC LIMIT $4 OFFSET $5`

	getUsersNickDesc = `SELECT users.id, email, password, login, name_surname, about, executor, img, AVG(score) AS rating
		FROM userservice.users AS users
		INNER JOIN userservice.reviews
		 ON users.id = reviews.to_user_id
		WHERE CASE WHEN $1 != 0 THEN (SELECT AVG(score) FROM userservice.reviews WHERE to_user_id = users.id) >= $1 ELSE true END
		AND CASE WHEN $2 != 0 THEN (SELECT AVG(score) FROM userservice.reviews WHERE to_user_id = users.id) <= $2 ELSE true END
		AND CASE WHEN $3 != '~' THEN to_tsvector(name_surname) @@ to_tsquery($3) ELSE true END
		GROUP BY users.id
		ORDER BY name_surname DESC LIMIT $4 OFFSET $5`

	ctxParam uint8 = 4;
)
func (r *Repository)GetUsers(ctx context.Context) ([]models.UserInfo, error){
	var userInfo []models.UserInfo
	param := ctx.Value(ctxParam).(map[string]interface{})
	limit := param["limit"].(int)
	offset := param["offset"].(int)
	desc := param["desc"].(bool)
	from := param["from"].(int)
	to := param["to"].(int)
	searchStr := param["search_str"].(string)
	sort := param["sort"].(string)
	if desc {
		if sort == "rating" {
			if err := r.Db.Select(&userInfo, getUsersRatingDesc, from, to, searchStr,  limit, offset); err != nil {
				customErr := errortools.SqlErrorChoice(err)
				return nil, errors.Wrap(customErr, err.Error())
			}
		}else{
			if err := r.Db.Select(&userInfo, getUsersNickDesc, from, to, searchStr,limit, offset); err != nil {
				customErr := errortools.SqlErrorChoice(err)
				return nil, errors.Wrap(customErr, err.Error())
			}
		}
	}else{
		if sort == "rating" {
			if err := r.Db.Select(&userInfo, getUsersRating, from, to, searchStr,  limit,offset); err != nil {
				customErr := errortools.SqlErrorChoice(err)
				return nil, errors.Wrap(customErr, err.Error())
			}
		}else{
			if err := r.Db.Select(&userInfo, getUsersNick, from, to, searchStr, limit, offset); err != nil {
				customErr := errortools.SqlErrorChoice(err)
				return nil, errors.Wrap(customErr, err.Error())
			}
		}
	}
	return userInfo, nil
}
