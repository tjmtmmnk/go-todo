package main

import (
	"github.com/tjmtmmnk/go-todo/graph/resolver"
	"github.com/tjmtmmnk/go-todo/graph/schema"
	"github.com/tjmtmmnk/go-todo/pkg/db/table"
	"github.com/tjmtmmnk/go-todo/pkg/dbx"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	dbConfig := &dbx.MySQLConnectionEnv{
		Host:     "127.0.0.1",
		Port:     "13306",
		User:     "root",
		DBName:   "devel",
		Password: "example",
	}
	appEnv, exist := os.LookupEnv("APP_ENV")
	if !exist {
		appEnv = "devel"
	}
	table.UseSchema(appEnv)

	err := dbConfig.InitDB()
	if err != nil {
		panic(err)
	}
	defer dbx.GetDB().Close()

	srv := handler.NewDefaultServer(schema.NewExecutableSchema(schema.Config{Resolvers: &resolver.Resolver{
		DB: dbx.GetDB(),
	}}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
