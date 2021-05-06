package usecase

import (
	"authorizationservice/internal/app/models"
	"authorizationservice/internal/app/session/mock"
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"testing"
)

//Проверка созданияю отзыва
func TestCreateSession(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	session := models.Session{
		UserId:   1,
		Executor: true,
	}
	mockSessionRepo := mock.NewMockRepository(ctrl)
	mockSessionRepo.EXPECT().Store(session, ctx).Times(1).Return(nil)
	mockTools := mock.NewMockSessionTools(ctrl)
	mockTools.EXPECT().BeforeCreate(session).Times(1).Return(session, nil)
	useCase := UseCase{
		sessionRepository: mockSessionRepo,
		tools:             mockTools,
	}

	_, err := useCase.Create(session.UserId, session.Executor, ctx)
	require.NoError(t, err)
}

func TestCreateSessionErr(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	session := models.Session{
		UserId:   1,
		Executor: true,
	}
	mockSessionRepo := mock.NewMockRepository(ctrl)
	mockSessionRepo.EXPECT().Store(session, ctx).Times(1).Return(errors.New("err"))
	mockTools := mock.NewMockSessionTools(ctrl)
	mockTools.EXPECT().BeforeCreate(session).Times(1).Return(session, nil)
	useCase := UseCase{
		sessionRepository: mockSessionRepo,
		tools:             mockTools,
	}

	_, err := useCase.Create(session.UserId, session.Executor, ctx)
	require.Error(t, err)
}

func TestCreateSessionErr2(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	session := models.Session{
		UserId:   1,
		Executor: true,
	}
	mockSessionRepo := mock.NewMockRepository(ctrl)
	mockTools := mock.NewMockSessionTools(ctrl)
	mockTools.EXPECT().BeforeCreate(session).Times(1).Return(session, errors.New("err"))
	useCase := UseCase{
		sessionRepository: mockSessionRepo,
		tools:             mockTools,
	}

	_, err := useCase.Create(session.UserId, session.Executor, ctx)
	require.Error(t, err)
}

func TestGetSession(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	session := models.Session{
		SessionID: "sadasdsdsa",
		UserId:    1,
		Executor:  true,
	}
	mockSessionRepo := mock.NewMockRepository(ctrl)
	mockSessionRepo.EXPECT().Get(session.SessionID, ctx)
	mockTools := mock.NewMockSessionTools(ctrl)
	useCase := UseCase{
		sessionRepository: mockSessionRepo,
		tools:             mockTools,
	}

	_, err := useCase.Get(session.SessionID, ctx)
	require.NoError(t, err)
}

func TestRemoveSession(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	session := models.Session{
		SessionID: "sadasdsdsa",
		UserId:    1,
		Executor:  true,
	}
	mockSessionRepo := mock.NewMockRepository(ctrl)
	mockSessionRepo.EXPECT().Remove(session.SessionID, ctx)
	mockTools := mock.NewMockSessionTools(ctrl)
	useCase := UseCase{
		sessionRepository: mockSessionRepo,
		tools:             mockTools,
	}

	err := useCase.Remove(session.SessionID, ctx)
	require.NoError(t, err)
}

func TestNewSession(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSessionRepo := mock.NewMockRepository(ctrl)

	uc := New(mockSessionRepo)

	require.Equal(t, uc.sessionRepository, mockSessionRepo)
}
