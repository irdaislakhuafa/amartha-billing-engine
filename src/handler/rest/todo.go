package rest

import (
	"github.com/gofiber/fiber/v2"
	"github.com/irdaislakhuafa/amartha-billing-engine/src/entity"
	"github.com/irdaislakhuafa/go-sdk/codes"
	"github.com/irdaislakhuafa/go-sdk/errors"
)

func (r *rest) CreateTodo(c *fiber.Ctx) error {
	body := entity.CreateTodoParams{}
	if err := c.BodyParser(&body); err != nil {
		return r.httpResError(c, errors.NewWithCode(codes.CodeBadRequest, err.Error()))
	}

	res, err := r.uc.Todo.Create(c.UserContext(), body)
	if err != nil {
		return r.httpResError(c, err)
	}

	return r.httpResSuccess(c, codes.CodeSuccess, res, nil)
}

func (r *rest) ListTodo(c *fiber.Ctx) error {
	queries := entity.ListTodoParams{}
	if err := c.QueryParser(&queries); err != nil {
		return r.httpResError(c, errors.NewWithCode(codes.CodeBadRequest, err.Error()))
	}

	res, p, err := r.uc.Todo.List(c.UserContext(), queries)
	if err != nil {
		return r.httpResError(c, err)
	}

	return r.httpResSuccess(c, codes.CodeSuccess, res, &p)
}
