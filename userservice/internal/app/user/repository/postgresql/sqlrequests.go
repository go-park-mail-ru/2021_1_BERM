package postgresql

const (
	CreateUserRequest = `INSERT INTO userservice.users (email, password, login, name_surname, about, executor)
       VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
	SelectUserByEmail = "SELECT * from userservice.users WHERE users.email=$1 "
	SelectUserByID = "SELECT * from userservice.users WHERE id=$1"
	UpdateUser = `UPDATE userservice.users SET
                password =:password,
                login =:login,
                name_surname =:name_surname,
                about=:about,
				 WHERE id = :id`
)
