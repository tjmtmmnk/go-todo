package resolver

import (
	"github.com/tjmtmmnk/go-todo/graph/schema"
	"github.com/tjmtmmnk/go-todo/pkg/dbx"
)

type Resolver struct {
	DB *dbx.DB
}

func (r *Resolver) Query() schema.QueryResolver {
	return &queryResolver{r}
}

func (r *Resolver) Todo() schema.TodoResolver {
	return &todoResolver{r}
}

func (r *Resolver) User() schema.UserResolver {
	//TODO implement me
	panic("implement me")
}

var _ schema.ResolverRoot = &Resolver{}

func (r *Resolver) Mutation() schema.MutationResolver {
	return &mutationResolver{r}
}
