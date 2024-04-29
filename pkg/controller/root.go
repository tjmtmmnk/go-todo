package controller

import (
	"fmt"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/labstack/echo/v4"
	"github.com/tjmtmmnk/go-todo/pkg/controller/session"
	"github.com/tjmtmmnk/go-todo/pkg/db/model"
	"github.com/tjmtmmnk/go-todo/pkg/db/table"
)

func (ctl *Controller) Root(c echo.Context) error {
	ctx := c.Request().Context()

	sess, err := session.ExtractFromContext(c)
	if err != nil {
		return err
	}

	var todoModels []model.Todos

	stmt := table.Todos.
		SELECT(table.Todos.AllColumns).
		FROM(table.Todos).
		WHERE(table.Todos.UserID.EQ(mysql.Uint64(sess.UserID)))

	err = stmt.QueryContext(ctx, ctl.db, &todoModels)
	if err != nil {
		return err
	}

	fmt.Println(todoModels)

	return nil
}
