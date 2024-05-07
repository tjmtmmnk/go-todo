package model

type Todo struct {
	ID        string  `json:"id"`
	ItemName  string  `json:"itemName"`
	Done      bool    `json:"done"`
	User      *User   `json:"user"`
	UserID    uint64  `json:"userId"`
	StartAt   *string `json:"startAt,omitempty"`
	EndAt     *string `json:"endAt,omitempty"`
	CreatedAt *string `json:"createdAt,omitempty"`
	UpdatedAt *string `json:"updatedAt,omitempty"`
}
