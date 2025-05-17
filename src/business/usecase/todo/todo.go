package todo

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/irdaislakhuafa/amartha-billing-engine/src/business/domain"
	"github.com/irdaislakhuafa/amartha-billing-engine/src/entity"
	entitygen "github.com/irdaislakhuafa/amartha-billing-engine/src/entity/gen"
	"github.com/irdaislakhuafa/amartha-billing-engine/src/utils/config"
	"github.com/irdaislakhuafa/amartha-billing-engine/src/utils/validation"
	"github.com/irdaislakhuafa/go-sdk/errors"
	"github.com/irdaislakhuafa/go-sdk/log"
)

type (
	Interface interface {
		Create(ctx context.Context, params entity.CreateTodoParams) (entitygen.Todo, error)
		List(ctx context.Context, params entity.ListTodoParams) ([]entitygen.Todo, entity.Pagination, error)
	}
	todo struct {
		log log.Interface
		cfg config.Config
		val *validator.Validate
		dom *domain.Domain
	}
)

func Init(log log.Interface, cfg config.Config, val *validator.Validate, dom *domain.Domain) Interface {
	return &todo{
		log: log,
		cfg: cfg,
		val: val,
		dom: dom,
	}
}

func (t *todo) Create(ctx context.Context, params entity.CreateTodoParams) (entitygen.Todo, error) {
	if err := t.val.StructCtx(ctx, params); err != nil {
		err = validation.ExtractError(err, params)
		return entitygen.Todo{}, errors.NewWithCode(errors.GetCode(err), err.Error())
	}

	result, err := t.dom.Todo.Create(ctx, params)
	if err != nil {
		return entitygen.Todo{}, errors.NewWithCode(errors.GetCode(err), err.Error())
	}

	return result, nil
}

func (t *todo) List(ctx context.Context, params entity.ListTodoParams) ([]entitygen.Todo, entity.Pagination, error) {
	if err := t.val.StructCtx(ctx, params); err != nil {
		err = validation.ExtractError(err, params)
		return nil, entity.Pagination{}, errors.NewWithCode(errors.GetCode(err), err.Error())
	}

	results, err := t.dom.Todo.List(ctx, params)
	if err != nil {
		return nil, entity.Pagination{}, errors.NewWithCode(errors.GetCode(err), err.Error())
	}

	totalItems, err := t.dom.Todo.Count(ctx, entity.CountTodoParams{
		IsDeleted: int8(params.IsDeleted),
		Search:    params.Search,
	})
	if err != nil {
		return nil, entity.Pagination{}, errors.NewWithCode(errors.GetCode(err), err.Error())
	}

	return results, entity.GenPagination(params.Page, params.Limit, int(totalItems), []string{}), nil
}
