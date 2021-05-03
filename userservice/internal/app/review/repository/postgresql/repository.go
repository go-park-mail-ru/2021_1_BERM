package postgresql

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

func (r *Repository) Create(review models.Review, ctx context.Context) (*models.Review, error) {
	var ID uint64
	err := r.Db.QueryRow(
		CreateReviewsRequest,
		review.UserId,
		review.ToUserId,
		review.OrderId,
		review.Description,
		review.Score).Scan(&ID)
	if err != nil {
		customErr := errortools.SqlErrorChoice(err)
		return nil, errors.Wrap(customErr, err.Error())
	}
	review.ID = ID
	return &review, nil
}

func (r *Repository) GetAll(userId uint64, ctx context.Context) ([]models.Review, error) {
	var reviews []models.Review
	err := r.Db.Select(&reviews, SelectAllReviewsByUseIDRequest, userId)
	if err != nil {
		customErr := errortools.SqlErrorChoice(err)
		return nil, errors.Wrap(customErr, err.Error())
	}
	return reviews, nil
}

func (r *Repository) GetAvgScoreByUserId(userId uint64, ctx context.Context) (uint8, error) {
	var rating uint8
	err := r.Db.Get(&rating, SelectAvgScore, userId)
	if errors.Unwrap(err).Error() == "converting NULL to uint8 is unsupported" {
		return 0, nil
	}
	if err != nil {
		customErr := errortools.SqlErrorChoice(err)
		return 0, errors.Wrap(customErr, err.Error())
	}
	return rating, nil
}
