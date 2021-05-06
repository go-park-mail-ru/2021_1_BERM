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