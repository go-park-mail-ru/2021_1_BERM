package vacancy_test

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"post/api"
	"post/internal/app/models"
	"post/internal/app/vacancy/mock"
	vacUseCase "post/internal/app/vacancy/usecase"
	"post/pkg/types"
	"testing"
)

const ctxUserID types.CtxKey = 2

func TestDeleteExecutor(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	vacancy := models.Vacancy{
		ExecutorID: 0,
	}
	ctx := context.Background()
	mockVacancyRepo := mock.NewMockRepository(ctrl)
	mockVacancyRepo.EXPECT().UpdateExecutor(vacancy, ctx).Times(1).Return(nil)

	mockUserRepo := mock.NewMockUserClient(ctrl)

	useCase := vacUseCase.NewUseCase(mockVacancyRepo, mockUserRepo)

	err := useCase.DeleteExecutor(vacancy, ctx)

	require.NoError(t, err)
}

func TestDeleteExecutorErr(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	vacancy := models.Vacancy{
		ExecutorID: 0,
	}
	ctx := context.Background()
	mockVacancyRepo := mock.NewMockRepository(ctrl)
	mockVacancyRepo.EXPECT().UpdateExecutor(vacancy, ctx).Times(1).Return(errors.New("Db dead"))
	mockUserRepo := mock.NewMockUserClient(ctrl)

	useCase := vacUseCase.NewUseCase(mockVacancyRepo, mockUserRepo)

	err := useCase.DeleteExecutor(vacancy, ctx)
	require.Error(t, err)
}

func TestCreateVacancy(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	vacancy := models.Vacancy{
		VacancyName: "Keke",
		CustomerID:  1,
		Salary:      1488,
		Description: "Aue jizn voram",
		Category:    "Back",
	}
	expectVacancy := &models.Vacancy{
		VacancyName: "Keke",
		CustomerID:  1,
		Salary:      1488,
		Description: "Aue jizn voram",
		Category:    "Back",
		Img:         "kek",
		Login:       "Mem",
		ID:          1,
	}
	var id uint64
	id = 1
	ctx := context.Background()
	mockVacancyRepo := mock.NewMockRepository(ctrl)
	mockUserRepo := mock.NewMockUserClient(ctrl)
	useCase := vacUseCase.NewUseCase(mockVacancyRepo, mockUserRepo)

	mockVacancyRepo.EXPECT().
		Create(vacancy, ctx).
		Times(1).
		Return(id, nil)
	mockUserRepo.EXPECT().
		GetUserById(ctx, &api.UserRequest{Id: vacancy.CustomerID}).
		Times(1).
		Return(&api.UserInfoResponse{Login: "Mem", Img: "kek"}, nil)

	respVacancy, err := useCase.Create(vacancy, ctx)

	require.Equal(t, expectVacancy, respVacancy)
	require.NoError(t, err)

	id = 0
	mockVacancyRepo.EXPECT().
		Create(vacancy, ctx).
		Times(1).
		Return(id, errors.New("DB error"))

	_, err = useCase.Create(vacancy, ctx)
	require.Error(t, err)

	id = 1
	mockVacancyRepo.EXPECT().
		Create(vacancy, ctx).
		Times(1).
		Return(id, nil)
	mockUserRepo.EXPECT().
		GetUserById(ctx, &api.UserRequest{Id: vacancy.CustomerID}).
		Times(1).
		Return(&api.UserInfoResponse{Login: "Mem", Img: "kek"}, errors.New("GRPC Err"))

	_, err = useCase.Create(vacancy, ctx)
	require.Error(t, err)
}

func TestFindByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	vacancy := &models.Vacancy{
		ID:          1,
		VacancyName: "Keke",
		CustomerID:  1,
		Salary:      1488,
		Description: "Aue jizn voram",
		Category:    "Back",
	}
	expectVacancy := &models.Vacancy{
		ID:          1,
		VacancyName: "Keke",
		CustomerID:  1,
		Salary:      1488,
		Description: "Aue jizn voram",
		Category:    "Back",
		Img:         "kek",
		Login:       "Mem",
	}

	id := uint64(1)
	ctx := context.Background()
	mockVacancyRepo := mock.NewMockRepository(ctrl)
	mockUserRepo := mock.NewMockUserClient(ctrl)
	useCase := vacUseCase.NewUseCase(mockVacancyRepo, mockUserRepo)

	mockVacancyRepo.EXPECT().
		FindByID(id, ctx).
		Times(1).
		Return(vacancy, nil)
	mockUserRepo.EXPECT().
		GetUserById(ctx, &api.UserRequest{Id: vacancy.CustomerID}).
		Times(1).
		Return(&api.UserInfoResponse{Login: "Mem", Img: "kek"}, nil)

	respVacancy, err := useCase.FindByID(id, ctx)

	require.Equal(t, respVacancy, expectVacancy)
	require.NoError(t, err)

	mockVacancyRepo.EXPECT().
		FindByID(id, ctx).
		Times(1).
		Return(vacancy, errors.New("DB err"))

	_, err = useCase.FindByID(id, ctx)

	require.Error(t, err)

	mockVacancyRepo.EXPECT().
		FindByID(id, ctx).
		Times(1).
		Return(vacancy, nil)
	mockUserRepo.EXPECT().
		GetUserById(ctx, &api.UserRequest{Id: vacancy.CustomerID}).
		Times(1).
		Return(&api.UserInfoResponse{Login: "Mem", Img: "kek"}, errors.New("GRPC err"))

	_, err = useCase.FindByID(id, ctx)

	require.Error(t, err)
}

func TestFindByUserID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	id := uint64(1)
	ctx := context.Background()
	mockVacancyRepo := mock.NewMockRepository(ctrl)
	mockUserRepo := mock.NewMockUserClient(ctrl)
	useCase := vacUseCase.NewUseCase(mockVacancyRepo, mockUserRepo)

	vacancies := []models.Vacancy{
		{
			ID:          1,
			VacancyName: "Keke",
			CustomerID:  1,
			Salary:      1488,
			Description: "Aue jizn voram",
			Category:    "Back",
		},
	}

	expectVacancies := []models.Vacancy{
		{
			ID:          1,
			VacancyName: "Keke",
			CustomerID:  1,
			Salary:      1488,
			Description: "Aue jizn voram",
			Category:    "Back",
			Login:       "Mem",
			Img:         "kek",
		},
	}

	mockUserRepo.EXPECT().
		GetUserById(ctx, &api.UserRequest{Id: id}).
		Times(2).
		Return(&api.UserInfoResponse{Login: "Mem", Img: "kek", Executor: true}, nil)
	mockVacancyRepo.EXPECT().
		FindByExecutorID(id, ctx).
		Times(1).
		Return(vacancies, nil)

	resVac, err := useCase.FindByUserID(id, ctx)

	require.Equal(t, resVac, expectVacancies)
	require.NoError(t, err)

	mockUserRepo.EXPECT().
		GetUserById(ctx, &api.UserRequest{Id: id}).
		Times(1).
		Return(&api.UserInfoResponse{Login: "Mem", Img: "kek", Executor: true}, errors.New("GRPC Err"))

	_, err = useCase.FindByUserID(id, ctx)

	require.Error(t, err)

	mockUserRepo.EXPECT().
		GetUserById(ctx, &api.UserRequest{Id: id}).
		Times(2).
		Return(&api.UserInfoResponse{Login: "Mem", Img: "kek", Executor: false}, nil)
	mockVacancyRepo.EXPECT().
		FindByCustomerID(id, ctx).
		Times(1).
		Return(vacancies, nil)

	resVac, err = useCase.FindByUserID(id, ctx)

	require.Equal(t, resVac, expectVacancies)
	require.NoError(t, err)

	expectEmptyVacancies := []models.Vacancy{}
	mockUserRepo.EXPECT().
		GetUserById(ctx, &api.UserRequest{Id: id}).
		Times(1).
		Return(&api.UserInfoResponse{Login: "Mem", Img: "kek", Executor: false}, nil)
	mockVacancyRepo.EXPECT().
		FindByCustomerID(id, ctx).
		Times(1).
		Return(nil, nil)

	resVac, err = useCase.FindByUserID(id, ctx)

	require.Equal(t, resVac, expectEmptyVacancies)
	require.NoError(t, err)

	mockUserRepo.EXPECT().
		GetUserById(ctx, &api.UserRequest{Id: id}).
		Times(1).
		Return(&api.UserInfoResponse{Login: "Mem", Img: "kek", Executor: false}, nil)
	mockVacancyRepo.EXPECT().
		FindByCustomerID(id, ctx).
		Times(1).
		Return(nil, errors.New("DB err"))

	_, err = useCase.FindByUserID(id, ctx)

	require.Error(t, err)
}

func TestChangeVacancy(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	vacancy := &models.Vacancy{
		ID:          1,
		VacancyName: "Keke",
		CustomerID:  1,
		Salary:      1488,
		Description: "Aue jizn voram",
		Category:    "Back",
	}
	oldVacancy := &models.Vacancy{
		ID:          1,
		VacancyName: "Keke",
		CustomerID:  1,
		Salary:      1488,
		Description: "Aue jizn voram",
		Category:    "Back",
	}
	expectVacancy := &models.Vacancy{
		VacancyName: "Keke",
		CustomerID:  1,
		Salary:      1488,
		Description: "Aue jizn voram",
		Category:    "Back",
		Img:         "kek",
		Login:       "Mem",
		ID:          1,
	}
	id := uint64(1)
	ctx := context.Background()
	mockVacancyRepo := mock.NewMockRepository(ctrl)
	mockUserRepo := mock.NewMockUserClient(ctrl)
	useCase := vacUseCase.NewUseCase(mockVacancyRepo, mockUserRepo)

	mockVacancyRepo.EXPECT().
		FindByID(id, ctx).
		Times(1).
		Return(oldVacancy, nil)
	mockVacancyRepo.EXPECT().
		Change(*vacancy, ctx).
		Times(1).
		Return(nil)
	mockUserRepo.EXPECT().
		GetUserById(ctx, &api.UserRequest{Id: vacancy.CustomerID}).
		Times(1).
		Return(&api.UserInfoResponse{Login: "Mem", Img: "kek"}, nil)

	resVacancy, err := useCase.ChangeVacancy(*vacancy, ctx)

	require.Equal(t, resVacancy, *expectVacancy)
	require.NoError(t, err)

	vacancyWithoutFields := &models.Vacancy{
		ID:         1,
		CustomerID: 1,
	}
	mockVacancyRepo.EXPECT().
		FindByID(id, ctx).
		Times(1).
		Return(oldVacancy, nil)
	mockVacancyRepo.EXPECT().
		Change(*vacancy, ctx).
		Times(1).
		Return(nil)
	mockUserRepo.EXPECT().
		GetUserById(ctx, &api.UserRequest{Id: vacancy.CustomerID}).
		Times(1).
		Return(&api.UserInfoResponse{Login: "Mem", Img: "kek"}, nil)

	resVacancy, err = useCase.ChangeVacancy(*vacancyWithoutFields, ctx)

	require.Equal(t, resVacancy, *expectVacancy)
	require.NoError(t, err)

	mockVacancyRepo.EXPECT().
		FindByID(id, ctx).
		Times(1).
		Return(oldVacancy, errors.New("DB err"))

	_, err = useCase.ChangeVacancy(*vacancyWithoutFields, ctx)

	require.Error(t, err)

	mockVacancyRepo.EXPECT().
		FindByID(id, ctx).
		Times(1).
		Return(oldVacancy, nil)
	mockVacancyRepo.EXPECT().
		Change(*vacancy, ctx).
		Times(1).
		Return(errors.New("DB err"))

	_, err = useCase.ChangeVacancy(*vacancy, ctx)
	require.Error(t, err)

	mockVacancyRepo.EXPECT().
		FindByID(id, ctx).
		Times(1).
		Return(oldVacancy, nil)
	mockVacancyRepo.EXPECT().
		Change(*vacancy, ctx).
		Times(1).
		Return(nil)
	mockUserRepo.EXPECT().
		GetUserById(ctx, &api.UserRequest{Id: vacancy.CustomerID}).
		Times(1).
		Return(&api.UserInfoResponse{Login: "Mem", Img: "kek"}, errors.New("GRPC err"))

	_, err = useCase.ChangeVacancy(*vacancy, ctx)

	require.Error(t, err)
}

func TestDeleteVacancy(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	var id = uint64(1)
	ctx := context.Background()
	mockVacancyRepo := mock.NewMockRepository(ctrl)
	mockUserRepo := mock.NewMockUserClient(ctrl)
	useCase := vacUseCase.NewUseCase(mockVacancyRepo, mockUserRepo)

	mockVacancyRepo.EXPECT().
		DeleteVacancy(id, ctx).
		Times(1).
		Return(nil)

	err := useCase.DeleteVacancy(id, ctx)

	require.NoError(t, err)

	mockVacancyRepo.EXPECT().
		DeleteVacancy(id, ctx).
		Times(1).
		Return(errors.New("DB err"))

	err = useCase.DeleteVacancy(id, ctx)

	require.Error(t, err)
}

func TestGetActualVacancies(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	id := uint64(1)
	mockVacancyRepo := mock.NewMockRepository(ctrl)
	mockUserRepo := mock.NewMockUserClient(ctrl)
	useCase := vacUseCase.NewUseCase(mockVacancyRepo, mockUserRepo)

	vacancies := []models.Vacancy{
		{
			ID:          1,
			VacancyName: "Keke",
			CustomerID:  1,
			Salary:      1488,
			Description: "Aue jizn voram",
			Category:    "Back",
		},
	}

	expectVacancies := []models.Vacancy{
		{
			ID:          1,
			VacancyName: "Keke",
			CustomerID:  1,
			Salary:      1488,
			Description: "Aue jizn voram",
			Category:    "Back",
			Login:       "Mem",
			Img:         "kek",
		},
	}

	ctx := context.WithValue(context.Background(), ctxUserID, uint64(1))
	mockUserRepo.EXPECT().
		GetUserById(context.Background(), &api.UserRequest{Id: id}).
		Times(1).
		Return(&api.UserInfoResponse{Login: "Mem", Img: "kek", Executor: true}, nil)
	mockUserRepo.EXPECT().
		GetUserById(ctx, &api.UserRequest{Id: id}).
		Times(1).
		Return(&api.UserInfoResponse{Login: "Mem", Img: "kek", Executor: true}, nil)
	mockVacancyRepo.EXPECT().
		GetActualVacancies(ctx).
		Times(1).
		Return(vacancies, nil)
	mockVacancyRepo.EXPECT().
		GetVacancyNum(ctx).
		Times(1).
		Return(uint64(1), nil)

	respVacancies, _, err := useCase.GetActualVacancies(ctx)

	require.Equal(t, respVacancies, expectVacancies)
	require.NoError(t, err)

	mockVacancyRepo.EXPECT().
		GetActualVacancies(ctx).
		Times(1).
		Return(vacancies, errors.New("DB err"))

	_, _, err = useCase.GetActualVacancies(ctx)

	require.Error(t, err)

	mockUserRepo.EXPECT().
		GetUserById(context.Background(), &api.UserRequest{Id: id}).
		Times(1).
		Return(&api.UserInfoResponse{Login: "Mem", Img: "kek", Executor: true}, errors.New("GRPC error"))
	mockVacancyRepo.EXPECT().
		GetActualVacancies(ctx).
		Times(1).
		Return(vacancies, nil)

	_, _, err = useCase.GetActualVacancies(ctx)

	require.Error(t, err)

	emptyVacancies := []models.Vacancy{}

	mockVacancyRepo.EXPECT().
		GetActualVacancies(ctx).
		Times(1).
		Return(nil, nil)

	respVacancies, _, err = useCase.GetActualVacancies(ctx)

	require.Equal(t, respVacancies, emptyVacancies)
	require.NoError(t, err)
}

func TestSelectExecutor(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()
	mockVacancyRepo := mock.NewMockRepository(ctrl)
	mockUserRepo := mock.NewMockUserClient(ctrl)
	useCase := vacUseCase.NewUseCase(mockVacancyRepo, mockUserRepo)

	vacancy := &models.Vacancy{
		ID:          1,
		VacancyName: "Keke",
		CustomerID:  1,
		Salary:      1488,
		Description: "Aue jizn voram",
		Category:    "Back",
		ExecutorID:  2,
	}

	mockVacancyRepo.EXPECT().
		UpdateExecutor(*vacancy, ctx).
		Times(1).
		Return(nil)
	mockUserRepo.EXPECT().
		GetUserById(ctx, &api.UserRequest{Id: vacancy.ExecutorID}).
		Times(1).
		Return(&api.UserInfoResponse{Login: "Mem", Img: "kek", Executor: true}, nil)

	err := useCase.SelectExecutor(*vacancy, ctx)

	require.NoError(t, err)

	errVacancy := &models.Vacancy{
		ID:          1,
		VacancyName: "Keke",
		CustomerID:  1,
		Salary:      1488,
		Description: "Aue jizn voram",
		Category:    "Back",
		ExecutorID:  1,
	}

	mockUserRepo.EXPECT().
		GetUserById(ctx, &api.UserRequest{Id: errVacancy.ExecutorID}).
		Times(1).
		Return(&api.UserInfoResponse{Login: "Mem", Img: "kek", Executor: true}, nil)

	err = useCase.SelectExecutor(*errVacancy, ctx)

	require.Error(t, err)

	mockUserRepo.EXPECT().
		GetUserById(ctx, &api.UserRequest{Id: errVacancy.ExecutorID}).
		Times(1).
		Return(&api.UserInfoResponse{Login: "Mem", Img: "kek", Executor: false}, nil)

	err = useCase.SelectExecutor(*errVacancy, ctx)

	require.Error(t, err)

	mockUserRepo.EXPECT().
		GetUserById(ctx, &api.UserRequest{Id: vacancy.ExecutorID}).
		Times(1).
		Return(&api.UserInfoResponse{Login: "Mem", Img: "kek", Executor: true}, errors.New("GRPC err"))

	err = useCase.SelectExecutor(*vacancy, ctx)

	require.Error(t, err)

	mockVacancyRepo.EXPECT().
		UpdateExecutor(*vacancy, ctx).
		Times(1).
		Return(errors.New("DB err"))
	mockUserRepo.EXPECT().
		GetUserById(ctx, &api.UserRequest{Id: vacancy.ExecutorID}).
		Times(1).
		Return(&api.UserInfoResponse{Login: "Mem", Img: "kek", Executor: true}, nil)

	err = useCase.SelectExecutor(*vacancy, ctx)

	require.Error(t, err)
}

func TestCloseVacancy(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()
	mockVacancyRepo := mock.NewMockRepository(ctrl)
	mockUserRepo := mock.NewMockUserClient(ctrl)
	useCase := vacUseCase.NewUseCase(mockVacancyRepo, mockUserRepo)

	id := uint64(1)

	vacancy := &models.Vacancy{
		ID:          1,
		VacancyName: "Keke",
		CustomerID:  1,
		Salary:      1488,
		Description: "Aue jizn voram",
		Category:    "Back",
		ExecutorID:  2,
	}

	mockVacancyRepo.EXPECT().
		FindByID(id, ctx).
		Times(1).
		Return(vacancy, nil)
	mockVacancyRepo.EXPECT().
		DeleteVacancy(id, ctx).
		Times(1).
		Return(nil)
	mockVacancyRepo.EXPECT().
		CreateArchive(*vacancy, ctx).
		Times(1).
		Return(id, nil)

	err := useCase.CloseVacancy(id, ctx)

	require.NoError(t, err)

	mockVacancyRepo.EXPECT().
		FindByID(id, ctx).
		Times(1).
		Return(vacancy, nil)
	mockVacancyRepo.EXPECT().
		DeleteVacancy(id, ctx).
		Times(1).
		Return(nil)
	mockVacancyRepo.EXPECT().
		CreateArchive(*vacancy, ctx).
		Times(1).
		Return(id, errors.New("DB err"))

	err = useCase.CloseVacancy(id, ctx)

	require.Error(t, err)

	mockVacancyRepo.EXPECT().
		FindByID(id, ctx).
		Times(1).
		Return(vacancy, nil)
	mockVacancyRepo.EXPECT().
		DeleteVacancy(id, ctx).
		Times(1).
		Return(errors.New("DB err"))

	err = useCase.CloseVacancy(id, ctx)

	require.Error(t, err)

	mockVacancyRepo.EXPECT().
		FindByID(id, ctx).
		Times(1).
		Return(vacancy, errors.New("DB err"))

	err = useCase.CloseVacancy(id, ctx)

	require.Error(t, err)
}

//func TestGetArchiveVacancies(t *testing.T) {
//	ctrl := gomock.NewController(t)
//	defer ctrl.Finish()
//	ctx := context.Background()
//	mockVacancyRepo := mock.NewMockRepository(ctrl)
//	mockUserRepo := mock.NewMockUserClient(ctrl)
//	useCase := vacUseCase.NewUseCase(mockVacancyRepo, mockUserRepo)
//
//	vacancy := []models.Vacancy{
//		{
//			ID:          1,
//			VacancyName: "Keke",
//			CustomerID:  1,
//			Salary:      1488,
//			Description: "Aue jizn voram",
//			Category:    "Back",
//			ExecutorID:  2,
//		},
//	}
//
//	expectVacancy := []models.Vacancy{
//		{
//			ID:          1,
//			VacancyName: "Keke",
//			CustomerID:  1,
//			Salary:      1488,
//			Description: "Aue jizn voram",
//			Category:    "Back",
//			ExecutorID:  2,
//			Img:         "kek",
//			Login:       "Mem",
//		},
//	}
//
//	mockVacancyRepo.EXPECT().
//		GetArchiveVacancies(ctx).
//		Times(1).
//		Return(vacancy, nil)
//	mockUserRepo.EXPECT().
//		GetUserById(ctx, &api.UserRequest{Id: vacancy[0].CustomerID}).
//		Times(1).
//		Return(&api.UserInfoResponse{Login: "Mem", Img: "kek"}, nil)
//
//	resVacancies, err := useCase.GetArchiveVacancies(ctx)
//
//	require.Equal(t, expectVacancy, resVacancies)
//	require.NoError(t, err)
//
//	mockVacancyRepo.EXPECT().
//		GetArchiveVacancies(ctx).
//		Times(1).
//		Return(vacancy, nil)
//	mockUserRepo.EXPECT().
//		GetUserById(ctx, &api.UserRequest{Id: vacancy[0].CustomerID}).
//		Times(1).
//		Return(&api.UserInfoResponse{Login: "Mem", Img: "kek"}, errors.New("GRPC err"))
//
//	resVacancies, err = useCase.GetArchiveVacancies(ctx)
//	require.Error(t, err)
//
//	mockVacancyRepo.EXPECT().
//		GetArchiveVacancies(ctx).
//		Times(1).
//		Return(nil, nil)
//
//	emptyVacancy := []models.Vacancy{}
//	resVacancies, err = useCase.GetArchiveVacancies(ctx)
//	require.Equal(t, emptyVacancy, resVacancies)
//	require.NoError(t, err)
//}

func TestSearchVacancies(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()
	mockVacancyrRepo := mock.NewMockRepository(ctrl)
	mockUserRepo := mock.NewMockUserClient(ctrl)
	useCase := vacUseCase.NewUseCase(mockVacancyrRepo, mockUserRepo)

	keyword := "Aue"
	vacancies := []models.Vacancy{
		{
			ID:          1,
			VacancyName: "Keke",
			CustomerID:  1,
			Salary:      1488,
			Description: "Aue jizn voram",
			Category:    "Back",
			ExecutorID:  2,
		},
	}

	expectVacancies := []models.Vacancy{
		{
			ID:          1,
			VacancyName: "Keke",
			CustomerID:  1,
			Salary:      1488,
			Description: "Aue jizn voram",
			Category:    "Back",
			ExecutorID:  2,
			Img:         "kek",
			Login:       "Mem",
		},
	}

	mockVacancyrRepo.EXPECT().
		SearchVacancy(keyword, ctx).
		Times(1).
		Return(vacancies, nil)
	mockUserRepo.EXPECT().
		GetUserById(ctx, &api.UserRequest{Id: vacancies[0].CustomerID}).
		Times(1).
		Return(&api.UserInfoResponse{Login: "Mem", Img: "kek"}, nil)

	resVacancies, err := useCase.SearchVacancy(keyword, ctx)

	require.Equal(t, expectVacancies, resVacancies)
	require.NoError(t, err)

	mockVacancyrRepo.EXPECT().
		SearchVacancy(keyword, ctx).
		Times(1).
		Return(nil, nil)

	resVacancies, err = useCase.SearchVacancy(keyword, ctx)

	emptyVacancies := []models.Vacancy{}
	require.Equal(t, emptyVacancies, resVacancies)
	require.NoError(t, err)

	mockVacancyrRepo.EXPECT().
		SearchVacancy(keyword, ctx).
		Times(1).
		Return(vacancies, errors.New("DB error"))

	_, err = useCase.SearchVacancy(keyword, ctx)

	require.Error(t, err)

	mockVacancyrRepo.EXPECT().
		SearchVacancy(keyword, ctx).
		Times(1).
		Return(vacancies, nil)
	mockUserRepo.EXPECT().
		GetUserById(ctx, &api.UserRequest{Id: vacancies[0].CustomerID}).
		Times(1).
		Return(&api.UserInfoResponse{Login: "Mem", Img: "kek"}, errors.New("GRPC err"))

	_, err = useCase.SearchVacancy(keyword, ctx)

	require.Error(t, err)
}
