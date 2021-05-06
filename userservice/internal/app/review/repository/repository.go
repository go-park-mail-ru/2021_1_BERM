package repository

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

func (r *Repository) GetAvgScoreByUserId(userId uint64, ctx context.Context) (*models.UserReviewInfo, error) {
	userReviewInfo := &models.UserReviewInfo{}
	err := r.Db.Get(userReviewInfo, SelectAvgScore, userId)
	if err != nil {
		if errors.Unwrap(err).Error() == "converting NULL to uint8 is unsupported" {
			return userReviewInfo, nil
		}
		customErr := errortools.SqlErrorChoice(err)
		return nil, errors.Wrap(customErr, err.Error())
	}
	return userReviewInfo, nil
}
