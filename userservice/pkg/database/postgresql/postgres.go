package postgresql

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	"user/Error"
)

type Postgres struct {
	db *sqlx.DB
}

func NewPostgres(dsn string) (*Postgres, error) {
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, err
	}
	return &Postgres{
		db: db,
	}, nil
}

func (p *Postgres) GetPostgres() *sqlx.DB {
	return p.db
}

func (p *Postgres) Close() error {
	err := p.db.Close()
	return err
}


func WrapPostgreError(err error) error{
	pqErr := &pq.Error{}
	if errors.As(err, pqErr){
		if pqErr.Code == PostgreDuplicateErrorCode{
			return &Error.Error{
				Err: err,
				InternalError: false,
				ErrorDescription: map[string]interface{}{
					"Err" : err.Error(),
				},
			}
		}
	}
	if errors.Is(err, sql.ErrNoRows){
		return &Error.Error{
			Err: err,
			InternalError: false,
			ErrorDescription: map[string]interface{}{
				"Err" : err.Error(),
			},
		}
	}
	return &Error.Error{
		Err: err,
		InternalError: true,
		ErrorDescription: Error.InternalServerErrorDescription,
	}
}