package loandelinquenthistories

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
		Create(ctx context.Context, params entity.CreateLoanDelinquentHistoryParams) (entity.LoanDelinquentHistory, error)
		List(ctx context.Context, params entity.ListLoanDelinquentHistoryParams) ([]entity.LoanDelinquentHistory, entity.Pagination, error)
		Get(ctx context.Context, params entity.GetLoanDelinquentHistoryParams) (entity.LoanDelinquentHistory, error)
		Update(ctx context.Context, params entity.UpdateLoanDelinquentHistoryParams) (entity.LoanDelinquentHistory, error)
		Delete(ctx context.Context, params entity.DeleteLoanDelinquentHistoryParams) (entity.LoanDelinquentHistory, error)
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
func (i *impl) Create(ctx context.Context, params entity.CreateLoanDelinquentHistoryParams) (entity.LoanDelinquentHistory, error) {
	args := entitygen.CreateLoanDelinquentHistoryParams{
		LoanTransactionID: params.LoanTransactionID,
		Bills:             params.Bills,
		CreatedAt:         time.Now(),
		CreatedBy:         convert.ToSafeValue[string](ctx.Value(ctxkey.USER_ID)),
	}

	row, err := i.queries.CreateLoanDelinquentHistory(ctx, args)
	if err != nil {
		return entity.LoanDelinquentHistory{}, errors.NewWithCode(codes.CodeSQLTxExec, err.Error())
	}

	result := entity.LoanDelinquentHistory{
		LoanTransactionID: args.LoanTransactionID,
		Bills:             args.Bills,
		Base: entity.Base{
			CreatedAt: args.CreatedAt,
			CreatedBy: args.CreatedBy,
		},
	}

	if result.ID, err = row.LastInsertId(); err != nil {
		return entity.LoanDelinquentHistory{}, errors.NewWithCode(codes.CodeSQLRead, err.Error())
	}

	return result, nil
}

// Delete implements Interface.
func (i *impl) Delete(ctx context.Context, params entity.DeleteLoanDelinquentHistoryParams) (entity.LoanDelinquentHistory, error) {
	args := entitygen.DeleteLoanDelinquentHistoryParams{
		IsDeleted: params.IsDeleted,
		DeletedAt: sql.NullTime{Time: time.Now(), Valid: true},
		DeletedBy: sql.NullString{String: convert.ToSafeValue[string](ctx.Value(ctxkey.USER_ID)), Valid: true},
		ID:        params.ID,
	}

	_, err := i.queries.DeleteLoanDelinquentHistory(ctx, args)
	if err != nil {
		return entity.LoanDelinquentHistory{}, errors.NewWithCode(codes.CodeSQLTxExec, err.Error())
	}

	result := entity.LoanDelinquentHistory{
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
func (i *impl) Get(ctx context.Context, params entity.GetLoanDelinquentHistoryParams) (entity.LoanDelinquentHistory, error) {
	ctx = sqlc.Build(ctx, func(b *sqlc.Builder) {
		b.And("id = ?", params.ID)
		b.And("is_deleted = ?", params.IsDeleted)
	})

	row, err := i.queries.GetLoanDelinquentHistory(ctx)
	if err != nil {
		return entity.LoanDelinquentHistory{}, errors.NewWithCode(codes.CodeSQLRead, err.Error())
	}

	result, err := i.rowToEntity(row)
	if err != nil {
		return entity.LoanDelinquentHistory{}, errors.NewWithCode(errors.GetCode(err), err.Error())
	}

	return result, nil
}

// List implements Interface.
func (i *impl) List(ctx context.Context, params entity.ListLoanDelinquentHistoryParams) ([]entity.LoanDelinquentHistory, entity.Pagination, error) {
	paramsBackup := params
	if err := params.Parse(); err != nil {
		return nil, entity.Pagination{}, errors.NewWithCode(errors.GetCode(err), err.Error())
	}

	ctx = sqlc.Build(ctx, func(b *sqlc.Builder) {
		b.And("is_deleted = ?", params.IsDeleted)
		if params.LoanTransactionID != 0 {
			b.And("loan_transaction_id = ?", params.LoanTransactionID)
		}
	})

	rows, err := i.queries.ListLoanDelinquentHistory(sqlc.Build(ctx, func(b *sqlc.Builder) {
		b.Limit(params.Limit)
		b.Offset(params.Page)
		b.Order(params.OrderBy + " " + params.OrderType)
	}))
	if err != nil {
		return []entity.LoanDelinquentHistory{}, entity.Pagination{}, errors.NewWithCode(codes.CodeSQLRead, err.Error())
	}

	var results []entity.LoanDelinquentHistory
	for _, row := range rows {
		result, err := i.rowToEntity(row)
		if err != nil {
			return []entity.LoanDelinquentHistory{}, entity.Pagination{}, errors.NewWithCode(errors.GetCode(err), err.Error())
		}
		results = append(results, result)
	}

	total, err := i.queries.CountLoanDelinquentHistory(ctx)
	if err != nil {
		return []entity.LoanDelinquentHistory{}, entity.Pagination{}, errors.NewWithCode(codes.CodeSQLRead, err.Error())
	}

	p := entity.GenPagination(paramsBackup.Page, paramsBackup.Limit, int(total), []string{params.OrderBy, params.OrderType})

	return results, p, nil
}

// Update implements Interface.
func (i *impl) Update(ctx context.Context, params entity.UpdateLoanDelinquentHistoryParams) (entity.LoanDelinquentHistory, error) {
	args := entitygen.UpdateLoanDelinquentHistoryParams{
		LoanTransactionID: params.LoanTransactionID,
		Bills:             params.Bills,
		UpdatedAt:         sql.NullTime{Time: time.Now(), Valid: true},
		UpdatedBy:         sql.NullString{String: convert.ToSafeValue[string](ctx.Value(ctxkey.USER_ID)), Valid: true},
		ID:                params.ID,
	}

	_, err := i.queries.UpdateLoanDelinquentHistory(ctx, args)
	if err != nil {
		return entity.LoanDelinquentHistory{}, errors.NewWithCode(codes.CodeSQLTxExec, err.Error())
	}

	result := entity.LoanDelinquentHistory{
		LoanTransactionID: args.LoanTransactionID,
		Bills:             args.Bills,
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
		queries: i.queries.WithTx(tx),
	}
}

func (i *impl) rowToEntity(row entitygen.LoanDelinquentHistory) (entity.LoanDelinquentHistory, error) {
	result := entity.LoanDelinquentHistory{
		LoanTransactionID: row.LoanTransactionID,
		Bills:             row.Bills,
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
