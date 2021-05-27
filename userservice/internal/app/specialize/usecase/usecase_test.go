package usecase_test

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/lib/pq"
	"github.com/stretchr/testify/require"
	"testing"
	"user/internal/app/specialize/mock"
	useCaseSpecialize "user/internal/app/specialize/usecase"
	customError "user/pkg/error"
)

//Проверка создания специализации
func TestCreateSpecialize(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	mockSpecializeMock := mock.NewMockRepository(ctrl)
	spec := "spec"
	mockSpecializeMock.EXPECT().Create(spec, ctx).Times(1).Return(uint64(1), nil)

	useCase := useCaseSpecialize.UseCase{
		SpecializeRepository: mockSpecializeMock,
	}
	ID, err := useCase.Create(spec, ctx)
	require.NoError(t, err)
	require.Equal(t, ID, uint64(1))
}

//Проверка удаления ассоциации специализации с юзером специализации
func TestRemoveSpecialize(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	mockSpecializeMock := mock.NewMockRepository(ctrl)
	spec := "spec"
	mockSpecializeMock.EXPECT().FindByName(spec, ctx).Times(1).Return(uint64(1), nil)
	mockSpecializeMock.EXPECT().RemoveAssociateSpecializationWithUser(uint64(1), uint64(1), ctx).Times(1).Return(nil)
	useCase := useCaseSpecialize.UseCase{
		SpecializeRepository: mockSpecializeMock,
	}
	err := useCase.Remove(1, spec, ctx)
	require.NoError(t, err)
}

//Проверка привязки специализации к юзеру
func TestAssociateSpecWithUsser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	mockSpecializeMock := mock.NewMockRepository(ctrl)
	spec := "spec"
	mockSpecializeMock.EXPECT().FindByName(spec, ctx).Times(1).Return(uint64(0), customError.ErrorDuplicate)
	mockSpecializeMock.EXPECT().Create(spec, ctx).Times(1).Return(uint64(1), nil)
	mockSpecializeMock.EXPECT().AssociateSpecializationWithUser(uint64(1), uint64(1), ctx).Times(1).Return(nil)
	useCase := useCaseSpecialize.UseCase{
		SpecializeRepository: mockSpecializeMock,
	}
	err := useCase.AssociateWithUser(1, spec, ctx)
	require.NoError(t, err)
}

//Проверка привязки существующей специализации к юзеру
func TestAssociateSpecWithUsserWithDuplicate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	mockSpecializeMock := mock.NewMockRepository(ctrl)
	spec := "spec"
	mockSpecializeMock.EXPECT().FindByName(spec, ctx).Times(1).Return(uint64(1), nil)
	mockSpecializeMock.EXPECT().AssociateSpecializationWithUser(uint64(1), uint64(1), ctx).Times(1).Return(&pq.Error{
		Code: "23505",
	})
	useCase := useCaseSpecialize.UseCase{
		SpecializeRepository: mockSpecializeMock,
	}
	err := useCase.AssociateWithUser(1, spec, ctx)
	require.Error(t, err)
}

func TestNewUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSpecializeRepo := mock.NewMockRepository(ctrl)
	s := useCaseSpecialize.New(mockSpecializeRepo)
	require.Equal(t, s.SpecializeRepository, mockSpecializeRepo)
}
