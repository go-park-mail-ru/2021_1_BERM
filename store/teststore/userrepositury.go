package teststore

import "fl_ru/model"

// UserRepository ...
type UserRepository struct {
	store *Store
	users map[uint64]*model.User
}

// Create ...
func (r *UserRepository) Create(u *model.User) error {
	if err := u.Validate(); err != nil {
		return err
	}

	if err := u.BeforeCreate(); err != nil {
		return err
	}

	u.Id = uint64(len(r.users) + 1)
	r.users[u.Id] = u

	return nil
}

// Find ...
func (r *UserRepository) Find(user *model.User) error {
	u, ok := r.users[user.Id]
	if !ok {
		return nil
	}
	*user = *u
	return nil
}

// FindByEmail ...
func (r *UserRepository) FindByEmail(user *model.User)  error {
	for _, u := range r.users {
		if u.Email == user.Email {
			return  nil
		}
	}

	return nil
}

func (r *UserRepository) ChangeUser(user *model.User) error{
	return nil
}