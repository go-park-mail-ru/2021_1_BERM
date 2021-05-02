package impl

import (
	"context"
	"github.com/pkg/errors"
	"user/internal/app/models"
	review "user/internal/app/review/repository"
	specialize "user/internal/app/specialize/repository"
	"user/internal/app/user/repository"
	"user/internal/app/user/tools"
	customErrors "user/pkg/error"
)

type UseCase struct {
	userRepository       repository.Repository
	specializeRepository specialize.Repository
	reviewsRepository    review.Repository
}

func (useCase *UseCase) SetImg(ID uint64, img string, ctx context.Context) error {
	err := useCase.userRepository.SetUserImg(ID, img, ctx)
	if err != nil {
		return err
	}
	return err
}

func (useCase *UseCase) Create(user models.NewUser, ctx context.Context) (map[string]interface{}, error) {
	if err := tools.ValidationCreateUser(&user); err != nil {
		return nil, errors.Wrap(err, "Validation error")
	}
	if err := tools.BeforeCreate(&user); err != nil {
		return nil, errors.Wrap(err, "Encrypt password error")
	}
	ID, err := useCase.userRepository.Create(&user, ctx)
	if err != nil {
		return nil, errors.Wrap(err, "Error in user repository.")
	}

	for _, spec := range user.Specializes {
		specID, err := useCase.specializeRepository.FindByName(spec, ctx)
		if err != nil {
			if err == customErrors.ErrorNoRows {
				specID, err = useCase.specializeRepository.Create(spec, ctx)
				if err != nil {
					return nil, errors.Wrap(err, "Error in data sourse")
				}
			}
		}
		err = useCase.specializeRepository.AssociateSpecializationWithUser(specID, ID, ctx)
		if err != nil {
			return nil, err
		}
	}

	return map[string]interface{}{
		"id":       ID,
		"executor": user.Executor,
	}, nil
}

func (useCase *UseCase) Verification(email string, password string, ctx context.Context) (map[string]interface{}, error) {
	user, err := useCase.userRepository.FindUserByEmail(email, ctx)
	if err != nil {
		return nil, errors.Wrap(err, "Error in user repository.")
	}
	if !tools.CompPass(user.Password, password) {
		return nil, customErrors.ErrorInvalidPassword
	}

	return map[string]interface{}{
		"id":       user.ID,
		"executor": user.Executor,
	}, nil
}

func (useCase *UseCase) GetById(ID uint64, ctx context.Context) (*models.UserInfo, error) {
	user, err := useCase.userRepository.FindUserByID(ID, ctx)
	if err != nil {
		return nil, err
	}
	if user.Executor {
		user.Specializes, err = useCase.specializeRepository.FindByUserID(user.ID, ctx)
		if err != nil {
			return nil, errors.Wrap(err, "Error in specialize repository.")
		}
	}
	rating, err := useCase.reviewsRepository.GetAvgScoreByUserId(ID, ctx)
	if err != nil {
		return nil, err
	}
	user.Rating = rating
	return user, nil
}

func (useCase *UseCase) Change(user models.ChangeUser, ctx context.Context) (map[string]interface{}, error) {
	err := tools.ValidationChangeUser(&user)
	if err != nil {
		return nil, err
	}
	oldUser, err := useCase.userRepository.FindUserByID(user.ID, ctx)
	if err != nil {
		return nil, errors.Wrap(err, "User db error")
	}

	if !tools.CompPass(oldUser.Password, user.Password) {
		return nil, customErrors.ErrorInvalidPassword
	}
	if user.NewPassword != "" {
		user.Password = user.NewPassword
	}
	if err := tools.BeforeChange(&user); err != nil {
		return nil, err
	}
	if user.Email == "" {
		user.Email = oldUser.Email
	}

	if user.About == "" {
		user.About = oldUser.About
	}

	if user.Password == "" {
		user.EncryptPassword = oldUser.Password
	}

	if user.Login == "" {
		user.Login = oldUser.Login
	}

	if user.NameSurname == "" {
		user.NameSurname = oldUser.NameSurname
	}

	for _, spec := range oldUser.Specializes {
		user.Specializes = append(user.Specializes, spec)
	}
	err = useCase.userRepository.Change(&user, ctx)
	if err != nil {
		return nil, err
	}
	for _, spec := range user.Specializes {
		specID, err := useCase.specializeRepository.FindByName(spec, ctx)
		if err != nil {
			if err == customErrors.ErrorNoRows{
				specID, err = useCase.specializeRepository.Create(spec, ctx)
				if err != nil {
					return nil, errors.Wrap(err, "Error in data sourse")
				}
			}
		}
		err = useCase.specializeRepository.AssociateSpecializationWithUser(specID, oldUser.ID, ctx)
		if err != nil {
			return nil, err
		}
	}
	return map[string]interface{}{
		"email":        user.Email,
		"about":        user.About,
		"executor":     user.Executor,
		"login":        user.Login,
		"name_surname": user.NameSurname,
	}, nil
}

func New(userRep repository.Repository, specRep specialize.Repository, reviewsRepository    review.Repository) *UseCase {
	return &UseCase{
		specializeRepository: specRep,
		userRepository:       userRep,
		reviewsRepository: reviewsRepository,
	}
}
