package postgresql

import (
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"strconv"
	"user/Error"
	"user/internal/app/models"
	"user/pkg/database/postgresql"
)

type Repository struct {
	Db *sqlx.DB
}

func (r *Repository) FindByUserID(userID uint64, ctx context.Context) (pq.StringArray, error) {
	rows := r.Db.QueryRow(SelectSpecializesByUserID, userID)
	var specializes pq.StringArray
	if err := rows.Scan(&specializes); err != nil {
		return nil, &Error.Error{
			Err:              err,
			InternalError:    true,
			ErrorDescription: Error.InternalServerErrorDescription,
		}
	}
	return specializes, nil
}

func (r *Repository) Create(specialize string, ctx context.Context) (uint64, error) {
	var ID uint64 = 0
	err := r.Db.QueryRow(
		CreateSpecializeRequest, specialize).Scan(&ID)
	if err != nil {
		return 0, postgresql.WrapPostgreError(err)
	}
	return ID, nil
}

func (r *Repository) FindById(ID uint64, ctx context.Context) (string, error) {
	spec := models.Specialize{}
	err := r.Db.Get(&spec, SelectSpecializesByID, ID)
	if err != nil {
		return "", postgresql.WrapPostgreError(err)
	}
	return spec.Name, nil
}

func (r *Repository) FindByName(spec string, ctx context.Context) (uint64, error) {
	specialize := models.Specialize{}
	err := r.Db.Get(&specialize, SelectSpecializesByName, spec)
	if err != nil {
		return 0, postgresql.WrapPostgreError(err)
	}
	return specialize.ID, nil
}

func (r *Repository) AssociateSpecializationWithUser(specId uint64, userId uint64, ctx context.Context) error {
	_, err := r.Db.NamedExec(
		CreateUserSpecializeRequest,
		map[string]interface{}{
			"userID": strconv.FormatUint(userId, 10),
			"specID": strconv.FormatUint(specId, 10),
		})
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) Remove(ID uint64, ctx context.Context) error {
	err := r.Db.QueryRow(DeleteSpecialize, ID).Err()
	if err != nil {
		return postgresql.WrapPostgreError(err)
	}
	return nil
}

func (r *Repository)RemoveAssociateSpecializationWithUser(specId uint64, userId uint64, ctx context.Context) error{
	err := r.Db.QueryRow(DeleteSpecialize, userId, specId).Err()
	if err != nil {
		return err
	}
	return nil
}
