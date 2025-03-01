package postgres

import (
	"database/sql"

	_ "github.com/lib/pq"
)

func NewDB(strConn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", strConn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
