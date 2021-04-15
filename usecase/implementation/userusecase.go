package implementation

import (
	"FL_2/model"
	"FL_2/store"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/microcosm-cc/bluemonday"
	"github.com/pkg/errors"
	"math/rand"
)

const (
	passwordSalt = "asdknj279312kasl0sshALkMnHG"
	saltLength = 8;
)

const (
	MinPswdLenght    int = 5
	MaxPswdLength    int = 300
	userUseCaseError     = "User use case error"
)

var(
	ErrBadPassword = errors.New("Bad password")
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

func (u *UserUseCase) beforeCreate(user *model.User) error {
	salt := make([]byte, saltLength)
	_, err := rand.Read(salt)
	if err != nil{
		return err
	}
	user.EncryptPassword = hashPass(salt, user.Password)
	return nil
}

func (u *UserUseCase) sanitize(user *model.User) {
	user.Password = ""
}

func (u *UserUseCase) UserVerification(email string, password string) (*model.User, error) {
	user, err := u.store.User().FindByEmail(email)
	if err != nil {
		return nil, errors.Wrap(err, userUseCaseError)
	}
	if !compPass(user.EncryptPassword, password) {
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
	var oldUser *model.User
	var err error
	if user.OldPassword != ""{
		oldUser, err = u.store.User().FindByID(user.ID)
		if  err != nil{
			return nil, errors.Wrap(err, userUseCaseError)
		}
		if !compPass(oldUser.EncryptPassword, user.OldPassword){
			return nil, ErrBadPassword
		}
	} else {
		oldUser, err = u.store.User().FindByID(user.ID)
	}
	err = u.beforeCreate(&user)
	if err != nil{
		return nil, err
	}

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