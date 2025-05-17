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
	"github.com/irdaislakhuafa/go-sdk/querybuilder/sqlc"
)

type (
	Interface interface {
		Create(ctx context.Context, params entity.CreateLoanBillingParams) (entity.LoanBilling, error)
		List(ctx context.Context, params entity.ListLoanBillingParams) ([]entity.LoanBilling, entity.Pagination, error)
		Get(ctx context.Context, params entity.GetLoanBillingParams) (entity.LoanBilling, error)
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
	args := entitygen.DeleteLoanBillingParams{
		IsDeleted: params.IsDeleted,
		DeletedAt: sql.NullTime{},
		DeletedBy: sql.NullString{},
		ID:        params.ID,
	}

	_, err := i.queries.DeleteLoanBilling(ctx, args)
	if err != nil {
		return entity.LoanBilling{}, errors.NewWithCode(codes.CodeSQLTxExec, err.Error())
	}

	result := entity.LoanBilling{
		Base: entity.Base{
			ID:        args.ID,
			DeletedAt: convert.SQLNullTimeToNil(args.DeletedAt),
			DeletedBy: convert.SQLNullStringToNil(args.DeletedBy),
			IsDeleted: args.IsDeleted,
		},
	}

	return result, nil
}

// Get implements Interface.
func (i *impl) Get(ctx context.Context, params entity.GetLoanBillingParams) (entity.LoanBilling, error) {
	ctx = sqlc.Build(ctx, func(b *sqlc.Builder) {
		b.And("id = ?", params.ID)
		b.And("is_deleted = ?", params.IsDeleted)
	})

	row, err := i.queries.GetLoanBilling(ctx)
	if err != nil {
		return entity.LoanBilling{}, errors.NewWithCode(codes.CodeSQLTxExec, err.Error())
	}

	result, err := i.rowToEntity(row)
	if err != nil {
		return entity.LoanBilling{}, errors.NewWithCode(errors.GetCode(err), err.Error())
	}

	return result, nil
}

// List implements Interface.
func (i *impl) List(ctx context.Context, params entity.ListLoanBillingParams) ([]entity.LoanBilling, entity.Pagination, error) {
	paramsBackup := params
	if err := params.Parse(); err != nil {
		return []entity.LoanBilling{}, entity.Pagination{}, errors.NewWithCode(errors.GetCode(err), err.Error())
	}

	ctx = sqlc.Build(ctx, func(b *sqlc.Builder) {
		b.And("is_deleted = ?", params.IsDeleted)
		if params.LoanTransactionID != 0 {
			b.And("loan_transaction_id = ?", params.LoanTransactionID)
		}

		if params.IsDeleted != 0 {
			b.And("is_deleted = ?", params.IsDeleted)
		}

		if !params.BillDateGTE.IsZero() {
			b.And("bill_date >= ?", params.BillDateGTE)
		}

		if !params.BillDateLTE.IsZero() {
			b.And("bill_date <= ?", params.BillDateLTE)
		}

		if !params.PrincipalAmountGTE.IsZero() {
			b.And("principal_amount >= ?", params.PrincipalAmountGTE)
		}

		if !params.PrincipalAmountLTE.IsZero() {
			b.And("principal_amount <= ?", params.PrincipalAmountLTE)
		}

		if !params.PrincipalAmountPaidGTE.IsZero() {
			b.And("principal_amount_paid >= ?", params.PrincipalAmountPaidGTE)
		}

		if !params.PrincipalAmountPaidLTE.IsZero() {
			b.And("principal_amount_paid <= ?", params.PrincipalAmountPaidLTE)
		}

		if !params.InterestAmountGTE.IsZero() {
			b.And("interest_amount >= ?", params.InterestAmountGTE)
		}

		if !params.InterestAmountLTE.IsZero() {
			b.And("interest_amount <= ?", params.InterestAmountLTE)
		}

		if !params.InterestAmountPaidGTE.IsZero() {
			b.And("interest_amount_paid >= ?", params.InterestAmountPaidGTE)
		}

		if !params.InterestAmountPaidLTE.IsZero() {
			b.And("interest_amount_paid <= ?", params.InterestAmountPaidLTE)
		}
	})

	rows, err := i.queries.ListLoanBilling(sqlc.Build(ctx, func(b *sqlc.Builder) {
		b.Limit(params.Limit)
		b.Offset(params.Page)
		b.Order(params.OrderBy + " " + params.OrderType)
	}))
	if err != nil {
		return []entity.LoanBilling{}, entity.Pagination{}, errors.NewWithCode(codes.CodeSQLRead, err.Error())
	}

	var results []entity.LoanBilling
	for _, row := range rows {
		result, err := i.rowToEntity(row)
		if err != nil {
			return []entity.LoanBilling{}, entity.Pagination{}, errors.NewWithCode(errors.GetCode(err), err.Error())
		}
		results = append(results, result)
	}

	total, err := i.queries.CountLoanBilling(ctx)
	if err != nil {
		return []entity.LoanBilling{}, entity.Pagination{}, errors.NewWithCode(codes.CodeSQLRead, err.Error())
	}

	p := entity.GenPagination(paramsBackup.Page, paramsBackup.Limit, int(total), []string{paramsBackup.OrderBy, paramsBackup.OrderType})

	return results, p, nil
}

// Update implements Interface.
func (i *impl) Update(ctx context.Context, params entity.UpdateLoanBillingParams) (entity.LoanBilling, error) {
	args := entitygen.UpdateLoanBillingParams{
		LoanTransactionID:   params.LoanTransactionID,
		BillDate:            params.BillDate,
		PrincipalAmount:     params.PrincipalAmount,
		PrincipalAmountPaid: params.PrincipalAmountPaid,
		InterestAmount:      params.InterestAmount,
		InterestAmountPaid:  params.InterestAmountPaid,
		UpdatedAt:           sql.NullTime{Time: time.Now(), Valid: true},
		UpdatedBy:           sql.NullString{String: convert.ToSafeValue[string](ctx.Value(ctxkey.USER_ID)), Valid: true},
		ID:                  params.ID,
	}

	_, err := i.queries.UpdateLoanBilling(ctx, args)
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
			ID:        args.ID,
			UpdatedAt: convert.SQLNullTimeToNil(args.UpdatedAt),
			UpdatedBy: convert.SQLNullStringToNil(args.UpdatedBy),
		},
	}

	return result, nil
}

// WithTx implements Interface.
func (i *impl) WithTx(ctx context.Context, tx *sql.Tx) Interface {
	return &impl{
		log:     i.log,
		queries: i.queries.WithTx(tx),
	}
}

func (i *impl) rowToEntity(row entitygen.LoanBilling) (entity.LoanBilling, error) {
	result := entity.LoanBilling{
		LoanTransactionID:   row.LoanTransactionID,
		BillDate:            row.BillDate,
		PrincipalAmount:     row.PrincipalAmount,
		PrincipalAmountPaid: row.PrincipalAmountPaid,
		InterestAmount:      row.InterestAmount,
		InterestAmountPaid:  row.InterestAmountPaid,
		Base: entity.Base{
			ID:        row.ID,
			CreatedAt: row.CreatedAt,
			CreatedBy: row.CreatedBy,
			UpdatedAt: convert.SQLNullTimeToNil(row.UpdatedAt),
			UpdatedBy: convert.SQLNullStringToNil(row.UpdatedBy),
			DeletedAt: convert.SQLNullTimeToNil(row.DeletedAt),
			DeletedBy: convert.SQLNullStringToNil(row.DeletedBy),
			IsDeleted: row.IsDeleted,
		},
	}

	return result, nil
}
