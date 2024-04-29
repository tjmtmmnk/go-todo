package controller

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/tjmtmmnk/go-todo/pkg/db/model"
	"github.com/tjmtmmnk/go-todo/pkg/db/table"
	"time"
)

type CreateTodoRequest struct {
	//UserID   uint64
	ItemName string `form:"item_name"`
	Done     bool   `form:"done"`
}

func (ctl *Controller) CreateTodo(c echo.Context) error {
	ctx := c.Request().Context()

	req := new(CreateTodoRequest)
	if err := c.Bind(req); err != nil {
		fmt.Println(err)
		return err
	}

	todoModel := model.Todos{
		ID:        2,
		UserID:    1,
		ItemName:  req.ItemName,
		Done:      req.Done,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	stmt := table.Todos.INSERT(table.Todos.AllColumns).MODEL(todoModel)
	_, err := stmt.ExecContext(ctx, ctl.db)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}
