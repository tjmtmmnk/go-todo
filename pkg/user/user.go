package user

import (
	"github.com/moznion/go-optional"
	"time"
)

type ID uint64

type Users struct {
	ID        ID                      `json:"id,omitempty"`
	Name      string                  `json:"name,omitempty"`
	Nickname  optional.Option[string] `json:"nickname,omitempty"`
	Password  HashedPassword          `json:"password,omitempty"`
	CreatedAt time.Time               `json:"created_at"`
	UpdatedAt time.Time               `json:"updated_at"`
}
