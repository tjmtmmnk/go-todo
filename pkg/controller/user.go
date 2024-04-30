package controller

import (
	"fmt"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
	"github.com/moznion/go-optional"
	"github.com/tjmtmmnk/go-todo/pkg/controller/session"
	"github.com/tjmtmmnk/go-todo/pkg/db/model"
	"github.com/tjmtmmnk/go-todo/pkg/db/table"
	"github.com/tjmtmmnk/go-todo/pkg/dbx"
	"github.com/tjmtmmnk/go-todo/pkg/timex"
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

func (ctl *Controller) CreateUser(c echo.Context) error {
	ctx := c.Request().Context()

	req := new(CreateUserRequest)
	if err := c.Bind(req); err != nil {
		fmt.Println(err)
		return err
	}

	passwordHash, err := user.ToHash(req.RawPassword)
	if err != nil {
		return err
	}

	userModel := model.Users{
		ID:        ctl.db.UUID(),
		Name:      req.Name,
		Nickname:  req.Nickname,
		Password:  string(passwordHash),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	err = dbx.Insert(ctx, ctl.db, table.Users, table.Users.AllColumns, userModel)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to insert")
	}

	return nil
}

func (ctl *Controller) GetUser(c echo.Context) error {
	ctx := c.Request().Context()

	sess, err := session.ExtractFromContext(c)
	if err != nil {
		return err
	}

	userModel, err := dbx.Single[model.Users](
		ctx,
		ctl.db,
		table.Users,
		table.Users.AllColumns,
		optional.Some(table.Users.ID.EQ(mysql.Uint64(sess.UserID))),
	)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to fetch user")
	}

	user := user.Users{
		ID:        user.ID(userModel.ID),
		Name:      userModel.Name,
		Nickname:  optional.FromNillable[string](userModel.Nickname),
		Password:  user.HashedPassword(userModel.Password),
		CreatedAt: userModel.CreatedAt.In(timex.JST),
		UpdatedAt: userModel.UpdatedAt.In(timex.JST),
	}

	return c.JSON(http.StatusOK, user)
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

	sess := session.Session{UserID: row.ID}

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
