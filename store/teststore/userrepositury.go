package teststore

import "fl_ru/model"

type UserRepository struct {
	store *Store
	users map[uint64]model.User
}

func (r *UserRepository) Create(u *model.User) error {
	u.Id = uint64(len(r.users) + 1)
	r.users[u.Id] = *u
	return nil
}

func (r *UserRepository) Find(user *model.User) error {
	u, ok := r.users[user.Id]
	if !ok {
		return nil
	}
	*user = u
	return nil
}

func (r *UserRepository) FindByEmail(user *model.User) error {
	for _, u := range r.users {
		if u.Email == user.Email {
			*user = u
			return nil
		}
	}

	return nil
}

func (r *UserRepository) ChangeUser(user *model.User) error {
	r.users[user.Id] = *user
	return nil
}
