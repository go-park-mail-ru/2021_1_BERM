package handlers

import (
	"authorizationservice/api"
	"authorizationservice/internal/app/models"
	"authorizationservice/internal/app/session/mock"
	"authorizationservice/pkg/metric"
	"context"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGRPCServer_Check(t *testing.T) {
	metric.New()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()


	sessionInfo := models.Session{
		SessionID: "kek",
		UserId: 1,
		Executor: true,
	}
	expectResponse := &api.SessionCheckResponse{
		ID: 1,
		Executor: true,
	}
	mockUseCase := mock.NewMockUseCase(ctrl)

	handle := NewGRPCServer(mockUseCase)

	mockUseCase.EXPECT().
		Get(sessionInfo.SessionID, context.Background()).
		Times(1).
		Return(&sessionInfo, nil)

	response, err := handle.Check(context.Background(), &api.SessionCheckRequest{SessionId: "kek"})
	require.NoError(t, err)
	require.Equal(t, expectResponse, response)
	metric.Destroy()
}

func TestGRPCServer_CheckErr(t *testing.T) {
	metric.New()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()


	sessionInfo := models.Session{
		SessionID: "kek",
		UserId: 1,
		Executor: true,
	}

	mockUseCase := mock.NewMockUseCase(ctrl)

	handle := NewGRPCServer(mockUseCase)

	mockUseCase.EXPECT().
		Get(sessionInfo.SessionID, context.Background()).
		Times(1).
		Return(&sessionInfo, errors.New("Err"))

	_, err := handle.Check(context.Background(), &api.SessionCheckRequest{SessionId: "kek"})
	require.Error(t, err)
	metric.Destroy()
}
