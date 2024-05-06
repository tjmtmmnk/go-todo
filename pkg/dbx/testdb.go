package dbx

import (
	"database/sql"
	"github.com/DATA-DOG/go-txdb"
	"github.com/tjmtmmnk/go-todo/pkg/db/table"
	"testing"
)

var env = &MySQLConnectionEnv{
	Host:     "localhost",
	Port:     "13306",
	User:     "root",
	DBName:   "test",
	Password: "example",
}

func init() {
	txdb.Register("txdb", "mysql", env.ToDSN())
}

func MustConnect(t *testing.T) *DB {
	_db, err := sql.Open("txdb", env.ToDSN())
	if err != nil {
		panic(err)
	}

	table.UseSchema("test")

	db = &DB{DB: _db}
	err = db.Ping()
	if err != nil {
		panic(err)
	}

	t.Cleanup(func() {
		if err := db.Close(); err != nil {
			panic(err)
		}
	})

	return db
}
