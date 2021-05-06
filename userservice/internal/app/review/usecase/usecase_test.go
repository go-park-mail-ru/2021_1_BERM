package usecase

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"testing"
	"user/internal/app/models"
	orderMock "user/internal/app/order/mock"
	reviewMock "user/internal/app/review/mock"
	userMock "user/internal/app/user/mock"
)

//Проверка созданияю отзыва
func TestCreateReview(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	review := &models.Review{
		ID:              1,
		UserId:          1,
		ToUserId:        2,
		OrderId:         1,
		Description:     "Збс делает",
		Score:           4,
		OrderName:       "Сделай что то",
		UserLogin:       "Lala@mail.ru",
		UserNameSurname: "Name surname",
	}

	ctx := context.Background()
	mockReviewRepo := reviewMock.NewMockRepository(ctrl)
	mockReviewRepo.EXPECT().Create(*review, ctx).Times(1).Return(review, nil)

	useCase := UseCase{
		reviewRepository: mockReviewRepo,
	}
	review, err := useCase.Create(*review, ctx)
	require.NoError(t, err)
}

func TestGetAllReviewsByUserID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	review := &models.Review{
		ID:              1,
		UserId:          1,
		ToUserId:        2,
		OrderId:         1,
		Description:     "Збс делает",
		Score:           4,
		OrderName:       "Сделай что то",
		UserLogin:       "Lala@mail.ru",
		UserNameSurname: "Name surname",
	}
	reviews := []models.Review{*review}
	ctx := context.Background()
	mockReviewRepo := reviewMock.NewMockRepository(ctrl)
	mockReviewRepo.EXPECT().GetAll(uint64(2), ctx).Times(1).Return(reviews, nil)

	mockUserRepository := userMock.NewMockRepository(ctrl)
	mockUserRepository.EXPECT().FindUserByID(review.UserId, ctx).Times(1).Return(&models.UserInfo{}, nil)
	mockUserRepository.EXPECT().FindUserByID(review.ToUserId, ctx).Times(1).Return(&models.UserInfo{}, nil)

	mockOrderRepository := orderMock.NewMockRepository(ctrl)
	mockOrderRepository.EXPECT().GetByID(review.OrderId, ctx).Times(1).Return(&models.OrderInfo{}, nil)
	useCase := UseCase{
		reviewRepository: mockReviewRepo,
		userRepository:   mockUserRepository,
		orderRepository:  mockOrderRepository,
	}
	userRewiews, err := useCase.GetAllReviewByUserId(2, ctx)
	require.NoError(t, err)
	require.Equal(t, userRewiews.Reviews, reviews)
}

func TestNewUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockReviewRepo := reviewMock.NewMockRepository(ctrl)
	mockUserRepo := userMock.NewMockRepository(ctrl)
	mockOrderRepo := orderMock.NewMockRepository(ctrl)
	r := New(mockReviewRepo, mockUserRepo, mockOrderRepo)
	require.Equal(t, r.reviewRepository, mockReviewRepo)
	require.Equal(t, r.userRepository, mockUserRepo)
	require.Equal(t, r.orderRepository, mockOrderRepo)
}
