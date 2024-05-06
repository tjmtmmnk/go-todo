package dbx

import (
	"context"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/google/go-cmp/cmp"
	"github.com/moznion/go-optional"
	"github.com/tjmtmmnk/go-todo/pkg/db/model"
	"github.com/tjmtmmnk/go-todo/pkg/db/table"
	"testing"
	"time"
)

func TestSingle(t *testing.T) {
	type args struct {
		ctx        context.Context
		db         qrm.Queryable
		selectArgs *SelectArgs
	}
	type testCase[T any] struct {
		name    string
		args    args
		want    *T
		wantErr bool
	}

	now := time.Now().UTC().Truncate(time.Second)
	db := MustConnect(t)

	testCases := []testCase[model.Todos]{
		{
			name: "can fetch 1row",
			args: args{
				ctx: context.Background(),
				db:  db,
				selectArgs: &SelectArgs{
					Table:      table.Todos,
					ColumnList: mysql.ProjectionList{table.Todos.AllColumns},
					Where:      optional.Some(table.Todos.UserID.EQ(mysql.Uint64(1))),
				},
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
	err := InsertByModel(
		context.Background(),
		db,
		&InsertArgs{
			table.Users,
			table.Users.AllColumns,
			userModel,
		},
	)
	if err != nil {
		panic(err)
	}

	todoModel1 := model.Todos{
		ID:        1,
		UserID:    userModel.ID,
		ItemName:  "",
		Done:      false,
		StartAt:   nil,
		EndAt:     nil,
		CreatedAt: now,
		UpdatedAt: now,
	}

	err = InsertByModel(
		context.Background(),
		db,
		&InsertArgs{
			table.Todos,
			table.Todos.AllColumns,
			todoModel1,
		},
	)
	if err != nil {
		panic(err)
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Single[model.Todos](tt.args.ctx, tt.args.db, tt.args.selectArgs)
			if (err != nil) != tt.wantErr {
				t.Errorf("Single() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Error(diff)
				return
			}
		})
	}

}
