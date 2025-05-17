package todo

import (
	"context"
	"database/sql"
	"time"

	"github.com/irdaislakhuafa/amartha-billing-engine/src/entity"
	entitygen "github.com/irdaislakhuafa/amartha-billing-engine/src/entity/gen"
	"github.com/irdaislakhuafa/amartha-billing-engine/src/utils/ctxkey"
	"github.com/irdaislakhuafa/go-sdk/codes"
	"github.com/irdaislakhuafa/go-sdk/convert"
	"github.com/irdaislakhuafa/go-sdk/errors"
	"github.com/irdaislakhuafa/go-sdk/log"
)

type (
	Interface interface {
		Create(ctx context.Context, params entity.CreateTodoParams) (entitygen.Todo, error)
		List(ctx context.Context, params entity.ListTodoParams) ([]entitygen.Todo, error)
		Count(ctx context.Context, params entity.CountTodoParams) (int64, error)
	}
	todo struct {
		log     log.Interface
		queries *entitygen.Queries
		db      *sql.DB
	}
)

func Init(log log.Interface, queries *entitygen.Queries, db *sql.DB) Interface {
	return &todo{
		log:     log,
		queries: queries,
		db:      db,
	}
}

func (t *todo) WithTx(ctx context.Context, tx *sql.Tx) Interface {
	return &todo{
		log:     t.log,
		queries: t.queries.WithTx(tx),
		db:      t.db,
	}
}

func (t *todo) Create(ctx context.Context, params entity.CreateTodoParams) (entitygen.Todo, error) {
	result := entitygen.Todo{
		Title:       params.Title,
		Description: params.Description,
		CreatedAt:   time.Now(),
		CreatedBy:   convert.ToSafeValue[string](ctx.Value(ctxkey.USER_ID)),
	}

	row, err := t.queries.CreateTodo(ctx, entitygen.CreateTodoParams{
		Title:       result.Title,
		Description: result.Description,
		CreatedAt:   result.CreatedAt,
		CreatedBy:   result.CreatedBy,
	})
	if err != nil {
		return entitygen.Todo{}, errors.NewWithCode(codes.CodeSQLTxExec, err.Error())
	}

	if result.ID, err = row.LastInsertId(); err != nil {
		return entitygen.Todo{}, errors.NewWithCode(codes.CodeSQLRead, err.Error())
	}

	return result, nil
}

func (t *todo) List(ctx context.Context, params entity.ListTodoParams) ([]entitygen.Todo, error) {
	if err := params.Parse(); err != nil {
		return nil, errors.NewWithCode(errors.GetCode(err), err.Error())
	}

	rows, err := t.queries.ListTodo(ctx, entitygen.ListTodoParams{
		IsDeleted: int8(params.IsDeleted),
		CONCAT:    params.Search,
		CONCAT_2:  params.Search,
		CONCAT_3:  params.Search,
		CONCAT_4:  params.Search,
		Limit:     int32(params.Limit),
		Offset:    int32(params.Page),
	})
	if err != nil {
		return nil, errors.NewWithCode(codes.CodeSQLRead, err.Error())
	}

	return rows, nil
}

func (r *todo) Count(ctx context.Context, params entity.CountTodoParams) (int64, error) {
	count, err := r.queries.CountTodo(ctx, entitygen.CountTodoParams{
		IsDeleted: params.IsDeleted,
		CONCAT:    params.Search,
		CONCAT_2:  params.Search,
	})
	if err != nil {
		return 0, errors.NewWithCode(codes.CodeSQLRead, err.Error())
	}

	return count, nil
}
