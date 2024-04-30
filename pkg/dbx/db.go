package dbx

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/go-jet/jet/v2/mysql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/moznion/go-optional"
)

type MySQLConnectionEnv struct {
	Host     string
	Port     string
	User     string
	DBName   string
	Password string
}

type DB struct {
	*sql.DB
}

func (mc *MySQLConnectionEnv) Connect() (*DB, error) {
	dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?parseTime=true&interpolateParams=true", mc.User, mc.Password, mc.Host, mc.Port, mc.DBName)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	return &DB{
		db,
	}, nil
}

func (db *DB) UUID() uint64 {
	var uuid uint64

	row := db.QueryRow("select uuid_short()")

	row.Scan(&uuid)

	return uuid
}

func Single[T any](
	ctx context.Context,
	db *DB,
	table mysql.Table,
	columnList mysql.ColumnList,
	where optional.Option[mysql.BoolExpression],
) (*T, error) {
	var (
		dest T
		stmt mysql.SelectStatement
	)
	if where.IsSome() {
		stmt = table.SELECT(columnList).FROM(table).WHERE(where.Unwrap()).LIMIT(1)
	} else {
		stmt = table.SELECT(columnList).FROM(table).LIMIT(1)
	}

	err := stmt.QueryContext(ctx, db, &dest)
	if err != nil {
		return nil, err
	}

	return &dest, nil
}

func Search[T any](
	ctx context.Context,
	db *DB,
	table mysql.Table,
	columnList mysql.ColumnList,
	where optional.Option[mysql.BoolExpression],
) ([]T, error) {
	var (
		dest []T
		stmt mysql.SelectStatement
	)
	if where.IsSome() {
		stmt = table.SELECT(columnList).FROM(table).WHERE(where.Unwrap())
	} else {
		stmt = table.SELECT(columnList).FROM(table)
	}

	err := stmt.QueryContext(ctx, db, &dest)
	if err != nil {
		return nil, err
	}

	return dest, nil
}

func Insert(
	ctx context.Context,
	db *DB,
	table mysql.Table,
	columnList mysql.ColumnList,
	model interface{},
) error {
	stmt := table.INSERT(columnList).MODEL(model)

	_, err := stmt.ExecContext(ctx, db)
	if err != nil {
		return err
	}

	return nil
}
