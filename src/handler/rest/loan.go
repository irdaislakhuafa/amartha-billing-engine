package rest

import (
	"github.com/gofiber/fiber/v2"
	"github.com/irdaislakhuafa/amartha-billing-engine/src/entity"
	"github.com/irdaislakhuafa/go-sdk/codes"
	"github.com/irdaislakhuafa/go-sdk/errors"
)

func (r *rest) ListLoan(c *fiber.Ctx) error {
	queries := entity.ListLoanParams{}
	if err := c.QueryParser(&queries); err != nil {
		return r.httpResError(c, errors.NewWithCode(codes.CodeBadRequest, err.Error()))
	}

	result, pagination, err := r.uc.Loan.List(c.UserContext(), queries)
	if err != nil {
		return r.httpResError(c, err)
	}

	return r.httpResSuccess(c, codes.CodeSuccess, result, &pagination)
}

func (r *rest) GetLoan(c *fiber.Ctx) error {
	params := entity.GetLoanParams{}
	if err := c.ParamsParser(&params); err != nil {
		return r.httpResError(c, errors.NewWithCode(codes.CodeBadRequest, err.Error()))
	}

	result, err := r.uc.Loan.Get(c.UserContext(), params)
	if err != nil {
		return r.httpResError(c, err)
	}

	return r.httpResSuccess(c, codes.CodeSuccess, result, nil)
}

func (r *rest) CreateLoan(c *fiber.Ctx) error {
	params := entity.CreateLoanParams{}
	if err := c.BodyParser(&params); err != nil {
		return r.httpResError(c, errors.NewWithCode(codes.CodeBadRequest, err.Error()))
	}

	result, err := r.uc.Loan.Create(c.UserContext(), params)
	if err != nil {
		return r.httpResError(c, err)
	}

	return r.httpResSuccess(c, codes.CodeSuccess, result, nil)
}

func (r *rest) UpdateLoan(c *fiber.Ctx) error {
	params := entity.UpdateLoanParams{}
	if err := c.BodyParser(&params); err != nil {
		return r.httpResError(c, errors.NewWithCode(codes.CodeBadRequest, err.Error()))
	}

	result, err := r.uc.Loan.Update(c.UserContext(), params)
	if err != nil {
		return r.httpResError(c, err)
	}

	return r.httpResSuccess(c, codes.CodeSuccess, result, nil)
}

func (r *rest) DeleteLoan(c *fiber.Ctx) error {
	params := entity.DeleteLoanParams{}
	if err := c.ParamsParser(&params); err != nil {
		return r.httpResError(c, errors.NewWithCode(codes.CodeBadRequest, err.Error()))
	}

	result, err := r.uc.Loan.Delete(c.UserContext(), params)
	if err != nil {
		return r.httpResError(c, err)
	}

	return r.httpResSuccess(c, codes.CodeSuccess, result, nil)
}
