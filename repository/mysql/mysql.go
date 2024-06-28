package mysql

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql" // This blank import is used for its init function

	"github.com/tusmasoma/connectHub-backend/config"
	"github.com/tusmasoma/connectHub-backend/internal/log"
	"github.com/tusmasoma/connectHub-backend/repository"
)

type SQLExecutor interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}

type transactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) repository.TransactionRepository {
	return &transactionRepository{
		db: db,
	}
}

func (tr *transactionRepository) Transaction(ctx context.Context, fn func(ctx context.Context) error) error {
	tx, err := tr.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return err
	}

	ctx = context.WithValue(ctx, CtxTxKey(), tx)

	var done bool
	defer func() {
		if !done {
			if err = tx.Rollback(); err != nil {
				log.Error("Failed to rollback transaction", log.Ferror(err))
			}
		}
	}()

	if err = fn(ctx); err != nil {
		return err
	}

	done = true
	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

type TxKey string

func CtxTxKey() TxKey {
	return "tx"
}

func TxFromCtx(ctx context.Context) *sql.Tx {
	tx, ok := ctx.Value(CtxTxKey()).(*sql.Tx)
	if !ok {
		return nil
	}
	return tx
}

func NewMySQLDB() (*sql.DB, error) {
	ctx := context.Background()

	conf, err := config.NewDBConfig(ctx)
	if err != nil {
		log.Error("Failed to load database config", log.Ferror(err))
		return nil, err
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=true",
		conf.User, conf.Password, conf.Host, conf.Port, conf.DBName)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Critical("Failed to connect to database", log.Fstring("dsn", dsn), log.Ferror(err))
		return nil, err
	}

	if err = db.Ping(); err != nil {
		log.Critical("Failed to ping database", log.Ferror(err))
		return nil, err
	}

	log.Info("Successfully connected to database", log.Fstring("dsn", dsn))
	return db, nil
}
