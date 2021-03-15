package tarantoolstore

import (
	"errors"
	"fl_ru/model"
	"github.com/tarantool/go-tarantool"
)

type UserRepository struct {
	store *Store
}

func (u *UserRepository) Create(user *model.User) error {
	resp, err := u.store.conn.Insert("user", userToTarantoolData(user))
	if err == nil {
		*user = *tarantoolDataToUser(resp.Tuples()[0])
	}

	return err
}

func (u *UserRepository) FindByEmail(user *model.User) error {
	resp, err := u.store.conn.Select("user", "email_key",
		0, 1, tarantool.IterEq, []interface{}{
			user.Email,
		})
	if err != nil {
		return err
	}
	if len(resp.Tuples()) == 0 {
		return errors.New("Bad password")
	}
	*user = *tarantoolDataToUser(resp.Tuples()[0])

	return nil
}

func (u *UserRepository) Find(user *model.User) error {
	resp, err := u.store.conn.Select("user", "primary",
		0, 1, tarantool.IterEq, []interface{}{
			user.ID,
		})
	if err != nil {
		return err
	}
	if len(resp.Tuples()) == 0 {
		return errors.New("Bad password")
	}
	*user = *tarantoolDataToUser(resp.Tuples()[0])

	return nil
}

func (u *UserRepository) ChangeUser(user *model.User) error {
	resp, err := u.store.conn.Update("user", "primary", []interface{}{user.ID}, userToTarantoolChangeData(user))
	if err != nil {
		return err
	}
	*user = *tarantoolDataToUser(resp.Tuples()[0])

	return nil
}

func userToTarantoolData(user *model.User) []interface{} {
	data := []interface{}{nil}
	if user.Email == "" {
		data = append(data, nil)
	} else {
		data = append(data, user.Email)
	}
	if len(user.Password) == 0 {
		data = append(data, nil)
	} else {
		data = append(data, user.Password)
	}
	if user.UserName == "" {
		data = append(data, nil)
	} else {
		data = append(data, user.UserName)
	}
	if len(user.FirstName) == 0 {
		data = append(data, nil)
	} else {
		data = append(data, user.FirstName)
	}
	if len(user.SecondName) == 0 {
		data = append(data, nil)
	} else {
		data = append(data, user.SecondName)
	}
	if len(user.Specializes) == 0 {
		user.Executor = false
	} else {
		user.Executor = true
	}
	data = append(data, user.Executor)
	if user.About == "" {
		data = append(data, nil)
	} else {
		data = append(data, user.About)
	}
	if len(user.Specializes) == 0 {
		data = append(data, nil)
	} else {
		data = append(data, user.Specializes)
	}
	if user.ImgURL == "" {
		data = append(data, nil)
	} else {
		data = append(data, user.ImgURL)
	}

	return data
}

func userToTarantoolChangeData(user *model.User) []interface{} {
	data := []interface{}{}
	if len(user.Password) != 0 {
		data = append(data, []interface{}{"=", 2, user.Password})
	}
	if len(user.UserName) != 0 {
		data = append(data, []interface{}{"=", 3, user.UserName})
	}
	if len(user.FirstName) != 0 {
		data = append(data, []interface{}{"=", 4, user.FirstName})
	}
	if len(user.SecondName) != 0 {
		data = append(data, []interface{}{"=", 5, user.SecondName})
	}
	if user.Executor {
		data = append(data, []interface{}{"=", 6, user.Executor})
	}
	if len(user.About) != 0 {
		data = append(data, []interface{}{"=", 7, user.About})
	}
	if user.Specializes != nil {
		data = append(data, []interface{}{"=", 8, user.Specializes})
	}
	if len(user.ImgURL) != 0 {
		data = append(data, []interface{}{"=", 9, user.ImgURL})
	}

	return data
}

func tarantoolDataToUser(data []interface{}) *model.User {
	u := &model.User{}
	u.ID, _ = data[0].(uint64)
	u.Email, _ = data[1].(string)
	u.Password, _ = data[2].(string)
	u.UserName, _ = data[3].(string)
	u.FirstName, _ = data[4].(string)
	u.SecondName, _ = data[5].(string)
	u.Executor, _ = data[6].(bool)
	u.About, _ = data[7].(string)
	specializes, _ := data[8].([]interface{})
	for _, elem := range specializes {
		specialize, _ := elem.(string)
		u.Specializes = append(u.Specializes, specialize)
	}
	u.ImgURL, _ = data[9].(string)

	return u
}
