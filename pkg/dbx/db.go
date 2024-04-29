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

func (mc *MySQLConnectionEnv) Connect() (*sql.DB, error) {
	dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?parseTime=true&interpolateParams=true", mc.User, mc.Password, mc.Host, mc.Port, mc.DBName)
	return sql.Open("mysql", dsn)
}

func UUID(db *sql.DB) uint64 {
	var uuid uint64

	row := db.QueryRow("select uuid_short()")

	row.Scan(&uuid)

	return uuid
}

func Single[T any](
	ctx context.Context,
	db *sql.DB,
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
	db *sql.DB,
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
