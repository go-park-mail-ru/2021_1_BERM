package vacancy

import (
	"context"
	"post/internal/app/models"
)

//go:generate mockgen -destination=./mock/mock_repository.go -package=mock post/internal/app/vacancy Repository
type Repository interface {
	Create(vacancy models.Vacancy, ctx context.Context) (uint64, error)
	FindByID(id uint64, ctx context.Context) (*models.Vacancy, error)
	GetActualVacancies(ctx context.Context) ([]models.Vacancy, error)
	Change(vacancy models.Vacancy, ctx context.Context) error
	DeleteVacancy(id uint64, ctx context.Context) error
	FindByExecutorID(executorID uint64, ctx context.Context) ([]models.Vacancy, error)
	FindByCustomerID(customerID uint64, ctx context.Context) ([]models.Vacancy, error)
	UpdateExecutor(vacancy models.Vacancy, ctx context.Context) error
	CreateArchive(vacancy models.Vacancy, ctx context.Context) (uint64, error)
	GetArchiveVacancies(ctx context.Context) ([]models.Vacancy, error)
	SearchVacancy(keyword string, ctx context.Context) ([]models.Vacancy, error)
	FindArchiveByID(id uint64, ctx context.Context) (*models.Vacancy, error)
	GetArchiveVacanciesByExecutorID(executorID uint64, ctx context.Context) ([]models.Vacancy, error)
	GetArchiveVacanciesByCustomerID(customerID uint64, ctx context.Context) ([]models.Vacancy, error)
	SuggestVacancyTitle(suggestWord string, ctx context.Context) ([]models.SuggestVacancyTittle, error)
}
