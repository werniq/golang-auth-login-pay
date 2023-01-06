package driver

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func OpenDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)

	if err != nil {
		log.Panic(err)
		return nil, err
	}

	err = db.Ping()

	if err != nil {
		log.Panic(err)
		return nil, err
	}

	return db, nil
}
