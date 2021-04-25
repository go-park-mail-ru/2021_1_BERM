package postgresql

const (
	CreateUserSpecializeRequest = `INSERT INTO user.user_specializes ( user_id, specialize_id )
			VALUES ( :userID, :specID )`
	CreateSpecializeRequest = `INSERT INTO user.specializes ( specialize_name ) 
			VALUES ( $1 )  RETURNING id`
	SelectSpecializesByUserID = "SELECT array_agg(specialize_name) AS specializes FROM user.specializes " +
		"INNER JOIN user.user_specializes us on specializes.id = us.specialize_id " +
		"WHERE user_id = $1"
	SelectSpecializesByID = "SELECT * FROM user.specializes WHERE id = $1"
	SelectSpecializesByName = "SELECT * FROM user.specializes WHERE specialize_name = $1"
)
