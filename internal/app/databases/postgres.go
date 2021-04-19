package databases

import (
	"github.com/jmoiron/sqlx"
)

type Postgres struct {
	db  *sqlx.DB
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

func (p *Postgres) ClosePostgres() error {
	return p.db.Close()
}
