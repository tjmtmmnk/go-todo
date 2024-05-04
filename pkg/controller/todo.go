package controller

import (
	"fmt"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/labstack/echo/v4"
	"github.com/moznion/go-optional"
	"github.com/samber/lo"
	"github.com/tjmtmmnk/go-todo/pkg/controller/session"
	"github.com/tjmtmmnk/go-todo/pkg/db/model"
	"github.com/tjmtmmnk/go-todo/pkg/db/table"
	"github.com/tjmtmmnk/go-todo/pkg/dbx"
	"github.com/tjmtmmnk/go-todo/pkg/timex"
	"github.com/tjmtmmnk/go-todo/pkg/todo"
	"net/http"
	"time"
)

type CreateTodoRequest struct {
	ItemName string             `form:"item_name"`
	Done     bool               `form:"done"`
	StartAt  timex.OptionalDate `form:"start_at"`
	EndAt    timex.OptionalDate `form:"end_at"`
}

func (ctl *Controller) CreateTodo(c echo.Context) error {
	ctx := c.Request().Context()

	sess, err := session.ExtractFromContext(c)
	if err != nil {
		return err
	}

	req := new(CreateTodoRequest)
	if err := c.Bind(req); err != nil {
		fmt.Println(err)
		return err
	}

	todoModel := model.Todos{
		ID:        dbx.GetDB().UUID(),
		UserID:    sess.UserID,
		ItemName:  req.ItemName,
		Done:      req.Done,
		StartAt:   timex.UnwrapAsUTCPtr(optional.Option[time.Time](req.StartAt)),
		EndAt:     timex.UnwrapAsUTCPtr(optional.Option[time.Time](req.EndAt)),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	err = dbx.Insert(ctx, table.Todos, table.Todos.AllColumns, todoModel)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to insert todo")
	}

	return nil
}

func (ctl *Controller) ListTodo(c echo.Context) error {
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
	fmt.Println(todoModels)

	todos := lo.Map(todoModels, func(t model.Todos, _ int) *todo.Todo {
		return &todo.Todo{
			ID:        t.ID,
			UserID:    t.UserID,
			ItemName:  t.ItemName,
			Done:      t.Done,
			CreatedAt: t.CreatedAt.In(timex.JST),
			UpdatedAt: t.UpdatedAt.In(timex.JST),
		}
	})

	return c.JSON(http.StatusOK, todos)
}
