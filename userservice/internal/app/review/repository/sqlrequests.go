package repository

const (
	CreateReviewsRequest = `INSERT INTO userservice.reviews(user_id,
					to_user_id, order_id, description, score)
			VALUES ($1, $2, $3, $4, $5) RETURNING id`
	SelectAllReviewsByUseIDRequest = "SELECT * FROM userservice.reviews WHERE to_user_id = $1"
	SelectAvgScore                 = "SELECT CASE WHEN SUM(score) IS NULL THEN 0 " +
		"ELSE AVG(score) END AS rating, COUNT(*) AS reviews_count " +
		"FROM userservice.reviews WHERE to_user_id = $1"
)
