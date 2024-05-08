package resolver

import (
	"context"
	"github.com/tjmtmmnk/go-todo/graph/model"
)

type mutationResolver struct{ *Resolver }

func (m mutationResolver) CreateTodo(ctx context.Context, input model.CreateTodoInput) (*model.TodoPayload, error) {
	//TODO implement me
	panic("implement me")
}

func (m mutationResolver) UpdateTodo(ctx context.Context, input model.UpdateTodoInput) (*model.TodoPayload, error) {
	//TODO implement me
	panic("implement me")
}

func (m mutationResolver) DeleteTodo(ctx context.Context, input model.DeleteTodoInput) (*model.DeleteTodoPayload, error) {
	//TODO implement me
	panic("implement me")
}

func (m mutationResolver) CreateUser(ctx context.Context, input model.CreateUserInput) (*model.User, error) {
	//TODO implement me
	panic("implement me")
}
