package usecase

import (
	"bytes"
	"crypto/rand"
	"ff/internal/app/image"
	"ff/internal/app/models"
	"ff/internal/app/user"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/microcosm-cc/bluemonday"
	"github.com/pkg/errors"
	"golang.org/x/crypto/argon2"
)

const (
	saltLength = 8
)

const (
	MinPswdLenght    int = 5
	MaxPswdLength    int = 300
	userUseCaseError     = "User use case errors"
)

var (
	ErrBadPassword = errors.New("Bad password")
)

type UserUseCase struct {
	userRepo   user.UserRepository
	imageStore image.ImageRepository
}

func (u *UserUseCase) Create(user *models.User) error {

	if err := u.validate(user); err != nil {
		return errors.Wrap(err, userUseCaseError)
	}
	if err := u.beforeCreate(user); err != nil {
		return errors.Wrap(err, userUseCaseError)
	}

	u.sanitizeUser(user)

	if user.Specializes != nil {
		user.Executor = true
	}
	id, err := u.userRepo.AddUser(*user)
	if err != nil {
		return errors.Wrap(err, userUseCaseError)
	}
	user.ID = id

	// если в таблице специализации нет данной специализации - добавляем ее в таблицу специализацй
	// а затем добаляем в талбицу соответствия юзер-специализация
	for _, spec := range user.Specializes {
		specialize, err := u.userRepo.FindSpecializeByName(spec)
		if err != nil {
			return errors.Wrap(err, userUseCaseError)
		}
		if specialize.Name == "" && specialize.ID == 0 {
			specialize.ID, err = u.userRepo.AddSpec(spec)
			if err != nil {
				return errors.Wrap(err, userUseCaseError)
			}

		}
		if err := u.userRepo.AddUserSpec(user.ID, specialize.ID); err != nil {
			return errors.Wrap(err, userUseCaseError)
		}
	}
	return err
}

func (u *UserUseCase) validate(user *models.User) error {
	return validation.ValidateStruct(
		user,
		validation.Field(&user.Email, validation.Required, is.Email),
		validation.Field(&user.Password, validation.Required, validation.Length(MinPswdLenght, MaxPswdLength)),
		validation.Field(&user.Login, validation.Required),
		validation.Field(&user.NameSurname, validation.Required),
	)
}

func (u *UserUseCase) beforeCreate(user *models.User) error {
	salt := make([]byte, saltLength)
	_, err := rand.Read(salt)
	if err != nil {
		return err
	}
	user.EncryptPassword = hashPass(salt, user.Password)
	return nil
}

func (u *UserUseCase) sanitize(user *models.User) {
	user.Password = ""
}

func (u *UserUseCase) UserVerification(email string, password string) (*models.User, error) {
	userByEmail, err := u.userRepo.FindUserByEmail(email)
	if err != nil {
		return nil, errors.Wrap(err, userUseCaseError)
	}
	if userByEmail.Executor {
		userByEmail.Specializes, err = u.userRepo.FindSpecializesByUserEmail(email)
		if err != nil {
			return nil, errors.Wrap(err, userUseCaseError)
		}
	}
	if !compPass(userByEmail.EncryptPassword, password) {
		return nil, errors.Wrap(err, userUseCaseError)
	}
	u.sanitize(userByEmail)
	getImage, err := u.imageStore.GetImage(userByEmail.Img)
	if err != nil {
		return nil, errors.Wrap(err, userUseCaseError)
	}
	userByEmail.Img = string(getImage)
	return userByEmail, err
}

func (u *UserUseCase) FindByID(id uint64) (*models.User, error) {
	findUserByID, err := u.userRepo.FindUserByID(id)
	if err != nil {
		return nil, errors.Wrap(err, userUseCaseError)
	}
	if findUserByID.Executor {
		findUserByID.Specializes, err = u.userRepo.FindSpecializesByUserID(id)
		if err != nil {
			return nil, errors.Wrap(err, userUseCaseError)
		}
	}
	u.sanitize(findUserByID)
	getImage, err := u.imageStore.GetImage(findUserByID.Img)
	if err != nil {
		return nil, errors.Wrap(err, userUseCaseError)
	}
	findUserByID.Img = string(getImage)
	return findUserByID, err
}

func (u *UserUseCase) ChangeUser(user models.User) (*models.User, error) {
	oldUser, err := u.userRepo.FindUserByID(user.ID)
	if err != nil {
		return nil, errors.Wrap(err, userUseCaseError)
	}

	if !compPass(oldUser.EncryptPassword, user.Password) {
		return nil, ErrBadPassword
	}
	if user.NewPassword != "" {
		user.Password = user.NewPassword
	}
	if err := u.beforeCreate(&user); err != nil {
		return nil, errors.Wrap(err, userUseCaseError)
	}
	if user.Email == "" {
		user.Email = oldUser.Email
	}

	if user.About == "" {
		user.About = oldUser.About
	}

	if user.Password == "" {
		user.Password = oldUser.Password
	}

	if user.Login == "" {
		user.Login = oldUser.Login
	}

	if user.Img == "" {
		user.Img = oldUser.Img
	}

	if user.NameSurname == "" {
		user.NameSurname = oldUser.NameSurname
	}

	if user.Rating == 0 {
		user.Rating = oldUser.Rating
	}

	user.Executor = oldUser.Executor

	for _, spec := range oldUser.Specializes {
		user.Specializes = append(user.Specializes, spec)
	}

	newUser, err := u.userRepo.ChangeUser(user)
	if err != nil {
		return nil, errors.Wrap(err, userUseCaseError)
	}
	u.sanitize(newUser)
	getImage, err := u.imageStore.GetImage(newUser.Img)
	if err != nil {
		return nil, errors.Wrap(err, userUseCaseError)
	}
	newUser.Img = string(getImage)
	return newUser, err
}

func (u *UserUseCase) AddSpecialize(specName string, userID uint64) error {
	specialize, err := u.userRepo.FindSpecializeByName(specName)

	if err != nil {
		return errors.Wrap(err, userUseCaseError)
	}

	if specialize.Name == "" && specialize.ID == 0 {
		specialize.ID, err = u.userRepo.AddSpec(specName)
		if err != nil {
			return errors.Wrap(err, userUseCaseError)
		}
	}

	haveSpec, err := u.userRepo.IsUserHaveSpec(specialize.ID, userID)
	if err != nil {
		return errors.Wrap(err, userUseCaseError)
	}
	if haveSpec {
		//TODO: Кидать 400 а не 500
		return errors.New("Spec duplicate")
	}

	if err = u.userRepo.AddUserSpec(userID, specialize.ID); err != nil {
		return errors.Wrap(err, userUseCaseError)
	}
	return nil
}

func (u *UserUseCase) DelSpecialize(specName string, userID uint64) error {
	specialize, err := u.userRepo.FindSpecializeByName(specName)
	if err != nil {
		return errors.Wrap(err, userUseCaseError)
	}
	err = u.userRepo.DelSpecialize(specialize.ID, userID)

	if err != nil {
		return errors.Wrap(err, userUseCaseError)
	}

	return nil
}

func (u *UserUseCase) sanitizeUser(user *models.User) {
	sanitizer := bluemonday.UGCPolicy()
	user.Img = sanitizer.Sanitize(user.Img)
	user.Email = sanitizer.Sanitize(user.Email)
	user.Login = sanitizer.Sanitize(user.Login)
	user.NameSurname = sanitizer.Sanitize(user.NameSurname)
	user.About = sanitizer.Sanitize(user.About)
}

func hashPass(salt []byte, plainPassword string) []byte {
	hashedPass := argon2.IDKey([]byte(plainPassword), salt, 1, 64*1024, 4, 32)
	return append(salt, hashedPass...)
}

func compPass(passHash []byte, plainPassword string) bool {
	salt := make([]byte, 8)
	copy(salt, passHash[0:8])

	userPassHash := hashPass(salt, plainPassword)
	return bytes.Equal(userPassHash, passHash)
}

