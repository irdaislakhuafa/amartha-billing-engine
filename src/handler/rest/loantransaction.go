package rest

import (
	"github.com/gofiber/fiber/v2"
	"github.com/irdaislakhuafa/amartha-billing-engine/src/entity"
	"github.com/irdaislakhuafa/go-sdk/codes"
	"github.com/irdaislakhuafa/go-sdk/errors"
)

func (r *rest) CreateLoanTransaction(c *fiber.Ctx) error {
	body := entity.CreateLoanTransactionParams{}
	if err := c.BodyParser(&body); err != nil {
		return r.httpResError(c, errors.NewWithCode(codes.CodeBadRequest, err.Error()))
	}

	result, err := r.uc.LoanTransaction.Create(c.UserContext(), body)
	if err != nil {
		return r.httpResError(c, err)
	}

	return r.httpResSuccess(c, codes.CodeSuccess, result, nil)
}

func (r *rest) CalculateOutstandingLoanTransaction(c *fiber.Ctx) error {
	params := entity.CalculateOutstandingLoanTransactionParams{}
	if err := c.ParamsParser(&params); err != nil {
		return r.httpResError(c, errors.NewWithCode(codes.CodeBadRequest, err.Error()))
	}

	result, err := r.uc.LoanTransaction.CalculateOutstanding(c.UserContext(), params)
	if err != nil {
		return r.httpResError(c, err)
	}

	return r.httpResSuccess(c, codes.CodeSuccess, result, nil)
}

func (r *rest) PayLoanTransaction(c *fiber.Ctx) error {
	body := entity.PayLoanTransactionParams{}
	if err := c.BodyParser(&body); err != nil {
		return r.httpResError(c, errors.NewWithCode(codes.CodeBadRequest, err.Error()))
	}

	result, err := r.uc.LoanTransaction.Pay(c.UserContext(), body)
	if err != nil {
		return r.httpResError(c, err)
	}

	return r.httpResSuccess(c, codes.CodeSuccess, result, nil)
}

func (r *rest) ListLoanTransaction(c *fiber.Ctx) error {
	queries := entity.ListLoanTransactionParams{}
	if err := c.QueryParser(&queries); err != nil {
		return r.httpResError(c, errors.NewWithCode(codes.CodeBadRequest, err.Error()))
	}

	results, pagination, err := r.uc.LoanTransaction.List(c.UserContext(), queries)
	if err != nil {
		return r.httpResError(c, err)
	}

	return r.httpResSuccess(c, codes.CodeSuccess, results, &pagination)
}

func (r *rest) ScheduleDelinquent(c *fiber.Ctx) error {
	err := r.uc.LoanTransaction.ScheduleDelinquent(c.UserContext())
	if err != nil {
		return r.httpResError(c, err)
	}

	return r.httpResSuccess(c, codes.CodeSuccess, "success", nil)
}
