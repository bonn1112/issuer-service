package database

import (
	"database/sql"
	"fmt"

	"github.com/lastrust/issuing-service/env"
)

func Open(conf env.DB) (*sql.DB, error) {
	db, err := sql.Open(
		"postgres",
		fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=disable",
			conf.User, conf.Password, conf.Host, conf.Port, conf.Database),
	)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(conf.MaxOpenConns)
	db.SetMaxIdleConns(conf.MaxIdleConns)

	if err = db.Ping(); err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
