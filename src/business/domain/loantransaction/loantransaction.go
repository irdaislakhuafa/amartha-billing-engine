package loantransaction

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
		Create(ctx context.Context, params entity.CreateLoanTransactionParams) (entity.LoanTransaction, error)
		List(ctx context.Context, params entity.ListLoanTransactionParams) ([]entity.LoanTransaction, entity.Pagination, error)
		Get(ctx context.Context, params entity.GetLoanTransactionParams) (entity.LoanTransaction, error)
		Update(ctx context.Context, params entity.UpdateLoanTransactionParams) (entity.LoanTransaction, error)
		Delete(ctx context.Context, params entity.DeleteLoanTransactionParams) (entity.LoanTransaction, error)
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

func (i *impl) Create(ctx context.Context, params entity.CreateLoanTransactionParams) (entity.LoanTransaction, error) {
	args := entitygen.CreateLoanTransactionParams{
		InvoiceNumber: params.InvoiceNumber,
		Notes:         params.Notes,
		UserID:        params.UserID,
		User:          params.User,
		LoanID:        params.LoanID,
		Loan:          params.Loan,
		Amount:        params.Amount,
		CreatedAt:     time.Now(),
		CreatedBy:     convert.ToSafeValue[string](ctx.Value(ctxkey.USER_ID)),
	}

	row, err := i.queries.CreateLoanTransaction(ctx, args)
	if err != nil {
		return entity.LoanTransaction{}, errors.NewWithCode(codes.CodeSQLTxExec, err.Error())
	}

	result := entity.LoanTransaction{
		InvoiceNumber: args.InvoiceNumber,
		Notes:         args.Notes,
		UserID:        args.UserID,
		User:          args.User,
		LoanID:        args.LoanID,
		Loan:          args.Loan,
		Amount:        args.Amount,
		Base: entity.Base{
			CreatedAt: time.Now(),
			CreatedBy: convert.ToSafeValue[string](ctx.Value(ctxkey.USER_ID)),
		},
	}

	if result.ID, err = row.LastInsertId(); err != nil {
		return entity.LoanTransaction{}, errors.NewWithCode(codes.CodeSQLRead, err.Error())
	}

	return result, nil
}

func (i *impl) List(ctx context.Context, params entity.ListLoanTransactionParams) ([]entity.LoanTransaction, entity.Pagination, error) {
	paramsBackup := params
	if err := params.Parse(); err != nil {
		return nil, entity.Pagination{}, errors.NewWithCode(errors.GetCode(err), err.Error())
	}

	ctx = sqlc.Build(ctx, func(b *sqlc.Builder) {
		b.And("is_deleted = ?", params.IsDeleted)

		if len(params.Invoices) > 0 {
			_, args := sqlc.GenQueryArgs(ctx, params.Invoices...)
			b.In("invoice_number", args...)
		}

		if len(params.LoanIDs) > 0 {
			_, args := sqlc.GenQueryArgs(ctx, params.LoanIDs...)
			b.In("loan_id", args...)
		}

		if len(params.UserIDs) > 0 {
			_, args := sqlc.GenQueryArgs(ctx, params.UserIDs)
			b.In("user_id", args...)
		}
	})

	rows, err := i.queries.ListLoanTransaction(ctx)
	if err != nil {
		return nil, entity.Pagination{}, errors.NewWithCode(errors.GetCode(err), err.Error())
	}

	results := []entity.LoanTransaction{}
	for _, row := range rows {
		result, err := i.rowToEntity(row)
		if err != nil {
			return nil, entity.Pagination{}, errors.NewWithCode(errors.GetCode(err), err.Error())
		}

		results = append(results, result)
	}

	total, err := i.queries.CountLoanTransaction(ctx)
	if err != nil {
		return nil, entity.Pagination{}, errors.NewWithCode(errors.GetCode(err), err.Error())
	}

	p := entity.GenPagination(paramsBackup.Page, paramsBackup.Limit, int(total), []string{params.OrderBy, params.OrderType})

	return results, p, nil
}

func (i *impl) Get(ctx context.Context, params entity.GetLoanTransactionParams) (entity.LoanTransaction, error) {
	ctx = sqlc.Build(ctx, func(b *sqlc.Builder) {
		b.And("is_deleted = ?", params.IsDeleted)
		b.And("id = ?", params.ID)
	})

	row, err := i.queries.GetLoanTransaction(ctx)
	if err != nil {
		return entity.LoanTransaction{}, errors.NewWithCode(codes.CodeSQLTxExec, err.Error())
	}

	result, err := i.rowToEntity(row)
	if err != nil {
		return entity.LoanTransaction{}, errors.NewWithCode(errors.GetCode(err), err.Error())
	}

	return result, nil
}

func (i *impl) Update(ctx context.Context, params entity.UpdateLoanTransactionParams) (entity.LoanTransaction, error) {
	args := entitygen.UpdateLoanTransactionParams{
		InvoiceNumber: params.InvoiceNumber,
		Notes:         params.Notes,
		UserID:        params.ID,
		User:          params.User,
		LoanID:        params.LoanID,
		Loan:          params.Loan,
		Amount:        params.Amount,
		UpdatedAt:     sql.NullTime{Time: time.Now(), Valid: true},
		UpdatedBy:     sql.NullString{String: convert.ToSafeValue[string](ctx.Value(ctxkey.USER_ID)), Valid: true},
		ID:            params.ID,
	}

	_, err := i.queries.UpdateLoanTransaction(ctx, args)
	if err != nil {
		return entity.LoanTransaction{}, errors.NewWithCode(codes.CodeSQLTxExec, err.Error())
	}

	result := entity.LoanTransaction{
		InvoiceNumber: args.InvoiceNumber,
		Notes:         args.Notes,
		UserID:        args.UserID,
		User:          args.User,
		LoanID:        args.LoanID,
		Loan:          args.Loan,
		Amount:        args.Amount,
		Base: entity.Base{
			ID:        args.ID,
			UpdatedAt: convert.SQLNullTimeToNil(args.UpdatedAt),
			UpdatedBy: convert.SQLNullStringToNil(args.UpdatedBy),
		},
	}

	return result, nil
}

func (i *impl) Delete(ctx context.Context, params entity.DeleteLoanTransactionParams) (entity.LoanTransaction, error) {
	args := entitygen.DeleteLoanTransactionParams{
		ID:        params.ID,
		IsDeleted: params.IsDeleted,
		DeletedAt: sql.NullTime{Time: time.Now(), Valid: true},
		DeletedBy: sql.NullString{String: convert.ToSafeValue[string](ctx.Value(ctxkey.USER_ID)), Valid: true},
	}

	_, err := i.queries.DeleteLoanTransaction(ctx, args)
	if err != nil {
		return entity.LoanTransaction{}, errors.NewWithCode(codes.CodeSQLTxExec, err.Error())
	}

	result := entity.LoanTransaction{
		Base: entity.Base{
			ID:        args.ID,
			IsDeleted: args.IsDeleted,
			DeletedAt: convert.SQLNullTimeToNil(args.DeletedAt),
			DeletedBy: convert.SQLNullStringToNil(args.DeletedBy),
		},
	}

	return result, nil
}

func (i *impl) WithTx(ctx context.Context, tx *sql.Tx) Interface {
	return &impl{
		queries: i.queries.WithTx(tx),
		log:     i.log,
	}
}

func (i *impl) rowToEntity(row entitygen.LoanTransaction) (entity.LoanTransaction, error) {
	result := entity.LoanTransaction{
		InvoiceNumber: row.InvoiceNumber,
		Notes:         row.Notes,
		UserID:        row.UserID,
		User:          row.User,
		LoanID:        row.LoanID,
		Loan:          row.Loan,
		Amount:        row.Amount,
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
