package usecase

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"user/internal/app/models"
	review2 "user/internal/app/review"
	specialize2 "user/internal/app/specialize"
	"user/internal/app/user"
	"user/internal/app/user/tools"
	"user/internal/app/user/tools/passwordencrypt"
	customErrors "user/pkg/error"
)

type UseCase struct {
	userRepository       user.Repository
	specializeRepository specialize2.Repository
	reviewsRepository    review2.Repository
	encrypter            tools.PasswordEncrypter
}

func (useCase *UseCase) SetImg(ID uint64, img string, ctx context.Context) error {
	err := useCase.userRepository.SetUserImg(ID, img, ctx)
	if err != nil {
		return err
	}
	return err
}

func (useCase *UseCase) Create(user models.NewUser, ctx context.Context) (*models.UserBasicInfo, error) {
	err := tools.ValidationCreateUser(&user)
	if err != nil {
		return nil, errors.Wrap(err, "Validation error")
	}
	if user, err = useCase.encrypter.BeforeCreate(user); err != nil {
		return nil, errors.Wrap(err, "Encrypt password error")
	}
	ID, err := useCase.userRepository.Create(user, ctx)
	if err != nil {
		return nil, errors.Wrap(err, "Error in user repository.")
	}

	for _, spec := range user.Specializes {
		specID, err := useCase.specializeRepository.FindByName(spec, ctx)
		if err != nil {
			if errors.Is(err, customErrors.ErrorNoRows) {
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

	return &models.UserBasicInfo{
		ID:       ID,
		Executor: user.Executor,
	}, nil
}

func (useCase *UseCase) Verification(email string, password string, ctx context.Context) (*models.UserBasicInfo, error) {
	user, err := useCase.userRepository.FindUserByEmail(email, ctx)
	if err != nil {
		return nil, errors.Wrap(err, "Error in user repository.")
	}
	if !useCase.encrypter.CompPass(user.Password, password) {
		return nil, customErrors.ErrorInvalidPassword
	}

	return &models.UserBasicInfo{
		ID:       user.ID,
		Executor: user.Executor,
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
	userReviewInfo, err := useCase.reviewsRepository.GetAvgScoreByUserId(ID, ctx)
	if err != nil {
		return nil, err
	}
	user.Rating = userReviewInfo.Rating
	user.ReviewCount = userReviewInfo.ReviewCount
	return user, nil
}

func (useCase *UseCase) Change(user models.ChangeUser, ctx context.Context) (*models.UserBasicInfo, error) {
	fmt.Println(user)
	err := tools.ValidationChangeUser(user)
	if err != nil {
		return nil, err
	}
	oldUser, err := useCase.userRepository.FindUserByID(user.ID, ctx)
	if err != nil {
		return nil, errors.Wrap(err, "User db error")
	}

	if !useCase.encrypter.CompPass(oldUser.Password, user.Password) {
		return nil, customErrors.ErrorInvalidPassword
	}
	if user.NewPassword != "" {
		user.Password = user.NewPassword
	}
	fmt.Println(user)
	if user, err = useCase.encrypter.BeforeChange(user); err != nil {
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
	err = useCase.userRepository.Change(user, ctx)
	if err != nil {
		return nil, err
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
		err = useCase.specializeRepository.AssociateSpecializationWithUser(specID, oldUser.ID, ctx)
		if err != nil {
			return nil, err
		}
	}

	return &models.UserBasicInfo{
		ID:          user.ID,
		Email:       user.Email,
		About:       user.About,
		Executor:    user.Executor,
		Login:       user.Login,
		NameSurname: user.NameSurname,
	}, nil
}

func New(userRep user.Repository, specRep specialize2.Repository, reviewsRepository review2.Repository) *UseCase {
	return &UseCase{
		specializeRepository: specRep,
		userRepository:       userRep,
		reviewsRepository:    reviewsRepository,
		encrypter: &passwordencrypt.PasswordEncrypter{},
	}
}
