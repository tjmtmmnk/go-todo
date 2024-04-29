package controller

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/tjmtmmnk/go-todo/pkg/controller/session"
)

func (ctl *Controller) Root(c echo.Context) error {
	sess, err := session.ExtractFromContext(c)
	if err != nil {
		return err
	}
	fmt.Println(sess.UserID)
	return nil
}
