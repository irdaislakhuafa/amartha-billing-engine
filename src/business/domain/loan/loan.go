package loan

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
	"github.com/shopspring/decimal"
)

type (
	Interface interface {
		Create(ctx context.Context, params entity.CreateLoanParams) (entity.Loan, error)
		List(ctx context.Context, params entity.ListLoanParams) ([]entity.Loan, entity.Pagination, error)
		Get(ctx context.Context, params entity.GetLoanParams) (entity.Loan, error)
		Update(ctx context.Context, params entity.UpdateLoanParams) (entity.Loan, error)
		Delete(ctx context.Context, params entity.DeleteLoanParams) (entity.Loan, error)
		WithTx(ctx context.Context, tx *sql.Tx) Interface
	}

	impl struct {
		queries *entitygen.Queries
		db      *sql.DB
		log     log.Interface
	}
)

func Init(queries *entitygen.Queries, db *sql.DB, log log.Interface) Interface {
	return &impl{
		queries: queries,
		db:      db,
		log:     log,
	}
}

func (i *impl) WithTx(ctx context.Context, tx *sql.Tx) Interface {
	return &impl{
		queries: i.queries.WithTx(tx),
		db:      i.db,
		log:     i.log,
	}
}

func (i *impl) Create(ctx context.Context, params entity.CreateLoanParams) (entity.Loan, error) {
	args := entitygen.CreateLoanParams{
		Name:              params.Name,
		Description:       params.Description,
		InterestRate:      decimal.NewFromFloat(params.InterestRate),
		RepaymentType:     params.RepaymentType,
		RepaymentDuration: int32(params.RepaymentDuration),
		CreatedAt:         time.Now(),
		CreatedBy:         convert.ToSafeValue[string](ctx.Value(ctxkey.USER_ID)),
	}

	row, err := i.queries.CreateLoan(ctx, args)
	if err != nil {
		return entity.Loan{}, errors.NewWithCode(codes.CodeSQLTxExec, err.Error())
	}

	result := entity.Loan{
		Name:              args.Name,
		Description:       args.Description,
		InterestRate:      args.InterestRate,
		RepaymentType:     args.RepaymentType,
		RepaymentDuration: args.RepaymentDuration,
		Base: entity.Base{
			CreatedAt: args.CreatedAt,
			CreatedBy: args.CreatedBy,
		},
	}

	if result.ID, err = row.LastInsertId(); err != nil {
		return entity.Loan{}, errors.NewWithCode(codes.CodeSQLRead, err.Error())
	}

	return result, nil
}

func (i *impl) List(ctx context.Context, params entity.ListLoanParams) ([]entity.Loan, entity.Pagination, error) {
	paramsBackup := params
	if err := params.Parse(); err != nil {
		return []entity.Loan{}, entity.Pagination{}, errors.NewWithCode(errors.GetCode(err), err.Error())
	}

	ctx = sqlc.Build(ctx, func(b *sqlc.Builder) {
		if len(params.IDs) > 0 {
			_, args := sqlc.GenQueryArgs(ctx, params.IDs...)
			b.In("id", args...)
		}

		if len(params.RepaymentTypes) > 0 {
			_, args := sqlc.GenQueryArgs(ctx, params.RepaymentTypes...)
			b.In("repayment_type", args...)
		}

		if params.Search != "" {
			params.Search = "%" + params.Search + "%"
			b.Where("name LIKE ?", params.Search)
		}

		b.And("is_deleted = ?", params.IsDeleted)
		b.Order(params.OrderBy + " " + params.OrderType)
	})

	rows, err := i.queries.ListLoan(sqlc.Build(ctx, func(b *sqlc.Builder) {
		b.Limit(params.Limit)
		b.Offset(params.Page)
	}))
	if err != nil {
		return []entity.Loan{}, entity.Pagination{}, errors.NewWithCode(codes.CodeSQLRead, err.Error())
	}

	var results []entity.Loan
	for _, row := range rows {
		result, err := i.rowToEntity(row)
		if err != nil {
			return []entity.Loan{}, entity.Pagination{}, errors.NewWithCode(errors.GetCode(err), err.Error())
		}
		results = append(results, result)
	}

	total, err := i.queries.CountLoan(ctx)
	if err != nil {
		return []entity.Loan{}, entity.Pagination{}, errors.NewWithCode(codes.CodeSQLRead, err.Error())
	}

	p := entity.GenPagination(paramsBackup.Page, paramsBackup.Limit, int(total), []string{params.OrderBy, params.OrderType})

	return results, p, nil
}

func (i *impl) Get(ctx context.Context, params entity.GetLoanParams) (entity.Loan, error) {
	ctx = sqlc.Build(ctx, func(b *sqlc.Builder) {
		b.And("id = ?", params.ID)
	})

	row, err := i.queries.GetLoan(ctx)
	if err != nil {
		if err == sql.ErrNoRows {
			return entity.Loan{}, errors.NewWithCode(codes.CodeSQLRecordDoesNotExist, err.Error())
		}
		return entity.Loan{}, errors.NewWithCode(codes.CodeSQLRead, err.Error())
	}

	result, err := i.rowToEntity(row)
	if err != nil {
		return entity.Loan{}, errors.NewWithCode(errors.GetCode(err), err.Error())
	}

	return result, nil
}

func (i *impl) Update(ctx context.Context, params entity.UpdateLoanParams) (entity.Loan, error) {
	args := entitygen.UpdateLoanParams{
		Name:              params.Name,
		Description:       params.Description,
		InterestRate:      decimal.NewFromFloat(params.InterestRate),
		RepaymentType:     params.RepaymentType,
		RepaymentDuration: int32(params.RepaymentDuration),
		UpdatedAt:         sql.NullTime{Time: time.Now(), Valid: true},
		UpdatedBy:         sql.NullString{String: convert.ToSafeValue[string](ctx.Value(ctxkey.USER_ID)), Valid: true},
		ID:                params.ID,
	}

	_, err := i.queries.UpdateLoan(ctx, args)
	if err != nil {
		return entity.Loan{}, errors.NewWithCode(codes.CodeSQLTxExec, err.Error())
	}

	result := entity.Loan{
		Name:              args.Name,
		Description:       args.Description,
		InterestRate:      args.InterestRate,
		RepaymentType:     args.RepaymentType,
		RepaymentDuration: args.RepaymentDuration,
		Base: entity.Base{
			ID:        args.ID,
			UpdatedAt: &args.UpdatedAt.Time,
			UpdatedBy: &args.UpdatedBy.String,
		},
	}

	return result, nil
}

func (i *impl) Delete(ctx context.Context, params entity.DeleteLoanParams) (entity.Loan, error) {
	args := entitygen.DeleteLoanParams{
		ID:        params.ID,
		IsDeleted: int8(params.IsDeleted),
		DeletedAt: sql.NullTime{Time: time.Now(), Valid: true},
		DeletedBy: sql.NullString{String: convert.ToSafeValue[string](ctx.Value(ctxkey.USER_ID)), Valid: true},
	}

	_, err := i.queries.DeleteLoan(ctx, args)
	if err != nil {
		return entity.Loan{}, errors.NewWithCode(codes.CodeSQLTxExec, err.Error())
	}

	result := entity.Loan{
		Base: entity.Base{
			ID:        args.ID,
			IsDeleted: args.IsDeleted,
			DeletedAt: &args.DeletedAt.Time,
			DeletedBy: &args.DeletedBy.String,
		},
	}

	return result, nil
}

func (i *impl) rowToEntity(row entitygen.Loan) (entity.Loan, error) {
	result := entity.Loan{
		Name:              row.Name,
		Description:       row.Description,
		InterestRate:      row.InterestRate,
		RepaymentType:     row.RepaymentType,
		RepaymentDuration: row.RepaymentDuration,
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
