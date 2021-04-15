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
	saltLength = 8
)

const (
	MinPswdLenght    int = 5
	MaxPswdLength    int = 300
	userUseCaseError     = "User use case error"
)

var (
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
	id, err := u.store.User().AddUser(user)
	if err != nil {
		return errors.Wrap(err, userUseCaseError)
	}
	user.ID = id

	// если в таблице специализации нет данной специализации - добавляем ее в таблицу специализацй
	// а затем добаляем в талбицу соответствия юзер-специализация
	for _, spec := range user.Specializes {
		specialize, err := u.store.User().FindSpecializeByName(spec)
		if err != nil {
			return errors.Wrap(err, userUseCaseError)
		}
		if specialize.Name == "" && specialize.ID == 0 {
			specialize.ID, err = u.store.User().AddSpec(spec)
			if err != nil {
				return errors.Wrap(err, userUseCaseError)
			}

		}
		if err := u.store.User().AddUserSpec(user.ID, specialize.ID); err != nil {
			return errors.Wrap(err, userUseCaseError)
		}
	}
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
	if err != nil {
		return err
	}
	user.EncryptPassword = HashPass(salt, user.Password)
	return nil
}

func (u *UserUseCase) sanitize(user *model.User) {
	user.Password = ""
}

func (u *UserUseCase) UserVerification(email string, password string) (*model.User, error) {
	user, err := u.store.User().FindUserByEmail(email)
	if err != nil {
		return nil, errors.Wrap(err, userUseCaseError)
	}
	if user.Executor {
		user.Specializes, err = u.store.User().FindSpecializesByUserEmail(email)
		if err != nil {
			return nil, errors.Wrap(err, userUseCaseError)
		}
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
	user, err := u.store.User().FindUserByID(id)
	if err != nil {
		return nil, errors.Wrap(err, userUseCaseError)
	}
	if user.Executor {
		user.Specializes, err = u.store.User().FindSpecializesByUserID(id)
		if err != nil {
			return nil, errors.Wrap(err, userUseCaseError)
		}
	}
	u.sanitize(user)
	image, err := u.mediaStore.Image().GetImage(user.Img)
	if err != nil {
		return nil, errors.Wrap(err, userUseCaseError)
	}
	user.Img = string(image)
	return user, err
}

func (u *UserUseCase) ChangeUser(user* model.User) (*model.User, error) {
	oldUser, err := u.store.User().FindUserByID(user.ID)
	if err != nil {
		return nil, errors.Wrap(err, userUseCaseError)
	}

	if !compPass(oldUser.EncryptPassword, user.Password) {
		return nil, ErrBadPassword
	}
	if user.NewPassword != "" {
		user.Password = user.NewPassword
	}
	if err := u.beforeCreate(user); err != nil {
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
	specialize, err := u.store.User().FindSpecializeByName(specName)

	if err != nil {
		return errors.Wrap(err, userUseCaseError)
	}

	if specialize.Name == "" && specialize.ID == 0 {
		specialize.ID, err = u.store.User().AddSpec(specName)
		if err != nil {
			return errors.Wrap(err, userUseCaseError)
		}
	}

	haveSpec, err := u.store.User().IsUserHaveSpec(specialize.ID, userID)
	if err != nil {
		return errors.Wrap(err, userUseCaseError)
	}
	if haveSpec {
		//TODO: Кидать 400 а не 500
		return errors.New("Spec duplicate")
	}

	if err = u.store.User().AddUserSpec(userID, specialize.ID); err != nil {
		return errors.Wrap(err, userUseCaseError)
	}
	return nil
}

func (u *UserUseCase) DelSpecialize(specName string, userID uint64) error {
	specialize, err := u.store.User().FindSpecializeByName(specName)
	if err != nil {
		return errors.Wrap(err, userUseCaseError)
	}
	err = u.store.User().DelSpecialize(specialize.ID, userID)

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
