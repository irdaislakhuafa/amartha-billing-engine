package loanbilling

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
		Create(ctx context.Context, params entity.CreateLoanBillingParams) (entity.LoanBilling, error)
		List(ctx context.Context, params any) ([]entity.LoanBilling, entity.Pagination, error)
		Get(ctx context.Context, params any) (entity.LoanBilling, error)
		Update(ctx context.Context, params entity.UpdateLoanBillingParams) (entity.LoanBilling, error)
		Delete(ctx context.Context, params entity.DeleteLoanBillingParams) (entity.LoanBilling, error)
		WithTx(ctx context.Context, tx *sql.Tx) Interface
	}
	impl struct {
		log     log.Interface
		queries *entitygen.Queries
	}
)

func Init(log log.Interface, queries *entitygen.Queries) Interface {
	return &impl{
		log:     log,
		queries: queries,
	}
}

// Create implements Interface.
func (i *impl) Create(ctx context.Context, params entity.CreateLoanBillingParams) (entity.LoanBilling, error) {
	args := entitygen.CreateLoanBillingParams{
		LoanTransactionID:   params.LoanTransactionID,
		BillDate:            params.BillDate,
		PrincipalAmount:     params.PrincipalAmount,
		PrincipalAmountPaid: params.PrincipalAmountPaid,
		InterestAmount:      params.InterestAmount,
		InterestAmountPaid:  params.InterestAmountPaid,
		CreatedAt:           time.Now(),
		CreatedBy:           convert.ToSafeValue[string](ctx.Value(ctxkey.USER_ID)),
	}

	row, err := i.queries.CreateLoanBilling(ctx, args)
	if err != nil {
		return entity.LoanBilling{}, errors.NewWithCode(codes.CodeSQLTxExec, err.Error())
	}

	result := entity.LoanBilling{
		LoanTransactionID:   args.LoanTransactionID,
		BillDate:            args.BillDate,
		PrincipalAmount:     args.PrincipalAmount,
		PrincipalAmountPaid: args.PrincipalAmountPaid,
		InterestAmount:      args.InterestAmount,
		InterestAmountPaid:  args.InterestAmountPaid,
		Base: entity.Base{
			CreatedAt: args.CreatedAt,
			CreatedBy: args.CreatedBy,
		},
	}

	if result.ID, err = row.LastInsertId(); err != nil {
		return entity.LoanBilling{}, errors.NewWithCode(codes.CodeSQLRead, err.Error())
	}

	return result, nil
}

// Delete implements Interface.
func (i *impl) Delete(ctx context.Context, params entity.DeleteLoanBillingParams) (entity.LoanBilling, error) {
	panic("unimplemented")
}

// Get implements Interface.
func (i *impl) Get(ctx context.Context, params any) (entity.LoanBilling, error) {
	panic("unimplemented")
}

// List implements Interface.
func (i *impl) List(ctx context.Context, params any) ([]entity.LoanBilling, entity.Pagination, error) {
	panic("unimplemented")
}

// Update implements Interface.
func (i *impl) Update(ctx context.Context, params entity.UpdateLoanBillingParams) (entity.LoanBilling, error) {
	panic("unimplemented")
}

// WithTx implements Interface.
func (i *impl) WithTx(ctx context.Context, tx *sql.Tx) Interface {
	panic("unimplemented")
}
