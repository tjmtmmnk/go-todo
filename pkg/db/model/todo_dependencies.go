//
// Code generated by go-jet DO NOT EDIT.
//
// WARNING: Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated
//

package model

import (
	"time"
)

type TodoDependencies struct {
	ID           uint64 `sql:"primary_key"`
	SourceTodoID uint64
	DestTodoID   uint64
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
