package implementation

import (
	"FL_2/model"
	"FL_2/store"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"golang.org/x/crypto/bcrypt"
)

const (
	passwordSalt = "asdknj279312kasl0sshALkMnHG"
)

const (
	MinPswdLenght int = 5
	MaxPswdLength int = 300
)



type UserUseCase struct {
	store 	   store.Store
	mediaStore store.MediaStore
}

func(u *UserUseCase)Create(user *model.User) error {

	if err := u.validate(user); err != nil {
		return err
	}
	if err := u.beforeCreate(user); err != nil {
		return err
	}

	id, err := u.store.User().Create(*user)
	user.ID = id
	return err
}


func(u *UserUseCase) validate(user *model.User) error {
	return validation.ValidateStruct(
		user,
		validation.Field(&user.Email, validation.Required, is.Email),
		validation.Field(&user.Password, validation.Required, validation.Length(MinPswdLenght, MaxPswdLength)),
		validation.Field(&user.Login, validation.Required),
		validation.Field(&user.NameSurname, validation.Required),
	)
}

func(u *UserUseCase) encryptPassword(password string, salt string) (string, error){
	b, err := bcrypt.GenerateFromPassword([]byte(password + salt), bcrypt.MinCost)
	if err != nil{
		return "", err
	}
	return string(b), nil
}


func(u *UserUseCase)beforeCreate(user *model.User) error {
	if len(user.Password) > 0 {
		encryptPassword, err := u.encryptPassword(user.Password, passwordSalt)
		user.Password = encryptPassword
		return err
	}
	return nil
}

func(u *UserUseCase)comparePassword(user *model.User, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password+passwordSalt)) == nil
}

func(u *UserUseCase)sanitize(user *model.User) {
	user.Password = ""
}


func(u *UserUseCase)UserVerification(email string, password string) (*model.User, error){
	user, err := u.store.User().FindByEmail(email)
	if err != nil{
		return nil, err;
	}
	if u.comparePassword(user, password) != true {
		return nil, err
	}
	u.sanitize(user)
	image, err := u.mediaStore.Image().GetImage(user.Img)
	if err != nil{
		return nil, err
	}
	user.Img = string(image)
	return user, err
}

func(u *UserUseCase)FindByID(id uint64) (*model.User, error){
	user, err := u.store.User().FindByID(id)
	if err != nil {
		return nil, err
	}
	u.sanitize(user)
	image, err := u.mediaStore.Image().GetImage(user.Img)
	if err != nil{
		return nil, err
	}
	user.Img = string(image)
	return user, err
}

func(u *UserUseCase)ChangeUser(user model.User) (*model.User, error){
	if err := u.beforeCreate(&user); err != nil {
		return nil, err
	}
	newUser, err := u.store.User().ChangeUser(user)
	if err != nil {
		return nil, err
	}
	u.sanitize(newUser)
	image, err := u.mediaStore.Image().GetImage(newUser.Img)
	if err != nil{
		return nil, err
	}
	newUser.Img = string(image)
	return newUser, err
}

func(u *UserUseCase)AddSpecialize(specName string, userID uint64) error{
	if err := u.store.User().AddSpecialize(specName, userID); err != nil {
		return err
	}
	return nil
}

func(u *UserUseCase)DelSpecialize(specName string, userID uint64) error{
	err := u.store.User().DelSpecialize(specName, userID)
	return err;
}
