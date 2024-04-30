package main

import (
	"context"
	"fmt"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/tjmtmmnk/go-todo/pkg/controller"
	"github.com/tjmtmmnk/go-todo/pkg/dbx"
	"log/slog"
	"os"
)

func main() {
	e := echo.New()

	fmt.Println()

	dbConfig := &dbx.MySQLConnectionEnv{
		Host:     "127.0.0.1",
		Port:     "13306",
		User:     "root",
		DBName:   "devel",
		Password: "example",
	}

	db, err := dbConfig.Connect()
	defer db.Close()
	if err != nil {
		panic(err)
	}

	controller := controller.NewController(db)

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	e.File("/", "public/root.html")
	e.File("/register", "public/register_user.html")
	e.File("/login", "public/login_user.html")
	e.GET("/todo", controller.ListTodo)
	e.GET("/user", controller.GetUser)
	e.POST("/todo", controller.CreateTodo)
	e.POST("/user", controller.CreateUser)
	e.POST("/login", controller.Login)

	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogStatus:   true,
		LogURI:      true,
		LogError:    true,
		HandleError: true, // forwards error to the global error handler, so it can decide appropriate status code
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			if v.Error == nil {
				logger.LogAttrs(context.Background(), slog.LevelInfo, "REQUEST",
					slog.String("uri", v.URI),
					slog.Int("status", v.Status),
				)
			} else {
				logger.LogAttrs(context.Background(), slog.LevelError, "REQUEST_ERROR",
					slog.String("uri", v.URI),
					slog.Int("status", v.Status),
					slog.String("err", v.Error.Error()),
				)
			}
			return nil
		},
	}))
	e.Use(session.Middleware(sessions.NewCookieStore([]byte("secret"))))

	e.Logger.Fatal(e.Start(":1323"))
}
