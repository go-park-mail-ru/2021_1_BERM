package usecase_test

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"testing"
	"user/internal/app/models"
	"user/internal/app/session/mock"
	sessionUseCase "user/internal/app/session/usecase"
)

//Проверка сессии
func TestCreateSpecialize(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	mockSessionRepository := mock.NewMockRepository(ctrl)
	sessionID := "wsdadkjSAHDBASDNjl"
	mockSessionRepository.EXPECT().Check(sessionID, ctx).Times(1).Return(&models.UserBasicInfo{ID: 1,
		Executor: true}, nil)

	useCase := sessionUseCase.UseCase{
		SessionRepository: mockSessionRepository,
	}
	u, err := useCase.Check(sessionID, ctx)
	require.NoError(t, err)
	require.Equal(t, u.ID, uint64(1))
	require.Equal(t, u.Executor, true)
}
