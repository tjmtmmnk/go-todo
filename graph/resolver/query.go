package resolver

import (
	"context"
	"github.com/tjmtmmnk/go-todo/graph/model"
	"github.com/tjmtmmnk/go-todo/graph/schema"
)

type queryResolver struct{ *Resolver }

var _ schema.QueryResolver = &queryResolver{}

func (r *queryResolver) User(ctx context.Context, id string) (*model.User, error) {
	panic("not implemented")
}
