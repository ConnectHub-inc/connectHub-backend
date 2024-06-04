package config

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql" // This blank import is used for its init function
)

func NewDB() (*sql.DB, error) {
	ctx := context.Background()

	conf, err := NewDBConfig(ctx)
	if err != nil {
		log.Printf("Failed to load database config: %s\n", err)
		return nil, err
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=true",
		conf.User, conf.Password, conf.Host, conf.Port, conf.DBName)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Printf("Database connection failed: %s\n", err)
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
