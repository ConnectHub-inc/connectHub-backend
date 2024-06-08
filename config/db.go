package config

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql" // This blank import is used for its init function

	"github.com/tusmasoma/connectHub-backend/internal/log"
)

func NewDB() (*sql.DB, error) {
	ctx := context.Background()

	conf, err := NewDBConfig(ctx)
	if err != nil {
		log.Error("Failed to load database config", log.Ferror(err))
		return nil, err
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=true",
		conf.User, conf.Password, conf.Host, conf.Port, conf.DBName)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Error("Failed to connect to database", log.Fstring("dsn", dsn), log.Ferror(err))
		return nil, err
	}

	if err = db.Ping(); err != nil {
		log.Error("Failed to ping database", log.Ferror(err))
		return nil, err
	}

	log.Info("Successfully connected to database", log.Fstring("dsn", dsn))
	return db, nil
}
