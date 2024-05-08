package resolver

import (
	"context"
	"github.com/tjmtmmnk/go-todo/graph/model"
)

type todoResolver struct{ *Resolver }

func (t todoResolver) ID(ctx context.Context, obj *model.Todo) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (t todoResolver) User(ctx context.Context, obj *model.Todo) (*model.User, error) {
	//TODO implement me
	panic("implement me")
}

func (t todoResolver) StartAt(ctx context.Context, obj *model.Todo) (*string, error) {
	//TODO implement me
	panic("implement me")
}

func (t todoResolver) EndAt(ctx context.Context, obj *model.Todo) (*string, error) {
	//TODO implement me
	panic("implement me")
}

func (t todoResolver) CreatedAt(ctx context.Context, obj *model.Todo) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (t todoResolver) UpdatedAt(ctx context.Context, obj *model.Todo) (string, error) {
	//TODO implement me
	panic("implement me")
}
