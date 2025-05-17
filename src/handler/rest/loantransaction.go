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
