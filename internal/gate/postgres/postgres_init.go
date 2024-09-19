package postgres

import (
	"database/sql"
	"log/slog"
)

type Handler struct {
	Db *sql.DB
}

func Init(log slog.Logger, url string) (Handler, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return Handler{
			Db: nil,
		}, err
	}

	if err := db.Ping(); err != nil {
		return Handler{
			Db: nil,
		}, err
	}

	log.Info("Connected to PostgreSQL")
	return Handler{
		Db: db,
	}, nil
}
