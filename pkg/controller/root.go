package controller

import (
	"fmt"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/labstack/echo/v4"
	"github.com/moznion/go-optional"
	"github.com/tjmtmmnk/go-todo/pkg/controller/session"
	"github.com/tjmtmmnk/go-todo/pkg/db/model"
	"github.com/tjmtmmnk/go-todo/pkg/db/table"
	"github.com/tjmtmmnk/go-todo/pkg/dbx"
)

func (ctl *Controller) Root(c echo.Context) error {
	ctx := c.Request().Context()

	sess, err := session.ExtractFromContext(c)
	if err != nil {
		return err
	}

	todoModels, err := dbx.Search[model.Todos](
		ctx,
		table.Todos,
		mysql.ProjectionList{table.Todos.AllColumns},
		optional.Some(table.Todos.UserID.EQ(mysql.Uint64(sess.UserID))),
	)
	if err != nil {
		return err
	}

	fmt.Println(todoModels)

	return nil
}
