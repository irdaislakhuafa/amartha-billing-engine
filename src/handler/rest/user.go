package rest

import (
	"github.com/gofiber/fiber/v2"
	"github.com/irdaislakhuafa/amartha-billing-engine/src/entity"
	"github.com/irdaislakhuafa/go-sdk/codes"
	"github.com/irdaislakhuafa/go-sdk/errors"
)

func (r *rest) Register(c *fiber.Ctx) error {
	body := entity.RegisterUserParams{}
	if err := c.BodyParser(&body); err != nil {
		return r.httpResError(c, errors.NewWithCode(codes.CodeBadRequest, err.Error()))
	}

	user, message, err := r.uc.User.Register(c.UserContext(), body)
	if err != nil {
		return r.httpResError(c, err)
	}

	return r.httpResSuccess(c, codes.CodeSuccess, map[string]any{
		"user":    user,
		"message": message,
	}, nil)
}

func (r *rest) Login(c *fiber.Ctx) error {
	body := entity.LoginUserParams{}
	if err := c.BodyParser(&body); err != nil {
		return r.httpResError(c, errors.NewWithCode(codes.CodeBadRequest, err.Error()))
	}

	user, token, err := r.uc.User.Login(c.UserContext(), body)
	if err != nil {
		return r.httpResError(c, err)
	}

	return r.httpResSuccess(c, codes.CodeSuccess, map[string]any{
		"user":  user,
		"token": token,
	}, nil)
}

func (r *rest) UpdateUser(c *fiber.Ctx) error {
	body := entity.UpdateUserParams{}
	if err := c.BodyParser(&body); err != nil {
		return r.httpResError(c, errors.NewWithCode(codes.CodeBadRequest, err.Error()))
	}

	result, err := r.uc.User.Update(c.UserContext(), body)
	if err != nil {
		return r.httpResError(c, err)
	}

	return r.httpResSuccess(c, codes.CodeSuccess, result, nil)
}

func (r *rest) ListUser(c *fiber.Ctx) error {
	queries := entity.ListUserParams{}
	if err := c.QueryParser(&queries); err != nil {
		return r.httpResError(c, errors.NewWithCode(codes.CodeBadRequest, err.Error()))
	}

	results, pagination, err := r.uc.User.List(c.UserContext(), queries)
	if err != nil {
		return r.httpResError(c, err)
	}

	return r.httpResSuccess(c, codes.CodeSuccess, results, &pagination)
}

func (r *rest) DeleteUser(c *fiber.Ctx) error {
	params := entity.DeleteUserParams{}
	if err := c.ParamsParser(&params); err != nil {
		return r.httpResError(c, errors.NewWithCode(codes.CodeBadRequest, err.Error()))
	}

	result, err := r.uc.User.Delete(c.UserContext(), params)
	if err != nil {
		return r.httpResError(c, err)
	}

	return r.httpResSuccess(c, codes.CodeSuccess, result, nil)
}

func (r *rest) GetUser(c *fiber.Ctx) error {
	params := entity.GetUserParams{}
	if err := c.ParamsParser(&params); err != nil {
		return r.httpResError(c, errors.NewWithCode(codes.CodeBadRequest, err.Error()))
	}

	result, err := r.uc.User.Get(c.UserContext(), params)
	if err != nil {
		return r.httpResError(c, err)
	}

	return r.httpResSuccess(c, codes.CodeSuccess, result, nil)
}
