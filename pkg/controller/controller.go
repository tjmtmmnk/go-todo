package controller

import (
	"github.com/tjmtmmnk/go-todo/pkg/dbx"
)

type Controller struct {
	db *dbx.DB
}

func NewController(db *dbx.DB) *Controller {
	return &Controller{
		db,
	}
}
