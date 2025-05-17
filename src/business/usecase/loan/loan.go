package loan

import (
	"context"
	"database/sql"

	"github.com/go-playground/validator/v10"
	"github.com/irdaislakhuafa/amartha-billing-engine/src/business/domain"
	"github.com/irdaislakhuafa/amartha-billing-engine/src/entity"
	"github.com/irdaislakhuafa/amartha-billing-engine/src/utils/config"
	"github.com/irdaislakhuafa/amartha-billing-engine/src/utils/errmessages"
	"github.com/irdaislakhuafa/amartha-billing-engine/src/utils/validation"
	"github.com/irdaislakhuafa/go-sdk/codes"
	"github.com/irdaislakhuafa/go-sdk/errors"
	"github.com/irdaislakhuafa/go-sdk/log"
)

type (
	Interface interface {
		Create(ctx context.Context, params entity.CreateLoanParams) (entity.Loan, error)
		List(ctx context.Context, params entity.ListLoanParams) ([]entity.Loan, entity.Pagination, error)
		Get(ctx context.Context, params entity.GetLoanParams) (entity.Loan, error)
		Update(ctx context.Context, params entity.UpdateLoanParams) (entity.Loan, error)
		Delete(ctx context.Context, params entity.DeleteLoanParams) (entity.Loan, error)
	}

	impl struct {
		cfg config.Config
		log log.Interface
		val *validator.Validate
		db  *sql.DB
		dom *domain.Domain
	}
)

func Init(
	cfg config.Config,
	log log.Interface,
	val *validator.Validate,
	db *sql.DB,
	dom *domain.Domain,
) Interface {
	return &impl{
		cfg: cfg,
		log: log,
		val: val,
		db:  db,
		dom: dom,
	}
}

func (i *impl) Create(ctx context.Context, params entity.CreateLoanParams) (entity.Loan, error) {
	if err := i.val.StructCtx(ctx, params); err != nil {
		err = validation.ExtractError(err, params)
		return entity.Loan{}, errors.NewWithCode(errors.GetCode(err), err.Error())
	}

	result, err := i.dom.Loan.Create(ctx, params)
	if err != nil {
		return entity.Loan{}, errors.NewWithCode(errors.GetCode(err), err.Error())
	}

	return result, nil
}

func (i *impl) List(ctx context.Context, params entity.ListLoanParams) ([]entity.Loan, entity.Pagination, error) {
	if err := i.val.StructCtx(ctx, params); err != nil {
		err = validation.ExtractError(err, params)
		return []entity.Loan{}, entity.Pagination{}, errors.NewWithCode(errors.GetCode(err), err.Error())
	}

	result, pagination, err := i.dom.Loan.List(ctx, params)
	if err != nil {
		return []entity.Loan{}, entity.Pagination{}, errors.NewWithCode(errors.GetCode(err), err.Error())
	}

	return result, pagination, nil
}

func (i *impl) Get(ctx context.Context, params entity.GetLoanParams) (entity.Loan, error) {
	if err := i.val.StructCtx(ctx, params); err != nil {
		err = validation.ExtractError(err, params)
		return entity.Loan{}, errors.NewWithCode(errors.GetCode(err), err.Error())
	}

	result, err := i.dom.Loan.Get(ctx, params)
	if err != nil {
		return entity.Loan{}, errors.NewWithCode(errors.GetCode(err), err.Error())
	}

	return result, nil
}

func (i *impl) Update(ctx context.Context, params entity.UpdateLoanParams) (entity.Loan, error) {
	if err := i.val.StructCtx(ctx, params); err != nil {
		err = validation.ExtractError(err, params)
		return entity.Loan{}, errors.NewWithCode(errors.GetCode(err), err.Error())
	}

	prev, err := i.dom.Loan.Get(ctx, entity.GetLoanParams{ID: params.ID})
	if err != nil {
		if errors.GetCode(err) == codes.CodeSQLRecordDoesNotExist {
			return entity.Loan{}, errors.NewWithCode(codes.CodeBadRequest, errmessages.LOAN_NOT_FOUND)
		}
		return entity.Loan{}, errors.NewWithCode(errors.GetCode(err), err.Error())
	}

	result, err := i.dom.Loan.Update(ctx, params)
	if err != nil {
		return entity.Loan{}, errors.NewWithCode(errors.GetCode(err), err.Error())
	}

	result.CreatedAt = prev.CreatedAt
	result.CreatedBy = prev.CreatedBy

	return result, nil
}

func (i *impl) Delete(ctx context.Context, params entity.DeleteLoanParams) (entity.Loan, error) {
	if err := i.val.StructCtx(ctx, params); err != nil {
		err = validation.ExtractError(err, params)
		return entity.Loan{}, errors.NewWithCode(errors.GetCode(err), err.Error())
	}

	prev, err := i.dom.Loan.Get(ctx, entity.GetLoanParams{ID: params.ID})
	if err != nil {
		if errors.GetCode(err) == codes.CodeSQLRecordDoesNotExist {
			return entity.Loan{}, errors.NewWithCode(codes.CodeBadRequest, errmessages.LOAN_NOT_FOUND)
		}

		return entity.Loan{}, errors.NewWithCode(errors.GetCode(err), err.Error())
	}

	del, err := i.dom.Loan.Delete(ctx, params)
	if err != nil {
		return entity.Loan{}, errors.NewWithCode(errors.GetCode(err), err.Error())
	}

	prev.DeletedAt = del.DeletedAt
	prev.DeletedBy = del.DeletedBy
	prev.IsDeleted = del.IsDeleted

	return prev, nil
}
