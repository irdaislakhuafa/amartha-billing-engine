package user

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
		Create(ctx context.Context, params entity.CreateUserParams) (entity.User, error)
		Get(ctx context.Context, params entity.GetUserParams) (entity.User, error)
		List(ctx context.Context, params entity.ListUserParams) ([]entity.User, entity.Pagination, error)
		Update(ctx context.Context, params entity.UpdateUserParams) (entity.User, error)
		Delete(ctx context.Context, params entity.DeleteUserParams) (entity.User, error)
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

func (i *impl) Create(ctx context.Context, params entity.CreateUserParams) (entity.User, error) {
	args := entitygen.CreateUserParams{
		Name:      params.Name,
		Email:     params.Email,
		Password:  params.Password,
		CreatedAt: time.Now(),
		CreatedBy: convert.ToSafeValue[string](ctx.Value(ctxkey.USER_ID)),
	}

	row, err := i.queries.CreateUser(ctx, args)
	if err != nil {
		return entity.User{}, errors.NewWithCode(codes.CodeSQLTxExec, err.Error())
	}

	result := entity.User{
		Name:     args.Name,
		Email:    args.Email,
		Password: args.Password,
		Base: entity.Base{
			CreatedAt: args.CreatedAt,
			CreatedBy: args.CreatedBy,
		},
	}

	if result.ID, err = row.LastInsertId(); err != nil {
		return entity.User{}, errors.NewWithCode(codes.CodeSQLRead, err.Error())
	}

	return result, nil
}

func (i *impl) Get(ctx context.Context, params entity.GetUserParams) (entity.User, error) {
	ctx = sqlc.Build(ctx, func(b *sqlc.Builder) {
		b.And("(id = ? OR email = ?)", params.ID, params.Email)
		b.And("is_deleted = ?", params.IsDeleted)
	})

	row, err := i.queries.GetUser(ctx)
	if err != nil {
		if err == sql.ErrNoRows {
			return entity.User{}, errors.NewWithCode(codes.CodeSQLRecordDoesNotExist, err.Error())
		}
		return entity.User{}, errors.NewWithCode(codes.CodeSQLRead, err.Error())
	}

	result, err := i.rowToEntity(row)
	if err != nil {
		return entity.User{}, errors.NewWithCode(errors.GetCode(err), err.Error())
	}

	return result, nil
}

func (i *impl) List(ctx context.Context, params entity.ListUserParams) ([]entity.User, entity.Pagination, error) {
	paramsBackup := params
	if err := params.Parse(); err != nil {
		return nil, entity.Pagination{}, errors.NewWithCode(errors.GetCode(err), err.Error())
	}

	ctx = sqlc.Build(ctx, func(b *sqlc.Builder) {
		b.And("is_deleted = ?", params.IsDeleted)
		if params.Search != "" {
			params.Search = "%" + params.Search + "%"
			b.And("(name LIKE ? OR email LIKE ?)", params.Search, params.Search)
		}

		if len(params.IDs) > 0 {
			_, args := sqlc.GenQueryArgs(ctx, params.IDs...)
			b.In("id", args...)
		}
	})

	rows, err := i.queries.ListUser(sqlc.Build(ctx, func(b *sqlc.Builder) {
		b.Limit(params.Limit)
		b.Offset(params.Page)
		b.Order(params.OrderBy + " " + params.OrderType)
	}))
	if err != nil {
		return []entity.User{}, entity.Pagination{}, errors.NewWithCode(codes.CodeSQLRead, err.Error())
	}

	var results []entity.User
	for _, row := range rows {
		result, err := i.rowToEntity(row)
		if err != nil {
			return []entity.User{}, entity.Pagination{}, errors.NewWithCode(errors.GetCode(err), err.Error())
		}
		results = append(results, result)
	}

	total, err := i.queries.CountUser(ctx)
	if err != nil {
		return []entity.User{}, entity.Pagination{}, errors.NewWithCode(codes.CodeSQLRead, err.Error())
	}

	p := entity.GenPagination(paramsBackup.Page, paramsBackup.Limit, int(total), []string{params.OrderBy, params.OrderType})

	return results, p, nil
}

func (i *impl) Update(ctx context.Context, params entity.UpdateUserParams) (entity.User, error) {
	args := entitygen.UpdateUserParams{
		Name:  params.Name,
		Email: params.Email,
		// Password:  params.Password,
		ID:              params.ID,
		DelinquentLevel: int32(params.DelinquentLevel),
		UpdatedAt:       sql.NullTime{Time: time.Now(), Valid: true},
		UpdatedBy:       sql.NullString{String: convert.ToSafeValue[string](ctx.Value(ctxkey.USER_ID)), Valid: true},
	}

	_, err := i.queries.UpdateUser(ctx, args)
	if err != nil {
		return entity.User{}, errors.NewWithCode(codes.CodeSQLTxExec, err.Error())
	}

	result := entity.User{
		Name:            args.Name,
		Email:           args.Email,
		DelinquentLevel: int(args.DelinquentLevel),
		// Password: args.Password,
		Base: entity.Base{
			ID:        args.ID,
			UpdatedAt: &args.UpdatedAt.Time,
			UpdatedBy: &args.UpdatedBy.String,
		},
	}

	return result, nil
}

func (i *impl) Delete(ctx context.Context, params entity.DeleteUserParams) (entity.User, error) {
	args := entitygen.DeleteUserParams{
		ID:        params.ID,
		IsDeleted: int8(params.IsDeleted),
		DeletedAt: sql.NullTime{Time: time.Now(), Valid: true},
		DeletedBy: sql.NullString{String: convert.ToSafeValue[string](ctx.Value(ctxkey.USER_ID)), Valid: true},
	}

	_, err := i.queries.DeleteUser(ctx, args)
	if err != nil {
		return entity.User{}, errors.NewWithCode(codes.CodeSQLTxExec, err.Error())
	}

	result := entity.User{
		Base: entity.Base{
			ID:        args.ID,
			IsDeleted: args.IsDeleted,
			DeletedAt: &args.DeletedAt.Time,
			DeletedBy: &args.DeletedBy.String,
		},
	}

	return result, nil
}

func (i *impl) WithTx(ctx context.Context, tx *sql.Tx) Interface {
	return &impl{
		log:     i.log,
		queries: i.queries.WithTx(tx),
	}
}

func (i *impl) rowToEntity(row entitygen.User) (entity.User, error) {
	result := entity.User{
		Name:            row.Name,
		Email:           row.Email,
		Password:        row.Password,
		DelinquentLevel: int(row.DelinquentLevel),
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
