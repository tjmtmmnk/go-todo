package controller

import (
	"fmt"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
	session2 "github.com/tjmtmmnk/go-todo/pkg/controller/session"
	"github.com/tjmtmmnk/go-todo/pkg/db/model"
	"github.com/tjmtmmnk/go-todo/pkg/db/table"
	"github.com/tjmtmmnk/go-todo/pkg/dbx"
	"github.com/tjmtmmnk/go-todo/pkg/user"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

type CreateUserRequest struct {
	Name        string  `form:"name"`
	Nickname    *string `form:"nickname"`
	RawPassword string  `form:"raw_password"`
}

func (ctl *Controller) CreateUser(ctx echo.Context) error {
	req := new(CreateUserRequest)
	if err := ctx.Bind(req); err != nil {
		fmt.Println(err)
		return err
	}

	passwordHash, err := user.ToHash(req.RawPassword)
	if err != nil {
		return err
	}

	userModel := model.Users{
		ID:        dbx.UUID(ctl.db),
		Name:      req.Name,
		Nickname:  req.Nickname,
		Password:  string(passwordHash),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	stmt := table.Users.INSERT(table.Users.AllColumns).MODEL(userModel)
	_, err = stmt.Exec(ctl.db)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

type LoginRequest struct {
	Name        string `form:"name"`
	RawPassword string `form:"raw_password"`
}

func (ctl *Controller) Login(c echo.Context) error {
	ctx := c.Request().Context()

	req := new(LoginRequest)
	if err := c.Bind(req); err != nil {
		fmt.Println(err)
		return err
	}

	var row struct {
		ID             uint64
		HashedPassword string
	}

	stmt := table.Users.
		SELECT(
			table.Users.ID.AS("id"),
			table.Users.Password.AS("hashed_password"),
		).
		FROM(table.Users).
		WHERE(table.Users.Name.EQ(mysql.String(req.Name)))

	err := stmt.QueryContext(ctx, ctl.db, &row)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to fetch password"+err.Error())
	}

	if bcrypt.CompareHashAndPassword([]byte(row.HashedPassword), []byte(req.RawPassword)) != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "wrong password")
	}

	sess := session2.Session{UserID: row.ID}

	err = sess.Save(c, &sessions.Options{
		Path:     "/",
		Domain:   "localhost",
		MaxAge:   int(24 * time.Hour * 7 / time.Microsecond),
		HttpOnly: true,
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to save session")
	}

	fmt.Println("login!!!")

	return nil
}
