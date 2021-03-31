package postgresstore

import (
	"fl_ru/model"
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
	return err
}

func (u *UserRepository) Create(user *model.User) (uint64, error) {
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
		return 0, err
	}

	for _, spec := range user.Specializes {
		rows, err := u.store.db.Queryx("SELECT * FROM specializes WHERE specialize_name = $1", spec)
		if err != nil {
			return 0, err
		}

		// если в таблице специализации нет данной специализации - добавляем ее в таблицу специализацй
		// а затем добаляем в талбицу соответствия юзер-специализация
		var specID uint64
		if rows.Next() == false {
			err := u.store.db.QueryRow(
				`INSERT INTO specializes (
    						specialize_name
    					) 
    					VALUES (
    						$1
    					)  RETURNING id`, spec).Scan(&specID)
			if err != nil {
				return 0, err
			}

			if err := u.insertToUserSpecTable(userID, specID); err != nil {
				return 0, err
			}
		} else {
			// в ином случае просто добавляем в таблицу соответствий
			specialize := model.Specialize{}
			err := rows.StructScan(&specialize)
			if err != nil {
				return 0, err
			}

			specID := specialize.ID

			if err := u.insertToUserSpecTable(userID, specID); err != nil {
				return 0, err
			}
		}
	}
	return userID, nil
}

func (u *UserRepository) FindByEmail(email string) (*model.User, error) {
	user := &model.User{}
	rows, err := u.store.db.Queryx("SELECT users.*, array_agg(specialize_name) AS specializes from users "+
		"INNER JOIN user_specializes ON users.id = user_specializes.user_id "+
		"INNER JOIN specializes ON user_specializes.specialize_id = specializes.id "+
		"WHERE users.email = $1 "+
		"GROUP BY users.id", email)

	if err != nil {
		return nil, err
	}
	for rows.Next() {
		if err := rows.StructScan(&user); err != nil {
			return nil, err
		}
	}

	return user, nil
}

func (u *UserRepository) FindById(id uint64) (*model.User, error) {
	user := &model.User{}
	rows, err := u.store.db.Queryx("SELECT users.*, array_agg(specialize_name) AS specializes from users "+
		"INNER JOIN user_specializes ON users.id = user_specializes.user_id "+
		"INNER JOIN specializes ON user_specializes.specialize_id = specializes.id "+
		"WHERE users.id = $1 "+
		"GROUP BY users.id", id)

	if err != nil {
		return nil, err
	}
	for rows.Next() {
		if err := rows.StructScan(&user); err != nil {
			return nil, err
		}
	}

	return user, nil
}

func (u *UserRepository) ChangeUser(user *model.User) (*model.User, error) {
	//userU, _:= u.FindById(user.ID)
	//
	//tx := u.store.db.MustBegin()
	//tx.NamedExec(`UPDATE users SET
    //             password =:password,
    //             login =:login,
    //             name_surname =:name_surname,
    //             about=:about,
    //             executor=:executor,
    //             img=:img,
    //             rating=:rating
	//			 WHERE id = :id`, &user)
	//if err := tx.Commit(); err != nil {
	//	return nil, err
	//}
	return &model.User{}, nil
}
