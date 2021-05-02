package postgresql

const (
	CreateReviewsRequest = `INSERT INTO userservice.reviews(user_id, to_user_id, order_id, description, score)
			VALUES ($1, $2, $3, $4, $5) RETURNING id`
	SelectAllReviewsByUseIDRequest = "SELECT * FROM userservice.reviews WHERE to_user_id = $1"
	SelectAvgScore = "SELECT AVG(score) FROM userservice.reviews WHERE to_user_id = $1"
)
