package user

import (
	"ff/internal/app/models"
	"github.com/lib/pq"
)

type UserRepository interface {
	AddUserSpec(userID uint64, specID uint64) error
	FindSpecializeByName(specName string) (models.Specialize, error)
	IsUserHaveSpec(specID uint64, userID uint64) (bool, error)
	AddSpec(specName string) (uint64, error)
	AddUser(user models.User) (uint64, error)
	FindUserByEmail(email string) (*models.User, error)
	FindSpecializesByUserEmail(email string) (pq.StringArray, error)
	FindUserByID(id uint64) (*models.User, error)
	FindSpecializesByUserID(id uint64) (pq.StringArray, error)
	ChangeUser(user models.User) (*models.User, error)
	DelSpecialize(specID uint64, userID uint64) error
}
