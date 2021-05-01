package tarantoolrepository

import (
	"authorizationservice/internal/models"
	"authorizationservice/internal/tools"
	customError "authorizationservice/pkg/error"
	"context"
	"github.com/tarantool/go-tarantool"
)

type Repository struct {
	conn *tarantool.Connection
}

func New(conn *tarantool.Connection) *Repository {
	return &Repository{
		conn: conn,
	}
}

func (r *Repository) Store(session *models.Session, ctx context.Context) error {
	_, err := r.conn.Insert("session", tools.EncodingSessionToTarantool(session))
	if err != nil {
		return err
	}
	return nil
}
func (r *Repository) Get(sessionID string, ctx context.Context) (*models.Session, error) {
	resp, err := r.conn.Select("session", "primary",
		0, 1, tarantool.IterEq, []interface{}{sessionID})
	if err != nil {
		return nil, err
	}
	if len(resp.Tuples()) == 0 {
		return nil, customError.ErrorNoRows
	}
	session := tools.DecodingTarantoolToSession(resp.Tuples()[0])
	return session, nil
}
func (r *Repository) Remove(sessionID string, ctx context.Context) error {
	_, err := r.conn.Delete("session", "primary", []interface{}{sessionID})
	if err != nil {
		return err
	}
	return nil
}
