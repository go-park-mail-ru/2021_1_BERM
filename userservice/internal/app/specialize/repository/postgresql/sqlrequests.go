package postgresql

const (
	CreateUserSpecializeRequest = `INSERT INTO userservice.user_specializes ( user_id, specialize_id )
			VALUES ( :userID, :specID )`
	CreateSpecializeRequest = `INSERT INTO userservice.specializes ( specialize_name ) 
			VALUES ( $1 )  RETURNING id`
	SelectSpecializesByUserID = "SELECT array_agg(specialize_name) AS specializes FROM userservice.specializes " +
		"INNER JOIN userservice.user_specializes us on specializes.id = us.specialize_id " +
		"WHERE user_id = $1"
	SelectSpecializesByID = "SELECT * FROM userservice.specializes WHERE id = $1"
	SelectSpecializesByName = "SELECT * FROM userservice.specializes WHERE specialize_name = $1"
)
