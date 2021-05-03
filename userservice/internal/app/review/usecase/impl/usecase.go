package impl

import (
	"context"
	"user/internal/app/models"
	orderRep "user/internal/app/order/repository"
	"user/internal/app/review/repository"
	repository2 "user/internal/app/user/repository"
)

type UseCase struct {
	reviewRepository repository.Repository
	userRepository   repository2.Repository
	orderRepository  orderRep.Repository
}

func New(reviewRepository repository.Repository, userRepository repository2.Repository, orderRepository orderRep.Repository) *UseCase {
	return &UseCase{
		reviewRepository: reviewRepository,
		userRepository: userRepository,
		orderRepository: orderRepository,
	}
}

func (useCase *UseCase) Create(review models.Review, ctx context.Context) (*models.Review, error) {
	revResp, err := useCase.reviewRepository.Create(review, ctx)
	if err != nil {
		return nil, err
	}
	return revResp, err
}

func (useCase *UseCase) GetAllReviewByUserId(userId uint64, ctx context.Context) ([]models.Review, error) {
	reviews, err := useCase.reviewRepository.GetAll(userId, ctx)
	if err != nil {
		return nil, err
	}
	for index, _ := range reviews {
		u, err := useCase.userRepository.FindUserByID(reviews[index].ToUserId, ctx)
		if err != nil {
			return nil, err
		}
		reviews[index].UserLogin = u.Login
		reviews[index].UserNameSurname = u.NameSurname
		oInf, err := useCase.orderRepository.GetByID(reviews[index].ToUserId, ctx)
		if err != nil {
			return nil, err
		}
		reviews[index].OrderName = oInf.OrderName
	}
	return reviews, err
}

func (useCase *UseCase) GetAvgScoreByUserId(userId uint64, ctx context.Context) (uint8, error) {
	revResp, err := useCase.reviewRepository.GetAvgScoreByUserId(userId, ctx)
	if err != nil {
		return 0, err
	}
	return revResp, err
}
