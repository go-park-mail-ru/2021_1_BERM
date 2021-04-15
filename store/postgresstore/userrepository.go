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

const (
	insertToUserSpecTable = `INSERT INTO ff.user_specializes (
                   			user_id, specialize_id
                        )
						VALUES (
							:userID, :specID
						)`

	insertToSpecTable = `INSERT INTO ff.specializes (
    						specialize_name
    					) 
    					VALUES (
    						$1
    					)  RETURNING id`

	insertUser = `INSERT INTO ff.users (
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
                ) RETURNING id`

	selectSpecializesByName = "SELECT * FROM ff.specializes WHERE specialize_name = $1"

	selectUserByEmail = "SELECT * from ff.users WHERE users.email=$1 "

	selectSpecializesByUserEmail = "SELECT array_agg(specialize_name) AS specializes FROM ff.specializes " +
		"INNER JOIN ff.user_specializes us on specializes.id = us.specialize_id " +
		"INNER JOIN ff.users u on us.user_id = u.id " +
		"WHERE u.email = $1"

	selectUserByID = "SELECT * from ff.users WHERE id=$1"

	selectSpecializesByUserID = "SELECT array_agg(specialize_name) AS specializes FROM ff.specializes " +
		"INNER JOIN ff.user_specializes us on specializes.id = us.specialize_id " +
		"WHERE user_id = $1"

	updateUser = `UPDATE ff.users SET 
                 password =:password,
                 login =:login,
                 name_surname =:name_surname,
                 about=:about,
                 executor=:executor,
                 img=:img,
                 rating=:rating
				 WHERE id = :id`

	deleteSpecializes = "DELETE FROM ff.user_specializes WHERE specialize_id=$1 AND user_id =$2"

	selectUserIDAndSpecID = "SELECT * FROM ff.user_specializes WHERE specialize_id=$1 AND user_id=$2"
)

func (u *UserRepository) AddUserSpec(userID uint64, specID uint64) error {
	_, err := u.store.Db.NamedExec(
		insertToUserSpecTable,
		map[string]interface{}{
			"userID": strconv.FormatUint(userID, 10),
			"specID": strconv.FormatUint(specID, 10),
		})
	if err != nil {
		return errors.Wrap(err, sqlDbSourceError)
	}
	return nil
}

func (u *UserRepository) FindSpecializeByName(specName string) (model.Specialize, error) {
	specialize := model.Specialize{}
	err := u.store.Db.Get(&specialize, selectSpecializesByName, specName)
	if errors.Is(err, sql.ErrNoRows) {
		return model.Specialize{
			ID:   0,
			Name: "",
		}, nil
	}
	if err != nil {
		return model.Specialize{}, errors.Wrap(err, sqlDbSourceError)
	}
	return specialize, nil
}

func (u *UserRepository) IsUserHaveSpec(specID uint64, userID uint64) (bool, error) {
	rows, err := u.store.Db.Queryx(selectUserIDAndSpecID, specID, userID)
	if err != nil {
		return true, errors.Wrap(err, sqlDbSourceError)
	}
	if rows.Next() != false {
		return true, nil
	}
	return false, nil
}

func (u *UserRepository) AddSpec(specName string) (uint64, error) {
	var specID uint64 = 0
	err := u.store.Db.QueryRow(
		insertToSpecTable, specName).Scan(&specID)
	if err != nil {
		return 0, errors.Wrap(err, sqlDbSourceError)
	}
	return specID, nil
}

func (u *UserRepository) AddUser(user *model.User) (uint64, error) {
	var userID uint64
	err := u.store.Db.QueryRow(
		insertUser,
		user.Email,
		user.EncryptPassword,
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
	return userID, nil
}

func (u *UserRepository) FindUserByEmail(email string) (*model.User, error) {
	user := model.User{}
	err := u.store.Db.Get(&user, selectUserByEmail, email)
	if err != nil {
		return nil, errors.Wrap(err, sqlDbSourceError)
	}
	return &user, nil
}

func (u *UserRepository) FindSpecializesByUserEmail(email string) (pq.StringArray, error) {
	user := model.User{}
	rows, err := u.store.Db.Queryx(selectSpecializesByUserEmail, email)
	if err != nil {
		return nil, errors.Wrap(err, sqlDbSourceError)
	}
	for rows.Next() {
		if err := rows.StructScan(&user); err != nil {
			return nil, errors.Wrap(err, sqlDbSourceError)
		}
	}
	return user.Specializes, nil
}

func (u *UserRepository) FindUserByID(id uint64) (*model.User, error) {
	user := model.User{}
	err := u.store.Db.Get(&user, selectUserByID, id)
	if err != nil {
		return nil, errors.Wrap(err, sqlDbSourceError)
	}
	return &user, nil
}

func (u *UserRepository) FindSpecializesByUserID(id uint64) (pq.StringArray, error) {
	user := model.User{}
	rows, err := u.store.Db.Queryx(selectSpecializesByUserID, id)
	if err != nil {
		return nil, errors.Wrap(err, sqlDbSourceError)
	}
	for rows.Next() {
		if err := rows.StructScan(&user); err != nil {
			return nil, errors.Wrap(err, sqlDbSourceError)
		}
	}
	return user.Specializes, nil
}


func (u *UserRepository) ChangeUser(user *model.User) (*model.User, error) {
	tx := u.store.db.MustBegin()
	_, err := tx.NamedExec(updateUser, user)

	if err != nil {
		return nil, errors.Wrap(err, sqlDbSourceError)
	}
	if err := tx.Commit(); err != nil {
		return nil, errors.Wrap(err, sqlDbSourceError)
	}
	return user, nil
}

func (u *UserRepository) DelSpecialize(specID uint64, userID uint64) error {
	_, err := u.store.Db.Queryx(deleteSpecializes, specID, userID)
	if err != nil {
		return errors.Wrap(err, sqlDbSourceError)
	}
	return nil
}
