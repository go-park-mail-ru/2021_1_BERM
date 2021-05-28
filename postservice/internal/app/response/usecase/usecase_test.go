package response_test

//
//func TestCreate(t *testing.T) {
//	ctrl := gomock.NewController(t)
//	defer ctrl.Finish()
//
//	response := &models.Response{
//		ID:              1,
//		Time:            123456789,
//		PostID:          1,
//		UserID:          2,
//		Rate:            5,
//		Text:            "Testiy eto otvratitelno",
//		OrderResponse:   true,
//		VacancyResponse: false,
//	}
//	expectResponse := &models.Response{
//		ID:              1,
//		Time:            123456789,
//		PostID:          1,
//		UserID:          2,
//		Rate:            5,
//		Text:            "Testiy eto otvratitelno",
//		OrderResponse:   true,
//		VacancyResponse: false,
//		UserImg:         "kek",
//		UserLogin:       "Mem",
//	}
//
//	id := uint64(1)
//	ctx := context.Background()
//	mockResponseRepo := mock.NewMockRepository(ctrl)
//	mockUserRepo := mock.NewMockUserClient(ctrl)
//	useCase := respUseCase.NewUseCase(mockResponseRepo, mockUserRepo)
//
//	mockResponseRepo.EXPECT().
//		Create(*expectResponse, ctx).
//		Times(1).
//		Return(id, nil)
//	mockUserRepo.EXPECT().
//		GetUserById(ctx, &api.UserRequest{Id: response.UserID}).
//		Times(1).
//		Return(&api.UserInfoResponse{Login: "Mem", Img: "kek"}, nil)
//
//	resResponse, err := useCase.Create(*response, ctx)
//
//	require.Equal(t, expectResponse, resResponse)
//	require.NoError(t, err)
//
//	mockResponseRepo.EXPECT().
//		Create(*expectResponse, ctx).
//		Times(1).
//		Return(id, errors.New("DB err"))
//	mockUserRepo.EXPECT().
//		GetUserById(ctx, &api.UserRequest{Id: response.UserID}).
//		Times(1).
//		Return(&api.UserInfoResponse{Login: "Mem", Img: "kek"}, nil)
//
//	_, err = useCase.Create(*response, ctx)
//
//	require.Error(t, err)
//
//	mockUserRepo.EXPECT().
//		GetUserById(ctx, &api.UserRequest{Id: response.UserID}).
//		Times(1).
//		Return(&api.UserInfoResponse{Login: "Mem", Img: "kek"}, errors.New("GRPC Err"))
//
//	_, err = useCase.Create(*response, ctx)
//
//	require.Error(t, err)
//}
//
//func TestFindByPostID(t *testing.T) {
//	ctrl := gomock.NewController(t)
//	defer ctrl.Finish()
//
//	response := []models.Response{
//		{
//			ID:              1,
//			Time:            123456789,
//			PostID:          1,
//			UserID:          2,
//			Rate:            5,
//			Text:            "Testiy eto otvratitelno",
//			OrderResponse:   true,
//			VacancyResponse: false,
//		},
//	}
//	expectResponse := []models.Response{
//		{
//			ID:              1,
//			Time:            123456789,
//			PostID:          1,
//			UserID:          2,
//			Rate:            5,
//			Text:            "Testiy eto otvratitelno",
//			OrderResponse:   true,
//			VacancyResponse: false,
//			UserImg:         "kek",
//			UserLogin:       "Mem",
//		},
//	}
//	var id = uint64(1)
//	ctx := context.Background()
//	mockResponseRepo := mock.NewMockRepository(ctrl)
//	mockUserRepo := mock.NewMockUserClient(ctrl)
//	useCase := respUseCase.NewUseCase(mockResponseRepo, mockUserRepo)
//
//	mockResponseRepo.EXPECT().
//		FindByOrderPostID(id, ctx).
//		Times(1).
//		Return(response, nil)
//	mockUserRepo.EXPECT().
//		GetUserById(ctx, &api.UserRequest{Id: response[0].UserID}).
//		Times(1).
//		Return(&api.UserInfoResponse{Login: "Mem", Img: "kek"}, nil)
//
//	resResponse, err := useCase.FindByPostID(id, true, false, ctx)
//
//	require.Equal(t, expectResponse, resResponse)
//	require.NoError(t, err)
//
//	response = []models.Response{
//		{
//			ID:              1,
//			Time:            123456789,
//			PostID:          1,
//			UserID:          2,
//			Rate:            5,
//			Text:            "Testiy eto otvratitelno",
//			OrderResponse:   false,
//			VacancyResponse: true,
//		},
//	}
//	expectResponse = []models.Response{
//		{
//			ID:              1,
//			Time:            123456789,
//			PostID:          1,
//			UserID:          2,
//			Rate:            5,
//			Text:            "Testiy eto otvratitelno",
//			OrderResponse:   false,
//			VacancyResponse: true,
//			UserImg:         "kek",
//			UserLogin:       "Mem",
//		},
//	}
//
//	mockResponseRepo.EXPECT().
//		FindByVacancyPostID(id, ctx).
//		Times(1).
//		Return(response, nil)
//	mockUserRepo.EXPECT().
//		GetUserById(ctx, &api.UserRequest{Id: response[0].UserID}).
//		Times(1).
//		Return(&api.UserInfoResponse{Login: "Mem", Img: "kek"}, nil)
//
//	resResponse, err = useCase.FindByPostID(id, false, true, ctx)
//
//	require.Equal(t, expectResponse, resResponse)
//	require.NoError(t, err)
//
//	mockResponseRepo.EXPECT().
//		FindByVacancyPostID(id, ctx).
//		Times(1).
//		Return(response, errors.New("DB err"))
//
//	_, err = useCase.FindByPostID(id, false, true, ctx)
//
//	require.Error(t, err)
//
//	mockResponseRepo.EXPECT().
//		FindByVacancyPostID(id, ctx).
//		Times(1).
//		Return(response, nil)
//	mockUserRepo.EXPECT().
//		GetUserById(ctx, &api.UserRequest{Id: response[0].UserID}).
//		Times(1).
//		Return(&api.UserInfoResponse{Login: "Mem", Img: "kek"}, errors.New("GRPC Err"))
//
//	_, err = useCase.FindByPostID(id, false, true, ctx)
//
//	require.Error(t, err)
//
//	mockResponseRepo.EXPECT().
//		FindByVacancyPostID(id, ctx).
//		Times(1).
//		Return(nil, nil)
//
//	resResponse, err = useCase.FindByPostID(id, false, true, ctx)
//
//	emptyResponses := []models.Response{}
//	require.Equal(t, emptyResponses, resResponse)
//	require.NoError(t, err)
//}
//
//func TestChange(t *testing.T) {
//	ctrl := gomock.NewController(t)
//	defer ctrl.Finish()
//
//	response := &models.Response{
//		ID:              1,
//		Time:            123456789,
//		PostID:          1,
//		UserID:          2,
//		Rate:            5,
//		Text:            "Testy eto otvratitelno",
//		OrderResponse:   true,
//		VacancyResponse: false,
//	}
//	expectResponse := &models.Response{
//		ID:              1,
//		Time:            123456789,
//		PostID:          1,
//		UserID:          2,
//		Rate:            5,
//		Text:            "Testy eto otvratitelno",
//		OrderResponse:   true,
//		VacancyResponse: false,
//		UserImg:         "kek",
//		UserLogin:       "Mem",
//	}
//	ctx := context.Background()
//	mockResponseRepo := mock.NewMockRepository(ctrl)
//	mockUserRepo := mock.NewMockUserClient(ctrl)
//	useCase := respUseCase.NewUseCase(mockResponseRepo, mockUserRepo)
//
//	mockResponseRepo.EXPECT().
//		ChangeOrderResponse(*response, ctx).
//		Times(1).
//		Return(response, nil)
//	mockUserRepo.EXPECT().
//		GetUserById(ctx, &api.UserRequest{Id: response.UserID}).
//		Times(1).
//		Return(&api.UserInfoResponse{Login: "Mem", Img: "kek"}, nil)
//
//	resResponse, err := useCase.Change(*response, ctx)
//
//	require.Equal(t, expectResponse, resResponse)
//	require.NoError(t, err)
//
//	response = &models.Response{
//		ID:              1,
//		Time:            123456789,
//		PostID:          1,
//		UserID:          2,
//		Rate:            5,
//		Text:            "Testy eto otvratitelno",
//		OrderResponse:   false,
//		VacancyResponse: true,
//	}
//	expectResponse = &models.Response{
//		ID:              1,
//		Time:            123456789,
//		PostID:          1,
//		UserID:          2,
//		Rate:            5,
//		Text:            "Testy eto otvratitelno",
//		OrderResponse:   false,
//		VacancyResponse: true,
//		UserImg:         "kek",
//		UserLogin:       "Mem",
//	}
//
//	mockResponseRepo.EXPECT().
//		ChangeVacancyResponse(*response, ctx).
//		Times(1).
//		Return(response, nil)
//	mockUserRepo.EXPECT().
//		GetUserById(ctx, &api.UserRequest{Id: response.UserID}).
//		Times(1).
//		Return(&api.UserInfoResponse{Login: "Mem", Img: "kek"}, nil)
//
//	resResponse, err = useCase.Change(*response, ctx)
//
//	require.Equal(t, expectResponse, resResponse)
//	require.NoError(t, err)
//
//	mockResponseRepo.EXPECT().
//		ChangeVacancyResponse(*response, ctx).
//		Times(1).
//		Return(response, errors.New("DB err"))
//
//	_, err = useCase.Change(*response, ctx)
//
//	require.Error(t, err)
//
//	mockResponseRepo.EXPECT().
//		ChangeVacancyResponse(*response, ctx).
//		Times(1).
//		Return(response, nil)
//	mockUserRepo.EXPECT().
//		GetUserById(ctx, &api.UserRequest{Id: response.UserID}).
//		Times(1).
//		Return(&api.UserInfoResponse{Login: "Mem", Img: "kek"}, errors.New("GRPC Err"))
//
//	_, err = useCase.Change(*response, ctx)
//
//	require.Error(t, err)
//}
//
//func TestDelete(t *testing.T) {
//	ctrl := gomock.NewController(t)
//	defer ctrl.Finish()
//
//	response := &models.Response{
//		ID:              1,
//		Time:            123456789,
//		PostID:          1,
//		UserID:          2,
//		Rate:            5,
//		Text:            "Testy eto otvratitelno",
//		OrderResponse:   true,
//		VacancyResponse: false,
//	}
//
//	ctx := context.Background()
//	mockResponseRepo := mock.NewMockRepository(ctrl)
//	mockUserRepo := mock.NewMockUserClient(ctrl)
//	useCase := respUseCase.NewUseCase(mockResponseRepo, mockUserRepo)
//
//	mockResponseRepo.EXPECT().
//		DeleteOrderResponse(*response, ctx).
//		Times(1).
//		Return(nil)
//
//	err := useCase.Delete(*response, ctx)
//
//	require.NoError(t, err)
//
//	response = &models.Response{
//		ID:              1,
//		Time:            123456789,
//		PostID:          1,
//		UserID:          2,
//		Rate:            5,
//		Text:            "Testy eto otvratitelno",
//		OrderResponse:   false,
//		VacancyResponse: true,
//	}
//
//	mockResponseRepo.EXPECT().
//		DeleteVacancyResponse(*response, ctx).
//		Times(1).
//		Return(nil)
//
//	err = useCase.Delete(*response, ctx)
//
//	require.NoError(t, err)
//
//	mockResponseRepo.EXPECT().
//		DeleteVacancyResponse(*response, ctx).
//		Times(1).
//		Return(errors.New("DB Err"))
//
//	err = useCase.Delete(*response, ctx)
//
//	require.Error(t, err)
//}
