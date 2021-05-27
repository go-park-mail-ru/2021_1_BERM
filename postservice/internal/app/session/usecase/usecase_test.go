package usecase_test

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"post/internal/app/models"
	"post/internal/app/session/mock"
	sessUCase "post/internal/app/session/usecase"
	"testing"
)

//Проверка сессии
func TestCreateSpecialize(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	mockSessionRepository := mock.NewMockRepository(ctrl)
	sessionID := "wsdadkjSAHDBASDNjl"
	mockSessionRepository.EXPECT().Check(ctx, sessionID).Times(1).Return(&models.UserBasicInfo{ID: 1,
		Executor: true}, nil)

	useCase := sessUCase.UseCase{
		SessionRepository: mockSessionRepository,
	}
	u, err := useCase.Check(sessionID, ctx)
	require.NoError(t, err)
	require.Equal(t, u.ID, uint64(1))
	require.Equal(t, u.Executor, true)
}
