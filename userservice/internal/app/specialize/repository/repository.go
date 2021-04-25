package repository

type Repository interface {
	Create(specialize string) (uint64, error)
	FindByUserID(userID uint64) ([]string, error)
	AssociateSpecializationWithUser(specId uint64, userId uint64)  error
	FindById(ID uint64) (string, error)
	FindByName(spec string) (uint64, error)
}
