package postgresstore

import (
	"FL_2/model"
	"database/sql"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	"strconv"
)

type UserRepository struct {
	store *Store
}

func (u *UserRepository) insertToUserSpecTable(userID uint64, specID uint64) error {
	_, err := u.store.db.NamedExec(
		`INSERT INTO user_specializes (
                   			user_id, specialize_id
                        )
						VALUES (
							:userID, :specID
						)`,
		map[string]interface{}{
			"userID": strconv.FormatUint(userID, 10),
			"specID": strconv.FormatUint(specID, 10),
		})
	if err != nil {
		return errors.Wrap(err, sqlDbSourceError)
	}
	return nil
}

func (u *UserRepository) insertToSpecTable(specName string) (uint64, error) {
	var specID uint64 = 0
	err := u.store.db.QueryRow(
		`INSERT INTO specializes (
    						specialize_name
    					) 
    					VALUES (
    						$1
    					)  RETURNING id`, specName).Scan(&specID)
	if err != nil {
		return 0, errors.Wrap(err, sqlDbSourceError)
	}
	return specID, nil
}

func (u *UserRepository) Create(user model.User) (uint64, error) {
	var userID uint64
	err := u.store.db.QueryRow(
		`INSERT INTO users (
                   email, 
                   password, 
                   login, 
                   name_surname, 
                   about, 
                   executor 
		)
        VALUES (
                $1, 
                $2, 
                $3,
				$4,
				$5,
				$6	
                ) RETURNING id`,
		user.Email,
		user.Password,
		user.Login,
		user.NameSurname,
		user.About,
		user.Executor).Scan(&userID)
	if err != nil {
		pqErr := &pq.Error{}
		if errors.As(err, &pqErr) {
			if pqErr.Code == duplicateErrorCode {
				return 0, errors.Wrap(&DuplicateSourceErr{
					Err: err,
				}, sqlDbSourceError)
			}
		}
		return 0, errors.Wrap(err, sqlDbSourceError)
	}
	for _, spec := range user.Specializes {
		rows, err := u.store.db.Queryx("SELECT * FROM specializes WHERE specialize_name = $1", spec)
		if err != nil {
			return 0, errors.Wrap(err, sqlDbSourceError)
		}

		// если в таблице специализации нет данной специализации - добавляем ее в таблицу специализацй
		// а затем добаляем в талбицу соответствия юзер-специализация
		if rows.Next() == false {
			specID, err := u.insertToSpecTable(spec)
			if err != nil {
				return 0, errors.Wrap(err, sqlDbSourceError)
			}
			if err := u.insertToUserSpecTable(userID, specID); err != nil {
				return 0, err
			}
		} else {
			// в ином случае просто добавляем в таблицу соответствий
			specialize := model.Specialize{}
			err := rows.StructScan(&specialize)
			if err != nil {
				return 0, errors.Wrap(err, sqlDbSourceError)
			}

			specID := specialize.ID

			if err := u.insertToUserSpecTable(userID, specID); err != nil {
				return 0, errors.Wrap(err, sqlDbSourceError)
			}
		}
	}
	return userID, nil
}

func (u *UserRepository) FindByEmail(email string) (*model.User, error) {
	user := model.User{}
	err := u.store.db.Get(&user, "SELECT * from users WHERE users.email=$1 ", email)
	if err != nil {
		return nil, errors.Wrap(err, sqlDbSourceError)
	}
	if user.Executor {
		rows, err := u.store.db.Queryx("SELECT array_agg(specialize_name) AS specializes FROM specializes "+
			"INNER JOIN user_specializes us on specializes.id = us.specialize_id "+
			"WHERE users.email = $1", email)
		if err != nil {
			return nil, errors.Wrap(err, sqlDbSourceError)
		}
		for rows.Next() {
			if err := rows.StructScan(&user); err != nil {
				return nil, errors.Wrap(err, sqlDbSourceError)
			}
		}
	}

	return &user, nil
}

func (u *UserRepository) FindByID(id uint64) (*model.User, error) {
	user := model.User{}
	err := u.store.db.Get(&user, "SELECT * from users WHERE id=$1 ", id)
	if err != nil {
		return nil, errors.Wrap(err, sqlDbSourceError)
	}
	if user.Executor {
		rows, err := u.store.db.Queryx("SELECT array_agg(specialize_name) AS specializes FROM specializes "+
			"INNER JOIN user_specializes us on specializes.id = us.specialize_id "+
			"WHERE user_id = $1", id)
		if err != nil {
			return nil, errors.Wrap(err, sqlDbSourceError)
		}
		for rows.Next() {
			if err := rows.StructScan(&user); err != nil {
				return nil, errors.Wrap(err, sqlDbSourceError)
			}
		}
	}

	return &user, nil
}

func (u *UserRepository) ChangeUser(user model.User) (*model.User, error) {
	oldUser, err := u.FindByID(user.ID)

	if err != nil {
		return nil, errors.Wrap(err, sqlDbSourceError)
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

	tx := u.store.db.MustBegin()
	_, err = tx.NamedExec(`UPDATE users SET 
                 password =:password,
                 login =:login,
                 name_surname =:name_surname,
                 about=:about,
                 executor=:executor,
                 img=:img,
                 rating=:rating
				 WHERE id = :id`, &user)
	if err != nil {
		return nil, errors.Wrap(err, sqlDbSourceError)
	}
	if err := tx.Commit(); err != nil {
		return nil, errors.Wrap(err, sqlDbSourceError)
	}
	return &user, nil
}

func (u *UserRepository) AddSpecialize(specName string, userID uint64) error {
	specialize := model.Specialize{}
	err := u.store.db.Get(&specialize, "SELECT * FROM specializes WHERE specialize_name=$1", specName)

	if !errors.Is(err,sql.ErrNoRows) || err != nil {
		return errors.Wrap(err, sqlDbSourceError)
	}
	var specID uint64
	if errors.Is(err,sql.ErrNoRows) {
		specID, err = u.insertToSpecTable(specName)
		if err != nil {
			return errors.Wrap(err, sqlDbSourceError)
		}
	}
	if err = u.insertToUserSpecTable(userID, specID); err != nil {
		return errors.Wrap(err, sqlDbSourceError)
	}

	return nil
}

func (u *UserRepository) DelSpecialize(specName string, userID uint64) error {
	specialize := model.Specialize{}
	err := u.store.db.Get(&specialize, "SELECT * FROM specializes WHERE specialize_name=$1", specName)
	if err != nil {
		return errors.Wrap(err, sqlDbSourceError)
	}

	_, err = u.store.db.Queryx("DELETE FROM user_specializes WHERE specialize_id=$1 AND user_id =$2", specialize.ID, userID)

	if err != nil {
		return errors.Wrap(err, sqlDbSourceError)
	}
	return nil
}
