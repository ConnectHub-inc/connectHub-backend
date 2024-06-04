package mysql

import (
	"context"
	"database/sql"
	"reflect"

	"github.com/doug-martin/goqu/v9"

	// Register MySQL dialect for goqu
	_ "github.com/doug-martin/goqu/v9/dialect/mysql"
)

type QueryCondition struct {
	Field string
	Value any
}

type base[T any] struct {
	db        *sql.DB
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
		return err
	}

	return nil
}

// structScanRowsは、複数行の結果をスキャンするためのメソッドです。
func (b *base[T]) structScanRows(rows *sql.Rows) ([]T, error) {
	var entitys []T
	for rows.Next() {
		var entity T
		v := reflect.ValueOf(&entity).Elem()
		t := v.Type()

		fields := make([]interface{}, t.NumField())
		for i := 0; i < t.NumField(); i++ {
			fields[i] = v.Field(i).Addr().Interface()
		}
		if err := rows.Scan(fields...); err != nil {
			return nil, err
		}
		entitys = append(entitys, entity)
	}
	return entitys, nil
}

func (b *base[T]) List(ctx context.Context, qcs []QueryCondition) ([]T, error) {
	var whereClauses []goqu.Expression
	for _, qc := range qcs {
		whereClauses = append(whereClauses, goqu.C(qc.Field).Eq(qc.Value))
	}

	query, _, err := b.dialect.From(b.tableName).Select("*").Where(whereClauses...).ToSQL()
	if err != nil {
		return nil, err
	}

	rows, err := b.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	entitys, err := b.structScanRows(rows)
	if err != nil {
		return nil, err
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return entitys, nil
}

func (b *base[T]) Get(ctx context.Context, id string) (*T, error) {
	var entity T
	query, _, err := b.dialect.From(b.tableName).Select("*").Where(goqu.C("id").Eq(id)).ToSQL()
	if err != nil {
		return nil, err
	}
	row := b.db.QueryRowContext(ctx, query)
	if err = b.structScanRow(&entity, row); err != nil {
		return nil, err
	}
	return &entity, nil
}

func (b *base[T]) Create(ctx context.Context, entity T) error {
	query, _, err := b.dialect.Insert(b.tableName).Rows(entity).ToSQL()
	if err != nil {
		return err
	}
	_, err = b.db.ExecContext(ctx, query)
	return err
}

func (b *base[T]) BatchCreate(ctx context.Context, entitys []T) error {
	query, _, err := b.dialect.Insert(b.tableName).Rows(entitys).ToSQL()
	if err != nil {
		return err
	}
	_, err = b.db.ExecContext(ctx, query)
	return err
}

func (b *base[T]) Update(ctx context.Context, id string, entity T) error {
	query, _, err := b.dialect.Update(b.tableName).Set(entity).Where(goqu.C("id").Eq(id)).ToSQL()
	if err != nil {
		return err
	}
	_, err = b.db.ExecContext(ctx, query)
	return err
}

func (b *base[T]) Delete(ctx context.Context, id string) error {
	query, _, err := b.dialect.Delete(b.tableName).Where(goqu.C("id").Eq(id)).ToSQL()
	if err != nil {
		return err
	}
	_, err = b.db.ExecContext(ctx, query)
	return err
}

func (b *base[T]) CreateOrUpdate(ctx context.Context, id string, qcs []QueryCondition, entity T) error {
	// TODO: アンチパターン(CreateOrUpdateは現状使わないこと)
	entitys, err := b.List(ctx, qcs)
	if err != nil {
		return err
	}
	if len(entitys) > 0 {
		return b.Update(ctx, id, entity)
	}
	return b.Create(ctx, entity)
}
