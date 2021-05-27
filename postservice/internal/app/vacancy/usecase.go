package vacancy

import (
	"context"
	"post/internal/app/models"
)

//go:generate mockgen -destination=./mock/mock_usecase.go -package=mock post/internal/app/vacancy UseCase
type UseCase interface {
	Create(vacancy models.Vacancy, ctx context.Context) (*models.Vacancy, error)
	FindByID(userID uint64, ctx context.Context) (*models.Vacancy, error)
	GetActualVacancies(ctx context.Context) ([]models.Vacancy, error)
	ChangeVacancy(vacancy models.Vacancy, ctx context.Context) (models.Vacancy, error)
	DeleteVacancy(id uint64, ctx context.Context) error
	FindByUserID(userID uint64, ctx context.Context) ([]models.Vacancy, error)
	SelectExecutor(vacancy models.Vacancy, ctx context.Context) error
	DeleteExecutor(vacancy models.Vacancy, ctx context.Context) error
	CloseVacancy(vacancyID uint64, ctx context.Context) error
	GetArchiveVacancies(info models.UserBasicInfo, ctx context.Context) ([]models.Vacancy, error)
	SearchVacancy(keyword string, ctx context.Context) ([]models.Vacancy, error)
	SuggestVacancyTitle(suggestWord string, ctx context.Context) ([]models.SuggestVacancyTittle, error)
}
