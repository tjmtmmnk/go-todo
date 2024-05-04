package dbx

import (
	"context"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/moznion/go-optional"
	"github.com/tjmtmmnk/go-todo/pkg/db/model"
	"github.com/tjmtmmnk/go-todo/pkg/db/table"
	"testing"
	"time"
)

func TestSingle(t *testing.T) {
	type args struct {
		ctx        context.Context
		table      mysql.Table
		columnList mysql.ProjectionList
		where      optional.Option[mysql.BoolExpression]
	}
	type testCase[T any] struct {
		name    string
		args    args
		want    *T
		wantErr bool
	}

	now := time.Now().UTC()
	InitTestDB()

	testCases := []testCase[model.Todos]{
		{
			name: "can fetch 1row",
			args: args{
				ctx:        context.Background(),
				table:      table.Todos,
				columnList: mysql.ProjectionList{table.Todos.AllColumns},
				where:      optional.Some(table.Todos.UserID.EQ(mysql.Uint64(1))),
			},
			want: &model.Todos{
				ID:        1,
				UserID:    1,
				ItemName:  "",
				Done:      false,
				StartAt:   nil,
				EndAt:     nil,
				CreatedAt: now,
				UpdatedAt: now,
			},
			wantErr: false,
		},
	}

	userModel := model.Users{
		ID:        1,
		Name:      "foo2",
		Nickname:  nil,
		Password:  "",
		CreatedAt: now,
		UpdatedAt: now,
	}
	_ = InsertByModel(
		context.Background(),
		table.Users,
		table.Users.AllColumns,
		userModel,
	)

	todoModel1 := model.Todos{
		ID:        GetDB().UUID(),
		UserID:    userModel.ID,
		ItemName:  "",
		Done:      false,
		StartAt:   nil,
		EndAt:     nil,
		CreatedAt: now,
		UpdatedAt: now,
	}

	_ = InsertByModel(
		context.Background(),
		table.Todos,
		table.Todos.AllColumns,
		todoModel1,
	)

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Single[model.Todos](tt.args.ctx, tt.args.table, tt.args.columnList, tt.args.where)
			if (err != nil) != tt.wantErr {
				t.Errorf("Single() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.UserID != tt.want.UserID {
				t.Error()
				return
			}
		})
	}
}
