package implementation

import (
	"FL_2/model"
	"FL_2/store"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/microcosm-cc/bluemonday"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

const (
	passwordSalt = "asdknj279312kasl0sshALkMnHG"
)

const (
	MinPswdLenght    int = 5
	MaxPswdLength    int = 300
	userUseCaseError     = "User use case error"
)

type UserUseCase struct {
	store      store.Store
	mediaStore store.MediaStore
}

func (u *UserUseCase) Create(user *model.User) error {

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
	id, err := u.store.User().Create(*user)
	if err != nil {
		return errors.Wrap(err, userUseCaseError)
	}
	user.ID = id
	return err
}

func (u *UserUseCase) validate(user *model.User) error {
	return validation.ValidateStruct(
		user,
		validation.Field(&user.Email, validation.Required, is.Email),
		validation.Field(&user.Password, validation.Required, validation.Length(MinPswdLenght, MaxPswdLength)),
		validation.Field(&user.Login, validation.Required),
		validation.Field(&user.NameSurname, validation.Required),
	)
}

func (u *UserUseCase) encryptPassword(password string, salt string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(password+salt), bcrypt.MinCost)
	if err != nil {
		return "", errors.Wrap(err, userUseCaseError)
	}
	return string(b), nil
}

func (u *UserUseCase) beforeCreate(user *model.User) error {
	if len(user.Password) > 0 {
		encryptPassword, err := u.encryptPassword(user.Password, passwordSalt)
		user.Password = encryptPassword
		return err
	}
	return nil
}

func (u *UserUseCase) comparePassword(user *model.User, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password+passwordSalt)) == nil
}

func (u *UserUseCase) sanitize(user *model.User) {
	user.Password = ""
}

func (u *UserUseCase) UserVerification(email string, password string) (*model.User, error) {
	user, err := u.store.User().FindByEmail(email)
	if err != nil {
		return nil, errors.Wrap(err, userUseCaseError)
	}
	if u.comparePassword(user, password) != true {
		return nil, errors.Wrap(err, userUseCaseError)
	}
	u.sanitize(user)
	image, err := u.mediaStore.Image().GetImage(user.Img)
	if err != nil {
		return nil, errors.Wrap(err, userUseCaseError)
	}
	user.Img = string(image)
	return user, err
}

func (u *UserUseCase) FindByID(id uint64) (*model.User, error) {
	user, err := u.store.User().FindByID(id)
	if err != nil {
		return nil, errors.Wrap(err, userUseCaseError)
	}
	u.sanitize(user)
	image, err := u.mediaStore.Image().GetImage(user.Img)
	if err != nil {
		return nil, errors.Wrap(err, userUseCaseError)
	}
	user.Img = string(image)
	return user, err
}

func (u *UserUseCase) ChangeUser(user model.User) (*model.User, error) {
	if err := u.beforeCreate(&user); err != nil {
		return nil, errors.Wrap(err, userUseCaseError)
	}
	u.sanitizeUser(&user)
	newUser, err := u.store.User().ChangeUser(user)
	if err != nil {
		return nil, errors.Wrap(err, userUseCaseError)
	}
	u.sanitize(newUser)
	image, err := u.mediaStore.Image().GetImage(newUser.Img)
	if err != nil {
		return nil, errors.Wrap(err, userUseCaseError)
	}
	newUser.Img = string(image)
	return newUser, err
}

func (u *UserUseCase) AddSpecialize(specName string, userID uint64) error {
	err := u.store.User().AddSpecialize(specName, userID);
	if  err != nil {
		return errors.Wrap(err, userUseCaseError)
	}
	return nil
}

func (u *UserUseCase) DelSpecialize(specName string, userID uint64) error {
	err := u.store.User().DelSpecialize(specName, userID)
	if err != nil {
		return errors.Wrap(err, userUseCaseError)
	}
	return nil
}

func (u *UserUseCase) sanitizeUser(user *model.User) {
	sanitizer := bluemonday.UGCPolicy()
	user.Img = sanitizer.Sanitize(user.Img)
	user.Email = sanitizer.Sanitize(user.Email)
	user.Login = sanitizer.Sanitize(user.Login)
	user.NameSurname = sanitizer.Sanitize(user.NameSurname)
	user.About = sanitizer.Sanitize(user.About)
}