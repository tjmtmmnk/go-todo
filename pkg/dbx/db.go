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

type SingleArgs struct {
	Table      mysql.Table
	ColumnList mysql.ProjectionList
	Where      optional.Option[mysql.BoolExpression]
}

func Single[T any](ctx context.Context, db qrm.Queryable, args *SingleArgs) (*T, error) {
	if len(args.ColumnList) == 0 {
		return nil, fmt.Errorf("column needed")
	}

	stmt := args.Table.SELECT(args.ColumnList[0], args.ColumnList[1:]...).FROM(args.Table)
	if args.Where.IsSome() {
		stmt = stmt.WHERE(args.Where.Unwrap())
	}
	stmt = stmt.LIMIT(1)

	dest := new(T)
	err := stmt.QueryContext(ctx, db, dest)
	if err != nil {
		return nil, err
	}

	return dest, nil
}

type SearchOpts struct {
	OrderBy   optional.Option[mysql.OrderByClause]
	Limit     optional.Option[int64]
	Offset    optional.Option[int64]
	ForUpdate optional.Option[mysql.RowLock]
	GroupBy   optional.Option[mysql.GroupByClause]
	Having    optional.Option[mysql.BoolExpression]
}

type SearchArgs struct {
	Table      mysql.Table
	ColumnList mysql.ProjectionList
	Where      optional.Option[mysql.BoolExpression]
	Opts       SearchOpts
}

func Search[T any](ctx context.Context, db qrm.Queryable, args *SearchArgs) ([]T, error) {
	if len(args.ColumnList) == 0 {
		return nil, fmt.Errorf("column needed")
	}

	stmt := args.Table.
		SELECT(args.ColumnList[0], args.ColumnList[1:]...).
		FROM(args.Table)

	if args.Where.IsSome() {
		stmt = stmt.WHERE(args.Where.Unwrap())
	}

	if args.Opts.OrderBy.IsSome() {
		stmt = stmt.ORDER_BY(args.Opts.OrderBy.Unwrap())
	}
	if args.Opts.Limit.IsSome() {
		stmt = stmt.LIMIT(args.Opts.Limit.Unwrap())
	}
	if args.Opts.Offset.IsSome() {
		stmt = stmt.OFFSET(args.Opts.Offset.Unwrap())
	}
	if args.Opts.ForUpdate.IsSome() {
		stmt = stmt.FOR(args.Opts.ForUpdate.Unwrap())
	}
	if args.Opts.GroupBy.IsSome() {
		stmt = stmt.GROUP_BY(args.Opts.GroupBy.Unwrap())
	}
	if args.Opts.Having.IsSome() {
		stmt = stmt.HAVING(args.Opts.Having.Unwrap())
	}

	fmt.Println(stmt.DebugSql())

	var dest []T
	err := stmt.QueryContext(ctx, db, &dest)
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
