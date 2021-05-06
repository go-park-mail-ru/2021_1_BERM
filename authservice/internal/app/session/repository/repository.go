package repository

import (
	models2 "authorizationservice/internal/app/models"
	tools2 "authorizationservice/internal/app/session/tools"
	sessiontools2 "authorizationservice/internal/app/session/tools/sessiontools"
	customError "authorizationservice/pkg/error"
	"context"
	"github.com/tarantool/go-tarantool"
)

type Repository struct {
	conn  *tarantool.Connection
	tools tools2.SessionTools
}

func New(conn *tarantool.Connection) *Repository {
	return &Repository{
		conn:  conn,
		tools: &sessiontools2.SessionTools{},
	}
}

func (r *Repository) Store(session models2.Session, ctx context.Context) error {
	_, err := r.conn.Insert("session", r.tools.EncodingSessionToTarantool(&session))
	if err != nil {
		return err
	}
	return nil
}
func (r *Repository) Get(sessionID string, ctx context.Context) (*models2.Session, error) {
	resp, err := r.conn.Select("session", "primary",
		0, 1, tarantool.IterEq, []interface{}{sessionID})
	if err != nil {
		return nil, err
	}
	if len(resp.Tuples()) == 0 {
		return nil, customError.ErrorNoRows
	}
	session := r.tools.DecodingTarantoolToSession(resp.Tuples()[0])
	return session, nil
}
func (r *Repository) Remove(sessionID string, ctx context.Context) error {
	_, err := r.conn.Delete("session", "primary", []interface{}{sessionID})
	if err != nil {
		return err
	}
	return nil
}
