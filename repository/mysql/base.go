package mysql

import (
	"context"
	"database/sql"
	"reflect"

	"github.com/doug-martin/goqu/v9"

	"github.com/tusmasoma/connectHub-backend/internal/log"
	"github.com/tusmasoma/connectHub-backend/repository"

	// Register MySQL dialect for goqu
	_ "github.com/doug-martin/goqu/v9/dialect/mysql"
)

type base[T any] struct {
	db        SQLExecutor
	dialect   *goqu.DialectWrapper
	tableName string
}

func newBase[T any](db *sql.DB, dialect *goqu.DialectWrapper, tableName string) *base[T] {
	return &base[T]{
		db:        db,
		dialect:   dialect,
		tableName: tableName,
	}
}

// structScanは、構造体のフィールドをスキャンするためのヘルパー関数です。
func (b *base[T]) structScanRow(entity *T, row *sql.Row) error {
	v := reflect.ValueOf(entity).Elem()
	t := v.Type()

	fields := make([]interface{}, t.NumField())
	for i := 0; i < t.NumField(); i++ {
		fields[i] = v.Field(i).Addr().Interface()
	}

	if err := row.Scan(fields...); err != nil {
		log.Error("Failed to scan row", log.Ferror(err))
		return err
	}

	return nil
}

// structScanRowsは、複数行の結果をスキャンするためのメソッドです。
func (b *base[T]) structScanRows(rows *sql.Rows) ([]T, error) {
	var entities []T
	for rows.Next() {
		var entity T
		v := reflect.ValueOf(&entity).Elem()
		t := v.Type()

		fields := make([]interface{}, t.NumField())
		for i := 0; i < t.NumField(); i++ {
			fields[i] = v.Field(i).Addr().Interface()
		}
		if err := rows.Scan(fields...); err != nil {
			log.Error("Failed to scan rows", log.Ferror(err))
			return nil, err
		}
		entities = append(entities, entity)
	}
	if err := rows.Err(); err != nil {
		log.Error("Rows iteration error", log.Ferror(err))
		return nil, err
	}
	return entities, nil
}

func (b *base[T]) List(ctx context.Context, qcs []repository.QueryCondition) ([]T, error) {
	executor := b.db
	if tx := TxFromCtx(ctx); tx != nil {
		executor = tx
	}

	var whereClauses []goqu.Expression
	for _, qc := range qcs {
		whereClauses = append(whereClauses, goqu.C(qc.Field).Eq(qc.Value))
	}

	query, _, err := b.dialect.From(b.tableName).Select("*").Where(whereClauses...).ToSQL()
	if err != nil {
		log.Error("Failed to generate SQL query", log.Ferror(err))
		return nil, err
	}

	rows, err := executor.QueryContext(ctx, query)
	if err != nil {
		log.Error("Failed to execute query", log.Ferror(err))
		return nil, err
	}
	defer rows.Close()

	entities, err := b.structScanRows(rows)
	if err != nil {
		return nil, err
	}
	return entities, nil
}

func (b *base[T]) Get(ctx context.Context, id string) (*T, error) {
	executor := b.db
	if tx := TxFromCtx(ctx); tx != nil {
		executor = tx
	}

	var entity T
	query, _, err := b.dialect.From(b.tableName).Select("*").Where(goqu.C("id").Eq(id)).ToSQL()
	if err != nil {
		log.Error("Failed to generate SQL query", log.Ferror(err))
		return nil, err
	}
	row := executor.QueryRowContext(ctx, query)
	if err = b.structScanRow(&entity, row); err != nil {
		return nil, err
	}
	return &entity, nil
}

func (b *base[T]) Create(ctx context.Context, entity T) error {
	executor := b.db
	if tx := TxFromCtx(ctx); tx != nil {
		executor = tx
	}

	query, _, err := b.dialect.Insert(b.tableName).Rows(entity).ToSQL()
	if err != nil {
		log.Error("Failed to generate SQL query", log.Ferror(err))
		return err
	}
	_, err = executor.ExecContext(ctx, query)
	if err != nil {
		log.Error("Failed to execute query", log.Ferror(err))
		return err
	}
	return nil
}

func (b *base[T]) BatchCreate(ctx context.Context, entities []T) error {
	executor := b.db
	if tx := TxFromCtx(ctx); tx != nil {
		executor = tx
	}
	if len(entities) == 0 {
		log.Warn("No entities to insert")
		return nil
	}
	query, _, err := b.dialect.Insert(b.tableName).Rows(entities).ToSQL()
	if err != nil {
		log.Error("Failed to generate SQL query", log.Ferror(err))
		return err
	}
	_, err = executor.ExecContext(ctx, query)
	if err != nil {
		log.Error("Failed to execute query", log.Ferror(err))
		return err
	}
	return nil
}

func (b *base[T]) Update(ctx context.Context, id string, entity T) error {
	executor := b.db
	if tx := TxFromCtx(ctx); tx != nil {
		executor = tx
	}

	query, _, err := b.dialect.Update(b.tableName).Set(entity).Where(goqu.C("id").Eq(id)).ToSQL()
	if err != nil {
		log.Error("Failed to generate SQL query", log.Ferror(err))
		return err
	}
	_, err = executor.ExecContext(ctx, query)
	if err != nil {
		log.Error("Failed to execute query", log.Ferror(err))
		return err
	}
	return nil
}

func (b *base[T]) Delete(ctx context.Context, id string) error {
	executor := b.db
	if tx := TxFromCtx(ctx); tx != nil {
		executor = tx
	}

	query, _, err := b.dialect.Delete(b.tableName).Where(goqu.C("id").Eq(id)).ToSQL()
	if err != nil {
		log.Error("Failed to generate SQL query", log.Ferror(err))
		return err
	}
	_, err = executor.ExecContext(ctx, query)
	if err != nil {
		log.Error("Failed to execute query", log.Ferror(err))
		return err
	}
	return nil
}

func (b *base[T]) CreateOrUpdate(ctx context.Context, id string, qcs []repository.QueryCondition, entity T) error {
	if tx := TxFromCtx(ctx); tx != nil {
		b.db = tx
	}

	// TODO: アンチパターン(CreateOrUpdateは現状使わないこと)
	entities, err := b.List(ctx, qcs)
	if err != nil {
		return err
	}
	if len(entities) > 0 {
		return b.Update(ctx, id, entity)
	}
	return b.Create(ctx, entity)
}
