package usecase

import (
	"context"
	"user/internal/app/models"
	"user/internal/app/order"
	"user/internal/app/review"
	"user/internal/app/user"
)

type UseCase struct {
	ReviewRepository review.Repository
	UserRepository   user.Repository
	OrderRepository  order.Repository
}

func New(reviewRepository review.Repository,
	userRepository user.Repository, orderRepository order.Repository) *UseCase {
	return &UseCase{
		ReviewRepository: reviewRepository,
		UserRepository:   userRepository,
		OrderRepository:  orderRepository,
	}
}

func (useCase *UseCase) Create(review models.Review, ctx context.Context) (*models.Review, error) {
	revResp, err := useCase.ReviewRepository.Create(review, ctx)
	if err != nil {
		return nil, err
	}
	return revResp, err
}

func (useCase *UseCase) GetAllReviewByUserId(userId uint64, ctx context.Context) (*models.UserReviews, error) {
	reviews, err := useCase.ReviewRepository.GetAll(userId, ctx)
	if err != nil {
		return nil, err
	}
	for index := range reviews {
		u, err := useCase.UserRepository.FindUserByID(reviews[index].UserId, ctx)
		if err != nil {
			return nil, err
		}
		reviews[index].UserLogin = u.Login
		reviews[index].UserNameSurname = u.NameSurname
		oInf, err := useCase.OrderRepository.GetByID(reviews[index].OrderId, ctx)
		if err != nil {
			return nil, err
		}
		reviews[index].OrderName = oInf.OrderName
	}
	u, err := useCase.UserRepository.FindUserByID(userId, ctx)
	if err != nil {
		return nil, err
	}
	if reviews == nil {
		return &models.UserReviews{
			Name:    u.NameSurname,
			Login:   u.Login,
			Reviews: []models.Review{},
		}, nil
	}
	return &models.UserReviews{
		Name:    u.NameSurname,
		Login:   u.Login,
		Reviews: reviews,
	}, nil
}
