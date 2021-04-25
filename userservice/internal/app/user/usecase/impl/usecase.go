package impl

import (
	"github.com/pkg/errors"
	"user/Error"
	"user/internal/app/models"
	specialize "user/internal/app/specialize/repository"
	"user/internal/app/user/repository"
	"user/internal/app/user/tools"
)

type UseCase struct {
	userRepository       repository.Repository
	specializeRepository specialize.Repository
}

func (useCase *UseCase) Create(user models.NewUser) (map[string]interface{}, error) {
	if err := tools.ValidationCreateUser(&user); err != nil {
		return nil, errors.Wrap(err, "Validation error")
	}
	if err := tools.BeforeCreate(&user); err != nil {
		return nil, errors.Wrap(err, "Encrypt password error")
	}
	ID, err := useCase.userRepository.Create(&user)
	if err != nil {
		return nil, errors.Wrap(err, "Error in user repository.")
	}

	for _, spec := range user.Specializes {
		specID, err := useCase.specializeRepository.FindByName(spec)
		if err != nil {
			newErr := &Error.Error{}
			if errors.As(err, newErr) {
				if !newErr.InternalError {
					specID, err = useCase.specializeRepository.Create(spec)
					if err != nil {
						return nil, errors.Wrap(err, "Error in data sourse")
					}
				}
			}
		}
		err = useCase.specializeRepository.AssociateSpecializationWithUser(specID, ID)
		if err != nil {
			return nil, err
		}
	}

	return map[string]interface{}{
		"id":       ID,
		"executor": user.Executor,
	}, nil
}

func (useCase *UseCase) Verification(email string, password string) (map[string]interface{}, error) {
	user, err := useCase.userRepository.FindUserByEmail(email)
	if err != nil {
		return nil, errors.Wrap(err, "Error in user repository.")
	}
	if !tools.CompPass(user.Password, password) {
		return nil, &Error.Error{
			Err:           err,
			InternalError: false,
			ErrorDescription: map[string]interface{}{
				"Error": "Invalid password",
			},
		}
	}

	return map[string]interface{}{
		"id":       user.ID,
		"executor": user.Executor,
	}, nil
}

func (useCase *UseCase) GetById(ID uint64) (*models.UserInfo, error) {
	user, err := useCase.userRepository.FindUserByID(ID)
	if err != nil {
		return nil, err
	}
	if user.Executor {
		user.Specializes, err = useCase.specializeRepository.FindByUserID(user.ID)
		if err != nil {
			return nil, errors.Wrap(err, "Error in specialize repository.")
		}
	}
	return user, nil
}

func (useCase *UseCase) Change(user models.ChangeUser) (map[string]interface{}, error) {
	err := tools.ValidationChangeUser(&user)
	if err != nil {
		return nil, err
	}
	oldUser, err := useCase.userRepository.FindUserByID(user.ID)
	if err != nil {
		return nil, errors.Wrap(err, "User db error")
	}

	if !tools.CompPass(oldUser.Password, user.Password) {
		return nil, &Error.Error{
			Err:           err,
			InternalError: false,
			ErrorDescription: map[string]interface{}{
				"Error": "Bad password",
			},
		}
	}
	if user.NewPassword != "" {
		user.Password = user.NewPassword
	}
	if err := tools.BeforeChange(&user); err != nil {
		return nil, &Error.Error{
			Err:           err,
			InternalError: true,
			ErrorDescription: map[string]interface{}{
				"Error": Error.InternalServerErrorDescription,
			},
		}
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
	err = useCase.userRepository.Change(&user)
	if err != nil {
		return nil, err
	}
	for _, spec := range user.Specializes {
		specID, err := useCase.specializeRepository.FindByName(spec)
		if err != nil {
			newErr := &Error.Error{}
			if errors.As(err, newErr) {
				if !newErr.InternalError {
					specID, err = useCase.specializeRepository.Create(spec)
					if err != nil {
						return nil, errors.Wrap(err, "Error in data sourse")
					}
				}
			}
		}
		err = useCase.specializeRepository.AssociateSpecializationWithUser(specID, oldUser.ID)
		if err != nil {
			return nil, err
		}
	}
	return map[string]interface{}{
		"email" : user.Email,
		"about" : user.About,
		"executor" : user.Executor,
		"login" : user.Login,
		"name_surname" : user.NameSurname,
	}, nil
}

func New(userRep repository.Repository, specRep specialize.Repository) *UseCase{
	return &UseCase{
		specializeRepository: specRep,
		userRepository: userRep,
	}
}