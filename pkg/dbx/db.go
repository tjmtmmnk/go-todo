package dbx

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	_ "github.com/go-sql-driver/mysql"
	"github.com/moznion/go-optional"
	"sync"
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

var (
	once  sync.Once
	db    *DB
	dbErr error
)

func (mc *MySQLConnectionEnv) ToDSN() string {
	return fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?parseTime=true&interpolateParams=true", mc.User, mc.Password, mc.Host, mc.Port, mc.DBName)
}

func (mc *MySQLConnectionEnv) InitDB() error {
	once.Do(func() {
		var _db *sql.DB

		_db, dbErr = sql.Open("mysql", mc.ToDSN())

		db = &DB{_db}

		db.SetConnMaxLifetime(0)
		db.SetMaxIdleConns(5)
		db.SetMaxOpenConns(1000)

		dbErr = db.Ping()
	})

	return dbErr
}

func GetDB() *DB {
	if db == nil {
		panic("db must be initialized")
	}
	return db
}

func (db *DB) UUID() uint64 {
	var uuid uint64

	row := db.QueryRow("select uuid_short()")

	row.Scan(&uuid)

	return uuid
}

type SelectArgs struct {
	Table      mysql.Table
	ColumnList mysql.ProjectionList
	Where      optional.Option[mysql.BoolExpression]
}

func Single[T any](ctx context.Context, db qrm.Queryable, args *SelectArgs) (*T, error) {
	if len(args.ColumnList) == 0 {
		return nil, fmt.Errorf("column needed")
	}

	var stmt mysql.SelectStatement
	if args.Where.IsSome() {
		stmt = args.Table.SELECT(args.ColumnList[0], args.ColumnList[1:]...).FROM(args.Table).WHERE(args.Where.Unwrap()).LIMIT(1)
	} else {
		stmt = args.Table.SELECT(args.ColumnList[0], args.ColumnList[1:]...).FROM(args.Table).LIMIT(1)
	}

	dest := new(T)
	err := stmt.QueryContext(ctx, db, dest)
	if err != nil {
		return nil, err
	}

	return dest, nil
}

func Search[T any](ctx context.Context, table mysql.Table, columnList mysql.ProjectionList, where optional.Option[mysql.BoolExpression]) ([]T, error) {
	if len(columnList) == 0 {
		return nil, fmt.Errorf("column needed")
	}

	var stmt mysql.SelectStatement
	if where.IsSome() {
		stmt = table.SELECT(columnList[0], columnList[1:]...).FROM(table).WHERE(where.Unwrap())
	} else {
		stmt = table.SELECT(columnList[0], columnList[1:]...).FROM(table)
	}

	var dest []T
	err := stmt.QueryContext(ctx, GetDB(), &dest)
	if err != nil {
		return nil, err
	}

	return dest, nil
}

type InsertArgs struct {
	Table      mysql.Table
	ColumnList mysql.ColumnList
	Model      interface{}
}

func InsertByModel(ctx context.Context, db qrm.Executable, args *InsertArgs) error {
	stmt := args.Table.INSERT(args.ColumnList).MODEL(args.Model)

	_, err := stmt.ExecContext(ctx, db)
	if err != nil {
		return err
	}

	return nil
}

func Update(ctx context.Context, table mysql.Table, set mysql.UpdateStatement, where optional.Option[mysql.BoolExpression]) error {
	var stmt mysql.UpdateStatement

	if where.IsSome() {
		stmt = table.UPDATE().SET(set).WHERE(where.Unwrap())
	} else {
		stmt = table.UPDATE().SET(set)
	}

	_, err := stmt.ExecContext(ctx, GetDB())
	if err != nil {
		return err
	}

	return nil
}
