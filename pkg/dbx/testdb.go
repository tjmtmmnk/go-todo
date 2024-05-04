package dbx

import (
	"database/sql"
	"github.com/tjmtmmnk/go-todo/pkg/db/table"
)

func InitTestDB() {
	env := &MySQLConnectionEnv{
		Host:     "localhost",
		Port:     "13306",
		User:     "root",
		DBName:   "test",
		Password: "example",
	}

	_db, err := sql.Open("mysql", env.ToDSN())
	if err != nil {
		panic(err)
	}

	table.UseSchema("test")

	db = &DB{_db}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
}
