package todo

import (
	"time"
)

type Todo struct {
	ID        uint64    `json:"id,omitempty"`
	UserID    uint64    `json:"user_id,omitempty"`
	ItemName  string    `json:"item_name,omitempty"`
	Done      bool      `json:"done,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
